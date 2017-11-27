package tsi1_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"testing"

	"github.com/influxdata/influxdb/models"
	"github.com/influxdata/influxdb/tsdb/index/tsi1"
	"github.com/influxdata/influxql"
)

// Bloom filter settings used in tests.
const M, K = 4096, 6

// Ensure index can iterate over all measurement names.
func TestIndex_ForEachMeasurementName(t *testing.T) {
	idx := MustOpenIndex()
	defer idx.Close()

	// Add series to index.
	if err := idx.CreateSeriesSliceIfNotExists([]Series{
		{Name: []byte("cpu"), Tags: models.NewTags(map[string]string{"region": "east"})},
		{Name: []byte("cpu"), Tags: models.NewTags(map[string]string{"region": "west"})},
		{Name: []byte("mem"), Tags: models.NewTags(map[string]string{"region": "east"})},
	}); err != nil {
		t.Fatal(err)
	}

	// Verify measurements are returned.
	idx.Run(t, func(t *testing.T) {
		var names []string
		if err := idx.ForEachMeasurementName(func(name []byte) error {
			names = append(names, string(name))
			return nil
		}); err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(names, []string{"cpu", "mem"}) {
			t.Fatalf("unexpected names: %#v", names)
		}
	})

	// Add more series.
	if err := idx.CreateSeriesSliceIfNotExists([]Series{
		{Name: []byte("disk")},
		{Name: []byte("mem")},
	}); err != nil {
		t.Fatal(err)
	}

	// Verify new measurements.
	idx.Run(t, func(t *testing.T) {
		var names []string
		if err := idx.ForEachMeasurementName(func(name []byte) error {
			names = append(names, string(name))
			return nil
		}); err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(names, []string{"cpu", "disk", "mem"}) {
			t.Fatalf("unexpected names: %#v", names)
		}
	})
}

// Ensure index can return whether a measurement exists.
func TestIndex_MeasurementExists(t *testing.T) {
	idx := MustOpenIndex()
	defer idx.Close()

	// Add series to index.
	if err := idx.CreateSeriesSliceIfNotExists([]Series{
		{Name: []byte("cpu"), Tags: models.NewTags(map[string]string{"region": "east"})},
		{Name: []byte("cpu"), Tags: models.NewTags(map[string]string{"region": "west"})},
	}); err != nil {
		t.Fatal(err)
	}

	// Verify measurement exists.
	idx.Run(t, func(t *testing.T) {
		if v, err := idx.MeasurementExists([]byte("cpu")); err != nil {
			t.Fatal(err)
		} else if !v {
			t.Fatal("expected measurement to exist")
		}
	})

	// Delete one series.
	if err := idx.DropSeries(models.MakeKey([]byte("cpu"), models.NewTags(map[string]string{"region": "east"}))); err != nil {
		t.Fatal(err)
	}

	// Verify measurement still exists.
	idx.Run(t, func(t *testing.T) {
		if v, err := idx.MeasurementExists([]byte("cpu")); err != nil {
			t.Fatal(err)
		} else if !v {
			t.Fatal("expected measurement to still exist")
		}
	})

	// Delete second series.
	if err := idx.DropSeries(models.MakeKey([]byte("cpu"), models.NewTags(map[string]string{"region": "west"}))); err != nil {
		t.Fatal(err)
	}

	// Verify measurement is now deleted.
	idx.Run(t, func(t *testing.T) {
		if v, err := idx.MeasurementExists([]byte("cpu")); err != nil {
			t.Fatal(err)
		} else if v {
			t.Fatal("expected measurement to be deleted")
		}
	})
}

// Ensure index can return a list of matching measurements.
func TestIndex_MeasurementNamesByExpr(t *testing.T) {
	idx := MustOpenIndex()
	defer idx.Close()

	// Add series to index.
	if err := idx.CreateSeriesSliceIfNotExists([]Series{
		{Name: []byte("cpu"), Tags: models.NewTags(map[string]string{"region": "east"})},
		{Name: []byte("cpu"), Tags: models.NewTags(map[string]string{"region": "west"})},
		{Name: []byte("disk"), Tags: models.NewTags(map[string]string{"region": "north"})},
		{Name: []byte("mem"), Tags: models.NewTags(map[string]string{"region": "west", "country": "us"})},
	}); err != nil {
		t.Fatal(err)
	}

	// Retrieve measurements by expression
	idx.Run(t, func(t *testing.T) {
		t.Run("EQ", func(t *testing.T) {
			names, err := idx.MeasurementNamesByExpr(influxql.MustParseExpr(`region = 'west'`))
			if err != nil {
				t.Fatal(err)
			} else if !reflect.DeepEqual(names, [][]byte{[]byte("cpu"), []byte("mem")}) {
				t.Fatalf("unexpected names: %v", names)
			}
		})

		t.Run("NEQ", func(t *testing.T) {
			names, err := idx.MeasurementNamesByExpr(influxql.MustParseExpr(`region != 'east'`))
			if err != nil {
				t.Fatal(err)
			} else if !reflect.DeepEqual(names, [][]byte{[]byte("disk"), []byte("mem")}) {
				t.Fatalf("unexpected names: %v", names)
			}
		})

		t.Run("EQREGEX", func(t *testing.T) {
			names, err := idx.MeasurementNamesByExpr(influxql.MustParseExpr(`region =~ /east|west/`))
			if err != nil {
				t.Fatal(err)
			} else if !reflect.DeepEqual(names, [][]byte{[]byte("cpu"), []byte("mem")}) {
				t.Fatalf("unexpected names: %v", names)
			}
		})

		t.Run("NEQREGEX", func(t *testing.T) {
			names, err := idx.MeasurementNamesByExpr(influxql.MustParseExpr(`country !~ /^u/`))
			if err != nil {
				t.Fatal(err)
			} else if !reflect.DeepEqual(names, [][]byte{[]byte("cpu"), []byte("disk")}) {
				t.Fatalf("unexpected names: %v", names)
			}
		})
	})
}

// Ensure index can return a list of matching measurements.
func TestIndex_MeasurementNamesByRegex(t *testing.T) {
	idx := MustOpenIndex()
	defer idx.Close()

	// Add series to index.
	if err := idx.CreateSeriesSliceIfNotExists([]Series{
		{Name: []byte("cpu")},
		{Name: []byte("disk")},
		{Name: []byte("mem")},
	}); err != nil {
		t.Fatal(err)
	}

	// Retrieve measurements by regex.
	idx.Run(t, func(t *testing.T) {
		names, err := idx.MeasurementNamesByRegex(regexp.MustCompile(`cpu|mem`))
		if err != nil {
			t.Fatal(err)
		} else if !reflect.DeepEqual(names, [][]byte{[]byte("cpu"), []byte("mem")}) {
			t.Fatalf("unexpected names: %v", names)
		}
	})
}

// Ensure index can delete a measurement and all related keys, values, & series.
func TestIndex_DropMeasurement(t *testing.T) {
	idx := MustOpenIndex()
	defer idx.Close()

	// Add series to index.
	if err := idx.CreateSeriesSliceIfNotExists([]Series{
		{Name: []byte("cpu"), Tags: models.NewTags(map[string]string{"region": "east"})},
		{Name: []byte("cpu"), Tags: models.NewTags(map[string]string{"region": "west"})},
		{Name: []byte("disk"), Tags: models.NewTags(map[string]string{"region": "north"})},
		{Name: []byte("mem"), Tags: models.NewTags(map[string]string{"region": "west", "country": "us"})},
	}); err != nil {
		t.Fatal(err)
	}

	// Drop measurement.
	if err := idx.DropMeasurement([]byte("cpu")); err != nil {
		t.Fatal(err)
	}

	// Verify data is gone in each stage.
	idx.Run(t, func(t *testing.T) {
		// Verify measurement is gone.
		if v, err := idx.MeasurementExists([]byte("cpu")); err != nil {
			t.Fatal(err)
		} else if v {
			t.Fatal("expected no measurement")
		}

		// Obtain file set to perform lower level checks.
		fs := idx.RetainFileSet()
		defer fs.Release()

		// Verify tags & values are gone.
		if e := fs.TagKeyIterator([]byte("cpu")).Next(); e != nil && !e.Deleted() {
			t.Fatal("expected deleted tag key")
		}
		if itr := fs.TagValueIterator([]byte("cpu"), []byte("region")); itr != nil {
			t.Fatal("expected nil tag value iterator")
		}

	})
}

func TestIndex_Open(t *testing.T) {
	// Opening a fresh index should set the MANIFEST version to current version.
	idx := NewIndex()
	t.Run("open new index", func(t *testing.T) {
		if err := idx.Open(); err != nil {
			t.Fatal(err)
		}

		// Check version set appropriately.
		if got, exp := idx.Manifest().Version, 1; got != exp {
			t.Fatalf("got index version %d, expected %d", got, exp)
		}
	})

	// Reopening an open index should return an error.
	t.Run("reopen open index", func(t *testing.T) {
		err := idx.Open()
		if err == nil {
			idx.Close()
			t.Fatal("didn't get an error on reopen, but expected one")
		}
		idx.Close()
	})

	// Opening an incompatible index should return an error.
	incompatibleVersions := []int{-1, 0, 2}
	for _, v := range incompatibleVersions {
		t.Run(fmt.Sprintf("incompatible index version: %d", v), func(t *testing.T) {
			idx = NewIndex()
			// Manually create a MANIFEST file for an incompatible index version.
			mpath := filepath.Join(idx.Path, tsi1.ManifestFileName)
			m := tsi1.NewManifest()
			m.Levels = nil
			m.Version = v // Set example MANIFEST version.
			if err := tsi1.WriteManifestFile(mpath, m); err != nil {
				t.Fatal(err)
			}

			// Log the MANIFEST file.
			data, err := ioutil.ReadFile(mpath)
			if err != nil {
				panic(err)
			}
			t.Logf("Incompatible MANIFEST: %s", data)

			// Opening this index should return an error because the MANIFEST has an
			// incompatible version.
			err = idx.Open()
			if err != tsi1.ErrIncompatibleVersion {
				idx.Close()
				t.Fatalf("got error %v, expected %v", err, tsi1.ErrIncompatibleVersion)
			}
		})
	}
}

func TestIndex_Manifest(t *testing.T) {
	t.Run("current MANIFEST", func(t *testing.T) {
		idx := MustOpenIndex()
		if got, exp := idx.Manifest().Version, tsi1.Version; got != exp {
			t.Fatalf("got MANIFEST version %d, expected %d", got, exp)
		}
	})
}

// Index is a test wrapper for tsi1.Index.
type Index struct {
	*tsi1.Index
}

// NewIndex returns a new instance of Index at a temporary path.
func NewIndex() *Index {
	idx := &Index{Index: tsi1.NewIndex()}
	idx.Path = MustTempDir()
	return idx
}

// MustOpenIndex returns a new, open index. Panic on error.
func MustOpenIndex() *Index {
	idx := NewIndex()
	if err := idx.Open(); err != nil {
		panic(err)
	}
	return idx
}

// Close closes and removes the index directory.
func (idx *Index) Close() error {
	defer os.RemoveAll(idx.Path)
	return idx.Index.Close()
}

// Reopen closes and opens the index.
func (idx *Index) Reopen() error {
	if err := idx.Index.Close(); err != nil {
		return err
	}

	path := idx.Path
	idx.Index = tsi1.NewIndex()
	idx.Path = path
	if err := idx.Open(); err != nil {
		return err
	}
	return nil
}

// Run executes a subtest for each of several different states:
//
// - Immediately
// - After reopen
// - After compaction
// - After reopen again
//
// The index should always respond in the same fashion regardless of
// how data is stored. This helper allows the index to be easily tested
// in all major states.
func (idx *Index) Run(t *testing.T, fn func(t *testing.T)) {
	// Invoke immediately.
	t.Run("state=initial", fn)

	// Reopen and invoke again.
	if err := idx.Reopen(); err != nil {
		t.Fatalf("reopen error: %s", err)
	}
	t.Run("state=reopen", fn)

	// TODO: Request a compaction.
	// if err := idx.Compact(); err != nil {
	// 	t.Fatalf("compact error: %s", err)
	// }
	// t.Run("state=post-compaction", fn)

	// Reopen and invoke again.
	if err := idx.Reopen(); err != nil {
		t.Fatalf("post-compaction reopen error: %s", err)
	}
	t.Run("state=post-compaction-reopen", fn)
}

// CreateSeriesSliceIfNotExists creates multiple series at a time.
func (idx *Index) CreateSeriesSliceIfNotExists(a []Series) error {
	for i, s := range a {
		if err := idx.CreateSeriesIfNotExists(nil, s.Name, s.Tags); err != nil {
			return fmt.Errorf("i=%d, name=%s, tags=%v, err=%s", i, s.Name, s.Tags, err)
		}
	}
	return nil
}

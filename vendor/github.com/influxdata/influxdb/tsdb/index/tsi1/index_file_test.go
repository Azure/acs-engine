package tsi1_test

import (
	"bytes"
	"testing"

	"github.com/influxdata/influxdb/models"
	"github.com/influxdata/influxdb/tsdb/index/tsi1"
)

// Ensure a simple index file can be built and opened.
func TestCreateIndexFile(t *testing.T) {
	f, err := CreateIndexFile([]Series{
		{Name: []byte("cpu"), Tags: models.NewTags(map[string]string{"region": "east"})},
		{Name: []byte("cpu"), Tags: models.NewTags(map[string]string{"region": "west"})},
		{Name: []byte("mem"), Tags: models.NewTags(map[string]string{"region": "east"})},
	})
	if err != nil {
		t.Fatal(err)
	}

	if e := f.TagValueElem([]byte("cpu"), []byte("region"), []byte("west")); e == nil {
		t.Fatal("expected element")
	} else if n := e.(*tsi1.TagBlockValueElem).SeriesN(); n != 1 {
		t.Fatalf("unexpected series count: %d", n)
	}
}

// Ensure index file generation can be successfully built.
func TestGenerateIndexFile(t *testing.T) {
	// Build generated index file.
	f, err := GenerateIndexFile(10, 3, 4)
	if err != nil {
		t.Fatal(err)
	}

	// Verify that tag/value series can be fetched.
	if e := f.TagValueElem([]byte("measurement0"), []byte("key0"), []byte("value0")); e == nil {
		t.Fatal("expected element")
	} else if n := e.(*tsi1.TagBlockValueElem).SeriesN(); n == 0 {
		t.Fatal("expected series")
	}
}

func BenchmarkIndexFile_TagValueSeries(b *testing.B) {
	b.Run("M=1,K=2,V=3", func(b *testing.B) {
		benchmarkIndexFile_TagValueSeries(b, MustFindOrGenerateIndexFile(1, 2, 3))
	})
	b.Run("M=10,K=5,V=5", func(b *testing.B) {
		benchmarkIndexFile_TagValueSeries(b, MustFindOrGenerateIndexFile(10, 5, 5))
	})
	b.Run("M=10,K=7,V=5", func(b *testing.B) {
		benchmarkIndexFile_TagValueSeries(b, MustFindOrGenerateIndexFile(10, 7, 7))
	})
}

func benchmarkIndexFile_TagValueSeries(b *testing.B, idx *tsi1.IndexFile) {
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		if e := idx.TagValueElem([]byte("measurement0"), []byte("key0"), []byte("value0")); e == nil {
			b.Fatal("expected element")
		} else if e.(*tsi1.TagBlockValueElem).SeriesN() == 0 {
			b.Fatal("expected series")
		}
	}
}

// CreateIndexFile creates an index file with a given set of series.
func CreateIndexFile(series []Series) (*tsi1.IndexFile, error) {
	lf, err := CreateLogFile(series)
	if err != nil {
		return nil, err
	}

	// Write index file to buffer.
	var buf bytes.Buffer
	if _, err := lf.CompactTo(&buf, M, K); err != nil {
		return nil, err
	}

	// Load index file from buffer.
	var f tsi1.IndexFile
	if err := f.UnmarshalBinary(buf.Bytes()); err != nil {
		return nil, err
	}
	return &f, nil
}

// GenerateIndexFile generates an index file from a set of series based on the count arguments.
// Total series returned will equal measurementN * tagN * valueN.
func GenerateIndexFile(measurementN, tagN, valueN int) (*tsi1.IndexFile, error) {
	// Generate a new log file first.
	lf, err := GenerateLogFile(measurementN, tagN, valueN)
	if err != nil {
		return nil, err
	}

	// Compact log file to buffer.
	var buf bytes.Buffer
	if _, err := lf.CompactTo(&buf, M, K); err != nil {
		return nil, err
	}

	// Load index file from buffer.
	var f tsi1.IndexFile
	if err := f.UnmarshalBinary(buf.Bytes()); err != nil {
		return nil, err
	}
	return &f, nil
}

func MustGenerateIndexFile(measurementN, tagN, valueN int) *tsi1.IndexFile {
	f, err := GenerateIndexFile(measurementN, tagN, valueN)
	if err != nil {
		panic(err)
	}
	return f
}

var indexFileCache struct {
	MeasurementN int
	TagN         int
	ValueN       int

	IndexFile *tsi1.IndexFile
}

// MustFindOrGenerateIndexFile returns a cached index file or generates one if it doesn't exist.
func MustFindOrGenerateIndexFile(measurementN, tagN, valueN int) *tsi1.IndexFile {
	// Use cache if fields match and the index file has been generated.
	if indexFileCache.MeasurementN == measurementN &&
		indexFileCache.TagN == tagN &&
		indexFileCache.ValueN == valueN &&
		indexFileCache.IndexFile != nil {
		return indexFileCache.IndexFile
	}

	// Generate and cache.
	indexFileCache.MeasurementN = measurementN
	indexFileCache.TagN = tagN
	indexFileCache.ValueN = valueN
	indexFileCache.IndexFile = MustGenerateIndexFile(measurementN, tagN, valueN)
	return indexFileCache.IndexFile
}

func pow(x, y int) int {
	r := 1
	for i := 0; i < y; i++ {
		r *= x
	}
	return r
}

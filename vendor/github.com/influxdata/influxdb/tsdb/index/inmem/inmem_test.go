package inmem_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/influxdata/influxdb/models"
	"github.com/influxdata/influxdb/tsdb"
	"github.com/influxdata/influxdb/tsdb/index/inmem"
)

func createData(lo, hi int) (keys, names [][]byte, tags []models.Tags) {
	for i := lo; i < hi; i++ {
		keys = append(keys, []byte(fmt.Sprintf("m0,tag0=t%d", i)))
		names = append(names, []byte("m0"))
		var t models.Tags
		t.Set([]byte("tag0"), []byte(fmt.Sprintf("%d", i)))
		tags = append(tags, t)
	}
	return
}

func BenchmarkShardIndex_CreateSeriesListIfNotExists_MaxValuesExceeded(b *testing.B) {
	sfile := mustOpenSeriesFile()
	defer sfile.Close()
	opt := tsdb.EngineOptions{InmemIndex: inmem.NewIndex("foo", sfile.SeriesFile)}
	opt.Config.MaxValuesPerTag = 10
	si := inmem.NewShardIndex(1, tsdb.NewSeriesIDSet(), opt)
	si.Open()
	keys, names, tags := createData(0, 10)
	si.CreateSeriesListIfNotExists(keys, names, tags)
	b.ReportAllocs()
	b.ResetTimer()

	keys, names, tags = createData(9, 5010)
	for i := 0; i < b.N; i++ {
		si.CreateSeriesListIfNotExists(keys, names, tags)
	}
}

func BenchmarkShardIndex_CreateSeriesListIfNotExists_MaxValuesNotExceeded(b *testing.B) {
	sfile := mustOpenSeriesFile()
	defer sfile.Close()
	opt := tsdb.EngineOptions{InmemIndex: inmem.NewIndex("foo", sfile.SeriesFile)}
	opt.Config.MaxValuesPerTag = 100000
	si := inmem.NewShardIndex(1, tsdb.NewSeriesIDSet(), opt)
	si.Open()
	keys, names, tags := createData(0, 10)
	si.CreateSeriesListIfNotExists(keys, names, tags)
	b.ReportAllocs()
	b.ResetTimer()

	keys, names, tags = createData(9, 5010)
	for i := 0; i < b.N; i++ {
		si.CreateSeriesListIfNotExists(keys, names, tags)
	}
}

func BenchmarkShardIndex_CreateSeriesListIfNotExists_NoMaxValues(b *testing.B) {
	sfile := mustOpenSeriesFile()
	defer sfile.Close()
	opt := tsdb.EngineOptions{InmemIndex: inmem.NewIndex("foo", sfile.SeriesFile)}
	si := inmem.NewShardIndex(1, tsdb.NewSeriesIDSet(), opt)
	si.Open()
	keys, names, tags := createData(0, 10)
	si.CreateSeriesListIfNotExists(keys, names, tags)
	b.ReportAllocs()
	b.ResetTimer()

	keys, names, tags = createData(9, 5010)
	for i := 0; i < b.N; i++ {
		si.CreateSeriesListIfNotExists(keys, names, tags)
	}
}

func BenchmarkShardIndex_CreateSeriesListIfNotExists_MaxSeriesExceeded(b *testing.B) {
	sfile := mustOpenSeriesFile()
	defer sfile.Close()
	opt := tsdb.EngineOptions{InmemIndex: inmem.NewIndex("foo", sfile.SeriesFile)}
	opt.Config.MaxValuesPerTag = 0
	opt.Config.MaxSeriesPerDatabase = 10
	si := inmem.NewShardIndex(1, tsdb.NewSeriesIDSet(), opt)
	si.Open()
	keys, names, tags := createData(0, 10)
	si.CreateSeriesListIfNotExists(keys, names, tags)
	b.ReportAllocs()
	b.ResetTimer()

	keys, names, tags = createData(9, 5010)
	for i := 0; i < b.N; i++ {
		si.CreateSeriesListIfNotExists(keys, names, tags)
	}
}

func TestIndex_Bytes(t *testing.T) {
	sfile := mustOpenSeriesFile()
	defer sfile.Close()
	opt := tsdb.EngineOptions{InmemIndex: inmem.NewIndex("foo", sfile.SeriesFile)}
	si := inmem.NewShardIndex(1, tsdb.NewSeriesIDSet(), opt).(*inmem.ShardIndex)

	indexBaseBytes := si.Bytes()

	name := []byte("name")
	err := si.CreateSeriesIfNotExists(name, name, models.Tags{})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	indexNewBytes := si.Bytes()
	if indexBaseBytes >= indexNewBytes {
		t.Errorf("index Bytes(): want >%d, got %d", indexBaseBytes, indexNewBytes)
	}
}

// seriesFileWrapper is a test wrapper for tsdb.seriesFileWrapper.
type seriesFileWrapper struct {
	*tsdb.SeriesFile
}

// newSeriesFileWrapper returns a new instance of seriesFileWrapper with a temporary file path.
func newSeriesFileWrapper() *seriesFileWrapper {
	dir, err := ioutil.TempDir("", "tsdb-series-file-")
	if err != nil {
		panic(err)
	}
	return &seriesFileWrapper{SeriesFile: tsdb.NewSeriesFile(dir)}
}

// mustOpenSeriesFile returns a new, open instance of seriesFileWrapper. Panic on error.
func mustOpenSeriesFile() *seriesFileWrapper {
	f := newSeriesFileWrapper()
	if err := f.Open(); err != nil {
		panic(err)
	}
	return f
}

// Close closes the log file and removes it from disk.
func (f *seriesFileWrapper) Close() error {
	defer os.RemoveAll(f.Path())
	return f.SeriesFile.Close()
}

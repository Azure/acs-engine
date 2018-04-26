package filesystem

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Filesystem provides methods which are runnable on a bare filesystem or a
// tar.gz file
type Filesystem interface {
	WriteFile(filename string, data []byte, perm os.FileMode) error
	Close() error
}

type filesystem struct {
	name string
}

var _ Filesystem = &filesystem{}

// NewFilesystem returns a Filesystem interface backed by a bare filesystem
func NewFilesystem(name string) (Filesystem, error) {
	err := os.RemoveAll(name)
	if err != nil {
		return nil, err
	}

	err = os.MkdirAll(name, 0777)
	if err != nil {
		return nil, err
	}

	return &filesystem{name}, nil
}

func (f *filesystem) mkdirAll(name string, perm os.FileMode) error {
	return os.MkdirAll(name, perm)
}

func (f *filesystem) WriteFile(filename string, data []byte, perm os.FileMode) error {
	err := f.mkdirAll(filepath.Dir(filepath.Join(f.name, filename)), 0755)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filepath.Join(f.name, filename), data, perm)
}

func (filesystem) Close() error {
	return nil
}

type tgzfile struct {
	gz   *gzip.Writer
	tw   *tar.Writer
	now  time.Time
	dirs map[string]struct{}
}

var _ Filesystem = &tgzfile{}

// NewTGZFile returns a Filesystem interface backed by a tar.gz file
func NewTGZFile(w io.Writer) (Filesystem, error) {
	gz := gzip.NewWriter(w)
	tw := &tgzfile{
		gz:   gz,
		tw:   tar.NewWriter(gz),
		now:  time.Now(),
		dirs: map[string]struct{}{},
	}

	return tw, nil
}

func (t *tgzfile) mkdirAll(name string, perm os.FileMode) error {
	parts := strings.Split(name, "/")
	for i := 1; i < len(parts); i++ {
		name = filepath.Join(parts[:i]...)
		if _, exists := t.dirs[name]; exists {
			continue
		}
		err := t.tw.WriteHeader(&tar.Header{
			Name:     name,
			Mode:     int64(perm),
			ModTime:  t.now,
			Typeflag: tar.TypeDir,
			Uname:    "root",
			Gname:    "root",
		})
		if err != nil {
			return err
		}
		t.dirs[name] = struct{}{}
	}
	return nil
}

func (t *tgzfile) WriteFile(filename string, data []byte, perm os.FileMode) error {
	err := t.mkdirAll(filepath.Dir(filename), 0755)
	if err != nil {
		return err
	}

	err = t.tw.WriteHeader(&tar.Header{
		Name:     filename,
		Mode:     int64(perm),
		Size:     int64(len(data)),
		ModTime:  t.now,
		Typeflag: tar.TypeReg,
		Uname:    "root",
		Gname:    "root",
	})
	if err != nil {
		return err
	}

	_, err = t.tw.Write(data)
	return err
}

func (t *tgzfile) Close() error {
	err := t.tw.Close()
	if err != nil {
		return err
	}
	return t.gz.Close()
}

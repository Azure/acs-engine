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

// Fileinfo is a struct that holds User, Group, and Mode
type Fileinfo struct {
	User  string
	Group string
	Mode  os.FileMode
}

// Closer implements Close()
type Closer interface {
	Close() error
}

// Reader provides read-related methods which are runnable on a bare filesystem
// or a tar.gz file
type Reader interface {
	ReadFile(filename string) ([]byte, error)
}

// Writer provides write-related methods which are runnable on a bare filesystem
// or a tar.gz file
type Writer interface {
	Mkdir(filename string, fileInfo Fileinfo) error
	WriteFile(filename string, data []byte, fileInfo Fileinfo) error
}

// WriteCloser implements Writer and Closer
type WriteCloser interface {
	Writer
	Closer
}

// Filesystem implements Reader, Writer and Closer
type Filesystem interface {
	Reader
	Writer
	Closer
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

// Mkdir makes a directory.  Note that it does not chown/chgrp as that would
// require elevated privileges
func (f *filesystem) Mkdir(name string, fileInfo Fileinfo) error {
	return os.Mkdir(name, fileInfo.Mode)
}

func (f *filesystem) mkdirAll(name string) error {
	return os.MkdirAll(name, 0755)
}

func (f *filesystem) WriteFile(filename string, data []byte, fileInfo Fileinfo) error {
	filePath := filepath.Join(f.name, filename)
	err := f.mkdirAll(filepath.Dir(filePath))
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filePath, data, fileInfo.Mode)
}

func (f *filesystem) ReadFile(filename string) ([]byte, error) {
	return ioutil.ReadFile(filepath.Join(f.name, filename))
}

func (filesystem) Close() error {
	return nil
}

type tgzwriter struct {
	gz   *gzip.Writer
	tw   *tar.Writer
	now  time.Time
	dirs map[string]struct{}
}

var _ WriteCloser = &tgzwriter{}

// NewTGZWriter returns a WriteCloser interface backed by a tar.gz file
func NewTGZWriter(w io.Writer) (WriteCloser, error) {
	gz := gzip.NewWriter(w)
	tw := &tgzwriter{
		gz:   gz,
		tw:   tar.NewWriter(gz),
		now:  time.Now(),
		dirs: map[string]struct{}{},
	}
	return tw, nil
}

func (t *tgzwriter) Mkdir(name string, fileInfo Fileinfo) error {
	if _, exists := t.dirs[name]; exists {
		return &os.PathError{Op: "mkdir", Path: name}
	}

	err := t.tw.WriteHeader(&tar.Header{
		Name:     name,
		Mode:     int64(fileInfo.Mode),
		ModTime:  t.now,
		Typeflag: tar.TypeDir,
		Uname:    fileInfo.User,
		Gname:    fileInfo.Group,
	})
	if err != nil {
		return err
	}
	t.dirs[name] = struct{}{}

	return nil
}

func (t *tgzwriter) mkdirAll(name string) error {
	parts := strings.Split(name, "/")
	for i := 1; i < len(parts); i++ {
		name = filepath.Join(parts[:i]...)
		if _, exists := t.dirs[name]; exists {
			continue
		}

		err := t.Mkdir(name, Fileinfo{Mode: 0755, User: "root", Group: "root"})
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *tgzwriter) WriteFile(filename string, data []byte, fileInfo Fileinfo) error {
	err := t.mkdirAll(filepath.Dir(filename))
	if err != nil {
		return err
	}

	err = t.tw.WriteHeader(&tar.Header{
		Name:     filename,
		Mode:     int64(fileInfo.Mode),
		Size:     int64(len(data)),
		ModTime:  t.now,
		Typeflag: tar.TypeReg,
		Uname:    fileInfo.User,
		Gname:    fileInfo.Group,
	})
	if err != nil {
		return err
	}

	_, err = t.tw.Write(data)
	return err
}

func (t *tgzwriter) Close() error {
	err := t.tw.Close()
	if err != nil {
		return err
	}
	return t.gz.Close()
}

type tgzreader struct {
	r io.ReadSeeker
}

var _ Reader = &tgzreader{}

// NewTGZReader returns a Reader interface backed by a tar.gz file
func NewTGZReader(r io.ReadSeeker) (Reader, error) {
	return &tgzreader{r: r}, nil
}

func (t *tgzreader) ReadFile(filename string) ([]byte, error) {
	_, err := t.r.Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}
	gz, err := gzip.NewReader(t.r)
	if err != nil {
		return nil, err
	}
	tr := tar.NewReader(gz)
	for {
		h, err := tr.Next()
		if err != nil {
			return nil, err
		}

		if h.Name == filename {
			return ioutil.ReadAll(tr)
		}
	}
}

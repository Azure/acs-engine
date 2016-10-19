package acsengine

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
)

// ACSEngineFileLoader represents the file loader used by ACSEngine
type ACSEngineFileLoader struct {
	filenameByteMap map[string][]byte
}

// Asset implements acsengine.AssetLoader
func (a *ACSEngineFileLoader) Asset(filename string) ([]byte, error) {
	if val, ok := a.filenameByteMap[filename]; ok {
		return val, nil
	}
	return nil, fmt.Errorf("file %s not found", filename)
}

// InitializeACSEngineFileLoader loads all files in the parts directory into memory
func InitializeACSEngineFileLoader(partsDirectory string) (*ACSEngineFileLoader, error) {
	a := &ACSEngineFileLoader{
		filenameByteMap: map[string][]byte{},
	}
	e := a.scanFilesInDirectory(partsDirectory, partsDirectory)
	if e != nil {
		return nil, e
	}
	return a, nil
}

func (a *ACSEngineFileLoader) scanFilesInDirectory(dir string, base string) error {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			if e := a.scanFilesInDirectory(filepath.Join(dir, file.Name()), base); e != nil {
				return e
			}
		} else {
			filename := filepath.Join(dir, file.Name())

			b, err := ioutil.ReadFile(filename)
			if err != nil {
				return err
			}
			rel, err2 := filepath.Rel(base, filename)
			if err2 != nil {
				return err
			}
			a.filenameByteMap[rel] = b
		}
	}
	return nil
}

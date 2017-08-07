package interpolatorwriter

import (
	"fmt"
	"github.com/Azure/acs-engine/pkg/interpolator"
	"io/ioutil"
	"os"
	"path"
	"sync"
)

// InterpolatorWriter represents a writer that will write interpolator data to disk
type InterpolatorWriter struct {
	outputDirectory string
	templateName    string
	parametersName  string
	interpolator    interpolator.Interpolator
	otherFiles      map[string][]byte
}

// NewInterpolatorWriter will return a new InterpolatorWriter that has been initialized
func NewInterpolatorWriter(outputDirectory, templateName, parametersName string, i interpolator.Interpolator) *InterpolatorWriter {
	return &InterpolatorWriter{
		outputDirectory: outputDirectory,
		templateName:    templateName,
		parametersName:  parametersName,
		interpolator:    i,
		otherFiles:      make(map[string][]byte),
	}
}

// AddFile will add an arbitrary file to also write to the configured directory
func (i *InterpolatorWriter) AddFile(name string, buffer []byte) {
	i.otherFiles[name] = buffer
}

// InterpolatorWriterMutex will prevent concurrent disk operations with an InterpolatorWriter
var InterpolatorWriterMutex sync.Mutex

// Write will write the contents of the InterpolatorWriter to disk
func (i *InterpolatorWriter) Write() error {
	InterpolatorWriterMutex.Lock()
	defer InterpolatorWriterMutex.Unlock()

	// Output directory
	ensureDirectory(i.outputDirectory)

	// Template
	templateBuffer, err := i.interpolator.GetTemplate()
	if err != nil || templateBuffer == nil {
		return fmt.Errorf("Error getting template buffer, or empty buffer: %v", err)
	}
	err = writeFile(templateBuffer, i.templateName, i.outputDirectory)
	if err != nil {
		return fmt.Errorf("Unable to write file [%s]: %v", i.templateName, err)
	}

	// Parameters
	parametersBuffer, err := i.interpolator.GetParameters()
	if err != nil || templateBuffer == nil {
		return fmt.Errorf("Error getting parameters buffer, or empty buffer: %v", err)
	}
	err = writeFile(parametersBuffer, i.parametersName, i.outputDirectory)
	if err != nil {
		return fmt.Errorf("Unable to write file [%s]: %v", i.templateName, err)
	}

	// Other files
	for name, buffer := range i.otherFiles {
		err := writeFile(buffer, name, i.outputDirectory)
		if err != nil {
			return fmt.Errorf("Unable to write file [%s]: %v", i.templateName, err)
		}
	}
	return nil
}

func writeFile(buffer []byte, name, directory string) error {
	absolutePath := path.Join(directory, name)
	if err := ioutil.WriteFile(absolutePath, buffer, 0600); err != nil {
		return err
	}
	return nil
}

func ensureDirectory(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0700)
		if err != nil {
			return fmt.Errorf("Error creating directory [%s]: %s", dir, err.Error())
		}
	}
	return nil
}

package acsengine

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Azure/acs-engine/pkg/api"
)

const (
	TestDataDir = "./testdata"
)

func TestExpected(t *testing.T) {
	// iterate the test data directory
	apiModelTestFiles := &[]APIModelTestFile{}
	if e := IterateTestFilesDirectory(TestDataDir, apiModelTestFiles); e != nil {
		t.Error(e.Error())
		return
	}

	for _, tuple := range *apiModelTestFiles {
		containerService, version, err := api.LoadContainerServiceFromFile(tuple.APIModelFilename)
		if err != nil {
			t.Errorf("Loading file %s got error: %s", tuple.APIModelFilename, err.Error())
			continue
		}
		expectedJson, e1 := ioutil.ReadFile(tuple.GetExpectedArmTemplateFilename())
		if e1 != nil {
			t.Error(e1.Error())
			continue
		}

		expectedParams, e2 := ioutil.ReadFile(tuple.GetExpectedArmTemplateParamsFilename())
		if e2 != nil {
			t.Error(e2.Error())
			continue
		}
		expectedJsonStr := strings.Replace(string(expectedJson), "\r", "", -1)
		expectedParamsStr := strings.Replace(string(expectedParams), "\r", "", -1)

		// test the output container service 3 times:
		// 1. first time tests loaded containerService
		// 2. second time tests generated containerService
		// 3. third time tests the generated containerService from the generated containerService
		templateGenerator, e3 := InitializeTemplateGenerator(false)
		if e3 != nil {
			t.Error(e3.Error())
			continue
		}
		for i := 0; i < 3; i++ {
			armTemplate, params, certsGenerated, err := templateGenerator.GenerateTemplate(containerService)
			if err != nil {
				t.Error(fmt.Errorf("error in file %s: %s", tuple.APIModelFilename, err.Error()))
				continue
			}
			ppArmTemplate, e1 := PrettyPrintArmTemplate(armTemplate)
			if e1 != nil {
				t.Error(armTemplate)
				t.Error(fmt.Errorf("error in file %s: %s", tuple.APIModelFilename, e1.Error()))
				break
			}

			ppParams, e2 := PrettyPrintJSON(params)
			if e2 != nil {
				t.Error(fmt.Errorf("error in file %s: %s", tuple.APIModelFilename, e2.Error()))
				continue
			}

			if certsGenerated == true {
				t.Errorf("cert generation unexpected for %s", containerService.Properties.OrchestratorProfile.OrchestratorType)
			}

			if !bytes.Equal([]byte(expectedJsonStr), []byte(ppArmTemplate)) {
				diffstr, differr := tuple.WriteArmTemplateErrFilename([]byte(ppArmTemplate))
				if differr != nil {
					diffstr += differr.Error()
				}
				t.Errorf("generated output different from expected for model %s: '%s'", tuple.GetExpectedArmTemplateFilename(), diffstr)
			}

			if !bytes.Equal([]byte(expectedParamsStr), []byte(ppParams)) {
				diffstr, differr := tuple.WriteArmTemplateParamsErrFilename([]byte(ppParams))
				if differr != nil {
					diffstr += differr.Error()
				}
				t.Errorf("generated parameters different from expected for model %s: '%s'", tuple.GetExpectedArmTemplateParamsFilename(), diffstr)
			}

			b, err := api.SerializeContainerService(containerService, version)
			if err != nil {
				t.Error(err)
			}
			containerService, version, err = api.DeserializeContainerService(b)
			if err != nil {
				t.Error(err)
			}
		}
	}
}

// APIModelTestFile holds the test file name and knows how to find the expected files
type APIModelTestFile struct {
	APIModelFilename string
}

// GetExpectedArmTemplateFilename returns the expected ARM template output for the model file
func (a *APIModelTestFile) GetExpectedArmTemplateFilename() string {
	j := strings.LastIndex(a.APIModelFilename, filepath.Ext(a.APIModelFilename))
	basename := a.APIModelFilename[:j]
	return fmt.Sprintf("%s_expected.json", basename)
}

// WriteArmTemplateErrFilename writes out an error file to sit parallel for comparison
func (a *APIModelTestFile) WriteArmTemplateErrFilename(contents []byte) (string, error) {
	filename := fmt.Sprintf("%s.err", a.GetExpectedArmTemplateFilename())
	if err := ioutil.WriteFile(filename, contents, 0600); err != nil {
		return "", err
	}
	return fmt.Sprintf("%s written for diff", filename), nil
}

// GetExpectedArmTemplateParamsFilename returns the expected ARM parameters output for the model file
func (a *APIModelTestFile) GetExpectedArmTemplateParamsFilename() string {
	j := strings.LastIndex(a.APIModelFilename, filepath.Ext(a.APIModelFilename))
	basename := a.APIModelFilename[:j]
	return fmt.Sprintf("%s_expected_params.json", basename)
}

// WriteArmTemplateParamsErrFilename writes out an error file to sit parallel for comparison
func (a *APIModelTestFile) WriteArmTemplateParamsErrFilename(contents []byte) (string, error) {
	filename := fmt.Sprintf("%s.err", a.GetExpectedArmTemplateParamsFilename())
	if err := ioutil.WriteFile(filename, contents, 0600); err != nil {
		return "", err
	}
	return fmt.Sprintf("%s written for diff", filename), nil
}

// IterateTestFilesDirectory iterates the test data directory adding api model files to the test file slice.
func IterateTestFilesDirectory(directory string, APIModelTestFiles *[]APIModelTestFile) error {
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			if e := IterateTestFilesDirectory(filepath.Join(directory, file.Name()), APIModelTestFiles); e != nil {
				return e
			}
		} else {
			if !strings.Contains(file.Name(), "_expected") {
				tuple := &APIModelTestFile{}
				tuple.APIModelFilename = filepath.Join(directory, file.Name())
				if _, ferr := os.Stat(tuple.GetExpectedArmTemplateFilename()); os.IsNotExist(ferr) {
					return fmt.Errorf("expected file '%s' is missing", tuple.GetExpectedArmTemplateFilename())
				}
				if _, ferr := os.Stat(tuple.GetExpectedArmTemplateParamsFilename()); os.IsNotExist(ferr) {
					return fmt.Errorf("expected file '%s' is missing", tuple.GetExpectedArmTemplateParamsFilename())
				}
				*APIModelTestFiles = append(*APIModelTestFiles, *tuple)
			}
		}
	}
	return nil
}

package acsengine

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Azure/acs-engine/pkg/api"
	"github.com/Azure/acs-engine/pkg/api/v20160330"
	"github.com/Azure/acs-engine/pkg/api/vlabs"
	. "github.com/onsi/gomega"
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
		containerService, version, err := api.LoadContainerServiceFromFile(tuple.APIModelFilename, true)
		if err != nil {
			t.Errorf("Loading file %s got error: %s", tuple.APIModelFilename, err.Error())
			continue
		}

		if version != vlabs.APIVersion && version != v20160330.APIVersion {
			// Set CertificateProfile here to avoid a new one generated.
			// Kubernetes template needs certificate profile to match expected template
			// API versions other than vlabs don't expose CertificateProfile
			// API versions after v20160330 supports Kubernetes
			containerService.Properties.CertificateProfile = &api.CertificateProfile{}
			addTestCertificateProfile(containerService.Properties.CertificateProfile)
		}

		isClassicMode := false
		if strings.Contains(tuple.APIModelFilename, "_classicmode") {
			isClassicMode = true
		}

		// test the output container service 3 times:
		// 1. first time tests loaded containerService
		// 2. second time tests generated containerService
		// 3. third time tests the generated containerService from the generated containerService
		templateGenerator, e3 := InitializeTemplateGenerator(isClassicMode)
		if e3 != nil {
			t.Error(e3.Error())
			continue
		}

		armTemplate, params, certsGenerated, err := templateGenerator.GenerateTemplate(containerService)
		if err != nil {
			t.Error(fmt.Errorf("error in file %s: %s", tuple.APIModelFilename, err.Error()))
			continue
		}

		expectedPpArmTemplate, e1 := PrettyPrintArmTemplate(armTemplate)
		if e1 != nil {
			t.Error(armTemplate)
			t.Error(fmt.Errorf("error in file %s: %s", tuple.APIModelFilename, e1.Error()))
			break
		}

		expectedPpParams, e2 := PrettyPrintJSON(params)
		if e2 != nil {
			t.Error(fmt.Errorf("error in file %s: %s", tuple.APIModelFilename, e2.Error()))
			continue
		}

		if certsGenerated == true {
			t.Errorf("cert generation unexpected for %s", containerService.Properties.OrchestratorProfile.OrchestratorType)
		}

		for i := 0; i < 3; i++ {
			armTemplate, params, certsGenerated, err := templateGenerator.GenerateTemplate(containerService)
			if err != nil {
				t.Error(fmt.Errorf("error in file %s: %s", tuple.APIModelFilename, err.Error()))
				continue
			}
			generatedPpArmTemplate, e1 := PrettyPrintArmTemplate(armTemplate)
			if e1 != nil {
				t.Error(fmt.Errorf("error in file %s: %s", tuple.APIModelFilename, e1.Error()))
				continue
			}

			generatedPpParams, e2 := PrettyPrintJSON(params)
			if e2 != nil {
				t.Error(fmt.Errorf("error in file %s: %s", tuple.APIModelFilename, e2.Error()))
				continue
			}

			if certsGenerated == true {
				t.Errorf("cert generation unexpected for %s", containerService.Properties.OrchestratorProfile.OrchestratorType)
			}

			if !bytes.Equal([]byte(expectedPpArmTemplate), []byte(generatedPpArmTemplate)) {
				diffstr, differr := tuple.WriteArmTemplateErrFilename([]byte(generatedPpArmTemplate))
				if differr != nil {
					diffstr += differr.Error()
				}
				t.Errorf("generated output different from expected for model %s: '%s'", tuple.APIModelFilename, diffstr)
			}

			if !bytes.Equal([]byte(expectedPpParams), []byte(generatedPpParams)) {
				diffstr, differr := tuple.WriteArmTemplateParamsErrFilename([]byte(generatedPpParams))
				if differr != nil {
					diffstr += differr.Error()
				}
				t.Errorf("generated parameters different from expected for model %s: '%s'", tuple.APIModelFilename, diffstr)
			}

			b, err := api.SerializeContainerService(containerService, version)
			if err != nil {
				t.Error(err)
			}
			containerService, version, err = api.DeserializeContainerService(b, true)
			if err != nil {
				t.Error(err)
			}
			if version != vlabs.APIVersion && version != v20160330.APIVersion {
				// Set CertificateProfile here to avoid a new one generated.
				// Kubernetes template needs certificate profile to match expected template
				// API versions other than vlabs don't expose CertificateProfile
				// API versions after v20160330 supports Kubernetes
				containerService.Properties.CertificateProfile = &api.CertificateProfile{}
				addTestCertificateProfile(containerService.Properties.CertificateProfile)
			}
		}
	}
}

// APIModelTestFile holds the test file name and knows how to find the expected files
type APIModelTestFile struct {
	APIModelFilename string
}

// WriteArmTemplateErrFilename writes out an error file to sit parallel for comparison
func (a *APIModelTestFile) WriteArmTemplateErrFilename(contents []byte) (string, error) {
	filename := fmt.Sprintf("%s_expected.err", a.APIModelFilename)
	if err := ioutil.WriteFile(filename, contents, 0600); err != nil {
		return "", err
	}
	return fmt.Sprintf("%s written for diff", filename), nil
}

// WriteArmTemplateParamsErrFilename writes out an error file to sit parallel for comparison
func (a *APIModelTestFile) WriteArmTemplateParamsErrFilename(contents []byte) (string, error) {
	filename := fmt.Sprintf("%s_expected_params.err", a.APIModelFilename)
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
			if !strings.Contains(file.Name(), "_expected") && strings.HasSuffix(file.Name(), ".json") {
				tuple := &APIModelTestFile{}
				tuple.APIModelFilename = filepath.Join(directory, file.Name())
				*APIModelTestFiles = append(*APIModelTestFiles, *tuple)
			}
		}
	}
	return nil
}

// addTestCertificateProfile add certificate artifacts for test purpose
func addTestCertificateProfile(api *api.CertificateProfile) {
	api.CaCertificate = "caCertificate"
	api.CaPrivateKey = "caPrivateKey"
	api.APIServerCertificate = "apiServerCertificate"
	api.APIServerPrivateKey = "apiServerPrivateKey"
	api.ClientCertificate = "clientCertificate"
	api.ClientPrivateKey = "clientPrivateKey"
	api.KubeConfigCertificate = "kubeConfigCertificate"
	api.KubeConfigPrivateKey = "kubeConfigPrivateKey"
}

func TestVersionOrdinal(t *testing.T) {
	RegisterTestingT(t)

	v170 := api.OrchestratorVersion("1.7.0")
	v166 := api.OrchestratorVersion("1.6.6")
	v162 := api.OrchestratorVersion("1.6.2")
	v160 := api.OrchestratorVersion("1.6.0")
	v153 := api.OrchestratorVersion("1.5.3")
	v16 := api.OrchestratorVersion("1.6")

	Expect(v166 < v170).To(BeTrue())
	Expect(v166 > v162).To(BeTrue())
	Expect(v162 < v166).To(BeTrue())
	Expect(v162 > v160).To(BeTrue())
	Expect(v160 < v162).To(BeTrue())
	Expect(v153 < v160).To(BeTrue())

	//testing with different version length
	Expect(v16 < v162).To(BeTrue())
	Expect(v16 > v153).To(BeTrue())

}

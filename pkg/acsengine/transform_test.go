// +build ignore

package acsengine

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/Sirupsen/logrus"
	. "github.com/onsi/gomega"
)

func TestNormalizeForVMSSScaling(t *testing.T) {
	RegisterTestingT(t)
	logger := logrus.New()
	fileContents, e := ioutil.ReadFile("./testFiles/vmssTemplate.json")
	Expect(e).To(BeNil())
	expectedFileContents, e := ioutil.ReadFile("./testFiles/vmssScaleTemplate_expected.json")
	Expect(e).To(BeNil())
	templateJSON := string(fileContents)
	var template interface{}
	json.Unmarshal([]byte(templateJSON), &template)
	templateMap := template.(map[string]interface{})
	e = NormalizeForVMSSScaling(logger, templateMap)
	Expect(e).To(BeNil())
	ValidateTemplate(templateMap, expectedFileContents, "vmssScaleTemplate")
}

func TestNormalizeForK8sVMASScalingUp(t *testing.T) {
	RegisterTestingT(t)
	logger := logrus.New()
	fileContents, e := ioutil.ReadFile("./testFiles/k8sVMASTemplate.json")
	Expect(e).To(BeNil())
	expectedFileContents, e := ioutil.ReadFile("./testFiles/k8sVMASTemplate_expected.json")
	Expect(e).To(BeNil())
	templateJSON := string(fileContents)
	var template interface{}
	json.Unmarshal([]byte(templateJSON), &template)
	templateMap := template.(map[string]interface{})
	e = NormalizeForK8sVMASScalingUp(logger, templateMap)
	Expect(e).To(BeNil())
	ValidateTemplate(templateMap, expectedFileContents, "k8sVMASTemplate")
}

func TestNormalizeResourcesForK8sMasterUpgrade(t *testing.T) {
	RegisterTestingT(t)
	logger := logrus.New()
	fileContents, e := ioutil.ReadFile("./testFiles/k8sVMASTemplate.json")
	Expect(e).To(BeNil())
	expectedFileContents, e := ioutil.ReadFile("./testFiles/k8sVMASTemplate_expected.json")
	Expect(e).To(BeNil())
	templateJSON := string(fileContents)
	var template interface{}
	json.Unmarshal([]byte(templateJSON), &template)
	templateMap := template.(map[string]interface{})
	e = NormalizeResourcesForK8sMasterUpgrade(logger, templateMap)
	Expect(e).To(BeNil())
	ValidateTemplate(templateMap, expectedFileContents, "k8sVMASTemplate")
}

func TestNormalizeResourcesForK8sAgentUpgrade(t *testing.T) {
	RegisterTestingT(t)
	logger := logrus.New()
	fileContents, e := ioutil.ReadFile("./testFiles/k8sVMASTemplate.json")
	Expect(e).To(BeNil())
	expectedFileContents, e := ioutil.ReadFile("./testFiles/k8sVMASTemplate_expected.json")
	Expect(e).To(BeNil())
	templateJSON := string(fileContents)
	var template interface{}
	json.Unmarshal([]byte(templateJSON), &template)
	templateMap := template.(map[string]interface{})
	e = NormalizeResourcesForK8sAgentUpgrade(logger, templateMap)
	Expect(e).To(BeNil())
	ValidateTemplate(templateMap, expectedFileContents, "k8sVMASTemplate")
}

func ValidateTemplate(templateMap map[string]interface{}, expectedFileContents []byte, testFileName string) {
	output, e := json.Marshal(templateMap)
	Expect(e).To(BeNil())
	prettyOutput, e := PrettyPrintArmTemplate(string(output))
	Expect(e).To(BeNil())
	prettyExpectedOutput, e := PrettyPrintArmTemplate(string(expectedFileContents))
	Expect(e).To(BeNil())
	if prettyOutput != prettyExpectedOutput {
		ioutil.WriteFile(fmt.Sprintf("./testFiles/%s.failure.json", testFileName), []byte(prettyOutput), 0600)
	}
	Expect(prettyOutput).To(Equal(prettyExpectedOutput))
}

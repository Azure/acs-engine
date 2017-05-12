package acsengine

import (
	"core/log"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/Azure/acs-engine/pkg/acsengine"
	. "github.com/onsi/gomega"
)

func TestNormalizeForVMSSScaling(t *testing.T) {
	RegisterTestingT(t)
	logger := log.IntializeTestLogger()
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
	logger := log.IntializeTestLogger()
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

func ValidateTemplate(templateMap map[string]interface{}, expectedFileContents []byte, testFileName string) {
	output, e := json.Marshal(templateMap)
	Expect(e).To(BeNil())
	prettyOutput, e := acsengine.PrettyPrintArmTemplate(string(output))
	Expect(e).To(BeNil())
	prettyExpectedOutput, e := acsengine.PrettyPrintArmTemplate(string(expectedFileContents))
	Expect(e).To(BeNil())
	if prettyOutput != prettyExpectedOutput {
		ioutil.WriteFile(fmt.Sprintf("./testFiles/%s.failure.json", testFileName), []byte(prettyOutput), 0600)
	}
	Expect(prettyOutput).To(Equal(prettyExpectedOutput))
}

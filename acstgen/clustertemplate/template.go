package clustertemplate

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"text/template"

	"./../api/vlabs"
)

// AcsClusterTemplate represents the full ACS template
type AcsClusterTemplate struct {
	// Parameters are the boiler plate parameters to the template
	Parameters string
	// Variables represent template variables
	Variables TemplateVariables
	// Resources represent template resources
	Resources TemplateResources
	// Outputs represent template outputs
	Outputs TemplateOutputs
}

// TemplateVariables represents the variables section of the template
type TemplateVariables struct {
	// Master describes the master
	Master string
	// Agents describes 0 or more agent pools
	Agents string
	// Diagnostics describes the diagnostics section
	Diagnostics string
}

// TemplateResources represents the resources section of the template
type TemplateResources struct {
	// Master describes the master
	Master string
	// Agents describes 0 or more agent pools
	Agents string
	// Diagnostics describes the diagnostics section
	Diagnostics string
}

// TemplateOutputs represents the output section of the template
type TemplateOutputs struct {
	// Master describes the master
	Master string
	// Agents describes 0 or more agent pools
	Agents string
	// Diagnostics describes the diagnostics section
	Diagnostics string
}

const (
	baseFile = "base.t"
)

var requiredFiles = []string{baseFile}

// VerifyFiles verifies that the required template files exist
func VerifyFiles(partsDirectory string) error {
	for _, file := range requiredFiles {
		templateFile := path.Join(partsDirectory, file)
		if _, err := os.Stat(templateFile); os.IsNotExist(err) {
			return fmt.Errorf("template file %s does not exist, did you specify the correct template directory?", templateFile)
		}
	}
	return nil
}

// GenerateTemplate generates the template from the API Model
func GenerateTemplate(acsCluster *vlabs.AcsCluster, partsDirectory string) (string, error) {
	var err error
	var templ *template.Template

	var acsTemplate = &AcsClusterTemplate{}
	if err = acsTemplate.generateContent(acsCluster); err != nil {
		return "", err
	}

	basetemplate, e := ioutil.ReadFile(path.Join(partsDirectory, baseFile))
	if e != nil {
		return "", fmt.Errorf("Error reading file %s: %s", basetemplate, e.Error())
	}
	s := string(basetemplate)
	if templ, err = template.New("acs template").Parse(s); err != nil {
		return "", err
	}
	var b bytes.Buffer
	if err = templ.Execute(&b, acsTemplate); err != nil {
		return "", err
	}

	return b.String(), nil
}

type templateObject interface {
	generateContent(acsCluster *vlabs.AcsCluster) error
}

func (a *AcsClusterTemplate) generateContent(acsCluster *vlabs.AcsCluster) error {
	a.Parameters = acsCluster.OrchestratorProfile.OrchestratorType
	if err := a.Variables.generateContent(acsCluster); err != nil {
		return fmt.Errorf("error generating TemplateVariables content: %s", err.Error())
	}
	if err := a.Resources.generateContent(acsCluster); err != nil {
		return fmt.Errorf("error generating TemplateResources content: %s", err.Error())
	}
	if err := a.Outputs.generateContent(acsCluster); err != nil {
		return fmt.Errorf("error generating TemplateOutputs content: %s", err.Error())
	}
	return nil
}

func (t *TemplateVariables) generateContent(acsCluster *vlabs.AcsCluster) error {
	t.Master = "TemplateVariables master"
	t.Agents = "TemplateVariables agents"
	t.Diagnostics = "TemplateVariables diagnostics"

	return nil
}

func (t *TemplateResources) generateContent(acsCluster *vlabs.AcsCluster) error {
	t.Master = "TemplateResources master"
	t.Agents = "TemplateResources agents"
	t.Diagnostics = "TemplateResources diagnostics"

	return nil
}

func (t *TemplateOutputs) generateContent(acsCluster *vlabs.AcsCluster) error {
	t.Master = "TemplateOutputs master"
	t.Agents = "TemplateOutputs agents"
	t.Diagnostics = "TemplateOutputs diagnostics"

	return nil
}

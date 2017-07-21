package agentpool

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Azure/acs-engine/pkg/acsengine"
	"github.com/Azure/acs-engine/pkg/api/kubernetesagentpool"
	"github.com/Azure/acs-engine/pkg/interpolator"
	"github.com/prometheus/common/log"
	txttemplate "text/template"
)

type Interpolator struct {
	agentPool         *kubernetesagentpool.AgentPool
	interpolated      bool
	template          []byte
	parameters        []byte
	templateDirectory string
}

func NewAgentPoolInterpolator(agentPool *kubernetesagentpool.AgentPool, templateDirectory string) interpolator.Interpolator {
	return &Interpolator{
		agentPool:         agentPool,
		templateDirectory: templateDirectory,
	}
}

func (i *Interpolator) Interpolate() error {

	// Init template
	templ := txttemplate.New("agentpool").Funcs(getTemplateFuncMap(i.agentPool))

	// Load files
	files, err := acsengine.AssetDir(i.templateDirectory)
	if err != nil {
		return fmt.Errorf("Unable to parse asset dir [kubernetes/agentpool]: %v", err)
	}

	// Parse files
	for _, file := range files {
		log.Infof("Loading file [%s]", file)
		bytes, err := acsengine.Asset(fmt.Sprintf("%s/%s", i.templateDirectory, file))
		if err != nil {
			return fmt.Errorf("Error reading file %s, Error: %s", file, err.Error())
		}
		if _, err = templ.New(file).Parse(string(bytes)); err != nil {
			return fmt.Errorf("Unable to parse template: %v", err)
		}
	}

	var b bytes.Buffer
	if err = templ.ExecuteTemplate(&b, "azuredeploy.json", i.agentPool.Properties); err != nil {
		return fmt.Errorf("Unable to execute template: %v", err)
	}

	var parametersMap map[string]interface{}
	parametersMap, err = getParameters(i.agentPool)
	if err != nil {
		return fmt.Errorf("Unable to get parameteres: %v", err)
	}
	var parameterBytes []byte
	parameterBytes, err = json.Marshal(parametersMap)
	var out bytes.Buffer
	err = json.Indent(&out, parameterBytes, "", "  ")
	if err != nil {
		return fmt.Errorf("Unable to pretty print json: %v", err)
	}
	parameterBytes = out.Bytes()


	if err != nil {
		return fmt.Errorf("Unable to marshal parameters map: %v", err)
	}

	// Cache the template
	i.template = b.Bytes()
	i.parameters = parameterBytes
	i.interpolated = true
	return nil
}

func (i *Interpolator) GetTemplate() ([]byte, error) {
	if i.interpolated == false {
		return []byte(""), fmt.Errorf("Unable to get template before calling Interpolate()")
	}
	return i.template, nil
}

func (i *Interpolator) GetParameters() ([]byte, error) {
	if i.interpolated == false {
		return []byte(""), fmt.Errorf("Unable to get template before calling Interpolate()")
	}
	return i.parameters, nil
}

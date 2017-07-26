package agentpool

import (
	"github.com/Azure/acs-engine/pkg/api/kubernetesagentpool"
	"text/template"
	"github.com/Azure/acs-engine/pkg/acsengine"
	"github.com/prometheus/common/log"
	"strings"
	//"encoding/base64"
	"encoding/base64"
)

// getTemplateFuncMap is where we can define functions used interpolating our code. Please try
// to use these as sparingly as possible, and if you catch yourself wanting to write a lot of them,
// it might be a sign that we need to pull your work into it's own implementation of an Interpolator
func getTemplateFuncMap(agentPool *kubernetesagentpool.AgentPool) map[string]interface{} {
	return template.FuncMap{
		"GetCustomAgentPoolData": func() string {

			customDataBytes, err := acsengine.Asset("kubernetes/agentpool/customData")
			if err != nil {
				log.Warnf("Unable to get customData: %v", err)
				return ""
			}
			kubeletServiceBytes, err := acsengine.Asset("kubernetes/agentpool/kubelet.service")
			if err != nil {
				log.Warnf("Unable to get kubelet.service: %v", err)
				return ""
			}
			base64KubeletService := base64.StdEncoding.EncodeToString(kubeletServiceBytes)
			fullStr := strings.Replace(string(customDataBytes), "KUBELET_SERVICE_BASE64", 	base64KubeletService, 1)
			str, err := formatJsonNewlineBytes(fullStr)
			if err != nil {
				log.Warnf("Unable to format bytes for ARM: %v", err)
				return ""
			}
			return str
		},
	}
}


func formatJsonNewlineBytes(str string) (string, error) {
	str = strings.Replace(str, "\n", "\\n", -1)
	//str = base64.StdEncoding.EncodeToString([]byte(str))
	return str, nil
}








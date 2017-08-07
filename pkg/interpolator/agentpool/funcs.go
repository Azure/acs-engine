package agentpool

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"github.com/Azure/acs-engine/pkg/acsengine"
	"github.com/Azure/acs-engine/pkg/api/kubernetesagentpool"
	"github.com/prometheus/common/log"
	"strings"
	"text/template"
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
			base64KubeletService := base64GzipStr(string(kubeletServiceBytes))
			fullStr := strings.Replace(string(customDataBytes), "KUBELET_SERVICE_BASE64", base64KubeletService, 1)
			provisionScriptBytes, err := acsengine.Asset("kubernetes/agentpool/provisionScript")
			if err != nil {
				log.Warnf("Unable to get provisionScript: %v", err)
				return ""
			}
			base64ProvisionScript := base64GzipStr(string(provisionScriptBytes))
			fullStr = strings.Replace(fullStr, "PROVISION_SCRIPT_BASE64", base64ProvisionScript, 1)
			str, err := formatJSONNewlineBytes(fullStr)
			if err != nil {
				log.Warnf("Unable to format bytes for ARM: %v", err)
				return ""
			}
			return str
		},
	}
}

func formatJSONNewlineBytes(str string) (string, error) {
	str = strings.Replace(str, "\n", "\\n", -1)
	return str, nil
}

func base64GzipStr(str string) string {
	var gzipB bytes.Buffer
	w := gzip.NewWriter(&gzipB)
	w.Write([]byte(str))
	w.Close()
	return base64.StdEncoding.EncodeToString(gzipB.Bytes())
}

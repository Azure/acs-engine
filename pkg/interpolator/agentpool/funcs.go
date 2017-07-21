package agentpool

import (
	"github.com/Azure/acs-engine/pkg/api/kubernetesagentpool"
	"text/template"
)

// getTemplateFuncMap is where we can define functions used interpolating our code. Please try
// to use these as sparingly as possible, and if you catch yourself wanting to write a lot of them,
// it might be a sign that we need to pull your work into it's own implementation of an Interpolator
func getTemplateFuncMap(agentPool *kubernetesagentpool.AgentPool) map[string]interface{} {
	return template.FuncMap{
	}
}

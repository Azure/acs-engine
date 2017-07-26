package kubernetesagentpool

import "github.com/Azure/acs-engine/pkg/agentPoolOnlyApi/v20170831"

func Convertv20170831ToAgentPool(hostedMaster *v20170831.HostedMaster) *AgentPool {
	return &AgentPool{}
}
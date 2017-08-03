package kubernetesagentpool

import (
	"encoding/json"
	"fmt"
	"github.com/Azure/acs-engine/pkg/agentPoolOnlyApi/v20170831"
)

// Convertv20170831ToAgentPool will convert a Convertv20170831 object into AgentPool=
func Convertv20170831ToAgentPool(hostedMaster *v20170831.HostedMaster) (*AgentPool, error) {

	hostedMasterBytes, err := json.Marshal(hostedMaster)
	if err != nil {
		return nil, fmt.Errorf("Unable to marshal hosted master: %v", err)
	}
	agentPool := &AgentPool{}
	err = json.Unmarshal(hostedMasterBytes, agentPool)
	if err != nil {
		return nil, fmt.Errorf("Unable to unmarshal onto AgentPool: %v", err)
	}
	return agentPool, nil
}

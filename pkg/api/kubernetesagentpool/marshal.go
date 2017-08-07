package kubernetesagentpool

import (
	"encoding/json"
	"fmt"
	"github.com/Azure/acs-engine/pkg/api"
	"github.com/Azure/acs-engine/pkg/api/kubernetesagentpool/v20170727"
	"io/ioutil"
)

// LoadAgentPoolFromFile will attempt to load an AgentPool struct from a given JSON file
func LoadAgentPoolFromFile(jsonFile string) (*AgentPool, string, error) {
	contents, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		return nil, "", fmt.Errorf("error reading file %s: %s", jsonFile, err.Error())
	}
	return DeserializeAgentPool(contents)
}

// DeserializeAgentPool loads an ACS Cluster API Model, validates it, and returns the unversioned representation
func DeserializeAgentPool(contents []byte) (*AgentPool, string, error) {
	m := &api.TypeMeta{}
	if err := json.Unmarshal(contents, &m); err != nil {
		return nil, "", err
	}
	agentPool, err := LoadAgentPool(contents, m.APIVersion)
	if err != nil {
		return nil, "", err
	}
	return agentPool, m.APIVersion, nil
}

// LoadAgentPool will attempt to load the versioned API and then convert to the AgentPool struct (if necessary)
func LoadAgentPool(contents []byte, version string) (*AgentPool, error) {
	switch version {
	case v20170727.APIVersion:
		// We know this is a 1:1 so this should marshal, otherwise err
		agentPool := &AgentPool{}
		if err := json.Unmarshal(contents, &agentPool); err != nil {
			return nil, fmt.Errorf("Unable to marshal on AgentPool: %v", err)
		}
		return agentPool, nil
	default:
		return nil, fmt.Errorf("Invalid API version: %s", version)
	}
}

package kubernetesagentpool

import (
	"encoding/json"
	"fmt"
	"github.com/Azure/acs-engine/pkg/api"
	"io/ioutil"
)

func LoadAgentPoolFromFile(jsonFile string) (*AgentPool, string, error) {
	contents, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		return nil, "", fmt.Errorf("error reading file %s: %s", jsonFile, err.Error())
	}
	return DeserializeAgentPool(contents)
}

// DeserializeContainerService loads an ACS Cluster API Model, validates it, and returns the unversioned representation
func DeserializeAgentPool(contents []byte) (*AgentPool, string, error) {
	m := &api.TypeMeta{}
	if err := json.Unmarshal(contents, &m); err != nil {
		return nil, "", err
	}
	version := m.APIVersion
	agentPool := &AgentPool{}
	if err := json.Unmarshal(contents, &agentPool); err != nil {
		return nil, "", err
	}
	//service, err := LoadAgentPool(contents, version)
	return agentPool, version, nil
}

//func LoadAgentPool(contents []byte, version string) (*AgentPool, error) {
//
//	fmt.Println(string(contents))
//	fmt.Println(version)
//	os.Exit(1)
//	return nil, nil
//}

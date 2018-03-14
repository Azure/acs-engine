package api

import (
	"testing"

	"github.com/Azure/acs-engine/pkg/api/agentPoolOnlyApi/v20180331"
)

func TestConvertFromV20180331AddonProfile(t *testing.T) {
	addonName := "AddonFoo"
	p := map[string]v20180331.AddonProfile{
		addonName: {
			Enabled: true,
			Config: map[string]string{
				"opt1": "value1",
			},
		},
	}
	api := convertV20180331AgentPoolOnlyAddonProfiles(p)

	if len(api) != 1 {
		t.Error("there has to be one addon")
	}
	if _, ok := api[addonName]; !ok {
		t.Error("addon is not found")
	}
	if api[addonName].Enabled != true {
		t.Error("addon should be enabled")
	}
	v, ok := api[addonName].Config["opt1"]
	if !ok {
		t.Error("Addon config opt1 is not found")
	}
	if v != "value1" {
		t.Error("addon config value does not match")
	}
}

func TestIfMasterProfileIsMissingThenApiModelIsAgentPoolOnly(t *testing.T) {
	json := `
	{
		"apiVersion": "vlabs",
		"properties": {
			"dnsPrefix": "dp",
			"fqdn": "fqdn",
			"agentPoolProfiles": [],
			"servicePrincipalProfile": {}
		}
	}
	`
	isAgentPool := isAgentPoolOnlyClusterJSON([]byte(json))
	if !isAgentPool {
		t.Error("Expected JSON without masterProfile to be interpreted as agent pool, but it was not")
	}
}

func TestIfMasterProfileIsPresentThenApiModelIsFullCluster(t *testing.T) {
	json := `
	{
		"apiVersion": "vlabs",
		"properties": {
			"orchestratorProfile": {},
			"masterProfile": {},
			"agentPoolProfiles": [],
			"servicePrincipalProfile": {}
		}
	}
	`
	isAgentPool := isAgentPoolOnlyClusterJSON([]byte(json))
	if isAgentPool {
		t.Error("Expected JSON with masterProfile not to be interpreted as agent pool, but it was")
	}
}

package api

import (
	"testing"

	"github.com/Azure/acs-engine/pkg/api/agentPoolOnlyApi/v20180331"
)

func TestConvertToV20180331AddonProfile(t *testing.T) {
	addonName := "AddonFoo"
	api := map[string]AddonProfile{
		addonName: {
			Enabled: true,
			Config: map[string]string{
				"opt1": "value1",
			},
		},
	}

	p := make(map[string]v20180331.AddonProfile)
	convertAddonsProfileToV20180331AgentPoolOnly(api, p)

	if len(p) != 1 {
		t.Error("there has to be one addon")
	}
	if _, ok := p[addonName]; !ok {
		t.Error("addon is not found")
	}
	if p[addonName].Enabled != true {
		t.Error("addon should be enabled")
	}
	v, ok := p[addonName].Config["opt1"]
	if !ok {
		t.Error("Addon config opt1 is not found")
	}
	if v != "value1" {
		t.Error("addon config value does not match")
	}
}

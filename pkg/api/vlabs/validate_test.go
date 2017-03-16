package vlabs

import (
	"testing"
)

func TestProperties_ValidateNetworkPolicy(t *testing.T) {
	p := &Properties{}
	p.OrchestratorProfile.OrchestratorType = Kubernetes

	for _, policy := range []string{"", "none", "calico"} {
		p.OrchestratorProfile.KubernetesConfig.NetworkPolicy = policy
		if err := p.ValidateNetworkPolicy(); err != nil {
			t.Errorf(
				"should not error on networkPolicy=\"%s\"",
				policy,
			)
		}
	}

	p.OrchestratorProfile.KubernetesConfig.NetworkPolicy = "not-existing"
	if err := p.ValidateNetworkPolicy(); err == nil {
		t.Errorf(
			"should error on invalid networkPolicy",
		)
	}

	p.OrchestratorProfile.KubernetesConfig.NetworkPolicy = "calico"
	p.AgentPoolProfiles = []AgentPoolProfile{
		{
			OSType: Windows,
		},
	}
	if err := p.ValidateNetworkPolicy(); err == nil {
		t.Errorf(
			"should error on calico for windows clusters",
		)
	}
}

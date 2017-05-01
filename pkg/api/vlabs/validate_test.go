package vlabs

import (
	"testing"
)

func TestProperties_ValidateNetworkPolicy(t *testing.T) {
	p := &Properties{}
	p.OrchestratorProfile = &OrchestratorProfile{}
	p.OrchestratorProfile.OrchestratorType = Kubernetes

	for _, policy := range NetworkPolicyValues {
		p.OrchestratorProfile.KubernetesConfig = &KubernetesConfig{}
		p.OrchestratorProfile.KubernetesConfig.NetworkPolicy = policy
		if err := p.validateNetworkPolicy(); err != nil {
			t.Errorf(
				"should not error on networkPolicy=\"%s\"",
				policy,
			)
		}
	}

	p.OrchestratorProfile.KubernetesConfig.NetworkPolicy = "not-existing"
	if err := p.validateNetworkPolicy(); err == nil {
		t.Errorf(
			"should error on invalid networkPolicy",
		)
	}

	p.OrchestratorProfile.KubernetesConfig.NetworkPolicy = "calico"
	p.AgentPoolProfiles = []*AgentPoolProfile{
		{
			OSType: Windows,
		},
	}
	if err := p.validateNetworkPolicy(); err == nil {
		t.Errorf(
			"should error on calico for windows clusters",
		)
	}
}

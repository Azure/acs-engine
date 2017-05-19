package vlabs

import (
	"testing"
)

func Test_OrchestratorProfile_Validate(t *testing.T) {
	o := &OrchestratorProfile{
		OrchestratorType: "DCOS",
		KubernetesConfig: &KubernetesConfig{},
	}

	if err := o.Validate(); err != nil {
		t.Errorf("should not error with empty object: %v", err)
	}

	o.KubernetesConfig.ClusterSubnet = "10.0.0.0/16"
	if err := o.Validate(); err == nil {
		t.Errorf("should error when KubernetesConfig populated for non-Kubernetes OrchestratorType")
	}
}

func Test_KubernetesConfig_Validate(t *testing.T) {
	c := KubernetesConfig{}

	if err := c.Validate(); err != nil {
		t.Errorf("should not error on empty KubernetesConfig: %v", err)
	}

	c.ClusterSubnet = "10.120.0.0/16"
	if err := c.Validate(); err != nil {
		t.Errorf("should not error on valid ClusterSubnet: %v", err)
	}

	c.ClusterSubnet = "10.16.x.0/invalid"
	if err := c.Validate(); err == nil {
		t.Error("should error on invalid ClusterSubnet")
	}
}

func Test_Properties_ValidateNetworkPolicy(t *testing.T) {
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

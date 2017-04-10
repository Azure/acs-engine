package vlabs

import (
	"testing"
)

func TestOrchestratorProfile_Validate(t *testing.T) {
	o := &OrchestratorProfile{
		OrchestratorType: "DCOS",
		KubernetesConfig: &KubernetesConfig{},
	}

	if err := o.Validate(); err != nil {
		t.Errorf("should not error with empty object: %v", err)
	}

	o.KubernetesConfig.ClusterCidr = "10.0.0.0/16"
	if err := o.Validate(); err == nil {
		t.Errorf("should error when KubernetesConfig populated for non-kube OrchestratorType")
	}
}

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
	p.AgentPoolProfiles = []AgentPoolProfile{
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

func TestKubernetesConfig_Validate(t *testing.T) {
	c := KubernetesConfig{}

	if err := c.Validate(); err != nil {
		t.Errorf("should not error on empty KubernetesConfig: %v", err)
	}

	c.ClusterCidr = "172.16.0.0/16"
	if err := c.Validate(); err != nil {
		t.Errorf("should not error on valid ClusterCidr: %v", err)
	}

	c.ClusterCidr = "172.16.0.0/a"
	if err := c.Validate(); err == nil {
		t.Error("should error on invalid ClusterCidr")
	}

	c = KubernetesConfig{
		DnsServiceIP: "192.168.0.10",
	}
	if err := c.Validate(); err == nil {
		t.Error("should error when DnsServiceIP but not ServiceCidr")
	}

	c = KubernetesConfig{
		ServiceCidr: "192.168.0.10/24",
	}
	if err := c.Validate(); err == nil {
		t.Error("should error when ServiceCidr but not DnsServiceIP")
	}

	c = KubernetesConfig{
		DnsServiceIP: "invalid",
		ServiceCidr:  "192.168.0.0/24",
	}
	if err := c.Validate(); err == nil {
		t.Error("should error when DnsServiceIP is invalid")
	}

	c = KubernetesConfig{
		DnsServiceIP: "192.168.1.10",
		ServiceCidr:  "192.168.0.0/not-a-len",
	}
	if err := c.Validate(); err == nil {
		t.Error("should error when ServiceCidr is invalid")
	}

	c = KubernetesConfig{
		DnsServiceIP: "192.168.1.10",
		ServiceCidr:  "192.168.0.0/24",
	}
	if err := c.Validate(); err == nil {
		t.Error("should error when DnsServiceIP is outside of ServiceCidr")
	}

	c = KubernetesConfig{
		DnsServiceIP: "172.99.255.255",
		ServiceCidr:  "172.99.0.1/16",
	}
	if err := c.Validate(); err == nil {
		t.Error("should error when DnsServiceIP is broadcast address of ServiceCidr")
	}

	c = KubernetesConfig{
		DnsServiceIP: "172.99.255.10",
		ServiceCidr:  "172.99.0.1/16",
	}
	if err := c.Validate(); err != nil {
		t.Error("should not error when DnsServiceIP and ServiceCidr are valid")
	}
}

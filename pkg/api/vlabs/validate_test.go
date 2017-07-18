package vlabs

import "testing"

const (
	ValidKubernetesNodeStatusUpdateFrequency        = "10s"
	ValidKubernetesCtrlMgrNodeMonitorGracePeriod    = "40s"
	ValidKubernetesCtrlMgrPodEvictionTimeout        = "5m0s"
	ValidKubernetesCtrlMgrRouteReconciliationPeriod = "10s"
	ValidKubernetesCloudProviderBackoff             = false
	ValidKubernetesCloudProviderBackoffRetries      = 6
	ValidKubernetesCloudProviderBackoffJitter       = 1
	ValidKubernetesCloudProviderBackoffDuration     = 5
	ValidKubernetesCloudProviderBackoffExponent     = 1.5
	ValidKubernetesCloudProviderRateLimit           = false
	ValidKubernetesCloudProviderRateLimitQPS        = 3
	ValidKubernetesCloudProviderRateLimitBucket     = 10
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
	// Tests that should pass across all versions
	for _, k8sVersion := range []OrchestratorVersion{Kubernetes153, Kubernetes157, Kubernetes160, Kubernetes162, Kubernetes166, Kubernetes170} {
		c := KubernetesConfig{}
		if err := c.Validate(k8sVersion); err != nil {
			t.Errorf("should not error on empty KubernetesConfig: %v, version %s", err, k8sVersion)
		}

		c = KubernetesConfig{
			ClusterSubnet:                    "10.120.0.0/16",
			DockerBridgeSubnet:               "10.120.1.0/16",
			NodeStatusUpdateFrequency:        ValidKubernetesNodeStatusUpdateFrequency,
			CtrlMgrNodeMonitorGracePeriod:    ValidKubernetesCtrlMgrNodeMonitorGracePeriod,
			CtrlMgrPodEvictionTimeout:        ValidKubernetesCtrlMgrPodEvictionTimeout,
			CtrlMgrRouteReconciliationPeriod: ValidKubernetesCtrlMgrRouteReconciliationPeriod,
			CloudProviderBackoff:             ValidKubernetesCloudProviderBackoff,
			CloudProviderBackoffRetries:      ValidKubernetesCloudProviderBackoffRetries,
			CloudProviderBackoffJitter:       ValidKubernetesCloudProviderBackoffJitter,
			CloudProviderBackoffDuration:     ValidKubernetesCloudProviderBackoffDuration,
			CloudProviderBackoffExponent:     ValidKubernetesCloudProviderBackoffExponent,
			CloudProviderRateLimit:           ValidKubernetesCloudProviderRateLimit,
			CloudProviderRateLimitQPS:        ValidKubernetesCloudProviderRateLimitQPS,
			CloudProviderRateLimitBucket:     ValidKubernetesCloudProviderRateLimitBucket,
		}
		if err := c.Validate(k8sVersion); err != nil {
			t.Errorf("should not error on a KubernetesConfig with valid param values: %v", err)
		}

		c = KubernetesConfig{
			ClusterSubnet: "10.16.x.0/invalid",
		}
		if err := c.Validate(k8sVersion); err == nil {
			t.Error("should error on invalid ClusterSubnet")
		}

		c = KubernetesConfig{
			DockerBridgeSubnet: "10.120.1.0/invalid",
		}
		if err := c.Validate(k8sVersion); err == nil {
			t.Error("should error on invalid DockerBridgeSubnet")
		}

		c = KubernetesConfig{
			NodeStatusUpdateFrequency: "invalid",
		}
		if err := c.Validate(k8sVersion); err == nil {
			t.Error("should error on invalid NodeStatusUpdateFrequency")
		}

		c = KubernetesConfig{
			CtrlMgrNodeMonitorGracePeriod: "invalid",
		}
		if err := c.Validate(k8sVersion); err == nil {
			t.Error("should error on invalid CtrlMgrNodeMonitorGracePeriod")
		}

		c = KubernetesConfig{
			NodeStatusUpdateFrequency:     "10s",
			CtrlMgrNodeMonitorGracePeriod: "30s",
		}
		if err := c.Validate(k8sVersion); err == nil {
			t.Error("should error when CtrlMgrRouteReconciliationPeriod is not sufficiently larger than NodeStatusUpdateFrequency")
		}

		c = KubernetesConfig{
			CtrlMgrPodEvictionTimeout: "invalid",
		}
		if err := c.Validate(k8sVersion); err == nil {
			t.Error("should error on invalid CtrlMgrPodEvictionTimeout")
		}

		c = KubernetesConfig{
			CtrlMgrRouteReconciliationPeriod: "invalid",
		}
		if err := c.Validate(k8sVersion); err == nil {
			t.Error("should error on invalid CtrlMgrRouteReconciliationPeriod")
		}
	}

	// Tests that apply to pre-1.6.6 versions
	for _, k8sVersion := range []OrchestratorVersion{Kubernetes153, Kubernetes157, Kubernetes160, Kubernetes162} {
		c := KubernetesConfig{
			CloudProviderBackoff:   true,
			CloudProviderRateLimit: true,
		}
		if err := c.Validate(k8sVersion); err == nil {
			t.Error("should error because backoff and rate limiting are not available before v1.6.6")
		}
	}

	// Tests that apply to 1.6.6 and later versions
	for _, k8sVersion := range []OrchestratorVersion{Kubernetes166, Kubernetes170} {
		c := KubernetesConfig{
			CloudProviderBackoff:   true,
			CloudProviderRateLimit: true,
		}
		if err := c.Validate(k8sVersion); err != nil {
			t.Error("should not error when basic backoff and rate limiting are set to true with no options")
		}
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

func Test_ServicePrincipalProfile_ValidateSecretOrKeyvaultSecretRef(t *testing.T) {

	t.Run("ServicePrincipalProfile with secret should pass", func(t *testing.T) {
		p := getK8sDefaultProperties()

		if err := p.Validate(); err != nil {
			t.Errorf("should not error %v", err)
		}
	})

	t.Run("ServicePrincipalProfile with KeyvaultSecretRef (with version) should pass", func(t *testing.T) {
		p := getK8sDefaultProperties()
		p.ServicePrincipalProfile.Secret = ""
		p.ServicePrincipalProfile.KeyvaultSecretRef = "/subscriptions/SUB-ID/resourceGroups/RG-NAME/providers/Microsoft.KeyVault/vaults/KV-NAME/secrets/secret-name/version"

		if err := p.Validate(); err != nil {
			t.Errorf("should not error %v", err)
		}
	})

	t.Run("ServicePrincipalProfile with KeyvaultSecretRef (without version) should pass", func(t *testing.T) {
		p := getK8sDefaultProperties()
		p.ServicePrincipalProfile.Secret = ""
		p.ServicePrincipalProfile.KeyvaultSecretRef = "/subscriptions/SUB-ID/resourceGroups/RG-NAME/providers/Microsoft.KeyVault/vaults/KV-NAME/secrets/secret-name>"

		if err := p.Validate(); err != nil {
			t.Errorf("should not error %v", err)
		}
	})

	t.Run("ServicePrincipalProfile with Secret and KeyvaultSecretRef should NOT pass", func(t *testing.T) {
		p := getK8sDefaultProperties()
		p.ServicePrincipalProfile.KeyvaultSecretRef = "/subscriptions/SUB-ID/resourceGroups/RG-NAME/providers/Microsoft.KeyVault/vaults/KV-NAME/secrets/secret-name/version"

		if err := p.Validate(); err == nil {
			t.Error("error should have occurred")
		}
	})

	t.Run("ServicePrincipalProfile with incorrect KeyvaultSecretRef format should NOT pass", func(t *testing.T) {
		p := getK8sDefaultProperties()
		p.ServicePrincipalProfile.Secret = ""
		p.ServicePrincipalProfile.KeyvaultSecretRef = "randomsecret"

		if err := p.Validate(); err == nil || err.Error() != "service principal client keyvault secret reference is of incorrect format" {
			t.Error("error should have occurred")
		}
	})
}

func getK8sDefaultProperties() *Properties {
	return &Properties{
		OrchestratorProfile: &OrchestratorProfile{
			OrchestratorType: Kubernetes,
		},
		MasterProfile: &MasterProfile{
			Count:     1,
			DNSPrefix: "foo",
			VMSize:    "Standard_DS2_v2",
		},
		AgentPoolProfiles: []*AgentPoolProfile{
			&AgentPoolProfile{
				Name:                "agentpool",
				VMSize:              "Standard_D2_v2",
				Count:               1,
				AvailabilityProfile: AvailabilitySet,
			},
		},
		LinuxProfile: &LinuxProfile{
			AdminUsername: "azureuser",
			SSH: struct {
				PublicKeys []PublicKey `json:"publicKeys"`
			}{
				PublicKeys: []PublicKey{{
					KeyData: "publickeydata",
				}},
			},
		},
		ServicePrincipalProfile: &ServicePrincipalProfile{
			ClientID: "clientID",
			Secret:   "clientSecret",
		},
	}
}

package vlabs

import (
	"testing"

	"github.com/Azure/acs-engine/pkg/api/common"
)

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

	o = &OrchestratorProfile{
		OrchestratorType: "Kubernetes",
		DcosConfig:       &DcosConfig{},
	}

	if err := o.Validate(); err != nil {
		t.Errorf("should not error with empty object: %v", err)
	}

	o.DcosConfig.DcosWindowsBootstrapURL = "http://www.microsoft.com"
	if err := o.Validate(); err == nil {
		t.Errorf("should error when DcosConfig populated for non-Kubernetes OrchestratorType")
	}
}

func Test_KubernetesConfig_Validate(t *testing.T) {
	// Tests that should pass across all releases
	for _, k8sRelease := range []string{common.KubernetesRelease1Dot5, common.KubernetesRelease1Dot6, common.KubernetesRelease1Dot7} {
		c := KubernetesConfig{}
		if err := c.Validate(k8sRelease); err != nil {
			t.Errorf("should not error on empty KubernetesConfig: %v, release %s", err, k8sRelease)
		}

		c = KubernetesConfig{
			ClusterSubnet:                    "10.120.0.0/16",
			DockerBridgeSubnet:               "10.120.1.0/16",
			MaxPods:                          42,
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
		if err := c.Validate(k8sRelease); err != nil {
			t.Errorf("should not error on a KubernetesConfig with valid param values: %v", err)
		}

		c = KubernetesConfig{
			ClusterSubnet: "10.16.x.0/invalid",
		}
		if err := c.Validate(k8sRelease); err == nil {
			t.Error("should error on invalid ClusterSubnet")
		}

		c = KubernetesConfig{
			DockerBridgeSubnet: "10.120.1.0/invalid",
		}
		if err := c.Validate(k8sRelease); err == nil {
			t.Error("should error on invalid DockerBridgeSubnet")
		}

		c = KubernetesConfig{
			MaxPods: KubernetesMinMaxPods - 1,
		}
		if err := c.Validate(k8sRelease); err == nil {
			t.Error("should error on invalid MaxPods")
		}

		c = KubernetesConfig{
			NodeStatusUpdateFrequency: "invalid",
		}
		if err := c.Validate(k8sRelease); err == nil {
			t.Error("should error on invalid NodeStatusUpdateFrequency")
		}

		c = KubernetesConfig{
			CtrlMgrNodeMonitorGracePeriod: "invalid",
		}
		if err := c.Validate(k8sRelease); err == nil {
			t.Error("should error on invalid CtrlMgrNodeMonitorGracePeriod")
		}

		c = KubernetesConfig{
			NodeStatusUpdateFrequency:     "10s",
			CtrlMgrNodeMonitorGracePeriod: "30s",
		}
		if err := c.Validate(k8sRelease); err == nil {
			t.Error("should error when CtrlMgrRouteReconciliationPeriod is not sufficiently larger than NodeStatusUpdateFrequency")
		}

		c = KubernetesConfig{
			CtrlMgrPodEvictionTimeout: "invalid",
		}
		if err := c.Validate(k8sRelease); err == nil {
			t.Error("should error on invalid CtrlMgrPodEvictionTimeout")
		}

		c = KubernetesConfig{
			CtrlMgrRouteReconciliationPeriod: "invalid",
		}
		if err := c.Validate(k8sRelease); err == nil {
			t.Error("should error on invalid CtrlMgrRouteReconciliationPeriod")
		}

		c = KubernetesConfig{
			DNSServiceIP: "192.168.0.10",
		}
		if err := c.Validate(k8sRelease); err == nil {
			t.Error("should error when DNSServiceIP but not ServiceCidr")
		}

		c = KubernetesConfig{
			ServiceCidr: "192.168.0.10/24",
		}
		if err := c.Validate(k8sRelease); err == nil {
			t.Error("should error when ServiceCidr but not DNSServiceIP")
		}

		c = KubernetesConfig{
			DNSServiceIP: "invalid",
			ServiceCidr:  "192.168.0.0/24",
		}
		if err := c.Validate(k8sRelease); err == nil {
			t.Error("should error when DNSServiceIP is invalid")
		}

		c = KubernetesConfig{
			DNSServiceIP: "192.168.1.10",
			ServiceCidr:  "192.168.0.0/not-a-len",
		}
		if err := c.Validate(k8sRelease); err == nil {
			t.Error("should error when ServiceCidr is invalid")
		}

		c = KubernetesConfig{
			DNSServiceIP: "192.168.1.10",
			ServiceCidr:  "192.168.0.0/24",
		}
		if err := c.Validate(k8sRelease); err == nil {
			t.Error("should error when DNSServiceIP is outside of ServiceCidr")
		}

		c = KubernetesConfig{
			DNSServiceIP: "172.99.255.255",
			ServiceCidr:  "172.99.0.1/16",
		}
		if err := c.Validate(k8sRelease); err == nil {
			t.Error("should error when DNSServiceIP is broadcast address of ServiceCidr")
		}

		c = KubernetesConfig{
			DNSServiceIP: "172.99.0.1",
			ServiceCidr:  "172.99.0.1/16",
		}
		if err := c.Validate(k8sRelease); err == nil {
			t.Error("should error when DNSServiceIP is first IP of ServiceCidr")
		}

		c = KubernetesConfig{
			DNSServiceIP: "172.99.255.10",
			ServiceCidr:  "172.99.0.1/16",
		}
		if err := c.Validate(k8sRelease); err != nil {
			t.Error("should not error when DNSServiceIP and ServiceCidr are valid")
		}
	}

	// Tests that apply to pre-1.6 releases
	for _, k8sRelease := range []string{common.KubernetesRelease1Dot5} {
		c := KubernetesConfig{
			CloudProviderBackoff:   true,
			CloudProviderRateLimit: true,
		}
		if err := c.Validate(k8sRelease); err == nil {
			t.Error("should error because backoff and rate limiting are not available before v1.6.6")
		}
	}

	// Tests that apply to 1.6 and later releases
	for _, k8sRelease := range []string{common.KubernetesRelease1Dot6, common.KubernetesRelease1Dot7} {
		c := KubernetesConfig{
			CloudProviderBackoff:   true,
			CloudProviderRateLimit: true,
		}
		if err := c.Validate(k8sRelease); err != nil {
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
		p.ServicePrincipalProfile.KeyvaultSecretRef = &KeyvaultSecretRef{
			VaultID:       "/subscriptions/SUB-ID/resourceGroups/RG-NAME/providers/Microsoft.KeyVault/vaults/KV-NAME",
			SecretName:    "secret-name",
			SecretVersion: "version",
		}
		if err := p.Validate(); err != nil {
			t.Errorf("should not error %v", err)
		}
	})

	t.Run("ServicePrincipalProfile with KeyvaultSecretRef (without version) should pass", func(t *testing.T) {
		p := getK8sDefaultProperties()
		p.ServicePrincipalProfile.Secret = ""
		p.ServicePrincipalProfile.KeyvaultSecretRef = &KeyvaultSecretRef{
			VaultID:    "/subscriptions/SUB-ID/resourceGroups/RG-NAME/providers/Microsoft.KeyVault/vaults/KV-NAME",
			SecretName: "secret-name",
		}

		if err := p.Validate(); err != nil {
			t.Errorf("should not error %v", err)
		}
	})

	t.Run("ServicePrincipalProfile with Secret and KeyvaultSecretRef should NOT pass", func(t *testing.T) {
		p := getK8sDefaultProperties()
		p.ServicePrincipalProfile.Secret = "secret"
		p.ServicePrincipalProfile.KeyvaultSecretRef = &KeyvaultSecretRef{
			VaultID:    "/subscriptions/SUB-ID/resourceGroups/RG-NAME/providers/Microsoft.KeyVault/vaults/KV-NAME",
			SecretName: "secret-name",
		}

		if err := p.Validate(); err == nil {
			t.Error("error should have occurred")
		}
	})

	t.Run("ServicePrincipalProfile with incorrect KeyvaultSecretRef format should NOT pass", func(t *testing.T) {
		p := getK8sDefaultProperties()
		p.ServicePrincipalProfile.Secret = ""
		p.ServicePrincipalProfile.KeyvaultSecretRef = &KeyvaultSecretRef{
			VaultID:    "randomID",
			SecretName: "secret-name",
		}

		if err := p.Validate(); err == nil || err.Error() != "service principal client keyvault secret reference is of incorrect format" {
			t.Error("error should have occurred")
		}
	})
}

func Test_AadProfile_Validate(t *testing.T) {
	t.Run("Valid aadProfile should pass", func(t *testing.T) {
		for _, aadProfile := range []AADProfile{
			{
				ClientAppID: "92444486-5bc3-4291-818b-d53ae480991b",
				ServerAppID: "403f018b-4d89-495b-b548-0cf9868cdb0a",
			},
			{
				ClientAppID: "92444486-5bc3-4291-818b-d53ae480991b",
				ServerAppID: "403f018b-4d89-495b-b548-0cf9868cdb0a",
				TenantID:    "feb784f6-7174-46da-aeae-da66e80c7a11",
			},
		} {
			if err := aadProfile.Validate(); err != nil {
				t.Errorf("should not error %v", err)
			}
		}
	})

	t.Run("Invalid aadProfiles should NOT pass", func(t *testing.T) {
		for _, aadProfile := range []AADProfile{
			{
				ClientAppID: "1",
				ServerAppID: "d",
			},
			{
				ClientAppID: "6a247d73-ae33-4559-8e5d-4001fdc17b15",
			},
			{
				ClientAppID: "92444486-5bc3-4291-818b-d53ae480991b",
				ServerAppID: "403f018b-4d89-495b-b548-0cf9868cdb0a",
				TenantID:    "1",
			},
			{},
		} {
			if err := aadProfile.Validate(); err == nil {
				t.Errorf("error should have occurred")
			}
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
			{
				Name:                "agentpool",
				VMSize:              "Standard_D2_v2",
				Count:               1,
				AvailabilityProfile: AvailabilitySet,
			},
		},
		LinuxProfile: &LinuxProfile{
			AdminUsername: "azureuser",
			SSH: struct {
				PublicKeys []PublicKey `json:"publicKeys" validate:"required,len=1"`
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

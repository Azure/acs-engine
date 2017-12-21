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

	o.KubernetesConfig.ClusterSubnet = "10.0.0.0/16"
	if err := o.Validate(false); err == nil {
		t.Errorf("should error when KubernetesConfig populated for non-Kubernetes OrchestratorType")
	}

	o = &OrchestratorProfile{
		OrchestratorType: "Kubernetes",
		DcosConfig:       &DcosConfig{},
	}

	if err := o.Validate(false); err != nil {
		t.Errorf("should not error with empty object: %v", err)
	}

	o.DcosConfig.DcosWindowsBootstrapURL = "http://www.microsoft.com"
	if err := o.Validate(false); err == nil {
		t.Errorf("should error when DcosConfig populated for non-Kubernetes OrchestratorType")
	}

	o.DcosConfig.DcosBootstrapURL = "http://www.microsoft.com"
	if err := o.Validate(false); err == nil {
		t.Errorf("should error when DcosConfig populated for non-Kubernetes OrchestratorType")
	}

	o = &OrchestratorProfile{
		OrchestratorType:    "Kubernetes",
		OrchestratorVersion: "1.7.3",
	}

	if err := o.Validate(false); err == nil {
		t.Errorf("should have failed on old patch version")
	}

	if err := o.Validate(true); err != nil {
		t.Errorf("should not have failed on old patch version during update valdiation")
	}
}

func Test_KubernetesConfig_Validate(t *testing.T) {
	// Tests that should pass across all versions
	for _, k8sVersion := range common.GetAllSupportedKubernetesVersions() {
		c := KubernetesConfig{}
		if err := c.Validate(k8sVersion); err != nil {
			t.Errorf("should not error on empty KubernetesConfig: %v, version %s", err, k8sVersion)
		}

		c = KubernetesConfig{
			ClusterSubnet:      "10.120.0.0/16",
			DockerBridgeSubnet: "10.120.1.0/16",
			MaxPods:            42,
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
			KubeletConfig: map[string]string{
				"--node-status-update-frequency": ValidKubernetesNodeStatusUpdateFrequency,
			},
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
			KubeletConfig: map[string]string{
				"--non-masquerade-cidr": "10.120.1.0/24",
			},
		}
		if err := c.Validate(k8sVersion); err != nil {
			t.Error("should not error on valid --non-masquerade-cidr")
		}

		c = KubernetesConfig{
			KubeletConfig: map[string]string{
				"--non-masquerade-cidr": "10.120.1.0/invalid",
			},
		}
		if err := c.Validate(k8sVersion); err == nil {
			t.Error("should error on invalid --non-masquerade-cidr")
		}

		c = KubernetesConfig{
			MaxPods: KubernetesMinMaxPods - 1,
		}
		if err := c.Validate(k8sVersion); err == nil {
			t.Error("should error on invalid MaxPods")
		}

		c = KubernetesConfig{
			KubeletConfig: map[string]string{
				"--node-status-update-frequency": "invalid",
			},
		}
		if err := c.Validate(k8sVersion); err == nil {
			t.Error("should error on invalid --node-status-update-frequency kubelet config")
		}

		c = KubernetesConfig{
			CtrlMgrNodeMonitorGracePeriod: "invalid",
		}
		if err := c.Validate(k8sVersion); err == nil {
			t.Error("should error on invalid CtrlMgrNodeMonitorGracePeriod")
		}

		c = KubernetesConfig{
			CtrlMgrNodeMonitorGracePeriod: "30s",
			KubeletConfig: map[string]string{
				"--node-status-update-frequency": "10s",
			},
		}
		if err := c.Validate(k8sVersion); err == nil {
			t.Error("should error when CtrlMgrRouteReconciliationPeriod is not sufficiently larger than --node-status-update-frequency kubelet config")
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

		c = KubernetesConfig{
			DNSServiceIP: "192.168.0.10",
		}
		if err := c.Validate(k8sVersion); err == nil {
			t.Error("should error when DNSServiceIP but not ServiceCidr")
		}

		c = KubernetesConfig{
			ServiceCidr: "192.168.0.10/24",
		}
		if err := c.Validate(k8sVersion); err == nil {
			t.Error("should error when ServiceCidr but not DNSServiceIP")
		}

		c = KubernetesConfig{
			DNSServiceIP: "invalid",
			ServiceCidr:  "192.168.0.0/24",
		}
		if err := c.Validate(k8sVersion); err == nil {
			t.Error("should error when DNSServiceIP is invalid")
		}

		c = KubernetesConfig{
			DNSServiceIP: "192.168.1.10",
			ServiceCidr:  "192.168.0.0/not-a-len",
		}
		if err := c.Validate(k8sVersion); err == nil {
			t.Error("should error when ServiceCidr is invalid")
		}

		c = KubernetesConfig{
			DNSServiceIP: "192.168.1.10",
			ServiceCidr:  "192.168.0.0/24",
		}
		if err := c.Validate(k8sVersion); err == nil {
			t.Error("should error when DNSServiceIP is outside of ServiceCidr")
		}

		c = KubernetesConfig{
			DNSServiceIP: "172.99.255.255",
			ServiceCidr:  "172.99.0.1/16",
		}
		if err := c.Validate(k8sVersion); err == nil {
			t.Error("should error when DNSServiceIP is broadcast address of ServiceCidr")
		}

		c = KubernetesConfig{
			DNSServiceIP: "172.99.0.1",
			ServiceCidr:  "172.99.0.1/16",
		}
		if err := c.Validate(k8sVersion); err == nil {
			t.Error("should error when DNSServiceIP is first IP of ServiceCidr")
		}

		c = KubernetesConfig{
			DNSServiceIP: "172.99.255.10",
			ServiceCidr:  "172.99.0.1/16",
		}
		if err := c.Validate(k8sVersion); err != nil {
			t.Error("should not error when DNSServiceIP and ServiceCidr are valid")
		}
	}

	// Tests that apply to pre-1.6 releases
	for _, k8sVersion := range []string{common.KubernetesVersion1Dot5Dot8} {
		c := KubernetesConfig{
			CloudProviderBackoff:   true,
			CloudProviderRateLimit: true,
		}
		if err := c.Validate(k8sVersion); err == nil {
			t.Error("should error because backoff and rate limiting are not available before v1.6.6")
		}
	}

	// Tests that apply to 1.6 and later releases
	for _, k8sVersion := range []string{common.KubernetesVersion1Dot6Dot11, common.KubernetesVersion1Dot6Dot12, common.KubernetesVersion1Dot6Dot13,
		common.KubernetesVersion1Dot7Dot7, common.KubernetesVersion1Dot7Dot9, common.KubernetesVersion1Dot7Dot10,
		common.KubernetesVersion1Dot8Dot1, common.KubernetesVersion1Dot8Dot2, common.KubernetesVersion1Dot8Dot4} {
		c := KubernetesConfig{
			CloudProviderBackoff:   true,
			CloudProviderRateLimit: true,
		}
		if err := c.Validate(k8sVersion); err != nil {
			t.Error("should not error when basic backoff and rate limiting are set to true with no options")
		}
	}

	trueVal := true
	// Tests that apply to pre-1.8 releases
	for _, k8sVersion := range []string{common.KubernetesVersion1Dot5Dot8, common.KubernetesVersion1Dot6Dot11, common.KubernetesVersion1Dot7Dot7} {
		c := KubernetesConfig{
			UseCloudControllerManager: &trueVal,
		}
		if err := c.Validate(k8sVersion); err == nil {
			t.Error("should error because UseCloudControllerManager is not available before v1.8")
		}
	}

	// Tests that apply to 1.8 and later releases
	for _, k8sVersion := range []string{common.KubernetesVersion1Dot8Dot1} {
		c := KubernetesConfig{
			UseCloudControllerManager: &trueVal,
		}
		if err := c.Validate(k8sVersion); err != nil {
			t.Error("should not error because UseCloudControllerManager is available since v1.8")
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

		if err := p.Validate(false); err != nil {
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
		if err := p.Validate(false); err != nil {
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

		if err := p.Validate(false); err != nil {
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

		if err := p.Validate(false); err == nil {
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

		if err := p.Validate(false); err == nil || err.Error() != "service principal client keyvault secret reference is of incorrect format" {
			t.Error("error should have occurred")
		}
	})
}

func TestValidateKubernetesLabelValue(t *testing.T) {

	validLabelValues := []string{"", "a", "a1", "this--valid--label--is--exactly--sixty--three--characters--long", "123456", "my-label_valid.com"}
	invalidLabelValues := []string{"a$$b", "-abc", "not.valid.", "This____long____label___is______sixty______four_____chararacters", "Label with spaces"}

	for _, l := range validLabelValues {
		if err := validateKubernetesLabelValue(l); err != nil {
			t.Fatalf("Label value %v should not return error: %v", l, err)
		}
	}

	for _, l := range invalidLabelValues {
		if err := validateKubernetesLabelValue(l); err == nil {
			t.Fatalf("Label value %v should return an error", l)
		}
	}
}

func TestValidateKubernetesLabelKey(t *testing.T) {

	validLabelKeys := []string{"a", "a1", "this--valid--label--is--exactly--sixty--three--characters--long", "123456", "my-label_valid.com", "foo.bar/name", "1.2321.324/key_name.foo", "valid.long.253.characters.label.key.prefix.12345678910.fooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooo/my-key"}
	invalidLabelKeys := []string{"", "a/b/c", ".startswithdot", "spaces in key", "foo/", "/name", "$.$/com", "too-long-254-characters-key-prefix-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------123/name", "wrong-slash\\foo"}

	for _, l := range validLabelKeys {
		if err := validateKubernetesLabelKey(l); err != nil {
			t.Fatalf("Label key %v should not return error: %v", l, err)
		}
	}

	for _, l := range invalidLabelKeys {
		if err := validateKubernetesLabelKey(l); err == nil {
			t.Fatalf("Label key %v should return an error", l)
		}
	}
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

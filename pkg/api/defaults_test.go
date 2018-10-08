package api

import (
	"encoding/base64"
	"reflect"
	"testing"

	"github.com/Azure/acs-engine/pkg/helpers"
)

func TestCertsAlreadyPresent(t *testing.T) {
	var cert *CertificateProfile

	result := certsAlreadyPresent(nil, 1)
	expected := map[string]bool{
		"ca":         false,
		"apiserver":  false,
		"client":     false,
		"kubeconfig": false,
		"etcd":       false,
	}

	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("certsAlreadyPresent() did not return false for all certs for a non-existent CertificateProfile")
	}
	cert = &CertificateProfile{}
	result = certsAlreadyPresent(cert, 1)
	expected = map[string]bool{
		"ca":         false,
		"apiserver":  false,
		"client":     false,
		"kubeconfig": false,
		"etcd":       false,
	}

	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("certsAlreadyPresent() did not return false for all certs for empty CertificateProfile")
	}
	cert = &CertificateProfile{
		APIServerCertificate: "a",
	}
	result = certsAlreadyPresent(cert, 1)
	expected = map[string]bool{
		"ca":         false,
		"apiserver":  false,
		"client":     false,
		"kubeconfig": false,
		"etcd":       false,
	}

	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("certsAlreadyPresent() did not return false for all certs for 1 cert in CertificateProfile")
	}

	cert = &CertificateProfile{
		APIServerCertificate:  "a",
		CaCertificate:         "c",
		CaPrivateKey:          "d",
		ClientCertificate:     "e",
		ClientPrivateKey:      "f",
		KubeConfigCertificate: "g",
		KubeConfigPrivateKey:  "h",
		EtcdClientCertificate: "i",
		EtcdClientPrivateKey:  "j",
		EtcdServerCertificate: "k",
		EtcdServerPrivateKey:  "l",
	}
	result = certsAlreadyPresent(cert, 3)
	expected = map[string]bool{
		"ca":         true,
		"apiserver":  false,
		"client":     true,
		"kubeconfig": true,
		"etcd":       false,
	}

	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("certsAlreadyPresent() did not return expected result for some certs in CertificateProfile")
	}
	cert = &CertificateProfile{
		APIServerCertificate:  "a",
		APIServerPrivateKey:   "b",
		CaCertificate:         "c",
		CaPrivateKey:          "d",
		ClientCertificate:     "e",
		ClientPrivateKey:      "f",
		KubeConfigCertificate: "g",
		KubeConfigPrivateKey:  "h",
		EtcdClientCertificate: "i",
		EtcdClientPrivateKey:  "j",
		EtcdServerCertificate: "k",
		EtcdServerPrivateKey:  "l",
		EtcdPeerCertificates:  []string{"0", "1", "2"},
		EtcdPeerPrivateKeys:   []string{"0", "1", "2"},
	}
	result = certsAlreadyPresent(cert, 3)
	expected = map[string]bool{
		"ca":         true,
		"apiserver":  true,
		"client":     true,
		"kubeconfig": true,
		"etcd":       true,
	}

	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("certsAlreadyPresent() did not return expected result for all certs in CertificateProfile")
	}
}

func TestSetMissingKubeletValues(t *testing.T) {
	config := &KubernetesConfig{}
	defaultKubeletConfig := map[string]string{
		"--network-plugin":               "1",
		"--pod-infra-container-image":    "2",
		"--max-pods":                     "3",
		"--eviction-hard":                "4",
		"--node-status-update-frequency": "5",
		"--image-gc-high-threshold":      "6",
		"--image-gc-low-threshold":       "7",
		"--non-masquerade-cidr":          "8",
		"--cloud-provider":               "9",
		"--pod-max-pids":                 "10",
	}
	setMissingKubeletValues(config, defaultKubeletConfig)
	for key, val := range defaultKubeletConfig {
		if config.KubeletConfig[key] != val {
			t.Fatalf("setMissingKubeletValue() did not return the expected value %s for key %s, instead returned: %s", val, key, config.KubeletConfig[key])
		}
	}

	config = &KubernetesConfig{
		KubeletConfig: map[string]string{
			"--network-plugin":            "a",
			"--pod-infra-container-image": "b",
			"--cloud-provider":            "c",
		},
	}
	expectedResult := map[string]string{
		"--network-plugin":               "a",
		"--pod-infra-container-image":    "b",
		"--max-pods":                     "3",
		"--eviction-hard":                "4",
		"--node-status-update-frequency": "5",
		"--image-gc-high-threshold":      "6",
		"--image-gc-low-threshold":       "7",
		"--non-masquerade-cidr":          "8",
		"--cloud-provider":               "c",
		"--pod-max-pids":                 "10",
	}
	setMissingKubeletValues(config, defaultKubeletConfig)
	for key, val := range expectedResult {
		if config.KubeletConfig[key] != val {
			t.Fatalf("setMissingKubeletValue() did not return the expected value %s for key %s, instead returned: %s", val, key, config.KubeletConfig[key])
		}
	}
	config = &KubernetesConfig{
		KubeletConfig: map[string]string{},
	}
	setMissingKubeletValues(config, defaultKubeletConfig)
	for key, val := range defaultKubeletConfig {
		if config.KubeletConfig[key] != val {
			t.Fatalf("setMissingKubeletValue() did not return the expected value %s for key %s, instead returned: %s", val, key, config.KubeletConfig[key])
		}
	}
}

func TestAddonsIndexByName(t *testing.T) {
	addonName := "testaddon"
	addons := []KubernetesAddon{
		getMockAddon(addonName),
	}
	i := getAddonsIndexByName(addons, addonName)
	if i != 0 {
		t.Fatalf("addonsIndexByName() did not return the expected index value 0, instead returned: %d", i)
	}
	i = getAddonsIndexByName(addons, "nonExistentAddonName")
	if i != -1 {
		t.Fatalf("addonsIndexByName() did not return -1 for a non-existent addon, instead returned: %d", i)
	}
}

func TestAssignDefaultAddonVals(t *testing.T) {
	addonName := "testaddon"
	customCPURequests := "60m"
	customMemoryRequests := "160Mi"
	customCPULimits := "40m"
	customMemoryLimits := "140Mi"
	// Verify that an addon with all custom values provided remains unmodified during default value assignment
	customAddon := KubernetesAddon{
		Name:    addonName,
		Enabled: helpers.PointerToBool(true),
		Containers: []KubernetesContainerSpec{
			{
				Name:           addonName,
				CPURequests:    customCPURequests,
				MemoryRequests: customMemoryRequests,
				CPULimits:      customCPULimits,
				MemoryLimits:   customMemoryLimits,
			},
		},
	}
	addonWithDefaults := getMockAddon(addonName)
	modifiedAddon := assignDefaultAddonVals(customAddon, addonWithDefaults)
	if modifiedAddon.Containers[0].Name != customAddon.Containers[0].Name {
		t.Fatalf("assignDefaultAddonVals() should not have modified Containers 'Name' value %s to %s,", customAddon.Containers[0].Name, modifiedAddon.Containers[0].Name)
	}
	if modifiedAddon.Containers[0].CPURequests != customAddon.Containers[0].CPURequests {
		t.Fatalf("assignDefaultAddonVals() should not have modified Containers 'CPURequests' value %s to %s,", customAddon.Containers[0].CPURequests, modifiedAddon.Containers[0].CPURequests)
	}
	if modifiedAddon.Containers[0].MemoryRequests != customAddon.Containers[0].MemoryRequests {
		t.Fatalf("assignDefaultAddonVals() should not have modified Containers 'MemoryRequests' value %s to %s,", customAddon.Containers[0].MemoryRequests, modifiedAddon.Containers[0].MemoryRequests)
	}
	if modifiedAddon.Containers[0].CPULimits != customAddon.Containers[0].CPULimits {
		t.Fatalf("assignDefaultAddonVals() should not have modified Containers 'CPULimits' value %s to %s,", customAddon.Containers[0].CPULimits, modifiedAddon.Containers[0].CPULimits)
	}
	if modifiedAddon.Containers[0].MemoryLimits != customAddon.Containers[0].MemoryLimits {
		t.Fatalf("assignDefaultAddonVals() should not have modified Containers 'MemoryLimits' value %s to %s,", customAddon.Containers[0].MemoryLimits, modifiedAddon.Containers[0].MemoryLimits)
	}

	// Verify that an addon with no custom values provided gets all the appropriate defaults
	customAddon = KubernetesAddon{
		Name:    addonName,
		Enabled: helpers.PointerToBool(true),
		Containers: []KubernetesContainerSpec{
			{
				Name: addonName,
			},
		},
	}
	modifiedAddon = assignDefaultAddonVals(customAddon, addonWithDefaults)
	if modifiedAddon.Containers[0].CPURequests != addonWithDefaults.Containers[0].CPURequests {
		t.Fatalf("assignDefaultAddonVals() should have assigned a default 'CPURequests' value of %s, instead assigned %s,", addonWithDefaults.Containers[0].CPURequests, modifiedAddon.Containers[0].CPURequests)
	}
	if modifiedAddon.Containers[0].MemoryRequests != addonWithDefaults.Containers[0].MemoryRequests {
		t.Fatalf("assignDefaultAddonVals() should have assigned a default 'MemoryRequests' value of %s, instead assigned %s,", addonWithDefaults.Containers[0].MemoryRequests, modifiedAddon.Containers[0].MemoryRequests)
	}
	if modifiedAddon.Containers[0].CPULimits != addonWithDefaults.Containers[0].CPULimits {
		t.Fatalf("assignDefaultAddonVals() should have assigned a default 'CPULimits' value of %s, instead assigned %s,", addonWithDefaults.Containers[0].CPULimits, modifiedAddon.Containers[0].CPULimits)
	}
	if modifiedAddon.Containers[0].MemoryLimits != addonWithDefaults.Containers[0].MemoryLimits {
		t.Fatalf("assignDefaultAddonVals() should have assigned a default 'MemoryLimits' value of %s, instead assigned %s,", addonWithDefaults.Containers[0].MemoryLimits, modifiedAddon.Containers[0].MemoryLimits)
	}

	// More checking to verify default interpolation
	customAddon = KubernetesAddon{
		Name:    addonName,
		Enabled: helpers.PointerToBool(true),
		Containers: []KubernetesContainerSpec{
			{
				Name:         addonName,
				CPURequests:  customCPURequests,
				MemoryLimits: customMemoryLimits,
			},
		},
	}
	modifiedAddon = assignDefaultAddonVals(customAddon, addonWithDefaults)
	if modifiedAddon.Containers[0].Name != customAddon.Containers[0].Name {
		t.Fatalf("assignDefaultAddonVals() should not have modified Containers 'Name' value %s to %s,", customAddon.Containers[0].Name, modifiedAddon.Containers[0].Name)
	}
	if modifiedAddon.Containers[0].MemoryRequests != addonWithDefaults.Containers[0].MemoryRequests {
		t.Fatalf("assignDefaultAddonVals() should have assigned a default 'MemoryRequests' value of %s, instead assigned %s,", addonWithDefaults.Containers[0].MemoryRequests, modifiedAddon.Containers[0].MemoryRequests)
	}
	if modifiedAddon.Containers[0].CPULimits != addonWithDefaults.Containers[0].CPULimits {
		t.Fatalf("assignDefaultAddonVals() should have assigned a default 'CPULimits' value of %s, instead assigned %s,", addonWithDefaults.Containers[0].CPULimits, modifiedAddon.Containers[0].CPULimits)
	}
	if modifiedAddon.Containers[0].MemoryLimits != customAddon.Containers[0].MemoryLimits {
		t.Fatalf("assignDefaultAddonVals() should not have modified Containers 'MemoryLimits' value %s to %s,", customAddon.Containers[0].MemoryLimits, modifiedAddon.Containers[0].MemoryLimits)
	}

	addonWithDefaults.Config = map[string]string{
		"os":    "Linux",
		"taint": "node.kubernetes.io/memory-pressure",
	}
	modifiedAddon = assignDefaultAddonVals(customAddon, addonWithDefaults)

	if modifiedAddon.Config["os"] != "Linux" {
		t.Error("assignDefaultAddonVals() should have added the default config property")
	}

	if modifiedAddon.Config["taint"] != "node.kubernetes.io/memory-pressure" {
		t.Error("assignDefaultAddonVals() should have added the default config property")
	}

}

func TestKubeletFeatureGatesEnsureFeatureGatesOnAgentsFor1_6_0(t *testing.T) {
	mockCS := getMockBaseContainerService("1.6.0")
	properties := mockCS.Properties

	// No KubernetesConfig.KubeletConfig set for MasterProfile or AgentProfile
	// so they will inherit the top-level config
	properties.OrchestratorProfile.KubernetesConfig = getKubernetesConfigWithFeatureGates("TopLevel=true")

	mockCS.setKubeletConfig()

	agentFeatureGates := properties.AgentPoolProfiles[0].KubernetesConfig.KubeletConfig["--feature-gates"]
	if agentFeatureGates != "TopLevel=true" {
		t.Fatalf("setKubeletConfig did not add 'TopLevel=true' for agent profile: expected 'TopLevel=true' got '%s'", agentFeatureGates)
	}

	// Verify that the TopLevel feature gate override has only been applied to the agents
	masterFeatureFates := properties.MasterProfile.KubernetesConfig.KubeletConfig["--feature-gates"]
	if masterFeatureFates != "TopLevel=true" {
		t.Fatalf("setKubeletConfig modified feature gates for master profile: expected 'TopLevel=true' got '%s'", agentFeatureGates)
	}
}

func TestKubeletFeatureGatesEnsureMasterAndAgentConfigUsedFor1_6_0(t *testing.T) {
	mockCS := getMockBaseContainerService("1.6.0")
	properties := mockCS.Properties

	// Set MasterProfile and AgentProfiles KubernetesConfig.KubeletConfig values
	// Verify that they are used instead of the top-level config
	properties.OrchestratorProfile.KubernetesConfig = getKubernetesConfigWithFeatureGates("TopLevel=true")
	properties.MasterProfile = &MasterProfile{KubernetesConfig: getKubernetesConfigWithFeatureGates("MasterLevel=true")}
	properties.AgentPoolProfiles[0].KubernetesConfig = getKubernetesConfigWithFeatureGates("AgentLevel=true")

	mockCS.setKubeletConfig()

	agentFeatureGates := properties.AgentPoolProfiles[0].KubernetesConfig.KubeletConfig["--feature-gates"]
	if agentFeatureGates != "AgentLevel=true" {
		t.Fatalf("setKubeletConfig agent profile: expected 'AgentLevel=true' got '%s'", agentFeatureGates)
	}

	// Verify that the TopLevel feature gate override has only been applied to the agents
	masterFeatureFates := properties.MasterProfile.KubernetesConfig.KubeletConfig["--feature-gates"]
	if masterFeatureFates != "MasterLevel=true" {
		t.Fatalf("setKubeletConfig master profile: expected 'MasterLevel=true' got '%s'", agentFeatureGates)
	}
}

func TestEtcdDiskSize(t *testing.T) {
	mockCS := getMockBaseContainerService("1.8.10")
	properties := mockCS.Properties
	properties.OrchestratorProfile.OrchestratorType = "Kubernetes"
	properties.MasterProfile.Count = 1
	mockCS.setOrchestratorDefaults(true)
	if properties.OrchestratorProfile.KubernetesConfig.EtcdDiskSizeGB != DefaultEtcdDiskSize {
		t.Fatalf("EtcdDiskSizeGB did not have the expected size, got %s, expected %s",
			properties.OrchestratorProfile.KubernetesConfig.EtcdDiskSizeGB, DefaultEtcdDiskSize)
	}

	mockCS = getMockBaseContainerService("1.8.10")
	properties = mockCS.Properties
	properties.OrchestratorProfile.OrchestratorType = "Kubernetes"
	properties.MasterProfile.Count = 5
	mockCS.setOrchestratorDefaults(true)
	if properties.OrchestratorProfile.KubernetesConfig.EtcdDiskSizeGB != DefaultEtcdDiskSizeGT3Nodes {
		t.Fatalf("EtcdDiskSizeGB did not have the expected size, got %s, expected %s",
			properties.OrchestratorProfile.KubernetesConfig.EtcdDiskSizeGB, DefaultEtcdDiskSizeGT3Nodes)
	}

	mockCS = getMockBaseContainerService("1.8.10")
	properties = mockCS.Properties
	properties.OrchestratorProfile.OrchestratorType = "Kubernetes"
	properties.MasterProfile.Count = 5
	properties.AgentPoolProfiles[0].Count = 6
	mockCS.setOrchestratorDefaults(true)
	if properties.OrchestratorProfile.KubernetesConfig.EtcdDiskSizeGB != DefaultEtcdDiskSizeGT10Nodes {
		t.Fatalf("EtcdDiskSizeGB did not have the expected size, got %s, expected %s",
			properties.OrchestratorProfile.KubernetesConfig.EtcdDiskSizeGB, DefaultEtcdDiskSizeGT10Nodes)
	}

	mockCS = getMockBaseContainerService("1.8.10")
	properties = mockCS.Properties
	properties.OrchestratorProfile.OrchestratorType = "Kubernetes"
	properties.MasterProfile.Count = 5
	properties.AgentPoolProfiles[0].Count = 16
	mockCS.setOrchestratorDefaults(true)
	if properties.OrchestratorProfile.KubernetesConfig.EtcdDiskSizeGB != DefaultEtcdDiskSizeGT20Nodes {
		t.Fatalf("EtcdDiskSizeGB did not have the expected size, got %s, expected %s",
			properties.OrchestratorProfile.KubernetesConfig.EtcdDiskSizeGB, DefaultEtcdDiskSizeGT20Nodes)
	}

	mockCS = getMockBaseContainerService("1.8.10")
	properties = mockCS.Properties
	properties.OrchestratorProfile.OrchestratorType = "Kubernetes"
	properties.MasterProfile.Count = 5
	properties.AgentPoolProfiles[0].Count = 50
	customEtcdDiskSize := "512"
	properties.OrchestratorProfile.KubernetesConfig.EtcdDiskSizeGB = customEtcdDiskSize
	mockCS.setOrchestratorDefaults(true)
	if properties.OrchestratorProfile.KubernetesConfig.EtcdDiskSizeGB != customEtcdDiskSize {
		t.Fatalf("EtcdDiskSizeGB did not have the expected size, got %s, expected %s",
			properties.OrchestratorProfile.KubernetesConfig.EtcdDiskSizeGB, customEtcdDiskSize)
	}
}

func TestGenerateEtcdEncryptionKey(t *testing.T) {
	key1 := generateEtcdEncryptionKey()
	key2 := generateEtcdEncryptionKey()
	if key1 == key2 {
		t.Fatalf("generateEtcdEncryptionKey should return a unique key each time, instead returned identical %s and %s", key1, key2)
	}
	for _, val := range []string{key1, key2} {
		_, err := base64.StdEncoding.DecodeString(val)
		if err != nil {
			t.Fatalf("generateEtcdEncryptionKey should return a base64 encoded key, instead returned %s", val)
		}
	}
}

func TestNetworkPolicyDefaults(t *testing.T) {
	mockCS := getMockBaseContainerService("1.8.10")
	properties := mockCS.Properties
	properties.OrchestratorProfile.OrchestratorType = "Kubernetes"
	properties.OrchestratorProfile.KubernetesConfig.NetworkPolicy = "calico"
	mockCS.setOrchestratorDefaults(true)
	if properties.OrchestratorProfile.KubernetesConfig.NetworkPlugin != "kubenet" {
		t.Fatalf("NetworkPlugin did not have the expected value, got %s, expected %s",
			properties.OrchestratorProfile.KubernetesConfig.NetworkPlugin, "kubenet")
	}

	mockCS = getMockBaseContainerService("1.8.10")
	properties = mockCS.Properties
	properties.OrchestratorProfile.OrchestratorType = "Kubernetes"
	properties.OrchestratorProfile.KubernetesConfig.NetworkPolicy = "cilium"
	mockCS.setOrchestratorDefaults(true)
	if properties.OrchestratorProfile.KubernetesConfig.NetworkPlugin != "cilium" {
		t.Fatalf("NetworkPlugin did not have the expected value, got %s, expected %s",
			properties.OrchestratorProfile.KubernetesConfig.NetworkPlugin, "cilium")
	}

	mockCS = getMockBaseContainerService("1.8.10")
	properties = mockCS.Properties
	properties.OrchestratorProfile.OrchestratorType = "Kubernetes"
	properties.OrchestratorProfile.KubernetesConfig.NetworkPolicy = "azure"
	mockCS.setOrchestratorDefaults(true)
	if properties.OrchestratorProfile.KubernetesConfig.NetworkPlugin != "azure" {
		t.Fatalf("NetworkPlugin did not have the expected value, got %s, expected %s",
			properties.OrchestratorProfile.KubernetesConfig.NetworkPlugin, "azure")
	}
	if properties.OrchestratorProfile.KubernetesConfig.NetworkPolicy != "" {
		t.Fatalf("NetworkPolicy did not have the expected value, got %s, expected %s",
			properties.OrchestratorProfile.KubernetesConfig.NetworkPolicy, "")
	}

	mockCS = getMockBaseContainerService("1.8.10")
	properties = mockCS.Properties
	properties.OrchestratorProfile.OrchestratorType = "Kubernetes"
	properties.OrchestratorProfile.KubernetesConfig.NetworkPolicy = "none"
	mockCS.setOrchestratorDefaults(true)
	if properties.OrchestratorProfile.KubernetesConfig.NetworkPlugin != "kubenet" {
		t.Fatalf("NetworkPlugin did not have the expected value, got %s, expected %s",
			properties.OrchestratorProfile.KubernetesConfig.NetworkPlugin, "kubenet")
	}
	if properties.OrchestratorProfile.KubernetesConfig.NetworkPolicy != "" {
		t.Fatalf("NetworkPolicy did not have the expected value, got %s, expected %s",
			properties.OrchestratorProfile.KubernetesConfig.NetworkPolicy, "")
	}
}

func TestStorageProfile(t *testing.T) {
	// Test ManagedDisks default configuration
	mockCS := getMockBaseContainerService("1.8.10")
	properties := mockCS.Properties
	properties.OrchestratorProfile.OrchestratorType = "Kubernetes"
	properties.MasterProfile.Count = 1
	properties.OrchestratorProfile.KubernetesConfig.PrivateCluster = &PrivateCluster{
		Enabled:        helpers.PointerToBool(true),
		JumpboxProfile: &PrivateJumpboxProfile{},
	}
	mockCS.SetPropertiesDefaults(false, false)
	if properties.MasterProfile.StorageProfile != ManagedDisks {
		t.Fatalf("MasterProfile.StorageProfile did not have the expected configuration, got %s, expected %s",
			properties.MasterProfile.StorageProfile, ManagedDisks)
	}
	if !properties.MasterProfile.IsManagedDisks() {
		t.Fatalf("MasterProfile.StorageProfile did not have the expected configuration, got %t, expected %t",
			false, true)
	}
	if properties.AgentPoolProfiles[0].StorageProfile != ManagedDisks {
		t.Fatalf("AgentPoolProfile.StorageProfile did not have the expected configuration, got %s, expected %s",
			properties.AgentPoolProfiles[0].StorageProfile, ManagedDisks)
	}
	if !properties.AgentPoolProfiles[0].IsManagedDisks() {
		t.Fatalf("AgentPoolProfile.IsManagedDisks() did not have the expected configuration, got %t, expected %t",
			false, true)
	}
	if properties.OrchestratorProfile.KubernetesConfig.PrivateCluster.JumpboxProfile.StorageProfile != ManagedDisks {
		t.Fatalf("MasterProfile.StorageProfile did not have the expected configuration, got %s, expected %s",
			properties.OrchestratorProfile.KubernetesConfig.PrivateCluster.JumpboxProfile.StorageProfile, ManagedDisks)
	}
	if !properties.AgentPoolProfiles[0].IsAvailabilitySets() {
		t.Fatalf("AgentPoolProfile[0].AvailabilityProfile did not have the expected configuration, got %s, expected %s",
			properties.AgentPoolProfiles[0].AvailabilityProfile, AvailabilitySet)
	}

	mockCS = getMockBaseContainerService("1.10.2")
	properties = mockCS.Properties
	properties.OrchestratorProfile.OrchestratorType = "Kubernetes"
	mockCS.SetPropertiesDefaults(false, false)
	if !properties.AgentPoolProfiles[0].IsVirtualMachineScaleSets() {
		t.Fatalf("AgentPoolProfile[0].AvailabilityProfile did not have the expected configuration, got %s, expected %s",
			properties.AgentPoolProfiles[0].AvailabilityProfile, VirtualMachineScaleSets)
	}

}

// TestMasterProfileDefaults covers tests for setMasterProfileDefaults
func TestMasterProfileDefaults(t *testing.T) {
	// this validates default masterProfile configuration
	mockCS := getMockBaseContainerService("1.10.3")
	properties := mockCS.Properties
	properties.OrchestratorProfile.OrchestratorType = "Kubernetes"
	properties.OrchestratorProfile.KubernetesConfig.ClusterSubnet = ""
	properties.OrchestratorProfile.KubernetesConfig.NetworkPlugin = "azure"
	properties.MasterProfile.AvailabilityProfile = ""
	properties.MasterProfile.Count = 3
	mockCS.Properties = properties
	mockCS.SetPropertiesDefaults(false, false)
	if properties.MasterProfile.IsVirtualMachineScaleSets() {
		t.Fatalf("Master VMAS, AzureCNI: MasterProfile AvailabilityProfile did not have the expected default configuration, got %s, expected %s",
			properties.MasterProfile.AvailabilityProfile, AvailabilitySet)
	}
	if properties.OrchestratorProfile.KubernetesConfig.ClusterSubnet != DefaultKubernetesSubnet {
		t.Fatalf("Master VMAS, AzureCNI: MasterProfile ClusterSubnet did not have the expected default configuration, got %s, expected %s",
			properties.OrchestratorProfile.KubernetesConfig.ClusterSubnet, DefaultKubernetesSubnet)
	}
	if properties.MasterProfile.Subnet != properties.OrchestratorProfile.KubernetesConfig.ClusterSubnet {
		t.Fatalf("Master VMAS, AzureCNI: MasterProfile Subnet did not have the expected default configuration, got %s, expected %s",
			properties.MasterProfile.Subnet, properties.OrchestratorProfile.KubernetesConfig.ClusterSubnet)
	}
	if properties.AgentPoolProfiles[0].Subnet != properties.MasterProfile.Subnet {
		t.Fatalf("Master VMAS, AzureCNI: AgentPoolProfiles Subnet did not have the expected default configuration, got %s, expected %s",
			properties.AgentPoolProfiles[0].Subnet, properties.MasterProfile.Subnet)
	}
	if properties.MasterProfile.FirstConsecutiveStaticIP != "10.255.255.5" {
		t.Fatalf("Master VMAS, AzureCNI: MasterProfile FirstConsecutiveStaticIP did not have the expected default configuration, got %s, expected %s",
			properties.MasterProfile.FirstConsecutiveStaticIP, "10.255.255.5")
	}

	// this validates default vmss masterProfile configuration
	mockCS = getMockBaseContainerService("1.10.3")
	properties = mockCS.Properties
	properties.OrchestratorProfile.OrchestratorType = "Kubernetes"
	properties.OrchestratorProfile.KubernetesConfig.NetworkPlugin = "azure"
	properties.OrchestratorProfile.KubernetesConfig.ClusterSubnet = ""
	properties.MasterProfile.AvailabilityProfile = VirtualMachineScaleSets
	mockCS.SetPropertiesDefaults(false, true)
	if !properties.MasterProfile.IsVirtualMachineScaleSets() {
		t.Fatalf("Master VMSS, AzureCNI: MasterProfile AvailabilityProfile did not have the expected default configuration, got %s, expected %s",
			properties.MasterProfile.AvailabilityProfile, VirtualMachineScaleSets)
	}
	if properties.OrchestratorProfile.KubernetesConfig.ClusterSubnet != DefaultKubernetesSubnet {
		t.Fatalf("Master VMSS, AzureCNI: MasterProfile ClusterSubnet did not have the expected default configuration, got %s, expected %s",
			properties.OrchestratorProfile.KubernetesConfig.ClusterSubnet, DefaultKubernetesSubnet)
	}
	if properties.MasterProfile.FirstConsecutiveStaticIP != DefaultFirstConsecutiveKubernetesStaticIPVMSS {
		t.Fatalf("Master VMSS, AzureCNI: MasterProfile FirstConsecutiveStaticIP did not have the expected default configuration, got %s, expected %s",
			properties.MasterProfile.FirstConsecutiveStaticIP, DefaultFirstConsecutiveKubernetesStaticIPVMSS)
	}
	if properties.MasterProfile.Subnet != DefaultKubernetesMasterSubnet {
		t.Fatalf("Master VMSS, AzureCNI: MasterProfile Subnet did not have the expected default configuration, got %s, expected %s",
			properties.MasterProfile.Subnet, DefaultKubernetesMasterSubnet)
	}
	if properties.MasterProfile.AgentSubnet != DefaultKubernetesAgentSubnetVMSS {
		t.Fatalf("Master VMSS, AzureCNI: MasterProfile AgentSubnet did not have the expected default configuration, got %s, expected %s",
			properties.MasterProfile.AgentSubnet, DefaultKubernetesAgentSubnetVMSS)
	}

	// this validates default masterProfile configuration and kubenet
	mockCS = getMockBaseContainerService("1.10.3")
	properties = mockCS.Properties
	properties.OrchestratorProfile.OrchestratorType = "Kubernetes"
	properties.OrchestratorProfile.KubernetesConfig.ClusterSubnet = ""
	properties.OrchestratorProfile.KubernetesConfig.NetworkPlugin = "kubenet"
	properties.MasterProfile.AvailabilityProfile = VirtualMachineScaleSets
	mockCS.SetPropertiesDefaults(false, true)
	if properties.OrchestratorProfile.KubernetesConfig.ClusterSubnet != DefaultKubernetesClusterSubnet {
		t.Fatalf("Master VMSS, kubenet: MasterProfile ClusterSubnet did not have the expected default configuration, got %s, expected %s",
			properties.OrchestratorProfile.KubernetesConfig.ClusterSubnet, DefaultKubernetesClusterSubnet)
	}
	if properties.MasterProfile.Subnet != DefaultKubernetesMasterSubnet {
		t.Fatalf("Master VMSS, kubenet: MasterProfile Subnet did not have the expected default configuration, got %s, expected %s",
			properties.MasterProfile.Subnet, DefaultKubernetesMasterSubnet)
	}
	if properties.MasterProfile.FirstConsecutiveStaticIP != DefaultFirstConsecutiveKubernetesStaticIPVMSS {
		t.Fatalf("Master VMSS, kubenet: MasterProfile FirstConsecutiveStaticIP did not have the expected default configuration, got %s, expected %s",
			properties.MasterProfile.FirstConsecutiveStaticIP, DefaultFirstConsecutiveKubernetesStaticIPVMSS)
	}
	if properties.MasterProfile.AgentSubnet != DefaultKubernetesAgentSubnetVMSS {
		t.Fatalf("Master VMSS, kubenet: MasterProfile AgentSubnet did not have the expected default configuration, got %s, expected %s",
			properties.MasterProfile.AgentSubnet, DefaultKubernetesAgentSubnetVMSS)
	}
	properties.MasterProfile.AvailabilityProfile = AvailabilitySet
	mockCS.SetPropertiesDefaults(false, true)
	if properties.MasterProfile.FirstConsecutiveStaticIP != DefaultFirstConsecutiveKubernetesStaticIP {
		t.Fatalf("Master VMAS, kubenet: MasterProfile FirstConsecutiveStaticIP did not have the expected default configuration, got %s, expected %s",
			properties.MasterProfile.FirstConsecutiveStaticIP, DefaultFirstConsecutiveKubernetesStaticIP)
	}

	// this validates default vmas masterProfile configuration, AzureCNI, and custom vnet
	mockCS = getMockBaseContainerService("1.10.3")
	properties = mockCS.Properties
	properties.OrchestratorProfile.OrchestratorType = "Kubernetes"
	properties.MasterProfile.VnetSubnetID = "/subscriptions/SUBSCRIPTION_ID/resourceGroups/RESOURCE_GROUP_NAME/providers/Microsoft.Network/virtualNetworks/ExampleCustomVNET/subnets/ExampleMasterSubnet"
	properties.MasterProfile.VnetCidr = "10.239.0.0/16"
	properties.MasterProfile.FirstConsecutiveStaticIP = "10.239.255.239"
	properties.OrchestratorProfile.KubernetesConfig.ClusterSubnet = ""
	properties.OrchestratorProfile.KubernetesConfig.NetworkPlugin = "azure"
	properties.MasterProfile.AvailabilityProfile = AvailabilitySet
	mockCS.SetPropertiesDefaults(false, true)
	if properties.MasterProfile.FirstConsecutiveStaticIP != "10.239.255.239" {
		t.Fatalf("Master VMAS, AzureCNI, customvnet: MasterProfile FirstConsecutiveStaticIP did not have the expected default configuration, got %s, expected %s",
			properties.MasterProfile.FirstConsecutiveStaticIP, "10.239.255.239")
	}

	// this validates default vmss masterProfile configuration, AzureCNI, and custom vnet
	mockCS = getMockBaseContainerService("1.10.3")
	properties = mockCS.Properties
	properties.OrchestratorProfile.OrchestratorType = "Kubernetes"
	properties.MasterProfile.VnetSubnetID = "/subscriptions/SUBSCRIPTION_ID/resourceGroups/RESOURCE_GROUP_NAME/providers/Microsoft.Network/virtualNetworks/ExampleCustomVNET/subnets/ExampleMasterSubnet"
	properties.MasterProfile.VnetCidr = "10.239.0.0/16"
	properties.OrchestratorProfile.KubernetesConfig.ClusterSubnet = ""
	properties.OrchestratorProfile.KubernetesConfig.NetworkPlugin = "azure"
	properties.MasterProfile.AvailabilityProfile = VirtualMachineScaleSets
	mockCS.SetPropertiesDefaults(false, true)
	if properties.MasterProfile.FirstConsecutiveStaticIP != "10.239.0.4" {
		t.Fatalf("Master VMSS, AzureCNI, customvnet: MasterProfile FirstConsecutiveStaticIP did not have the expected default configuration, got %s, expected %s",
			properties.MasterProfile.FirstConsecutiveStaticIP, "10.239.0.4")
	}

	// this validates default configurations for LoadBalancerSku and ExcludeMasterFromStandardLB
	mockCS = getMockBaseContainerService("1.11.6")
	properties = mockCS.Properties
	properties.OrchestratorProfile.OrchestratorType = "Kubernetes"
	properties.OrchestratorProfile.KubernetesConfig.LoadBalancerSku = "Standard"
	mockCS.SetPropertiesDefaults(false, false)
	if *properties.OrchestratorProfile.KubernetesConfig.ExcludeMasterFromStandardLB != DefaultExcludeMasterFromStandardLB {
		t.Fatalf("OrchestratorProfile.KubernetesConfig.ExcludeMasterFromStandardLB did not have the expected configuration, got %t, expected %t",
			*properties.OrchestratorProfile.KubernetesConfig.ExcludeMasterFromStandardLB, DefaultExcludeMasterFromStandardLB)
	}
}

func TestAgentPoolProfile(t *testing.T) {
	mockCS := getMockBaseContainerService("1.10")
	properties := mockCS.Properties
	properties.OrchestratorProfile.OrchestratorType = "Kubernetes"
	properties.MasterProfile.Count = 1
	mockCS.SetPropertiesDefaults(false, false)
	if properties.AgentPoolProfiles[0].ScaleSetPriority != "" {
		t.Fatalf("AgentPoolProfiles[0].ScaleSetPriority did not have the expected configuration, got %s, expected %s",
			properties.AgentPoolProfiles[0].ScaleSetPriority, "")
	}
	if properties.AgentPoolProfiles[0].ScaleSetEvictionPolicy != "" {
		t.Fatalf("AgentPoolProfiles[0].ScaleSetEvictionPolicy did not have the expected configuration, got %s, expected %s",
			properties.AgentPoolProfiles[0].ScaleSetEvictionPolicy, "")
	}
	properties.AgentPoolProfiles[0].ScaleSetPriority = ScaleSetPriorityLow
	mockCS.SetPropertiesDefaults(false, false)
	if properties.AgentPoolProfiles[0].ScaleSetEvictionPolicy != ScaleSetEvictionPolicyDelete {
		t.Fatalf("AgentPoolProfile[0].ScaleSetEvictionPolicy did not have the expected configuration, got %s, expected %s",
			properties.AgentPoolProfiles[0].ScaleSetEvictionPolicy, ScaleSetEvictionPolicyDelete)
	}
}

// TestSetComponentsNetworkDefaults covers tests for setMasterProfileDefaults and setAgentProfileDefaults
// TODO: Currently this test covers only Distro setting. Extend test cases to cover network configuration too.
func TestSetComponentsNetworkDefaults(t *testing.T) {

	var tests = []struct {
		name                string              // test case name
		orchestratorProfile OrchestratorProfile // orchestrator to be tested
		expectedDistro      Distro              // expected result default disto to be used
	}{
		{
			"default_kubernetes",
			OrchestratorProfile{
				OrchestratorType: Kubernetes,
			},
			AKS,
		},
		{
			"default_openshift",
			OrchestratorProfile{
				OrchestratorType: OpenShift,
			},
			"",
		},
		{
			"default_swarm",
			OrchestratorProfile{
				OrchestratorType: Swarm,
			},
			Ubuntu,
		},
		{
			"default_swarmmode",
			OrchestratorProfile{
				OrchestratorType: SwarmMode,
			},
			Ubuntu,
		},
		{
			"default_dcos",
			OrchestratorProfile{
				OrchestratorType: DCOS,
			},
			Ubuntu,
		},
	}

	for _, test := range tests {
		mockAPI := getMockAPIProperties("1.0.0")
		mockAPI.OrchestratorProfile = &test.orchestratorProfile
		mockAPI.setMasterProfileDefaults(false)
		mockAPI.setAgentProfileDefaults(false, false)
		if mockAPI.MasterProfile.Distro != test.expectedDistro {
			t.Fatalf("setMasterProfileDefaults() test case %v did not return right Distro configurations %v != %v", test.name, mockAPI.MasterProfile.Distro, test.expectedDistro)
		}
		for _, agent := range mockAPI.AgentPoolProfiles {
			if agent.Distro != test.expectedDistro {
				t.Fatalf("setAgentProfileDefaults() test case %v did not return right Distro configurations %v != %v", test.name, agent.Distro, test.expectedDistro)
			}
		}
	}
}

func TestIsAzureCNINetworkmonitorAddon(t *testing.T) {
	mockCS := getMockBaseContainerService("1.10.3")
	properties := mockCS.Properties
	properties.OrchestratorProfile.OrchestratorType = "Kubernetes"
	properties.MasterProfile.Count = 1
	properties.OrchestratorProfile.KubernetesConfig.Addons = []KubernetesAddon{
		{
			Name: AzureCNINetworkMonitoringAddonName,
			Containers: []KubernetesContainerSpec{
				{
					Name:           AzureCNINetworkMonitoringAddonName,
					CPURequests:    "50m",
					MemoryRequests: "150Mi",
					CPULimits:      "50m",
					MemoryLimits:   "150Mi",
				},
			},
			Enabled: helpers.PointerToBool(true),
		},
	}
	mockCS.setOrchestratorDefaults(true)

	i := getAddonsIndexByName(properties.OrchestratorProfile.KubernetesConfig.Addons, AzureCNINetworkMonitoringAddonName)
	if !helpers.IsTrueBoolPointer(properties.OrchestratorProfile.KubernetesConfig.Addons[i].Enabled) {
		t.Fatalf("Azure CNI networkmonitor addon should be present")
	}

	mockCS = getMockBaseContainerService("1.10.3")
	properties = mockCS.Properties
	properties.OrchestratorProfile.OrchestratorType = "Kubernetes"
	properties.MasterProfile.Count = 1
	properties.OrchestratorProfile.KubernetesConfig.NetworkPlugin = "azure"
	mockCS.setOrchestratorDefaults(true)

	i = getAddonsIndexByName(properties.OrchestratorProfile.KubernetesConfig.Addons, AzureCNINetworkMonitoringAddonName)
	if !helpers.IsTrueBoolPointer(properties.OrchestratorProfile.KubernetesConfig.Addons[i].Enabled) {
		t.Fatalf("Azure CNI networkmonitor addon should be present by default if Azure CNI is set")
	}
}

// TestSetVMSSDefaultsAndZones covers tests for setVMSSDefaultsForAgents and masters
func TestSetVMSSDefaultsAndZones(t *testing.T) {
	// masters with vmss and no zones
	mockCS := getMockBaseContainerService("1.12.0")
	properties := mockCS.Properties
	properties.OrchestratorProfile.OrchestratorType = "Kubernetes"
	properties.MasterProfile.AvailabilityProfile = VirtualMachineScaleSets
	mockCS.SetPropertiesDefaults(false, false)
	if properties.MasterProfile.HasAvailabilityZones() {
		t.Fatalf("MasterProfile.HasAvailabilityZones did not have the expected return, got %t, expected %t",
			properties.MasterProfile.HasAvailabilityZones(), false)
	}
	if properties.OrchestratorProfile.KubernetesConfig.LoadBalancerSku != DefaultLoadBalancerSku {
		t.Fatalf("OrchestratorProfile.KubernetesConfig.LoadBalancerSku did not have the expected configuration, got %s, expected %s",
			properties.OrchestratorProfile.KubernetesConfig.LoadBalancerSku, DefaultLoadBalancerSku)
	}
	// masters with vmss and zones
	mockCS = getMockBaseContainerService("1.12.0")
	properties = mockCS.Properties
	properties.OrchestratorProfile.OrchestratorType = "Kubernetes"
	properties.MasterProfile.AvailabilityProfile = VirtualMachineScaleSets
	properties.MasterProfile.AvailabilityZones = []string{"1", "2"}
	mockCS.SetPropertiesDefaults(false, false)
	if *properties.MasterProfile.SinglePlacementGroup != DefaultSinglePlacementGroup {
		t.Fatalf("MasterProfile.SinglePlacementGroup default did not have the expected configuration, got %t, expected %t",
			*properties.MasterProfile.SinglePlacementGroup, DefaultSinglePlacementGroup)
	}
	if !properties.MasterProfile.HasAvailabilityZones() {
		t.Fatalf("MasterProfile.HasAvailabilityZones did not have the expected return, got %t, expected %t",
			properties.MasterProfile.HasAvailabilityZones(), true)
	}
	if properties.OrchestratorProfile.KubernetesConfig.LoadBalancerSku != "Standard" {
		t.Fatalf("OrchestratorProfile.KubernetesConfig.LoadBalancerSku did not have the expected configuration, got %s, expected %s",
			properties.OrchestratorProfile.KubernetesConfig.LoadBalancerSku, "Standard")
	}
	if *properties.OrchestratorProfile.KubernetesConfig.ExcludeMasterFromStandardLB != DefaultExcludeMasterFromStandardLB {
		t.Fatalf("OrchestratorProfile.KubernetesConfig.ExcludeMasterFromStandardLB did not have the expected configuration, got %t, expected %t",
			*properties.OrchestratorProfile.KubernetesConfig.ExcludeMasterFromStandardLB, DefaultExcludeMasterFromStandardLB)
	}
	// agents with vmss and no zones
	mockCS = getMockBaseContainerService("1.12.0")
	properties = mockCS.Properties
	properties.OrchestratorProfile.OrchestratorType = "Kubernetes"
	properties.AgentPoolProfiles[0].Count = 4
	mockCS.SetPropertiesDefaults(false, false)
	if properties.AgentPoolProfiles[0].HasAvailabilityZones() {
		t.Fatalf("AgentPoolProfiles[0].HasAvailabilityZones did not have the expected return, got %t, expected %t",
			properties.AgentPoolProfiles[0].HasAvailabilityZones(), false)
	}
	if properties.OrchestratorProfile.KubernetesConfig.LoadBalancerSku != DefaultLoadBalancerSku {
		t.Fatalf("OrchestratorProfile.KubernetesConfig.LoadBalancerSku did not have the expected configuration, got %s, expected %s",
			properties.OrchestratorProfile.KubernetesConfig.LoadBalancerSku, DefaultLoadBalancerSku)
	}
	// agents with vmss and zones
	mockCS = getMockBaseContainerService("1.12.0")
	properties = mockCS.Properties
	properties.OrchestratorProfile.OrchestratorType = "Kubernetes"
	properties.AgentPoolProfiles[0].Count = 4
	properties.AgentPoolProfiles[0].AvailabilityZones = []string{"1", "2"}
	mockCS.SetPropertiesDefaults(false, false)
	if !properties.AgentPoolProfiles[0].IsVirtualMachineScaleSets() {
		t.Fatalf("AgentPoolProfile[0].AvailabilityProfile did not have the expected configuration, got %s, expected %s",
			properties.AgentPoolProfiles[0].AvailabilityProfile, VirtualMachineScaleSets)
	}
	if !properties.AgentPoolProfiles[0].HasAvailabilityZones() {
		t.Fatalf("AgentPoolProfiles[0].HasAvailabilityZones did not have the expected return, got %t, expected %t",
			properties.AgentPoolProfiles[0].HasAvailabilityZones(), true)
	}
	if *properties.AgentPoolProfiles[0].SinglePlacementGroup != DefaultSinglePlacementGroup {
		t.Fatalf("AgentPoolProfile[0].SinglePlacementGroup default did not have the expected configuration, got %t, expected %t",
			*properties.AgentPoolProfiles[0].SinglePlacementGroup, DefaultSinglePlacementGroup)
	}
	if properties.OrchestratorProfile.KubernetesConfig.LoadBalancerSku != "Standard" {
		t.Fatalf("OrchestratorProfile.KubernetesConfig.LoadBalancerSku did not have the expected configuration, got %s, expected %s",
			properties.OrchestratorProfile.KubernetesConfig.LoadBalancerSku, "Standard")
	}
	if *properties.OrchestratorProfile.KubernetesConfig.ExcludeMasterFromStandardLB != DefaultExcludeMasterFromStandardLB {
		t.Fatalf("OrchestratorProfile.KubernetesConfig.ExcludeMasterFromStandardLB did not have the expected configuration, got %t, expected %t",
			*properties.OrchestratorProfile.KubernetesConfig.ExcludeMasterFromStandardLB, DefaultExcludeMasterFromStandardLB)
	}

	properties.AgentPoolProfiles[0].Count = 110
	mockCS.SetPropertiesDefaults(false, false)
	if helpers.IsTrueBoolPointer(properties.AgentPoolProfiles[0].SinglePlacementGroup) {
		t.Fatalf("AgentPoolProfile[0].SinglePlacementGroup did not have the expected configuration, got %t, expected %t",
			*properties.AgentPoolProfiles[0].SinglePlacementGroup, false)
	}

	if !*properties.AgentPoolProfiles[0].SinglePlacementGroup && properties.AgentPoolProfiles[0].StorageProfile != ManagedDisks {
		t.Fatalf("AgentPoolProfile[0].StorageProfile did not have the expected configuration, got %s, expected %s",
			properties.AgentPoolProfiles[0].StorageProfile, ManagedDisks)
	}

}

func TestAzureCNIVersionString(t *testing.T) {
	mockCS := getMockBaseContainerService("1.10.3")
	properties := mockCS.Properties
	properties.OrchestratorProfile.OrchestratorType = "Kubernetes"
	properties.MasterProfile.Count = 1
	properties.OrchestratorProfile.KubernetesConfig.NetworkPlugin = "azure"
	mockCS.setOrchestratorDefaults(true)

	if properties.OrchestratorProfile.KubernetesConfig.AzureCNIVersion != AzureCniPluginVerLinux {
		t.Fatalf("Azure CNI Version string not the expected value, got %s, expected %s", properties.OrchestratorProfile.KubernetesConfig.AzureCNIVersion, AzureCniPluginVerLinux)
	}

	mockCS = getMockBaseContainerService("1.10.3")
	properties = mockCS.Properties
	properties.OrchestratorProfile.OrchestratorType = "Kubernetes"
	properties.MasterProfile.Count = 1
	properties.AgentPoolProfiles[0].OSType = "Windows"
	properties.OrchestratorProfile.KubernetesConfig.NetworkPlugin = "azure"
	mockCS.setOrchestratorDefaults(true)

	if properties.OrchestratorProfile.KubernetesConfig.AzureCNIVersion != AzureCniPluginVerWindows {
		t.Fatalf("Azure CNI Version string not the expected value, got %s, expected %s", properties.OrchestratorProfile.KubernetesConfig.AzureCNIVersion, AzureCniPluginVerWindows)
	}

	mockCS = getMockBaseContainerService("1.10.3")
	properties = mockCS.Properties
	properties.OrchestratorProfile.OrchestratorType = "Kubernetes"
	properties.MasterProfile.Count = 1
	properties.OrchestratorProfile.KubernetesConfig.NetworkPlugin = "kubenet"
	mockCS.setOrchestratorDefaults(true)

	if properties.OrchestratorProfile.KubernetesConfig.AzureCNIVersion != "" {
		t.Fatalf("Azure CNI Version string not the expected value, got %s, expected %s", properties.OrchestratorProfile.KubernetesConfig.AzureCNIVersion, "")
	}
}

func TestDefaultDisableRbac(t *testing.T) {
	mockCS := getMockBaseContainerService("1.10.3")
	properties := mockCS.Properties
	properties.OrchestratorProfile.OrchestratorType = "Kubernetes"
	properties.OrchestratorProfile.KubernetesConfig.EnableRbac = helpers.PointerToBool(false)
	mockCS.setOrchestratorDefaults(true)

	if properties.OrchestratorProfile.KubernetesConfig.EnableAggregatedAPIs {
		t.Fatalf("got unexpected EnableAggregatedAPIs config value for EnableRbac=false: %t",
			properties.OrchestratorProfile.KubernetesConfig.EnableAggregatedAPIs)
	}
}

func TestDefaultCloudProvider(t *testing.T) {
	mockCS := getMockBaseContainerService("1.10.3")
	properties := mockCS.Properties
	properties.OrchestratorProfile.OrchestratorType = "Kubernetes"
	mockCS.setOrchestratorDefaults(true)

	if !helpers.IsTrueBoolPointer(properties.OrchestratorProfile.KubernetesConfig.CloudProviderBackoff) {
		t.Fatalf("got unexpected CloudProviderBackoff expected true, got %t",
			helpers.IsTrueBoolPointer(properties.OrchestratorProfile.KubernetesConfig.CloudProviderBackoff))
	}

	if !helpers.IsTrueBoolPointer(properties.OrchestratorProfile.KubernetesConfig.CloudProviderRateLimit) {
		t.Fatalf("got unexpected CloudProviderBackoff expected true, got %t",
			helpers.IsTrueBoolPointer(properties.OrchestratorProfile.KubernetesConfig.CloudProviderBackoff))
	}

	mockCS = getMockBaseContainerService("1.10.3")
	properties = mockCS.Properties
	properties.OrchestratorProfile.OrchestratorType = "Kubernetes"
	properties.OrchestratorProfile.KubernetesConfig.CloudProviderBackoff = helpers.PointerToBool(false)
	properties.OrchestratorProfile.KubernetesConfig.CloudProviderRateLimit = helpers.PointerToBool(false)
	mockCS.setOrchestratorDefaults(true)

	if !helpers.IsFalseBoolPointer(properties.OrchestratorProfile.KubernetesConfig.CloudProviderBackoff) {
		t.Fatalf("got unexpected CloudProviderBackoff expected true, got %t",
			helpers.IsTrueBoolPointer(properties.OrchestratorProfile.KubernetesConfig.CloudProviderBackoff))
	}

	if !helpers.IsFalseBoolPointer(properties.OrchestratorProfile.KubernetesConfig.CloudProviderRateLimit) {
		t.Fatalf("got unexpected CloudProviderBackoff expected true, got %t",
			helpers.IsTrueBoolPointer(properties.OrchestratorProfile.KubernetesConfig.CloudProviderBackoff))
	}
}

func TestSetOpenShiftCertDefaults(t *testing.T) {
	cs := &ContainerService{
		Properties: &Properties{
			AzProfile: &AzProfile{
				TenantID:       "sampleTenantID",
				SubscriptionID: "foobarsubscription",
				ResourceGroup:  "sampleRG",
				Location:       "westus2",
			},
			ServicePrincipalProfile: &ServicePrincipalProfile{
				ClientID: "barClientID",
				Secret:   "bazSecret",
			},
			MasterProfile: &MasterProfile{
				Count:     1,
				DNSPrefix: "myprefix1",
				VMSize:    "Standard_DS2_v2",
			},
			OrchestratorProfile: &OrchestratorProfile{
				OrchestratorType:    OpenShift,
				OrchestratorVersion: "3.9.0",
				OpenShiftConfig:     &OpenShiftConfig{},
			},
		},
	}

	cs.Properties.setMasterProfileDefaults(false)

	result, err := cs.Properties.setDefaultCerts()
	if !result {
		t.Error("expected setOpenShiftDefaultCerts to return true")
	}

	if err != nil {
		t.Errorf("unexpected error thrown while executing setOpenShiftDefaultCerts %s", err.Error())
	}

	cs = &ContainerService{
		Properties: &Properties{
			AzProfile: &AzProfile{
				TenantID:       "sampleTenantID",
				SubscriptionID: "foobarsubscription",
				ResourceGroup:  "sampleRG",
				Location:       "westus2",
			},
			ServicePrincipalProfile: &ServicePrincipalProfile{
				ClientID: "barClientID",
				Secret:   "bazSecret",
			},
			MasterProfile: &MasterProfile{
				Count:               1,
				DNSPrefix:           "myprefix1",
				VMSize:              "Standard_DS2_v2",
				AvailabilityProfile: VirtualMachineScaleSets,
			},
			OrchestratorProfile: &OrchestratorProfile{
				OrchestratorType:    OpenShift,
				OrchestratorVersion: "3.7.0",
				OpenShiftConfig:     &OpenShiftConfig{},
			},
		},
	}

	cs.Properties.setMasterProfileDefaults(false)
	result, err = cs.Properties.setDefaultCerts()

	if !result {
		t.Error("expected setOpenShiftDefaultCerts to return true")
	}

	if err != nil {
		t.Errorf("unexpected error thrown while executing setOpenShiftDefaultCerts %s", err.Error())
	}

}

func getMockBaseContainerService(orchestratorVersion string) ContainerService {
	mockAPIProperties := getMockAPIProperties(orchestratorVersion)
	return ContainerService{
		Properties: &mockAPIProperties,
	}
}

func getMockAPIProperties(orchestratorVersion string) Properties {
	return Properties{
		ProvisioningState: "",
		OrchestratorProfile: &OrchestratorProfile{
			OrchestratorVersion: orchestratorVersion,
			KubernetesConfig:    &KubernetesConfig{},
		},
		MasterProfile: &MasterProfile{},
		AgentPoolProfiles: []*AgentPoolProfile{
			{},
		}}
}

func getKubernetesConfigWithFeatureGates(featureGates string) *KubernetesConfig {
	return &KubernetesConfig{
		KubeletConfig: map[string]string{"--feature-gates": featureGates},
	}
}

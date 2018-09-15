package acsengine

import (
	"encoding/base64"
	"reflect"
	"testing"

	"github.com/Azure/acs-engine/pkg/api"
	"github.com/Azure/acs-engine/pkg/helpers"
)

func TestCertsAlreadyPresent(t *testing.T) {
	var cert *api.CertificateProfile

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
	cert = &api.CertificateProfile{}
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
	cert = &api.CertificateProfile{
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

	cert = &api.CertificateProfile{
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
	cert = &api.CertificateProfile{
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
	config := &api.KubernetesConfig{}
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

	config = &api.KubernetesConfig{
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
	config = &api.KubernetesConfig{
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
	addons := []api.KubernetesAddon{
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

func TestGetAddonContainersIndexByName(t *testing.T) {
	addonName := "testaddon"
	containers := getMockAddon(addonName).Containers
	i := getAddonContainersIndexByName(containers, addonName)
	if i != 0 {
		t.Fatalf("getAddonContainersIndexByName() did not return the expected index value 0, instead returned: %d", i)
	}
	i = getAddonContainersIndexByName(containers, "nonExistentContainerName")
	if i != -1 {
		t.Fatalf("getAddonContainersIndexByName() did not return the expected index value 0, instead returned: %d", i)
	}
}

func TestAssignDefaultAddonVals(t *testing.T) {
	addonName := "testaddon"
	customCPURequests := "60m"
	customMemoryRequests := "160Mi"
	customCPULimits := "40m"
	customMemoryLimits := "140Mi"
	// Verify that an addon with all custom values provided remains unmodified during default value assignment
	customAddon := api.KubernetesAddon{
		Name:    addonName,
		Enabled: helpers.PointerToBool(true),
		Containers: []api.KubernetesContainerSpec{
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
	customAddon = api.KubernetesAddon{
		Name:    addonName,
		Enabled: helpers.PointerToBool(true),
		Containers: []api.KubernetesContainerSpec{
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
	customAddon = api.KubernetesAddon{
		Name:    addonName,
		Enabled: helpers.PointerToBool(true),
		Containers: []api.KubernetesContainerSpec{
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

}

func TestKubeletFeatureGatesEnsureFeatureGatesOnAgentsFor1_6_0(t *testing.T) {
	mockCS := getMockBaseContainerService("1.6.0")
	properties := mockCS.Properties

	// No KubernetesConfig.KubeletConfig set for MasterProfile or AgentProfile
	// so they will inherit the top-level config
	properties.OrchestratorProfile.KubernetesConfig = getKubernetesConfigWithFeatureGates("TopLevel=true")

	setKubeletConfig(&mockCS)

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
	properties.MasterProfile = &api.MasterProfile{KubernetesConfig: getKubernetesConfigWithFeatureGates("MasterLevel=true")}
	properties.AgentPoolProfiles[0].KubernetesConfig = getKubernetesConfigWithFeatureGates("AgentLevel=true")

	setKubeletConfig(&mockCS)

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
	setOrchestratorDefaults(&mockCS, true)
	if properties.OrchestratorProfile.KubernetesConfig.EtcdDiskSizeGB != DefaultEtcdDiskSize {
		t.Fatalf("EtcdDiskSizeGB did not have the expected size, got %s, expected %s",
			properties.OrchestratorProfile.KubernetesConfig.EtcdDiskSizeGB, DefaultEtcdDiskSize)
	}

	mockCS = getMockBaseContainerService("1.8.10")
	properties = mockCS.Properties
	properties.OrchestratorProfile.OrchestratorType = "Kubernetes"
	properties.MasterProfile.Count = 5
	setOrchestratorDefaults(&mockCS, true)
	if properties.OrchestratorProfile.KubernetesConfig.EtcdDiskSizeGB != DefaultEtcdDiskSizeGT3Nodes {
		t.Fatalf("EtcdDiskSizeGB did not have the expected size, got %s, expected %s",
			properties.OrchestratorProfile.KubernetesConfig.EtcdDiskSizeGB, DefaultEtcdDiskSizeGT3Nodes)
	}

	mockCS = getMockBaseContainerService("1.8.10")
	properties = mockCS.Properties
	properties.OrchestratorProfile.OrchestratorType = "Kubernetes"
	properties.MasterProfile.Count = 5
	properties.AgentPoolProfiles[0].Count = 6
	setOrchestratorDefaults(&mockCS, true)
	if properties.OrchestratorProfile.KubernetesConfig.EtcdDiskSizeGB != DefaultEtcdDiskSizeGT10Nodes {
		t.Fatalf("EtcdDiskSizeGB did not have the expected size, got %s, expected %s",
			properties.OrchestratorProfile.KubernetesConfig.EtcdDiskSizeGB, DefaultEtcdDiskSizeGT10Nodes)
	}

	mockCS = getMockBaseContainerService("1.8.10")
	properties = mockCS.Properties
	properties.OrchestratorProfile.OrchestratorType = "Kubernetes"
	properties.MasterProfile.Count = 5
	properties.AgentPoolProfiles[0].Count = 16
	setOrchestratorDefaults(&mockCS, true)
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
	setOrchestratorDefaults(&mockCS, true)
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
	setOrchestratorDefaults(&mockCS, true)
	if properties.OrchestratorProfile.KubernetesConfig.NetworkPlugin != "kubenet" {
		t.Fatalf("NetworkPlugin did not have the expected value, got %s, expected %s",
			properties.OrchestratorProfile.KubernetesConfig.NetworkPlugin, "kubenet")
	}

	mockCS = getMockBaseContainerService("1.8.10")
	properties = mockCS.Properties
	properties.OrchestratorProfile.OrchestratorType = "Kubernetes"
	properties.OrchestratorProfile.KubernetesConfig.NetworkPolicy = "cilium"
	setOrchestratorDefaults(&mockCS, true)
	if properties.OrchestratorProfile.KubernetesConfig.NetworkPlugin != "cilium" {
		t.Fatalf("NetworkPlugin did not have the expected value, got %s, expected %s",
			properties.OrchestratorProfile.KubernetesConfig.NetworkPlugin, "cilium")
	}

	mockCS = getMockBaseContainerService("1.8.10")
	properties = mockCS.Properties
	properties.OrchestratorProfile.OrchestratorType = "Kubernetes"
	properties.OrchestratorProfile.KubernetesConfig.NetworkPolicy = "azure"
	setOrchestratorDefaults(&mockCS, true)
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
	setOrchestratorDefaults(&mockCS, true)
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
	properties.OrchestratorProfile.KubernetesConfig.PrivateCluster = &api.PrivateCluster{
		Enabled:        helpers.PointerToBool(true),
		JumpboxProfile: &api.PrivateJumpboxProfile{},
	}
	setPropertiesDefaults(&mockCS, false, false)
	if properties.MasterProfile.StorageProfile != api.ManagedDisks {
		t.Fatalf("MasterProfile.StorageProfile did not have the expected configuration, got %s, expected %s",
			properties.MasterProfile.StorageProfile, api.ManagedDisks)
	}
	if !properties.MasterProfile.IsManagedDisks() {
		t.Fatalf("MasterProfile.StorageProfile did not have the expected configuration, got %t, expected %t",
			false, true)
	}
	if properties.AgentPoolProfiles[0].StorageProfile != api.ManagedDisks {
		t.Fatalf("AgentPoolProfile.StorageProfile did not have the expected configuration, got %s, expected %s",
			properties.AgentPoolProfiles[0].StorageProfile, api.ManagedDisks)
	}
	if !properties.AgentPoolProfiles[0].IsManagedDisks() {
		t.Fatalf("AgentPoolProfile.IsManagedDisks() did not have the expected configuration, got %t, expected %t",
			false, true)
	}
	if properties.OrchestratorProfile.KubernetesConfig.PrivateCluster.JumpboxProfile.StorageProfile != api.ManagedDisks {
		t.Fatalf("MasterProfile.StorageProfile did not have the expected configuration, got %s, expected %s",
			properties.OrchestratorProfile.KubernetesConfig.PrivateCluster.JumpboxProfile.StorageProfile, api.ManagedDisks)
	}
	if !properties.AgentPoolProfiles[0].IsAvailabilitySets() {
		t.Fatalf("AgentPoolProfile[0].AvailabilityProfile did not have the expected configuration, got %s, expected %s",
			properties.AgentPoolProfiles[0].AvailabilityProfile, api.AvailabilitySet)
	}

	mockCS = getMockBaseContainerService("1.10.2")
	properties = mockCS.Properties
	properties.OrchestratorProfile.OrchestratorType = "Kubernetes"
	setPropertiesDefaults(&mockCS, false, false)
	if !properties.AgentPoolProfiles[0].IsVirtualMachineScaleSets() {
		t.Fatalf("AgentPoolProfile[0].AvailabilityProfile did not have the expected configuration, got %s, expected %s",
			properties.AgentPoolProfiles[0].AvailabilityProfile, api.VirtualMachineScaleSets)
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
	setPropertiesDefaults(&mockCS, false, false)
	if properties.MasterProfile.IsVirtualMachineScaleSets() {
		t.Fatalf("Master VMAS, AzureCNI: MasterProfile AvailabilityProfile did not have the expected default configuration, got %s, expected %s",
			properties.MasterProfile.AvailabilityProfile, api.AvailabilitySet)
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
	properties.MasterProfile.AvailabilityProfile = api.VirtualMachineScaleSets
	setPropertiesDefaults(&mockCS, false, true)
	if !properties.MasterProfile.IsVirtualMachineScaleSets() {
		t.Fatalf("Master VMSS, AzureCNI: MasterProfile AvailabilityProfile did not have the expected default configuration, got %s, expected %s",
			properties.MasterProfile.AvailabilityProfile, api.VirtualMachineScaleSets)
	}
	if properties.OrchestratorProfile.KubernetesConfig.ClusterSubnet != DefaultKubernetesSubnet {
		t.Fatalf("Master VMSS, AzureCNI: MasterProfile ClusterSubnet did not have the expected default configuration, got %s, expected %s",
			properties.OrchestratorProfile.KubernetesConfig.ClusterSubnet, DefaultKubernetesSubnet)
	}
	if properties.MasterProfile.FirstConsecutiveStaticIP != api.DefaultFirstConsecutiveKubernetesStaticIPVMSS {
		t.Fatalf("Master VMSS, AzureCNI: MasterProfile FirstConsecutiveStaticIP did not have the expected default configuration, got %s, expected %s",
			properties.MasterProfile.FirstConsecutiveStaticIP, api.DefaultFirstConsecutiveKubernetesStaticIPVMSS)
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
	properties.MasterProfile.AvailabilityProfile = api.VirtualMachineScaleSets
	setPropertiesDefaults(&mockCS, false, true)
	if properties.OrchestratorProfile.KubernetesConfig.ClusterSubnet != DefaultKubernetesClusterSubnet {
		t.Fatalf("Master VMSS, kubenet: MasterProfile ClusterSubnet did not have the expected default configuration, got %s, expected %s",
			properties.OrchestratorProfile.KubernetesConfig.ClusterSubnet, DefaultKubernetesClusterSubnet)
	}
	if properties.MasterProfile.Subnet != DefaultKubernetesMasterSubnet {
		t.Fatalf("Master VMSS, kubenet: MasterProfile Subnet did not have the expected default configuration, got %s, expected %s",
			properties.MasterProfile.Subnet, DefaultKubernetesMasterSubnet)
	}
	if properties.MasterProfile.FirstConsecutiveStaticIP != api.DefaultFirstConsecutiveKubernetesStaticIPVMSS {
		t.Fatalf("Master VMSS, kubenet: MasterProfile FirstConsecutiveStaticIP did not have the expected default configuration, got %s, expected %s",
			properties.MasterProfile.FirstConsecutiveStaticIP, api.DefaultFirstConsecutiveKubernetesStaticIPVMSS)
	}
	if properties.MasterProfile.AgentSubnet != DefaultKubernetesAgentSubnetVMSS {
		t.Fatalf("Master VMSS, kubenet: MasterProfile AgentSubnet did not have the expected default configuration, got %s, expected %s",
			properties.MasterProfile.AgentSubnet, DefaultKubernetesAgentSubnetVMSS)
	}
	properties.MasterProfile.AvailabilityProfile = api.AvailabilitySet
	setPropertiesDefaults(&mockCS, false, true)
	if properties.MasterProfile.FirstConsecutiveStaticIP != api.DefaultFirstConsecutiveKubernetesStaticIP {
		t.Fatalf("Master VMAS, kubenet: MasterProfile FirstConsecutiveStaticIP did not have the expected default configuration, got %s, expected %s",
			properties.MasterProfile.FirstConsecutiveStaticIP, api.DefaultFirstConsecutiveKubernetesStaticIP)
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
	properties.MasterProfile.AvailabilityProfile = api.AvailabilitySet
	setPropertiesDefaults(&mockCS, false, true)
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
	properties.MasterProfile.AvailabilityProfile = api.VirtualMachineScaleSets
	setPropertiesDefaults(&mockCS, false, true)
	if properties.MasterProfile.FirstConsecutiveStaticIP != "10.239.0.4" {
		t.Fatalf("Master VMSS, AzureCNI, customvnet: MasterProfile FirstConsecutiveStaticIP did not have the expected default configuration, got %s, expected %s",
			properties.MasterProfile.FirstConsecutiveStaticIP, "10.239.0.4")
	}

}

func TestAgentPoolProfile(t *testing.T) {
	mockCS := getMockBaseContainerService("1.10")
	properties := mockCS.Properties
	properties.OrchestratorProfile.OrchestratorType = "Kubernetes"
	properties.MasterProfile.Count = 1
	setPropertiesDefaults(&mockCS, false, false)
	if properties.AgentPoolProfiles[0].ScaleSetPriority != "" {
		t.Fatalf("AgentPoolProfiles[0].ScaleSetPriority did not have the expected configuration, got %s, expected %s",
			properties.AgentPoolProfiles[0].ScaleSetPriority, "")
	}
	if properties.AgentPoolProfiles[0].ScaleSetEvictionPolicy != "" {
		t.Fatalf("AgentPoolProfiles[0].ScaleSetEvictionPolicy did not have the expected configuration, got %s, expected %s",
			properties.AgentPoolProfiles[0].ScaleSetEvictionPolicy, "")
	}
	properties.AgentPoolProfiles[0].ScaleSetPriority = api.ScaleSetPriorityLow
	setPropertiesDefaults(&mockCS, false, false)
	if properties.AgentPoolProfiles[0].ScaleSetEvictionPolicy != api.ScaleSetEvictionPolicyDelete {
		t.Fatalf("AgentPoolProfile[0].ScaleSetEvictionPolicy did not have the expected configuration, got %s, expected %s",
			properties.AgentPoolProfiles[0].ScaleSetEvictionPolicy, api.ScaleSetEvictionPolicyDelete)
	}
}

// TestSetComponentsNetworkDefaults covers tests for setMasterProfileDefaults and setAgentProfileDefaults
// TODO: Currently this test covers only api.Distro setting. Extend test cases to cover network configuration too.
func TestSetComponentsNetworkDefaults(t *testing.T) {

	var tests = []struct {
		name                string                  // test case name
		orchestratorProfile api.OrchestratorProfile // orchestrator to be tested
		expectedDistro      api.Distro              // expected result default disto to be used
	}{
		{
			"ubuntu_kubernetes",
			api.OrchestratorProfile{
				OrchestratorType: api.Kubernetes,
			},
			api.AKS,
		},
		{
			"rhel_openshift",
			api.OrchestratorProfile{
				OrchestratorType: api.OpenShift,
			},
			"",
		},
	}

	for _, test := range tests {
		mockAPI := getMockAPIProperties("1.0.0")
		mockAPI.OrchestratorProfile = &test.orchestratorProfile
		setMasterProfileDefaults(&mockAPI, false)
		setAgentProfileDefaults(&mockAPI, false, false)
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
	properties.OrchestratorProfile.KubernetesConfig.Addons = []api.KubernetesAddon{
		getMockAddon(AzureCNINetworkMonitoringAddonName),
	}
	setOrchestratorDefaults(&mockCS, true)

	i := getAddonsIndexByName(properties.OrchestratorProfile.KubernetesConfig.Addons, AzureCNINetworkMonitoringAddonName)
	if !helpers.IsTrueBoolPointer(properties.OrchestratorProfile.KubernetesConfig.Addons[i].Enabled) {
		t.Fatalf("Azure CNI networkmonitor addon should be present")
	}

	mockCS = getMockBaseContainerService("1.10.3")
	properties = mockCS.Properties
	properties.OrchestratorProfile.OrchestratorType = "Kubernetes"
	properties.MasterProfile.Count = 1
	properties.OrchestratorProfile.KubernetesConfig.NetworkPlugin = "azure"
	setOrchestratorDefaults(&mockCS, true)

	i = getAddonsIndexByName(properties.OrchestratorProfile.KubernetesConfig.Addons, AzureCNINetworkMonitoringAddonName)
	if !helpers.IsTrueBoolPointer(properties.OrchestratorProfile.KubernetesConfig.Addons[i].Enabled) {
		t.Fatalf("Azure CNI networkmonitor addon should be present by default if Azure CNI is set")
	}
}

// TestSetVMSSDefaults covers tests for setVMSSDefaults
func TestSetVMSSDefaults(t *testing.T) {
	mockCS := getMockBaseContainerService("1.10.3")
	properties := mockCS.Properties
	properties.OrchestratorProfile.OrchestratorType = "Kubernetes"
	properties.AgentPoolProfiles[0].Count = 4
	setPropertiesDefaults(&mockCS, false, false)
	if !properties.AgentPoolProfiles[0].IsVirtualMachineScaleSets() {
		t.Fatalf("AgentPoolProfile[0].AvailabilityProfile did not have the expected configuration, got %s, expected %s",
			properties.AgentPoolProfiles[0].AvailabilityProfile, api.VirtualMachineScaleSets)
	}

	if *properties.AgentPoolProfiles[0].SinglePlacementGroup != api.DefaultSinglePlacementGroup {
		t.Fatalf("AgentPoolProfile[0].SinglePlacementGroup default did not have the expected configuration, got %t, expected %t",
			*properties.AgentPoolProfiles[0].SinglePlacementGroup, api.DefaultSinglePlacementGroup)
	}

	if properties.AgentPoolProfiles[0].HasAvailabilityZones() {
		if properties.OrchestratorProfile.KubernetesConfig.LoadBalancerSku != "Standard" {
			t.Fatalf("OrchestratorProfile.KubernetesConfig.LoadBalancerSku did not have the expected configuration, got %s, expected %s",
				properties.OrchestratorProfile.KubernetesConfig.LoadBalancerSku, "Standard")
		}
		if properties.OrchestratorProfile.KubernetesConfig.ExcludeMasterFromStandardLB != helpers.PointerToBool(api.DefaultExcludeMasterFromStandardLB) {
			t.Fatalf("OrchestratorProfile.KubernetesConfig.ExcludeMasterFromStandardLB did not have the expected configuration, got %t, expected %t",
				*properties.OrchestratorProfile.KubernetesConfig.ExcludeMasterFromStandardLB, api.DefaultExcludeMasterFromStandardLB)
		}
	}

	properties.AgentPoolProfiles[0].Count = 110
	setPropertiesDefaults(&mockCS, false, false)
	if helpers.IsTrueBoolPointer(properties.AgentPoolProfiles[0].SinglePlacementGroup) {
		t.Fatalf("AgentPoolProfile[0].SinglePlacementGroup did not have the expected configuration, got %t, expected %t",
			*properties.AgentPoolProfiles[0].SinglePlacementGroup, false)
	}

	if !*properties.AgentPoolProfiles[0].SinglePlacementGroup && properties.AgentPoolProfiles[0].StorageProfile != api.ManagedDisks {
		t.Fatalf("AgentPoolProfile[0].StorageProfile did not have the expected configuration, got %s, expected %s",
			properties.AgentPoolProfiles[0].StorageProfile, api.ManagedDisks)
	}

}

func TestAzureCNIVersionString(t *testing.T) {
	mockCS := getMockBaseContainerService("1.10.3")
	properties := mockCS.Properties
	properties.OrchestratorProfile.OrchestratorType = "Kubernetes"
	properties.MasterProfile.Count = 1
	properties.OrchestratorProfile.KubernetesConfig.NetworkPlugin = "azure"
	setOrchestratorDefaults(&mockCS, true)

	if properties.OrchestratorProfile.KubernetesConfig.AzureCNIVersion != AzureCniPluginVerLinux {
		t.Fatalf("Azure CNI Version string not the expected value, got %s, expected %s", properties.OrchestratorProfile.KubernetesConfig.AzureCNIVersion, AzureCniPluginVerLinux)
	}

	mockCS = getMockBaseContainerService("1.10.3")
	properties = mockCS.Properties
	properties.OrchestratorProfile.OrchestratorType = "Kubernetes"
	properties.MasterProfile.Count = 1
	properties.AgentPoolProfiles[0].OSType = "Windows"
	properties.OrchestratorProfile.KubernetesConfig.NetworkPlugin = "azure"
	setOrchestratorDefaults(&mockCS, true)

	if properties.OrchestratorProfile.KubernetesConfig.AzureCNIVersion != AzureCniPluginVerWindows {
		t.Fatalf("Azure CNI Version string not the expected value, got %s, expected %s", properties.OrchestratorProfile.KubernetesConfig.AzureCNIVersion, AzureCniPluginVerWindows)
	}

	mockCS = getMockBaseContainerService("1.10.3")
	properties = mockCS.Properties
	properties.OrchestratorProfile.OrchestratorType = "Kubernetes"
	properties.MasterProfile.Count = 1
	properties.OrchestratorProfile.KubernetesConfig.NetworkPlugin = "kubenet"
	setOrchestratorDefaults(&mockCS, true)

	if properties.OrchestratorProfile.KubernetesConfig.AzureCNIVersion != "" {
		t.Fatalf("Azure CNI Version string not the expected value, got %s, expected %s", properties.OrchestratorProfile.KubernetesConfig.AzureCNIVersion, "")
	}
}

func TestDefaultDisableRbac(t *testing.T) {
	mockCS := getMockBaseContainerService("1.10.3")
	properties := mockCS.Properties
	properties.OrchestratorProfile.OrchestratorType = "Kubernetes"
	properties.OrchestratorProfile.KubernetesConfig.EnableRbac = helpers.PointerToBool(false)
	setOrchestratorDefaults(&mockCS, true)

	if properties.OrchestratorProfile.KubernetesConfig.EnableAggregatedAPIs {
		t.Fatalf("got unexpected EnableAggregatedAPIs config value for EnableRbac=false: %t",
			properties.OrchestratorProfile.KubernetesConfig.EnableAggregatedAPIs)
	}
}

func getMockAddon(name string) api.KubernetesAddon {
	return api.KubernetesAddon{
		Name:    name,
		Enabled: helpers.PointerToBool(true),
		Containers: []api.KubernetesContainerSpec{
			{
				Name:           name,
				CPURequests:    "50m",
				MemoryRequests: "150Mi",
				CPULimits:      "50m",
				MemoryLimits:   "150Mi",
			},
		},
	}
}

func getMockBaseContainerService(orchestratorVersion string) api.ContainerService {
	mockAPIProperties := getMockAPIProperties(orchestratorVersion)
	return api.ContainerService{
		Properties: &mockAPIProperties,
	}
}

func getMockAPIProperties(orchestratorVersion string) api.Properties {
	return api.Properties{
		ProvisioningState: "",
		OrchestratorProfile: &api.OrchestratorProfile{
			OrchestratorVersion: orchestratorVersion,
			KubernetesConfig:    &api.KubernetesConfig{},
		},
		MasterProfile: &api.MasterProfile{},
		AgentPoolProfiles: []*api.AgentPoolProfile{
			{},
		}}
}

func getKubernetesConfigWithFeatureGates(featureGates string) *api.KubernetesConfig {
	return &api.KubernetesConfig{
		KubeletConfig: map[string]string{"--feature-gates": featureGates},
	}
}

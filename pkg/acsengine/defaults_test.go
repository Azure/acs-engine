package acsengine

import (
	"reflect"
	"testing"

	"github.com/Azure/acs-engine/pkg/api"
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
		Enabled: pointerToBool(true),
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
		Enabled: pointerToBool(true),
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
		Enabled: pointerToBool(true),
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

func TestPointerToBool(t *testing.T) {
	boolVar := true
	ret := pointerToBool(boolVar)
	if *ret != boolVar {
		t.Fatalf("expected pointerToBool(true) to return *true, instead returned %#v", ret)
	}
}

func getMockAddon(name string) api.KubernetesAddon {
	return api.KubernetesAddon{
		Name:    name,
		Enabled: pointerToBool(true),
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

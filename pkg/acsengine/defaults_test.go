package acsengine

import (
	"testing"

	"github.com/Azure/acs-engine/pkg/api"
	. "github.com/onsi/gomega"
)

func TestAllCertsAlreadyPresent(t *testing.T) {
	RegisterTestingT(t)
	var cert *api.CertificateProfile

	Expect(allCertAlreadyPresent(nil)).To(BeFalse())

	cert = &api.CertificateProfile{}
	Expect(allCertAlreadyPresent(cert)).To(BeFalse())

	cert = &api.CertificateProfile{
		APIServerCertificate: "a",
	}
	Expect(allCertAlreadyPresent(cert)).To(BeFalse())

	cert = &api.CertificateProfile{
		APIServerCertificate: "a",
		CaCertificate: "c",
		CaPrivateKey: "d",
		ClientCertificate: "e",
		ClientPrivateKey: "f",
		KubeConfigCertificate: "g",
		KubeConfigPrivateKey: "h"	
	}
	Expect(allCertAlreadyPresent(cert)).To(BeFalse())

	cert = &api.CertificateProfile{
		APIServerCertificate: "a",
		APIServerPrivateKey: "b",
		CaCertificate: "c",
		CaPrivateKey: "d",
		ClientCertificate: "e",
		ClientPrivateKey: "f",
		KubeConfigCertificate: "g",
		KubeConfigPrivateKey: "h"	
	}
	Expect(allCertAlreadyPresent(cert)).To(BeTrue())


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

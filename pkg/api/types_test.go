package api

import (
	"log"
	"testing"
)

const exampleCustomHyperkubeImage = `example.azurecr.io/example/hyperkube-amd64:custom`

const exampleAPIModel = `{
		"apiVersion": "vlabs",
	"properties": {
		"orchestratorProfile": {
			"orchestratorType": "Kubernetes",
			"kubernetesConfig": {
				"customHyperkubeImage": "` + exampleCustomHyperkubeImage + `"
			}
		},
		"masterProfile": { "count": 1, "dnsPrefix": "", "vmSize": "Standard_D2_v2" },
		"agentPoolProfiles": [ { "name": "linuxpool1", "count": 2, "vmSize": "Standard_D2_v2", "availabilityProfile": "AvailabilitySet" } ],
		"windowsProfile": { "adminUsername": "azureuser", "adminPassword": "replacepassword1234$" },
		"linuxProfile": { "adminUsername": "azureuser", "ssh": { "publicKeys": [ { "keyData": "" } ] }
		},
		"servicePrincipalProfile": { "clientId": "", "secret": "" }
	}
}
`

func TestIsDCOS(t *testing.T) {
	dCOSProfile := &OrchestratorProfile{
		OrchestratorType: "DCOS",
	}
	if !dCOSProfile.IsDCOS() {
		t.Fatalf("unable to detect DCOS orchestrator profile from OrchestratorType=%s", dCOSProfile.OrchestratorType)
	}
	kubernetesProfile := &OrchestratorProfile{
		OrchestratorType: "Kubernetes",
	}
	if kubernetesProfile.IsDCOS() {
		t.Fatalf("unexpectedly detected DCOS orchestrator profile from OrchestratorType=%s", kubernetesProfile.OrchestratorType)
	}
}

func TestCustomHyperkubeImageField(t *testing.T) {
	log.Println(exampleAPIModel)
	apiloader := &Apiloader{
		Translator: nil,
	}
	apimodel, _, err := apiloader.DeserializeContainerService([]byte(exampleAPIModel), false, false, nil)
	if err != nil {
		t.Fatalf("unexpectedly error deserializing the example apimodel: %s", err)
	}

	actualCustomHyperkubeImage := apimodel.Properties.OrchestratorProfile.KubernetesConfig.CustomHyperkubeImage
	if actualCustomHyperkubeImage != exampleCustomHyperkubeImage {
		t.Fatalf("kubernetesConfig->customHyperkubeImage field value was unexpected: got(%s), expected(%s)", actualCustomHyperkubeImage, exampleCustomHyperkubeImage)
	}
}

func TestKubernetesAddon(t *testing.T) {
	addon := getMockAddon("addon")
	if addon.IsEnabled(true) != true {
		t.Fatalf("KubernetesAddon.IsEnabled(true) should always return true when Enabled property is not specified")
	}

	if addon.IsEnabled(false) != false {
		t.Fatalf("KubernetesAddon.IsEnabled(false) should always return false when Enabled property is not specified")
	}
	e := true
	addon.Enabled = &e
	if addon.IsEnabled(false) != true {
		t.Fatalf("KubernetesAddon.IsEnabled(false) should always return true when Enabled property is set to true")
	}
	if addon.IsEnabled(true) != true {
		t.Fatalf("KubernetesAddon.IsEnabled(true) should always return true when Enabled property is set to true")
	}
	e = false
	addon.Enabled = &e
	if addon.IsEnabled(false) != false {
		t.Fatalf("KubernetesAddon.IsEnabled(false) should always return false when Enabled property is set to false")
	}
	if addon.IsEnabled(true) != false {
		t.Fatalf("KubernetesAddon.IsEnabled(true) should always return false when Enabled property is set to false")
	}
}

func TestIsTillerEnabled(t *testing.T) {
	c := KubernetesConfig{
		Addons: []KubernetesAddon{
			getMockAddon("addon"),
		},
	}
	e := c.IsTillerEnabled()
	if e != DefaultTillerAddonEnabled {
		t.Fatalf("KubernetesConfig.IsTillerEnabled() should return %t when no tiller addon has been specified, instead returned %t", DefaultTillerAddonEnabled, e)
	}
	c.Addons = append(c.Addons, getMockAddon("tiller"))
	e = c.IsTillerEnabled()
	if e != true {
		t.Fatalf("KubernetesConfig.IsTillerEnabled() should return true when a custom tiller addon has been specified, instead returned %t", e)
	}
	f := false
	c = KubernetesConfig{
		Addons: []KubernetesAddon{
			{
				Name:    "tiller",
				Enabled: &f,
			},
		},
	}
	e = c.IsTillerEnabled()
	if e != false {
		t.Fatalf("KubernetesConfig.IsTillerEnabled() should return false when a custom tiller addon has been specified as disabled, instead returned %t", e)
	}
}

func TestIsDashboardEnabled(t *testing.T) {
	c := KubernetesConfig{
		Addons: []KubernetesAddon{
			getMockAddon("addon"),
		},
	}
	e := c.IsDashboardEnabled()
	if e != DefaultDashboardAddonEnabled {
		t.Fatalf("KubernetesConfig.IsDashboardEnabled() should return %t when no kubernetes-dashboard addon has been specified, instead returned %t", DefaultDashboardAddonEnabled, e)
	}
	c.Addons = append(c.Addons, getMockAddon("kubernetes-dashboard"))
	e = c.IsDashboardEnabled()
	if e != true {
		t.Fatalf("KubernetesConfig.IsDashboardEnabled() should return true when a custom kubernetes-dashboard addon has been specified, instead returned %t", e)
	}
	f := false
	c = KubernetesConfig{
		Addons: []KubernetesAddon{
			{
				Name:    "kubernetes-dashboard",
				Enabled: &f,
			},
		},
	}
	e = c.IsDashboardEnabled()
	if e != false {
		t.Fatalf("KubernetesConfig.IsDashboardEnabled() should return false when a custom kubernetes-dashboard addon has been specified as disabled, instead returned %t", e)
	}
}

func getMockAddon(name string) KubernetesAddon {
	return KubernetesAddon{
		Name: name,
		Containers: []KubernetesContainerSpec{
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

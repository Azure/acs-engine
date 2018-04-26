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
	if !addon.IsEnabled(true) {
		t.Fatalf("KubernetesAddon.IsEnabled(true) should always return true when Enabled property is not specified")
	}

	if addon.IsEnabled(false) {
		t.Fatalf("KubernetesAddon.IsEnabled(false) should always return false when Enabled property is not specified")
	}
	e := true
	addon.Enabled = &e
	if !addon.IsEnabled(false) {
		t.Fatalf("KubernetesAddon.IsEnabled(false) should always return true when Enabled property is set to true")
	}
	if !addon.IsEnabled(true) {
		t.Fatalf("KubernetesAddon.IsEnabled(true) should always return true when Enabled property is set to true")
	}
	e = false
	addon.Enabled = &e
	if addon.IsEnabled(false) {
		t.Fatalf("KubernetesAddon.IsEnabled(false) should always return false when Enabled property is set to false")
	}
	if addon.IsEnabled(true) {
		t.Fatalf("KubernetesAddon.IsEnabled(true) should always return false when Enabled property is set to false")
	}
}

func TestIsTillerEnabled(t *testing.T) {
	c := KubernetesConfig{
		Addons: []KubernetesAddon{
			getMockAddon("addon"),
		},
	}
	enabled := c.IsTillerEnabled()
	if enabled != DefaultTillerAddonEnabled {
		t.Fatalf("KubernetesConfig.IsTillerEnabled() should return %t when no tiller addon has been specified, instead returned %t", DefaultTillerAddonEnabled, enabled)
	}
	c.Addons = append(c.Addons, getMockAddon(DefaultTillerAddonName))
	enabled = c.IsTillerEnabled()
	if !enabled {
		t.Fatalf("KubernetesConfig.IsTillerEnabled() should return true when a custom tiller addon has been specified, instead returned %t", enabled)
	}
	b := false
	c = KubernetesConfig{
		Addons: []KubernetesAddon{
			{
				Name:    DefaultTillerAddonName,
				Enabled: &b,
			},
		},
	}
	enabled = c.IsTillerEnabled()
	if enabled {
		t.Fatalf("KubernetesConfig.IsTillerEnabled() should return false when a custom tiller addon has been specified as disabled, instead returned %t", enabled)
	}
}

func TestIsACIConnectorEnabled(t *testing.T) {
	c := KubernetesConfig{
		Addons: []KubernetesAddon{
			getMockAddon("addon"),
		},
	}
	enabled := c.IsACIConnectorEnabled()
	if enabled != DefaultACIConnectorAddonEnabled {
		t.Fatalf("KubernetesConfig.IsACIConnectorEnabled() should return %t when no ACI connector addon has been specified, instead returned %t", DefaultACIConnectorAddonEnabled, enabled)
	}
	c.Addons = append(c.Addons, getMockAddon(DefaultACIConnectorAddonName))
	enabled = c.IsACIConnectorEnabled()
	if enabled {
		t.Fatalf("KubernetesConfig.IsACIConnectorEnabled() should return true when ACI connector has been specified, instead returned %t", enabled)
	}
	b := true
	c = KubernetesConfig{
		Addons: []KubernetesAddon{
			{
				Name:    DefaultACIConnectorAddonName,
				Enabled: &b,
			},
		},
	}
	enabled = c.IsACIConnectorEnabled()
	if !enabled {
		t.Fatalf("KubernetesConfig.IsACIConnectorEnabled() should return false when ACI connector addon has been specified as disabled, instead returned %t", enabled)
	}
}

func TestIsDashboardEnabled(t *testing.T) {
	c := KubernetesConfig{
		Addons: []KubernetesAddon{
			getMockAddon("addon"),
		},
	}
	enabled := c.IsDashboardEnabled()
	if enabled != DefaultDashboardAddonEnabled {
		t.Fatalf("KubernetesConfig.IsDashboardEnabled() should return %t when no kubernetes-dashboard addon has been specified, instead returned %t", DefaultDashboardAddonEnabled, enabled)
	}
	c.Addons = append(c.Addons, getMockAddon(DefaultDashboardAddonName))
	enabled = c.IsDashboardEnabled()
	if !enabled {
		t.Fatalf("KubernetesConfig.IsDashboardEnabled() should return true when a custom kubernetes-dashboard addon has been specified, instead returned %t", enabled)
	}
	b := false
	c = KubernetesConfig{
		Addons: []KubernetesAddon{
			{
				Name:    DefaultDashboardAddonName,
				Enabled: &b,
			},
		},
	}
	enabled = c.IsDashboardEnabled()
	if enabled {
		t.Fatalf("KubernetesConfig.IsDashboardEnabled() should return false when a custom kubernetes-dashboard addon has been specified as disabled, instead returned %t", enabled)
	}
}

func TestIsReschedulerEnabled(t *testing.T) {
	c := KubernetesConfig{
		Addons: []KubernetesAddon{
			getMockAddon("addon"),
		},
	}
	enabled := c.IsReschedulerEnabled()
	if enabled != DefaultReschedulerAddonEnabled {
		t.Fatalf("KubernetesConfig.IsReschedulerEnabled() should return %t when no rescheduler addon has been specified, instead returned %t", DefaultReschedulerAddonEnabled, enabled)
	}
	c.Addons = append(c.Addons, getMockAddon(DefaultReschedulerAddonName))
	enabled = c.IsReschedulerEnabled()
	if enabled {
		t.Fatalf("KubernetesConfig.IsReschedulerEnabled() should return true when a custom rescheduler addon has been specified, instead returned %t", enabled)
	}
	b := true
	c = KubernetesConfig{
		Addons: []KubernetesAddon{
			{
				Name:    DefaultReschedulerAddonName,
				Enabled: &b,
			},
		},
	}
	enabled = c.IsReschedulerEnabled()
	if !enabled {
		t.Fatalf("KubernetesConfig.IsReschedulerEnabled() should return false when a custom rescheduler addon has been specified as enabled, instead returned %t", enabled)
	}
}

func TestIsMetricsServerEnabled(t *testing.T) {
	v := "1.8.0"
	o := OrchestratorProfile{
		OrchestratorType:    "Kubernetes",
		OrchestratorVersion: v,
		KubernetesConfig: &KubernetesConfig{Addons: []KubernetesAddon{
			getMockAddon("addon"),
		},
		},
	}
	enabled := o.IsMetricsServerEnabled()
	if enabled != DefaultMetricsServerAddonEnabled {
		t.Fatalf("KubernetesConfig.IsMetricsServerEnabled() should return %t for kubernetes version %s when no metrics-server addon has been specified, instead returned %t", DefaultMetricsServerAddonEnabled, v, enabled)
	}

	o.KubernetesConfig.Addons = append(o.KubernetesConfig.Addons, getMockAddon(DefaultMetricsServerAddonName))
	enabled = o.IsMetricsServerEnabled()
	if enabled != DefaultMetricsServerAddonEnabled {
		t.Fatalf("KubernetesConfig.IsMetricsServerEnabled() should return %t for kubernetes version %s when the metrics-server addon has been specified, instead returned %t", DefaultMetricsServerAddonEnabled, v, enabled)
	}

	b := true
	o = OrchestratorProfile{
		OrchestratorType:    "Kubernetes",
		OrchestratorVersion: v,
		KubernetesConfig: &KubernetesConfig{Addons: []KubernetesAddon{
			{
				Name:    DefaultMetricsServerAddonName,
				Enabled: &b,
			},
		},
		},
	}
	enabled = o.IsMetricsServerEnabled()
	if !enabled {
		t.Fatalf("KubernetesConfig.IsMetricsServerEnabled() should return true for kubernetes version %s when the metrics-server addon has been specified as enabled, instead returned %t", v, enabled)
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

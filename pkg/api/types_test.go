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

func TestIsAzureCNI(t *testing.T) {
	k := &KubernetesConfig{
		NetworkPlugin: "azure",
	}

	o := &OrchestratorProfile{
		KubernetesConfig: k,
	}
	if !o.IsAzureCNI() {
		t.Fatalf("unable to detect orchestrator profile is using Azure CNI from NetworkPlugin=%s", o.KubernetesConfig.NetworkPlugin)
	}

	k = &KubernetesConfig{
		NetworkPlugin: "none",
	}

	o = &OrchestratorProfile{
		KubernetesConfig: k,
	}
	if o.IsAzureCNI() {
		t.Fatalf("unable to detect orchestrator profile is not using Azure CNI from NetworkPlugin=%s", o.KubernetesConfig.NetworkPlugin)
	}

	o = &OrchestratorProfile{}
	if o.IsAzureCNI() {
		t.Fatalf("unable to detect orchestrator profile is not using Azure CNI from nil KubernetesConfig")
	}
}

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

func TestIsClusterAutoscalerEnabled(t *testing.T) {
	c := KubernetesConfig{
		Addons: []KubernetesAddon{
			getMockAddon("addon"),
		},
	}
	enabled := c.IsClusterAutoscalerEnabled()
	if enabled != DefaultClusterAutoscalerAddonEnabled {
		t.Fatalf("KubernetesConfig.IsAutoscalerEnabled() should return %t when no cluster autoscaler addon has been specified, instead returned %t", DefaultClusterAutoscalerAddonEnabled, enabled)
	}
	c.Addons = append(c.Addons, getMockAddon(DefaultClusterAutoscalerAddonName))
	enabled = c.IsClusterAutoscalerEnabled()
	if enabled {
		t.Fatalf("KubernetesConfig.IsClusterAutoscalerEnabled() should return true when cluster autoscaler has been specified, instead returned %t", enabled)
	}
	b := true
	c = KubernetesConfig{
		Addons: []KubernetesAddon{
			{
				Name:    DefaultClusterAutoscalerAddonName,
				Enabled: &b,
			},
		},
	}
	enabled = c.IsClusterAutoscalerEnabled()
	if !enabled {
		t.Fatalf("KubernetesConfig.IsClusterAutoscalerEnabled() should return false when cluster autoscaler addon has been specified as disabled, instead returned %t", enabled)
	}
}

func TestIsNVIDIADevicePluginEnabled(t *testing.T) {
	p := Properties{
		AgentPoolProfiles: []*AgentPoolProfile{
			{
				Name:   "agentpool",
				VMSize: "Standard_N",
				Count:  1,
			},
		},
		OrchestratorProfile: &OrchestratorProfile{
			OrchestratorType:    Kubernetes,
			OrchestratorVersion: "1.9.0",
			KubernetesConfig: &KubernetesConfig{
				Addons: []KubernetesAddon{
					getMockAddon("addon"),
				},
			},
		},
	}
	enabled := p.IsNVIDIADevicePluginEnabled()
	if enabled == isNSeriesSKU(&p) {
		t.Fatalf("KubernetesConfig.IsNVIDIADevicePluginEnabled() should return false with N-series VMs with < k8s 1.10, instead returned %t", enabled)
	}

	o := p.OrchestratorProfile
	o.OrchestratorVersion = "1.10.0"
	enabled = p.IsNVIDIADevicePluginEnabled()
	if enabled != isNSeriesSKU(&p) {
		t.Fatalf("KubernetesConfig.IsNVIDIADevicePluginEnabled() should return %t with N-series VMs with k8s >= 1.10, instead returned %t", isNSeriesSKU(&p), enabled)
	}

	b := false
	c := p.OrchestratorProfile.KubernetesConfig
	c.Addons = []KubernetesAddon{
		{
			Name:    DefaultNVIDIADevicePluginAddonName,
			Enabled: &b,
		},
	}
	enabled = p.IsNVIDIADevicePluginEnabled()
	if enabled {
		t.Fatalf("KubernetesConfig.IsNVIDIADevicePluginEnabled() should return false when explicitly disabled")
	}
}

func TestIsContainerMonitoringEnabled(t *testing.T) {
	v := "1.9.0"
	o := OrchestratorProfile{
		OrchestratorType:    "Kubernetes",
		OrchestratorVersion: v,
		KubernetesConfig: &KubernetesConfig{Addons: []KubernetesAddon{
			getMockAddon("addon"),
		},
		},
	}
	enabled := o.IsContainerMonitoringEnabled()
	if enabled != DefaultContainerMonitoringAddonEnabled {
		t.Fatalf("KubernetesConfig.IsContainerMonitoringEnabled() should return %t for kubernetes version %s when no container-monitoring addon has been specified, instead returned %t", DefaultContainerMonitoringAddonEnabled, v, enabled)
	}

	b := true
	cm := getMockAddon(ContainerMonitoringAddonName)
	cm.Enabled = &b
	o.KubernetesConfig.Addons = append(o.KubernetesConfig.Addons, cm)
	enabled = o.IsContainerMonitoringEnabled()
	if enabled != true {
		t.Fatalf("KubernetesConfig.IsContainerMonitoringEnabled() should return %t for kubernetes version %s when the container-monitoring addon has been specified, instead returned %t", true, v, enabled)
	}

	b = false
	o = OrchestratorProfile{
		OrchestratorType:    "Kubernetes",
		OrchestratorVersion: v,
		KubernetesConfig: &KubernetesConfig{Addons: []KubernetesAddon{
			{
				Name:    ContainerMonitoringAddonName,
				Enabled: &b,
			},
		},
		},
	}
	enabled = o.IsContainerMonitoringEnabled()
	if enabled {
		t.Fatalf("KubernetesConfig.IsContainerMonitoringEnabled() should return false when a custom container monitoring addon has been specified as disabled, instead returned %t", enabled)
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

package acsengine

import (
	"strconv"

	"github.com/Azure/acs-engine/pkg/api"
	"github.com/Azure/acs-engine/pkg/api/common"
	"github.com/Azure/acs-engine/pkg/helpers"
)

func setAddonsConfig(cs *api.ContainerService) {
	o := cs.Properties.OrchestratorProfile
	defaultTillerAddonsConfig := api.KubernetesAddon{
		Name:    DefaultTillerAddonName,
		Enabled: helpers.PointerToBool(api.DefaultTillerAddonEnabled),
		Containers: []api.KubernetesContainerSpec{
			{
				Name:           DefaultTillerAddonName,
				CPURequests:    "50m",
				MemoryRequests: "150Mi",
				CPULimits:      "50m",
				MemoryLimits:   "150Mi",
			},
		},
		Config: map[string]string{
			"max-history": strconv.Itoa(DefaultTillerMaxHistory),
		},
	}

	defaultACIConnectorAddonsConfig := api.KubernetesAddon{
		Name:    DefaultACIConnectorAddonName,
		Enabled: helpers.PointerToBool(api.DefaultACIConnectorAddonEnabled),
		Config: map[string]string{
			"region":   "westus",
			"nodeName": "aci-connector",
			"os":       "Linux",
			"taint":    "azure.com/aci",
		},
		Containers: []api.KubernetesContainerSpec{
			{
				Name:           DefaultACIConnectorAddonName,
				CPURequests:    "50m",
				MemoryRequests: "150Mi",
				CPULimits:      "50m",
				MemoryLimits:   "150Mi",
			},
		},
	}

	defaultClusterAutoscalerAddonsConfig := api.KubernetesAddon{
		Name:    DefaultClusterAutoscalerAddonName,
		Enabled: helpers.PointerToBool(api.DefaultClusterAutoscalerAddonEnabled),
		Config: map[string]string{
			"minNodes": "1",
			"maxNodes": "5",
		},
		Containers: []api.KubernetesContainerSpec{
			{
				Name:           DefaultClusterAutoscalerAddonName,
				CPURequests:    "100m",
				MemoryRequests: "300Mi",
				CPULimits:      "100m",
				MemoryLimits:   "300Mi",
			},
		},
	}

	defaultBlobfuseFlexVolumeAddonsConfig := api.KubernetesAddon{
		Name:    DefaultBlobfuseFlexVolumeAddonName,
		Enabled: helpers.PointerToBool(common.IsKubernetesVersionGe(o.OrchestratorVersion, "1.8.0") && api.DefaultBlobfuseFlexVolumeAddonEnabled),
		Containers: []api.KubernetesContainerSpec{
			{
				Name:           DefaultBlobfuseFlexVolumeAddonName,
				CPURequests:    "50m",
				MemoryRequests: "10Mi",
				CPULimits:      "50m",
				MemoryLimits:   "10Mi",
			},
		},
	}

	defaultSMBFlexVolumeAddonsConfig := api.KubernetesAddon{
		Name:    DefaultSMBFlexVolumeAddonName,
		Enabled: helpers.PointerToBool(common.IsKubernetesVersionGe(o.OrchestratorVersion, "1.8.0") && api.DefaultSMBFlexVolumeAddonEnabled),
		Containers: []api.KubernetesContainerSpec{
			{
				Name:           DefaultSMBFlexVolumeAddonName,
				CPURequests:    "50m",
				MemoryRequests: "10Mi",
				CPULimits:      "50m",
				MemoryLimits:   "10Mi",
			},
		},
	}

	defaultKeyVaultFlexVolumeAddonsConfig := api.KubernetesAddon{
		Name:    DefaultKeyVaultFlexVolumeAddonName,
		Enabled: helpers.PointerToBool(api.DefaultKeyVaultFlexVolumeAddonEnabled),
		Containers: []api.KubernetesContainerSpec{
			{
				Name:           DefaultKeyVaultFlexVolumeAddonName,
				CPURequests:    "50m",
				MemoryRequests: "10Mi",
				CPULimits:      "50m",
				MemoryLimits:   "10Mi",
			},
		},
	}

	defaultDashboardAddonsConfig := api.KubernetesAddon{
		Name:    DefaultDashboardAddonName,
		Enabled: helpers.PointerToBool(api.DefaultDashboardAddonEnabled),
		Containers: []api.KubernetesContainerSpec{
			{
				Name:           DefaultDashboardAddonName,
				CPURequests:    "300m",
				MemoryRequests: "150Mi",
				CPULimits:      "300m",
				MemoryLimits:   "150Mi",
			},
		},
	}

	defaultReschedulerAddonsConfig := api.KubernetesAddon{
		Name:    DefaultReschedulerAddonName,
		Enabled: helpers.PointerToBool(api.DefaultReschedulerAddonEnabled),
		Containers: []api.KubernetesContainerSpec{
			{
				Name:           DefaultReschedulerAddonName,
				CPURequests:    "10m",
				MemoryRequests: "100Mi",
				CPULimits:      "10m",
				MemoryLimits:   "100Mi",
			},
		},
	}

	defaultMetricsServerAddonsConfig := api.KubernetesAddon{
		Name:    DefaultMetricsServerAddonName,
		Enabled: k8sVersionMetricsServerAddonEnabled(o),
		Containers: []api.KubernetesContainerSpec{
			{
				Name: DefaultMetricsServerAddonName,
			},
		},
	}

	defaultNVIDIADevicePluginAddonsConfig := api.KubernetesAddon{
		Name:    NVIDIADevicePluginAddonName,
		Enabled: helpers.PointerToBool(api.IsNSeriesSKU(cs.Properties) && common.IsKubernetesVersionGe(o.OrchestratorVersion, "1.10.0")),
		Containers: []api.KubernetesContainerSpec{
			{
				Name: NVIDIADevicePluginAddonName,
				// from https://github.com/kubernetes/kubernetes/blob/master/cluster/addons/device-plugins/nvidia-gpu/daemonset.yaml#L44
				CPURequests:    "50m",
				MemoryRequests: "10Mi",
				CPULimits:      "50m",
				MemoryLimits:   "10Mi",
			},
		},
	}

	defaultContainerMonitoringAddonsConfig := api.KubernetesAddon{
		Name:    ContainerMonitoringAddonName,
		Enabled: helpers.PointerToBool(api.DefaultContainerMonitoringAddonEnabled),
		Config: map[string]string{
			"omsAgentVersion":       "1.6.0-42",
			"dockerProviderVersion": "2.0.0-3",
		},
		Containers: []api.KubernetesContainerSpec{
			{
				Name:           "omsagent",
				Image:          "microsoft/oms:acsenginelogfixnew",
				CPURequests:    "50m",
				MemoryRequests: "200Mi",
				CPULimits:      "150m",
				MemoryLimits:   "750Mi",
			},
		},
	}

	defaultAzureCNINetworkMonitorAddonsConfig := api.KubernetesAddon{
		Name:    AzureCNINetworkMonitoringAddonName,
		Enabled: azureCNINetworkMonitorAddonEnabled(o),
		Containers: []api.KubernetesContainerSpec{
			{
				Name: AzureCNINetworkMonitoringAddonName,
			},
		},
	}

	defaultAzureNetworkPolicyAddonsConfig := api.KubernetesAddon{
		Name:    AzureNetworkPolicyAddonName,
		Enabled: azureNetworkPolicyAddonEnabled(o),
		Containers: []api.KubernetesContainerSpec{
			{
				Name: AzureNetworkPolicyAddonName,
			},
		},
	}

	defaultAddons := []api.KubernetesAddon{
		defaultTillerAddonsConfig,
		defaultACIConnectorAddonsConfig,
		defaultClusterAutoscalerAddonsConfig,
		defaultBlobfuseFlexVolumeAddonsConfig,
		defaultSMBFlexVolumeAddonsConfig,
		defaultKeyVaultFlexVolumeAddonsConfig,
		defaultDashboardAddonsConfig,
		defaultReschedulerAddonsConfig,
		defaultMetricsServerAddonsConfig,
		defaultNVIDIADevicePluginAddonsConfig,
		defaultContainerMonitoringAddonsConfig,
		defaultAzureCNINetworkMonitorAddonsConfig,
		defaultAzureNetworkPolicyAddonsConfig,
	}
	// Add default addons specification, if no user-provided spec exists
	if o.KubernetesConfig.Addons == nil {
		o.KubernetesConfig.Addons = defaultAddons
	} else {
		for _, addon := range defaultAddons {
			i := getAddonsIndexByName(o.KubernetesConfig.Addons, addon.Name)
			if i < 0 {
				o.KubernetesConfig.Addons = append(o.KubernetesConfig.Addons, addon)
			}
		}
	}

	for _, addon := range defaultAddons {
		synthesizeAddonsConfig(o.KubernetesConfig.Addons, addon, false)
	}
}

func getAddonsIndexByName(addons []api.KubernetesAddon, name string) int {
	for i := range addons {
		if addons[i].Name == name {
			return i
		}
	}
	return -1
}

func getAddonContainersIndexByName(containers []api.KubernetesContainerSpec, name string) int {
	for i := range containers {
		if containers[i].Name == name {
			return i
		}
	}
	return -1
}

// assignDefaultAddonVals will assign default values to addon from defaults, for each property in addon that has a zero value
func assignDefaultAddonVals(addon, defaults api.KubernetesAddon) api.KubernetesAddon {
	if addon.Enabled == nil {
		addon.Enabled = defaults.Enabled
	}
	for i := range defaults.Containers {
		c := getAddonContainersIndexByName(addon.Containers, defaults.Containers[i].Name)
		if c < 0 {
			addon.Containers = append(addon.Containers, defaults.Containers[i])
		} else {
			if addon.Containers[c].Image == "" {
				addon.Containers[c].Image = defaults.Containers[i].Image
			}
			if addon.Containers[c].CPURequests == "" {
				addon.Containers[c].CPURequests = defaults.Containers[i].CPURequests
			}
			if addon.Containers[c].MemoryRequests == "" {
				addon.Containers[c].MemoryRequests = defaults.Containers[i].MemoryRequests
			}
			if addon.Containers[c].CPULimits == "" {
				addon.Containers[c].CPULimits = defaults.Containers[i].CPULimits
			}
			if addon.Containers[c].MemoryLimits == "" {
				addon.Containers[c].MemoryLimits = defaults.Containers[i].MemoryLimits
			}
		}
	}
	for key, val := range defaults.Config {
		if addon.Config == nil {
			addon.Config = make(map[string]string)
		}
		if v, ok := addon.Config[key]; !ok || v == "" {
			addon.Config[key] = val
		}
	}
	return addon
}

func synthesizeAddonsConfig(addons []api.KubernetesAddon, addon api.KubernetesAddon, enableIfNil bool) {
	i := getAddonsIndexByName(addons, addon.Name)
	if i >= 0 {
		if addons[i].IsEnabled(enableIfNil) {
			addons[i] = assignDefaultAddonVals(addons[i], addon)
		}
	}
}

func k8sVersionMetricsServerAddonEnabled(o *api.OrchestratorProfile) *bool {
	return helpers.PointerToBool(common.IsKubernetesVersionGe(o.OrchestratorVersion, "1.9.0"))
}

func azureNetworkPolicyAddonEnabled(o *api.OrchestratorProfile) *bool {
	return helpers.PointerToBool(o.KubernetesConfig.NetworkPlugin == NetworkPluginAzure && o.KubernetesConfig.NetworkPolicy == NetworkPolicyAzure)
}

func azureCNINetworkMonitorAddonEnabled(o *api.OrchestratorProfile) *bool {
	return helpers.PointerToBool(o.IsAzureCNI())
}

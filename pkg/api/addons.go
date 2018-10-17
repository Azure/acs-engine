package api

import (
	"strconv"

	"github.com/Azure/acs-engine/pkg/api/common"
	"github.com/Azure/acs-engine/pkg/helpers"
)

func (cs *ContainerService) setAddonsConfig() {
	o := cs.Properties.OrchestratorProfile
	defaultTillerAddonsConfig := KubernetesAddon{
		Name:    DefaultTillerAddonName,
		Enabled: helpers.PointerToBool(DefaultTillerAddonEnabled),
		Containers: []KubernetesContainerSpec{
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

	defaultACIConnectorAddonsConfig := KubernetesAddon{
		Name:    DefaultACIConnectorAddonName,
		Enabled: helpers.PointerToBool(DefaultACIConnectorAddonEnabled),
		Config: map[string]string{
			"region":   "westus",
			"nodeName": "aci-connector",
			"os":       "Linux",
			"taint":    "azure.com/aci",
		},
		Containers: []KubernetesContainerSpec{
			{
				Name:           DefaultACIConnectorAddonName,
				CPURequests:    "50m",
				MemoryRequests: "150Mi",
				CPULimits:      "50m",
				MemoryLimits:   "150Mi",
			},
		},
	}

	defaultClusterAutoscalerAddonsConfig := KubernetesAddon{
		Name:    DefaultClusterAutoscalerAddonName,
		Enabled: helpers.PointerToBool(DefaultClusterAutoscalerAddonEnabled),
		Config: map[string]string{
			"minNodes": "1",
			"maxNodes": "5",
		},
		Containers: []KubernetesContainerSpec{
			{
				Name:           DefaultClusterAutoscalerAddonName,
				CPURequests:    "100m",
				MemoryRequests: "300Mi",
				CPULimits:      "100m",
				MemoryLimits:   "300Mi",
			},
		},
	}

	defaultBlobfuseFlexVolumeAddonsConfig := KubernetesAddon{
		Name:    DefaultBlobfuseFlexVolumeAddonName,
		Enabled: helpers.PointerToBool(common.IsKubernetesVersionGe(o.OrchestratorVersion, "1.8.0") && DefaultBlobfuseFlexVolumeAddonEnabled),
		Containers: []KubernetesContainerSpec{
			{
				Name:           DefaultBlobfuseFlexVolumeAddonName,
				CPURequests:    "50m",
				MemoryRequests: "10Mi",
				CPULimits:      "50m",
				MemoryLimits:   "10Mi",
			},
		},
	}

	defaultSMBFlexVolumeAddonsConfig := KubernetesAddon{
		Name:    DefaultSMBFlexVolumeAddonName,
		Enabled: helpers.PointerToBool(common.IsKubernetesVersionGe(o.OrchestratorVersion, "1.8.0") && DefaultSMBFlexVolumeAddonEnabled),
		Containers: []KubernetesContainerSpec{
			{
				Name:           DefaultSMBFlexVolumeAddonName,
				CPURequests:    "50m",
				MemoryRequests: "10Mi",
				CPULimits:      "50m",
				MemoryLimits:   "10Mi",
			},
		},
	}

	defaultKeyVaultFlexVolumeAddonsConfig := KubernetesAddon{
		Name:    DefaultKeyVaultFlexVolumeAddonName,
		Enabled: helpers.PointerToBool(DefaultKeyVaultFlexVolumeAddonEnabled),
		Containers: []KubernetesContainerSpec{
			{
				Name:           DefaultKeyVaultFlexVolumeAddonName,
				CPURequests:    "50m",
				MemoryRequests: "10Mi",
				CPULimits:      "50m",
				MemoryLimits:   "10Mi",
			},
		},
	}

	defaultDashboardAddonsConfig := KubernetesAddon{
		Name:    DefaultDashboardAddonName,
		Enabled: helpers.PointerToBool(DefaultDashboardAddonEnabled),
		Containers: []KubernetesContainerSpec{
			{
				Name:           DefaultDashboardAddonName,
				CPURequests:    "300m",
				MemoryRequests: "150Mi",
				CPULimits:      "300m",
				MemoryLimits:   "150Mi",
			},
		},
	}

	defaultReschedulerAddonsConfig := KubernetesAddon{
		Name:    DefaultReschedulerAddonName,
		Enabled: helpers.PointerToBool(DefaultReschedulerAddonEnabled),
		Containers: []KubernetesContainerSpec{
			{
				Name:           DefaultReschedulerAddonName,
				CPURequests:    "10m",
				MemoryRequests: "100Mi",
				CPULimits:      "10m",
				MemoryLimits:   "100Mi",
			},
		},
	}

	defaultMetricsServerAddonsConfig := KubernetesAddon{
		Name:    DefaultMetricsServerAddonName,
		Enabled: k8sVersionMetricsServerAddonEnabled(o),
		Containers: []KubernetesContainerSpec{
			{
				Name: DefaultMetricsServerAddonName,
			},
		},
	}

	defaultNVIDIADevicePluginAddonsConfig := KubernetesAddon{
		Name:    NVIDIADevicePluginAddonName,
		Enabled: helpers.PointerToBool(IsNSeriesSKU(cs.Properties) && common.IsKubernetesVersionGe(o.OrchestratorVersion, "1.10.0")),
		Containers: []KubernetesContainerSpec{
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

	defaultContainerMonitoringAddonsConfig := KubernetesAddon{
		Name:    ContainerMonitoringAddonName,
		Enabled: helpers.PointerToBool(DefaultContainerMonitoringAddonEnabled),
		Config: map[string]string{
			"omsAgentVersion":       "1.6.0-42",
			"dockerProviderVersion": "2.0.0-3",
		},
		Containers: []KubernetesContainerSpec{
			{
				Name:           "omsagent",
				Image:          "microsoft/oms:ciprod10162018-2",
				CPURequests:    "50m",
				MemoryRequests: "200Mi",
				CPULimits:      "150m",
				MemoryLimits:   "750Mi",
			},
		},
	}

	defaultIPMasqAgentAddonsConfig := KubernetesAddon{
		Name:    IPMASQAgentAddonName,
		Enabled: helpers.PointerToBool(IPMasqAgentAddonEnabled),
		Containers: []KubernetesContainerSpec{
			{
				Name:           IPMASQAgentAddonName,
				CPURequests:    "50m",
				MemoryRequests: "50Mi",
				CPULimits:      "50m",
				MemoryLimits:   "250Mi",
			},
		},
	}

	defaultAzureCNINetworkMonitorAddonsConfig := KubernetesAddon{
		Name:    AzureCNINetworkMonitoringAddonName,
		Enabled: azureCNINetworkMonitorAddonEnabled(o),
		Containers: []KubernetesContainerSpec{
			{
				Name: AzureCNINetworkMonitoringAddonName,
			},
		},
	}

	defaultAzureNetworkPolicyAddonsConfig := KubernetesAddon{
		Name:    AzureNetworkPolicyAddonName,
		Enabled: azureNetworkPolicyAddonEnabled(o),
		Containers: []KubernetesContainerSpec{
			{
				Name: AzureNetworkPolicyAddonName,
			},
		},
	}

	defaultAddons := []KubernetesAddon{
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
		defaultIPMasqAgentAddonsConfig,
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

func getAddonsIndexByName(addons []KubernetesAddon, name string) int {
	for i := range addons {
		if addons[i].Name == name {
			return i
		}
	}
	return -1
}

// assignDefaultAddonVals will assign default values to addon from defaults, for each property in addon that has a zero value
func assignDefaultAddonVals(addon, defaults KubernetesAddon) KubernetesAddon {
	if addon.Enabled == nil {
		addon.Enabled = defaults.Enabled
	}
	for i := range defaults.Containers {
		c := addon.GetAddonContainersIndexByName(defaults.Containers[i].Name)
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

func synthesizeAddonsConfig(addons []KubernetesAddon, addon KubernetesAddon, enableIfNil bool) {
	i := getAddonsIndexByName(addons, addon.Name)
	if i >= 0 {
		if addons[i].IsEnabled(enableIfNil) {
			addons[i] = assignDefaultAddonVals(addons[i], addon)
		}
	}
}

func k8sVersionMetricsServerAddonEnabled(o *OrchestratorProfile) *bool {
	return helpers.PointerToBool(common.IsKubernetesVersionGe(o.OrchestratorVersion, "1.9.0"))
}

func azureNetworkPolicyAddonEnabled(o *OrchestratorProfile) *bool {
	return helpers.PointerToBool(o.KubernetesConfig.NetworkPlugin == NetworkPluginAzure && o.KubernetesConfig.NetworkPolicy == NetworkPolicyAzure)
}

func azureCNINetworkMonitorAddonEnabled(o *OrchestratorProfile) *bool {
	return helpers.PointerToBool(o.IsAzureCNI())
}

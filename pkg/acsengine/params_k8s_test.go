package acsengine

import (
	"path"
	"testing"

	"github.com/Azure/acs-engine/pkg/helpers"

	"github.com/Azure/acs-engine/pkg/api"
	"github.com/Azure/acs-engine/pkg/i18n"
	"github.com/leonelquinteros/gotext"
)

func TestAssignKubernetesParameters(t *testing.T) {
	// Initialize locale for translation
	locale := gotext.NewLocale(path.Join("..", "..", "translations"), "en_US")
	i18n.Initialize(locale)

	apiloader := &api.Apiloader{
		Translator: &i18n.Translator{
			Locale: locale,
		},
	}
	// iterate the test data directory
	apiModelTestFiles := &[]APIModelTestFile{}
	if e := IterateTestFilesDirectory(TestDataDir, apiModelTestFiles); e != nil {
		t.Error(e.Error())
		return
	}

	for _, tuple := range *apiModelTestFiles {
		containerService, _, err := apiloader.LoadContainerServiceFromFile(tuple.APIModelFilename, true, false, nil)
		if err != nil {
			t.Errorf("Loading file %s got error: %s", tuple.APIModelFilename, err.Error())
			continue
		}

		parametersMap := paramsMap{}
		containerService.Location = "eatsus"
		cloudSpecConfig := containerService.GetCloudSpecConfig()
		assignKubernetesParameters(containerService.Properties, parametersMap, cloudSpecConfig, DefaultGeneratorCode)
		for k, v := range parametersMap {
			switch val := v.(paramsMap)["value"].(type) {
			case *bool:
				t.Errorf("got a pointer to bool in paramsMap value, this is dangerous!: %s: %v", k, val)
			}
		}
	}
}

func TestKubernetesParamsAddons(t *testing.T) {
	tests := map[string]struct {
		addon          api.KubernetesAddon
		expectedParams []string
	}{
		"nvidia-device-plugin": {
			api.KubernetesAddon{
				Name:    api.NVIDIADevicePluginAddonName,
				Enabled: helpers.PointerToBool(true),
				Containers: []api.KubernetesContainerSpec{
					{
						Name:           api.NVIDIADevicePluginAddonName,
						CPURequests:    "50m",
						MemoryRequests: "150Mi",
						CPULimits:      "50m",
						MemoryLimits:   "150Mi",
					},
				},
			},
			[]string{
				"kubernetesNVIDIADevicePluginCPURequests",
				"kubernetesNVIDIADevicePluginMemoryRequests",
				"kubernetesNVIDIADevicePluginCPULimit",
				"kubernetesNVIDIADevicePluginMemoryLimit",
				"kubernetesNVIDIADevicePluginSpec",
			},
		},

		"container-monitoring": {
			api.KubernetesAddon{
				Name:    api.ContainerMonitoringAddonName,
				Enabled: helpers.PointerToBool(true),
				Containers: []api.KubernetesContainerSpec{
					{
						Name:           "omsagent",
						CPURequests:    "50m",
						MemoryRequests: "150Mi",
						CPULimits:      "50m",
						MemoryLimits:   "150Mi",
					},
				},
			},
			[]string{
				"omsAgentVersion",
				"omsAgentDockerProviderVersion",
				"omsAgentWorkspaceGuid",
				"omsAgentWorkspaceKey",
				"kubernetesOMSAgentCPURequests",
				"kubernetesOMSAgentCPULimit",
				"kubernetesOMSAgentMemoryLimit",
				"kubernetesOMSAgentMemoryRequests",
				"omsAgentImage",
			},
		},

		"kubernetes-dashboard": {
			api.KubernetesAddon{
				Name:    api.DefaultDashboardAddonName,
				Enabled: helpers.PointerToBool(true),
				Containers: []api.KubernetesContainerSpec{
					{
						Name:           api.DefaultDashboardAddonName,
						CPURequests:    "50m",
						MemoryRequests: "150Mi",
						CPULimits:      "50m",
						MemoryLimits:   "150Mi",
					},
				},
			},
			[]string{
				"kubernetesDashboardCPURequests",
				"kubernetesDashboardCPULimit",
				"kubernetesDashboardMemoryRequests",
				"kubernetesDashboardMemoryLimit",
				"kubernetesDashboardSpec",
			},
		},

		"rescheduler": {
			api.KubernetesAddon{
				Name:    api.DefaultReschedulerAddonName,
				Enabled: helpers.PointerToBool(true),
				Containers: []api.KubernetesContainerSpec{
					{
						Name:           api.DefaultReschedulerAddonName,
						CPURequests:    "50m",
						MemoryRequests: "150Mi",
						CPULimits:      "50m",
						MemoryLimits:   "150Mi",
					},
				},
			},
			[]string{
				"kubernetesReschedulerCPURequests",
				"kubernetesReschedulerCPULimit",
				"kubernetesReschedulerMemoryRequests",
				"kubernetesReschedulerMemoryLimit",
				"kubernetesReschedulerSpec",
			},
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			properties := api.GetK8sDefaultProperties(false)
			properties.OrchestratorProfile.KubernetesConfig.Addons = []api.KubernetesAddon{test.addon}
			parametersMap := paramsMap{}
			assignKubernetesParameters(properties, parametersMap, api.AzureCloudSpec, DefaultGeneratorCode)

			for _, expectedParam := range test.expectedParams {
				if !isKeyPresent(expectedParam, parametersMap) {
					t.Errorf("expected key %s to be present in the map", expectedParam)
				}
			}
		})
	}

}

func isKeyPresent(key string, paramMap map[string]interface{}) bool {
	_, ok := paramMap[key]
	return ok
}

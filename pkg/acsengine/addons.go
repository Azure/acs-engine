package acsengine

import (
	"fmt"
	"strings"

	"github.com/Azure/acs-engine/pkg/api"
	"github.com/Azure/acs-engine/pkg/api/common"
	"github.com/Azure/acs-engine/pkg/helpers"
)

type kubernetesFeatureSetting struct {
	sourceFile      string
	destinationFile string
	isEnabled       bool
}

func kubernetesAddonSettingsInit(profile *api.Properties) []kubernetesFeatureSetting {
	return []kubernetesFeatureSetting{
		{
			"kubernetesmasteraddons-heapster-deployment.yaml",
			"kube-heapster-deployment.yaml",
			true,
		},
		{
			"kubernetesmasteraddons-kube-dns-deployment.yaml",
			"kube-dns-deployment.yaml",
			true,
		},
		{
			"kubernetesmasteraddons-kube-proxy-daemonset.yaml",
			"kube-proxy-daemonset.yaml",
			true,
		},
		{
			"kubernetesmasteraddons-nvidia-device-plugin-daemonset.yaml",
			"nvidia-device-plugin.yaml",
			profile.IsNVIDIADevicePluginEnabled(),
		},
		{
			"kubernetesmasteraddons-kubernetes-dashboard-deployment.yaml",
			"kubernetes-dashboard-deployment.yaml",
			profile.OrchestratorProfile.KubernetesConfig.IsDashboardEnabled(),
		},
		{
			"kubernetesmasteraddons-unmanaged-azure-storage-classes.yaml",
			"azure-storage-classes.yaml",
			profile.AgentPoolProfiles[0].StorageProfile != api.ManagedDisks,
		},
		{
			"kubernetesmasteraddons-managed-azure-storage-classes.yaml",
			"azure-storage-classes.yaml",
			profile.AgentPoolProfiles[0].StorageProfile == api.ManagedDisks,
		},
		{
			"kubernetesmasteraddons-tiller-deployment.yaml",
			"kube-tiller-deployment.yaml",
			profile.OrchestratorProfile.KubernetesConfig.IsTillerEnabled(),
		},
		{
			"kubernetesmasteraddons-aci-connector-deployment.yaml",
			"aci-connector-deployment.yaml",
			profile.OrchestratorProfile.KubernetesConfig.IsACIConnectorEnabled(),
		},
		{
			"kubernetesmasteraddons-cluster-autoscaler-deployment.yaml",
			"cluster-autoscaler-deployment.yaml",
			profile.OrchestratorProfile.KubernetesConfig.IsClusterAutoscalerEnabled(),
		},
		{
			"kubernetesmasteraddons-kube-rescheduler-deployment.yaml",
			"kube-rescheduler-deployment.yaml",
			profile.OrchestratorProfile.KubernetesConfig.IsReschedulerEnabled(),
		},
		{
			"kubernetesmasteraddons-azure-npm-daemonset.yaml",
			"azure-npm-daemonset.yaml",
			profile.OrchestratorProfile.KubernetesConfig.NetworkPolicy == NetworkPolicyAzure && profile.OrchestratorProfile.KubernetesConfig.NetworkPlugin == NetworkPluginAzure,
		},
		{
			"kubernetesmasteraddons-calico-daemonset.yaml",
			"calico-daemonset.yaml",
			profile.OrchestratorProfile.KubernetesConfig.NetworkPolicy == NetworkPolicyCalico,
		},
		{
			"kubernetesmasteraddons-cilium-daemonset.yaml",
			"cilium-daemonset.yaml",
			profile.OrchestratorProfile.KubernetesConfig.NetworkPolicy == NetworkPolicyCilium,
		},
		{
			"kubernetesmasteraddons-flannel-daemonset.yaml",
			"flannel-daemonset.yaml",
			profile.OrchestratorProfile.KubernetesConfig.NetworkPlugin == NetworkPluginFlannel,
		},
		{
			"kubernetesmasteraddons-aad-default-admin-group-rbac.yaml",
			"aad-default-admin-group-rbac.yaml",
			profile.AADProfile != nil && profile.AADProfile.AdminGroupID != "",
		},
		{
			"kubernetesmasteraddons-azure-cloud-provider-deployment.yaml",
			"azure-cloud-provider-deployment.yaml",
			true,
		},
		{
			"kubernetesmasteraddons-metrics-server-deployment.yaml",
			"kube-metrics-server-deployment.yaml",
			profile.OrchestratorProfile.IsMetricsServerEnabled(),
		},
		{
			"kubernetesmasteraddons-omsagent-daemonset.yaml",
			"omsagent-daemonset.yaml",
			profile.OrchestratorProfile.IsContainerMonitoringEnabled(),
		},
		{
			"azure-cni-networkmonitor.yaml",
			"azure-cni-networkmonitor.yaml",
			profile.OrchestratorProfile.IsAzureCNI(),
		},
		{
			"kubernetesmaster-audit-policy.yaml",
			"audit-policy.yaml",
			common.IsKubernetesVersionGe(profile.OrchestratorProfile.OrchestratorVersion, "1.8.0"),
		},
		{
			"kubernetesmasteraddons-keyvault-flexvolume-installer.yaml",
			"keyvault-flexvolume-installer.yaml",
			profile.OrchestratorProfile.KubernetesConfig.IsKeyVaultFlexVolumeEnabled(),
		},
	}
}

func kubernetesManifestSettingsInit(profile *api.Properties) []kubernetesFeatureSetting {
	return []kubernetesFeatureSetting{
		{
			"kubernetesmaster-kube-scheduler.yaml",
			"kube-scheduler.yaml",
			true,
		},
		{
			"kubernetesmaster-kube-controller-manager.yaml",
			"kube-controller-manager.yaml",
			true,
		},
		{
			"kubernetesmaster-cloud-controller-manager.yaml",
			"cloud-controller-manager.yaml",
			profile.OrchestratorProfile.KubernetesConfig.UseCloudControllerManager != nil && *profile.OrchestratorProfile.KubernetesConfig.UseCloudControllerManager,
		},
		{
			"kubernetesmaster-pod-security-policy.yaml",
			"pod-security-policy.yaml",
			helpers.IsTrueBoolPointer(profile.OrchestratorProfile.KubernetesConfig.EnablePodSecurityPolicy),
		},
		{
			"kubernetesmaster-kube-apiserver.yaml",
			"kube-apiserver.yaml",
			true,
		},
		{
			"kubernetesmaster-kube-addon-manager.yaml",
			"kube-addon-manager.yaml",
			true,
		},
	}
}

func kubernetesArtifactSettingsInitMaster(profile *api.Properties) []kubernetesFeatureSetting {
	return []kubernetesFeatureSetting{
		{
			"kuberneteskubelet.service",
			"kubelet.service",
			true,
		},
		{
			"kubernetesazurekms.service",
			"kms.service",
			true,
		},
	}
}

func kubernetesArtifactSettingsInitAgent(profile *api.Properties) []kubernetesFeatureSetting {
	return []kubernetesFeatureSetting{
		{
			"kuberneteskubelet.service",
			"kubelet.service",
			true,
		},
	}
}

func substituteConfigString(input string, kubernetesFeatureSettings []kubernetesFeatureSetting, sourcePath string, destinationPath string, placeholder string, orchestratorVersion string) string {
	var config string

	versions := strings.Split(orchestratorVersion, ".")
	for _, setting := range kubernetesFeatureSettings {
		if setting.isEnabled {
			config += buildConfigString(
				setting.sourceFile,
				setting.destinationFile,
				sourcePath,
				destinationPath,
				versions[0]+"."+versions[1])
		}
	}

	return strings.Replace(input, placeholder, config, -1)
}

func buildConfigString(sourceFile string, destinationFile string, sourcePath string, destinationPath string, version string) string {
	sourceFileFullPath := sourcePath + "/" + sourceFile
	sourceFileFullPathVersioned := sourcePath + "/" + version + "/" + sourceFile

	// Test to check if the versioned file can be read.
	_, err := Asset(sourceFileFullPathVersioned)
	if err == nil {
		sourceFileFullPath = sourceFileFullPathVersioned
	}

	contents := []string{
		fmt.Sprintf("- path: %s/%s", destinationPath, destinationFile),
		"  permissions: \\\"0644\\\"",
		"  encoding: gzip",
		"  owner: \\\"root\\\"",
		"  content: !!binary |",
		fmt.Sprintf("    %s\\n\\n", getBase64CustomScript(sourceFileFullPath)),
	}

	return strings.Join(contents, "\\n")
}

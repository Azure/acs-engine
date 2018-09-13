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

type kubernetesAddonSetting struct {
	kubernetesFeatureSetting
	rawScript string
}

func kubernetesAddonSettingsInit(profile *api.Properties) []kubernetesAddonSetting {
	return []kubernetesAddonSetting{
		{

			kubernetesFeatureSetting{
				"kubernetesmasteraddons-heapster-deployment.yaml",
				"kube-heapster-deployment.yaml",
				true,
			},
			profile.OrchestratorProfile.KubernetesConfig.GetAddonScript(DefaultKubeHeapsterDeploymentAddonName),
		},
		{
			kubernetesFeatureSetting{
				"kubernetesmasteraddons-kube-dns-deployment.yaml",
				"kube-dns-deployment.yaml",
				true,
			},
			profile.OrchestratorProfile.KubernetesConfig.GetAddonScript(DefaultKubeDNSDeploymentAddonName),
		},
		{
			kubernetesFeatureSetting{
				"kubernetesmasteraddons-kube-proxy-daemonset.yaml",
				"kube-proxy-daemonset.yaml",
				true,
			},
			profile.OrchestratorProfile.KubernetesConfig.GetAddonScript(DefaultKubeProxyAddonName),
		},
		{
			kubernetesFeatureSetting{
				"kubernetesmasteraddons-nvidia-device-plugin-daemonset.yaml",
				"nvidia-device-plugin.yaml",
				profile.IsNVIDIADevicePluginEnabled(),
			},
			profile.OrchestratorProfile.KubernetesConfig.GetAddonScript(NVIDIADevicePluginAddonName),
		},
		{
			kubernetesFeatureSetting{
				"kubernetesmasteraddons-kubernetes-dashboard-deployment.yaml",
				"kubernetes-dashboard-deployment.yaml",
				profile.OrchestratorProfile.KubernetesConfig.IsDashboardEnabled(),
			},
			profile.OrchestratorProfile.KubernetesConfig.GetAddonScript(DefaultDashboardAddonName),
		},
		{
			kubernetesFeatureSetting{
				"kubernetesmasteraddons-unmanaged-azure-storage-classes.yaml",
				"azure-storage-classes.yaml",
				profile.AgentPoolProfiles[0].StorageProfile != api.ManagedDisks,
			},
			profile.OrchestratorProfile.KubernetesConfig.GetAddonScript(DefaultAzureStorageClassesAddonName),
		},
		{
			kubernetesFeatureSetting{
				"kubernetesmasteraddons-managed-azure-storage-classes.yaml",
				"azure-storage-classes.yaml",
				profile.AgentPoolProfiles[0].StorageProfile == api.ManagedDisks,
			},
			profile.OrchestratorProfile.KubernetesConfig.GetAddonScript(DefaultAzureStorageClassesAddonName),
		},
		{
			kubernetesFeatureSetting{
				"kubernetesmasteraddons-tiller-deployment.yaml",
				"kube-tiller-deployment.yaml",
				profile.OrchestratorProfile.KubernetesConfig.IsTillerEnabled(),
			},
			profile.OrchestratorProfile.KubernetesConfig.GetAddonScript(DefaultTillerAddonName),
		},
		{
			kubernetesFeatureSetting{
				"kubernetesmasteraddons-aad-pod-identity-deployment.yaml",
				"aad-pod-identity-deployment.yaml",
				profile.OrchestratorProfile.KubernetesConfig.IsAADPodIdentityEnabled(),
			},
			profile.OrchestratorProfile.KubernetesConfig.GetAddonScript(DefaultAADPodIdentityAddonName),
		},
		{
			kubernetesFeatureSetting{
				"kubernetesmasteraddons-aci-connector-deployment.yaml",
				"aci-connector-deployment.yaml",
				profile.OrchestratorProfile.KubernetesConfig.IsACIConnectorEnabled(),
			},
			profile.OrchestratorProfile.KubernetesConfig.GetAddonScript(DefaultACIConnectorAddonName),
		},
		{
			kubernetesFeatureSetting{
				"kubernetesmasteraddons-cluster-autoscaler-deployment.yaml",
				"cluster-autoscaler-deployment.yaml",
				profile.OrchestratorProfile.KubernetesConfig.IsClusterAutoscalerEnabled(),
			},
			profile.OrchestratorProfile.KubernetesConfig.GetAddonScript(DefaultClusterAutoscalerAddonName),
		},
		{
			kubernetesFeatureSetting{
				"kubernetesmasteraddons-kube-rescheduler-deployment.yaml",
				"kube-rescheduler-deployment.yaml",
				profile.OrchestratorProfile.KubernetesConfig.IsReschedulerEnabled(),
			},
			profile.OrchestratorProfile.KubernetesConfig.GetAddonScript(DefaultReschedulerAddonName),
		},
		{

			kubernetesFeatureSetting{
				"kubernetesmasteraddons-azure-npm-daemonset.yaml",
				"azure-npm-daemonset.yaml",
				profile.OrchestratorProfile.KubernetesConfig.NetworkPolicy == NetworkPolicyAzure && profile.OrchestratorProfile.KubernetesConfig.NetworkPlugin == NetworkPluginAzure,
			},
			profile.OrchestratorProfile.KubernetesConfig.GetAddonScript(DefaultAzureNpmDaemonSetAddonName),
		},
		{

			kubernetesFeatureSetting{
				"kubernetesmasteraddons-calico-daemonset.yaml",
				"calico-daemonset.yaml",
				profile.OrchestratorProfile.KubernetesConfig.NetworkPolicy == NetworkPolicyCalico,
			},
			profile.OrchestratorProfile.KubernetesConfig.GetAddonScript(DefaultCalicoDaemonSetAddonName),
		},
		{
			kubernetesFeatureSetting{
				"kubernetesmasteraddons-cilium-daemonset.yaml",
				"cilium-daemonset.yaml",
				profile.OrchestratorProfile.KubernetesConfig.NetworkPolicy == NetworkPolicyCilium,
			},
			profile.OrchestratorProfile.KubernetesConfig.GetAddonScript(DefaultCiliumDaemonSetAddonName),
		},
		{
			kubernetesFeatureSetting{
				"kubernetesmasteraddons-flannel-daemonset.yaml",
				"flannel-daemonset.yaml",
				profile.OrchestratorProfile.KubernetesConfig.NetworkPlugin == NetworkPluginFlannel,
			},
			profile.OrchestratorProfile.KubernetesConfig.GetAddonScript(DefaultFlannelDaemonSetAddonName),
		},
		{
			kubernetesFeatureSetting{
				"kubernetesmasteraddons-aad-default-admin-group-rbac.yaml",
				"aad-default-admin-group-rbac.yaml",
				profile.AADProfile != nil && profile.AADProfile.AdminGroupID != "",
			},
			profile.OrchestratorProfile.KubernetesConfig.GetAddonScript(DefaultAADAdminGroupRBACAddonName),
		},
		{
			kubernetesFeatureSetting{
				"kubernetesmasteraddons-azure-cloud-provider-deployment.yaml",
				"azure-cloud-provider-deployment.yaml",
				true,
			},
			profile.OrchestratorProfile.KubernetesConfig.GetAddonScript(DefaultAzureCloudProviderDeploymentAddonName),
		},
		{
			kubernetesFeatureSetting{
				"kubernetesmasteraddons-metrics-server-deployment.yaml",
				"kube-metrics-server-deployment.yaml",
				profile.OrchestratorProfile.IsMetricsServerEnabled(),
			},
			profile.OrchestratorProfile.KubernetesConfig.GetAddonScript(DefaultMetricsServerAddonName),
		},
		{
			kubernetesFeatureSetting{
				"kubernetesmasteraddons-omsagent-daemonset.yaml",
				"omsagent-daemonset.yaml",
				profile.OrchestratorProfile.KubernetesConfig.IsContainerMonitoringEnabled(),
			},
			profile.OrchestratorProfile.KubernetesConfig.GetAddonScript(ContainerMonitoringAddonName),
		},
		{
			kubernetesFeatureSetting{
				"azure-cni-networkmonitor.yaml",
				"azure-cni-networkmonitor.yaml",
				profile.OrchestratorProfile.IsAzureCNI(),
			},
			profile.OrchestratorProfile.KubernetesConfig.GetAddonScript(DefaultAzureCNINetworkMonitorAddonName),
		},
		{
			kubernetesFeatureSetting{
				"kubernetesmaster-audit-policy.yaml",
				"audit-policy.yaml",
				common.IsKubernetesVersionGe(profile.OrchestratorProfile.OrchestratorVersion, "1.8.0"),
			},
			profile.OrchestratorProfile.KubernetesConfig.GetAddonScript(DefaultAuditPolicyAddonName),
		},
		{
			kubernetesFeatureSetting{
				"kubernetesmasteraddons-blobfuse-flexvolume-installer.yaml",
				"blobfuse-flexvolume-installer.yaml",
				profile.OrchestratorProfile.KubernetesConfig.IsBlobfuseFlexVolumeEnabled(),
			},
			profile.OrchestratorProfile.KubernetesConfig.GetAddonScript(DefaultBlobfuseFlexVolumeAddonName),
		},
		{
			kubernetesFeatureSetting{
				"kubernetesmasteraddons-smb-flexvolume-installer.yaml",
				"smb-flexvolume-installer.yaml",
				profile.OrchestratorProfile.KubernetesConfig.IsSMBFlexVolumeEnabled(),
			},
			profile.OrchestratorProfile.KubernetesConfig.GetAddonScript(DefaultSMBFlexVolumeAddonName),
		},
		{
			kubernetesFeatureSetting{
				"kubernetesmasteraddons-keyvault-flexvolume-installer.yaml",
				"keyvault-flexvolume-installer.yaml",
				profile.OrchestratorProfile.KubernetesConfig.IsKeyVaultFlexVolumeEnabled(),
			},
			profile.OrchestratorProfile.KubernetesConfig.GetAddonScript(DefaultKeyVaultFlexVolumeAddonName),
		},
		{
			kubernetesFeatureSetting{
				"kubernetesmasteraddons-elb-svc.yaml",
				"elb-svc.yaml",
				profile.OrchestratorProfile.KubernetesConfig.LoadBalancerSku == "Standard",
			},
			profile.OrchestratorProfile.KubernetesConfig.GetAddonScript(DefaultELBSVCAddonName),
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
				getCustomScriptFromFile(setting.sourceFile,
					sourcePath,
					versions[0]+"."+versions[1]),
				setting.destinationFile,
				destinationPath)
		}
	}

	return strings.Replace(input, placeholder, config, -1)
}

func substituteAddonConfigString(input string, kubernetesAddonSettings []kubernetesAddonSetting, sourcePath string, destinationPath string, placeholder string, orchestratorVersion string) string {
	var config string

	versions := strings.Split(orchestratorVersion, ".")
	for _, setting := range kubernetesAddonSettings {
		if setting.isEnabled {
			var cscript string
			if setting.rawScript != "" {
				cscript = setting.rawScript
			} else {
				cscript = getCustomScriptFromFile(setting.sourceFile,
					sourcePath,
					versions[0]+"."+versions[1])
			}
			config += buildConfigString(
				cscript,
				setting.destinationFile,
				destinationPath)
		}
	}

	return strings.Replace(input, placeholder, config, -1)
}

func buildConfigString(configString, destinationFile, destinationPath string) string {

	contents := []string{
		fmt.Sprintf("- path: %s/%s", destinationPath, destinationFile),
		"  permissions: \\\"0644\\\"",
		"  encoding: gzip",
		"  owner: \\\"root\\\"",
		"  content: !!binary |",
		fmt.Sprintf("    %s\\n\\n", configString),
	}

	return strings.Join(contents, "\\n")
}

func getCustomScriptFromFile(sourceFile, sourcePath, version string) string {
	sourceFileFullPath := sourcePath + "/" + sourceFile
	sourceFileFullPathVersioned := sourcePath + "/" + version + "/" + sourceFile

	// Test to check if the versioned file can be read.
	_, err := Asset(sourceFileFullPathVersioned)
	if err == nil {
		sourceFileFullPath = sourceFileFullPathVersioned
	}
	return getBase64CustomScript(sourceFileFullPath)
}

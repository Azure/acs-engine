package acsengine

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Azure/acs-engine/pkg/api"
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
			"kubernetesmasteraddons-kube-rescheduler-deployment.yaml",
			"kube-rescheduler-deployment.yaml",
			profile.OrchestratorProfile.KubernetesConfig.IsReschedulerEnabled(),
		},
		{
			"kubernetesmasteraddons-calico-daemonset.yaml",
			"calico-daemonset.yaml",
			profile.OrchestratorProfile.KubernetesConfig.NetworkPolicy == "calico",
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

func kubernetesArtifactSettingsInit(profile *api.Properties) []kubernetesFeatureSetting {
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
	sourceFileNoSuffix := strings.TrimSuffix(sourceFile, filepath.Ext(sourceFile))
	sourceFileVersioned := sourcePath + "/" + sourceFile
	sourcePathVersioned := "parts/" + sourcePath + "/" + version

	filepath.Walk(sourcePathVersioned, func(path string, info os.FileInfo, err error) error {
		if err == nil {
			if strings.HasPrefix(info.Name(), sourceFileNoSuffix) {
				sourceFileVersioned = sourcePath + "/" + version + "/" + info.Name()
			}
			return nil
		}
		return err
	})

	contents := []string{
		fmt.Sprintf("- path: %s/%s", destinationPath, destinationFile),
		"  permissions: \\\"0644\\\"",
		"  encoding: gzip",
		"  owner: \\\"root\\\"",
		"  content: !!binary |",
		fmt.Sprintf("    %s\\n\\n", getBase64CustomScript(sourceFileVersioned)),
	}

	return strings.Join(contents, "\\n")
}

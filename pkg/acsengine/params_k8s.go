package acsengine

import (
	"strconv"
	"strings"

	"github.com/Azure/acs-engine/pkg/api"
	"github.com/Azure/acs-engine/pkg/helpers"
)

func assignKubernetesParameters(properties *api.Properties, parametersMap paramsMap,
	cloudSpecConfig AzureEnvironmentSpecConfig, generatorCode string) {
	addValue(parametersMap, "generatorCode", generatorCode)
	if properties.OrchestratorProfile.IsKubernetes() ||
		properties.OrchestratorProfile.IsOpenShift() {
		k8sVersion := properties.OrchestratorProfile.OrchestratorVersion

		dockerEngineVersion := KubeConfigs[k8sVersion]["dockerEngineVersion"]

		if properties.OrchestratorProfile.KubernetesConfig != nil {
			if helpers.IsTrueBoolPointer(properties.OrchestratorProfile.KubernetesConfig.UseCloudControllerManager) {
				kubernetesCcmSpec := properties.OrchestratorProfile.KubernetesConfig.KubernetesImageBase + KubeConfigs[k8sVersion]["ccm"]
				if properties.OrchestratorProfile.KubernetesConfig.CustomCcmImage != "" {
					kubernetesCcmSpec = properties.OrchestratorProfile.KubernetesConfig.CustomCcmImage
				}

				addValue(parametersMap, "kubernetesCcmImageSpec", kubernetesCcmSpec)
			}

			kubernetesHyperkubeSpec := properties.OrchestratorProfile.KubernetesConfig.KubernetesImageBase + KubeConfigs[k8sVersion]["hyperkube"]
			if properties.OrchestratorProfile.KubernetesConfig.CustomHyperkubeImage != "" {
				kubernetesHyperkubeSpec = properties.OrchestratorProfile.KubernetesConfig.CustomHyperkubeImage
			}

			addValue(parametersMap, "kubeDNSServiceIP", properties.OrchestratorProfile.KubernetesConfig.DNSServiceIP)
			addValue(parametersMap, "kubeServiceCidr", properties.OrchestratorProfile.KubernetesConfig.ServiceCIDR)
			addValue(parametersMap, "kubernetesHyperkubeSpec", kubernetesHyperkubeSpec)
			addValue(parametersMap, "kubernetesAddonManagerSpec", cloudSpecConfig.KubernetesSpecConfig.KubernetesImageBase+KubeConfigs[k8sVersion]["addonmanager"])
			addValue(parametersMap, "kubernetesAddonResizerSpec", cloudSpecConfig.KubernetesSpecConfig.KubernetesImageBase+KubeConfigs[k8sVersion]["addonresizer"])
			addValue(parametersMap, "kubernetesDNSMasqSpec", cloudSpecConfig.KubernetesSpecConfig.KubernetesImageBase+KubeConfigs[k8sVersion]["dnsmasq"])
			addValue(parametersMap, "kubernetesExecHealthzSpec", cloudSpecConfig.KubernetesSpecConfig.KubernetesImageBase+KubeConfigs[k8sVersion]["exechealthz"])
			addValue(parametersMap, "kubernetesHeapsterSpec", cloudSpecConfig.KubernetesSpecConfig.KubernetesImageBase+KubeConfigs[k8sVersion]["heapster"])
			tillerAddon := getAddonByName(properties.OrchestratorProfile.KubernetesConfig.Addons, DefaultTillerAddonName)
			c := getAddonContainersIndexByName(tillerAddon.Containers, DefaultTillerAddonName)
			if c > -1 {
				addValue(parametersMap, "kubernetesTillerCPURequests", tillerAddon.Containers[c].CPURequests)
				addValue(parametersMap, "kubernetesTillerCPULimit", tillerAddon.Containers[c].CPULimits)
				addValue(parametersMap, "kubernetesTillerMemoryRequests", tillerAddon.Containers[c].MemoryRequests)
				addValue(parametersMap, "kubernetesTillerMemoryLimit", tillerAddon.Containers[c].MemoryLimits)
				addValue(parametersMap, "kubernetesTillerMaxHistory", tillerAddon.Config["max-history"])
				if tillerAddon.Containers[c].Image != "" {
					addValue(parametersMap, "kubernetesTillerSpec", tillerAddon.Containers[c].Image)
				} else {
					addValue(parametersMap, "kubernetesTillerSpec", cloudSpecConfig.KubernetesSpecConfig.TillerImageBase+KubeConfigs[k8sVersion][DefaultTillerAddonName])
				}
			}
			aadPodIdentityAddon := getAddonByName(properties.OrchestratorProfile.KubernetesConfig.Addons, DefaultAADPodIdentityAddonName)
			c = getAddonContainersIndexByName(aadPodIdentityAddon.Containers, DefaultAADPodIdentityAddonName)
			if c > -1 {
				addValue(parametersMap, "kubernetesAADPodIdentityEnabled", helpers.IsTrueBoolPointer(aadPodIdentityAddon.Enabled))
			}

			aciConnectorAddon := getAddonByName(properties.OrchestratorProfile.KubernetesConfig.Addons, DefaultACIConnectorAddonName)
			c = getAddonContainersIndexByName(aciConnectorAddon.Containers, DefaultACIConnectorAddonName)
			if c > -1 {
				addValue(parametersMap, "kubernetesACIConnectorEnabled", helpers.IsTrueBoolPointer(aciConnectorAddon.Enabled))
				addValue(parametersMap, "kubernetesACIConnectorNodeName", aciConnectorAddon.Config["nodeName"])
				addValue(parametersMap, "kubernetesACIConnectorOS", aciConnectorAddon.Config["os"])
				addValue(parametersMap, "kubernetesACIConnectorTaint", aciConnectorAddon.Config["taint"])
				addValue(parametersMap, "kubernetesACIConnectorRegion", aciConnectorAddon.Config["region"])
				addValue(parametersMap, "kubernetesACIConnectorCPURequests", aciConnectorAddon.Containers[c].CPURequests)
				addValue(parametersMap, "kubernetesACIConnectorCPULimit", aciConnectorAddon.Containers[c].CPULimits)
				addValue(parametersMap, "kubernetesACIConnectorMemoryRequests", aciConnectorAddon.Containers[c].MemoryRequests)
				addValue(parametersMap, "kubernetesACIConnectorMemoryLimit", aciConnectorAddon.Containers[c].MemoryLimits)
				if aciConnectorAddon.Containers[c].Image != "" {
					addValue(parametersMap, "kubernetesACIConnectorSpec", aciConnectorAddon.Containers[c].Image)
				} else {
					addValue(parametersMap, "kubernetesACIConnectorSpec", cloudSpecConfig.KubernetesSpecConfig.ACIConnectorImageBase+KubeConfigs[k8sVersion][DefaultACIConnectorAddonName])
				}
			}
			clusterAutoscalerAddon := getAddonByName(properties.OrchestratorProfile.KubernetesConfig.Addons, DefaultClusterAutoscalerAddonName)
			c = getAddonContainersIndexByName(clusterAutoscalerAddon.Containers, DefaultClusterAutoscalerAddonName)
			if c > -1 {
				addValue(parametersMap, "kubernetesClusterAutoscalerAzureCloud", cloudSpecConfig.CloudName)
				addValue(parametersMap, "kubernetesClusterAutoscalerCPURequests", clusterAutoscalerAddon.Containers[c].CPURequests)
				addValue(parametersMap, "kubernetesClusterAutoscalerCPULimit", clusterAutoscalerAddon.Containers[c].CPULimits)
				addValue(parametersMap, "kubernetesClusterAutoscalerMemoryRequests", clusterAutoscalerAddon.Containers[c].MemoryRequests)
				addValue(parametersMap, "kubernetesClusterAutoscalerMemoryLimit", clusterAutoscalerAddon.Containers[c].MemoryLimits)
				addValue(parametersMap, "kubernetesClusterAutoscalerMinNodes", clusterAutoscalerAddon.Config["minNodes"])
				addValue(parametersMap, "kubernetesClusterAutoscalerMaxNodes", clusterAutoscalerAddon.Config["maxNodes"])
				addValue(parametersMap, "kubernetesClusterAutoscalerEnabled", helpers.IsTrueBoolPointer(clusterAutoscalerAddon.Enabled))
				addValue(parametersMap, "kubernetesClusterAutoscalerUseManagedIdentity", strings.ToLower(strconv.FormatBool(properties.OrchestratorProfile.KubernetesConfig.UseManagedIdentity)))
				if clusterAutoscalerAddon.Containers[c].Image != "" {
					addValue(parametersMap, "kubernetesClusterAutoscalerSpec", clusterAutoscalerAddon.Containers[c].Image)
				} else {
					addValue(parametersMap, "kubernetesClusterAutoscalerSpec", cloudSpecConfig.KubernetesSpecConfig.KubernetesImageBase+KubeConfigs[k8sVersion][DefaultClusterAutoscalerAddonName])
				}
			}
			kvFlexVolumeInstallerAddon := getAddonByName(properties.OrchestratorProfile.KubernetesConfig.Addons, DefaultKeyVaultFlexVolumeAddonName)
			c = getAddonContainersIndexByName(kvFlexVolumeInstallerAddon.Containers, DefaultKeyVaultFlexVolumeAddonName)
			if c > -1 {
				addValue(parametersMap, "kubernetesKeyVaultFlexVolumeInstallerCPURequests", kvFlexVolumeInstallerAddon.Containers[c].CPURequests)
				addValue(parametersMap, "kubernetesKeyVaultFlexVolumeInstallerCPULimit", kvFlexVolumeInstallerAddon.Containers[c].CPULimits)
				addValue(parametersMap, "kubernetesKeyVaultFlexVolumeInstallerMemoryRequests", kvFlexVolumeInstallerAddon.Containers[c].MemoryRequests)
				addValue(parametersMap, "kubernetesKeyVaultFlexVolumeInstallerMemoryLimit", kvFlexVolumeInstallerAddon.Containers[c].MemoryLimits)
			}
			dashboardAddon := getAddonByName(properties.OrchestratorProfile.KubernetesConfig.Addons, DefaultDashboardAddonName)
			c = getAddonContainersIndexByName(dashboardAddon.Containers, DefaultDashboardAddonName)
			if c > -1 {
				addValue(parametersMap, "kubernetesDashboardCPURequests", dashboardAddon.Containers[c].CPURequests)
				addValue(parametersMap, "kubernetesDashboardCPULimit", dashboardAddon.Containers[c].CPULimits)
				addValue(parametersMap, "kubernetesDashboardMemoryRequests", dashboardAddon.Containers[c].MemoryRequests)
				addValue(parametersMap, "kubernetesDashboardMemoryLimit", dashboardAddon.Containers[c].MemoryLimits)
				if dashboardAddon.Containers[c].Image != "" {
					addValue(parametersMap, "kubernetesDashboardSpec", dashboardAddon.Containers[c].Image)
				} else {
					addValue(parametersMap, "kubernetesDashboardSpec", cloudSpecConfig.KubernetesSpecConfig.KubernetesImageBase+KubeConfigs[k8sVersion][DefaultDashboardAddonName])
				}
			}
			reschedulerAddon := getAddonByName(properties.OrchestratorProfile.KubernetesConfig.Addons, DefaultReschedulerAddonName)
			c = getAddonContainersIndexByName(reschedulerAddon.Containers, DefaultReschedulerAddonName)
			if c > -1 {
				addValue(parametersMap, "kubernetesReschedulerCPURequests", reschedulerAddon.Containers[c].CPURequests)
				addValue(parametersMap, "kubernetesReschedulerCPULimit", reschedulerAddon.Containers[c].CPULimits)
				addValue(parametersMap, "kubernetesReschedulerMemoryRequests", reschedulerAddon.Containers[c].MemoryRequests)
				addValue(parametersMap, "kubernetesReschedulerMemoryLimit", reschedulerAddon.Containers[c].MemoryLimits)
				if reschedulerAddon.Containers[c].Image != "" {
					addValue(parametersMap, "kubernetesReschedulerSpec", dashboardAddon.Containers[c].Image)
				} else {
					addValue(parametersMap, "kubernetesReschedulerSpec", cloudSpecConfig.KubernetesSpecConfig.KubernetesImageBase+KubeConfigs[k8sVersion][DefaultReschedulerAddonName])
				}
			}
			metricsServerAddon := getAddonByName(properties.OrchestratorProfile.KubernetesConfig.Addons, DefaultMetricsServerAddonName)
			c = getAddonContainersIndexByName(metricsServerAddon.Containers, DefaultMetricsServerAddonName)
			if c > -1 {
				if metricsServerAddon.Containers[c].Image != "" {
					addValue(parametersMap, "kubernetesMetricsServerSpec", metricsServerAddon.Containers[c].Image)
				} else {
					addValue(parametersMap, "kubernetesMetricsServerSpec", cloudSpecConfig.KubernetesSpecConfig.KubernetesImageBase+KubeConfigs[k8sVersion][DefaultMetricsServerAddonName])
				}
			}
			nvidiaDevicePluginAddon := getAddonByName(properties.OrchestratorProfile.KubernetesConfig.Addons, NVIDIADevicePluginAddonName)
			c = getAddonContainersIndexByName(nvidiaDevicePluginAddon.Containers, NVIDIADevicePluginAddonName)
			if c > -1 {
				addValue(parametersMap, "kubernetesNVIDIADevicePluginCPURequests", nvidiaDevicePluginAddon.Containers[c].CPURequests)
				addValue(parametersMap, "kubernetesNVIDIADevicePluginCPULimit", nvidiaDevicePluginAddon.Containers[c].CPULimits)
				addValue(parametersMap, "kubernetesNVIDIADevicePluginMemoryRequests", nvidiaDevicePluginAddon.Containers[c].MemoryRequests)
				addValue(parametersMap, "kubernetesNVIDIADevicePluginMemoryLimit", nvidiaDevicePluginAddon.Containers[c].MemoryLimits)
				if nvidiaDevicePluginAddon.Containers[c].Image != "" {
					addValue(parametersMap, "kubernetesNVIDIADevicePluginSpec", nvidiaDevicePluginAddon.Containers[c].Image)
				} else {
					addValue(parametersMap, "kubernetesNVIDIADevicePluginSpec", cloudSpecConfig.KubernetesSpecConfig.NVIDIAImageBase+KubeConfigs[k8sVersion][NVIDIADevicePluginAddonName])
				}
			}
			containerMonitoringAddon := getAddonByName(properties.OrchestratorProfile.KubernetesConfig.Addons, ContainerMonitoringAddonName)
			c = getAddonContainersIndexByName(containerMonitoringAddon.Containers, "omsagent")
			if c > -1 {
				addValue(parametersMap, "omsAgentVersion", containerMonitoringAddon.Config["omsAgentVersion"])
				addValue(parametersMap, "omsAgentDockerProviderVersion", containerMonitoringAddon.Config["dockerProviderVersion"])
				addValue(parametersMap, "omsAgentWorkspaceGuid", containerMonitoringAddon.Config["workspaceGuid"])
				addValue(parametersMap, "omsAgentWorkspaceKey", containerMonitoringAddon.Config["workspaceKey"])
				addValue(parametersMap, "kubernetesOMSAgentCPURequests", containerMonitoringAddon.Containers[c].CPURequests)
				addValue(parametersMap, "kubernetesOMSAgentCPULimit", containerMonitoringAddon.Containers[c].CPULimits)
				addValue(parametersMap, "kubernetesOMSAgentMemoryRequests", containerMonitoringAddon.Containers[c].MemoryRequests)
				addValue(parametersMap, "kubernetesOMSAgentMemoryLimit", containerMonitoringAddon.Containers[c].MemoryLimits)
				if containerMonitoringAddon.Containers[c].Image != "" {
					addValue(parametersMap, "omsAgentImage", containerMonitoringAddon.Containers[c].Image)
				} else {
					addValue(parametersMap, "omsAgentImage", cloudSpecConfig.KubernetesSpecConfig.KubernetesImageBase+KubeConfigs[k8sVersion][ContainerMonitoringAddonName])
				}
			}
			if properties.OrchestratorProfile.IsAzureCNI() {
				azureCNINetworkmonitorAddon := getAddonByName(properties.OrchestratorProfile.KubernetesConfig.Addons, AzureCNINetworkMonitoringAddonName)
				c = getAddonContainersIndexByName(azureCNINetworkmonitorAddon.Containers, AzureCNINetworkMonitoringAddonName)
				if c > -1 {
					if azureCNINetworkmonitorAddon.Containers[c].Image != "" {
						addValue(parametersMap, "AzureCNINetworkMonitorImageURL", azureCNINetworkmonitorAddon.Containers[c].Image)
					} else {
						addValue(parametersMap, "AzureCNINetworkMonitorImageURL", cloudSpecConfig.KubernetesSpecConfig.AzureCNIImageBase+KubeConfigs[k8sVersion][AzureCNINetworkMonitoringAddonName])
					}
				}
			}
			addValue(parametersMap, "kubernetesKubeDNSSpec", cloudSpecConfig.KubernetesSpecConfig.KubernetesImageBase+KubeConfigs[k8sVersion]["dns"])
			addValue(parametersMap, "kubernetesPodInfraContainerSpec", cloudSpecConfig.KubernetesSpecConfig.KubernetesImageBase+KubeConfigs[k8sVersion]["pause"])
			addValue(parametersMap, "cloudProviderBackoff", strconv.FormatBool(properties.OrchestratorProfile.KubernetesConfig.CloudProviderBackoff))
			addValue(parametersMap, "cloudProviderBackoffRetries", strconv.Itoa(properties.OrchestratorProfile.KubernetesConfig.CloudProviderBackoffRetries))
			addValue(parametersMap, "cloudProviderBackoffExponent", strconv.FormatFloat(properties.OrchestratorProfile.KubernetesConfig.CloudProviderBackoffExponent, 'f', -1, 64))
			addValue(parametersMap, "cloudProviderBackoffDuration", strconv.Itoa(properties.OrchestratorProfile.KubernetesConfig.CloudProviderBackoffDuration))
			addValue(parametersMap, "cloudProviderBackoffJitter", strconv.FormatFloat(properties.OrchestratorProfile.KubernetesConfig.CloudProviderBackoffJitter, 'f', -1, 64))
			addValue(parametersMap, "cloudProviderRatelimit", strconv.FormatBool(properties.OrchestratorProfile.KubernetesConfig.CloudProviderRateLimit))
			addValue(parametersMap, "cloudProviderRatelimitQPS", strconv.FormatFloat(properties.OrchestratorProfile.KubernetesConfig.CloudProviderRateLimitQPS, 'f', -1, 64))
			addValue(parametersMap, "cloudProviderRatelimitBucket", strconv.Itoa(properties.OrchestratorProfile.KubernetesConfig.CloudProviderRateLimitBucket))
			addValue(parametersMap, "kubeClusterCidr", properties.OrchestratorProfile.KubernetesConfig.ClusterSubnet)
			addValue(parametersMap, "kubernetesNonMasqueradeCidr", properties.OrchestratorProfile.KubernetesConfig.KubeletConfig["--non-masquerade-cidr"])
			addValue(parametersMap, "kubernetesKubeletClusterDomain", properties.OrchestratorProfile.KubernetesConfig.KubeletConfig["--cluster-domain"])
			addValue(parametersMap, "dockerBridgeCidr", properties.OrchestratorProfile.KubernetesConfig.DockerBridgeSubnet)
			addValue(parametersMap, "networkPolicy", properties.OrchestratorProfile.KubernetesConfig.NetworkPolicy)
			addValue(parametersMap, "networkPlugin", properties.OrchestratorProfile.KubernetesConfig.NetworkPlugin)
			addValue(parametersMap, "containerRuntime", properties.OrchestratorProfile.KubernetesConfig.ContainerRuntime)
			addValue(parametersMap, "cniPluginsURL", cloudSpecConfig.KubernetesSpecConfig.CNIPluginsDownloadURL)
			addValue(parametersMap, "vnetCniLinuxPluginsURL", cloudSpecConfig.KubernetesSpecConfig.VnetCNILinuxPluginsDownloadURL)
			addValue(parametersMap, "vnetCniWindowsPluginsURL", cloudSpecConfig.KubernetesSpecConfig.VnetCNIWindowsPluginsDownloadURL)
			addValue(parametersMap, "gchighthreshold", properties.OrchestratorProfile.KubernetesConfig.GCHighThreshold)
			addValue(parametersMap, "gclowthreshold", properties.OrchestratorProfile.KubernetesConfig.GCLowThreshold)
			addValue(parametersMap, "etcdDownloadURLBase", cloudSpecConfig.KubernetesSpecConfig.EtcdDownloadURLBase)
			addValue(parametersMap, "etcdVersion", properties.OrchestratorProfile.KubernetesConfig.EtcdVersion)
			addValue(parametersMap, "etcdDiskSizeGB", properties.OrchestratorProfile.KubernetesConfig.EtcdDiskSizeGB)
			addValue(parametersMap, "etcdEncryptionKey", properties.OrchestratorProfile.KubernetesConfig.EtcdEncryptionKey)
			if properties.OrchestratorProfile.KubernetesConfig.PrivateJumpboxProvision() {
				addValue(parametersMap, "jumpboxVMName", properties.OrchestratorProfile.KubernetesConfig.PrivateCluster.JumpboxProfile.Name)
				addValue(parametersMap, "jumpboxVMSize", properties.OrchestratorProfile.KubernetesConfig.PrivateCluster.JumpboxProfile.VMSize)
				addValue(parametersMap, "jumpboxUsername", properties.OrchestratorProfile.KubernetesConfig.PrivateCluster.JumpboxProfile.Username)
				addValue(parametersMap, "jumpboxOSDiskSizeGB", properties.OrchestratorProfile.KubernetesConfig.PrivateCluster.JumpboxProfile.OSDiskSizeGB)
				addValue(parametersMap, "jumpboxPublicKey", properties.OrchestratorProfile.KubernetesConfig.PrivateCluster.JumpboxProfile.PublicKey)
				addValue(parametersMap, "jumpboxStorageProfile", properties.OrchestratorProfile.KubernetesConfig.PrivateCluster.JumpboxProfile.StorageProfile)
			}

			if properties.OrchestratorProfile.KubernetesConfig.DockerEngineVersion != "" {
				dockerEngineVersion = properties.OrchestratorProfile.KubernetesConfig.DockerEngineVersion
			}
		}

		if properties.OrchestratorProfile.KubernetesConfig == nil ||
			!properties.OrchestratorProfile.KubernetesConfig.UseManagedIdentity {

			addValue(parametersMap, "servicePrincipalClientId", properties.ServicePrincipalProfile.ClientID)
			if properties.ServicePrincipalProfile.KeyvaultSecretRef != nil {
				addKeyvaultReference(parametersMap, "servicePrincipalClientSecret",
					properties.ServicePrincipalProfile.KeyvaultSecretRef.VaultID,
					properties.ServicePrincipalProfile.KeyvaultSecretRef.SecretName,
					properties.ServicePrincipalProfile.KeyvaultSecretRef.SecretVersion)
			} else {
				addValue(parametersMap, "servicePrincipalClientSecret", properties.ServicePrincipalProfile.Secret)
			}

			if properties.OrchestratorProfile.KubernetesConfig != nil && helpers.IsTrueBoolPointer(properties.OrchestratorProfile.KubernetesConfig.EnableEncryptionWithExternalKms) && !properties.OrchestratorProfile.KubernetesConfig.UseManagedIdentity && properties.ServicePrincipalProfile.ObjectID != "" {
				addValue(parametersMap, "servicePrincipalObjectId", properties.ServicePrincipalProfile.ObjectID)
			}
		}

		if properties.HostedMasterProfile != nil {
			addValue(parametersMap, "orchestratorName", "aks")
		} else if properties.OrchestratorProfile.IsOpenShift() {
			addValue(parametersMap, "orchestratorName", DefaultOpenshiftOrchestratorName)
		} else {
			addValue(parametersMap, "orchestratorName", DefaultOrchestratorName)
		}

		/**
		 The following parameters could be either a plain text, or referenced to a secret in a keyvault:
		 - apiServerCertificate
		 - apiServerPrivateKey
		 - caCertificate
		 - clientCertificate
		 - clientPrivateKey
		 - kubeConfigCertificate
		 - kubeConfigPrivateKey
		 - servicePrincipalClientSecret
		 - etcdClientCertificate
		 - etcdClientPrivateKey
		 - etcdServerCertificate
		 - etcdServerPrivateKey
		 - etcdPeerCertificates
		 - etcdPeerPrivateKeys

		 To refer to a keyvault secret, the value of the parameter in the api model file should be formatted as:

		 "<PARAMETER>": "/subscriptions/<SUB_ID>/resourceGroups/<RG_NAME>/providers/Microsoft.KeyVault/vaults/<KV_NAME>/secrets/<NAME>[/<VERSION>]"
		 where:
		   <SUB_ID> is the subscription ID of the keyvault
		   <RG_NAME> is the resource group of the keyvault
		   <KV_NAME> is the name of the keyvault
		   <NAME> is the name of the secret.
		   <VERSION> (optional) is the version of the secret (default: the latest version)

		 This will generate a reference block in the parameters file:

		 "reference": {
		   "keyVault": {
		     "id": "/subscriptions/<SUB_ID>/resourceGroups/<RG_NAME>/providers/Microsoft.KeyVault/vaults/<KV_NAME>"
		   },
		   "secretName": "<NAME>"
		   "secretVersion": "<VERSION>"
		}
		**/

		if properties.CertificateProfile != nil {
			addSecret(parametersMap, "apiServerCertificate", properties.CertificateProfile.APIServerCertificate, true)
			addSecret(parametersMap, "apiServerPrivateKey", properties.CertificateProfile.APIServerPrivateKey, true)
			addSecret(parametersMap, "caCertificate", properties.CertificateProfile.CaCertificate, true)
			addSecret(parametersMap, "caPrivateKey", properties.CertificateProfile.CaPrivateKey, true)
			addSecret(parametersMap, "clientCertificate", properties.CertificateProfile.ClientCertificate, true)
			addSecret(parametersMap, "clientPrivateKey", properties.CertificateProfile.ClientPrivateKey, true)
			addSecret(parametersMap, "kubeConfigCertificate", properties.CertificateProfile.KubeConfigCertificate, true)
			addSecret(parametersMap, "kubeConfigPrivateKey", properties.CertificateProfile.KubeConfigPrivateKey, true)
			if properties.MasterProfile != nil {
				addSecret(parametersMap, "etcdServerCertificate", properties.CertificateProfile.EtcdServerCertificate, true)
				addSecret(parametersMap, "etcdServerPrivateKey", properties.CertificateProfile.EtcdServerPrivateKey, true)
				addSecret(parametersMap, "etcdClientCertificate", properties.CertificateProfile.EtcdClientCertificate, true)
				addSecret(parametersMap, "etcdClientPrivateKey", properties.CertificateProfile.EtcdClientPrivateKey, true)
				for i, pc := range properties.CertificateProfile.EtcdPeerCertificates {
					addSecret(parametersMap, "etcdPeerCertificate"+strconv.Itoa(i), pc, true)
				}
				for i, pk := range properties.CertificateProfile.EtcdPeerPrivateKeys {
					addSecret(parametersMap, "etcdPeerPrivateKey"+strconv.Itoa(i), pk, true)
				}
			}
		}

		if properties.HostedMasterProfile != nil && properties.HostedMasterProfile.FQDN != "" {
			addValue(parametersMap, "kubernetesEndpoint", properties.HostedMasterProfile.FQDN)
		}

		if !properties.OrchestratorProfile.IsOpenShift() {
			addValue(parametersMap, "dockerEngineDownloadRepo", cloudSpecConfig.DockerSpecConfig.DockerEngineRepo)
			addValue(parametersMap, "dockerEngineVersion", dockerEngineVersion)
		}

		if properties.AADProfile != nil {
			addValue(parametersMap, "aadTenantId", properties.AADProfile.TenantID)
			if properties.AADProfile.AdminGroupID != "" {
				addValue(parametersMap, "aadAdminGroupId", properties.AADProfile.AdminGroupID)
			}
		}
	}
}

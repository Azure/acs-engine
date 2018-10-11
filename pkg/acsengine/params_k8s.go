package acsengine

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/Azure/acs-engine/pkg/api"
	"github.com/Azure/acs-engine/pkg/helpers"
)

func assignKubernetesParameters(properties *api.Properties, parametersMap paramsMap,
	cloudSpecConfig api.AzureEnvironmentSpecConfig, generatorCode string) {
	addValue(parametersMap, "generatorCode", generatorCode)

	orchestratorProfile := properties.OrchestratorProfile

	if orchestratorProfile.IsKubernetes() ||
		orchestratorProfile.IsOpenShift() {
		k8sComponents := api.K8sComponentsByVersionMap[orchestratorProfile.OrchestratorVersion]

		dockerEngineVersion := k8sComponents["dockerEngineVersion"]

		kubernetesConfig := orchestratorProfile.KubernetesConfig

		if kubernetesConfig != nil {
			if helpers.IsTrueBoolPointer(kubernetesConfig.UseCloudControllerManager) {
				kubernetesCcmSpec := kubernetesConfig.KubernetesImageBase + k8sComponents["ccm"]
				if kubernetesConfig.CustomCcmImage != "" {
					kubernetesCcmSpec = kubernetesConfig.CustomCcmImage
				}

				addValue(parametersMap, "kubernetesCcmImageSpec", kubernetesCcmSpec)
			}

			kubernetesHyperkubeSpec := kubernetesConfig.KubernetesImageBase + k8sComponents["hyperkube"]
			if kubernetesConfig.CustomHyperkubeImage != "" {
				kubernetesHyperkubeSpec = kubernetesConfig.CustomHyperkubeImage
			}

			addValue(parametersMap, "kubeDNSServiceIP", kubernetesConfig.DNSServiceIP)
			addValue(parametersMap, "kubeServiceCidr", kubernetesConfig.ServiceCIDR)
			addValue(parametersMap, "kubernetesHyperkubeSpec", kubernetesHyperkubeSpec)
			addValue(parametersMap, "kubernetesAddonManagerSpec", cloudSpecConfig.KubernetesSpecConfig.KubernetesImageBase+k8sComponents["addonmanager"])
			addValue(parametersMap, "kubernetesAddonResizerSpec", cloudSpecConfig.KubernetesSpecConfig.KubernetesImageBase+k8sComponents["addonresizer"])
			addValue(parametersMap, "kubernetesDNSMasqSpec", cloudSpecConfig.KubernetesSpecConfig.KubernetesImageBase+k8sComponents["dnsmasq"])
			addValue(parametersMap, "kubernetesExecHealthzSpec", cloudSpecConfig.KubernetesSpecConfig.KubernetesImageBase+k8sComponents["exechealthz"])
			addValue(parametersMap, "kubernetesDNSSidecarSpec", cloudSpecConfig.KubernetesSpecConfig.KubernetesImageBase+k8sComponents["k8s-dns-sidecar"])
			addValue(parametersMap, "kubernetesHeapsterSpec", cloudSpecConfig.KubernetesSpecConfig.KubernetesImageBase+k8sComponents["heapster"])
			if kubernetesConfig.IsTillerEnabled() {
				tillerAddon := kubernetesConfig.GetAddonByName(DefaultTillerAddonName)
				c := tillerAddon.GetAddonContainersIndexByName(DefaultTillerAddonName)
				if c > -1 {
					addValue(parametersMap, "kubernetesTillerCPURequests", tillerAddon.Containers[c].CPURequests)
					addValue(parametersMap, "kubernetesTillerCPULimit", tillerAddon.Containers[c].CPULimits)
					addValue(parametersMap, "kubernetesTillerMemoryRequests", tillerAddon.Containers[c].MemoryRequests)
					addValue(parametersMap, "kubernetesTillerMemoryLimit", tillerAddon.Containers[c].MemoryLimits)
					addValue(parametersMap, "kubernetesTillerMaxHistory", tillerAddon.Config["max-history"])
					if tillerAddon.Containers[c].Image != "" {
						addValue(parametersMap, "kubernetesTillerSpec", tillerAddon.Containers[c].Image)
					} else {
						addValue(parametersMap, "kubernetesTillerSpec", cloudSpecConfig.KubernetesSpecConfig.TillerImageBase+k8sComponents[DefaultTillerAddonName])
					}
				}
			}
			if kubernetesConfig.IsAADPodIdentityEnabled() {
				aadPodIdentityAddon := kubernetesConfig.GetAddonByName(DefaultAADPodIdentityAddonName)
				c := aadPodIdentityAddon.GetAddonContainersIndexByName(DefaultAADPodIdentityAddonName)
				if c > -1 {
					addValue(parametersMap, "kubernetesAADPodIdentityEnabled", helpers.IsTrueBoolPointer(aadPodIdentityAddon.Enabled))
				}
			}
			if kubernetesConfig.IsACIConnectorEnabled() {
				aciConnectorAddon := kubernetesConfig.GetAddonByName(DefaultACIConnectorAddonName)
				c := aciConnectorAddon.GetAddonContainersIndexByName(DefaultACIConnectorAddonName)
				if c > -1 {
					addValue(parametersMap, "kubernetesACIConnectorEnabled", true)
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
						addValue(parametersMap, "kubernetesACIConnectorSpec", cloudSpecConfig.KubernetesSpecConfig.ACIConnectorImageBase+k8sComponents[DefaultACIConnectorAddonName])
					}
				}
			} else {
				addValue(parametersMap, "kubernetesACIConnectorEnabled", false)
			}
			if kubernetesConfig.IsClusterAutoscalerEnabled() {
				clusterAutoscalerAddon := kubernetesConfig.GetAddonByName(DefaultClusterAutoscalerAddonName)
				c := clusterAutoscalerAddon.GetAddonContainersIndexByName(DefaultClusterAutoscalerAddonName)
				if c > -1 {
					addValue(parametersMap, "kubernetesClusterAutoscalerAzureCloud", cloudSpecConfig.CloudName)
					addValue(parametersMap, "kubernetesClusterAutoscalerCPURequests", clusterAutoscalerAddon.Containers[c].CPURequests)
					addValue(parametersMap, "kubernetesClusterAutoscalerCPULimit", clusterAutoscalerAddon.Containers[c].CPULimits)
					addValue(parametersMap, "kubernetesClusterAutoscalerMemoryRequests", clusterAutoscalerAddon.Containers[c].MemoryRequests)
					addValue(parametersMap, "kubernetesClusterAutoscalerMemoryLimit", clusterAutoscalerAddon.Containers[c].MemoryLimits)
					addValue(parametersMap, "kubernetesClusterAutoscalerMinNodes", clusterAutoscalerAddon.Config["minNodes"])
					addValue(parametersMap, "kubernetesClusterAutoscalerMaxNodes", clusterAutoscalerAddon.Config["maxNodes"])
					addValue(parametersMap, "kubernetesClusterAutoscalerEnabled", true)
					addValue(parametersMap, "kubernetesClusterAutoscalerUseManagedIdentity", strings.ToLower(strconv.FormatBool(kubernetesConfig.UseManagedIdentity)))
					if clusterAutoscalerAddon.Containers[c].Image != "" {
						addValue(parametersMap, "kubernetesClusterAutoscalerSpec", clusterAutoscalerAddon.Containers[c].Image)
					} else {
						addValue(parametersMap, "kubernetesClusterAutoscalerSpec", cloudSpecConfig.KubernetesSpecConfig.KubernetesImageBase+k8sComponents[DefaultClusterAutoscalerAddonName])
					}
				}
			} else {
				addValue(parametersMap, "kubernetesClusterAutoscalerEnabled", false)
			}
			flexVolumeDriverConfig := map[string]string{}
			bfFlexVolumeInstallerAddon := kubernetesConfig.GetAddonByName(DefaultBlobfuseFlexVolumeAddonName)
			c := bfFlexVolumeInstallerAddon.GetAddonContainersIndexByName(DefaultBlobfuseFlexVolumeAddonName)
			if c > -1 {
				flexVolumeDriverConfig["kubernetesBlobfuseFlexVolumeInstallerCPURequests"] = bfFlexVolumeInstallerAddon.Containers[c].CPURequests
				flexVolumeDriverConfig["kubernetesBlobfuseFlexVolumeInstallerCPULimit"] = bfFlexVolumeInstallerAddon.Containers[c].CPULimits
				flexVolumeDriverConfig["kubernetesBlobfuseFlexVolumeInstallerMemoryRequests"] = bfFlexVolumeInstallerAddon.Containers[c].MemoryRequests
				flexVolumeDriverConfig["kubernetesBlobfuseFlexVolumeInstallerMemoryLimit"] = bfFlexVolumeInstallerAddon.Containers[c].MemoryLimits
			}
			smbFlexVolumeInstallerAddon := kubernetesConfig.GetAddonByName(DefaultSMBFlexVolumeAddonName)
			c = smbFlexVolumeInstallerAddon.GetAddonContainersIndexByName(DefaultSMBFlexVolumeAddonName)
			if c > -1 {
				flexVolumeDriverConfig["kubernetesSMBFlexVolumeInstallerCPURequests"] = smbFlexVolumeInstallerAddon.Containers[c].CPURequests
				flexVolumeDriverConfig["kubernetesSMBFlexVolumeInstallerCPULimit"] = smbFlexVolumeInstallerAddon.Containers[c].CPULimits
				flexVolumeDriverConfig["kubernetesSMBFlexVolumeInstallerMemoryRequests"] = smbFlexVolumeInstallerAddon.Containers[c].MemoryRequests
				flexVolumeDriverConfig["kubernetesSMBFlexVolumeInstallerMemoryLimit"] = smbFlexVolumeInstallerAddon.Containers[c].MemoryLimits
			}
			addValue(parametersMap, "flexVolumeDriverConfig", flexVolumeDriverConfig)
			if kubernetesConfig.IsKeyVaultFlexVolumeEnabled() {
				kvFlexVolumeInstallerAddon := kubernetesConfig.GetAddonByName(DefaultKeyVaultFlexVolumeAddonName)
				c := kvFlexVolumeInstallerAddon.GetAddonContainersIndexByName(DefaultKeyVaultFlexVolumeAddonName)
				if c > -1 {
					addValue(parametersMap, "kubernetesKeyVaultFlexVolumeInstallerCPURequests", kvFlexVolumeInstallerAddon.Containers[c].CPURequests)
					addValue(parametersMap, "kubernetesKeyVaultFlexVolumeInstallerCPULimit", kvFlexVolumeInstallerAddon.Containers[c].CPULimits)
					addValue(parametersMap, "kubernetesKeyVaultFlexVolumeInstallerMemoryRequests", kvFlexVolumeInstallerAddon.Containers[c].MemoryRequests)
					addValue(parametersMap, "kubernetesKeyVaultFlexVolumeInstallerMemoryLimit", kvFlexVolumeInstallerAddon.Containers[c].MemoryLimits)
				}
			}
			if kubernetesConfig.IsDashboardEnabled() {
				dashboardAddon := kubernetesConfig.GetAddonByName(DefaultDashboardAddonName)
				dashboardAddon.GetAddonContainersIndexByName(DefaultDashboardAddonName)
				if c > -1 {
					addValue(parametersMap, "kubernetesDashboardCPURequests", dashboardAddon.Containers[c].CPURequests)
					addValue(parametersMap, "kubernetesDashboardCPULimit", dashboardAddon.Containers[c].CPULimits)
					addValue(parametersMap, "kubernetesDashboardMemoryRequests", dashboardAddon.Containers[c].MemoryRequests)
					addValue(parametersMap, "kubernetesDashboardMemoryLimit", dashboardAddon.Containers[c].MemoryLimits)
					if dashboardAddon.Containers[c].Image != "" {
						addValue(parametersMap, "kubernetesDashboardSpec", dashboardAddon.Containers[c].Image)
					} else {
						addValue(parametersMap, "kubernetesDashboardSpec", cloudSpecConfig.KubernetesSpecConfig.KubernetesImageBase+k8sComponents[DefaultDashboardAddonName])
					}
				}
			}
			if kubernetesConfig.IsReschedulerEnabled() {
				reschedulerAddon := kubernetesConfig.GetAddonByName(DefaultReschedulerAddonName)
				c := reschedulerAddon.GetAddonContainersIndexByName(DefaultReschedulerAddonName)
				if c > -1 {
					addValue(parametersMap, "kubernetesReschedulerCPURequests", reschedulerAddon.Containers[c].CPURequests)
					addValue(parametersMap, "kubernetesReschedulerCPULimit", reschedulerAddon.Containers[c].CPULimits)
					addValue(parametersMap, "kubernetesReschedulerMemoryRequests", reschedulerAddon.Containers[c].MemoryRequests)
					addValue(parametersMap, "kubernetesReschedulerMemoryLimit", reschedulerAddon.Containers[c].MemoryLimits)
					if reschedulerAddon.Containers[c].Image != "" {
						addValue(parametersMap, "kubernetesReschedulerSpec", reschedulerAddon.Containers[c].Image)
					} else {
						addValue(parametersMap, "kubernetesReschedulerSpec", cloudSpecConfig.KubernetesSpecConfig.KubernetesImageBase+k8sComponents[DefaultReschedulerAddonName])
					}
				}
			}
			if properties.OrchestratorProfile.IsMetricsServerEnabled() {
				metricsServerAddon := kubernetesConfig.GetAddonByName(DefaultMetricsServerAddonName)
				c = metricsServerAddon.GetAddonContainersIndexByName(DefaultMetricsServerAddonName)
				if c > -1 {
					if metricsServerAddon.Containers[c].Image != "" {
						addValue(parametersMap, "kubernetesMetricsServerSpec", metricsServerAddon.Containers[c].Image)
					} else {
						addValue(parametersMap, "kubernetesMetricsServerSpec", cloudSpecConfig.KubernetesSpecConfig.KubernetesImageBase+k8sComponents[DefaultMetricsServerAddonName])
					}
				}
			}
			if properties.IsNVIDIADevicePluginEnabled() {
				nvidiaDevicePluginAddon := kubernetesConfig.GetAddonByName(NVIDIADevicePluginAddonName)
				c := nvidiaDevicePluginAddon.GetAddonContainersIndexByName(NVIDIADevicePluginAddonName)
				if c > -1 {
					addValue(parametersMap, "kubernetesNVIDIADevicePluginCPURequests", nvidiaDevicePluginAddon.Containers[c].CPURequests)
					addValue(parametersMap, "kubernetesNVIDIADevicePluginCPULimit", nvidiaDevicePluginAddon.Containers[c].CPULimits)
					addValue(parametersMap, "kubernetesNVIDIADevicePluginMemoryRequests", nvidiaDevicePluginAddon.Containers[c].MemoryRequests)
					addValue(parametersMap, "kubernetesNVIDIADevicePluginMemoryLimit", nvidiaDevicePluginAddon.Containers[c].MemoryLimits)
					if nvidiaDevicePluginAddon.Containers[c].Image != "" {
						addValue(parametersMap, "kubernetesNVIDIADevicePluginSpec", nvidiaDevicePluginAddon.Containers[c].Image)
					} else {
						addValue(parametersMap, "kubernetesNVIDIADevicePluginSpec", cloudSpecConfig.KubernetesSpecConfig.NVIDIAImageBase+k8sComponents[NVIDIADevicePluginAddonName])
					}
				}
			}
			if kubernetesConfig.IsContainerMonitoringEnabled() {
				containerMonitoringAddon := kubernetesConfig.GetAddonByName(ContainerMonitoringAddonName)
				c := containerMonitoringAddon.GetAddonContainersIndexByName("omsagent")
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
						addValue(parametersMap, "omsAgentImage", cloudSpecConfig.KubernetesSpecConfig.KubernetesImageBase+k8sComponents[ContainerMonitoringAddonName])
					}
				}
			}
			if kubernetesConfig.IsIPMasqAgentEnabled() {
				ipMasqAgentAddon := kubernetesConfig.GetAddonByName(IPMASQAgentAddonName)
				i := ipMasqAgentAddon.GetAddonContainersIndexByName(IPMASQAgentAddonName)
				if i > -1 {
					addValue(parametersMap, "kubernetesIPMasqAgentCPURequests", ipMasqAgentAddon.Containers[c].CPURequests)
					addValue(parametersMap, "kubernetesIPMasqAgentMemoryRequests", ipMasqAgentAddon.Containers[c].MemoryRequests)
					addValue(parametersMap, "kubernetesIPMasqAgentCPULimit", ipMasqAgentAddon.Containers[c].CPULimits)
					addValue(parametersMap, "kubernetesIPMasqAgentMemoryLimit", ipMasqAgentAddon.Containers[c].MemoryLimits)
				}
			}
			if kubernetesConfig.LoadBalancerSku == "Standard" {
				random := rand.New(rand.NewSource(time.Now().UnixNano()))
				elbsvcName := random.Int()
				addValue(parametersMap, "kuberneteselbsvcname", fmt.Sprintf("%d", elbsvcName))
			}

			if properties.OrchestratorProfile.IsAzureCNI() {
				azureCNINetworkmonitorAddon := kubernetesConfig.GetAddonByName(AzureCNINetworkMonitoringAddonName)
				c = azureCNINetworkmonitorAddon.GetAddonContainersIndexByName(AzureCNINetworkMonitoringAddonName)
				if c > -1 {
					if azureCNINetworkmonitorAddon.Containers[c].Image != "" {
						addValue(parametersMap, "AzureCNINetworkMonitorImageURL", azureCNINetworkmonitorAddon.Containers[c].Image)
					} else {
						addValue(parametersMap, "AzureCNINetworkMonitorImageURL", cloudSpecConfig.KubernetesSpecConfig.AzureCNIImageBase+k8sComponents[AzureCNINetworkMonitoringAddonName])
					}
				}
			}
			addValue(parametersMap, "kubernetesCoreDNSSpec", "coredns/coredns:1.2.2")
			addValue(parametersMap, "kubernetesKubeDNSSpec", cloudSpecConfig.KubernetesSpecConfig.KubernetesImageBase+k8sComponents["dns"])
			addValue(parametersMap, "kubernetesPodInfraContainerSpec", cloudSpecConfig.KubernetesSpecConfig.KubernetesImageBase+k8sComponents["pause"])
			addValue(parametersMap, "cloudproviderConfig", api.CloudProviderConfig{
				CloudProviderBackoff:         kubernetesConfig.CloudProviderBackoff,
				CloudProviderBackoffRetries:  kubernetesConfig.CloudProviderBackoffRetries,
				CloudProviderBackoffJitter:   strconv.FormatFloat(kubernetesConfig.CloudProviderBackoffJitter, 'f', -1, 64),
				CloudProviderBackoffDuration: kubernetesConfig.CloudProviderBackoffDuration,
				CloudProviderBackoffExponent: strconv.FormatFloat(kubernetesConfig.CloudProviderBackoffExponent, 'f', -1, 64),
				CloudProviderRateLimit:       kubernetesConfig.CloudProviderRateLimit,
				CloudProviderRateLimitQPS:    strconv.FormatFloat(kubernetesConfig.CloudProviderRateLimitQPS, 'f', -1, 64),
				CloudProviderRateLimitBucket: kubernetesConfig.CloudProviderRateLimitBucket,
			})
			addValue(parametersMap, "kubeClusterCidr", kubernetesConfig.ClusterSubnet)
			if !properties.IsHostedMasterProfile() {
				if properties.OrchestratorProfile.IsAzureCNI() {
					if properties.MasterProfile != nil && properties.MasterProfile.IsCustomVNET() {
						addValue(parametersMap, "kubernetesNonMasqueradeCidr", properties.MasterProfile.VnetCidr)
					} else {
						addValue(parametersMap, "kubernetesNonMasqueradeCidr", DefaultVNETCIDR)
					}
				} else {
					addValue(parametersMap, "kubernetesNonMasqueradeCidr", properties.OrchestratorProfile.KubernetesConfig.ClusterSubnet)
				}
			}
			addValue(parametersMap, "kubernetesKubeletClusterDomain", kubernetesConfig.KubeletConfig["--cluster-domain"])
			addValue(parametersMap, "dockerBridgeCidr", kubernetesConfig.DockerBridgeSubnet)
			addValue(parametersMap, "networkPolicy", kubernetesConfig.NetworkPolicy)
			addValue(parametersMap, "networkPlugin", kubernetesConfig.NetworkPlugin)
			addValue(parametersMap, "containerRuntime", kubernetesConfig.ContainerRuntime)
			addValue(parametersMap, "containerdDownloadURLBase", cloudSpecConfig.KubernetesSpecConfig.ContainerdDownloadURLBase)
			addValue(parametersMap, "cniPluginsURL", cloudSpecConfig.KubernetesSpecConfig.CNIPluginsDownloadURL)
			addValue(parametersMap, "vnetCniLinuxPluginsURL", cloudSpecConfig.KubernetesSpecConfig.VnetCNILinuxPluginsDownloadURL)
			addValue(parametersMap, "vnetCniWindowsPluginsURL", cloudSpecConfig.KubernetesSpecConfig.VnetCNIWindowsPluginsDownloadURL)
			addValue(parametersMap, "gchighthreshold", kubernetesConfig.GCHighThreshold)
			addValue(parametersMap, "gclowthreshold", kubernetesConfig.GCLowThreshold)
			addValue(parametersMap, "etcdDownloadURLBase", cloudSpecConfig.KubernetesSpecConfig.EtcdDownloadURLBase)
			addValue(parametersMap, "etcdVersion", kubernetesConfig.EtcdVersion)
			addValue(parametersMap, "etcdDiskSizeGB", kubernetesConfig.EtcdDiskSizeGB)
			addValue(parametersMap, "etcdEncryptionKey", kubernetesConfig.EtcdEncryptionKey)
			if kubernetesConfig.PrivateJumpboxProvision() {
				addValue(parametersMap, "jumpboxVMName", kubernetesConfig.PrivateCluster.JumpboxProfile.Name)
				addValue(parametersMap, "jumpboxVMSize", kubernetesConfig.PrivateCluster.JumpboxProfile.VMSize)
				addValue(parametersMap, "jumpboxUsername", kubernetesConfig.PrivateCluster.JumpboxProfile.Username)
				addValue(parametersMap, "jumpboxOSDiskSizeGB", kubernetesConfig.PrivateCluster.JumpboxProfile.OSDiskSizeGB)
				addValue(parametersMap, "jumpboxPublicKey", kubernetesConfig.PrivateCluster.JumpboxProfile.PublicKey)
				addValue(parametersMap, "jumpboxStorageProfile", kubernetesConfig.PrivateCluster.JumpboxProfile.StorageProfile)
			}

			if kubernetesConfig.DockerEngineVersion != "" {
				dockerEngineVersion = kubernetesConfig.DockerEngineVersion
			}

			addValue(parametersMap, "enableAggregatedAPIs", kubernetesConfig.EnableAggregatedAPIs)
		}

		if kubernetesConfig == nil ||
			!kubernetesConfig.UseManagedIdentity {

			addValue(parametersMap, "servicePrincipalClientId", properties.ServicePrincipalProfile.ClientID)
			if properties.ServicePrincipalProfile.KeyvaultSecretRef != nil {
				addKeyvaultReference(parametersMap, "servicePrincipalClientSecret",
					properties.ServicePrincipalProfile.KeyvaultSecretRef.VaultID,
					properties.ServicePrincipalProfile.KeyvaultSecretRef.SecretName,
					properties.ServicePrincipalProfile.KeyvaultSecretRef.SecretVersion)
			} else {
				addValue(parametersMap, "servicePrincipalClientSecret", properties.ServicePrincipalProfile.Secret)
			}

			if kubernetesConfig != nil && helpers.IsTrueBoolPointer(kubernetesConfig.EnableEncryptionWithExternalKms) && !kubernetesConfig.UseManagedIdentity && properties.ServicePrincipalProfile.ObjectID != "" {
				addValue(parametersMap, "servicePrincipalObjectId", properties.ServicePrincipalProfile.ObjectID)
			}
		}

		addValue(parametersMap, "orchestratorName", properties.K8sOrchestratorName())

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

		certificateProfile := properties.CertificateProfile
		if certificateProfile != nil {
			addSecret(parametersMap, "apiServerCertificate", certificateProfile.APIServerCertificate, true)
			addSecret(parametersMap, "apiServerPrivateKey", certificateProfile.APIServerPrivateKey, true)
			addSecret(parametersMap, "caCertificate", certificateProfile.CaCertificate, true)
			addSecret(parametersMap, "caPrivateKey", certificateProfile.CaPrivateKey, true)
			addSecret(parametersMap, "clientCertificate", certificateProfile.ClientCertificate, true)
			addSecret(parametersMap, "clientPrivateKey", certificateProfile.ClientPrivateKey, true)
			addSecret(parametersMap, "kubeConfigCertificate", certificateProfile.KubeConfigCertificate, true)
			addSecret(parametersMap, "kubeConfigPrivateKey", certificateProfile.KubeConfigPrivateKey, true)
			if properties.MasterProfile != nil {
				addSecret(parametersMap, "etcdServerCertificate", certificateProfile.EtcdServerCertificate, true)
				addSecret(parametersMap, "etcdServerPrivateKey", certificateProfile.EtcdServerPrivateKey, true)
				addSecret(parametersMap, "etcdClientCertificate", certificateProfile.EtcdClientCertificate, true)
				addSecret(parametersMap, "etcdClientPrivateKey", certificateProfile.EtcdClientPrivateKey, true)
				for i, pc := range certificateProfile.EtcdPeerCertificates {
					addSecret(parametersMap, "etcdPeerCertificate"+strconv.Itoa(i), pc, true)
				}
				for i, pk := range certificateProfile.EtcdPeerPrivateKeys {
					addSecret(parametersMap, "etcdPeerPrivateKey"+strconv.Itoa(i), pk, true)
				}
			}
		}

		if properties.HostedMasterProfile != nil && properties.HostedMasterProfile.FQDN != "" {
			addValue(parametersMap, "kubernetesEndpoint", properties.HostedMasterProfile.FQDN)
		}

		if !orchestratorProfile.IsOpenShift() {
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

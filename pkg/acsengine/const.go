package acsengine

const (
	// DefaultOpenShiftMasterSubnet is the default value for master subnet for Openshift.
	DefaultOpenShiftMasterSubnet = "10.0.0.0/24"
	// DefaultOpenShiftFirstConsecutiveStaticIP is the default static ip address for master 0 for Openshift.
	DefaultOpenShiftFirstConsecutiveStaticIP = "10.0.0.11"
	// DefaultMasterSubnet specifies the default master subnet for DCOS or Swarm
	DefaultMasterSubnet = "172.16.0.0/24"
	// DefaultFirstConsecutiveStaticIP specifies the static IP address on master 0 for DCOS or Swarm
	DefaultFirstConsecutiveStaticIP = "172.16.0.5"
	// DefaultSwarmWindowsMasterSubnet specifies the default master subnet for a Swarm Windows cluster
	DefaultSwarmWindowsMasterSubnet = "192.168.255.0/24"
	// DefaultSwarmWindowsFirstConsecutiveStaticIP specifies the static IP address on master 0 for a Swarm WIndows cluster
	DefaultSwarmWindowsFirstConsecutiveStaticIP = "192.168.255.5"
	// DefaultDCOSMasterSubnet specifies the default master subnet for a DCOS cluster
	DefaultDCOSMasterSubnet = "192.168.255.0/24"
	// DefaultDCOSFirstConsecutiveStaticIP  specifies the static IP address on master 0 for a DCOS cluster
	DefaultDCOSFirstConsecutiveStaticIP = "192.168.255.5"
	// DefaultDCOSBootstrapStaticIP specifies the static IP address on bootstrap for a DCOS cluster
	DefaultDCOSBootstrapStaticIP = "192.168.255.240"
	// DefaultKubernetesMasterSubnet specifies the default subnet for masters and agents.
	DefaultKubernetesMasterSubnet = "10.240.0.0/16"
	// DefaultKubernetesClusterSubnet specifies the default subnet for pods.
	DefaultKubernetesClusterSubnet = "10.244.0.0/16"
	// DefaultDockerBridgeSubnet specifies the default subnet for the docker bridge network for masters and agents.
	DefaultDockerBridgeSubnet = "172.17.0.1/16"
	// DefaultFirstConsecutiveKubernetesStaticIP specifies the static IP address on Kubernetes master 0
	DefaultFirstConsecutiveKubernetesStaticIP = "10.240.255.5"
	// DefaultAgentSubnetTemplate specifies a default agent subnet
	DefaultAgentSubnetTemplate = "10.%d.0.0/16"
	// DefaultKubernetesSubnet specifies the default subnet used for all masters, agents and pods
	// when VNET integration is enabled.
	DefaultKubernetesSubnet = "10.240.0.0/12"
	// DefaultKubernetesFirstConsecutiveStaticIPOffset specifies the IP address offset of master 0
	// when VNET integration is enabled.
	DefaultKubernetesFirstConsecutiveStaticIPOffset = 5
	// DefaultKubernetesMaxPods is the maximum number of pods to run on a node.
	DefaultKubernetesMaxPods = 110
	// DefaultKubernetesMaxPodsVNETIntegrated is the maximum number of pods to run on a node when VNET integration is enabled.
	DefaultKubernetesMaxPodsVNETIntegrated = 30
	// DefaultKubernetesClusterDomain is the dns suffix used in the cluster (used as a SAN in the PKI generation)
	DefaultKubernetesClusterDomain = "cluster.local"
	// DefaultInternalLbStaticIPOffset specifies the offset of the internal LoadBalancer's IP
	// address relative to the first consecutive Kubernetes static IP
	DefaultInternalLbStaticIPOffset = 10
	// NetworkPolicyNone is the string expression for the deprecated NetworkPolicy usage pattern "none"
	NetworkPolicyNone = "none"
	// NetworkPolicyCalico is the string expression for calico network policy config option
	NetworkPolicyCalico = "calico"
	// NetworkPolicyCilium is the string expression for cilium network policy config option
	NetworkPolicyCilium = "cilium"
	// NetworkPolicyAzure is the string expression for Azure CNI network policy manager
	NetworkPolicyAzure = "azure"
	// NetworkPluginAzure is the string expression for Azure CNI plugin
	NetworkPluginAzure = "azure"
	// NetworkPluginKubenet is the string expression for kubenet network plugin
	NetworkPluginKubenet = "kubenet"
	// NetworkPluginFlannel is the string expression for flannel network policy config option
	NetworkPluginFlannel = "flannel"
	// DefaultNetworkPlugin defines the network plugin to use by default
	DefaultNetworkPlugin = NetworkPluginKubenet
	// DefaultNetworkPolicy defines the network policy implementation to use by default
	DefaultNetworkPolicy = ""
	// DefaultNetworkPluginWindows defines the network plugin implementation to use by default for clusters with Windows agent pools
	DefaultNetworkPluginWindows = NetworkPluginKubenet
	// DefaultNetworkPolicyWindows defines the network policy implementation to use by default for clusters with Windows agent pools
	DefaultNetworkPolicyWindows = ""
	// DefaultContainerRuntime is docker
	DefaultContainerRuntime = "docker"
	// DefaultKubernetesNodeStatusUpdateFrequency is 10s, see --node-status-update-frequency at https://kubernetes.io/docs/admin/kubelet/
	DefaultKubernetesNodeStatusUpdateFrequency = "10s"
	// DefaultKubernetesHardEvictionThreshold is memory.available<100Mi,nodefs.available<10%,nodefs.inodesFree<5%, see --eviction-hard at https://kubernetes.io/docs/admin/kubelet/
	DefaultKubernetesHardEvictionThreshold = "memory.available<100Mi,nodefs.available<10%,nodefs.inodesFree<5%"
	// DefaultKubernetesCtrlMgrNodeMonitorGracePeriod is 40s, see --node-monitor-grace-period at https://kubernetes.io/docs/admin/kube-controller-manager/
	DefaultKubernetesCtrlMgrNodeMonitorGracePeriod = "40s"
	// DefaultKubernetesCtrlMgrPodEvictionTimeout is 5m0s, see --pod-eviction-timeout at https://kubernetes.io/docs/admin/kube-controller-manager/
	DefaultKubernetesCtrlMgrPodEvictionTimeout = "5m0s"
	// DefaultKubernetesCtrlMgrRouteReconciliationPeriod is 10s, see --route-reconciliation-period at https://kubernetes.io/docs/admin/kube-controller-manager/
	DefaultKubernetesCtrlMgrRouteReconciliationPeriod = "10s"
	// DefaultKubernetesCtrlMgrTerminatedPodGcThreshold is set to 5000, see --terminated-pod-gc-threshold at https://kubernetes.io/docs/admin/kube-controller-manager/ and https://github.com/kubernetes/kubernetes/issues/22680
	DefaultKubernetesCtrlMgrTerminatedPodGcThreshold = "5000"
	// DefaultKubernetesCtrlMgrUseSvcAccountCreds is "true", see --use-service-account-credentials at https://kubernetes.io/docs/admin/kube-controller-manager/
	DefaultKubernetesCtrlMgrUseSvcAccountCreds = "false"
	// DefaultKubernetesCloudProviderBackoff is false to disable cloudprovider backoff implementation for API calls
	DefaultKubernetesCloudProviderBackoff = false
	// DefaultKubernetesCloudProviderBackoffRetries is 6, takes effect if DefaultKubernetesCloudProviderBackoff is true
	DefaultKubernetesCloudProviderBackoffRetries = 6
	// DefaultKubernetesCloudProviderBackoffJitter is 1, takes effect if DefaultKubernetesCloudProviderBackoff is true
	DefaultKubernetesCloudProviderBackoffJitter = 1.0
	// DefaultKubernetesCloudProviderBackoffDuration is 5, takes effect if DefaultKubernetesCloudProviderBackoff is true
	DefaultKubernetesCloudProviderBackoffDuration = 5
	// DefaultKubernetesCloudProviderBackoffExponent is 1.5, takes effect if DefaultKubernetesCloudProviderBackoff is true
	DefaultKubernetesCloudProviderBackoffExponent = 1.5
	// DefaultKubernetesCloudProviderRateLimit is false to disable cloudprovider rate limiting implementation for API calls
	DefaultKubernetesCloudProviderRateLimit = false
	// DefaultKubernetesCloudProviderRateLimitQPS is 3, takes effect if DefaultKubernetesCloudProviderRateLimit is true
	DefaultKubernetesCloudProviderRateLimitQPS = 3.0
	// DefaultKubernetesCloudProviderRateLimitBucket is 10, takes effect if DefaultKubernetesCloudProviderRateLimit is true
	DefaultKubernetesCloudProviderRateLimitBucket = 10
	// DefaultTillerAddonName is the name of the tiller addon deployment
	DefaultTillerAddonName = "tiller"
	// DefaultTillerMaxHistory limits the maximum number of revisions saved per release. Use 0 for no limit.
	DefaultTillerMaxHistory = 0
	// DefaultACIConnectorAddonName is the name of the tiller addon deployment
	DefaultACIConnectorAddonName = "aci-connector"
	// DefaultDashboardAddonName is the name of the kubernetes-dashboard addon deployment
	DefaultDashboardAddonName = "kubernetes-dashboard"
	// DefaultClusterAutoscalerAddonName is the name of the autoscaler addon deployment
	DefaultClusterAutoscalerAddonName = "cluster-autoscaler"
	// DefaultKubernetesDNSServiceIP specifies the IP address that kube-dns
	// listens on by default. must by in the default Service CIDR range.
	DefaultKubernetesDNSServiceIP = "10.0.0.10"
	// DefaultKubernetesServiceCIDR specifies the IP subnet that kubernetes will
	// create Service IPs within.
	DefaultKubernetesServiceCIDR = "10.0.0.0/16"
	//DefaultKubernetesGCHighThreshold specifies the value for  for the image-gc-high-threshold kubelet flag
	DefaultKubernetesGCHighThreshold = 85
	//DefaultKubernetesGCLowThreshold specifies the value for the image-gc-low-threshold kubelet flag
	DefaultKubernetesGCLowThreshold = 80
	// DefaultGeneratorCode specifies the source generator of the cluster template.
	DefaultGeneratorCode = "acsengine"
	// DefaultOrchestratorName specifies the 3 character orchestrator code of the cluster template and affects resource naming.
	DefaultOrchestratorName = "k8s"
	// DefaultOpenshiftOrchestratorName specifies the 3 character orchestrator code of the cluster template and affects resource naming.
	DefaultOpenshiftOrchestratorName = "ocp"
	// DefaultEtcdVersion specifies the default etcd version to install
	DefaultEtcdVersion = "3.2.23"
	// DefaultEtcdDiskSize specifies the default size for Kubernetes master etcd disk volumes in GB
	DefaultEtcdDiskSize = "256"
	// DefaultEtcdDiskSizeGT3Nodes = size for Kubernetes master etcd disk volumes in GB if > 3 nodes
	DefaultEtcdDiskSizeGT3Nodes = "512"
	// DefaultEtcdDiskSizeGT10Nodes = size for Kubernetes master etcd disk volumes in GB if > 10 nodes
	DefaultEtcdDiskSizeGT10Nodes = "1024"
	// DefaultEtcdDiskSizeGT20Nodes = size for Kubernetes master etcd disk volumes in GB if > 20 nodes
	DefaultEtcdDiskSizeGT20Nodes = "2048"
	// DefaultReschedulerAddonName is the name of the rescheduler addon deployment
	DefaultReschedulerAddonName = "rescheduler"
	// DefaultMetricsServerAddonName is the name of the kubernetes Metrics server addon deployment
	DefaultMetricsServerAddonName = "metrics-server"
	// NVIDIADevicePluginAddonName is the name of the kubernetes NVIDIA Device Plugin daemon set
	NVIDIADevicePluginAddonName = "nvidia-device-plugin"
	// ContainerMonitoringAddonName is the name of the kubernetes Container Monitoring addon deployment
	ContainerMonitoringAddonName = "container-monitoring"
	// AzureCNINetworkMonitoringAddonName is the name of the Azure CNI networkmonitor addon
	AzureCNINetworkMonitoringAddonName = "azure-cni-networkmonitor"
	// AzureNetworkPolicyAddonName is the name of the Azure CNI networkmonitor addon
	AzureNetworkPolicyAddonName = "azure-npm-daemonset"
	// DefaultKubernetesKubeletMaxPods is the max pods per kubelet
	DefaultKubernetesKubeletMaxPods = 110
	// DefaultMasterEtcdServerPort is the default etcd server port for Kubernetes master nodes
	DefaultMasterEtcdServerPort = 2380
	// DefaultMasterEtcdClientPort is the default etcd client port for Kubernetes master nodes
	DefaultMasterEtcdClientPort = 2379
	// DefaultKubeletEventQPS is 0, see --event-qps at https://kubernetes.io/docs/reference/generated/kubelet/
	DefaultKubeletEventQPS = "0"
	// DefaultKubeletCadvisorPort is 0, see --cadvisor-port at https://kubernetes.io/docs/reference/generated/kubelet/
	DefaultKubeletCadvisorPort = "0"
	// DefaultJumpboxDiskSize specifies the default size for private cluster jumpbox OS disk in GB
	DefaultJumpboxDiskSize = 30
	// DefaultJumpboxUsername specifies the default admin username for the private cluster jumpbox
	DefaultJumpboxUsername = "azureuser"
	// DefaultKubeletPodMaxPIDs specifies the default max pid authorized by pods
	DefaultKubeletPodMaxPIDs = 100
)

const (
	// DCOSMaster represents the master node type
	DCOSMaster DCOSNodeType = "DCOSMaster"
	// DCOSPrivateAgent represents the private agent node type
	DCOSPrivateAgent DCOSNodeType = "DCOSPrivateAgent"
	// DCOSPublicAgent represents the public agent node type
	DCOSPublicAgent DCOSNodeType = "DCOSPublicAgent"
)

const (
	//DefaultExtensionsRootURL  Root URL for extensions
	DefaultExtensionsRootURL = "https://raw.githubusercontent.com/Azure/acs-engine/master/"
	// DefaultDockerEngineRepo for grabbing docker engine packages
	DefaultDockerEngineRepo = "https://download.docker.com/linux/ubuntu"
	// DefaultDockerComposeURL for grabbing docker images
	DefaultDockerComposeURL = "https://github.com/docker/compose/releases/download"

	//AzureEdgeDCOSBootstrapDownloadURL is the azure edge CDN download url
	AzureEdgeDCOSBootstrapDownloadURL = "https://dcosio.azureedge.net/dcos/%s/bootstrap/%s.bootstrap.tar.xz"
	//AzureChinaCloudDCOSBootstrapDownloadURL is the China specific DCOS package download url.
	AzureChinaCloudDCOSBootstrapDownloadURL = "https://acsengine.blob.core.chinacloudapi.cn/dcos/%s.bootstrap.tar.xz"
	//AzureEdgeDCOSWindowsBootstrapDownloadURL
)

const (
	//DefaultConfigurationScriptRootURL  Root URL for configuration script (used for script extension on RHEL)
	DefaultConfigurationScriptRootURL = "https://raw.githubusercontent.com/Azure/acs-engine/master/parts/"
)

const (
	kubernetesMasterCustomDataYaml           = "k8s/kubernetesmastercustomdata.yml"
	kubernetesCustomScript                   = "k8s/kubernetescustomscript.sh"
	kubernetesProvisionSourceScript          = "k8s/kubernetesprovisionsource.sh"
	kubernetesMountetcd                      = "k8s/kubernetes_mountetcd.sh"
	kubernetesCustomSearchDomainsScript      = "k8s/setup-custom-search-domains.sh"
	kubernetesMasterGenerateProxyCertsScript = "k8s/kubernetesmastergenerateproxycertscript.sh"
	kubernetesAgentCustomDataYaml            = "k8s/kubernetesagentcustomdata.yml"
	kubernetesJumpboxCustomDataYaml          = "k8s/kubernetesjumpboxcustomdata.yml"
	kubeConfigJSON                           = "k8s/kubeconfig.json"
	kubernetesWindowsAgentCustomDataPS1      = "k8s/kuberneteswindowssetup.ps1"
	// OpenShift custom scripts
	openshiftNodeScript     = "openshift/unstable/openshiftnodescript.sh"
	openshiftMasterScript   = "openshift/unstable/openshiftmasterscript.sh"
	openshift39NodeScript   = "openshift/release-3.9/openshiftnodescript.sh"
	openshift39MasterScript = "openshift/release-3.9/openshiftmasterscript.sh"
)

const (
	dcosCustomData188       = "dcos/dcoscustomdata188.t"
	dcosCustomData190       = "dcos/dcoscustomdata190.t"
	dcosCustomData198       = "dcos/dcoscustomdata198.t"
	dcosCustomData110       = "dcos/dcoscustomdata110.t"
	dcosProvision           = "dcos/dcosprovision.sh"
	dcosWindowsProvision    = "dcos/dcosWindowsProvision.ps1"
	dcosProvisionSource     = "dcos/dcosprovisionsource.sh"
	dcos2Provision          = "dcos/bstrap/dcosprovision.sh"
	dcos2BootstrapProvision = "dcos/bstrap/bootstrapprovision.sh"
	dcos2CustomData1110     = "dcos/bstrap/dcos1.11.0.customdata.t"
	dcos2CustomData1112     = "dcos/bstrap/dcos1.11.2.customdata.t"
)

const (
	swarmProvision            = "swarm/configure-swarm-cluster.sh"
	swarmWindowsProvision     = "swarm/Install-ContainerHost-And-Join-Swarm.ps1"
	swarmModeProvision        = "swarm/configure-swarmmode-cluster.sh"
	swarmModeWindowsProvision = "swarm/Join-SwarmMode-cluster.ps1"
)

const (
	agentOutputs                  = "agentoutputs.t"
	agentParams                   = "agentparams.t"
	classicParams                 = "classicparams.t"
	dcosAgentResourcesVMAS        = "dcos/dcosagentresourcesvmas.t"
	dcosWindowsAgentResourcesVMAS = "dcos/dcosWindowsAgentResourcesVmas.t"
	dcosAgentResourcesVMSS        = "dcos/dcosagentresourcesvmss.t"
	dcosWindowsAgentResourcesVMSS = "dcos/dcosWindowsAgentResourcesVmss.t"
	dcosAgentVars                 = "dcos/dcosagentvars.t"
	dcosBaseFile                  = "dcos/dcosbase.t"
	dcosParams                    = "dcos/dcosparams.t"
	dcosMasterResources           = "dcos/dcosmasterresources.t"
	dcosMasterVars                = "dcos/dcosmastervars.t"
	dcos2BaseFile                 = "dcos/bstrap/dcosbase.t"
	dcos2BootstrapVars            = "dcos/bstrap/bootstrapvars.t"
	dcos2BootstrapParams          = "dcos/bstrap/bootstrapparams.t"
	dcos2BootstrapResources       = "dcos/bstrap/bootstrapresources.t"
	dcos2BootstrapCustomdata      = "dcos/bstrap/bootstrapcustomdata.yml"
	dcos2MasterVars               = "dcos/bstrap/dcosmastervars.t"
	dcos2MasterResources          = "dcos/bstrap/dcosmasterresources.t"
	iaasOutputs                   = "iaasoutputs.t"
	kubernetesBaseFile            = "k8s/kubernetesbase.t"
	kubernetesAgentResourcesVMAS  = "k8s/kubernetesagentresourcesvmas.t"
	kubernetesAgentResourcesVMSS  = "k8s/kubernetesagentresourcesvmss.t"
	kubernetesAgentVars           = "k8s/kubernetesagentvars.t"
	kubernetesMasterResources     = "k8s/kubernetesmasterresources.t"
	kubernetesMasterVars          = "k8s/kubernetesmastervars.t"
	kubernetesParams              = "k8s/kubernetesparams.t"
	kubernetesWinAgentVars        = "k8s/kuberneteswinagentresourcesvmas.t"
	kubernetesWinAgentVarsVMSS    = "k8s/kuberneteswinagentresourcesvmss.t"
	masterOutputs                 = "masteroutputs.t"
	masterParams                  = "masterparams.t"
	openshiftInfraResources       = "openshift/infraresources.t"
	swarmBaseFile                 = "swarm/swarmbase.t"
	swarmParams                   = "swarm/swarmparams.t"
	swarmAgentResourcesVMAS       = "swarm/swarmagentresourcesvmas.t"
	swarmAgentResourcesVMSS       = "swarm/swarmagentresourcesvmss.t"
	swarmAgentResourcesClassic    = "swarm/swarmagentresourcesclassic.t"
	swarmAgentVars                = "swarm/swarmagentvars.t"
	swarmMasterResources          = "swarm/swarmmasterresources.t"
	swarmMasterVars               = "swarm/swarmmastervars.t"
	swarmWinAgentResourcesVMAS    = "swarm/swarmwinagentresourcesvmas.t"
	swarmWinAgentResourcesVMSS    = "swarm/swarmwinagentresourcesvmss.t"
	windowsParams                 = "windowsparams.t"
)

const (
	azurePublicCloud       = "AzurePublicCloud"
	azureChinaCloud        = "AzureChinaCloud"
	azureGermanCloud       = "AzureGermanCloud"
	azureUSGovernmentCloud = "AzureUSGovernmentCloud"
)

package acsengine

const (
	// DefaultMasterSubnet specifies the default master subnet for DCOS or Swarm
	DefaultMasterSubnet = "172.16.0.0/24"
	// DefaultFirstConsecutiveStaticIP specifies the static IP address on master 0 for DCOS or Swarm
	DefaultFirstConsecutiveStaticIP = "172.16.0.5"
	// DefaultSwarmWindowsMasterSubnet specifies the default master subnet for a Swarm Windows cluster
	DefaultSwarmWindowsMasterSubnet = "192.168.255.0/24"
	// DefaultSwarmWindowsFirstConsecutiveStaticIP specifies the static IP address on master 0 for a Swarm WIndows cluster
	DefaultSwarmWindowsFirstConsecutiveStaticIP = "192.168.255.5"
	// DefaultKubernetesMasterSubnet specifies the default subnet for masters and agents.
	DefaultKubernetesMasterSubnet = "10.240.0.0/16"
	// DefaultKubernetesClusterSubnet specifies the default subnet for pods.
	DefaultKubernetesClusterSubnet = "10.244.0.0/16"
	// DefaultDockerBridgeSubnet specifies the default subnet for the docker bridge network for masters and agents.
	DefaultDockerBridgeSubnet = "172.17.0.1/16"
	// DefaultNonMasqueradeCidr specifies the subnet that should not be masqueraded on host
	DefaultNonMasqueradeCidr = "10.0.0.0/8"
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
	// NetworkPolicyNone is the string expression for no network policy
	NetworkPolicyNone = "none"
	// NetworkPluginKubenet is the string expression for kubenet network plugin
	NetworkPluginKubenet = "kubenet"
	// DefaultNetworkPolicy defines the network policy to use by default
	DefaultNetworkPolicy = NetworkPolicyNone
	// DefaultNetworkPolicyWindows defines the network policy to use by default for clusters with Windows agent pools
	DefaultNetworkPolicyWindows = NetworkPolicyNone
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
	// DefaultACIConnectorAddonName is the name of the tiller addon deployment
	DefaultACIConnectorAddonName = "aci-connector"
	// DefaultDashboardAddonName is the name of the kubernetes-dashboard addon deployment
	DefaultDashboardAddonName = "kubernetes-dashboard"
	// DefaultTillerImage defines the Helm Tiller deployment version on Kubernetes Clusters
	// TODO deprecate this usage, we should be favoring a more frequent upgrade cycle that pins fresh tiller versions to specific k8s versions
	DefaultTillerImage = "tiller:v2.6.2"
	// DefaultACIConnectorImage defines the ACI Connector deployment version on Kubernetes Clusters
	DefaultACIConnectorImage = "virtual-kubelet:latest"
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
	// DefaultEtcdVersion specifies the default etcd version to install
	DefaultEtcdVersion = "3.2.11"
	// DefaultEtcdDiskSize specifies the default size for Kubernetes master etcd disk volumes in GB
	DefaultEtcdDiskSize = "128"
	// DefaultReschedulerImage defines the rescheduler deployment version on Kubernetes Clusters
	DefaultReschedulerImage = "rescheduler:v0.3.1"
	// DefaultReschedulerAddonName is the name of the rescheduler addon deployment
	DefaultReschedulerAddonName = "rescheduler"
	// DefaultKubernetesKubeletMaxPods is the max pods per kubelet
	DefaultKubernetesKubeletMaxPods = 110
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

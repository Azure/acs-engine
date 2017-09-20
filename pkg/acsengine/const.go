package acsengine

import (
	"strconv"

	"github.com/Azure/acs-engine/pkg/api"
)

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
	// DefaultNetworkPolicy is disabling network policy enforcement
	DefaultNetworkPolicy = "none"
	// DefaultKubernetesNodeStatusUpdateFrequency is 10s, see --node-status-update-frequency at https://kubernetes.io/docs/admin/kubelet/
	DefaultKubernetesNodeStatusUpdateFrequency = "10s"
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
	// DefaultTillerImage defines the Helm Tiller deployment version on Kubernetes Clusters
	DefaultTillerImage = "tiller:v2.6.0"
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
)

const (
	// DCOSMaster represents the master node type
	DCOSMaster DCOSNodeType = "DCOSMaster"
	// DCOSPrivateAgent represents the private agent node type
	DCOSPrivateAgent DCOSNodeType = "DCOSPrivateAgent"
	// DCOSPublicAgent represents the public agent node type
	DCOSPublicAgent DCOSNodeType = "DCOSPublicAgent"
)

// KubeConfigs represents Docker images used for Kubernetes components based on Kubernetes releases (major.minor)
// For instance, Kubernetes release "1.7" would contain the version "1.7.2"
var KubeConfigs = map[string]map[string]string{
	api.KubernetesRelease1Dot7: {
		"hyperkube":       "hyperkube-amd64:v1.7.5",
		"dashboard":       "kubernetes-dashboard-amd64:v1.6.3",
		"exechealthz":     "exechealthz-amd64:1.2",
		"addonresizer":    "addon-resizer:1.7",
		"heapster":        "heapster-amd64:v1.4.2",
		"dns":             "k8s-dns-kube-dns-amd64:1.14.4",
		"addonmanager":    "kube-addon-manager-amd64:v6.4-beta.2",
		"dnsmasq":         "k8s-dns-dnsmasq-nanny-amd64:1.14.4",
		"pause":           "pause-amd64:3.0",
		"tiller":          DefaultTillerImage,
		"windowszip":      "v1.7.5-3intwinnat.zip",
		"nodestatusfreq":  DefaultKubernetesNodeStatusUpdateFrequency,
		"nodegraceperiod": DefaultKubernetesCtrlMgrNodeMonitorGracePeriod,
		"podeviction":     DefaultKubernetesCtrlMgrPodEvictionTimeout,
		"routeperiod":     DefaultKubernetesCtrlMgrRouteReconciliationPeriod,
		"backoffretries":  strconv.Itoa(DefaultKubernetesCloudProviderBackoffRetries),
		"backoffjitter":   strconv.FormatFloat(DefaultKubernetesCloudProviderBackoffJitter, 'f', -1, 64),
		"backoffduration": strconv.Itoa(DefaultKubernetesCloudProviderBackoffDuration),
		"backoffexponent": strconv.FormatFloat(DefaultKubernetesCloudProviderBackoffExponent, 'f', -1, 64),
		"ratelimitqps":    strconv.FormatFloat(DefaultKubernetesCloudProviderRateLimitQPS, 'f', -1, 64),
		"ratelimitbucket": strconv.Itoa(DefaultKubernetesCloudProviderRateLimitBucket),
		"gchighthreshold": strconv.Itoa(DefaultKubernetesGCHighThreshold),
		"gclowthreshold":  strconv.Itoa(DefaultKubernetesGCLowThreshold),
	},
	api.KubernetesRelease1Dot6: {
		"hyperkube":       "hyperkube-amd64:v1.6.9",
		"dashboard":       "kubernetes-dashboard-amd64:v1.6.3",
		"exechealthz":     "exechealthz-amd64:1.2",
		"addonresizer":    "addon-resizer:1.7",
		"heapster":        "heapster-amd64:v1.3.0",
		"dns":             "k8s-dns-kube-dns-amd64:1.14.4",
		"addonmanager":    "kube-addon-manager-amd64:v6.4-beta.2",
		"dnsmasq":         "k8s-dns-dnsmasq-nanny-amd64:1.14.4",
		"pause":           "pause-amd64:3.0",
		"tiller":          DefaultTillerImage,
		"windowszip":      "v1.6.9-3intwinnat.zip",
		"nodestatusfreq":  DefaultKubernetesNodeStatusUpdateFrequency,
		"nodegraceperiod": DefaultKubernetesCtrlMgrNodeMonitorGracePeriod,
		"podeviction":     DefaultKubernetesCtrlMgrPodEvictionTimeout,
		"routeperiod":     DefaultKubernetesCtrlMgrRouteReconciliationPeriod,
		"backoffretries":  strconv.Itoa(DefaultKubernetesCloudProviderBackoffRetries),
		"backoffjitter":   strconv.FormatFloat(DefaultKubernetesCloudProviderBackoffJitter, 'f', -1, 64),
		"backoffduration": strconv.Itoa(DefaultKubernetesCloudProviderBackoffDuration),
		"backoffexponent": strconv.FormatFloat(DefaultKubernetesCloudProviderBackoffExponent, 'f', -1, 64),
		"ratelimitqps":    strconv.FormatFloat(DefaultKubernetesCloudProviderRateLimitQPS, 'f', -1, 64),
		"ratelimitbucket": strconv.Itoa(DefaultKubernetesCloudProviderRateLimitBucket),
		"gchighthreshold": strconv.Itoa(DefaultKubernetesGCHighThreshold),
		"gclowthreshold":  strconv.Itoa(DefaultKubernetesGCLowThreshold),
	},
	api.KubernetesRelease1Dot5: {
		"hyperkube":       "hyperkube-amd64:v1.5.7",
		"dashboard":       "kubernetes-dashboard-amd64:v1.5.1",
		"exechealthz":     "exechealthz-amd64:1.2",
		"addonresizer":    "addon-resizer:1.6",
		"heapster":        "heapster:v1.2.0",
		"dns":             "kubedns-amd64:1.7",
		"addonmanager":    "kube-addon-manager-amd64:v6.4-beta.2",
		"dnsmasq":         "kube-dnsmasq-amd64:1.3",
		"pause":           "pause-amd64:3.0",
		"tiller":          "tiller:v2.5.1",
		"windowszip":      "v1.5.7intwinnat.zip",
		"nodestatusfreq":  DefaultKubernetesNodeStatusUpdateFrequency,
		"nodegraceperiod": DefaultKubernetesCtrlMgrNodeMonitorGracePeriod,
		"podeviction":     DefaultKubernetesCtrlMgrPodEvictionTimeout,
		"routeperiod":     DefaultKubernetesCtrlMgrRouteReconciliationPeriod,
		"gchighthreshold": strconv.Itoa(DefaultKubernetesGCHighThreshold),
		"gclowthreshold":  strconv.Itoa(DefaultKubernetesGCLowThreshold),
	},
}

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

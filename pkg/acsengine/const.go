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
	// DefaultAgentIPAddressCount is the default number of IP addresses per network interface on agents
	DefaultAgentIPAddressCount = 1
	// DefaultAgentMultiIPAddressCount is the default number of IP addresses per network interface on agents,
	// when VNET integration is enabled. It can be overridden per pool by setting the pool's IPAdddressCount property.
	DefaultAgentMultiIPAddressCount = 128
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
)

const (
	// DCOSMaster represents the master node type
	DCOSMaster DCOSNodeType = "DCOSMaster"
	// DCOSPrivateAgent represents the private agent node type
	DCOSPrivateAgent DCOSNodeType = "DCOSPrivateAgent"
	// DCOSPublicAgent represents the public agent node type
	DCOSPublicAgent DCOSNodeType = "DCOSPublicAgent"
)

// KubeImages represents Docker images used for Kubernetes components based on Kubernetes version
var KubeImages = map[api.OrchestratorVersion]map[string]string{
	api.Kubernetes170: {
		"hyperkube":    "hyperkube-amd64:v1.7.0",
		"dashboard":    "kubernetes-dashboard-amd64:v1.6.1",
		"exechealthz":  "exechealthz-amd64:1.2",
		"addonresizer": "addon-resizer:2.0",
		"heapster":     "heapster:v1.4.0",
		"dns":          "k8s-dns-kube-dns-amd64:1.14.4",
		"addonmanager": "kube-addon-manager-amd64:v6.1.2",
		"dnsmasq":      "k8s-dns-dnsmasq-amd64:1.14.4",
		"pause":        "pause-amd64:3.0",
		"windowszip":   "v1.7.0intwinnat.zip",
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
	},
	api.Kubernetes166: {
		"hyperkube":       "hyperkube-amd64:v1.6.6",
		"dashboard":       "kubernetes-dashboard-amd64:v1.6.1",
		"exechealthz":     "exechealthz-amd64:1.2",
		"addonresizer":    "addon-resizer:1.7",
		"heapster":        "heapster:v1.3.0",
		"dns":             "k8s-dns-kube-dns-amd64:1.14.4",
		"addonmanager":    "kube-addon-manager-amd64:v6.4-beta.2",
		"dnsmasq":         "k8s-dns-dnsmasq-amd64:1.13.0",
		"pause":           "pause-amd64:3.0",
		"windowszip":      "v1.6.6intwinnat.zip",
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
	},
	api.Kubernetes162: {
		"hyperkube":       "hyperkube-amd64:v1.6.2",
		"dashboard":       "kubernetes-dashboard-amd64:v1.6.0",
		"exechealthz":     "exechealthz-amd64:1.2",
		"addonresizer":    "addon-resizer:1.6",
		"heapster":        "heapster:v1.2.0",
		"dns":             "k8s-dns-kube-dns-amd64:1.13.0",
		"addonmanager":    "kube-addon-manager-amd64:v6.4",
		"dnsmasq":         "k8s-dns-dnsmasq-amd64:1.13.0",
		"pause":           "pause-amd64:3.0",
		"windowszip":      "v1.6.2intwinnat.zip",
		"nodestatusfreq":  DefaultKubernetesNodeStatusUpdateFrequency,
		"nodegraceperiod": DefaultKubernetesCtrlMgrNodeMonitorGracePeriod,
		"podeviction":     DefaultKubernetesCtrlMgrPodEvictionTimeout,
		"routeperiod":     DefaultKubernetesCtrlMgrRouteReconciliationPeriod,
	},

	api.Kubernetes160: {
		"hyperkube":       "hyperkube-amd64:v1.6.0",
		"dashboard":       "kubernetes-dashboard-amd64:v1.6.0",
		"exechealthz":     "exechealthz-amd64:1.2",
		"addonresizer":    "addon-resizer:1.6",
		"heapster":        "heapster:v1.2.0",
		"dns":             "k8s-dns-kube-dns-amd64:1.13.0",
		"addonmanager":    "kube-addon-manager-amd64:v6.4",
		"dnsmasq":         "k8s-dns-dnsmasq-amd64:1.13.0",
		"pause":           "pause-amd64:3.0",
		"windowszip":      "v1.6.0intwinnat.zip",
		"nodestatusfreq":  DefaultKubernetesNodeStatusUpdateFrequency,
		"nodegraceperiod": DefaultKubernetesCtrlMgrNodeMonitorGracePeriod,
		"podeviction":     DefaultKubernetesCtrlMgrPodEvictionTimeout,
		"routeperiod":     DefaultKubernetesCtrlMgrRouteReconciliationPeriod,
	},

	api.Kubernetes157: {
		"hyperkube":       "hyperkube-amd64:v1.5.7",
		"dashboard":       "kubernetes-dashboard-amd64:v1.5.1",
		"exechealthz":     "exechealthz-amd64:1.2",
		"addonresizer":    "addon-resizer:1.6",
		"heapster":        "heapster:v1.2.0",
		"dns":             "kubedns-amd64:1.7",
		"addonmanager":    "kube-addon-manager-amd64:v6.2",
		"dnsmasq":         "kube-dnsmasq-amd64:1.3",
		"pause":           "pause-amd64:3.0",
		"windowszip":      "v1.5.7intwinnat.zip",
		"nodestatusfreq":  DefaultKubernetesNodeStatusUpdateFrequency,
		"nodegraceperiod": DefaultKubernetesCtrlMgrNodeMonitorGracePeriod,
		"podeviction":     DefaultKubernetesCtrlMgrPodEvictionTimeout,
		"routeperiod":     DefaultKubernetesCtrlMgrRouteReconciliationPeriod,
	},

	api.Kubernetes153: {
		"hyperkube":       "hyperkube-amd64:v1.5.3",
		"dashboard":       "kubernetes-dashboard-amd64:v1.5.1",
		"exechealthz":     "exechealthz-amd64:1.2",
		"addonresizer":    "addon-resizer:1.6",
		"heapster":        "heapster:v1.2.0",
		"dns":             "kubedns-amd64:1.7",
		"addonmanager":    "kube-addon-manager-amd64:v6.2",
		"dnsmasq":         "kube-dnsmasq-amd64:1.3",
		"pause":           "pause-amd64:3.0",
		"windowszip":      "v1.5.3intwinnat.zip",
		"nodestatusfreq":  DefaultKubernetesNodeStatusUpdateFrequency,
		"nodegraceperiod": DefaultKubernetesCtrlMgrNodeMonitorGracePeriod,
		"podeviction":     DefaultKubernetesCtrlMgrPodEvictionTimeout,
		"routeperiod":     DefaultKubernetesCtrlMgrRouteReconciliationPeriod,
	},
}

const (
	//MsecndDCOSBootstrapDownloadURL Azure CDN to download DCOS1.7.3
	MsecndDCOSBootstrapDownloadURL = "https://az837203.vo.msecnd.net/dcos/%s/bootstrap/%s.bootstrap.tar.xz"
	//AzureEdgeDCOSBootstrapDownloadURL is the azure edge CDN download url
	AzureEdgeDCOSBootstrapDownloadURL = "https://dcosio.azureedge.net/dcos/%s/bootstrap/%s.bootstrap.tar.xz"
	//AzureChinaCloudDCOSBootstrapDownloadURL is the China specific DCOS package download url.
	AzureChinaCloudDCOSBootstrapDownloadURL = "https://acsengine.blob.core.chinacloudapi.cn/dcos/%s.bootstrap.tar.xz"
)

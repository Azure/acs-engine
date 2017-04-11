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
	// DefaultKubernetesMasterSubnet specifies the default kubernetes master subnet
	DefaultKubernetesMasterSubnet = "10.240.0.0/16"
	// DefaultKubernetesSubnet specifies the default Kubernetes subnet when VNET integration is enabled.
	DefaultKubernetesSubnet = "10.240.0.0/12"
	// DefaultFirstConsecutiveKubernetesStaticIP specifies the static IP address on Kubernetes master 0
	DefaultFirstConsecutiveKubernetesStaticIP = "10.240.255.5"
	// DefaultAgentSubnetTemplate specifies a default agent subnet
	DefaultAgentSubnetTemplate = "10.%d.0.0/16"
	// DefaultAgentIPAddressCount is the default number of IP addresses per network interface on agents
	DefaultAgentIPAddressCount = 1
	// DefaultAgentMultiIPAddressCount is the default number of IP addresses per network interface on agents,
	// when VNET integration is enabled. It can be overriden per pool by setting the pool's IPAdddressCount property.
	DefaultAgentMultiIPAddressCount = 128
	// DefaultKubernetesClusterDomain is the dns suffix used in the cluster (used as a SAN in the PKI generation)
	DefaultKubernetesClusterDomain = "cluster.local"
	// DefaultInternalLbStaticIPOffset specifies the offset of the internal LoadBalancer's IP
	// address relative to the first consecutive Kubernetes static IP
	DefaultInternalLbStaticIPOffset = 10
	// DefaultNetworkPolicy is disabling network policy enforcement
	DefaultNetworkPolicy = "none"
	// DefaultKubernetesDnsServiceIP specifies the IP address that kube-dns
	// listens on by default. must by in the default Service CIDR range.
	DefaultKubernetesDnsServiceIP = "10.0.0.10"
	// DefaultKubernetesServiceCIDR specifies the IP subnet that kubernetes will
	// create Service IPs within.
	DefaultKubernetesServiceCIDR = "10.0.0.0/16"
	// DefaultKubernetesClusterCIDR specifies the subnet that Pod IPs will be
	// assigned from.
	DefaultKubernetesClusterCIDR = "10.244.0.0/16"
)

const (
	// Master represents the master node type
	DCOSMaster DCOSNodeType = "DCOSMaster"
	// PrivateAgent represents the private agent node type
	DCOSPrivateAgent DCOSNodeType = "DCOSPrivateAgent"
	// PublicAgent represents the public agent node type
	DCOSPublicAgent DCOSNodeType = "DCOSPublicAgent"
)

const (
	KubernetesHyperkubeImageName         = "hyperkube-amd64:v1.5.3"
	KubernetesDashboardImageName         = "kubernetes-dashboard-amd64:v1.5.1"
	KubernetesExechealthzImageName       = "exechealthz-amd64:1.2"
	KubernetesAddonResizerImageName      = "addon-resizer:1.6"
	KubernetesHeapsterImageName          = "heapster:v1.2.0"
	KubernetesDNSImageName               = "kubedns-amd64:1.7"
	KubernetesAddonManagerImageName      = "kube-addon-manager-amd64:v6.2"
	KubernetesDNSMasqImageName           = "kube-dnsmasq-amd64:1.3"
	KubernetesPodInfraContainerImageName = "pause-amd64:3.0"
)

const (
	//MsecndDCOSBootstrapDownloadURL Azure CDN to download DCOS1.7.3
	MsecndDCOSBootstrapDownloadURL = "https://az837203.vo.msecnd.net/dcos/%s/bootstrap/%s.bootstrap.tar.xz"
	//AzureEdgeDCOSBootstrapDownloadURL is the azure edge CDN download url
	AzureEdgeDCOSBootstrapDownloadURL = "https://dcosio.azureedge.net/dcos/%s/bootstrap/%s.bootstrap.tar.xz"
	//AzureChinaCloudDCOSBootstrapDownloadURL is the China specific DCOS package download url.
	AzureChinaCloudDCOSBootstrapDownloadURL = "https://acsengine.blob.core.chinacloudapi.cn/dcos/%s.bootstrap.tar.xz"
)

package acsengine

//DockerSpecConfig is default docker install script URL.
type DockerSpecConfig struct {
	DefaultDockerInstallScriptURL string
}

//KubernetesSpecConfig is the kubernetes container images used.
type KubernetesSpecConfig struct {
	DefaultKubernetesHyperkubeSpec         string
	DefaultKubernetesDashboardSpec         string
	DefaultKubernetesExechealthzSpec       string
	DefaultKubernetesAddonResizerSpec      string
	DefaultKubernetesHeapsterSpec          string
	DefaultKubernetesDNSSpec               string
	DefaultKubernetesAddonManagerSpec      string
	DefaultKubernetesDNSMasqSpec           string
	DefaultKubernetesPodInfraContainerSpec string
	DefaultKubectlDownloadURL              string
}

//AzureEnvironmentSpecConfig is the overall configuration differences in different cloud environments.
type AzureEnvironmentSpecConfig struct {
	DockerSpecConfig     DockerSpecConfig
	KubernetesSpecConfig KubernetesSpecConfig
}

var (
	AzurePublicCloud = AzureEnvironmentSpecConfig{
		//DockerConfigAzurePublicCloud specify the default script location of docker installer script
		DockerSpecConfig: DockerSpecConfig{
			DefaultDockerInstallScriptURL: "https://get.docker.com/",
		},
		//KubeConfigAzurePublicCloud is the default kubernetes container image url.
		KubernetesSpecConfig: KubernetesSpecConfig{
			DefaultKubernetesHyperkubeSpec:         "gcr.io/google_containers/hyperkube-amd64:v1.4.6",
			DefaultKubernetesDashboardSpec:         "gcr.io/google_containers/kubernetes-dashboard-amd64:v1.5.0",
			DefaultKubernetesExechealthzSpec:       "gcr.io/google_containers/exechealthz-amd64:1.2",
			DefaultKubernetesAddonResizerSpec:      "gcr.io/google_containers/addon-resizer:1.6",
			DefaultKubernetesHeapsterSpec:          "gcr.io/google_containers/heapster:v1.2.0",
			DefaultKubernetesDNSSpec:               "gcr.io/google_containers/kubedns-amd64:1.7",
			DefaultKubernetesAddonManagerSpec:      "gcr.io/google_containers/kube-addon-manager-amd64:v5.1",
			DefaultKubernetesDNSMasqSpec:           "gcr.io/google_containers/kube-dnsmasq-amd64:1.3",
			DefaultKubernetesPodInfraContainerSpec: "gcr.io/google_containers/pause-amd64:3.0",
			DefaultKubectlDownloadURL:              "https://storage.googleapis.com/kubernetes-release/release/" + DefaultKubectlVersion + "/bin/linux/amd64/kubectl",
		},
	}

	AzureChinaCloud = AzureEnvironmentSpecConfig{
		//DockerConfigAzureChinaCloud specify the docker install script download URL in China.
		DockerSpecConfig: DockerSpecConfig{
			DefaultDockerInstallScriptURL: "https://acsengine.blob.core.chinacloudapi.cn/docker/install-docker",
		},
		//KubeConfigAzureChinaCloud - Due to Chinese firewall issue, the default containers from google is blocked, use the Chinese local mirror instead
		KubernetesSpecConfig: KubernetesSpecConfig{
			DefaultKubernetesHyperkubeSpec:         "mirror.azure.cn:5000/google_containers/hyperkube-amd64:v1.4.6",
			DefaultKubernetesDashboardSpec:         "mirror.azure.cn:5000/google_containers/kubernetes-dashboard-amd64:v1.5.0",
			DefaultKubernetesExechealthzSpec:       "mirror.azure.cn:5000/google_containers/exechealthz-amd64:1.2",
			DefaultKubernetesAddonResizerSpec:      "mirror.azure.cn:5000/google_containers/addon-resizer:1.6",
			DefaultKubernetesHeapsterSpec:          "mirror.azure.cn:5000/google_containers/heapster:v1.2.0",
			DefaultKubernetesDNSSpec:               "mirror.azure.cn:5000/google_containers/kubedns-amd64:1.7",
			DefaultKubernetesAddonManagerSpec:      "mirror.azure.cn:5000/google_containers/kube-addon-manager-amd64:v5.1",
			DefaultKubernetesDNSMasqSpec:           "mirror.azure.cn:5000/google_containers/kube-dnsmasq-amd64:1.3",
			DefaultKubernetesPodInfraContainerSpec: "mirror.azure.cn:5000/google_containers/pause-amd64:3.0",
			DefaultKubectlDownloadURL:              "https://acsengine.blob.core.chinacloudapi.cn/kubernetes/kubectl",
		},
	}
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
	// DefaultKubernetesMasterSubnet specifies the default kubernetes master subnet
	DefaultKubernetesMasterSubnet = "10.240.0.0/16"
	// DefaultFirstConsecutiveKubernetesStaticIP specifies the static IP address on Kubernetes master 0
	DefaultFirstConsecutiveKubernetesStaticIP = "10.240.255.5"
	// DefaultAgentSubnetTemplate specifies a default agent subnet
	DefaultAgentSubnetTemplate = "10.%d.0.0/16"
	// DefaultKubernetesClusterDomain is the dns suffix used in the cluster (used as a SAN in the PKI generation)
	DefaultKubernetesClusterDomain = "cluster.local"
	//DefaultKubectlVersion is the default version of kubectl
	DefaultKubectlVersion = "v1.4.6"
)

const (
	// Master represents the master node type
	DCOSMaster DCOSNodeType = "DCOSMaster"
	// PrivateAgent represents the private agent node type
	DCOSPrivateAgent DCOSNodeType = "DCOSPrivateAgent"
	// PublicAgent represents the public agent node type
	DCOSPublicAgent DCOSNodeType = "DCOSPublicAgent"
)

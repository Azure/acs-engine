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
	// DefaultFirstConsecutiveKubernetesStaticIP specifies the static IP address on Kubernetes master 0
	DefaultFirstConsecutiveKubernetesStaticIP = "10.240.255.5"
	// DefaultAgentSubnetTemplate specifies a default agent subnet
	DefaultAgentSubnetTemplate = "10.%d.0.0/16"
	// DefaultKubernetesClusterDomain is the dns suffix used in the cluster (used as a SAN in the PKI generation)
	DefaultKubernetesClusterDomain = "cluster.local"
	// DefaultKubernetesHyperkubeSpec is the default version used for Kubernetes setup
	// The latest stable version can be found here: https://storage.googleapis.com/kubernetes-release/release/stable.txt
	DefaultKubernetesHyperkubeSpec = "gcr.io/google_containers/hyperkube-amd64:v1.5.1"
	// DefaultKubectlVersion is the version used for kubectl
	// The latest stable version can be found here: https://storage.googleapis.com/kubernetes-release/release/stable.txt
	DefaultKubectlVersion = "v1.5.1"
)

const (
	// Master represents the master node type
	DCOSMaster DCOSNodeType = "DCOSMaster"
	// PrivateAgent represents the private agent node type
	DCOSPrivateAgent DCOSNodeType = "DCOSPrivateAgent"
	// PublicAgent represents the public agent node type
	DCOSPublicAgent DCOSNodeType = "DCOSPublicAgent"
)

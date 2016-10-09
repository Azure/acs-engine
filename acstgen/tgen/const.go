package tgen

const (
	// BaseLBPriority specifies the base lb priority.
	BaseLBPriority = 200
	// DefaultMasterSubnet specifies the default master subnet for DCOS or Swarm
	DefaultMasterSubnet = "172.16.0.0/24"
	// DefaultFirstConsecutiveStaticIP specifies the static IP address on master 0 for DCOS or Swarm
	DefaultFirstConsecutiveStaticIP = "172.16.0.5"
	// DefaultKubernetesMasterSubnet specifies the default kubernetes master subnet
	DefaultKubernetesMasterSubnet = "10.240.0.0/16"
	// DefaultFirstConsecutiveKubernetesStaticIP specifies the static IP address on Kubernetes master 0
	DefaultFirstConsecutiveKubernetesStaticIP = "10.240.255.5"
	// DefaultAgentSubnetTemplate specifies a default agent subnet
	DefaultAgentSubnetTemplate = "10.%d.0.0/24"
	// DefaultKubernetesClusterDomain is the dns suffix used in the cluster (used as a SAN in the PKI generation)
	DefaultKubernetesClusterDomain = "cluster.local"
	// KubernetesHyperkubeSpec is the hyperkube version used for Kubernetes setup
	// The latest stable version can be found here: https://storage.googleapis.com/kubernetes-release/release/stable.txt
	KubernetesHyperkubeSpec = "gcr.io/google_containers/hyperkube-amd64:v1.4.0"
)

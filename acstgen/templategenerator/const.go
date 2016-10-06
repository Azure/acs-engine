package templategenerator

const (
	// BaseLBPriority specifies the base lb priority.
	BaseLBPriority = 200
	// DefaultMasterSubnet specifies the default master subnet
	DefaultMasterSubnet = "172.16.0.0/24"
	// DefaultFirstConsecutiveStaticIP specifies the static IP address on master 0
	DefaultFirstConsecutiveStaticIP = "172.16.0.5"
	// DefaultKubernetesMasterSubnet specifies the default kubernetes master subnet
	DefaultKubernetesMasterSubnet = "10.240.0.0/16"
	// DefaultFirstConsecutiveKubernetesStaticIP specifies the static IP address on Kubernetes master 0
	DefaultFirstConsecutiveKubernetesStaticIP = "10.240.255.5"
	// DefaultAgentSubnetTemplate specifies a default agent subnet
	DefaultAgentSubnetTemplate = "10.%d.0.0/24"
	// DefaultKubernetesClusterDomain is the dns suffix used in the cluster (used as a SAN in the PKI generation)
	DefaultKubernetesClusterDomain = "cluster.local"
)

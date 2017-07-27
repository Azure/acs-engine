package api

// the orchestrators supported by vlabs
const (
	// Mesos is the string constant for MESOS orchestrator type
	Mesos string = "Mesos"
	// DCOS is the string constant for DCOS orchestrator type and defaults to DCOS188
	DCOS string = "DCOS"
	// Swarm is the string constant for the Swarm orchestrator type
	Swarm string = "Swarm"
	// Kubernetes is the string constant for the Kubernetes orchestrator type
	Kubernetes string = "Kubernetes"
	// SwarmMode is the string constant for the Swarm Mode orchestrator type
	SwarmMode string = "SwarmMode"
)

// the OSTypes supported by vlabs
const (
	Windows OSType = "Windows"
	Linux   OSType = "Linux"
)

// validation values
const (
	// MinAgentCount are the minimum number of agents per agent pool
	MinAgentCount = 1
	// MaxAgentCount are the maximum number of agents per agent pool
	MaxAgentCount = 100
	// MinPort specifies the minimum tcp port to open
	MinPort = 1
	// MaxPort specifies the maximum tcp port to open
	MaxPort = 65535
	// MaxDisks specifies the maximum attached disks to add to the cluster
	MaxDisks = 4
)

// Availability profiles
const (
	// AvailabilitySet means that the vms are in an availability set
	AvailabilitySet = "AvailabilitySet"
	// VirtualMachineScaleSets means that the vms are in a virtual machine scaleset
	VirtualMachineScaleSets = "VirtualMachineScaleSets"
)

// storage profiles
const (
	// StorageAccount means that the nodes use raw storage accounts for their os and attached volumes
	StorageAccount = "StorageAccount"
	// ManagedDisks means that the nodes use managed disks for their os and attached volumes
	ManagedDisks = "ManagedDisks"
)

const (
	// Kubernetes153 is the string constant for Kubernetes 1.5.3
	Kubernetes153 string = "1.5.3"
	// Kubernetes157 is the string constant for Kubernetes 1.5.7
	Kubernetes157 string = "1.5.7"
	// Kubernetes160 is the string constant for Kubernetes 1.6.0
	Kubernetes160 string = "1.6.0"
	// Kubernetes162 is the string constant for Kubernetes 1.6.2
	Kubernetes162 string = "1.6.2"
	// Kubernetes166 is the string constant for Kubernetes 1.6.6
	Kubernetes166 string = "1.6.6"
	// Kubernetes170 is the string constant for Kubernetes 1.7.0
	Kubernetes170 string = "1.7.0"
	// Kubernetes171 is the string constant for Kubernetes 1.7.1
	Kubernetes171 string = "1.7.1"
	// KubernetesDefaultVersion is the string constant for current Kubernetes version
	KubernetesDefaultVersion string = Kubernetes166
)

const (
	// DCOS190 is the string constant for DCOS 1.9.0
	DCOS190 string = "1.9.0"
	// DCOS188 is the string constant for DCOS 1.8.8
	DCOS188 string = "1.8.8"
	// DCOS187 is the string constant for DCOS 1.8.7
	DCOS187 string = "1.8.7"
	// DCOS184 is the string constant for DCOS 1.8.4
	DCOS184 string = "1.8.4"
	// DCOS173 is the string constant for DCOS 1.7.3
	DCOS173 string = "1.7.3"
	// DCOSLatest is the string constant for latest DCOS version
	DCOSLatest string = DCOS190
)

// To identify programmatically generated public agent pools
const publicAgentPoolSuffix = "-public"

package api

// the orchestrators supported by vlabs
const (
	// Mesos is the string constant for MESOS orchestrator type
	Mesos OrchestratorType = "Mesos"
	// DCOS is the string constant for DCOS orchestrator type and defaults to DCOS187
	DCOS OrchestratorType = "DCOS"
	// DCOS187 is the string constant for DCOS 1.8.7 orchestrator type
	DCOS187 OrchestratorType = "DCOS187"
	// DCOS184 is the string constant for DCOS 1.8.4 orchestrator type
	DCOS184 OrchestratorType = "DCOS184"
	// DCOS173 is the string constant for DCOS 1.7.3 orchestrator type
	DCOS173 OrchestratorType = "DCOS173"
	// Swarm is the string constant for the Swarm orchestrator type
	Swarm OrchestratorType = "Swarm"
	// Kubernetes is the string constant for the Kubernetes orchestrator type
	Kubernetes OrchestratorType = "Kubernetes"
)

const (
	Windows OSType = "Windows"
	Linux   OSType = "Linux"
)

// subscription states
const (
	// Registered means the subscription is entitled to use the namespace
	Registered SubscriptionState = iota
	// Unregistered means the subscription is not entitled to use the namespace
	Unregistered
	// Suspended means the subscription has been suspended from the system
	Suspended
	// Deleted means the subscription has been deleted
	Deleted
	// Warned means the subscription has been warned
	Warned
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

// ClassicAgentPoolProfileType profiles
const (
	SwarmPublic ClassicAgentPoolProfileType = "SwarmPublic"
	DCOSPublic  ClassicAgentPoolProfileType = "DCOSPublic"
	DCOSPrivate ClassicAgentPoolProfileType = "DCOSPrivate"
	NotClassic  ClassicAgentPoolProfileType = ""
)

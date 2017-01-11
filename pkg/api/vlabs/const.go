package vlabs

const (
	// APIVersion is the version of this API
	APIVersion = "vlabs"
)

// the orchestrators supported by vlabs
const (
	// DCOS is the string constant for DCOS orchestrator type and defaults to DCOS187
	DCOS = "DCOS"
	// DCOS187 is the string constant for DCOS 1.8.7 orchestrator type
	DCOS187 OrchestratorType = "DCOS187"
	// DCOS184 is the string constant for DCOS 1.8.4 orchestrator type
	DCOS184 = "DCOS184"
	// DCOS173 is the string constant for DCOS 1.7.3 orchestrator type
	DCOS173 = "DCOS173"
	// Swarm is the string constant for the Swarm orchestrator type
	Swarm = "Swarm"
	// Kubernetes is the string constant for the Kubernetes orchestrator type
	Kubernetes = "Kubernetes"
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
	// MinDiskSizeGB specifies the minimum attached disk size
	MinDiskSizeGB = 1
	// MaxDiskSizeGB specifies the maximum attached disk size
	MaxDiskSizeGB = 1024
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
	// StorageAccountClassic means that we follow the older versions (09-30-2016, 03-30-2016)
	// storage account naming conventions
	StorageAccountClassic = "StorageAccountClassic"
	// StorageAccount means that the nodes use raw storage accounts for their os and attached volumes
	StorageAccount = "StorageAccount"
	// ManagedDisks means that the nodes use managed disks for their os and attached volumes
	ManagedDisks = "ManagedDisks"
)

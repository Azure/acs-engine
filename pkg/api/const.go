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
	// DockerCE is the string constant for the Swarm Mode orchestrator type
	DockerCE OrchestratorType = "DockerCE"
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
	// StorageAccountClassic means that we follow the older versions (09-30-2016, 03-30-2016)
	// storage account naming conventions
	StorageAccountClassic = "StorageAccountClassic"
	// StorageAccount means that the nodes use raw storage accounts for their os and attached volumes
	StorageAccount = "StorageAccount"
	// ManagedDisks means that the nodes use managed disks for their os and attached volumes
	ManagedDisks = "ManagedDisks"
)

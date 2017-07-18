package v20170831

const (
	// APIVersion is the version of this API
	APIVersion = "2017-08-31"
)

const (
	Windows OSType = "Windows"
	Linux   OSType = "Linux"
)

// validation values
const (
	// MinAgentCount are the minimum number of agents
	MinAgentCount = 1
	// MaxAgentCount are the maximum number of agents
	MaxAgentCount = 100
	// MinDiskSizeGB specifies the minimum attached disk size
	MinDiskSizeGB = 1
	// MaxDiskSizeGB specifies the maximum attached disk size
	MaxDiskSizeGB = 1023
)

const (
	// Kubernetes166 is the string constant for Kubernetes 1.6.6
	Kubernetes166 OrchestratorVersion = "1.6.6"
	// KubernetesLatest is the string constant for latest Kubernetes version
	KubernetesLatest OrchestratorVersion = Kubernetes166
)

// storage profiles
const (
	// StorageAccount means that the nodes use raw storage accounts for their os and attached volumes
	StorageAccount = "StorageAccount"
	// ManagedDisks means that the nodes use managed disks for their os and attached volumes
	ManagedDisks = "ManagedDisks"
)

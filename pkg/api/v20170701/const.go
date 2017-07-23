package v20170701

const (
	// APIVersion is the version of this API
	APIVersion = "2017-07-01"
)

// the orchestrators supported by 2017-07-01
const (
	// DCOS is the string constant for DCOS orchestrator type and defaults to DCOS187
	DCOS OrchestratorType = "DCOS"
	// Swarm is the string constant for the Swarm orchestrator type
	Swarm OrchestratorType = "Swarm"
	// Kubernetes is the string constant for the Kubernetes orchestrator type
	Kubernetes OrchestratorType = "Kubernetes"
	// DockerCE is the string constant for the Docker CE orchestrator type
	DockerCE OrchestratorType = "DockerCE"
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
	// MinPort specifies the minimum tcp port to open
	MinPort = 1
	// MaxPort specifies the maximum tcp port to open
	MaxPort = 65535
	// MinDiskSizeGB specifies the minimum attached disk size
	MinDiskSizeGB = 1
	// MaxDiskSizeGB specifies the maximum attached disk size
	MaxDiskSizeGB = 1023
)

// storage profiles
const (
	// StorageAccount means that the nodes use raw storage accounts for their os and attached volumes
	StorageAccount = "StorageAccount"
	// ManagedDisks means that the nodes use managed disks for their os and attached volumes
	ManagedDisks = "ManagedDisks"
)

const (
	// DCOS190 is the string constant for DCOS 1.9.0
	DCOS190 OrchestratorVersion = "1.9.0"
	// DCOS188 is the string constant for DCOS 1.8.8
	DCOS188 OrchestratorVersion = "1.8.8"
	// DCOS187 is the string constant for DCOS 1.8.7
	DCOS187 OrchestratorVersion = "1.8.7"
	// DCOS184 is the string constant for DCOS 1.8.4
	DCOS184 OrchestratorVersion = "1.8.4"
	// DCOS173 is the string constant for DCOS 1.7.3
	DCOS173 OrchestratorVersion = "1.7.3"
	// DCOSLatest is the string constant for latest DCOS version
	DCOSLatest OrchestratorVersion = DCOS190
)

const (
	// Kubernetes153 is the string constant for Kubernetes 1.5.3
	Kubernetes153 OrchestratorVersion = "1.5.3"
	// Kubernetes157 is the string constant for Kubernetes 1.5.3
	Kubernetes157 OrchestratorVersion = "1.5.7"
	// Kubernetes160 is the string constant for Kubernetes 1.6.0
	Kubernetes160 OrchestratorVersion = "1.6.0"
	// Kubernetes162 is the string constant for Kubernetes 1.6.2
	Kubernetes162 OrchestratorVersion = "1.6.2"
	// Kubernetes166 is the string constant for Kubernetes 1.6.6
	Kubernetes166 OrchestratorVersion = "1.6.6"
	// KubernetesDefaultVersion is the string constant for current Kubernetes version
	KubernetesDefaultVersion OrchestratorVersion = Kubernetes166
)

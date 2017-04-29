package v20170131

const (
	// APIVersion is the version of this API
	APIVersion = "2017-01-31"
)

// the orchestrators supported by 2017-01-31
const (
	// Mesos is the string constant for the Mesos orchestrator type
	Mesos OrchestratorType = "Mesos"
	// DCOS is the string constant for DCOS orchestrator type and defaults to DCOS187
	DCOS OrchestratorType = "DCOS"
	// Swarm is the string constant for the Swarm orchestrator type
	Swarm OrchestratorType = "Swarm"
	// Kubernetes is the string constant for the Kubernetes orchestrator type
	Kubernetes OrchestratorType = "Kubernetes"
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
)

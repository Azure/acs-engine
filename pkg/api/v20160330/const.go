package v20160330

const (
	// APIVersion is the version of this API
	APIVersion = "2016-03-30"
)

// v20160330 supports orchestrators Mesos, Swarm, DCOS
const (
	Mesos OrchestratorType = "Mesos"
	Swarm OrchestratorType = "Swarm"
	DCOS  OrchestratorType = "DCOS"
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
	// MinAgentCount are the minimum number of agents
	MinAgentCount = 1
	// MaxAgentCount are the maximum number of agents
	MaxAgentCount = 100
)

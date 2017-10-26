package v20160330

const (
	// APIVersion is the version of this API
	APIVersion = "2016-03-30"
)

// v20160330 supports orchestrators Mesos, Swarm, DCOS
const (
	Mesos string = "Mesos"
	Swarm string = "Swarm"
	DCOS  string = "DCOS"
)

// v20160330 supports OSTypes Windows and Linux
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

package vlabs

const (
	// DCOS is the string constant for DCOS orchestrator type and defaults to DCOS184
	DCOS = "DCOS"
	// DCOS184 is the string constant for DCOS 1.8.4 orchestrator type
	DCOS184 = "DCOS184"
	// DCOS173 is the string constant for DCOS 1.7.3 orchestrator type
	DCOS173 = "DCOS173"
	// SWARM is the string constant for the Swarm orchestrator type
	SWARM = "Swarm"
	// MinAgentCount are the minimum number of agents
	MinAgentCount = 1
	// MaxAgentCount are the maximum number of agents
	MaxAgentCount = 100
	// MinPort specifies the minimum tcp port to open
	MinPort = 1
	// MaxPort specifies the maximum tcp port to open
	MaxPort = 65535
	// DefaultMasterSubnet
	DefaultMasterSubnet = "172.16.0.0/24"
)

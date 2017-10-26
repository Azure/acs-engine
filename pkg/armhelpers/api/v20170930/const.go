package v20170930

const (
	// APIVersion is the version of this API
	APIVersion = "2017-09-30"
)

// the orchestrators supported by 2017-07-01
const (
	// DCOS is the string constant for DCOS orchestrator type and defaults to DCOS187
	DCOS string = "DCOS"
	// Swarm is the string constant for the Swarm orchestrator type
	Swarm string = "Swarm"
	// Kubernetes is the string constant for the Kubernetes orchestrator type
	Kubernetes string = "Kubernetes"
	// DockerCE is the string constant for the Docker CE orchestrator type
	DockerCE string = "DockerCE"
)

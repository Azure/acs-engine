package vlabs

// AcsCluster represents the ACS cluster definition
type AcsCluster struct {
	OrchestratorProfile OrchestratorProfile `json:"orchestratorProfile"`
	MasterProfile       MasterProfile       `json:"masterProfile"`
	AgentPoolProfiles   []AgentPoolProfile  `json:"agentPoolProfiles"`
	LinuxProfile        LinuxProfile        `json:"linuxProfile"`
}

// OrchestratorProfile represents the type of orchestrator
type OrchestratorProfile struct {
	OrchestratorType string `json:"orchestratorType"`
}

// MasterProfile represents the definition of the master cluster
type MasterProfile struct {
	Count     int    `json:"count"`
	DNSPrefix string `json:"dnsPrefix"`
	VMSize    string `json:"vmSize"`
	Subnet    string `json:"subnet,omitempty"`
}

// AgentPoolProfile represents an agent pool definition
type AgentPoolProfile struct {
	Name        string `json:"name"`
	Count       int    `json:"count"`
	VMSize      string `json:"vmSize"`
	Subnet      string `json:"subnet,omitempty"`
	IsStateless bool   `json:"isStateless,omitempty"`
	DNSPrefix   string `json:"dnsPrefix,omitempty"`
	Ports       []int  `json:"ports,omitempty"`
}

// LinuxProfile represents the linux parameters passed to the cluster
type LinuxProfile struct {
	AdminUsername string `json:"adminUsername"`
	SSH           struct {
		PublicKeys []struct {
			KeyData string `json:"keyData"`
		} `json:"publicKeys"`
	} `json:"ssh"`
}

// APIObject defines the required functionality of an api object
type APIObject interface {
	SetDefaults()
	Validate() error
}

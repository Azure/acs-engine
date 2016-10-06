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
	OrchestratorType             string `json:"orchestratorType"`
	ServicePrincipalClientID     string `json:"servicePrincipalClientID,omitempty"`
	ServicePrincipalClientSecret string `json:"servicePrincipalClientSecret,omitempty"`
	ApiServerCertificate         string `json:"apiServerCertificate,omitempty"`
	ApiServerPrivateKey          string `json:"apiServerPrivateKey,omitempty"`
	CaCertificate                string `json:"caCertificate,omitempty"`
	ClientCertificate            string `json:"clientCertificate,omitempty"`
	ClientPrivateKey             string `json:"clientPrivateKey,omitempty"`
	ClusterID                    string `json:"clusterid,omitempty"`
	// caPrivateKey is an internal field only set if generation required
	caPrivateKey string
}

// MasterProfile represents the definition of the master cluster
type MasterProfile struct {
	Count                    int    `json:"count"`
	DNSPrefix                string `json:"dnsPrefix"`
	VMSize                   string `json:"vmSize"`
	VnetSubnetID             string `json:"vnetSubnetID,omitempty"`
	FirstConsecutiveStaticIP string `json:"firstConsecutiveStaticIP,omitempty"`
	// subnet is internal
	subnet string
}

// AgentPoolProfile represents an agent pool definition
type AgentPoolProfile struct {
	Name         string `json:"name"`
	Count        int    `json:"count"`
	VMSize       string `json:"vmSize"`
	DNSPrefix    string `json:"dnsPrefix,omitempty"`
	Ports        []int  `json:"ports,omitempty"`
	IsStateful   bool   `json:"isStateful,omitempty"`
	DiskSizesGB  []int  `json:"diskSizesGB,omitempty"`
	VnetSubnetID string `json:"vnetSubnetID,omitempty"`
	// subnet is internal
	subnet string
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
	Validate() error
}

// GetCAPrivateKey returns the ca private key
func (o *OrchestratorProfile) GetCAPrivateKey() string {
	return o.caPrivateKey
}

// SetCAPrivateKey returns the ca private key
func (o *OrchestratorProfile) SetCAPrivateKey(caPrivateKey string) {
	o.caPrivateKey = caPrivateKey
}

// IsCustomVNET returns true if the customer brought their own VNET
func (m *MasterProfile) IsCustomVNET() bool {
	return len(m.VnetSubnetID) > 0
}

// GetSubnet returns the read-only subnet for the master
func (m *MasterProfile) GetSubnet() string {
	return m.subnet
}

// SetSubnet sets the read-only subnet for the master
func (m *MasterProfile) SetSubnet(subnet string) {
	m.subnet = subnet
}

// IsCustomVNET returns true if the customer brought their own VNET
func (a *AgentPoolProfile) IsCustomVNET() bool {
	return len(a.VnetSubnetID) > 0
}

// HasDisks returns true if the customer specified disks
func (a *AgentPoolProfile) HasDisks() bool {
	return len(a.DiskSizesGB) > 0
}

// GetSubnet returns the read-only subnet for the agent pool
func (a *AgentPoolProfile) GetSubnet() string {
	return a.subnet
}

// SetSubnet sets the read-only subnet for the agent pool
func (a *AgentPoolProfile) SetSubnet(subnet string) {
	a.subnet = subnet
}

package vlabs

// AcsCluster represents the ACS cluster definition
type AcsCluster struct {
	OrchestratorProfile     OrchestratorProfile     `json:"orchestratorProfile"`
	MasterProfile           MasterProfile           `json:"masterProfile"`
	AgentPoolProfiles       []AgentPoolProfile      `json:"agentPoolProfiles"`
	WindowsProfile          WindowsProfile          `json:"windowsProfile"`
	LinuxProfile            LinuxProfile            `json:"linuxProfile"`
	ServicePrincipalProfile ServicePrincipalProfile `json:"servicePrincipalProfile"`
	CertificateProfile      CertificateProfile      `json:"certificateProfile"`
	// classic mode is used to output parameters and outputs
	classicMode bool
}

// OrchestratorProfile represents the type of orchestrator
type OrchestratorProfile struct {
	OrchestratorType string `json:"orchestratorType"`
}

// ServicePrincipalProfile contains the client and secret used by the cluster for Azure Resource CRUD
type ServicePrincipalProfile struct {
	ClientID string `json:"servicePrincipalClientID,omitempty"`
	Secret   string `json:"servicePrincipalClientSecret,omitempty"`
}

// CertificateProfile represents the definition of the master cluster
type CertificateProfile struct {
	// CaCertificate is the certificate authority certificate.
	CaCertificate string `json:"caCertificate,omitempty"`
	// ApiServerCertificate is the rest api server certificate, and signed by the CA
	APIServerCertificate string `json:"apiServerCertificate,omitempty"`
	// ApiServerPrivateKey is the rest api server private key, and signed by the CA
	APIServerPrivateKey string `json:"apiServerPrivateKey,omitempty"`
	// ClientCertificate is the certificate used by the client kubelet services and signed by the CA
	ClientCertificate string `json:"clientCertificate,omitempty"`
	// ClientPrivateKey is the private key used by the client kubelet services and signed by the CA
	ClientPrivateKey string `json:"clientPrivateKey,omitempty"`
	// KubeConfigCertificate is the client certificate used for kubectl cli and signed by the CA
	KubeConfigCertificate string `json:"kubeConfigCertificate,omitempty"`
	// KubeConfigPrivateKey is the client private key used for kubectl cli and signed by the CA
	KubeConfigPrivateKey string `json:"kubeConfigPrivateKey,omitempty"`
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
	OSType       string `json:"osType,omitempty"`
	Ports        []int  `json:"ports,omitempty"`
	StorageType  string `json:"storageType,omitempty"`
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

// WindowsProfile represents the windows parameters passed to the cluster
type WindowsProfile struct {
	AdminUsername string `json:"adminUsername"`
	AdminPassword string `json:"adminPassword"`
}

// APIObject defines the required functionality of an api object
type APIObject interface {
	Validate() error
}

// GetClassicMode gets the classic mode for deciding to output classic parameters
func (a *AcsCluster) GetClassicMode() bool {
	return a.classicMode
}

// SetClassicMode toggles classic parameters and outputs
func (a *AcsCluster) SetClassicMode(isClassicMode bool) {
	a.classicMode = isClassicMode
}

// HasWindows returns true if the cluster contains windows
func (a *AcsCluster) HasWindows() bool {
	for _, agentPoolProfile := range a.AgentPoolProfiles {
		if agentPoolProfile.OSType == OSTypeWindows {
			return true
		}
	}
	return false
}

// GetCAPrivateKey returns the ca private key
func (c *CertificateProfile) GetCAPrivateKey() string {
	return c.caPrivateKey
}

// SetCAPrivateKey sets the ca private key
func (c *CertificateProfile) SetCAPrivateKey(caPrivateKey string) {
	c.caPrivateKey = caPrivateKey
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

// IsWindows returns true if the agent pool is windows
func (a *AgentPoolProfile) IsWindows() bool {
	return a.OSType == OSTypeWindows
}

// IsVolumeBasedStorage returns true if the customer specified disks
func (a *AgentPoolProfile) IsVolumeBasedStorage() bool {
	return a.StorageType == StorageVolumes
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

package v20170727

const (

	// APIVersion is the unique string to identify this API
	APIVersion = "v20170727"
)

// AgentPool represents a Kubernetes Agent Pool
type AgentPool struct {
	ID       string                `json:"id,omitempty"`
	Location string                `json:"location,omitempty"`
	Name     string                `json:"name,omitempty"`
	Plan     *ResourcePurchasePlan `json:"plan,omitempty"`
	Tags     map[string]string     `json:"tags,omitempty"`
	Type     string                `json:"type,omitempty"`

	Properties *Properties `json:"properties"`
}

// Properties represents all data needed to define agent pools for Kubernetes
type Properties struct {
	KubernetesVersion       string                   `json:"kubernetesVersion"`
	KubernetesEndpoint      string                   `json:"kubernetesEndpoint"`
	DNSPrefix               string                   `json:"dnsPrefix,omitempty"`
	Version                 string                   `json:"version,omitempty"`
	AgentPoolProfiles       []*AgentPoolProfile      `json:"agentPoolProfiles,omitempty"`
	WindowsProfile          *WindowsProfile          `json:"windowsProfile,omitempty"`
	LinuxProfile            *LinuxProfile            `json:"linuxProfile,omitempty"`
	NetworkProfile          *NetworkProfile          `json:"networkProfile,omitempty"`
	JumpBoxProfile          *JumpBoxProfile          `json:"jumpboxProfile,omitempty"`
	ServicePrincipalProfile *ServicePrincipalProfile `json:"servicePrincipalProfile,omitempty"`
	CertificateProfile      *CertificateProfile      `json:"certificateProfile,omitempty"`
}

// CertificateProfile represents the TLS material for connecting to the Kubernetes API server
type CertificateProfile struct {
	// CaCertificate is the certificate authority certificate.
	CaCertificate string `json:"caCertificate,omitempty"`
	// CaPrivateKey is the certificate authority key.
	CaPrivateKey string `json:"caPrivateKey,omitempty"`
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
}

// ServicePrincipalProfile represents a service principal in Azure
type ServicePrincipalProfile struct {
	ClientID string `json:"servicePrincipalClientID,omitempty"`
	Secret   string `json:"servicePrincipalClientSecret,omitempty"`
}

// JumpBoxProfile represents the jumpbox that will be created with agent pools
type JumpBoxProfile struct {
	PublicIPAddressID string `json:"publicIpAddressId,omitempty"`
	// internalAddress must be inside the VNET and k8s-subnet
	InternalAddress string `json:"internalAddress,omitempty"`
	VMSize          string `json:"vmSize,omitempty"`
	Count           int    `json:"count,omitempty"`
}

// NetworkProfile represents the network that will be configured with agent pools
type NetworkProfile struct {
	AgentCidr        string `json:"agentCIDR,omitempty"`
	VnetSubnetID     string `json:"vnetSubnetID,omitempty"`
	KubeDNSServiceIP string `json:"kubeDnsServiceIP,omitempty"`
}

// AgentPoolProfile represents a single agent pool
type AgentPoolProfile struct {
	Name         string `json:"name,omitempty"`
	Count        int    `json:"count,omitempty"`
	VMSize       string `json:"vmSize,omitempty"`
	OSType       string `json:"osType,omitempty"`
	OSDiskSizeGb int    `json:"osDiskSizeGb,omitempty"`
}

// WindowsProfile represents the windows parameters passed to the cluster
type WindowsProfile struct {
	AdminUsername string            `json:"adminUsername,omitempty"`
	AdminPassword string            `json:"adminPassword,omitempty"`
	Secrets       []KeyVaultSecrets `json:"secrets,omitempty"`
}

// LinuxProfile represents the linux parameters passed to the cluster
type LinuxProfile struct {
	AdminUsername string `json:"adminUsername"`
	SSH           struct {
		PublicKeys []struct {
			KeyData string `json:"keyData"`
		} `json:"publicKeys"`
	} `json:"ssh"`
	Secrets []KeyVaultSecrets `json:"secrets,omitempty"`
}

// KeyVaultSecrets represents key vault secrets in Azure
type KeyVaultSecrets struct {
	SourceVault       *KeyVaultID           `json:"sourceVault,omitempty"`
	VaultCertificates []KeyVaultCertificate `json:"vaultCertificates,omitempty"`
}

// KeyVaultID represents the ID for a key vault in Azure
type KeyVaultID struct {
	ID string `json:"id,omitempty"`
}

// KeyVaultCertificate represents the TLS certificate for key vault in Azure
type KeyVaultCertificate struct {
	CertificateURL   string `json:"certificateUrl,omitempty"`
	CertificateStore string `json:"certificateStore,omitempty"`
}

// ResourcePurchasePlan represents the resource purchase plan in Azure
type ResourcePurchasePlan struct {
	Name          string `json:"name,omitempty"`
	Product       string `json:"product,omitempty"`
	PromotionCode string `json:"promotionCode,omitempty"`
	Publisher     string `json:"publisher,omitempty"`
}

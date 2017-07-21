package kubernetesagentpool

const (
	APIVersion = "kubernetesagentpool"
)

type AgentPool struct {
	ID       string                `json:"id,omitempty"`
	Location string                `json:"location,omitempty"`
	Name     string                `json:"name,omitempty"`
	Plan     *ResourcePurchasePlan `json:"plan,omitempty"`
	Tags     map[string]string     `json:"tags,omitempty"`
	Type     string                `json:"type,omitempty"`

	Properties *Properties `json:"properties"`
}

type Properties struct {
	KubernetesVersion       string                   `json:"kubernetesEndpoint"`
	KubernetesEndpoint      string                   `json:"kubernetesEndpoint"`
	DnsPrefix               string                   `json:"dnsPrefix,omitempty"`
	Version                 string                   `json:"version,omitempty"`
	AgentPoolProfiles       []*AgentPoolProfile      `json:"agentPoolProfiles,omitempty"`
	WindowsProfile          *WindowsProfile          `json:"windowsProfile,omitempty"`
	LinuxProfile            *LinuxProfile            `json:"linuxProfile,omitempty"`
	NetworkProfile          *NetworkProfile          `json:"networkProfile,omitempty"`
	JumpBoxProfile          *JumpBoxProfile          `json:"jumpboxProfile,omitempty"`
	ServicePrincipalProfile *ServicePrincipalProfile `json:"servicePrincipalProfile,omitempty"`
	CertificateProfile      *CertificateProfile      `json:"certificateProfile,omitempty"`
}

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

type ServicePrincipalProfile struct {
	ClientID string `json:"servicePrincipalClientID,omitempty"`
	Secret   string `json:"servicePrincipalClientSecret,omitempty"`
}

type JumpBoxProfile struct {
	PublicIpAddressId string `json:"publicIpAddressId,omitempty"`
	VMSize            string `json:"vmSize,omitempty"`
	Count             int    `json:"count,omitempty"`
}

type NetworkProfile struct {
	PodCIDR          string `json:"podCidr,omitempty"`
	ServiceCIDR      string `json:"serviceCidr,omitempty"`
	VnetSubnetId     string `json:"vnetSubnetID,omitempty"`
	KubeDnsServiceIp string `json:"kubeDnsServiceIP,omitempty"`
}

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

type KeyVaultSecrets struct {
	SourceVault       *KeyVaultID           `json:"sourceVault,omitempty"`
	VaultCertificates []KeyVaultCertificate `json:"vaultCertificates,omitempty"`
}

type KeyVaultID struct {
	ID string `json:"id,omitempty"`
}

type KeyVaultCertificate struct {
	CertificateURL   string `json:"certificateUrl,omitempty"`
	CertificateStore string `json:"certificateStore,omitempty"`
}

type ResourcePurchasePlan struct {
	Name          string `json:"name,omitempty"`
	Product       string `json:"product,omitempty"`
	PromotionCode string `json:"promotionCode,omitempty"`
	Publisher     string `json:"publisher,omitempty"`
}

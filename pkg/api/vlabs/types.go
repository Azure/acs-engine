package vlabs

import (
	"encoding/json"
	"fmt"
	"strings"
)

// ResourcePurchasePlan defines resource plan as required by ARM
// for billing purposes.
type ResourcePurchasePlan struct {
	Name          string `json:"name,omitempty"`
	Product       string `json:"product,omitempty"`
	PromotionCode string `json:"promotionCode,omitempty"`
	Publisher     string `json:"publisher,omitempty"`
}

// ContainerService complies with the ARM model of
// resource definition in a JSON template.
type ContainerService struct {
	ID       string                `json:"id,omitempty"`
	Location string                `json:"location,omitempty"`
	Name     string                `json:"name,omitempty"`
	Plan     *ResourcePurchasePlan `json:"plan,omitempty"`
	Tags     map[string]string     `json:"tags,omitempty"`
	Type     string                `json:"type,omitempty"`

	Properties *Properties `json:"properties"`
}

// Properties represents the ACS cluster definition
type Properties struct {
	ProvisioningState       ProvisioningState        `json:"provisioningState,omitempty"`
	OrchestratorProfile     *OrchestratorProfile     `json:"orchestratorProfile,omitempty" validate:"required"`
	MasterProfile           *MasterProfile           `json:"masterProfile,omitempty" validate:"required"`
	AgentPoolProfiles       []*AgentPoolProfile      `json:"agentPoolProfiles,omitempty" validate:"dive,required"`
	LinuxProfile            *LinuxProfile            `json:"linuxProfile,omitempty" validate:"required"`
	ExtensionProfiles       []*ExtensionProfile      `json:"extensionProfiles,omitempty"`
	WindowsProfile          *WindowsProfile          `json:"windowsProfile,omitempty"`
	ServicePrincipalProfile *ServicePrincipalProfile `json:"servicePrincipalProfile,omitempty"`
	CertificateProfile      *CertificateProfile      `json:"certificateProfile,omitempty"`
	AADProfile              *AADProfile              `json:"aadProfile,omitempty"`
}

// ServicePrincipalProfile contains the client and secret used by the cluster for Azure Resource CRUD
// The 'Secret' and 'KeyvaultSecretRef' parameters are mutually exclusive
// The 'Secret' parameter should be a secret in plain text.
// The 'KeyvaultSecretRef' parameter is a reference to a secret in a keyvault.
type ServicePrincipalProfile struct {
	ClientID          string             `json:"clientId,omitempty"`
	Secret            string             `json:"secret,omitempty"`
	KeyvaultSecretRef *KeyvaultSecretRef `json:"keyvaultSecretRef,omitempty"`
}

// KeyvaultSecretRef is a reference to a secret in a keyvault.
// The format of 'VaultID' value should be
// "/subscriptions/<SUB_ID>/resourceGroups/<RG_NAME>/providers/Microsoft.KeyVault/vaults/<KV_NAME>"
// where:
//    <SUB_ID> is the subscription ID of the keyvault
//    <RG_NAME> is the resource group of the keyvault
//    <KV_NAME> is the name of the keyvault
// The 'SecretName' is the name of the secret in the keyvault
// The 'SecretVersion' (optional) is the version of the secret (default: the latest version)
type KeyvaultSecretRef struct {
	VaultID       string `json:"vaultID" validate:"required"`
	SecretName    string `json:"secretName" validate:"required"`
	SecretVersion string `json:"version,omitempty"`
}

// CertificateProfile represents the definition of the master cluster
// The JSON parameters could be either a plain text, or referenced to a secret in a keyvault.
// In the latter case, the format of the parameter's value should be
// "/subscriptions/<SUB_ID>/resourceGroups/<RG_NAME>/providers/Microsoft.KeyVault/vaults/<KV_NAME>/secrets/<NAME>[/<VERSION>]"
// where:
//    <SUB_ID> is the subscription ID of the keyvault
//    <RG_NAME> is the resource group of the keyvault
//    <KV_NAME> is the name of the keyvault
//    <NAME> is the name of the secret
//    <VERSION> (optional) is the version of the secret (default: the latest version)
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

// LinuxProfile represents the linux parameters passed to the cluster
type LinuxProfile struct {
	AdminUsername string `json:"adminUsername" validate:"required"`
	SSH           struct {
		PublicKeys []PublicKey `json:"publicKeys" validate:"required,len=1"`
	} `json:"ssh" validate:"required"`
	Secrets []KeyVaultSecrets `json:"secrets,omitempty"`
}

// PublicKey represents an SSH key for LinuxProfile
type PublicKey struct {
	KeyData string `json:"keyData"`
}

// WindowsProfile represents the windows parameters passed to the cluster
type WindowsProfile struct {
	AdminUsername string            `json:"adminUsername,omitempty"`
	AdminPassword string            `json:"adminPassword,omitempty"`
	Secrets       []KeyVaultSecrets `json:"secrets,omitempty"`
}

// ProvisioningState represents the current state of container service resource.
type ProvisioningState string

const (
	// Creating means ContainerService resource is being created.
	Creating ProvisioningState = "Creating"
	// Updating means an existing ContainerService resource is being updated
	Updating ProvisioningState = "Updating"
	// Failed means resource is in failed state
	Failed ProvisioningState = "Failed"
	// Succeeded means resource created succeeded during last create/update
	Succeeded ProvisioningState = "Succeeded"
	// Deleting means resource is in the process of being deleted
	Deleting ProvisioningState = "Deleting"
	// Migrating means resource is being migrated from one subscription or
	// resource group to another
	Migrating ProvisioningState = "Migrating"
)

// OrchestratorProfile contains Orchestrator properties
type OrchestratorProfile struct {
	OrchestratorType    string            `json:"orchestratorType" validate:"required"`
	OrchestratorRelease string            `json:"orchestratorRelease"`
	OrchestratorVersion string            `json:"orchestratorVersion" validate:"len=0"`
	KubernetesConfig    *KubernetesConfig `json:"kubernetesConfig,omitempty"`
}

// UnmarshalJSON unmarshal json using the default behavior
// And do fields manipulation, such as populating default value
func (o *OrchestratorProfile) UnmarshalJSON(b []byte) error {
	// Need to have a alias type to avoid circular unmarshal
	type aliasOrchestratorProfile OrchestratorProfile
	op := aliasOrchestratorProfile{}
	if e := json.Unmarshal(b, &op); e != nil {
		return e
	}
	*o = OrchestratorProfile(op)
	// Unmarshal OrchestratorType, format it as well
	orchestratorType := o.OrchestratorType
	switch {
	case strings.EqualFold(orchestratorType, DCOS):
		o.OrchestratorType = DCOS
	case strings.EqualFold(orchestratorType, Swarm):
		o.OrchestratorType = Swarm
	case strings.EqualFold(orchestratorType, Kubernetes):
		o.OrchestratorType = Kubernetes
	case strings.EqualFold(orchestratorType, SwarmMode):
		o.OrchestratorType = SwarmMode
	default:
		return fmt.Errorf("OrchestratorType has unknown orchestrator: %s", orchestratorType)
	}
	return nil
}

// KubernetesConfig contains the Kubernetes config structure, containing
// Kubernetes specific configuration
type KubernetesConfig struct {
	KubernetesImageBase              string  `json:"kubernetesImageBase,omitempty"`
	ClusterSubnet                    string  `json:"clusterSubnet,omitempty"`
	DNSServiceIP                     string  `json:"dnsServiceIP,omitempty"`
	ServiceCidr                      string  `json:"serviceCidr,omitempty"`
	NetworkPolicy                    string  `json:"networkPolicy,omitempty"`
	MaxPods                          int     `json:"maxPods,omitempty"`
	DockerBridgeSubnet               string  `json:"DockerBridgeSubnet,omitempty"`
	NodeStatusUpdateFrequency        string  `json:"nodeStatusUpdateFrequency,omitempty"`
	CtrlMgrNodeMonitorGracePeriod    string  `json:"ctrlMgrNodeMonitorGracePeriod,omitempty"`
	CtrlMgrPodEvictionTimeout        string  `json:"ctrlMgrPodEvictionTimeout,omitempty"`
	CtrlMgrRouteReconciliationPeriod string  `json:"ctrlMgrRouteReconciliationPeriod,omitempty"`
	CloudProviderBackoff             bool    `json:"cloudProviderBackoff,omitempty"`
	CloudProviderBackoffRetries      int     `json:"cloudProviderBackoffRetries,omitempty"`
	CloudProviderBackoffJitter       float64 `json:"cloudProviderBackoffJitter,omitempty"`
	CloudProviderBackoffDuration     int     `json:"cloudProviderBackoffDuration,omitempty"`
	CloudProviderBackoffExponent     float64 `json:"cloudProviderBackoffExponent,omitempty"`
	CloudProviderRateLimit           bool    `json:"cloudProviderRateLimit,omitempty"`
	CloudProviderRateLimitQPS        float64 `json:"cloudProviderRateLimitQPS,omitempty"`
	CloudProviderRateLimitBucket     int     `json:"cloudProviderRateLimitBucket,omitempty"`
	UseManagedIdentity               bool    `json:"useManagedIdentity,omitempty"`
	CustomHyperkubeImage             string  `json:"customHyperkubeImage,omitempty"`
	UseInstanceMetadata              bool    `json:"useInstanceMetadata,omitempty"`
	EnableRbac                       bool    `json:"enableRbac,omitempty"`
}

// MasterProfile represents the definition of the master cluster
type MasterProfile struct {
	Count                    int         `json:"count" validate:"required,eq=1|eq=3|eq=5"`
	DNSPrefix                string      `json:"dnsPrefix" validate:"required"`
	VMSize                   string      `json:"vmSize" validate:"required"`
	OSDiskSizeGB             int         `json:"osDiskSizeGB,omitempty" validate:"min=0,max=1023"`
	VnetSubnetID             string      `json:"vnetSubnetID,omitempty"`
	VnetCidr                 string      `json:"vnetCidr,omitempty"`
	FirstConsecutiveStaticIP string      `json:"firstConsecutiveStaticIP,omitempty"`
	IPAddressCount           int         `json:"ipAddressCount,omitempty" validate:"min=0,max=256"`
	StorageProfile           string      `json:"storageProfile,omitempty" validate:"eq=StorageAccount|eq=ManagedDisks|len=0"`
	HTTPSourceAddressPrefix  string      `json:"HTTPSourceAddressPrefix,omitempty"`
	OAuthEnabled             bool        `json:"oauthEnabled"`
	PreProvisionExtension    *Extension  `json:"preProvisionExtension"`
	Extensions               []Extension `json:"extensions"`

	// subnet is internal
	subnet string

	// Master LB public endpoint/FQDN with port
	// The format will be FQDN:2376
	// Not used during PUT, returned as part of GET
	FQDN string `json:"fqdn,omitempty"`
}

// ClassicAgentPoolProfileType represents types of classic profiles
type ClassicAgentPoolProfileType string

// ExtensionProfile represents an extension definition
type ExtensionProfile struct {
	Name                string `json:"name"`
	Version             string `json:"version"`
	ExtensionParameters string `json:"extensionParameters"`
	RootURL             string `json:"rootURL"`
	// This is only needed for preprovision extensions and it needs to be a bash script
	Script   string `json:"script"`
	URLQuery string `json:"urlQuery"`
}

// Extension represents an extension definition in the master or agentPoolProfile
type Extension struct {
	Name        string `json:"name"`
	SingleOrAll string `json:"singleOrAll"`
	Template    string `json:"template"`
}

// AgentPoolProfile represents an agent pool definition
type AgentPoolProfile struct {
	Name                string `json:"name" validate:"required"`
	Count               int    `json:"count" validate:"required,min=1,max=100"`
	VMSize              string `json:"vmSize" validate:"required"`
	OSDiskSizeGB        int    `json:"osDiskSizeGB,omitempty" validate:"min=0,max=1023"`
	DNSPrefix           string `json:"dnsPrefix,omitempty"`
	OSType              OSType `json:"osType,omitempty"`
	Ports               []int  `json:"ports,omitempty" validate:"dive,min=1,max=65535"`
	AvailabilityProfile string `json:"availabilityProfile"`
	StorageProfile      string `json:"storageProfile" validate:"eq=StorageAccount|eq=ManagedDisks|len=0"`
	DiskSizesGB         []int  `json:"diskSizesGB,omitempty" validate:"max=4,dive,min=1,max=1023"`
	VnetSubnetID        string `json:"vnetSubnetID,omitempty"`
	IPAddressCount      int    `json:"ipAddressCount,omitempty" validate:"min=0,max=256"`

	// subnet is internal
	subnet string

	FQDN                  string            `json:"fqdn"`
	CustomNodeLabels      map[string]string `json:"customNodeLabels,omitempty"`
	PreProvisionExtension *Extension        `json:"preProvisionExtension"`
	Extensions            []Extension       `json:"extensions"`
}

// AADProfile specifies attributes for AAD integration
type AADProfile struct {
	// The client AAD application ID.
	ClientAppID string `json:"clientAppID,omitempty"`
	// The server AAD application ID.
	ServerAppID string `json:"serverAppID,omitempty"`
	// The AAD tenant ID to use for authentication.
	// If not specified, will use the tenant of the deployment subscription.
	// Optional
	TenantID string `json:"tenantID,omitempty"`
}

// KeyVaultSecrets specifies certificates to install on the pool
// of machines from a given key vault
// the key vault specified must have been granted read permissions to CRP
type KeyVaultSecrets struct {
	SourceVault       *KeyVaultID           `json:"sourceVault,omitempty"`
	VaultCertificates []KeyVaultCertificate `json:"vaultCertificates,omitempty"`
}

// KeyVaultID specifies a key vault
type KeyVaultID struct {
	ID string `json:"id,omitempty"`
}

// KeyVaultCertificate specifies a certificate to install
// On Linux, the certificate file is placed under the /var/lib/waagent directory
// with the file name <UppercaseThumbprint>.crt for the X509 certificate file
// and <UppercaseThumbprint>.prv for the private key. Both of these files are .pem formatted.
// On windows the certificate will be saved in the specified store.
type KeyVaultCertificate struct {
	CertificateURL   string `json:"certificateUrl,omitempty"`
	CertificateStore string `json:"certificateStore,omitempty"`
}

// OSType represents OS types of agents
type OSType string

// HasWindows returns true if the cluster contains windows
func (p *Properties) HasWindows() bool {
	for _, agentPoolProfile := range p.AgentPoolProfiles {
		if agentPoolProfile.OSType == Windows {
			return true
		}
	}
	return false
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

// IsManagedDisks returns true if the master specified managed disks
func (m *MasterProfile) IsManagedDisks() bool {
	return m.StorageProfile == ManagedDisks
}

// IsStorageAccount returns true if the master specified storage account
func (m *MasterProfile) IsStorageAccount() bool {
	return m.StorageProfile == StorageAccount
}

// IsCustomVNET returns true if the customer brought their own VNET
func (a *AgentPoolProfile) IsCustomVNET() bool {
	return len(a.VnetSubnetID) > 0
}

// IsWindows returns true if the agent pool is windows
func (a *AgentPoolProfile) IsWindows() bool {
	return a.OSType == Windows
}

// IsLinux returns true if the agent pool is linux
func (a *AgentPoolProfile) IsLinux() bool {
	return a.OSType == Linux
}

// IsAvailabilitySets returns true if the customer specified disks
func (a *AgentPoolProfile) IsAvailabilitySets() bool {
	return a.AvailabilityProfile == AvailabilitySet
}

// IsManagedDisks returns true if the customer specified managed disks
func (a *AgentPoolProfile) IsManagedDisks() bool {
	return a.StorageProfile == ManagedDisks
}

// IsStorageAccount returns true if the customer specified storage account
func (a *AgentPoolProfile) IsStorageAccount() bool {
	return a.StorageProfile == StorageAccount
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

// IsSwarmMode returns true if this template is for Swarm Mode orchestrator
func (o *OrchestratorProfile) IsSwarmMode() bool {
	return o.OrchestratorType == SwarmMode
}

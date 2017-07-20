package vlabs

import (
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
	OrchestratorProfile     *OrchestratorProfile     `json:"orchestratorProfile,omitempty"`
	MasterProfile           *MasterProfile           `json:"masterProfile,omitempty"`
	AgentPoolProfiles       []*AgentPoolProfile      `json:"agentPoolProfiles,omitempty"`
	LinuxProfile            *LinuxProfile            `json:"linuxProfile,omitempty"`
	WindowsProfile          *WindowsProfile          `json:"windowsProfile,omitempty"`
	ServicePrincipalProfile *ServicePrincipalProfile `json:"servicePrincipalProfile,omitempty"`
	CertificateProfile      *CertificateProfile      `json:"certificateProfile,omitempty"`
}

// ServicePrincipalProfile contains the client and secret used by the cluster for Azure Resource CRUD
// The 'Secret' parameter should be a secret in plain text.
// The 'KeyvaultSecretRef' parameter is a reference to a secret in a keyvault.
// The format of the parameter's value should be
// "/subscriptions/<SUB_ID>/resourceGroups/<RG_NAME>/providers/Microsoft.KeyVault/vaults/<KV_NAME>/secrets/<NAME>[/<VERSION>]"
// where:
//    <SUB_ID> is the subscription ID of the keyvault
//    <RG_NAME> is the resource group of the keyvault
//    <KV_NAME> is the name of the keyvault
//    <NAME> is the name of the secret.
//    <VERSION> (optional) is the version of the secret (default: the latest version)
type ServicePrincipalProfile struct {
	ClientID          string `json:"servicePrincipalClientID,omitempty"`
	Secret            string `json:"servicePrincipalClientSecret,omitempty"`
	KeyvaultSecretRef string `json:"servicePrincipalClientKeyvaultSecretRef,omitempty"`
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
	AdminUsername string `json:"adminUsername"`
	SSH           struct {
		PublicKeys []PublicKey `json:"publicKeys"`
	} `json:"ssh"`
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
	OrchestratorType    OrchestratorType    `json:"orchestratorType"`
	OrchestratorVersion OrchestratorVersion `json:"orchestratorVersion"`
	KubernetesConfig    *KubernetesConfig   `json:"kubernetesConfig,omitempty"`
}

// KubernetesConfig contains the Kubernetes config structure, containing
// Kubernetes specific configuration
type KubernetesConfig struct {
	KubernetesImageBase              string  `json:"kubernetesImageBase,omitempty"`
	ClusterSubnet                    string  `json:"clusterSubnet,omitempty"`
	NetworkPolicy                    string  `json:"networkPolicy,omitempty"`
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
}

// MasterProfile represents the definition of the master cluster
type MasterProfile struct {
	Count                    int    `json:"count"`
	DNSPrefix                string `json:"dnsPrefix"`
	VMSize                   string `json:"vmSize"`
	OSDiskSizeGB             int    `json:"osDiskSizeGB,omitempty"`
	VnetSubnetID             string `json:"vnetSubnetID,omitempty"`
	FirstConsecutiveStaticIP string `json:"firstConsecutiveStaticIP,omitempty"`
	IPAddressCount           int    `json:"ipAddressCount,omitempty"`
	StorageProfile           string `json:"storageProfile,omitempty"`
	HttpSourceAddressPrefix  string `json:"httpSourceAddressPrefix,omitempty"`
	OAuthEnabled             bool   `json:"oauthEnabled"`

	// subnet is internal
	subnet string

	// Master LB public endpoint/FQDN with port
	// The format will be FQDN:2376
	// Not used during PUT, returned as part of GET
	FQDN string `json:"fqdn,omitempty"`
}

// ClassicAgentPoolProfileType represents types of classic profiles
type ClassicAgentPoolProfileType string

// AgentPoolProfile represents an agent pool definition
type AgentPoolProfile struct {
	Name                string `json:"name"`
	Count               int    `json:"count"`
	VMSize              string `json:"vmSize"`
	OSDiskSizeGB        int    `json:"osDiskSizeGB,omitempty"`
	DNSPrefix           string `json:"dnsPrefix,omitempty"`
	OSType              OSType `json:"osType,omitempty"`
	Ports               []int  `json:"ports,omitempty"`
	AvailabilityProfile string `json:"availabilityProfile"`
	StorageProfile      string `json:"storageProfile"`
	DiskSizesGB         []int  `json:"diskSizesGB,omitempty"`
	VnetSubnetID        string `json:"vnetSubnetID,omitempty"`
	IPAddressCount      int    `json:"ipAddressCount,omitempty"`

	// subnet is internal
	subnet string

	FQDN             string            `json:"fqdn"`
	CustomNodeLabels map[string]string `json:"customNodeLabels,omitempty"`
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

// OrchestratorType defines orchestrators supported by ACS
type OrchestratorType string

// OrchestratorVersion defines the version for orchestratorType
type OrchestratorVersion string

// UnmarshalText decodes OrchestratorType text, do a case insensitive comparison with
// the defined OrchestratorType constant and set to it if they equal
func (o *OrchestratorType) UnmarshalText(text []byte) error {
	s := string(text)
	switch {
	case strings.EqualFold(s, string(DCOS)):
		*o = DCOS
	case strings.EqualFold(s, string(Swarm)):
		*o = Swarm
	case strings.EqualFold(s, string(Kubernetes)):
		*o = Kubernetes
	case strings.EqualFold(s, string(SwarmMode)):
		*o = SwarmMode
	default:
		return fmt.Errorf("OrchestratorType has unknown orchestrator: %s", s)
	}

	return nil
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

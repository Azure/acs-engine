package api

import (
	neturl "net/url"

	"github.com/Azure/acs-engine/pkg/api/agentPoolOnlyApi/v20170831"
	"github.com/Azure/acs-engine/pkg/api/agentPoolOnlyApi/v20180331"
	"github.com/Azure/acs-engine/pkg/api/common"
	"github.com/Azure/acs-engine/pkg/api/v20160330"
	"github.com/Azure/acs-engine/pkg/api/v20160930"
	"github.com/Azure/acs-engine/pkg/api/v20170131"
	"github.com/Azure/acs-engine/pkg/api/v20170701"
	"github.com/Azure/acs-engine/pkg/api/vlabs"
	"github.com/Masterminds/semver"
)

// TypeMeta describes an individual API model object
type TypeMeta struct {
	// APIVersion is on every object
	APIVersion string `json:"apiVersion"`
}

// ResourcePurchasePlan defines resource plan as required by ARM
// for billing purposes.
type ResourcePurchasePlan struct {
	Name          string `json:"name"`
	Product       string `json:"product"`
	PromotionCode string `json:"promotionCode"`
	Publisher     string `json:"publisher"`
}

// ContainerService complies with the ARM model of
// resource definition in a JSON template.
type ContainerService struct {
	ID       string                `json:"id"`
	Location string                `json:"location"`
	Name     string                `json:"name"`
	Plan     *ResourcePurchasePlan `json:"plan,omitempty"`
	Tags     map[string]string     `json:"tags"`
	Type     string                `json:"type"`

	Properties *Properties `json:"properties,omitempty"`
}

// Properties represents the ACS cluster definition
type Properties struct {
	ProvisioningState       ProvisioningState        `json:"provisioningState,omitempty"`
	OrchestratorProfile     *OrchestratorProfile     `json:"orchestratorProfile,omitempty"`
	MasterProfile           *MasterProfile           `json:"masterProfile,omitempty"`
	AgentPoolProfiles       []*AgentPoolProfile      `json:"agentPoolProfiles,omitempty"`
	LinuxProfile            *LinuxProfile            `json:"linuxProfile,omitempty"`
	WindowsProfile          *WindowsProfile          `json:"windowsProfile,omitempty"`
	ExtensionProfiles       []*ExtensionProfile      `json:"extensionProfiles"`
	DiagnosticsProfile      *DiagnosticsProfile      `json:"diagnosticsProfile,omitempty"`
	JumpboxProfile          *JumpboxProfile          `json:"jumpboxProfile,omitempty"`
	ServicePrincipalProfile *ServicePrincipalProfile `json:"servicePrincipalProfile,omitempty"`
	CertificateProfile      *CertificateProfile      `json:"certificateProfile,omitempty"`
	AADProfile              *AADProfile              `json:"aadProfile,omitempty"`
	CustomProfile           *CustomProfile           `json:"customProfile,omitempty"`
	HostedMasterProfile     *HostedMasterProfile     `json:"hostedMasterProfile,omitempty"`
	AddonProfiles           map[string]AddonProfile  `json:"addonProfiles,omitempty"`
	AzProfile               *AzProfile               `json:"azProfile,omitempty"`
}

// AddonProfile represents an addon for managed cluster
type AddonProfile struct {
	Enabled bool              `json:"enabled"`
	Config  map[string]string `json:"config"`
}

// AzProfile holds the azure context for where the cluster resides
type AzProfile struct {
	TenantID       string `json:"tenantId,omitempty"`
	SubscriptionID string `json:"subscriptionId,omitempty"`
	ResourceGroup  string `json:"resourceGroup,omitempty"`
	Location       string `json:"location,omitempty"`
}

// ServicePrincipalProfile contains the client and secret used by the cluster for Azure Resource CRUD
type ServicePrincipalProfile struct {
	ClientID          string             `json:"clientId"`
	Secret            string             `json:"secret,omitempty"`
	KeyvaultSecretRef *KeyvaultSecretRef `json:"keyvaultSecretRef,omitempty"`
}

// KeyvaultSecretRef specifies path to the Azure keyvault along with secret name and (optionaly) version
// for Service Principal's secret
type KeyvaultSecretRef struct {
	VaultID       string `json:"vaultID"`
	SecretName    string `json:"secretName"`
	SecretVersion string `json:"version,omitempty"`
}

// CertificateProfile represents the definition of the master cluster
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
	// EtcdServerCertificate is the server certificate for etcd, and signed by the CA
	EtcdServerCertificate string `json:"etcdServerCertificate,omitempty"`
	// EtcdServerPrivateKey is the server private key for etcd, and signed by the CA
	EtcdServerPrivateKey string `json:"etcdServerPrivateKey,omitempty"`
	// EtcdClientCertificate is etcd client certificate, and signed by the CA
	EtcdClientCertificate string `json:"etcdClientCertificate,omitempty"`
	// EtcdClientPrivateKey is the etcd client private key, and signed by the CA
	EtcdClientPrivateKey string `json:"etcdClientPrivateKey,omitempty"`
	// EtcdPeerCertificates is list of etcd peer certificates, and signed by the CA
	EtcdPeerCertificates []string `json:"etcdPeerCertificates,omitempty"`
	// EtcdPeerPrivateKeys is list of etcd peer private keys, and signed by the CA
	EtcdPeerPrivateKeys []string `json:"etcdPeerPrivateKeys,omitempty"`
}

// LinuxProfile represents the linux parameters passed to the cluster
type LinuxProfile struct {
	AdminUsername string `json:"adminUsername"`
	SSH           struct {
		PublicKeys []PublicKey `json:"publicKeys"`
	} `json:"ssh"`
	Secrets       []KeyVaultSecrets `json:"secrets,omitempty"`
	Distro        Distro            `json:"distro,omitempty"`
	ScriptRootURL string            `json:"scriptroot,omitempty"`
}

// PublicKey represents an SSH key for LinuxProfile
type PublicKey struct {
	KeyData string `json:"keyData"`
}

// WindowsProfile represents the windows parameters passed to the cluster
type WindowsProfile struct {
	AdminUsername         string            `json:"adminUsername"`
	AdminPassword         string            `json:"adminPassword"`
	ImageVersion          string            `json:"imageVersion"`
	WindowsImageSourceURL string            `json:"windowsImageSourceURL"`
	WindowsPublisher      string            `json:"windowsPublisher"`
	WindowsOffer          string            `json:"windowsOffer"`
	WindowsSku            string            `json:"windowsSku"`
	Secrets               []KeyVaultSecrets `json:"secrets,omitempty"`
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
	// Upgrading means an existing ContainerService resource is being upgraded
	Upgrading ProvisioningState = "Upgrading"
)

// OrchestratorProfile contains Orchestrator properties
type OrchestratorProfile struct {
	OrchestratorType    string            `json:"orchestratorType"`
	OrchestratorVersion string            `json:"orchestratorVersion"`
	KubernetesConfig    *KubernetesConfig `json:"kubernetesConfig,omitempty"`
	OpenShiftConfig     *OpenShiftConfig  `json:"openshiftConfig,omitempty"`
	DcosConfig          *DcosConfig       `json:"dcosConfig,omitempty"`
}

// OrchestratorVersionProfile contains information of a supported orchestrator version:
type OrchestratorVersionProfile struct {
	// Orchestrator type and version
	OrchestratorProfile
	// Whether this orchestrator version is deployed by default if orchestrator release is not specified
	Default bool `json:"default,omitempty"`
	// List of available upgrades for this orchestrator version
	Upgrades []*OrchestratorProfile `json:"upgrades,omitempty"`
}

// KubernetesContainerSpec defines configuration for a container spec
type KubernetesContainerSpec struct {
	Name           string `json:"name,omitempty"`
	Image          string `json:"image,omitempty"`
	CPURequests    string `json:"cpuRequests,omitempty"`
	MemoryRequests string `json:"memoryRequests,omitempty"`
	CPULimits      string `json:"cpuLimits,omitempty"`
	MemoryLimits   string `json:"memoryLimits,omitempty"`
}

// KubernetesAddon defines a list of addons w/ configuration to include with the cluster deployment
type KubernetesAddon struct {
	Name       string                    `json:"name,omitempty"`
	Enabled    *bool                     `json:"enabled,omitempty"`
	Containers []KubernetesContainerSpec `json:"containers,omitempty"`
	Config     map[string]string         `json:"config,omitempty"`
}

// IsEnabled returns if the addon is explicitly enabled, or the user-provided default if non explicitly enabled
func (a *KubernetesAddon) IsEnabled(ifNil bool) bool {
	if a.Enabled == nil {
		return ifNil
	}
	return *a.Enabled
}

// PrivateCluster defines the configuration for a private cluster
type PrivateCluster struct {
	Enabled        *bool                  `json:"enabled,omitempty"`
	JumpboxProfile *PrivateJumpboxProfile `json:"jumpboxProfile,omitempty"`
}

// PrivateJumpboxProfile represents a jumpbox definition
type PrivateJumpboxProfile struct {
	Name           string `json:"name" validate:"required"`
	VMSize         string `json:"vmSize" validate:"required"`
	OSDiskSizeGB   int    `json:"osDiskSizeGB,omitempty" validate:"min=0,max=1023"`
	Username       string `json:"username,omitempty"`
	PublicKey      string `json:"publicKey" validate:"required"`
	StorageProfile string `json:"storageProfile,omitempty"`
}

// CloudProviderConfig contains the KubernetesConfig properties specific to the Cloud Provider
// TODO use this when strict JSON checking accommodates struct embedding
type CloudProviderConfig struct {
	CloudProviderBackoff         bool    `json:"cloudProviderBackoff,omitempty"`
	CloudProviderBackoffRetries  int     `json:"cloudProviderBackoffRetries,omitempty"`
	CloudProviderBackoffJitter   float64 `json:"cloudProviderBackoffJitter,omitempty"`
	CloudProviderBackoffDuration int     `json:"cloudProviderBackoffDuration,omitempty"`
	CloudProviderBackoffExponent float64 `json:"cloudProviderBackoffExponent,omitempty"`
	CloudProviderRateLimit       bool    `json:"cloudProviderRateLimit,omitempty"`
	CloudProviderRateLimitQPS    float64 `json:"cloudProviderRateLimitQPS,omitempty"`
	CloudProviderRateLimitBucket int     `json:"cloudProviderRateLimitBucket,omitempty"`
}

// KubernetesConfigDeprecated are properties that are no longer operable and will be ignored
// TODO use this when strict JSON checking accommodates struct embedding
type KubernetesConfigDeprecated struct {
	NonMasqueradeCidr                string `json:"nonMasqueradeCidr,omitempty"`
	NodeStatusUpdateFrequency        string `json:"nodeStatusUpdateFrequency,omitempty"`
	HardEvictionThreshold            string `json:"hardEvictionThreshold,omitempty"`
	CtrlMgrNodeMonitorGracePeriod    string `json:"ctrlMgrNodeMonitorGracePeriod,omitempty"`
	CtrlMgrPodEvictionTimeout        string `json:"ctrlMgrPodEvictionTimeout,omitempty"`
	CtrlMgrRouteReconciliationPeriod string `json:"ctrlMgrRouteReconciliationPeriod,omitempty"`
}

// KubernetesConfig contains the Kubernetes config structure, containing
// Kubernetes specific configuration
type KubernetesConfig struct {
	KubernetesImageBase              string            `json:"kubernetesImageBase,omitempty"`
	ClusterSubnet                    string            `json:"clusterSubnet,omitempty"`
	NetworkPolicy                    string            `json:"networkPolicy,omitempty"`
	ContainerRuntime                 string            `json:"containerRuntime,omitempty"`
	MaxPods                          int               `json:"maxPods,omitempty"`
	DockerBridgeSubnet               string            `json:"dockerBridgeSubnet,omitempty"`
	DNSServiceIP                     string            `json:"dnsServiceIP,omitempty"`
	ServiceCIDR                      string            `json:"serviceCidr,omitempty"`
	UseManagedIdentity               bool              `json:"useManagedIdentity,omitempty"`
	CustomHyperkubeImage             string            `json:"customHyperkubeImage,omitempty"`
	DockerEngineVersion              string            `json:"dockerEngineVersion,omitempty"`
	CustomCcmImage                   string            `json:"customCcmImage,omitempty"` // Image for cloud-controller-manager
	UseCloudControllerManager        *bool             `json:"useCloudControllerManager,omitempty"`
	UseInstanceMetadata              *bool             `json:"useInstanceMetadata,omitempty"`
	EnableRbac                       *bool             `json:"enableRbac,omitempty"`
	EnableSecureKubelet              *bool             `json:"enableSecureKubelet,omitempty"`
	EnableAggregatedAPIs             bool              `json:"enableAggregatedAPIs,omitempty"`
	PrivateCluster                   *PrivateCluster   `json:"privateCluster,omitempty"`
	GCHighThreshold                  int               `json:"gchighthreshold,omitempty"`
	GCLowThreshold                   int               `json:"gclowthreshold,omitempty"`
	EtcdVersion                      string            `json:"etcdVersion,omitempty"`
	EtcdDiskSizeGB                   string            `json:"etcdDiskSizeGB,omitempty"`
	EnableDataEncryptionAtRest       *bool             `json:"enableDataEncryptionAtRest,omitempty"`
	EnablePodSecurityPolicy          *bool             `json:"enablePodSecurityPolicy,omitempty"`
	Addons                           []KubernetesAddon `json:"addons,omitempty"`
	KubeletConfig                    map[string]string `json:"kubeletConfig,omitempty"`
	ControllerManagerConfig          map[string]string `json:"controllerManagerConfig,omitempty"`
	CloudControllerManagerConfig     map[string]string `json:"cloudControllerManagerConfig,omitempty"`
	APIServerConfig                  map[string]string `json:"apiServerConfig,omitempty"`
	SchedulerConfig                  map[string]string `json:"schedulerConfig,omitempty"`
	CloudProviderBackoff             bool              `json:"cloudProviderBackoff,omitempty"`
	CloudProviderBackoffRetries      int               `json:"cloudProviderBackoffRetries,omitempty"`
	CloudProviderBackoffJitter       float64           `json:"cloudProviderBackoffJitter,omitempty"`
	CloudProviderBackoffDuration     int               `json:"cloudProviderBackoffDuration,omitempty"`
	CloudProviderBackoffExponent     float64           `json:"cloudProviderBackoffExponent,omitempty"`
	CloudProviderRateLimit           bool              `json:"cloudProviderRateLimit,omitempty"`
	CloudProviderRateLimitQPS        float64           `json:"cloudProviderRateLimitQPS,omitempty"`
	CloudProviderRateLimitBucket     int               `json:"cloudProviderRateLimitBucket,omitempty"`
	NonMasqueradeCidr                string            `json:"nonMasqueradeCidr,omitempty"`
	NodeStatusUpdateFrequency        string            `json:"nodeStatusUpdateFrequency,omitempty"`
	HardEvictionThreshold            string            `json:"hardEvictionThreshold,omitempty"`
	CtrlMgrNodeMonitorGracePeriod    string            `json:"ctrlMgrNodeMonitorGracePeriod,omitempty"`
	CtrlMgrPodEvictionTimeout        string            `json:"ctrlMgrPodEvictionTimeout,omitempty"`
	CtrlMgrRouteReconciliationPeriod string            `json:"ctrlMgrRouteReconciliationPeriod,omitempty"`
}

// DcosConfig Configuration for DC/OS
type DcosConfig struct {
	DcosBootstrapURL         string `json:"dcosBootstrapURL,omitempty"`
	DcosWindowsBootstrapURL  string `json:"dcosWindowsBootstrapURL,omitempty"`
	Registry                 string `json:"registry,omitempty"`
	RegistryUser             string `json:"registryUser,omitempty"`
	RegistryPass             string `json:"registryPassword,omitempty"`
	DcosRepositoryURL        string `json:"dcosRepositoryURL,omitempty"`        // For CI use, you need to specify
	DcosClusterPackageListID string `json:"dcosClusterPackageListID,omitempty"` // all three of these items
	DcosProviderPackageID    string `json:"dcosProviderPackageID,omitempty"`    // repo url is the location of the build,
}

// OpenShiftConfig holds configuration for OpenShift
type OpenShiftConfig struct {
	KubernetesConfig *KubernetesConfig `json:"kubernetesConfig,omitempty"`

	// ClusterUsername and ClusterPassword are temporary before AAD
	// authentication is enabled, and will be removed subsequently.
	ClusterUsername string `json:"clusterUsername,omitempty"`
	ClusterPassword string `json:"clusterPassword,omitempty"`

	ConfigBundles          map[string][]byte `json:"-"`
	ExternalMasterHostname string            `json:"-"`
	RouterLBHostname       string            `json:"-"`
}

// MasterProfile represents the definition of the master cluster
type MasterProfile struct {
	Count                    int               `json:"count"`
	DNSPrefix                string            `json:"dnsPrefix"`
	VMSize                   string            `json:"vmSize"`
	OSDiskSizeGB             int               `json:"osDiskSizeGB,omitempty"`
	VnetSubnetID             string            `json:"vnetSubnetID,omitempty"`
	VnetCidr                 string            `json:"vnetCidr,omitempty"`
	FirstConsecutiveStaticIP string            `json:"firstConsecutiveStaticIP,omitempty"`
	Subnet                   string            `json:"subnet"`
	IPAddressCount           int               `json:"ipAddressCount,omitempty"`
	StorageProfile           string            `json:"storageProfile,omitempty"`
	HTTPSourceAddressPrefix  string            `json:"HTTPSourceAddressPrefix,omitempty"`
	OAuthEnabled             bool              `json:"oauthEnabled"`
	PreprovisionExtension    *Extension        `json:"preProvisionExtension"`
	Extensions               []Extension       `json:"extensions"`
	Distro                   Distro            `json:"distro,omitempty"`
	KubernetesConfig         *KubernetesConfig `json:"kubernetesConfig,omitempty"`
	ImageRef                 *ImageReference   `json:"imageReference,omitempty"`

	// Master LB public endpoint/FQDN with port
	// The format will be FQDN:2376
	// Not used during PUT, returned as part of GET
	FQDN string `json:"fqdn,omitempty"`
}

// ImageReference represents a reference to an Image resource in Azure.
type ImageReference struct {
	Name          string `json:"name,omitempty"`
	ResourceGroup string `json:"resourceGroup,omitempty"`
}

// ExtensionProfile represents an extension definition
type ExtensionProfile struct {
	Name                           string             `json:"name"`
	Version                        string             `json:"version"`
	ExtensionParameters            string             `json:"extensionParameters,omitempty"`
	ExtensionParametersKeyVaultRef *KeyvaultSecretRef `json:"parametersKeyvaultSecretRef,omitempty"`
	RootURL                        string             `json:"rootURL,omitempty"`
	// This is only needed for preprovision extensions and it needs to be a bash script
	Script   string `json:"script,omitempty"`
	URLQuery string `json:"urlQuery,omitempty"`
}

// Extension represents an extension definition in the master or agentPoolProfile
type Extension struct {
	Name        string `json:"name"`
	SingleOrAll string `json:"singleOrAll"`
	Template    string `json:"template"`
}

// AgentPoolProfile represents an agent pool definition
type AgentPoolProfile struct {
	Name                string               `json:"name"`
	Count               int                  `json:"count"`
	VMSize              string               `json:"vmSize"`
	OSDiskSizeGB        int                  `json:"osDiskSizeGB,omitempty"`
	DNSPrefix           string               `json:"dnsPrefix,omitempty"`
	OSType              OSType               `json:"osType,omitempty"`
	Ports               []int                `json:"ports,omitempty"`
	AvailabilityProfile string               `json:"availabilityProfile"`
	StorageProfile      string               `json:"storageProfile,omitempty"`
	DiskSizesGB         []int                `json:"diskSizesGB,omitempty"`
	VnetSubnetID        string               `json:"vnetSubnetID,omitempty"`
	Subnet              string               `json:"subnet"`
	IPAddressCount      int                  `json:"ipAddressCount,omitempty"`
	Distro              Distro               `json:"distro,omitempty"`
	Role                AgentPoolProfileRole `json:"role,omitempty"`

	FQDN                  string            `json:"fqdn,omitempty"`
	CustomNodeLabels      map[string]string `json:"customNodeLabels,omitempty"`
	PreprovisionExtension *Extension        `json:"preProvisionExtension"`
	Extensions            []Extension       `json:"extensions"`
	KubernetesConfig      *KubernetesConfig `json:"kubernetesConfig,omitempty"`
	ImageRef              *ImageReference   `json:"imageReference,omitempty"`
}

// AgentPoolProfileRole represents an agent role
type AgentPoolProfileRole string

// DiagnosticsProfile setting to enable/disable capturing
// diagnostics for VMs hosting container cluster.
type DiagnosticsProfile struct {
	VMDiagnostics *VMDiagnostics `json:"vmDiagnostics"`
}

// VMDiagnostics contains settings to on/off boot diagnostics collection
// in RD Host
type VMDiagnostics struct {
	Enabled bool `json:"enabled"`

	// Specifies storage account Uri where Boot Diagnostics (CRP &
	// VMSS BootDiagostics) and VM Diagnostics logs (using Linux
	// Diagnostics Extension) will be stored. Uri will be of standard
	// blob domain. i.e. https://storageaccount.blob.core.windows.net/
	// This field is readonly as ACS RP will create a storage account
	// for the customer.
	StorageURL *neturl.URL `json:"storageUrl"`
}

// JumpboxProfile describes properties of the jumpbox setup
// in the ACS container cluster.
type JumpboxProfile struct {
	OSType    OSType `json:"osType"`
	DNSPrefix string `json:"dnsPrefix"`

	// Jumpbox public endpoint/FQDN with port
	// The format will be FQDN:2376
	// Not used during PUT, returned as part of GET
	FQDN string `json:"fqdn,omitempty"`
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

// Distro represents Linux distro to use for Linux VMs
type Distro string

// HostedMasterProfile defines properties for a hosted master
type HostedMasterProfile struct {
	// Master public endpoint/FQDN with port
	// The format will be FQDN:2376
	// Not used during PUT, returned as part of GETFQDN
	FQDN      string `json:"fqdn,omitempty"`
	DNSPrefix string `json:"dnsPrefix"`
	// Subnet holds the CIDR which defines the Azure Subnet in which
	// Agents will be provisioned. This is stored on the HostedMasterProfile
	// and will become `masterSubnet` in the compiled template.
	Subnet string `json:"subnet"`
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
	// The Azure Active Directory Group Object ID that will be assigned the
	// cluster-admin RBAC role.
	// Optional
	AdminGroupID string `json:"adminGroupID,omitempty"`
}

// CustomProfile specifies custom properties that are used for
// cluster instantiation.  Should not be used by most users.
type CustomProfile struct {
	Orchestrator string `json:"orchestrator,omitempty"`
}

// VlabsARMContainerService is the type we read and write from file
// needed because the json that is sent to ARM and acs-engine
// is different from the json that the ACS RP Api gets from ARM
type VlabsARMContainerService struct {
	TypeMeta
	*vlabs.ContainerService
}

// V20160330ARMContainerService is the type we read and write from file
// needed because the json that is sent to ARM and acs-engine
// is different from the json that the ACS RP Api gets from ARM
type V20160330ARMContainerService struct {
	TypeMeta
	*v20160330.ContainerService
}

// V20160930ARMContainerService is the type we read and write from file
// needed because the json that is sent to ARM and acs-engine
// is different from the json that the ACS RP Api gets from ARM
type V20160930ARMContainerService struct {
	TypeMeta
	*v20160930.ContainerService
}

// V20170131ARMContainerService is the type we read and write from file
// needed because the json that is sent to ARM and acs-engine
// is different from the json that the ACS RP Api gets from ARM
type V20170131ARMContainerService struct {
	TypeMeta
	*v20170131.ContainerService
}

// V20170701ARMContainerService is the type we read and write from file
// needed because the json that is sent to ARM and acs-engine
// is different from the json that the ACS RP Api gets from ARM
type V20170701ARMContainerService struct {
	TypeMeta
	*v20170701.ContainerService
}

// V20170831ARMManagedContainerService is the type we read and write from file
// needed because the json that is sent to ARM and acs-engine
// is different from the json that the ACS RP Api gets from ARM
type V20170831ARMManagedContainerService struct {
	TypeMeta
	*v20170831.ManagedCluster
}

// V20180331ARMManagedContainerService is the type we read and write from file
// needed because the json that is sent to ARM and acs-engine
// is different from the json that the ACS RP Api gets from ARM
type V20180331ARMManagedContainerService struct {
	TypeMeta
	*v20180331.ManagedCluster
}

// HasWindows returns true if the cluster contains windows
func (p *Properties) HasWindows() bool {
	for _, agentPoolProfile := range p.AgentPoolProfiles {
		if agentPoolProfile.OSType == Windows {
			return true
		}
	}
	return false
}

// HasManagedDisks returns true if the cluster contains Managed Disks
func (p *Properties) HasManagedDisks() bool {
	if p.MasterProfile != nil && p.MasterProfile.StorageProfile == ManagedDisks {
		return true
	}
	for _, agentPoolProfile := range p.AgentPoolProfiles {
		if agentPoolProfile.StorageProfile == ManagedDisks {
			return true
		}
	}
	if p.OrchestratorProfile.KubernetesConfig.PrivateJumpboxProvision() && p.OrchestratorProfile.KubernetesConfig.PrivateCluster.JumpboxProfile.StorageProfile == ManagedDisks {
		return true
	}
	return false
}

// HasStorageAccountDisks returns true if the cluster contains Storage Account Disks
func (p *Properties) HasStorageAccountDisks() bool {
	if p.OrchestratorProfile.OrchestratorType == OpenShift {
		return true
	}
	if p.MasterProfile != nil && p.MasterProfile.StorageProfile == StorageAccount {
		return true
	}
	for _, agentPoolProfile := range p.AgentPoolProfiles {
		if agentPoolProfile.StorageProfile == StorageAccount {
			return true
		}
	}
	if p.OrchestratorProfile.KubernetesConfig.PrivateJumpboxProvision() && p.OrchestratorProfile.KubernetesConfig.PrivateCluster.JumpboxProfile.StorageProfile == StorageAccount {
		return true
	}
	return false
}

// TotalNodes returns the total number of nodes in the cluster configuration
func (p *Properties) TotalNodes() int {
	var totalNodes int
	if p.MasterProfile != nil {
		totalNodes = p.MasterProfile.Count
	}
	for _, pool := range p.AgentPoolProfiles {
		totalNodes = totalNodes + pool.Count
	}
	return totalNodes
}

// IsCustomVNET returns true if the customer brought their own VNET
func (m *MasterProfile) IsCustomVNET() bool {
	return len(m.VnetSubnetID) > 0
}

// IsManagedDisks returns true if the master specified managed disks
func (m *MasterProfile) IsManagedDisks() bool {
	return m.StorageProfile == ManagedDisks
}

// IsStorageAccount returns true if the master specified storage account
func (m *MasterProfile) IsStorageAccount() bool {
	return m.StorageProfile == StorageAccount
}

// IsRHEL returns true if the master specified a RHEL distro
func (m *MasterProfile) IsRHEL() bool {
	return m.Distro == RHEL
}

// IsCoreOS returns true if the master specified a CoreOS distro
func (m *MasterProfile) IsCoreOS() bool {
	return m.Distro == CoreOS
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

// IsRHEL returns true if the agent pool specified a RHEL distro
func (a *AgentPoolProfile) IsRHEL() bool {
	return a.OSType == Linux && a.Distro == RHEL
}

// IsCoreOS returns true if the agent specified a CoreOS distro
func (a *AgentPoolProfile) IsCoreOS() bool {
	return a.OSType == Linux && a.Distro == CoreOS
}

// IsAvailabilitySets returns true if the customer specified disks
func (a *AgentPoolProfile) IsAvailabilitySets() bool {
	return a.AvailabilityProfile == AvailabilitySet
}

// IsManagedDisks returns true if the customer specified disks
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

// HasSecrets returns true if the customer specified secrets to install
func (w *WindowsProfile) HasSecrets() bool {
	return len(w.Secrets) > 0
}

// HasCustomImage returns true if there is a custom windows os image url specified
func (w *WindowsProfile) HasCustomImage() bool {
	return len(w.WindowsImageSourceURL) > 0
}

// HasSecrets returns true if the customer specified secrets to install
func (l *LinuxProfile) HasSecrets() bool {
	return len(l.Secrets) > 0
}

// IsSwarmMode returns true if this template is for Swarm Mode orchestrator
func (o *OrchestratorProfile) IsSwarmMode() bool {
	return o.OrchestratorType == SwarmMode
}

// IsKubernetes returns true if this template is for Kubernetes orchestrator
func (o *OrchestratorProfile) IsKubernetes() bool {
	return o.OrchestratorType == Kubernetes
}

// IsOpenShift returns true if this template is for OpenShift orchestrator
func (o *OrchestratorProfile) IsOpenShift() bool {
	return o.OrchestratorType == OpenShift
}

// IsDCOS returns true if this template is for DCOS orchestrator
func (o *OrchestratorProfile) IsDCOS() bool {
	return o.OrchestratorType == DCOS
}

// IsAzureCNI returns true if Azure VNET integration is enabled
func (o *OrchestratorProfile) IsAzureCNI() bool {
	switch o.OrchestratorType {
	case Kubernetes:
		return o.KubernetesConfig.NetworkPolicy == "azure"
	default:
		return false
	}
}

// HasAadProfile  returns true if the has aad profile
func (p *Properties) HasAadProfile() bool {
	return p.AADProfile != nil
}

// GetAPIServerEtcdAPIVersion Used to set apiserver's etcdapi version
func (o *OrchestratorProfile) GetAPIServerEtcdAPIVersion() string {
	ret := "etcd3"
	if o.KubernetesConfig != nil {
		// if we are here, version has already been validated..
		etcdversion, _ := semver.NewVersion(o.KubernetesConfig.EtcdVersion)
		if etcdversion != nil && 2 == etcdversion.Major() {
			return "etcd2"
		}
	}
	return ret
}

// IsMetricsServerEnabled checks if the metrics server addon is enabled
func (o *OrchestratorProfile) IsMetricsServerEnabled() bool {
	var metricsServerAddon KubernetesAddon
	k := o.KubernetesConfig
	for i := range k.Addons {
		if k.Addons[i].Name == DefaultMetricsServerAddonName {
			metricsServerAddon = k.Addons[i]
		}
	}
	return metricsServerAddon.IsEnabled(DefaultMetricsServerAddonEnabled || common.IsKubernetesVersionGe(o.OrchestratorVersion, "1.9.0"))
}

// IsTillerEnabled checks if the tiller addon is enabled
func (k *KubernetesConfig) IsTillerEnabled() bool {
	var tillerAddon KubernetesAddon
	for i := range k.Addons {
		if k.Addons[i].Name == DefaultTillerAddonName {
			tillerAddon = k.Addons[i]
		}
	}
	return tillerAddon.IsEnabled(DefaultTillerAddonEnabled)
}

// IsACIConnectorEnabled checks if the ACI Connector addon is enabled
func (k *KubernetesConfig) IsACIConnectorEnabled() bool {
	var aciConnectorAddon KubernetesAddon
	for i := range k.Addons {
		if k.Addons[i].Name == DefaultACIConnectorAddonName {
			aciConnectorAddon = k.Addons[i]
		}
	}
	return aciConnectorAddon.IsEnabled(DefaultACIConnectorAddonEnabled)
}

// IsDashboardEnabled checks if the kubernetes-dashboard addon is enabled
func (k *KubernetesConfig) IsDashboardEnabled() bool {
	var dashboardAddon KubernetesAddon
	for i := range k.Addons {
		if k.Addons[i].Name == DefaultDashboardAddonName {
			dashboardAddon = k.Addons[i]
		}
	}
	return dashboardAddon.IsEnabled(DefaultDashboardAddonEnabled)
}

// IsReschedulerEnabled checks if the rescheduler addon is enabled
func (k *KubernetesConfig) IsReschedulerEnabled() bool {
	var reschedulerAddon KubernetesAddon
	for i := range k.Addons {
		if k.Addons[i].Name == DefaultReschedulerAddonName {
			reschedulerAddon = k.Addons[i]
		}
	}
	return reschedulerAddon.IsEnabled(DefaultReschedulerAddonEnabled)
}

// PrivateJumpboxProvision checks if a private cluster has jumpbox auto-provisioning
func (k *KubernetesConfig) PrivateJumpboxProvision() bool {
	if k != nil && k.PrivateCluster != nil && *k.PrivateCluster.Enabled && k.PrivateCluster.JumpboxProfile != nil {
		return true
	}
	return false
}

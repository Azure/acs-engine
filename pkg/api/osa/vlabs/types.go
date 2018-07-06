package vlabs

// OpenShiftCluster complies with the ARM model of resource definition in a JSON
// template.
type OpenShiftCluster struct {
	ID         string                `json:"id,omitempty"`
	Location   string                `json:"location,omitempty"`
	Name       string                `json:"name,omitempty"`
	Plan       *ResourcePurchasePlan `json:"plan,omitempty"`
	Tags       map[string]string     `json:"tags,omitempty"`
	Type       string                `json:"type,omitempty"`
	Properties *Properties           `json:"properties,omitempty"`
}

// ResourcePurchasePlan defines the resource plan as required by ARM for billing
// purposes.
type ResourcePurchasePlan struct {
	Name          string `json:"name,omitempty"`
	Product       string `json:"product,omitempty"`
	PromotionCode string `json:"promotionCode,omitempty"`
	Publisher     string `json:"publisher,omitempty"`
}

// Properties represents the cluster definition.
type Properties struct {
	// ProvisioningState (out): current state of the OSA resource.
	ProvisioningState ProvisioningState `json:"provisioningState,omitempty"`

	// OpenShiftVersion (in): OpenShift version to be created/updated, e.g.
	// `v3.10`.
	OpenShiftVersion string `json:"openShiftVersion,omitempty"`

	// TODO: the following DNS-related fields are basically workable but don't
	// feel very elegant or future-proof (e.g. multiple routers, non-wildcard
	// setups).

	// PublicHostname (in,optional): Optional user-specified FQDN for OpenShift
	// API server.  If specified, after OSA cluster creation, user must create a
	// PublicHostname CNAME record forwarding to the returned FQDN value.
	PublicHostname string `json:"publicHostname,omitempty"`

	// FQDN (out): Auto-allocated FQDN for OpenShift API server.
	FQDN string `json:"fqdn,omitempty"`

	// RoutingConfigSubdomain (in,optional/out): DNS subdomain for OpenShift
	// router.  If specified, after OSA cluster creation, user must create a
	// (wildcard) *.RoutingConfigSubdomain CNAME record forwarding to the
	// returned RoutingConfigFQDN value.  If not specified, OSA will
	// auto-allocate and setup a RoutingConfigSubdomain and return it.
	RoutingConfigSubdomain string `json:"routingConfigSubdomain,omitempty"`

	// RoutingConfigFQDN (out): Auto-allocated FQDN for the OpenShift router.
	RoutingConfigFQDN string `json:"routingConfigFqdn,omitempty"`

	// TODO: need to clarify the external API for AgentPoolProfiles.  Will we
	// allow users to specify any aspects of non-`compute` pools?
	// AgentPoolProfiles (in): configuration of OpenShift cluster VMs.
	AgentPoolProfiles AgentPoolProfiles `json:"agentPoolProfiles,omitempty"`

	// TODO: is this compatible with MSI?
	// ServicePrincipalProfile (in): Service principal for OpenShift cluster.
	ServicePrincipalProfile ServicePrincipalProfile `json:"servicePrincipalProfile,omitempty"`
}

// ProvisioningState represents the current state of the OSA resource.
type ProvisioningState string

const (
	// Creating means the OSA resource is being created.
	Creating ProvisioningState = "Creating"
	// Updating means the existing OSA resource is being updated.
	Updating ProvisioningState = "Updating"
	// Failed means the OSA resource is in failed state.
	Failed ProvisioningState = "Failed"
	// Succeeded means the last create/update succeeded.
	Succeeded ProvisioningState = "Succeeded"
	// Deleting means the OSA resource is being deleted.
	Deleting ProvisioningState = "Deleting"
	// Migrating means the OSA resource is being migrated from one subscription
	// or resource group to another.
	Migrating ProvisioningState = "Migrating"
	// Upgrading means the existing OAS resource is being upgraded.
	Upgrading ProvisioningState = "Upgrading"
)

// AgentPoolProfiles represents all the AgentPoolProfiles.
type AgentPoolProfiles []AgentPoolProfile

// AgentPoolProfile represents configuration of OpenShift cluster VMs.
type AgentPoolProfile struct {
	Name   string               `json:"name,omitempty"`
	Role   AgentPoolProfileRole `json:"role,omitempty"`
	Count  int                  `json:"count,omitempty"`
	VMSize string               `json:"vmSize,omitempty"`

	// VnetSubnetID is expected to be empty or match
	// `^/subscriptions/[^/]+
	//   /resourceGroups/[^/]+
	//   /providers/Microsoft.Network
	//   /virtualNetworks/[^/]+
	//   /subnets/[^/]+$`
	VnetSubnetID string `json:"vnetSubnetID,omitempty"`
	OSType       OSType `json:"osType,omitempty"`
}

// AgentPoolProfileRole represents the role of the AgentPoolProfile.
type AgentPoolProfileRole string

const (
	// AgentPoolProfileRoleCompute is the compute role.
	AgentPoolProfileRoleCompute AgentPoolProfileRole = "compute"
	// AgentPoolProfileRoleInfra is the infra role.
	AgentPoolProfileRoleInfra AgentPoolProfileRole = "infra"
	// AgentPoolProfileRoleMaster is the master role.
	AgentPoolProfileRoleMaster AgentPoolProfileRole = "master"
)

// OSType represents the OS type of VMs in an AgentPool.
type OSType string

const (
	// OSTypeLinux is Linux.
	OSTypeLinux OSType = "Linux"
	// OSTypeWindows is Windows.
	OSTypeWindows OSType = "Windows"
)

// ServicePrincipalProfile contains the client and secret used by the cluster
// for Azure Resource CRUD.
type ServicePrincipalProfile struct {
	ClientID string `json:"clientId,omitempty"`
	Secret   string `json:"secret,omitempty"`
}

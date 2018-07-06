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
	Properties Properties            `json:"properties,omitempty"`
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
	ProvisioningState ProvisioningState `json:"provisioningState,omitempty"`
	OpenShiftVersion  string            `json:"openShiftVersion,omitempty"`
	// TODO: need to clarify external API for PublicHostname and decide what to
	// implement when.  Allow users to specify a PublicHostname then have them
	// create a CNAME to a FQDN we return?  Allow users to not specify and we
	// return a FQDN?  In which case, how will the plugin know the FQDN?
	PublicHostname string `json:"publicHostname,omitempty"`
	// FQDN string `json:"fqdn,omitempty"` // TODO: do we need to add this?
	// TODO: need to clarify external API for RoutingConfigSubdomain.  Do we
	// create one and return it if it's not provided?  Will this be transparent
	// to the plugin?
	RoutingConfigSubdomain string `json:"routingConfigSubdomain,omitempty"`
	// TODO: need to clarify the external API for AgentPoolProfiles.  Will we
	// allow users to specify an `infra` pool?
	AgentPoolProfiles       AgentPoolProfiles       `json:"agentPoolProfiles,omitempty"`
	ServicePrincipalProfile ServicePrincipalProfile `json:"servicePrincipalProfile,omitempty"`
}

// ProvisioningState represents the current state of container service resource.
type ProvisioningState string

const (
	// Creating means ContainerService resource is being created.
	Creating ProvisioningState = "Creating"
	// Updating means an existing ContainerService resource is being updated.
	Updating ProvisioningState = "Updating"
	// Failed means resource is in failed state.
	Failed ProvisioningState = "Failed"
	// Succeeded means last create/update succeeded.
	Succeeded ProvisioningState = "Succeeded"
	// Deleting means resource is in the process of being deleted.
	Deleting ProvisioningState = "Deleting"
	// Migrating means resource is being migrated from one subscription or
	// resource group to another.
	Migrating ProvisioningState = "Migrating"
	// Upgrading means an existing resource is being upgraded.
	Upgrading ProvisioningState = "Upgrading"
)

// AgentPoolProfiles represents all the AgentPoolProfiles.
type AgentPoolProfiles []AgentPoolProfile

// AgentPoolProfile represents configuration of VMs running agent daemons that
// register with the master and offer resources to host applications in
// containers.
type AgentPoolProfile struct {
	Name         string               `json:"name,omitempty"`
	Role         AgentPoolProfileRole `json:"role,omitempty"` // TODO: should we expose this?
	Count        int                  `json:"count,omitempty"`
	VMSize       string               `json:"vmSize,omitempty"`
	VnetSubnetID string               `json:"vnetSubnetID,omitempty"`
	// OSDiskSizeGB int `json:"osDiskSizeGB,omitempty"` // TODO: do we need to add this?
	// AvailabilityProfile string `json:"availabilityProfile,omitempty"` // TODO: do we need to add this?
	// StorageProfile string `json:"storageProfile,omitempty"` // TODO: do we need to add this?
	// OSType OSType `json:"osType,omitempty"` // TODO: do we need to add this?
}

// AgentPoolProfileRole representes the role of the AgentPoolProfile.
// TODO: should we expose this?
type AgentPoolProfileRole string

const (
	// AgentPoolProfileRoleEmpty is the empty role
	AgentPoolProfileRoleEmpty AgentPoolProfileRole = ""
	// AgentPoolProfileRoleInfra is the infra role
	AgentPoolProfileRoleInfra AgentPoolProfileRole = "infra"
)

// ServicePrincipalProfile contains the client and secret used by the cluster
// for Azure Resource CRUD.
type ServicePrincipalProfile struct {
	ClientID string `json:"clientId,omitempty"`
	Secret   string `json:"secret,omitempty"`
}

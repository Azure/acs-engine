package v20160330

import (
	neturl "net/url"
)

// SubscriptionState represents the state of the subscription
type SubscriptionState int

// Subscription represents the customer subscription
type Subscription struct {
	ID    string
	State SubscriptionState
}

// OperationStatus represents the operation status
// And it should be written as the response body
// for operation status query
type OperationStatus struct {
	ID                 string                `json:"id,omitempty"`
	Name               string                `json:"name,omitempty"`
	Status             string                `json:"status,omitempty"`
	StartTime          string                `json:"startTime,omitempty"`          // “startTime”: “<DateLiteral per ISO8601>”
	EndTime            string                `json:"endTime,omitempty"`            // “endTime”: “<DateLiteral per ISO8601>”
	PercentageComplete float64               `json:"percentageComplete,omitempty"` // “percentComplete”: <double between 0 and 100>
	Properties         map[string]string     `json:"properties,omitempty"`         // Only return on succeeded case
	Error              ContainerServiceError `json:"error,omitempty"`
	SubscriptionID     string                `json:"subscriptionid,omitempty"`
	Location           string                `json:"location,omitempty"`
}

// ContainerServiceError conform the OData v4 error format
type ContainerServiceError struct {
	Code    string                 `json:"code,omitempty"`
	Message string                 `json:"message,omitempty"`
	Details *ContainerServiceError `json:"details,omitempty"`
}

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
	ID       string               `json:"id,omitempty"`
	Location string               `json:"location,omitempty"`
	Name     string               `json:"name,omitempty"`
	Plan     ResourcePurchasePlan `json:"plan,omitempty"`
	Tags     map[string]string    `json:"tags,omitempty"`
	Type     string               `json:"type,omitempty"`

	Properties Properties `json:"properties"`
}

// Properties is currently incomplete. More fields will be added later.
type Properties struct {
	ProvisioningState   ProvisioningState   `json:"provisioningState"`
	OrchestratorProfile OrchestratorProfile `json:"orchestratorProfile"`

	MasterProfile MasterProfile `json:"masterProfile"`

	AgentPoolProfiles []AgentPoolProfile `json:"agentPoolProfiles"`

	LinuxProfile LinuxProfile `json:"linuxProfile"`

	WindowsProfile WindowsProfile `json:"windowsProfile"`

	// TODO: This field is versioned to "2016-03-30"
	DiagnosticsProfile DiagnosticsProfile `json:"diagnosticsProfile"`

	// JumpboxProfile has made it into the new ACS RP stack for
	// backward compability.
	// TODO: Version this field so that newer versions don't
	// allow jumpbox creation
	JumpboxProfile JumpboxProfile `json:"jumpboxProfile"`
}

// LinuxProfile represents the Linux configuration passed to the cluster
type LinuxProfile struct {
	AdminUsername string `json:"adminUsername"`

	SSH struct {
		PublicKeys []struct {
			KeyData string `json:"keyData"`
		} `json:"publicKeys"`
	} `json:"ssh"`
}

// WindowsProfile represents the Windows configuration passed to the cluster
type WindowsProfile struct {
	AdminUsername string `json:"adminUsername"`

	AdminPassword string `json:"adminPassword"`
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
	OrchestratorType OrchestratorType `json:"orchestratorType"`
}

// MasterProfile represents the definition of master cluster
type MasterProfile struct {
	Count     int    `json:"count"`
	DNSPrefix string `json:"dnsPrefix"`

	// Master LB public endpoint/FQDN with port
	// The format will be FQDN:2376
	// Not used during PUT, returned as part of GET
	FQDN string `json:"fqdn,omitempty"`

	// subnet is internal
	subnet string
}

// AgentPoolProfile represents configuration of VMs running agent
// daemons that register with the master and offer resources to
// host applications in containers.
type AgentPoolProfile struct {
	Name      string `json:"name"`
	Count     int    `json:"count"`
	VMSize    string `json:"vmSize"`
	DNSPrefix string `json:"dnsPrefix"`
	FQDN      string `json:"fqdn,omitempty"`
	OSType    OSType `json:"osType"` // TODO: This field is versioned to "2016-03-30"
	// subnet is internal
	subnet string
}

// JumpboxProfile dscribes properties of the jumpbox setup
// in the ACS container cluster.
type JumpboxProfile struct {
	OSType    OSType `json:"osType"`
	DNSPrefix string `json:"dnsPrefix"`

	// Jumpbox public endpoint/FQDN with port
	// The format will be FQDN:2376
	// Not used during PUT, returned as part of GET
	FQDN string `json:"fqdn,omitempty"`
}

// DiagnosticsProfile setting to enable/disable capturing
// diagnostics for VMs hosting container cluster.
type DiagnosticsProfile struct {
	VMDiagnostics VMDiagnostics `json:"vmDiagnostics"`
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
	StorageURL neturl.URL `json:"storageUrl"`
}

// OrchestratorType defines orchestrators supported by ACS
type OrchestratorType string

// OSType represents OS types of agents
type OSType string

// HasWindows returns true if the cluster contains windows
func (a *Properties) HasWindows() bool {
	for _, agentPoolProfile := range a.AgentPoolProfiles {
		if agentPoolProfile.OSType == Windows {
			return true
		}
	}
	return false
}

// GetSubnet returns the read-only subnet for the master
func (m *MasterProfile) GetSubnet() string {
	return m.subnet
}

// SetSubnet sets the read-only subnet for the master
func (m *MasterProfile) SetSubnet(subnet string) {
	m.subnet = subnet
}

// IsWindows returns true if the agent pool is windows
func (a *AgentPoolProfile) IsWindows() bool {
	return a.OSType == Windows
}

// GetSubnet returns the read-only subnet for the agent pool
func (a *AgentPoolProfile) GetSubnet() string {
	return a.subnet
}

// SetSubnet sets the read-only subnet for the agent pool
func (a *AgentPoolProfile) SetSubnet(subnet string) {
	a.subnet = subnet
}

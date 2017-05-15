package armhelpers

import (
	"github.com/Azure/azure-sdk-for-go/arm/compute"
	"github.com/Azure/azure-sdk-for-go/arm/resources/resources"
)

// ACSEngineClient is the interface used to talk to an Azure environment.
// This interface exposes just the subset of Azure APIs and clients needed for
// ACS-Engine.
type ACSEngineClient interface {
	DeploymentsClient() DeploymentsClient // wraps the deployment client
	VirtualMachinesClient() VirtualMachinesClient
}

// DeploymentsClient exposes methods needed for handling Deployments
type DeploymentsClient interface {
	DeployTemplate(resourceGroup, name string, template, parameters map[string]interface{}, cancel <-chan struct{}) (*resources.DeploymentExtended, error)
}

// VirtualMachinesClient exposes methods needed for handling VirtualMachines
type VirtualMachinesClient interface {
	// List lists VM resources
	ListVirtualMachines(resourceGroup string) (compute.VirtualMachineListResult, error)
}

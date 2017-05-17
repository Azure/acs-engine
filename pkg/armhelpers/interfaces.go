package armhelpers

import (
	"github.com/Azure/azure-sdk-for-go/arm/compute"
	"github.com/Azure/azure-sdk-for-go/arm/resources/resources"
)

// ACSEngineClient is the interface used to talk to an Azure environment.
// This interface exposes just the subset of Azure APIs and clients needed for
// ACS-Engine.
type ACSEngineClient interface {
	// DeployTemplate can deploy a template into Azure ARM
	DeployTemplate(resourceGroup, name string, template, parameters map[string]interface{}, cancel <-chan struct{}) (*resources.DeploymentExtended, error)

	// List lists VM resources
	ListVirtualMachines(resourceGroup string) (compute.VirtualMachineListResult, error)
}

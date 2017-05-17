package armhelpers

import (
	"github.com/Azure/azure-sdk-for-go/arm/compute"
	"github.com/Azure/azure-sdk-for-go/arm/resources/resources"
	"github.com/Azure/go-autorest/autorest"
)

// ACSEngineClient is the interface used to talk to an Azure environment.
// This interface exposes just the subset of Azure APIs and clients needed for
// ACS-Engine.
type ACSEngineClient interface {
	//
	// RESOURCES

	// DeployTemplate can deploy a template into Azure ARM
	DeployTemplate(resourceGroup, name string, template, parameters map[string]interface{}, cancel <-chan struct{}) (*resources.DeploymentExtended, error)

	// EnsureResourceGroup ensures the specified resource group exists in the specified location
	EnsureResourceGroup(resourceGroup, location string) (*resources.Group, error)

	//
	// COMPUTE

	// List lists VM resources
	ListVirtualMachines(resourceGroup string) (compute.VirtualMachineListResult, error)

	// GetVirtualMachine retrieves the specified virtual machine.
	GetVirtualMachine(resourceGroup, name string) (compute.VirtualMachine, error)

	// DeleteVirtualMachine deletes the specified virtual machine.
	DeleteVirtualMachine(resourceGroup, name string, cancel <-chan struct{}) (<-chan compute.OperationStatusResponse, <-chan error)

	//
	// STORAGE

	// GetStorageClient uses SRP to retrieve keys, and then an authenticated client for talking to the specified storage
	// account.
	GetStorageClient(resourceGroup, accountName string) (ACSStorageClient, error)

	//
	// NETWORK

	// DeleteNetworkInterface deletes the specified network interface.
	DeleteNetworkInterface(resourceGroup, nicName string, cancel <-chan struct{}) (<-chan autorest.Response, <-chan error)
}

// ACSStorageClient interface models the azure storage client
type ACSStorageClient interface {
	// DeleteBlob deletes the specified blob in the specified container.
	DeleteBlob(container, blob string) error
}

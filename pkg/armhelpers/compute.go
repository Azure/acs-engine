package armhelpers

import (
	"github.com/Azure/azure-sdk-for-go/arm/compute"
)

// ListVirtualMachines returns (the first page of) the machines in the specified resource group.
func (az *AzureClient) ListVirtualMachines(resourceGroup string) (compute.VirtualMachineListResult, error) {
	return az.virtualMachinesClient.List(resourceGroup)
}

// GetVirtualMachine returns the specified machine in the specified resource group.
func (az *AzureClient) GetVirtualMachine(resourceGroup, name string) (compute.VirtualMachine, error) {
	return az.virtualMachinesClient.Get(resourceGroup, name, "")
}

// DeleteVirtualMachine handles deletion of a CRP/VMAS VM (aka, not a VMSS VM).
func (az *AzureClient) DeleteVirtualMachine(resourceGroup, name string, cancel <-chan struct{}) (<-chan compute.OperationStatusResponse, <-chan error) {
	return az.virtualMachinesClient.Delete(resourceGroup, name, cancel)
}

// ListVirtualMachineScaleSets returns (the first page of) the vmss resources in the specified resource group.
func (az *AzureClient) ListVirtualMachineScaleSets(resourceGroup string) (compute.VirtualMachineScaleSetListResult, error) {
	return az.virtualMachineScaleSetsClient.List(resourceGroup)
}

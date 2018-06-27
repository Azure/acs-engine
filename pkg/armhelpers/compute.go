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

// ListVirtualMachineScaleSetVMs returns the list of VMs per VMSS
func (az *AzureClient) ListVirtualMachineScaleSetVMs(resourceGroup, virtualMachineScaleSet string) (compute.VirtualMachineScaleSetVMListResult, error) {
	return az.virtualMachineScaleSetVMsClient.List(resourceGroup, virtualMachineScaleSet, "", "", "")
}

// DeleteVirtualMachineScaleSetVM deletes a VM in a VMSS
func (az *AzureClient) DeleteVirtualMachineScaleSetVM(resourceGroup, virtualMachineScaleSet, instanceID string, cancel <-chan struct{}) (<-chan compute.OperationStatusResponse, <-chan error) {
	return az.virtualMachineScaleSetVMsClient.Delete(resourceGroup, virtualMachineScaleSet, instanceID, cancel)
}

// SetVirtualMachineScaleSetCapacity sets the VMSS capacity
func (az *AzureClient) SetVirtualMachineScaleSetCapacity(resourceGroup, virtualMachineScaleSet string, sku compute.Sku, location string, cancel <-chan struct{}) (<-chan compute.VirtualMachineScaleSet, <-chan error) {
	return az.virtualMachineScaleSetsClient.CreateOrUpdate(
		resourceGroup,
		virtualMachineScaleSet,
		compute.VirtualMachineScaleSet{
			Location: &location,
			Sku:      &sku,
		},
		cancel,
	)
}

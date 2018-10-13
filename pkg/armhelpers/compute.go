package armhelpers

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2018-04-01/compute"
)

// ListVirtualMachines returns (the first page of) the machines in the specified resource group.
func (az *AzureClient) ListVirtualMachines(ctx context.Context, resourceGroup string) (VirtualMachineListResultPage, error) {
	page, err := az.virtualMachinesClient.List(ctx, resourceGroup)
	return &page, err
}

// GetVirtualMachine returns the specified machine in the specified resource group.
func (az *AzureClient) GetVirtualMachine(ctx context.Context, resourceGroup, name string) (compute.VirtualMachine, error) {
	return az.virtualMachinesClient.Get(ctx, resourceGroup, name, "")
}

// DeleteVirtualMachine handles deletion of a CRP/VMAS VM (aka, not a VMSS VM).
func (az *AzureClient) DeleteVirtualMachine(ctx context.Context, resourceGroup, name string) error {
	future, err := az.virtualMachinesClient.Delete(ctx, resourceGroup, name)
	if err != nil {
		return err
	}

	if err = future.WaitForCompletionRef(ctx, az.virtualMachinesClient.Client); err != nil {
		return err
	}

	_, err = future.Result(az.virtualMachinesClient)
	return err
}

// ListVirtualMachineScaleSets returns (the first page of) the vmss resources in the specified resource group.
func (az *AzureClient) ListVirtualMachineScaleSets(ctx context.Context, resourceGroup string) (compute.VirtualMachineScaleSetListResultPage, error) {
	return az.virtualMachineScaleSetsClient.List(ctx, resourceGroup)
}

// ListVirtualMachineScaleSetVMs returns the list of VMs per VMSS
func (az *AzureClient) ListVirtualMachineScaleSetVMs(ctx context.Context, resourceGroup, virtualMachineScaleSet string) (compute.VirtualMachineScaleSetVMListResultPage, error) {
	return az.virtualMachineScaleSetVMsClient.List(ctx, resourceGroup, virtualMachineScaleSet, "", "", "")
}

// DeleteVirtualMachineScaleSetVM deletes a VM in a VMSS
func (az *AzureClient) DeleteVirtualMachineScaleSetVM(ctx context.Context, resourceGroup, virtualMachineScaleSet, instanceID string) error {
	future, err := az.virtualMachineScaleSetVMsClient.Delete(ctx, resourceGroup, virtualMachineScaleSet, instanceID)
	if err != nil {
		return err
	}

	if err = future.WaitForCompletionRef(ctx, az.virtualMachineScaleSetVMsClient.Client); err != nil {
		return err
	}

	_, err = future.Result(az.virtualMachineScaleSetVMsClient)
	return err
}

// SetVirtualMachineScaleSetCapacity sets the VMSS capacity
func (az *AzureClient) SetVirtualMachineScaleSetCapacity(ctx context.Context, resourceGroup, virtualMachineScaleSet string, sku compute.Sku, location string) error {
	future, err := az.virtualMachineScaleSetsClient.CreateOrUpdate(
		ctx,
		resourceGroup,
		virtualMachineScaleSet,
		compute.VirtualMachineScaleSet{
			Location: &location,
			Sku:      &sku,
		})
	if err != nil {
		return err
	}

	if err = future.WaitForCompletionRef(ctx, az.virtualMachineScaleSetsClient.Client); err != nil {
		return err
	}

	_, err = future.Result(az.virtualMachineScaleSetsClient)
	return err
}

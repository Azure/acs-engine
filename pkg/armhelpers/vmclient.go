package armhelpers

import (
	"github.com/Azure/azure-sdk-for-go/arm/compute"
)

// ListVirtualMachines returns (the first page of) the machines in the specified resource group.
func (az *AzureClient) ListVirtualMachines(resourceGroup string) (compute.VirtualMachineListResult, error) {
	return az.virtualMachinesClient.List(resourceGroup)
}

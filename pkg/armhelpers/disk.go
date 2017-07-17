package armhelpers

import "github.com/Azure/azure-sdk-for-go/arm/disk"

// DeleteManagedDisk deletes a managed disk.
func (az *AzureClient) DeleteManagedDisk(resourceGroupName string, diskName string, cancel <-chan struct{}) (<-chan disk.OperationStatusResponse, <-chan error) {
	return az.disksClient.Delete(resourceGroupName, diskName, cancel)
}

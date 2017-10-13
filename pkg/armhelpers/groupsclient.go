package armhelpers

import (
	"github.com/Azure/azure-sdk-for-go/arm/resources/resources"
	"github.com/Azure/go-autorest/autorest"
)

// EnsureResourceGroup ensures the named resouce group exists in the given location.
func (az *AzureClient) EnsureResourceGroup(name, location string, managedBy *string) (resourceGroup *resources.Group, err error) {
	response, err := az.groupsClient.CreateOrUpdate(name, resources.Group{
		Name:      &name,
		Location:  &location,
		ManagedBy: managedBy,
	})
	if err != nil {
		return &response, err
	}

	return &response, nil
}

// CheckResourceGroupExistence return if the resource group exists
func (az *AzureClient) CheckResourceGroupExistence(name string) (result autorest.Response, err error) {
	return az.groupsClient.CheckExistence(name)
}

// DeleteResourceGroup delete the named resource group
func (az *AzureClient) DeleteResourceGroup(name string, cancel chan struct{}) (<-chan autorest.Response, <-chan error) {
	return az.groupsClient.Delete(name, cancel)
}

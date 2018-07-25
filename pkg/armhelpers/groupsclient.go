package armhelpers

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2018-02-01/resources"
	"github.com/Azure/go-autorest/autorest"
)

// EnsureResourceGroup ensures the named resouce group exists in the given location.
func (az *AzureClient) EnsureResourceGroup(ctx context.Context, name, location string, managedBy *string) (resourceGroup *resources.Group, err error) {
	var tags map[string]*string
	group, err := az.groupsClient.Get(ctx, name)
	if err == nil {
		tags = group.Tags
	}

	response, err := az.groupsClient.CreateOrUpdate(ctx, name, resources.Group{
		Name:      &name,
		Location:  &location,
		ManagedBy: managedBy,
		Tags:      tags,
	})
	if err != nil {
		return &response, err
	}

	return &response, nil
}

// CheckResourceGroupExistence return if the resource group exists
func (az *AzureClient) CheckResourceGroupExistence(ctx context.Context, name string) (result autorest.Response, err error) {
	return az.groupsClient.CheckExistence(ctx, name)
}

// DeleteResourceGroup delete the named resource group
func (az *AzureClient) DeleteResourceGroup(ctx context.Context, name string) (resources.GroupsDeleteFuture, error) {
	return az.groupsClient.Delete(ctx, name)
}

package armhelpers

import (
	"net/http"

	"github.com/Azure/azure-sdk-for-go/arm/resources/resources"
	"github.com/prometheus/common/log"
)

// EnsureResourceGroup ensures the named resouce group exists in the given location.
func (az *AzureClient) EnsureResourceGroup(name, location string) (resourceGroup *resources.Group, err error) {
	log.Debugf("Ensuring resource group exists. resourcegroup=%q", name)
	response, err := az.groupsClient.CreateOrUpdate(name, resources.Group{
		Name:     &name,
		Location: &location,
	})
	if err != nil {
		return &response, err
	}

	return &response, nil
}

// DeleteResourceGroup delete the named resource group
func (az *AzureClient) DeleteResourceGroup(name string, cancel chan struct{}) error {
	resCh, errCh := az.groupsClient.Delete(name, cancel)
	res := <-resCh
	err := <-errCh
	// When resourceGroup not exists, we will get 404
	// Explictly set the error to nil for this scenario
	if res.Response != nil && res.StatusCode == http.StatusNotFound {
		return nil
	}
	return err
}

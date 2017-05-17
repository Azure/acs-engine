package armhelpers

import (
	"github.com/Azure/go-autorest/autorest"
)

// DeleteNetworkInterface deletes the specified network interface.
func (az *AzureClient) DeleteNetworkInterface(resourceGroup, nicName string, cancel <-chan struct{}) (<-chan autorest.Response, <-chan error) {
	return az.interfacesClient.Delete(resourceGroup, nicName, cancel)
}

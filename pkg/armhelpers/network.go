package armhelpers

import (
	"github.com/Azure/go-autorest/autorest"
)

// DeleteNetworkInterface deletes the specified network interface.
func (az *AzureClient) DeleteNetworkInterface(resourceGroup, nicName string, cancel <-chan struct{}) (<-chan autorest.Response, <-chan error) {
	client := az.interfacesClient
	resultChan := make(chan autorest.Response, 1)
	errChan := make(chan error, 1)
	go func() {
		var err error
		var result autorest.Response
		defer func() {
			resultChan <- result
			errChan <- err
			close(resultChan)
			close(errChan)
		}()
		req, err := client.DeletePreparer(resourceGroup, nicName, cancel)
		if err != nil {
			err = autorest.NewErrorWithError(err, "network.InterfacesClient", "Delete", nil, "Failure preparing request")
			return
		}
		az.addAcceptLanguages(req)

		resp, err := client.DeleteSender(req)
		if err != nil {
			result.Response = resp
			err = autorest.NewErrorWithError(err, "network.InterfacesClient", "Delete", resp, "Failure sending request")
			return
		}

		result, err = client.DeleteResponder(resp)
		if err != nil {
			err = autorest.NewErrorWithError(err, "network.InterfacesClient", "Delete", resp, "Failure responding to request")
		}
	}()
	return resultChan, errChan
}

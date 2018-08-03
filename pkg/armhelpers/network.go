package armhelpers

import (
	"context"
)

// DeleteNetworkInterface deletes the specified network interface.
func (az *AzureClient) DeleteNetworkInterface(ctx context.Context, resourceGroup, nicName string) error {
	future, err := az.interfacesClient.Delete(ctx, resourceGroup, nicName)
	if err != nil {
		return err
	}

	if err = future.WaitForCompletion(ctx, az.interfacesClient.Client); err != nil {
		return err
	}

	_, err = future.Result(az.interfacesClient)
	return err
}

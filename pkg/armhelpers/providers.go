package armhelpers

import (
	"github.com/Azure/azure-sdk-for-go/arm/resources/resources"
	"github.com/Azure/go-autorest/autorest/to"
)

// ListProviders returns all the providers for a given AzureClient
func (az *AzureClient) ListProviders() (resources.ProviderListResult, error) {
	return az.providersClient.List(to.Int32Ptr(100), "")
}

package armhelpers

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2018-02-01/storage"
	azStorage "github.com/Azure/azure-sdk-for-go/storage"
	"github.com/Azure/go-autorest/autorest/to"
)

// AzureStorageClient implements the StorageClient interface and wraps the Azure storage client.
type AzureStorageClient struct {
	client *azStorage.Client
}

// GetStorageClient returns an authenticated client for the specified account.
func (az *AzureClient) GetStorageClient(ctx context.Context, resourceGroup, accountName string) (ACSStorageClient, error) {
	keys, err := az.getStorageKeys(ctx, resourceGroup, accountName)
	if err != nil {
		return nil, err
	}

	client, err := azStorage.NewBasicClientOnSovereignCloud(accountName, to.String(keys[0].Value), az.environment)
	if err != nil {
		return nil, err
	}

	return &AzureStorageClient{
		client: &client,
	}, nil
}

func (az *AzureClient) getStorageKeys(ctx context.Context, resourceGroup, accountName string) ([]storage.AccountKey, error) {
	storageKeysResult, err := az.storageAccountsClient.ListKeys(ctx, resourceGroup, accountName)
	if err != nil {
		return nil, err
	}

	return *storageKeysResult.Keys, nil
}

// DeleteBlob deletes the specified blob
// TODO(colemick): why doesn't SDK give a way to just delete a blob by URI?
// it's what it ends up doing internally anyway...
func (as *AzureStorageClient) DeleteBlob(vhdContainer, vhdBlob string) error {
	bs := as.client.GetBlobService()
	containerRef := bs.GetContainerReference(vhdContainer)
	blobRef := containerRef.GetBlobReference(vhdBlob)

	return blobRef.Delete(&azStorage.DeleteBlobOptions{})
}

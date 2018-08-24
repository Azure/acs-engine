package armhelpers

import (
	"context"

	"net/http"

	"github.com/Azure/azure-sdk-for-go/services/preview/msi/mgmt/2015-08-31-preview/msi"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/to"
	log "github.com/sirupsen/logrus"
)

func (az *AzureClient) IsUserAssignedIDPresent(resourceGroup string, userAssignedID string) (bool, error) {
	log.Infof("Checking if the %s is present in rg %s", userAssignedID, resourceGroup)
	_, err := az.msiClient.Get(context.Background(), resourceGroup, userAssignedID)
	if err != nil {
		detailedError, castOk := err.(autorest.DetailedError)
		if castOk && detailedError.StatusCode == http.StatusNotFound {
			log.Errorf("Not found")
			return false, nil
		}
		log.Errorf("Error: %+v", err)
		// TODO(mandatory): There could be other errors like network errors, which we need to distinguish
		// and bail out, instead of going ahead and trying to go to creating the new
		// identity.
		return false, err
	}
	log.Infof("Found the userAssignedID")
	return true, nil
}

func (az *AzureClient) CreateUserAssignedID(location string, resourceGroup string, userAssignedID string) (id *msi.Identity, err error) {
	idCreated, err := az.msiClient.CreateOrUpdate(context.Background(), resourceGroup, userAssignedID, msi.Identity{
		Location: to.StringPtr(location),
	})
	if err != nil {
		log.Error(err)
		return nil, err
	}
	log.Infof("Created %s in rg %s", userAssignedID, resourceGroup)
	return &idCreated, nil
}

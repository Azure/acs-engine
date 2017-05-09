//------------------------------------------------------------
// Copyright (c) Microsoft Corporation.  All rights reserved.
//------------------------------------------------------------

package armhelpers

import (
	"net/http"

	"github.com/Azure/azure-sdk-for-go/arm/resources/resources"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/satori/uuid"
)

type armResourceGroupClient struct {
	client *resources.GroupsClient
}

// NewResourceGroupClient returns an instance of ResourceGroupClient that talks to the ARM API
// subscriptionID and service principal are required, if baseURI is empty a default is used.
func NewResourceGroupClient(subscriptionID uuid.UUID, token *adal.ServicePrincipalToken, baseURI string) ResourceGroupClient {
	var client resources.GroupsClient
	if len(baseURI) == 0 {
		client = resources.NewGroupsClient(subscriptionID.String())
	} else {
		client = resources.NewGroupsClientWithBaseURI(baseURI, subscriptionID.String())
	}
	client.Authorizer = autorest.NewBearerAuthorizer(token)

	return &armResourceGroupClient{
		client: &client,
	}
}

func (c *armResourceGroupClient) CreateOrUpdate(resourceGroup, location string) error {
	_, err := c.client.CreateOrUpdate(resourceGroup, resources.Group{
		Location: &location,
	})
	return err
}

func (c *armResourceGroupClient) Delete(resourceGroup string, cancel chan struct{}) error {
	resChan, errChan := c.client.Delete(resourceGroup, cancel)
	res := <-resChan

	// When resourceGroup not exists, we will get 404
	// Explictly set the error to nil for this scenario
	if res.Response != nil && res.StatusCode == http.StatusNotFound {
		return nil
	}
	return <-errChan
}

func (c *armResourceGroupClient) Get(resourceGroup string) (*resources.Group, error) {
	group, err := c.client.Get(resourceGroup)
	// Return nil when resourceGroup not found
	if group.Response.Response != nil && group.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	return &group, err
}

func (c *armResourceGroupClient) ListResources(resourceGroup string, filter string, expand string, top *int32) (*resources.ListResult, error) {
	listres, err := c.client.ListResources(resourceGroup, filter, expand, top)

	return &listres, err
}

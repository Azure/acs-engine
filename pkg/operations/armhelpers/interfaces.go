//------------------------------------------------------------
// Copyright (c) Microsoft Corporation.  All rights reserved.
//------------------------------------------------------------

package armhelpers

import "github.com/Azure/azure-sdk-for-go/arm/resources/resources"

// ResourceGroupClient is an interface for managing resource groups
type ResourceGroupClient interface {
	// CreateOrUpdate a resource group, returns an error if one occurs
	CreateOrUpdate(resourceGroup string, location string) error
	// Delete a resource group, returns an error if one occurs
	Delete(resourceGroup string, cancel chan struct{}) error
	// Get a resource group, returns nil, if it doesn't exist, returns an error if one occurs.
	Get(resourceGroup string) (*resources.Group, error)
	// // ListResources in a resource group
	ListResources(resourceGroup string, filter string, expand string, top *int32) (*resources.ListResult, error)
}

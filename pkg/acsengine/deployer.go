//------------------------------------------------------------
// Copyright (c) Microsoft Corporation.  All rights reserved.
//------------------------------------------------------------

// This package contains ARM specific implementations of the TemplateDeployer interface.
// Expected use:
//    deployer := NewARMDeployer(...)
//    deployer.CreateOrUpderDeployment(...)

package acsengine

import (
	"net/http"

	"github.com/Azure/azure-sdk-for-go/arm/resources/resources"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/Sirupsen/logrus"
	"github.com/satori/go.uuid"
)

// TemplateDeployer is an interface that knows how to deploy templates
type TemplateDeployer interface {
	// Validate validates if a deployment is valid (e.g. it can be run successfully)
	Validate(resourceGroup, name string, d *resources.Deployment) (*resources.DeploymentValidateResult, error)
	// CreateOrUpdate a deployment of a template, returns an error if one occurs
	CreateOrUpdate(resourceGroup, name string, d *resources.Deployment, cancel chan struct{}) error
	// Delete a deployment of a template, returns an error if one occurs
	Delete(resourceGroup, name string, cancel chan struct{}) error
	// Get a deployment, returns nil if it doesn't exist, returns an error if one occurs
	Get(resourceGroup, name string) (*resources.DeploymentExtended, error)
}

type armDeployer struct {
	client resources.DeploymentsClient
}

// TODO: We should construct an entire "AzureClient" (like azkube) that handles
// construct ahead-of-time of all of the various clients. This will help with mocking
// and will help with not having to have each chunk of code that needs a new client to
// handle the token/Authorizer/etc. (and can encapsulate the environment[aka baseURI] and subId handling)

// NewARMDeployer creates a concrete instance of TemplateDeployer that uses the ARM API
// subscriptionID and service principal are required, if baseURI is empty a default is used.
// TODO: Add regional language so error messages are produced in the correct language
func NewARMDeployer(subscriptionID uuid.UUID, token *adal.ServicePrincipalToken, baseURI string) TemplateDeployer {
	var c resources.DeploymentsClient
	if len(baseURI) == 0 {
		c = resources.NewDeploymentsClient(subscriptionID.String())
	} else {
		c = resources.NewDeploymentsClientWithBaseURI(baseURI, subscriptionID.String())
	}
	c.Authorizer = autorest.NewBearerAuthorizer(token)

	return &armDeployer{
		client: c,
	}
}

func (a *armDeployer) Validate(resourceGroup, name string, d *resources.Deployment) (*resources.DeploymentValidateResult, error) {
	res, err := a.client.Validate(resourceGroup, name, *d)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// CreateOrUpdate implements the TemplateDeployer interface
func (a *armDeployer) CreateOrUpdate(resourceGroup, name string, d *resources.Deployment, cancel chan struct{}) error {
	_, err := a.client.CreateOrUpdate(resourceGroup, name, *d, cancel)
	if e := <-err; e != nil {
		logrus.Errorf("Error creating deployment: %s", e.Error())
		return e
	}
	return nil
}

// Delete implements the TemplateDeployer interface
func (a *armDeployer) Delete(resourceGroup, name string, cancel chan struct{}) error {
	_, err := a.client.Delete(resourceGroup, name, cancel)
	if e := <-err; e != nil {
		logrus.Errorf("Error deleting deployment: %s", e.Error())
		return e
	}
	return nil
}

// Get implements the TemplateDeployer interface
func (a *armDeployer) Get(resourceGroup, name string) (*resources.DeploymentExtended, error) {
	res, err := a.client.Get(resourceGroup, name)
	if res.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	return &res, err
}

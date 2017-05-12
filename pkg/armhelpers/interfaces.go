package armhelpers

import (
	"github.com/Azure/azure-sdk-for-go/arm/resources/resources"
	"github.com/Azure/go-autorest/autorest"
)

// the master client interface
// if this is mocked out, then the entire context of the target environment is mocked:
// once one of the NewAzureClient... returns, it has returned a fully initialized UberClient
// all other code that needs to do work in subscription/AAD tenant will use this
// currently only one impl: AzureClient
type UberClient interface {
	TemplateDeployer() TemplateDeployer // wraps the deployment client
}

// TemplateDeployer is an interface that knows how to deploy templates
type TemplateDeployer interface {
	// Validate validates if a deployment is valid (e.g. it can be run successfully)
	Validate(resourceGroup, name string, d resources.Deployment) (resources.DeploymentValidateResult, error)
	// CreateOrUpdate a deployment of a template, returns an error if one occurs
	CreateOrUpdate(resourceGroup, name string, d resources.Deployment, cancel <-chan struct{}) (<-chan resources.DeploymentExtended, <-chan error)
	// Delete a deployment of a template, returns an error if one occurs
	Delete(resourceGroup, name string, cancel <-chan struct{}) (<-chan autorest.Response, <-chan error)
	// Get a deployment, returns nil if it doesn't exist, returns an error if one occurs
	Get(resourceGroup, name string) (resources.DeploymentExtended, error)
}

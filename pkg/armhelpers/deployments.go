package armhelpers

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2018-05-01/resources"
	"github.com/Azure/go-autorest/autorest"
	log "github.com/sirupsen/logrus"
)

// DeployTemplate implements the TemplateDeployer interface for the AzureClient client
func (az *AzureClient) DeployTemplate(ctx context.Context, resourceGroupName, deploymentName string, template map[string]interface{}, parameters map[string]interface{}) (resources.DeploymentExtended, error) {
	deployment := resources.Deployment{
		Properties: &resources.DeploymentProperties{
			Template:   &template,
			Parameters: &parameters,
			Mode:       resources.Incremental,
		},
	}

	log.Infof("Starting ARM Deployment (%s). This will take some time...", deploymentName)
	future, err := az.deploymentsClient.CreateOrUpdate(ctx, resourceGroupName, deploymentName, deployment)

	// wait till future completes
	e1 := future.WaitForCompletion(ctx, az.deploymentsClient.Client)
	if err != nil && e1 != nil {
		return resources.DeploymentExtended{}, err
	}

	// get future results
	deploymentResult, e2 := future.Result(az.deploymentsClient)
	if err != nil && e2 != nil {
		return resources.DeploymentExtended{}, err
	}

	outcomeText := "Succeeded"
	if err != nil {
		outcomeText = fmt.Sprintf("Error: %v", err)
	}

	log.Infof("Finished ARM Deployment (%s). %s", deploymentName, outcomeText)
	return deploymentResult, err
}

// ValidateTemplate validate the template and parameters
func (az *AzureClient) ValidateTemplate(
	ctx context.Context,
	resourceGroupName string,
	deploymentName string,
	template map[string]interface{},
	parameters map[string]interface{}) (result resources.DeploymentValidateResult, err error) {
	deployment := resources.Deployment{
		Properties: &resources.DeploymentProperties{
			Template:   &template,
			Parameters: &parameters,
			Mode:       resources.Incremental,
		},
	}
	return az.deploymentsClient.Validate(ctx, resourceGroupName, deploymentName, deployment)
}

// GetDeployment returns the template deployment
func (az *AzureClient) GetDeployment(ctx context.Context, resourceGroupName, deploymentName string) (result resources.DeploymentExtended, err error) {
	return az.deploymentsClient.Get(ctx, resourceGroupName, deploymentName)
}

// CheckDeploymentExistence returns if the deployment already exists
func (az *AzureClient) CheckDeploymentExistence(ctx context.Context, resourceGroupName string, deploymentName string) (result autorest.Response, err error) {
	return az.deploymentsClient.CheckExistence(ctx, resourceGroupName, deploymentName)
}

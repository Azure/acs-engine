package armhelpers

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/arm/resources/resources"
	"github.com/Azure/go-autorest/autorest"
	log "github.com/sirupsen/logrus"
)

// DeployTemplate implements the TemplateDeployer interface for the AzureClient client
func (az *AzureClient) DeployTemplate(resourceGroupName, deploymentName string, template map[string]interface{}, parameters map[string]interface{}, cancel <-chan struct{}) (*resources.DeploymentExtended, error) {
	deployment := resources.Deployment{
		Properties: &resources.DeploymentProperties{
			Template:   &template,
			Parameters: &parameters,
			Mode:       resources.Incremental,
		},
	}

	log.Infof("Starting ARM Deployment (%s). This will take some time...", deploymentName)

	resChan, errChan := az.deploymentsClient.CreateOrUpdate(
		resourceGroupName,
		deploymentName,
		deployment,
		cancel)

	err := <-errChan
	res, ok := <-resChan
	if !ok {
		// This path is taken when validation is failed before calling ARM
		return nil, err
	}

	outcomeText := "Succeeded"
	if err != nil {
		outcomeText = fmt.Sprintf("Error: %v", err)
	}
	log.Infof("Finished ARM Deployment (%s). %s", deploymentName, outcomeText)

	return &res, err
}

// ValidateTemplate validate the template and parameters
func (az *AzureClient) ValidateTemplate(
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
	return az.deploymentsClient.Validate(resourceGroupName, deploymentName, deployment)
}

// GetDeployment returns the template deployment
func (az *AzureClient) GetDeployment(resourceGroupName, deploymentName string) (result resources.DeploymentExtended, err error) {
	return az.deploymentsClient.Get(resourceGroupName, deploymentName)
}

// CheckDeploymentExistence returns if the deployment already exists
func (az *AzureClient) CheckDeploymentExistence(resourceGroupName string, deploymentName string) (result autorest.Response, err error) {
	return az.deploymentsClient.CheckExistence(resourceGroupName, deploymentName)
}

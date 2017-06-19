package armhelpers

import (
	"github.com/Azure/azure-sdk-for-go/arm/resources/resources"
	"github.com/Azure/go-autorest/autorest"
	"github.com/prometheus/common/log"
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

	log.Infof("Starting ARM Deployment. This will take some time. deployment=%q", deploymentName)

	resChan, errChan := az.deploymentsClient.CreateOrUpdate(
		resourceGroupName,
		deploymentName,
		deployment,
		cancel)
	if err := <-errChan; err != nil {
		return nil, err
	}
	res := <-resChan

	log.Infof("Finished ARM Deployment. deployment=%q. res=%q", deploymentName, res)

	return &res, nil
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

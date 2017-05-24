package armhelpers

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/arm/resources/resources"
	"github.com/prometheus/common/log"
)

// DeployTemplate implements the TemplateDeployer interface for the AzureClient client
func (az *AzureClient) DeployTemplate(resourceGroupName, deploymentName string, template map[string]interface{}, parameters map[string]interface{}, cancel <-chan struct{}) (*resources.DeploymentExtended, error) {
	// this is needed because either ARM or the SDK can't distinguish between past
	// deployments and current deployments with the same deploymentName.
	uniqueSuffix := fmt.Sprintf("-%d", time.Now().Unix())
	deploymentName = deploymentName + uniqueSuffix

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

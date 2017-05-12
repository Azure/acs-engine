package armhelpers

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/arm/resources/resources"
	"github.com/Azure/go-autorest/autorest"
	"github.com/prometheus/common/log"
)

type AzureTemplateDeployer struct{}

func (az *AzureClient) Validate(resourceGroup, name string, d resources.Deployment) (resources.DeploymentValidateResult, error) {
	return az.DeploymentsClient.Validate(resourceGroup, name, d)
}

func (az *AzureClient) CreateOrUpdate(resourceGroup, name string, d resources.Deployment, cancel chan struct{}) (<-chan resources.DeploymentExtended, <-chan error) {
	return az.DeploymentsClient.CreateOrUpdate(resourceGroup, name, d, cancel)
}

func (az *AzureClient) Delete(resourceGroup, name string, cancel chan struct{}) (<-chan autorest.Response, <-chan error) {
	return az.DeploymentsClient.Delete(resourceGroup, name, cancel)
}

func (az *AzureClient) Get(resourceGroup, name string) (resources.DeploymentExtended, error) {
	return az.DeploymentsClient.Get(resourceGroup, name)
}

func DeployTemplate(deployer TemplateDeployer, resourceGroupName, deploymentName string, template map[string]interface{}, parameters map[string]interface{}) (response *resources.DeploymentExtended, err error) {
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

	resChan, errChan := deployer.CreateOrUpdate(
		resourceGroupName,
		deploymentName,
		deployment,
		nil)
	if err := <-errChan; err != nil {
		return nil, err
	}
	res := <-resChan

	log.Infof("Finished ARM Deployment. deployment=%q. res=%q", deploymentName, res)

	return nil, nil
}

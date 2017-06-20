package armhelpers

import (
	"github.com/Azure/azure-sdk-for-go/arm/resources/resources"
)

// ListDeploymentOperations gets all deployments operations for a deployment.
func (az *AzureClient) ListDeploymentOperations(resourceGroupName string, deploymentName string, top *int32) (result resources.DeploymentOperationsListResult, err error) {
	return az.deploymentOperationsClient.List(resourceGroupName, deploymentName, top)
}

// ListDeploymentOperationsNextResults retrieves the next set of results, if any.
func (az *AzureClient) ListDeploymentOperationsNextResults(lastResults resources.DeploymentOperationsListResult) (result resources.DeploymentOperationsListResult, err error) {
	return az.deploymentOperationsClient.ListNextResults(lastResults)
}

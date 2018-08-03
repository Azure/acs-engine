package armhelpers

import (
	"context"
)

// ListDeploymentOperations gets all deployments operations for a deployment.
func (az *AzureClient) ListDeploymentOperations(ctx context.Context, resourceGroupName string, deploymentName string, top *int32) (DeploymentOperationsListResultPage, error) {
	list, err := az.deploymentOperationsClient.List(ctx, resourceGroupName, deploymentName, top)
	return &list, err
}

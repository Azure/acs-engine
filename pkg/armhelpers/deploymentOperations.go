package armhelpers

import (
	"github.com/Azure/azure-sdk-for-go/arm/resources/resources"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/validation"
)

// ListDeploymentOperations gets all deployments operations for a deployment.
func (az *AzureClient) ListDeploymentOperations(resourceGroupName string, deploymentName string, top *int32) (result resources.DeploymentOperationsListResult, err error) {
	client := az.deploymentOperationsClient
	if err := validation.Validate([]validation.Validation{
		{TargetValue: resourceGroupName,
			Constraints: []validation.Constraint{{Target: "resourceGroupName", Name: validation.MaxLength, Rule: 90, Chain: nil},
				{Target: "resourceGroupName", Name: validation.MinLength, Rule: 1, Chain: nil},
				{Target: "resourceGroupName", Name: validation.Pattern, Rule: `^[-\w\._\(\)]+$`, Chain: nil}}},
		{TargetValue: deploymentName,
			Constraints: []validation.Constraint{{Target: "deploymentName", Name: validation.MaxLength, Rule: 64, Chain: nil},
				{Target: "deploymentName", Name: validation.MinLength, Rule: 1, Chain: nil},
				{Target: "deploymentName", Name: validation.Pattern, Rule: `^[-\w\._\(\)]+$`, Chain: nil}}}}); err != nil {
		return result, validation.NewErrorWithValidationError(err, "resources.DeploymentOperationsClient", "List")
	}

	req, err := client.ListPreparer(resourceGroupName, deploymentName, top)
	if err != nil {
		err = autorest.NewErrorWithError(err, "resources.DeploymentOperationsClient", "List", nil, "Failure preparing request")
		return
	}
	az.addAcceptLanguages(req)

	resp, err := client.ListSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "resources.DeploymentOperationsClient", "List", resp, "Failure sending request")
		return
	}

	result, err = client.ListResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "resources.DeploymentOperationsClient", "List", resp, "Failure responding to request")
	}

	return
}

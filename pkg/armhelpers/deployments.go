package armhelpers

import (
	"github.com/Azure/azure-sdk-for-go/arm/resources/resources"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/validation"
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
	client := az.deploymentsClient

	log.Infof("Starting ARM Deployment. This will take some time. deployment=%q", deploymentName)

	resChan := make(chan resources.DeploymentExtended, 1)
	errChan := make(chan error, 1)
	if err := validation.Validate([]validation.Validation{
		{TargetValue: resourceGroupName,
			Constraints: []validation.Constraint{{Target: "resourceGroupName", Name: validation.MaxLength, Rule: 90, Chain: nil},
				{Target: "resourceGroupName", Name: validation.MinLength, Rule: 1, Chain: nil},
				{Target: "resourceGroupName", Name: validation.Pattern, Rule: `^[-\w\._\(\)]+$`, Chain: nil}}},
		{TargetValue: deploymentName,
			Constraints: []validation.Constraint{{Target: "deploymentName", Name: validation.MaxLength, Rule: 64, Chain: nil},
				{Target: "deploymentName", Name: validation.MinLength, Rule: 1, Chain: nil},
				{Target: "deploymentName", Name: validation.Pattern, Rule: `^[-\w\._\(\)]+$`, Chain: nil}}},
		{TargetValue: parameters,
			Constraints: []validation.Constraint{{Target: "parameters.Properties", Name: validation.Null, Rule: true,
				Chain: []validation.Constraint{{Target: "parameters.Properties.TemplateLink", Name: validation.Null, Rule: false,
					Chain: []validation.Constraint{{Target: "parameters.Properties.TemplateLink.URI", Name: validation.Null, Rule: true, Chain: nil}}},
					{Target: "parameters.Properties.ParametersLink", Name: validation.Null, Rule: false,
						Chain: []validation.Constraint{{Target: "parameters.Properties.ParametersLink.URI", Name: validation.Null, Rule: true, Chain: nil}}},
				}}}}}); err != nil {
		return nil, validation.NewErrorWithValidationError(err, "resources.DeploymentsClient", "CreateOrUpdate")
	}

	go func() {
		var err error
		var result resources.DeploymentExtended
		defer func() {
			resChan <- result
			errChan <- err
			close(resChan)
			close(errChan)
		}()
		req, err := client.CreateOrUpdatePreparer(resourceGroupName, deploymentName, deployment, cancel)
		if err != nil {
			err = autorest.NewErrorWithError(err, "resources.DeploymentsClient", "CreateOrUpdate", nil, "Failure preparing request")
			return
		}
		az.addAcceptLanguages(req)

		resp, err := client.CreateOrUpdateSender(req)
		if err != nil {
			result.Response = autorest.Response{Response: resp}
			err = autorest.NewErrorWithError(err, "resources.DeploymentsClient", "CreateOrUpdate", resp, "Failure sending request")
			return
		}

		result, err = client.CreateOrUpdateResponder(resp)
		if err != nil {
			err = autorest.NewErrorWithError(err, "resources.DeploymentsClient", "CreateOrUpdate", resp, "Failure responding to request")
		}
	}()
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
	client := az.deploymentsClient

	if err := validation.Validate([]validation.Validation{
		{TargetValue: resourceGroupName,
			Constraints: []validation.Constraint{{Target: "resourceGroupName", Name: validation.MaxLength, Rule: 90, Chain: nil},
				{Target: "resourceGroupName", Name: validation.MinLength, Rule: 1, Chain: nil},
				{Target: "resourceGroupName", Name: validation.Pattern, Rule: `^[-\w\._\(\)]+$`, Chain: nil}}},
		{TargetValue: deploymentName,
			Constraints: []validation.Constraint{{Target: "deploymentName", Name: validation.MaxLength, Rule: 64, Chain: nil},
				{Target: "deploymentName", Name: validation.MinLength, Rule: 1, Chain: nil},
				{Target: "deploymentName", Name: validation.Pattern, Rule: `^[-\w\._\(\)]+$`, Chain: nil}}},
		{TargetValue: parameters,
			Constraints: []validation.Constraint{{Target: "parameters.Properties", Name: validation.Null, Rule: true,
				Chain: []validation.Constraint{{Target: "parameters.Properties.TemplateLink", Name: validation.Null, Rule: false,
					Chain: []validation.Constraint{{Target: "parameters.Properties.TemplateLink.URI", Name: validation.Null, Rule: true, Chain: nil}}},
					{Target: "parameters.Properties.ParametersLink", Name: validation.Null, Rule: false,
						Chain: []validation.Constraint{{Target: "parameters.Properties.ParametersLink.URI", Name: validation.Null, Rule: true, Chain: nil}}},
				}}}}}); err != nil {
		return result, validation.NewErrorWithValidationError(err, "resources.DeploymentsClient", "Validate")
	}

	req, err := client.ValidatePreparer(resourceGroupName, deploymentName, deployment)
	if err != nil {
		err = autorest.NewErrorWithError(err, "resources.DeploymentsClient", "Validate", nil, "Failure preparing request")
		return
	}
	az.addAcceptLanguages(req)

	resp, err := client.ValidateSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "resources.DeploymentsClient", "Validate", resp, "Failure sending request")
		return
	}

	result, err = client.ValidateResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "resources.DeploymentsClient", "Validate", resp, "Failure responding to request")
	}

	return
}

// GetDeployment returns the template deployment
func (az *AzureClient) GetDeployment(resourceGroupName, deploymentName string) (result resources.DeploymentExtended, err error) {
	client := az.deploymentsClient
	if err := validation.Validate([]validation.Validation{
		{TargetValue: resourceGroupName,
			Constraints: []validation.Constraint{{Target: "resourceGroupName", Name: validation.MaxLength, Rule: 90, Chain: nil},
				{Target: "resourceGroupName", Name: validation.MinLength, Rule: 1, Chain: nil},
				{Target: "resourceGroupName", Name: validation.Pattern, Rule: `^[-\w\._\(\)]+$`, Chain: nil}}},
		{TargetValue: deploymentName,
			Constraints: []validation.Constraint{{Target: "deploymentName", Name: validation.MaxLength, Rule: 64, Chain: nil},
				{Target: "deploymentName", Name: validation.MinLength, Rule: 1, Chain: nil},
				{Target: "deploymentName", Name: validation.Pattern, Rule: `^[-\w\._\(\)]+$`, Chain: nil}}}}); err != nil {
		return result, validation.NewErrorWithValidationError(err, "resources.DeploymentsClient", "Get")
	}

	req, err := client.GetPreparer(resourceGroupName, deploymentName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "resources.DeploymentsClient", "Get", nil, "Failure preparing request")
		return
	}
	az.addAcceptLanguages(req)

	resp, err := client.GetSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "resources.DeploymentsClient", "Get", resp, "Failure sending request")
		return
	}

	result, err = client.GetResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "resources.DeploymentsClient", "Get", resp, "Failure responding to request")
	}

	return
}

// CheckDeploymentExistence returns if the deployment already exists
func (az *AzureClient) CheckDeploymentExistence(resourceGroupName string, deploymentName string) (result autorest.Response, err error) {
	client := az.deploymentsClient
	if err := validation.Validate([]validation.Validation{
		{TargetValue: resourceGroupName,
			Constraints: []validation.Constraint{{Target: "resourceGroupName", Name: validation.MaxLength, Rule: 90, Chain: nil},
				{Target: "resourceGroupName", Name: validation.MinLength, Rule: 1, Chain: nil},
				{Target: "resourceGroupName", Name: validation.Pattern, Rule: `^[-\w\._\(\)]+$`, Chain: nil}}},
		{TargetValue: deploymentName,
			Constraints: []validation.Constraint{{Target: "deploymentName", Name: validation.MaxLength, Rule: 64, Chain: nil},
				{Target: "deploymentName", Name: validation.MinLength, Rule: 1, Chain: nil},
				{Target: "deploymentName", Name: validation.Pattern, Rule: `^[-\w\._\(\)]+$`, Chain: nil}}}}); err != nil {
		return result, validation.NewErrorWithValidationError(err, "resources.DeploymentsClient", "CheckExistence")
	}

	req, err := client.CheckExistencePreparer(resourceGroupName, deploymentName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "resources.DeploymentsClient", "CheckExistence", nil, "Failure preparing request")
		return
	}
	az.addAcceptLanguages(req)

	resp, err := client.CheckExistenceSender(req)
	if err != nil {
		result.Response = resp
		err = autorest.NewErrorWithError(err, "resources.DeploymentsClient", "CheckExistence", resp, "Failure sending request")
		return
	}

	result, err = client.CheckExistenceResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "resources.DeploymentsClient", "CheckExistence", resp, "Failure responding to request")
	}

	return
}

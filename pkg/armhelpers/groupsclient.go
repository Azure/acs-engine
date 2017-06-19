package armhelpers

import (
	"github.com/Azure/azure-sdk-for-go/arm/resources/resources"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/validation"
	"github.com/prometheus/common/log"
)

// EnsureResourceGroup ensures the named resouce group exists in the given location.
func (az *AzureClient) EnsureResourceGroup(name, location string) (resourceGroup *resources.Group, err error) {
	client := az.groupsClient
	log.Debugf("Ensuring resource group exists. resourcegroup=%q", name)
	parameters := resources.Group{
		Name:     &name,
		Location: &location,
	}
	if err := validation.Validate([]validation.Validation{
		{TargetValue: name,
			Constraints: []validation.Constraint{{Target: "resourceGroupName", Name: validation.MaxLength, Rule: 90, Chain: nil},
				{Target: "resourceGroupName", Name: validation.MinLength, Rule: 1, Chain: nil},
				{Target: "resourceGroupName", Name: validation.Pattern, Rule: `^[-\w\._\(\)]+$`, Chain: nil}}},
		{TargetValue: parameters,
			Constraints: []validation.Constraint{{Target: "parameters.Location", Name: validation.Null, Rule: true, Chain: nil}}}}); err != nil {
		return resourceGroup, validation.NewErrorWithValidationError(err, "resources.GroupsClient", "CreateOrUpdate")
	}

	req, err := client.CreateOrUpdatePreparer(name, parameters)
	if err != nil {
		err = autorest.NewErrorWithError(err, "resources.GroupsClient", "CreateOrUpdate", nil, "Failure preparing request")
		return
	}
	az.addAcceptLanguages(req)

	resp, err := client.CreateOrUpdateSender(req)
	if err != nil {
		resourceGroup = &resources.Group{
			Response: autorest.Response{Response: resp},
		}
		err = autorest.NewErrorWithError(err, "resources.GroupsClient", "CreateOrUpdate", resp, "Failure sending request")
		return
	}

	result, err := client.CreateOrUpdateResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "resources.GroupsClient", "CreateOrUpdate", resp, "Failure responding to request")
	}
	resourceGroup = &result
	return resourceGroup, err
}

// DeleteResourceGroup delete the named resource group
func (az *AzureClient) DeleteResourceGroup(name string, cancel chan struct{}) (<-chan autorest.Response, <-chan error) {
	client := az.groupsClient
	resultChan := make(chan autorest.Response, 1)
	errChan := make(chan error, 1)
	if err := validation.Validate([]validation.Validation{
		{TargetValue: name,
			Constraints: []validation.Constraint{{Target: "resourceGroupName", Name: validation.MaxLength, Rule: 90, Chain: nil},
				{Target: "resourceGroupName", Name: validation.MinLength, Rule: 1, Chain: nil},
				{Target: "resourceGroupName", Name: validation.Pattern, Rule: `^[-\w\._\(\)]+$`, Chain: nil}}}}); err != nil {
		errChan <- validation.NewErrorWithValidationError(err, "resources.GroupsClient", "Delete")
		close(errChan)
		close(resultChan)
		return resultChan, errChan
	}

	go func() {
		var err error
		var result autorest.Response
		defer func() {
			resultChan <- result
			errChan <- err
			close(resultChan)
			close(errChan)
		}()
		req, err := client.DeletePreparer(name, cancel)
		if err != nil {
			err = autorest.NewErrorWithError(err, "resources.GroupsClient", "Delete", nil, "Failure preparing request")
			return
		}
		az.addAcceptLanguages(req)

		resp, err := client.DeleteSender(req)
		if err != nil {
			result.Response = resp
			err = autorest.NewErrorWithError(err, "resources.GroupsClient", "Delete", resp, "Failure sending request")
			return
		}

		result, err = client.DeleteResponder(resp)
		if err != nil {
			err = autorest.NewErrorWithError(err, "resources.GroupsClient", "Delete", resp, "Failure responding to request")
		}
	}()
	return resultChan, errChan
}

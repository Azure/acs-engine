package apimanagement

// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Code generated by Microsoft (R) AutoRest Code Generator 1.0.1.0
// Changes may cause incorrect behavior and will be lost if the code is
// regenerated.

import (
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/validation"
	"net/http"
)

// QuotaByPeriodKeysClient is the composite Swagger for ApiManagement Client
type QuotaByPeriodKeysClient struct {
	ManagementClient
}

// NewQuotaByPeriodKeysClient creates an instance of the
// QuotaByPeriodKeysClient client.
func NewQuotaByPeriodKeysClient(subscriptionID string) QuotaByPeriodKeysClient {
	return NewQuotaByPeriodKeysClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewQuotaByPeriodKeysClientWithBaseURI creates an instance of the
// QuotaByPeriodKeysClient client.
func NewQuotaByPeriodKeysClientWithBaseURI(baseURI string, subscriptionID string) QuotaByPeriodKeysClient {
	return QuotaByPeriodKeysClient{NewWithBaseURI(baseURI, subscriptionID)}
}

// Get gets the value of the quota counter associated with the counter-key in
// the policy for the specific period in service instance.
//
// resourceGroupName is the name of the resource group. serviceName is the name
// of the API Management service. quotaCounterKey is quota counter key
// identifier. quotaPeriodKey is quota period key identifier.
func (client QuotaByPeriodKeysClient) Get(resourceGroupName string, serviceName string, quotaCounterKey string, quotaPeriodKey string) (result QuotaCounterContract, err error) {
	if err := validation.Validate([]validation.Validation{
		{TargetValue: serviceName,
			Constraints: []validation.Constraint{{Target: "serviceName", Name: validation.MaxLength, Rule: 50, Chain: nil},
				{Target: "serviceName", Name: validation.MinLength, Rule: 1, Chain: nil},
				{Target: "serviceName", Name: validation.Pattern, Rule: `^[a-zA-Z](?:[a-zA-Z0-9-]*[a-zA-Z0-9])?$`, Chain: nil}}}}); err != nil {
		return result, validation.NewErrorWithValidationError(err, "apimanagement.QuotaByPeriodKeysClient", "Get")
	}

	req, err := client.GetPreparer(resourceGroupName, serviceName, quotaCounterKey, quotaPeriodKey)
	if err != nil {
		err = autorest.NewErrorWithError(err, "apimanagement.QuotaByPeriodKeysClient", "Get", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "apimanagement.QuotaByPeriodKeysClient", "Get", resp, "Failure sending request")
		return
	}

	result, err = client.GetResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "apimanagement.QuotaByPeriodKeysClient", "Get", resp, "Failure responding to request")
	}

	return
}

// GetPreparer prepares the Get request.
func (client QuotaByPeriodKeysClient) GetPreparer(resourceGroupName string, serviceName string, quotaCounterKey string, quotaPeriodKey string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"quotaCounterKey":   autorest.Encode("path", quotaCounterKey),
		"quotaPeriodKey":    autorest.Encode("path", quotaPeriodKey),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"serviceName":       autorest.Encode("path", serviceName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2016-10-10"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ApiManagement/service/{serviceName}/quotas/{quotaCounterKey}/{quotaPeriodKey}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{})
}

// GetSender sends the Get request. The method will close the
// http.Response Body if it receives an error.
func (client QuotaByPeriodKeysClient) GetSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req)
}

// GetResponder handles the response to the Get request. The method always
// closes the http.Response Body.
func (client QuotaByPeriodKeysClient) GetResponder(resp *http.Response) (result QuotaCounterContract, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// Update updates an existing quota counter value in the specified service
// instance.
//
// resourceGroupName is the name of the resource group. serviceName is the name
// of the API Management service. quotaCounterKey is quota counter key
// identifier. quotaPeriodKey is quota period key identifier. parameters is the
// value of the Quota counter to be applied on the specified period.
func (client QuotaByPeriodKeysClient) Update(resourceGroupName string, serviceName string, quotaCounterKey string, quotaPeriodKey string, parameters QuotaCounterValueContract) (result autorest.Response, err error) {
	if err := validation.Validate([]validation.Validation{
		{TargetValue: serviceName,
			Constraints: []validation.Constraint{{Target: "serviceName", Name: validation.MaxLength, Rule: 50, Chain: nil},
				{Target: "serviceName", Name: validation.MinLength, Rule: 1, Chain: nil},
				{Target: "serviceName", Name: validation.Pattern, Rule: `^[a-zA-Z](?:[a-zA-Z0-9-]*[a-zA-Z0-9])?$`, Chain: nil}}}}); err != nil {
		return result, validation.NewErrorWithValidationError(err, "apimanagement.QuotaByPeriodKeysClient", "Update")
	}

	req, err := client.UpdatePreparer(resourceGroupName, serviceName, quotaCounterKey, quotaPeriodKey, parameters)
	if err != nil {
		err = autorest.NewErrorWithError(err, "apimanagement.QuotaByPeriodKeysClient", "Update", nil, "Failure preparing request")
		return
	}

	resp, err := client.UpdateSender(req)
	if err != nil {
		result.Response = resp
		err = autorest.NewErrorWithError(err, "apimanagement.QuotaByPeriodKeysClient", "Update", resp, "Failure sending request")
		return
	}

	result, err = client.UpdateResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "apimanagement.QuotaByPeriodKeysClient", "Update", resp, "Failure responding to request")
	}

	return
}

// UpdatePreparer prepares the Update request.
func (client QuotaByPeriodKeysClient) UpdatePreparer(resourceGroupName string, serviceName string, quotaCounterKey string, quotaPeriodKey string, parameters QuotaCounterValueContract) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"quotaCounterKey":   autorest.Encode("path", quotaCounterKey),
		"quotaPeriodKey":    autorest.Encode("path", quotaPeriodKey),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"serviceName":       autorest.Encode("path", serviceName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2016-10-10"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsJSON(),
		autorest.AsPatch(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ApiManagement/service/{serviceName}/quotas/{quotaCounterKey}/{quotaPeriodKey}", pathParameters),
		autorest.WithJSON(parameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{})
}

// UpdateSender sends the Update request. The method will close the
// http.Response Body if it receives an error.
func (client QuotaByPeriodKeysClient) UpdateSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req)
}

// UpdateResponder handles the response to the Update request. The method always
// closes the http.Response Body.
func (client QuotaByPeriodKeysClient) UpdateResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusNoContent),
		autorest.ByClosing())
	result.Response = resp
	return
}

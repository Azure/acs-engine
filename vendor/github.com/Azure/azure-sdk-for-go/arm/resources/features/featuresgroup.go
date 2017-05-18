package features

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
	"net/http"
)

// GroupClient is the azure Feature Exposure Control (AFEC) provides a
// mechanism for the resource providers to control feature exposure to users.
// Resource providers typically use this mechanism to provide public/private
// preview for new features prior to making them generally available. Users
// need to explicitly register for AFEC features to get access to such
// functionality.
type GroupClient struct {
	ManagementClient
}

// NewGroupClient creates an instance of the GroupClient client.
func NewGroupClient(subscriptionID string) GroupClient {
	return NewGroupClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewGroupClientWithBaseURI creates an instance of the GroupClient client.
func NewGroupClientWithBaseURI(baseURI string, subscriptionID string) GroupClient {
	return GroupClient{NewWithBaseURI(baseURI, subscriptionID)}
}

// Get gets the preview feature with the specified name.
//
// resourceProviderNamespace is the resource provider namespace for the
// feature. featureName is the name of the feature to get.
func (client GroupClient) Get(resourceProviderNamespace string, featureName string) (result Result, err error) {
	req, err := client.GetPreparer(resourceProviderNamespace, featureName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "features.GroupClient", "Get", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "features.GroupClient", "Get", resp, "Failure sending request")
		return
	}

	result, err = client.GetResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "features.GroupClient", "Get", resp, "Failure responding to request")
	}

	return
}

// GetPreparer prepares the Get request.
func (client GroupClient) GetPreparer(resourceProviderNamespace string, featureName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"featureName":               autorest.Encode("path", featureName),
		"resourceProviderNamespace": autorest.Encode("path", resourceProviderNamespace),
		"subscriptionId":            autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2015-12-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/providers/Microsoft.Features/providers/{resourceProviderNamespace}/features/{featureName}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{})
}

// GetSender sends the Get request. The method will close the
// http.Response Body if it receives an error.
func (client GroupClient) GetSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req)
}

// GetResponder handles the response to the Get request. The method always
// closes the http.Response Body.
func (client GroupClient) GetResponder(resp *http.Response) (result Result, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// List gets all the preview features in a provider namespace that are
// available through AFEC for the subscription.
//
// resourceProviderNamespace is the namespace of the resource provider for
// getting features.
func (client GroupClient) List(resourceProviderNamespace string) (result OperationsListResult, err error) {
	req, err := client.ListPreparer(resourceProviderNamespace)
	if err != nil {
		err = autorest.NewErrorWithError(err, "features.GroupClient", "List", nil, "Failure preparing request")
		return
	}

	resp, err := client.ListSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "features.GroupClient", "List", resp, "Failure sending request")
		return
	}

	result, err = client.ListResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "features.GroupClient", "List", resp, "Failure responding to request")
	}

	return
}

// ListPreparer prepares the List request.
func (client GroupClient) ListPreparer(resourceProviderNamespace string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"resourceProviderNamespace": autorest.Encode("path", resourceProviderNamespace),
		"subscriptionId":            autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2015-12-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/providers/Microsoft.Features/providers/{resourceProviderNamespace}/features", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{})
}

// ListSender sends the List request. The method will close the
// http.Response Body if it receives an error.
func (client GroupClient) ListSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req)
}

// ListResponder handles the response to the List request. The method always
// closes the http.Response Body.
func (client GroupClient) ListResponder(resp *http.Response) (result OperationsListResult, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// ListNextResults retrieves the next set of results, if any.
func (client GroupClient) ListNextResults(lastResults OperationsListResult) (result OperationsListResult, err error) {
	req, err := lastResults.OperationsListResultPreparer()
	if err != nil {
		return result, autorest.NewErrorWithError(err, "features.GroupClient", "List", nil, "Failure preparing next results request")
	}
	if req == nil {
		return
	}

	resp, err := client.ListSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "features.GroupClient", "List", resp, "Failure sending next results request")
	}

	result, err = client.ListResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "features.GroupClient", "List", resp, "Failure responding to next results request")
	}

	return
}

// ListAll gets all the preview features that are available through AFEC for
// the subscription.
func (client GroupClient) ListAll() (result OperationsListResult, err error) {
	req, err := client.ListAllPreparer()
	if err != nil {
		err = autorest.NewErrorWithError(err, "features.GroupClient", "ListAll", nil, "Failure preparing request")
		return
	}

	resp, err := client.ListAllSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "features.GroupClient", "ListAll", resp, "Failure sending request")
		return
	}

	result, err = client.ListAllResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "features.GroupClient", "ListAll", resp, "Failure responding to request")
	}

	return
}

// ListAllPreparer prepares the ListAll request.
func (client GroupClient) ListAllPreparer() (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"subscriptionId": autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2015-12-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/providers/Microsoft.Features/features", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{})
}

// ListAllSender sends the ListAll request. The method will close the
// http.Response Body if it receives an error.
func (client GroupClient) ListAllSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req)
}

// ListAllResponder handles the response to the ListAll request. The method always
// closes the http.Response Body.
func (client GroupClient) ListAllResponder(resp *http.Response) (result OperationsListResult, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// ListAllNextResults retrieves the next set of results, if any.
func (client GroupClient) ListAllNextResults(lastResults OperationsListResult) (result OperationsListResult, err error) {
	req, err := lastResults.OperationsListResultPreparer()
	if err != nil {
		return result, autorest.NewErrorWithError(err, "features.GroupClient", "ListAll", nil, "Failure preparing next results request")
	}
	if req == nil {
		return
	}

	resp, err := client.ListAllSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "features.GroupClient", "ListAll", resp, "Failure sending next results request")
	}

	result, err = client.ListAllResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "features.GroupClient", "ListAll", resp, "Failure responding to next results request")
	}

	return
}

// Register registers the preview feature for the subscription.
//
// resourceProviderNamespace is the namespace of the resource provider.
// featureName is the name of the feature to register.
func (client GroupClient) Register(resourceProviderNamespace string, featureName string) (result Result, err error) {
	req, err := client.RegisterPreparer(resourceProviderNamespace, featureName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "features.GroupClient", "Register", nil, "Failure preparing request")
		return
	}

	resp, err := client.RegisterSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "features.GroupClient", "Register", resp, "Failure sending request")
		return
	}

	result, err = client.RegisterResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "features.GroupClient", "Register", resp, "Failure responding to request")
	}

	return
}

// RegisterPreparer prepares the Register request.
func (client GroupClient) RegisterPreparer(resourceProviderNamespace string, featureName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"featureName":               autorest.Encode("path", featureName),
		"resourceProviderNamespace": autorest.Encode("path", resourceProviderNamespace),
		"subscriptionId":            autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2015-12-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsPost(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/providers/Microsoft.Features/providers/{resourceProviderNamespace}/features/{featureName}/register", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{})
}

// RegisterSender sends the Register request. The method will close the
// http.Response Body if it receives an error.
func (client GroupClient) RegisterSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req)
}

// RegisterResponder handles the response to the Register request. The method always
// closes the http.Response Body.
func (client GroupClient) RegisterResponder(resp *http.Response) (result Result, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

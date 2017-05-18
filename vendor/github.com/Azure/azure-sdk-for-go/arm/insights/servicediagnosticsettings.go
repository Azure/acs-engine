package insights

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

// ServiceDiagnosticSettingsClient is the composite Swagger for Insights
// Management Client
type ServiceDiagnosticSettingsClient struct {
	ManagementClient
}

// NewServiceDiagnosticSettingsClient creates an instance of the
// ServiceDiagnosticSettingsClient client.
func NewServiceDiagnosticSettingsClient(subscriptionID string) ServiceDiagnosticSettingsClient {
	return NewServiceDiagnosticSettingsClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewServiceDiagnosticSettingsClientWithBaseURI creates an instance of the
// ServiceDiagnosticSettingsClient client.
func NewServiceDiagnosticSettingsClientWithBaseURI(baseURI string, subscriptionID string) ServiceDiagnosticSettingsClient {
	return ServiceDiagnosticSettingsClient{NewWithBaseURI(baseURI, subscriptionID)}
}

// CreateOrUpdate create or update new diagnostic settings for the specified
// resource.
//
// resourceURI is the identifier of the resource. parameters is parameters
// supplied to the operation.
func (client ServiceDiagnosticSettingsClient) CreateOrUpdate(resourceURI string, parameters ServiceDiagnosticSettingsResource) (result ServiceDiagnosticSettingsResource, err error) {
	req, err := client.CreateOrUpdatePreparer(resourceURI, parameters)
	if err != nil {
		err = autorest.NewErrorWithError(err, "insights.ServiceDiagnosticSettingsClient", "CreateOrUpdate", nil, "Failure preparing request")
		return
	}

	resp, err := client.CreateOrUpdateSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "insights.ServiceDiagnosticSettingsClient", "CreateOrUpdate", resp, "Failure sending request")
		return
	}

	result, err = client.CreateOrUpdateResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "insights.ServiceDiagnosticSettingsClient", "CreateOrUpdate", resp, "Failure responding to request")
	}

	return
}

// CreateOrUpdatePreparer prepares the CreateOrUpdate request.
func (client ServiceDiagnosticSettingsClient) CreateOrUpdatePreparer(resourceURI string, parameters ServiceDiagnosticSettingsResource) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"resourceUri": autorest.Encode("path", resourceURI),
	}

	const APIVersion = "2015-07-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsJSON(),
		autorest.AsPut(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/{resourceUri}/providers/microsoft.insights/diagnosticSettings/service", pathParameters),
		autorest.WithJSON(parameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{})
}

// CreateOrUpdateSender sends the CreateOrUpdate request. The method will close the
// http.Response Body if it receives an error.
func (client ServiceDiagnosticSettingsClient) CreateOrUpdateSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req)
}

// CreateOrUpdateResponder handles the response to the CreateOrUpdate request. The method always
// closes the http.Response Body.
func (client ServiceDiagnosticSettingsClient) CreateOrUpdateResponder(resp *http.Response) (result ServiceDiagnosticSettingsResource, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// Get gets the active diagnostic settings for the specified resource.
//
// resourceURI is the identifier of the resource.
func (client ServiceDiagnosticSettingsClient) Get(resourceURI string) (result ServiceDiagnosticSettingsResource, err error) {
	req, err := client.GetPreparer(resourceURI)
	if err != nil {
		err = autorest.NewErrorWithError(err, "insights.ServiceDiagnosticSettingsClient", "Get", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "insights.ServiceDiagnosticSettingsClient", "Get", resp, "Failure sending request")
		return
	}

	result, err = client.GetResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "insights.ServiceDiagnosticSettingsClient", "Get", resp, "Failure responding to request")
	}

	return
}

// GetPreparer prepares the Get request.
func (client ServiceDiagnosticSettingsClient) GetPreparer(resourceURI string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"resourceUri": autorest.Encode("path", resourceURI),
	}

	const APIVersion = "2015-07-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/{resourceUri}/providers/microsoft.insights/diagnosticSettings/service", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{})
}

// GetSender sends the Get request. The method will close the
// http.Response Body if it receives an error.
func (client ServiceDiagnosticSettingsClient) GetSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req)
}

// GetResponder handles the response to the Get request. The method always
// closes the http.Response Body.
func (client ServiceDiagnosticSettingsClient) GetResponder(resp *http.Response) (result ServiceDiagnosticSettingsResource, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

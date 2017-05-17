package web

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

// RecommendationsClient is the composite Swagger for WebSite Management Client
type RecommendationsClient struct {
	ManagementClient
}

// NewRecommendationsClient creates an instance of the RecommendationsClient
// client.
func NewRecommendationsClient(subscriptionID string) RecommendationsClient {
	return NewRecommendationsClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewRecommendationsClientWithBaseURI creates an instance of the
// RecommendationsClient client.
func NewRecommendationsClientWithBaseURI(baseURI string, subscriptionID string) RecommendationsClient {
	return RecommendationsClient{NewWithBaseURI(baseURI, subscriptionID)}
}

// DisableAllForWebApp disable all recommendations for an app.
//
// resourceGroupName is name of the resource group to which the resource
// belongs. siteName is name of the app.
func (client RecommendationsClient) DisableAllForWebApp(resourceGroupName string, siteName string) (result autorest.Response, err error) {
	if err := validation.Validate([]validation.Validation{
		{TargetValue: resourceGroupName,
			Constraints: []validation.Constraint{{Target: "resourceGroupName", Name: validation.MaxLength, Rule: 90, Chain: nil},
				{Target: "resourceGroupName", Name: validation.MinLength, Rule: 1, Chain: nil},
				{Target: "resourceGroupName", Name: validation.Pattern, Rule: `^[-\w\._\(\)]+[^\.]$`, Chain: nil}}}}); err != nil {
		return result, validation.NewErrorWithValidationError(err, "web.RecommendationsClient", "DisableAllForWebApp")
	}

	req, err := client.DisableAllForWebAppPreparer(resourceGroupName, siteName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "web.RecommendationsClient", "DisableAllForWebApp", nil, "Failure preparing request")
		return
	}

	resp, err := client.DisableAllForWebAppSender(req)
	if err != nil {
		result.Response = resp
		err = autorest.NewErrorWithError(err, "web.RecommendationsClient", "DisableAllForWebApp", resp, "Failure sending request")
		return
	}

	result, err = client.DisableAllForWebAppResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "web.RecommendationsClient", "DisableAllForWebApp", resp, "Failure responding to request")
	}

	return
}

// DisableAllForWebAppPreparer prepares the DisableAllForWebApp request.
func (client RecommendationsClient) DisableAllForWebAppPreparer(resourceGroupName string, siteName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"siteName":          autorest.Encode("path", siteName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2016-03-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsPost(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Web/sites/{siteName}/recommendations/disable", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{})
}

// DisableAllForWebAppSender sends the DisableAllForWebApp request. The method will close the
// http.Response Body if it receives an error.
func (client RecommendationsClient) DisableAllForWebAppSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req)
}

// DisableAllForWebAppResponder handles the response to the DisableAllForWebApp request. The method always
// closes the http.Response Body.
func (client RecommendationsClient) DisableAllForWebAppResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusNoContent),
		autorest.ByClosing())
	result.Response = resp
	return
}

// GetRuleDetailsByWebApp get a recommendation rule for an app.
//
// resourceGroupName is name of the resource group to which the resource
// belongs. siteName is name of the app. name is name of the recommendation.
// updateSeen is specify <code>true</code> to update the last-seen timestamp of
// the recommendation object.
func (client RecommendationsClient) GetRuleDetailsByWebApp(resourceGroupName string, siteName string, name string, updateSeen *bool) (result RecommendationRule, err error) {
	if err := validation.Validate([]validation.Validation{
		{TargetValue: resourceGroupName,
			Constraints: []validation.Constraint{{Target: "resourceGroupName", Name: validation.MaxLength, Rule: 90, Chain: nil},
				{Target: "resourceGroupName", Name: validation.MinLength, Rule: 1, Chain: nil},
				{Target: "resourceGroupName", Name: validation.Pattern, Rule: `^[-\w\._\(\)]+[^\.]$`, Chain: nil}}}}); err != nil {
		return result, validation.NewErrorWithValidationError(err, "web.RecommendationsClient", "GetRuleDetailsByWebApp")
	}

	req, err := client.GetRuleDetailsByWebAppPreparer(resourceGroupName, siteName, name, updateSeen)
	if err != nil {
		err = autorest.NewErrorWithError(err, "web.RecommendationsClient", "GetRuleDetailsByWebApp", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetRuleDetailsByWebAppSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "web.RecommendationsClient", "GetRuleDetailsByWebApp", resp, "Failure sending request")
		return
	}

	result, err = client.GetRuleDetailsByWebAppResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "web.RecommendationsClient", "GetRuleDetailsByWebApp", resp, "Failure responding to request")
	}

	return
}

// GetRuleDetailsByWebAppPreparer prepares the GetRuleDetailsByWebApp request.
func (client RecommendationsClient) GetRuleDetailsByWebAppPreparer(resourceGroupName string, siteName string, name string, updateSeen *bool) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"name":              autorest.Encode("path", name),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"siteName":          autorest.Encode("path", siteName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2016-03-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}
	if updateSeen != nil {
		queryParameters["updateSeen"] = autorest.Encode("query", *updateSeen)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Web/sites/{siteName}/recommendations/{name}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{})
}

// GetRuleDetailsByWebAppSender sends the GetRuleDetailsByWebApp request. The method will close the
// http.Response Body if it receives an error.
func (client RecommendationsClient) GetRuleDetailsByWebAppSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req)
}

// GetRuleDetailsByWebAppResponder handles the response to the GetRuleDetailsByWebApp request. The method always
// closes the http.Response Body.
func (client RecommendationsClient) GetRuleDetailsByWebAppResponder(resp *http.Response) (result RecommendationRule, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// List list all recommendations for a subscription.
//
// featured is specify <code>true</code> to return only the most critical
// recommendations. The default is <code>false</code>, which returns all
// recommendations. filter is filter is specified by using OData syntax.
// Example: $filter=channels eq 'Api' or channel eq 'Notification' and
// startTime eq '2014-01-01T00:00:00Z' and endTime eq '2014-12-31T23:59:59Z'
// and timeGrain eq duration'[PT1H|PT1M|P1D]
func (client RecommendationsClient) List(featured *bool, filter string) (result ListRecommendation, err error) {
	req, err := client.ListPreparer(featured, filter)
	if err != nil {
		err = autorest.NewErrorWithError(err, "web.RecommendationsClient", "List", nil, "Failure preparing request")
		return
	}

	resp, err := client.ListSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "web.RecommendationsClient", "List", resp, "Failure sending request")
		return
	}

	result, err = client.ListResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "web.RecommendationsClient", "List", resp, "Failure responding to request")
	}

	return
}

// ListPreparer prepares the List request.
func (client RecommendationsClient) ListPreparer(featured *bool, filter string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"subscriptionId": autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2016-03-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}
	if featured != nil {
		queryParameters["featured"] = autorest.Encode("query", *featured)
	}
	if len(filter) > 0 {
		queryParameters["$filter"] = filter
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/providers/Microsoft.Web/recommendations", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{})
}

// ListSender sends the List request. The method will close the
// http.Response Body if it receives an error.
func (client RecommendationsClient) ListSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req)
}

// ListResponder handles the response to the List request. The method always
// closes the http.Response Body.
func (client RecommendationsClient) ListResponder(resp *http.Response) (result ListRecommendation, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Value),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// ListHistoryForWebApp get past recommendations for an app, optionally
// specified by the time range.
//
// resourceGroupName is name of the resource group to which the resource
// belongs. siteName is name of the app. filter is filter is specified by using
// OData syntax. Example: $filter=channels eq 'Api' or channel eq
// 'Notification' and startTime eq '2014-01-01T00:00:00Z' and endTime eq
// '2014-12-31T23:59:59Z' and timeGrain eq duration'[PT1H|PT1M|P1D]
func (client RecommendationsClient) ListHistoryForWebApp(resourceGroupName string, siteName string, filter string) (result ListRecommendation, err error) {
	if err := validation.Validate([]validation.Validation{
		{TargetValue: resourceGroupName,
			Constraints: []validation.Constraint{{Target: "resourceGroupName", Name: validation.MaxLength, Rule: 90, Chain: nil},
				{Target: "resourceGroupName", Name: validation.MinLength, Rule: 1, Chain: nil},
				{Target: "resourceGroupName", Name: validation.Pattern, Rule: `^[-\w\._\(\)]+[^\.]$`, Chain: nil}}}}); err != nil {
		return result, validation.NewErrorWithValidationError(err, "web.RecommendationsClient", "ListHistoryForWebApp")
	}

	req, err := client.ListHistoryForWebAppPreparer(resourceGroupName, siteName, filter)
	if err != nil {
		err = autorest.NewErrorWithError(err, "web.RecommendationsClient", "ListHistoryForWebApp", nil, "Failure preparing request")
		return
	}

	resp, err := client.ListHistoryForWebAppSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "web.RecommendationsClient", "ListHistoryForWebApp", resp, "Failure sending request")
		return
	}

	result, err = client.ListHistoryForWebAppResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "web.RecommendationsClient", "ListHistoryForWebApp", resp, "Failure responding to request")
	}

	return
}

// ListHistoryForWebAppPreparer prepares the ListHistoryForWebApp request.
func (client RecommendationsClient) ListHistoryForWebAppPreparer(resourceGroupName string, siteName string, filter string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"siteName":          autorest.Encode("path", siteName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2016-03-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}
	if len(filter) > 0 {
		queryParameters["$filter"] = filter
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Web/sites/{siteName}/recommendationHistory", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{})
}

// ListHistoryForWebAppSender sends the ListHistoryForWebApp request. The method will close the
// http.Response Body if it receives an error.
func (client RecommendationsClient) ListHistoryForWebAppSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req)
}

// ListHistoryForWebAppResponder handles the response to the ListHistoryForWebApp request. The method always
// closes the http.Response Body.
func (client RecommendationsClient) ListHistoryForWebAppResponder(resp *http.Response) (result ListRecommendation, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Value),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// ListRecommendedRulesForWebApp get all recommendations for an app.
//
// resourceGroupName is name of the resource group to which the resource
// belongs. siteName is name of the app. featured is specify <code>true</code>
// to return only the most critical recommendations. The default is
// <code>false</code>, which returns all recommendations. filter is return only
// channels specified in the filter. Filter is specified by using OData syntax.
// Example: $filter=channels eq 'Api' or channel eq 'Notification'
func (client RecommendationsClient) ListRecommendedRulesForWebApp(resourceGroupName string, siteName string, featured *bool, filter string) (result ListRecommendation, err error) {
	if err := validation.Validate([]validation.Validation{
		{TargetValue: resourceGroupName,
			Constraints: []validation.Constraint{{Target: "resourceGroupName", Name: validation.MaxLength, Rule: 90, Chain: nil},
				{Target: "resourceGroupName", Name: validation.MinLength, Rule: 1, Chain: nil},
				{Target: "resourceGroupName", Name: validation.Pattern, Rule: `^[-\w\._\(\)]+[^\.]$`, Chain: nil}}}}); err != nil {
		return result, validation.NewErrorWithValidationError(err, "web.RecommendationsClient", "ListRecommendedRulesForWebApp")
	}

	req, err := client.ListRecommendedRulesForWebAppPreparer(resourceGroupName, siteName, featured, filter)
	if err != nil {
		err = autorest.NewErrorWithError(err, "web.RecommendationsClient", "ListRecommendedRulesForWebApp", nil, "Failure preparing request")
		return
	}

	resp, err := client.ListRecommendedRulesForWebAppSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "web.RecommendationsClient", "ListRecommendedRulesForWebApp", resp, "Failure sending request")
		return
	}

	result, err = client.ListRecommendedRulesForWebAppResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "web.RecommendationsClient", "ListRecommendedRulesForWebApp", resp, "Failure responding to request")
	}

	return
}

// ListRecommendedRulesForWebAppPreparer prepares the ListRecommendedRulesForWebApp request.
func (client RecommendationsClient) ListRecommendedRulesForWebAppPreparer(resourceGroupName string, siteName string, featured *bool, filter string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"siteName":          autorest.Encode("path", siteName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2016-03-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}
	if featured != nil {
		queryParameters["featured"] = autorest.Encode("query", *featured)
	}
	if len(filter) > 0 {
		queryParameters["$filter"] = filter
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Web/sites/{siteName}/recommendations", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{})
}

// ListRecommendedRulesForWebAppSender sends the ListRecommendedRulesForWebApp request. The method will close the
// http.Response Body if it receives an error.
func (client RecommendationsClient) ListRecommendedRulesForWebAppSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req)
}

// ListRecommendedRulesForWebAppResponder handles the response to the ListRecommendedRulesForWebApp request. The method always
// closes the http.Response Body.
func (client RecommendationsClient) ListRecommendedRulesForWebAppResponder(resp *http.Response) (result ListRecommendation, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Value),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// ResetAllFilters reset all recommendation opt-out settings for a
// subscription.
func (client RecommendationsClient) ResetAllFilters() (result autorest.Response, err error) {
	req, err := client.ResetAllFiltersPreparer()
	if err != nil {
		err = autorest.NewErrorWithError(err, "web.RecommendationsClient", "ResetAllFilters", nil, "Failure preparing request")
		return
	}

	resp, err := client.ResetAllFiltersSender(req)
	if err != nil {
		result.Response = resp
		err = autorest.NewErrorWithError(err, "web.RecommendationsClient", "ResetAllFilters", resp, "Failure sending request")
		return
	}

	result, err = client.ResetAllFiltersResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "web.RecommendationsClient", "ResetAllFilters", resp, "Failure responding to request")
	}

	return
}

// ResetAllFiltersPreparer prepares the ResetAllFilters request.
func (client RecommendationsClient) ResetAllFiltersPreparer() (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"subscriptionId": autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2016-03-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsPost(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/providers/Microsoft.Web/recommendations/reset", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{})
}

// ResetAllFiltersSender sends the ResetAllFilters request. The method will close the
// http.Response Body if it receives an error.
func (client RecommendationsClient) ResetAllFiltersSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req)
}

// ResetAllFiltersResponder handles the response to the ResetAllFilters request. The method always
// closes the http.Response Body.
func (client RecommendationsClient) ResetAllFiltersResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusNoContent),
		autorest.ByClosing())
	result.Response = resp
	return
}

// ResetAllFiltersForWebApp reset all recommendation opt-out settings for an
// app.
//
// resourceGroupName is name of the resource group to which the resource
// belongs. siteName is name of the app.
func (client RecommendationsClient) ResetAllFiltersForWebApp(resourceGroupName string, siteName string) (result autorest.Response, err error) {
	if err := validation.Validate([]validation.Validation{
		{TargetValue: resourceGroupName,
			Constraints: []validation.Constraint{{Target: "resourceGroupName", Name: validation.MaxLength, Rule: 90, Chain: nil},
				{Target: "resourceGroupName", Name: validation.MinLength, Rule: 1, Chain: nil},
				{Target: "resourceGroupName", Name: validation.Pattern, Rule: `^[-\w\._\(\)]+[^\.]$`, Chain: nil}}}}); err != nil {
		return result, validation.NewErrorWithValidationError(err, "web.RecommendationsClient", "ResetAllFiltersForWebApp")
	}

	req, err := client.ResetAllFiltersForWebAppPreparer(resourceGroupName, siteName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "web.RecommendationsClient", "ResetAllFiltersForWebApp", nil, "Failure preparing request")
		return
	}

	resp, err := client.ResetAllFiltersForWebAppSender(req)
	if err != nil {
		result.Response = resp
		err = autorest.NewErrorWithError(err, "web.RecommendationsClient", "ResetAllFiltersForWebApp", resp, "Failure sending request")
		return
	}

	result, err = client.ResetAllFiltersForWebAppResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "web.RecommendationsClient", "ResetAllFiltersForWebApp", resp, "Failure responding to request")
	}

	return
}

// ResetAllFiltersForWebAppPreparer prepares the ResetAllFiltersForWebApp request.
func (client RecommendationsClient) ResetAllFiltersForWebAppPreparer(resourceGroupName string, siteName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"siteName":          autorest.Encode("path", siteName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2016-03-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsPost(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Web/sites/{siteName}/recommendations/reset", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{})
}

// ResetAllFiltersForWebAppSender sends the ResetAllFiltersForWebApp request. The method will close the
// http.Response Body if it receives an error.
func (client RecommendationsClient) ResetAllFiltersForWebAppSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req)
}

// ResetAllFiltersForWebAppResponder handles the response to the ResetAllFiltersForWebApp request. The method always
// closes the http.Response Body.
func (client RecommendationsClient) ResetAllFiltersForWebAppResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusNoContent),
		autorest.ByClosing())
	result.Response = resp
	return
}

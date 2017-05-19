package scheduler

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

// JobsClient is the client for the Jobs methods of the Scheduler service.
type JobsClient struct {
	ManagementClient
}

// NewJobsClient creates an instance of the JobsClient client.
func NewJobsClient(subscriptionID string) JobsClient {
	return NewJobsClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewJobsClientWithBaseURI creates an instance of the JobsClient client.
func NewJobsClientWithBaseURI(baseURI string, subscriptionID string) JobsClient {
	return JobsClient{NewWithBaseURI(baseURI, subscriptionID)}
}

// CreateOrUpdate provisions a new job or updates an existing job.
//
// resourceGroupName is the resource group name. jobCollectionName is the job
// collection name. jobName is the job name. job is the job definition.
func (client JobsClient) CreateOrUpdate(resourceGroupName string, jobCollectionName string, jobName string, job JobDefinition) (result JobDefinition, err error) {
	req, err := client.CreateOrUpdatePreparer(resourceGroupName, jobCollectionName, jobName, job)
	if err != nil {
		err = autorest.NewErrorWithError(err, "scheduler.JobsClient", "CreateOrUpdate", nil, "Failure preparing request")
		return
	}

	resp, err := client.CreateOrUpdateSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "scheduler.JobsClient", "CreateOrUpdate", resp, "Failure sending request")
		return
	}

	result, err = client.CreateOrUpdateResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "scheduler.JobsClient", "CreateOrUpdate", resp, "Failure responding to request")
	}

	return
}

// CreateOrUpdatePreparer prepares the CreateOrUpdate request.
func (client JobsClient) CreateOrUpdatePreparer(resourceGroupName string, jobCollectionName string, jobName string, job JobDefinition) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"jobCollectionName": autorest.Encode("path", jobCollectionName),
		"jobName":           autorest.Encode("path", jobName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2016-03-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsJSON(),
		autorest.AsPut(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Scheduler/jobCollections/{jobCollectionName}/jobs/{jobName}", pathParameters),
		autorest.WithJSON(job),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{})
}

// CreateOrUpdateSender sends the CreateOrUpdate request. The method will close the
// http.Response Body if it receives an error.
func (client JobsClient) CreateOrUpdateSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req)
}

// CreateOrUpdateResponder handles the response to the CreateOrUpdate request. The method always
// closes the http.Response Body.
func (client JobsClient) CreateOrUpdateResponder(resp *http.Response) (result JobDefinition, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusCreated),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// Delete deletes a job.
//
// resourceGroupName is the resource group name. jobCollectionName is the job
// collection name. jobName is the job name.
func (client JobsClient) Delete(resourceGroupName string, jobCollectionName string, jobName string) (result autorest.Response, err error) {
	req, err := client.DeletePreparer(resourceGroupName, jobCollectionName, jobName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "scheduler.JobsClient", "Delete", nil, "Failure preparing request")
		return
	}

	resp, err := client.DeleteSender(req)
	if err != nil {
		result.Response = resp
		err = autorest.NewErrorWithError(err, "scheduler.JobsClient", "Delete", resp, "Failure sending request")
		return
	}

	result, err = client.DeleteResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "scheduler.JobsClient", "Delete", resp, "Failure responding to request")
	}

	return
}

// DeletePreparer prepares the Delete request.
func (client JobsClient) DeletePreparer(resourceGroupName string, jobCollectionName string, jobName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"jobCollectionName": autorest.Encode("path", jobCollectionName),
		"jobName":           autorest.Encode("path", jobName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2016-03-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsDelete(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Scheduler/jobCollections/{jobCollectionName}/jobs/{jobName}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{})
}

// DeleteSender sends the Delete request. The method will close the
// http.Response Body if it receives an error.
func (client JobsClient) DeleteSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req)
}

// DeleteResponder handles the response to the Delete request. The method always
// closes the http.Response Body.
func (client JobsClient) DeleteResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByClosing())
	result.Response = resp
	return
}

// Get gets a job.
//
// resourceGroupName is the resource group name. jobCollectionName is the job
// collection name. jobName is the job name.
func (client JobsClient) Get(resourceGroupName string, jobCollectionName string, jobName string) (result JobDefinition, err error) {
	req, err := client.GetPreparer(resourceGroupName, jobCollectionName, jobName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "scheduler.JobsClient", "Get", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "scheduler.JobsClient", "Get", resp, "Failure sending request")
		return
	}

	result, err = client.GetResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "scheduler.JobsClient", "Get", resp, "Failure responding to request")
	}

	return
}

// GetPreparer prepares the Get request.
func (client JobsClient) GetPreparer(resourceGroupName string, jobCollectionName string, jobName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"jobCollectionName": autorest.Encode("path", jobCollectionName),
		"jobName":           autorest.Encode("path", jobName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2016-03-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Scheduler/jobCollections/{jobCollectionName}/jobs/{jobName}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{})
}

// GetSender sends the Get request. The method will close the
// http.Response Body if it receives an error.
func (client JobsClient) GetSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req)
}

// GetResponder handles the response to the Get request. The method always
// closes the http.Response Body.
func (client JobsClient) GetResponder(resp *http.Response) (result JobDefinition, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// List lists all jobs under the specified job collection.
//
// resourceGroupName is the resource group name. jobCollectionName is the job
// collection name. top is the number of jobs to request, in the of range of
// [1..100]. skip is the (0-based) index of the job history list from which to
// begin requesting entries. filter is the filter to apply on the job state.
func (client JobsClient) List(resourceGroupName string, jobCollectionName string, top *int32, skip *int32, filter string) (result JobListResult, err error) {
	if err := validation.Validate([]validation.Validation{
		{TargetValue: top,
			Constraints: []validation.Constraint{{Target: "top", Name: validation.Null, Rule: false,
				Chain: []validation.Constraint{{Target: "top", Name: validation.InclusiveMaximum, Rule: 100, Chain: nil},
					{Target: "top", Name: validation.InclusiveMinimum, Rule: 1, Chain: nil},
				}}}}}); err != nil {
		return result, validation.NewErrorWithValidationError(err, "scheduler.JobsClient", "List")
	}

	req, err := client.ListPreparer(resourceGroupName, jobCollectionName, top, skip, filter)
	if err != nil {
		err = autorest.NewErrorWithError(err, "scheduler.JobsClient", "List", nil, "Failure preparing request")
		return
	}

	resp, err := client.ListSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "scheduler.JobsClient", "List", resp, "Failure sending request")
		return
	}

	result, err = client.ListResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "scheduler.JobsClient", "List", resp, "Failure responding to request")
	}

	return
}

// ListPreparer prepares the List request.
func (client JobsClient) ListPreparer(resourceGroupName string, jobCollectionName string, top *int32, skip *int32, filter string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"jobCollectionName": autorest.Encode("path", jobCollectionName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2016-03-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}
	if top != nil {
		queryParameters["$top"] = autorest.Encode("query", *top)
	}
	if skip != nil {
		queryParameters["$skip"] = autorest.Encode("query", *skip)
	}
	if len(filter) > 0 {
		queryParameters["$filter"] = autorest.Encode("query", filter)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Scheduler/jobCollections/{jobCollectionName}/jobs", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{})
}

// ListSender sends the List request. The method will close the
// http.Response Body if it receives an error.
func (client JobsClient) ListSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req)
}

// ListResponder handles the response to the List request. The method always
// closes the http.Response Body.
func (client JobsClient) ListResponder(resp *http.Response) (result JobListResult, err error) {
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
func (client JobsClient) ListNextResults(lastResults JobListResult) (result JobListResult, err error) {
	req, err := lastResults.JobListResultPreparer()
	if err != nil {
		return result, autorest.NewErrorWithError(err, "scheduler.JobsClient", "List", nil, "Failure preparing next results request")
	}
	if req == nil {
		return
	}

	resp, err := client.ListSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "scheduler.JobsClient", "List", resp, "Failure sending next results request")
	}

	result, err = client.ListResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "scheduler.JobsClient", "List", resp, "Failure responding to next results request")
	}

	return
}

// ListJobHistory lists job history.
//
// resourceGroupName is the resource group name. jobCollectionName is the job
// collection name. jobName is the job name. top is the number of job history
// to request, in the of range of [1..100]. skip is the (0-based) index of the
// job history list from which to begin requesting entries. filter is the
// filter to apply on the job state.
func (client JobsClient) ListJobHistory(resourceGroupName string, jobCollectionName string, jobName string, top *int32, skip *int32, filter string) (result JobHistoryListResult, err error) {
	if err := validation.Validate([]validation.Validation{
		{TargetValue: top,
			Constraints: []validation.Constraint{{Target: "top", Name: validation.Null, Rule: false,
				Chain: []validation.Constraint{{Target: "top", Name: validation.InclusiveMaximum, Rule: 100, Chain: nil},
					{Target: "top", Name: validation.InclusiveMinimum, Rule: 1, Chain: nil},
				}}}}}); err != nil {
		return result, validation.NewErrorWithValidationError(err, "scheduler.JobsClient", "ListJobHistory")
	}

	req, err := client.ListJobHistoryPreparer(resourceGroupName, jobCollectionName, jobName, top, skip, filter)
	if err != nil {
		err = autorest.NewErrorWithError(err, "scheduler.JobsClient", "ListJobHistory", nil, "Failure preparing request")
		return
	}

	resp, err := client.ListJobHistorySender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "scheduler.JobsClient", "ListJobHistory", resp, "Failure sending request")
		return
	}

	result, err = client.ListJobHistoryResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "scheduler.JobsClient", "ListJobHistory", resp, "Failure responding to request")
	}

	return
}

// ListJobHistoryPreparer prepares the ListJobHistory request.
func (client JobsClient) ListJobHistoryPreparer(resourceGroupName string, jobCollectionName string, jobName string, top *int32, skip *int32, filter string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"jobCollectionName": autorest.Encode("path", jobCollectionName),
		"jobName":           autorest.Encode("path", jobName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2016-03-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}
	if top != nil {
		queryParameters["$top"] = autorest.Encode("query", *top)
	}
	if skip != nil {
		queryParameters["$skip"] = autorest.Encode("query", *skip)
	}
	if len(filter) > 0 {
		queryParameters["$filter"] = autorest.Encode("query", filter)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Scheduler/jobCollections/{jobCollectionName}/jobs/{jobName}/history", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{})
}

// ListJobHistorySender sends the ListJobHistory request. The method will close the
// http.Response Body if it receives an error.
func (client JobsClient) ListJobHistorySender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req)
}

// ListJobHistoryResponder handles the response to the ListJobHistory request. The method always
// closes the http.Response Body.
func (client JobsClient) ListJobHistoryResponder(resp *http.Response) (result JobHistoryListResult, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// ListJobHistoryNextResults retrieves the next set of results, if any.
func (client JobsClient) ListJobHistoryNextResults(lastResults JobHistoryListResult) (result JobHistoryListResult, err error) {
	req, err := lastResults.JobHistoryListResultPreparer()
	if err != nil {
		return result, autorest.NewErrorWithError(err, "scheduler.JobsClient", "ListJobHistory", nil, "Failure preparing next results request")
	}
	if req == nil {
		return
	}

	resp, err := client.ListJobHistorySender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "scheduler.JobsClient", "ListJobHistory", resp, "Failure sending next results request")
	}

	result, err = client.ListJobHistoryResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "scheduler.JobsClient", "ListJobHistory", resp, "Failure responding to next results request")
	}

	return
}

// Patch patches an existing job.
//
// resourceGroupName is the resource group name. jobCollectionName is the job
// collection name. jobName is the job name. job is the job definition.
func (client JobsClient) Patch(resourceGroupName string, jobCollectionName string, jobName string, job JobDefinition) (result JobDefinition, err error) {
	req, err := client.PatchPreparer(resourceGroupName, jobCollectionName, jobName, job)
	if err != nil {
		err = autorest.NewErrorWithError(err, "scheduler.JobsClient", "Patch", nil, "Failure preparing request")
		return
	}

	resp, err := client.PatchSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "scheduler.JobsClient", "Patch", resp, "Failure sending request")
		return
	}

	result, err = client.PatchResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "scheduler.JobsClient", "Patch", resp, "Failure responding to request")
	}

	return
}

// PatchPreparer prepares the Patch request.
func (client JobsClient) PatchPreparer(resourceGroupName string, jobCollectionName string, jobName string, job JobDefinition) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"jobCollectionName": autorest.Encode("path", jobCollectionName),
		"jobName":           autorest.Encode("path", jobName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2016-03-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsJSON(),
		autorest.AsPatch(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Scheduler/jobCollections/{jobCollectionName}/jobs/{jobName}", pathParameters),
		autorest.WithJSON(job),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{})
}

// PatchSender sends the Patch request. The method will close the
// http.Response Body if it receives an error.
func (client JobsClient) PatchSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req)
}

// PatchResponder handles the response to the Patch request. The method always
// closes the http.Response Body.
func (client JobsClient) PatchResponder(resp *http.Response) (result JobDefinition, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// Run runs a job.
//
// resourceGroupName is the resource group name. jobCollectionName is the job
// collection name. jobName is the job name.
func (client JobsClient) Run(resourceGroupName string, jobCollectionName string, jobName string) (result autorest.Response, err error) {
	req, err := client.RunPreparer(resourceGroupName, jobCollectionName, jobName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "scheduler.JobsClient", "Run", nil, "Failure preparing request")
		return
	}

	resp, err := client.RunSender(req)
	if err != nil {
		result.Response = resp
		err = autorest.NewErrorWithError(err, "scheduler.JobsClient", "Run", resp, "Failure sending request")
		return
	}

	result, err = client.RunResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "scheduler.JobsClient", "Run", resp, "Failure responding to request")
	}

	return
}

// RunPreparer prepares the Run request.
func (client JobsClient) RunPreparer(resourceGroupName string, jobCollectionName string, jobName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"jobCollectionName": autorest.Encode("path", jobCollectionName),
		"jobName":           autorest.Encode("path", jobName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2016-03-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsPost(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Scheduler/jobCollections/{jobCollectionName}/jobs/{jobName}/run", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{})
}

// RunSender sends the Run request. The method will close the
// http.Response Body if it receives an error.
func (client JobsClient) RunSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req)
}

// RunResponder handles the response to the Run request. The method always
// closes the http.Response Body.
func (client JobsClient) RunResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByClosing())
	result.Response = resp
	return
}

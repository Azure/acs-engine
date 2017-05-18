package recoveryservicesbackup

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

// JobDetailsClient is the client for the JobDetails methods of the
// Recoveryservicesbackup service.
type JobDetailsClient struct {
	ManagementClient
}

// NewJobDetailsClient creates an instance of the JobDetailsClient client.
func NewJobDetailsClient(subscriptionID string) JobDetailsClient {
	return NewJobDetailsClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewJobDetailsClientWithBaseURI creates an instance of the JobDetailsClient
// client.
func NewJobDetailsClientWithBaseURI(baseURI string, subscriptionID string) JobDetailsClient {
	return JobDetailsClient{NewWithBaseURI(baseURI, subscriptionID)}
}

// Get gets exteded information associated with the job.
//
// vaultName is the name of the recovery services vault. resourceGroupName is
// the name of the resource group where the recovery services vault is present.
// jobName is name of the job whose details are to be fetched.
func (client JobDetailsClient) Get(vaultName string, resourceGroupName string, jobName string) (result JobResource, err error) {
	req, err := client.GetPreparer(vaultName, resourceGroupName, jobName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "recoveryservicesbackup.JobDetailsClient", "Get", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "recoveryservicesbackup.JobDetailsClient", "Get", resp, "Failure sending request")
		return
	}

	result, err = client.GetResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "recoveryservicesbackup.JobDetailsClient", "Get", resp, "Failure responding to request")
	}

	return
}

// GetPreparer prepares the Get request.
func (client JobDetailsClient) GetPreparer(vaultName string, resourceGroupName string, jobName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"jobName":           autorest.Encode("path", jobName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
		"vaultName":         autorest.Encode("path", vaultName),
	}

	const APIVersion = "2016-12-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/Subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.RecoveryServices/vaults/{vaultName}/backupJobs/{jobName}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{})
}

// GetSender sends the Get request. The method will close the
// http.Response Body if it receives an error.
func (client JobDetailsClient) GetSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req)
}

// GetResponder handles the response to the Get request. The method always
// closes the http.Response Body.
func (client JobDetailsClient) GetResponder(resp *http.Response) (result JobResource, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

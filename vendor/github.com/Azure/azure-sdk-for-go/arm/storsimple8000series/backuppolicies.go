package storsimple8000series

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
// Code generated by Microsoft (R) AutoRest Code Generator 1.1.0.0
// Changes may cause incorrect behavior and will be lost if the code is
// regenerated.

import (
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/validation"
	"net/http"
)

// BackupPoliciesClient is the client for the BackupPolicies methods of the
// Storsimple8000series service.
type BackupPoliciesClient struct {
	ManagementClient
}

// NewBackupPoliciesClient creates an instance of the BackupPoliciesClient
// client.
func NewBackupPoliciesClient(subscriptionID string) BackupPoliciesClient {
	return NewBackupPoliciesClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewBackupPoliciesClientWithBaseURI creates an instance of the
// BackupPoliciesClient client.
func NewBackupPoliciesClientWithBaseURI(baseURI string, subscriptionID string) BackupPoliciesClient {
	return BackupPoliciesClient{NewWithBaseURI(baseURI, subscriptionID)}
}

// BackupNow backup the backup policy now. This method may poll for completion.
// Polling can be canceled by passing the cancel channel argument. The channel
// will be used to cancel polling and any outstanding HTTP requests.
//
// deviceName is the device name backupPolicyName is the backup policy name.
// backupType is the backup Type. This can be cloudSnapshot or localSnapshot.
// resourceGroupName is the resource group name managerName is the manager name
func (client BackupPoliciesClient) BackupNow(deviceName string, backupPolicyName string, backupType string, resourceGroupName string, managerName string, cancel <-chan struct{}) (<-chan autorest.Response, <-chan error) {
	resultChan := make(chan autorest.Response, 1)
	errChan := make(chan error, 1)
	if err := validation.Validate([]validation.Validation{
		{TargetValue: managerName,
			Constraints: []validation.Constraint{{Target: "managerName", Name: validation.MaxLength, Rule: 50, Chain: nil},
				{Target: "managerName", Name: validation.MinLength, Rule: 2, Chain: nil}}}}); err != nil {
		errChan <- validation.NewErrorWithValidationError(err, "storsimple8000series.BackupPoliciesClient", "BackupNow")
		close(errChan)
		close(resultChan)
		return resultChan, errChan
	}

	go func() {
		var err error
		var result autorest.Response
		defer func() {
			if err != nil {
				errChan <- err
			}
			resultChan <- result
			close(resultChan)
			close(errChan)
		}()
		req, err := client.BackupNowPreparer(deviceName, backupPolicyName, backupType, resourceGroupName, managerName, cancel)
		if err != nil {
			err = autorest.NewErrorWithError(err, "storsimple8000series.BackupPoliciesClient", "BackupNow", nil, "Failure preparing request")
			return
		}

		resp, err := client.BackupNowSender(req)
		if err != nil {
			result.Response = resp
			err = autorest.NewErrorWithError(err, "storsimple8000series.BackupPoliciesClient", "BackupNow", resp, "Failure sending request")
			return
		}

		result, err = client.BackupNowResponder(resp)
		if err != nil {
			err = autorest.NewErrorWithError(err, "storsimple8000series.BackupPoliciesClient", "BackupNow", resp, "Failure responding to request")
		}
	}()
	return resultChan, errChan
}

// BackupNowPreparer prepares the BackupNow request.
func (client BackupPoliciesClient) BackupNowPreparer(deviceName string, backupPolicyName string, backupType string, resourceGroupName string, managerName string, cancel <-chan struct{}) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"backupPolicyName":  backupPolicyName,
		"deviceName":        deviceName,
		"managerName":       managerName,
		"resourceGroupName": resourceGroupName,
		"subscriptionId":    client.SubscriptionID,
	}

	const APIVersion = "2017-06-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
		"backupType":  backupType,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsPost(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.StorSimple/managers/{managerName}/devices/{deviceName}/backupPolicies/{backupPolicyName}/backup", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{Cancel: cancel})
}

// BackupNowSender sends the BackupNow request. The method will close the
// http.Response Body if it receives an error.
func (client BackupPoliciesClient) BackupNowSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client,
		req,
		azure.DoPollForAsynchronous(client.PollingDelay))
}

// BackupNowResponder handles the response to the BackupNow request. The method always
// closes the http.Response Body.
func (client BackupPoliciesClient) BackupNowResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusAccepted, http.StatusNoContent),
		autorest.ByClosing())
	result.Response = resp
	return
}

// CreateOrUpdate creates or updates the backup policy. This method may poll
// for completion. Polling can be canceled by passing the cancel channel
// argument. The channel will be used to cancel polling and any outstanding
// HTTP requests.
//
// deviceName is the device name backupPolicyName is the name of the backup
// policy to be created/updated. parameters is the backup policy.
// resourceGroupName is the resource group name managerName is the manager name
func (client BackupPoliciesClient) CreateOrUpdate(deviceName string, backupPolicyName string, parameters BackupPolicy, resourceGroupName string, managerName string, cancel <-chan struct{}) (<-chan BackupPolicy, <-chan error) {
	resultChan := make(chan BackupPolicy, 1)
	errChan := make(chan error, 1)
	if err := validation.Validate([]validation.Validation{
		{TargetValue: parameters,
			Constraints: []validation.Constraint{{Target: "parameters.BackupPolicyProperties", Name: validation.Null, Rule: true,
				Chain: []validation.Constraint{{Target: "parameters.BackupPolicyProperties.VolumeIds", Name: validation.Null, Rule: true, Chain: nil}}}}},
		{TargetValue: managerName,
			Constraints: []validation.Constraint{{Target: "managerName", Name: validation.MaxLength, Rule: 50, Chain: nil},
				{Target: "managerName", Name: validation.MinLength, Rule: 2, Chain: nil}}}}); err != nil {
		errChan <- validation.NewErrorWithValidationError(err, "storsimple8000series.BackupPoliciesClient", "CreateOrUpdate")
		close(errChan)
		close(resultChan)
		return resultChan, errChan
	}

	go func() {
		var err error
		var result BackupPolicy
		defer func() {
			if err != nil {
				errChan <- err
			}
			resultChan <- result
			close(resultChan)
			close(errChan)
		}()
		req, err := client.CreateOrUpdatePreparer(deviceName, backupPolicyName, parameters, resourceGroupName, managerName, cancel)
		if err != nil {
			err = autorest.NewErrorWithError(err, "storsimple8000series.BackupPoliciesClient", "CreateOrUpdate", nil, "Failure preparing request")
			return
		}

		resp, err := client.CreateOrUpdateSender(req)
		if err != nil {
			result.Response = autorest.Response{Response: resp}
			err = autorest.NewErrorWithError(err, "storsimple8000series.BackupPoliciesClient", "CreateOrUpdate", resp, "Failure sending request")
			return
		}

		result, err = client.CreateOrUpdateResponder(resp)
		if err != nil {
			err = autorest.NewErrorWithError(err, "storsimple8000series.BackupPoliciesClient", "CreateOrUpdate", resp, "Failure responding to request")
		}
	}()
	return resultChan, errChan
}

// CreateOrUpdatePreparer prepares the CreateOrUpdate request.
func (client BackupPoliciesClient) CreateOrUpdatePreparer(deviceName string, backupPolicyName string, parameters BackupPolicy, resourceGroupName string, managerName string, cancel <-chan struct{}) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"backupPolicyName":  backupPolicyName,
		"deviceName":        deviceName,
		"managerName":       managerName,
		"resourceGroupName": resourceGroupName,
		"subscriptionId":    client.SubscriptionID,
	}

	const APIVersion = "2017-06-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsJSON(),
		autorest.AsPut(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.StorSimple/managers/{managerName}/devices/{deviceName}/backupPolicies/{backupPolicyName}", pathParameters),
		autorest.WithJSON(parameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{Cancel: cancel})
}

// CreateOrUpdateSender sends the CreateOrUpdate request. The method will close the
// http.Response Body if it receives an error.
func (client BackupPoliciesClient) CreateOrUpdateSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client,
		req,
		azure.DoPollForAsynchronous(client.PollingDelay))
}

// CreateOrUpdateResponder handles the response to the CreateOrUpdate request. The method always
// closes the http.Response Body.
func (client BackupPoliciesClient) CreateOrUpdateResponder(resp *http.Response) (result BackupPolicy, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusAccepted),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// Delete deletes the backup policy. This method may poll for completion.
// Polling can be canceled by passing the cancel channel argument. The channel
// will be used to cancel polling and any outstanding HTTP requests.
//
// deviceName is the device name backupPolicyName is the name of the backup
// policy. resourceGroupName is the resource group name managerName is the
// manager name
func (client BackupPoliciesClient) Delete(deviceName string, backupPolicyName string, resourceGroupName string, managerName string, cancel <-chan struct{}) (<-chan autorest.Response, <-chan error) {
	resultChan := make(chan autorest.Response, 1)
	errChan := make(chan error, 1)
	if err := validation.Validate([]validation.Validation{
		{TargetValue: managerName,
			Constraints: []validation.Constraint{{Target: "managerName", Name: validation.MaxLength, Rule: 50, Chain: nil},
				{Target: "managerName", Name: validation.MinLength, Rule: 2, Chain: nil}}}}); err != nil {
		errChan <- validation.NewErrorWithValidationError(err, "storsimple8000series.BackupPoliciesClient", "Delete")
		close(errChan)
		close(resultChan)
		return resultChan, errChan
	}

	go func() {
		var err error
		var result autorest.Response
		defer func() {
			if err != nil {
				errChan <- err
			}
			resultChan <- result
			close(resultChan)
			close(errChan)
		}()
		req, err := client.DeletePreparer(deviceName, backupPolicyName, resourceGroupName, managerName, cancel)
		if err != nil {
			err = autorest.NewErrorWithError(err, "storsimple8000series.BackupPoliciesClient", "Delete", nil, "Failure preparing request")
			return
		}

		resp, err := client.DeleteSender(req)
		if err != nil {
			result.Response = resp
			err = autorest.NewErrorWithError(err, "storsimple8000series.BackupPoliciesClient", "Delete", resp, "Failure sending request")
			return
		}

		result, err = client.DeleteResponder(resp)
		if err != nil {
			err = autorest.NewErrorWithError(err, "storsimple8000series.BackupPoliciesClient", "Delete", resp, "Failure responding to request")
		}
	}()
	return resultChan, errChan
}

// DeletePreparer prepares the Delete request.
func (client BackupPoliciesClient) DeletePreparer(deviceName string, backupPolicyName string, resourceGroupName string, managerName string, cancel <-chan struct{}) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"backupPolicyName":  backupPolicyName,
		"deviceName":        deviceName,
		"managerName":       managerName,
		"resourceGroupName": resourceGroupName,
		"subscriptionId":    client.SubscriptionID,
	}

	const APIVersion = "2017-06-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsDelete(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.StorSimple/managers/{managerName}/devices/{deviceName}/backupPolicies/{backupPolicyName}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{Cancel: cancel})
}

// DeleteSender sends the Delete request. The method will close the
// http.Response Body if it receives an error.
func (client BackupPoliciesClient) DeleteSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client,
		req,
		azure.DoPollForAsynchronous(client.PollingDelay))
}

// DeleteResponder handles the response to the Delete request. The method always
// closes the http.Response Body.
func (client BackupPoliciesClient) DeleteResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusAccepted, http.StatusNoContent),
		autorest.ByClosing())
	result.Response = resp
	return
}

// Get gets the properties of the specified backup policy name.
//
// deviceName is the device name backupPolicyName is the name of backup policy
// to be fetched. resourceGroupName is the resource group name managerName is
// the manager name
func (client BackupPoliciesClient) Get(deviceName string, backupPolicyName string, resourceGroupName string, managerName string) (result BackupPolicy, err error) {
	if err := validation.Validate([]validation.Validation{
		{TargetValue: managerName,
			Constraints: []validation.Constraint{{Target: "managerName", Name: validation.MaxLength, Rule: 50, Chain: nil},
				{Target: "managerName", Name: validation.MinLength, Rule: 2, Chain: nil}}}}); err != nil {
		return result, validation.NewErrorWithValidationError(err, "storsimple8000series.BackupPoliciesClient", "Get")
	}

	req, err := client.GetPreparer(deviceName, backupPolicyName, resourceGroupName, managerName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "storsimple8000series.BackupPoliciesClient", "Get", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "storsimple8000series.BackupPoliciesClient", "Get", resp, "Failure sending request")
		return
	}

	result, err = client.GetResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "storsimple8000series.BackupPoliciesClient", "Get", resp, "Failure responding to request")
	}

	return
}

// GetPreparer prepares the Get request.
func (client BackupPoliciesClient) GetPreparer(deviceName string, backupPolicyName string, resourceGroupName string, managerName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"backupPolicyName":  backupPolicyName,
		"deviceName":        deviceName,
		"managerName":       managerName,
		"resourceGroupName": resourceGroupName,
		"subscriptionId":    client.SubscriptionID,
	}

	const APIVersion = "2017-06-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.StorSimple/managers/{managerName}/devices/{deviceName}/backupPolicies/{backupPolicyName}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{})
}

// GetSender sends the Get request. The method will close the
// http.Response Body if it receives an error.
func (client BackupPoliciesClient) GetSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req)
}

// GetResponder handles the response to the Get request. The method always
// closes the http.Response Body.
func (client BackupPoliciesClient) GetResponder(resp *http.Response) (result BackupPolicy, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// ListByDevice gets all the backup policies in a device.
//
// deviceName is the device name resourceGroupName is the resource group name
// managerName is the manager name
func (client BackupPoliciesClient) ListByDevice(deviceName string, resourceGroupName string, managerName string) (result BackupPolicyList, err error) {
	if err := validation.Validate([]validation.Validation{
		{TargetValue: managerName,
			Constraints: []validation.Constraint{{Target: "managerName", Name: validation.MaxLength, Rule: 50, Chain: nil},
				{Target: "managerName", Name: validation.MinLength, Rule: 2, Chain: nil}}}}); err != nil {
		return result, validation.NewErrorWithValidationError(err, "storsimple8000series.BackupPoliciesClient", "ListByDevice")
	}

	req, err := client.ListByDevicePreparer(deviceName, resourceGroupName, managerName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "storsimple8000series.BackupPoliciesClient", "ListByDevice", nil, "Failure preparing request")
		return
	}

	resp, err := client.ListByDeviceSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "storsimple8000series.BackupPoliciesClient", "ListByDevice", resp, "Failure sending request")
		return
	}

	result, err = client.ListByDeviceResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "storsimple8000series.BackupPoliciesClient", "ListByDevice", resp, "Failure responding to request")
	}

	return
}

// ListByDevicePreparer prepares the ListByDevice request.
func (client BackupPoliciesClient) ListByDevicePreparer(deviceName string, resourceGroupName string, managerName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"deviceName":        deviceName,
		"managerName":       managerName,
		"resourceGroupName": resourceGroupName,
		"subscriptionId":    client.SubscriptionID,
	}

	const APIVersion = "2017-06-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.StorSimple/managers/{managerName}/devices/{deviceName}/backupPolicies", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{})
}

// ListByDeviceSender sends the ListByDevice request. The method will close the
// http.Response Body if it receives an error.
func (client BackupPoliciesClient) ListByDeviceSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req)
}

// ListByDeviceResponder handles the response to the ListByDevice request. The method always
// closes the http.Response Body.
func (client BackupPoliciesClient) ListByDeviceResponder(resp *http.Response) (result BackupPolicyList, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

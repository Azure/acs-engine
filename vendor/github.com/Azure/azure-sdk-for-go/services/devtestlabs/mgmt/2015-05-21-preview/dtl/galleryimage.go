package dtl

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
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

import (
	"context"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"net/http"
)

// GalleryImageClient is the azure DevTest Labs REST API version 2015-05-21-preview.
type GalleryImageClient struct {
	BaseClient
}

// NewGalleryImageClient creates an instance of the GalleryImageClient client.
func NewGalleryImageClient(subscriptionID string) GalleryImageClient {
	return NewGalleryImageClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewGalleryImageClientWithBaseURI creates an instance of the GalleryImageClient client.
func NewGalleryImageClientWithBaseURI(baseURI string, subscriptionID string) GalleryImageClient {
	return GalleryImageClient{NewWithBaseURI(baseURI, subscriptionID)}
}

// List list gallery images.
//
// resourceGroupName is the name of the resource group. labName is the name of the lab. filter is the filter to apply
// on the operation.
func (client GalleryImageClient) List(ctx context.Context, resourceGroupName string, labName string, filter string, top *int32, orderBy string) (result ResponseWithContinuationGalleryImagePage, err error) {
	result.fn = client.listNextResults
	req, err := client.ListPreparer(ctx, resourceGroupName, labName, filter, top, orderBy)
	if err != nil {
		err = autorest.NewErrorWithError(err, "dtl.GalleryImageClient", "List", nil, "Failure preparing request")
		return
	}

	resp, err := client.ListSender(req)
	if err != nil {
		result.rwcgi.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "dtl.GalleryImageClient", "List", resp, "Failure sending request")
		return
	}

	result.rwcgi, err = client.ListResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "dtl.GalleryImageClient", "List", resp, "Failure responding to request")
	}

	return
}

// ListPreparer prepares the List request.
func (client GalleryImageClient) ListPreparer(ctx context.Context, resourceGroupName string, labName string, filter string, top *int32, orderBy string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"labName":           autorest.Encode("path", labName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2015-05-21-preview"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}
	if len(filter) > 0 {
		queryParameters["$filter"] = autorest.Encode("query", filter)
	}
	if top != nil {
		queryParameters["$top"] = autorest.Encode("query", *top)
	}
	if len(orderBy) > 0 {
		queryParameters["$orderBy"] = autorest.Encode("query", orderBy)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DevTestLab/labs/{labName}/galleryimages", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// ListSender sends the List request. The method will close the
// http.Response Body if it receives an error.
func (client GalleryImageClient) ListSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		azure.DoRetryWithRegistration(client.Client))
}

// ListResponder handles the response to the List request. The method always
// closes the http.Response Body.
func (client GalleryImageClient) ListResponder(resp *http.Response) (result ResponseWithContinuationGalleryImage, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// listNextResults retrieves the next set of results, if any.
func (client GalleryImageClient) listNextResults(lastResults ResponseWithContinuationGalleryImage) (result ResponseWithContinuationGalleryImage, err error) {
	req, err := lastResults.responseWithContinuationGalleryImagePreparer()
	if err != nil {
		return result, autorest.NewErrorWithError(err, "dtl.GalleryImageClient", "listNextResults", nil, "Failure preparing next results request")
	}
	if req == nil {
		return
	}
	resp, err := client.ListSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "dtl.GalleryImageClient", "listNextResults", resp, "Failure sending next results request")
	}
	result, err = client.ListResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "dtl.GalleryImageClient", "listNextResults", resp, "Failure responding to next results request")
	}
	return
}

// ListComplete enumerates all values, automatically crossing page boundaries as required.
func (client GalleryImageClient) ListComplete(ctx context.Context, resourceGroupName string, labName string, filter string, top *int32, orderBy string) (result ResponseWithContinuationGalleryImageIterator, err error) {
	result.page, err = client.List(ctx, resourceGroupName, labName, filter, top, orderBy)
	return
}

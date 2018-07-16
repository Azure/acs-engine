package operationalinsights

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
	"encoding/json"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/Azure/go-autorest/autorest/to"
	"net/http"
)

// SearchSortEnum enumerates the values for search sort enum.
type SearchSortEnum string

const (
	// Asc ...
	Asc SearchSortEnum = "asc"
	// Desc ...
	Desc SearchSortEnum = "desc"
)

// StorageInsightState enumerates the values for storage insight state.
type StorageInsightState string

const (
	// ERROR ...
	ERROR StorageInsightState = "ERROR"
	// OK ...
	OK StorageInsightState = "OK"
)

// CoreSummary the core summary of a search.
type CoreSummary struct {
	// Status - The status of a core summary.
	Status *string `json:"Status,omitempty"`
	// NumberOfDocuments - The number of documents of a core summary.
	NumberOfDocuments *int64 `json:"NumberOfDocuments,omitempty"`
}

// LinkTarget metadata for a workspace that isn't linked to an Azure subscription.
type LinkTarget struct {
	// CustomerID - The GUID that uniquely identifies the workspace.
	CustomerID *string `json:"customerId,omitempty"`
	// DisplayName - The display name of the workspace.
	DisplayName *string `json:"accountName,omitempty"`
	// WorkspaceName - The DNS valid workspace name.
	WorkspaceName *string `json:"workspaceName,omitempty"`
	// Location - The location of the workspace.
	Location *string `json:"location,omitempty"`
}

// ListLinkTarget ...
type ListLinkTarget struct {
	autorest.Response `json:"-"`
	Value             *[]LinkTarget `json:"value,omitempty"`
}

// ProxyResource common properties of proxy resource.
type ProxyResource struct {
	// ID - Resource ID.
	ID *string `json:"id,omitempty"`
	// Name - Resource name.
	Name *string `json:"name,omitempty"`
	// Type - Resource type.
	Type *string `json:"type,omitempty"`
	// Tags - Resource tags
	Tags *map[string]*string `json:"tags,omitempty"`
}

// Resource the resource definition.
type Resource struct {
	// ID - Resource Id
	ID *string `json:"id,omitempty"`
	// Name - Resource name
	Name *string `json:"name,omitempty"`
	// Type - Resource type
	Type *string `json:"type,omitempty"`
	// Location - Resource location
	Location *string `json:"location,omitempty"`
	// Tags - Resource tags
	Tags *map[string]*string `json:"tags,omitempty"`
}

// SavedSearch value object for saved search results.
type SavedSearch struct {
	autorest.Response `json:"-"`
	// ID - The id of the saved search.
	ID *string `json:"id,omitempty"`
	// Etag - The etag of the saved search.
	Etag *string `json:"etag,omitempty"`
	// SavedSearchProperties - Gets or sets properties of the saved search.
	*SavedSearchProperties `json:"properties,omitempty"`
}

// UnmarshalJSON is the custom unmarshaler for SavedSearch struct.
func (ss *SavedSearch) UnmarshalJSON(body []byte) error {
	var m map[string]*json.RawMessage
	err := json.Unmarshal(body, &m)
	if err != nil {
		return err
	}
	var v *json.RawMessage

	v = m["id"]
	if v != nil {
		var ID string
		err = json.Unmarshal(*m["id"], &ID)
		if err != nil {
			return err
		}
		ss.ID = &ID
	}

	v = m["etag"]
	if v != nil {
		var etag string
		err = json.Unmarshal(*m["etag"], &etag)
		if err != nil {
			return err
		}
		ss.Etag = &etag
	}

	v = m["properties"]
	if v != nil {
		var properties SavedSearchProperties
		err = json.Unmarshal(*m["properties"], &properties)
		if err != nil {
			return err
		}
		ss.SavedSearchProperties = &properties
	}

	return nil
}

// SavedSearchesListResult the saved search operation response.
type SavedSearchesListResult struct {
	autorest.Response `json:"-"`
	// Metadata - The metadata from search results.
	Metadata *SearchMetadata `json:"__metadata,omitempty"`
	// Value - The array of result values.
	Value *[]SavedSearch `json:"value,omitempty"`
}

// SavedSearchProperties value object for saved search results.
type SavedSearchProperties struct {
	// Category - The category of the saved search. This helps the user to find a saved search faster.
	Category *string `json:"Category,omitempty"`
	// DisplayName - Saved search display name.
	DisplayName *string `json:"DisplayName,omitempty"`
	// Query - The query expression for the saved search. Please see https://docs.microsoft.com/en-us/azure/log-analytics/log-analytics-search-reference for reference.
	Query *string `json:"Query,omitempty"`
	// Version - The version number of the query lanuage. Only verion 1 is allowed here.
	Version *int64 `json:"Version,omitempty"`
	// Tags - The tags attached to the saved search.
	Tags *[]Tag `json:"Tags,omitempty"`
}

// SearchError details for a search error.
type SearchError struct {
	// Type - The error type.
	Type *string `json:"type,omitempty"`
	// Message - The error message.
	Message *string `json:"message,omitempty"`
}

// SearchGetSchemaResponse the get schema operation response.
type SearchGetSchemaResponse struct {
	autorest.Response `json:"-"`
	// Metadata - The metadata from search results.
	Metadata *SearchMetadata `json:"__metadata,omitempty"`
	// Value - The array of result values.
	Value *[]SearchSchemaValue `json:"value,omitempty"`
}

// SearchHighlight highlight details.
type SearchHighlight struct {
	// Pre - The string that is put before a matched result.
	Pre *string `json:"pre,omitempty"`
	// Post - The string that is put after a matched result.
	Post *string `json:"post,omitempty"`
}

// SearchMetadata metadata for search results.
type SearchMetadata struct {
	// SearchID - The request id of the search.
	SearchID *string `json:"RequestId,omitempty"`
	// ResultType - The search result type.
	ResultType *string `json:"resultType,omitempty"`
	// Total - The total number of search results.
	Total *int64 `json:"total,omitempty"`
	// Top - The number of top search results.
	Top *int64 `json:"top,omitempty"`
	// ID - The id of the search results request.
	ID *string `json:"id,omitempty"`
	// CoreSummaries - The core summaries.
	CoreSummaries *[]CoreSummary `json:"CoreSummaries,omitempty"`
	// Status - The status of the search results.
	Status *string `json:"Status,omitempty"`
	// StartTime - The start time for the search.
	StartTime *date.Time `json:"StartTime,omitempty"`
	// LastUpdated - The time of last update.
	LastUpdated *date.Time `json:"LastUpdated,omitempty"`
	// ETag - The ETag of the search results.
	ETag *string `json:"ETag,omitempty"`
	// Sort - How the results are sorted.
	Sort *[]SearchSort `json:"sort,omitempty"`
	// RequestTime - The request time.
	RequestTime *int64 `json:"requestTime,omitempty"`
	// AggregatedValueField - The aggregated value field.
	AggregatedValueField *string `json:"aggregatedValueField,omitempty"`
	// AggregatedGroupingFields - The aggregated grouping fields.
	AggregatedGroupingFields *string `json:"aggregatedGroupingFields,omitempty"`
	// Sum - The sum of all aggregates returned in the result set.
	Sum *int64 `json:"sum,omitempty"`
	// Max - The max of all aggregates returned in the result set.
	Max *int64 `json:"max,omitempty"`
	// Schema - The schema.
	Schema *SearchMetadataSchema `json:"schema,omitempty"`
}

// SearchMetadataSchema schema metadata for search.
type SearchMetadataSchema struct {
	// Name - The name of the metadata schema.
	Name *string `json:"name,omitempty"`
	// Version - The version of the metadata schema.
	Version *int32 `json:"version,omitempty"`
}

// SearchParameters parameters specifying the search query and range.
type SearchParameters struct {
	// Top - The number to get from the top.
	Top *int64 `json:"top,omitempty"`
	// Highlight - The highlight that looks for all occurences of a string.
	Highlight *SearchHighlight `json:"highlight,omitempty"`
	// Query - The query to search.
	Query *string `json:"query,omitempty"`
	// Start - The start date filter, so the only query results returned are after this date.
	Start *date.Time `json:"start,omitempty"`
	// End - The end date filter, so the only query results returned are before this date.
	End *date.Time `json:"end,omitempty"`
}

// SearchResultsResponse the get search result operation response.
type SearchResultsResponse struct {
	autorest.Response `json:"-"`
	// ID - The id of the search, which includes the full url.
	ID *string `json:"id,omitempty"`
	// Metadata - The metadata from search results.
	Metadata *SearchMetadata `json:"__metadata,omitempty"`
	// Value - The array of result values.
	Value *[]map[string]interface{} `json:"value,omitempty"`
	// Error - The error.
	Error *SearchError `json:"error,omitempty"`
}

// SearchSchemaValue value object for schema results.
type SearchSchemaValue struct {
	// Name - The name of the schema.
	Name *string `json:"name,omitempty"`
	// DisplayName - The display name of the schema.
	DisplayName *string `json:"displayName,omitempty"`
	// Type - The type.
	Type *string `json:"type,omitempty"`
	// Indexed - The boolean that indicates the field is searchable as free text.
	Indexed *bool `json:"indexed,omitempty"`
	// Stored - The boolean that indicates whether or not the field is stored.
	Stored *bool `json:"stored,omitempty"`
	// Facet - The boolean that indicates whether or not the field is a facet.
	Facet *bool `json:"facet,omitempty"`
	// OwnerType - The array of workflows containing the field.
	OwnerType *[]string `json:"ownerType,omitempty"`
}

// SearchSort the sort parameters for search.
type SearchSort struct {
	// Name - The name of the field the search query is sorted on.
	Name *string `json:"name,omitempty"`
	// Order - The sort order of the search. Possible values include: 'Asc', 'Desc'
	Order SearchSortEnum `json:"order,omitempty"`
}

// StorageAccount describes a storage account connection.
type StorageAccount struct {
	// ID - The Azure Resource Manager ID of the storage account resource.
	ID *string `json:"id,omitempty"`
	// Key - The storage account key.
	Key *string `json:"key,omitempty"`
}

// StorageInsight the top level storage insight resource container.
type StorageInsight struct {
	autorest.Response `json:"-"`
	// ID - Resource ID.
	ID *string `json:"id,omitempty"`
	// Name - Resource name.
	Name *string `json:"name,omitempty"`
	// Type - Resource type.
	Type *string `json:"type,omitempty"`
	// Tags - Resource tags
	Tags *map[string]*string `json:"tags,omitempty"`
	// StorageInsightProperties - Storage insight properties.
	*StorageInsightProperties `json:"properties,omitempty"`
	// ETag - The ETag of the storage insight.
	ETag *string `json:"eTag,omitempty"`
}

// UnmarshalJSON is the custom unmarshaler for StorageInsight struct.
func (si *StorageInsight) UnmarshalJSON(body []byte) error {
	var m map[string]*json.RawMessage
	err := json.Unmarshal(body, &m)
	if err != nil {
		return err
	}
	var v *json.RawMessage

	v = m["properties"]
	if v != nil {
		var properties StorageInsightProperties
		err = json.Unmarshal(*m["properties"], &properties)
		if err != nil {
			return err
		}
		si.StorageInsightProperties = &properties
	}

	v = m["eTag"]
	if v != nil {
		var eTag string
		err = json.Unmarshal(*m["eTag"], &eTag)
		if err != nil {
			return err
		}
		si.ETag = &eTag
	}

	v = m["id"]
	if v != nil {
		var ID string
		err = json.Unmarshal(*m["id"], &ID)
		if err != nil {
			return err
		}
		si.ID = &ID
	}

	v = m["name"]
	if v != nil {
		var name string
		err = json.Unmarshal(*m["name"], &name)
		if err != nil {
			return err
		}
		si.Name = &name
	}

	v = m["type"]
	if v != nil {
		var typeVar string
		err = json.Unmarshal(*m["type"], &typeVar)
		if err != nil {
			return err
		}
		si.Type = &typeVar
	}

	v = m["tags"]
	if v != nil {
		var tags map[string]*string
		err = json.Unmarshal(*m["tags"], &tags)
		if err != nil {
			return err
		}
		si.Tags = &tags
	}

	return nil
}

// StorageInsightListResult the list storage insights operation response.
type StorageInsightListResult struct {
	autorest.Response `json:"-"`
	// Value - Gets or sets a list of storage insight instances.
	Value *[]StorageInsight `json:"value,omitempty"`
	// OdataNextLink - The link (url) to the next page of results.
	OdataNextLink *string `json:"@odata.nextLink,omitempty"`
}

// StorageInsightListResultIterator provides access to a complete listing of StorageInsight values.
type StorageInsightListResultIterator struct {
	i    int
	page StorageInsightListResultPage
}

// Next advances to the next value.  If there was an error making
// the request the iterator does not advance and the error is returned.
func (iter *StorageInsightListResultIterator) Next() error {
	iter.i++
	if iter.i < len(iter.page.Values()) {
		return nil
	}
	err := iter.page.Next()
	if err != nil {
		iter.i--
		return err
	}
	iter.i = 0
	return nil
}

// NotDone returns true if the enumeration should be started or is not yet complete.
func (iter StorageInsightListResultIterator) NotDone() bool {
	return iter.page.NotDone() && iter.i < len(iter.page.Values())
}

// Response returns the raw server response from the last page request.
func (iter StorageInsightListResultIterator) Response() StorageInsightListResult {
	return iter.page.Response()
}

// Value returns the current value or a zero-initialized value if the
// iterator has advanced beyond the end of the collection.
func (iter StorageInsightListResultIterator) Value() StorageInsight {
	if !iter.page.NotDone() {
		return StorageInsight{}
	}
	return iter.page.Values()[iter.i]
}

// IsEmpty returns true if the ListResult contains no values.
func (silr StorageInsightListResult) IsEmpty() bool {
	return silr.Value == nil || len(*silr.Value) == 0
}

// storageInsightListResultPreparer prepares a request to retrieve the next set of results.
// It returns nil if no more results exist.
func (silr StorageInsightListResult) storageInsightListResultPreparer() (*http.Request, error) {
	if silr.OdataNextLink == nil || len(to.String(silr.OdataNextLink)) < 1 {
		return nil, nil
	}
	return autorest.Prepare(&http.Request{},
		autorest.AsJSON(),
		autorest.AsGet(),
		autorest.WithBaseURL(to.String(silr.OdataNextLink)))
}

// StorageInsightListResultPage contains a page of StorageInsight values.
type StorageInsightListResultPage struct {
	fn   func(StorageInsightListResult) (StorageInsightListResult, error)
	silr StorageInsightListResult
}

// Next advances to the next page of values.  If there was an error making
// the request the page does not advance and the error is returned.
func (page *StorageInsightListResultPage) Next() error {
	next, err := page.fn(page.silr)
	if err != nil {
		return err
	}
	page.silr = next
	return nil
}

// NotDone returns true if the page enumeration should be started or is not yet complete.
func (page StorageInsightListResultPage) NotDone() bool {
	return !page.silr.IsEmpty()
}

// Response returns the raw server response from the last page request.
func (page StorageInsightListResultPage) Response() StorageInsightListResult {
	return page.silr
}

// Values returns the slice of values for the current page or nil if there are no values.
func (page StorageInsightListResultPage) Values() []StorageInsight {
	if page.silr.IsEmpty() {
		return nil
	}
	return *page.silr.Value
}

// StorageInsightProperties storage insight properties.
type StorageInsightProperties struct {
	// Containers - The names of the blob containers that the workspace should read
	Containers *[]string `json:"containers,omitempty"`
	// Tables - The names of the Azure tables that the workspace should read
	Tables *[]string `json:"tables,omitempty"`
	// StorageAccount - The storage account connection details
	StorageAccount *StorageAccount `json:"storageAccount,omitempty"`
	// Status - The status of the storage insight
	Status *StorageInsightStatus `json:"status,omitempty"`
}

// StorageInsightStatus the status of the storage insight.
type StorageInsightStatus struct {
	// State - The state of the storage insight connection to the workspace. Possible values include: 'OK', 'ERROR'
	State StorageInsightState `json:"state,omitempty"`
	// Description - Description of the state of the storage insight.
	Description *string `json:"description,omitempty"`
}

// Tag a tag of a saved search.
type Tag struct {
	// Name - The tag name.
	Name *string `json:"Name,omitempty"`
	// Value - The tag value.
	Value *string `json:"Value,omitempty"`
}

// WorkspacesGetSearchResultsFuture an abstraction for monitoring and retrieving the results of a long-running
// operation.
type WorkspacesGetSearchResultsFuture struct {
	azure.Future
	req *http.Request
}

// Result returns the result of the asynchronous operation.
// If the operation has not completed it will return an error.
func (future WorkspacesGetSearchResultsFuture) Result(client WorkspacesClient) (srr SearchResultsResponse, err error) {
	var done bool
	done, err = future.Done(client)
	if err != nil {
		return
	}
	if !done {
		return srr, autorest.NewError("operationalinsights.WorkspacesGetSearchResultsFuture", "Result", "asynchronous operation has not completed")
	}
	if future.PollingMethod() == azure.PollingLocation {
		srr, err = client.GetSearchResultsResponder(future.Response())
		return
	}
	var resp *http.Response
	resp, err = autorest.SendWithSender(client, autorest.ChangeToGet(future.req),
		autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
	if err != nil {
		return
	}
	srr, err = client.GetSearchResultsResponder(resp)
	return
}

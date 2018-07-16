// +build go1.9

// Copyright 2017 Microsoft Corporation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// This code was auto-generated by:
// github.com/Azure/azure-sdk-for-go/tools/profileBuilder
// commit ID: 2014fbbf031942474ad27a5a66dffaed5347f3fb

package notificationhubs

import original "github.com/Azure/azure-sdk-for-go/services/notificationhubs/mgmt/2017-04-01/notificationhubs"

type NameClient = original.NameClient
type NamespacesClient = original.NamespacesClient
type Client = original.Client

const (
	DefaultBaseURI = original.DefaultBaseURI
)

type BaseClient = original.BaseClient
type HubsClient = original.HubsClient
type AccessRights = original.AccessRights

const (
	Listen AccessRights = original.Listen
	Manage AccessRights = original.Manage
	Send   AccessRights = original.Send
)

type NamespaceType = original.NamespaceType

const (
	Messaging       NamespaceType = original.Messaging
	NotificationHub NamespaceType = original.NotificationHub
)

type SkuName = original.SkuName

const (
	Basic    SkuName = original.Basic
	Free     SkuName = original.Free
	Standard SkuName = original.Standard
)

type AdmCredential = original.AdmCredential
type AdmCredentialProperties = original.AdmCredentialProperties
type ApnsCredential = original.ApnsCredential
type ApnsCredentialProperties = original.ApnsCredentialProperties
type BaiduCredential = original.BaiduCredential
type BaiduCredentialProperties = original.BaiduCredentialProperties
type CheckAvailabilityParameters = original.CheckAvailabilityParameters
type CheckAvailabilityResult = original.CheckAvailabilityResult
type CheckNameAvailabilityRequestParameters = original.CheckNameAvailabilityRequestParameters
type CheckNameAvailabilityResponse = original.CheckNameAvailabilityResponse
type CreateOrUpdateParameters = original.CreateOrUpdateParameters
type GcmCredential = original.GcmCredential
type GcmCredentialProperties = original.GcmCredentialProperties
type ListResult = original.ListResult
type ListResultIterator = original.ListResultIterator
type ListResultPage = original.ListResultPage
type MpnsCredential = original.MpnsCredential
type MpnsCredentialProperties = original.MpnsCredentialProperties
type NamespaceCreateOrUpdateParameters = original.NamespaceCreateOrUpdateParameters
type NamespaceListResult = original.NamespaceListResult
type NamespaceListResultIterator = original.NamespaceListResultIterator
type NamespaceListResultPage = original.NamespaceListResultPage
type NamespacePatchParameters = original.NamespacePatchParameters
type NamespaceProperties = original.NamespaceProperties
type NamespaceResource = original.NamespaceResource
type NamespacesDeleteFuture = original.NamespacesDeleteFuture
type PnsCredentialsProperties = original.PnsCredentialsProperties
type PnsCredentialsResource = original.PnsCredentialsResource
type PolicykeyResource = original.PolicykeyResource
type Properties = original.Properties
type Resource = original.Resource
type ResourceListKeys = original.ResourceListKeys
type ResourceType = original.ResourceType
type SharedAccessAuthorizationRuleCreateOrUpdateParameters = original.SharedAccessAuthorizationRuleCreateOrUpdateParameters
type SharedAccessAuthorizationRuleListResult = original.SharedAccessAuthorizationRuleListResult
type SharedAccessAuthorizationRuleListResultIterator = original.SharedAccessAuthorizationRuleListResultIterator
type SharedAccessAuthorizationRuleListResultPage = original.SharedAccessAuthorizationRuleListResultPage
type SharedAccessAuthorizationRuleProperties = original.SharedAccessAuthorizationRuleProperties
type SharedAccessAuthorizationRuleResource = original.SharedAccessAuthorizationRuleResource
type Sku = original.Sku
type SubResource = original.SubResource
type WnsCredential = original.WnsCredential
type WnsCredentialProperties = original.WnsCredentialProperties

func NewNamespacesClient(subscriptionID string) NamespacesClient {
	return original.NewNamespacesClient(subscriptionID)
}
func NewNamespacesClientWithBaseURI(baseURI string, subscriptionID string) NamespacesClient {
	return original.NewNamespacesClientWithBaseURI(baseURI, subscriptionID)
}
func NewClient(subscriptionID string) Client {
	return original.NewClient(subscriptionID)
}
func NewClientWithBaseURI(baseURI string, subscriptionID string) Client {
	return original.NewClientWithBaseURI(baseURI, subscriptionID)
}
func UserAgent() string {
	return original.UserAgent() + " profiles/preview"
}
func Version() string {
	return original.Version()
}
func New(subscriptionID string) BaseClient {
	return original.New(subscriptionID)
}
func NewWithBaseURI(baseURI string, subscriptionID string) BaseClient {
	return original.NewWithBaseURI(baseURI, subscriptionID)
}
func NewHubsClient(subscriptionID string) HubsClient {
	return original.NewHubsClient(subscriptionID)
}
func NewHubsClientWithBaseURI(baseURI string, subscriptionID string) HubsClient {
	return original.NewHubsClientWithBaseURI(baseURI, subscriptionID)
}
func NewNameClient(subscriptionID string) NameClient {
	return original.NewNameClient(subscriptionID)
}
func NewNameClientWithBaseURI(baseURI string, subscriptionID string) NameClient {
	return original.NewNameClientWithBaseURI(baseURI, subscriptionID)
}

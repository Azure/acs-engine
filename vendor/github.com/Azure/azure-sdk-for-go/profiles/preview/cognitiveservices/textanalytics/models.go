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

package textanalytics

import original "github.com/Azure/azure-sdk-for-go/services/cognitiveservices/v2.0/textanalytics"

type BaseClient = original.BaseClient
type AzureRegions = original.AzureRegions

const (
	Australiaeast  AzureRegions = original.Australiaeast
	Brazilsouth    AzureRegions = original.Brazilsouth
	Eastasia       AzureRegions = original.Eastasia
	Eastus         AzureRegions = original.Eastus
	Eastus2        AzureRegions = original.Eastus2
	Northeurope    AzureRegions = original.Northeurope
	Southcentralus AzureRegions = original.Southcentralus
	Southeastasia  AzureRegions = original.Southeastasia
	Westcentralus  AzureRegions = original.Westcentralus
	Westeurope     AzureRegions = original.Westeurope
	Westus         AzureRegions = original.Westus
	Westus2        AzureRegions = original.Westus2
)

type BatchInput = original.BatchInput
type DetectedLanguage = original.DetectedLanguage
type ErrorRecord = original.ErrorRecord
type ErrorResponse = original.ErrorResponse
type Input = original.Input
type InternalError = original.InternalError
type KeyPhraseBatchResult = original.KeyPhraseBatchResult
type KeyPhraseBatchResultItem = original.KeyPhraseBatchResultItem
type LanguageBatchResult = original.LanguageBatchResult
type LanguageBatchResultItem = original.LanguageBatchResultItem
type MultiLanguageBatchInput = original.MultiLanguageBatchInput
type MultiLanguageInput = original.MultiLanguageInput
type SentimentBatchResult = original.SentimentBatchResult
type SentimentBatchResultItem = original.SentimentBatchResultItem

func UserAgent() string {
	return original.UserAgent() + " profiles/preview"
}
func Version() string {
	return original.Version()
}
func New(azureRegion AzureRegions) BaseClient {
	return original.New(azureRegion)
}
func NewWithoutDefaults(azureRegion AzureRegions) BaseClient {
	return original.NewWithoutDefaults(azureRegion)
}

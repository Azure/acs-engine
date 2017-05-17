package mobileengagement

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
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/Azure/go-autorest/autorest/to"
	"net/http"
)

// AudienceOperators enumerates the values for audience operators.
type AudienceOperators string

const (
	// EQ specifies the eq state for audience operators.
	EQ AudienceOperators = "EQ"
	// GE specifies the ge state for audience operators.
	GE AudienceOperators = "GE"
	// GT specifies the gt state for audience operators.
	GT AudienceOperators = "GT"
	// LE specifies the le state for audience operators.
	LE AudienceOperators = "LE"
	// LT specifies the lt state for audience operators.
	LT AudienceOperators = "LT"
)

// CampaignFeedbacks enumerates the values for campaign feedbacks.
type CampaignFeedbacks string

const (
	// Actioned specifies the actioned state for campaign feedbacks.
	Actioned CampaignFeedbacks = "actioned"
	// Exited specifies the exited state for campaign feedbacks.
	Exited CampaignFeedbacks = "exited"
	// Pushed specifies the pushed state for campaign feedbacks.
	Pushed CampaignFeedbacks = "pushed"
	// Replied specifies the replied state for campaign feedbacks.
	Replied CampaignFeedbacks = "replied"
)

// CampaignKinds enumerates the values for campaign kinds.
type CampaignKinds string

const (
	// Announcements specifies the announcements state for campaign kinds.
	Announcements CampaignKinds = "announcements"
	// DataPushes specifies the data pushes state for campaign kinds.
	DataPushes CampaignKinds = "dataPushes"
	// NativePushes specifies the native pushes state for campaign kinds.
	NativePushes CampaignKinds = "nativePushes"
	// Polls specifies the polls state for campaign kinds.
	Polls CampaignKinds = "polls"
)

// CampaignStates enumerates the values for campaign states.
type CampaignStates string

const (
	// Draft specifies the draft state for campaign states.
	Draft CampaignStates = "draft"
	// Finished specifies the finished state for campaign states.
	Finished CampaignStates = "finished"
	// InProgress specifies the in progress state for campaign states.
	InProgress CampaignStates = "in-progress"
	// Queued specifies the queued state for campaign states.
	Queued CampaignStates = "queued"
	// Scheduled specifies the scheduled state for campaign states.
	Scheduled CampaignStates = "scheduled"
)

// CampaignType enumerates the values for campaign type.
type CampaignType string

const (
	// Announcement specifies the announcement state for campaign type.
	Announcement CampaignType = "Announcement"
	// DataPush specifies the data push state for campaign type.
	DataPush CampaignType = "DataPush"
	// NativePush specifies the native push state for campaign type.
	NativePush CampaignType = "NativePush"
	// Poll specifies the poll state for campaign type.
	Poll CampaignType = "Poll"
)

// CampaignTypes enumerates the values for campaign types.
type CampaignTypes string

const (
	// OnlyNotif specifies the only notif state for campaign types.
	OnlyNotif CampaignTypes = "only_notif"
	// Textbase64 specifies the textbase 64 state for campaign types.
	Textbase64 CampaignTypes = "text/base64"
	// Texthtml specifies the texthtml state for campaign types.
	Texthtml CampaignTypes = "text/html"
	// Textplain specifies the textplain state for campaign types.
	Textplain CampaignTypes = "text/plain"
)

// DeliveryTimes enumerates the values for delivery times.
type DeliveryTimes string

const (
	// Any specifies the any state for delivery times.
	Any DeliveryTimes = "any"
	// Background specifies the background state for delivery times.
	Background DeliveryTimes = "background"
	// Session specifies the session state for delivery times.
	Session DeliveryTimes = "session"
)

// ExportFormat enumerates the values for export format.
type ExportFormat string

const (
	// CsvBlob specifies the csv blob state for export format.
	CsvBlob ExportFormat = "CsvBlob"
	// JSONBlob specifies the json blob state for export format.
	JSONBlob ExportFormat = "JsonBlob"
)

// ExportState enumerates the values for export state.
type ExportState string

const (
	// ExportStateFailed specifies the export state failed state for export
	// state.
	ExportStateFailed ExportState = "Failed"
	// ExportStateQueued specifies the export state queued state for export
	// state.
	ExportStateQueued ExportState = "Queued"
	// ExportStateStarted specifies the export state started state for export
	// state.
	ExportStateStarted ExportState = "Started"
	// ExportStateSucceeded specifies the export state succeeded state for
	// export state.
	ExportStateSucceeded ExportState = "Succeeded"
)

// ExportType enumerates the values for export type.
type ExportType string

const (
	// ExportTypeActivity specifies the export type activity state for export
	// type.
	ExportTypeActivity ExportType = "Activity"
	// ExportTypeCrash specifies the export type crash state for export type.
	ExportTypeCrash ExportType = "Crash"
	// ExportTypeError specifies the export type error state for export type.
	ExportTypeError ExportType = "Error"
	// ExportTypeEvent specifies the export type event state for export type.
	ExportTypeEvent ExportType = "Event"
	// ExportTypeJob specifies the export type job state for export type.
	ExportTypeJob ExportType = "Job"
	// ExportTypePush specifies the export type push state for export type.
	ExportTypePush ExportType = "Push"
	// ExportTypeSession specifies the export type session state for export
	// type.
	ExportTypeSession ExportType = "Session"
	// ExportTypeTag specifies the export type tag state for export type.
	ExportTypeTag ExportType = "Tag"
	// ExportTypeToken specifies the export type token state for export type.
	ExportTypeToken ExportType = "Token"
)

// JobStates enumerates the values for job states.
type JobStates string

const (
	// JobStatesFailed specifies the job states failed state for job states.
	JobStatesFailed JobStates = "Failed"
	// JobStatesQueued specifies the job states queued state for job states.
	JobStatesQueued JobStates = "Queued"
	// JobStatesStarted specifies the job states started state for job states.
	JobStatesStarted JobStates = "Started"
	// JobStatesSucceeded specifies the job states succeeded state for job
	// states.
	JobStatesSucceeded JobStates = "Succeeded"
)

// NotificationTypes enumerates the values for notification types.
type NotificationTypes string

const (
	// Popup specifies the popup state for notification types.
	Popup NotificationTypes = "popup"
	// System specifies the system state for notification types.
	System NotificationTypes = "system"
)

// ProvisioningStates enumerates the values for provisioning states.
type ProvisioningStates string

const (
	// Creating specifies the creating state for provisioning states.
	Creating ProvisioningStates = "Creating"
	// Succeeded specifies the succeeded state for provisioning states.
	Succeeded ProvisioningStates = "Succeeded"
)

// PushModes enumerates the values for push modes.
type PushModes string

const (
	// Manual specifies the manual state for push modes.
	Manual PushModes = "manual"
	// OneShot specifies the one shot state for push modes.
	OneShot PushModes = "one-shot"
	// RealTime specifies the real time state for push modes.
	RealTime PushModes = "real-time"
)

// AnnouncementFeedbackCriterion is used to target devices who received an
// announcement.
type AnnouncementFeedbackCriterion struct {
	ContentID *int32            `json:"content-id,omitempty"`
	Action    CampaignFeedbacks `json:"action,omitempty"`
}

// APIError is
type APIError struct {
	Error *APIErrorError `json:"error,omitempty"`
}

// APIErrorError is
type APIErrorError struct {
	Code    *string `json:"code,omitempty"`
	Message *string `json:"message,omitempty"`
}

// App is the Mobile Engagement App resource.
type App struct {
	ID             *string             `json:"id,omitempty"`
	Name           *string             `json:"name,omitempty"`
	Type           *string             `json:"type,omitempty"`
	Location       *string             `json:"location,omitempty"`
	Tags           *map[string]*string `json:"tags,omitempty"`
	*AppProperties `json:"properties,omitempty"`
}

// AppCollection is the AppCollection resource.
type AppCollection struct {
	ID                       *string             `json:"id,omitempty"`
	Name                     *string             `json:"name,omitempty"`
	Type                     *string             `json:"type,omitempty"`
	Location                 *string             `json:"location,omitempty"`
	Tags                     *map[string]*string `json:"tags,omitempty"`
	*AppCollectionProperties `json:"properties,omitempty"`
}

// AppCollectionListResult is the list AppCollections operation response.
type AppCollectionListResult struct {
	autorest.Response `json:"-"`
	Value             *[]AppCollection `json:"value,omitempty"`
	NextLink          *string          `json:"nextLink,omitempty"`
}

// AppCollectionListResultPreparer prepares a request to retrieve the next set of results. It returns
// nil if no more results exist.
func (client AppCollectionListResult) AppCollectionListResultPreparer() (*http.Request, error) {
	if client.NextLink == nil || len(to.String(client.NextLink)) <= 0 {
		return nil, nil
	}
	return autorest.Prepare(&http.Request{},
		autorest.AsJSON(),
		autorest.AsGet(),
		autorest.WithBaseURL(to.String(client.NextLink)))
}

// AppCollectionNameAvailability is
type AppCollectionNameAvailability struct {
	autorest.Response    `json:"-"`
	Name                 *string `json:"name,omitempty"`
	Available            *bool   `json:"available,omitempty"`
	UnavailabilityReason *string `json:"unavailabilityReason,omitempty"`
}

// AppCollectionProperties is
type AppCollectionProperties struct {
	ProvisioningState ProvisioningStates `json:"provisioningState,omitempty"`
}

// AppInfoFilter is send only to users who have some app info set. This is a
// special filter that is automatically added if your campaign contains appInfo
// parameters. It is not intended to be public and should not be used as it
// could be removed or replaced by the API.
type AppInfoFilter struct {
	AppInfo *[]string `json:"appInfo,omitempty"`
}

// ApplicationVersionCriterion is used to target devices based on the version
// of the application they are using.
type ApplicationVersionCriterion struct {
	Name *string `json:"name,omitempty"`
}

// AppListResult is the list Apps operation response.
type AppListResult struct {
	autorest.Response `json:"-"`
	Value             *[]App  `json:"value,omitempty"`
	NextLink          *string `json:"nextLink,omitempty"`
}

// AppListResultPreparer prepares a request to retrieve the next set of results. It returns
// nil if no more results exist.
func (client AppListResult) AppListResultPreparer() (*http.Request, error) {
	if client.NextLink == nil || len(to.String(client.NextLink)) <= 0 {
		return nil, nil
	}
	return autorest.Prepare(&http.Request{},
		autorest.AsJSON(),
		autorest.AsGet(),
		autorest.WithBaseURL(to.String(client.NextLink)))
}

// AppProperties is
type AppProperties struct {
	BackendID *string `json:"backendId,omitempty"`
	Platform  *string `json:"platform,omitempty"`
	AppState  *string `json:"appState,omitempty"`
}

// BooleanTagCriterion is target devices based on a boolean tag value.
type BooleanTagCriterion struct {
	Name  *string `json:"name,omitempty"`
	Value *bool   `json:"value,omitempty"`
}

// Campaign is
type Campaign struct {
	NotificationTitle     *string                           `json:"notificationTitle,omitempty"`
	NotificationMessage   *string                           `json:"notificationMessage,omitempty"`
	NotificationImage     *[]byte                           `json:"notificationImage,omitempty"`
	NotificationOptions   *NotificationOptions              `json:"notificationOptions,omitempty"`
	Title                 *string                           `json:"title,omitempty"`
	Body                  *string                           `json:"body,omitempty"`
	ActionButtonText      *string                           `json:"actionButtonText,omitempty"`
	ExitButtonText        *string                           `json:"exitButtonText,omitempty"`
	ActionURL             *string                           `json:"actionUrl,omitempty"`
	Payload               *map[string]interface{}           `json:"payload,omitempty"`
	Name                  *string                           `json:"name,omitempty"`
	Audience              *CampaignAudience                 `json:"audience,omitempty"`
	Category              *string                           `json:"category,omitempty"`
	PushMode              PushModes                         `json:"pushMode,omitempty"`
	Type                  CampaignTypes                     `json:"type,omitempty"`
	DeliveryTime          DeliveryTimes                     `json:"deliveryTime,omitempty"`
	DeliveryActivities    *[]string                         `json:"deliveryActivities,omitempty"`
	StartTime             *string                           `json:"startTime,omitempty"`
	EndTime               *string                           `json:"endTime,omitempty"`
	Timezone              *string                           `json:"timezone,omitempty"`
	NotificationType      NotificationTypes                 `json:"notificationType,omitempty"`
	NotificationIcon      *bool                             `json:"notificationIcon,omitempty"`
	NotificationCloseable *bool                             `json:"notificationCloseable,omitempty"`
	NotificationVibrate   *bool                             `json:"notificationVibrate,omitempty"`
	NotificationSound     *bool                             `json:"notificationSound,omitempty"`
	NotificationBadge     *bool                             `json:"notificationBadge,omitempty"`
	Localization          *map[string]*CampaignLocalization `json:"localization,omitempty"`
	Questions             *[]PollQuestion                   `json:"questions,omitempty"`
}

// CampaignAudience is specify which users will be targeted by this campaign.
// By default, all users will be targeted. If you set `pushMode` property to
// `manual`, the only thing you can specify in the audience is the push quota
// filter. An audience is a boolean expression made of criteria (variables)
// operators (`not`, `and` or `or`) and parenthesis. Additionally, a set of
// filters can be added to an audience. 65535 bytes max as per JSON encoding.
type CampaignAudience struct {
	Expression *string                `json:"expression,omitempty"`
	Criteria   *map[string]*Criterion `json:"criteria,omitempty"`
	Filters    *[]Filter              `json:"filters,omitempty"`
}

// CampaignListResult is
type CampaignListResult struct {
	State         CampaignStates `json:"state,omitempty"`
	ID            *int32         `json:"id,omitempty"`
	Name          *string        `json:"name,omitempty"`
	ActivatedDate *date.Time     `json:"activatedDate,omitempty"`
	FinishedDate  *date.Time     `json:"finishedDate,omitempty"`
	StartTime     *date.Time     `json:"startTime,omitempty"`
	EndTime       *date.Time     `json:"endTime,omitempty"`
	Timezone      *string        `json:"timezone,omitempty"`
}

// CampaignLocalization is
type CampaignLocalization struct {
	NotificationTitle   *string                 `json:"notificationTitle,omitempty"`
	NotificationMessage *string                 `json:"notificationMessage,omitempty"`
	NotificationImage   *[]byte                 `json:"notificationImage,omitempty"`
	NotificationOptions *NotificationOptions    `json:"notificationOptions,omitempty"`
	Title               *string                 `json:"title,omitempty"`
	Body                *string                 `json:"body,omitempty"`
	ActionButtonText    *string                 `json:"actionButtonText,omitempty"`
	ExitButtonText      *string                 `json:"exitButtonText,omitempty"`
	ActionURL           *string                 `json:"actionUrl,omitempty"`
	Payload             *map[string]interface{} `json:"payload,omitempty"`
}

// CampaignPushParameters is
type CampaignPushParameters struct {
	DeviceIds *[]string `json:"deviceIds,omitempty"`
	Data      *Campaign `json:"data,omitempty"`
}

// CampaignPushResult is
type CampaignPushResult struct {
	autorest.Response `json:"-"`
	InvalidDeviceIds  *[]string `json:"invalidDeviceIds,omitempty"`
}

// CampaignResult is
type CampaignResult struct {
	autorest.Response     `json:"-"`
	NotificationTitle     *string                           `json:"notificationTitle,omitempty"`
	NotificationMessage   *string                           `json:"notificationMessage,omitempty"`
	NotificationImage     *[]byte                           `json:"notificationImage,omitempty"`
	NotificationOptions   *NotificationOptions              `json:"notificationOptions,omitempty"`
	Title                 *string                           `json:"title,omitempty"`
	Body                  *string                           `json:"body,omitempty"`
	ActionButtonText      *string                           `json:"actionButtonText,omitempty"`
	ExitButtonText        *string                           `json:"exitButtonText,omitempty"`
	ActionURL             *string                           `json:"actionUrl,omitempty"`
	Payload               *map[string]interface{}           `json:"payload,omitempty"`
	Name                  *string                           `json:"name,omitempty"`
	Audience              *CampaignAudience                 `json:"audience,omitempty"`
	Category              *string                           `json:"category,omitempty"`
	PushMode              PushModes                         `json:"pushMode,omitempty"`
	Type                  CampaignTypes                     `json:"type,omitempty"`
	DeliveryTime          DeliveryTimes                     `json:"deliveryTime,omitempty"`
	DeliveryActivities    *[]string                         `json:"deliveryActivities,omitempty"`
	StartTime             *string                           `json:"startTime,omitempty"`
	EndTime               *string                           `json:"endTime,omitempty"`
	Timezone              *string                           `json:"timezone,omitempty"`
	NotificationType      NotificationTypes                 `json:"notificationType,omitempty"`
	NotificationIcon      *bool                             `json:"notificationIcon,omitempty"`
	NotificationCloseable *bool                             `json:"notificationCloseable,omitempty"`
	NotificationVibrate   *bool                             `json:"notificationVibrate,omitempty"`
	NotificationSound     *bool                             `json:"notificationSound,omitempty"`
	NotificationBadge     *bool                             `json:"notificationBadge,omitempty"`
	Localization          *map[string]*CampaignLocalization `json:"localization,omitempty"`
	Questions             *[]PollQuestion                   `json:"questions,omitempty"`
	ID                    *int32                            `json:"id,omitempty"`
	State                 CampaignStates                    `json:"state,omitempty"`
	ActivatedDate         *date.Time                        `json:"activatedDate,omitempty"`
	FinishedDate          *date.Time                        `json:"finishedDate,omitempty"`
}

// CampaignsListResult is the campaigns list result.
type CampaignsListResult struct {
	autorest.Response `json:"-"`
	Value             *[]CampaignListResult `json:"value,omitempty"`
	NextLink          *string               `json:"nextLink,omitempty"`
}

// CampaignsListResultPreparer prepares a request to retrieve the next set of results. It returns
// nil if no more results exist.
func (client CampaignsListResult) CampaignsListResultPreparer() (*http.Request, error) {
	if client.NextLink == nil || len(to.String(client.NextLink)) <= 0 {
		return nil, nil
	}
	return autorest.Prepare(&http.Request{},
		autorest.AsJSON(),
		autorest.AsGet(),
		autorest.WithBaseURL(to.String(client.NextLink)))
}

// CampaignState is
type CampaignState struct {
	autorest.Response `json:"-"`
	State             CampaignStates `json:"state,omitempty"`
}

// CampaignStateResult is
type CampaignStateResult struct {
	autorest.Response `json:"-"`
	State             CampaignStates `json:"state,omitempty"`
	ID                *int32         `json:"id,omitempty"`
}

// CampaignStatisticsResult is
type CampaignStatisticsResult struct {
	autorest.Response           `json:"-"`
	Queued                      *int32                              `json:"queued,omitempty"`
	Pushed                      *int32                              `json:"pushed,omitempty"`
	PushedNative                *int32                              `json:"pushed-native,omitempty"`
	PushedNativeGoogle          *int32                              `json:"pushed-native-google,omitempty"`
	PushedNativeAdm             *int32                              `json:"pushed-native-adm,omitempty"`
	Delivered                   *int32                              `json:"delivered,omitempty"`
	Dropped                     *int32                              `json:"dropped,omitempty"`
	SystemNotificationDisplayed *int32                              `json:"system-notification-displayed,omitempty"`
	InAppNotificationDisplayed  *int32                              `json:"in-app-notification-displayed,omitempty"`
	ContentDisplayed            *int32                              `json:"content-displayed,omitempty"`
	SystemNotificationActioned  *int32                              `json:"system-notification-actioned,omitempty"`
	SystemNotificationExited    *int32                              `json:"system-notification-exited,omitempty"`
	InAppNotificationActioned   *int32                              `json:"in-app-notification-actioned,omitempty"`
	InAppNotificationExited     *int32                              `json:"in-app-notification-exited,omitempty"`
	ContentActioned             *int32                              `json:"content-actioned,omitempty"`
	ContentExited               *int32                              `json:"content-exited,omitempty"`
	Answers                     *map[string]*map[string]interface{} `json:"answers,omitempty"`
}

// CampaignTestNewParameters is
type CampaignTestNewParameters struct {
	DeviceID *string   `json:"deviceId,omitempty"`
	Lang     *string   `json:"lang,omitempty"`
	Data     *Campaign `json:"data,omitempty"`
}

// CampaignTestSavedParameters is
type CampaignTestSavedParameters struct {
	DeviceID *string `json:"deviceId,omitempty"`
	Lang     *string `json:"lang,omitempty"`
}

// CarrierCountryCriterion is used to target devices based on their carrier
// country.
type CarrierCountryCriterion struct {
	Name *string `json:"name,omitempty"`
}

// CarrierNameCriterion is used to target devices based on their carrier name.
type CarrierNameCriterion struct {
	Name *string `json:"name,omitempty"`
}

// Criterion is
type Criterion struct {
}

// DatapushFeedbackCriterion is used to target devices who received a data
// push.
type DatapushFeedbackCriterion struct {
	ContentID *int32            `json:"content-id,omitempty"`
	Action    CampaignFeedbacks `json:"action,omitempty"`
}

// DateRangeExportTaskParameter is
type DateRangeExportTaskParameter struct {
	ContainerURL *string      `json:"containerUrl,omitempty"`
	Description  *string      `json:"description,omitempty"`
	StartDate    *date.Date   `json:"startDate,omitempty"`
	EndDate      *date.Date   `json:"endDate,omitempty"`
	ExportFormat ExportFormat `json:"exportFormat,omitempty"`
}

// DateTagCriterion is target devices based on a date tag value.
type DateTagCriterion struct {
	Name  *string           `json:"name,omitempty"`
	Value *date.Date        `json:"value,omitempty"`
	Op    AudienceOperators `json:"op,omitempty"`
}

// Device is
type Device struct {
	autorest.Response `json:"-"`
	DeviceID          *string             `json:"deviceId,omitempty"`
	Meta              *DeviceMeta         `json:"meta,omitempty"`
	Info              *DeviceInfo         `json:"info,omitempty"`
	Location          *DeviceLocation     `json:"location,omitempty"`
	AppInfo           *map[string]*string `json:"appInfo,omitempty"`
}

// DeviceInfo is
type DeviceInfo struct {
	PhoneModel             *string `json:"phoneModel,omitempty"`
	PhoneManufacturer      *string `json:"phoneManufacturer,omitempty"`
	FirmwareVersion        *string `json:"firmwareVersion,omitempty"`
	FirmwareName           *string `json:"firmwareName,omitempty"`
	AndroidAPILevel        *int32  `json:"androidAPILevel,omitempty"`
	CarrierCountry         *string `json:"carrierCountry,omitempty"`
	Locale                 *string `json:"locale,omitempty"`
	CarrierName            *string `json:"carrierName,omitempty"`
	NetworkType            *string `json:"networkType,omitempty"`
	NetworkSubtype         *string `json:"networkSubtype,omitempty"`
	ApplicationVersionName *string `json:"applicationVersionName,omitempty"`
	ApplicationVersionCode *int32  `json:"applicationVersionCode,omitempty"`
	TimeZoneOffset         *int32  `json:"timeZoneOffset,omitempty"`
	ServiceVersion         *string `json:"serviceVersion,omitempty"`
}

// DeviceLocation is
type DeviceLocation struct {
	Countrycode *string `json:"countrycode,omitempty"`
	Region      *string `json:"region,omitempty"`
	Locality    *string `json:"locality,omitempty"`
}

// DeviceManufacturerCriterion is used to target devices based on the device
// manufacturer.
type DeviceManufacturerCriterion struct {
	Name *string `json:"name,omitempty"`
}

// DeviceMeta is
type DeviceMeta struct {
	FirstSeen         *int64 `json:"firstSeen,omitempty"`
	LastSeen          *int64 `json:"lastSeen,omitempty"`
	LastInfo          *int64 `json:"lastInfo,omitempty"`
	LastLocation      *int64 `json:"lastLocation,omitempty"`
	NativePushEnabled *bool  `json:"nativePushEnabled,omitempty"`
}

// DeviceModelCriterion is used to target devices based on the device model.
type DeviceModelCriterion struct {
	Name *string `json:"name,omitempty"`
}

// DeviceQueryResult is
type DeviceQueryResult struct {
	DeviceID *string             `json:"deviceId,omitempty"`
	Meta     *DeviceMeta         `json:"meta,omitempty"`
	AppInfo  *map[string]*string `json:"appInfo,omitempty"`
}

// DevicesQueryResult is the campaigns list result.
type DevicesQueryResult struct {
	autorest.Response `json:"-"`
	Value             *[]DeviceQueryResult `json:"value,omitempty"`
	NextLink          *string              `json:"nextLink,omitempty"`
}

// DevicesQueryResultPreparer prepares a request to retrieve the next set of results. It returns
// nil if no more results exist.
func (client DevicesQueryResult) DevicesQueryResultPreparer() (*http.Request, error) {
	if client.NextLink == nil || len(to.String(client.NextLink)) <= 0 {
		return nil, nil
	}
	return autorest.Prepare(&http.Request{},
		autorest.AsJSON(),
		autorest.AsGet(),
		autorest.WithBaseURL(to.String(client.NextLink)))
}

// DeviceTagsParameters is
type DeviceTagsParameters struct {
	Tags         *map[string]map[string]*string `json:"tags,omitempty"`
	DeleteOnNull *bool                          `json:"deleteOnNull,omitempty"`
}

// DeviceTagsResult is
type DeviceTagsResult struct {
	autorest.Response `json:"-"`
	InvalidIds        *[]string `json:"invalidIds,omitempty"`
}

// EngageActiveUsersFilter is send only to users who have used the app in the
// last {threshold} days.
type EngageActiveUsersFilter struct {
	Threshold *int32 `json:"threshold,omitempty"`
}

// EngageIdleUsersFilter is send only to users who haven't used the app in the
// last {threshold} days.
type EngageIdleUsersFilter struct {
	Threshold *int32 `json:"threshold,omitempty"`
}

// EngageNewUsersFilter is send only to users whose first app use is less than
// {threshold} days old.
type EngageNewUsersFilter struct {
	Threshold *int32 `json:"threshold,omitempty"`
}

// EngageOldUsersFilter is send only to users whose first app use is more than
// {threshold} days old.
type EngageOldUsersFilter struct {
	Threshold *int32 `json:"threshold,omitempty"`
}

// EngageSubsetFilter is send only to a maximum of max users.
type EngageSubsetFilter struct {
	Max *int32 `json:"max,omitempty"`
}

// ExportOptions is options to control export generation.
type ExportOptions struct {
	ExportUserID *bool `json:"exportUserId,omitempty"`
}

// ExportTaskListResult is gets a paged list of ExportTasks.
type ExportTaskListResult struct {
	autorest.Response `json:"-"`
	Value             *[]ExportTaskResult `json:"value,omitempty"`
	NextLink          *string             `json:"nextLink,omitempty"`
}

// ExportTaskListResultPreparer prepares a request to retrieve the next set of results. It returns
// nil if no more results exist.
func (client ExportTaskListResult) ExportTaskListResultPreparer() (*http.Request, error) {
	if client.NextLink == nil || len(to.String(client.NextLink)) <= 0 {
		return nil, nil
	}
	return autorest.Prepare(&http.Request{},
		autorest.AsJSON(),
		autorest.AsGet(),
		autorest.WithBaseURL(to.String(client.NextLink)))
}

// ExportTaskParameter is
type ExportTaskParameter struct {
	ContainerURL *string      `json:"containerUrl,omitempty"`
	Description  *string      `json:"description,omitempty"`
	ExportFormat ExportFormat `json:"exportFormat,omitempty"`
}

// ExportTaskResult is
type ExportTaskResult struct {
	autorest.Response `json:"-"`
	ID                *string     `json:"id,omitempty"`
	Description       *string     `json:"description,omitempty"`
	State             ExportState `json:"state,omitempty"`
	DateCreated       *date.Time  `json:"dateCreated,omitempty"`
	DateCompleted     *date.Time  `json:"dateCompleted,omitempty"`
	ExportType        ExportType  `json:"exportType,omitempty"`
	ErrorDetails      *string     `json:"errorDetails,omitempty"`
}

// FeedbackByCampaignParameter is
type FeedbackByCampaignParameter struct {
	ContainerURL *string      `json:"containerUrl,omitempty"`
	Description  *string      `json:"description,omitempty"`
	CampaignType CampaignType `json:"campaignType,omitempty"`
	CampaignIds  *[]int32     `json:"campaignIds,omitempty"`
	ExportFormat ExportFormat `json:"exportFormat,omitempty"`
}

// FeedbackByDateRangeParameter is
type FeedbackByDateRangeParameter struct {
	ContainerURL        *string      `json:"containerUrl,omitempty"`
	Description         *string      `json:"description,omitempty"`
	CampaignType        CampaignType `json:"campaignType,omitempty"`
	CampaignWindowStart *date.Time   `json:"campaignWindowStart,omitempty"`
	CampaignWindowEnd   *date.Time   `json:"campaignWindowEnd,omitempty"`
	ExportFormat        ExportFormat `json:"exportFormat,omitempty"`
}

// Filter is
type Filter struct {
}

// FirmwareVersionCriterion is used to target devices based on their firmware
// version.
type FirmwareVersionCriterion struct {
	Name *string `json:"name,omitempty"`
}

// GeoFencingCriterion is used to target devices based on a specific region. A
// center point (defined by a latitude and longitude) and a radius form the
// boundary for the region. This criterion will be met when the user crosses
// the boundaries of the region.
type GeoFencingCriterion struct {
	Lat        *float64 `json:"lat,omitempty"`
	Lon        *float64 `json:"lon,omitempty"`
	Radius     *int32   `json:"radius,omitempty"`
	Expiration *int32   `json:"expiration,omitempty"`
}

// ImportTask is
type ImportTask struct {
	StorageURL *string `json:"storageUrl,omitempty"`
}

// ImportTaskListResult is gets a paged list of import tasks.
type ImportTaskListResult struct {
	autorest.Response `json:"-"`
	Value             *[]ImportTaskResult `json:"value,omitempty"`
	NextLink          *string             `json:"nextLink,omitempty"`
}

// ImportTaskListResultPreparer prepares a request to retrieve the next set of results. It returns
// nil if no more results exist.
func (client ImportTaskListResult) ImportTaskListResultPreparer() (*http.Request, error) {
	if client.NextLink == nil || len(to.String(client.NextLink)) <= 0 {
		return nil, nil
	}
	return autorest.Prepare(&http.Request{},
		autorest.AsJSON(),
		autorest.AsGet(),
		autorest.WithBaseURL(to.String(client.NextLink)))
}

// ImportTaskResult is
type ImportTaskResult struct {
	autorest.Response `json:"-"`
	StorageURL        *string    `json:"storageUrl,omitempty"`
	ID                *string    `json:"id,omitempty"`
	State             JobStates  `json:"state,omitempty"`
	DateCreated       *date.Time `json:"dateCreated,omitempty"`
	DateCompleted     *date.Time `json:"dateCompleted,omitempty"`
	ErrorDetails      *string    `json:"errorDetails,omitempty"`
}

// IntegerTagCriterion is target devices based on an integer tag value.
type IntegerTagCriterion struct {
	Name  *string           `json:"name,omitempty"`
	Value *int32            `json:"value,omitempty"`
	Op    AudienceOperators `json:"op,omitempty"`
}

// LanguageCriterion is used to target devices based on the language of their
// device.
type LanguageCriterion struct {
	Name *string `json:"name,omitempty"`
}

// LocationCriterion is used to target devices based on their last know area.
type LocationCriterion struct {
	Country  *string `json:"country,omitempty"`
	Region   *string `json:"region,omitempty"`
	Locality *string `json:"locality,omitempty"`
}

// NativePushEnabledFilter is engage only users with native push enabled.
type NativePushEnabledFilter struct {
}

// NetworkTypeCriterion is used to target devices based their network type.
type NetworkTypeCriterion struct {
	Name *string `json:"name,omitempty"`
}

// NotificationOptions is
type NotificationOptions struct {
	BigText    *string `json:"bigText,omitempty"`
	BigPicture *string `json:"bigPicture,omitempty"`
	Sound      *string `json:"sound,omitempty"`
	ActionText *string `json:"actionText,omitempty"`
}

// PollAnswerFeedbackCriterion is used to target devices who answered X to a
// given question.
type PollAnswerFeedbackCriterion struct {
	ContentID *int32 `json:"content-id,omitempty"`
	ChoiceID  *int32 `json:"choice-id,omitempty"`
}

// PollFeedbackCriterion is used to target devices who received a poll.
type PollFeedbackCriterion struct {
	ContentID *int32            `json:"content-id,omitempty"`
	Action    CampaignFeedbacks `json:"action,omitempty"`
}

// PollQuestion is
type PollQuestion struct {
	Title        *string                               `json:"title,omitempty"`
	ID           *int32                                `json:"id,omitempty"`
	Localization *map[string]*PollQuestionLocalization `json:"localization,omitempty"`
	Choices      *[]PollQuestionChoice                 `json:"choices,omitempty"`
}

// PollQuestionChoice is
type PollQuestionChoice struct {
	Title        *string                                     `json:"title,omitempty"`
	ID           *int32                                      `json:"id,omitempty"`
	Localization *map[string]*PollQuestionChoiceLocalization `json:"localization,omitempty"`
	IsDefault    *bool                                       `json:"isDefault,omitempty"`
}

// PollQuestionChoiceLocalization is
type PollQuestionChoiceLocalization struct {
	Title *string `json:"title,omitempty"`
}

// PollQuestionLocalization is
type PollQuestionLocalization struct {
	Title *string `json:"title,omitempty"`
}

// PushQuotaFilter is engage only users for whom the push quota is not reached.
type PushQuotaFilter struct {
}

// Resource is
type Resource struct {
	ID       *string             `json:"id,omitempty"`
	Name     *string             `json:"name,omitempty"`
	Type     *string             `json:"type,omitempty"`
	Location *string             `json:"location,omitempty"`
	Tags     *map[string]*string `json:"tags,omitempty"`
}

// ScreenSizeCriterion is used to target devices based on the screen resolution
// of their device.
type ScreenSizeCriterion struct {
	Name *string `json:"name,omitempty"`
}

// SegmentCriterion is target devices based on an existing segment.
type SegmentCriterion struct {
	ID      *int32 `json:"id,omitempty"`
	Exclude *bool  `json:"exclude,omitempty"`
}

// StringTagCriterion is target devices based on a string tag value.
type StringTagCriterion struct {
	Name  *string `json:"name,omitempty"`
	Value *string `json:"value,omitempty"`
}

// SupportedPlatformsListResult is
type SupportedPlatformsListResult struct {
	autorest.Response `json:"-"`
	Platforms         *[]string `json:"platforms,omitempty"`
}

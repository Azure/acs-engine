package resourcehealth

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

// AvailabilityStateValues enumerates the values for availability state values.
type AvailabilityStateValues string

const (
	// Available specifies the available state for availability state values.
	Available AvailabilityStateValues = "Available"
	// Unavailable specifies the unavailable state for availability state
	// values.
	Unavailable AvailabilityStateValues = "Unavailable"
	// Unknown specifies the unknown state for availability state values.
	Unknown AvailabilityStateValues = "Unknown"
)

// ReasonChronicityTypes enumerates the values for reason chronicity types.
type ReasonChronicityTypes string

const (
	// Persistent specifies the persistent state for reason chronicity types.
	Persistent ReasonChronicityTypes = "Persistent"
	// Transient specifies the transient state for reason chronicity types.
	Transient ReasonChronicityTypes = "Transient"
)

// AvailabilityStatus is availabilityStatus of a resource.
type AvailabilityStatus struct {
	autorest.Response `json:"-"`
	ID                *string                       `json:"id,omitempty"`
	Name              *string                       `json:"name,omitempty"`
	Type              *string                       `json:"type,omitempty"`
	Location          *string                       `json:"location,omitempty"`
	Properties        *AvailabilityStatusProperties `json:"properties,omitempty"`
}

// AvailabilityStatusProperties is properties of availability state.
type AvailabilityStatusProperties struct {
	AvailabilityState        AvailabilityStateValues                            `json:"availabilityState,omitempty"`
	Summary                  *string                                            `json:"summary,omitempty"`
	DetailedStatus           *string                                            `json:"detailedStatus,omitempty"`
	ReasonType               *string                                            `json:"reasonType,omitempty"`
	RootCauseAttributionTime *date.Time                                         `json:"rootCauseAttributionTime,omitempty"`
	ResolutionETA            *date.Time                                         `json:"resolutionETA,omitempty"`
	OccuredTime              *date.Time                                         `json:"occuredTime,omitempty"`
	ReasonChronicity         ReasonChronicityTypes                              `json:"reasonChronicity,omitempty"`
	ReportedTime             *date.Time                                         `json:"reportedTime,omitempty"`
	RecentlyResolvedState    *AvailabilityStatusPropertiesRecentlyResolvedState `json:"recentlyResolvedState,omitempty"`
	RecommendedActions       *[]RecommendedAction                               `json:"recommendedActions,omitempty"`
	ServiceImpactingEvents   *[]ServiceImpactingEvent                           `json:"serviceImpactingEvents,omitempty"`
}

// AvailabilityStatusPropertiesRecentlyResolvedState is an annotation
// describing a change in the availabilityState to Available from Unavailable
// with a reasonType of type Unplanned
type AvailabilityStatusPropertiesRecentlyResolvedState struct {
	UnavailableOccurredTime *date.Time `json:"unavailableOccurredTime,omitempty"`
	ResolvedTime            *date.Time `json:"resolvedTime,omitempty"`
	UnavailabilitySummary   *string    `json:"unavailabilitySummary,omitempty"`
}

// AvailabilityStatusListResult is the List availabilityStatus operation
// response.
type AvailabilityStatusListResult struct {
	autorest.Response `json:"-"`
	Value             *[]AvailabilityStatus `json:"value,omitempty"`
	NextLink          *string               `json:"nextLink,omitempty"`
}

// AvailabilityStatusListResultPreparer prepares a request to retrieve the next set of results. It returns
// nil if no more results exist.
func (client AvailabilityStatusListResult) AvailabilityStatusListResultPreparer() (*http.Request, error) {
	if client.NextLink == nil || len(to.String(client.NextLink)) <= 0 {
		return nil, nil
	}
	return autorest.Prepare(&http.Request{},
		autorest.AsJSON(),
		autorest.AsGet(),
		autorest.WithBaseURL(to.String(client.NextLink)))
}

// ErrorResponse is error details.
type ErrorResponse struct {
	Code    *string `json:"code,omitempty"`
	Message *string `json:"message,omitempty"`
	Details *string `json:"details,omitempty"`
}

// Operation is operation available in the resourcehealth resource provider.
type Operation struct {
	Name    *string           `json:"name,omitempty"`
	Display *OperationDisplay `json:"display,omitempty"`
}

// OperationDisplay is properties of the operation.
type OperationDisplay struct {
	Provider    *string `json:"provider,omitempty"`
	Resource    *string `json:"resource,omitempty"`
	Operation   *string `json:"operation,omitempty"`
	Description *string `json:"description,omitempty"`
}

// OperationListResult is lists the operations response.
type OperationListResult struct {
	autorest.Response `json:"-"`
	Value             *[]Operation `json:"value,omitempty"`
}

// RecommendedAction is lists actions the user can take based on the current
// availabilityState of the resource.
type RecommendedAction struct {
	Action        *string `json:"action,omitempty"`
	ActionURL     *string `json:"actionUrl,omitempty"`
	ActionURLText *string `json:"actionUrlText,omitempty"`
}

// ServiceImpactingEvent is lists the service impacting events that may be
// affecting the health of the resource.
type ServiceImpactingEvent struct {
	EventStartTime              *date.Time                               `json:"eventStartTime,omitempty"`
	EventStatusLastModifiedTime *date.Time                               `json:"eventStatusLastModifiedTime,omitempty"`
	CorrelationID               *string                                  `json:"correlationId,omitempty"`
	Status                      *ServiceImpactingEventStatus             `json:"status,omitempty"`
	IncidentProperties          *ServiceImpactingEventIncidentProperties `json:"incidentProperties,omitempty"`
}

// ServiceImpactingEventIncidentProperties is properties of the service
// impacting event.
type ServiceImpactingEventIncidentProperties struct {
	Title        *string `json:"title,omitempty"`
	Service      *string `json:"service,omitempty"`
	Region       *string `json:"region,omitempty"`
	IncidentType *string `json:"incidentType,omitempty"`
}

// ServiceImpactingEventStatus is status of the service impacting event.
type ServiceImpactingEventStatus struct {
	Value *string `json:"value,omitempty"`
}

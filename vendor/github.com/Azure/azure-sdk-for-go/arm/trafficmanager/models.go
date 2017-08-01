package trafficmanager

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
)

// CheckTrafficManagerRelativeDNSNameAvailabilityParameters is parameters
// supplied to check Traffic Manager name operation.
type CheckTrafficManagerRelativeDNSNameAvailabilityParameters struct {
	Name *string `json:"name,omitempty"`
	Type *string `json:"type,omitempty"`
}

// DNSConfig is class containing DNS settings in a Traffic Manager profile.
type DNSConfig struct {
	RelativeName *string `json:"relativeName,omitempty"`
	Fqdn         *string `json:"fqdn,omitempty"`
	TTL          *int64  `json:"ttl,omitempty"`
}

// Endpoint is class representing a Traffic Manager endpoint.
type Endpoint struct {
	autorest.Response   `json:"-"`
	ID                  *string `json:"id,omitempty"`
	Name                *string `json:"name,omitempty"`
	Type                *string `json:"type,omitempty"`
	*EndpointProperties `json:"properties,omitempty"`
}

// EndpointProperties is class representing a Traffic Manager endpoint
// properties.
type EndpointProperties struct {
	TargetResourceID      *string `json:"targetResourceId,omitempty"`
	Target                *string `json:"target,omitempty"`
	EndpointStatus        *string `json:"endpointStatus,omitempty"`
	Weight                *int64  `json:"weight,omitempty"`
	Priority              *int64  `json:"priority,omitempty"`
	EndpointLocation      *string `json:"endpointLocation,omitempty"`
	EndpointMonitorStatus *string `json:"endpointMonitorStatus,omitempty"`
	MinChildEndpoints     *int64  `json:"minChildEndpoints,omitempty"`
}

// MonitorConfig is class containing endpoint monitoring settings in a Traffic
// Manager profile.
type MonitorConfig struct {
	ProfileMonitorStatus *string `json:"profileMonitorStatus,omitempty"`
	Protocol             *string `json:"protocol,omitempty"`
	Port                 *int64  `json:"port,omitempty"`
	Path                 *string `json:"path,omitempty"`
}

// NameAvailability is class representing a Traffic Manager Name Availability
// response.
type NameAvailability struct {
	autorest.Response `json:"-"`
	Name              *string `json:"name,omitempty"`
	Type              *string `json:"type,omitempty"`
	NameAvailable     *bool   `json:"nameAvailable,omitempty"`
	Reason            *string `json:"reason,omitempty"`
	Message           *string `json:"message,omitempty"`
}

// Profile is class representing a Traffic Manager profile.
type Profile struct {
	autorest.Response  `json:"-"`
	ID                 *string             `json:"id,omitempty"`
	Name               *string             `json:"name,omitempty"`
	Type               *string             `json:"type,omitempty"`
	Location           *string             `json:"location,omitempty"`
	Tags               *map[string]*string `json:"tags,omitempty"`
	*ProfileProperties `json:"properties,omitempty"`
}

// ProfileListResult is the list Traffic Manager profiles operation response.
type ProfileListResult struct {
	autorest.Response `json:"-"`
	Value             *[]Profile `json:"value,omitempty"`
}

// ProfileProperties is class representing the Traffic Manager profile
// properties.
type ProfileProperties struct {
	ProfileStatus        *string        `json:"profileStatus,omitempty"`
	TrafficRoutingMethod *string        `json:"trafficRoutingMethod,omitempty"`
	DNSConfig            *DNSConfig     `json:"dnsConfig,omitempty"`
	MonitorConfig        *MonitorConfig `json:"monitorConfig,omitempty"`
	Endpoints            *[]Endpoint    `json:"endpoints,omitempty"`
}

// Resource is
type Resource struct {
	ID       *string             `json:"id,omitempty"`
	Name     *string             `json:"name,omitempty"`
	Type     *string             `json:"type,omitempty"`
	Location *string             `json:"location,omitempty"`
	Tags     *map[string]*string `json:"tags,omitempty"`
}

// SubResource is
type SubResource struct {
	ID *string `json:"id,omitempty"`
}

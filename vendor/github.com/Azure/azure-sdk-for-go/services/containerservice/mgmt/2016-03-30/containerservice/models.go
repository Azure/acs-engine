package containerservice

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
	"net/http"
)

// OchestratorTypes enumerates the values for ochestrator types.
type OchestratorTypes string

const (
	// DCOS ...
	DCOS OchestratorTypes = "DCOS"
	// Swarm ...
	Swarm OchestratorTypes = "Swarm"
)

// VMSizeTypes enumerates the values for vm size types.
type VMSizeTypes string

const (
	// StandardA0 ...
	StandardA0 VMSizeTypes = "Standard_A0"
	// StandardA1 ...
	StandardA1 VMSizeTypes = "Standard_A1"
	// StandardA10 ...
	StandardA10 VMSizeTypes = "Standard_A10"
	// StandardA11 ...
	StandardA11 VMSizeTypes = "Standard_A11"
	// StandardA2 ...
	StandardA2 VMSizeTypes = "Standard_A2"
	// StandardA3 ...
	StandardA3 VMSizeTypes = "Standard_A3"
	// StandardA4 ...
	StandardA4 VMSizeTypes = "Standard_A4"
	// StandardA5 ...
	StandardA5 VMSizeTypes = "Standard_A5"
	// StandardA6 ...
	StandardA6 VMSizeTypes = "Standard_A6"
	// StandardA7 ...
	StandardA7 VMSizeTypes = "Standard_A7"
	// StandardA8 ...
	StandardA8 VMSizeTypes = "Standard_A8"
	// StandardA9 ...
	StandardA9 VMSizeTypes = "Standard_A9"
	// StandardD1 ...
	StandardD1 VMSizeTypes = "Standard_D1"
	// StandardD11 ...
	StandardD11 VMSizeTypes = "Standard_D11"
	// StandardD11V2 ...
	StandardD11V2 VMSizeTypes = "Standard_D11_v2"
	// StandardD12 ...
	StandardD12 VMSizeTypes = "Standard_D12"
	// StandardD12V2 ...
	StandardD12V2 VMSizeTypes = "Standard_D12_v2"
	// StandardD13 ...
	StandardD13 VMSizeTypes = "Standard_D13"
	// StandardD13V2 ...
	StandardD13V2 VMSizeTypes = "Standard_D13_v2"
	// StandardD14 ...
	StandardD14 VMSizeTypes = "Standard_D14"
	// StandardD14V2 ...
	StandardD14V2 VMSizeTypes = "Standard_D14_v2"
	// StandardD1V2 ...
	StandardD1V2 VMSizeTypes = "Standard_D1_v2"
	// StandardD2 ...
	StandardD2 VMSizeTypes = "Standard_D2"
	// StandardD2V2 ...
	StandardD2V2 VMSizeTypes = "Standard_D2_v2"
	// StandardD3 ...
	StandardD3 VMSizeTypes = "Standard_D3"
	// StandardD3V2 ...
	StandardD3V2 VMSizeTypes = "Standard_D3_v2"
	// StandardD4 ...
	StandardD4 VMSizeTypes = "Standard_D4"
	// StandardD4V2 ...
	StandardD4V2 VMSizeTypes = "Standard_D4_v2"
	// StandardD5V2 ...
	StandardD5V2 VMSizeTypes = "Standard_D5_v2"
	// StandardDS1 ...
	StandardDS1 VMSizeTypes = "Standard_DS1"
	// StandardDS11 ...
	StandardDS11 VMSizeTypes = "Standard_DS11"
	// StandardDS12 ...
	StandardDS12 VMSizeTypes = "Standard_DS12"
	// StandardDS13 ...
	StandardDS13 VMSizeTypes = "Standard_DS13"
	// StandardDS14 ...
	StandardDS14 VMSizeTypes = "Standard_DS14"
	// StandardDS2 ...
	StandardDS2 VMSizeTypes = "Standard_DS2"
	// StandardDS3 ...
	StandardDS3 VMSizeTypes = "Standard_DS3"
	// StandardDS4 ...
	StandardDS4 VMSizeTypes = "Standard_DS4"
	// StandardG1 ...
	StandardG1 VMSizeTypes = "Standard_G1"
	// StandardG2 ...
	StandardG2 VMSizeTypes = "Standard_G2"
	// StandardG3 ...
	StandardG3 VMSizeTypes = "Standard_G3"
	// StandardG4 ...
	StandardG4 VMSizeTypes = "Standard_G4"
	// StandardG5 ...
	StandardG5 VMSizeTypes = "Standard_G5"
	// StandardGS1 ...
	StandardGS1 VMSizeTypes = "Standard_GS1"
	// StandardGS2 ...
	StandardGS2 VMSizeTypes = "Standard_GS2"
	// StandardGS3 ...
	StandardGS3 VMSizeTypes = "Standard_GS3"
	// StandardGS4 ...
	StandardGS4 VMSizeTypes = "Standard_GS4"
	// StandardGS5 ...
	StandardGS5 VMSizeTypes = "Standard_GS5"
)

// AgentPoolProfile profile for the container service agent pool.
type AgentPoolProfile struct {
	// Name - Unique name of the agent pool profile in the context of the subscription and resource group.
	Name *string `json:"name,omitempty"`
	// Count - Number of agents (VMs) to host docker containers. Allowed values must be in the range of 1 to 100 (inclusive). The default value is 1.
	Count *int32 `json:"count,omitempty"`
	// VMSize - Size of agent VMs. Possible values include: 'StandardA0', 'StandardA1', 'StandardA2', 'StandardA3', 'StandardA4', 'StandardA5', 'StandardA6', 'StandardA7', 'StandardA8', 'StandardA9', 'StandardA10', 'StandardA11', 'StandardD1', 'StandardD2', 'StandardD3', 'StandardD4', 'StandardD11', 'StandardD12', 'StandardD13', 'StandardD14', 'StandardD1V2', 'StandardD2V2', 'StandardD3V2', 'StandardD4V2', 'StandardD5V2', 'StandardD11V2', 'StandardD12V2', 'StandardD13V2', 'StandardD14V2', 'StandardG1', 'StandardG2', 'StandardG3', 'StandardG4', 'StandardG5', 'StandardDS1', 'StandardDS2', 'StandardDS3', 'StandardDS4', 'StandardDS11', 'StandardDS12', 'StandardDS13', 'StandardDS14', 'StandardGS1', 'StandardGS2', 'StandardGS3', 'StandardGS4', 'StandardGS5'
	VMSize VMSizeTypes `json:"vmSize,omitempty"`
	// DNSPrefix - DNS prefix to be used to create the FQDN for the agent pool.
	DNSPrefix *string `json:"dnsPrefix,omitempty"`
	// Fqdn - FDQN for the agent pool.
	Fqdn *string `json:"fqdn,omitempty"`
}

// ContainerService container service.
type ContainerService struct {
	autorest.Response `json:"-"`
	// ID - Resource Id
	ID *string `json:"id,omitempty"`
	// Name - Resource name
	Name *string `json:"name,omitempty"`
	// Type - Resource type
	Type *string `json:"type,omitempty"`
	// Location - Resource location
	Location *string `json:"location,omitempty"`
	// Tags - Resource tags
	Tags        *map[string]*string `json:"tags,omitempty"`
	*Properties `json:"properties,omitempty"`
}

// UnmarshalJSON is the custom unmarshaler for ContainerService struct.
func (cs *ContainerService) UnmarshalJSON(body []byte) error {
	var m map[string]*json.RawMessage
	err := json.Unmarshal(body, &m)
	if err != nil {
		return err
	}
	var v *json.RawMessage

	v = m["properties"]
	if v != nil {
		var properties Properties
		err = json.Unmarshal(*m["properties"], &properties)
		if err != nil {
			return err
		}
		cs.Properties = &properties
	}

	v = m["id"]
	if v != nil {
		var ID string
		err = json.Unmarshal(*m["id"], &ID)
		if err != nil {
			return err
		}
		cs.ID = &ID
	}

	v = m["name"]
	if v != nil {
		var name string
		err = json.Unmarshal(*m["name"], &name)
		if err != nil {
			return err
		}
		cs.Name = &name
	}

	v = m["type"]
	if v != nil {
		var typeVar string
		err = json.Unmarshal(*m["type"], &typeVar)
		if err != nil {
			return err
		}
		cs.Type = &typeVar
	}

	v = m["location"]
	if v != nil {
		var location string
		err = json.Unmarshal(*m["location"], &location)
		if err != nil {
			return err
		}
		cs.Location = &location
	}

	v = m["tags"]
	if v != nil {
		var tags map[string]*string
		err = json.Unmarshal(*m["tags"], &tags)
		if err != nil {
			return err
		}
		cs.Tags = &tags
	}

	return nil
}

// ContainerServicesCreateOrUpdateFuture an abstraction for monitoring and retrieving the results of a long-running
// operation.
type ContainerServicesCreateOrUpdateFuture struct {
	azure.Future
	req *http.Request
}

// Result returns the result of the asynchronous operation.
// If the operation has not completed it will return an error.
func (future ContainerServicesCreateOrUpdateFuture) Result(client ContainerServicesClient) (cs ContainerService, err error) {
	var done bool
	done, err = future.Done(client)
	if err != nil {
		return
	}
	if !done {
		return cs, autorest.NewError("containerservice.ContainerServicesCreateOrUpdateFuture", "Result", "asynchronous operation has not completed")
	}
	if future.PollingMethod() == azure.PollingLocation {
		cs, err = client.CreateOrUpdateResponder(future.Response())
		return
	}
	var resp *http.Response
	resp, err = autorest.SendWithSender(client, autorest.ChangeToGet(future.req),
		autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
	if err != nil {
		return
	}
	cs, err = client.CreateOrUpdateResponder(resp)
	return
}

// ContainerServicesDeleteFuture an abstraction for monitoring and retrieving the results of a long-running operation.
type ContainerServicesDeleteFuture struct {
	azure.Future
	req *http.Request
}

// Result returns the result of the asynchronous operation.
// If the operation has not completed it will return an error.
func (future ContainerServicesDeleteFuture) Result(client ContainerServicesClient) (ar autorest.Response, err error) {
	var done bool
	done, err = future.Done(client)
	if err != nil {
		return
	}
	if !done {
		return ar, autorest.NewError("containerservice.ContainerServicesDeleteFuture", "Result", "asynchronous operation has not completed")
	}
	if future.PollingMethod() == azure.PollingLocation {
		ar, err = client.DeleteResponder(future.Response())
		return
	}
	var resp *http.Response
	resp, err = autorest.SendWithSender(client, autorest.ChangeToGet(future.req),
		autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
	if err != nil {
		return
	}
	ar, err = client.DeleteResponder(resp)
	return
}

// DiagnosticsProfile ...
type DiagnosticsProfile struct {
	// VMDiagnostics - Profile for the container service VM diagnostic agent.
	VMDiagnostics *VMDiagnostics `json:"vmDiagnostics,omitempty"`
}

// LinuxProfile profile for Linux VMs in the container service cluster.
type LinuxProfile struct {
	// AdminUsername - The administrator username to use for all Linux VMs
	AdminUsername *string `json:"adminUsername,omitempty"`
	// SSH - The ssh key configuration for Linux VMs.
	SSH *SSHConfiguration `json:"ssh,omitempty"`
}

// ListResult the response from the List Container Services operation.
type ListResult struct {
	autorest.Response `json:"-"`
	// Value - the list of container services.
	Value *[]ContainerService `json:"value,omitempty"`
}

// MasterProfile profile for the container service master.
type MasterProfile struct {
	// Count - Number of masters (VMs) in the container service cluster. Allowed values are 1, 3, and 5. The default value is 1.
	Count *int32 `json:"count,omitempty"`
	// DNSPrefix - DNS prefix to be used to create the FQDN for master.
	DNSPrefix *string `json:"dnsPrefix,omitempty"`
	// Fqdn - FDQN for the master.
	Fqdn *string `json:"fqdn,omitempty"`
}

// OrchestratorProfile profile for the container service orchestrator.
type OrchestratorProfile struct {
	// OrchestratorType - The orchestrator to use to manage container service cluster resources. Valid values are Swarm, DCOS, and Custom. Possible values include: 'Swarm', 'DCOS'
	OrchestratorType OchestratorTypes `json:"orchestratorType,omitempty"`
}

// Properties properties of the container service.
type Properties struct {
	// ProvisioningState - the current deployment or provisioning state, which only appears in the response.
	ProvisioningState *string `json:"provisioningState,omitempty"`
	// OrchestratorProfile - Properties of the orchestrator.
	OrchestratorProfile *OrchestratorProfile `json:"orchestratorProfile,omitempty"`
	// MasterProfile - Properties of master agents.
	MasterProfile *MasterProfile `json:"masterProfile,omitempty"`
	// AgentPoolProfiles - Properties of the agent pool.
	AgentPoolProfiles *[]AgentPoolProfile `json:"agentPoolProfiles,omitempty"`
	// WindowsProfile - Properties of Windows VMs.
	WindowsProfile *WindowsProfile `json:"windowsProfile,omitempty"`
	// LinuxProfile - Properties of Linux VMs.
	LinuxProfile *LinuxProfile `json:"linuxProfile,omitempty"`
	// DiagnosticsProfile - Properties of the diagnostic agent.
	DiagnosticsProfile *DiagnosticsProfile `json:"diagnosticsProfile,omitempty"`
}

// Resource the Resource model definition.
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

// SSHConfiguration SSH configuration for Linux-based VMs running on Azure.
type SSHConfiguration struct {
	// PublicKeys - the list of SSH public keys used to authenticate with Linux-based VMs.
	PublicKeys *[]SSHPublicKey `json:"publicKeys,omitempty"`
}

// SSHPublicKey contains information about SSH certificate public key data.
type SSHPublicKey struct {
	// KeyData - Certificate public key used to authenticate with VMs through SSH. The certificate must be in PEM format with or without headers.
	KeyData *string `json:"keyData,omitempty"`
}

// VMDiagnostics profile for diagnostics on the container service VMs.
type VMDiagnostics struct {
	// Enabled - Whether the VM diagnostic agent is provisioned on the VM.
	Enabled *bool `json:"enabled,omitempty"`
	// StorageURI - The URI of the storage account where diagnostics are stored.
	StorageURI *string `json:"storageUri,omitempty"`
}

// WindowsProfile profile for Windows VMs in the container service cluster.
type WindowsProfile struct {
	// AdminUsername - The administrator username to use for Windows VMs
	AdminUsername *string `json:"adminUsername,omitempty"`
	// AdminPassword - The administrator password to use for Windows VMs
	AdminPassword *string `json:"adminPassword,omitempty"`
}

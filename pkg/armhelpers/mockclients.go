package armhelpers

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/arm/compute"
	"github.com/Azure/azure-sdk-for-go/arm/resources/resources"
	"github.com/Azure/go-autorest/autorest"
)

//MockACSEngineClient is an implemetnation of ACSEngineClient where all requests error out
type MockACSEngineClient struct {
	FailDeployTemplate         bool
	FailEnsureResourceGroup    bool
	FailListVirtualMachines    bool
	FailGetVirtualMachine      bool
	FailDeleteVirtualMachine   bool
	FailGetStorageClient       bool
	FailDeleteNetworkInterface bool
}

//DeployTemplate mock
func (mc *MockACSEngineClient) DeployTemplate(resourceGroup, name string, template, parameters map[string]interface{}, cancel <-chan struct{}) (*resources.DeploymentExtended, error) {
	if mc.FailDeployTemplate {
		return nil, fmt.Errorf("DeployTemplate failed")
	}

	return nil, nil
}

//EnsureResourceGroup mock
func (mc *MockACSEngineClient) EnsureResourceGroup(resourceGroup, location string) (*resources.Group, error) {
	if mc.FailEnsureResourceGroup {
		return nil, fmt.Errorf("EnsureResourceGroup failed")
	}

	return nil, nil
}

//ListVirtualMachines mock
func (mc *MockACSEngineClient) ListVirtualMachines(resourceGroup string) (compute.VirtualMachineListResult, error) {
	if mc.FailListVirtualMachines {
		return compute.VirtualMachineListResult{}, fmt.Errorf("ListVirtualMachines failed")
	}

	vm1Name := "k8s-master-12345678-0"

	creationSourceString := "creationSource"
	orchestratorString := "orchestrator"
	resourceNameSuffixString := "resourceNameSuffix"

	creationSource := "acsengine-k8s-master-12345678-0"
	orchestrator := "Kubernetes:1.5.3"
	resourceNameSuffix := "12345678"

	tags := map[string]*string{
		creationSourceString:     &creationSource,
		orchestratorString:       &orchestrator,
		resourceNameSuffixString: &resourceNameSuffix,
	}

	vm1 := compute.VirtualMachine{
		Name: &vm1Name,
		Tags: &tags,
		VirtualMachineProperties: &compute.VirtualMachineProperties{
			StorageProfile: &compute.StorageProfile{
				OsDisk: &compute.OSDisk{
					Vhd: &compute.VirtualHardDisk{
						URI: &validOSDiskResourceName},
				},
			},
			NetworkProfile: &compute.NetworkProfile{
				NetworkInterfaces: &[]compute.NetworkInterfaceReference{
					{
						ID: &validNicResourceName,
					},
				},
			},
		},
	}

	vmr := compute.VirtualMachineListResult{}
	vmr.Value = &[]compute.VirtualMachine{vm1}

	return vmr, nil
}

//GetVirtualMachine mock
func (mc *MockACSEngineClient) GetVirtualMachine(resourceGroup, name string) (compute.VirtualMachine, error) {
	if mc.FailGetVirtualMachine {
		return compute.VirtualMachine{}, fmt.Errorf("GetVirtualMachine failed")
	}

	vm1Name := "k8s-master-12345678-0"

	creationSourceString := "creationSource"
	orchestratorString := "orchestrator"
	resourceNameSuffixString := "resourceNameSuffix"

	creationSource := "acsengine-k8s-master-12345678-0"
	orchestrator := "Kubernetes:1.5.3"
	resourceNameSuffix := "12345678"

	tags := map[string]*string{
		creationSourceString:     &creationSource,
		orchestratorString:       &orchestrator,
		resourceNameSuffixString: &resourceNameSuffix,
	}

	return compute.VirtualMachine{
		Name: &vm1Name,
		Tags: &tags,
		VirtualMachineProperties: &compute.VirtualMachineProperties{
			StorageProfile: &compute.StorageProfile{
				OsDisk: &compute.OSDisk{
					Vhd: &compute.VirtualHardDisk{
						URI: &validOSDiskResourceName},
				},
			},
			NetworkProfile: &compute.NetworkProfile{
				NetworkInterfaces: &[]compute.NetworkInterfaceReference{
					{
						ID: &validNicResourceName,
					},
				},
			},
		},
	}, nil
}

//DeleteVirtualMachine mock
func (mc *MockACSEngineClient) DeleteVirtualMachine(resourceGroup, name string, cancel <-chan struct{}) (<-chan compute.OperationStatusResponse, <-chan error) {
	if mc.FailDeleteVirtualMachine {
		errChan := make(chan error)
		respChan := make(chan compute.OperationStatusResponse)
		go func() {
			defer func() {
				close(errChan)
			}()
			defer func() {
				close(respChan)
			}()
			errChan <- fmt.Errorf("DeleteVirtualMachine failed")
			time.Sleep(1 * time.Second)
		}()
		return respChan, errChan
	}

	errChan := make(chan error)
	respChan := make(chan compute.OperationStatusResponse)
	go func() {
		defer func() {
			close(errChan)
		}()
		defer func() {
			close(respChan)
		}()
		errChan <- nil
		respChan <- compute.OperationStatusResponse{}
		time.Sleep(1 * time.Second)
	}()
	return respChan, errChan
}

//GetStorageClient mock
func (mc *MockACSEngineClient) GetStorageClient(resourceGroup, accountName string) (ACSStorageClient, error) {
	if mc.FailGetStorageClient {
		return nil, fmt.Errorf("GetStorageClient failed")
	}

	return &MockStorageClient{}, nil
}

//DeleteNetworkInterface mock
func (mc *MockACSEngineClient) DeleteNetworkInterface(resourceGroup, nicName string, cancel <-chan struct{}) (<-chan autorest.Response, <-chan error) {
	if mc.FailDeleteNetworkInterface {
		errChan := make(chan error)
		respChan := make(chan autorest.Response)
		go func() {
			defer func() {
				close(errChan)
			}()
			defer func() {
				close(respChan)
			}()
			errChan <- fmt.Errorf("DeleteNetworkInterface failed")
			time.Sleep(1 * time.Second)
		}()
		return respChan, errChan
	}

	errChan := make(chan error)
	respChan := make(chan autorest.Response)
	go func() {
		defer func() {
			close(errChan)
		}()
		defer func() {
			close(respChan)
		}()
		errChan <- nil
		respChan <- autorest.Response{}
		time.Sleep(1 * time.Second)
	}()
	return respChan, errChan
}

var validOSDiskResourceName = "https://00k71r4u927seqiagnt0.blob.core.windows.net/osdisk/k8s-agentpool1-12345678-0-osdisk.vhd"
var validNicResourceName = "/subscriptions/DEC923E3-1EF1-4745-9516-37906D56DEC4/resourceGroups/acsK8sTest/providers/Microsoft.Network/networkInterfaces/k8s-agent-12345678-nic-0"

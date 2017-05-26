package armhelpers

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/arm/compute"
	"github.com/Azure/azure-sdk-for-go/arm/resources/resources"
	"github.com/Azure/go-autorest/autorest"
)

//FailingMockClient is an implemetnation of ACSEngineClient where all requests error out
type FailingMockClient struct{}

//DeployTemplate mock
func (fmc *FailingMockClient) DeployTemplate(resourceGroup, name string, template, parameters map[string]interface{}, cancel <-chan struct{}) (*resources.DeploymentExtended, error) {
	return nil, fmt.Errorf("failed")
}

//EnsureResourceGroup mock
func (fmc *FailingMockClient) EnsureResourceGroup(resourceGroup, location string) (*resources.Group, error) {
	return nil, fmt.Errorf("failed")
}

//ListVirtualMachines mock
func (fmc *FailingMockClient) ListVirtualMachines(resourceGroup string) (compute.VirtualMachineListResult, error) {
	return compute.VirtualMachineListResult{}, fmt.Errorf("failed")
}

//GetVirtualMachine mock
func (fmc *FailingMockClient) GetVirtualMachine(resourceGroup, name string) (compute.VirtualMachine, error) {
	return compute.VirtualMachine{}, fmt.Errorf("failed")

}

//DeleteVirtualMachine mock
func (fmc *FailingMockClient) DeleteVirtualMachine(resourceGroup, name string, cancel <-chan struct{}) (<-chan compute.OperationStatusResponse, <-chan error) {
	errChan := make(chan error)
	respChan := make(chan compute.OperationStatusResponse)
	go func() {
		defer func() {
			close(errChan)
		}()
		defer func() {
			close(respChan)
		}()
		errChan <- fmt.Errorf("failed")
		time.Sleep(1 * time.Second)
	}()
	return respChan, errChan
}

//GetStorageClient mock
func (fmc *FailingMockClient) GetStorageClient(resourceGroup, accountName string) (ACSStorageClient, error) {
	return nil, fmt.Errorf("failed")
}

//DeleteNetworkInterface mock
func (fmc *FailingMockClient) DeleteNetworkInterface(resourceGroup, nicName string, cancel <-chan struct{}) (<-chan autorest.Response, <-chan error) {
	errChan := make(chan error)
	respChan := make(chan autorest.Response)
	go func() {
		defer func() {
			close(errChan)
		}()
		defer func() {
			close(respChan)
		}()
		errChan <- fmt.Errorf("failed")
		time.Sleep(1 * time.Second)
	}()
	return respChan, errChan
}

//MockClient is an implementation of ACSEngineClient where all requests return a valid response
type MockClient struct{}

//DeployTemplate mock
func (mc *MockClient) DeployTemplate(resourceGroup, name string, template, parameters map[string]interface{}, cancel <-chan struct{}) (*resources.DeploymentExtended, error) {
	return nil, nil
}

//EnsureResourceGroup mock
func (mc *MockClient) EnsureResourceGroup(resourceGroup, location string) (*resources.Group, error) {
	return nil, nil
}

//ListVirtualMachines mock
func (mc *MockClient) ListVirtualMachines(resourceGroup string) (compute.VirtualMachineListResult, error) {
	return compute.VirtualMachineListResult{}, nil
}

var validOsDiskURI = "https://osdisk.storage.com/container/blob/disk.vhd"
var validNicID = "/subscriptions/subid/resourceGroups/acs-k8s-int/providers/Microsoft.Network/networkInterfaces/k8s-agent-F8EADCCF-nic-0"

//GetVirtualMachine mock
func (mc *MockClient) GetVirtualMachine(resourceGroup, name string) (compute.VirtualMachine, error) {
	return compute.VirtualMachine{
		VirtualMachineProperties: &compute.VirtualMachineProperties{
			StorageProfile: &compute.StorageProfile{
				OsDisk: &compute.OSDisk{
					Vhd: &compute.VirtualHardDisk{
						URI: &validOsDiskURI},
				},
			},
			NetworkProfile: &compute.NetworkProfile{
				NetworkInterfaces: &[]compute.NetworkInterfaceReference{
					compute.NetworkInterfaceReference{
						ID: &validNicID,
					},
				},
			},
		},
	}, nil

}

//DeleteVirtualMachine mock
func (mc *MockClient) DeleteVirtualMachine(resourceGroup, name string, cancel <-chan struct{}) (<-chan compute.OperationStatusResponse, <-chan error) {
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
func (mc *MockClient) GetStorageClient(resourceGroup, accountName string) (ACSStorageClient, error) {
	return &MockStorageClient{}, nil
}

//DeleteNetworkInterface mock
func (mc *MockClient) DeleteNetworkInterface(resourceGroup, nicName string, cancel <-chan struct{}) (<-chan autorest.Response, <-chan error) {
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

//MockStorageClient mock implementation of StorageClient
type MockStorageClient struct{}

//DeleteBlob mock
func (msc *MockStorageClient) DeleteBlob(container, blob string) error {
	return nil
}

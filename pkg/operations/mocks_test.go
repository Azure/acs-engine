package operations

import (
	"fmt"
	"time"

	"github.com/Azure/acs-engine/pkg/armhelpers"
	"github.com/Azure/azure-sdk-for-go/arm/compute"
	"github.com/Azure/azure-sdk-for-go/arm/resources/resources"
	"github.com/Azure/go-autorest/autorest"
)

type failingMockClient struct{}

func (fmc *failingMockClient) DeployTemplate(resourceGroup, name string, template, parameters map[string]interface{}, cancel <-chan struct{}) (*resources.DeploymentExtended, error) {
	return nil, fmt.Errorf("failed")
}

func (fmc *failingMockClient) EnsureResourceGroup(resourceGroup, location string) (*resources.Group, error) {
	return nil, fmt.Errorf("failed")
}

func (fmc *failingMockClient) ListVirtualMachines(resourceGroup string) (compute.VirtualMachineListResult, error) {
	return compute.VirtualMachineListResult{}, fmt.Errorf("failed")
}

func (fmc *failingMockClient) GetVirtualMachine(resourceGroup, name string) (compute.VirtualMachine, error) {
	return compute.VirtualMachine{}, fmt.Errorf("failed")

}

func (fmc *failingMockClient) DeleteVirtualMachine(resourceGroup, name string, cancel <-chan struct{}) (<-chan compute.OperationStatusResponse, <-chan error) {
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

func (fmc *failingMockClient) GetStorageClient(resourceGroup, accountName string) (armhelpers.ACSStorageClient, error) {
	return nil, fmt.Errorf("failed")
}

func (fmc *failingMockClient) DeleteNetworkInterface(resourceGroup, nicName string, cancel <-chan struct{}) (<-chan autorest.Response, <-chan error) {
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

type mockClient struct{}

func (mc *mockClient) DeployTemplate(resourceGroup, name string, template, parameters map[string]interface{}, cancel <-chan struct{}) (*resources.DeploymentExtended, error) {
	return nil, nil
}

func (mc *mockClient) EnsureResourceGroup(resourceGroup, location string) (*resources.Group, error) {
	return nil, nil
}

func (mc *mockClient) ListVirtualMachines(resourceGroup string) (compute.VirtualMachineListResult, error) {
	return compute.VirtualMachineListResult{}, nil
}

var validOsDiskURI = "https://osdisk.storage.com/container/blob/disk.vhd"
var validNicID = "/subscriptions/subid/resourceGroups/acs-k8s-int/providers/Microsoft.Network/networkInterfaces/k8s-agent-F8EADCCF-nic-0"

func (mc *mockClient) GetVirtualMachine(resourceGroup, name string) (compute.VirtualMachine, error) {
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

func (mc *mockClient) DeleteVirtualMachine(resourceGroup, name string, cancel <-chan struct{}) (<-chan compute.OperationStatusResponse, <-chan error) {
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

func (mc *mockClient) GetStorageClient(resourceGroup, accountName string) (armhelpers.ACSStorageClient, error) {
	return &mockStorageClient{}, nil
}

func (mc *mockClient) DeleteNetworkInterface(resourceGroup, nicName string, cancel <-chan struct{}) (<-chan autorest.Response, <-chan error) {
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

type mockStorageClient struct{}

func (msc *mockStorageClient) DeleteBlob(container, blob string) error {
	return nil
}

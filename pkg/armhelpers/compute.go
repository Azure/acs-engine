package armhelpers

import (
	"github.com/Azure/azure-sdk-for-go/arm/compute"
	"github.com/Azure/go-autorest/autorest"
)

// ListVirtualMachines returns (the first page of) the machines in the specified resource group.
func (az *AzureClient) ListVirtualMachines(resourceGroup string) (result compute.VirtualMachineListResult, err error) {
	client := az.virtualMachinesClient
	req, err := client.ListPreparer(resourceGroup)
	if err != nil {
		err = autorest.NewErrorWithError(err, "compute.VirtualMachinesClient", "List", nil, "Failure preparing request")
		return
	}
	az.addAcceptLanguages(req)

	resp, err := client.ListSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "compute.VirtualMachinesClient", "List", resp, "Failure sending request")
		return
	}

	result, err = client.ListResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "compute.VirtualMachinesClient", "List", resp, "Failure responding to request")
	}

	return
}

// GetVirtualMachine returns the specified machine in the specified resource group.
func (az *AzureClient) GetVirtualMachine(resourceGroup, name string) (result compute.VirtualMachine, err error) {
	client := az.virtualMachinesClient
	req, err := client.GetPreparer(resourceGroup, name, "")
	if err != nil {
		err = autorest.NewErrorWithError(err, "compute.VirtualMachinesClient", "Get", nil, "Failure preparing request")
		return
	}
	az.addAcceptLanguages(req)

	resp, err := client.GetSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "compute.VirtualMachinesClient", "Get", resp, "Failure sending request")
		return
	}

	result, err = client.GetResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "compute.VirtualMachinesClient", "Get", resp, "Failure responding to request")
	}

	return
}

// DeleteVirtualMachine handles deletion of a CRP/VMAS VM (aka, not a VMSS VM).
func (az *AzureClient) DeleteVirtualMachine(resourceGroup, name string, cancel <-chan struct{}) (<-chan compute.OperationStatusResponse, <-chan error) {
	client := az.virtualMachinesClient
	resultChan := make(chan compute.OperationStatusResponse, 1)
	errChan := make(chan error, 1)
	go func() {
		var err error
		var result compute.OperationStatusResponse
		defer func() {
			resultChan <- result
			errChan <- err
			close(resultChan)
			close(errChan)
		}()
		req, err := client.DeletePreparer(resourceGroup, name, cancel)
		if err != nil {
			err = autorest.NewErrorWithError(err, "compute.VirtualMachinesClient", "Delete", nil, "Failure preparing request")
			return
		}
		az.addAcceptLanguages(req)

		resp, err := client.DeleteSender(req)
		if err != nil {
			result.Response = autorest.Response{Response: resp}
			err = autorest.NewErrorWithError(err, "compute.VirtualMachinesClient", "Delete", resp, "Failure sending request")
			return
		}

		result, err = client.DeleteResponder(resp)
		if err != nil {
			err = autorest.NewErrorWithError(err, "compute.VirtualMachinesClient", "Delete", resp, "Failure responding to request")
		}
	}()
	return resultChan, errChan
}

// ListVirtualMachineScaleSets returns (the first page of) the vmss resources in the specified resource group.
func (az *AzureClient) ListVirtualMachineScaleSets(resourceGroup string) (result compute.VirtualMachineScaleSetListResult, err error) {
	client := az.virtualMachineScaleSetsClient
	req, err := client.ListPreparer(resourceGroup)
	if err != nil {
		err = autorest.NewErrorWithError(err, "compute.VirtualMachineScaleSetsClient", "List", nil, "Failure preparing request")
		return
	}
	az.addAcceptLanguages(req)

	resp, err := client.ListSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "compute.VirtualMachineScaleSetsClient", "List", resp, "Failure sending request")
		return
	}

	result, err = client.ListResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "compute.VirtualMachineScaleSetsClient", "List", resp, "Failure responding to request")
	}

	return
}

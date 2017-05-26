package operations

import (
	"container/list"

	"github.com/Azure/acs-engine/pkg/armhelpers"
)

//VMScalingErrorDetails give the index in the agent pool that failed and the accompanying error
type VMScalingErrorDetails struct {
	Name  string
	Error error
}

// ScaleDownVMs removes the vms in the provided list. Returns a list with details on each failure.
// all items in the list will always be of type *VMScalingErrorDetails
func ScaleDownVMs(az armhelpers.ACSEngineClient, resourceGroup string, vmNames ...string) *list.List {
	numVmsToDelete := len(vmNames)
	errChan := make(chan *VMScalingErrorDetails, numVmsToDelete)
	defer close(errChan)
	for _, vmName := range vmNames {
		go func(vmName string) {
			err := CleanDeleteVirtualMachine(az, resourceGroup, vmName)
			if err != nil {
				errChan <- &VMScalingErrorDetails{Name: vmName, Error: err}
				return
			}
			errChan <- nil
		}(vmName)
	}
	failedVMDeletions := &list.List{}
	for i := 0; i < numVmsToDelete; i++ {
		errDetails := <-errChan
		if errDetails != nil {
			failedVMDeletions.PushBack(errDetails)
		}
	}
	if failedVMDeletions.Len() > 0 {
		return failedVMDeletions
	}
	return nil
}

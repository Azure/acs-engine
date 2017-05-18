package operations

import (
	"fmt"

	"github.com/Azure/acs-engine/pkg/armhelpers"
	"github.com/prometheus/common/log"
)

// CleanDeleteVirtualMachine deletes a VM and any associated OS disk
func CleanDeleteVirtualMachine(az armhelpers.ACSEngineClient, resourceGroup, name string) error {
	log.Infof("fetching VM: %s/%s", resourceGroup, name)
	vm, err := az.GetVirtualMachine(resourceGroup, name)
	if err != nil {
		log.Errorf("failed to get VM: %s/%s: %s", resourceGroup, name, err.Error())
		return err
	}

	// NOTE: This code assumes a non-managed disk!
	vhd := vm.VirtualMachineProperties.StorageProfile.OsDisk.Vhd
	if vhd == nil {
		log.Warnf("found an OS Disk with no VHD URI. This is probably a VM with a managed disk")
		return fmt.Errorf("os disk does not have a VHD URI")
	}
	accountName, vhdContainer, vhdBlob, err := armhelpers.SplitBlobURI(*vhd.URI)
	if err != nil {
		return err
	}

	nicID := (*vm.VirtualMachineProperties.NetworkProfile.NetworkInterfaces)[0].ID
	nicName, err := armhelpers.ResourceName(*nicID)
	if err != nil {
		return err
	}

	log.Infof("found os disk storage reference: %s %s %s", accountName, vhdContainer, vhdBlob)
	log.Infof("found nic name for VM (%s/%s): %s", resourceGroup, name, nicName)

	log.Infof("deleting VM: %s/%s", resourceGroup, name)
	_, deleteErrChan := az.DeleteVirtualMachine(resourceGroup, name, nil)

	as, err := az.GetStorageClient(resourceGroup, accountName)
	if err != nil {
		return err
	}

	log.Infof("waiting for vm deletion: %s/%s", resourceGroup, name)
	if err := <-deleteErrChan; err != nil {
		return err
	}

	log.Infof("deleting nic: %s/%s", resourceGroup, nicName)
	_, nicErrChan := az.DeleteNetworkInterface(resourceGroup, nicName, nil)
	if err != nil {
		return err
	}

	log.Infof("deleting blob: %s/%s", vhdContainer, vhdBlob)
	if err = as.DeleteBlob(vhdContainer, vhdBlob); err != nil {
		return err
	}

	log.Infof("waiting for nic deletion: %s/%s", resourceGroup, nicName)
	if nicErr := <-nicErrChan; nicErr != nil {
		return nicErr
	}

	return nil
}

package operations

import (
	"fmt"

	"github.com/Azure/acs-engine/pkg/armhelpers"
	log "github.com/Sirupsen/logrus"
)

// CleanDeleteVirtualMachine deletes a VM and any associated OS disk
func CleanDeleteVirtualMachine(az armhelpers.ACSEngineClient, logger *log.Entry, resourceGroup, name string) error {
	logger.Infof("fetching VM: %s/%s", resourceGroup, name)
	vm, err := az.GetVirtualMachine(resourceGroup, name)
	if err != nil {
		logger.Errorf("failed to get VM: %s/%s: %s", resourceGroup, name, err.Error())
		return err
	}

	// NOTE: This code assumes a non-managed disk!
	vhd := vm.VirtualMachineProperties.StorageProfile.OsDisk.Vhd
	if vhd == nil {
		logger.Warnf("found an OS Disk with no VHD URI. This is probably a VM with a managed disk")
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

	logger.Infof("found os disk storage reference: %s %s %s", accountName, vhdContainer, vhdBlob)
	logger.Infof("found nic name for VM (%s/%s): %s", resourceGroup, name, nicName)

	logger.Infof("deleting VM: %s/%s", resourceGroup, name)
	_, deleteErrChan := az.DeleteVirtualMachine(resourceGroup, name, nil)

	as, err := az.GetStorageClient(resourceGroup, accountName)
	if err != nil {
		return err
	}

	logger.Infof("waiting for vm deletion: %s/%s", resourceGroup, name)
	if err := <-deleteErrChan; err != nil {
		return err
	}

	logger.Infof("deleting nic: %s/%s", resourceGroup, nicName)
	_, nicErrChan := az.DeleteNetworkInterface(resourceGroup, nicName, nil)
	if err != nil {
		return err
	}

	logger.Infof("deleting blob: %s/%s", vhdContainer, vhdBlob)
	if err = as.DeleteBlob(vhdContainer, vhdBlob); err != nil {
		return err
	}

	logger.Infof("waiting for nic deletion: %s/%s", resourceGroup, nicName)
	if nicErr := <-nicErrChan; nicErr != nil {
		return nicErr
	}

	return nil
}

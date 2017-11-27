package operations

import (
	"fmt"

	"github.com/Azure/acs-engine/pkg/armhelpers"
	log "github.com/sirupsen/logrus"
)

// CleanDeleteVirtualMachine deletes a VM and any associated OS disk
func CleanDeleteVirtualMachine(az armhelpers.ACSEngineClient, logger *log.Entry, resourceGroup, name string) error {
	logger.Infof("fetching VM: %s/%s", resourceGroup, name)
	vm, err := az.GetVirtualMachine(resourceGroup, name)
	if err != nil {
		logger.Errorf("failed to get VM: %s/%s: %s", resourceGroup, name, err.Error())
		return err
	}

	vhd := vm.VirtualMachineProperties.StorageProfile.OsDisk.Vhd
	managedDisk := vm.VirtualMachineProperties.StorageProfile.OsDisk.ManagedDisk
	if vhd == nil && managedDisk == nil {
		logger.Errorf("failed to get a valid os disk URI for VM: %s/%s", resourceGroup, name)

		return fmt.Errorf("os disk does not have a VHD URI")
	}

	osDiskName := vm.VirtualMachineProperties.StorageProfile.OsDisk.Name

	var nicName string
	nicID := (*vm.VirtualMachineProperties.NetworkProfile.NetworkInterfaces)[0].ID
	if nicID == nil {
		logger.Warnf("NIC ID is not set for VM (%s/%s)", resourceGroup, name)
	} else {
		nicName, err = armhelpers.ResourceName(*nicID)
		if err != nil {
			return err
		}
		logger.Infof("found nic name for VM (%s/%s): %s", resourceGroup, name, nicName)
	}
	logger.Infof("deleting VM: %s/%s", resourceGroup, name)
	_, deleteErrChan := az.DeleteVirtualMachine(resourceGroup, name, nil)

	logger.Infof("waiting for vm deletion: %s/%s", resourceGroup, name)
	if err := <-deleteErrChan; err != nil {
		return err
	}

	if len(nicName) > 0 {
		logger.Infof("deleting nic: %s/%s", resourceGroup, nicName)
		_, nicErrChan := az.DeleteNetworkInterface(resourceGroup, nicName, nil)

		logger.Infof("waiting for nic deletion: %s/%s", resourceGroup, nicName)
		if nicErr := <-nicErrChan; nicErr != nil {
			return nicErr
		}
	}

	if vhd != nil {
		accountName, vhdContainer, vhdBlob, err := armhelpers.SplitBlobURI(*vhd.URI)
		if err != nil {
			return err
		}

		logger.Infof("found os disk storage reference: %s %s %s", accountName, vhdContainer, vhdBlob)

		as, err := az.GetStorageClient(resourceGroup, accountName)
		if err != nil {
			return err
		}

		logger.Infof("deleting blob: %s/%s", vhdContainer, vhdBlob)
		if err = as.DeleteBlob(vhdContainer, vhdBlob); err != nil {
			return err
		}
	} else if managedDisk != nil {
		if osDiskName == nil {
			logger.Warnf("osDisk is not set for VM %s/%s", resourceGroup, name)
		} else {
			logger.Infof("deleting managed disk: %s/%s", resourceGroup, *osDiskName)
			_, diskErrChan := az.DeleteManagedDisk(resourceGroup, *osDiskName, nil)

			if err := <-diskErrChan; err != nil {
				return err
			}
		}
	}

	return nil
}

package operations

import (
	"fmt"

	"github.com/Azure/acs-engine/pkg/armhelpers"
	"github.com/Azure/acs-engine/pkg/armhelpers/utils"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

const (
	// AADRoleResourceGroupScopeTemplate is a template for a roleDefinition scope
	AADRoleResourceGroupScopeTemplate = "/subscriptions/%s/resourceGroups/%s"
)

// CleanDeleteVirtualMachine deletes a VM and any associated OS disk
func CleanDeleteVirtualMachine(az armhelpers.ACSEngineClient, logger *log.Entry, subscriptionID, resourceGroup, name string) error {
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

		return errors.New("os disk does not have a VHD URI")
	}

	osDiskName := vm.VirtualMachineProperties.StorageProfile.OsDisk.Name

	var nicName string
	nicID := (*vm.VirtualMachineProperties.NetworkProfile.NetworkInterfaces)[0].ID
	if nicID == nil {
		logger.Warnf("NIC ID is not set for VM (%s/%s)", resourceGroup, name)
	} else {
		nicName, err = utils.ResourceName(*nicID)
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
		accountName, vhdContainer, vhdBlob, err := utils.SplitBlobURI(*vhd.URI)
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

	if vm.Identity != nil {
		// Role assignments are not deleted if the VM is destroyed, so we must cleanup ourselves!
		// The role assignments should only be relevant if managed identities are used,
		// but always cleaning them up is easier than adding rule based logic here and there.
		scope := fmt.Sprintf(AADRoleResourceGroupScopeTemplate, subscriptionID, resourceGroup)
		logger.Infof("fetching roleAssignments: %s with principal %s", scope, *vm.Identity.PrincipalID)
		vmRoleAssignments, listRoleAssingmentsError := az.ListRoleAssignmentsForPrincipal(scope, *vm.Identity.PrincipalID)
		if listRoleAssingmentsError != nil {
			logger.Errorf("failed to list role assignments: %s/%s: %s", scope, *vm.Identity.PrincipalID, listRoleAssingmentsError.Error())
			return listRoleAssingmentsError
		}

		for _, roleAssignment := range *vmRoleAssignments.Value {
			logger.Infof("deleting role assignment: %s", *roleAssignment.ID)
			_, deleteRoleAssignmentErr := az.DeleteRoleAssignmentByID(*roleAssignment.ID)
			if deleteRoleAssignmentErr != nil {
				logger.Errorf("failed to delete role assignment: %s: %s", *roleAssignment.ID, deleteRoleAssignmentErr.Error())
				return deleteRoleAssignmentErr
			}
		}
	}

	return nil
}

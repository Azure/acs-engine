package acsengine

import (
	"fmt"
	"log"
	"strings"

	"github.com/Azure/acs-engine/pkg/i18n"
	"github.com/Azure/azure-sdk-for-go/arm/compute"
	"github.com/Sirupsen/logrus"
)

const (
	//Field names
	customDataFieldName            = "customData"
	dependsOnFieldName             = "dependsOn"
	hardwareProfileFieldName       = "hardwareProfile"
	imageReferenceFieldName        = "imageReference"
	nameFieldName                  = "name"
	osProfileFieldName             = "osProfile"
	propertiesFieldName            = "properties"
	resourcesFieldName             = "resources"
	storageProfileFieldName        = "storageProfile"
	typeFieldName                  = "type"
	virtualMachineProfileFieldName = "virtualMachineProfile"
	vmSizeFieldName                = "vmSize"
	dataDisksFieldName             = "dataDisks"
	createOptionFieldName          = "createOption"
	tagsFieldName                  = "tags"
	managedDiskFieldName           = "managedDisk"

	// ARM resource Types
	nsgResourceType  = "Microsoft.Network/networkSecurityGroups"
	vmResourceType   = "Microsoft.Compute/virtualMachines"
	vmssResourceType = "Microsoft.Compute/virtualMachineScaleSets"
	vmExtensionType  = "Microsoft.Compute/virtualMachines/extensions"

	// resource ids
	nsgID = "nsgID"
)

// Transformer represents the object that transforms template
type Transformer struct {
	Translator *i18n.Translator
}

// NormalizeForVMSSScaling takes a template and removes elements that are unwanted in a VMSS scale up/down case
func (t *Transformer) NormalizeForVMSSScaling(logger *logrus.Entry, templateMap map[string]interface{}) error {
	if err := t.NormalizeMasterResourcesForScaling(logger, templateMap); err != nil {
		return err
	}

	resources := templateMap[resourcesFieldName].([]interface{})
	for _, resource := range resources {
		resourceMap, ok := resource.(map[string]interface{})
		if !ok {
			logger.Warnf("Template improperly formatted")
			continue
		}

		resourceType, ok := resourceMap[typeFieldName].(string)
		if !ok || resourceType != vmssResourceType {
			continue
		}

		resourceProperties, ok := resourceMap[propertiesFieldName].(map[string]interface{})
		if !ok {
			logger.Warnf("Template improperly formatted")
			continue
		}

		virtualMachineProfile, ok := resourceProperties[virtualMachineProfileFieldName].(map[string]interface{})
		if !ok {
			logger.Warnf("Template improperly formatted")
			continue
		}

		if !t.removeCustomData(logger, virtualMachineProfile) || !t.removeImageReference(logger, virtualMachineProfile) {
			continue
		}
	}
	return nil
}

// NormalizeForK8sVMASScalingUp takes a template and removes elements that are unwanted in a K8s VMAS scale up/down case
func (t *Transformer) NormalizeForK8sVMASScalingUp(logger *logrus.Entry, templateMap map[string]interface{}) error {
	if err := t.NormalizeMasterResourcesForScaling(logger, templateMap); err != nil {
		return err
	}
	nsgIndex := -1
	resources := templateMap[resourcesFieldName].([]interface{})
	for index, resource := range resources {
		resourceMap, ok := resource.(map[string]interface{})
		if !ok {
			logger.Warnf("Template improperly formatted for resource")
			continue
		}

		resourceType, ok := resourceMap[typeFieldName].(string)
		if ok && resourceType == nsgResourceType {
			if nsgIndex != -1 {
				err := t.Translator.Errorf("Found 2 resources with type %s in the template. There should only be 1", nsgResourceType)
				logger.Errorf(err.Error())
				return err
			}
			nsgIndex = index
		}

		dependencies, ok := resourceMap[dependsOnFieldName].([]interface{})
		if !ok {
			logger.Warnf("%s field not found for type: %s. Continue...", dependsOnFieldName, resourceType)
			continue
		}

		for dIndex := len(dependencies) - 1; dIndex >= 0; dIndex-- {
			dependency := dependencies[dIndex].(string)
			if strings.Contains(dependency, nsgResourceType) || strings.Contains(dependency, nsgID) {
				dependencies = append(dependencies[:dIndex], dependencies[dIndex+1:]...)
			}
		}

		resourceMap[dependsOnFieldName] = dependencies
	}
	if nsgIndex == -1 {
		err := t.Translator.Errorf("Found no resources with type %s in the template. There should have been 1", nsgResourceType)
		logger.Errorf(err.Error())
		return err
	}

	templateMap[resourcesFieldName] = append(resources[:nsgIndex], resources[nsgIndex+1:]...)

	return nil
}

// NormalizeMasterResourcesForScaling takes a template and removes elements that are unwanted in any scale up/down case
func (t *Transformer) NormalizeMasterResourcesForScaling(logger *logrus.Entry, templateMap map[string]interface{}) error {
	resources := templateMap[resourcesFieldName].([]interface{})
	//update master nodes resources
	for _, resource := range resources {
		resourceMap, ok := resource.(map[string]interface{})
		if !ok {
			logger.Warnf("Template improperly formatted")
			continue
		}

		resourceType, ok := resourceMap[typeFieldName].(string)
		if !ok || resourceType != vmResourceType {
			continue
		}

		resourceName, ok := resourceMap[nameFieldName].(string)
		if !ok {
			logger.Warnf("Template improperly formatted")
			continue
		}

		// make sure this is only modifying the master vms
		if !strings.Contains(resourceName, "variables('masterVMNamePrefix')") {
			continue
		}

		resourceProperties, ok := resourceMap[propertiesFieldName].(map[string]interface{})
		if !ok {
			logger.Warnf("Template improperly formatted")
			continue
		}

		hardwareProfile, ok := resourceProperties[hardwareProfileFieldName].(map[string]interface{})
		if !ok {
			logger.Warnf("Template improperly formatted")
			continue
		}

		if hardwareProfile[vmSizeFieldName] != nil {
			delete(hardwareProfile, vmSizeFieldName)
		}

		if !t.removeCustomData(logger, resourceProperties) || !t.removeImageReference(logger, resourceProperties) {
			continue
		}
	}

	return nil
}

func (t *Transformer) removeCustomData(logger *logrus.Entry, resourceProperties map[string]interface{}) bool {
	osProfile, ok := resourceProperties[osProfileFieldName].(map[string]interface{})
	if !ok {
		logger.Warnf("Template improperly formatted")
		return ok
	}

	if osProfile[customDataFieldName] != nil {
		delete(osProfile, customDataFieldName)
	}
	return ok
}

func (t *Transformer) removeImageReference(logger *logrus.Entry, resourceProperties map[string]interface{}) bool {
	storageProfile, ok := resourceProperties[storageProfileFieldName].(map[string]interface{})
	if !ok {
		logger.Warnf("Template improperly formatted. Could not find: %s", storageProfileFieldName)
		return ok
	}

	if storageProfile[imageReferenceFieldName] != nil {
		delete(storageProfile, imageReferenceFieldName)
	}
	return ok
}

// NormalizeResourcesForK8sMasterUpgrade takes a template and removes elements that are unwanted in any scale up/down case
func (t *Transformer) NormalizeResourcesForK8sMasterUpgrade(logger *logrus.Entry, templateMap map[string]interface{}, isMasterManagedDisk bool, agentPoolsToPreserve map[string]bool) error {
	resources := templateMap[resourcesFieldName].([]interface{})
	logger.Infoln(fmt.Sprintf("Resource count before running NormalizeResourcesForK8sMasterUpgrade: %d", len(resources)))

	filteredResources := resources[:0]

	// remove agent nodes resources if needed and set dataDisk createOption to attach
	for _, resource := range resources {
		filteredResources = append(filteredResources, resource)
		resourceMap, ok := resource.(map[string]interface{})
		if !ok {
			logger.Warnf("Template improperly formatted for field name: %s", resourcesFieldName)
			continue
		}

		resourceType, ok := resourceMap[typeFieldName].(string)
		if !ok {
			continue
		}

		if !(resourceType == vmResourceType || resourceType == vmExtensionType) {
			continue
		}

		resourceName, ok := resourceMap[nameFieldName].(string)
		if !ok {
			logger.Warnf("Template improperly formatted for field name: %s", nameFieldName)
			continue
		}

		if strings.EqualFold(resourceType, vmResourceType) &&
			strings.Contains(resourceName, "variables('masterVMNamePrefix')") {
			resourceProperties, ok := resourceMap[propertiesFieldName].(map[string]interface{})
			if !ok {
				logger.Warnf("Template improperly formatted for field name: %s, resource name: %s", propertiesFieldName, resourceName)
				continue
			}

			storageProfile, ok := resourceProperties[storageProfileFieldName].(map[string]interface{})
			if !ok {
				logger.Warnf("Template improperly formatted: %s", storageProfileFieldName)
				continue
			}

			dataDisks := storageProfile[dataDisksFieldName].([]interface{})
			dataDisk, _ := dataDisks[0].(map[string]interface{})
			dataDisk[createOptionFieldName] = "attach"

			if isMasterManagedDisk {
				managedDisk := compute.ManagedDiskParameters{}
				id := "[concat('/subscriptions/', variables('subscriptionId'), '/resourceGroups/', variables('resourceGroup'),'/providers/Microsoft.Compute/disks/', variables('masterVMNamePrefix'), copyIndex(variables('masterOffset')),'-etcddisk')]"
				managedDisk.ID = &id
				var diskInterface interface{}
				diskInterface = &managedDisk
				dataDisk[managedDiskFieldName] = diskInterface
			}
		}

		tags, _ := resourceMap[tagsFieldName].(map[string]interface{})
		poolName := fmt.Sprint(tags["poolName"]) // poolName tag exists on agents only

		if resourceType == vmResourceType {
			logger.Infoln(fmt.Sprintf("Evaluating if agent pool: %s, resource: %s needs to be removed", poolName, resourceName))
			// Not an agent (could be a master VM)
			if tags["poolName"] == nil || strings.Contains(resourceName, "variables('masterVMNamePrefix')") {
				continue
			}

			logger.Infoln(fmt.Sprintf("agentPoolsToPreserve: %v...", agentPoolsToPreserve))

			if agentPoolsToPreserve == nil || len(agentPoolsToPreserve) == 0 || agentPoolsToPreserve[poolName] != true {
				logger.Infoln(fmt.Sprintf("Removing agent pool: %s, resource: %s from template", poolName, resourceName))
				if len(filteredResources) > 0 {
					filteredResources = filteredResources[:len(filteredResources)-1]
				}
			}
		} else if resourceType == vmExtensionType {
			logger.Infoln(fmt.Sprintf("Evaluating if extension: %s needs to be removed", resourceName))
			if strings.Contains(resourceName, "variables('masterVMNamePrefix')") {
				continue
			}

			logger.Infoln(fmt.Sprintf("agentPoolsToPreserve: %v...", agentPoolsToPreserve))

			removeExtension := true
			for poolName, preserve := range agentPoolsToPreserve {
				if strings.Contains(resourceName, "variables('"+poolName) && preserve == true {
					removeExtension = false
				}
			}

			if removeExtension == true {
				logger.Infoln(fmt.Sprintf("Removing extension: %s from template", resourceName))
				if len(filteredResources) > 0 {
					filteredResources = filteredResources[:len(filteredResources)-1]
				}
			}
		}
	}

	templateMap[resourcesFieldName] = filteredResources

	logger.Infoln(fmt.Sprintf("Resource count after running NormalizeResourcesForK8sMasterUpgrade: %d",
		len(templateMap[resourcesFieldName].([]interface{}))))
	return nil
}

// NormalizeResourcesForK8sAgentUpgrade takes a template and removes elements that are unwanted in any scale up/down case
func (t *Transformer) NormalizeResourcesForK8sAgentUpgrade(logger *logrus.Entry, templateMap map[string]interface{}, isMasterManagedDisk bool, agentPoolsToPreserve map[string]bool) error {
	logger.Infoln(fmt.Sprintf("Running NormalizeResourcesForK8sMasterUpgrade...."))
	if err := t.NormalizeResourcesForK8sMasterUpgrade(logger, templateMap, isMasterManagedDisk, agentPoolsToPreserve); err != nil {
		log.Fatalln(err)
		return err
	}

	logger.Infoln(fmt.Sprintf("Running NormalizeForK8sVMASScalingUp...."))
	if err := t.NormalizeForK8sVMASScalingUp(logger, templateMap); err != nil {
		log.Fatalln(err)
		return err
	}

	return nil
}

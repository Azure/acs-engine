package acsengine

import (
	"fmt"
	"strings"

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

	// ARM resource Types
	nsgResourceType  = "Microsoft.Network/networkSecurityGroups"
	vmResourceType   = "Microsoft.Compute/virtualMachines"
	vmssResourceType = "Microsoft.Compute/virtualMachineScaleSets"
)

// NormalizeForVMSSScaling takes a template and removes elements that are unwanted in a VMSS scale up/down case
func NormalizeForVMSSScaling(logger *logrus.Entry, templateMap map[string]interface{}) error {
	if err := NormalizeMasterResourcesForScaling(logger, templateMap); err != nil {
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

		if !removeCustomData(logger, virtualMachineProfile) || !removeImageReference(logger, virtualMachineProfile) {
			continue
		}
	}
	return nil
}

// NormalizeForK8sVMASScalingUp takes a template and removes elements that are unwanted in a K8s VMAS scale up/down case
func NormalizeForK8sVMASScalingUp(logger *logrus.Entry, templateMap map[string]interface{}) error {
	if err := NormalizeMasterResourcesForScaling(logger, templateMap); err != nil {
		return err
	}
	nsgIndex := -1
	resources := templateMap[resourcesFieldName].([]interface{})
	for index, resource := range resources {
		resourceMap, ok := resource.(map[string]interface{})
		if !ok {
			logger.Warnf("Template improperly formatted")
			continue
		}

		resourceType, ok := resourceMap[typeFieldName].(string)
		if ok && resourceType == nsgResourceType {
			if nsgIndex != -1 {
				err := fmt.Errorf("Found 2 resources with type %s in the template. There should only be 1", nsgResourceType)
				logger.Errorf(err.Error())
				return err
			}
			nsgIndex = index
		}

		dependencies, ok := resourceMap[dependsOnFieldName].([]interface{})
		if !ok {
			logger.Warnf("Template improperly formatted")
			continue
		}

		for dIndex := len(dependencies) - 1; dIndex >= 0; dIndex-- {
			dependency := dependencies[dIndex].(string)
			if strings.Contains(dependency, nsgResourceType) {
				dependencies = append(dependencies[:dIndex], dependencies[dIndex+1:]...)
			}
		}

		resourceMap[dependsOnFieldName] = dependencies
	}
	if nsgIndex == -1 {
		err := fmt.Errorf("Found no resources with type %s in the template. There should have been 1", nsgResourceType)
		logger.Errorf(err.Error())
		return err
	}

	templateMap[resourcesFieldName] = append(resources[:nsgIndex], resources[nsgIndex+1:]...)

	return nil
}

// NormalizeMasterResourcesForScaling takes a template and removes elements that are unwanted in any scale up/down case
func NormalizeMasterResourcesForScaling(logger *logrus.Entry, templateMap map[string]interface{}) error {
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

		if !removeCustomData(logger, resourceProperties) || !removeImageReference(logger, resourceProperties) {
			continue
		}
	}

	return nil
}

func removeCustomData(logger *logrus.Entry, resourceProperties map[string]interface{}) bool {
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

func removeImageReference(logger *logrus.Entry, resourceProperties map[string]interface{}) bool {
	storageProfile, ok := resourceProperties[storageProfileFieldName].(map[string]interface{})
	if !ok {
		logger.Warnf("Template improperly formatted")
		return ok
	}

	if storageProfile[imageReferenceFieldName] != nil {
		delete(storageProfile, imageReferenceFieldName)
	}
	return ok
}

package kubernetesupgrade

import (
	"encoding/json"
	"time"

	"github.com/Azure/acs-engine/pkg/acsengine"
	"github.com/Azure/acs-engine/pkg/acsengine/transform"
	"github.com/Azure/acs-engine/pkg/api"
	"github.com/Azure/acs-engine/pkg/armhelpers"
	"github.com/Azure/acs-engine/pkg/armhelpers/utils"
	"github.com/Azure/acs-engine/pkg/i18n"
	"github.com/Azure/acs-engine/pkg/operations"
	"github.com/sirupsen/logrus"
	"k8s.io/api/core/v1"
)

// Upgrader holds information on upgrading an ACS cluster
type Upgrader struct {
	Translator *i18n.Translator
	logger     *logrus.Entry
	ClusterTopology
	Client           armhelpers.ACSEngineClient
	kubeConfig       string
	stepTimeout      *time.Duration
	ACSEngineVersion string
}

type vmStatus int

const (
	defaultTimeout            = time.Minute * 10
	vmStatusUpgraded vmStatus = iota
	vmStatusNotUpgraded
	vmStatusIgnored
)

type vmInfo struct {
	name   string
	status vmStatus
}

// Init initializes an upgrader struct
func (ku *Upgrader) Init(translator *i18n.Translator, logger *logrus.Entry, clusterTopology ClusterTopology, client armhelpers.ACSEngineClient, kubeConfig string, stepTimeout *time.Duration, acsEngineVersion string) {
	ku.Translator = translator
	ku.logger = logger
	ku.ClusterTopology = clusterTopology
	ku.Client = client
	ku.kubeConfig = kubeConfig
	ku.stepTimeout = stepTimeout
	ku.ACSEngineVersion = acsEngineVersion
}

// RunUpgrade runs the upgrade pipeline
func (ku *Upgrader) RunUpgrade() error {
	if err := ku.upgradeMasterNodes(); err != nil {
		return err
	}

	if err := ku.upgradeAgentScaleSets(); err != nil {
		return err
	}

	return ku.upgradeAgentPools()
}

// Validate will run validation post upgrade
func (ku *Upgrader) Validate() error {
	return nil
}

func (ku *Upgrader) upgradeMasterNodes() error {
	if ku.ClusterTopology.DataModel.Properties.MasterProfile == nil {
		return nil
	}
	ku.logger.Infof("Master nodes StorageProfile: %s", ku.ClusterTopology.DataModel.Properties.MasterProfile.StorageProfile)
	// Upgrade Master VMs
	templateMap, parametersMap, err := ku.generateUpgradeTemplate(ku.ClusterTopology.DataModel, ku.ACSEngineVersion)
	if err != nil {
		return ku.Translator.Errorf("error generating upgrade template: %s", err.Error())
	}

	ku.logger.Infof("Prepping master nodes for upgrade...")

	transformer := &transform.Transformer{
		Translator: ku.Translator,
	}
	if err := transformer.NormalizeResourcesForK8sMasterUpgrade(ku.logger, templateMap, ku.DataModel.Properties.MasterProfile.IsManagedDisks(), nil); err != nil {
		ku.logger.Errorf(err.Error())
		return err
	}

	upgradeMasterNode := UpgradeMasterNode{
		Translator: ku.Translator,
		logger:     ku.logger,
	}
	upgradeMasterNode.TemplateMap = templateMap
	upgradeMasterNode.ParametersMap = parametersMap
	upgradeMasterNode.UpgradeContainerService = ku.ClusterTopology.DataModel
	upgradeMasterNode.ResourceGroup = ku.ClusterTopology.ResourceGroup
	upgradeMasterNode.SubscriptionID = ku.ClusterTopology.SubscriptionID
	upgradeMasterNode.Client = ku.Client
	upgradeMasterNode.kubeConfig = ku.kubeConfig
	if ku.stepTimeout == nil {
		upgradeMasterNode.timeout = defaultTimeout
	} else {
		upgradeMasterNode.timeout = *ku.stepTimeout
	}

	expectedMasterCount := ku.ClusterTopology.DataModel.Properties.MasterProfile.Count
	mastersUpgradedCount := len(*ku.ClusterTopology.UpgradedMasterVMs)
	mastersToUgradeCount := expectedMasterCount - mastersUpgradedCount

	ku.logger.Infof("Total expected master count: %d", expectedMasterCount)
	ku.logger.Infof("Master nodes that need to be upgraded: %d", mastersToUgradeCount)
	ku.logger.Infof("Master nodes that have been upgraded: %d", mastersUpgradedCount)

	ku.logger.Infof("Starting upgrade of master nodes...")

	masterNodesInCluster := len(*ku.ClusterTopology.MasterVMs) + mastersUpgradedCount
	ku.logger.Infof("masterNodesInCluster: %d", masterNodesInCluster)
	if masterNodesInCluster > expectedMasterCount {
		return ku.Translator.Errorf("Total count of master VMs: %d exceeded expected count: %d", masterNodesInCluster, expectedMasterCount)
	}

	upgradedMastersIndex := make(map[int]bool)

	for _, vm := range *ku.ClusterTopology.UpgradedMasterVMs {
		ku.logger.Infof("Master VM: %s is upgraded to expected orchestrator version", *vm.Name)
		masterIndex, _ := utils.GetVMNameIndex(vm.StorageProfile.OsDisk.OsType, *vm.Name)
		upgradedMastersIndex[masterIndex] = true
	}

	for _, vm := range *ku.ClusterTopology.MasterVMs {
		ku.logger.Infof("Upgrading Master VM: %s", *vm.Name)

		masterIndex, _ := utils.GetVMNameIndex(vm.StorageProfile.OsDisk.OsType, *vm.Name)

		err := upgradeMasterNode.DeleteNode(vm.Name, false)
		if err != nil {
			ku.logger.Infof("Error deleting master VM: %s, err: %v", *vm.Name, err)
			return err
		}

		err = upgradeMasterNode.CreateNode("master", masterIndex)
		if err != nil {
			ku.logger.Infof("Error creating upgraded master VM: %s", *vm.Name)
			return err
		}

		err = upgradeMasterNode.Validate(vm.Name)
		if err != nil {
			ku.logger.Infof("Error validating upgraded master VM: %s", *vm.Name)
			return err
		}

		upgradedMastersIndex[masterIndex] = true
	}

	// This condition is possible if the previous upgrade operation failed during master
	// VM upgrade when a master VM was deleted but creation of upgraded master did not run.
	if masterNodesInCluster < expectedMasterCount {
		ku.logger.Infof(
			"Found missing master VMs in the cluster. Reconstructing names of missing master VMs for recreation during upgrade...")
	}

	mastersToCreate := expectedMasterCount - masterNodesInCluster
	ku.logger.Infof("Expected master count: %d, Creating %d more master VMs", expectedMasterCount, mastersToCreate)

	// NOTE: this is NOT completely idempotent because it assumes that
	// the OS disk has been deleted
	for i := 0; i < mastersToCreate; i++ {
		masterIndexToCreate := 0
		for upgradedMastersIndex[masterIndexToCreate] {
			masterIndexToCreate++
		}

		ku.logger.Infof("Creating upgraded master VM with index: %d", masterIndexToCreate)

		err = upgradeMasterNode.CreateNode("master", masterIndexToCreate)
		if err != nil {
			ku.logger.Infof("Error creating upgraded master VM with index: %d", masterIndexToCreate)
			return err
		}

		tempVMName := ""
		err = upgradeMasterNode.Validate(&tempVMName)
		if err != nil {
			ku.logger.Infof("Error validating upgraded master VM with index: %d", masterIndexToCreate)
			return err
		}

		upgradedMastersIndex[masterIndexToCreate] = true
	}

	return nil
}

func (ku *Upgrader) upgradeAgentPools() error {
	for _, agentPool := range ku.ClusterTopology.AgentPools {
		// Upgrade Agent VMs
		templateMap, parametersMap, err := ku.generateUpgradeTemplate(ku.ClusterTopology.DataModel, ku.ACSEngineVersion)
		if err != nil {
			ku.logger.Errorf("Error generating upgrade template: %v", err)
			return ku.Translator.Errorf("Error generating upgrade template: %s", err.Error())
		}

		ku.logger.Infof("Prepping agent pool '%s' for upgrade...", *agentPool.Name)

		preservePools := map[string]bool{*agentPool.Name: true}
		transformer := &transform.Transformer{
			Translator: ku.Translator,
		}
		var isMasterManagedDisk bool
		if ku.DataModel.Properties.MasterProfile != nil {
			isMasterManagedDisk = ku.DataModel.Properties.MasterProfile.IsManagedDisks()
		}
		if err := transformer.NormalizeResourcesForK8sAgentUpgrade(ku.logger, templateMap, isMasterManagedDisk, preservePools); err != nil {
			ku.logger.Errorf(err.Error())
			return ku.Translator.Errorf("Error generating upgrade template: %s", err.Error())
		}

		var agentCount, agentPoolIndex int
		var agentOsType api.OSType
		var agentPoolName string
		for indx, app := range ku.ClusterTopology.DataModel.Properties.AgentPoolProfiles {
			if app.Name == *agentPool.Name {
				agentCount = app.Count
				agentOsType = app.OSType
				agentPoolName = app.Name
				agentPoolIndex = indx
				break
			}
		}

		if agentCount == 0 {
			ku.logger.Infof("Agent pool '%s' is empty", *agentPool.Name)
			return nil
		}

		upgradeAgentNode := UpgradeAgentNode{
			Translator: ku.Translator,
			logger:     ku.logger,
		}
		upgradeAgentNode.TemplateMap = templateMap
		upgradeAgentNode.ParametersMap = parametersMap
		upgradeAgentNode.UpgradeContainerService = ku.ClusterTopology.DataModel
		upgradeAgentNode.SubscriptionID = ku.ClusterTopology.SubscriptionID
		upgradeAgentNode.ResourceGroup = ku.ClusterTopology.ResourceGroup
		upgradeAgentNode.Client = ku.Client
		upgradeAgentNode.kubeConfig = ku.kubeConfig
		if ku.stepTimeout == nil {
			upgradeAgentNode.timeout = defaultTimeout
		} else {
			upgradeAgentNode.timeout = *ku.stepTimeout
		}

		agentVMs := make(map[int]*vmInfo)
		// Go over upgraded VMs and verify provisioning state
		// per https://docs.microsoft.com/en-us/rest/api/compute/virtualmachines/virtualmachines-state :
		//  - Creating: Indicates the virtual Machine is being created.
		//  - Updating: Indicates that there is an update operation in progress on the Virtual Machine.
		//  - Succeeded: Indicates that the operation executed on the virtual machine succeeded.
		//  - Deleting: Indicates that the virtual machine is being deleted.
		//  - Failed: Indicates that the update operation on the Virtual Machine failed.
		// Delete VMs in 'bad' state. Such VMs will be re-created later in this function.
		upgradedCount := 0
		for _, vm := range *agentPool.UpgradedAgentVMs {
			ku.logger.Infof("Agent VM: %s, pool name: %s on expected orchestrator version", *vm.Name, *agentPool.Name)
			var vmProvisioningState string
			if vm.VirtualMachineProperties != nil && vm.VirtualMachineProperties.ProvisioningState != nil {
				vmProvisioningState = *vm.VirtualMachineProperties.ProvisioningState
			}
			agentIndex, _ := utils.GetVMNameIndex(vm.StorageProfile.OsDisk.OsType, *vm.Name)

			switch vmProvisioningState {
			case "Creating", "Updating", "Succeeded":
				agentVMs[agentIndex] = &vmInfo{*vm.Name, vmStatusUpgraded}
				upgradedCount++

			case "Failed":
				ku.logger.Infof("Deleting agent VM %s in provisioning state %s", *vm.Name, vmProvisioningState)
				err := upgradeAgentNode.DeleteNode(vm.Name, false)
				if err != nil {
					ku.logger.Errorf("Error deleting agent VM %s: %v", *vm.Name, err)
					return err
				}

			case "Deleting":
				fallthrough
			default:
				ku.logger.Infof("Ignoring agent VM %s in provisioning state %s", *vm.Name, vmProvisioningState)
				agentVMs[agentIndex] = &vmInfo{*vm.Name, vmStatusIgnored}
			}
		}

		for _, vm := range *agentPool.AgentVMs {
			agentIndex, _ := utils.GetVMNameIndex(vm.StorageProfile.OsDisk.OsType, *vm.Name)
			agentVMs[agentIndex] = &vmInfo{*vm.Name, vmStatusNotUpgraded}
		}
		toBeUpgradedCount := len(*agentPool.AgentVMs)

		ku.logger.Infof("Starting upgrade of %d agent nodes (out of %d) in pool identifier: %s, name: %s...",
			toBeUpgradedCount, agentCount, *agentPool.Identifier, *agentPool.Name)

		// Create missing nodes to match agentCount. This could be due to previous upgrade failure
		// If there are nodes that need to be upgraded, create one extra node, which will be used to take on the load from upgrading nodes.
		if toBeUpgradedCount > 0 {
			agentCount++
		}
		for upgradedCount+toBeUpgradedCount < agentCount {
			agentIndex := getAvailableIndex(agentVMs)

			vmName, err := utils.GetK8sVMName(agentOsType, ku.DataModel.Properties.HostedMasterProfile != nil,
				ku.NameSuffix, agentPoolName, agentPoolIndex, agentIndex)
			if err != nil {
				ku.logger.Errorf("Error reconstructing agent VM name with index %d: %v", agentIndex, err)
				return err
			}
			ku.logger.Infof("Creating new agent node %s (index %d)", vmName, agentIndex)

			err = upgradeAgentNode.CreateNode(*agentPool.Name, agentIndex)
			if err != nil {
				ku.logger.Errorf("Error creating agent node %s (index %d): %v", vmName, agentIndex, err)
				return err
			}

			err = upgradeAgentNode.Validate(&vmName)
			if err != nil {
				ku.logger.Infof("Error validating agent node %s (index %d): %v", vmName, agentIndex, err)
				return err
			}

			agentVMs[agentIndex] = &vmInfo{vmName, vmStatusUpgraded}
			upgradedCount++
		}

		if toBeUpgradedCount == 0 {
			ku.logger.Infof("No nodes to upgrade")
			return nil
		}

		// Upgrade nodes in agent pool
		upgradedCount = 0
		for agentIndex, vm := range agentVMs {
			if vm.status != vmStatusNotUpgraded {
				continue
			}
			ku.logger.Infof("Upgrading Agent VM: %s, pool name: %s", vm.name, *agentPool.Name)

			err := upgradeAgentNode.DeleteNode(&vm.name, true)
			if err != nil {
				ku.logger.Errorf("Error deleting agent VM %s: %v", vm.name, err)
				return err
			}

			// do not create last node in favor of already created extra node.
			if upgradedCount == toBeUpgradedCount-1 {
				ku.logger.Infof("Skipping creation of VM %s (index %d)", vm.name, agentIndex)
				delete(agentVMs, agentIndex)
			} else {
				err = upgradeAgentNode.CreateNode(*agentPool.Name, agentIndex)
				if err != nil {
					ku.logger.Errorf("Error creating upgraded agent VM %s: %v", vm.name, err)
					return err
				}

				err = upgradeAgentNode.Validate(&vm.name)
				if err != nil {
					ku.logger.Errorf("Error validating upgraded agent VM %s: %v", vm.name, err)
					return err
				}
				vm.status = vmStatusUpgraded
			}
			upgradedCount++
		}
	}

	return nil
}

func (ku *Upgrader) upgradeAgentScaleSets() error {
	for _, vmssToUpgrade := range ku.ClusterTopology.AgentPoolScaleSetsToUpgrade {
		ku.logger.Infof("Upgrading VMSS %s", vmssToUpgrade.Name)

		if len(vmssToUpgrade.VMsToUpgrade) == 0 {
			ku.logger.Infof("No VMs to upgrade for VMSS %s, skipping", vmssToUpgrade.Name)
			continue
		}

		newCapacity := *vmssToUpgrade.Sku.Capacity + 1
		ku.logger.Infof(
			"VMSS %s current capacity is %d and new capacity will be %d while each node is swapped",
			vmssToUpgrade.Name,
			*vmssToUpgrade.Sku.Capacity,
			newCapacity,
		)

		*vmssToUpgrade.Sku.Capacity = newCapacity

		for _, vmToUpgrade := range vmssToUpgrade.VMsToUpgrade {
			success, failure := ku.Client.SetVirtualMachineScaleSetCapacity(
				ku.ClusterTopology.ResourceGroup,
				vmssToUpgrade.Name,
				vmssToUpgrade.Sku,
				vmssToUpgrade.Location,
				make(chan struct{}),
			)

			select {
			case <-success:
				ku.logger.Infof("Successfully set capacity for VMSS %s", vmssToUpgrade.Name)
			case err := <-failure:
				ku.logger.Errorf("Failure to set capacity for VMSS %s", vmssToUpgrade.Name)
				return err
			}

			// Before we can delete the node we should safely and responsibly drain it
			var kubeAPIServerURL string
			getClientTimeout := 10 * time.Second

			if ku.DataModel.Properties.HostedMasterProfile != nil {
				kubeAPIServerURL = ku.DataModel.Properties.HostedMasterProfile.FQDN
			} else {
				kubeAPIServerURL = ku.DataModel.Properties.MasterProfile.FQDN
			}
			client, err := ku.Client.GetKubernetesClient(
				kubeAPIServerURL,
				ku.kubeConfig,
				interval,
				getClientTimeout,
			)
			if err != nil {
				ku.logger.Errorf("Error getting Kubernetes client: %v", err)
				return err
			}

			ku.logger.Infof("Draining node %s", vmToUpgrade.Name)
			err = operations.SafelyDrainNodeWithClient(
				client,
				ku.logger,
				vmToUpgrade.Name,
				time.Minute,
			)
			if err != nil {
				ku.logger.Errorf("Error draining VM in VMSS: %v", err)
				return err
			}

			ku.logger.Infof(
				"Deleting VM %s in VMSS %s",
				vmToUpgrade.Name,
				vmssToUpgrade.Name,
			)

			// At this point we have our buffer node that will replace the node to delete
			// so we can just remove this current node then
			res, failure := ku.Client.DeleteVirtualMachineScaleSetVM(
				ku.ClusterTopology.ResourceGroup,
				vmssToUpgrade.Name,
				vmToUpgrade.InstanceID,
				make(chan struct{}),
			)

			select {
			case <-res:
				ku.logger.Infof(
					"Successfully deleted VM %s in VMSS %s",
					vmToUpgrade.Name,
					vmssToUpgrade.Name,
				)
			case err := <-failure:
				ku.logger.Errorf(
					"Failed to delete VM %s in VMSS %s",
					vmToUpgrade.Name,
					vmssToUpgrade,
				)
				return err
			}
		}
		ku.logger.Infof("Completed upgrading VMSS %s", vmssToUpgrade)
	}

	ku.logger.Infoln("Completed upgrading all VMSS")

	return nil
}

func (ku *Upgrader) generateUpgradeTemplate(upgradeContainerService *api.ContainerService, acsengineVersion string) (map[string]interface{}, map[string]interface{}, error) {
	var err error
	ctx := acsengine.Context{
		Translator: ku.Translator,
	}
	templateGenerator, err := acsengine.InitializeTemplateGenerator(ctx, false)
	if err != nil {
		return nil, nil, ku.Translator.Errorf("failed to initialize template generator: %s", err.Error())
	}

	var templateJSON string
	var parametersJSON string
	if templateJSON, parametersJSON, _, err = templateGenerator.GenerateTemplate(upgradeContainerService, acsengine.DefaultGeneratorCode, true, acsengineVersion); err != nil {
		return nil, nil, ku.Translator.Errorf("error generating upgrade template: %s", err.Error())
	}

	var template interface{}
	var parameters interface{}
	json.Unmarshal([]byte(templateJSON), &template)
	json.Unmarshal([]byte(parametersJSON), &parameters)
	templateMap := template.(map[string]interface{})
	parametersMap := parameters.(map[string]interface{})

	return templateMap, parametersMap, nil
}

// return unused index within the range of agent indices, or subsequent index
func getAvailableIndex(vms map[int]*vmInfo) int {
	maxIndex := 0

	for indx := range vms {
		if indx > maxIndex {
			maxIndex = indx
		}
	}

	for indx := 0; indx < maxIndex; indx++ {
		if _, found := vms[indx]; !found {
			return indx
		}
	}

	return maxIndex + 1
}

// isNodeReady returns true if a node is ready; false otherwise.
// Copied from: https://github.com/kubernetes/kubernetes/blob/886e04f1fffbb04faf8a9f9ee141143b2684ae68/pkg/api/v1/node/util.go#L40
func isNodeReady(node *v1.Node) bool {
	for _, c := range node.Status.Conditions {
		if c.Type == v1.NodeReady {
			return c.Status == v1.ConditionTrue
		}
	}
	return false
}

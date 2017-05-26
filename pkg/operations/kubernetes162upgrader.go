package operations

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"

	"strings"

	"github.com/Azure/acs-engine/pkg/acsengine"
	"github.com/Azure/acs-engine/pkg/api"
	"github.com/Azure/acs-engine/pkg/armhelpers"
	log "github.com/Sirupsen/logrus"
)

// Compiler to verify QueueMessageProcessor implements OperationsProcessor
var _ UpgradeWorkFlow = &Kubernetes162upgrader{}

// Kubernetes162upgrader upgrades a Kubernetes 1.5.3 cluster to 1.6.2
type Kubernetes162upgrader struct {
	ClusterTopology
	GoalStateDataModel *api.ContainerService

	Client armhelpers.ACSEngineClient
}

// ClusterPreflightCheck does preflight check
func (ku *Kubernetes162upgrader) ClusterPreflightCheck() error {
	// Check that current cluster is 1.5.3
	if ku.DataModel.Properties.OrchestratorProfile.OrchestratorVersion != api.Kubernetes153 {
		return fmt.Errorf("Upgrade to Kubernetes 1.6.2 is not supported from version: %s",
			ku.DataModel.Properties.OrchestratorProfile.OrchestratorVersion)
	}

	return nil
}

// RunUpgrade runs the upgrade pipeline
func (ku *Kubernetes162upgrader) RunUpgrade() error {
	if err := ku.ClusterPreflightCheck(); err != nil {
		return err
	}

	ku.GoalStateDataModel = ku.ClusterTopology.DataModel
	ku.GoalStateDataModel.Properties.OrchestratorProfile.OrchestratorVersion = api.Kubernetes162

	if err := ku.upgradeMasterNodes(); err != nil {
		return err
	}

	if err := ku.upgradeAgentPools(); err != nil {
		return err
	}

	return nil
}

// Validate will run validation post upgrade
func (ku *Kubernetes162upgrader) Validate() error {
	return nil
}

func (ku *Kubernetes162upgrader) upgradeMasterNodes() error {
	// Upgrade Master VMs
	templateMap, parametersMap, err := ku.generateUpgradeTemplate(ku.GoalStateDataModel)
	if err != nil {
		return fmt.Errorf("error generating upgrade template: %s", err.Error())
	}

	log.Infoln(fmt.Sprintf("Prepping master nodes for upgrade..."))

	if err := acsengine.NormalizeResourcesForK8sMasterUpgrade(log.NewEntry(log.New()), templateMap, nil); err != nil {
		log.Fatalln(err)
		return err
	}

	upgradeMasterNode := UpgradeMasterNode{}
	upgradeMasterNode.TemplateMap = templateMap
	upgradeMasterNode.ParametersMap = parametersMap
	upgradeMasterNode.UpgradeContainerService = ku.GoalStateDataModel
	upgradeMasterNode.ResourceGroup = ku.ClusterTopology.ResourceGroup
	upgradeMasterNode.Client = ku.Client

	expectedMasterCount := ku.GoalStateDataModel.Properties.MasterProfile.Count
	mastersUpgradedCount := len(*ku.ClusterTopology.UpgradedMasterVMs)
	mastersToUgradeCount := expectedMasterCount - mastersUpgradedCount

	log.Infoln(fmt.Sprintf("Total expected master count: %d", expectedMasterCount))
	log.Infoln(fmt.Sprintf("Master nodes that need to be upgraded: %d", mastersToUgradeCount))
	log.Infoln(fmt.Sprintf("Master nodes that have been upgraded: %d", mastersUpgradedCount))

	log.Infoln(fmt.Sprintf("Starting master nodes upgrade..."))

	masterNodesInCluster := len(*ku.ClusterTopology.MasterVMs) + mastersUpgradedCount
	log.Infoln(fmt.Sprintf("masterNodesInCluster: %d", masterNodesInCluster))
	if masterNodesInCluster > expectedMasterCount {
		return fmt.Errorf("Total count of master VMs: %d exceeded expected count: %d", masterNodesInCluster, expectedMasterCount)
	}

	masterVMsToUgradeStatus := make(VMStatusSlice, 0, masterNodesInCluster)

	var sampleMasterVMName *string
	// TODO pass suffix and set this
	temp := "k8s-master-12345678-0"
	sampleMasterVMName = &temp

	// Add Master VMs that need to be upgraded to masterVMsToUgradeStatus
	for _, vm := range *ku.ClusterTopology.MasterVMs {
		vmStatus := VMStatus{vm.Name, true, true}
		sampleMasterVMName = vm.Name
		masterVMsToUgradeStatus = append(masterVMsToUgradeStatus, &vmStatus)

		log.Infoln(fmt.Sprintf("Adding master VM name: %s with tag: %s to masterVMsToUgradeStatus list",
			*vm.Name, *(*vm.Tags)["orchestrator"]))
	}

	// Add Master VMs that have already been upgraded to masterVMsToUgradeStatus
	for _, vm := range *ku.ClusterTopology.UpgradedMasterVMs {
		vmStatus := VMStatus{vm.Name, false, false}
		masterVMsToUgradeStatus = append(masterVMsToUgradeStatus, &vmStatus)

		log.Infoln(fmt.Sprintf("Adding master VM name: %s with tag: %s to masterVMsToUgradeStatus list",
			*vm.Name, *(*vm.Tags)["orchestrator"]))
	}

	// This condition is possible if the previous upgrade operation failed during master
	// VM upgrade when a master VM was deleted but creation of upgraded master did not run.
	if masterNodesInCluster < expectedMasterCount {
		log.Infoln(fmt.Sprintf(
			"Found missing master VMs in the cluster. Reconstructing names of missing master VMs for recreation during upgrade..."))

		// Note that this assumes that VM numbers were in the pool were consecutive
		availableMasterVMs := make([]bool, expectedMasterCount)
		for _, masterVMStatus := range masterVMsToUgradeStatus {
			availableMasterVMs[masterVMStatus.VMNumber()] = true
		}

		// orchestrator, type, suffix, _, err := armhelpers.LinuxVMNameParts(*sampleMasterVMName)
		orch, poolType, suffix, _, err := armhelpers.LinuxVMNameParts(*sampleMasterVMName)
		if err != nil {
			log.Fatalln(err)
			return err
		}

		for i := 0; i < len(availableMasterVMs); i++ {
			if availableMasterVMs[i] == false {
				var vmNameArray = []string{orch, poolType, suffix, strconv.Itoa(i)}
				var vmName = strings.Join(vmNameArray, "-")
				vmStatus := VMStatus{&vmName, false, true}
				log.Infoln(fmt.Sprintf(
					"Adding missing master VM: %s for recreation during upgrade...", vmName))
				masterVMsToUgradeStatus = append(masterVMsToUgradeStatus, &vmStatus)
			}
		}
	}

	sort.Sort(sort.Reverse(masterVMsToUgradeStatus))

	for _, vm := range masterVMsToUgradeStatus {
		log.Infoln(fmt.Sprintf("Upgrading Master VM: %s", *vm.Name))

		if vm.Delete == true {
			log.Infoln(fmt.Sprintf("Deleting Master VM: %s", *vm.Name))
			// 1.	Shutdown and delete one master VM at a time while preserving the persistent disk backing etcd.
			upgradeMasterNode.DeleteNode(vm.Name)
		}

		if vm.Upgrade == true {
			log.Infoln(fmt.Sprintf("Creating upgraded Master VM: %s", *vm.Name))
			// 2.	Call CreateVMWithRetries
			upgradeMasterNode.CreateNode("master", vm.VMNumber())
		}

		upgradeMasterNode.Validate()
	}

	return nil
}

func (ku *Kubernetes162upgrader) upgradeAgentPools() error {
	for _, agentPool := range ku.ClusterTopology.AgentPools {
		// Upgrade Agent VMs
		templateMap, parametersMap, err := ku.generateUpgradeTemplate(ku.GoalStateDataModel)
		if err != nil {
			return fmt.Errorf("error generating upgrade template: %s", err.Error())
		}

		log.Infoln(fmt.Sprintf("Prepping agent pool: %s for upgrade...", *agentPool.Name))

		preservePools := map[string]bool{*agentPool.Name: true}
		if err := acsengine.NormalizeResourcesForK8sAgentUpgrade(log.NewEntry(log.New()), templateMap, preservePools); err != nil {
			log.Fatalln(err)
			return err
		}

		var agentCount int
		for _, app := range ku.GoalStateDataModel.Properties.AgentPoolProfiles {
			if app.Name == *agentPool.Name {
				agentCount = app.Count
				break
			}
		}

		upgradeAgentNode := UpgradeAgentNode{}
		upgradeAgentNode.TemplateMap = templateMap
		upgradeAgentNode.ParametersMap = parametersMap
		upgradeAgentNode.UpgradeContainerService = ku.GoalStateDataModel
		upgradeAgentNode.ResourceGroup = ku.ClusterTopology.ResourceGroup
		upgradeAgentNode.Client = ku.Client

		upgradedAgentsIndex := make(map[int]bool)

		for _, vm := range *agentPool.UpgradedAgentVMs {
			log.Infoln(fmt.Sprintf("Agent VM: %s, pool name: %s is upgraded to expected orchestrator version", *vm.Name, *agentPool.Name))
			agentIndex, _ := armhelpers.GetVMNameIndex(vm.StorageProfile.OsDisk.OsType, *vm.Name)
			upgradedAgentsIndex[agentIndex] = true
		}

		log.Infoln(fmt.Sprintf("Starting upgrade of agent nodes in pool identifier: %s, name: %s...", *agentPool.Identifier, *agentPool.Name))

		for _, vm := range *agentPool.AgentVMs {
			log.Infoln(fmt.Sprintf("Upgrading Agent VM: %s, pool name: %s", *vm.Name, *agentPool.Name))

			agentIndex, _ := armhelpers.GetVMNameIndex(vm.StorageProfile.OsDisk.OsType, *vm.Name)

			// 1.	Shutdown and delete one agent VM at a time
			upgradeAgentNode.DeleteNode(vm.Name)

			// 2.	Call CreateVMWithRetries
			upgradeAgentNode.CreateNode(*agentPool.Name, agentIndex)

			upgradeAgentNode.Validate()

			upgradedAgentsIndex[agentIndex] = true
		}

		agentsToCreate := agentCount - len(upgradedAgentsIndex)
		log.Infoln(fmt.Sprintf("Expected agent count in the pool: %d, Creating %d more agents", agentCount, agentsToCreate))

		// NOTE: this is NOT completely idempotent because it assumes that
		// the OS disk has been deleted
		for i := 0; i < agentsToCreate; i++ {
			agentIndexToCreate := 0
			for upgradedAgentsIndex[agentIndexToCreate] == true {
				agentIndexToCreate++
			}

			log.Infoln(fmt.Sprintf("Creating upgraded Agent VM with index: %d, pool name: %s", agentIndexToCreate, *agentPool.Name))
			upgradeAgentNode.CreateNode(*agentPool.Name, agentIndexToCreate)

			upgradeAgentNode.Validate()

			upgradedAgentsIndex[agentIndexToCreate] = true
		}
	}

	return nil
}

func (ku *Kubernetes162upgrader) generateUpgradeTemplate(upgradeContainerService *api.ContainerService) (map[string]interface{}, map[string]interface{}, error) {
	var err error
	templateGenerator, err := acsengine.InitializeTemplateGenerator(false)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to initialize template generator: %s", err.Error())
	}

	var templateJSON string
	var parametersJSON string
	if templateJSON, parametersJSON, _, err = templateGenerator.GenerateTemplate(upgradeContainerService); err != nil {
		return nil, nil, fmt.Errorf("error generating upgrade template: %s", err.Error())
	}

	var template interface{}
	var parameters interface{}
	json.Unmarshal([]byte(templateJSON), &template)
	json.Unmarshal([]byte(parametersJSON), &parameters)
	templateMap := template.(map[string]interface{})
	parametersMap := parameters.(map[string]interface{})

	return templateMap, parametersMap, nil
}

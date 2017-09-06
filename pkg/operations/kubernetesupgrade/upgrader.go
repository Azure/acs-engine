package kubernetesupgrade

import (
	"encoding/json"
	"fmt"

	"github.com/Azure/acs-engine/pkg/acsengine"
	"github.com/Azure/acs-engine/pkg/api"
	"github.com/Azure/acs-engine/pkg/armhelpers"
	"github.com/Azure/acs-engine/pkg/i18n"
	log "github.com/sirupsen/logrus"
)

// Upgrader holds information on upgrading an ACS cluster
type Upgrader struct {
	Translator *i18n.Translator
	ClusterTopology
	Client armhelpers.ACSEngineClient
}

// Init initializes an upgrader struct
func (ku *Upgrader) Init(translator *i18n.Translator, clusterTopology ClusterTopology, client armhelpers.ACSEngineClient) {
	ku.Translator = translator
	ku.ClusterTopology = clusterTopology
	ku.Client = client
}

// RunUpgrade runs the upgrade pipeline
func (ku *Upgrader) RunUpgrade() error {
	if err := ku.upgradeMasterNodes(); err != nil {
		return err
	}

	if err := ku.upgradeAgentPools(); err != nil {
		return err
	}

	return nil
}

// Validate will run validation post upgrade
func (ku *Upgrader) Validate() error {
	return nil
}

func (ku *Upgrader) upgradeMasterNodes() error {
	log.Infoln(fmt.Sprintf("Master nodes StorageProfile: %s", ku.ClusterTopology.DataModel.Properties.MasterProfile.StorageProfile))
	// Upgrade Master VMs
	templateMap, parametersMap, err := ku.generateUpgradeTemplate(ku.ClusterTopology.DataModel)
	if err != nil {
		return ku.Translator.Errorf("error generating upgrade template: %s", err.Error())
	}

	log.Infoln(fmt.Sprintf("Prepping master nodes for upgrade..."))

	transformer := &acsengine.Transformer{
		Translator: ku.Translator,
	}
	if err := transformer.NormalizeResourcesForK8sMasterUpgrade(log.NewEntry(log.New()), templateMap, ku.DataModel.Properties.MasterProfile.IsManagedDisks(), nil); err != nil {
		log.Fatalln(err)
		return err
	}

	upgradeMasterNode := UpgradeMasterNode{
		Translator: ku.Translator,
	}
	upgradeMasterNode.TemplateMap = templateMap
	upgradeMasterNode.ParametersMap = parametersMap
	upgradeMasterNode.UpgradeContainerService = ku.ClusterTopology.DataModel
	upgradeMasterNode.ResourceGroup = ku.ClusterTopology.ResourceGroup
	upgradeMasterNode.Client = ku.Client

	expectedMasterCount := ku.ClusterTopology.DataModel.Properties.MasterProfile.Count
	mastersUpgradedCount := len(*ku.ClusterTopology.UpgradedMasterVMs)
	mastersToUgradeCount := expectedMasterCount - mastersUpgradedCount

	log.Infoln(fmt.Sprintf("Total expected master count: %d", expectedMasterCount))
	log.Infoln(fmt.Sprintf("Master nodes that need to be upgraded: %d", mastersToUgradeCount))
	log.Infoln(fmt.Sprintf("Master nodes that have been upgraded: %d", mastersUpgradedCount))

	log.Infoln(fmt.Sprintf("Starting upgrade of master nodes..."))

	masterNodesInCluster := len(*ku.ClusterTopology.MasterVMs) + mastersUpgradedCount
	log.Infoln(fmt.Sprintf("masterNodesInCluster: %d", masterNodesInCluster))
	if masterNodesInCluster > expectedMasterCount {
		return ku.Translator.Errorf("Total count of master VMs: %d exceeded expected count: %d", masterNodesInCluster, expectedMasterCount)
	}

	upgradedMastersIndex := make(map[int]bool)

	for _, vm := range *ku.ClusterTopology.UpgradedMasterVMs {
		log.Infoln(fmt.Sprintf("Master VM: %s is upgraded to expected orchestrator version", *vm.Name))
		masterIndex, _ := armhelpers.GetVMNameIndex(vm.StorageProfile.OsDisk.OsType, *vm.Name)
		upgradedMastersIndex[masterIndex] = true
	}

	for _, vm := range *ku.ClusterTopology.MasterVMs {
		log.Infoln(fmt.Sprintf("Upgrading Master VM: %s", *vm.Name))

		masterIndex, _ := armhelpers.GetVMNameIndex(vm.StorageProfile.OsDisk.OsType, *vm.Name)

		err := upgradeMasterNode.DeleteNode(vm.Name)
		if err != nil {
			log.Infoln(fmt.Sprintf("Error deleting master VM: %s, err: %v", *vm.Name, err))
			return err
		}

		err = upgradeMasterNode.CreateNode("master", masterIndex)
		if err != nil {
			log.Infoln(fmt.Sprintf("Error creating upgraded master VM: %s", *vm.Name))
			return err
		}

		err = upgradeMasterNode.Validate()
		if err != nil {
			log.Infoln(fmt.Sprintf("Error validating upgraded master VM: %s", *vm.Name))
			return err
		}

		upgradedMastersIndex[masterIndex] = true
	}

	// This condition is possible if the previous upgrade operation failed during master
	// VM upgrade when a master VM was deleted but creation of upgraded master did not run.
	if masterNodesInCluster < expectedMasterCount {
		log.Infoln(fmt.Sprintf(
			"Found missing master VMs in the cluster. Reconstructing names of missing master VMs for recreation during upgrade..."))
	}

	mastersToCreate := expectedMasterCount - masterNodesInCluster
	log.Infoln(fmt.Sprintf("Expected master count: %d, Creating %d more master VMs", expectedMasterCount, mastersToCreate))

	// NOTE: this is NOT completely idempotent because it assumes that
	// the OS disk has been deleted
	for i := 0; i < mastersToCreate; i++ {
		masterIndexToCreate := 0
		for upgradedMastersIndex[masterIndexToCreate] == true {
			masterIndexToCreate++
		}

		log.Infoln(fmt.Sprintf("Creating upgraded master VM with index: %d", masterIndexToCreate))

		err = upgradeMasterNode.CreateNode("master", masterIndexToCreate)
		if err != nil {
			log.Infoln(fmt.Sprintf("Error creating upgraded master VM with index: %d", masterIndexToCreate))
			return err
		}

		err = upgradeMasterNode.Validate()
		if err != nil {
			log.Infoln(fmt.Sprintf("Error validating upgraded master VM with index: %d", masterIndexToCreate))
			return err
		}

		upgradedMastersIndex[masterIndexToCreate] = true
	}

	return nil
}

func (ku *Upgrader) upgradeAgentPools() error {
	for _, agentPool := range ku.ClusterTopology.AgentPools {
		// Upgrade Agent VMs
		templateMap, parametersMap, err := ku.generateUpgradeTemplate(ku.ClusterTopology.DataModel)
		if err != nil {
			return ku.Translator.Errorf("error generating upgrade template: %s", err.Error())
		}

		log.Infoln(fmt.Sprintf("Prepping agent pool: %s for upgrade...", *agentPool.Name))

		preservePools := map[string]bool{*agentPool.Name: true}
		transformer := &acsengine.Transformer{
			Translator: ku.Translator,
		}
		if err := transformer.NormalizeResourcesForK8sAgentUpgrade(log.NewEntry(log.New()), templateMap, ku.DataModel.Properties.MasterProfile.IsManagedDisks(), preservePools); err != nil {
			log.Fatalln(err)
			return err
		}

		var agentCount int
		for _, app := range ku.ClusterTopology.DataModel.Properties.AgentPoolProfiles {
			if app.Name == *agentPool.Name {
				agentCount = app.Count
				break
			}
		}

		upgradeAgentNode := UpgradeAgentNode{
			Translator: ku.Translator,
		}
		upgradeAgentNode.TemplateMap = templateMap
		upgradeAgentNode.ParametersMap = parametersMap
		upgradeAgentNode.UpgradeContainerService = ku.ClusterTopology.DataModel
		upgradeAgentNode.ResourceGroup = ku.ClusterTopology.ResourceGroup
		upgradeAgentNode.Client = ku.Client

		upgradedAgentsIndex := make(map[int]bool)

		for _, vm := range *agentPool.UpgradedAgentVMs {
			log.Infoln(fmt.Sprintf("Agent VM: %s, pool name: %s on expected orchestrator version", *vm.Name, *agentPool.Name))
			agentIndex, _ := armhelpers.GetVMNameIndex(vm.StorageProfile.OsDisk.OsType, *vm.Name)
			upgradedAgentsIndex[agentIndex] = true
		}

		log.Infoln(fmt.Sprintf("Starting upgrade of agent nodes in pool identifier: %s, name: %s...",
			*agentPool.Identifier, *agentPool.Name))

		for _, vm := range *agentPool.AgentVMs {
			log.Infoln(fmt.Sprintf("Upgrading Agent VM: %s, pool name: %s", *vm.Name, *agentPool.Name))

			agentIndex, _ := armhelpers.GetVMNameIndex(vm.StorageProfile.OsDisk.OsType, *vm.Name)

			err := upgradeAgentNode.DeleteNode(vm.Name)
			if err != nil {
				log.Infoln(fmt.Sprintf("Error deleting agent VM: %s", *vm.Name))
				return err
			}

			err = upgradeAgentNode.CreateNode(*agentPool.Name, agentIndex)
			if err != nil {
				log.Infoln(fmt.Sprintf("Error creating upgraded agent VM: %s", *vm.Name))
				return err
			}

			err = upgradeAgentNode.Validate()
			if err != nil {
				log.Infoln(fmt.Sprintf("Error validating upgraded agent VM: %s", *vm.Name))
				return err
			}

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

			err = upgradeAgentNode.CreateNode(*agentPool.Name, agentIndexToCreate)
			if err != nil {
				log.Infoln(fmt.Sprintf("Error creating upgraded agent VM with index: %d", agentIndexToCreate))
				return err
			}

			err = upgradeAgentNode.Validate()
			if err != nil {
				log.Infoln(fmt.Sprintf("Error validating upgraded agent VM with index: %d", agentIndexToCreate))
				return err
			}

			upgradedAgentsIndex[agentIndexToCreate] = true
		}
	}

	return nil
}

func (ku *Upgrader) generateUpgradeTemplate(upgradeContainerService *api.ContainerService) (map[string]interface{}, map[string]interface{}, error) {
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
	if templateJSON, parametersJSON, _, err = templateGenerator.GenerateTemplate(upgradeContainerService); err != nil {
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

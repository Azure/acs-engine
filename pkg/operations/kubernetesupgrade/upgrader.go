package kubernetesupgrade

import (
	"encoding/json"

	"github.com/Azure/acs-engine/pkg/acsengine"
	"github.com/Azure/acs-engine/pkg/api"
	"github.com/Azure/acs-engine/pkg/armhelpers"
	"github.com/Azure/acs-engine/pkg/i18n"
	"github.com/sirupsen/logrus"
)

// Upgrader holds information on upgrading an ACS cluster
type Upgrader struct {
	Translator  *i18n.Translator
	TraceLogger *logrus.Entry
	ClusterTopology
	Client     armhelpers.ACSEngineClient
	kubeConfig string
}

// Init initializes an upgrader struct
func (ku *Upgrader) Init(translator *i18n.Translator, traceLogger *logrus.Entry, clusterTopology ClusterTopology, client armhelpers.ACSEngineClient, kubeConfig string) {
	ku.Translator = translator
	ku.TraceLogger = traceLogger
	ku.ClusterTopology = clusterTopology
	ku.Client = client
	ku.kubeConfig = kubeConfig
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
	if ku.ClusterTopology.DataModel.Properties.MasterProfile == nil {
		return nil
	}
	ku.TraceLogger.Infof("Master nodes StorageProfile: %s\n", ku.ClusterTopology.DataModel.Properties.MasterProfile.StorageProfile)
	// Upgrade Master VMs
	templateMap, parametersMap, err := ku.generateUpgradeTemplate(ku.ClusterTopology.DataModel)
	if err != nil {
		return ku.Translator.Errorf("error generating upgrade template: %s", err.Error())
	}

	ku.TraceLogger.Infof("Prepping master nodes for upgrade...\n")

	transformer := &acsengine.Transformer{
		Translator: ku.Translator,
	}
	if err := transformer.NormalizeResourcesForK8sMasterUpgrade(ku.TraceLogger, templateMap, ku.DataModel.Properties.MasterProfile.IsManagedDisks(), nil); err != nil {
		ku.TraceLogger.Errorf(err.Error())
		return err
	}

	upgradeMasterNode := UpgradeMasterNode{
		Translator:  ku.Translator,
		TraceLogger: ku.TraceLogger,
	}
	upgradeMasterNode.TemplateMap = templateMap
	upgradeMasterNode.ParametersMap = parametersMap
	upgradeMasterNode.UpgradeContainerService = ku.ClusterTopology.DataModel
	upgradeMasterNode.ResourceGroup = ku.ClusterTopology.ResourceGroup
	upgradeMasterNode.Client = ku.Client

	expectedMasterCount := ku.ClusterTopology.DataModel.Properties.MasterProfile.Count
	mastersUpgradedCount := len(*ku.ClusterTopology.UpgradedMasterVMs)
	mastersToUgradeCount := expectedMasterCount - mastersUpgradedCount

	ku.TraceLogger.Infof("Total expected master count: %d\n", expectedMasterCount)
	ku.TraceLogger.Infof("Master nodes that need to be upgraded: %d\n", mastersToUgradeCount)
	ku.TraceLogger.Infof("Master nodes that have been upgraded: %d\n", mastersUpgradedCount)

	ku.TraceLogger.Infof("Starting upgrade of master nodes...\n")

	masterNodesInCluster := len(*ku.ClusterTopology.MasterVMs) + mastersUpgradedCount
	ku.TraceLogger.Infof("masterNodesInCluster: %d\n", masterNodesInCluster)
	if masterNodesInCluster > expectedMasterCount {
		return ku.Translator.Errorf("Total count of master VMs: %d exceeded expected count: %d", masterNodesInCluster, expectedMasterCount)
	}

	upgradedMastersIndex := make(map[int]bool)

	for _, vm := range *ku.ClusterTopology.UpgradedMasterVMs {
		ku.TraceLogger.Infof("Master VM: %s is upgraded to expected orchestrator version\n", *vm.Name)
		masterIndex, _ := armhelpers.GetVMNameIndex(vm.StorageProfile.OsDisk.OsType, *vm.Name)
		upgradedMastersIndex[masterIndex] = true
	}

	for _, vm := range *ku.ClusterTopology.MasterVMs {
		ku.TraceLogger.Infof("Upgrading Master VM: %s\n", *vm.Name)

		masterIndex, _ := armhelpers.GetVMNameIndex(vm.StorageProfile.OsDisk.OsType, *vm.Name)

		err := upgradeMasterNode.DeleteNode(vm.Name)
		if err != nil {
			ku.TraceLogger.Infof("Error deleting master VM: %s, err: %v\n", *vm.Name, err)
			return err
		}

		err = upgradeMasterNode.CreateNode("master", masterIndex)
		if err != nil {
			ku.TraceLogger.Infof("Error creating upgraded master VM: %s\n", *vm.Name)
			return err
		}

		err = upgradeMasterNode.Validate()
		if err != nil {
			ku.TraceLogger.Infof("Error validating upgraded master VM: %s\n", *vm.Name)
			return err
		}

		upgradedMastersIndex[masterIndex] = true
	}

	// This condition is possible if the previous upgrade operation failed during master
	// VM upgrade when a master VM was deleted but creation of upgraded master did not run.
	if masterNodesInCluster < expectedMasterCount {
		ku.TraceLogger.Infof(
			"Found missing master VMs in the cluster. Reconstructing names of missing master VMs for recreation during upgrade...\n")
	}

	mastersToCreate := expectedMasterCount - masterNodesInCluster
	ku.TraceLogger.Infof("Expected master count: %d, Creating %d more master VMs\n", expectedMasterCount, mastersToCreate)

	// NOTE: this is NOT completely idempotent because it assumes that
	// the OS disk has been deleted
	for i := 0; i < mastersToCreate; i++ {
		masterIndexToCreate := 0
		for upgradedMastersIndex[masterIndexToCreate] == true {
			masterIndexToCreate++
		}

		ku.TraceLogger.Infof("Creating upgraded master VM with index: %d\n", masterIndexToCreate)

		err = upgradeMasterNode.CreateNode("master", masterIndexToCreate)
		if err != nil {
			ku.TraceLogger.Infof("Error creating upgraded master VM with index: %d\n", masterIndexToCreate)
			return err
		}

		err = upgradeMasterNode.Validate()
		if err != nil {
			ku.TraceLogger.Infof("Error validating upgraded master VM with index: %d\n", masterIndexToCreate)
			return err
		}

		upgradedMastersIndex[masterIndexToCreate] = true
	}

	return nil
}

func (ku *Upgrader) upgradeAgentPools() error {
	// Unused until safely drain node is being called
	// var kubeAPIServerURL string
	// if ku.DataModel.Properties.MasterProfile != nil {
	// 	kubeAPIServerURL = ku.DataModel.Properties.MasterProfile.FQDN
	// }
	// if ku.DataModel.Properties.HostedMasterProfile != nil {
	// 	kubeAPIServerURL = ku.DataModel.Properties.HostedMasterProfile.FQDN
	// }
	for _, agentPool := range ku.ClusterTopology.AgentPools {
		// Upgrade Agent VMs
		templateMap, parametersMap, err := ku.generateUpgradeTemplate(ku.ClusterTopology.DataModel)
		if err != nil {
			return ku.Translator.Errorf("error generating upgrade template: %s", err.Error())
		}

		ku.TraceLogger.Infof("Prepping agent pool: %s for upgrade...\n", *agentPool.Name)

		preservePools := map[string]bool{*agentPool.Name: true}
		transformer := &acsengine.Transformer{
			Translator: ku.Translator,
		}
		var isMasterManagedDisk bool
		if ku.DataModel.Properties.MasterProfile != nil {
			isMasterManagedDisk = ku.DataModel.Properties.MasterProfile.IsManagedDisks()
		}
		if err := transformer.NormalizeResourcesForK8sAgentUpgrade(ku.TraceLogger, templateMap, isMasterManagedDisk, preservePools); err != nil {
			ku.TraceLogger.Errorf(err.Error())
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
			Translator:  ku.Translator,
			TraceLogger: ku.TraceLogger,
		}
		upgradeAgentNode.TemplateMap = templateMap
		upgradeAgentNode.ParametersMap = parametersMap
		upgradeAgentNode.UpgradeContainerService = ku.ClusterTopology.DataModel
		upgradeAgentNode.ResourceGroup = ku.ClusterTopology.ResourceGroup
		upgradeAgentNode.Client = ku.Client

		upgradedAgentsIndex := make(map[int]bool)

		for _, vm := range *agentPool.UpgradedAgentVMs {
			ku.TraceLogger.Infof("Agent VM: %s, pool name: %s on expected orchestrator version\n", *vm.Name, *agentPool.Name)
			agentIndex, _ := armhelpers.GetVMNameIndex(vm.StorageProfile.OsDisk.OsType, *vm.Name)
			upgradedAgentsIndex[agentIndex] = true
		}

		ku.TraceLogger.Infof("Starting upgrade of agent nodes in pool identifier: %s, name: %s...\n",
			*agentPool.Identifier, *agentPool.Name)

		for _, vm := range *agentPool.AgentVMs {
			ku.TraceLogger.Infof("Upgrading Agent VM: %s, pool name: %s\n", *vm.Name, *agentPool.Name)

			agentIndex, _ := armhelpers.GetVMNameIndex(vm.StorageProfile.OsDisk.OsType, *vm.Name)

			// Currently in a sinlge node cluster the api server will not be running when this point is reached on the first node so it will always fail.
			// err := operations.SafelyDrainNode(ku.Client, log.New().WithField("operation", "upgrade"), kubeAPIServerURL, ku.kubeConfig, *vm.Name)
			// if err != nil {
			// 	ku.TraceLogger.Infof("Error draining agent VM: %s", *vm.Name))
			// 	return err
			// }

			err := upgradeAgentNode.DeleteNode(vm.Name)
			if err != nil {
				ku.TraceLogger.Infof("Error deleting agent VM: %s\n", *vm.Name)
				return err
			}

			err = upgradeAgentNode.CreateNode(*agentPool.Name, agentIndex)
			if err != nil {
				ku.TraceLogger.Infof("Error creating upgraded agent VM: %s\n", *vm.Name)
				return err
			}

			err = upgradeAgentNode.Validate()
			if err != nil {
				ku.TraceLogger.Infof("Error validating upgraded agent VM: %s\n", *vm.Name)
				return err
			}

			upgradedAgentsIndex[agentIndex] = true
		}

		agentsToCreate := agentCount - len(upgradedAgentsIndex)
		ku.TraceLogger.Infof("Expected agent count in the pool: %d, Creating %d more agents\n", agentCount, agentsToCreate)

		// NOTE: this is NOT completely idempotent because it assumes that
		// the OS disk has been deleted
		for i := 0; i < agentsToCreate; i++ {
			agentIndexToCreate := 0
			for upgradedAgentsIndex[agentIndexToCreate] == true {
				agentIndexToCreate++
			}

			ku.TraceLogger.Infof("Creating upgraded Agent VM with index: %d, pool name: %s\n", agentIndexToCreate, *agentPool.Name)

			err = upgradeAgentNode.CreateNode(*agentPool.Name, agentIndexToCreate)
			if err != nil {
				ku.TraceLogger.Infof("Error creating upgraded agent VM with index: %d\n", agentIndexToCreate)
				return err
			}

			err = upgradeAgentNode.Validate()
			if err != nil {
				ku.TraceLogger.Infof("Error validating upgraded agent VM with index: %d\n", agentIndexToCreate)
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

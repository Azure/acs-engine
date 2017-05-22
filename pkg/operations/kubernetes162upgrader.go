package operations

import (
	"encoding/json"
	"fmt"
	"sort"

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

	upgradeContainerService := ku.ClusterTopology.DataModel
	upgradeContainerService.Properties.OrchestratorProfile.OrchestratorVersion = api.Kubernetes162

	// Upgrade Master VMs
	templateMap, parametersMap, err := ku.generateUpgradeTemplate(upgradeContainerService)
	if err != nil {
		return fmt.Errorf("error generating upgrade template: %s", err.Error())
	}

	if err := acsengine.NormalizeResourcesForK8sMasterUpgrade(log.NewEntry(log.New()), templateMap, true); err != nil {
		log.Fatalln(err)
		return err
	}

	upgradeMasterNode := UpgradeMasterNode{}
	upgradeMasterNode.TemplateMap = templateMap
	upgradeMasterNode.ParametersMap = parametersMap
	upgradeMasterNode.UpgradeContainerService = upgradeContainerService
	upgradeMasterNode.ResourceGroup = ku.ClusterTopology.ResourceGroup
	upgradeMasterNode.Client = ku.Client

	// Sort by VM Name (e.g.: k8s-master-22551669-0) offset no. in descending order
	sort.Sort(sort.Reverse(armhelpers.ByVMNameOffset(*ku.ClusterTopology.MasterVMs)))

	log.Infoln(fmt.Sprintf("Starting master nodes upgrade..."))

	masterLoopCount := 1
	for _, vm := range *ku.ClusterTopology.MasterVMs {
		log.Infoln(fmt.Sprintf("Upgrading Master VM: %s", *vm.Name))

		// 1.	Shutdown and delete one master VM at a time while preserving the persistent disk backing etcd.
		upgradeMasterNode.DeleteNode(vm.Name)
		// 2.	Call CreateVMWithRetries
		upgradeMasterNode.CreateNode("master", masterLoopCount)

		upgradeMasterNode.Validate()

		masterLoopCount++
	}

	// Upgrade Agent VMs
	templateMap, parametersMap, err = ku.generateUpgradeTemplate(upgradeContainerService)
	if err != nil {
		return fmt.Errorf("error generating upgrade template: %s", err.Error())
	}

	if err := acsengine.NormalizeResourcesForK8sAgentUpgrade(log.NewEntry(log.New()), templateMap); err != nil {
		log.Fatalln(err)
		return err
	}

	upgradeAgentNode := UpgradeAgentNode{}
	upgradeAgentNode.TemplateMap = templateMap
	upgradeAgentNode.ParametersMap = parametersMap
	upgradeAgentNode.UpgradeContainerService = upgradeContainerService
	upgradeAgentNode.ResourceGroup = ku.ClusterTopology.ResourceGroup
	upgradeAgentNode.Client = ku.Client

	sort.Sort(sort.Reverse(armhelpers.ByVMNameOffset(*ku.ClusterTopology.AgentVMs)))

	log.Infoln(fmt.Sprintf("Starting agent nodes upgrade..."))

	// TODO: Upgrade one agent pool at a time
	// TODO: Enable upgrade of Windows agent pools
	agentLoopCount := 1
	for _, vm := range *ku.ClusterTopology.AgentVMs {
		// 1.	Shutdown and delete one agent VM at a time
		upgradeAgentNode.DeleteNode(vm.Name)

		// 2.	Call CreateVMWithRetries
		_, poolName, _, _, _ := armhelpers.VMNameParts(*vm.Name)
		log.Infoln(fmt.Sprintf("Upgrading Agent VM: %s, pool name: %s", *vm.Name, poolName))

		upgradeAgentNode.CreateNode(poolName, agentLoopCount)

		upgradeAgentNode.Validate()

		agentLoopCount++
	}

	return nil
}

// Validate will run validation post upgrade
func (ku *Kubernetes162upgrader) Validate() error {
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

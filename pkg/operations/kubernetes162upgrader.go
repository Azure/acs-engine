package operations

import (
	"encoding/json"
	"fmt"
	"path"

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
		return fmt.Errorf("Upgrade to Kubernetes 1.6.2 is not supported from version: %s", ku.DataModel.Properties.OrchestratorProfile.OrchestratorVersion)
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

	templateGenerator, err := acsengine.InitializeTemplateGenerator(false)
	if err != nil {
		return fmt.Errorf("failed to initialize template generator: %s", err.Error())
	}

	var templateJSON string
	var parametersJSON string
	if templateJSON, parametersJSON, _, err = templateGenerator.GenerateTemplate(upgradeContainerService); err != nil {
		return fmt.Errorf("error generating upgrade template: %s", err.Error())
	}

	var template interface{}
	var parameters interface{}
	json.Unmarshal([]byte(templateJSON), &template)
	json.Unmarshal([]byte(parametersJSON), &parameters)
	templateMap := template.(map[string]interface{})
	// parametersMap := parameters.(map[string]interface{})

	log.Infoln(fmt.Sprintf("RunUpgrade 2"))

	templateVariables := templateMap["variables"].(map[string]interface{})

	masterCount, _ := templateVariables["masterCount"]
	masterCountInt := int(masterCount.(float64))
	log.Infoln(fmt.Sprintf("Master count: %d", masterCountInt))

	masterOffset, _ := templateVariables["masterOffset"]
	log.Infoln(fmt.Sprintf("Master offset: %v", masterOffset))

	loopCount := 1

	// for _, vm := range *ku.ClusterTopology.MasterVMs {
	// upgradeMasterNode := UpgradeMasterNode{}

	// 1.	Shutdown and delete one master VM at a time while preserving the persistent disk backing etcd.

	// 2.	Call CreateVMWithRetries

	// log.Infoln(fmt.Sprintf("Master VM: %v", vm))

	templateVariables["masterOffset"] = masterCountInt - loopCount
	masterOffset, _ = templateVariables["masterOffset"]
	log.Infoln(fmt.Sprintf("Master offset: %v", masterOffset))

	if e := acsengine.NormalizeResourcesForK8sMasterUpgrade(log.NewEntry(log.New()), templateMap); e != nil {
		log.Fatalln(err)
	}

	// ************************
	output, _ := json.Marshal(templateMap)
	var templateapp, parametersapp string
	if templateapp, err = acsengine.PrettyPrintArmTemplate(string(output)); err != nil {
		log.Fatalf("error pretty printing template: %s \n", err.Error())
	}
	if parametersapp, err = acsengine.PrettyPrintJSON(parametersJSON); err != nil {
		log.Fatalf("error pretty printing template parameters: %s \n", err.Error())
	}
	outputDirectory := path.Join("_output", upgradeContainerService.Properties.MasterProfile.DNSPrefix, "Upgrade")
	if err = acsengine.WriteArtifacts(upgradeContainerService, "vlabs", templateapp, parametersapp, outputDirectory, false, false); err != nil {
		log.Fatalf("error writing artifacts: %s \n", err.Error())
	}

	// ************************

	// loopCount++

	// var random *rand.Rand
	// deploymentSuffix := random.Int31()

	// _, err = ku.Client.DeployTemplate(
	// 	ku.ClusterTopology.ResourceGroup,
	// 	fmt.Sprintf("%s-%d", ku.ResourceGroup, deploymentSuffix),
	// 	templateMap,
	// 	parametersMap,
	// 	nil)

	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// }

	return nil
}

// Validate will run validation post upgrade
func (ku *Kubernetes162upgrader) Validate() error {
	return nil
}

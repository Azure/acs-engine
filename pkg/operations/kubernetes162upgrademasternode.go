package operations

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"path"
	"time"

	"github.com/Azure/acs-engine/pkg/acsengine"
	"github.com/Azure/acs-engine/pkg/api"
	"github.com/Azure/acs-engine/pkg/armhelpers"
	log "github.com/Sirupsen/logrus"
)

// Compiler to verify QueueMessageProcessor implements OperationsProcessor
var _ UpgradeNode = &UpgradeMasterNode{}

// UpgradeMasterNode upgrades a Kubernetes 1.5.3 master node to 1.6.2
type UpgradeMasterNode struct {
	TemplateMap             map[string]interface{}
	ParametersMap           map[string]interface{}
	UpgradeContainerService *api.ContainerService
	ResourceGroup           string
	Client                  armhelpers.ACSEngineClient
}

// DeleteNode takes state/resources of the master/agent node from ListNodeResources
// backs up/preserves state as needed by a specific version of Kubernetes and then deletes
// the node
func (kmn *UpgradeMasterNode) DeleteNode(vmName *string) error {
	if err := CleanDeleteVirtualMachine(kmn.Client, kmn.ResourceGroup, *vmName); err != nil {
		log.Fatalln(err)
		return err
	}

	return nil
}

// CreateNode creates a new master/agent node with the targeted version of Kubernetes
func (kmn *UpgradeMasterNode) CreateNode(countForOffset int) error {
	templateVariables := kmn.TemplateMap["variables"].(map[string]interface{})
	masterCount, _ := templateVariables["masterCount"]
	masterCountInt := int(masterCount.(float64))

	// Call CreateVMWithRetries
	templateVariables["masterOffset"] = masterCountInt - countForOffset
	masterOffset, _ := templateVariables["masterOffset"]
	log.Infoln(fmt.Sprintf("Master offset: %v", masterOffset))

	if err := acsengine.NormalizeResourcesForK8sMasterUpgrade(log.NewEntry(log.New()), kmn.TemplateMap); err != nil {
		log.Fatalln(err)
		return err
	}

	// ***********Save master update template*************
	updatedTemplateJSON, _ := json.Marshal(kmn.TemplateMap)
	parametersJSON, _ := json.Marshal(kmn.ParametersMap)

	templateapp, err := acsengine.PrettyPrintArmTemplate(string(updatedTemplateJSON))
	if err != nil {
		log.Fatalf("error pretty printing template: %s \n", err.Error())
	}
	parametersapp, e := acsengine.PrettyPrintJSON(string(parametersJSON))
	if e != nil {
		log.Fatalf("error pretty printing template parameters: %s \n", e.Error())
	}
	outputDirectory := path.Join("_output", kmn.UpgradeContainerService.Properties.MasterProfile.DNSPrefix, "Upgrade")
	if err := acsengine.WriteArtifacts(kmn.UpgradeContainerService, "vlabs", templateapp, parametersapp, outputDirectory, false, false); err != nil {
		log.Fatalf("error writing artifacts: %s \n", err.Error())
	}
	// ************************

	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	deploymentSuffix := random.Int31()

	_, err = kmn.Client.DeployTemplate(
		kmn.ResourceGroup,
		fmt.Sprintf("%s-%d", kmn.ResourceGroup, deploymentSuffix),
		kmn.TemplateMap,
		kmn.ParametersMap,
		nil)

	if err != nil {
		log.Fatalln(err)
	}

	return nil
}

// Validate will verify the that master/agent node has been upgraded as expected.
func (kmn *UpgradeMasterNode) Validate() error {
	return nil
}

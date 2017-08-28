package kubernetesupgrade

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/Azure/acs-engine/pkg/api"
	"github.com/Azure/acs-engine/pkg/armhelpers"
	"github.com/Azure/acs-engine/pkg/i18n"
	"github.com/Azure/acs-engine/pkg/operations"
	log "github.com/sirupsen/logrus"
)

// Compiler to verify QueueMessageProcessor implements OperationsProcessor
var _ UpgradeNode = &UpgradeMasterNode{}

// UpgradeMasterNode upgrades a Kubernetes 1.5 master node to 1.6
type UpgradeMasterNode struct {
	Translator              *i18n.Translator
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
	if err := operations.CleanDeleteVirtualMachine(kmn.Client, log.NewEntry(log.New()), kmn.ResourceGroup, *vmName); err != nil {
		return err
	}

	return nil
}

// CreateNode creates a new master/agent node with the targeted version of Kubernetes
func (kmn *UpgradeMasterNode) CreateNode(poolName string, masterNo int) error {
	templateVariables := kmn.TemplateMap["variables"].(map[string]interface{})

	templateVariables["masterOffset"] = masterNo
	masterOffsetVar, _ := templateVariables["masterOffset"]
	log.Infoln(fmt.Sprintf("Master offset: %v", masterOffsetVar))

	templateVariables["masterCount"] = masterNo + 1
	masterOffset, _ := templateVariables["masterCount"]
	log.Infoln(fmt.Sprintf("Master pool set count to: %v temporarily during upgrade...", masterOffset))

	WriteTemplate(kmn.Translator, kmn.UpgradeContainerService, kmn.TemplateMap, kmn.ParametersMap)

	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	deploymentSuffix := random.Int31()

	_, err := kmn.Client.DeployTemplate(
		kmn.ResourceGroup,
		fmt.Sprintf("%s-%d", kmn.ResourceGroup, deploymentSuffix),
		kmn.TemplateMap,
		kmn.ParametersMap,
		nil)

	if err != nil {
		return err
	}

	return nil
}

// Validate will verify the that master/agent node has been upgraded as expected.
func (kmn *UpgradeMasterNode) Validate() error {
	return nil
}

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
var _ UpgradeNode = &UpgradeAgentNode{}

// UpgradeAgentNode upgrades a Kubernetes 1.5 agent node to 1.6
type UpgradeAgentNode struct {
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
func (kan *UpgradeAgentNode) DeleteNode(vmName *string) error {
	if err := operations.CleanDeleteVirtualMachine(kan.Client, log.NewEntry(log.New()), kan.ResourceGroup, *vmName); err != nil {
		return err
	}

	return nil
}

// CreateNode creates a new master/agent node with the targeted version of Kubernetes
func (kan *UpgradeAgentNode) CreateNode(poolName string, agentNo int) error {
	poolCountParameter := kan.ParametersMap[poolName+"Count"].(map[string]interface{})
	poolCountParameter["value"] = agentNo + 1
	agentCount, _ := poolCountParameter["value"]
	log.Infoln(fmt.Sprintf("Agent pool: %s, set count to: %d temporarily during upgrade. Upgrading agent: %d",
		poolName, agentCount, agentNo))

	poolOffsetVarName := poolName + "Offset"
	templateVariables := kan.TemplateMap["variables"].(map[string]interface{})
	templateVariables[poolOffsetVarName] = agentNo

	WriteTemplate(kan.Translator, kan.UpgradeContainerService, kan.TemplateMap, kan.ParametersMap)

	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	deploymentSuffix := random.Int31()

	_, err := kan.Client.DeployTemplate(
		kan.ResourceGroup,
		fmt.Sprintf("%s-%d", kan.ResourceGroup, deploymentSuffix),
		kan.TemplateMap,
		kan.ParametersMap,
		nil)

	if err != nil {
		return err
	}

	return nil
}

// Validate will verify the that master/agent node has been upgraded as expected.
func (kan *UpgradeAgentNode) Validate() error {
	return nil
}

package operations

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/Azure/acs-engine/pkg/api"
	"github.com/Azure/acs-engine/pkg/armhelpers"

	log "github.com/Sirupsen/logrus"
)

// Compiler to verify QueueMessageProcessor implements OperationsProcessor
var _ UpgradeNode = &UpgradeAgentNode{}

// UpgradeAgentNode upgrades a Kubernetes 1.5.3 agent node to 1.6.2
type UpgradeAgentNode struct {
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
	if err := CleanDeleteVirtualMachine(kan.Client, kan.ResourceGroup, *vmName); err != nil {
		log.Fatalln(err)
		return err
	}

	return nil
}

// CreateNode creates a new master/agent node with the targeted version of Kubernetes
func (kan *UpgradeAgentNode) CreateNode(poolName string, countForOffset int) error {
	poolCountParameter := kan.ParametersMap[poolName+"Count"].(map[string]interface{})
	agentCount, _ := poolCountParameter["value"]
	agentCountInt := int(agentCount.(float64))
	log.Infoln(fmt.Sprintf("Agent pool: %s, count: %d", poolName, agentCountInt))

	poolOffsetVarName := poolName + "Offset"
	templateVariables := kan.TemplateMap["variables"].(map[string]interface{})
	templateVariables[poolOffsetVarName] = agentCountInt - countForOffset
	agentOffset, _ := templateVariables[poolOffsetVarName]
	log.Infoln(fmt.Sprintf("Agent offset: %v", agentOffset))

	WriteTemplate(kan.UpgradeContainerService, kan.TemplateMap, kan.ParametersMap)

	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	deploymentSuffix := random.Int31()

	_, err := kan.Client.DeployTemplate(
		kan.ResourceGroup,
		fmt.Sprintf("%s-%d", kan.ResourceGroup, deploymentSuffix),
		kan.TemplateMap,
		kan.ParametersMap,
		nil)

	if err != nil {
		log.Fatalln(err)
	}

	return nil
}

// Validate will verify the that master/agent node has been upgraded as expected.
func (kan *UpgradeAgentNode) Validate() error {
	return nil
}

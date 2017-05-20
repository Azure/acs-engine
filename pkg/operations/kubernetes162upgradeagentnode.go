package operations

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/Azure/acs-engine/pkg/acsengine"
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
func (kma *UpgradeAgentNode) DeleteNode(vmName *string) error {
	if err := CleanDeleteVirtualMachine(kma.Client, kma.ResourceGroup, *vmName); err != nil {
		log.Fatalln(err)
		return err
	}

	return nil
}

// CreateNode creates a new master/agent node with the targeted version of Kubernetes
func (kma *UpgradeAgentNode) CreateNode(poolName string, countForOffset int) error {
	poolCountParameter := kma.ParametersMap[poolName+"Count"].(map[string]interface{})
	agentCount, _ := poolCountParameter["value"]
	agentCountInt := int(agentCount.(float64))
	log.Infoln(fmt.Sprintf("Agent pool: %s, count: %d", poolName, agentCountInt))

	poolOffsetVarName := poolName + "Offset"
	templateVariables := kma.TemplateMap["variables"].(map[string]interface{})
	templateVariables[poolOffsetVarName] = agentCountInt - countForOffset
	agentOffset, _ := templateVariables[poolOffsetVarName]
	log.Infoln(fmt.Sprintf("Agent offset: %v", agentOffset))

	if err := acsengine.NormalizeResourcesForK8sMasterUpgrade(log.NewEntry(log.New()), kma.TemplateMap); err != nil {
		log.Fatalln(err)
		return err
	}

	WriteTemplate(kma.UpgradeContainerService, kma.TemplateMap, kma.ParametersMap)

	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	deploymentSuffix := random.Int31()

	_, err := kma.Client.DeployTemplate(
		kma.ResourceGroup,
		fmt.Sprintf("%s-%d", kma.ResourceGroup, deploymentSuffix),
		kma.TemplateMap,
		kma.ParametersMap,
		nil)

	if err != nil {
		log.Fatalln(err)
	}

	return nil
}

// Validate will verify the that master/agent node has been upgraded as expected.
func (kma *UpgradeAgentNode) Validate() error {
	return nil
}

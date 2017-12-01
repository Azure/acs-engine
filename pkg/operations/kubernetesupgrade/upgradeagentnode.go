package kubernetesupgrade

import (
	"fmt"
	"math/rand"
	"time"

	"k8s.io/client-go/pkg/api/v1/node"

	"github.com/Azure/acs-engine/pkg/api"
	"github.com/Azure/acs-engine/pkg/armhelpers"
	"github.com/Azure/acs-engine/pkg/i18n"
	"github.com/Azure/acs-engine/pkg/operations"
	"github.com/sirupsen/logrus"
)

const (
	interval = time.Second * 1
	retry    = time.Second * 5
	timeout  = time.Minute * 10
)

// Compiler to verify QueueMessageProcessor implements OperationsProcessor
var _ UpgradeNode = &UpgradeAgentNode{}

// UpgradeAgentNode upgrades a Kubernetes 1.5 agent node to 1.6
type UpgradeAgentNode struct {
	Translator              *i18n.Translator
	logger                  *logrus.Entry
	TemplateMap             map[string]interface{}
	ParametersMap           map[string]interface{}
	UpgradeContainerService *api.ContainerService
	ResourceGroup           string
	Client                  armhelpers.ACSEngineClient
	kubeConfig              string
}

// DeleteNode takes state/resources of the master/agent node from ListNodeResources
// backs up/preserves state as needed by a specific version of Kubernetes and then deletes
// the node
// The 'drain' flag is used to invoke 'cordon and drain' flow.
func (kan *UpgradeAgentNode) DeleteNode(vmName *string, drain bool) error {
	if drain {
		var kubeAPIServerURL string

		if kan.UpgradeContainerService.Properties.HostedMasterProfile != nil {
			kubeAPIServerURL = kan.UpgradeContainerService.Properties.HostedMasterProfile.FQDN
		} else {
			kubeAPIServerURL = kan.UpgradeContainerService.Properties.MasterProfile.FQDN
		}

		err := operations.SafelyDrainNode(kan.Client, kan.logger, kubeAPIServerURL, kan.kubeConfig, *vmName, time.Minute)
		if err != nil {
			kan.logger.Warningf("Error draining agent VM %s. Proceeding with deletion. Error: %v", *vmName, err)
			// Proceed with deletion anyways
		}
	}
	if err := operations.CleanDeleteVirtualMachine(kan.Client, kan.logger, kan.ResourceGroup, *vmName); err != nil {
		return err
	}
	return nil
}

// CreateNode creates a new master/agent node with the targeted version of Kubernetes
func (kan *UpgradeAgentNode) CreateNode(poolName string, agentNo int) error {
	poolCountParameter := kan.ParametersMap[poolName+"Count"].(map[string]interface{})
	poolCountParameter["value"] = agentNo + 1
	agentCount, _ := poolCountParameter["value"]
	kan.logger.Infof("Agent pool: %s, set count to: %d temporarily during upgrade. Upgrading agent: %d",
		poolName, agentCount, agentNo)

	poolOffsetVarName := poolName + "Offset"
	templateVariables := kan.TemplateMap["variables"].(map[string]interface{})
	templateVariables[poolOffsetVarName] = agentNo

	// Debug function - keep commented out
	// WriteTemplate(kan.Translator, kan.UpgradeContainerService, kan.TemplateMap, kan.ParametersMap)

	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	deploymentSuffix := random.Int31()
	deploymentName := fmt.Sprintf("agent-%s-%d", time.Now().Format("06-01-02T15.04.05"), deploymentSuffix)

	_, err := kan.Client.DeployTemplate(
		kan.ResourceGroup,
		deploymentName,
		kan.TemplateMap,
		kan.ParametersMap,
		nil)

	if err != nil {
		return err
	}

	return nil
}

// Validate will verify that agent node has been upgraded as expected.
func (kan *UpgradeAgentNode) Validate(vmName *string) error {
	if vmName == nil || *vmName == "" {
		kan.logger.Warningf("VM name was empty. Skipping node condition check")
		return nil
	}

	var masterURL string
	if kan.UpgradeContainerService.Properties.HostedMasterProfile != nil {
		masterURL = kan.UpgradeContainerService.Properties.HostedMasterProfile.FQDN
	} else {
		masterURL = kan.UpgradeContainerService.Properties.MasterProfile.FQDN
	}

	client, err := kan.Client.GetKubernetesClient(masterURL, kan.kubeConfig, interval, timeout)
	if err != nil {
		return err
	}

	retryTimer := time.NewTimer(time.Millisecond)
	timeoutTimer := time.NewTimer(timeout)
	for {
		select {
		case <-timeoutTimer.C:
			retryTimer.Stop()
			err := fmt.Errorf("Node was not ready within %v", timeout)
			kan.logger.Errorf(err.Error())
			return err
		case <-retryTimer.C:
			agentNode, err := client.GetNode(*vmName)
			if err != nil {
				kan.logger.Infof("Agent VM: %s status error: %v", *vmName, err)
				retryTimer.Reset(retry)
			} else if node.IsNodeReady(agentNode) {
				kan.logger.Infof("Agent VM: %s is ready", *vmName)
				timeoutTimer.Stop()
				return nil
			} else {
				kan.logger.Infof("Agent VM: %s not ready yet...", *vmName)
				retryTimer.Reset(retry)
			}
		}
	}
}

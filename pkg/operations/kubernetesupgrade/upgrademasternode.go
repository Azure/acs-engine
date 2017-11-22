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

// Compiler to verify QueueMessageProcessor implements OperationsProcessor
var _ UpgradeNode = &UpgradeMasterNode{}

// UpgradeMasterNode upgrades a Kubernetes 1.5 master node to 1.6
type UpgradeMasterNode struct {
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
// the node.
// The 'drain' flag is not used for deleting master nodes.
func (kmn *UpgradeMasterNode) DeleteNode(vmName *string, drain bool) error {
	if err := operations.CleanDeleteVirtualMachine(kmn.Client, kmn.logger, kmn.ResourceGroup, *vmName); err != nil {
		return err
	}

	return nil
}

// CreateNode creates a new master/agent node with the targeted version of Kubernetes
func (kmn *UpgradeMasterNode) CreateNode(poolName string, masterNo int) error {
	templateVariables := kmn.TemplateMap["variables"].(map[string]interface{})

	templateVariables["masterOffset"] = masterNo
	masterOffsetVar, _ := templateVariables["masterOffset"]
	kmn.logger.Infof("Master offset: %v\n", masterOffsetVar)

	templateVariables["masterCount"] = masterNo + 1
	masterOffset, _ := templateVariables["masterCount"]
	kmn.logger.Infof("Master pool set count to: %v temporarily during upgrade...\n", masterOffset)

	// Debug function - keep commented out
	// WriteTemplate(kmn.Translator, kmn.UpgradeContainerService, kmn.TemplateMap, kmn.ParametersMap)

	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	deploymentSuffix := random.Int31()
	deploymentName := fmt.Sprintf("master-%s-%d", time.Now().Format("06-01-02T15.04.05"), deploymentSuffix)

	_, err := kmn.Client.DeployTemplate(
		kmn.ResourceGroup,
		deploymentName,
		kmn.TemplateMap,
		kmn.ParametersMap,
		nil)

	if err != nil {
		return err
	}

	return nil
}

// Validate will verify the that master node has been upgraded as expected.
func (kmn *UpgradeMasterNode) Validate(vmName *string) error {
	if vmName == nil || *vmName == "" {
		kmn.logger.Warningf("VM name was empty. Skipping node condition check")
		return nil
	}

	if kmn.UpgradeContainerService.Properties.MasterProfile == nil {
		kmn.logger.Warningf("Master profile was empty. Skipping node condition check")
		return nil
	}

	masterURL := kmn.UpgradeContainerService.Properties.MasterProfile.FQDN

	client, err := kmn.Client.GetKubernetesClient(masterURL, kmn.kubeConfig, interval, timeout)
	if err != nil {
		return err
	}

	ch := make(chan struct{}, 1)
	go func() {
		for {
			masterNode, err := client.GetNode(*vmName)
			if err != nil {
				kmn.logger.Infof("Master VM: %s status error: %v\n", *vmName, err)
				time.Sleep(time.Second * 5)
			} else if node.IsNodeReady(masterNode) {
				kmn.logger.Infof("Master VM: %s is ready", *vmName)
				ch <- struct{}{}
			} else {
				kmn.logger.Infof("Master VM: %s not ready yet...", *vmName)
				time.Sleep(time.Second * 5)
			}
		}
	}()

	for {
		select {
		case <-ch:
			return nil
		case <-time.After(timeout):
			kmn.logger.Errorf("Node was not ready within %v", timeout)
			return fmt.Errorf("Node was not ready within %v", timeout)
		}
	}
}

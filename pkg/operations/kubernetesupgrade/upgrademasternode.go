package kubernetesupgrade

import (
	"fmt"
	"math/rand"
	"time"

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
}

// DeleteNode takes state/resources of the master/agent node from ListNodeResources
// backs up/preserves state as needed by a specific version of Kubernetes and then deletes
// the node
func (kmn *UpgradeMasterNode) DeleteNode(vmName *string) error {
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

// Validate will verify the that master/agent node has been upgraded as expected.
func (kmn *UpgradeMasterNode) Validate(vmName *string) error {
	if vmName == nil || *vmName == "" {
		kan.logger.Warningf(fmt.Sprintf("VM name was empty. Skipping node condition check"))
		return nil
	}

	var masterURL string
	if kan.UpgradeContainerService.Properties.HostedMasterProfile != nil {
		masterURL = kan.UpgradeContainerService.Properties.HostedMasterProfile.FQDN
	} else {
		masterURL = kan.UpgradeContainerService.Properties.MasterProfile.FQDN
	}

	if masterURL == "" {
		kan.Translator.Errorf("Control plane FQDN was not set.")
	}

	client, err := kan.Client.GetKubernetesClient(masterURL, kan.kubeConfig, interval, timeout)
	if err != nil {
		return err
	}
	return nil
}

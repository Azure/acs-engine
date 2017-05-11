package operations

import (
	"fmt"

	"strings"

	"github.com/Azure/acs-engine/pkg/api"
	"github.com/Azure/acs-engine/pkg/api/vlabs"
	"github.com/Azure/acs-engine/pkg/operations/armhelpers"
	"github.com/Azure/azure-sdk-for-go/arm/compute"
	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/prometheus/common/log"
	"github.com/satori/go.uuid"
)

// ClusterTopology contains resources of the cluster the upgrade operation
// is targeting
type ClusterTopology struct {
	APIModel  *api.ContainerService
	MasterVMs *[]compute.VirtualMachine
	AgentVMs  *[]compute.VirtualMachine
}

// UpgradeCluster upgrades a cluster with Orchestrator version X
// (or X.X or X.X.X) to version y (or Y.Y or X.X.X). RIght now
// upgrades are supported for Kubernetes cluster only.
type UpgradeCluster struct {
	ClusterTopology
	AzureClients armhelpers.AzureClients
}

// UpgradeCluster runs the workflow to upgrade a Kubernetes cluster.
// UpgradeContainerService contains target state of the cluster that
// the operation will drive towards.
func (uc *UpgradeCluster) UpgradeCluster(subscriptionID uuid.UUID, resourceGroup string,
	cs *api.ContainerService, ucs *vlabs.UpgradeContainerService, token *adal.ServicePrincipalToken) {
	uc.ClusterTopology = ClusterTopology{}
	uc.APIModel = cs

	AzureClients := armhelpers.AzureClients{
		SubscriptionID: subscriptionID.String(),
	}
	AzureClients.Create(token)

	if err := uc.getUpgradableResources(subscriptionID, resourceGroup); err != nil {
		// Bail
		return
	}
}

func (uc *UpgradeCluster) getUpgradableResources(subscriptionID uuid.UUID, resourceGroup string) error {
	vmListResult, err := uc.AzureClients.VMClient.List(resourceGroup)
	if err != nil {
		return err
	}

	orchestratorTypeVersion := fmt.Sprintf("%s:%s", uc.APIModel.Properties.OrchestratorProfile.OrchestratorType,
		uc.APIModel.Properties.OrchestratorProfile.OrchestratorVersion)

	for _, vm := range *vmListResult.Value {
		if *(*vm.Tags)["orchestrator"] == orchestratorTypeVersion {
			if strings.Contains(*(vm.Name), "k8s-master-") {
				// TODO: *vm.Tags["resourceNameSuffix"] ==  Read VM NAME SUFFIX from temp parameter
				*uc.MasterVMs = append(*uc.MasterVMs, vm)
			}
			// TODO: Add logic to separate out VMs in various agent pookls
			if strings.Contains(*(vm.Name), "k8s-agentpool-") {
				// TODO: *vm.Tags["resourceNameSuffix"] ==  Read VM NAME SUFFIX from temp parameter
				*uc.AgentVMs = append(*uc.AgentVMs, vm)
			}
		}
	}

	log.Infoln("Master VMs: %+v", *uc.MasterVMs)
	log.Infoln("Agent VMs: %+v", *uc.AgentVMs)

	return nil
}

// UpgradeWorkFlow outlines various individual high level steps
// that need to be run (one or more times) in the upgrade workflow.
type UpgradeWorkFlow interface {
	ClusterPreflightCheck()

	// upgrade masters
	// upgrade agent nodes
	RunUpgrade() error

	Validate() error
}

// UpgradeNode drives work flow of deleting and replacing a master or agent node to a
// specified target version of Kubernetes
type UpgradeNode interface {
	// DeleteNode takes state/resources of the master/agent node from ListNodeResources
	// backs up/preserves state as needed by a specific version of Kubernetes and then deletes
	// the node
	DeleteNode() error

	// CreateNode creates a new master/agent node with the targeted version of Kubernetes
	CreateNode() error

	// Validate will verify the that master/agent node has been upgraded as expected.
	Validate() error
}

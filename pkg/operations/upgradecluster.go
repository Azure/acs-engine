package operations

import (
	"fmt"

	"strings"

	"github.com/Azure/acs-engine/pkg/api"
	"github.com/Azure/acs-engine/pkg/operations/armhelpers"
	"github.com/Azure/azure-sdk-for-go/arm/compute"
	log "github.com/Sirupsen/logrus"
	"github.com/satori/go.uuid"
)

// ClusterTopology contains resources of the cluster the upgrade operation
// is targeting
type ClusterTopology struct {
	DataModel *api.ContainerService
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
	cs *api.ContainerService, ucs *api.UpgradeContainerService) error {
	uc.ClusterTopology = ClusterTopology{}
	uc.DataModel = cs
	uc.MasterVMs = &[]compute.VirtualMachine{}
	uc.AgentVMs = &[]compute.VirtualMachine{}

	if err := uc.getUpgradableResources(subscriptionID, resourceGroup); err != nil {
		return fmt.Errorf("Error while querying ARM for resources: %+v", err)
	}

	switch ucs.OrchestratorProfile.OrchestratorVersion {
	case api.Kubernetes162:
		upgrader := Kubernetes162upgrader{}
		upgrader.RunUpgrade()
	default:
		return fmt.Errorf("Upgrade to Kubernetes 1.6.2 is not supported from version: %s", uc.DataModel.Properties.OrchestratorProfile.OrchestratorVersion)
	}

	return nil
}

func (uc *UpgradeCluster) getUpgradableResources(subscriptionID uuid.UUID, resourceGroup string) error {
	vmListResult, err := uc.AzureClients.VMClient.List(resourceGroup)
	if err != nil {
		return err
	}

	orchestratorTypeVersion := fmt.Sprintf("%s:%s", uc.DataModel.Properties.OrchestratorProfile.OrchestratorType,
		uc.DataModel.Properties.OrchestratorProfile.OrchestratorVersion)

	for _, vm := range *vmListResult.Value {
		if *(*vm.Tags)["orchestrator"] == orchestratorTypeVersion {
			if strings.Contains(*(vm.Name), "k8s-master-") {
				log.Infoln(fmt.Sprintf("Master VM name: %s", *vm.Name))
				// TODO: *vm.Tags["resourceNameSuffix"] ==  Read VM NAME SUFFIX from temp parameter
				*uc.MasterVMs = append(*uc.MasterVMs, vm)
			}
			// TODO: Add logic to separate out VMs in various agent pookls
			if strings.Contains(*(vm.Name), "k8s-agentpool") {
				log.Infoln(fmt.Sprintf("Agent VM name: %s", *vm.Name))
				// TODO: *vm.Tags["resourceNameSuffix"] ==  Read VM NAME SUFFIX from temp parameter
				*uc.AgentVMs = append(*uc.AgentVMs, vm)
			}
		}
	}

	return nil
}

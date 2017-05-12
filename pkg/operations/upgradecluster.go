package operations

import (
	"fmt"
	"strings"

	"github.com/Azure/acs-engine/pkg/api"
	"github.com/Azure/acs-engine/pkg/armhelpers"
	"github.com/Azure/go-autorest/autorest/to"

	"github.com/Azure/azure-sdk-for-go/arm/compute"
	"github.com/prometheus/common/log"
)

// ClusterTopology contains resources of the cluster the upgrade operation
// is targeting
type ClusterTopology struct {
	APIModel      *api.ContainerService
	ResourceGroup string
}

// UpgradeCluster upgrades a cluster with Orchestrator version X
// (or X.X or X.X.X) to version y (or Y.Y or X.X.X). RIght now
// upgrades are supported for Kubernetes cluster only.
type UpgradeCluster struct {
	ClusterTopology
	UpgradeModel *api.UpgradeContainerService
	AzureClient  *armhelpers.AzureClient
}

// UpgradeCluster runs the workflow to upgrade a Kubernetes cluster.
// UpgradeContainerService contains target state of the cluster that
// the operation will drive towards.
func (uc *UpgradeCluster) UpgradeCluster(resourceGroup string,
	cs *api.ContainerService, ucs *api.UpgradeContainerService) error {
	uc.ClusterTopology = ClusterTopology{
		APIModel:      cs,
		ResourceGroup: resourceGroup,
	}
	uc.UpgradeModel = ucs
	uc.APIModel = cs

	masterVMs, agentVMs, err := uc.getUpgradableResources(resourceGroup)
	if err != nil {
		log.Errorln("Error while querying ARM for resources: %+v", err)
		return err
	}

	for _, masterVM := range masterVMs {
		log.Infoln("Upgrade master:", to.String(masterVM.Name))

		uc.AzureClient.TemplateDeployer.DeployTemplate()

		// TODO:
		// generate template
		// take extra step of removing the agent, since apparently we can't fully suppress it
		// -- alternatively, we put something in the base template to skip agnet pools entirely if the thing passed ot it has none...
	}

	for _, agentVM := range agentVMs {
		log.Infoln("Upgrade agent:", to.String(agentVM.Name))
	}

	return nil
}

func (uc *UpgradeCluster) getUpgradableResources(resourceGroup string) ([]compute.VirtualMachine, []compute.VirtualMachine, error) {
	vmListResult, err := uc.AzureClient.VirtualMachinesClient.List(resourceGroup)
	if err != nil {
		return nil, nil, err
	}

	orchestratorTypeVersion := fmt.Sprintf("%s:%s", uc.APIModel.Properties.OrchestratorProfile.OrchestratorType,
		uc.APIModel.Properties.OrchestratorProfile.OrchestratorVersion)

	masterVMs := []compute.VirtualMachine{}
	agentVMs := []compute.VirtualMachine{}

	for _, vm := range *vmListResult.Value {
		if *(*vm.Tags)["orchestrator"] == orchestratorTypeVersion {
			if strings.Contains(*(vm.Name), "k8s-master-") {
				log.Infoln(fmt.Sprintf("Master VM name: %s", *vm.Name))
				// TODO: *vm.Tags["resourceNameSuffix"] ==  Read VM NAME SUFFIX from temp parameter
				// TODO(colemick): presumably we will filter based on the suffix hash here?
				masterVMs = append(masterVMs, vm)
			}
			// TODO: Add logic to separate out VMs in various agent pookls
			if strings.Contains(*(vm.Name), "k8s-agentpool") {
				log.Infoln(fmt.Sprintf("Agent VM name: %s", *vm.Name))
				// TODO: *vm.Tags["resourceNameSuffix"] ==  Read VM NAME SUFFIX from temp parameter
				agentVMs = append(agentVMs, vm)
			}
		}
	}

	return masterVMs, agentVMs, nil
}

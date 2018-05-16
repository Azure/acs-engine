package dcosupgrade

import (
	"fmt"
	"strings"

	"github.com/Azure/acs-engine/pkg/api"
	"github.com/Azure/acs-engine/pkg/armhelpers"
	"github.com/Azure/acs-engine/pkg/i18n"
	"github.com/Azure/azure-sdk-for-go/arm/compute"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

// ClusterTopology contains resources of the cluster the upgrade operation
// is targeting
type ClusterTopology struct {
	DataModel          *api.ContainerService
	Location           string
	ResourceGroup      string
	CurrentDcosVersion string
	NameSuffix         string
	SSHKey             []byte

	AgentPoolsToUpgrade map[string]bool
	AgentPools          map[string]*AgentPoolTopology

	BootstrapVMs      *[]compute.VirtualMachine
	MasterVMs         *[]compute.VirtualMachine
	AgentVMs          *[]compute.VirtualMachine
	UpgradedMasterVMs *[]compute.VirtualMachine
}

// AgentPoolTopology contains agent VMs in a single pool
type AgentPoolTopology struct {
	Identifier       *string
	Name             *string
	AgentVMs         *[]compute.VirtualMachine
	UpgradedAgentVMs *[]compute.VirtualMachine
}

// UpgradeCluster upgrades a cluster with Orchestrator version X.X to version Y.Y.
// Right now upgrades are supported for Kubernetes cluster only.
type UpgradeCluster struct {
	Translator *i18n.Translator
	Logger     *logrus.Entry
	ClusterTopology
	Client armhelpers.ACSEngineClient
}

// UpgradeCluster runs the workflow to upgrade a DCOS cluster.
func (uc *UpgradeCluster) UpgradeCluster(subscriptionID uuid.UUID, resourceGroup, currentDcosVersion string,
	cs *api.ContainerService, nameSuffix string, sshKey []byte) error {
	uc.ClusterTopology = ClusterTopology{}
	uc.ResourceGroup = resourceGroup
	uc.CurrentDcosVersion = currentDcosVersion
	uc.DataModel = cs
	uc.NameSuffix = nameSuffix
	uc.SSHKey = sshKey
	uc.BootstrapVMs = &[]compute.VirtualMachine{}
	uc.MasterVMs = &[]compute.VirtualMachine{}
	uc.AgentVMs = &[]compute.VirtualMachine{}
	uc.UpgradedMasterVMs = &[]compute.VirtualMachine{}
	uc.AgentPools = make(map[string]*AgentPoolTopology)
	uc.AgentPoolsToUpgrade = make(map[string]bool)

	for _, pool := range cs.Properties.AgentPoolProfiles {
		uc.AgentPoolsToUpgrade[pool.Name] = true
	}
	uc.AgentPoolsToUpgrade["master"] = true

	if err := uc.getClusterNodeStatus(subscriptionID, resourceGroup); err != nil {
		return uc.Translator.Errorf("Error while querying ARM for resources: %+v", err)
	}

	upgradeVersion := uc.DataModel.Properties.OrchestratorProfile.OrchestratorVersion
	uc.Logger.Infof("Upgrading to DCOS version %s", upgradeVersion)

	if err := uc.runUpgrade(); err != nil {
		return err
	}

	uc.Logger.Infof("Cluster upgraded successfully to DCOS %s", upgradeVersion)
	return nil
}

func (uc *UpgradeCluster) getClusterNodeStatus(subscriptionID uuid.UUID, resourceGroup string) error {
	vmListResult, err := uc.Client.ListVirtualMachines(resourceGroup)
	if err != nil {
		return err
	}
	bootstrapPrefix := fmt.Sprintf("bootstrap-%s-", uc.NameSuffix)
	masterPrefix := fmt.Sprintf("dcos-master-%s-", uc.NameSuffix)

	for _, vm := range *vmListResult.Value {

		if strings.HasPrefix(*(vm.Name), bootstrapPrefix) {
			uc.Logger.Infof("Bootstrap VM name: %s", *vm.Name)
			*uc.BootstrapVMs = append(*uc.BootstrapVMs, vm)
		} else if strings.HasPrefix(*(vm.Name), masterPrefix) {
			uc.Logger.Infof("Master VM name: %s", *vm.Name)
			*uc.MasterVMs = append(*uc.MasterVMs, vm)
		} else {
			uc.Logger.Infof("Agent VM name: %s", *vm.Name)
			*uc.AgentVMs = append(*uc.AgentVMs, vm)
		}
	}
	return nil
}

package kubernetesupgrade

import (
	"encoding/json"
	"fmt"
	"path"
	"strconv"

	"strings"

	"github.com/Azure/acs-engine/pkg/acsengine"
	"github.com/Azure/acs-engine/pkg/api"
	"github.com/Azure/acs-engine/pkg/armhelpers"
	"github.com/Azure/azure-sdk-for-go/arm/compute"
	log "github.com/Sirupsen/logrus"
	"github.com/satori/go.uuid"
)

// ClusterTopology contains resources of the cluster the upgrade operation
// is targeting
type ClusterTopology struct {
	DataModel     *api.ContainerService
	ResourceGroup string
	NameSuffix    string

	AgentPools map[string]*AgentPoolTopology

	MasterVMs         *[]compute.VirtualMachine
	UpgradedMasterVMs *[]compute.VirtualMachine
}

// AgentPoolTopology contains agent VMs in a single pool
type AgentPoolTopology struct {
	Identifier       *string
	Name             *string
	AgentVMs         *[]compute.VirtualMachine
	UpgradedAgentVMs *[]compute.VirtualMachine
}

// UpgradeCluster upgrades a cluster with Orchestrator version X
// (or X.X or X.X.X) to version y (or Y.Y or X.X.X). RIght now
// upgrades are supported for Kubernetes cluster only.
type UpgradeCluster struct {
	ClusterTopology
	Client armhelpers.ACSEngineClient

	UpgradeModel *api.UpgradeContainerService
}

// MasterVMNamePrefix is the prefix for all master VM names for Kubernetes clusters
const MasterVMNamePrefix = "k8s-master-"

// UpgradeCluster runs the workflow to upgrade a Kubernetes cluster.
// UpgradeContainerService contains target state of the cluster that
// the operation will drive towards.
func (uc *UpgradeCluster) UpgradeCluster(subscriptionID uuid.UUID, resourceGroup string,
	cs *api.ContainerService, ucs *api.UpgradeContainerService, nameSuffix string) error {
	uc.ClusterTopology = ClusterTopology{}
	uc.ResourceGroup = resourceGroup
	uc.DataModel = cs
	uc.NameSuffix = nameSuffix
	uc.MasterVMs = &[]compute.VirtualMachine{}
	uc.UpgradedMasterVMs = &[]compute.VirtualMachine{}
	uc.AgentPools = make(map[string]*AgentPoolTopology)

	uc.UpgradeModel = ucs

	if err := uc.getClusterNodeStatus(subscriptionID, resourceGroup); err != nil {
		return fmt.Errorf("Error while querying ARM for resources: %+v", err)
	}

	switch ucs.OrchestratorProfile.OrchestratorVersion {
	case api.Kubernetes162:
		log.Infoln(fmt.Sprintf("Upgrading to Kubernetes 1.6.2"))
		upgrader := Kubernetes162upgrader{}
		upgrader.ClusterTopology = uc.ClusterTopology
		upgrader.Client = uc.Client
		if err := upgrader.RunUpgrade(); err != nil {
			return err
		}
	default:
		return fmt.Errorf("Upgrade to Kubernetes version: %s is not supported from version: %s",
			ucs.OrchestratorProfile.OrchestratorVersion,
			uc.DataModel.Properties.OrchestratorProfile.OrchestratorVersion)
	}

	log.Infoln(fmt.Sprintf("Cluster upraded sucessfully to Kubernetes version: %s",
		ucs.OrchestratorProfile.OrchestratorVersion))
	return nil
}

func (uc *UpgradeCluster) getClusterNodeStatus(subscriptionID uuid.UUID, resourceGroup string) error {
	vmListResult, err := uc.Client.ListVirtualMachines(resourceGroup)
	if err != nil {
		return err
	}

	orchestratorTypeVersion := fmt.Sprintf("%s:%s", uc.DataModel.Properties.OrchestratorProfile.OrchestratorType,
		uc.DataModel.Properties.OrchestratorProfile.OrchestratorVersion)
	targetOrchestratorTypeVersion := fmt.Sprintf("%s:%s", uc.UpgradeModel.OrchestratorProfile.OrchestratorType,
		uc.UpgradeModel.OrchestratorProfile.OrchestratorVersion)

	for _, vm := range *vmListResult.Value {
		if vm.Tags == nil {
			log.Infoln(fmt.Sprintf("No tags found for VM: %s skipping.", *vm.Name))
			continue
		}

		vmOrchestratorTypeAndVersion := *(*vm.Tags)["orchestrator"]
		if vmOrchestratorTypeAndVersion == orchestratorTypeVersion {
			if strings.Contains(*(vm.Name), MasterVMNamePrefix) {
				if !strings.Contains(*(vm.Name), uc.NameSuffix) {
					log.Infoln(fmt.Sprintf("Skipping VM: %s for upgrade as it does not belong to cluster with expected name suffix: %s",
						*vm.Name, uc.NameSuffix))
					continue
				}
				log.Infoln(fmt.Sprintf("Master VM name: %s, orchestrator: %s (MasterVMs)", *vm.Name, vmOrchestratorTypeAndVersion))
				*uc.MasterVMs = append(*uc.MasterVMs, vm)
			} else {
				uc.addVMToAgentPool(vm, true)
			}
		} else if vmOrchestratorTypeAndVersion == targetOrchestratorTypeVersion {
			if !strings.Contains(*(vm.Name), uc.NameSuffix) {
				log.Infoln(fmt.Sprintf("Not adding VM: %s to upgraded list as it does not belong to cluster with expected name suffix: %s",
					*vm.Name, uc.NameSuffix))
				continue
			}
			if strings.Contains(*(vm.Name), MasterVMNamePrefix) {
				log.Infoln(fmt.Sprintf("Master VM name: %s, orchestrator: %s (UpgradedMasterVMs)", *vm.Name, vmOrchestratorTypeAndVersion))
				*uc.UpgradedMasterVMs = append(*uc.UpgradedMasterVMs, vm)
			} else {
				uc.addVMToAgentPool(vm, false)
			}
		}
	}

	return nil
}

func (uc *UpgradeCluster) addVMToAgentPool(vm compute.VirtualMachine, isUpgradableVM bool) error {
	var poolIdentifier string
	var poolPrefix string
	var err error
	if vm.StorageProfile.OsDisk.OsType == compute.Linux {
		_, poolIdentifier, poolPrefix, _, err = armhelpers.LinuxVMNameParts(*vm.Name)
		if err != nil {
			log.Errorln(err)
			return err
		}
	} else if vm.StorageProfile.OsDisk.OsType == compute.Windows {
		poolPrefix, acsStr, poolIndex, _, err := armhelpers.WindowsVMNameParts(*vm.Name)
		if err != nil {
			log.Errorln(err)
			return err
		}

		poolIdentifier = poolPrefix + acsStr + strconv.Itoa(poolIndex)
	}

	if !strings.Contains(uc.NameSuffix, poolPrefix) {
		log.Infoln(fmt.Sprintf("Skipping VM: %s for upgrade as it does not belong to cluster with expected name suffix: %s",
			*vm.Name, uc.NameSuffix))
		return nil
	}

	if uc.AgentPools[poolIdentifier] == nil {
		uc.AgentPools[poolIdentifier] =
			&AgentPoolTopology{&poolIdentifier, (*vm.Tags)["poolName"], &[]compute.VirtualMachine{}, &[]compute.VirtualMachine{}}
	}

	if isUpgradableVM {
		log.Infoln(fmt.Sprintf("Adding Agent VM: %s, orchestrator: %s to pool: %s (AgentVMs)",
			*vm.Name, *(*vm.Tags)["orchestrator"], poolIdentifier))
		*uc.AgentPools[poolIdentifier].AgentVMs = append(*uc.AgentPools[poolIdentifier].AgentVMs, vm)
	} else {
		log.Infoln(fmt.Sprintf("Adding Agent VM: %s, orchestrator: %s to pool: %s (UpgradedAgentVMs)",
			*vm.Name, *(*vm.Tags)["orchestrator"], poolIdentifier))
		*uc.AgentPools[poolIdentifier].UpgradedAgentVMs = append(*uc.AgentPools[poolIdentifier].UpgradedAgentVMs, vm)
	}

	return nil
}

// WriteTemplate writes upgrade template to a folder
func WriteTemplate(upgradeContainerService *api.ContainerService,
	templateMap map[string]interface{}, parametersMap map[string]interface{}) {
	updatedTemplateJSON, _ := json.Marshal(templateMap)
	parametersJSON, _ := json.Marshal(parametersMap)

	templateapp, err := acsengine.PrettyPrintArmTemplate(string(updatedTemplateJSON))
	if err != nil {
		log.Fatalf("error pretty printing template: %s \n", err.Error())
	}
	parametersapp, e := acsengine.PrettyPrintJSON(string(parametersJSON))
	if e != nil {
		log.Fatalf("error pretty printing template parameters: %s \n", e.Error())
	}
	outputDirectory := path.Join("_output", upgradeContainerService.Properties.MasterProfile.DNSPrefix, "Upgrade")
	if err := acsengine.WriteArtifacts(upgradeContainerService, "vlabs", templateapp, parametersapp, outputDirectory, false, false); err != nil {
		log.Fatalf("error writing artifacts: %s \n", err.Error())
	}
}

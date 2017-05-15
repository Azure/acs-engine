package Kubernetes162Upgrade

import (
	"fmt"

	"github.com/Azure/acs-engine/pkg/api"
	"github.com/Azure/acs-engine/pkg/operations"
	"github.com/Azure/acs-engine/pkg/operations/armhelpers"
)

// Compiler to verify QueueMessageProcessor implements OperationsProcessor
var _ operations.UpgradeWorkFlow = &UpgradeCluster{}

// UpgradeCluster upgrades a Kubernetes 1.5.3 cluster to 1.6.2
type UpgradeCluster struct {
	operations.ClusterTopology
	AzureClients armhelpers.AzureClients
}

// ClusterPreflightCheck dpes preflight check
func (mp *UpgradeCluster) ClusterPreflightCheck() error {
	// Check that current cluster is 1.5.3
	if mp.APIModel.Properties.OrchestratorProfile.OrchestratorVersion != api.Kubernetes153 {
		return fmt.Errorf("Upgrade to Kubernetes 1.6.2 is not supported from version: %s", mp.APIModel.Properties.OrchestratorProfile.OrchestratorVersion)
	}

	return nil
}

// RunUpgrade runs the upgrade pipeline
func (mp *UpgradeCluster) RunUpgrade() error {
	if err := mp.ClusterPreflightCheck(); err != nil {
		return err
	}

	// 1.	Shutdown and delete one master VM at a time while preserving the persistent disk backing etcd.
	// 2.	Call CreateVMWithRetries
	return nil
}

// Validate will run validation post upgrade
func (mp *UpgradeCluster) Validate() error {
	return nil
}

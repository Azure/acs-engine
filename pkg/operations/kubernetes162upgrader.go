package operations

import (
	"fmt"

	"github.com/Azure/acs-engine/pkg/api"
)

// Compiler to verify QueueMessageProcessor implements OperationsProcessor
var _ UpgradeWorkFlow = &Kubernetes162upgrader{}

// Kubernetes162upgrader upgrades a Kubernetes 1.5.3 cluster to 1.6.2
type Kubernetes162upgrader struct {
	ClusterTopology
}

// ClusterPreflightCheck dpes preflight check
func (mp *Kubernetes162upgrader) ClusterPreflightCheck() error {
	// Check that current cluster is 1.5.3
	if mp.DataModel.Properties.OrchestratorProfile.OrchestratorVersion != api.Kubernetes153 {
		return fmt.Errorf("Upgrade to Kubernetes 1.6.2 is not supported from version: %s", mp.DataModel.Properties.OrchestratorProfile.OrchestratorVersion)
	}

	return nil
}

// RunUpgrade runs the upgrade pipeline
func (mp *Kubernetes162upgrader) RunUpgrade() error {
	if err := mp.ClusterPreflightCheck(); err != nil {
		return err
	}

	// 1.	Shutdown and delete one master VM at a time while preserving the persistent disk backing etcd.
	// 2.	Call CreateVMWithRetries
	return nil
}

// Validate will run validation post upgrade
func (mp *Kubernetes162upgrader) Validate() error {
	return nil
}

package kubernetesupgrade

import (
	"fmt"

	"github.com/Azure/acs-engine/pkg/api"
)

// Compiler to verify QueueMessageProcessor implements OperationsProcessor
var _ UpgradeWorkFlow = &Kubernetes17upgrader{}

// Kubernetes17upgrader upgrades a Kubernetes 1.6 cluster to 1.7
type Kubernetes17upgrader struct {
	Kubernetes16upgrader
}

// ClusterPreflightCheck does preflight check
func (ku *Kubernetes17upgrader) ClusterPreflightCheck() error {
	// Check that current cluster is 1.6 or 1.7
	switch ku.DataModel.Properties.OrchestratorProfile.OrchestratorRelease {
	case api.KubernetesRelease1Dot6:
	case api.KubernetesRelease1Dot7:
	default:
		return fmt.Errorf("Upgrade to Kubernetes 1.7 is not supported from orchestrator release: %s",
			ku.DataModel.Properties.OrchestratorProfile.OrchestratorRelease)
	}
	return nil
}

// RunUpgrade runs the upgrade pipeline
func (ku *Kubernetes17upgrader) RunUpgrade() error {
	return ku.Kubernetes16upgrader.RunUpgrade()
}

// Validate will run validation post upgrade
func (ku *Kubernetes17upgrader) Validate() error {
	return nil
}

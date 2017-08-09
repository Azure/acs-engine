package kubernetesupgrade

import (
	"fmt"

	"github.com/Azure/acs-engine/pkg/api"
)

// Compiler to verify QueueMessageProcessor implements OperationsProcessor
var _ UpgradeWorkFlow = &Kubernetes16upgrader{}

// Kubernetes16upgrader upgrades a Kubernetes 1.5 cluster to 1.6
type Kubernetes16upgrader struct {
	Upgrader
}

// ClusterPreflightCheck does preflight check
func (ku *Kubernetes16upgrader) ClusterPreflightCheck() error {
	// Check that current cluster is 1.5 or 1.6
	switch ku.DataModel.Properties.OrchestratorProfile.OrchestratorRelease {
	case api.KubernetesRelease1Dot5:
	case api.KubernetesRelease1Dot6:
	default:
		return fmt.Errorf("Upgrade to Kubernetes 1.6 is not supported from orchestrator release: %s",
			ku.DataModel.Properties.OrchestratorProfile.OrchestratorRelease)
	}
	return nil
}

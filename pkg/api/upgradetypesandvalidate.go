package api

import (
	"fmt"

	"github.com/Azure/acs-engine/pkg/api/common"
)

// UpgradeContainerService API model
type UpgradeContainerService OrchestratorProfile

// Validate implements APIObject
func (ucs *UpgradeContainerService) Validate() error {
	switch ucs.OrchestratorType {
	case DCOS, SwarmMode, Swarm:
		return fmt.Errorf("Upgrade is not supported for orchestrator: %s", ucs.OrchestratorType)
	case Kubernetes:
		switch ucs.OrchestratorRelease {
		case common.KubernetesRelease1Dot6:
		case common.KubernetesRelease1Dot7:
		default:
			return fmt.Errorf("Upgrade is not supported to orchestrator release: %s", ucs.OrchestratorRelease)
		}
	}
	return nil
}

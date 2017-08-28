package vlabs

import (
	"fmt"

	"github.com/Azure/acs-engine/pkg/api/common"
)

// UpgradeContainerService API model
type UpgradeContainerService struct {
	OrchestratorProfile *OrchestratorProfile `json:"orchestratorProfile,omitempty"`
}

// Validate implements APIObject
func (ucs *UpgradeContainerService) Validate() error {
	switch ucs.OrchestratorProfile.OrchestratorType {
	case DCOS, SwarmMode, Swarm:
		return fmt.Errorf("Upgrade is not supported for orchestrator: %s", ucs.OrchestratorProfile.OrchestratorType)
	case Kubernetes:
		switch ucs.OrchestratorProfile.OrchestratorRelease {
		case common.KubernetesRelease1Dot6:
		case common.KubernetesRelease1Dot7:
		default:
			return fmt.Errorf("Upgrade is not supported to orchestrator release: %s", ucs.OrchestratorProfile.OrchestratorRelease)
		}
	}

	return nil
}

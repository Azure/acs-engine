package vlabs

import "fmt"

// UpgradeContainerService API model
type UpgradeContainerService struct {
	OrchestratorProfile *OrchestratorProfile `json:"orchestratorProfile,omitempty"`
}

// Validate implements APIObject
func (ucs *UpgradeContainerService) Validate() error {
	switch ucs.OrchestratorProfile.OrchestratorType {
	case DCOS:
	case Swarm:
	case SwarmMode:
		return fmt.Errorf("Upgrade is not supported for orchestrator: %s", ucs.OrchestratorProfile.OrchestratorType)
	case Kubernetes:
		switch ucs.OrchestratorProfile.OrchestratorVersion {
		case Kubernetes162:
		case Kubernetes160:
		default:
			return fmt.Errorf("Invalid orchestrator version: %s", ucs.OrchestratorProfile.OrchestratorVersion)
		}
	}

	return nil
}

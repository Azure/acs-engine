package vlabs

import (
	"fmt"

	"github.com/Azure/acs-engine/pkg/api/common"
	validator "gopkg.in/go-playground/validator.v9"
)

// UpgradeContainerService API model
type UpgradeContainerService struct {
	OrchestratorProfile *OrchestratorProfile `json:"orchestratorProfile,omitempty" validate:"required"`
}

func handleUpgradeValidationErrors(e validator.ValidationErrors) error {
	err := e[0]
	ns := err.Namespace()
	switch ns {
	case "UpgradeContainerService.OrchestratorProfile", "UpgradeContainerService.OrchestratorProfile.OrchestratorType":
		return fmt.Errorf("missing %s", ns)
	case "UpgradeContainerService.OrchestratorProfile.OrchestratorVersion":
		return fmt.Errorf("OrchestratorVersion is a readyonly field, leave it empty")
	}
	return nil
}

// Validate implements APIObject
func (ucs *UpgradeContainerService) Validate() error {
	if e := validate.Struct(ucs); e != nil {
		return handleUpgradeValidationErrors(e.(validator.ValidationErrors))
	}
	switch ucs.OrchestratorProfile.OrchestratorType {
	case DCOS:
	case Swarm:
	case SwarmMode:
		return fmt.Errorf("Upgrade is not supported for orchestrator: %s", ucs.OrchestratorProfile.OrchestratorType)
	case Kubernetes:
		switch ucs.OrchestratorProfile.OrchestratorVersionHint {
		case common.KubernetesVersionHint16:
		default:
			return fmt.Errorf("Invalid orchestrator version: %s", ucs.OrchestratorProfile.OrchestratorVersion)
		}
	}

	return nil
}

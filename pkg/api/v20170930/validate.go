package v20170930

import (
	"fmt"
	"strings"

	"github.com/Azure/acs-engine/pkg/api/common"
)

// Validate implements APIObject
func (o *OrchestratorVersionProfile) Validate() error {
	switch {
	case strings.EqualFold(o.OrchestratorType, Kubernetes):
		o.OrchestratorType = Kubernetes
	case strings.EqualFold(o.OrchestratorType, DCOS):
		o.OrchestratorType = DCOS
	case strings.EqualFold(o.OrchestratorType, Swarm):
		o.OrchestratorType = Swarm
	case strings.EqualFold(o.OrchestratorType, DockerCE):
		o.OrchestratorType = DockerCE
	default:
		return fmt.Errorf("Unsupported orchestrator '%s'", o.OrchestratorType)
	}
	return nil
}

// ValidateForUpgrade validates upgrade input data
func (o *OrchestratorProfile) ValidateForUpgrade() error {
	switch o.OrchestratorType {
	case DCOS, DockerCE, Swarm:
		return fmt.Errorf("Upgrade is not supported for orchestrator %s", o.OrchestratorType)
	case Kubernetes:
		switch o.OrchestratorVersion {
		case common.KubernetesVersion1Dot6Dot13:
		case common.KubernetesVersion1Dot7Dot10:
		default:
			return fmt.Errorf("Upgrade to Kubernetes %s is not supported", o.OrchestratorVersion)
		}
	}
	return nil
}

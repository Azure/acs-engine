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
		release, err := common.GetReleaseFromVersion(o.OrchestratorVersion)
		if err != nil {
			return fmt.Errorf("OrchestratorVersion %s is not a valid version", o.OrchestratorVersion)
		}
		switch release {
		case common.KubernetesRelease1Dot6, common.KubernetesRelease1Dot7:
			if o.OrchestratorVersion != common.KubeReleaseToVersion[release] {
				return fmt.Errorf("Upgrade to Kubernetes version %s is not supported, we support upgrade to version %s",
					o.OrchestratorVersion,
					common.KubeReleaseToVersion[release])
			}
		default:
			return fmt.Errorf("Upgrade to Kubernetes %s is not supported", release)
		}
	}
	return nil
}

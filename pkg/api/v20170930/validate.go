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
		if _, ok := common.KubeReleaseToVersion[o.OrchestratorRelease]; !ok {
			return fmt.Errorf("Unsupported Kubernetes release '%s'", o.OrchestratorRelease)
		}
	case strings.EqualFold(o.OrchestratorType, DCOS):
		o.OrchestratorType = DCOS
		if _, ok := common.DCOSReleaseToVersion[o.OrchestratorRelease]; !ok {
			return fmt.Errorf("Unsupported DCOS release '%s'", o.OrchestratorRelease)
		}
	case strings.EqualFold(o.OrchestratorType, Swarm):
		o.OrchestratorType = Swarm
	case strings.EqualFold(o.OrchestratorType, DockerCE):
		o.OrchestratorType = DockerCE
	default:
		return fmt.Errorf("Unsupported orchestrator '%s'", o.OrchestratorType)
	}
	return nil
}

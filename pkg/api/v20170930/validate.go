package v20170930

import (
	"fmt"
	"strings"

	"github.com/Azure/acs-engine/pkg/api/common"
)

// Validate implements APIObject
func (o *OrchestratorInfo) Validate() error {
	switch {
	case strings.EqualFold(o.Orchestrator, Kubernetes):
		o.Orchestrator = Kubernetes
		if _, ok := common.KubeReleaseToVersion[o.Release]; !ok {
			return fmt.Errorf("Unsupported Kubernetes release '%s'", o.Release)
		}
	case strings.EqualFold(o.Orchestrator, DCOS):
		o.Orchestrator = DCOS
		if _, ok := common.DCOSReleaseToVersion[o.Release]; !ok {
			return fmt.Errorf("Unsupported Kubernetes release '%s'", o.Release)
		}
	case strings.EqualFold(o.Orchestrator, Swarm):
		o.Orchestrator = Swarm
	case strings.EqualFold(o.Orchestrator, DockerCE):
		o.Orchestrator = DockerCE
	default:
		return fmt.Errorf("Unsupported orchestrator '%s'", o.Orchestrator)
	}
	return nil
}

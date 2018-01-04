package v20170930

import (
	"fmt"
	"strings"
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

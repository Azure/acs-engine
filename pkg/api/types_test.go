package api

import (
	"testing"
)

func TestIsDCOS(t *testing.T) {
	dCOSProfile := &OrchestratorProfile{
		OrchestratorType:    "DCOS",
		OrchestratorVersion: "vlabs",
	}
	if !dCOSProfile.IsDCOS() {
		t.Fatalf("unable to detect DCOS orchestrator profile from OrchestratorType=%s", dCOSProfile.OrchestratorType)
	}
	kubernetesProfile := &OrchestratorProfile{
		OrchestratorType:    "Kubernetes",
		OrchestratorVersion: "vlabs",
	}
	if kubernetesProfile.IsDCOS() {
		t.Fatalf("unexpectedly detected DCOS orchestrator profile from OrchestratorType=%s", kubernetesProfile.OrchestratorType)
	}
}

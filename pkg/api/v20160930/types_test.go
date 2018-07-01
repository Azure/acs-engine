package v20160930

import (
	"encoding/json"
	"testing"
)

func TestIsDCOS(t *testing.T) {
	dCOSProfile := &OrchestratorProfile{
		OrchestratorType: "DCOS",
	}
	if !dCOSProfile.IsDCOS() {
		t.Fatalf("unable to detect DCOS orchestrator profile from OrchestratorType=%s", dCOSProfile.OrchestratorType)
	}
	kubernetesProfile := &OrchestratorProfile{
		OrchestratorType: "Kubernetes",
	}
	if kubernetesProfile.IsDCOS() {
		t.Fatalf("unexpectedly detected DCOS orchestrator profile from OrchestratorType=%s", kubernetesProfile.OrchestratorType)
	}
}

func TestMasterProfile(t *testing.T) {
	MasterProfileText := "{\"count\" : 0}"
	mp := &MasterProfile{}
	if e := json.Unmarshal([]byte(MasterProfileText), mp); e != nil {
		t.Fatalf("unexpectedly detected unmarshal failure for MasterProfile, %+v", e)
	}

	if mp.Count != 1 {
		t.Fatalf("unexpectedly detected MasterProfile.Count != 1 after unmarshal")
	}
}

func TestAgentPoolProfile(t *testing.T) {
	// With osType not specified
	AgentPoolProfileText := "{\"count\" : 0}"
	ap := &AgentPoolProfile{}
	if e := json.Unmarshal([]byte(AgentPoolProfileText), ap); e != nil {
		t.Fatalf("unexpectedly detected unmarshal failure for AgentPoolProfile, %+v", e)
	}

	if ap.Count != 1 {
		t.Fatalf("unexpectedly detected AgentPoolProfile.Count != 1 after unmarshal")
	}

	if !ap.IsLinux() {
		t.Fatalf("unexpectedly detected AgentPoolProfile.OSType != Linux after unmarshal")
	}

	// With osType specified
	AgentPoolProfileText = `{ "name": "linuxpool1", "osType" : "Windows", "count": 1, "vmSize": "Standard_D2_v2", 
"availabilityProfile": "AvailabilitySet", "storageProfile" : "ManagedDisks", "vnetSubnetID" : "12345" }`
	ap = &AgentPoolProfile{}
	if e := json.Unmarshal([]byte(AgentPoolProfileText), ap); e != nil {
		t.Fatalf("unexpectedly detected unmarshal failure for AgentPoolProfile, %+v", e)
	}

	if ap.Count != 1 {
		t.Fatalf("unexpectedly detected AgentPoolProfile.Count != 1 after unmarshal")
	}

	if !ap.IsWindows() {
		t.Fatalf("unexpectedly detected AgentPoolProfile.OSType != Windows after unmarshal")
	}
}

func TestOrchestratorProfile(t *testing.T) {
	OrchestratorProfileText := `{ "orchestratorType": "Mesos" }`
	op := &OrchestratorProfile{}
	if e := json.Unmarshal([]byte(OrchestratorProfileText), op); e != nil {
		t.Fatalf("unexpectedly detected unmarshal failure for OrchestratorProfile, %+v", e)
	}

	OrchestratorProfileText = `{ "orchestratorType": "Swarm" }`
	op = &OrchestratorProfile{}
	if e := json.Unmarshal([]byte(OrchestratorProfileText), op); e != nil {
		t.Fatalf("unexpectedly detected unmarshal failure for OrchestratorProfile, %+v", e)
	}

	OrchestratorProfileText = `{ "orchestratorType": "DCOS" }`
	op = &OrchestratorProfile{}
	if e := json.Unmarshal([]byte(OrchestratorProfileText), op); e != nil {
		t.Fatalf("unexpectedly detected unmarshal failure for OrchestratorProfile, %+v", e)
	}

	OrchestratorProfileText = `{ "orchestratorType": "Kubernetes" }`
	op = &OrchestratorProfile{}
	if e := json.Unmarshal([]byte(OrchestratorProfileText), op); e != nil {
		t.Fatalf("unexpectedly detected unmarshal failure for OrchestratorProfile, %+v", e)

	}
}

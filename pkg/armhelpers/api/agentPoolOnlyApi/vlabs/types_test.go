package vlabs

import (
	"encoding/json"
	"testing"
)

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

	if ap.OSType != Linux {
		t.Fatalf("unexpectedly detected AgentPoolProfile.OSType != Linux after unmarshal")
	}

	if !ap.IsStorageAccount() {
		t.Fatalf("unexpectedly detected AgentPoolProfile.StorageProfile != StorageAccount after unmarshal")
	}
}

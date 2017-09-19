package v20170701

import (
	"encoding/json"
	"testing"
)

func TestMasterProfile(t *testing.T) {
	MasterProfileText := "{\"count\" : 0}"
	mp := &MasterProfile{}
	if e := json.Unmarshal([]byte(MasterProfileText), mp); e != nil {
		t.Fatalf("unexpectedly detected unmarshal failure for MasterProfile, %+v", e)
	}

	if mp.Count != 1 {
		t.Fatalf("unexpectedly detected MasterProfile.Count != 1 after unmarshal")
	}

	if mp.FirstConsecutiveStaticIP != "10.240.255.5" {
		t.Fatalf("unexpectedly detected MasterProfile.FirstConsecutiveStaticIP != 10.240.255.5 after unmarshal")
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

	if ap.OSType != Linux {
		t.Fatalf("unexpectedly detected AgentPoolProfile.OSType != Linux after unmarshal")
	}
}

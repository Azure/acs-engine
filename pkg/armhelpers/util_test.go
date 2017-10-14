package armhelpers

import (
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/arm/compute"
)

func Test_SplitBlobURI(t *testing.T) {
	expectedAccountName := "vhdstorage8h8pjybi9hbsl6"
	expectedContainerName := "vhds"
	expectedBlobPath := "osdisks/disk1234.vhd"
	accountName, containerName, blobPath, err := SplitBlobURI("https://vhdstorage8h8pjybi9hbsl6.blob.core.windows.net/vhds/osdisks/disk1234.vhd")
	if accountName != expectedAccountName {
		t.Fatalf("incorrect account name. expected=%s actual=%s", expectedAccountName, accountName)
	}
	if containerName != expectedContainerName {
		t.Fatalf("incorrect account name. expected=%s actual=%s", expectedContainerName, containerName)
	}
	if blobPath != expectedBlobPath {
		t.Fatalf("incorrect account name. expected=%s actual=%s", expectedBlobPath, blobPath)
	}
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
}

func Test_LinuxVMNameParts(t *testing.T) {
	data := []struct {
		orchestrator, poolIdentifier, nameSuffix string
		agentIndex                               int
	}{
		{"k8s", "agentpool1", "38988164", 10},
		{"k8s", "agent-pool1", "38988164", 8},
		{"k8s", "agent-pool-1", "38988164", 0},
	}

	for _, el := range data {
		vmName := fmt.Sprintf("%s-%s-%s-%d", el.orchestrator, el.poolIdentifier, el.nameSuffix, el.agentIndex)
		orchestrator, poolIdentifier, nameSuffix, agentIndex, err := LinuxVMNameParts(vmName)
		if orchestrator != el.orchestrator {
			t.Fatalf("incorrect orchestrator. expected=%s actual=%s", el.orchestrator, orchestrator)
		}
		if poolIdentifier != el.poolIdentifier {
			t.Fatalf("incorrect poolIdentifier. expected=%s actual=%s", el.poolIdentifier, poolIdentifier)
		}
		if nameSuffix != el.nameSuffix {
			t.Fatalf("incorrect nameSuffix. expected=%s actual=%s", el.nameSuffix, nameSuffix)
		}
		if agentIndex != el.agentIndex {
			t.Fatalf("incorrect agentIndex. expected=%d actual=%d", el.agentIndex, agentIndex)
		}
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
	}
}

func Test_WindowsVMNameParts(t *testing.T) {
	expectedPoolPrefix := "38988"
	expectedAcs := "acs"
	expectedPoolIndex := 903
	expectedAgentIndex := 12

	poolPrefix, acs, poolIndex, agentIndex, err := WindowsVMNameParts("38988acs90312")
	if poolPrefix != expectedPoolPrefix {
		t.Fatalf("incorrect poolPrefix. expected=%s actual=%s", expectedPoolPrefix, poolPrefix)
	}
	if acs != expectedAcs {
		t.Fatalf("incorrect acs string. expected=%s actual=%s", expectedAcs, acs)
	}
	if poolIndex != expectedPoolIndex {
		t.Fatalf("incorrect poolIndex. expected=%d actual=%d", expectedPoolIndex, poolIndex)
	}
	if agentIndex != expectedAgentIndex {
		t.Fatalf("incorrect agentIndex. expected=%d actual=%d", expectedAgentIndex, agentIndex)
	}
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
}

func Test_GetVMNameIndexLinux(t *testing.T) {
	expectedAgentIndex := 65

	agentIndex, err := GetVMNameIndex(compute.Linux, "k8s-agentpool1-38988164-65")

	if agentIndex != expectedAgentIndex {
		t.Fatalf("incorrect agentIndex. expected=%d actual=%d", expectedAgentIndex, agentIndex)
	}
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
}

func Test_GetVMNameIndexWindows(t *testing.T) {
	expectedAgentIndex := 20

	agentIndex, err := GetVMNameIndex(compute.Windows, "38988acs90320")

	if agentIndex != expectedAgentIndex {
		t.Fatalf("incorrect agentIndex. expected=%d actual=%d", expectedAgentIndex, agentIndex)
	}
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
}

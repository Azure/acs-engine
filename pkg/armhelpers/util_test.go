package armhelpers

import (
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
	expectedOrchestrator := "k8s"
	expectedPoolIdentifier := "agentpool1"
	expectedNameSuffix := "38988164"
	expectedAgentIndex := 10

	orchestrator, poolIdentifier, nameSuffix, agentIndex, err := LinuxVMNameParts("k8s-agentpool1-38988164-10")
	if orchestrator != expectedOrchestrator {
		t.Fatalf("incorrect orchestrator. expected=%s actual=%s", expectedOrchestrator, orchestrator)
	}
	if poolIdentifier != expectedPoolIdentifier {
		t.Fatalf("incorrect poolIdentifier. expected=%s actual=%s", expectedPoolIdentifier, poolIdentifier)
	}
	if nameSuffix != expectedNameSuffix {
		t.Fatalf("incorrect nameSuffix. expected=%s actual=%s", expectedNameSuffix, nameSuffix)
	}
	if agentIndex != expectedAgentIndex {
		t.Fatalf("incorrect agentIndex. expected=%d actual=%d", expectedAgentIndex, agentIndex)
	}
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
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

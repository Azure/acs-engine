package utils

import (
	"fmt"
	"testing"

	"github.com/Azure/acs-engine/pkg/api"
	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2018-04-01/compute"
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
		poolIdentifier, nameSuffix string
		agentIndex                 int
	}{
		{"agentpool1", "38988164", 10},
		{"agent-pool1", "38988164", 8},
		{"agent-pool-1", "38988164", 0},
	}

	for _, el := range data {
		vmName := fmt.Sprintf("k8s-%s-%s-%d", el.poolIdentifier, el.nameSuffix, el.agentIndex)
		poolIdentifier, nameSuffix, agentIndex, err := K8sLinuxVMNameParts(vmName)
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

func Test_VmssNameParts(t *testing.T) {
	data := []struct {
		poolIdentifier, nameSuffix string
	}{
		{"agentpool1", "38988164"},
		{"agent-pool1", "38988164"},
		{"agent-pool-1", "38988164"},
	}

	for _, el := range data {
		vmssName := fmt.Sprintf("swarmm-%s-%s-vmss", el.poolIdentifier, el.nameSuffix)
		poolIdentifier, nameSuffix, err := VmssNameParts(vmssName)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		if poolIdentifier != el.poolIdentifier {
			t.Fatalf("incorrect poolIdentifier. expected=%s actual=%s", el.poolIdentifier, poolIdentifier)
		}
		if nameSuffix != el.nameSuffix {
			t.Fatalf("incorrect nameSuffix. expected=%s actual=%s", el.nameSuffix, nameSuffix)
		}
	}
}

func Test_WindowsVMNameParts(t *testing.T) {
	expectedPoolPrefix := "38988"
	expectedAcs := "k8s"
	expectedPoolIndex := 3
	expectedAgentIndex := 12

	poolPrefix, acs, poolIndex, agentIndex, err := WindowsVMNameParts("38988k8s90312")
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

func Test_WindowsVMSSNameParts(t *testing.T) {
	expectedPoolPrefix := "38988"
	expectedAcs := "k8s"
	expectedPoolIndex := 3

	poolPrefix, acs, poolIndex, err := WindowsVMSSNameParts("38988k8s903")
	if poolPrefix != expectedPoolPrefix {
		t.Fatalf("incorrect poolPrefix. expected=%s actual=%s", expectedPoolPrefix, poolPrefix)
	}
	if acs != expectedAcs {
		t.Fatalf("incorrect acs string. expected=%s actual=%s", expectedAcs, acs)
	}
	if poolIndex != expectedPoolIndex {
		t.Fatalf("incorrect poolIndex. expected=%d actual=%d", expectedPoolIndex, poolIndex)
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

	agentIndex, err := GetVMNameIndex(compute.Windows, "38988k8s90320")

	if agentIndex != expectedAgentIndex {
		t.Fatalf("incorrect agentIndex. expected=%d actual=%d", expectedAgentIndex, agentIndex)
	}
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
}

func Test_GetK8sVMName(t *testing.T) {

	for _, s := range []struct {
		osType                     api.OSType
		isAKS                      bool
		nameSuffix, agentPoolName  string
		agentPoolIndex, agentIndex int
		expected                   string
		expectedErr                bool
	}{
		{api.Linux, true, "35953384", "agentpool1", 0, 2, "aks-agentpool1-35953384-2", false},
		{api.Windows, false, "35953384", "agentpool1", 0, 2, "35953k8s9002", false},
		{"macOS", false, "35953384", "agentpool1", 0, 2, "", true},
	} {
		vmName, err := GetK8sVMName(s.osType, s.isAKS, s.nameSuffix, s.agentPoolName, s.agentPoolIndex, s.agentIndex)

		if !s.expectedErr {
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
		}
		if vmName != s.expected {
			t.Fatalf("vmName %s, expected %s", vmName, s.expected)
		}
	}
}

func Test_ResourceName(t *testing.T) {
	s := "https://vhdstorage8h8pjybi9hbsl6.blob.core.windows.net/vhds/osdisks/disk1234.vhd"
	expected := "disk1234.vhd"
	r, err := ResourceName(s)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if r != expected {
		t.Fatalf("resourceName %s, expected %s", r, expected)
	}
}

func Test_ResourceNameInvalid(t *testing.T) {
	s := "https://vhdstorage8h8pjybi9hbsl6.blob.core.windows.net/vhds/osdisks/"
	expectedMsg := "resource name was missing from identifier"
	_, err := ResourceName(s)
	if err == nil || err.Error() != expectedMsg {
		t.Fatalf("expected error with message: %s", expectedMsg)
	}
}

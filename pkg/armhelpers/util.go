package armhelpers

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/Azure/azure-sdk-for-go/arm/compute"
)

// ResourceName returns the last segment (the resource name) for the specified resource identifier.
func ResourceName(ID string) (string, error) {
	parts := strings.Split(ID, "/")
	name := parts[len(parts)-1]
	if len(name) == 0 {
		return "", fmt.Errorf("resource name was missing from identifier")
	}

	return name, nil
}

// SplitBlobURI returns a decomposed blob URI parts: accountName, containerName, blobName.
func SplitBlobURI(URI string) (string, string, string, error) {
	uri, err := url.Parse(URI)
	if err != nil {
		return "", "", "", err
	}

	accountName := strings.Split(uri.Host, ".")[0]
	urlParts := strings.Split(uri.Path, "/")

	containerName := urlParts[1]
	blobPath := strings.Join(urlParts[2:], "/")

	return accountName, containerName, blobPath, nil
}

// LinuxVMNameParts returns parts of Linux VM name e.g: k8s-agentpool1-11290731-0
func LinuxVMNameParts(vmName string) (orchestrator string, poolType string, nameSuffix string, agentIndex int, err error) {
	vmNameParts := strings.Split(vmName, "-")

	if len(vmNameParts) != 4 {
		return "", "", "", -1, fmt.Errorf("resource name was missing from identifier")
	}

	vmNum, err := strconv.Atoi(vmNameParts[3])

	if err != nil {
		return "", "", "", -1, fmt.Errorf("Error parsing VM Name: %v", err)
	}

	return vmNameParts[0], vmNameParts[1], vmNameParts[2], vmNum, nil
}

// WindowsVMNameParts returns parts of Windows VM name e.g: 50621acs9000
func WindowsVMNameParts(vmName string) (poolPrefix string, acsStr string, poolIdentifier int, agentIndex int, err error) {
	poolPrefix = strings.Split(vmName, "acs")[0]
	poolInfo := strings.Split(vmName, "acs")[1]

	poolIdentifier, err = strconv.Atoi(poolInfo[:3])
	if err != nil {
		return "", "", -1, -1, fmt.Errorf("Error parsing VM Name: %v", err)
	}

	agentIndex, _ = strconv.Atoi(poolInfo[3:])
	fmt.Printf("%d\n", agentIndex)

	if err != nil {
		return "", "", -1, -1, fmt.Errorf("Error parsing VM Name: %v", err)
	}

	return poolPrefix, "acs", poolIdentifier, agentIndex, nil
}

// ByVMNameOffset implements sort.Interface for []VirtualMachine based on
// the Name field.
type ByVMNameOffset []compute.VirtualMachine

func (vm ByVMNameOffset) Len() int      { return len(vm) }
func (vm ByVMNameOffset) Swap(i, j int) { vm[i], vm[j] = vm[j], vm[i] }
func (vm ByVMNameOffset) Less(i, j int) bool {
	vm1NameParts := strings.Split(*vm[i].Name, "-")
	vm1Num := vm1NameParts[len(vm1NameParts)-1]

	vm2NameParts := strings.Split(*vm[j].Name, "-")
	vm2Num := vm2NameParts[len(vm2NameParts)-1]

	return vm1Num < vm2Num
}

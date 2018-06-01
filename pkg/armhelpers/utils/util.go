package utils

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/Azure/acs-engine/pkg/api"
	"github.com/Azure/azure-sdk-for-go/arm/compute"
	log "github.com/sirupsen/logrus"
)

const (
	// TODO: merge with the RP code
	k8sLinuxVMNamingFormat         = "^[0-9a-zA-Z]{3}-(.+)-([0-9a-fA-F]{8})-{0,2}([0-9]+)$"
	k8sLinuxVMAgentPoolNameIndex   = 1
	k8sLinuxVMAgentClusterIDIndex  = 2
	k8sLinuxVMAgentIndexArrayIndex = 3

	k8sWindowsVMNamingFormat               = "^([a-fA-F0-9]{5})([0-9a-zA-Z]{3})([a-zA-Z0-9]{4,6})$"
	k8sWindowsVMAgentPoolPrefixIndex       = 1
	k8sWindowsVMAgentOrchestratorNameIndex = 2
	k8sWindowsVMAgentPoolInfoIndex         = 3
	// here there are 2 capture groups
	//  the first is the agent pool name, which can contain -s
	//  the second group is the Cluster ID for the cluster
	vmssNamingFormat       = "^[0-9a-zA-Z]+-(.+)-([0-9a-fA-F]{8})-vmss$"
	vmssAgentPoolNameIndex = 1
	vmssClusterIDIndex     = 2

	windowsVmssNamingFormat                   = "^([a-fA-F0-9]{5})([0-9a-zA-Z]{3})([a-zA-Z0-9]{3})$"
	windowsVmssAgentPoolNameIndex             = 1
	windowsVmssAgentPoolOrchestratorNameIndex = 2
	windowsVmssAgentPoolIndex                 = 3
)

var vmnameLinuxRegexp *regexp.Regexp
var vmssnameRegexp *regexp.Regexp
var vmnameWindowsRegexp *regexp.Regexp
var vmssnameWindowsRegexp *regexp.Regexp

func init() {
	vmnameLinuxRegexp = regexp.MustCompile(k8sLinuxVMNamingFormat)
	vmnameWindowsRegexp = regexp.MustCompile(k8sWindowsVMNamingFormat)
	vmssnameRegexp = regexp.MustCompile(vmssNamingFormat)
	vmssnameWindowsRegexp = regexp.MustCompile(windowsVmssNamingFormat)
}

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

// K8sLinuxVMNameParts returns parts of Linux VM name e.g: k8s-agentpool1-11290731-0
func K8sLinuxVMNameParts(vmName string) (poolIdentifier, nameSuffix string, agentIndex int, err error) {
	vmNameParts := vmnameLinuxRegexp.FindStringSubmatch(vmName)
	if len(vmNameParts) != 4 {
		return "", "", -1, fmt.Errorf("resource name was missing from identifier")
	}

	vmNum, err := strconv.Atoi(vmNameParts[k8sLinuxVMAgentIndexArrayIndex])

	if err != nil {
		return "", "", -1, fmt.Errorf("Error parsing VM Name: %v", err)
	}

	return vmNameParts[k8sLinuxVMAgentPoolNameIndex], vmNameParts[k8sLinuxVMAgentClusterIDIndex], vmNum, nil
}

// VmssNameParts returns parts of Linux VM name e.g: k8s-agentpool1-11290731-0
func VmssNameParts(vmssName string) (poolIdentifier, nameSuffix string, err error) {
	vmssNameParts := vmssnameRegexp.FindStringSubmatch(vmssName)
	if len(vmssNameParts) != 3 {
		return "", "", fmt.Errorf("resource name was missing from identifier")
	}

	return vmssNameParts[vmssAgentPoolNameIndex], vmssNameParts[vmssClusterIDIndex], nil
}

// WindowsVMNameParts returns parts of Windows VM name e.g: 50621k8s9000
func WindowsVMNameParts(vmName string) (poolPrefix string, acsStr string, poolIndex int, agentIndex int, err error) {
	vmNameParts := vmnameWindowsRegexp.FindStringSubmatch(vmName)
	if len(vmNameParts) != 4 {
		return "", "", -1, -1, fmt.Errorf("resource name was missing from identifier")
	}

	poolPrefix = vmNameParts[k8sWindowsVMAgentPoolPrefixIndex]
	acsStr = vmNameParts[k8sWindowsVMAgentOrchestratorNameIndex]
	poolInfo := vmNameParts[k8sWindowsVMAgentPoolInfoIndex]

	poolIndex, err = strconv.Atoi(poolInfo[:3])
	if err != nil {
		return "", "", -1, -1, fmt.Errorf("Error parsing VM Name: %v", err)
	}
	poolIndex -= 900
	agentIndex, _ = strconv.Atoi(poolInfo[3:])
	fmt.Printf("%d\n", agentIndex)

	if err != nil {
		return "", "", -1, -1, fmt.Errorf("Error parsing VM Name: %v", err)
	}

	return poolPrefix, acsStr, poolIndex, agentIndex, nil
}

// WindowsVMSSNameParts returns parts of Windows VM name e.g: 50621k8s900
func WindowsVMSSNameParts(vmssName string) (poolPrefix string, acsStr string, poolIndex int, err error) {
	vmssNameParts := vmssnameWindowsRegexp.FindStringSubmatch(vmssName)
	if len(vmssNameParts) != 4 {
		return "", "", -1, fmt.Errorf("resource name was missing from identifier")
	}

	poolPrefix = vmssNameParts[windowsVmssAgentPoolNameIndex]
	acsStr = vmssNameParts[windowsVmssAgentPoolOrchestratorNameIndex]
	poolInfo := vmssNameParts[windowsVmssAgentPoolIndex]

	poolIndex, err = strconv.Atoi(poolInfo)
	if err != nil {
		return "", "", -1, fmt.Errorf("Error parsing VM Name: %v", err)
	}
	poolIndex -= 900

	return poolPrefix, acsStr, poolIndex, nil
}

// GetVMNameIndex return VM index of a node in the Kubernetes cluster
func GetVMNameIndex(osType compute.OperatingSystemTypes, vmName string) (int, error) {
	var agentIndex int
	var err error
	if osType == compute.Linux {
		_, _, agentIndex, err = K8sLinuxVMNameParts(vmName)
		if err != nil {
			log.Errorln(err)
			return 0, err
		}
	} else if osType == compute.Windows {
		_, _, _, agentIndex, err = WindowsVMNameParts(vmName)
		if err != nil {
			log.Errorln(err)
			return 0, err
		}
	}

	return agentIndex, nil
}

// GetK8sVMName reconstructs VM name
func GetK8sVMName(osType api.OSType, isAKS bool, nameSuffix, agentPoolName string, agentPoolIndex, agentIndex int) (string, error) {
	prefix := "k8s"
	if isAKS {
		prefix = "aks"
	}
	if osType == api.Linux {
		return fmt.Sprintf("%s-%s-%s-%d", prefix, agentPoolName, nameSuffix, agentIndex), nil
	}
	if osType == api.Windows {
		return fmt.Sprintf("%s%s%d%d", nameSuffix[:5], prefix, 900+agentPoolIndex, agentIndex), nil
	}
	return "", fmt.Errorf("Failed to reconstruct VM Name")
}

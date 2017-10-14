package armhelpers

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/Azure/azure-sdk-for-go/arm/compute"
	log "github.com/sirupsen/logrus"
)

const (
	// TODO: merge with the RP code
	k8sLinuxVMNamingFormat         = "^k8s-(.+)-([0-9a-fA-F]{8})-{0,2}([0-9]+)$"
	k8sLinuxVMAgentPoolNameIndex   = 1
	k8sLinuxVMAgentClusterIDIndex  = 2
	k8sLinuxVMAgentIndexArrayIndex = 3
)

var vmnameLinuxRegexp *regexp.Regexp

func init() {
	vmnameLinuxRegexp = regexp.MustCompile(k8sLinuxVMNamingFormat)
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

// WindowsVMNameParts returns parts of Windows VM name e.g: 50621acs9000
func WindowsVMNameParts(vmName string) (poolPrefix string, acsStr string, poolIndex int, agentIndex int, err error) {
	poolPrefix = strings.Split(vmName, "acs")[0]
	poolInfo := strings.Split(vmName, "acs")[1]

	poolIndex, err = strconv.Atoi(poolInfo[:3])
	if err != nil {
		return "", "", -1, -1, fmt.Errorf("Error parsing VM Name: %v", err)
	}

	agentIndex, _ = strconv.Atoi(poolInfo[3:])
	fmt.Printf("%d\n", agentIndex)

	if err != nil {
		return "", "", -1, -1, fmt.Errorf("Error parsing VM Name: %v", err)
	}

	return poolPrefix, "acs", poolIndex, agentIndex, nil
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

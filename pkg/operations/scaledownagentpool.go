package operations

import (
	"container/list"
	"fmt"

	"github.com/Azure/acs-engine/pkg/api"
	"github.com/Azure/acs-engine/pkg/armhelpers"
)

var vmNameFormats = make(map[api.OrchestratorType]map[api.OSType]string)

const (
	//These formats assume the arguments to Sprintf are agent pool index, name suffix, vm index in that order
	windowsAgentVMFormat           = "%.5[2]sacs9%02[1]d%[3]d"
	k8sLinuxAgentVMFormat          = "k8s-agent%d-%s-%d"
	dcosLinuxAgentVMFormat         = "dcos-agent%d-%s-%d"
	mesosLinuxAgentVMFormat        = "mesos-agent%d-%s-%d"
	swarmLinuxAgentVMFormat        = "swarm-agent%d-%s-%d"
	swarmModeLinuxAgentVMFormat    = "swarmm-agent%d-%s-%d"
	k8sClassicLinuxAgentVMFormat   = "k8s-agent-%[2]s-%[3]d"
	dcosClassicLinuxAgentVMFormat  = "dcos-agent-%[2]s-%[3]d"
	swarmClassicLinuxAgentVMFormat = "swarm-agent-%[2]s-%[3]d"
)

func init() {
	vmNameFormats[api.Kubernetes] = make(map[api.OSType]string)
	vmNameFormats[api.Kubernetes][api.Linux] = k8sLinuxAgentVMFormat
	vmNameFormats[api.Kubernetes][api.Windows] = windowsAgentVMFormat
	vmNameFormats[api.DCOS] = make(map[api.OSType]string)
	vmNameFormats[api.DCOS][api.Linux] = dcosLinuxAgentVMFormat
	// Mesos is deprecated in favor of DCOS never supported VMAS agents so we shouldn't need a name format
	// This line is just to prevent a nil dereference
	vmNameFormats[api.Mesos] = make(map[api.OSType]string)
	vmNameFormats[api.Swarm] = make(map[api.OSType]string)
	vmNameFormats[api.Swarm][api.Linux] = swarmLinuxAgentVMFormat
	vmNameFormats[api.Swarm][api.Windows] = windowsAgentVMFormat
	vmNameFormats[api.SwarmMode] = make(map[api.OSType]string)
	vmNameFormats[api.SwarmMode][api.Linux] = swarmModeLinuxAgentVMFormat
	vmNameFormats[api.SwarmMode][api.Windows] = windowsAgentVMFormat
}

//VMScalingErrorDetails give the index in the agent pool that failed and the accompanying error
type VMScalingErrorDetails struct {
	Index int
	Error error
}

// ScaleDownVMASAgentPool removes the vms at the specified indexes from the specified agent pool. Returns a list with details on each failure.
// all items in the list will always be of type *VMScalingErrorDetails
func ScaleDownVMASAgentPool(az armhelpers.ACSEngineClient, clusterSuffix, resourceGroup string, agentPoolIndex int,
	orchestratorType api.OrchestratorType, osType api.OSType, indexesToRemove ...int) *list.List {
	numVmsToDelete := len(indexesToRemove)
	errChan := make(chan *VMScalingErrorDetails, numVmsToDelete)
	defer close(errChan)
	for _, vmIndex := range indexesToRemove {
		go func(vmIndex int) {
			vmName, err := getVMName(orchestratorType, osType, agentPoolIndex, vmIndex, clusterSuffix)
			if err != nil {
				errChan <- &VMScalingErrorDetails{Index: vmIndex, Error: err}
				return
			}
			err = CleanDeleteVirtualMachine(az, resourceGroup, vmName)
			if err != nil {
				errChan <- &VMScalingErrorDetails{Index: vmIndex, Error: err}
				return
			}
			errChan <- nil
		}(vmIndex)
	}
	failedVMDeletions := &list.List{}
	for i := 0; i < numVmsToDelete; i++ {
		errDetails := <-errChan
		if errDetails != nil {
			failedVMDeletions.PushBack(errDetails)
		}
	}
	if failedVMDeletions.Len() > 0 {
		return failedVMDeletions
	}
	return nil
}

func getVMName(orchestratorType api.OrchestratorType, osType api.OSType, agentpoolIndex, vmIndex int, clusterSuffix string) (string, error) {
	format, ok := vmNameFormats[orchestratorType][osType]
	if !ok {
		return "", fmt.Errorf("Unable to find vm name format for Orchestrator %s and OS %s", orchestratorType, osType)
	}
	return fmt.Sprintf(format, agentpoolIndex, clusterSuffix, vmIndex), nil
}

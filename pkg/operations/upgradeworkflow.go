package operations

import "github.com/Azure/acs-engine/pkg/armhelpers"

// UpgradeWorkFlow outlines various individual high level steps
// that need to be run (one or more times) in the upgrade workflow.
type UpgradeWorkFlow interface {
	ClusterPreflightCheck() error

	// upgrade masters
	// upgrade agent nodes
	RunUpgrade() error

	Validate() error
}

// UpgradeNode drives work flow of deleting and replacing a master or agent node to a
// specified target version of Kubernetes
type UpgradeNode interface {
	// DeleteNode takes state/resources of the master/agent node from ListNodeResources
	// backs up/preserves state as needed by a specific version of Kubernetes and then deletes
	// the node
	DeleteNode(*string) error

	// CreateNode creates a new master/agent node with the targeted version of Kubernetes
	CreateNode(string, int) error

	// Validate will verify the that master/agent node has been upgraded as expected.
	Validate() error
}

// VMStatusSlice for intermmediate upgrade data gathering
type VMStatusSlice []*VMStatus

// Len is part of sort.Interface.
func (vm VMStatusSlice) Len() int {
	return len(vm)
}

// Swap is part of sort.Interface.
func (vm VMStatusSlice) Swap(i, j int) {
	vm[i], vm[j] = vm[j], vm[i]
}

// Less is part of sort.Interface. We use count as the value to sort by
func (vm VMStatusSlice) Less(i, j int) bool {
	_, _, _, vm1Num, _ := armhelpers.LinuxVMNameParts(*vm[i].Name)
	_, _, _, vm2Num, _ := armhelpers.LinuxVMNameParts(*vm[j].Name)

	return vm1Num < vm2Num
}

// VMStatus for intermmediate upgrade data gathering
type VMStatus struct {
	Name    *string
	Delete  bool
	Upgrade bool
}

// VMNumber parse VM name and return the VM number from the name
func (vms VMStatus) VMNumber() int {
	_, _, _, number, err := armhelpers.LinuxVMNameParts(*vms.Name)

	if err != nil {
		return -1
	}

	return number
}

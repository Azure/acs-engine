package common

// validation values
const (
	// MinAgentCount are the minimum number of agents per agent pool
	MinAgentCount = 1
	// MaxAgentCount are the maximum number of agents per agent pool
	MaxAgentCount = 100
	// MinPort specifies the minimum tcp port to open
	MinPort = 1
	// MaxPort specifies the maximum tcp port to open
	MaxPort = 65535
	// MaxDisks specifies the maximum attached disks to add to the cluster
	MaxDisks = 4
	// MinDiskSizeGB specifies the minimum attached disk size
	MinDiskSizeGB = 1
	// MaxDiskSizeGB specifies the maximum attached disk size
	MaxDiskSizeGB = 1023
	// MinIPAddressCount specifies the minimum number of IP addresses per network interface
	MinIPAddressCount = 1
	// MaxIPAddressCount specifies the maximum number of IP addresses per network interface
	MaxIPAddressCount = 256
)

// Availability profiles
const (
	// AvailabilitySet means that the vms are in an availability set
	AvailabilitySet = "AvailabilitySet"
	// VirtualMachineScaleSets means that the vms are in a virtual machine scaleset
	VirtualMachineScaleSets = "VirtualMachineScaleSets"
)

// storage profiles
const (
	// StorageAccount means that the nodes use raw storage accounts for their os and attached volumes
	StorageAccount = "StorageAccount"
	// ManagedDisks means that the nodes use managed disks for their os and attached volumes
	ManagedDisks = "ManagedDisks"
)

const (
	KubernetesVersionHint17 string = "1.7"
	KubernetesVersionHint16 string = "1.6"
	KubernetesVersionHint15 string = "1.5"
	// Why not 1.7?
	KubernetesDefaultVersionHint string = KubernetesVersionHint16
)

// KubeHintToVersion is the hint to actual version map
var KubeHintToVersion = map[string]string{
	KubernetesVersionHint17: "1.7.1",
	KubernetesVersionHint16: "1.6.6",
	KubernetesVersionHint15: "1.5.7",
}

const (
	DCOSVersionHint19      string = "1.9"
	DCOSVersionHint18      string = "1.8"
	DCOSVersionHint17      string = "1.7"
	DCOSDefaultVersionHint string = DCOSVersionHint19
)

// DCOSHintToVersion is the hint to actual version map
var DCOSHintToVersion = map[string]string{
	DCOSVersionHint19: "1.9.0",
	DCOSVersionHint18: "1.8.8",
	DCOSVersionHint17: "1.7.3",
}

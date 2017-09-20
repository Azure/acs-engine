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
	// KubernetesRelease1Dot7 is the major.minor string prefix for 1.7 versions of kubernetes
	KubernetesRelease1Dot7 string = "1.7"
	// KubernetesRelease1Dot6 is the major.minor string prefix for 1.6 versions of kubernetes
	KubernetesRelease1Dot6 string = "1.6"
	// KubernetesRelease1Dot5 is the major.minor string prefix for 1.5 versions of kubernetes
	KubernetesRelease1Dot5 string = "1.5"
	// KubernetesDefaultRelease is the default major.minor version for kubernetes
	KubernetesDefaultRelease string = KubernetesRelease1Dot7
)

// KubeReleaseToVersion maps a major.minor release to an full major.minor.patch version
var KubeReleaseToVersion = map[string]string{
	KubernetesRelease1Dot7: "1.7.5",
	KubernetesRelease1Dot6: "1.6.9",
	KubernetesRelease1Dot5: "1.5.7",
}

const (
	// DCOSRelease1Dot10 is the major.minor string prefix for 1.9 versions of DCOS
	DCOSRelease1Dot10 string = "1.10"
	// DCOSRelease1Dot9 is the major.minor string prefix for 1.9 versions of DCOS
	DCOSRelease1Dot9 string = "1.9"
	// DCOSRelease1Dot8 is the major.minor string prefix for 1.8 versions of DCOS
	DCOSRelease1Dot8 string = "1.8"
	// DCOSDefaultRelease is the default major.minor version for DCOS
	DCOSDefaultRelease string = DCOSRelease1Dot9
)

// DCOSReleaseToVersion maps a major.minor release to an full major.minor.patch version
var DCOSReleaseToVersion = map[string]string{
	DCOSRelease1Dot10: "1.10.0",
	DCOSRelease1Dot9:  "1.9.0",
	DCOSRelease1Dot8:  "1.8.8",
}

package common

// the orchestrators supported
const (
	// Mesos is the string constant for MESOS orchestrator type
	Mesos string = "Mesos"
	// DCOS is the string constant for DCOS orchestrator type and defaults to DCOS188
	DCOS string = "DCOS"
	// Swarm is the string constant for the Swarm orchestrator type
	Swarm string = "Swarm"
	// Kubernetes is the string constant for the Kubernetes orchestrator type
	Kubernetes string = "Kubernetes"
	// SwarmMode is the string constant for the Swarm Mode orchestrator type
	SwarmMode string = "SwarmMode"
	// OpenShift is the string constant for the OpenShift orchestrator type
	OpenShift string = "OpenShift"
)

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
	// KubernetesDefaultRelease is the default Kubernetes release
	KubernetesDefaultRelease string = "1.8"
	// KubernetesDefaultReleaseWindows is the default Kubernetes release
	KubernetesDefaultReleaseWindows string = "1.9"
)

const (
	// DCOSVersion1Dot11Dot2 is the major.minor.patch string for 1.11.0 versions of DCOS
	DCOSVersion1Dot11Dot2 string = "1.11.2"
	// DCOSVersion1Dot11Dot0 is the major.minor.patch string for 1.11.0 versions of DCOS
	DCOSVersion1Dot11Dot0 string = "1.11.0"
	// DCOSVersion1Dot10Dot0 is the major.minor.patch string for 1.10.0 versions of DCOS
	DCOSVersion1Dot10Dot0 string = "1.10.0"
	// DCOSVersion1Dot9Dot0 is the major.minor.patch string for 1.9.0 versions of DCOS
	DCOSVersion1Dot9Dot0 string = "1.9.0"
	// DCOSVersion1Dot9Dot8 is the major.minor.patch string for 1.9.8 versions of DCOS
	DCOSVersion1Dot9Dot8 string = "1.9.8"
	// DCOSVersion1Dot8Dot8 is the major.minor.patch string for 1.8.8 versions of DCOS
	DCOSVersion1Dot8Dot8 string = "1.8.8"
	// DCOSDefaultVersion is the default major.minor.patch version for DCOS
	DCOSDefaultVersion string = DCOSVersion1Dot11Dot0
)

// AllDCOSSupportedVersions maintain a list of available dcos versions in acs-engine
var AllDCOSSupportedVersions = []string{
	DCOSVersion1Dot11Dot2,
	DCOSVersion1Dot11Dot0,
	DCOSVersion1Dot10Dot0,
	DCOSVersion1Dot9Dot8,
	DCOSVersion1Dot9Dot0,
	DCOSVersion1Dot8Dot8,
}

const (
	// OpenShiftVersion3Dot9Dot0 is the major.minor.patch string for the 3.9.0 version of OpenShift
	OpenShiftVersion3Dot9Dot0 string = "3.9.0"
	// OpenShiftVersionUnstable is used for development in acs-engine and should not be used by end-users.
	OpenShiftVersionUnstable string = "unstable"
	// OpenShiftDefaultVersion is the default major.minor.patch version for OpenShift
	OpenShiftDefaultVersion string = OpenShiftVersion3Dot9Dot0
)

// GetAllSupportedOpenShiftVersions returns a slice of all supported OpenShift versions.
func GetAllSupportedOpenShiftVersions() []string {
	return []string{OpenShiftVersion3Dot9Dot0, OpenShiftVersionUnstable}
}

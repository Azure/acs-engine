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
	// KubernetesVersion1Dot8Dot0 is the major.minor.patch string for the 1.8.0 version of kubernetes
	KubernetesVersion1Dot8Dot0 string = "1.8.0"
	// KubernetesVersion1Dot8Dot1 is the major.minor.patch string for the 1.8.1 version of kubernetes
	KubernetesVersion1Dot8Dot1 string = "1.8.1"
	// KubernetesVersion1Dot8Dot2 is the major.minor.patch string for the 1.8.2 version of kubernetes
	KubernetesVersion1Dot8Dot2 string = "1.8.2"
	// KubernetesVersion1Dot8Dot4 is the major.minor.patch string for the 1.8.4 version of kubernetes
	KubernetesVersion1Dot8Dot4 string = "1.8.4"
	// KubernetesVersion1Dot7Dot0 is the major.minor.patch string for the 1.7.0 version of kubernetes
	KubernetesVersion1Dot7Dot0 string = "1.7.0"
	// KubernetesVersion1Dot7Dot1 is the major.minor.patch string for the 1.7.1 version of kubernetes
	KubernetesVersion1Dot7Dot1 string = "1.7.1"
	// KubernetesVersion1Dot7Dot2 is the major.minor.patch string for the 1.7.2 version of kubernetes
	KubernetesVersion1Dot7Dot2 string = "1.7.2"
	// KubernetesVersion1Dot7Dot4 is the major.minor.patch string for the 1.7.4 version of kubernetes
	KubernetesVersion1Dot7Dot4 string = "1.7.4"
	// KubernetesVersion1Dot7Dot5 is the major.minor.patch string for the 1.7.5 version of kubernetes
	KubernetesVersion1Dot7Dot5 string = "1.7.5"
	// KubernetesVersion1Dot7Dot7 is the major.minor.patch string for the 1.7.7 version of kubernetes
	KubernetesVersion1Dot7Dot7 string = "1.7.7"
	// KubernetesVersion1Dot7Dot9 is the major.minor.patch string for the 1.7.9 version of kubernetes
	KubernetesVersion1Dot7Dot9 string = "1.7.9"
	// KubernetesVersion1Dot7Dot10 is the major.minor.patch string for the 1.7.10 version of kubernetes
	KubernetesVersion1Dot7Dot10 string = "1.7.10"
	// KubernetesVersion1Dot6Dot6 is the major.minor.patch string for the 1.6.6 version of kubernetes
	KubernetesVersion1Dot6Dot6 string = "1.6.6"
	// KubernetesVersion1Dot6Dot9 is the major.minor.patch string for the 1.6.9 version of kubernetes
	KubernetesVersion1Dot6Dot9 string = "1.6.9"
	// KubernetesVersion1Dot6Dot11 is the major.minor.patch string for the 1.6.11 version of kubernetes
	KubernetesVersion1Dot6Dot11 string = "1.6.11"
	// KubernetesVersion1Dot6Dot12 is the major.minor.patch string for the 1.6.12 version of kubernetes
	KubernetesVersion1Dot6Dot12 string = "1.6.12"
	// KubernetesVersion1Dot6Dot13 is the major.minor.patch string for the 1.6.13 version of kubernetes
	KubernetesVersion1Dot6Dot13 string = "1.6.13"
	// KubernetesVersion1Dot5Dot7 is the major.minor.patch string for the 1.5.7 version of kubernetes
	KubernetesVersion1Dot5Dot7 string = "1.5.7"
	// KubernetesVersion1Dot5Dot8 is the major.minor.patch string for the 1.5.8 version of kubernetes
	KubernetesVersion1Dot5Dot8 string = "1.5.8"
	// KubernetesDefaultVersion is the default major.minor.patch version for kubernetes
	KubernetesDefaultVersion string = KubernetesVersion1Dot7Dot9
)

// AllKubernetesSupportedVersions is a whitelist map of supported Kubernetes version strings
var AllKubernetesSupportedVersions = map[string]bool{
	KubernetesVersion1Dot5Dot7:  true,
	KubernetesVersion1Dot5Dot8:  true,
	KubernetesVersion1Dot6Dot6:  true,
	KubernetesVersion1Dot6Dot9:  true,
	KubernetesVersion1Dot6Dot11: true,
	KubernetesVersion1Dot6Dot12: true,
	KubernetesVersion1Dot6Dot13: true,
	KubernetesVersion1Dot7Dot0:  true,
	KubernetesVersion1Dot7Dot1:  true,
	KubernetesVersion1Dot7Dot2:  true,
	KubernetesVersion1Dot7Dot4:  true,
	KubernetesVersion1Dot7Dot5:  true,
	KubernetesVersion1Dot7Dot7:  true,
	KubernetesVersion1Dot7Dot9:  true,
	KubernetesVersion1Dot7Dot10: true,
	KubernetesVersion1Dot8Dot0:  true,
	KubernetesVersion1Dot8Dot1:  true,
	KubernetesVersion1Dot8Dot2:  true,
	KubernetesVersion1Dot8Dot4:  true,
}

// GetSupportedKubernetesVersion verifies that a passed-in version string is supported, or returns a default version string if not
func GetSupportedKubernetesVersion(version string) string {
	if k8sVersion := version; AllKubernetesSupportedVersions[k8sVersion] {
		return k8sVersion
	}
	return KubernetesDefaultVersion
}

// GetAllSupportedKubernetesVersions returns a slice of all supported Kubernetes versions
func GetAllSupportedKubernetesVersions() []string {
	versions := make([]string, 0, len(AllKubernetesSupportedVersions))
	for k := range AllKubernetesSupportedVersions {
		versions = append(versions, k)
	}
	return versions
}

// AllKubernetesWindowsSupportedVersions maintain a set of available k8s Windows versions in acs-engine
var AllKubernetesWindowsSupportedVersions = map[string]bool{
	KubernetesVersion1Dot7Dot2:  true,
	KubernetesVersion1Dot7Dot4:  true,
	KubernetesVersion1Dot7Dot5:  true,
	KubernetesVersion1Dot7Dot7:  true,
	KubernetesVersion1Dot7Dot9:  true,
	KubernetesVersion1Dot7Dot10: true,
	KubernetesVersion1Dot8Dot0:  true,
	KubernetesVersion1Dot8Dot1:  true,
	KubernetesVersion1Dot8Dot2:  true,
	KubernetesVersion1Dot8Dot4:  true,
}

const (
	// DCOSVersion1Dot10Dot0 is the major.minor.patch string for 1.10.0 versions of DCOS
	DCOSVersion1Dot10Dot0 string = "1.10.0"
	// DCOSVersion1Dot9Dot0 is the major.minor.patch string for 1.9.0 versions of DCOS
	DCOSVersion1Dot9Dot0 string = "1.9.0"
	// DCOSVersion1Dot8Dot8 is the major.minor.patch string for 1.8.8 versions of DCOS
	DCOSVersion1Dot8Dot8 string = "1.8.8"
	// DCOSDefaultVersion is the default major.minor.patch version for DCOS
	DCOSDefaultVersion string = DCOSVersion1Dot9Dot0
)

// AllDCOSSupportedVersions maintain a list of available dcos versions in acs-engine
var AllDCOSSupportedVersions = []string{
	DCOSVersion1Dot10Dot0,
	DCOSVersion1Dot9Dot0,
	DCOSVersion1Dot8Dot8,
}

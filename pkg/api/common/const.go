package common

import (
	"github.com/Masterminds/semver"
)

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
	// KubernetesDefaultVersion is the default Kubernetes version
	KubernetesDefaultVersion string = "1.8.9"
)

// AllKubernetesSupportedVersions is a whitelist map of supported Kubernetes version strings
var AllKubernetesSupportedVersions = map[string]bool{
	"1.6.6":         true,
	"1.6.9":         true,
	"1.6.11":        true,
	"1.6.12":        true,
	"1.6.13":        true,
	"1.7.0":         true,
	"1.7.1":         true,
	"1.7.2":         true,
	"1.7.4":         true,
	"1.7.5":         true,
	"1.7.7":         true,
	"1.7.9":         true,
	"1.7.10":        true,
	"1.7.12":        true,
	"1.7.13":        true,
	"1.7.14":        true,
	"1.7.15":        true,
	"1.8.0":         true,
	"1.8.1":         true,
	"1.8.2":         true,
	"1.8.4":         true,
	"1.8.6":         true,
	"1.8.7":         true,
	"1.8.8":         true,
	"1.8.9":         true,
	"1.8.10":        true,
	"1.9.0":         true,
	"1.9.1":         true,
	"1.9.2":         true,
	"1.9.3":         true,
	"1.9.4":         true,
	"1.9.5":         true,
	"1.9.6":         true,
	"1.10.0-beta.2": true,
	"1.10.0-beta.4": true,
	"1.10.0-rc.1":   true,
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

// GetVersionsGt returns a list of versions greater than a semver string given a list of versions
func GetVersionsGt(versions []string, version string) []string {
	// Try to get latest version matching the release
	var ret []string
	for _, v := range versions {
		sv, _ := semver.NewVersion(v)
		cons, _ := semver.NewConstraint(">" + version)
		if cons.Check(sv) {
			ret = append(ret, v)
		}
	}
	return ret
}

// AllKubernetesWindowsSupportedVersions maintain a set of available k8s Windows versions in acs-engine
var AllKubernetesWindowsSupportedVersions = getAllKubernetesWindowsSupportedVersionsMap()

func getAllKubernetesWindowsSupportedVersionsMap() map[string]bool {
	ret := make(map[string]bool)
	for k, v := range AllKubernetesSupportedVersions {
		ret[k] = v
	}
	for _, version := range []string{
		"1.6.6",
		"1.6.9",
		"1.6.11",
		"1.6.12",
		"1.6.13",
		"1.7.0",
		"1.7.1",
		"1.10.0-beta.2",
		"1.10.0-beta.4",
		"1.10.0-rc.1"} {
		ret[version] = false
	}
	return ret
}

// GetAllSupportedKubernetesVersionsWindows returns a slice of all supported Kubernetes versions on Windows
func GetAllSupportedKubernetesVersionsWindows() []string {
	versions := make([]string, 0, len(AllKubernetesWindowsSupportedVersions))
	for k := range AllKubernetesWindowsSupportedVersions {
		versions = append(versions, k)
	}
	return versions
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

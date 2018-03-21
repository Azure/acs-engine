package common

import "github.com/Masterminds/semver"

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

package common

import (
	"fmt"
	"strings"

	"github.com/Masterminds/semver"
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
	"1.7.16":        true,
	"1.8.0":         true,
	"1.8.1":         true,
	"1.8.2":         true,
	"1.8.4":         true,
	"1.8.6":         true,
	"1.8.7":         true,
	"1.8.8":         true,
	"1.8.9":         true,
	"1.8.10":        true,
	"1.8.11":        true,
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
	"1.10.0":        true,
	"1.10.1":        true,
}

// GetDefaultKubernetesVersion returns the default Kubernetes version, that is the latest patch of the default release
func GetDefaultKubernetesVersion() string {
	return GetLatestPatchVersion(KubernetesDefaultRelease, GetAllSupportedKubernetesVersions())
}

// GetSupportedKubernetesVersion verifies that a passed-in version string is supported, or returns a default version string if not
func GetSupportedKubernetesVersion(version string) string {
	if k8sVersion := version; AllKubernetesSupportedVersions[k8sVersion] {
		return k8sVersion
	}
	return GetDefaultKubernetesVersion()
}

// GetAllSupportedKubernetesVersions returns a slice of all supported Kubernetes versions
func GetAllSupportedKubernetesVersions() []string {
	var versions []string
	for ver, supported := range AllKubernetesSupportedVersions {
		if supported {
			versions = append(versions, ver)
		}
	}
	return versions
}

// GetVersionsGt returns a list of versions greater than a semver string given a list of versions
// inclusive=true means that we test for equality as well
// preReleases=true means that we include pre-release versions in the list
func GetVersionsGt(versions []string, version string, inclusive, preReleases bool) []string {
	// Try to get latest version matching the release
	var ret []string
	var cons *semver.Constraints
	minVersion, _ := semver.NewVersion(version)
	for _, v := range versions {
		sv, _ := semver.NewVersion(v)
		if preReleases && minVersion.Prerelease() == "" {
			sv, _ = semver.NewVersion(fmt.Sprintf("%d.%d.%d", sv.Major(), sv.Minor(), sv.Patch()))
		}
		if inclusive {
			cons, _ = semver.NewConstraint(">=" + version)
		} else {
			cons, _ = semver.NewConstraint(">" + version)
		}
		if cons.Check(sv) {
			ret = append(ret, v)
		}
	}
	return ret
}

// GetVersionsLt returns a list of versions less than than a semver string given a list of versions
// inclusive=true means that we test for equality as well
// preReleases=true means that we include pre-release versions in the list
func GetVersionsLt(versions []string, version string, inclusive, preReleases bool) []string {
	// Try to get latest version matching the release
	var ret []string
	var cons *semver.Constraints
	for _, v := range versions {
		sv, _ := semver.NewVersion(v)
		if preReleases {
			sv, _ = semver.NewVersion(fmt.Sprintf("%d.%d.%d", sv.Major(), sv.Minor(), sv.Patch()))
		}
		if inclusive {
			cons, _ = semver.NewConstraint("<=" + version)
		} else {
			cons, _ = semver.NewConstraint("<" + version)
		}
		if cons.Check(sv) {
			ret = append(ret, v)
		}
	}
	return ret
}

// GetVersionsBetween returns a list of versions between a min and max
// inclusive=true means that we test for equality on both bounds
// preReleases=true means that we include pre-release versions in the list
func GetVersionsBetween(versions []string, versionMin, versionMax string, inclusive, preReleases bool) []string {
	var ret []string
	if minV, _ := semver.NewVersion(versionMin); minV.Prerelease() != "" {
		preReleases = true
	}
	greaterThan := GetVersionsGt(versions, versionMin, inclusive, preReleases)
	lessThan := GetVersionsLt(versions, versionMax, inclusive, preReleases)
	for _, lv := range lessThan {
		for _, gv := range greaterThan {
			if lv == gv {
				ret = append(ret, lv)
			}
		}
	}
	return ret
}

// GetMaxVersion gets the highest semver version
// preRelease=true means accept a pre-release version as a max value
func GetMaxVersion(versions []string, preRelease bool) string {
	if len(versions) < 1 {
		return ""
	}
	highest, _ := semver.NewVersion("0.0.0")
	highestPreRelease, _ := semver.NewVersion("0.0.0-alpha.0")
	var preReleaseVersions []*semver.Version
	for _, v := range versions {
		sv, _ := semver.NewVersion(v)
		if sv.Prerelease() != "" {
			preReleaseVersions = append(preReleaseVersions, sv)
		} else {
			if sv.Compare(highest) == 1 {
				highest = sv
			}
		}
	}
	if preRelease {
		for _, sv := range preReleaseVersions {
			if sv.Compare(highestPreRelease) == 1 {
				highestPreRelease = sv
			}
		}
		switch highestPreRelease.Compare(highest) {
		case 1:
			return highestPreRelease.String()
		default:
			return highest.String()
		}

	}
	return highest.String()
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
	var versions []string
	for ver, supported := range AllKubernetesWindowsSupportedVersions {
		if supported {
			versions = append(versions, ver)
		}
	}
	return versions
}

// GetSupportedVersions get supported version list for a certain orchestrator
func GetSupportedVersions(orchType string, hasWindows bool) (versions []string, defaultVersion string) {
	switch orchType {
	case Kubernetes:
		if hasWindows {
			return GetAllSupportedKubernetesVersionsWindows(), GetDefaultKubernetesVersion()
		}
		return GetAllSupportedKubernetesVersions(), GetDefaultKubernetesVersion()

	case DCOS:
		return AllDCOSSupportedVersions, DCOSDefaultVersion
	default:
		return nil, ""
	}
}

//GetValidPatchVersion gets the current valid patch version for the minor version of the passed in version
func GetValidPatchVersion(orchType, orchVer string) string {
	if orchVer == "" {
		return RationalizeReleaseAndVersion(
			orchType,
			"",
			"",
			false)
	}

	// check if the current version is valid, this allows us to have multiple supported patch versions in the future if we need it
	version := RationalizeReleaseAndVersion(
		orchType,
		"",
		orchVer,
		false)

	if version == "" {
		sv, err := semver.NewVersion(orchVer)
		if err != nil {
			return ""
		}
		sr := fmt.Sprintf("%d.%d", sv.Major(), sv.Minor())

		version = RationalizeReleaseAndVersion(
			orchType,
			sr,
			"",
			false)
	}
	return version
}

// RationalizeReleaseAndVersion return a version when it can be rationalized from the input, otherwise ""
func RationalizeReleaseAndVersion(orchType, orchRel, orchVer string, hasWindows bool) (version string) {
	// ignore "v" prefix in orchestrator version and release: "v1.8.0" is equivalent to "1.8.0", "v1.9" is equivalent to "1.9"
	orchVer = strings.TrimPrefix(orchVer, "v")
	orchRel = strings.TrimPrefix(orchRel, "v")
	supportedVersions, defaultVersion := GetSupportedVersions(orchType, hasWindows)
	if supportedVersions == nil {
		return ""
	}

	if orchRel == "" && orchVer == "" {
		return defaultVersion
	}

	if orchVer == "" {
		// Try to get latest version matching the release
		version = GetLatestPatchVersion(orchRel, supportedVersions)
		return version
	} else if orchRel == "" {
		// Try to get version the same with orchVer
		version = ""
		for _, ver := range supportedVersions {
			if ver == orchVer {
				version = ver
				break
			}
		}
		return version
	}
	// Try to get latest version matching the release
	version = ""
	for _, ver := range supportedVersions {
		sv, _ := semver.NewVersion(ver)
		sr := fmt.Sprintf("%d.%d", sv.Major(), sv.Minor())
		if sr == orchRel && ver == orchVer {
			version = ver
			break
		}
	}
	return version
}

// IsKubernetesVersionGe returns if a semver string is >= to a compare-against semver string (suppport "-" suffixes)
func IsKubernetesVersionGe(actualVersion, version string) bool {
	orchestratorVersion, _ := semver.NewVersion(strings.Split(actualVersion, "-")[0]) // to account for -alpha and -beta suffixes
	constraint, _ := semver.NewConstraint(">=" + version)
	return constraint.Check(orchestratorVersion)
}

// GetLatestPatchVersion gets the most recent patch version from a list of semver versions given a major.minor string
func GetLatestPatchVersion(majorMinor string, versionsList []string) (version string) {
	// Try to get latest version matching the release
	version = ""
	for _, ver := range versionsList {
		sv, err := semver.NewVersion(ver)
		if err != nil {
			return
		}
		sr := fmt.Sprintf("%d.%d", sv.Major(), sv.Minor())
		if sr == majorMinor {
			if version == "" {
				version = ver
			} else {
				cons, _ := semver.NewConstraint(">" + version)
				if cons.Check(sv) {
					version = ver
				}
			}
		}
	}
	return version
}

package common

import (
	"fmt"
	"sort"
	"strings"

	"github.com/blang/semver"
)

// AllKubernetesSupportedVersions is a whitelist map of supported Kubernetes version strings
var AllKubernetesSupportedVersions = map[string]bool{
	"1.6.6":          true,
	"1.6.9":          true,
	"1.6.11":         true,
	"1.6.12":         true,
	"1.6.13":         true,
	"1.7.0":          true,
	"1.7.1":          true,
	"1.7.2":          true,
	"1.7.4":          true,
	"1.7.5":          true,
	"1.7.7":          true,
	"1.7.9":          true,
	"1.7.10":         true,
	"1.7.12":         true,
	"1.7.13":         true,
	"1.7.14":         true,
	"1.7.15":         true,
	"1.7.16":         true,
	"1.8.0":          true,
	"1.8.1":          true,
	"1.8.2":          true,
	"1.8.4":          true,
	"1.8.6":          true,
	"1.8.7":          true,
	"1.8.8":          true,
	"1.8.9":          true,
	"1.8.10":         true,
	"1.8.11":         true,
	"1.8.12":         true,
	"1.8.13":         true,
	"1.8.14":         true,
	"1.8.15":         true,
	"1.9.0":          true,
	"1.9.1":          true,
	"1.9.2":          true,
	"1.9.3":          true,
	"1.9.4":          true,
	"1.9.5":          true,
	"1.9.6":          true,
	"1.9.7":          true,
	"1.9.8":          true,
	"1.9.9":          true,
	"1.9.10":         true,
	"1.10.0-beta.2":  true,
	"1.10.0-beta.4":  true,
	"1.10.0-rc.1":    true,
	"1.10.0":         true,
	"1.10.1":         true,
	"1.10.2":         true,
	"1.10.3":         true,
	"1.10.4":         true,
	"1.10.5":         true,
	"1.10.6":         true,
	"1.11.0-alpha.1": true,
	"1.11.0-alpha.2": true,
	"1.11.0-beta.1":  true,
	"1.11.0-beta.2":  true,
	"1.11.0-rc.1":    true,
	"1.11.0-rc.2":    true,
	"1.11.0-rc.3":    true,
	"1.11.0":         true,
	"1.11.1":         true,
	"1.11.2":         true,
	"1.12.0-alpha.1": true,
}

// GetDefaultKubernetesVersion returns the default Kubernetes version, that is the latest patch of the default release
func GetDefaultKubernetesVersion() string {
	return GetLatestPatchVersion(KubernetesDefaultRelease, GetAllSupportedKubernetesVersions())
}

// GetDefaultKubernetesVersionWindows returns the default Kubernetes version for Windows, that is the latest patch of the default release
func GetDefaultKubernetesVersionWindows() string {
	return GetLatestPatchVersion(KubernetesDefaultReleaseWindows, GetAllSupportedKubernetesVersionsWindows())
}

// GetSupportedKubernetesVersion verifies that a passed-in version string is supported, or returns a default version string if not
func GetSupportedKubernetesVersion(version string, hasWindows bool) string {
	var k8sVersion string
	if hasWindows {
		k8sVersion = GetDefaultKubernetesVersionWindows()
		if AllKubernetesWindowsSupportedVersions[version] {
			k8sVersion = version
		}
	} else {
		k8sVersion = GetDefaultKubernetesVersion()
		if AllKubernetesSupportedVersions[version] {
			k8sVersion = version
		}
	}
	return k8sVersion
}

// GetAllSupportedKubernetesVersions returns a slice of all supported Kubernetes versions
func GetAllSupportedKubernetesVersions() []string {
	var versions []string
	for ver, supported := range AllKubernetesSupportedVersions {
		if supported {
			versions = append(versions, ver)
		}
	}
	sort.Slice(versions, func(i, j int) bool {
		return IsKubernetesVersionGe(versions[j], versions[i])
	})
	return versions
}

// GetVersionsGt returns a list of versions greater than a semver string given a list of versions
// inclusive=true means that we test for equality as well
// preReleases=true means that we include pre-release versions in the list
func GetVersionsGt(versions []string, version string, inclusive, preReleases bool) []string {
	// Try to get latest version matching the release
	var ret []string
	minVersion, _ := semver.Make(version)
	for _, v := range versions {
		sv, _ := semver.Make(v)
		if !preReleases && len(sv.Pre) != 0 {
			continue
		}
		if (inclusive && sv.GTE(minVersion)) || (!inclusive && sv.GT(minVersion)) {
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
	minVersion, _ := semver.Make(version)
	for _, v := range versions {
		sv, _ := semver.Make(v)
		if !preReleases && len(sv.Pre) != 0 {
			continue
		}
		if (inclusive && sv.LTE(minVersion)) || (!inclusive && sv.LT(minVersion)) {
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
	if minV, _ := semver.Make(versionMin); len(minV.Pre) != 0 {
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
	highest, _ := semver.Make("0.0.0")
	highestPreRelease, _ := semver.Make("0.0.0-alpha.0")
	var preReleaseVersions []semver.Version
	for _, v := range versions {
		sv, _ := semver.Make(v)
		if len(sv.Pre) != 0 {
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
		"1.8.13",
		"1.8.14",
		"1.8.15",
		"1.10.0-beta.2",
		"1.10.0-beta.4",
		"1.10.0-rc.1",
		"1.11.0-alpha.1",
		"1.11.0-alpha.2"} {
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
	sort.Slice(versions, func(i, j int) bool {
		return IsKubernetesVersionGe(versions[j], versions[i])
	})
	return versions
}

// GetSupportedVersions get supported version list for a certain orchestrator
func GetSupportedVersions(orchType string, hasWindows bool) (versions []string, defaultVersion string) {
	switch orchType {
	case Kubernetes:
		if hasWindows {
			return GetAllSupportedKubernetesVersionsWindows(), GetDefaultKubernetesVersionWindows()
		}
		return GetAllSupportedKubernetesVersions(), GetDefaultKubernetesVersion()

	case OpenShift:
		return GetAllSupportedOpenShiftVersions(), string(OpenShiftDefaultVersion)

	case DCOS:
		return AllDCOSSupportedVersions, DCOSDefaultVersion
	default:
		return nil, ""
	}
}

//GetValidPatchVersion gets the current valid patch version for the minor version of the passed in version
func GetValidPatchVersion(orchType, orchVer string, hasWindows bool) string {
	if orchVer == "" {
		return RationalizeReleaseAndVersion(
			orchType,
			"",
			"",
			hasWindows)
	}

	// check if the current version is valid, this allows us to have multiple supported patch versions in the future if we need it
	version := RationalizeReleaseAndVersion(
		orchType,
		"",
		orchVer,
		hasWindows)

	if version == "" {
		sv, err := semver.Make(orchVer)
		if err != nil {
			return ""
		}
		sr := fmt.Sprintf("%d.%d", sv.Major, sv.Minor)

		version = RationalizeReleaseAndVersion(
			orchType,
			sr,
			"",
			hasWindows)
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
		sv, _ := semver.Make(ver)
		sr := fmt.Sprintf("%d.%d", sv.Major, sv.Minor)
		if sr == orchRel && ver == orchVer {
			version = ver
			break
		}
	}
	return version
}

// IsKubernetesVersionGe returns true if actualVersion is greater than or equal to version
func IsKubernetesVersionGe(actualVersion, version string) bool {
	v1, _ := semver.Make(actualVersion)
	v2, _ := semver.Make(version)
	return v1.GE(v2)
}

// GetLatestPatchVersion gets the most recent patch version from a list of semver versions given a major.minor string
func GetLatestPatchVersion(majorMinor string, versionsList []string) (version string) {
	// Try to get latest version matching the release
	version = ""
	for _, ver := range versionsList {
		sv, err := semver.Make(ver)
		if err != nil {
			return
		}
		sr := fmt.Sprintf("%d.%d", sv.Major, sv.Minor)
		if sr == majorMinor {
			if version == "" {
				version = ver
			} else {
				current, _ := semver.Make(version)
				if sv.GT(current) {
					version = ver
				}
			}
		}
	}
	return version
}

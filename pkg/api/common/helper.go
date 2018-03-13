package common

import (
	"fmt"
	"strings"

	"github.com/Masterminds/semver"

	validator "gopkg.in/go-playground/validator.v9"
)

// HandleValidationErrors is the helper function to catch validator.ValidationError
// based on Namespace of the error, and return customized error message.
func HandleValidationErrors(e validator.ValidationErrors) error {
	err := e[0]
	ns := err.Namespace()
	switch ns {
	case "Properties.OrchestratorProfile", "Properties.OrchestratorProfile.OrchestratorType",
		"Properties.MasterProfile", "Properties.MasterProfile.DNSPrefix", "Properties.MasterProfile.VMSize",
		"Properties.LinuxProfile", "Properties.ServicePrincipalProfile.ClientID",
		"Properties.WindowsProfile.AdminUsername",
		"Properties.WindowsProfile.AdminPassword":
		return fmt.Errorf("missing %s", ns)
	case "Properties.MasterProfile.Count":
		return fmt.Errorf("MasterProfile count needs to be 1, 3, or 5")
	case "Properties.MasterProfile.OSDiskSizeGB":
		return fmt.Errorf("Invalid os disk size of %d specified.  The range of valid values are [%d, %d]", err.Value().(int), MinDiskSizeGB, MaxDiskSizeGB)
	case "Properties.MasterProfile.IPAddressCount":
		return fmt.Errorf("MasterProfile.IPAddressCount needs to be in the range [%d,%d]", MinIPAddressCount, MaxIPAddressCount)
	case "Properties.MasterProfile.StorageProfile":
		return fmt.Errorf("Unknown storageProfile '%s'. Specify either %s or %s", err.Value().(string), StorageAccount, ManagedDisks)
	default:
		if strings.HasPrefix(ns, "Properties.AgentPoolProfiles") {
			switch {
			case strings.HasSuffix(ns, ".Name") || strings.HasSuffix(ns, "VMSize"):
				return fmt.Errorf("missing %s", ns)
			case strings.HasSuffix(ns, ".Count"):
				return fmt.Errorf("AgentPoolProfile count needs to be in the range [%d,%d]", MinAgentCount, MaxAgentCount)
			case strings.HasSuffix(ns, ".OSDiskSizeGB"):
				return fmt.Errorf("Invalid os disk size of %d specified.  The range of valid values are [%d, %d]", err.Value().(int), MinDiskSizeGB, MaxDiskSizeGB)
			case strings.Contains(ns, ".Ports"):
				return fmt.Errorf("AgentPoolProfile Ports must be in the range[%d, %d]", MinPort, MaxPort)
			case strings.HasSuffix(ns, ".StorageProfile"):
				return fmt.Errorf("Unknown storageProfile '%s'. Specify either %s or %s", err.Value().(string), StorageAccount, ManagedDisks)
			case strings.Contains(ns, ".DiskSizesGB"):
				return fmt.Errorf("A maximum of %d disks may be specified, The range of valid disk size values are [%d, %d]", MaxDisks, MinDiskSizeGB, MaxDiskSizeGB)
			case strings.HasSuffix(ns, ".IPAddressCount"):
				return fmt.Errorf("AgentPoolProfile.IPAddressCount needs to be in the range [%d,%d]", MinIPAddressCount, MaxIPAddressCount)
			default:
				break
			}
		}
	}
	return fmt.Errorf("Namespace %s is not caught, %+v", ns, e)
}

// GetSupportedVersions get supported version list for a certain orchestrator
func GetSupportedVersions(orchType string, hasWindows bool) (versions []string, defaultVersion string) {
	switch orchType {
	case Kubernetes:
		if hasWindows {
			return GetAllSupportedKubernetesVersionsWindows(), string(KubernetesDefaultVersion)
		}
		return GetAllSupportedKubernetesVersions(), string(KubernetesDefaultVersion)

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
		sv, _ := semver.NewVersion(ver)
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

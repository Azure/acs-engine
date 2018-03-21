package common

import (
	"testing"
)

func Test_GetValidPatchVersion(t *testing.T) {
	v := GetValidPatchVersion(Kubernetes, "")
	if v != KubernetesDefaultVersion {
		t.Errorf("It is not the default Kubernetes version")
	}

	for version, enabled := range AllKubernetesSupportedVersions {
		if enabled {
			v = GetValidPatchVersion(Kubernetes, version)
			if v != version {
				t.Errorf("Expected version %s, instead got version %s", version, v)
			}
		}
	}
}

func TestGetLatestPatchVersion(t *testing.T) {
	expected := "1.1.2"
	version := GetLatestPatchVersion("1.1", []string{"1.1.1", expected})
	if version != expected {
		t.Errorf("GetLatestPatchVersion returned the wrong latest version, expected %s, got %s", expected, version)
	}

	expected = "1.1.2"
	version = GetLatestPatchVersion("1.1", []string{"1.1.0", expected})
	if version != expected {
		t.Errorf("GetLatestPatchVersion returned the wrong latest version, expected %s, got %s", expected, version)
	}

	expected = "1.2.0"
	version = GetLatestPatchVersion("1.2", []string{"1.1.0", "1.3.0", expected})
	if version != expected {
		t.Errorf("GetLatestPatchVersion returned the wrong latest version, expected %s, got %s", expected, version)
	}
}

func Test_RationalizeReleaseAndVersion(t *testing.T) {
	version := RationalizeReleaseAndVersion(Kubernetes, "", "", false)
	if version != KubernetesDefaultVersion {
		t.Errorf("It is not the default Kubernetes version")
	}

	latest1Dot6Version := GetLatestPatchVersion("1.6", GetAllSupportedKubernetesVersions())
	version = RationalizeReleaseAndVersion(Kubernetes, "1.6", "", false)
	if version != latest1Dot6Version {
		t.Errorf("It is not Kubernetes version %s", latest1Dot6Version)
	}

	expectedVersion := "1.6.11"
	version = RationalizeReleaseAndVersion(Kubernetes, "", expectedVersion, false)
	if version != expectedVersion {
		t.Errorf("It is not Kubernetes version %s", expectedVersion)
	}

	version = RationalizeReleaseAndVersion(Kubernetes, "1.6", expectedVersion, false)
	if version != expectedVersion {
		t.Errorf("It is not Kubernetes version %s", expectedVersion)
	}

	version = RationalizeReleaseAndVersion(Kubernetes, "", "v"+expectedVersion, false)
	if version != expectedVersion {
		t.Errorf("It is not Kubernetes version %s", expectedVersion)
	}

	version = RationalizeReleaseAndVersion(Kubernetes, "v1.6", "v"+expectedVersion, false)
	if version != expectedVersion {
		t.Errorf("It is not Kubernetes version %s", expectedVersion)
	}

	version = RationalizeReleaseAndVersion(Kubernetes, "1.1", "", false)
	if version != "" {
		t.Errorf("It is not empty string")
	}

	version = RationalizeReleaseAndVersion(Kubernetes, "1.1", "1.6.6", false)
	if version != "" {
		t.Errorf("It is not empty string")
	}
}

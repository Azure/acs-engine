package common

import (
	"testing"
)

func Test_GetAllSupportedKubernetesVersions(t *testing.T) {
	responseFromGetter := GetAllSupportedKubernetesVersions()

	if len(AllKubernetesSupportedVersions) != len(responseFromGetter) {
		t.Errorf("GetAllSupportedKubernetesVersions() returned %d items, expected %d", len(responseFromGetter), len(AllKubernetesSupportedVersions))
	}

	for _, version := range responseFromGetter {
		if !AllKubernetesSupportedVersions[version] {
			t.Errorf("GetAllSupportedKubernetesVersions() returned a version %s that was not in the definitive AllKubernetesSupportedVersions map", version)
		}
	}
}

func Test_GetSupportedKubernetesVersion(t *testing.T) {
	versions := GetAllSupportedKubernetesVersions()
	for _, version := range versions {
		supportedVersion := GetSupportedKubernetesVersion(version)
		if supportedVersion != version {
			t.Errorf("GetSupportedKubernetesVersion(%s) should return the same passed-in string, instead returned %s", version, supportedVersion)
		}
	}

	defaultVersion := GetSupportedKubernetesVersion("")
	if defaultVersion != KubernetesDefaultVersion {
		t.Errorf("GetSupportedKubernetesVersion(\"\") should return the default version %s, instead returned %s", KubernetesDefaultVersion, defaultVersion)
	}
}

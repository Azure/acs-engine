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

func TestGetVersionsGt(t *testing.T) {
	versions := []string{"1.1.0", "1.2.0", "1.2.1"}
	expected := []string{"1.2.0", "1.2.1"}
	expectedMap := map[string]bool{
		"1.2.0": true,
		"1.2.1": true,
	}
	v := GetVersionsGt(versions, "1.1.0")
	errStr := "GetVersionsGt returned an unexpected list of strings"
	if len(v) != len(expected) {
		t.Errorf(errStr)
	}
	for _, ver := range v {
		if !expectedMap[ver] {
			t.Errorf(errStr)
		}
	}
}

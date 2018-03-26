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
	versions := []string{"1.1.0-rc.1", "1.2.0-rc.1", "1.2.0", "1.2.1"}
	expected := []string{"1.1.0-rc.1", "1.2.0-rc.1", "1.2.0", "1.2.1"}
	expectedMap := map[string]bool{
		"1.1.0-rc.1": true,
		"1.2.0-rc.1": true,
		"1.2.0":      true,
		"1.2.1":      true,
	}
	v := GetVersionsGt(versions, "1.1.0-alpha.1", false, false)
	errStr := "GetVersionsGt returned an unexpected list of strings"
	if len(v) != len(expected) {
		t.Errorf(errStr)
	}
	for _, ver := range v {
		if !expectedMap[ver] {
			t.Errorf(errStr)
		}
	}

	versions = []string{"1.1.0", "1.2.0", "1.2.1"}
	expected = []string{"1.1.0", "1.2.0", "1.2.1"}
	expectedMap = map[string]bool{
		"1.1.0": true,
		"1.2.0": true,
		"1.2.1": true,
	}
	v = GetVersionsGt(versions, "1.1.0", true, false)
	if len(v) != len(expected) {
		t.Errorf(errStr)
	}
	for _, ver := range v {
		if !expectedMap[ver] {
			t.Errorf(errStr)
		}
	}
}

func TestGetVersionsLt(t *testing.T) {
	versions := []string{"1.1.0", "1.2.0-rc.1", "1.2.0", "1.2.1"}
	expected := []string{"1.2.0", "1.2.1"}
	// less than comparisons exclude pre-release versions from the result
	expectedMap := map[string]bool{
		"1.1.0": true,
		"1.2.0": true,
	}
	v := GetVersionsLt(versions, "1.2.1", false, false)
	errStr := "GetVersionsLt returned an unexpected list of strings"
	if len(v) != len(expected) {
		t.Errorf(errStr)
	}
	for _, ver := range v {
		if !expectedMap[ver] {
			t.Errorf(errStr)
		}
	}

	versions = []string{"1.1.0", "1.2.0", "1.2.1"}
	expected = []string{"1.1.0", "1.2.0", "1.2.1"}
	expectedMap = map[string]bool{
		"1.1.0": true,
		"1.2.0": true,
		"1.2.1": true,
	}
	v = GetVersionsLt(versions, "1.2.1", true, false)
	if len(v) != len(expected) {
		t.Errorf(errStr)
	}
	for _, ver := range v {
		if !expectedMap[ver] {
			t.Errorf(errStr)
		}
	}
}

func TestGetVersionsBetween(t *testing.T) {
	versions := []string{"1.1.0", "1.2.0", "1.2.1"}
	expected := []string{"1.2.0"}
	expectedMap := map[string]bool{
		"1.2.0": true,
	}
	v := GetVersionsBetween(versions, "1.1.0", "1.2.1", false, false)
	errStr := "GetVersionsBetween returned an unexpected list of strings"
	if len(v) != len(expected) {
		t.Errorf(errStr)
	}
	for _, ver := range v {
		if !expectedMap[ver] {
			t.Errorf(errStr)
		}
	}

	versions = []string{"1.1.0", "1.2.0", "1.2.1"}
	expected = []string{"1.1.0", "1.2.0", "1.2.1"}
	expectedMap = map[string]bool{
		"1.1.0": true,
		"1.2.0": true,
		"1.2.1": true,
	}
	v = GetVersionsBetween(versions, "1.1.0", "1.2.1", true, false)
	if len(v) != len(expected) {
		t.Errorf(errStr)
	}
	for _, ver := range v {
		if !expectedMap[ver] {
			t.Errorf(errStr)
		}
	}

	versions = []string{"1.9.6", "1.10.0-beta.2", "1.10.0-beta.4", "1.10.0-rc.1"}
	expected = []string{"1.10.0-beta.2", "1.10.0-beta.4", "1.10.0-rc.1"}
	expectedMap = map[string]bool{
		"1.10.0-beta.2": true,
		"1.10.0-beta.4": true,
		"1.10.0-rc.1":   true,
	}
	v = GetVersionsBetween(versions, "1.9.6", "1.11.0", false, true)
	if len(v) != len(expected) {
		t.Errorf(errStr)
	}
	for _, ver := range v {
		if !expectedMap[ver] {
			t.Errorf(errStr)
		}
	}
	v = GetVersionsBetween(versions, "1.9.6", "1.11.0", false, false)
	if len(v) != 0 {
		t.Errorf(errStr)
	}

	versions = []string{"1.9.6", "1.10.0-beta.2", "1.10.0-beta.4", "1.10.0-rc.1"}
	expected = []string{"1.10.0-beta.4", "1.10.0-rc.1"}
	expectedMap = map[string]bool{
		"1.10.0-beta.4": true,
		"1.10.0-rc.1":   true,
	}
	v = GetVersionsBetween(versions, "1.10.0-beta.2", "1.12.0", false, false)
	if len(v) != len(expected) {
		t.Errorf(errStr)
	}
	for _, ver := range v {
		if !expectedMap[ver] {
			t.Errorf(errStr)
		}
	}

	versions = []string{"1.10.0", "1.10.0-beta.2", "1.10.0-beta.4", "1.10.0-rc.1"}
	v = GetVersionsBetween(versions, "1.10.0", "1.12.0", false, false)
	if len(v) != 0 {
		t.Errorf(errStr)
	}

	versions = []string{"1.9.6", "1.10.0-beta.2", "1.10.0-beta.4", "1.10.0-rc.1"}
	expectedMap = map[string]bool{
		"1.9.6":         true,
		"1.10.0-beta.2": true,
		"1.10.0-beta.4": true,
		"1.10.0-rc.1":   true,
	}
	v = GetVersionsBetween(versions, "1.9.5", "1.12.0", false, true)
	if len(v) != len(versions) {
		t.Errorf(errStr)
	}
	for _, ver := range v {
		if !expectedMap[ver] {
			t.Errorf(errStr)
		}
	}

	versions = []string{"1.9.6", "1.10.0", "1.10.1", "1.10.2"}
	expected = []string{"1.10.0", "1.10.1", "1.10.2"}
	expectedMap = map[string]bool{
		"1.10.0": true,
		"1.10.1": true,
		"1.10.2": true,
	}
	v = GetVersionsBetween(versions, "1.10.0-rc.1", "1.12.0", false, true)
	if len(v) != len(expected) {
		t.Errorf(errStr)
	}
	for _, ver := range v {
		if !expectedMap[ver] {
			t.Errorf(errStr)
		}
	}
}

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

	expected = "1.2.0-rc.3"
	version = GetLatestPatchVersion("1.2", []string{"1.2.0-alpha.1", "1.2.0-beta.1", "1.2.0-rc.3", expected})
	if version != expected {
		t.Errorf("GetLatestPatchVersion returned the wrong latest version, expected %s, got %s", expected, version)
	}

	expected = ""
	version = GetLatestPatchVersion("1.2", []string{"1.1.0", "1.1.1", "1.1.2", expected})
	if version != expected {
		t.Errorf("GetLatestPatchVersion returned the wrong latest version, expected %s, got %s", expected, version)
	}
}

func TestGetMaxVersion(t *testing.T) {
	expected := "1.0.3"
	versions := []string{"1.0.1", "1.0.2", expected}
	max := GetMaxVersion(versions, false)
	if max != expected {
		t.Errorf("GetMaxVersion returned the wrong max version, expected %s, got %s", expected, max)
	}

	expected = "1.2.3"
	versions = []string{"1.0.1", "1.1.2", expected}
	max = GetMaxVersion(versions, false)
	if max != expected {
		t.Errorf("GetMaxVersion returned the wrong max version, expected %s, got %s", expected, max)
	}

	expected = "1.1.2"
	versions = []string{"1.0.1", expected, "1.2.3-alpha.1"}
	max = GetMaxVersion(versions, false)
	if max != expected {
		t.Errorf("GetMaxVersion returned the wrong max version, expected %s, got %s", expected, max)
	}

	expected = "1.2.3-alpha.1"
	versions = []string{"1.0.1", "1.1.2", expected}
	max = GetMaxVersion(versions, true)
	if max != expected {
		t.Errorf("GetMaxVersion returned the wrong max version, expected %s, got %s", expected, max)
	}

	expected = ""
	versions = []string{}
	max = GetMaxVersion(versions, false)
	if max != expected {
		t.Errorf("GetMaxVersion returned the wrong max version, expected %s, got %s", expected, max)
	}

	expected = ""
	versions = []string{}
	max = GetMaxVersion(versions, true)
	if max != expected {
		t.Errorf("GetMaxVersion returned the wrong max version, expected %s, got %s", expected, max)
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

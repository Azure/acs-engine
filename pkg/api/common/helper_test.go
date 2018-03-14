package common

import "testing"

func Test_GetValidPatchVersion(t *testing.T) {
	version := GetValidPatchVersion(Kubernetes, "")
	if version != KubernetesDefaultVersion {
		t.Errorf("It is not the default Kubernetes version")
	}

	version = GetValidPatchVersion(Kubernetes, "1.6.3")
	if version != KubernetesVersion1Dot6Dot13 {
		t.Errorf("It is not Kubernetes version %s", KubernetesVersion1Dot6Dot13)
	}

	version = GetValidPatchVersion(Kubernetes, "1.7.3")
	if version != KubernetesVersion1Dot7Dot14 {
		t.Errorf("It is not Kubernetes version %s", KubernetesVersion1Dot7Dot14)
	}

	version = GetValidPatchVersion(Kubernetes, "1.8.7")
	if version != KubernetesVersion1Dot8Dot7 {
		t.Errorf("It is not Kubernetes version %s", KubernetesVersion1Dot8Dot7)
	}

	version = GetValidPatchVersion(Kubernetes, "1.9.1")
	if version != KubernetesVersion1Dot9Dot1 {
		t.Errorf("It is not Kubernetes version %s", KubernetesVersion1Dot9Dot1)
	}

	version = GetValidPatchVersion(Kubernetes, "1.9.2")
	if version != KubernetesVersion1Dot9Dot2 {
		t.Errorf("It is not Kubernetes version %s", KubernetesVersion1Dot9Dot2)
	}

	version = GetValidPatchVersion(Kubernetes, "1.10.0")
	if version != KubernetesVersion1Dot10Dot0 {
		t.Errorf("It is not Kubernetes version %s", KubernetesVersion1Dot10Dot0)
	}
}

func TestGetLatestPatchVersion(t *testing.T) {
	version := GetLatestPatchVersion("1.6", GetAllSupportedKubernetesVersions())
	if version != KubernetesVersion1Dot6Dot13 {
		t.Errorf("It is not Kubernetes version %s", KubernetesVersion1Dot6Dot13)
	}

	version = GetLatestPatchVersion("1.7", GetAllSupportedKubernetesVersions())
	if version != KubernetesVersion1Dot7Dot14 {
		t.Errorf("It is not Kubernetes version %s", KubernetesVersion1Dot7Dot14)
	}

	version = GetLatestPatchVersion("1.8", GetAllSupportedKubernetesVersions())
	if version != KubernetesVersion1Dot8Dot9 {
		t.Errorf("It is not Kubernetes version %s", KubernetesVersion1Dot8Dot9)
	}

	version = GetLatestPatchVersion("1.9", GetAllSupportedKubernetesVersions())
	if version != KubernetesVersion1Dot9Dot4 {
		t.Errorf("It is not Kubernetes version %s", KubernetesVersion1Dot9Dot4)
	}

	version = GetLatestPatchVersion("1.10", GetAllSupportedKubernetesVersions())
	if version != KubernetesVersion1Dot10Dot0 {
		t.Errorf("It is not Kubernetes version %s", KubernetesVersion1Dot10Dot0)
	}

	expected := "99.1.2"
	version = GetLatestPatchVersion("99.1", []string{"99.1.1", "99.1.2-beta.5", expected})
	if version != expected {
		t.Errorf("GetLatestPatchVersion returned the wrong latest version, expected %s", expected)
	}
}

func Test_RationalizeReleaseAndVersion(t *testing.T) {
	version := RationalizeReleaseAndVersion(Kubernetes, "", "", false)
	if version != KubernetesDefaultVersion {
		t.Errorf("It is not the default Kubernetes version")
	}

	version = RationalizeReleaseAndVersion(Kubernetes, "1.6", "", false)
	if version != KubernetesVersion1Dot6Dot13 {
		t.Errorf("It is not Kubernetes version %s", KubernetesVersion1Dot6Dot13)
	}

	version = RationalizeReleaseAndVersion(Kubernetes, "", "1.6.11", false)
	if version != KubernetesVersion1Dot6Dot11 {
		t.Errorf("It is not Kubernetes version %s", KubernetesVersion1Dot6Dot11)
	}

	version = RationalizeReleaseAndVersion(Kubernetes, "1.6", "1.6.11", false)
	if version != KubernetesVersion1Dot6Dot11 {
		t.Errorf("It is not Kubernetes version %s", KubernetesVersion1Dot6Dot11)
	}

	version = RationalizeReleaseAndVersion(Kubernetes, "1.7", "", false)
	if version != KubernetesVersion1Dot7Dot14 {
		t.Errorf("It is not Kubernetes version %s", KubernetesVersion1Dot7Dot14)
	}

	version = RationalizeReleaseAndVersion(Kubernetes, "", "1.7.14", false)
	if version != KubernetesVersion1Dot7Dot14 {
		t.Errorf("It is not Kubernetes version %s", KubernetesVersion1Dot7Dot14)
	}

	version = RationalizeReleaseAndVersion(Kubernetes, "1.7", "1.7.14", false)
	if version != KubernetesVersion1Dot7Dot14 {
		t.Errorf("It is not Kubernetes version %s", KubernetesVersion1Dot7Dot14)
	}

	version = RationalizeReleaseAndVersion(Kubernetes, "", "1.6.7", false)
	if version != "" {
		t.Errorf("It is not empty string")
	}

	version = RationalizeReleaseAndVersion(Kubernetes, "1.1", "", false)
	if version != "" {
		t.Errorf("It is not empty string")
	}

	version = RationalizeReleaseAndVersion(Kubernetes, "1.1", "1.6.6", false)
	if version != "" {
		t.Errorf("It is not empty string")
	}

	version = RationalizeReleaseAndVersion(Kubernetes, "", "v1.8.8", false)
	if version != KubernetesVersion1Dot8Dot8 {
		t.Errorf("It is not Kubernetes version %s", KubernetesVersion1Dot8Dot8)
	}

	version = RationalizeReleaseAndVersion(Kubernetes, "v1.9", "", false)
	if version != KubernetesVersion1Dot9Dot4 {
		t.Errorf("It is not Kubernetes version %s", KubernetesVersion1Dot9Dot4)
	}

	version = RationalizeReleaseAndVersion(Kubernetes, "1.10", "", false)
	if version != KubernetesVersion1Dot10Dot0 {
		t.Errorf("It is not Kubernetes version %s", KubernetesVersion1Dot10Dot0)
	}
}

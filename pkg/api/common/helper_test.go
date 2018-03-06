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
	if version != KubernetesVersion1Dot7Dot13 {
		t.Errorf("It is not Kubernetes version %s", KubernetesVersion1Dot7Dot13)
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
}

func Test_RationalizeReleaseAndVersion(t *testing.T) {
	version := RationalizeReleaseAndVersion(Kubernetes, "", "")
	if version != KubernetesDefaultVersion {
		t.Errorf("It is not the default Kubernetes version")
	}

	version = RationalizeReleaseAndVersion(Kubernetes, "1.6", "")
	if version != KubernetesVersion1Dot6Dot13 {
		t.Errorf("It is not Kubernetes version %s", KubernetesVersion1Dot6Dot13)
	}

	version = RationalizeReleaseAndVersion(Kubernetes, "1.9", "")
	if version != KubernetesVersion1Dot9Dot3 {
		t.Errorf("It is not Kubernetes version %s", KubernetesVersion1Dot9Dot3)
	}

	version = RationalizeReleaseAndVersion(Kubernetes, "", "1.6.11")
	if version != KubernetesVersion1Dot6Dot11 {
		t.Errorf("It is not Kubernetes version %s", KubernetesVersion1Dot6Dot11)
	}

	version = RationalizeReleaseAndVersion(Kubernetes, "1.6", "1.6.11")
	if version != KubernetesVersion1Dot6Dot11 {
		t.Errorf("It is not Kubernetes version %s", KubernetesVersion1Dot6Dot11)
	}

	version = RationalizeReleaseAndVersion(Kubernetes, "", "1.6.7")
	if version != "" {
		t.Errorf("It is not empty string")
	}

	version = RationalizeReleaseAndVersion(Kubernetes, "1.1", "")
	if version != "" {
		t.Errorf("It is not empty string")
	}

	version = RationalizeReleaseAndVersion(Kubernetes, "1.1", "1.6.6")
	if version != "" {
		t.Errorf("It is not empty string")
	}

	version = RationalizeReleaseAndVersion(Kubernetes, "", "v1.8.8")
	if version != KubernetesVersion1Dot8Dot8 {
		t.Errorf("It is not Kubernetes version %s", KubernetesVersion1Dot8Dot8)
	}

	version = RationalizeReleaseAndVersion(Kubernetes, "v1.9", "")
	if version != KubernetesVersion1Dot9Dot3 {
		t.Errorf("It is not Kubernetes version %s", KubernetesVersion1Dot9Dot3)
	}
}

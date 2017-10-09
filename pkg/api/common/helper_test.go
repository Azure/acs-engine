package common

import "testing"

func Test_RationalizeReleaseAndVersion(t *testing.T) {
	version := ""
	version = RationalizeReleaseAndVersion(Kubernetes, "", "")
	if version != KubernetesDefaultVersion {
		t.Errorf("It is not the default Kubernetes version")
	}

	version = RationalizeReleaseAndVersion(Kubernetes, "1.6", "")
	if version != KubernetesVersion1Dot6Dot11 {
		t.Errorf("It is not Kubernetes version %s", KubernetesVersion1Dot6Dot11)
	}

	version = RationalizeReleaseAndVersion(Kubernetes, "", "1.6.11")
	if version != KubernetesVersion1Dot6Dot11 {
		t.Errorf("It is not Kubernetes version %s", KubernetesVersion1Dot6Dot11)
	}

	version = RationalizeReleaseAndVersion(Kubernetes, "1.6", "1.6.11")
	if version != KubernetesVersion1Dot6Dot11 {
		t.Errorf("It is not Kubernetes version %s", KubernetesVersion1Dot6Dot11)
	}

	version = RationalizeReleaseAndVersion(Kubernetes, "", "1.6.6")
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

	version = RationalizeReleaseAndVersion(Kubernetes, "", "1.5.8")
	if version != "1.5.8" {
		t.Errorf("It is not 1.5.8")
	}

}

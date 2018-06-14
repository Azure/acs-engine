package api

import (
	"testing"

	"github.com/Azure/acs-engine/pkg/api/common"
	"github.com/blang/semver"
	. "github.com/onsi/gomega"
)

func TestInvalidVersion(t *testing.T) {
	RegisterTestingT(t)

	invalid := []string{
		"invalid number",
		"invalid.number",
		"a4.b7.c3",
		"31.29.",
		".17.02",
		"43.156.89.",
		"1.2.a"}

	for _, v := range invalid {
		_, e := semver.Make(v)
		Expect(e).NotTo(BeNil())
	}
}

func TestVersionCompare(t *testing.T) {
	RegisterTestingT(t)

	type record struct {
		v1, v2    string
		isGreater bool
	}
	records := []record{
		{"37.48.59", "37.48.59", false},
		{"17.4.5", "3.1.1", true},
		{"9.6.5", "9.45.5", false},
		{"2.3.8", "2.3.24", false}}

	for _, r := range records {
		ver, e := semver.Make(r.v1)
		Expect(e).To(BeNil())
		constraint, e := semver.Make(r.v2)
		Expect(e).To(BeNil())
		Expect(r.isGreater).To(Equal(ver.GT(constraint)))
	}
}

func TestOrchestratorUpgradeInfo(t *testing.T) {
	RegisterTestingT(t)
	// 1.6.9 is upgradable to 1.6.x and 1.7.x
	deployedVersion := "1.6.9"
	nextNextMinorVersion := "1.8.0"
	csOrch := &OrchestratorProfile{
		OrchestratorType:    Kubernetes,
		OrchestratorVersion: deployedVersion,
	}
	v := common.GetVersionsBetween(common.GetAllSupportedKubernetesVersions(), deployedVersion, nextNextMinorVersion, false, true)
	orch, e := GetOrchestratorVersionProfile(csOrch)
	Expect(e).To(BeNil())
	Expect(len(orch.Upgrades)).To(Equal(len(v)))

	// 1.7.0 is upgradable to 1.7.x and 1.8.x
	deployedVersion = "1.7.0"
	nextNextMinorVersion = "1.9.0"
	csOrch = &OrchestratorProfile{
		OrchestratorType:    Kubernetes,
		OrchestratorVersion: deployedVersion,
	}
	v = common.GetVersionsBetween(common.GetAllSupportedKubernetesVersions(), deployedVersion, nextNextMinorVersion, false, true)
	orch, e = GetOrchestratorVersionProfile(csOrch)
	Expect(e).To(BeNil())
	Expect(len(orch.Upgrades)).To(Equal(len(v)))

	// 1.7.15 is upgradable to 1.8.x
	deployedVersion = "1.7.15"
	nextNextMinorVersion = "1.9.0"
	csOrch = &OrchestratorProfile{
		OrchestratorType:    Kubernetes,
		OrchestratorVersion: deployedVersion,
	}
	v = common.GetVersionsBetween(common.GetAllSupportedKubernetesVersions(), deployedVersion, nextNextMinorVersion, false, true)
	orch, e = GetOrchestratorVersionProfile(csOrch)
	Expect(e).To(BeNil())
	Expect(len(orch.Upgrades)).To(Equal(len(v)))

	// 1.8.4 is upgradable to 1.8.x and 1.9.x
	deployedVersion = "1.8.4"
	nextNextMinorVersion = "1.10.0"
	csOrch = &OrchestratorProfile{
		OrchestratorType:    Kubernetes,
		OrchestratorVersion: deployedVersion,
	}
	v = common.GetVersionsBetween(common.GetAllSupportedKubernetesVersions(), deployedVersion, nextNextMinorVersion, false, true)
	orch, e = GetOrchestratorVersionProfile(csOrch)
	Expect(e).To(BeNil())
	Expect(len(orch.Upgrades)).To(Equal(len(v)))

	// 1.9.6 is upgradable to 1.10.x
	deployedVersion = "1.9.6"
	nextNextMinorVersion = "1.11.0"
	csOrch = &OrchestratorProfile{
		OrchestratorType:    Kubernetes,
		OrchestratorVersion: deployedVersion,
	}
	v = common.GetVersionsBetween(common.GetAllSupportedKubernetesVersions(), deployedVersion, nextNextMinorVersion, false, true)
	orch, e = GetOrchestratorVersionProfile(csOrch)
	Expect(e).To(BeNil())
	Expect(len(orch.Upgrades)).To(Equal(len(v)))

	// 1.10.0-beta.2 is upgradable to newer pre-release versions in 1.10.n release channel and official 1.10.n releases
	deployedVersion = "1.10.0-beta.2"
	nextNextMinorVersion = "1.12.0"
	csOrch = &OrchestratorProfile{
		OrchestratorType:    Kubernetes,
		OrchestratorVersion: deployedVersion,
	}
	v = common.GetVersionsBetween(common.GetAllSupportedKubernetesVersions(), deployedVersion, nextNextMinorVersion, false, true)
	orch, e = GetOrchestratorVersionProfile(csOrch)
	Expect(e).To(BeNil())
	Expect(len(orch.Upgrades)).To(Equal(len(v)))

	// The latest version is not upgradable
	csOrch = &OrchestratorProfile{
		OrchestratorType:    Kubernetes,
		OrchestratorVersion: common.GetMaxVersion(common.GetAllSupportedKubernetesVersions(), true),
	}
	orch, e = GetOrchestratorVersionProfile(csOrch)
	Expect(e).To(BeNil())
	Expect(len(orch.Upgrades)).To(Equal(0))
}

func TestGetOrchestratorVersionProfileListV20170930(t *testing.T) {
	RegisterTestingT(t)
	// v20170930 - all orchestrators
	list, e := GetOrchestratorVersionProfileListV20170930("", "")
	Expect(e).To(BeNil())
	numSwarmVersions := 1
	numDockerCEVersions := 1

	totalNumVersions := numSwarmVersions +
		numDockerCEVersions +
		len(common.GetAllSupportedKubernetesVersions()) +
		len(common.AllDCOSSupportedVersions) +
		len(common.GetAllSupportedOpenShiftVersions()) - 1

	Expect(len(list.Properties.Orchestrators)).To(Equal(totalNumVersions))

	// v20170930 - kubernetes only
	list, e = GetOrchestratorVersionProfileListV20170930(common.Kubernetes, "")
	Expect(e).To(BeNil())
	Expect(len(list.Properties.Orchestrators)).To(Equal(len(common.GetAllSupportedKubernetesVersions())))
}

func TestKubernetesInfo(t *testing.T) {
	RegisterTestingT(t)

	invalid := []string{
		"invalid number",
		"invalid.number",
		"a4.b7.c3",
		"31.29.",
		".17.02",
		"43.156.89.",
		"1.2.a",
		"1.5.9",
		"1.6.8"}

	for _, v := range invalid {
		csOrch := &OrchestratorProfile{
			OrchestratorType:    Kubernetes,
			OrchestratorVersion: v,
		}

		_, e := kubernetesInfo(csOrch)
		Expect(e).NotTo(BeNil())
	}

}

func TestOpenshiftInfo(t *testing.T) {
	RegisterTestingT(t)

	invalid := []string{
		"invalid number",
		"invalid.number",
		"a4.b7.c3",
		"31.29.",
		".17.02",
		"43.156.89.",
		"1.2.a",
		"3.8.9",
		"3.9.2"}

	for _, v := range invalid {
		csOrch := &OrchestratorProfile{
			OrchestratorType:    OpenShift,
			OrchestratorVersion: v,
		}

		_, e := openShiftInfo(csOrch)
		Expect(e).NotTo(BeNil())
	}

	// test good value
	csOrch := &OrchestratorProfile{
		OrchestratorType:    OpenShift,
		OrchestratorVersion: common.OpenShiftDefaultVersion,
	}

	_, e := openShiftInfo(csOrch)
	Expect(e).To(BeNil())
}

func TestDcosInfo(t *testing.T) {
	RegisterTestingT(t)
	invalid := []string{
		"invalid number",
		"invalid.number",
		"a4.b7.c3",
		"31.29.",
		".17.02",
		"43.156.89.",
		"1.2.a"}

	for _, v := range invalid {
		csOrch := &OrchestratorProfile{
			OrchestratorType:    DCOS,
			OrchestratorVersion: v,
		}

		_, e := dcosInfo(csOrch)
		Expect(e).NotTo(BeNil())
	}

	// test good value
	csOrch := &OrchestratorProfile{
		OrchestratorType:    DCOS,
		OrchestratorVersion: common.DCOSDefaultVersion,
	}

	_, e := dcosInfo(csOrch)
	Expect(e).To(BeNil())
}

func TestSwarmInfo(t *testing.T) {
	RegisterTestingT(t)
	invalid := []string{
		"swarm:1.1.1",
		"swarm:1.1.2",
	}

	for _, v := range invalid {
		csOrch := &OrchestratorProfile{
			OrchestratorType:    Swarm,
			OrchestratorVersion: v,
		}

		_, e := swarmInfo(csOrch)
		Expect(e).NotTo(BeNil())
	}

	// test good value
	csOrch := &OrchestratorProfile{
		OrchestratorType:    Swarm,
		OrchestratorVersion: common.SwarmVersion,
	}

	_, e := swarmInfo(csOrch)
	Expect(e).To(BeNil())
}

func TestDockerceInfoInfo(t *testing.T) {
	RegisterTestingT(t)
	invalid := []string{
		"17.02.1",
		"43.156.89",
	}

	for _, v := range invalid {
		csOrch := &OrchestratorProfile{
			OrchestratorType:    SwarmMode,
			OrchestratorVersion: v,
		}

		_, e := dockerceInfo(csOrch)
		Expect(e).NotTo(BeNil())
	}

	// test good value
	csOrch := &OrchestratorProfile{
		OrchestratorType:    SwarmMode,
		OrchestratorVersion: common.DockerCEVersion,
	}

	_, e := dockerceInfo(csOrch)
	Expect(e).To(BeNil())
}

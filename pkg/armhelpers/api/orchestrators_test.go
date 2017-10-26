package api

import (
	"testing"

	"github.com/Azure/acs-engine/pkg/api/common"
	"github.com/Masterminds/semver"
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
		_, e := semver.NewVersion(v)
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
		ver, e := semver.NewVersion(r.v1)
		Expect(e).To(BeNil())
		constraint, e := semver.NewConstraint(">" + r.v2)
		Expect(e).To(BeNil())
		Expect(r.isGreater).To(Equal(constraint.Check(ver)))
	}
}

func TestOrchestratorUpgradeInfo(t *testing.T) {
	RegisterTestingT(t)
	// 1.5.3 is upgradable to 1.6.x
	csOrch := &OrchestratorProfile{
		OrchestratorType:    Kubernetes,
		OrchestratorVersion: "1.5.3",
	}
	orch, e := GetOrchestratorVersionProfile(csOrch)
	Expect(e).To(BeNil())
	Expect(len(orch.Upgrades)).To(Equal(1))
	Expect(orch.Upgrades[0].OrchestratorVersion).To(Equal(common.KubernetesVersion1Dot6Dot11))

	// 1.6.8 is upgradable to 1.6.x and 1.7.x
	csOrch = &OrchestratorProfile{
		OrchestratorType:    Kubernetes,
		OrchestratorVersion: "1.6.8",
	}
	orch, e = GetOrchestratorVersionProfile(csOrch)
	Expect(e).To(BeNil())
	Expect(len(orch.Upgrades)).To(Equal(2))
	Expect(orch.Upgrades[0].OrchestratorVersion).To(Equal(common.KubernetesVersion1Dot6Dot11))
	Expect(orch.Upgrades[1].OrchestratorVersion).To(Equal(common.KubernetesVersion1Dot7Dot7))

	// 1.7.0 is upgradable to 1.7.x and 1.8.x
	csOrch = &OrchestratorProfile{
		OrchestratorType:    Kubernetes,
		OrchestratorVersion: "1.7.0",
	}
	orch, e = GetOrchestratorVersionProfile(csOrch)
	Expect(e).To(BeNil())
	Expect(len(orch.Upgrades)).To(Equal(2))
	Expect(orch.Upgrades[0].OrchestratorVersion).To(Equal(common.KubernetesVersion1Dot7Dot7))
	Expect(orch.Upgrades[1].OrchestratorVersion).To(Equal(common.KubernetesVersion1Dot8Dot2))

	// 1.7.9 is upgradable to 1.8.x
	csOrch = &OrchestratorProfile{
		OrchestratorType:    Kubernetes,
		OrchestratorVersion: "1.7.9",
	}
	orch, e = GetOrchestratorVersionProfile(csOrch)
	Expect(e).To(BeNil())
	Expect(len(orch.Upgrades)).To(Equal(1))
	Expect(orch.Upgrades[0].OrchestratorVersion).To(Equal(common.KubernetesVersion1Dot8Dot2))

	// 1.8.2 is not upgradable
	csOrch = &OrchestratorProfile{
		OrchestratorType:    Kubernetes,
		OrchestratorVersion: KubernetesVersion1Dot8Dot2,
	}
	orch, e = GetOrchestratorVersionProfile(csOrch)
	Expect(e).To(BeNil())
	Expect(len(orch.Upgrades)).To(Equal(0))

	// v20170930
	list, e := GetOrchestratorVersionProfileListV20170930("", "")
	Expect(e).To(BeNil())
	Expect(len(list.Properties.Orchestrators)).NotTo(Equal(0))
}

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
	// 1.6.8 is upgradable to 1.6.x and 1.7.x
	csOrch := &OrchestratorProfile{
		OrchestratorType:    Kubernetes,
		OrchestratorVersion: "1.6.8",
	}
	orch, e := GetOrchestratorVersionProfile(csOrch)
	Expect(e).To(BeNil())
	// 1.6.9, 1.6.11, 1.6.12, 1.6.13, 1.7.0, 1.7.1, 1.7.2, 1.7.4, 1.7.5, 1.7.7, 1.7.9, 1.7.10, 1.7.12, 1.7.13, 1.7.14
	Expect(len(orch.Upgrades)).To(Equal(15))

	// 1.7.0 is upgradable to 1.7.x and 1.8.x
	csOrch = &OrchestratorProfile{
		OrchestratorType:    Kubernetes,
		OrchestratorVersion: "1.7.0",
	}
	orch, e = GetOrchestratorVersionProfile(csOrch)
	Expect(e).To(BeNil())
	// 1.7.1, 1.7.2, 1.7.4, 1.7.5, 1.7.7, 1.7.9, 1.7.10, 1.7.12, 1.7.13, 1.7.14, 1.8.0, 1.8.1, 1.8.2, 1.8.4, 1.8.6, 1.8.7, 1.8.8, 1.8.9
	Expect(len(orch.Upgrades)).To(Equal(18))

	// 1.7.14 is upgradable to 1.8.x
	csOrch = &OrchestratorProfile{
		OrchestratorType:    Kubernetes,
		OrchestratorVersion: "1.7.14",
	}
	orch, e = GetOrchestratorVersionProfile(csOrch)
	Expect(e).To(BeNil())
	// 1.8.0, 1.8.1, 1.8.2, 1.8.4, 1.8.6, 1.8.7, 1.8.8, 1.8.9
	Expect(len(orch.Upgrades)).To(Equal(8))

	// 1.8.4 is upgradable to 1.8.x and 1.9.x
	csOrch = &OrchestratorProfile{
		OrchestratorType:    Kubernetes,
		OrchestratorVersion: "1.8.4",
	}
	orch, e = GetOrchestratorVersionProfile(csOrch)
	Expect(e).To(BeNil())
	// 1.8.6, 1.8.7, 1.8.8, 1.8.9, 1.9.0, 1.9.1, 1.9.2, 1.9.3, 1.9.4
	Expect(len(orch.Upgrades)).To(Equal(9))

	// 1.9.4 is upgradable to 1.10.x
	csOrch = &OrchestratorProfile{
		OrchestratorType:    Kubernetes,
		OrchestratorVersion: "1.9.4",
	}
	orch, e = GetOrchestratorVersionProfile(csOrch)
	Expect(e).To(BeNil())
	// 1.10.0-beta.2
	Expect(len(orch.Upgrades)).To(Equal(1))

	// 1.10.x is not upgradable
	csOrch = &OrchestratorProfile{
		OrchestratorType:    Kubernetes,
		OrchestratorVersion: "1.10.0-beta.2",
	}
	orch, e = GetOrchestratorVersionProfile(csOrch)
	Expect(e).To(BeNil())
	Expect(len(orch.Upgrades)).To(Equal(0))

	// v20170930 - all orchestrators
	list, e := GetOrchestratorVersionProfileListV20170930("", "")
	Expect(e).To(BeNil())
	Expect(len(list.Properties.Orchestrators)).To(Equal(35))

	// v20170930 - kubernetes only
	list, e = GetOrchestratorVersionProfileListV20170930(common.Kubernetes, "")
	Expect(e).To(BeNil())
	Expect(len(list.Properties.Orchestrators)).To(Equal(30))
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
		"1.2.a"}

	for _, v := range invalid {
		csOrch := &OrchestratorProfile{
			OrchestratorType:    Kubernetes,
			OrchestratorVersion: v,
		}

		_, e := kubernetesInfo(csOrch)
		Expect(e).NotTo(BeNil())
	}

}

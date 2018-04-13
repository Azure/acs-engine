package v20180331

import (
	"testing"

	"github.com/Azure/acs-engine/pkg/helpers"
)

func TestMerge_DNSPrefix(t *testing.T) {
	newMC := &ManagedCluster{
		Properties: &Properties{
			DNSPrefix: "newprefix",
		},
	}

	existingMC := &ManagedCluster{
		Properties: &Properties{
			DNSPrefix:  "oldprefix",
			EnableRBAC: helpers.PointerToBool(false),
		},
	}

	e := newMC.Merge(existingMC)
	if e == nil {
		t.Error("expect error to not be nil")
	}

	newMC = &ManagedCluster{
		Properties: &Properties{},
	}

	existingMC = &ManagedCluster{
		Properties: &Properties{
			DNSPrefix:  "oldprefix",
			EnableRBAC: helpers.PointerToBool(false),
		},
	}

	e = newMC.Merge(existingMC)
	if e != nil {
		t.Error("expect error to be nil")
	}

	if newMC.Properties.DNSPrefix != "oldprefix" {
		t.Error("expect dnsPrefix to be oldprefix when update with empty input")
	}

	newMC = &ManagedCluster{
		Properties: &Properties{},
	}

	existingMC = &ManagedCluster{
		Properties: &Properties{
			DNSPrefix:  "",
			EnableRBAC: helpers.PointerToBool(false),
		},
	}

	e = newMC.Merge(existingMC)
	if e == nil {
		t.Error("expect error to not be nil")
	}
}

func TestMerge_EnableRBAC(t *testing.T) {
	newMC := &ManagedCluster{
		Properties: &Properties{
			EnableRBAC: nil,
		},
	}

	existingMC := &ManagedCluster{
		Properties: &Properties{
			DNSPrefix:  "something",
			EnableRBAC: helpers.PointerToBool(false),
		},
	}

	e := newMC.Merge(existingMC)
	if e != nil {
		t.Error("expect error to be nil")
	}
	if newMC.Properties.EnableRBAC == nil || *newMC.Properties.EnableRBAC != false {
		t.Error("expect EnableRBAC to be same with existing when omit in updating")
	}

	newMC = &ManagedCluster{
		Properties: &Properties{
			EnableRBAC: nil,
		},
	}

	existingMC = &ManagedCluster{
		Properties: &Properties{
			DNSPrefix:  "something",
			EnableRBAC: helpers.PointerToBool(true),
		},
	}

	e = newMC.Merge(existingMC)
	if e != nil {
		t.Error("expect error to be nil")
	}
	if newMC.Properties.EnableRBAC == nil || *newMC.Properties.EnableRBAC != true {
		t.Error("expect EnableRBAC to be same with existing when omit in updating")
	}

	newMC = &ManagedCluster{
		Properties: &Properties{
			EnableRBAC: nil,
		},
	}

	existingMC = &ManagedCluster{
		Properties: &Properties{
			DNSPrefix:  "something",
			EnableRBAC: nil,
		},
	}

	e = newMC.Merge(existingMC)
	if e == nil {
		t.Error("expect error not to be nil")
	}

	newMC = &ManagedCluster{
		Properties: &Properties{
			EnableRBAC: helpers.PointerToBool(true),
		},
	}

	existingMC = &ManagedCluster{
		Properties: &Properties{
			DNSPrefix:  "something",
			EnableRBAC: helpers.PointerToBool(true),
		},
	}

	e = newMC.Merge(existingMC)
	if e != nil {
		t.Error("expect error to be nil")
	}
	if newMC.Properties.EnableRBAC == nil || *newMC.Properties.EnableRBAC != true {
		t.Error("expect EnableRBAC to be true")
	}

	newMC = &ManagedCluster{
		Properties: &Properties{
			EnableRBAC: helpers.PointerToBool(false),
		},
	}

	existingMC = &ManagedCluster{
		Properties: &Properties{
			DNSPrefix:  "something",
			EnableRBAC: helpers.PointerToBool(true),
		},
	}

	e = newMC.Merge(existingMC)
	if e == nil {
		t.Error("expect error to be nil")
	}

}

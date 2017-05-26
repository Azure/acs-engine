package operations

import (
	"testing"

	"github.com/Azure/acs-engine/pkg/armhelpers"
	log "github.com/Sirupsen/logrus"
	. "github.com/onsi/gomega"
)

func Test_ScaleDownVMs_ErrorPath(t *testing.T) {
	RegisterTestingT(t)
	errs := ScaleDownVMs(&armhelpers.FailingMockClient{}, log.NewEntry(log.New()), "rg", "vm1", "vm2", "vm3", "vm5")
	Expect(errs.Len()).To(Equal(4))
	for e := errs.Front(); e != nil; e = e.Next() {
		output := e.Value.(*VMScalingErrorDetails)
		Expect(output.Name).To(ContainSubstring("vm"))
		Expect(output.Error).To(Not(BeNil()))
	}
}

func Test_ScaleDownVMs_HappyPath(t *testing.T) {
	RegisterTestingT(t)
	errs := ScaleDownVMs(&armhelpers.MockClient{}, log.NewEntry(log.New()), "rg", "k8s-agent-F8EADCCF-0", "k8s-agent-F8EADCCF-3", "k8s-agent-F8EADCCF-2", "k8s-agent-F8EADCCF-4")
	Expect(errs).To(BeNil())
}

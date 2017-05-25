package operations

import (
	"testing"

	"github.com/Azure/acs-engine/pkg/api"
	. "github.com/onsi/gomega"
)

func Test_getVMName(t *testing.T) {
	RegisterTestingT(t)
	suffix := "F8EADCCF"
	agentIndex := 12
	shortPoolIndex := 3
	longPoolIndex := 45

	result, err := getVMName(api.Kubernetes, api.Linux, shortPoolIndex, agentIndex, suffix)
	Expect(err).To(BeNil())
	Expect(result).To(Equal("k8s-agent3-F8EADCCF-12"))

	result, err = getVMName(api.Kubernetes, api.Windows, longPoolIndex, agentIndex, suffix)
	Expect(err).To(BeNil())
	Expect(result).To(Equal("F8EADacs94512"))

	result, err = getVMName(api.Kubernetes, api.Windows, shortPoolIndex, agentIndex, suffix)
	Expect(err).To(BeNil())
	Expect(result).To(Equal("F8EADacs90312"))

	result, err = getVMName(api.DCOS, api.Linux, shortPoolIndex, agentIndex, suffix)
	Expect(err).To(BeNil())
	Expect(result).To(Equal("dcos-agent3-F8EADCCF-12"))

	result, err = getVMName(api.DCOS, api.Windows, shortPoolIndex, agentIndex, suffix)
	Expect(err).To(Not(BeNil()))
	Expect(result).To(Equal(""))

	result, err = getVMName(api.Swarm, api.Linux, shortPoolIndex, agentIndex, suffix)
	Expect(err).To(BeNil())
	Expect(result).To(Equal("swarm-agent3-F8EADCCF-12"))

	result, err = getVMName(api.Swarm, api.Windows, longPoolIndex, agentIndex, suffix)
	Expect(err).To(BeNil())
	Expect(result).To(Equal("F8EADacs94512"))

	result, err = getVMName(api.Swarm, api.Windows, shortPoolIndex, agentIndex, suffix)
	Expect(err).To(BeNil())
	Expect(result).To(Equal("F8EADacs90312"))

	result, err = getVMName(api.SwarmMode, api.Linux, shortPoolIndex, agentIndex, suffix)
	Expect(err).To(BeNil())
	Expect(result).To(Equal("swarmm-agent3-F8EADCCF-12"))

	result, err = getVMName(api.SwarmMode, api.Windows, longPoolIndex, agentIndex, suffix)
	Expect(err).To(BeNil())
	Expect(result).To(Equal("F8EADacs94512"))

	result, err = getVMName(api.SwarmMode, api.Windows, shortPoolIndex, agentIndex, suffix)
	Expect(err).To(BeNil())
	Expect(result).To(Equal("F8EADacs90312"))
}

func Test_ScaleDownVMASAgentPool_ErrorPath(t *testing.T) {
	RegisterTestingT(t)
	errs := ScaleDownVMASAgentPool(&failingMockClient{}, "F8EADCCF", "rg", 1, api.Kubernetes, api.Linux,
		2, 3, 4, 6)
	Expect(errs.Len()).To(Equal(4))
	for e := errs.Front(); e != nil; e = e.Next() {
		output := e.Value.(*VMScalingErrorDetails)
		Expect(output.Index).To(Not(Equal(0)))
		Expect(output.Error).To(Not(BeNil()))
	}
}

func Test_ScaleDownVMASAgentPool_HappyPath(t *testing.T) {
	RegisterTestingT(t)
	errs := ScaleDownVMASAgentPool(&mockClient{}, "F8EADCCF", "rg", 1, api.Kubernetes, api.Linux,
		2, 3, 4, 6)
	Expect(errs).To(BeNil())
}

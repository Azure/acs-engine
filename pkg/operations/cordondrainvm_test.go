package operations

import (
	"github.com/Azure/acs-engine/pkg/armhelpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/pkg/api/v1"
)

var _ = Describe("Safely Drain node operation tests", func() {
	It("Should return error messages for failure to create kubernetes client", func() {
		err := SafelyDrainNode(&armhelpers.MockACSEngineClient{FailGetKubernetesClient: true}, log.NewEntry(log.New()), "http://bad.com/", "bad", "node")
		Expect(err).Should(HaveOccurred())
	})
	It("Should return error messages for Failure to get node ", func() {
		mockClient := &armhelpers.MockACSEngineClient{MockKubernetesClient: &armhelpers.MockKubernetesClient{}}
		mockClient.MockKubernetesClient.FailGetNode = true
		err := SafelyDrainNode(mockClient, log.NewEntry(log.New()), "http://bad.com/", "bad", "node")
		Expect(err).Should(HaveOccurred())
	})
	It("Should return error messages for Failure to update node ", func() {
		mockClient := &armhelpers.MockACSEngineClient{MockKubernetesClient: &armhelpers.MockKubernetesClient{}}
		mockClient.MockKubernetesClient.FailUpdateNode = true
		err := SafelyDrainNode(mockClient, log.NewEntry(log.New()), "http://bad.com/", "bad", "node")
		Expect(err).Should(HaveOccurred())
	})
	It("Should return error messages for Failure to list pods ", func() {
		mockClient := &armhelpers.MockACSEngineClient{MockKubernetesClient: &armhelpers.MockKubernetesClient{}}
		mockClient.MockKubernetesClient.FailListPods = true
		err := SafelyDrainNode(mockClient, log.NewEntry(log.New()), "http://bad.com/", "bad", "node")
		Expect(err).Should(HaveOccurred())
	})
	It("Should return error messages for Failure to check support eviction ", func() {
		mockClient := &armhelpers.MockACSEngineClient{MockKubernetesClient: &armhelpers.MockKubernetesClient{}}
		mockClient.MockKubernetesClient.PodsList = &v1.PodList{Items: []v1.Pod{{}}}
		mockClient.MockKubernetesClient.FailSupportEviction = true
		err := SafelyDrainNode(mockClient, log.NewEntry(log.New()), "http://bad.com/", "bad", "node")
		Expect(err).Should(HaveOccurred())
	})
	It("Should return error messages for Failure to delete pod ", func() {
		mockClient := &armhelpers.MockACSEngineClient{MockKubernetesClient: &armhelpers.MockKubernetesClient{}}
		mockClient.MockKubernetesClient.PodsList = &v1.PodList{Items: []v1.Pod{{}}}
		mockClient.MockKubernetesClient.FailDeletePod = true
		err := SafelyDrainNode(mockClient, log.NewEntry(log.New()), "http://bad.com/", "bad", "node")
		Expect(err).Should(HaveOccurred())
	})
	It("Should return error messages for Failure to Evict Pod ", func() {
		mockClient := &armhelpers.MockACSEngineClient{MockKubernetesClient: &armhelpers.MockKubernetesClient{}}
		mockClient.MockKubernetesClient.PodsList = &v1.PodList{Items: []v1.Pod{{}}}
		mockClient.MockKubernetesClient.ShouldSupportEviction = true
		mockClient.MockKubernetesClient.FailEvictPod = true
		err := SafelyDrainNode(mockClient, log.NewEntry(log.New()), "http://bad.com/", "bad", "node")
		Expect(err).Should(HaveOccurred())
	})
	It("Should return error messages for Failure to wait for delete in delete path ", func() {
		mockClient := &armhelpers.MockACSEngineClient{MockKubernetesClient: &armhelpers.MockKubernetesClient{}}
		mockClient.MockKubernetesClient.PodsList = &v1.PodList{Items: []v1.Pod{{}}}
		mockClient.MockKubernetesClient.ShouldSupportEviction = true
		mockClient.MockKubernetesClient.FailWaitForDelete = true
		err := SafelyDrainNode(mockClient, log.NewEntry(log.New()), "http://bad.com/", "bad", "node")
		Expect(err).Should(HaveOccurred())
	})
	It("Should return error messages for Failure to wait for delete in eviction path ", func() {
		mockClient := &armhelpers.MockACSEngineClient{MockKubernetesClient: &armhelpers.MockKubernetesClient{}}
		mockClient.MockKubernetesClient.PodsList = &v1.PodList{Items: []v1.Pod{{}}}
		mockClient.MockKubernetesClient.ShouldSupportEviction = false
		mockClient.MockKubernetesClient.FailWaitForDelete = true
		err := SafelyDrainNode(mockClient, log.NewEntry(log.New()), "http://bad.com/", "bad", "node")
		Expect(err).Should(HaveOccurred())
	})
	It("Should not return error in valid eviction path ", func() {
		mockClient := &armhelpers.MockACSEngineClient{MockKubernetesClient: &armhelpers.MockKubernetesClient{}}
		mockClient.MockKubernetesClient.PodsList = &v1.PodList{Items: []v1.Pod{{}}}
		mockClient.MockKubernetesClient.ShouldSupportEviction = true
		err := SafelyDrainNode(mockClient, log.NewEntry(log.New()), "http://bad.com/", "bad", "node")
		Expect(err).ShouldNot(HaveOccurred())
	})
	It("Should not return error in valid delete path ", func() {
		mockClient := &armhelpers.MockACSEngineClient{MockKubernetesClient: &armhelpers.MockKubernetesClient{}}
		mockClient.MockKubernetesClient.PodsList = &v1.PodList{Items: []v1.Pod{{}}}
		mockClient.MockKubernetesClient.ShouldSupportEviction = false
		err := SafelyDrainNode(mockClient, log.NewEntry(log.New()), "http://bad.com/", "bad", "node")
		Expect(err).ShouldNot(HaveOccurred())
	})
})

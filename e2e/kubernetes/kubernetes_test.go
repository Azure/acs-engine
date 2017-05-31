package kubernetes

import (
	"fmt"
	"os/exec"

	"github.com/Azure/acs-engine/e2e"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Azure Container Cluster using the Kubernetes Orchestrator", func() {

	XIt("should be logged into the correct account", func() {
		acct, err := e2e.NewAccount()
		if err != nil {
			Fail("Unable to correctly build Account!")
		}
		acct.Login()
		acct.SetSubscription()

		current, err := e2e.GetCurrentAccount()

		Expect(err).NotTo(HaveOccurred())
		Expect(current.User.ID).To(Equal(acct.User.ID))
		Expect(current.TenantID).To(Equal(acct.TenantID))
		Expect(current.SubscriptionID).To(Equal(acct.SubscriptionID))
	})

	It("should have have the appropriate node count", func() {
		nodeList, err := GetNodes()
		Expect(err).NotTo(HaveOccurred())
		Expect(len(nodeList.Nodes)).To(Equal(4))
	})

	It("should be running the expected version", func() {
		version, err := GetVersion()
		Expect(err).NotTo(HaveOccurred())
		Expect(version).To(Equal("v1.5.3"))
	})

	It("should have kube-dns running", func() {
		running, err := AreAllPodsRunning("kube-dns", "kube-system")
		Expect(err).NotTo(HaveOccurred())
		Expect(running).To(Equal(true))
	})

	It("should have kube-dashboard running", func() {
		running, err := AreAllPodsRunning("kubernetes-dashboard", "kube-system")
		Expect(err).NotTo(HaveOccurred())
		Expect(running).To(Equal(true))
	})

	It("should have kube-proxy running", func() {
		running, err := AreAllPodsRunning("kube-proxy", "kube-system")
		Expect(err).NotTo(HaveOccurred())
		Expect(running).To(Equal(true))
	})

	XIt("should be able to access the dashboard from each node", func() {
		c, err := GetConfig()
		Expect(err).NotTo(HaveOccurred())

		sl, err := GetServices("kube-system")
		Expect(err).NotTo(HaveOccurred())
		port := sl.GetNodePortForPort("kubernetes-dashboard", 80)

		keyPath := "~/.ssh/id_rsa"
		master := fmt.Sprintf("azureuser@%s", c.GetServerName())
		nodeList, err := GetNodes()
		Expect(err).NotTo(HaveOccurred())

		for _, node := range nodeList.Nodes {
			dashboardURL := fmt.Sprintf("http://%s:%v", node.NodeStatus.GetAddressByType("InternalIP").Address, port)
			curlCMD := fmt.Sprintf("curl --max-time 60 %s", dashboardURL)
			_, err := exec.Command("ssh", "-i", keyPath, "-o", "ConnectTimeout=10", "-o", "StrictHostKeyChecking=no", "-o", "UserKnownHostsFile=/dev/null", master, curlCMD).Output()
			Expect(err).NotTo(HaveOccurred())
		}
	})

	// It("should be able to deploy an nginx service", func() {
	// 	_, err := exec.Command("kubectl", "create", "namespace", "deis").Output()
	// 	Expect(err).NotTo(HaveOccurred())
	// 	_, err = exec.Command("kubectl", "run", "--image", "library/nginx:latest", "nginx", "--namespace", "deis").Output()
	// 	Expect(err).NotTo(HaveOccurred())
	// })

})

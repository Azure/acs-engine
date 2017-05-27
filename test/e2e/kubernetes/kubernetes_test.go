package kubernetes

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os/exec"
	"time"

	"github.com/Azure/acs-engine/test/e2e/azure"
	"github.com/Azure/acs-engine/test/e2e/kubernetes/deployment"
	"github.com/Azure/acs-engine/test/e2e/kubernetes/namespace"
	"github.com/Azure/acs-engine/test/e2e/kubernetes/node"
	"github.com/Azure/acs-engine/test/e2e/kubernetes/pod"
	"github.com/Azure/acs-engine/test/e2e/kubernetes/service"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	r      = rand.New(rand.NewSource(time.Now().UnixNano()))
	suffix = r.Intn(99999)
	name   = fmt.Sprintf("test-%v", suffix)
	ns     *namespace.Namespace
	err    error
)

var _ = BeforeSuite(func() {
	ns, err = namespace.Create(name)
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	ns.Delete()
})

var _ = Describe("Azure Container Cluster using the Kubernetes Orchestrator", func() {

	It("should be logged into the correct account", func() {
		acct, err := azure.NewAccount()
		if err != nil {
			Fail("Unable to correctly build Account!")
		}
		acct.Login()
		acct.SetSubscription()

		current, err := azure.GetCurrentAccount()

		Expect(err).NotTo(HaveOccurred())
		Expect(current.User.ID).To(Equal(acct.User.ID))
		Expect(current.TenantID).To(Equal(acct.TenantID))
		Expect(current.SubscriptionID).To(Equal(acct.SubscriptionID))
	})

	It("should have have the appropriate node count", func() {
		nodeList, err := node.Get()
		Expect(err).NotTo(HaveOccurred())
		Expect(len(nodeList.Nodes)).To(Equal(4))
	})

	It("should be running the expected version", func() {
		version, err := node.Version()
		Expect(err).NotTo(HaveOccurred())
		Expect(version).To(Equal("v1.5.3"))
	})

	It("should have kube-dns running", func() {
		running, err := pod.AreAllPodsRunning("kube-dns", "kube-system")
		Expect(err).NotTo(HaveOccurred())
		Expect(running).To(Equal(true))
	})

	It("should have kube-dashboard running", func() {
		running, err := pod.AreAllPodsRunning("kubernetes-dashboard", "kube-system")
		Expect(err).NotTo(HaveOccurred())
		Expect(running).To(Equal(true))
	})

	It("should have kube-proxy running", func() {
		running, err := pod.AreAllPodsRunning("kube-proxy", "kube-system")
		Expect(err).NotTo(HaveOccurred())
		Expect(running).To(Equal(true))
	})

	It("should be able to access the dashboard from each node", func() {
		c, err := GetConfig()
		Expect(err).NotTo(HaveOccurred())

		s, err := service.Get("kubernetes-dashboard", "kube-system")
		Expect(err).NotTo(HaveOccurred())
		port := s.GetNodePort(80)

		keyPath := "~/.ssh/id_rsa"
		master := fmt.Sprintf("azureuser@%s", c.GetServerName())
		nodeList, err := node.Get()
		Expect(err).NotTo(HaveOccurred())

		for _, node := range nodeList.Nodes {
			dashboardURL := fmt.Sprintf("http://%s:%v", node.Status.GetAddressByType("InternalIP").Address, port)
			curlCMD := fmt.Sprintf("curl --max-time 60 %s", dashboardURL)
			_, err := exec.Command("ssh", "-i", keyPath, "-o", "ConnectTimeout=10", "-o", "StrictHostKeyChecking=no", "-o", "UserKnownHostsFile=/dev/null", master, curlCMD).CombinedOutput()
			Expect(err).NotTo(HaveOccurred())
		}
	})

	It("should be able to deploy an nginx service", func() {
		d, err := deployment.Create("library/nginx:latest", name, name)
		Expect(err).NotTo(HaveOccurred())
		err = d.Expose(80)
		Expect(err).NotTo(HaveOccurred())

		s, err := service.Get(name, name)
		Expect(err).NotTo(HaveOccurred())
		s, err = s.WaitForExternalIP(360, 5)
		Expect(err).NotTo(HaveOccurred())
		Expect(s.Status.LoadBalancer.Ingress).NotTo(BeEmpty())

		url := fmt.Sprintf("http://%s", s.Status.LoadBalancer.Ingress[0]["ip"])
		resp, err := http.Get(url)
		Expect(err).NotTo(HaveOccurred())
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		Expect(err).NotTo(HaveOccurred())
		Expect(string(body)).To(MatchRegexp("(Welcome to nginx)"))
	})

})

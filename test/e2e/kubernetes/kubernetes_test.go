package kubernetes

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/Azure/acs-engine/test/e2e/azure"
	"github.com/Azure/acs-engine/test/e2e/config"
	"github.com/Azure/acs-engine/test/e2e/engine"
	"github.com/Azure/acs-engine/test/e2e/kubernetes/deployment"
	"github.com/Azure/acs-engine/test/e2e/kubernetes/node"
	"github.com/Azure/acs-engine/test/e2e/kubernetes/pod"
	"github.com/Azure/acs-engine/test/e2e/kubernetes/service"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	cfg  config.Config
	err  error
	acct azure.Account
)

var _ = BeforeSuite(func() {
	c, err := config.ParseConfig()
	Expect(err).NotTo(HaveOccurred())
	cfg = *c // We have to do this because golang anon functions and scoping and stuff

	a, err := azure.NewAccount()
	Expect(err).NotTo(HaveOccurred())
	acct = *a // We have to do this because golang anon functions and scoping and stuff

	acct.Login()
	acct.SetSubscription()

	if cfg.Name == "" {
		cfg.Name = cfg.GenerateName()
		log.Printf("Cluster name:%s\n", cfg.Name)
		// Lets modify our template and call acs-engine generate on it
		e, err := engine.Build(cfg.ClusterDefinition, "_output", cfg.Name)
		Expect(err).NotTo(HaveOccurred())

		err = e.Generate()
		Expect(err).NotTo(HaveOccurred())

		err = acct.CreateGroup(cfg.Name, cfg.Location)
		Expect(err).NotTo(HaveOccurred())

		// Lets start by just using the normal az group deployment cli for creating a cluster
		log.Println("Creating deployment this make take a few minutes...")
		err = acct.CreateDeployment(cfg.Name, e)
		Expect(err).NotTo(HaveOccurred())
	}

	err = os.Setenv("KUBECONFIG", cfg.GetKubeConfig())
	Expect(err).NotTo(HaveOccurred())

	log.Println("Waiting on nodes to go into ready state...")
	ready := node.WaitOnReady(10*time.Second, 10*time.Minute)
	Expect(ready).To(BeTrue())
})

var _ = AfterSuite(func() {
	if cfg.CleanUpOnExit {
		log.Printf("Deleting Group:%s\n", cfg.Name)
		acct.DeleteGroup()
	}

	if _, err := os.Stat(cfg.GetKubeConfig()); os.IsExist(err) {
		if svc, _ := service.Get(cfg.Name, "default"); svc != nil {
			svc.Delete()
		}

		if d, _ := deployment.Get(cfg.Name, "default"); d != nil {
			d.Delete()
		}
	}
})

var _ = Describe("Azure Container Cluster using the Kubernetes Orchestrator", func() {

	It("should be logged into the correct account", func() {
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

	It("should be running the expected default version", func() {
		version, err := node.Version()
		Expect(err).NotTo(HaveOccurred())
		Expect(version).To(Equal("v1.6.6"))
	})

	It("should have kube-dns running", func() {
		pod.WaitOnReady("kube-dns", "kube-system", 5*time.Second, 3*time.Minute)
		running, err := pod.AreAllPodsRunning("kube-dns", "kube-system")
		Expect(err).NotTo(HaveOccurred())
		Expect(running).To(Equal(true))
	})

	It("should have kube-dashboard running", func() {
		pod.WaitOnReady("kubernetes-dashboard", "kube-system", 5*time.Second, 3*time.Minute)
		running, err := pod.AreAllPodsRunning("kubernetes-dashboard", "kube-system")
		Expect(err).NotTo(HaveOccurred())
		Expect(running).To(Equal(true))
	})

	It("should have kube-proxy running", func() {
		pod.WaitOnReady("kube-proxy", "kube-system", 5*time.Second, 3*time.Minute)
		running, err := pod.AreAllPodsRunning("kube-proxy", "kube-system")
		Expect(err).NotTo(HaveOccurred())
		Expect(running).To(Equal(true))
	})

	It("should be able to access the dashboard from each node", func() {
		pod.WaitOnReady("kubernetes-dashboard", "kube-system", 5*time.Second, 3*time.Minute)
		kubeConfig, err := GetConfig()
		Expect(err).NotTo(HaveOccurred())

		sshKeyPath, err := cfg.GetSSHKeyPath()
		Expect(err).NotTo(HaveOccurred())

		s, err := service.Get("kubernetes-dashboard", "kube-system")
		Expect(err).NotTo(HaveOccurred())
		port := s.GetNodePort(80)

		master := fmt.Sprintf("azureuser@%s", kubeConfig.GetServerName())
		nodeList, err := node.Get()
		Expect(err).NotTo(HaveOccurred())

		for _, node := range nodeList.Nodes {
			dashboardURL := fmt.Sprintf("http://%s:%v", node.Status.GetAddressByType("InternalIP").Address, port)
			curlCMD := fmt.Sprintf("curl --max-time 60 %s", dashboardURL)
			_, err := exec.Command("ssh", "-i", sshKeyPath, "-o", "ConnectTimeout=10", "-o", "StrictHostKeyChecking=no", "-o", "UserKnownHostsFile=/dev/null", master, curlCMD).CombinedOutput()
			Expect(err).NotTo(HaveOccurred())
		}
	})

	It("should be able to deploy an nginx service", func() {
		d, err := deployment.Create("library/nginx:latest", cfg.Name, "default")
		Expect(err).NotTo(HaveOccurred())
		err = d.Expose(80)
		Expect(err).NotTo(HaveOccurred())

		s, err := service.Get(cfg.Name, "default")
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

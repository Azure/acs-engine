package kubernetes

import (
	"fmt"
	"log"
	"math/rand"
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
	acct azure.Account
	eng  engine.Engine
	err  error
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
		eng = *e

		err = eng.Generate()
		Expect(err).NotTo(HaveOccurred())

		err = acct.CreateGroup(cfg.Name, cfg.Location)
		Expect(err).NotTo(HaveOccurred())

		// Lets start by just using the normal az group deployment cli for creating a cluster
		log.Println("Creating deployment this make take a few minutes...")
		err = acct.CreateDeployment(cfg.Name, &eng)
		Expect(err).NotTo(HaveOccurred())
	} else {
		e, err := engine.Build(cfg.ClusterDefinition, "_output", cfg.Name)
		Expect(err).NotTo(HaveOccurred())
		eng = *e
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
})

var _ = Describe("Azure Container Cluster using the Kubernetes Orchestrator", func() {
	Context("regardless of agent pool type", func() {
		It("should be logged into the correct account", func() {
			current, err := azure.GetCurrentAccount()

			Expect(err).NotTo(HaveOccurred())
			Expect(current.User.ID).To(Equal(acct.User.ID))
			Expect(current.TenantID).To(Equal(acct.TenantID))
			Expect(current.SubscriptionID).To(Equal(acct.SubscriptionID))
		})

		It("should have have the appropriate node count", func() {
			expectedCount := eng.ClusterDefinition.Properties.MasterProfile.Count
			for _, pool := range eng.ClusterDefinition.Properties.AgentPoolProfiles {
				expectedCount = expectedCount + pool.Count
			}
			nodeList, err := node.Get()
			Expect(err).NotTo(HaveOccurred())
			Expect(len(nodeList.Nodes)).To(Equal(expectedCount))
		})

		It("should be running the expected default version", func() {
			version, err := node.Version()
			Expect(err).NotTo(HaveOccurred())
			Expect(version).To(Equal("v1.6.9"))
		})

		It("should have kube-dns running", func() {
			running, err := pod.WaitOnReady("kube-dns", "kube-system", 5*time.Second, 10*time.Minute)
			Expect(err).NotTo(HaveOccurred())
			Expect(running).To(Equal(true))
		})

		It("should have kube-dashboard running", func() {
			running, err := pod.WaitOnReady("kubernetes-dashboard", "kube-system", 5*time.Second, 10*time.Minute)
			Expect(err).NotTo(HaveOccurred())
			Expect(running).To(Equal(true))
		})

		It("should have kube-proxy running", func() {
			running, err := pod.WaitOnReady("kube-proxy", "kube-system", 5*time.Second, 10*time.Minute)
			Expect(err).NotTo(HaveOccurred())
			Expect(running).To(Equal(true))
		})

		It("should have heapster running", func() {
			running, err := pod.WaitOnReady("heapster", "kube-system", 5*time.Second, 10*time.Minute)
			Expect(err).NotTo(HaveOccurred())
			Expect(running).To(Equal(true))
		})

		It("should have kube-addon-manager running", func() {
			running, err := pod.WaitOnReady("kube-addon-manager", "kube-system", 5*time.Second, 10*time.Minute)
			Expect(err).NotTo(HaveOccurred())
			Expect(running).To(Equal(true))
		})

		It("should have kube-apiserver running", func() {
			running, err := pod.WaitOnReady("kube-apiserver", "kube-system", 5*time.Second, 10*time.Minute)
			Expect(err).NotTo(HaveOccurred())
			Expect(running).To(Equal(true))
		})

		It("should have kube-controller-manager running", func() {
			running, err := pod.WaitOnReady("kube-controller-manager", "kube-system", 5*time.Second, 10*time.Minute)
			Expect(err).NotTo(HaveOccurred())
			Expect(running).To(Equal(true))
		})

		It("should have kube-scheduler running", func() {
			running, err := pod.WaitOnReady("kube-scheduler", "kube-system", 5*time.Second, 10*time.Minute)
			Expect(err).NotTo(HaveOccurred())
			Expect(running).To(Equal(true))
		})

		It("should have tiller running", func() {
			running, err := pod.WaitOnReady("tiller", "kube-system", 5*time.Second, 10*time.Minute)
			Expect(err).NotTo(HaveOccurred())
			Expect(running).To(Equal(true))
		})

		It("should be able to access the dashboard from each node", func() {
			running, err := pod.WaitOnReady("kube-proxy", "kube-system", 5*time.Second, 10*time.Minute)
			Expect(err).NotTo(HaveOccurred())
			Expect(running).To(Equal(true))

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
				output, err := exec.Command("ssh", "-i", sshKeyPath, "-o", "ConnectTimeout=10", "-o", "StrictHostKeyChecking=no", "-o", "UserKnownHostsFile=/dev/null", master, curlCMD).CombinedOutput()
				if err != nil {
					log.Printf("\n\nOutput:%s\n\n", string(output))
				}
			}
		})
	})

	Context("with a linux agent pool", func() {
		It("should be able to deploy an nginx service", func() {
			if eng.HasLinuxAgents() {
				r := rand.New(rand.NewSource(time.Now().UnixNano()))
				deploymentName := fmt.Sprintf("%s-%v", cfg.Name, r.Intn(99999))
				d, err := deployment.CreateLinuxDeploy("library/nginx:latest", deploymentName, "default")
				Expect(err).NotTo(HaveOccurred())

				running, err := pod.WaitOnReady(deploymentName, "default", 5*time.Second, 10*time.Minute)
				Expect(err).NotTo(HaveOccurred())
				Expect(running).To(Equal(true))

				err = d.Expose(80, 80)
				Expect(err).NotTo(HaveOccurred())

				s, err := service.Get(deploymentName, "default")
				Expect(err).NotTo(HaveOccurred())
				s, err = s.WaitForExternalIP(10*time.Minute, 5*time.Second)
				Expect(err).NotTo(HaveOccurred())
				Expect(s.Status.LoadBalancer.Ingress).NotTo(BeEmpty())

				valid := s.Validate("(Welcome to nginx)", 5, 5*time.Second)
				Expect(valid).To(BeTrue())
			} else {
				Skip("No linux agent was provisioned for this Cluster Definition")
			}
		})
	})

	Context("with a windows agent pool", func() {
		It("should be able to deploy an iis webserver", func() {
			if eng.HasWindowsAgents() {
				r := rand.New(rand.NewSource(time.Now().UnixNano()))
				deploymentName := fmt.Sprintf("%s-%v", cfg.Name, r.Intn(99999))
				d, err := deployment.CreateWindowsDeploy("microsoft/iis", deploymentName, "default", 80)
				Expect(err).NotTo(HaveOccurred())

				running, err := pod.WaitOnReady(deploymentName, "default", 5*time.Second, 15*time.Minute)
				Expect(err).NotTo(HaveOccurred())
				Expect(running).To(Equal(true))

				err = d.Expose(80, 80)
				Expect(err).NotTo(HaveOccurred())

				s, err := service.Get(deploymentName, "default")
				Expect(err).NotTo(HaveOccurred())
				s, err = s.WaitForExternalIP(10*time.Minute, 5*time.Second)
				Expect(err).NotTo(HaveOccurred())
				Expect(s.Status.LoadBalancer.Ingress).NotTo(BeEmpty())

				valid := s.Validate("(IIS Windows Server)", 5, 5*time.Second)
				Expect(valid).To(BeTrue())
			} else {
				Skip("No windows agent was provisioned for this Cluster Definition")
			}
		})
	})
})

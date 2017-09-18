package kubernetes

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/Azure/acs-engine/pkg/api"
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
	cfg config.Config
	eng engine.Engine
	err error
)

var _ = BeforeSuite(func() {
	cwd, _ := os.Getwd()
	rootPath := filepath.Join(cwd, "../../..") // The current working dir of these tests is down a few levels from the root of the project. We should traverse up that path so we can find the _output dir
	c, err := config.ParseConfig()
	c.CurrentWorkingDir = rootPath
	Expect(err).NotTo(HaveOccurred())
	cfg = *c // We have to do this because golang anon functions and scoping and stuff

	engCfg, err := engine.ParseConfig(c.CurrentWorkingDir, c.ClusterDefinition, c.Name)
	Expect(err).NotTo(HaveOccurred())
	cs, err := engine.Parse(engCfg.ClusterDefinitionTemplate)
	Expect(err).NotTo(HaveOccurred())
	eng = engine.Engine{
		Config:            engCfg,
		ClusterDefinition: cs,
	}
})

var _ = Describe("Azure Container Cluster using the Kubernetes Orchestrator", func() {
	Context("regardless of agent pool type", func() {

		It("should have have the appropriate node count", func() {
			nodeList, err := node.Get()
			Expect(err).NotTo(HaveOccurred())
			Expect(len(nodeList.Nodes)).To(Equal(eng.NodeCount()))
		})

		It("should be running the expected version", func() {
			version, err := node.Version()
			Expect(err).NotTo(HaveOccurred())

			if eng.ClusterDefinition.Properties.OrchestratorProfile.OrchestratorRelease != "" {
				Expect(version).To(MatchRegexp("v" + api.KubernetesReleaseToVersion[eng.ClusterDefinition.Properties.OrchestratorProfile.OrchestratorRelease]))
			} else {
				Expect(version).To(Equal("v" + api.KubernetesReleaseToVersion[api.KubernetesDefaultRelease]))
			}
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
			running, err := pod.WaitOnReady("kubernetes-dashboard", "kube-system", 5*time.Second, 10*time.Minute)
			Expect(err).NotTo(HaveOccurred())
			Expect(running).To(Equal(true))

			kubeConfig, err := GetConfig()
			Expect(err).NotTo(HaveOccurred())
			sshKeyPath := cfg.GetSSHKeyPath()

			s, err := service.Get("kubernetes-dashboard", "kube-system")
			Expect(err).NotTo(HaveOccurred())
			port := s.GetNodePort(80)

			master := fmt.Sprintf("azureuser@%s", kubeConfig.GetServerName())
			nodeList, err := node.Get()
			Expect(err).NotTo(HaveOccurred())

			for _, node := range nodeList.Nodes {
				success := false
				for i := 0; i < 20; i++ {
					dashboardURL := fmt.Sprintf("http://%s:%v", node.Status.GetAddressByType("InternalIP").Address, port)
					curlCMD := fmt.Sprintf("curl --max-time 60 %s", dashboardURL)
					output, err := exec.Command("ssh", "-i", sshKeyPath, "-o", "ConnectTimeout=10", "-o", "StrictHostKeyChecking=no", "-o", "UserKnownHostsFile=/dev/null", master, curlCMD).CombinedOutput()
					if err != nil {
						log.Printf("Error on iteration:%v for node (%s)\n", i, node.Metadata.Name)
						log.Printf("Command:ssh -i %s -o ConnectTimeout=10 -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null %s %s\n", sshKeyPath, master, curlCMD)
						log.Printf("\nOutput:%s\n", string(output))
					} else {
						success = true
					}
					time.Sleep(5 * time.Second)
				}
				Expect(success).To(BeTrue())
			}
		})
	})

	Context("with a linux agent pool", func() {
		It("should be able to deploy an nginx service", func() {
			if eng.HasLinuxAgents() {
				r := rand.New(rand.NewSource(time.Now().UnixNano()))
				deploymentName := fmt.Sprintf("nginx-%s-%v", cfg.Name, r.Intn(99999))
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
				deploymentName := fmt.Sprintf("iis-%s-%v", cfg.Name, r.Intn(99999))
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

package kubernetes

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"time"

	"github.com/Azure/acs-engine/pkg/api"
	"github.com/Azure/acs-engine/test/e2e/config"
	"github.com/Azure/acs-engine/test/e2e/engine"
	"github.com/Azure/acs-engine/test/e2e/kubernetes/deployment"
	"github.com/Azure/acs-engine/test/e2e/kubernetes/node"
	"github.com/Azure/acs-engine/test/e2e/kubernetes/pod"
	"github.com/Azure/acs-engine/test/e2e/kubernetes/service"
	"github.com/Azure/acs-engine/test/e2e/remote"
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
	Describe("regardless of agent pool type", func() {

		It("should have have the appropriate node count", func() {
			nodeList, err := node.Get()
			Expect(err).NotTo(HaveOccurred())
			Expect(len(nodeList.Nodes)).To(Equal(eng.NodeCount()))
		})

		It("should be running the expected version", func() {
			version, err := node.Version()
			Expect(err).NotTo(HaveOccurred())

			if eng.ClusterDefinition.Properties.OrchestratorProfile.OrchestratorVersion != "" {
				Expect(version).To(MatchRegexp("v" + eng.ClusterDefinition.Properties.OrchestratorProfile.OrchestratorVersion))
			} else {
				Expect(version).To(Equal("v" + api.KubernetesDefaultVersion))
			}
		})

		/* The master nodes are hidden behind a load balancer. Therefore, we will create an ssh connection and then continue to reuse that connection for subsequent commands. We will iterate the nodes first to make sure that we ssh onto each host from a given master and then the inner loop will verify that we cannot connect to another master's etcd instance. If we see a "Host key verification failed" error this is an indication that we are trying to ssh onto a host that we are already on. Then we will just execute the etcdctl command locally. */
		It("should not expose etcd to the internet", func() {
			hostKeyRegex, err := regexp.Compile("Host key verification failed")
			Expect(err).NotTo(HaveOccurred())

			nodes, err := node.GetByPrefix("k8s-master")
			Expect(err).NotTo(HaveOccurred())
			Expect(len(nodes)).NotTo(Equal(0))

			conn, err := remote.NewConnection(fmt.Sprintf("%s.%s.cloudapp.azure.com", cfg.Name, cfg.Location), "22", eng.ClusterDefinition.Properties.LinuxProfile.AdminUsername, cfg.GetSSHKeyPath())
			Expect(err).NotTo(HaveOccurred())

			hostname, err := conn.Execute("hostname")
			Expect(err).NotTo(HaveOccurred())
			for _, n := range nodes {
				for _, nprime := range nodes {
					// I am doing this to validate that we always run these commands from the same host
					host, err := conn.Execute("hostname")
					Expect(err).NotTo(HaveOccurred())
					Expect(hostname).To(Equal(host))

					if n.Metadata.Name != nprime.Metadata.Name {
						etcdCmd := fmt.Sprintf("etcdctl --endpoint=http://%s:2379 ls /registry/secrets/kube-system", nprime.Status.GetAddressByType("InternalIP").Address)
						cmd := fmt.Sprintf("ssh %s@%s %s", eng.ClusterDefinition.Properties.LinuxProfile.AdminUsername, n.Metadata.Name, etcdCmd)

						out, err := conn.Execute(cmd)
						matched := hostKeyRegex.MatchString(string(out))
						if !matched {
							Expect(err).To(HaveOccurred())
							Expect(out).To(MatchRegexp("connection refused"))
						} else {
							out, err := conn.Execute(etcdCmd)
							Expect(err).To(HaveOccurred())
							Expect(out).To(MatchRegexp("connection refused"))
						}
					}
				}
			}
		})

		It("should have kube-dns running", func() {
			running, err := pod.WaitOnReady("kube-dns", "kube-system", 5*time.Second, cfg.Timeout)
			Expect(err).NotTo(HaveOccurred())
			Expect(running).To(Equal(true))
		})

		It("should have kube-dashboard running", func() {
			running, err := pod.WaitOnReady("kubernetes-dashboard", "kube-system", 5*time.Second, cfg.Timeout)
			Expect(err).NotTo(HaveOccurred())
			Expect(running).To(Equal(true))
		})

		It("should have kube-proxy running", func() {
			running, err := pod.WaitOnReady("kube-proxy", "kube-system", 5*time.Second, cfg.Timeout)
			Expect(err).NotTo(HaveOccurred())
			Expect(running).To(Equal(true))
		})

		It("should have heapster running", func() {
			running, err := pod.WaitOnReady("heapster", "kube-system", 5*time.Second, cfg.Timeout)
			Expect(err).NotTo(HaveOccurred())
			Expect(running).To(Equal(true))
		})

		It("should have kube-addon-manager running", func() {
			running, err := pod.WaitOnReady("kube-addon-manager", "kube-system", 5*time.Second, cfg.Timeout)
			Expect(err).NotTo(HaveOccurred())
			Expect(running).To(Equal(true))
		})

		It("should have kube-apiserver running", func() {
			running, err := pod.WaitOnReady("kube-apiserver", "kube-system", 5*time.Second, cfg.Timeout)
			Expect(err).NotTo(HaveOccurred())
			Expect(running).To(Equal(true))
		})

		It("should have kube-controller-manager running", func() {
			running, err := pod.WaitOnReady("kube-controller-manager", "kube-system", 5*time.Second, cfg.Timeout)
			Expect(err).NotTo(HaveOccurred())
			Expect(running).To(Equal(true))
		})

		It("should have kube-scheduler running", func() {
			running, err := pod.WaitOnReady("kube-scheduler", "kube-system", 5*time.Second, cfg.Timeout)
			Expect(err).NotTo(HaveOccurred())
			Expect(running).To(Equal(true))
		})

		It("should have tiller running", func() {
			running, err := pod.WaitOnReady("tiller", "kube-system", 5*time.Second, cfg.Timeout)
			Expect(err).NotTo(HaveOccurred())
			Expect(running).To(Equal(true))
		})

		It("should be able to access the dashboard from each node", func() {
			running, err := pod.WaitOnReady("kubernetes-dashboard", "kube-system", 5*time.Second, cfg.Timeout)
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
				for i := 0; i < 60; i++ {
					dashboardURL := fmt.Sprintf("http://%s:%v", node.Status.GetAddressByType("InternalIP").Address, port)
					curlCMD := fmt.Sprintf("curl --max-time 60 %s", dashboardURL)
					_, err := exec.Command("ssh", "-i", sshKeyPath, "-o", "ConnectTimeout=10", "-o", "StrictHostKeyChecking=no", "-o", "UserKnownHostsFile=/dev/null", master, curlCMD).CombinedOutput()
					if err == nil {
						success = true
						break
					}
					time.Sleep(10 * time.Second)
				}
				Expect(success).To(BeTrue())
			}
		})
	})

	Describe("with a linux agent pool", func() {
		It("should be able to deploy an nginx service", func() {
			if eng.HasLinuxAgents() {
				r := rand.New(rand.NewSource(time.Now().UnixNano()))
				deploymentName := fmt.Sprintf("nginx-%s-%v", cfg.Name, r.Intn(99999))
				nginxDeploy, err := deployment.CreateLinuxDeploy("library/nginx:latest", deploymentName, "default")
				Expect(err).NotTo(HaveOccurred())

				running, err := pod.WaitOnReady(deploymentName, "default", 5*time.Second, cfg.Timeout)
				Expect(err).NotTo(HaveOccurred())
				Expect(running).To(Equal(true))

				err = nginxDeploy.Expose(80, 80)
				Expect(err).NotTo(HaveOccurred())

				s, err := service.Get(deploymentName, "default")
				Expect(err).NotTo(HaveOccurred())
				s, err = s.WaitForExternalIP(cfg.Timeout, 5*time.Second)
				Expect(err).NotTo(HaveOccurred())
				Expect(s.Status.LoadBalancer.Ingress).NotTo(BeEmpty())

				valid := s.Validate("(Welcome to nginx)", 5, 5*time.Second)
				Expect(valid).To(BeTrue())

				nginxPods, err := nginxDeploy.Pods()
				Expect(err).NotTo(HaveOccurred())
				Expect(len(nginxPods)).ToNot(BeZero())
				for _, nginxPod := range nginxPods {
					pass, err := nginxPod.CheckLinuxOutboundConnection(5*time.Second, cfg.Timeout)
					Expect(err).NotTo(HaveOccurred())
					Expect(pass).To(BeTrue())
				}
			} else {
				Skip("No linux agent was provisioned for this Cluster Definition")
			}
		})
	})

	Describe("with a windows agent pool", func() {
		It("should be able to deploy an iis webserver", func() {
			if eng.HasWindowsAgents() {
				r := rand.New(rand.NewSource(time.Now().UnixNano()))
				deploymentName := fmt.Sprintf("iis-%s-%v", cfg.Name, r.Intn(99999))
				iisDeploy, err := deployment.CreateWindowsDeploy("microsoft/iis", deploymentName, "default", 80)
				Expect(err).NotTo(HaveOccurred())

				running, err := pod.WaitOnReady(deploymentName, "default", 5*time.Second, cfg.Timeout)
				Expect(err).NotTo(HaveOccurred())
				Expect(running).To(Equal(true))

				err = iisDeploy.Expose(80, 80)
				Expect(err).NotTo(HaveOccurred())

				s, err := service.Get(deploymentName, "default")
				Expect(err).NotTo(HaveOccurred())
				s, err = s.WaitForExternalIP(cfg.Timeout, 5*time.Second)
				Expect(err).NotTo(HaveOccurred())
				Expect(s.Status.LoadBalancer.Ingress).NotTo(BeEmpty())

				valid := s.Validate("(IIS Windows Server)", 5, 5*time.Second)
				Expect(valid).To(BeTrue())

				iisPods, err := iisDeploy.Pods()
				Expect(err).NotTo(HaveOccurred())
				Expect(len(iisPods)).ToNot(BeZero())
				for _, iisPod := range iisPods {
					pass, err := iisPod.CheckWindowsOutboundConnection(5*time.Second, cfg.Timeout)
					Expect(err).NotTo(HaveOccurred())
					Expect(pass).To(BeTrue())
				}
			} else {
				Skip("No windows agent was provisioned for this Cluster Definition")
			}
		})
	})
})

package kubernetes

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/Azure/acs-engine/pkg/api/common"
	"github.com/Azure/acs-engine/test/e2e/config"
	"github.com/Azure/acs-engine/test/e2e/engine"
	"github.com/Azure/acs-engine/test/e2e/kubernetes/deployment"
	"github.com/Azure/acs-engine/test/e2e/kubernetes/node"
	"github.com/Azure/acs-engine/test/e2e/kubernetes/persistentvolumeclaims"
	"github.com/Azure/acs-engine/test/e2e/kubernetes/pod"
	"github.com/Azure/acs-engine/test/e2e/kubernetes/service"
	"github.com/Azure/acs-engine/test/e2e/kubernetes/storageclass"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	WorkloadDir = "workloads"
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
				Expect(version).To(Equal("v" + common.KubernetesDefaultVersion))
			}
		})

		It("should have kube-dns running", func() {
			running, err := pod.WaitOnReady("kube-dns", "kube-system", 3, 30*time.Second, cfg.Timeout)
			Expect(err).NotTo(HaveOccurred())
			Expect(running).To(Equal(true))
		})

		It("should have kube-dashboard running", func() {
			running, err := pod.WaitOnReady("kubernetes-dashboard", "kube-system", 3, 30*time.Second, cfg.Timeout)
			Expect(err).NotTo(HaveOccurred())
			Expect(running).To(Equal(true))
		})

		It("should have kube-proxy running", func() {
			running, err := pod.WaitOnReady("kube-proxy", "kube-system", 3, 30*time.Second, cfg.Timeout)
			Expect(err).NotTo(HaveOccurred())
			Expect(running).To(Equal(true))
		})

		It("should have heapster running", func() {
			running, err := pod.WaitOnReady("heapster", "kube-system", 3, 30*time.Second, cfg.Timeout)
			Expect(err).NotTo(HaveOccurred())
			Expect(running).To(Equal(true))
		})

		It("should have kube-addon-manager running", func() {
			running, err := pod.WaitOnReady("kube-addon-manager", "kube-system", 3, 30*time.Second, cfg.Timeout)
			Expect(err).NotTo(HaveOccurred())
			Expect(running).To(Equal(true))
		})

		It("should have kube-apiserver running", func() {
			running, err := pod.WaitOnReady("kube-apiserver", "kube-system", 3, 30*time.Second, cfg.Timeout)
			Expect(err).NotTo(HaveOccurred())
			Expect(running).To(Equal(true))
		})

		It("should have kube-controller-manager running", func() {
			running, err := pod.WaitOnReady("kube-controller-manager", "kube-system", 3, 30*time.Second, cfg.Timeout)
			Expect(err).NotTo(HaveOccurred())
			Expect(running).To(Equal(true))
		})

		It("should have kube-scheduler running", func() {
			running, err := pod.WaitOnReady("kube-scheduler", "kube-system", 3, 30*time.Second, cfg.Timeout)
			Expect(err).NotTo(HaveOccurred())
			Expect(running).To(Equal(true))
		})

		It("should have tiller running", func() {
			running, err := pod.WaitOnReady("tiller", "kube-system", 3, 30*time.Second, cfg.Timeout)
			Expect(err).NotTo(HaveOccurred())
			Expect(running).To(Equal(true))
		})

		It("should be able to access the dashboard from each node", func() {
			running, err := pod.WaitOnReady("kubernetes-dashboard", "kube-system", 3, 30*time.Second, cfg.Timeout)
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

				running, err := pod.WaitOnReady(deploymentName, "default", 3, 30*time.Second, cfg.Timeout)
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
				iisDeploy, err := deployment.CreateWindowsDeploy("microsoft/iis:windowsservercore-1709", deploymentName, "default", 80, -1)
				Expect(err).NotTo(HaveOccurred())

				running, err := pod.WaitOnReady(deploymentName, "default", 3, 30*time.Second, cfg.Timeout)
				Expect(err).NotTo(HaveOccurred())
				Expect(running).To(Equal(true))

				err = iisDeploy.Expose(80, 80)
				Expect(err).NotTo(HaveOccurred())

				s, err := service.Get(deploymentName, "default")
				Expect(err).NotTo(HaveOccurred())
				s, err = s.WaitForExternalIP(cfg.Timeout, 5*time.Second)
				Expect(err).NotTo(HaveOccurred())
				Expect(s.Status.LoadBalancer.Ingress).NotTo(BeEmpty())

				valid := s.Validate("(IIS Windows Server)", 10, 10*time.Second)
				Expect(valid).To(BeTrue())

				iisPods, err := iisDeploy.Pods()
				Expect(err).NotTo(HaveOccurred())
				Expect(len(iisPods)).ToNot(BeZero())
				for _, iisPod := range iisPods {
					pass, err := iisPod.CheckWindowsOutboundConnection(10*time.Second, cfg.Timeout)
					Expect(err).NotTo(HaveOccurred())
					Expect(pass).To(BeTrue())
				}

				err = iisDeploy.Delete()
				Expect(err).NotTo(HaveOccurred())
			} else {
				Skip("No windows agent was provisioned for this Cluster Definition")
			}
		})

		It("should be able to reach hostport in an iis webserver", func() {
			if eng.HasWindowsAgents() {
				r := rand.New(rand.NewSource(time.Now().UnixNano()))
				hostport := 8123
				deploymentName := fmt.Sprintf("iis-%s-%v", cfg.Name, r.Intn(99999))
				iisDeploy, err := deployment.CreateWindowsDeploy("microsoft/iis:windowsservercore-1709", deploymentName, "default", 80, hostport)
				Expect(err).NotTo(HaveOccurred())

				running, err := pod.WaitOnReady(deploymentName, "default", 3, 30*time.Second, cfg.Timeout)
				Expect(err).NotTo(HaveOccurred())
				Expect(running).To(Equal(true))

				iisPods, err := iisDeploy.Pods()
				Expect(err).NotTo(HaveOccurred())
				Expect(len(iisPods)).ToNot(BeZero())

				kubeConfig, err := GetConfig()
				Expect(err).NotTo(HaveOccurred())
				master := fmt.Sprintf("azureuser@%s", kubeConfig.GetServerName())
				sshKeyPath := cfg.GetSSHKeyPath()

				for _, iisPod := range iisPods {
					valid := iisPod.ValidateHostPort("(IIS Windows Server)", 10, 10*time.Second, master, sshKeyPath)
					Expect(valid).To(BeTrue())
				}

				err = iisDeploy.Delete()
				Expect(err).NotTo(HaveOccurred())
			} else {
				Skip("No windows agent was provisioned for this Cluster Definition")
			}
		})

		It("should be able to attach azure file", func() {
			if eng.HasWindowsAgents() {
				if eng.OrchestratorVersion1Dot8AndUp() {
					storageclassName := "azurefile" // should be the same as in storageclass-azurefile.yaml
					sc, err := storageclass.CreateStorageClassFromFile(filepath.Join(WorkloadDir, "storageclass-azurefile.yaml"), storageclassName)
					Expect(err).NotTo(HaveOccurred())
					ready, err := sc.WaitOnReady(5*time.Second, cfg.Timeout)
					Expect(err).NotTo(HaveOccurred())
					Expect(ready).To(Equal(true))

					pvcName := "pvc-azurefile" // should be the same as in pvc-azurefile.yaml
					pvc, err := persistentvolumeclaims.CreatePersistentVolumeClaimsFromFile(filepath.Join(WorkloadDir, "pvc-azurefile.yaml"), pvcName, "default")
					Expect(err).NotTo(HaveOccurred())
					ready, err = pvc.WaitOnReady("default", 5*time.Second, cfg.Timeout)
					Expect(err).NotTo(HaveOccurred())
					Expect(ready).To(Equal(true))

					podName := "iis-azurefile" // should be the same as in iis-azurefile.yaml
					iisPod, err := pod.CreatePodFromFile(filepath.Join(WorkloadDir, "iis-azurefile.yaml"), podName, "default")
					Expect(err).NotTo(HaveOccurred())
					ready, err = iisPod.WaitOnReady(5*time.Second, cfg.Timeout)
					Expect(err).NotTo(HaveOccurred())
					Expect(ready).To(Equal(true))

					valid, err := iisPod.ValidateAzureFile("mnt\\azure", 10, 10*time.Second)
					Expect(valid).To(BeTrue())
					Expect(err).NotTo(HaveOccurred())

					err = iisPod.Delete()
					Expect(err).NotTo(HaveOccurred())
				} else {
					Skip("Kubernetes version needs to be 1.8 and up for Azure File test")
				}
			} else {
				Skip("No windows agent was provisioned for this Cluster Definition")
			}
		})
	})
})

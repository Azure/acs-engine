package kubernetes

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"time"

	"github.com/Azure/acs-engine/pkg/api/common"
	"github.com/Azure/acs-engine/test/e2e/config"
	"github.com/Azure/acs-engine/test/e2e/engine"
	"github.com/Azure/acs-engine/test/e2e/kubernetes/deployment"
	"github.com/Azure/acs-engine/test/e2e/kubernetes/job"
	"github.com/Azure/acs-engine/test/e2e/kubernetes/node"
	"github.com/Azure/acs-engine/test/e2e/kubernetes/pod"
	"github.com/Azure/acs-engine/test/e2e/kubernetes/service"
	"github.com/Azure/acs-engine/test/e2e/kubernetes/util"
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
	csInput, err := engine.ParseInput(engCfg.ClusterDefinitionTemplate)
	Expect(err).NotTo(HaveOccurred())
	csGenerated, err := engine.ParseOutput(engCfg.GeneratedDefinitionPath + "/apimodel.json")
	Expect(err).NotTo(HaveOccurred())
	eng = engine.Engine{
		Config:             engCfg,
		ClusterDefinition:  csInput,
		ExpandedDefinition: csGenerated,
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

			var expectedVersion string
			if eng.ClusterDefinition.Properties.OrchestratorProfile.OrchestratorRelease != "" ||
				eng.ClusterDefinition.Properties.OrchestratorProfile.OrchestratorVersion != "" {
				expectedVersion = common.RationalizeReleaseAndVersion(
					common.Kubernetes,
					eng.ClusterDefinition.Properties.OrchestratorProfile.OrchestratorRelease,
					eng.ClusterDefinition.Properties.OrchestratorProfile.OrchestratorVersion)
			} else {
				expectedVersion = common.RationalizeReleaseAndVersion(
					common.Kubernetes,
					eng.Config.OrchestratorRelease,
					eng.Config.OrchestratorVersion)
			}
			Expect(version).To(Equal("v" + expectedVersion))
		})

		It("should have kube-dns running", func() {
			running, err := pod.WaitOnReady("kube-dns", "kube-system", 3, 30*time.Second, cfg.Timeout)
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
			if hasTiller, tillerAddon := eng.HasAddon("tiller"); hasTiller {
				running, err := pod.WaitOnReady("tiller", "kube-system", 3, 30*time.Second, cfg.Timeout)
				Expect(err).NotTo(HaveOccurred())
				Expect(running).To(Equal(true))
				pods, err := pod.GetAllByPrefix("tiller-deploy", "kube-system")
				Expect(err).NotTo(HaveOccurred())
				By("Ensuring that the correct max-history has been applied")
				maxHistory := tillerAddon.Config["max-history"]
				// There is only one tiller pod and one container in that pod
				actualTillerMaxHistory, err := pods[0].Spec.Containers[0].GetEnvironmentVariable("TILLER_HISTORY_MAX")
				Expect(err).NotTo(HaveOccurred())
				Expect(actualTillerMaxHistory).To(Equal(maxHistory))
				By("Ensuring that the correct resources have been applied")
				err = pods[0].Spec.Containers[0].ValidateResources(tillerAddon.Containers[0])
				Expect(err).NotTo(HaveOccurred())
			} else {
				Skip("tiller disabled for this cluster, will not test")
			}
		})

		It("should be able to access the dashboard from each node", func() {
			if hasDashboard, dashboardAddon := eng.HasAddon("kubernetes-dashboard"); hasDashboard {
				By("Ensuring that the kubernetes-dashboard pod is Running")

				running, err := pod.WaitOnReady("kubernetes-dashboard", "kube-system", 3, 30*time.Second, cfg.Timeout)
				Expect(err).NotTo(HaveOccurred())
				Expect(running).To(Equal(true))

				By("Ensuring that the kubernetes-dashboard service is Running")

				s, err := service.Get("kubernetes-dashboard", "kube-system")
				Expect(err).NotTo(HaveOccurred())

				if !eng.HasWindowsAgents() {
					By("Gathering connection information to determine whether or not to connect via HTTP or HTTPS")
					dashboardPort := 80
					version, err := node.Version()
					Expect(err).NotTo(HaveOccurred())
					re := regexp.MustCompile("v1.9")
					if re.FindString(version) != "" {
						dashboardPort = 443
					}
					port := s.GetNodePort(dashboardPort)

					kubeConfig, err := GetConfig()
					Expect(err).NotTo(HaveOccurred())
					master := fmt.Sprintf("azureuser@%s", kubeConfig.GetServerName())

					sshKeyPath := cfg.GetSSHKeyPath()

					if dashboardPort == 80 {
						By("Ensuring that we can connect via HTTP to the dashboard on any one node")
					} else {
						By("Ensuring that we can connect via HTTPS to the dashboard on any one node")
					}
					nodeList, err := node.Get()
					Expect(err).NotTo(HaveOccurred())
					for _, node := range nodeList.Nodes {
						success := false
						for i := 0; i < 60; i++ {
							dashboardURL := fmt.Sprintf("http://%s:%v", node.Status.GetAddressByType("InternalIP").Address, port)
							curlCMD := fmt.Sprintf("curl --max-time 60 %s", dashboardURL)
							cmd := exec.Command("ssh", "-i", sshKeyPath, "-o", "ConnectTimeout=10", "-o", "StrictHostKeyChecking=no", "-o", "UserKnownHostsFile=/dev/null", master, curlCMD)
							util.PrintCommand(cmd)
							out, err := cmd.CombinedOutput()
							if err == nil {
								success = true
								break
							}
							if i > 58 {
								log.Printf("Error while connecting to Windows dashboard:%s\n", err)
								log.Println(string(out))
							}
							time.Sleep(10 * time.Second)
						}
						Expect(success).To(BeTrue())
					}
					By("Ensuring that the correct resources have been applied")
					// Assuming one dashboard pod
					pods, err := pod.GetAllByPrefix("kubernetes-dashboard", "kube-system")
					Expect(err).NotTo(HaveOccurred())
					for i, c := range dashboardAddon.Containers {
						err := pods[0].Spec.Containers[i].ValidateResources(c)
						Expect(err).NotTo(HaveOccurred())
					}
				}
			} else {
				Skip("kubernetes-dashboard disabled for this cluster, will not test")
			}
		})

		It("should have aci-connector running", func() {
			if hasACIConnector, ACIConnectorAddon := eng.HasAddon("aci-connector"); hasACIConnector {
				running, err := pod.WaitOnReady("aci-connector", "kube-system", 3, 30*time.Second, cfg.Timeout)
				Expect(err).NotTo(HaveOccurred())
				Expect(running).To(Equal(true))
				By("Ensuring that the correct resources have been applied")
				// Assuming one aci-connector pod
				pods, err := pod.GetAllByPrefix("aci-connector", "kube-system")
				Expect(err).NotTo(HaveOccurred())
				for i, c := range ACIConnectorAddon.Containers {
					err := pods[0].Spec.Containers[i].ValidateResources(c)
					Expect(err).NotTo(HaveOccurred())
				}

			} else {
				Skip("aci-connector disabled for this cluster, will not test")
			}
		})

		It("should have rescheduler running", func() {
			if hasRescheduler, reschedulerAddon := eng.HasAddon("rescheduler"); hasRescheduler {
				running, err := pod.WaitOnReady("rescheduler", "kube-system", 3, 30*time.Second, cfg.Timeout)
				Expect(err).NotTo(HaveOccurred())
				Expect(running).To(Equal(true))
				By("Ensuring that the correct resources have been applied")
				// Assuming one rescheduler pod
				pods, err := pod.GetAllByPrefix("rescheduler", "kube-system")
				Expect(err).NotTo(HaveOccurred())
				for i, c := range reschedulerAddon.Containers {
					err := pods[0].Spec.Containers[i].ValidateResources(c)
					Expect(err).NotTo(HaveOccurred())
				}
			} else {
				Skip("rescheduler disabled for this cluster, will not test")
			}
		})
	})

	Describe("with a linux agent pool", func() {
		It("should be able to autoscale", func() {
			if eng.HasLinuxAgents() {
				By("Creating a test php-apache deployment with request limit thresholds")
				// Inspired by http://blog.kubernetes.io/2016/07/autoscaling-in-kubernetes.html
				r := rand.New(rand.NewSource(time.Now().UnixNano()))
				phpApacheName := fmt.Sprintf("php-apache-%s-%v", cfg.Name, r.Intn(99999))
				phpApacheDeploy, err := deployment.CreateLinuxDeploy("gcr.io/google_containers/hpa-example", phpApacheName, "default", "--requests=cpu=50m,memory=50M")
				if err != nil {
					fmt.Println(err)
				}
				Expect(err).NotTo(HaveOccurred())

				By("Ensuring that one php-apache pod is running before autoscale configuration or load applied")
				running, err := pod.WaitOnReady(phpApacheName, "default", 3, 30*time.Second, cfg.Timeout)
				Expect(err).NotTo(HaveOccurred())
				Expect(running).To(Equal(true))

				phpPods, err := phpApacheDeploy.Pods()
				Expect(err).NotTo(HaveOccurred())
				// We should have exactly 1 pod to begin
				Expect(len(phpPods)).To(Equal(1))

				By("Exposing TCP 80 internally on the php-apache deployment")
				err = phpApacheDeploy.Expose("ClusterIP", 80, 80)
				Expect(err).NotTo(HaveOccurred())
				s, err := service.Get(phpApacheName, "default")
				Expect(err).NotTo(HaveOccurred())

				By("Assigning hpa configuration to the php-apache deployment")
				// Apply autoscale characteristics to deployment
				err = phpApacheDeploy.CreateDeploymentHPA(5, 1, 10)
				Expect(err).NotTo(HaveOccurred())

				By("Sending load to the php-apache service by creating a 3 replica deployment")
				// Launch a simple busybox pod that wget's continuously to the apache serviceto simulate load
				commandString := fmt.Sprintf("while true; do wget -q -O- http://%s.default.svc.cluster.local; done", phpApacheName)
				loadTestName := fmt.Sprintf("load-test-%s-%v", cfg.Name, r.Intn(99999))
				numLoadTestPods := 3
				loadTestDeploy, err := deployment.RunLinuxDeploy("busybox", loadTestName, "default", commandString, numLoadTestPods)
				Expect(err).NotTo(HaveOccurred())

				By("Ensuring there are 3 load test pods")
				running, err = pod.WaitOnReady(loadTestName, "default", 3, 30*time.Second, cfg.Timeout)
				Expect(err).NotTo(HaveOccurred())
				Expect(running).To(Equal(true))

				// We should have three load tester pods running
				loadTestPods, err := loadTestDeploy.Pods()
				Expect(err).NotTo(HaveOccurred())
				Expect(len(loadTestPods)).To(Equal(numLoadTestPods))

				By("Waiting 3 minutes for load to take effect")
				// Wait 3 minutes for autoscaler to respond to load
				time.Sleep(3 * time.Minute)

				By("Ensuring we have more than 1 apache-php pods due to hpa enforcement")
				phpPods, err = phpApacheDeploy.Pods()
				Expect(err).NotTo(HaveOccurred())
				// We should have > 1 pods after autoscale effects
				Expect(len(phpPods) > 1).To(BeTrue())

				By("Cleaning up after ourselves")
				err = loadTestDeploy.Delete()
				Expect(err).NotTo(HaveOccurred())
				err = phpApacheDeploy.Delete()
				Expect(err).NotTo(HaveOccurred())
				err = s.Delete()
				Expect(err).NotTo(HaveOccurred())
			} else {
				Skip("This flavor/version of Kubernetes doesn't support hpa autoscale")
			}
		})

		It("should be able to deploy an nginx service", func() {
			if eng.HasLinuxAgents() {
				By("Creating a nginx deployment")
				r := rand.New(rand.NewSource(time.Now().UnixNano()))
				deploymentName := fmt.Sprintf("nginx-%s-%v", cfg.Name, r.Intn(99999))
				nginxDeploy, err := deployment.CreateLinuxDeploy("library/nginx:latest", deploymentName, "default", "")
				Expect(err).NotTo(HaveOccurred())

				By("Ensure there is a Running nginx pod")
				running, err := pod.WaitOnReady(deploymentName, "default", 3, 30*time.Second, cfg.Timeout)
				Expect(err).NotTo(HaveOccurred())
				Expect(running).To(Equal(true))

				By("Exposing TCP 80 LB on the nginx deployment")
				err = nginxDeploy.Expose("LoadBalancer", 80, 80)
				Expect(err).NotTo(HaveOccurred())

				By("Ensuring we can connect to the service")
				s, err := service.Get(deploymentName, "default")
				Expect(err).NotTo(HaveOccurred())

				By("Ensuring the service root URL returns the expected payload")
				valid := s.Validate("(Welcome to nginx)", 5, 30*time.Second, cfg.Timeout)
				Expect(valid).To(BeTrue())

				By("Ensuring we have outbound internet access from the nginx pods")
				nginxPods, err := nginxDeploy.Pods()
				Expect(err).NotTo(HaveOccurred())
				Expect(len(nginxPods)).ToNot(BeZero())
				for _, nginxPod := range nginxPods {
					pass, err := nginxPod.CheckLinuxOutboundConnection(5*time.Second, cfg.Timeout)
					Expect(err).NotTo(HaveOccurred())
					Expect(pass).To(BeTrue())
				}

				By("Cleaning up after ourselves")
				err = nginxDeploy.Delete()
				Expect(err).NotTo(HaveOccurred())
				err = s.Delete()
				Expect(err).NotTo(HaveOccurred())
			} else {
				Skip("No linux agent was provisioned for this Cluster Definition")
			}
		})
	})

	Describe("with a GPU-enabled agent pool", func() {
		It("should be able to run a nvidia-gpu job", func() {
			if eng.HasGPUNodes() {
				j, err := job.CreateJobFromFile(filepath.Join(WorkloadDir, "nvidia-smi.yaml"), "nvidia-smi", "default")
				Expect(err).NotTo(HaveOccurred())
				ready, err := j.WaitOnReady(30*time.Second, cfg.Timeout)
				delErr := j.Delete()
				if delErr != nil {
					fmt.Printf("could not delete job %s\n", j.Metadata.Name)
					fmt.Println(delErr)
				}
				Expect(err).NotTo(HaveOccurred())
				Expect(ready).To(Equal(true))
			}
		})
	})

	Describe("with a windows agent pool", func() {
		// TODO stabilize this test
		/*It("should be able to deploy an iis webserver", func() {
			if eng.HasWindowsAgents() {
				r := rand.New(rand.NewSource(time.Now().UnixNano()))
				deploymentName := fmt.Sprintf("iis-%s-%v", cfg.Name, r.Intn(99999))
				iisDeploy, err := deployment.CreateWindowsDeploy("microsoft/iis:windowsservercore-1709", deploymentName, "default", 80, -1)
				Expect(err).NotTo(HaveOccurred())

				running, err := pod.WaitOnReady(deploymentName, "default", 3, 30*time.Second, cfg.Timeout)
				Expect(err).NotTo(HaveOccurred())
				Expect(running).To(Equal(true))

				err = iisDeploy.Expose("LoadBalancer", 80, 80)
				Expect(err).NotTo(HaveOccurred())

				s, err := service.Get(deploymentName, "default")
				Expect(err).NotTo(HaveOccurred())

				valid := s.Validate("(IIS Windows Server)", 10, 10*time.Second, cfg.Timeout)
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
				err = s.Delete()
				Expect(err).NotTo(HaveOccurred())
			} else {
				Skip("No windows agent was provisioned for this Cluster Definition")
			}
		})*/

		// TODO stabilize this test
		/*It("should be able to reach hostport in an iis webserver", func() {
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
		})*/

		// TODO stabilize this test
		/*It("should be able to attach azure file", func() {
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
		})*/
	})
})

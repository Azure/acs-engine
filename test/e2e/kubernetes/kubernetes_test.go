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
	"github.com/Azure/acs-engine/test/e2e/kubernetes/networkpolicy"
	"github.com/Azure/acs-engine/test/e2e/kubernetes/node"
	"github.com/Azure/acs-engine/test/e2e/kubernetes/persistentvolumeclaims"
	"github.com/Azure/acs-engine/test/e2e/kubernetes/pod"
	"github.com/Azure/acs-engine/test/e2e/kubernetes/service"
	"github.com/Azure/acs-engine/test/e2e/kubernetes/storageclass"
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

		It("should have functional DNS", func() {
			if !eng.HasWindowsAgents() {
				if !eng.HasNetworkPolicy("calico") {
					var err error
					var p *pod.Pod
					p, err = pod.CreatePodFromFile(filepath.Join(WorkloadDir, "dns-liveness.yaml"), "dns-liveness", "default")
					if cfg.SoakClusterName == "" {
						Expect(err).NotTo(HaveOccurred())
					} else {
						if err != nil {
							p, err = pod.Get("dns-liveness", "default")
							Expect(err).NotTo(HaveOccurred())
						}
					}
					running, err := p.WaitOnReady(5*time.Second, 2*time.Minute)
					Expect(err).NotTo(HaveOccurred())
					Expect(running).To(Equal(true))
				}

				kubeConfig, err := GetConfig()
				Expect(err).NotTo(HaveOccurred())
				master := fmt.Sprintf("azureuser@%s", kubeConfig.GetServerName())
				sshKeyPath := cfg.GetSSHKeyPath()

				ifconfigCmd := fmt.Sprintf("ifconfig -a -v")
				cmd := exec.Command("ssh", "-i", sshKeyPath, "-o", "ConnectTimeout=10", "-o", "StrictHostKeyChecking=no", "-o", "UserKnownHostsFile=/dev/null", master, ifconfigCmd)
				util.PrintCommand(cmd)
				out, err := cmd.CombinedOutput()
				log.Printf("%s\n", out)
				if err != nil {
					log.Printf("Error while querying DNS: %s\n", out)
				}

				resolvCmd := fmt.Sprintf("cat /etc/resolv.conf")
				cmd = exec.Command("ssh", "-i", sshKeyPath, "-o", "ConnectTimeout=10", "-o", "StrictHostKeyChecking=no", "-o", "UserKnownHostsFile=/dev/null", master, resolvCmd)
				util.PrintCommand(cmd)
				out, err = cmd.CombinedOutput()
				log.Printf("%s\n", out)
				if err != nil {
					log.Printf("Error while querying DNS: %s\n", out)
				}

				By("Ensuring that we have a valid connection to our resolver")
				digCmd := fmt.Sprintf("dig +short +search +answer `hostname`")
				cmd = exec.Command("ssh", "-i", sshKeyPath, "-o", "ConnectTimeout=10", "-o", "StrictHostKeyChecking=no", "-o", "UserKnownHostsFile=/dev/null", master, digCmd)
				util.PrintCommand(cmd)
				out, err = cmd.CombinedOutput()
				if err != nil {
					log.Printf("Error while querying DNS: %s\n", out)
				}

				nodeList, err := node.Get()
				Expect(err).NotTo(HaveOccurred())
				for _, node := range nodeList.Nodes {
					By("Ensuring that we get a DNS lookup answer response for each node hostname")
					digCmd := fmt.Sprintf("dig +short +search +answer %s | grep -v -e '^$'", node.Metadata.Name)
					cmd = exec.Command("ssh", "-i", sshKeyPath, "-o", "ConnectTimeout=10", "-o", "StrictHostKeyChecking=no", "-o", "UserKnownHostsFile=/dev/null", master, digCmd)
					util.PrintCommand(cmd)
					out, err = cmd.CombinedOutput()
					if err != nil {
						log.Printf("Error while querying DNS: %s\n", out)
					}
					Expect(err).NotTo(HaveOccurred())
				}

				By("Ensuring that we get a DNS lookup answer response for external names")
				digCmd = fmt.Sprintf("dig +short +search www.bing.com | grep -v -e '^$'")
				cmd = exec.Command("ssh", "-i", sshKeyPath, "-o", "ConnectTimeout=10", "-o", "StrictHostKeyChecking=no", "-o", "UserKnownHostsFile=/dev/null", master, digCmd)
				util.PrintCommand(cmd)
				out, err = cmd.CombinedOutput()
				if err != nil {
					log.Printf("Error while querying DNS: %s\n", out)
				}
				digCmd = fmt.Sprintf("dig +short +search google.com | grep -v -e '^$'")
				cmd = exec.Command("ssh", "-i", sshKeyPath, "-o", "ConnectTimeout=10", "-o", "StrictHostKeyChecking=no", "-o", "UserKnownHostsFile=/dev/null", master, digCmd)
				util.PrintCommand(cmd)
				out, err = cmd.CombinedOutput()
				if err != nil {
					log.Printf("Error while querying DNS: %s\n", out)
				}

				By("Ensuring that we get a DNS lookup answer response for external names using external resolver")
				digCmd = fmt.Sprintf("dig +short +search www.bing.com @8.8.8.8 | grep -v -e '^$'")
				cmd = exec.Command("ssh", "-i", sshKeyPath, "-o", "ConnectTimeout=10", "-o", "StrictHostKeyChecking=no", "-o", "UserKnownHostsFile=/dev/null", master, digCmd)
				util.PrintCommand(cmd)
				out, err = cmd.CombinedOutput()
				if err != nil {
					log.Printf("Error while querying DNS: %s\n", out)
				}
				digCmd = fmt.Sprintf("dig +short +search google.com @8.8.8.8 | grep -v -e '^$'")
				cmd = exec.Command("ssh", "-i", sshKeyPath, "-o", "ConnectTimeout=10", "-o", "StrictHostKeyChecking=no", "-o", "UserKnownHostsFile=/dev/null", master, digCmd)
				util.PrintCommand(cmd)
				out, err = cmd.CombinedOutput()
				if err != nil {
					log.Printf("Error while querying DNS: %s\n", out)
				}

				j, err := job.CreateJobFromFile(filepath.Join(WorkloadDir, "validate-dns.yaml"), "validate-dns", "default")
				Expect(err).NotTo(HaveOccurred())
				ready, err := j.WaitOnReady(5*time.Second, cfg.Timeout)
				delErr := j.Delete()
				if delErr != nil {
					fmt.Printf("could not delete job %s\n", j.Metadata.Name)
					fmt.Println(delErr)
				}
				Expect(err).NotTo(HaveOccurred())
				Expect(ready).To(Equal(true))
			}
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
					dashboardPort := 443
					version, err := node.Version()
					Expect(err).NotTo(HaveOccurred())
					re := regexp.MustCompile("1.(5|6|7|8).")
					if re.FindString(version) != "" {
						dashboardPort = 80
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

		It("should have cluster-autoscaler running", func() {
			if hasClusterAutoscaler, clusterAutoscalerAddon := eng.HasAddon("autoscaler"); hasClusterAutoscaler {
				running, err := pod.WaitOnReady("cluster-autoscaler", "kube-system", 3, 30*time.Second, cfg.Timeout)
				Expect(err).NotTo(HaveOccurred())
				Expect(running).To(Equal(true))
				By("Ensuring that the correct resources have been applied")
				pods, err := pod.GetAllByPrefix("cluster-autoscaler", "kube-system")
				Expect(err).NotTo(HaveOccurred())
				for i, c := range clusterAutoscalerAddon.Containers {
					err := pods[0].Spec.Containers[i].ValidateResources(c)
					Expect(err).NotTo(HaveOccurred())
				}
			} else {
				Skip("cluster autoscaler disabled for this cluster, will not test")
			}
		})

		It("should have cluster-omsagent daemonset running", func() {
			if hasContainerMonitoring, clusterContainerMonitoringAddon := eng.HasAddon("container-monitoring"); hasContainerMonitoring {
				running, err := pod.WaitOnReady("omsagent-", "kube-system", 3, 30*time.Second, cfg.Timeout)
				Expect(err).NotTo(HaveOccurred())
				Expect(running).To(Equal(true))
				By("Ensuring that the correct resources have been applied")
				pods, err := pod.GetAllByPrefix("omsagent-", "kube-system")
				Expect(err).NotTo(HaveOccurred())
				for i, c := range clusterContainerMonitoringAddon.Containers {
					err := pods[0].Spec.Containers[i].ValidateResources(c)
					Expect(err).NotTo(HaveOccurred())
				}
			} else {
				Skip("container monitoring disabled for this cluster, will not test")
			}
		})

		It("should have cluster-omsagent replicaset running", func() {
			if hasContainerMonitoring, clusterContainerMonitoringAddon := eng.HasAddon("container-monitoring"); hasContainerMonitoring {
				running, err := pod.WaitOnReady("omsagent-rs", "kube-system", 3, 30*time.Second, cfg.Timeout)
				Expect(err).NotTo(HaveOccurred())
				Expect(running).To(Equal(true))
				By("Ensuring that the correct resources have been applied")
				pods, err := pod.GetAllByPrefix("omsagent-rs", "kube-system")
				Expect(err).NotTo(HaveOccurred())
				for i, c := range clusterContainerMonitoringAddon.Containers {
					err := pods[0].Spec.Containers[i].ValidateResources(c)
					Expect(err).NotTo(HaveOccurred())
				}
			} else {
				Skip("container monitoring disabled for this cluster, will not test")
			}
		})

		It("should be successfully running kubepodinventory plugin - ContainerMonitoring", func() {
			if hasContainerMonitoring, _ := eng.HasAddon("container-monitoring"); hasContainerMonitoring {
				running, err := pod.WaitOnReady("omsagent-rs", "kube-system", 3, 30*time.Second, cfg.Timeout)
				Expect(err).NotTo(HaveOccurred())
				Expect(running).To(Equal(true))
				By("Ensuring that the kubepodinventory plugin is writing data successfully")
				pods, err := pod.GetAllByPrefix("omsagent-rs", "kube-system")
				Expect(err).NotTo(HaveOccurred())
				_, err = pods[0].Exec("grep", "\"in_kube_podinventory::emit-stream : Success\"", "/var/opt/microsoft/omsagent/log/omsagent.log")
				Expect(err).NotTo(HaveOccurred())
			} else {
				Skip("container monitoring disabled for this cluster, will not test")
			}
		})

		It("should be successfully running kubenodeinventory plugin - ContainerMonitoring", func() {
			if hasContainerMonitoring, _ := eng.HasAddon("container-monitoring"); hasContainerMonitoring {
				running, err := pod.WaitOnReady("omsagent-rs", "kube-system", 3, 30*time.Second, cfg.Timeout)
				Expect(err).NotTo(HaveOccurred())
				Expect(running).To(Equal(true))
				By("Ensuring that the kubenodeinventory plugin is writing data successfully")
				pods, err := pod.GetAllByPrefix("omsagent-rs", "kube-system")
				Expect(err).NotTo(HaveOccurred())
				_, err = pods[0].Exec("grep", "\"in_kube_nodeinventory::emit-stream : Success\"", "/var/opt/microsoft/omsagent/log/omsagent.log")
				Expect(err).NotTo(HaveOccurred())
			} else {
				Skip("container monitoring disabled for this cluster, will not test")
			}
		})

		It("should be successfully running cadvisor_perf plugin - ContainerMonitoring", func() {
			if hasContainerMonitoring, _ := eng.HasAddon("container-monitoring"); hasContainerMonitoring {
				running, err := pod.WaitOnReady("omsagent-", "kube-system", 3, 30*time.Second, cfg.Timeout)
				Expect(err).NotTo(HaveOccurred())
				Expect(running).To(Equal(true))
				By("Ensuring that the cadvisor_perf plugin is writing data successfully")
				pods, err := pod.GetAllByPrefix("omsagent-", "kube-system")
				Expect(err).NotTo(HaveOccurred())
				_, err = pods[0].Exec("grep", "\"in_cadvisor_perf::emit-stream : Success\"", "/var/opt/microsoft/omsagent/log/omsagent.log")
				Expect(err).NotTo(HaveOccurred())
			} else {
				Skip("container monitoring disabled for this cluster, will not test")
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

		It("should have nvidia-device-plugin running", func() {
			if eng.HasGPUNodes() {
				if hasNVIDIADevicePlugin, NVIDIADevicePluginAddon := eng.HasAddon("nvidia-device-plugin"); hasNVIDIADevicePlugin {
					running, err := pod.WaitOnReady("nvidia-device-plugin", "kube-system", 3, 30*time.Second, cfg.Timeout)
					Expect(err).NotTo(HaveOccurred())
					Expect(running).To(Equal(true))
					pods, err := pod.GetAllByPrefix("nvidia-device-plugin", "kube-system")
					Expect(err).NotTo(HaveOccurred())
					for i, c := range NVIDIADevicePluginAddon.Containers {
						err := pods[0].Spec.Containers[i].ValidateResources(c)
						Expect(err).NotTo(HaveOccurred())
					}
				} else {
					Skip("nvidia-device-plugin disabled for this cluster, will not test")
				}
			}
		})
	})

	Describe("with a linux agent pool", func() {
		It("should be able to produce a working ILB connection", func() {
			if eng.HasLinuxAgents() {
				By("Creating a nginx deployment")
				r := rand.New(rand.NewSource(time.Now().UnixNano()))
				serviceName := "ingress-nginx"
				deploymentName := fmt.Sprintf("ingress-nginx-%s-%v", cfg.Name, r.Intn(99999))
				deploy, err := deployment.CreateLinuxDeploy("library/nginx:latest", deploymentName, "default", "--labels=app="+serviceName)
				Expect(err).NotTo(HaveOccurred())

				s, err := service.CreateServiceFromFile(filepath.Join(WorkloadDir, "ingress-nginx-ilb.yaml"), serviceName, "default")
				Expect(err).NotTo(HaveOccurred())
				svc, err := s.WaitForExternalIP(cfg.Timeout, 5*time.Second)
				Expect(err).NotTo(HaveOccurred())

				By("Ensuring the ILB IP is assigned to the service")
				curlDeploymentName := fmt.Sprintf("long-running-pod-%s-%v", cfg.Name, r.Intn(99999))
				curlDeploy, err := deployment.CreateLinuxDeploy("library/nginx:latest", curlDeploymentName, "default", "")
				Expect(err).NotTo(HaveOccurred())
				running, err := pod.WaitOnReady(curlDeploymentName, "default", 3, 30*time.Second, cfg.Timeout)
				Expect(err).NotTo(HaveOccurred())
				Expect(running).To(Equal(true))
				curlPods, err := curlDeploy.Pods()
				Expect(err).NotTo(HaveOccurred())
				for i, curlPod := range curlPods {
					if i < 1 {
						pass, err := curlPod.ValidateCurlConnection(svc.Status.LoadBalancer.Ingress[0]["ip"], 5*time.Second, cfg.Timeout)
						Expect(err).NotTo(HaveOccurred())
						Expect(pass).To(BeTrue())
					}
				}
				By("Cleaning up after ourselves")
				err = curlDeploy.Delete()
				Expect(err).NotTo(HaveOccurred())
				err = deploy.Delete()
				Expect(err).NotTo(HaveOccurred())
				err = s.Delete()
				Expect(err).NotTo(HaveOccurred())
			} else {
				Skip("No linux agent was provisioned for this Cluster Definition")
			}
		})

		It("should be able to autoscale", func() {
			if eng.HasLinuxAgents() {
				By("Creating a test php-apache deployment with request limit thresholds")
				// Inspired by http://blog.kubernetes.io/2016/07/autoscaling-in-kubernetes.html
				r := rand.New(rand.NewSource(time.Now().UnixNano()))
				phpApacheName := fmt.Sprintf("php-apache-%s-%v", cfg.Name, r.Intn(99999))
				phpApacheDeploy, err := deployment.CreateLinuxDeploy("k8s-gcrio.azureedge.net/hpa-example", phpApacheName, "default", "--requests=cpu=50m,memory=50M")
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

				By("Ensuring we have more than 1 apache-php pods due to hpa enforcement")
				_, err = phpApacheDeploy.WaitForReplicas(2, 5*time.Second, cfg.Timeout)
				Expect(err).NotTo(HaveOccurred())

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
				version := common.RationalizeReleaseAndVersion(
					common.Kubernetes,
					eng.ClusterDefinition.Properties.OrchestratorProfile.OrchestratorRelease,
					eng.ClusterDefinition.Properties.OrchestratorProfile.OrchestratorVersion,
					eng.HasWindowsAgents())
				if common.IsKubernetesVersionGe(version, "1.10.0") {
					j, err := job.CreateJobFromFile(filepath.Join(WorkloadDir, "cuda-vector-add.yaml"), "cuda-vector-add", "default")
					Expect(err).NotTo(HaveOccurred())
					ready, err := j.WaitOnReady(30*time.Second, cfg.Timeout)
					delErr := j.Delete()
					if delErr != nil {
						fmt.Printf("could not delete job %s\n", j.Metadata.Name)
						fmt.Println(delErr)
					}
					Expect(err).NotTo(HaveOccurred())
					Expect(ready).To(Equal(true))
				} else {
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
			}
		})
	})

	Describe("after the cluster has been up for awhile", func() {
		It("dns-liveness pod should not have any restarts", func() {
			if !eng.HasWindowsAgents() && !eng.HasNetworkPolicy("calico") {
				pod, err := pod.Get("dns-liveness", "default")
				Expect(err).NotTo(HaveOccurred())
				running, err := pod.WaitOnReady(5*time.Second, 3*time.Minute)
				Expect(err).NotTo(HaveOccurred())
				Expect(running).To(Equal(true))
				restarts := pod.Status.ContainerStatuses[0].RestartCount
				if cfg.SoakClusterName == "" {
					err = pod.Delete()
					Expect(err).NotTo(HaveOccurred())
					Expect(restarts).To(Equal(0))
				} else {
					log.Printf("%d DNS livenessProbe restarts since this cluster was created...\n", restarts)
				}
			}
		})
	})

	Describe("with calico network policy enabled", func() {
		It("should apply a network policy and deny outbound internet access to nginx pod", func() {
			if eng.HasNetworkPolicy("calico") || eng.HasNetworkPolicy("azure") {
				namespace := "default"
				By("Creating a nginx deployment")
				r := rand.New(rand.NewSource(time.Now().UnixNano()))
				deploymentName := fmt.Sprintf("nginx-%s-%v", cfg.Name, r.Intn(99999))
				nginxDeploy, err := deployment.CreateLinuxDeploy("library/nginx:latest", deploymentName, namespace, "")
				Expect(err).NotTo(HaveOccurred())

				By("Ensure there is a Running nginx pod")
				running, err := pod.WaitOnReady(deploymentName, namespace, 3, 30*time.Second, cfg.Timeout)
				Expect(err).NotTo(HaveOccurred())
				Expect(running).To(Equal(true))

				By("Ensuring we have outbound internet access from the nginx pods")
				nginxPods, err := nginxDeploy.Pods()
				Expect(err).NotTo(HaveOccurred())
				Expect(len(nginxPods)).ToNot(BeZero())
				for _, nginxPod := range nginxPods {
					pass, err := nginxPod.CheckLinuxOutboundConnection(5*time.Second, cfg.Timeout)
					Expect(err).NotTo(HaveOccurred())
					Expect(pass).To(BeTrue())
				}

				By("Applying a network policy to deny egress access")
				networkPolicyName := "calico-policy"
				err = networkpolicy.CreateNetworkPolicyFromFile(filepath.Join(WorkloadDir, "calico-policy.yaml"), networkPolicyName, namespace)
				Expect(err).NotTo(HaveOccurred())

				By("Ensuring we no longer have outbound internet access from the nginx pods")
				for _, nginxPod := range nginxPods {
					pass, err := nginxPod.CheckLinuxOutboundConnection(5*time.Second, 3*time.Minute)
					Expect(err).Should(HaveOccurred())
					Expect(pass).To(BeFalse())
				}

				By("Cleaning up after ourselves")
				networkpolicy.DeleteNetworkPolicy(networkPolicyName, namespace)
				// TODO delete networkpolicy
				// Expect(err).NotTo(HaveOccurred())
				err = nginxDeploy.Delete()
				Expect(err).NotTo(HaveOccurred())
			} else {
				Skip("Calico network policy was not provisioned for this Cluster Definition")
			}
		})
	})

	Describe("with a windows agent pool", func() {
		It("should be able to deploy an iis webserver", func() {
			if eng.HasWindowsAgents() {
				iisImage := "microsoft/iis:windowsservercore-1803" // BUG: This should be set based on the host OS version

				By("Creating a deployment with 1 pod running IIS")
				r := rand.New(rand.NewSource(time.Now().UnixNano()))
				deploymentName := fmt.Sprintf("iis-%s-%v", cfg.Name, r.Intn(99999))
				iisDeploy, err := deployment.CreateWindowsDeploy(iisImage, deploymentName, "default", 80, -1)
				Expect(err).NotTo(HaveOccurred())

				By("Waiting on pod to be Ready")
				running, err := pod.WaitOnReady(deploymentName, "default", 3, 30*time.Second, cfg.Timeout)
				Expect(err).NotTo(HaveOccurred())
				Expect(running).To(Equal(true))

				By("Exposing a LoadBalancer for the pod")
				err = iisDeploy.Expose("LoadBalancer", 80, 80)
				Expect(err).NotTo(HaveOccurred())
				s, err := service.Get(deploymentName, "default")
				Expect(err).NotTo(HaveOccurred())

				By("Verifying that the service is reachable and returns the default IIS start page")
				valid := s.Validate("(IIS Windows Server)", 10, 10*time.Second, cfg.Timeout)
				Expect(valid).To(BeTrue())

				By("Checking that each pod can reach http://www.bing.com")
				iisPods, err := iisDeploy.Pods()
				Expect(err).NotTo(HaveOccurred())
				Expect(len(iisPods)).ToNot(BeZero())
				for _, iisPod := range iisPods {
					pass, err := iisPod.CheckWindowsOutboundConnection(10*time.Second, cfg.Timeout)
					Expect(err).NotTo(HaveOccurred())
					Expect(pass).To(BeTrue())
				}

				By("Verifying pods & services can be deleted")
				err = iisDeploy.Delete()
				Expect(err).NotTo(HaveOccurred())
				err = s.Delete()
				Expect(err).NotTo(HaveOccurred())
			} else {
				Skip("No windows agent was provisioned for this Cluster Definition")
			}
		})

		It("Should not have any unready or crashing pods right after deployment", func() {
			if eng.HasWindowsAgents() {
				By("Checking ready status of each pod in kube-system")
				pods, err := pod.GetAll("kube-system")
				Expect(err).NotTo(HaveOccurred())
				Expect(len(pods.Pods)).ToNot(BeZero())
				for _, currentPod := range pods.Pods {
					log.Printf("Checking %s", currentPod.Metadata.Name)
					Expect(currentPod.Status.ContainerStatuses[0].Ready).To(BeTrue())
					Expect(currentPod.Status.ContainerStatuses[0].RestartCount).To(BeNumerically("<", 3))
				}
			}
		})

		// Windows Bug 18213017: Kubernetes Hostport mappings don't work
		/*
			It("should be able to reach hostport in an iis webserver", func() {
				if eng.HasWindowsAgents() {
					r := rand.New(rand.NewSource(time.Now().UnixNano()))
					hostport := 8123
					deploymentName := fmt.Sprintf("iis-%s-%v", cfg.Name, r.Intn(99999))
					iisDeploy, err := deployment.CreateWindowsDeploy(iisImage, deploymentName, "default", 80, hostport)
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

		It("should be able to attach azure file", func() {
			if eng.HasWindowsAgents() {
				if common.IsKubernetesVersionGe(eng.ClusterDefinition.ContainerService.Properties.OrchestratorProfile.OrchestratorVersion, "1.11") {
					// Failure in 1.11+ - https://github.com/kubernetes/kubernetes/issues/65845
					Skip("Kubernetes 1.11 has a known issue creating Azure PersistentVolumeClaims")
				} else if common.IsKubernetesVersionGe(eng.ClusterDefinition.ContainerService.Properties.OrchestratorProfile.OrchestratorVersion, "1.8") {
					By("Creating an AzureFile storage class")
					storageclassName := "azurefile" // should be the same as in storageclass-azurefile.yaml
					sc, err := storageclass.CreateStorageClassFromFile(filepath.Join(WorkloadDir, "storageclass-azurefile.yaml"), storageclassName)
					Expect(err).NotTo(HaveOccurred())
					ready, err := sc.WaitOnReady(5*time.Second, cfg.Timeout)
					Expect(err).NotTo(HaveOccurred())
					Expect(ready).To(Equal(true))

					By("Creating a persistent volume claim")
					pvcName := "pvc-azurefile" // should be the same as in pvc-azurefile.yaml
					pvc, err := persistentvolumeclaims.CreatePersistentVolumeClaimsFromFile(filepath.Join(WorkloadDir, "pvc-azurefile.yaml"), pvcName, "default")
					Expect(err).NotTo(HaveOccurred())
					ready, err = pvc.WaitOnReady("default", 5*time.Second, cfg.Timeout)
					Expect(err).NotTo(HaveOccurred())
					Expect(ready).To(Equal(true))

					By("Launching an IIS pod using the volume claim")
					podName := "iis-azurefile"                                                                                 // should be the same as in iis-azurefile.yaml
					iisPod, err := pod.CreatePodFromFile(filepath.Join(WorkloadDir, "iis-azurefile.yaml"), podName, "default") // BUG: this should support OS versioning
					Expect(err).NotTo(HaveOccurred())
					ready, err = iisPod.WaitOnReady(5*time.Second, cfg.Timeout)
					Expect(err).NotTo(HaveOccurred())
					Expect(ready).To(Equal(true))

					By("Checking that the pod can access volume")
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

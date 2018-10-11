package kubernetes

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/Azure/acs-engine/pkg/api/common"
	"github.com/Azure/acs-engine/test/e2e/config"
	"github.com/Azure/acs-engine/test/e2e/engine"
	"github.com/Azure/acs-engine/test/e2e/kubernetes/deployment"
	"github.com/Azure/acs-engine/test/e2e/kubernetes/job"
	"github.com/Azure/acs-engine/test/e2e/kubernetes/namespace"
	"github.com/Azure/acs-engine/test/e2e/kubernetes/networkpolicy"
	"github.com/Azure/acs-engine/test/e2e/kubernetes/node"
	"github.com/Azure/acs-engine/test/e2e/kubernetes/persistentvolume"
	"github.com/Azure/acs-engine/test/e2e/kubernetes/persistentvolumeclaims"
	"github.com/Azure/acs-engine/test/e2e/kubernetes/pod"
	"github.com/Azure/acs-engine/test/e2e/kubernetes/service"
	"github.com/Azure/acs-engine/test/e2e/kubernetes/storageclass"
	"github.com/Azure/acs-engine/test/e2e/kubernetes/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	WorkloadDir           = "workloads"
	PolicyDir             = "workloads/policies"
	deleteResourceRetries = 10
)

var (
	cfg                         config.Config
	eng                         engine.Engine
	masterSSHPort               string
	masterSSHPrivateKeyFilepath string
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
	masterNodes, err := node.GetByPrefix("k8s-master")
	Expect(err).NotTo(HaveOccurred())
	masterName := masterNodes[0].Metadata.Name
	if strings.Contains(masterName, "vmss") {
		masterSSHPort = "50001"
	} else {
		masterSSHPort = "22"
	}
	masterSSHPrivateKeyFilepath = cfg.GetSSHKeyPath()
	// TODO
	// If no user-configurable stability iteration value is passed in, run stability tests once
	/*if cfg.StabilityIterations == 0 {
		cfg.StabilityIterations = 1
	}*/
})

var _ = Describe("Azure Container Cluster using the Kubernetes Orchestrator", func() {
	Describe("regardless of agent pool type", func() {
		It("should display the installed Ubuntu version on the master node", func() {
			kubeConfig, err := GetConfig()
			Expect(err).NotTo(HaveOccurred())
			master := fmt.Sprintf("azureuser@%s", kubeConfig.GetServerName())

			lsbReleaseCmd := fmt.Sprintf("lsb_release -a && uname -r")
			var cmd *exec.Cmd
			cmd = exec.Command("ssh", "-i", masterSSHPrivateKeyFilepath, "-p", masterSSHPort, "-o", "ConnectTimeout=10", "-o", "StrictHostKeyChecking=no", "-o", "UserKnownHostsFile=/dev/null", master, lsbReleaseCmd)
			util.PrintCommand(cmd)
			out, err := cmd.CombinedOutput()
			log.Printf("%s\n", out)
			if err != nil {
				log.Printf("Error while getting Ubuntu image version: %s\n", err)
			}

			kernelVerCmd := fmt.Sprintf("cat /proc/version")
			cmd = exec.Command("ssh", "-i", masterSSHPrivateKeyFilepath, "-p", masterSSHPort, "-o", "ConnectTimeout=10", "-o", "StrictHostKeyChecking=no", "-o", "UserKnownHostsFile=/dev/null", master, kernelVerCmd)
			util.PrintCommand(cmd)
			out, err = cmd.CombinedOutput()
			log.Printf("%s\n", out)
			if err != nil {
				log.Printf("Error while getting LinuxKernel version: %s\n", err)
			}
		})

		It("should have have the appropriate node count", func() {
			nodeList, err := node.Get()
			Expect(err).NotTo(HaveOccurred())
			Expect(len(nodeList.Nodes)).To(Equal(eng.NodeCount()))
		})

		It("should have DNS pod running", func() {
			var err error
			var running bool
			if common.IsKubernetesVersionGe(eng.ExpandedDefinition.Properties.OrchestratorProfile.OrchestratorVersion, "1.12.0") {
				By("Ensuring that coredns is running")
				running, err = pod.WaitOnReady("coredns", "kube-system", 3, 30*time.Second, cfg.Timeout)

			} else {
				By("Ensuring that kube-dns is running")
				running, err = pod.WaitOnReady("kube-dns", "kube-system", 3, 30*time.Second, cfg.Timeout)
			}
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

		It("should have ip-masq-agent running", func() {
			if hasIPMasqAgent, IPMasqAgentAddon := eng.HasAddon("ip-masq-agent"); hasIPMasqAgent {
				running, err := pod.WaitOnReady("azure-ip-masq-agent", "kube-system", 3, 30*time.Second, cfg.Timeout)
				Expect(err).NotTo(HaveOccurred())
				Expect(running).To(Equal(true))
				By("Ensuring that the correct resources have been applied")
				pods, err := pod.GetAllByPrefix("azure-ip-masq-agent", "kube-system")
				Expect(err).NotTo(HaveOccurred())
				for _, p := range pods {
					for i, c := range IPMasqAgentAddon.Containers {
						err := p.Spec.Containers[i].ValidateResources(c)
						Expect(err).NotTo(HaveOccurred())
					}
				}
			} else {
				Skip("ip-masq-agent disabled for this cluster, will not test")
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

		It("should have keyvault-flexvolume running", func() {
			if hasKeyVaultFlexVolume, KeyVaultFlexVolumeAddon := eng.HasAddon("keyvault-flexvolume"); hasKeyVaultFlexVolume {
				running, err := pod.WaitOnReady("keyvault-flexvolume", "kv", 3, 30*time.Second, cfg.Timeout)
				Expect(err).NotTo(HaveOccurred())
				Expect(running).To(Equal(true))
				By("Ensuring that the correct resources have been applied")
				pods, err := pod.GetAllByPrefix("keyvault-flexvolume", "kv")
				Expect(err).NotTo(HaveOccurred())
				for i, c := range KeyVaultFlexVolumeAddon.Containers {
					err := pods[0].Spec.Containers[i].ValidateResources(c)
					Expect(err).NotTo(HaveOccurred())
				}

			} else {
				Skip("keyvault-flexvolume disabled for this cluster, will not test")
			}
		})

		It("should have blobfuse-flexvolume running", func() {
			if hasBlobfuseFlexVolume, BlobfuseFlexVolumeAddon := eng.HasAddon("blobfuse-flexvolume"); hasBlobfuseFlexVolume {
				running, err := pod.WaitOnReady("blobfuse-flexvol-installer", "flex", 3, 30*time.Second, 2*time.Minute)
				Expect(err).NotTo(HaveOccurred())
				Expect(running).To(Equal(true))
				By("Ensuring that the correct resources have been applied")
				pods, err := pod.GetAllByPrefix("blobfuse-flexvol-installer", "flex")
				Expect(err).NotTo(HaveOccurred())
				for i, c := range BlobfuseFlexVolumeAddon.Containers {
					err := pods[0].Spec.Containers[i].ValidateResources(c)
					Expect(err).NotTo(HaveOccurred())
				}

			} else {
				Skip("blobfuse-flexvolume disabled for this cluster, will not test")
			}
		})

		It("should have smb-flexvolume running", func() {
			if hasSMBFlexVolume, SMBFlexVolumeAddon := eng.HasAddon("smb-flexvolume"); hasSMBFlexVolume {
				running, err := pod.WaitOnReady("smb-flexvol-installer", "flex", 3, 30*time.Second, 2*time.Minute)
				Expect(err).NotTo(HaveOccurred())
				Expect(running).To(Equal(true))
				By("Ensuring that the correct resources have been applied")
				pods, err := pod.GetAllByPrefix("smb-flexvol-installer", "flex")
				Expect(err).NotTo(HaveOccurred())
				for i, c := range SMBFlexVolumeAddon.Containers {
					err := pods[0].Spec.Containers[i].ValidateResources(c)
					Expect(err).NotTo(HaveOccurred())
				}

			} else {
				Skip("smb-flexvolume disabled for this cluster, will not test")
			}
		})

		It("should have cluster-omsagent daemonset running", func() {
			if hasContainerMonitoring, clusterContainerMonitoringAddon := eng.HasAddon("container-monitoring"); hasContainerMonitoring {
				running, err := pod.WaitOnReady("omsagent-", "kube-system", 3, 30*time.Second, cfg.Timeout)
				Expect(err).NotTo(HaveOccurred())
				Expect(running).To(Equal(true))
				By("Ensuring that the correct resources have been applied")
				pods, err := pod.GetAllByPrefix("omsagent", "kube-system")
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
				pass, err := pods[0].ValidateOmsAgentLogs("kubePodInventoryEmitStreamSuccess", 30*time.Second, cfg.Timeout)
				Expect(err).NotTo(HaveOccurred())
				Expect(pass).To(BeTrue())
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
				pass, err := pods[0].ValidateOmsAgentLogs("kubeNodeInventoryEmitStreamSuccess", 30*time.Second, cfg.Timeout)
				Expect(err).NotTo(HaveOccurred())
				Expect(pass).To(BeTrue())
			} else {
				Skip("container monitoring disabled for this cluster, will not test")
			}
		})

		It("should be successfully running cadvisor_perf plugin - ContainerMonitoring", func() {
			if hasContainerMonitoring, _ := eng.HasAddon("container-monitoring"); hasContainerMonitoring {
				running, err := pod.WaitOnReady("omsagent", "kube-system", 3, 30*time.Second, cfg.Timeout)
				Expect(err).NotTo(HaveOccurred())
				Expect(running).To(Equal(true))
				By("Ensuring that the cadvisor_perf plugin is writing data successfully")
				pods, err := pod.GetAllByPrefix("omsagent", "kube-system")
				Expect(err).NotTo(HaveOccurred())
				pass, err := pods[0].ValidateOmsAgentLogs("cAdvisorPerfEmitStreamSuccess", 30*time.Second, cfg.Timeout)
				Expect(err).NotTo(HaveOccurred())
				Expect(pass).To(BeTrue())
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

		It("should have stable external container networking", func() {
			name := fmt.Sprintf("alpine-%s", cfg.Name)
			command := fmt.Sprintf("nc -vz 8.8.8.8 53 || nc -vz 8.8.4.4 53")
			successes, err := pod.RunCommandMultipleTimes(pod.RunLinuxPod, "alpine", name, command, cfg.StabilityIterations)
			Expect(err).NotTo(HaveOccurred())
			Expect(successes).To(Equal(cfg.StabilityIterations))
		})

		It("should have stable internal container networking", func() {
			name := fmt.Sprintf("alpine-%s", cfg.Name)
			var command string
			if common.IsKubernetesVersionGe(eng.ExpandedDefinition.Properties.OrchestratorProfile.OrchestratorVersion, "1.12.0") {
				command = fmt.Sprintf("nc -vz kubernetes 443 && nc -vz kubernetes.default.svc 443 && nc -vz kubernetes.default.svc.cluster.local 443")
			} else {
				command = fmt.Sprintf("nc -vz kubernetes 443")
			}
			successes, err := pod.RunCommandMultipleTimes(pod.RunLinuxPod, "alpine", name, command, cfg.StabilityIterations)
			Expect(err).NotTo(HaveOccurred())
			Expect(successes).To(Equal(cfg.StabilityIterations))
		})

		It("should be able to launch a long-running container networking DNS liveness pod", func() {
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
		})

		It("should have functional host OS DNS", func() {
			kubeConfig, err := GetConfig()
			Expect(err).NotTo(HaveOccurred())
			master := fmt.Sprintf("azureuser@%s", kubeConfig.GetServerName())

			ifconfigCmd := fmt.Sprintf("ifconfig -a -v")
			var cmd *exec.Cmd
			cmd = exec.Command("ssh", "-i", masterSSHPrivateKeyFilepath, "-p", masterSSHPort, "-o", "ConnectTimeout=10", "-o", "StrictHostKeyChecking=no", "-o", "UserKnownHostsFile=/dev/null", master, ifconfigCmd)
			util.PrintCommand(cmd)
			out, err := cmd.CombinedOutput()
			log.Printf("%s\n", out)
			if err != nil {
				log.Printf("Error while querying DNS: %s\n", err)
			}

			resolvCmd := fmt.Sprintf("cat /etc/resolv.conf")
			cmd = exec.Command("ssh", "-i", masterSSHPrivateKeyFilepath, "-p", masterSSHPort, "-o", "ConnectTimeout=10", "-o", "StrictHostKeyChecking=no", "-o", "UserKnownHostsFile=/dev/null", master, resolvCmd)
			util.PrintCommand(cmd)
			out, err = cmd.CombinedOutput()
			log.Printf("%s\n", out)
			if err != nil {
				log.Printf("Error while querying DNS: %s\n", err)
			}

			By("Ensuring that we have a valid connection to our resolver")
			digCmd := fmt.Sprintf("dig +short +search +answer `hostname`")
			cmd = exec.Command("ssh", "-i", masterSSHPrivateKeyFilepath, "-p", masterSSHPort, "-o", "ConnectTimeout=10", "-o", "StrictHostKeyChecking=no", "-o", "UserKnownHostsFile=/dev/null", master, digCmd)
			util.PrintCommand(cmd)
			out, err = cmd.CombinedOutput()
			log.Printf("%s\n", out)
			if err != nil {
				log.Printf("Error while querying DNS: %s\n", err)
			}

			nodeList, err := node.Get()
			Expect(err).NotTo(HaveOccurred())
			for _, node := range nodeList.Nodes {
				By("Ensuring that we get a DNS lookup answer response for each node hostname")
				digCmd := fmt.Sprintf("dig +short +search +answer %s | grep -v -e '^$'", node.Metadata.Name)

				cmd = exec.Command("ssh", "-i", masterSSHPrivateKeyFilepath, "-p", masterSSHPort, "-o", "ConnectTimeout=10", "-o", "StrictHostKeyChecking=no", "-o", "UserKnownHostsFile=/dev/null", master, digCmd)
				util.PrintCommand(cmd)
				out, err = cmd.CombinedOutput()
				log.Printf("%s\n", out)
				if err != nil {
					log.Printf("Error while querying DNS: %s\n", err)
				}
				Expect(err).NotTo(HaveOccurred())
			}

			By("Ensuring that we get a DNS lookup answer response for external names")
			digCmd = fmt.Sprintf("dig +short +search www.bing.com | grep -v -e '^$'")
			cmd = exec.Command("ssh", "-i", masterSSHPrivateKeyFilepath, "-p", masterSSHPort, "-o", "ConnectTimeout=10", "-o", "StrictHostKeyChecking=no", "-o", "UserKnownHostsFile=/dev/null", master, digCmd)
			util.PrintCommand(cmd)
			out, err = cmd.CombinedOutput()
			if err != nil {
				log.Printf("Error while querying DNS: %s\n", out)
			}
			digCmd = fmt.Sprintf("dig +short +search google.com | grep -v -e '^$'")
			cmd = exec.Command("ssh", "-i", masterSSHPrivateKeyFilepath, "-p", masterSSHPort, "-o", "ConnectTimeout=10", "-o", "StrictHostKeyChecking=no", "-o", "UserKnownHostsFile=/dev/null", master, digCmd)
			util.PrintCommand(cmd)
			out, err = cmd.CombinedOutput()
			log.Printf("%s\n", out)
			if err != nil {
				log.Printf("Error while querying DNS: %s\n", err)
			}

			By("Ensuring that we get a DNS lookup answer response for external names using external resolver")
			digCmd = fmt.Sprintf("dig +short +search www.bing.com @8.8.8.8 | grep -v -e '^$'")
			cmd = exec.Command("ssh", "-i", masterSSHPrivateKeyFilepath, "-p", masterSSHPort, "-o", "ConnectTimeout=10", "-o", "StrictHostKeyChecking=no", "-o", "UserKnownHostsFile=/dev/null", master, digCmd)
			util.PrintCommand(cmd)
			out, err = cmd.CombinedOutput()
			log.Printf("%s\n", out)
			if err != nil {
				log.Printf("Error while querying DNS: %s\n", err)
			}
			digCmd = fmt.Sprintf("dig +short +search google.com @8.8.8.8 | grep -v -e '^$'")
			cmd = exec.Command("ssh", "-i", masterSSHPrivateKeyFilepath, "-p", masterSSHPort, "-o", "ConnectTimeout=10", "-o", "StrictHostKeyChecking=no", "-o", "UserKnownHostsFile=/dev/null", master, digCmd)
			util.PrintCommand(cmd)
			out, err = cmd.CombinedOutput()
			log.Printf("%s\n", out)
			if err != nil {
				log.Printf("Error while querying DNS: %s\n", err)
			}
		})

		It("should have functional container networking DNS", func() {
			By("Ensuring that we have functional DNS resolution from a container")
			j, err := job.CreateJobFromFile(filepath.Join(WorkloadDir, "validate-dns.yaml"), "validate-dns", "default")
			Expect(err).NotTo(HaveOccurred())
			ready, err := j.WaitOnReady(5*time.Second, cfg.Timeout)
			delErr := j.Delete(deleteResourceRetries)
			if delErr != nil {
				fmt.Printf("could not delete job %s\n", j.Metadata.Name)
				fmt.Println(delErr)
			}
			Expect(err).NotTo(HaveOccurred())
			Expect(ready).To(Equal(true))

			By("Ensuring that we have stable external DNS resolution from a container")
			name := fmt.Sprintf("alpine-%s", cfg.Name)
			command := fmt.Sprintf("nc -vz bbc.co.uk 80 || nc -vz google.com 443 || nc -vz microsoft.com 80")
			successes, err := pod.RunCommandMultipleTimes(pod.RunLinuxPod, "alpine", name, command, cfg.StabilityIterations)
			Expect(err).NotTo(HaveOccurred())
			Expect(successes).To(Equal(cfg.StabilityIterations))
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
							address := node.Status.GetAddressByType("InternalIP")
							if address == nil {
								log.Printf("One of our nodes does not have an InternalIP value!: %s\n", node.Metadata.Name)
							}
							Expect(address).NotTo(BeNil())
							dashboardURL := fmt.Sprintf("http://%s:%v", address.Address, port)
							curlCMD := fmt.Sprintf("curl --max-time 60 %s", dashboardURL)
							var cmd *exec.Cmd
							cmd = exec.Command("ssh", "-i", masterSSHPrivateKeyFilepath, "-p", masterSSHPort, "-o", "ConnectTimeout=10", "-o", "StrictHostKeyChecking=no", "-o", "UserKnownHostsFile=/dev/null", master, curlCMD)
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
				err = curlDeploy.Delete(deleteResourceRetries)
				Expect(err).NotTo(HaveOccurred())
				err = deploy.Delete(deleteResourceRetries)
				Expect(err).NotTo(HaveOccurred())
				err = s.Delete(deleteResourceRetries)
				Expect(err).NotTo(HaveOccurred())
			} else {
				Skip("No linux agent was provisioned for this Cluster Definition")
			}
		})

		It("should be able to autoscale", func() {
			if eng.HasLinuxAgents() && eng.ExpandedDefinition.Properties.OrchestratorProfile.KubernetesConfig.EnableAggregatedAPIs {
				By("Creating a test php-apache deployment with request limit thresholds")
				// Inspired by http://blog.kubernetes.io/2016/07/autoscaling-in-kubernetes.html
				r := rand.New(rand.NewSource(time.Now().UnixNano()))
				phpApacheName := fmt.Sprintf("php-apache-%s-%v", cfg.Name, r.Intn(99999))
				phpApacheDeploy, err := deployment.CreateLinuxDeploy("k8s.gcr.io/hpa-example", phpApacheName, "default", "--requests=cpu=10m,memory=10M")
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
				err = loadTestDeploy.Delete(deleteResourceRetries)
				Expect(err).NotTo(HaveOccurred())
				err = phpApacheDeploy.Delete(deleteResourceRetries)
				Expect(err).NotTo(HaveOccurred())
				err = s.Delete(deleteResourceRetries)
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
				err = nginxDeploy.Delete(deleteResourceRetries)
				Expect(err).NotTo(HaveOccurred())
				err = s.Delete(deleteResourceRetries)
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
					false,
					eng.HasWindowsAgents())
				if common.IsKubernetesVersionGe(version, "1.10.0") {
					j, err := job.CreateJobFromFile(filepath.Join(WorkloadDir, "cuda-vector-add.yaml"), "cuda-vector-add", "default")
					Expect(err).NotTo(HaveOccurred())
					ready, err := j.WaitOnReady(30*time.Second, cfg.Timeout)
					delErr := j.Delete(deleteResourceRetries)
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
					delErr := j.Delete(deleteResourceRetries)
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

	Describe("with zoned master profile", func() {
		It("should be labeled with zones for each masternode", func() {
			if eng.ExpandedDefinition.Properties.MasterProfile.HasAvailabilityZones() {
				nodeList, err := node.Get()
				Expect(err).NotTo(HaveOccurred())
				for _, node := range nodeList.Nodes {
					role := node.Metadata.Labels["kubernetes.io/role"]
					if role == "master" {
						By("Ensuring that we get zones for each master node")
						zones := node.Metadata.Labels["failure-domain.beta.kubernetes.io/zone"]
						contains := strings.Contains(zones, "-")
						Expect(contains).To(Equal(true))
					}
				}
			} else {
				Skip("Availability zones was not configured for master profile for this Cluster Definition")
			}
		})
	})

	Describe("with all zoned agent pools", func() {
		It("should be labeled with zones for each node", func() {
			if eng.ExpandedDefinition.Properties.HasZonesForAllAgentPools() {
				nodeList, err := node.Get()
				Expect(err).NotTo(HaveOccurred())
				for _, node := range nodeList.Nodes {
					role := node.Metadata.Labels["kubernetes.io/role"]
					if role == "agent" {
						By("Ensuring that we get zones for each agent node")
						zones := node.Metadata.Labels["failure-domain.beta.kubernetes.io/zone"]
						contains := strings.Contains(zones, "-")
						Expect(contains).To(Equal(true))
					}
				}
			} else {
				Skip("Availability zones was not configured for this Cluster Definition")
			}
		})

		It("should create pv with zone labels and node affinity", func() {
			if eng.ExpandedDefinition.Properties.HasZonesForAllAgentPools() {
				By("Creating a persistent volume claim")
				pvcName := "azure-managed-disk" // should be the same as in pvc-premium.yaml
				pvc, err := persistentvolumeclaims.CreatePersistentVolumeClaimsFromFile(filepath.Join(WorkloadDir, "pvc-premium.yaml"), pvcName, "default")
				Expect(err).NotTo(HaveOccurred())
				ready, err := pvc.WaitOnReady("default", 5*time.Second, cfg.Timeout)
				Expect(err).NotTo(HaveOccurred())
				Expect(ready).To(Equal(true))

				pvList, err := persistentvolume.Get()
				Expect(err).NotTo(HaveOccurred())
				pvZone := ""
				for _, pv := range pvList.PersistentVolumes {
					By("Ensuring that we get zones for the pv")
					// zone is chosen by round-robin across all zones
					pvZone = pv.Metadata.Labels["failure-domain.beta.kubernetes.io/zone"]
					fmt.Printf("pvZone: %s\n", pvZone)
					contains := strings.Contains(pvZone, "-")
					Expect(contains).To(Equal(true))
					// VolumeScheduling feature gate is set to true by default starting v1.10+
					for _, expression := range pv.Spec.NodeAffinity.Required.NodeSelectorTerms[0].MatchExpressions {
						if expression.Key == "failure-domain.beta.kubernetes.io/zone" {
							By("Ensuring that we get nodeAffinity for each pv")
							value := expression.Values[0]
							fmt.Printf("NodeAffinity value: %s\n", value)
							contains := strings.Contains(value, "-")
							Expect(contains).To(Equal(true))
						}
					}
				}

				By("Launching a pod using the volume claim")
				podName := "zone-pv-pod" // should be the same as in pod-pvc.yaml
				testPod, err := pod.CreatePodFromFile(filepath.Join(WorkloadDir, "pod-pvc.yaml"), podName, "default")
				Expect(err).NotTo(HaveOccurred())
				ready, err = testPod.WaitOnReady(5*time.Second, cfg.Timeout)
				Expect(err).NotTo(HaveOccurred())
				Expect(ready).To(Equal(true))

				By("Checking that the pod can access volume")
				valid, err := testPod.ValidatePVC("/mnt/azure", 10, 10*time.Second)
				Expect(valid).To(BeTrue())
				Expect(err).NotTo(HaveOccurred())

				By("Ensuring that attached volume pv has the same zone as the zone of the node")
				nodeName := testPod.Spec.NodeName
				nodeList, err := node.GetByPrefix(nodeName)
				Expect(err).NotTo(HaveOccurred())
				nodeZone := nodeList[0].Metadata.Labels["failure-domain.beta.kubernetes.io/zone"]
				fmt.Printf("pvZone: %s\n", pvZone)
				fmt.Printf("nodeZone: %s\n", nodeZone)
				Expect(nodeZone == pvZone).To(Equal(true))

				By("Cleaning up after ourselves")
				err = testPod.Delete(deleteResourceRetries)
				Expect(err).NotTo(HaveOccurred())
				err = pvc.Delete(deleteResourceRetries)
				Expect(err).NotTo(HaveOccurred())
			} else {
				Skip("Availability zones was not configured for this Cluster Definition")
			}
		})
	})

	Describe("after the cluster has been up for awhile", func() {
		It("dns-liveness pod should not have any restarts", func() {
			if !eng.HasNetworkPolicy("calico") {
				pod, err := pod.Get("dns-liveness", "default")
				Expect(err).NotTo(HaveOccurred())
				running, err := pod.WaitOnReady(5*time.Second, 3*time.Minute)
				Expect(err).NotTo(HaveOccurred())
				Expect(running).To(Equal(true))
				restarts := pod.Status.ContainerStatuses[0].RestartCount
				if cfg.SoakClusterName == "" {
					err = pod.Delete(deleteResourceRetries)
					Expect(err).NotTo(HaveOccurred())
					Expect(restarts).To(Equal(0))
				} else {
					log.Printf("%d DNS livenessProbe restarts since this cluster was created...\n", restarts)
				}
			}
		})
	})

	Describe("with calico or azure network policy enabled", func() {
		It("should apply various network policies and enforce access to nginx pod", func() {
			if eng.HasNetworkPolicy("calico") || eng.HasNetworkPolicy("azure") {
				nsClientOne, nsClientTwo, nsServer := "client-one", "client-two", "server"
				By("Creating namespaces")
				_, err := namespace.Create(nsClientOne)
				Expect(err).NotTo(HaveOccurred())
				_, err = namespace.Create(nsClientTwo)
				Expect(err).NotTo(HaveOccurred())
				_, err = namespace.Create(nsServer)
				Expect(err).NotTo(HaveOccurred())
				By("Creating client and server nginx deployments")
				r := rand.New(rand.NewSource(time.Now().UnixNano()))
				randInt := r.Intn(99999)
				clientOneDeploymentName := fmt.Sprintf("nginx-%s-%v", cfg.Name, randInt)
				clientTwoDeploymentName := fmt.Sprintf("nginx-%s-%v", cfg.Name, randInt+100000)
				serverDeploymentName := fmt.Sprintf("nginx-%s-%v", cfg.Name, randInt+200000)
				clientOneDeploy, err := deployment.CreateLinuxDeploy("library/nginx:latest", clientOneDeploymentName, nsClientOne, "--labels=role=client-one")
				Expect(err).NotTo(HaveOccurred())
				clientTwoDeploy, err := deployment.CreateLinuxDeploy("library/nginx:latest", clientTwoDeploymentName, nsClientTwo, "--labels=role=client-two")
				Expect(err).NotTo(HaveOccurred())
				serverDeploy, err := deployment.CreateLinuxDeploy("library/nginx:latest", serverDeploymentName, nsServer, "--labels=role=server")
				Expect(err).NotTo(HaveOccurred())

				By("Ensure there is a Running nginx client one pod")
				running, err := pod.WaitOnReady(clientOneDeploymentName, nsClientOne, 3, 30*time.Second, cfg.Timeout)
				Expect(err).NotTo(HaveOccurred())
				Expect(running).To(Equal(true))

				By("Ensure there is a Running nginx client two pod")
				running, err = pod.WaitOnReady(clientTwoDeploymentName, nsClientTwo, 3, 30*time.Second, cfg.Timeout)
				Expect(err).NotTo(HaveOccurred())
				Expect(running).To(Equal(true))

				By("Ensure there is a Running nginx server pod")
				running, err = pod.WaitOnReady(serverDeploymentName, nsServer, 3, 30*time.Second, cfg.Timeout)
				Expect(err).NotTo(HaveOccurred())
				Expect(running).To(Equal(true))

				By("Ensuring we have outbound internet access from the nginx client one pods")
				clientOnePods, err := clientOneDeploy.Pods()
				Expect(err).NotTo(HaveOccurred())
				Expect(len(clientOnePods)).ToNot(BeZero())
				for _, clientOnePod := range clientOnePods {
					pass, err := clientOnePod.CheckLinuxOutboundConnection(5*time.Second, cfg.Timeout)
					Expect(err).NotTo(HaveOccurred())
					Expect(pass).To(BeTrue())
				}

				By("Ensuring we have outbound internet access from the nginx client one pods")
				clientTwoPods, err := clientTwoDeploy.Pods()
				Expect(err).NotTo(HaveOccurred())
				Expect(len(clientTwoPods)).ToNot(BeZero())
				for _, clientTwoPod := range clientTwoPods {
					pass, err := clientTwoPod.CheckLinuxOutboundConnection(5*time.Second, cfg.Timeout)
					Expect(err).NotTo(HaveOccurred())
					Expect(pass).To(BeTrue())
				}

				By("Ensuring we have outbound internet access from the nginx server pods")
				serverPods, err := serverDeploy.Pods()
				Expect(err).NotTo(HaveOccurred())
				Expect(len(serverPods)).ToNot(BeZero())
				for _, serverPod := range serverPods {
					pass, err := serverPod.CheckLinuxOutboundConnection(5*time.Second, cfg.Timeout)
					Expect(err).NotTo(HaveOccurred())
					Expect(pass).To(BeTrue())
				}

				var (
					networkPolicyName string
					namespace         string
				)

				By("Applying a network policy to deny egress access")
				networkPolicyName, namespace = "client-one-deny-egress", nsClientOne
				err = networkpolicy.CreateNetworkPolicyFromFile(filepath.Join(PolicyDir, "client-one-deny-egress-policy.yaml"), networkPolicyName, namespace)
				Expect(err).NotTo(HaveOccurred())

				By("Ensuring we no longer have outbound internet access from the nginx client pods")
				for _, clientOnePod := range clientOnePods {
					pass, err := clientOnePod.CheckLinuxOutboundConnection(5*time.Second, 3*time.Minute)
					Expect(err).Should(HaveOccurred())
					Expect(pass).To(BeFalse())
				}

				By("Cleaning up after ourselves")
				networkpolicy.DeleteNetworkPolicy(networkPolicyName, namespace)

				By("Applying a network policy to deny ingress access")
				networkPolicyName, namespace = "client-one-deny-ingress", nsServer
				err = networkpolicy.CreateNetworkPolicyFromFile(filepath.Join(PolicyDir, "client-one-deny-ingress-policy.yaml"), networkPolicyName, namespace)
				Expect(err).NotTo(HaveOccurred())

				By("Ensuring we no longer have inbound internet access from the nginx server pods")
				for _, clientOnePod := range clientOnePods {
					for _, serverPod := range serverPods {
						pass, err := clientOnePod.ValidateCurlConnection(serverPod.Status.PodIP, 5*time.Second, 3*time.Minute)
						Expect(err).Should(HaveOccurred())
						Expect(pass).To(BeFalse())
					}
				}

				By("Cleaning up after ourselves")
				networkpolicy.DeleteNetworkPolicy(networkPolicyName, namespace)
				err = clientOneDeploy.Delete(deleteResourceRetries)
				Expect(err).NotTo(HaveOccurred())
				err = clientTwoDeploy.Delete(deleteResourceRetries)
				Expect(err).NotTo(HaveOccurred())
				err = serverDeploy.Delete(deleteResourceRetries)
				Expect(err).NotTo(HaveOccurred())
			} else {
				Skip("Calico or Azure network policy was not provisioned for this Cluster Definition")
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
					pass, err := iisPod.CheckWindowsOutboundConnection("www.bing.com", 10*time.Second, cfg.Timeout)
					Expect(err).NotTo(HaveOccurred())
					Expect(pass).To(BeTrue())
				}

				By("Verifying pods & services can be deleted")
				err = iisDeploy.Delete(deleteResourceRetries)
				Expect(err).NotTo(HaveOccurred())
				err = s.Delete(deleteResourceRetries)
				Expect(err).NotTo(HaveOccurred())
			} else {
				Skip("No windows agent was provisioned for this Cluster Definition")
			}
		})

		It("should be able to resolve DNS across windows and linux deployments", func() {
			if eng.HasWindowsAgents() {
				iisImage := "microsoft/iis:windowsservercore-1803" // BUG: This should be set based on the host OS version
				windowsServerImage := "microsoft/windowsservercore:1803"

				By("Creating a deployment running IIS")
				r := rand.New(rand.NewSource(time.Now().UnixNano()))
				windowsDeploymentName := fmt.Sprintf("iis-dns-%s-%v", cfg.Name, r.Intn(99999))
				windowsIISDeployment, err := deployment.CreateWindowsDeploy(iisImage, windowsDeploymentName, "default", 80, -1)
				Expect(err).NotTo(HaveOccurred())

				By("Creating a nginx deployment")
				nginxDeploymentName := fmt.Sprintf("nginx-dns-%s-%v", cfg.Name, r.Intn(99999))
				linuxNginxDeploy, err := deployment.CreateLinuxDeploy("library/nginx:latest", nginxDeploymentName, "default", "")
				Expect(err).NotTo(HaveOccurred())

				By("Ensure there is a Running nginx pod")
				running, err := pod.WaitOnReady(nginxDeploymentName, "default", 3, 30*time.Second, cfg.Timeout)
				Expect(err).NotTo(HaveOccurred())
				Expect(running).To(Equal(true))

				By("Ensure there is a Running iis pod")
				running, err = pod.WaitOnReady(windowsDeploymentName, "default", 3, 30*time.Second, cfg.Timeout)
				Expect(err).NotTo(HaveOccurred())
				Expect(running).To(Equal(true))

				By("Exposing a internal service for the linux nginx deployment")
				err = linuxNginxDeploy.Expose("ClusterIP", 80, 80)
				Expect(err).NotTo(HaveOccurred())
				linuxService, err := service.Get(nginxDeploymentName, "default")
				Expect(err).NotTo(HaveOccurred())

				By("Exposing a internal service for the windows iis deployment")
				err = windowsIISDeployment.Expose("ClusterIP", 80, 80)
				Expect(err).NotTo(HaveOccurred())
				windowsService, err := service.Get(windowsDeploymentName, "default")
				Expect(err).NotTo(HaveOccurred())

				By("Connecting to Windows from another Windows deployment")
				name := fmt.Sprintf("windows-2-windows-%s", cfg.Name)
				command := fmt.Sprintf("iwr -UseBasicParsing -TimeoutSec 60 %s", windowsService.Metadata.Name)
				successes, err := pod.RunCommandMultipleTimes(pod.RunWindowsPod, windowsServerImage, name, command, cfg.StabilityIterations)
				Expect(err).NotTo(HaveOccurred())
				Expect(successes).To(Equal(cfg.StabilityIterations))

				By("Connecting to Linux from Windows deployment")
				name = fmt.Sprintf("windows-2-linux-%s", cfg.Name)
				command = fmt.Sprintf("iwr -UseBasicParsing -TimeoutSec 60 %s", linuxService.Metadata.Name)
				successes, err = pod.RunCommandMultipleTimes(pod.RunWindowsPod, windowsServerImage, name, command, cfg.StabilityIterations)
				Expect(err).NotTo(HaveOccurred())
				Expect(successes).To(Equal(cfg.StabilityIterations))

				By("Connecting to Windows from Linux deployment")
				name = fmt.Sprintf("linux-2-windows-%s", cfg.Name)
				command = fmt.Sprintf("wget %s", windowsService.Metadata.Name)
				successes, err = pod.RunCommandMultipleTimes(pod.RunLinuxPod, "alpine", name, command, cfg.StabilityIterations)
				Expect(err).NotTo(HaveOccurred())
				Expect(successes).To(Equal(cfg.StabilityIterations))

				By("Cleaning up after ourselves")
				err = windowsIISDeployment.Delete(deleteResourceRetries)
				Expect(err).NotTo(HaveOccurred())
				err = linuxNginxDeploy.Delete(deleteResourceRetries)
				Expect(err).NotTo(HaveOccurred())
				err = windowsService.Delete(deleteResourceRetries)
				Expect(err).NotTo(HaveOccurred())
				err = linuxService.Delete(deleteResourceRetries)
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

					for _, iisPod := range iisPods {
						valid := iisPod.ValidateHostPort("(IIS Windows Server)", 10, 10*time.Second, master, masterSSHPrivateKeyFilepath)
						Expect(valid).To(BeTrue())
					}

					err = iisDeploy.Delete(kubectlOutput)
					Expect(err).NotTo(HaveOccurred())
				} else {
					Skip("No windows agent was provisioned for this Cluster Definition")
				}
			})*/

		It("should be able to attach azure file", func() {
			if eng.HasWindowsAgents() {
				if eng.ExpandedDefinition.Properties.OrchestratorProfile.OrchestratorVersion == "1.11.0" {
					// Failure in 1.11.0 - https://github.com/kubernetes/kubernetes/issues/65845, fixed in 1.11.1
					Skip("Kubernetes 1.11.0 has a known issue creating Azure PersistentVolumeClaims")
				} else if common.IsKubernetesVersionGe(eng.ExpandedDefinition.Properties.OrchestratorProfile.OrchestratorVersion, "1.8.0") {
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

					err = iisPod.Delete(deleteResourceRetries)
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

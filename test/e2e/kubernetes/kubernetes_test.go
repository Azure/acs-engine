package kubernetes

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
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
			Expect(version).To(Equal("v1.6.6"))
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
			running, err := pod.WaitOnReady("heapster", "kube-system", 5*time.Second, 3*time.Minute)
			Expect(err).NotTo(HaveOccurred())
			Expect(running).To(Equal(true))
		})

		It("should have kube-addon-manager running", func() {
			running, err := pod.WaitOnReady("kube-addon-manager", "kube-system", 5*time.Second, 3*time.Minute)
			Expect(err).NotTo(HaveOccurred())
			Expect(running).To(Equal(true))
		})

		It("should have kube-apiserver running", func() {
			running, err := pod.WaitOnReady("kube-apiserver", "kube-system", 5*time.Second, 3*time.Minute)
			Expect(err).NotTo(HaveOccurred())
			Expect(running).To(Equal(true))
		})

		It("should have kube-controller-manager running", func() {
			running, err := pod.WaitOnReady("kube-controller-manager", "kube-system", 5*time.Second, 3*time.Minute)
			Expect(err).NotTo(HaveOccurred())
			Expect(running).To(Equal(true))
		})

		It("should have kube-scheduler running", func() {
			running, err := pod.WaitOnReady("kube-scheduler", "kube-system", 5*time.Second, 3*time.Minute)
			Expect(err).NotTo(HaveOccurred())
			Expect(running).To(Equal(true))
		})

		It("should have tiller running", func() {
			running, err := pod.WaitOnReady("tiller", "kube-system", 5*time.Second, 3*time.Minute)
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
				err = d.Expose(80)
				Expect(err).NotTo(HaveOccurred())

				s, err := service.Get(deploymentName, "default")
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
			} else {
				Skip("No linux agent was provisioned for this Cluster Definition")
			}
		})
	})

	Context("with a windows agent pool", func() {
		It("should be able to deploy a powershell webserver", func() {
			if eng.HasWindowsAgents() {
				command := `powershell.exe -command "<#code used from https://gist.github.com/wagnerandrade/5424431#> ; $$listener = New-Object System.Net.HttpListener ; $$listener.Prefixes.Add('http://*:80/') ; $$listener.Start() ; $$callerCounts = @{} ; Write-Host('Listening at http://*:80/') ; while ($$listener.IsListening) { ;$$context = $$listener.GetContext() ;$$requestUrl = $$context.Request.Url ;$$clientIP = $$context.Request.RemoteEndPoint.Address ;$$response = $$context.Response ;Write-Host '' ;Write-Host('> {0}' -f $$requestUrl) ;  ;$$count = 1 ;$$k=$$callerCounts.Get_Item($$clientIP) ;if ($$k -ne $$null) { $$count += $$k } ;$$callerCounts.Set_Item($$clientIP, $$count) ;$$header='<html><body><H1>Windows Container Web Server</H1>' ;$$callerCountsString='' ;$$callerCounts.Keys | % { $$callerCountsString+='<p>IP {0} callerCount {1} ' -f $$_,$$callerCounts.Item($$_) } ;$$footer='</body></html>' ;$$content='{0}{1}{2}' -f $$header,$$callerCountsString,$$footer ;Write-Output $$content ;$$buffer = [System.Text.Encoding]::UTF8.GetBytes($$content) ;$$response.ContentLength64 = $$buffer.Length ;$$response.OutputStream.Write($$buffer, 0, $$buffer.Length) ;$$response.Close() ;$$responseStatus = $$response.StatusCode ;Write-Host('< {0}' -f $$responseStatus)  } ; "`
				r := rand.New(rand.NewSource(time.Now().UnixNano()))
				deploymentName := fmt.Sprintf("%s-%v", cfg.Name, r.Intn(99999))
				d, err := deployment.CreateWindowsDeploy("microsoft/windowsservercore", deploymentName, "default", command)
				Expect(err).NotTo(HaveOccurred())
				err = d.Expose(80)
				Expect(err).NotTo(HaveOccurred())

				s, err := service.Get(deploymentName, "default")
				Expect(err).NotTo(HaveOccurred())
				s, err = s.WaitForExternalIP(360, 5)
				Expect(err).NotTo(HaveOccurred())
				Expect(s.Status.LoadBalancer.Ingress).NotTo(BeEmpty())
			} else {
				Skip("No windows agent was provisioned for this Cluster Definition")
			}
		})
	})
})

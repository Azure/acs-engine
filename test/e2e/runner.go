package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"time"

	"github.com/Azure/acs-engine/test/e2e/azure"
	"github.com/Azure/acs-engine/test/e2e/config"
	"github.com/Azure/acs-engine/test/e2e/dcos"
	"github.com/Azure/acs-engine/test/e2e/engine"
	"github.com/Azure/acs-engine/test/e2e/kubernetes/node"
	"github.com/Azure/acs-engine/test/e2e/metrics"
)

var (
	cfg  *config.Config
	acct *azure.Account
	eng  *engine.Engine
	rgs  []string
	err  error
	pt   *metrics.Point
)

func main() {
	cwd, _ := os.Getwd()
	cfg, err = config.ParseConfig()
	if err != nil {
		log.Fatalf("Error while trying to parse configuration: %s\n", err)
	}
	cfg.CurrentWorkingDir = cwd

	acct, err = azure.NewAccount()
	if err != nil {
		log.Fatalf("Error while trying to setup azure account: %s\n", err)
	}

	err := acct.Login()
	if err != nil {
		log.Fatal("Error while trying to login to azure account!")
	}

	err = acct.SetSubscription()
	if err != nil {
		log.Fatal("Error while trying to set azure subscription!")
	}
	pt = metrics.BuildPoint(cfg.Orchestrator, cfg.Location, cfg.ClusterDefinition, acct.SubscriptionID)

	// If an interrupt/kill signal is sent we will run the clean up procedure
	trap()

	// Only provision a cluster if there isnt a name present
	if cfg.Name == "" {
		for i := 1; i <= cfg.ProvisionRetries; i++ {
			pt.SetProvisionStart()
			success := provisionCluster()
			rgs = append(rgs, cfg.Name)
			if success {
				pt.RecordProvisionSuccess()
				break
			} else if i == cfg.ProvisionRetries {
				pt.RecordProvisionError()
				teardown()
				log.Fatalf("Exceeded Provision retry count!")
			}
			pt.RecordProvisionError()
		}
	} else {
		engCfg, err := engine.ParseConfig(cfg.CurrentWorkingDir, cfg.ClusterDefinition, cfg.Name)
		if err != nil {
			teardown()
			log.Fatalf("Error trying to parse Engine config:%s\n", err)
		}
		cs, err := engine.Parse(engCfg.ClusterDefinitionTemplate)
		if err != nil {
			teardown()
			log.Fatalf("Error trying to parse engine template into memory:%s\n", err)

		}
		eng = &engine.Engine{
			Config:            engCfg,
			ClusterDefinition: cs,
		}
	}

	pt.SetNodeWaitStart()
	if cfg.IsKubernetes() {
		os.Setenv("KUBECONFIG", cfg.GetKubeConfig())
		log.Printf("Kubeconfig:%s\n", cfg.GetKubeConfig())
		log.Println("Waiting on nodes to go into ready state...")
		ready := node.WaitOnReady(eng.NodeCount(), 10*time.Second, cfg.Timeout)
		if ready == false {
			pt.RecordNodeWait()
			teardown()
			log.Fatalf("Error: Not all nodes in ready state!")
		}
	}

	if cfg.IsDCOS() {
		host := fmt.Sprintf("%s.%s.cloudapp.azure.com", cfg.Name, cfg.Location)
		user := eng.ClusterDefinition.Properties.LinuxProfile.AdminUsername
		log.Printf("SSH Key: %s\n", cfg.GetSSHKeyPath())
		log.Printf("Master Node: %s@%s\n", user, host)
		log.Printf("SSH Command: ssh -i %s -p 2200 %s@%s", cfg.GetSSHKeyPath(), user, host)
		cluster := dcos.NewCluster(cfg, eng)
		err = cluster.InstallDCOSClient()
		if err != nil {
			teardown()
			log.Fatalf("Error trying to install dcos client:%s\n", err)
		}
		ready := cluster.WaitForNodes(eng.NodeCount(), 10*time.Second, cfg.Timeout)
		if ready == false {
			pt.RecordNodeWait()
			teardown()
			log.Fatal("Error: Not all nodes in healthy state!")
		}
	}
	pt.RecordNodeWait()

	if !cfg.SkipTest {
		pt.SetTestStart()
		err := runGinkgo(cfg.Orchestrator)
		if err != nil {
			pt.RecordTestError()
			teardown()
			os.Exit(1)
		} else {
			pt.RecordTestSuccess()
		}
	}

	teardown()
	os.Exit(0)
}

func trap() {
	// If an interrupt/kill signal is sent we will run the clean up procedure
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)
	go func() {
		for sig := range c {
			log.Printf("Received Signal:%s ... Clean Up On Exit?:%v\n", sig.String(), cfg.CleanUpOnExit)
			teardown()
			os.Exit(1)
		}
	}()
}

func teardown() {
	pt.RecordTotalTime()
	pt.Write()
	if cfg.CleanUpOnExit {
		for _, rg := range rgs {
			log.Printf("Deleting Group:%s\n", rg)
			acct.DeleteGroup(rg)
		}
	}
}

func runGinkgo(orchestrator string) error {
	testDir := fmt.Sprintf("test/e2e/%s", orchestrator)
	cmd := exec.Command("ginkgo", "-nodes", "10", "-slowSpecThreshold", "180", "-r", testDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Start()
	if err != nil {
		log.Printf("Error while trying to start ginkgo:%s\n", err)
		return err
	}

	err = cmd.Wait()
	if err != nil {
		return err
	}
	return nil
}

func provisionCluster() bool {
	cfg.Name = cfg.GenerateName()
	os.Setenv("NAME", cfg.Name)
	log.Printf("Cluster name:%s\n", cfg.Name)

	outputPath := filepath.Join(cfg.CurrentWorkingDir, "_output")
	os.RemoveAll(outputPath)
	os.Mkdir(outputPath, 0755)

	out, err := exec.Command("ssh-keygen", "-f", cfg.GetSSHKeyPath(), "-q", "-N", "", "-b", "2048", "-t", "rsa").CombinedOutput()
	if err != nil {
		log.Printf("Error while trying to generate ssh key:%s\n\nOutput:%s\n", err, out)
		return false
	}
	exec.Command("chmod", "0600", cfg.GetSSHKeyPath()+"*")

	publicSSHKey, err := cfg.ReadPublicSSHKey()
	if err != nil {
		log.Printf("Error while trying to read public ssh key: %s\n", err)
		return false
	}
	os.Setenv("PUBLIC_SSH_KEY", publicSSHKey)
	os.Setenv("DNS_PREFIX", cfg.Name)

	err = acct.CreateGroup(cfg.Name, cfg.Location)
	if err != nil {
		log.Printf("Error while trying to create resource group: %s\n", err)
		return false
	}

	subnetID := ""
	vnetName := fmt.Sprintf("%sCustomVnet", cfg.Name)
	subnetName := fmt.Sprintf("%sCustomSubnet", cfg.Name)
	if cfg.CreateVNET {
		err = acct.CreateVnet(vnetName, "10.239.0.0/16", subnetName, "10.239.0.0/16")
		if err != nil {
			return false
		}
		subnetID = fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/virtualNetworks/%s/subnets/%s", acct.SubscriptionID, acct.ResourceGroup.Name, vnetName, subnetName)
	}

	// Lets modify our template and call acs-engine generate on it
	eng, err = engine.Build(cfg, subnetID)
	if err != nil {
		log.Printf("Error while trying to build cluster definition: %s\n", err)
		return false
	}

	err = eng.Write()
	if err != nil {
		log.Printf("Error while trying to write Engine Template to disk:%s\n", err)
		return false
	}

	err = eng.Generate()
	if err != nil {
		log.Printf("Error while trying to generate acs-engine template: %s\n", err)
		return false
	}

	// Lets start by just using the normal az group deployment cli for creating a cluster
	log.Println("Creating deployment this make take a few minutes...")
	err = acct.CreateDeployment(cfg.Name, eng)
	if err != nil {
		return false
	}

	if cfg.CreateVNET {
		err = acct.UpdateRouteTables(subnetName, vnetName)
		if err != nil {
			return false
		}
	}

	return true
}

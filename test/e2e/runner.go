package main

import (
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"time"

	"github.com/Azure/acs-engine/test/e2e/azure"
	"github.com/Azure/acs-engine/test/e2e/config"
	"github.com/Azure/acs-engine/test/e2e/engine"
	"github.com/Azure/acs-engine/test/e2e/kubernetes/node"
)

var (
	cfg  *config.Config
	acct *azure.Account
	err  error
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

	acct.Login()
	acct.SetSubscription()

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

	if cfg.Name == "" {
		cfg.Name = cfg.GenerateName()
		log.Printf("Cluster name:%s\n", cfg.Name)

		outputPath := filepath.Join(cfg.CurrentWorkingDir, "_output")
		os.RemoveAll(outputPath)
		os.Mkdir(outputPath, 0755)

		out, err := exec.Command("ssh-keygen", "-f", cfg.GetSSHKeyPath(), "-q", "-N", "", "-b", "2048", "-t", "rsa").CombinedOutput()
		if err != nil {
			log.Printf("Error while trying to generate ssh key:%s\n", err)
			log.Printf("Output:%s\n", out)
			os.Exit(1)
		}
		exec.Command("chmod", "0600", cfg.GetSSHKeyPath()+"*")
		exec.Command("ssh-add", cfg.GetSSHKeyPath())

		publicSSHKey, err := cfg.ReadPublicSSHKey()
		if err != nil {
			log.Fatalf("Error while trying to read public ssh key: %s\n", err)
		}
		os.Setenv("PUBLIC_SSH_KEY", publicSSHKey)
		os.Setenv("DNS_PREFIX", cfg.Name)

		// Lets modify our template and call acs-engine generate on it
		e, err := engine.Build(cfg.CurrentWorkingDir, cfg.ClusterDefinition, "_output", cfg.Name)
		if err != nil {
			teardown()
			log.Fatalf("Error while trying to build cluster definition: %s\n", err)
		}

		err = e.Generate()
		if err != nil {
			teardown()
			log.Fatalf("Error while trying to generate acs-engine template: %s\n", err)
		}

		err = acct.CreateGroup(cfg.Name, cfg.Location)
		if err != nil {
			teardown()
			log.Fatalf("Error while trying to create resource group: %s\n", err)
		}

		// Lets start by just using the normal az group deployment cli for creating a cluster
		log.Println("Creating deployment this make take a few minutes...")
		err = acct.CreateDeployment(cfg.Name, e)
		if err != nil {
			teardown()
			log.Fatalf("Error while trying to create deployment: %s\n", err)
		}
	}

	os.Setenv("NAME", cfg.Name)
	os.Setenv("KUBECONFIG", cfg.GetKubeConfig())
	log.Printf("Kubeconfig:%s\n", cfg.GetKubeConfig())

	log.Println("Waiting on nodes to go into ready state...")
	ready := node.WaitOnReady(10*time.Second, 10*time.Minute)
	if ready == false {
		teardown()
		log.Fatalf("Error: Not all nodes in ready state!")
	}

	cmd := exec.Command("ginkgo", "-nodes", "10", "-slowSpecThreshold", "180", "-r", "test/e2e/", cfg.Orchestrator)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Start()
	if err != nil {
		log.Printf("Error while trying to start ginkgo:%s\n", err)
		teardown()
		os.Exit(1)
	}

	err = cmd.Wait()
	if err != nil {
		teardown()
		os.Exit(1)
	}
}

func teardown() {
	if cfg.CleanUpOnExit {
		log.Printf("Deleting Group:%s\n", cfg.Name)
		acct.DeleteGroup()
	}
	exec.Command("ssh-add", "-d", cfg.GetSSHKeyPath())
}

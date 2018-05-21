package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"time"

	"github.com/Azure/acs-engine/test/e2e/azure"
	"github.com/Azure/acs-engine/test/e2e/config"
	"github.com/Azure/acs-engine/test/e2e/engine"
	"github.com/Azure/acs-engine/test/e2e/metrics"
	"github.com/Azure/acs-engine/test/e2e/runner"
)

var (
	cfg            *config.Config
	acct           *azure.Account
	eng            *engine.Engine
	rgs            []string
	err            error
	pt             *metrics.Point
	cliProvisioner *runner.CLIProvisioner
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

	cliProvisioner, err = runner.BuildCLIProvisioner(cfg, acct, pt)
	if err != nil {
		log.Fatalf("Error while trying to build CLI Provisioner:%s", err)
	}

	sa := acct.StorageAccount

	// Soak test specific setup
	if cfg.SoakClusterName != "" {
		sa.Name = "acsesoaktests" + cfg.Location
		sa.ResourceGroup.Name = "acse-test-infrastructure-storage"
		sa.ResourceGroup.Location = cfg.Location
		err = sa.CreateStorageAccount()
		if err != nil {
			log.Fatalf("Error while trying to create storage account: %s\n", err)
		}
		err = sa.SetConnectionString()
		if err != nil {
			log.Fatalf("Error while trying to set storage account connection string: %s\n", err)
		}
		provision := true
		rg := cfg.SoakClusterName
		err = acct.SetResourceGroup(rg)
		if err != nil {
			log.Printf("Error while trying to set RG:%s\n", err)
		} else {
			// set expiration time to 7 days = 168h for now
			d, err := time.ParseDuration("168h")
			if err != nil {
				log.Fatalf("Unexpected error parsing duration: %s", err)
			}
			provision = acct.IsClusterExpired(d)
		}
		if provision || cfg.ForceDeploy {
			log.Printf("Soak cluster %s does not exist or has expired\n", rg)
			log.Printf("Deleting Resource Group:%s\n", rg)
			acct.DeleteGroup(rg, true)
			log.Printf("Deleting Storage files:%s\n", rg)
			sa.DeleteFiles(cfg.SoakClusterName)
			cfg.Name = ""
		} else {
			log.Printf("Soak cluster %s exists, downloading output files from storage...\n", rg)
			err = sa.DownloadFiles(cfg.SoakClusterName, "_output")
			if err != nil {
				log.Printf("Error while trying to download _output dir: %s, will provision a new cluster.\n", err)
				log.Printf("Deleting Resource Group:%s\n", rg)
				acct.DeleteGroup(rg, true)
				log.Printf("Deleting Storage files:%s\n", rg)
				sa.DeleteFiles(cfg.SoakClusterName)
				cfg.Name = ""
			} else {
				cfg.SetSSHKeyPermissions()
			}
		}
	}
	// Only provision a cluster if there isn't a name present
	if cfg.Name == "" {
		err = cliProvisioner.Run()
		rgs = cliProvisioner.ResourceGroups
		eng = cliProvisioner.Engine
		if err != nil {
			if cfg.CleanUpIfFail {
				teardown()
			}
			log.Fatalf("Error while trying to provision cluster:%s", err)
		}
		if cfg.SoakClusterName != "" {
			err = sa.CreateFileShare(cfg.SoakClusterName)
			if err != nil {
				log.Printf("Error while trying to create file share:%s\n", err)
			}
			err = sa.UploadFiles(filepath.Join(cfg.CurrentWorkingDir, "_output"), cfg.SoakClusterName)
			if err != nil {
				log.Fatalf("Error while trying to upload _output dir:%s\n", err)
			}
		}
	} else {
		engCfg, err := engine.ParseConfig(cfg.CurrentWorkingDir, cfg.ClusterDefinition, cfg.Name)
		cfg.SetKubeConfig()
		if err != nil {
			teardown()
			log.Fatalf("Error trying to parse Engine config:%s\n", err)
		}
		cs, err := engine.ParseInput(engCfg.ClusterDefinitionTemplate)
		if err != nil {
			teardown()
			log.Fatalf("Error trying to parse engine template into memory:%s\n", err)
		}
		eng = &engine.Engine{
			Config:            engCfg,
			ClusterDefinition: cs,
		}
		cliProvisioner.Engine = eng
	}

	if !cfg.SkipTest {
		g, err := runner.BuildGinkgoRunner(cfg, pt)
		if err != nil {
			teardown()
			log.Fatalf("Error: Unable to parse ginkgo configuration!")
		}
		err = g.Run()
		if err != nil {
			teardown()
			os.Exit(1)
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
	if cliProvisioner.Config.IsKubernetes() && cfg.SoakClusterName == "" {
		hostname := fmt.Sprintf("%s.%s.cloudapp.azure.com", cfg.Name, cfg.Location)
		logsPath := filepath.Join(cfg.CurrentWorkingDir, "_logs", hostname)
		err := os.MkdirAll(logsPath, 0755)
		if err != nil {
			log.Printf("cannot create directory for logs: %s", err)
		}
		err = cliProvisioner.FetchProvisioningMetrics(logsPath, cfg, acct)
		if err != nil {
			log.Printf("cliProvisioner.FetchProvisioningMetrics error: %s\n", err)
		}
	}
	if cliProvisioner.Config.IsOpenShift() {
		hostname := fmt.Sprintf("%s.%s.cloudapp.azure.com", cfg.Name, cfg.Location)
		logsPath := filepath.Join(cfg.CurrentWorkingDir, "_logs", hostname)
		err := os.MkdirAll(logsPath, 0755)
		if err != nil {
			log.Printf("cannot create directory for logs: %s", err)
		}
		err = cliProvisioner.FetchOpenShiftInfraLogs(logsPath)
		if err != nil {
			log.Printf("cliProvisioner.FetchOpenShiftInfraLogs error: %s", err)
		}
	}
	if !cfg.RetainSSH {
		creds := filepath.Join(cfg.CurrentWorkingDir, "_output/", "*ssh*")
		files, err := filepath.Glob(creds)
		if err != nil {
			log.Printf("failed to get ssh files using %s: %s\n", creds, err)
		}
		for _, file := range files {
			err := os.Remove(file)
			if err != nil {
				log.Printf("failed to delete file %s: %s\n", file, err)
			}
		}
	}
	if cfg.CleanUpOnExit {
		for _, rg := range rgs {
			log.Printf("Deleting Group:%s\n", rg)
			acct.DeleteGroup(rg, false)
		}
	}
}

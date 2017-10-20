package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/Azure/acs-engine/test/e2e/azure"
	"github.com/Azure/acs-engine/test/e2e/config"
	"github.com/Azure/acs-engine/test/e2e/engine"
	"github.com/Azure/acs-engine/test/e2e/metrics"
	"github.com/Azure/acs-engine/test/e2e/runner"
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
		cliProvisioner, err := runner.BuildCLIProvisioner(cfg, acct, pt)
		if err != nil {
			log.Fatalf("Error while trying to build CLI Provisioner:%s", err)
		}
		err = cliProvisioner.Run()
		rgs = cliProvisioner.ResourceGroups
		eng = cliProvisioner.Engine
		if err != nil {
			teardown()
			log.Fatalf("Error while trying to provision cluster:%s", err)
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
	if cfg.CleanUpOnExit {
		for _, rg := range rgs {
			log.Printf("Deleting Group:%s\n", rg)
			acct.DeleteGroup(rg)
		}
	}
}

package runner

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/kelseyhightower/envconfig"
)

// Ginkgo contains all of the information needed to run the ginkgo suite of tests
type Ginkgo struct {
	Orchestrator string `envconfig:"ORCHESTRATOR" default:"kubernetes"`
}

// ParseGinkgoConfig creates a new TestRunner object
func ParseGinkgoConfig() (*Ginkgo, error) {
	g := new(Ginkgo)
	if err := envconfig.Process("ginkgo", g); err != nil {
		return nil, err
	}
	return g, nil
}

// Run will execute an orchestrator suite of tests
func (g *Ginkgo) Run() error {
	testDir := fmt.Sprintf("test/e2e/%s", g.Orchestrator)
	cmd := exec.Command("ginkgo", "-nodes", "10", "-slowSpecThreshold", "180", "-r", testDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
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

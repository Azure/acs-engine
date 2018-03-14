package config

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/kelseyhightower/envconfig"
)

// Config holds global test configuration
type Config struct {
	SkipTest          bool          `envconfig:"SKIP_TEST" default:"false"`
	Orchestrator      string        `envconfig:"ORCHESTRATOR" default:"kubernetes"`
	Name              string        `envconfig:"NAME"`                                                                  // Name allows you to set the name of a cluster already created
	Location          string        `envconfig:"LOCATION"`                                                              // Location where you want to create the cluster
	Regions           []string      `envconfig:"REGIONS"`                                                               // A whitelist of availableregions
	ClusterDefinition string        `envconfig:"CLUSTER_DEFINITION" required:"true" default:"examples/kubernetes.json"` // ClusterDefinition is the path on disk to the json template these are normally located in examples/
	CleanUpOnExit     bool          `envconfig:"CLEANUP_ON_EXIT" default:"true"`                                        // if set the tests will not clean up rgs when tests finish
	RetainSSH         bool          `envconfig:"RETAIN_SSH" default:"true"`
	Timeout           time.Duration `envconfig:"TIMEOUT" default:"10m"`
	CurrentWorkingDir string
	SoakClusterName   string `envconfig:"SOAK_CLUSTER_NAME"`
}

const (
	kubernetesOrchestrator = "kubernetes"
	dcosOrchestrator       = "dcos"
	swarmModeOrchestrator  = "swarmmode"
	swarmOrchestrator      = "swarm"
)

// ParseConfig will parse needed environment variables for running the tests
func ParseConfig() (*Config, error) {
	c := new(Config)
	if err := envconfig.Process("config", c); err != nil {
		return nil, err
	}
	if c.Location == "" {
		c.SetRandomRegion()
	}
	return c, nil
}

// GetKubeConfig returns the absolute path to the kubeconfig for c.Location
func (c *Config) GetKubeConfig() string {
	file := fmt.Sprintf("kubeconfig.%s.json", c.Location)
	kubeconfig := filepath.Join(c.CurrentWorkingDir, "_output", c.Name, "kubeconfig", file)
	return kubeconfig
}

// SetKubeConfig will set the KUBECONIFG env var
func (c *Config) SetKubeConfig() {
	os.Setenv("KUBECONFIG", c.GetKubeConfig())
	log.Printf("\nKubeconfig:%s\n", c.GetKubeConfig())
}

// GetSSHKeyPath will return the absolute path to the ssh private key
func (c *Config) GetSSHKeyPath() string {
	return filepath.Join(c.CurrentWorkingDir, "_output", c.Name+"-ssh")
}

// SetEnvVars will determine if we need to
func (c *Config) SetEnvVars() error {
	envFile := fmt.Sprintf("%s/%s.env", c.CurrentWorkingDir, c.ClusterDefinition)
	if _, err := os.Stat(envFile); err == nil {
		file, err := os.Open(envFile)
		if err != nil {
			return err
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			log.Printf("Setting the following:%s\n", line)
			env := strings.Split(line, "=")
			if len(env) > 0 {
				os.Setenv(env[0], env[1])
			}
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}
	return nil
}

// ReadPublicSSHKey will read the contents of the public ssh key on disk into a string
func (c *Config) ReadPublicSSHKey() (string, error) {
	file := c.GetSSHKeyPath() + ".pub"
	contents, err := ioutil.ReadFile(file)
	if err != nil {
		log.Printf("Error while trying to read public ssh key at (%s):%s\n", file, err)
		return "", err
	}
	return string(contents), nil
}

// IsKubernetes will return true if the ORCHESTRATOR env var is set to kubernetes or not set at all
func (c *Config) IsKubernetes() bool {
	return c.Orchestrator == kubernetesOrchestrator
}

// IsDCOS will return true if the ORCHESTRATOR env var is set to dcos
func (c *Config) IsDCOS() bool {
	return c.Orchestrator == dcosOrchestrator
}

// IsSwarmMode will return true if the ORCHESTRATOR env var is set to dcos
func (c *Config) IsSwarmMode() bool {
	return c.Orchestrator == swarmModeOrchestrator
}

// IsSwarm will return true if the ORCHESTRATOR env var is set to dcos
func (c *Config) IsSwarm() bool {
	return c.Orchestrator == swarmOrchestrator
}

// SetRandomRegion sets Location to a random region
func (c *Config) SetRandomRegion() {
	var regions []string
	if c.Regions == nil {
		regions = []string{"eastus", "southcentralus", "westcentralus", "southeastasia", "westus2", "westeurope"}
	} else {
		regions = c.Regions
	}
	log.Printf("Picking Random Region from list %s\n", regions)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	c.Location = regions[r.Intn(len(regions))]
	os.Setenv("LOCATION", c.Location)
	log.Printf("Picked Random Region:%s\n", c.Location)
}

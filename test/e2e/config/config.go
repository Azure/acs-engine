package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"path/filepath"
	"time"

	"github.com/kelseyhightower/envconfig"
)

// Config holds global test configuration
type Config struct {
	Orchestrator      string `envconfig:"ORCHESTRATOR" default:"kubernetes"`
	Name              string `envconfig:"NAME"`                                                                  // Name allows you to set the name of a cluster already created
	Location          string `envconfig:"LOCATION" required:"true" default:"southcentralus"`                     // Location where you want to create the cluster
	ClusterDefinition string `envconfig:"CLUSTER_DEFINITION" required:"true" default:"examples/kubernetes.json"` // ClusterDefinition is the path on disk to the json template these are normally located in examples/
	CleanUpOnExit     bool   `envconfig:"CLEANUP_ON_EXIT" default:"true"`                                        // if set the tests will not clean up rgs when tests finish
	CurrentWorkingDir string
}

// ParseConfig will parse needed environment variables for running the tests
func ParseConfig() (*Config, error) {
	c := new(Config)
	if err := envconfig.Process("config", c); err != nil {
		return nil, err
	}
	return c, nil
}

// GenerateName will generate a new name if one has not been set
func (c *Config) GenerateName() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	suffix := r.Intn(99999)
	prefix := fmt.Sprintf("%s-%s", c.Orchestrator, c.Location)
	return fmt.Sprintf("%s-%v", prefix, suffix)
}

// GetKubeConfig returns the absolute path to the kubeconfig for c.Location
func (c *Config) GetKubeConfig() string {
	file := fmt.Sprintf("kubeconfig.%s.json", c.Location)
	kubeconfig := filepath.Join(c.CurrentWorkingDir, "_output", c.Name, "kubeconfig", file)
	return kubeconfig
}

// GetSSHKeyPath will return the absolute path to the ssh private key
func (c *Config) GetSSHKeyPath() string {
	return filepath.Join(c.CurrentWorkingDir, "_output", c.Name+"-ssh")
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

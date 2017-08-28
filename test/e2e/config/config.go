package config

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/kelseyhightower/envconfig"
)

// Config holds global test configuration
type Config struct {
	Name              string `envconfig:"NAME"`                               // Name allows you to set the name of a cluster already created
	Location          string `envconfig:"LOCATION" required:"true"`           // Location where you want to create the cluster
	ClusterDefinition string `envconfig:"CLUSTER_DEFINITION" required:"true"` // ClusterDefinition is the path on disk to the json template these are normally located in examples/
	CleanUpOnExit     bool   `envconfig:"CLEANUP_ON_EXIT" default:"true"`     // if set the tests will not clean up rgs when tests finish
	SSHKeyName        string `envconfig:"SSH_KEY_NAME"`                       // not absolute path
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
	prefix := fmt.Sprintf("k8s-%s", c.Location)
	return fmt.Sprintf("%s-%v", prefix, suffix)
}

// GetKubeConfig returns the absolute path to the kubeconfig for c.Location
func (c *Config) GetKubeConfig() string {
	cwd, _ := os.Getwd()
	file := fmt.Sprintf("kubeconfig.%s.json", c.Location)
	kubeconfig := filepath.Join(cwd, "../../../_output", c.Name, "kubeconfig", file)
	return kubeconfig
}

// GetSSHKeyPath will return the absolute path to the ssh private key
func (c *Config) GetSSHKeyPath() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		log.Printf("Error while trying to get the current working directory: %s\n", err)
		return "", err
	}
	sshKeyPath := filepath.Join(cwd, "../../../_output", c.SSHKeyName)
	return sshKeyPath, nil
}

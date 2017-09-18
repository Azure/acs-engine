package engine

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/Azure/acs-engine/pkg/api"
	"github.com/Azure/acs-engine/test/e2e/config"
	"github.com/kelseyhightower/envconfig"
)

// Config represents the configuration values of a template stored as env vars
type Config struct {
	ClientID              string `envconfig:"CLIENT_ID"`
	ClientSecret          string `envconfig:"CLIENT_SECRET"`
	MasterDNSPrefix       string `envconfig:"DNS_PREFIX"`
	PublicSSHKey          string `envconfig:"PUBLIC_SSH_KEY"`
	WindowsAdminPasssword string `envconfig:"WINDOWS_ADMIN_PASSWORD"`
	OrchestratorRelease   string `envconfig:"ORCHESTRATOR_RELEASE"`
	OutputDirectory       string `envconfig:"OUTPUT_DIR" default:"_output"`

	ClusterDefinitionPath     string // The original template we want to use to build the cluster from.
	ClusterDefinitionTemplate string // This is the template after we splice in the environment variables
	GeneratedDefinitionPath   string // Holds the contents of running acs-engine generate
	OutputPath                string // This is the root output path
	DefinitionName            string // Unique cluster name
	GeneratedTemplatePath     string // azuredeploy.json path
	GeneratedParametersPath   string // azuredeploy.parameters.json path
}

// Engine holds necessary information to interact with acs-engine cli
type Engine struct {
	Config            *Config
	ClusterDefinition *api.VlabsARMContainerService // Holds the parsed ClusterDefinition
}

// ParseConfig will return a new engine config struct taking values from env vars
func ParseConfig(cwd, clusterDefinition, name string) (*Config, error) {
	c := new(Config)
	if err := envconfig.Process("config", c); err != nil {
		return nil, err
	}

	clusterDefinitionTemplate := fmt.Sprintf("%s/%s.json", c.OutputDirectory, name)
	generatedDefinitionPath := fmt.Sprintf("%s/%s", c.OutputDirectory, name)
	c.DefinitionName = name
	c.ClusterDefinitionPath = filepath.Join(cwd, clusterDefinition)
	c.ClusterDefinitionTemplate = filepath.Join(cwd, clusterDefinitionTemplate)
	c.OutputPath = filepath.Join(cwd, c.OutputDirectory)
	c.GeneratedDefinitionPath = filepath.Join(cwd, generatedDefinitionPath)
	c.GeneratedTemplatePath = filepath.Join(cwd, generatedDefinitionPath, "azuredeploy.json")
	c.GeneratedParametersPath = filepath.Join(cwd, generatedDefinitionPath, "azuredeploy.parameters.json")
	return c, nil
}

// Build takes a template path and will inject values based on provided environment variables
// it will then serialize the structs back into json and save it to outputPath
func Build(cfg *config.Config, subnetID string) (*Engine, error) {
	config, err := ParseConfig(cfg.CurrentWorkingDir, cfg.ClusterDefinition, cfg.Name)
	if err != nil {
		log.Printf("Error while trying to build Engine Configuration:%s\n", err)
	}

	cs, err := Parse(config.ClusterDefinitionPath)
	if err != nil {
		return nil, err
	}

	if config.ClientID != "" && config.ClientSecret != "" {
		cs.ContainerService.Properties.ServicePrincipalProfile.ClientID = config.ClientID
		cs.ContainerService.Properties.ServicePrincipalProfile.Secret = config.ClientSecret
	}

	if config.MasterDNSPrefix != "" {
		cs.ContainerService.Properties.MasterProfile.DNSPrefix = config.MasterDNSPrefix
	}

	if config.PublicSSHKey != "" {
		cs.ContainerService.Properties.LinuxProfile.SSH.PublicKeys[0].KeyData = config.PublicSSHKey
	}

	if config.WindowsAdminPasssword != "" {
		cs.ContainerService.Properties.WindowsProfile.AdminPassword = config.WindowsAdminPasssword
	}

	if config.OrchestratorRelease != "" {
		cs.ContainerService.Properties.OrchestratorProfile.OrchestratorRelease = config.OrchestratorRelease
	}

	if cfg.CreateVNET {
		cs.ContainerService.Properties.MasterProfile.VnetSubnetID = subnetID
		for _, p := range cs.ContainerService.Properties.AgentPoolProfiles {
			p.VnetSubnetID = subnetID
		}
	}

	return &Engine{
		Config:            config,
		ClusterDefinition: cs,
	}, nil
}

// NodeCount returns the number of nodes that should be provisioned for a given cluster definition
func (e *Engine) NodeCount() int {
	expectedCount := e.ClusterDefinition.Properties.MasterProfile.Count
	for _, pool := range e.ClusterDefinition.Properties.AgentPoolProfiles {
		expectedCount = expectedCount + pool.Count
	}
	return expectedCount
}

// HasLinuxAgents will return true if there is at least 1 linux agent pool
func (e *Engine) HasLinuxAgents() bool {
	for _, ap := range e.ClusterDefinition.Properties.AgentPoolProfiles {
		if ap.OSType == "" || ap.OSType == "Linux" {
			return true
		}
	}
	return false
}

// HasWindowsAgents will return true is there is at least 1 windows agent pool
func (e *Engine) HasWindowsAgents() bool {
	for _, ap := range e.ClusterDefinition.Properties.AgentPoolProfiles {
		if ap.OSType == "Windows" {
			return true
		}
	}
	return false
}

// Write will write the cluster definition to disk
func (e *Engine) Write() error {
	json, err := json.Marshal(e.ClusterDefinition)
	if err != nil {
		log.Printf("Error while trying to serialize Container Service object to json:%s\n%+v\n", err, e.ClusterDefinition)
		return err
	}
	err = ioutil.WriteFile(e.Config.ClusterDefinitionTemplate, json, 0777)
	if err != nil {
		log.Printf("Error while trying to write container service definition to file (%s):%s\n%s\n", e.Config.ClusterDefinitionTemplate, err, string(json))
	}
	return nil
}

// Parse takes a template path and will parse that into a api.VlabsARMContainerService
func Parse(path string) (*api.VlabsARMContainerService, error) {
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("Error while trying to read cluster definition at (%s):%s\n", path, err)
		return nil, err
	}
	cs := api.VlabsARMContainerService{}
	if err = json.Unmarshal(contents, &cs); err != nil {
		log.Printf("Error while trying to unmarshal container service json:%s\n%s\n", err, string(contents))
		return nil, err
	}
	return &cs, nil
}

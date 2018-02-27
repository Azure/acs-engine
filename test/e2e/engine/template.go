package engine

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"github.com/Azure/acs-engine/pkg/api"
	"github.com/Azure/acs-engine/pkg/api/vlabs"
	"github.com/Azure/acs-engine/pkg/helpers"
	"github.com/Azure/acs-engine/pkg/i18n"
	"github.com/Azure/acs-engine/test/e2e/config"
	"github.com/kelseyhightower/envconfig"
)

// Config represents the configuration values of a template stored as env vars
type Config struct {
	ClientID              string `envconfig:"CLIENT_ID"`
	ClientSecret          string `envconfig:"CLIENT_SECRET"`
	MasterDNSPrefix       string `envconfig:"DNS_PREFIX"`
	AgentDNSPrefix        string `envconfig:"DNS_PREFIX"`
	PublicSSHKey          string `envconfig:"PUBLIC_SSH_KEY"`
	WindowsAdminPasssword string `envconfig:"WINDOWS_ADMIN_PASSWORD"`
	OrchestratorRelease   string `envconfig:"ORCHESTRATOR_RELEASE"`
	OrchestratorVersion   string `envconfig:"ORCHESTRATOR_VERSION"`
	OutputDirectory       string `envconfig:"OUTPUT_DIR" default:"_output"`
	CreateVNET            bool   `envconfig:"CREATE_VNET" default:"false"`

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
	Config             *Config
	ClusterDefinition  *api.VlabsARMContainerService // Holds the parsed ClusterDefinition
	ExpandedDefinition *api.ContainerService         // Holds the expanded ClusterDefinition
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

	cs, err := ParseInput(config.ClusterDefinitionPath)
	if err != nil {
		return nil, err
	}

	if config.ClientID != "" && config.ClientSecret != "" {
		cs.ContainerService.Properties.ServicePrincipalProfile = &vlabs.ServicePrincipalProfile{
			ClientID: config.ClientID,
			Secret:   config.ClientSecret,
		}
	}

	if config.MasterDNSPrefix != "" {
		cs.ContainerService.Properties.MasterProfile.DNSPrefix = config.MasterDNSPrefix
	}

	if !cfg.IsKubernetes() && config.AgentDNSPrefix != "" {
		for idx, pool := range cs.ContainerService.Properties.AgentPoolProfiles {
			pool.DNSPrefix = fmt.Sprintf("%v-%v", config.AgentDNSPrefix, idx)
		}
	}

	if config.PublicSSHKey != "" {
		cs.ContainerService.Properties.LinuxProfile.SSH.PublicKeys[0].KeyData = config.PublicSSHKey
	}

	if config.WindowsAdminPasssword != "" {
		cs.ContainerService.Properties.WindowsProfile.AdminPassword = config.WindowsAdminPasssword
	}

	// If the parsed api model input has no expressed version opinion, we check if ENV does have an opinion
	if cs.ContainerService.Properties.OrchestratorProfile.OrchestratorRelease == "" &&
		cs.ContainerService.Properties.OrchestratorProfile.OrchestratorVersion == "" {
		// First, prefer the release string if ENV declares it
		if config.OrchestratorRelease != "" {
			cs.ContainerService.Properties.OrchestratorProfile.OrchestratorRelease = config.OrchestratorRelease
			// Or, choose the version string if ENV declares it
		} else if config.OrchestratorVersion != "" {
			cs.ContainerService.Properties.OrchestratorProfile.OrchestratorVersion = config.OrchestratorVersion
			// If ENV similarly has no version opinion, we will rely upon the acs-engine default
		} else {
			log.Println("No orchestrator version specified, will use the default.")
		}
	}

	if config.CreateVNET {
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
	for _, ap := range e.ExpandedDefinition.Properties.AgentPoolProfiles {
		if ap.OSType == "" || ap.OSType == "Linux" {
			return true
		}
	}
	return false
}

// HasWindowsAgents will return true is there is at least 1 windows agent pool
func (e *Engine) HasWindowsAgents() bool {
	for _, ap := range e.ExpandedDefinition.Properties.AgentPoolProfiles {
		if ap.OSType == "Windows" {
			return true
		}
	}
	return false
}

// HasGPUNodes will return true if the VM SKU is GPU-enabled
func (e *Engine) HasGPUNodes() bool {
	for _, ap := range e.ExpandedDefinition.Properties.AgentPoolProfiles {
		if strings.Contains(ap.VMSize, "Standard_N") {
			return true
		}
	}
	return false
}

// HasAddon will return true if an addon is enabled
func (e *Engine) HasAddon(name string) (bool, api.KubernetesAddon) {
	for _, addon := range e.ExpandedDefinition.Properties.OrchestratorProfile.KubernetesConfig.Addons {
		if addon.Name == name {
			return *addon.Enabled, addon
		}
	}
	return false, api.KubernetesAddon{}
}

// OrchestratorVersion1Dot8AndUp will return true if the orchestrator version is 1.8 and up
func (e *Engine) OrchestratorVersion1Dot8AndUp() bool {
	return e.ClusterDefinition.ContainerService.Properties.OrchestratorProfile.OrchestratorVersion >= "1.8"
}

// Write will write the cluster definition to disk
func (e *Engine) Write() error {
	json, err := helpers.JSONMarshal(e.ClusterDefinition, false)
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

// ParseInput takes a template path and will parse that into a api.VlabsARMContainerService
func ParseInput(path string) (*api.VlabsARMContainerService, error) {
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

// ParseOutput takes the generated api model and will parse that into a api.ContainerService
func ParseOutput(path string) (*api.ContainerService, error) {
	locale, err := i18n.LoadTranslations()
	if err != nil {
		return nil, fmt.Errorf(fmt.Sprintf("error loading translation files: %s", err.Error()))
	}
	apiloader := &api.Apiloader{
		Translator: &i18n.Translator{
			Locale: locale,
		},
	}
	containerService, _, err := apiloader.LoadContainerServiceFromFile(path, true, false, nil)
	if err != nil {
		return nil, err
	}
	return containerService, nil
}

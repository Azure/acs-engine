package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"time"

	"github.com/Azure/acs-engine/pkg/acsengine"
	"github.com/Azure/acs-engine/pkg/api"
	"github.com/Azure/acs-engine/pkg/armhelpers"
	"github.com/Azure/acs-engine/pkg/helpers"
	"github.com/Azure/acs-engine/pkg/i18n"
	"github.com/Azure/acs-engine/pkg/operations/kubernetesupgrade"
	"github.com/leonelquinteros/gotext"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	upgradeName             = "upgrade"
	upgradeShortDescription = "Upgrade an existing Kubernetes cluster"
	upgradeLongDescription  = "Upgrade an existing Kubernetes cluster, one minor version at a time"
)

type upgradeCmd struct {
	authArgs

	// user input
	resourceGroupName   string
	deploymentDirectory string
	upgradeVersion      string
	location            string
	timeoutInMinutes    int

	// derived
	containerService    *api.ContainerService
	apiVersion          string
	client              armhelpers.ACSEngineClient
	locale              *gotext.Locale
	nameSuffix          string
	agentPoolsToUpgrade []string
	timeout             *time.Duration
}

// NewUpgradeCmd run a command to upgrade a Kubernetes cluster
func newUpgradeCmd() *cobra.Command {
	uc := upgradeCmd{}

	upgradeCmd := &cobra.Command{
		Use:   upgradeName,
		Short: upgradeShortDescription,
		Long:  upgradeLongDescription,
		RunE: func(cmd *cobra.Command, args []string) error {
			return uc.run(cmd, args)
		},
	}

	f := upgradeCmd.Flags()
	f.StringVarP(&uc.location, "location", "l", "", "location the cluster is deployed in (required)")
	f.StringVarP(&uc.resourceGroupName, "resource-group", "g", "", "the resource group where the cluster is deployed (required)")
	f.StringVar(&uc.deploymentDirectory, "deployment-dir", "", "the location of the output from `generate` (required)")
	f.StringVar(&uc.upgradeVersion, "upgrade-version", "", "desired kubernetes version (required)")
	f.IntVar(&uc.timeoutInMinutes, "vm-timeout", -1, "how long to wait for each vm to be upgraded in minutes")
	addAuthFlags(&uc.authArgs, f)

	return upgradeCmd
}

func (uc *upgradeCmd) validate(cmd *cobra.Command) error {
	log.Infoln("validating...")

	var err error

	uc.locale, err = i18n.LoadTranslations()
	if err != nil {
		return fmt.Errorf("error loading translation files: %s", err.Error())
	}

	if uc.resourceGroupName == "" {
		cmd.Usage()
		return fmt.Errorf("--resource-group must be specified")
	}

	if uc.location == "" {
		cmd.Usage()
		return fmt.Errorf("--location must be specified")
	}
	uc.location = helpers.NormalizeAzureRegion(uc.location)

	if uc.timeoutInMinutes != -1 {
		timeout := time.Duration(uc.timeoutInMinutes) * time.Minute
		uc.timeout = &timeout
	}

	// TODO(colemick): add in the cmd annotation to help enable autocompletion
	if uc.upgradeVersion == "" {
		cmd.Usage()
		return fmt.Errorf("--upgrade-version must be specified")
	}

	if uc.deploymentDirectory == "" {
		cmd.Usage()
		return fmt.Errorf("--deployment-dir must be specified")
	}
	return nil
}

func (uc *upgradeCmd) loadCluster(cmd *cobra.Command) error {
	var err error

	if err = uc.authArgs.validateAuthArgs(); err != nil {
		return fmt.Errorf("%s", err.Error())
	}

	if uc.client, err = uc.authArgs.getClient(); err != nil {
		return fmt.Errorf("Failed to get client: %s", err.Error())
	}

	_, err = uc.client.EnsureResourceGroup(uc.resourceGroupName, uc.location, nil)
	if err != nil {
		return fmt.Errorf("Error ensuring resource group: %s", err.Error())
	}

	// load apimodel from the deployment directory
	apiModelPath := path.Join(uc.deploymentDirectory, "apimodel.json")

	if _, err = os.Stat(apiModelPath); os.IsNotExist(err) {
		return fmt.Errorf("specified api model does not exist (%s)", apiModelPath)
	}

	apiloader := &api.Apiloader{
		Translator: &i18n.Translator{
			Locale: uc.locale,
		},
	}
	uc.containerService, uc.apiVersion, err = apiloader.LoadContainerServiceFromFile(apiModelPath, true, true, nil)
	if err != nil {
		return fmt.Errorf("error parsing the api model: %s", err.Error())
	}

	if uc.containerService.Location == "" {
		uc.containerService.Location = uc.location
	} else if uc.containerService.Location != uc.location {
		return fmt.Errorf("--location does not match api model location")
	}

	// get available upgrades for container service
	orchestratorInfo, err := api.GetOrchestratorVersionProfile(uc.containerService.Properties.OrchestratorProfile)
	if err != nil {
		return fmt.Errorf("error getting list of available upgrades: %s", err.Error())
	}
	// add the current version if upgrade has failed
	orchestratorInfo.Upgrades = append(orchestratorInfo.Upgrades, &api.OrchestratorProfile{
		OrchestratorType:    uc.containerService.Properties.OrchestratorProfile.OrchestratorType,
		OrchestratorVersion: uc.containerService.Properties.OrchestratorProfile.OrchestratorVersion})

	// validate desired upgrade version and set goal state
	found := false
	for _, up := range orchestratorInfo.Upgrades {
		if up.OrchestratorVersion == uc.upgradeVersion {
			uc.containerService.Properties.OrchestratorProfile.OrchestratorVersion = uc.upgradeVersion
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("version %s is not supported", uc.upgradeVersion)
	}

	// Read name suffix to identify nodes in the resource group that belong
	// to this cluster.
	// TODO: Also update to read  namesuffix from the parameters file as
	// user could have specified a name suffix instead of using the default
	// value generated by ACS Engine
	templatePath := path.Join(uc.deploymentDirectory, "azuredeploy.json")
	contents, _ := ioutil.ReadFile(templatePath)

	var template interface{}
	json.Unmarshal(contents, &template)

	templateMap := template.(map[string]interface{})
	templateParameters := templateMap["parameters"].(map[string]interface{})

	nameSuffixParam := templateParameters["nameSuffix"].(map[string]interface{})
	uc.nameSuffix = nameSuffixParam["defaultValue"].(string)
	log.Infoln(fmt.Sprintf("Name suffix: %s", uc.nameSuffix))

	uc.agentPoolsToUpgrade = []string{}
	log.Infoln(fmt.Sprintf("Gathering agent pool names..."))
	for _, agentPool := range uc.containerService.Properties.AgentPoolProfiles {
		uc.agentPoolsToUpgrade = append(uc.agentPoolsToUpgrade, agentPool.Name)
	}
	return nil
}

func (uc *upgradeCmd) run(cmd *cobra.Command, args []string) error {
	err := uc.validate(cmd)
	if err != nil {
		log.Fatalf("error validating upgrade command: %v", err)
	}

	err = uc.loadCluster(cmd)
	if err != nil {
		log.Fatalf("error loading existing cluster: %v", err)
	}

	upgradeCluster := kubernetesupgrade.UpgradeCluster{
		Translator: &i18n.Translator{
			Locale: uc.locale,
		},
		Logger:      log.NewEntry(log.New()),
		Client:      uc.client,
		StepTimeout: uc.timeout,
	}

	kubeConfig, err := acsengine.GenerateKubeConfig(uc.containerService.Properties, uc.location)
	if err != nil {
		log.Fatalf("failed to generate kube config: %v", err) // TODO: cleanup
	}

	if err = upgradeCluster.UpgradeCluster(uc.authArgs.SubscriptionID, kubeConfig, uc.resourceGroupName,
		uc.containerService, uc.nameSuffix, uc.agentPoolsToUpgrade, BuildTag); err != nil {
		log.Fatalf("Error upgrading cluster: %v\n", err)
	}

	apiloader := &api.Apiloader{
		Translator: &i18n.Translator{
			Locale: uc.locale,
		},
	}
	b, err := apiloader.SerializeContainerService(uc.containerService, uc.apiVersion)
	if err != nil {
		return err
	}

	f := acsengine.FileSaver{
		Translator: &i18n.Translator{
			Locale: uc.locale,
		},
	}

	return f.SaveFile(uc.deploymentDirectory, "apimodel.json", b)
}

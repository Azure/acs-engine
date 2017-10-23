package cmd

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path"
	"time"

	"github.com/Azure/acs-engine/pkg/acsengine"
	"github.com/Azure/acs-engine/pkg/api"
	"github.com/Azure/acs-engine/pkg/armhelpers"
	"github.com/Azure/acs-engine/pkg/i18n"
	"github.com/leonelquinteros/gotext"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type scaleCmd struct {
	authArgs

	// user input
	resourceGroupName    string
	deploymentDirectory  string
	newDesiredAgentCount int
	containerService     *api.ContainerService
	apiVersion           string
	location             string
	agentPoolToScale     string
	classicMode          bool

	// derived
	apiModelPath   string
	agentPool      *api.AgentPoolProfile
	client         armhelpers.ACSEngineClient
	locale         *gotext.Locale
	nameSuffix     string
	agentPoolIndex int
}

const (
	scaleName             = "scale"
	scaleShortDescription = "scale a deployed cluster"
	scaleLongDescription  = "scale a deployed cluster"
)

// NewScaleCmd run a command to upgrade a Kubernetes cluster
func newScaleCmd() *cobra.Command {
	sc := scaleCmd{}

	scaleCmd := &cobra.Command{
		Use:   scaleName,
		Short: scaleShortDescription,
		Long:  scaleLongDescription,
		RunE: func(cmd *cobra.Command, args []string) error {
			return sc.run(cmd, args)
		},
	}

	f := scaleCmd.Flags()
	f.StringVar(&sc.location, "location", "", "location the cluster is deployed in")
	f.StringVar(&sc.resourceGroupName, "resource-group", "", "the resource group where the cluster is deployed")
	f.StringVar(&sc.deploymentDirectory, "deployment-dir", "", "the location of the output from `generate`")
	f.IntVar(&sc.newDesiredAgentCount, "new-agent-count", 0, "desired number of agents")
	f.BoolVar(&sc.classicMode, "classic-mode", false, "enable classic parameters and outputs")
	f.StringVar(&sc.agentPoolToScale, "agent-pool", "", "agent pool to scale")

	addAuthFlags(&sc.authArgs, f)

	return scaleCmd
}

func (sc *scaleCmd) validate(cmd *cobra.Command, args []string) {
	log.Infoln("validating...")
	var err error

	sc.locale, err = i18n.LoadTranslations()
	if err != nil {
		log.Fatalf("error loading translation files: %s", err.Error())
	}

	if sc.resourceGroupName == "" {
		cmd.Usage()
		log.Fatal("--resource-group must be specified")
	}

	if sc.location == "" {
		cmd.Usage()
		log.Fatal("--location must be specified")
	}

	if sc.newDesiredAgentCount == 0 {
		cmd.Usage()
		log.Fatal("--new-agent-count must be specified")
	}

	if sc.client, err = sc.authArgs.getClient(); err != nil {
		log.Error("Failed to get client:", err)
	}

	if sc.deploymentDirectory == "" {
		cmd.Usage()
		log.Fatal("--deployment-dir must be specified")
	}

	_, err = sc.client.EnsureResourceGroup(sc.resourceGroupName, sc.location, nil)
	if err != nil {
		log.Fatalln(err)
	}

	// load apimodel from the deployment directory
	sc.apiModelPath = path.Join(sc.deploymentDirectory, "apimodel.json")

	if _, err = os.Stat(sc.apiModelPath); os.IsNotExist(err) {
		log.Fatalf("specified api model does not exist (%s)", sc.apiModelPath)
	}

	apiloader := &api.Apiloader{
		Translator: &i18n.Translator{
			Locale: sc.locale,
		},
	}
	sc.containerService, sc.apiVersion, err = apiloader.LoadContainerServiceFromFile(sc.apiModelPath, true, nil)
	if err != nil {
		log.Fatalf("error parsing the api model: %s", err.Error())
	}

	if sc.agentPoolToScale == "" {
		agentPoolCount := len(sc.containerService.Properties.AgentPoolProfiles)
		if agentPoolCount > 1 {
			log.Fatal("--agent-pool is required if more than one agent pool is defined in the container service")
		} else if agentPoolCount == 1 {
			sc.agentPool = sc.containerService.Properties.AgentPoolProfiles[0]
			sc.agentPoolIndex = 0
			sc.agentPoolToScale = sc.containerService.Properties.AgentPoolProfiles[0].Name
		} else {
			log.Fatal("No agent pools found to scale")
		}
	} else {
		agentPoolIndex := -1
		for i, pool := range sc.containerService.Properties.AgentPoolProfiles {
			if pool.Name == sc.agentPoolToScale {
				agentPoolIndex = i
				sc.agentPool = pool
				sc.agentPoolIndex = i
			}
		}
		if agentPoolIndex == -1 {
			log.Fatalf("Agent pool %s wasn't in the deployed api model", sc.agentPoolToScale)
		}
	}
}

func (sc *scaleCmd) run(cmd *cobra.Command, args []string) error {
	sc.validate(cmd, args)

	if sc.agentPool.IsAvailabilitySets() && sc.agentPool.Count > sc.newDesiredAgentCount {
		// TODO add scale down of VMS code path
		log.Fatalln("Scaling down availability sets is not currently supported")
	}

	ctx := acsengine.Context{
		Translator: &i18n.Translator{
			Locale: sc.locale,
		},
	}
	templateGenerator, err := acsengine.InitializeTemplateGenerator(ctx, sc.classicMode)
	if err != nil {
		log.Fatalln("failed to initialize template generator: %s", err.Error())
	}

	sc.containerService.Properties.AgentPoolProfiles = []*api.AgentPoolProfile{sc.agentPool}

	template, parameters, _, err := templateGenerator.GenerateTemplate(sc.containerService, acsengine.DefaultGeneratorCode)
	if err != nil {
		log.Fatalf("error generating template %s: %s", sc.apiModelPath, err.Error())
		os.Exit(1)
	}

	if template, err = acsengine.PrettyPrintArmTemplate(template); err != nil {
		log.Fatalf("error pretty printing template: %s \n", err.Error())
	}

	templateJSON := make(map[string]interface{})
	parametersJSON := make(map[string]interface{})

	err = json.Unmarshal([]byte(template), &templateJSON)
	if err != nil {
		log.Fatalln(err)
	}

	err = json.Unmarshal([]byte(parameters), &parametersJSON)
	if err != nil {
		log.Fatalln(err)
	}

	transformer := acsengine.Transformer{Translator: ctx.Translator}
	orchestratorInfo := sc.containerService.Properties.OrchestratorProfile
	addValue(parametersJSON, sc.agentPool.Name+"Count", sc.newDesiredAgentCount)

	switch orchestratorInfo.OrchestratorType {
	case api.Kubernetes:
		err = transformer.NormalizeForK8sVMASScalingUp(log.New().WithField("source", "scaling command line"), templateJSON)
		if err != nil {
			log.Fatalf("error tranforming the template for scaling template %s: %s", sc.apiModelPath, err.Error())
			os.Exit(1)
		}
		if sc.agentPool.IsAvailabilitySets() {
			addValue(parametersJSON, fmt.Sprintf("%sOffset", sc.agentPool.Name), sc.agentPool.Count)
		}
		break
	case api.Swarm:
	case api.SwarmMode:
	case api.DCOS:
		if sc.agentPool.IsAvailabilitySets() {
			log.Fatalf("scaling isn't supported for orchestrator %s, with availability sets", orchestratorInfo.OrchestratorType)
			os.Exit(1)
		}
		transformer.NormalizeForVMSSScaling(log.New().WithField("source", "scaling command line"), templateJSON)
	}

	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	deploymentSuffix := random.Int31()

	_, err = sc.client.DeployTemplate(
		sc.resourceGroupName,
		fmt.Sprintf("%s-%d", sc.resourceGroupName, deploymentSuffix),
		templateJSON,
		parametersJSON,
		nil)
	if err != nil {
		log.Fatalln(err)
	}

	apiloader := &api.Apiloader{
		Translator: &i18n.Translator{
			Locale: sc.locale,
		},
	}
	var apiVersion string
	sc.containerService, apiVersion, err = apiloader.LoadContainerServiceFromFile(sc.apiModelPath, true, nil)
	sc.containerService.Properties.AgentPoolProfiles[sc.agentPoolIndex].Count = sc.newDesiredAgentCount

	b, e := apiloader.SerializeContainerService(sc.containerService, apiVersion)

	if err != nil {
		return err
	}

	f := acsengine.FileSaver{
		Translator: &i18n.Translator{
			Locale: sc.locale,
		},
	}

	if e = f.SaveFile(sc.deploymentDirectory, "apimodel.json", b); e != nil {
		return e
	}

	return nil
}

type paramsMap map[string]interface{}

func addValue(m paramsMap, k string, v interface{}) {
	m[k] = paramsMap{
		"value": v,
	}
}

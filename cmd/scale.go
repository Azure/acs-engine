package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"sort"
	"strings"
	"time"

	"github.com/Azure/acs-engine/pkg/acsengine"
	"github.com/Azure/acs-engine/pkg/api"
	"github.com/Azure/acs-engine/pkg/armhelpers"
	"github.com/Azure/acs-engine/pkg/i18n"
	"github.com/Azure/acs-engine/pkg/operations"
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
	masterFQDN     string
	logger         *log.Entry
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
	f.IntVar(&sc.newDesiredAgentCount, "new-node-count", 0, "desired number of nodes")
	f.BoolVar(&sc.classicMode, "classic-mode", false, "enable classic parameters and outputs")
	f.StringVar(&sc.agentPoolToScale, "node-pool", "", "node pool to scale")
	f.StringVar(&sc.masterFQDN, "master-FQDN", "", "FQDN for the master load balancer, Needed to scale down Kubernetes agent pools")

	addAuthFlags(&sc.authArgs, f)

	return scaleCmd
}

func (sc *scaleCmd) validate(cmd *cobra.Command, args []string) {
	log.Infoln("validating...")
	sc.logger = log.New().WithField("source", "scaling command line")
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
		log.Fatal("--new-node-count must be specified")
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
	sc.containerService, sc.apiVersion, err = apiloader.LoadContainerServiceFromFile(sc.apiModelPath, true, true, nil)
	if err != nil {
		log.Fatalf("error parsing the api model: %s", err.Error())
	}

	if sc.agentPoolToScale == "" {
		agentPoolCount := len(sc.containerService.Properties.AgentPoolProfiles)
		if agentPoolCount > 1 {
			log.Fatal("--node-pool is required if more than one agent pool is defined in the container service")
		} else if agentPoolCount == 1 {
			sc.agentPool = sc.containerService.Properties.AgentPoolProfiles[0]
			sc.agentPoolIndex = 0
			sc.agentPoolToScale = sc.containerService.Properties.AgentPoolProfiles[0].Name
		} else {
			log.Fatal("No node pools found to scale")
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
			log.Fatalf("node pool %s wasn't in the deployed api model", sc.agentPoolToScale)
		}
	}

	templatePath := path.Join(sc.deploymentDirectory, "azuredeploy.json")
	contents, _ := ioutil.ReadFile(templatePath)

	var template interface{}
	json.Unmarshal(contents, &template)

	templateMap := template.(map[string]interface{})
	templateParameters := templateMap["parameters"].(map[string]interface{})

	nameSuffixParam := templateParameters["nameSuffix"].(map[string]interface{})
	sc.nameSuffix = nameSuffixParam["defaultValue"].(string)
	log.Infoln(fmt.Sprintf("Name suffix: %s", sc.nameSuffix))
}

func (sc *scaleCmd) run(cmd *cobra.Command, args []string) error {
	sc.validate(cmd, args)

	orchestratorInfo := sc.containerService.Properties.OrchestratorProfile
	var currentNodeCount, highestUsedIndex int
	indexes := make([]int, 0)
	indexToVM := make(map[int]string)
	if sc.agentPool.IsAvailabilitySets() {
		//TODO handle when there is a nextLink in the response and get more nodes
		vms, err := sc.client.ListVirtualMachines(sc.resourceGroupName)
		if err != nil {
			log.Fatalln("failed to get vms in the resource group. Error: %s", err.Error())
		}
		for _, vm := range *vms.Value {

			poolName, nameSuffix, index, err := armhelpers.K8sLinuxVMNameParts(*vm.Name)
			if err != nil || !strings.EqualFold(poolName, sc.agentPoolToScale) || !strings.EqualFold(nameSuffix, sc.nameSuffix) {
				continue
			}

			indexToVM[index] = *vm.Name
			indexes = append(indexes, index)
		}
		sortedIndexes := sort.IntSlice(indexes)
		sortedIndexes.Sort()
		indexes = []int(sortedIndexes)
		currentNodeCount = len(indexes)

		if currentNodeCount == sc.newDesiredAgentCount {
			return nil
		}
		highestUsedIndex = indexes[len(indexes)-1]

		// Scale down Scenario
		if currentNodeCount > sc.newDesiredAgentCount {
			if sc.masterFQDN == "" {
				cmd.Usage()
				log.Fatal("master-FQDN is required to scale down a kubernetes cluster's agent pool")
			}

			vmsToDelete := make([]string, 0)
			for i := currentNodeCount - 1; i >= sc.newDesiredAgentCount; i-- {
				vmsToDelete = append(vmsToDelete, indexToVM[i])
			}

			if orchestratorInfo.OrchestratorType == api.Kubernetes {
				err = sc.drainNodes(vmsToDelete)
				if err != nil {
					log.Errorf("Got error %+v, while draining the nodes to be deleted", err)
					return err
				}
			}

			errList := operations.ScaleDownVMs(sc.client, sc.logger, sc.resourceGroupName, vmsToDelete...)
			if errList != nil {
				errorMessage := ""
				for element := errList.Front(); element != nil; element = element.Next() {
					vmError, ok := element.Value.(*operations.VMScalingErrorDetails)
					if ok {
						error := fmt.Sprintf("Node '%s' failed to delete with error: '%s'", vmError.Name, vmError.Error.Error())
						errorMessage = errorMessage + error
					}
				}
				return fmt.Errorf(errorMessage)
			}

			return nil
		}
	} else {
		vmssList, err := sc.client.ListVirtualMachineScaleSets(sc.resourceGroupName)
		if err != nil {
			log.Fatalln("failed to get vmss list in the resource group. Error: %s", err.Error())
		}
		for _, vmss := range *vmssList.Value {
			poolName, nameSuffix, err := armhelpers.VmssNameParts(*vmss.Name)
			if err != nil || !strings.EqualFold(poolName, sc.agentPoolToScale) || !strings.EqualFold(nameSuffix, sc.nameSuffix) {
				continue
			}
			currentNodeCount = int(*vmss.Sku.Capacity)
			highestUsedIndex = 0
		}
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
	// Our templates generate a range of nodes based on a count and offset, it is possible for there to be holes in the template
	// So we need to set the count in the template to get enough nodes for the range, if there are holes that number will be larger than the desired count
	countForTemplate := sc.newDesiredAgentCount
	if highestUsedIndex != 0 {
		countForTemplate += highestUsedIndex + 1 - currentNodeCount
	}
	addValue(parametersJSON, sc.agentPool.Name+"Count", countForTemplate)

	switch orchestratorInfo.OrchestratorType {
	case api.Kubernetes:
		err = transformer.NormalizeForK8sVMASScalingUp(sc.logger, templateJSON)
		if err != nil {
			log.Fatalf("error tranforming the template for scaling template %s: %s", sc.apiModelPath, err.Error())
			os.Exit(1)
		}
		if sc.agentPool.IsAvailabilitySets() {
			addValue(parametersJSON, fmt.Sprintf("%sOffset", sc.agentPool.Name), highestUsedIndex+1)
		}
		break
	case api.Swarm:
	case api.SwarmMode:
	case api.DCOS:
		if sc.agentPool.IsAvailabilitySets() {
			log.Fatalf("scaling isn't supported for orchestrator %s, with availability sets", orchestratorInfo.OrchestratorType)
			os.Exit(1)
		}
		transformer.NormalizeForVMSSScaling(sc.logger, templateJSON)
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
	sc.containerService, apiVersion, err = apiloader.LoadContainerServiceFromFile(sc.apiModelPath, false, true, nil)
	sc.containerService.Properties.AgentPoolProfiles[sc.agentPoolIndex].Count = sc.newDesiredAgentCount

	b, e := apiloader.SerializeContainerService(sc.containerService, apiVersion)

	if e != nil {
		return e
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

func (sc *scaleCmd) drainNodes(vmsToDelete []string) error {
	kubeConfig, err := acsengine.GenerateKubeConfig(sc.containerService.Properties, sc.location)
	if err != nil {
		log.Fatalf("failed to generate kube config") // TODO: cleanup
	}
	var errorMessage string
	masterURL := sc.masterFQDN
	if !strings.HasPrefix(masterURL, "https://") {
		masterURL = fmt.Sprintf("https://%s", masterURL)
	}
	numVmsToDrain := len(vmsToDelete)
	errChan := make(chan *operations.VMScalingErrorDetails, numVmsToDrain)
	defer close(errChan)
	for _, vmName := range vmsToDelete {
		go func(vmName string) {
			e := operations.SafelyDrainNode(sc.client, sc.logger,
				masterURL, kubeConfig, vmName, time.Duration(60)*time.Minute)
			if e != nil {
				log.Errorf("Failed to drain node %s, got error %s", vmName, e.Error())
				errChan <- &operations.VMScalingErrorDetails{Error: e, Name: vmName}
				return
			}
			errChan <- nil
		}(vmName)
	}

	for i := 0; i < numVmsToDrain; i++ {
		errDetails := <-errChan
		if errDetails != nil {
			error := fmt.Sprintf("Node '%s' failed to drain with error: '%s'", errDetails.Name, errDetails.Error.Error())
			errorMessage = errorMessage + error
			return fmt.Errorf(error)
		}
	}

	return nil
}

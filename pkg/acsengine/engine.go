package acsengine

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"hash/fnv"
	"io/ioutil"
	"math/rand"
	"net/http"
	"regexp"
	"runtime/debug"
	"strconv"
	"strings"
	"text/template"

	//log "github.com/Sirupsen/logrus"
	"github.com/Azure/acs-engine/pkg/api"
	"github.com/Azure/acs-engine/pkg/i18n"
	"github.com/Masterminds/semver"
	"github.com/ghodss/yaml"
)

const (
	kubernetesMasterCustomDataYaml      = "kubernetesmastercustomdata.yml"
	kubernetesMasterCustomScript        = "kubernetesmastercustomscript.sh"
	kubernetesAgentCustomDataYaml       = "kubernetesagentcustomdata.yml"
	kubeConfigJSON                      = "kubeconfig.json"
	kubernetesWindowsAgentCustomDataPS1 = "kuberneteswindowssetup.ps1"
)

const (
	dcosCustomData173    = "dcoscustomdata173.t"
	dcosCustomData188    = "dcoscustomdata188.t"
	dcosCustomData190    = "dcoscustomdata190.t"
	dcosProvision        = "dcosprovision.sh"
	dcosWindowsProvision = "dcosWindowsProvision.ps1"
)

const (
	swarmProvision        = "configure-swarm-cluster.sh"
	swarmWindowsProvision = "Install-ContainerHost-And-Join-Swarm.ps1"

	swarmModeProvision        = "configure-swarmmode-cluster.sh"
	swarmModeWindowsProvision = "Join-SwarmMode-cluster.ps1"
)

const (
	agentOutputs                  = "agentoutputs.t"
	agentParams                   = "agentparams.t"
	classicParams                 = "classicparams.t"
	dcosAgentResourcesVMAS        = "dcosagentresourcesvmas.t"
	dcosWindowsAgentResourcesVMAS = "dcosWindowsAgentResourcesVmas.t"
	dcosAgentResourcesVMSS        = "dcosagentresourcesvmss.t"
	dcosWindowsAgentResourcesVMSS = "dcosWindowsAgentResourcesVmss.t"
	dcosAgentVars                 = "dcosagentvars.t"
	dcosBaseFile                  = "dcosbase.t"
	dcosParams                    = "dcosparams.t"
	dcosMasterResources           = "dcosmasterresources.t"
	dcosMasterVars                = "dcosmastervars.t"
	iaasOutputs                   = "iaasoutputs.t"
	kubernetesBaseFile            = "kubernetesbase.t"
	kubernetesAgentResourcesVMAS  = "kubernetesagentresourcesvmas.t"
	kubernetesAgentVars           = "kubernetesagentvars.t"
	kubernetesMasterResources     = "kubernetesmasterresources.t"
	kubernetesMasterVars          = "kubernetesmastervars.t"
	kubernetesParams              = "kubernetesparams.t"
	kubernetesWinAgentVars        = "kuberneteswinagentresourcesvmas.t"
	kubernetesKubeletService      = "kuberneteskubelet.service"
	masterOutputs                 = "masteroutputs.t"
	masterParams                  = "masterparams.t"
	swarmBaseFile                 = "swarmbase.t"
	swarmAgentResourcesVMAS       = "swarmagentresourcesvmas.t"
	swarmAgentResourcesVMSS       = "swarmagentresourcesvmss.t"
	swarmAgentResourcesClassic    = "swarmagentresourcesclassic.t"
	swarmAgentVars                = "swarmagentvars.t"
	swarmMasterResources          = "swarmmasterresources.t"
	swarmMasterVars               = "swarmmastervars.t"
	swarmWinAgentResourcesVMAS    = "swarmwinagentresourcesvmas.t"
	swarmWinAgentResourcesVMSS    = "swarmwinagentresourcesvmss.t"
	windowsParams                 = "windowsparams.t"
)

const (
	azurePublicCloud       = "AzurePublicCloud"
	azureChinaCloud        = "AzureChinaCloud"
	azureGermanCloud       = "AzureGermanCloud"
	azureUSGovernmentCloud = "AzureUSGovernmentCloud"
)

var kubernetesManifestYamls = map[string]string{
	"MASTER_KUBERNETES_SCHEDULER_B64_GZIP_STR":          "kubernetesmaster-kube-scheduler.yaml",
	"MASTER_KUBERNETES_CONTROLLER_MANAGER_B64_GZIP_STR": "kubernetesmaster-kube-controller-manager.yaml",
	"MASTER_KUBERNETES_APISERVER_B64_GZIP_STR":          "kubernetesmaster-kube-apiserver.yaml",
	"MASTER_KUBERNETES_ADDON_MANAGER_B64_GZIP_STR":      "kubernetesmaster-kube-addon-manager.yaml",
}

var kubernetesAritfacts = map[string]string{
	"MASTER_PROVISION_B64_GZIP_STR": kubernetesMasterCustomScript,
	"KUBELET_SERVICE_B64_GZIP_STR":  kubernetesKubeletService,
}

var kubernetesAritfacts15 = map[string]string{
	"MASTER_PROVISION_B64_GZIP_STR": kubernetesMasterCustomScript,
	"KUBELET_SERVICE_B64_GZIP_STR":  "kuberneteskubelet1.5.service",
}

var kubernetesAddonYamls = map[string]string{
	"MASTER_ADDON_HEAPSTER_DEPLOYMENT_B64_GZIP_STR":             "kubernetesmasteraddons-heapster-deployment.yaml",
	"MASTER_ADDON_KUBE_DNS_DEPLOYMENT_B64_GZIP_STR":             "kubernetesmasteraddons-kube-dns-deployment.yaml",
	"MASTER_ADDON_KUBE_PROXY_DAEMONSET_B64_GZIP_STR":            "kubernetesmasteraddons-kube-proxy-daemonset.yaml",
	"MASTER_ADDON_KUBERNETES_DASHBOARD_DEPLOYMENT_B64_GZIP_STR": "kubernetesmasteraddons-kubernetes-dashboard-deployment.yaml",
	"MASTER_ADDON_AZURE_STORAGE_CLASSES_B64_GZIP_STR":           "kubernetesmasteraddons-azure-storage-classes.yaml",
	"MASTER_ADDON_TILLER_DEPLOYMENT_B64_GZIP_STR":               "kubernetesmasteraddons-tiller-deployment.yaml",
}

var kubernetesAddonYamls15 = map[string]string{
	"MASTER_ADDON_HEAPSTER_DEPLOYMENT_B64_GZIP_STR":             "kubernetesmasteraddons-heapster-deployment1.5.yaml",
	"MASTER_ADDON_KUBE_DNS_DEPLOYMENT_B64_GZIP_STR":             "kubernetesmasteraddons-kube-dns-deployment1.5.yaml",
	"MASTER_ADDON_KUBE_PROXY_DAEMONSET_B64_GZIP_STR":            "kubernetesmasteraddons-kube-proxy-daemonset1.5.yaml",
	"MASTER_ADDON_KUBERNETES_DASHBOARD_DEPLOYMENT_B64_GZIP_STR": "kubernetesmasteraddons-kubernetes-dashboard-deployment1.5.yaml",
	"MASTER_ADDON_AZURE_STORAGE_CLASSES_B64_GZIP_STR":           "kubernetesmasteraddons-azure-storage-classes.yaml",
	"MASTER_ADDON_TILLER_DEPLOYMENT_B64_GZIP_STR":               "kubernetesmasteraddons-tiller-deployment1.5.yaml",
}

var calicoAddonYamls = map[string]string{
	"MASTER_ADDON_CALICO_DAEMONSET_B64_GZIP_STR": "kubernetesmasteraddons-calico-daemonset.yaml",
}

var calicoAddonYamls15 = map[string]string{
	"MASTER_ADDON_CALICO_DAEMONSET_B64_GZIP_STR": "kubernetesmasteraddons-calico-daemonset1.5.yaml",
}

var commonTemplateFiles = []string{agentOutputs, agentParams, classicParams, masterOutputs, iaasOutputs, masterParams, windowsParams}
var dcosTemplateFiles = []string{dcosBaseFile, dcosAgentResourcesVMAS, dcosAgentResourcesVMSS, dcosAgentVars, dcosMasterResources, dcosMasterVars, dcosParams, dcosWindowsAgentResourcesVMAS, dcosWindowsAgentResourcesVMSS}
var kubernetesTemplateFiles = []string{kubernetesBaseFile, kubernetesAgentResourcesVMAS, kubernetesAgentVars, kubernetesMasterResources, kubernetesMasterVars, kubernetesParams, kubernetesWinAgentVars}
var swarmTemplateFiles = []string{swarmBaseFile, swarmAgentResourcesVMAS, swarmAgentVars, swarmAgentResourcesVMSS, swarmAgentResourcesClassic, swarmBaseFile, swarmMasterResources, swarmMasterVars, swarmWinAgentResourcesVMAS, swarmWinAgentResourcesVMSS}
var swarmModeTemplateFiles = []string{swarmBaseFile, swarmAgentResourcesVMAS, swarmAgentVars, swarmAgentResourcesVMSS, swarmAgentResourcesClassic, swarmBaseFile, swarmMasterResources, swarmMasterVars, swarmWinAgentResourcesVMAS, swarmWinAgentResourcesVMSS}

/**
 The following parameters could be either a plain text, or referenced to a secret in a keyvault:
 - apiServerCertificate
 - apiServerPrivateKey
 - caCertificate
 - clientCertificate
 - clientPrivateKey
 - kubeConfigCertificate
 - kubeConfigPrivateKey
 - servicePrincipalClientSecret

 To refer to a keyvault secret, the value of the parameter in the api model file should be formatted as:

 "<PARAMETER>": "/subscriptions/<SUB_ID>/resourceGroups/<RG_NAME>/providers/Microsoft.KeyVault/vaults/<KV_NAME>/secrets/<NAME>[/<VERSION>]"
 where:
   <SUB_ID> is the subscription ID of the keyvault
   <RG_NAME> is the resource group of the keyvault
   <KV_NAME> is the name of the keyvault
   <NAME> is the name of the secret.
   <VERSION> (optional) is the version of the secret (default: the latest version)

 This will generate a reference block in the parameters file:

 "reference": {
   "keyVault": {
     "id": "/subscriptions/<SUB_ID>/resourceGroups/<RG_NAME>/providers/Microsoft.KeyVault/vaults/<KV_NAME>"
   },
   "secretName": "<NAME>"
   "secretVersion": "<VERSION>"
}
**/

// KeyVaultID represents a KeyVault instance on Azure
type KeyVaultID struct {
	ID string `json:"id"`
}

// KeyVaultRef represents a reference to KeyVault instance on Azure
type KeyVaultRef struct {
	KeyVault      KeyVaultID `json:"keyVault"`
	SecretName    string     `json:"secretName"`
	SecretVersion string     `json:"secretVersion,omitempty"`
}

type paramsMap map[string]interface{}

var keyvaultSecretPathRe *regexp.Regexp

func init() {
	keyvaultSecretPathRe = regexp.MustCompile(`^(/subscriptions/\S+/resourceGroups/\S+/providers/Microsoft.KeyVault/vaults/\S+)/secrets/([^/\s]+)(/(\S+))?$`)
}

func (t *TemplateGenerator) verifyFiles() error {
	allFiles := commonTemplateFiles
	allFiles = append(allFiles, dcosTemplateFiles...)
	allFiles = append(allFiles, kubernetesTemplateFiles...)
	allFiles = append(allFiles, swarmTemplateFiles...)
	for _, file := range allFiles {
		if _, err := Asset(file); err != nil {
			return t.Translator.Errorf("template file %s does not exist", file)
		}
	}
	return nil
}

// TemplateGenerator represents the object that performs the template generation.
type TemplateGenerator struct {
	ClassicMode bool
	Translator  *i18n.Translator
}

// InitializeTemplateGenerator creates a new template generator object
func InitializeTemplateGenerator(ctx Context, classicMode bool) (*TemplateGenerator, error) {
	t := &TemplateGenerator{
		ClassicMode: classicMode,
		Translator:  ctx.Translator,
	}

	if err := t.verifyFiles(); err != nil {
		return nil, err
	}

	return t, nil
}

// GenerateTemplate generates the template from the API Model
func (t *TemplateGenerator) GenerateTemplate(containerService *api.ContainerService) (templateRaw string, parametersRaw string, certsGenerated bool, err error) {
	// named return values are used in order to set err in case of a panic
	templateRaw = ""
	parametersRaw = ""
	err = nil

	var templ *template.Template

	properties := containerService.Properties

	if certsGenerated, err = SetPropertiesDefaults(containerService); err != nil {
		return templateRaw, parametersRaw, certsGenerated, err
	}

	templ = template.New("acs template").Funcs(t.getTemplateFuncMap(containerService))

	files, baseFile, e := t.prepareTemplateFiles(properties)
	if e != nil {
		return "", "", false, e
	}

	for _, file := range files {
		bytes, e := Asset(file)
		if e != nil {
			err = t.Translator.Errorf("Error reading file %s, Error: %s", file, e.Error())
			return templateRaw, parametersRaw, certsGenerated, err
		}
		if _, err = templ.New(file).Parse(string(bytes)); err != nil {
			return templateRaw, parametersRaw, certsGenerated, err
		}
	}
	// template generation may have panics in the called functions.  This catches those panics
	// and ensures the panic is returned as an error
	defer func() {
		if r := recover(); r != nil {
			s := debug.Stack()
			err = fmt.Errorf("%v - %s", r, s)

			// invalidate the template and the parameters
			templateRaw = ""
			parametersRaw = ""
		}
	}()

	var b bytes.Buffer
	if err = templ.ExecuteTemplate(&b, baseFile, properties); err != nil {
		return templateRaw, parametersRaw, certsGenerated, err
	}
	templateRaw = b.String()

	var parametersMap paramsMap
	if parametersMap, err = getParameters(containerService, t.ClassicMode); err != nil {
		return templateRaw, parametersRaw, certsGenerated, err
	}
	var parameterBytes []byte
	if parameterBytes, err = json.Marshal(parametersMap); err != nil {
		return templateRaw, parametersRaw, certsGenerated, err
	}
	parametersRaw = string(parameterBytes)

	return templateRaw, parametersRaw, certsGenerated, err
}

// GenerateClusterID creates a unique 8 string cluster ID
func GenerateClusterID(properties *api.Properties) string {
	uniqueNameSuffixSize := 8
	// the name suffix uniquely identifies the cluster and is generated off a hash
	// from the master dns name
	h := fnv.New64a()
	if properties.MasterProfile != nil {
		h.Write([]byte(properties.MasterProfile.DNSPrefix))
	} else if properties.HostedMasterProfile != nil {
		h.Write([]byte(properties.HostedMasterProfile.DNSPrefix))
	} else {
		h.Write([]byte(properties.AgentPoolProfiles[0].Name))
	}
	rand.Seed(int64(h.Sum64()))
	return fmt.Sprintf("%08d", rand.Uint32())[:uniqueNameSuffixSize]
}

// GenerateKubeConfig returns a JSON string representing the KubeConfig
func GenerateKubeConfig(properties *api.Properties, location string) (string, error) {
	b, err := Asset(kubeConfigJSON)
	if err != nil {
		return "", fmt.Errorf("error reading kube config template file %s: %s", kubeConfigJSON, err.Error())
	}
	kubeconfig := string(b)
	// variable replacement
	kubeconfig = strings.Replace(kubeconfig, "{{WrapAsVerbatim \"variables('caCertificate')\"}}", base64.StdEncoding.EncodeToString([]byte(properties.CertificateProfile.CaCertificate)), -1)
	kubeconfig = strings.Replace(kubeconfig, "{{WrapAsVerbatim \"reference(concat('Microsoft.Network/publicIPAddresses/', variables('masterPublicIPAddressName'))).dnsSettings.fqdn\"}}", FormatAzureProdFQDN(properties.MasterProfile.DNSPrefix, location), -1)
	kubeconfig = strings.Replace(kubeconfig, "{{WrapAsVariable \"resourceGroup\"}}", properties.MasterProfile.DNSPrefix, -1)
	kubeconfig = strings.Replace(kubeconfig, "{{WrapAsVerbatim \"variables('kubeConfigCertificate')\"}}", base64.StdEncoding.EncodeToString([]byte(properties.CertificateProfile.KubeConfigCertificate)), -1)
	kubeconfig = strings.Replace(kubeconfig, "{{WrapAsVerbatim \"variables('kubeConfigPrivateKey')\"}}", base64.StdEncoding.EncodeToString([]byte(properties.CertificateProfile.KubeConfigPrivateKey)), -1)

	return kubeconfig, nil
}

func (t *TemplateGenerator) prepareTemplateFiles(properties *api.Properties) ([]string, string, error) {
	var files []string
	var baseFile string
	switch properties.OrchestratorProfile.OrchestratorType {
	case api.DCOS:
		files = append(commonTemplateFiles, dcosTemplateFiles...)
		baseFile = dcosBaseFile
	case api.Swarm:
		files = append(commonTemplateFiles, swarmTemplateFiles...)
		baseFile = swarmBaseFile
	case api.Kubernetes:
		files = append(commonTemplateFiles, kubernetesTemplateFiles...)
		baseFile = kubernetesBaseFile
	case api.SwarmMode:
		files = append(commonTemplateFiles, swarmModeTemplateFiles...)
		baseFile = swarmBaseFile
	default:
		return nil, "", t.Translator.Errorf("orchestrator '%s' is unsupported", properties.OrchestratorProfile.OrchestratorType)
	}

	return files, baseFile, nil
}

//GetCloudSpecConfig returns the kubenernetes container images url configurations based on the deploy target environment
//for example: if the target is the public azure, then the default container image url should be gcrio.azureedge.net/google_container/...
//if the target is azure china, then the default container image should be mirror.azure.cn:5000/google_container/...
func GetCloudSpecConfig(location string) AzureEnvironmentSpecConfig {
	switch GetCloudTargetEnv(location) {
	case azureChinaCloud:
		return AzureChinaCloudSpec
	//TODO - add cloud specs for germany and usgov
	default:
		return AzureCloudSpec
	}
}

// GetCloudTargetEnv determines and returns whether the region is a sovereign cloud which
// have their own data compliance regulations (China/Germany/USGov) or standard
//  Azure public cloud
func GetCloudTargetEnv(location string) string {
	loc := strings.ToLower(strings.Join(strings.Fields(location), ""))
	switch {
	case loc == "chinaeast" || loc == "chinanorth":
		return azureChinaCloud
	case loc == "germanynortheast" || loc == "germanycentral":
		return azureGermanCloud
	case strings.HasPrefix(loc, "usgov") || strings.HasPrefix(loc, "usdod"):
		return azureUSGovernmentCloud
	default:
		return azurePublicCloud
	}
}

func getParameters(cs *api.ContainerService, isClassicMode bool) (paramsMap, error) {
	properties := cs.Properties
	location := cs.Location
	parametersMap := paramsMap{}

	// Master Parameters
	addValue(parametersMap, "location", location)
	addValue(parametersMap, "targetEnvironment", GetCloudTargetEnv(location))
	addValue(parametersMap, "linuxAdminUsername", properties.LinuxProfile.AdminUsername)
	// masterEndpointDNSNamePrefix is the basis for storage account creation across dcos, swarm, and k8s
	if properties.MasterProfile != nil {
		// MasterProfile exists, uses master DNS prefix
		addValue(parametersMap, "masterEndpointDNSNamePrefix", properties.MasterProfile.DNSPrefix)
	} else if properties.HostedMasterProfile != nil {
		// Agents only, use cluster DNS prefix
		addValue(parametersMap, "masterEndpointDNSNamePrefix", properties.HostedMasterProfile.DNSPrefix)
	}
	if properties.MasterProfile != nil {
		if properties.MasterProfile.IsCustomVNET() {
			addValue(parametersMap, "masterVnetSubnetID", properties.MasterProfile.VnetSubnetID)
		} else {
			addValue(parametersMap, "masterSubnet", properties.MasterProfile.Subnet)
		}
		addValue(parametersMap, "firstConsecutiveStaticIP", properties.MasterProfile.FirstConsecutiveStaticIP)
		addValue(parametersMap, "masterVMSize", properties.MasterProfile.VMSize)
		if isClassicMode {
			addValue(parametersMap, "masterCount", properties.MasterProfile.Count)
		}
	}
	addValue(parametersMap, "sshRSAPublicKey", properties.LinuxProfile.SSH.PublicKeys[0].KeyData)
	for i, s := range properties.LinuxProfile.Secrets {
		addValue(parametersMap, fmt.Sprintf("linuxKeyVaultID%d", i), s.SourceVault.ID)
		for j, c := range s.VaultCertificates {
			addValue(parametersMap, fmt.Sprintf("linuxKeyVaultID%dCertificateURL%d", i, j), c.CertificateURL)
		}
	}

	cloudSpecConfig := GetCloudSpecConfig(location)
	// Kubernetes Parameters
	if properties.OrchestratorProfile.OrchestratorType == api.Kubernetes {
		KubernetesRelease := properties.OrchestratorProfile.OrchestratorRelease

		kubernetesHyperkubeSpec := properties.OrchestratorProfile.KubernetesConfig.KubernetesImageBase + KubeConfigs[KubernetesRelease]["hyperkube"]
		if properties.OrchestratorProfile.KubernetesConfig.CustomHyperkubeImage != "" {
			kubernetesHyperkubeSpec = properties.OrchestratorProfile.KubernetesConfig.CustomHyperkubeImage
		}

		if properties.CertificateProfile != nil {
			addSecret(parametersMap, "apiServerCertificate", properties.CertificateProfile.APIServerCertificate, true)
			addSecret(parametersMap, "apiServerPrivateKey", properties.CertificateProfile.APIServerPrivateKey, true)
			addSecret(parametersMap, "caCertificate", properties.CertificateProfile.CaCertificate, true)
			addSecret(parametersMap, "caPrivateKey", properties.CertificateProfile.CaPrivateKey, true)
			addSecret(parametersMap, "clientCertificate", properties.CertificateProfile.ClientCertificate, true)
			addSecret(parametersMap, "clientPrivateKey", properties.CertificateProfile.ClientPrivateKey, true)
			addSecret(parametersMap, "kubeConfigCertificate", properties.CertificateProfile.KubeConfigCertificate, true)
			addSecret(parametersMap, "kubeConfigPrivateKey", properties.CertificateProfile.KubeConfigPrivateKey, true)
		}
		if properties.HostedMasterProfile != nil && properties.HostedMasterProfile.FQDN != "" {
			addValue(parametersMap, "kubernetesEndpoint", properties.HostedMasterProfile.FQDN)
		}
		addValue(parametersMap, "dockerEngineDownloadRepo", cloudSpecConfig.DockerSpecConfig.DockerEngineRepo)
		addValue(parametersMap, "kubernetesHyperkubeSpec", kubernetesHyperkubeSpec)
		addValue(parametersMap, "kubernetesAddonManagerSpec", cloudSpecConfig.KubernetesSpecConfig.KubernetesImageBase+KubeConfigs[KubernetesRelease]["addonmanager"])
		addValue(parametersMap, "kubernetesAddonResizerSpec", cloudSpecConfig.KubernetesSpecConfig.KubernetesImageBase+KubeConfigs[KubernetesRelease]["addonresizer"])
		addValue(parametersMap, "kubernetesDashboardSpec", cloudSpecConfig.KubernetesSpecConfig.KubernetesImageBase+KubeConfigs[KubernetesRelease]["dashboard"])
		addValue(parametersMap, "kubernetesDNSMasqSpec", cloudSpecConfig.KubernetesSpecConfig.KubernetesImageBase+KubeConfigs[KubernetesRelease]["dnsmasq"])
		addValue(parametersMap, "kubernetesExecHealthzSpec", cloudSpecConfig.KubernetesSpecConfig.KubernetesImageBase+KubeConfigs[KubernetesRelease]["exechealthz"])
		addValue(parametersMap, "kubernetesHeapsterSpec", cloudSpecConfig.KubernetesSpecConfig.KubernetesImageBase+KubeConfigs[KubernetesRelease]["heapster"])
		addValue(parametersMap, "kubernetesTillerSpec", cloudSpecConfig.KubernetesSpecConfig.TillerImageBase+KubeConfigs[KubernetesRelease]["tiller"])
		addValue(parametersMap, "kubernetesKubeDNSSpec", cloudSpecConfig.KubernetesSpecConfig.KubernetesImageBase+KubeConfigs[KubernetesRelease]["dns"])
		addValue(parametersMap, "kubernetesPodInfraContainerSpec", cloudSpecConfig.KubernetesSpecConfig.KubernetesImageBase+KubeConfigs[KubernetesRelease]["pause"])
		addValue(parametersMap, "kubernetesNodeStatusUpdateFrequency", properties.OrchestratorProfile.KubernetesConfig.NodeStatusUpdateFrequency)
		addValue(parametersMap, "kubernetesCtrlMgrNodeMonitorGracePeriod", properties.OrchestratorProfile.KubernetesConfig.CtrlMgrNodeMonitorGracePeriod)
		addValue(parametersMap, "kubernetesCtrlMgrPodEvictionTimeout", properties.OrchestratorProfile.KubernetesConfig.CtrlMgrPodEvictionTimeout)
		addValue(parametersMap, "kubernetesCtrlMgrRouteReconciliationPeriod", properties.OrchestratorProfile.KubernetesConfig.CtrlMgrRouteReconciliationPeriod)
		addValue(parametersMap, "cloudProviderBackoff", strconv.FormatBool(properties.OrchestratorProfile.KubernetesConfig.CloudProviderBackoff))
		addValue(parametersMap, "cloudProviderBackoffRetries", strconv.Itoa(properties.OrchestratorProfile.KubernetesConfig.CloudProviderBackoffRetries))
		addValue(parametersMap, "cloudProviderBackoffExponent", strconv.FormatFloat(properties.OrchestratorProfile.KubernetesConfig.CloudProviderBackoffExponent, 'f', -1, 64))
		addValue(parametersMap, "cloudProviderBackoffDuration", strconv.Itoa(properties.OrchestratorProfile.KubernetesConfig.CloudProviderBackoffDuration))
		addValue(parametersMap, "cloudProviderBackoffJitter", strconv.FormatFloat(properties.OrchestratorProfile.KubernetesConfig.CloudProviderBackoffJitter, 'f', -1, 64))
		addValue(parametersMap, "cloudProviderRatelimit", strconv.FormatBool(properties.OrchestratorProfile.KubernetesConfig.CloudProviderRateLimit))
		addValue(parametersMap, "cloudProviderRatelimitQPS", strconv.FormatFloat(properties.OrchestratorProfile.KubernetesConfig.CloudProviderRateLimitQPS, 'f', -1, 64))
		addValue(parametersMap, "cloudProviderRatelimitBucket", strconv.Itoa(properties.OrchestratorProfile.KubernetesConfig.CloudProviderRateLimitBucket))
		addValue(parametersMap, "kubeClusterCidr", properties.OrchestratorProfile.KubernetesConfig.ClusterSubnet)
		addValue(parametersMap, "dockerBridgeCidr", properties.OrchestratorProfile.KubernetesConfig.DockerBridgeSubnet)
		addValue(parametersMap, "networkPolicy", properties.OrchestratorProfile.KubernetesConfig.NetworkPolicy)
		if properties.OrchestratorProfile.KubernetesConfig == nil ||
			!properties.OrchestratorProfile.KubernetesConfig.UseManagedIdentity {
			addValue(parametersMap, "servicePrincipalClientId", properties.ServicePrincipalProfile.ClientID)
			if properties.ServicePrincipalProfile.KeyvaultSecretRef != "" {
				addSecret(parametersMap, "servicePrincipalClientSecret", properties.ServicePrincipalProfile.KeyvaultSecretRef, false)
			} else {
				addValue(parametersMap, "servicePrincipalClientSecret", properties.ServicePrincipalProfile.Secret)
			}
		}
	}

	if strings.HasPrefix(properties.OrchestratorProfile.OrchestratorType, api.DCOS) {
		dcosBootstrapURL := cloudSpecConfig.DCOSSpecConfig.DCOS188BootstrapDownloadURL
		switch properties.OrchestratorProfile.OrchestratorType {
		case api.DCOS:
			switch properties.OrchestratorProfile.OrchestratorRelease {
			case api.DCOSRelease1Dot7:
				dcosBootstrapURL = cloudSpecConfig.DCOSSpecConfig.DCOS173BootstrapDownloadURL
			case api.DCOSRelease1Dot8:
				dcosBootstrapURL = cloudSpecConfig.DCOSSpecConfig.DCOS188BootstrapDownloadURL
			case api.DCOSRelease1Dot9:
				dcosBootstrapURL = cloudSpecConfig.DCOSSpecConfig.DCOS190BootstrapDownloadURL
			}
		}
		addValue(parametersMap, "dcosBootstrapURL", dcosBootstrapURL)

		dcosWindowsBootstrapURL := cloudSpecConfig.DCOSSpecConfig.DCOSWindowsBootstrapDownloadURL
		addValue(parametersMap, "dcosWindowsBootstrapURL", dcosWindowsBootstrapURL)
	}

	// Agent parameters
	for _, agentProfile := range properties.AgentPoolProfiles {
		addValue(parametersMap, fmt.Sprintf("%sCount", agentProfile.Name), agentProfile.Count)
		addValue(parametersMap, fmt.Sprintf("%sVMSize", agentProfile.Name), agentProfile.VMSize)
		if agentProfile.IsCustomVNET() {
			addValue(parametersMap, fmt.Sprintf("%sVnetSubnetID", agentProfile.Name), agentProfile.VnetSubnetID)
		} else {
			addValue(parametersMap, fmt.Sprintf("%sSubnet", agentProfile.Name), agentProfile.Subnet)
		}
		if len(agentProfile.Ports) > 0 {
			addValue(parametersMap, fmt.Sprintf("%sEndpointDNSNamePrefix", agentProfile.Name), agentProfile.DNSPrefix)
		}
	}

	// Windows parameters
	if properties.HasWindows() {
		addValue(parametersMap, "windowsAdminUsername", properties.WindowsProfile.AdminUsername)
		addSecret(parametersMap, "windowsAdminPassword", properties.WindowsProfile.AdminPassword, false)
		if properties.OrchestratorProfile.OrchestratorType == api.Kubernetes {
			KubernetesRelease := properties.OrchestratorProfile.OrchestratorRelease
			addValue(parametersMap, "kubeBinariesSASURL", cloudSpecConfig.KubernetesSpecConfig.KubeBinariesSASURLBase+KubeConfigs[KubernetesRelease]["windowszip"])
			addValue(parametersMap, "kubeBinariesVersion", api.KubernetesReleaseToVersion[KubernetesRelease])
		}
		for i, s := range properties.WindowsProfile.Secrets {
			addValue(parametersMap, fmt.Sprintf("windowsKeyVaultID%d", i), s.SourceVault.ID)
			for j, c := range s.VaultCertificates {
				addValue(parametersMap, fmt.Sprintf("windowsKeyVaultID%dCertificateURL%d", i, j), c.CertificateURL)
				addValue(parametersMap, fmt.Sprintf("windowsKeyVaultID%dCertificateStore%d", i, j), c.CertificateStore)
			}
		}
	}

	return parametersMap, nil
}

func addValue(m paramsMap, k string, v interface{}) {
	m[k] = paramsMap{
		"value": v,
	}
}

func addSecret(m paramsMap, k string, v interface{}, encode bool) {
	str, ok := v.(string)
	if !ok {
		addValue(m, k, v)
		return
	}
	parts := keyvaultSecretPathRe.FindStringSubmatch(str)
	if parts == nil || len(parts) != 5 {
		if encode {
			addValue(m, k, base64.StdEncoding.EncodeToString([]byte(str)))
		} else {
			addValue(m, k, str)
		}
		return
	}

	m[k] = paramsMap{
		"reference": &KeyVaultRef{
			KeyVault: KeyVaultID{
				ID: parts[1],
			},
			SecretName:    parts[2],
			SecretVersion: parts[4],
		},
	}
}

// getStorageAccountType returns the support managed disk storage tier for a give VM size
func getStorageAccountType(sizeName string) (string, error) {
	spl := strings.Split(sizeName, "_")
	if len(spl) < 2 {
		return "", fmt.Errorf("Invalid sizeName: %s", sizeName)
	}
	capability := spl[1]
	if strings.Contains(strings.ToLower(capability), "s") {
		return "Premium_LRS", nil
	}
	return "Standard_LRS", nil
}

// getTemplateFuncMap returns all functions used in template generation
func (t *TemplateGenerator) getTemplateFuncMap(cs *api.ContainerService) template.FuncMap {
	return template.FuncMap{
		"IsHostedMaster": func() bool {
			return cs.Properties.HostedMasterProfile != nil
		},
		"IsDCOS19": func() bool {
			return cs.Properties.OrchestratorProfile.OrchestratorType == api.DCOS &&
				cs.Properties.OrchestratorProfile.OrchestratorRelease == api.DCOSRelease1Dot9
		},
		"IsKubernetesVersionGe": func(version string) bool {
			orchestratorVersion, _ := semver.NewVersion(cs.Properties.OrchestratorProfile.OrchestratorVersion)
			constraint, _ := semver.NewConstraint(">=" + version)
			return cs.Properties.OrchestratorProfile.OrchestratorType == api.Kubernetes && constraint.Check(orchestratorVersion)
		},
		"GetKubernetesLabels": func(profile *api.AgentPoolProfile) string {
			var buf bytes.Buffer
			buf.WriteString(fmt.Sprintf("role=agent,agentpool=%s", profile.Name))
			if profile.StorageProfile == api.ManagedDisks {
				storagetier, _ := getStorageAccountType(profile.VMSize)
				buf.WriteString(fmt.Sprintf(",storageprofile=managed,storagetier=%s", storagetier))
			}
			for k, v := range profile.CustomNodeLabels {
				buf.WriteString(fmt.Sprintf(",%s=%s", k, v))
			}

			return buf.String()
		},
		"RequiresFakeAgentOutput": func() bool {
			return cs.Properties.OrchestratorProfile.OrchestratorType == api.Kubernetes
		},
		"IsSwarmMode": func() bool {
			return cs.Properties.OrchestratorProfile.IsSwarmMode()
		},
		"IsKubernetes": func() bool {
			return cs.Properties.OrchestratorProfile.IsKubernetes()
		},
		"IsPublic": func(ports []int) bool {
			return len(ports) > 0
		},
		"IsVNETIntegrated": func() bool {
			return cs.Properties.OrchestratorProfile.IsVNETIntegrated()
		},
		"UseManagedIdentity": func() bool {
			return cs.Properties.OrchestratorProfile.KubernetesConfig.UseManagedIdentity
		},
		"UseInstanceMetadata": func() bool {
			return cs.Properties.OrchestratorProfile.KubernetesConfig.UseInstanceMetadata
		},
		"GetVNETSubnetDependencies": func() string {
			return getVNETSubnetDependencies(cs.Properties)
		},
		"GetLBRules": func(name string, ports []int) string {
			return getLBRules(name, ports)
		},
		"GetProbes": func(ports []int) string {
			return getProbes(ports)
		},
		"GetSecurityRules": func(ports []int) string {
			return getSecurityRules(ports)
		},
		"GetUniqueNameSuffix": func() string {
			return GenerateClusterID(cs.Properties)
		},
		"GetVNETAddressPrefixes": func() string {
			return getVNETAddressPrefixes(cs.Properties)
		},
		"GetVNETSubnets": func(addNSG bool) string {
			return getVNETSubnets(cs.Properties, addNSG)
		},
		"GetDataDisks": func(profile *api.AgentPoolProfile) string {
			return getDataDisks(profile)
		},
		"GetDCOSMasterCustomData": func() string {
			masterProvisionScript := getDCOSMasterProvisionScript()
			masterAttributeContents := getDCOSMasterCustomNodeLabels()
			str := getSingleLineDCOSCustomData(
				cs.Properties.OrchestratorProfile.OrchestratorType,
				cs.Properties.OrchestratorProfile.OrchestratorRelease,
				cs.Properties.MasterProfile.Count, masterProvisionScript, masterAttributeContents)

			return fmt.Sprintf("\"customData\": \"[base64(concat('#cloud-config\\n\\n', '%s'))]\",", str)
		},
		"GetDCOSAgentCustomData": func(profile *api.AgentPoolProfile) string {
			agentProvisionScript := getDCOSAgentProvisionScript(profile)
			attributeContents := getDCOSAgentCustomNodeLabels(profile)
			str := getSingleLineDCOSCustomData(
				cs.Properties.OrchestratorProfile.OrchestratorType,
				cs.Properties.OrchestratorProfile.OrchestratorRelease,
				cs.Properties.MasterProfile.Count, agentProvisionScript, attributeContents)

			return fmt.Sprintf("\"customData\": \"[base64(concat('#cloud-config\\n\\n', '%s'))]\",", str)
		},
		"GetDCOSWindowsAgentCustomData": func(profile *api.AgentPoolProfile) string {
			str := getBase64CustomScript(dcosWindowsProvision)
			return fmt.Sprintf("\"customData\": \"%s\"", str)
		},
		"GetMasterAllowedSizes": func() string {
			if t.ClassicMode {
				return GetClassicAllowedSizes()
			} else if cs.Properties.OrchestratorProfile.OrchestratorType == api.DCOS {
				return GetDCOSMasterAllowedSizes()
			}
			return GetMasterAgentAllowedSizes()
		},
		"GetAgentAllowedSizes": func() string {
			if t.ClassicMode {
				return GetClassicAllowedSizes()
			} else if cs.Properties.OrchestratorProfile.OrchestratorType == api.Kubernetes {
				return GetKubernetesAgentAllowedSizes()
			}
			return GetMasterAgentAllowedSizes()
		},
		"GetSwarmVersions": func() string {
			return GetSwarmVersions(SwarmVersion, SwarmDockerComposeVersion)
		},
		"GetSwarmModeVersions": func() string {
			return GetSwarmVersions(DockerCEVersion, DockerCEDockerComposeVersion)
		},
		"GetSizeMap": func() string {
			if t.ClassicMode {
				return GetClassicSizeMap()
			}
			return GetSizeMap()
		},
		"GetClassicMode": func() bool {
			return t.ClassicMode
		},
		"Base64": func(s string) string {
			return base64.StdEncoding.EncodeToString([]byte(s))
		},
		"GetDefaultInternalLbStaticIPOffset": func() int {
			return DefaultInternalLbStaticIPOffset
		},
		"GetKubernetesMasterCustomScript": func() string {
			return getBase64CustomScript(kubernetesMasterCustomScript)
		},
		"GetKubernetesMasterCustomData": func(profile *api.Properties) string {
			str, e := t.getSingleLineForTemplate(kubernetesMasterCustomDataYaml, cs, profile)
			if e != nil {
				return ""
			}

			for placeholder, filename := range kubernetesManifestYamls {
				manifestTextContents := getBase64CustomScript(filename)
				str = strings.Replace(str, placeholder, manifestTextContents, -1)
			}

			// add artifacts and addons
			var artifiacts map[string]string
			if profile.OrchestratorProfile.OrchestratorRelease == api.KubernetesRelease1Dot5 {
				artifiacts = kubernetesAritfacts15
			} else {
				artifiacts = kubernetesAritfacts
			}
			for placeholder, filename := range artifiacts {
				addonTextContents := getBase64CustomScript(filename)
				str = strings.Replace(str, placeholder, addonTextContents, -1)
			}

			var addonYamls map[string]string
			if profile.OrchestratorProfile.OrchestratorRelease == api.KubernetesRelease1Dot5 {
				addonYamls = kubernetesAddonYamls15
			} else {
				addonYamls = kubernetesAddonYamls
			}
			for placeholder, filename := range addonYamls {
				addonTextContents := getBase64CustomScript(filename)
				str = strings.Replace(str, placeholder, addonTextContents, -1)
			}

			// add calico manifests
			if profile.OrchestratorProfile.KubernetesConfig.NetworkPolicy == "calico" {
				if profile.OrchestratorProfile.OrchestratorRelease == api.KubernetesRelease1Dot5 {
					calicoAddonYamls = calicoAddonYamls15
				}
				for placeholder, filename := range calicoAddonYamls {
					addonTextContents := getBase64CustomScript(filename)
					str = strings.Replace(str, placeholder, addonTextContents, -1)
				}
			}

			// return the custom data
			return fmt.Sprintf("\"customData\": \"[base64(concat('%s'))]\",", str)
		},
		"GetKubernetesAgentCustomData": func(profile *api.AgentPoolProfile) string {
			str, e := t.getSingleLineForTemplate(kubernetesAgentCustomDataYaml, cs, profile)
			if e != nil {
				return ""
			}

			// add artifacts
			var artifiacts map[string]string
			if cs.Properties.OrchestratorProfile.OrchestratorVersion == api.KubernetesRelease1Dot5 {
				artifiacts = kubernetesAritfacts15
			} else {
				artifiacts = kubernetesAritfacts
			}
			for placeholder, filename := range artifiacts {
				addonTextContents := getBase64CustomScript(filename)
				str = strings.Replace(str, placeholder, addonTextContents, -1)
			}

			return fmt.Sprintf("\"customData\": \"[base64(concat('%s'))]\",", str)
		},
		"WriteLinkedTemplatesForExtensions": func() string {
			extensions := getLinkedTemplatesForExtensions(cs.Properties)
			return extensions
		},
		"GetKubernetesB64Provision": func() string {
			return getBase64CustomScript(kubernetesMasterCustomScript)
		},
		"GetMasterSwarmCustomData": func() string {
			files := []string{swarmProvision}
			str := buildYamlFileWithWriteFiles(files)
			str = escapeSingleLine(str)
			return fmt.Sprintf("\"customData\": \"[base64('%s')]\",", str)
		},
		"GetAgentSwarmCustomData": func() string {
			files := []string{swarmProvision}
			str := buildYamlFileWithWriteFiles(files)
			str = escapeSingleLine(str)
			return fmt.Sprintf("\"customData\": \"[base64(concat('%s',variables('agentRunCmdFile'),variables('agentRunCmd')))]\",", str)
		},
		"GetLocation": func() string {
			return cs.Location
		},
		"GetWinAgentSwarmCustomData": func() string {
			str := getBase64CustomScript(swarmWindowsProvision)
			return fmt.Sprintf("\"customData\": \"%s\"", str)
		},
		"GetWinAgentSwarmModeCustomData": func() string {
			str := getBase64CustomScript(swarmModeWindowsProvision)
			return fmt.Sprintf("\"customData\": \"%s\"", str)
		},
		"GetKubernetesWindowsAgentCustomData": func(profile *api.AgentPoolProfile) string {
			str, e := t.getSingleLineForTemplate(kubernetesWindowsAgentCustomDataPS1, cs, profile)
			if e != nil {
				return ""
			}
			return fmt.Sprintf("\"customData\": \"[base64(concat('%s'))]\",", str)
		},
		"GetMasterSwarmModeCustomData": func() string {
			files := []string{swarmModeProvision}
			str := buildYamlFileWithWriteFiles(files)
			str = escapeSingleLine(str)
			return fmt.Sprintf("\"customData\": \"[base64('%s')]\",", str)
		},
		"GetAgentSwarmModeCustomData": func() string {
			files := []string{swarmModeProvision}
			str := buildYamlFileWithWriteFiles(files)
			str = escapeSingleLine(str)
			return fmt.Sprintf("\"customData\": \"[base64(concat('%s',variables('agentRunCmdFile'),variables('agentRunCmd')))]\",", str)
		},
		"GetKubernetesSubnets": func() string {
			return getKubernetesSubnets(cs.Properties)
		},
		"GetKubernetesPodStartIndex": func() string {
			return fmt.Sprintf("%d", getKubernetesPodStartIndex(cs.Properties))
		},
		"WrapAsVariable": func(s string) string {
			return fmt.Sprintf("',variables('%s'),'", s)
		},
		"WrapAsVerbatim": func(s string) string {
			return fmt.Sprintf("',%s,'", s)
		},
		"AnyAgentUsesAvailablilitySets": func() bool {
			for _, agentProfile := range cs.Properties.AgentPoolProfiles {
				if agentProfile.IsAvailabilitySets() {
					return true
				}
			}
			return false
		},
		"HasLinuxAgents": func() bool {
			for _, agentProfile := range cs.Properties.AgentPoolProfiles {
				if agentProfile.IsLinux() {
					return true
				}
			}
			return false
		},
		"HasLinuxSecrets": func() bool {
			return cs.Properties.LinuxProfile.HasSecrets()
		},
		"HasWindowsSecrets": func() bool {
			return cs.Properties.WindowsProfile.HasSecrets()
		},
		"PopulateClassicModeDefaultValue": func(attr string) string {
			var val string
			if !t.ClassicMode {
				val = ""
			} else {
				kubernetesRelease := cs.Properties.OrchestratorProfile.OrchestratorRelease
				cloudSpecConfig := GetCloudSpecConfig(cs.Location)
				switch attr {
				case "kubernetesHyperkubeSpec":
					val = cs.Properties.OrchestratorProfile.KubernetesConfig.KubernetesImageBase + KubeConfigs[kubernetesRelease]["hyperkube"]
					if cs.Properties.OrchestratorProfile.KubernetesConfig.CustomHyperkubeImage != "" {
						val = cs.Properties.OrchestratorProfile.KubernetesConfig.CustomHyperkubeImage
					}
				case "kubernetesAddonManagerSpec":
					val = cloudSpecConfig.KubernetesSpecConfig.KubernetesImageBase + KubeConfigs[kubernetesRelease]["addonmanager"]
				case "kubernetesAddonResizerSpec":
					val = cloudSpecConfig.KubernetesSpecConfig.KubernetesImageBase + KubeConfigs[kubernetesRelease]["addonresizer"]
				case "kubernetesDashboardSpec":
					val = cloudSpecConfig.KubernetesSpecConfig.KubernetesImageBase + KubeConfigs[kubernetesRelease]["dashboard"]
				case "kubernetesDNSMasqSpec":
					val = cloudSpecConfig.KubernetesSpecConfig.KubernetesImageBase + KubeConfigs[kubernetesRelease]["dnsmasq"]
				case "kubernetesExecHealthzSpec":
					val = cloudSpecConfig.KubernetesSpecConfig.KubernetesImageBase + KubeConfigs[kubernetesRelease]["exechealthz"]
				case "kubernetesHeapsterSpec":
					val = cloudSpecConfig.KubernetesSpecConfig.KubernetesImageBase + KubeConfigs[kubernetesRelease]["heapster"]
				case "kubernetesTillerSpec":
					val = cloudSpecConfig.KubernetesSpecConfig.TillerImageBase + KubeConfigs[kubernetesRelease]["tiller"]
				case "kubernetesKubeDNSSpec":
					val = cloudSpecConfig.KubernetesSpecConfig.KubernetesImageBase + KubeConfigs[kubernetesRelease]["dns"]
				case "kubernetesPodInfraContainerSpec":
					val = cloudSpecConfig.KubernetesSpecConfig.KubernetesImageBase + KubeConfigs[kubernetesRelease]["pause"]
				case "kubernetesNodeStatusUpdateFrequency":
					val = cs.Properties.OrchestratorProfile.KubernetesConfig.NodeStatusUpdateFrequency
				case "kubernetesCtrlMgrNodeMonitorGracePeriod":
					val = cs.Properties.OrchestratorProfile.KubernetesConfig.CtrlMgrNodeMonitorGracePeriod
				case "kubernetesCtrlMgrPodEvictionTimeout":
					val = cs.Properties.OrchestratorProfile.KubernetesConfig.CtrlMgrPodEvictionTimeout
				case "kubernetesCtrlMgrRouteReconciliationPeriod":
					val = cs.Properties.OrchestratorProfile.KubernetesConfig.CtrlMgrRouteReconciliationPeriod
				case "cloudProviderBackoff":
					val = strconv.FormatBool(cs.Properties.OrchestratorProfile.KubernetesConfig.CloudProviderBackoff)
				case "cloudProviderBackoffRetries":
					val = strconv.Itoa(cs.Properties.OrchestratorProfile.KubernetesConfig.CloudProviderBackoffRetries)
				case "cloudProviderBackoffExponent":
					val = strconv.FormatFloat(cs.Properties.OrchestratorProfile.KubernetesConfig.CloudProviderBackoffExponent, 'f', -1, 64)
				case "cloudProviderBackoffDuration":
					val = strconv.Itoa(cs.Properties.OrchestratorProfile.KubernetesConfig.CloudProviderBackoffDuration)
				case "cloudProviderBackoffJitter":
					val = strconv.FormatFloat(cs.Properties.OrchestratorProfile.KubernetesConfig.CloudProviderBackoffJitter, 'f', -1, 64)
				case "cloudProviderRatelimit":
					val = strconv.FormatBool(cs.Properties.OrchestratorProfile.KubernetesConfig.CloudProviderRateLimit)
				case "cloudProviderRatelimitQPS":
					val = strconv.FormatFloat(cs.Properties.OrchestratorProfile.KubernetesConfig.CloudProviderRateLimitQPS, 'f', -1, 64)
				case "cloudProviderRatelimitBucket":
					val = strconv.Itoa(cs.Properties.OrchestratorProfile.KubernetesConfig.CloudProviderRateLimitBucket)
				case "kubeBinariesSASURL":
					val = cloudSpecConfig.KubernetesSpecConfig.KubeBinariesSASURLBase + KubeConfigs[kubernetesRelease]["windowszip"]
				case "kubeClusterCidr":
					val = "10.244.0.0/16"
				case "kubeBinariesVersion":
					val = api.KubernetesReleaseToVersion[cs.Properties.OrchestratorProfile.OrchestratorRelease]
				case "caPrivateKey":
					// The base64 encoded "NotAvailable"
					val = "Tm90QXZhaWxhYmxlCg=="
				case "dockerBridgeCidr":
					val = DefaultDockerBridgeSubnet
				default:
					val = ""
				}
			}
			return fmt.Sprintf("\"defaultValue\": \"%s\",", val)
		},
		// inspired by http://stackoverflow.com/questions/18276173/calling-a-template-with-several-pipeline-parameters/18276968#18276968
		"dict": func(values ...interface{}) (map[string]interface{}, error) {
			if len(values)%2 != 0 {
				return nil, errors.New("invalid dict call")
			}
			dict := make(map[string]interface{}, len(values)/2)
			for i := 0; i < len(values); i += 2 {
				key, ok := values[i].(string)
				if !ok {
					return nil, errors.New("dict keys must be strings")
				}
				dict[key] = values[i+1]
			}
			return dict, nil
		},
		"loop": func(min, max int) []int {
			var s []int
			for i := min; i <= max; i++ {
				s = append(s, i)
			}
			return s
		},
	}
}

func getPackageGUID(orchestratorType string, orchestratorRelease string, masterCount int) string {
	if orchestratorType == api.DCOS {
		switch orchestratorRelease {
		case api.DCOSRelease1Dot9:
			switch masterCount {
			case 1:
				return "bcc883b7a3191412cf41824bdee06c1142187a0b"
			case 3:
				return "dcff7e24c0c1827bebeb7f1a806f558054481b33"
			case 5:
				return "b41bfa84137a6374b2ff5eb1655364d7302bd257"
			}
		case api.DCOSRelease1Dot8:
			switch masterCount {
			case 1:
				return "441385ce2f5942df7e29075c12fb38fa5e92cbba"
			case 3:
				return "b1cd359287504efb780257bd12cc3a63704e42d4"
			case 5:
				return "d9b61156dfcc9383e014851529738aa550ef57d9"
			}
		case api.DCOSRelease1Dot7:
			switch masterCount {
			case 1:
				return "6b604c1331c2b8b52bb23d1ea8a8d17e0f2b7428"
			case 3:
				return "6af5097e7956962a3d4318d28fbf280a47305485"
			case 5:
				return "376e07e0dbad2af3da2c03bc92bb07e84b3dafd5"
			}
		}
	}
	return ""
}

func getDCOSCustomDataPublicIPStr(orchestratorType string, masterCount int) string {
	if orchestratorType == api.DCOS {
		var buf bytes.Buffer
		for i := 0; i < masterCount; i++ {
			buf.WriteString(fmt.Sprintf("reference(variables('masterVMNic')[%d]).ipConfigurations[0].properties.privateIPAddress,", i))
			if i < (masterCount - 1) {
				buf.WriteString(`'\\\", \\\"', `)
			}
		}
		return buf.String()
	}
	return ""
}

func getDCOSMasterCustomNodeLabels() string {
	// return empty string for DCOS since no attribtutes needed on master
	return ""
}

func getDCOSAgentCustomNodeLabels(profile *api.AgentPoolProfile) string {
	var buf bytes.Buffer
	buf.WriteString("")
	if len(profile.CustomNodeLabels) > 0 {
		buf.WriteString("MESOS_ATTRIBUTES=")
		for k, v := range profile.CustomNodeLabels {
			buf.WriteString(fmt.Sprintf("%s:%s;", k, v))
		}
	}
	return buf.String()
}

func getVNETAddressPrefixes(properties *api.Properties) string {
	visitedSubnets := make(map[string]bool)
	var buf bytes.Buffer
	buf.WriteString(`"[variables('masterSubnet')]"`)
	visitedSubnets[properties.MasterProfile.Subnet] = true
	for _, profile := range properties.AgentPoolProfiles {
		if _, ok := visitedSubnets[profile.Subnet]; !ok {
			buf.WriteString(fmt.Sprintf(",\n            \"[variables('%sSubnet')]\"", profile.Name))
		}
	}
	return buf.String()
}

func getVNETSubnetDependencies(properties *api.Properties) string {
	agentString := `        "[concat('Microsoft.Network/networkSecurityGroups/', variables('%sNSGName'))]"`
	var buf bytes.Buffer
	for index, agentProfile := range properties.AgentPoolProfiles {
		if index > 0 {
			buf.WriteString(",\n")
		}
		buf.WriteString(fmt.Sprintf(agentString, agentProfile.Name))
	}
	return buf.String()
}

func getVNETSubnets(properties *api.Properties, addNSG bool) string {
	masterString := `{
            "name": "[variables('masterSubnetName')]",
            "properties": {
              "addressPrefix": "[variables('masterSubnet')]"
            }
          }`
	agentString := `          {
            "name": "[variables('%sSubnetName')]",
            "properties": {
              "addressPrefix": "[variables('%sSubnet')]"
            }
          }`
	agentStringNSG := `          {
            "name": "[variables('%sSubnetName')]",
            "properties": {
              "addressPrefix": "[variables('%sSubnet')]",
              "networkSecurityGroup": {
                "id": "[resourceId('Microsoft.Network/networkSecurityGroups', variables('%sNSGName'))]"
              }
            }
          }`
	var buf bytes.Buffer
	buf.WriteString(masterString)
	for _, agentProfile := range properties.AgentPoolProfiles {
		buf.WriteString(",\n")
		if addNSG {
			buf.WriteString(fmt.Sprintf(agentStringNSG, agentProfile.Name, agentProfile.Name, agentProfile.Name))
		} else {
			buf.WriteString(fmt.Sprintf(agentString, agentProfile.Name, agentProfile.Name))
		}

	}
	return buf.String()
}

func getLBRule(name string, port int) string {
	return fmt.Sprintf(`	          {
            "name": "LBRule%d",
            "properties": {
              "backendAddressPool": {
                "id": "[concat(variables('%sLbID'), '/backendAddressPools/', variables('%sLbBackendPoolName'))]"
              },
              "backendPort": %d,
              "enableFloatingIP": false,
              "frontendIPConfiguration": {
                "id": "[variables('%sLbIPConfigID')]"
              },
              "frontendPort": %d,
              "idleTimeoutInMinutes": 5,
              "loadDistribution": "Default",
              "probe": {
                "id": "[concat(variables('%sLbID'),'/probes/tcp%dProbe')]"
              },
              "protocol": "tcp"
            }
          }`, port, name, name, port, name, port, name, port)
}

func getLBRules(name string, ports []int) string {
	var buf bytes.Buffer
	for index, port := range ports {
		if index > 0 {
			buf.WriteString(",\n")
		}
		buf.WriteString(getLBRule(name, port))
	}
	return buf.String()
}

func getProbe(port int) string {
	return fmt.Sprintf(`          {
            "name": "tcp%dProbe",
            "properties": {
              "intervalInSeconds": "5",
              "numberOfProbes": "2",
              "port": %d,
              "protocol": "tcp"
            }
          }`, port, port)
}

func getProbes(ports []int) string {
	var buf bytes.Buffer
	for index, port := range ports {
		if index > 0 {
			buf.WriteString(",\n")
		}
		buf.WriteString(getProbe(port))
	}
	return buf.String()
}

func getSecurityRule(port int, portIndex int) string {
	// BaseLBPriority specifies the base lb priority.
	BaseLBPriority := 200
	return fmt.Sprintf(`          {
            "name": "Allow_%d",
            "properties": {
              "access": "Allow",
              "description": "Allow traffic from the Internet to port %d",
              "destinationAddressPrefix": "*",
              "destinationPortRange": "%d",
              "direction": "Inbound",
              "priority": %d,
              "protocol": "*",
              "sourceAddressPrefix": "Internet",
              "sourcePortRange": "*"
            }
          }`, port, port, port, BaseLBPriority+portIndex)
}

func getDataDisks(a *api.AgentPoolProfile) string {
	if !a.HasDisks() {
		return ""
	}
	var buf bytes.Buffer
	buf.WriteString("\"dataDisks\": [\n")
	dataDisks := `            {
              "createOption": "Empty",
              "diskSizeGB": "%d",
              "lun": %d,
              "name": "[concat(variables('%sVMNamePrefix'), copyIndex(),'-datadisk%d')]",
              "vhd": {
                "uri": "[concat('http://',variables('storageAccountPrefixes')[mod(add(add(div(copyIndex(),variables('maxVMsPerStorageAccount')),variables('%sStorageAccountOffset')),variables('dataStorageAccountPrefixSeed')),variables('storageAccountPrefixesCount'))],variables('storageAccountPrefixes')[div(add(add(div(copyIndex(),variables('maxVMsPerStorageAccount')),variables('%sStorageAccountOffset')),variables('dataStorageAccountPrefixSeed')),variables('storageAccountPrefixesCount'))],variables('%sDataAccountName'),'.blob.core.windows.net/vhds/',variables('%sVMNamePrefix'),copyIndex(), '--datadisk%d.vhd')]"
              }
            }`
	managedDataDisks := `            {
              "diskSizeGB": "%d",
              "lun": %d,
              "createOption": "Empty"
            }`
	for i, diskSize := range a.DiskSizesGB {
		if i > 0 {
			buf.WriteString(",\n")
		}
		if a.StorageProfile == api.StorageAccount {
			buf.WriteString(fmt.Sprintf(dataDisks, diskSize, i, a.Name, i, a.Name, a.Name, a.Name, a.Name, i))
		} else if a.StorageProfile == api.ManagedDisks {
			buf.WriteString(fmt.Sprintf(managedDataDisks, diskSize, i))
		}
	}
	buf.WriteString("\n          ],")
	return buf.String()
}

func getSecurityRules(ports []int) string {
	var buf bytes.Buffer
	for index, port := range ports {
		if index > 0 {
			buf.WriteString(",\n")
		}
		buf.WriteString(getSecurityRule(port, index))
	}
	return buf.String()
}

// getSingleLineForTemplate returns the file as a single line for embedding in an arm template
func (t *TemplateGenerator) getSingleLineForTemplate(textFilename string, cs *api.ContainerService, profile interface{}) (string, error) {
	b, err := Asset(textFilename)
	if err != nil {
		return "", t.Translator.Errorf("yaml file %s does not exist", textFilename)
	}

	// use go templates to process the text filename
	templ := template.New("customdata template").Funcs(t.getTemplateFuncMap(cs))
	if _, err = templ.New(textFilename).Parse(string(b)); err != nil {
		return "", t.Translator.Errorf("error parsing file %s: %v", textFilename, err)
	}

	var buffer bytes.Buffer
	if err = templ.ExecuteTemplate(&buffer, textFilename, profile); err != nil {
		return "", t.Translator.Errorf("error executing template for file %s: %v", textFilename, err)
	}
	expandedTemplate := buffer.String()

	textStr := escapeSingleLine(string(expandedTemplate))

	return textStr, nil
}

func escapeSingleLine(escapedStr string) string {
	// template.JSEscapeString leaves undesirable chars that don't work with pretty print
	escapedStr = strings.Replace(escapedStr, "\\", "\\\\", -1)
	escapedStr = strings.Replace(escapedStr, "\r\n", "\\n", -1)
	escapedStr = strings.Replace(escapedStr, "\n", "\\n", -1)
	escapedStr = strings.Replace(escapedStr, "\"", "\\\"", -1)
	return escapedStr
}

// getBase64CustomScript will return a base64 of the CSE
func getBase64CustomScript(csFilename string) string {
	b, err := Asset(csFilename)
	if err != nil {
		// this should never happen and this is a bug
		panic(fmt.Sprintf("BUG: %s", err.Error()))
	}
	// translate the parameters
	csStr := string(b)
	csStr = strings.Replace(csStr, "\r\n", "\n", -1)
	return getBase64CustomScriptFromStr(csStr)
}

// getBase64CustomScript will return a base64 of the CSE
func getBase64CustomScriptFromStr(str string) string {
	var gzipB bytes.Buffer
	w := gzip.NewWriter(&gzipB)
	w.Write([]byte(str))
	w.Close()
	return base64.StdEncoding.EncodeToString(gzipB.Bytes())
}

func getDCOSAgentProvisionScript(profile *api.AgentPoolProfile) string {
	// add the provision script

	var scriptname string
	if profile.OSType == api.Windows {
		scriptname = dcosWindowsProvision
	} else {
		scriptname = dcosProvision
	}

	bp, err1 := Asset(scriptname)
	if err1 != nil {
		panic(fmt.Sprintf("BUG: %s", err1.Error()))
	}

	provisionScript := string(bp)
	if strings.Contains(provisionScript, "'") {
		panic(fmt.Sprintf("BUG: %s may not contain character '", dcosProvision))
	}

	// the embedded roleFileContents
	var roleFileContents string
	if len(profile.Ports) > 0 {
		// public agents
		roleFileContents = "touch /etc/mesosphere/roles/slave_public"
	} else {
		roleFileContents = "touch /etc/mesosphere/roles/slave"
	}

	provisionScript = strings.Replace(provisionScript, "ROLESFILECONTENTS", roleFileContents, -1)

	return provisionScript
}

func getDCOSMasterProvisionScript() string {
	// add the provision script
	bp, err1 := Asset(dcosProvision)
	if err1 != nil {
		panic(fmt.Sprintf("BUG: %s", err1.Error()))
	}

	provisionScript := string(bp)
	if strings.Contains(provisionScript, "'") {
		panic(fmt.Sprintf("BUG: %s may not contain character '", dcosProvision))
	}

	// the embedded roleFileContents
	roleFileContents := `touch /etc/mesosphere/roles/master
touch /etc/mesosphere/roles/azure_master`
	provisionScript = strings.Replace(provisionScript, "ROLESFILECONTENTS", roleFileContents, -1)

	return provisionScript
}

// getSingleLineForTemplate returns the file as a single line for embedding in an arm template
func getSingleLineDCOSCustomData(orchestratorType string, orchestratorRelease string, masterCount int, provisionContent string, attributeContents string) string {
	yamlFilename := ""
	switch orchestratorType {
	case api.DCOS:
		switch orchestratorRelease {
		case api.DCOSRelease1Dot7:
			yamlFilename = dcosCustomData173
		case api.DCOSRelease1Dot8:
			yamlFilename = dcosCustomData188
		case api.DCOSRelease1Dot9:
			yamlFilename = dcosCustomData190
		}
	default:
		// it is a bug to get here
		panic(fmt.Sprintf("BUG: invalid orchestrator %s", orchestratorType))
	}

	b, err := Asset(yamlFilename)
	if err != nil {
		panic(fmt.Sprintf("BUG: %s", err.Error()))
	}

	// transform the provision script content
	provisionContent = strings.Replace(provisionContent, "\r\n", "\n", -1)
	provisionContent = strings.Replace(provisionContent, "\n", "\n\n    ", -1)

	yamlStr := string(b)
	yamlStr = strings.Replace(yamlStr, "PROVISION_STR", provisionContent, -1)
	yamlStr = strings.Replace(yamlStr, "ATTRIBUTES_STR", attributeContents, -1)

	// convert to json
	jsonBytes, err4 := yaml.YAMLToJSON([]byte(yamlStr))
	if err4 != nil {
		panic(fmt.Sprintf("BUG: %s", err4.Error()))
	}
	yamlStr = string(jsonBytes)

	// convert to one line
	yamlStr = strings.Replace(yamlStr, "\\", "\\\\", -1)
	yamlStr = strings.Replace(yamlStr, "\r\n", "\\n", -1)
	yamlStr = strings.Replace(yamlStr, "\n", "\\n", -1)
	yamlStr = strings.Replace(yamlStr, "\"", "\\\"", -1)

	// variable replacement
	rVariable, e1 := regexp.Compile("{{{([^}]*)}}}")
	if e1 != nil {
		panic(fmt.Sprintf("BUG: %s", e1.Error()))
	}
	yamlStr = rVariable.ReplaceAllString(yamlStr, "',variables('$1'),'")

	// replace the internal values
	guid := getPackageGUID(orchestratorType, orchestratorRelease, masterCount)
	yamlStr = strings.Replace(yamlStr, "DCOSGUID", guid, -1)
	publicIPStr := getDCOSCustomDataPublicIPStr(orchestratorType, masterCount)
	yamlStr = strings.Replace(yamlStr, "DCOSCUSTOMDATAPUBLICIPSTR", publicIPStr, -1)

	return yamlStr
}

func buildYamlFileWithWriteFiles(files []string) string {
	clusterYamlFile := `#cloud-config

write_files:
%s
`
	writeFileBlock := ` -  encoding: gzip
    content: !!binary |
        %s
    path: /opt/azure/containers/%s
    permissions: "0744"
`

	filelines := ""
	for _, file := range files {
		b64GzipString := getBase64CustomScript(file)
		filelines = filelines + fmt.Sprintf(writeFileBlock, b64GzipString, file)
	}
	return fmt.Sprintf(clusterYamlFile, filelines)
}

func getKubernetesSubnets(properties *api.Properties) string {
	subnetString := `{
            "name": "podCIDR%d",
            "properties": {
              "addressPrefix": "10.244.%d.0/24",
              "networkSecurityGroup": {
                "id": "[variables('nsgID')]"
              },
              "routeTable": {
                "id": "[variables('routeTableID')]"
              }
            }
          }`
	var buf bytes.Buffer

	cidrIndex := getKubernetesPodStartIndex(properties)
	for _, agentProfile := range properties.AgentPoolProfiles {
		if agentProfile.OSType == api.Windows {
			for i := 0; i < agentProfile.Count; i++ {
				buf.WriteString(",\n")
				buf.WriteString(fmt.Sprintf(subnetString, cidrIndex, cidrIndex))
				cidrIndex++
			}
		}
	}
	return buf.String()
}

func getKubernetesPodStartIndex(properties *api.Properties) int {
	nodeCount := 0
	nodeCount += properties.MasterProfile.Count
	for _, agentProfile := range properties.AgentPoolProfiles {
		if agentProfile.OSType != api.Windows {
			nodeCount += agentProfile.Count
		}
	}

	return nodeCount + 1
}

// getLinkedTemplatesForExtensions returns the
// Microsoft.Resources/deployments for each extension
//func getLinkedTemplatesForExtensions(properties api.Properties) string {
func getLinkedTemplatesForExtensions(properties *api.Properties) string {
	var result string

	var extensions = properties.ExtensionsProfile

	//This is temporary - to show you how to access the Extensions in the MasterProfile
	var masterProfileExtensions = properties.MasterProfile.Extensions
	var orchestratorType = properties.OrchestratorProfile.OrchestratorType

	for err, extensionProfile := range extensions {
		_ = err

		masterOptedForExtension, singleOrAll := validateProfileOptedForExtension(extensionProfile.Name, masterProfileExtensions)
		if masterOptedForExtension {
			result += ","
			dta, e := getMasterLinkedTemplateText(properties.MasterProfile, orchestratorType, extensionProfile, singleOrAll)
			if e != nil {
				fmt.Printf(e.Error())
				return ""
			}
			result += dta
		}

		for _, agentPoolProfile := range properties.AgentPoolProfiles {
			var poolProfileExtensions = agentPoolProfile.Extensions
			poolOptedForExtension, singleOrAll := validateProfileOptedForExtension(extensionProfile.Name, poolProfileExtensions)
			if poolOptedForExtension {
				result += ","
				dta, e := getAgentPoolLinkedTemplateText(agentPoolProfile, orchestratorType, extensionProfile, singleOrAll)
				if e != nil {
					fmt.Printf(e.Error())
					return ""
				}
				result += dta
			}

		}
	}

	return result
}

func getMasterLinkedTemplateText(masterProfile *api.MasterProfile, orchestratorType string, extensionProfile api.ExtensionProfile, singleOrAll string) (string, error) {
	extTargetVMNamePrefix := "variables('masterVMNamePrefix')"
	loopCount := "[sub(variables('masterCount'), variables('masterOffset'))]"
	if strings.EqualFold(singleOrAll, "single") {
		loopCount = "1"
	}
	return internalGetPoolLinkedTemplateText(extTargetVMNamePrefix, orchestratorType, loopCount,
		"variables('masterOffset')", extensionProfile)
}

func getAgentPoolLinkedTemplateText(agentPoolProfile *api.AgentPoolProfile, orchestratorType string, extensionProfile api.ExtensionProfile, singleOrAll string) (string, error) {
	extTargetVMNamePrefix := fmt.Sprintf("variables('%sVMNamePrefix')", agentPoolProfile.Name)
	loopCount := fmt.Sprintf("[variables('%sCount'))]", agentPoolProfile.Name)
	loopOffset := ""

	// Availability sets can have an offset since we don't redeploy vms.
	// So we don't want to rerun these extensions in scale up scenarios.
	if agentPoolProfile.IsAvailabilitySets() {
		loopCount = fmt.Sprintf("[sub(variables('%sCount'), variables('%sOffset'))]",
			agentPoolProfile.Name, agentPoolProfile.Name)
		loopOffset = fmt.Sprintf("variables('%sOffset')", agentPoolProfile.Name)

	}

	if strings.EqualFold(singleOrAll, "single") {
		loopCount = "1"
	}

	return internalGetPoolLinkedTemplateText(extTargetVMNamePrefix, orchestratorType, loopCount,
		loopOffset, extensionProfile)
}

func internalGetPoolLinkedTemplateText(extTargetVMNamePrefix, orchestratorType, loopCount, loopOffset string, extensionProfile api.ExtensionProfile) (string, error) {
	dta, e := getLinkedTemplateText(orchestratorType, extensionProfile.Name, extensionProfile.Version, extensionProfile.RootURL)
	if e != nil {
		return "", e
	}
	parmetersString := strings.Replace(extensionProfile.ExtensionParameters, "EXTENSION_LOOP_INDEX",
		"copyIndex(EXTENSION_LOOP_OFFSET)", -1)
	dta = strings.Replace(dta, "EXTENSION_PARAMETERS_REPLACE", parmetersString, -1)

	if strings.TrimSpace(extensionProfile.RootURL) == "" {
		dta = strings.Replace(dta, "EXTENSION_URL_REPLACE", DefaultExtensionsRootURL, -1)
	} else {
		dta = strings.Replace(dta, "EXTENSION_URL_REPLACE", extensionProfile.RootURL, -1)
	}
	dta = strings.Replace(dta, "EXTENSION_TARGET_VM_NAME_PREFIX", extTargetVMNamePrefix, -1)
	dta = strings.Replace(dta, "EXTENSION_LOOP_COUNT", loopCount, -1)
	dta = strings.Replace(dta, "EXTENSION_LOOP_OFFSET", loopOffset, -1)
	return dta, nil
}

func validateProfileOptedForExtension(extensionName string, profileExtensions []api.Extension) (bool, string) {
	for _, extension := range profileExtensions {
		if extensionName == extension.Name {
			return true, extension.SingleOrAll
		}
	}
	return false, ""
}

// getLinkedTemplateText returns the string data from
// template-link.json in the following directory:
// extensionsRootURL/extensions/extensionName/version
// It returns an error if the extension cannot be found
// or loaded.  getLinkedTemplateText calls getLinkedTemplateTextForURL,
// passing the default rootURL.
func getLinkedTemplateText(orchestratorType string, extensionName string, version string, rootURL string) (string, error) {
	if strings.TrimSpace(rootURL) == "" {
		return getLinkedTemplateTextForURL(DefaultExtensionsRootURL, orchestratorType, extensionName, version)
	}

	return getLinkedTemplateTextForURL(rootURL, orchestratorType, extensionName, version)
}

// getLinkedTemplateTextForURL returns the string data from
// template-link.json in the following directory:
// extensionsRootURL/extensions/extensionName/version
// It returns an error if the extension cannot be found
// or loaded.  getLinkedTemplateTextForURL provides the ability
// to pass a root extensions url for testing
func getLinkedTemplateTextForURL(rootURL string, orchestrator string, extensionName string, version string) (string, error) {
	supportsExtension, err := orchestratorSupportsExtension(rootURL, orchestrator, extensionName, version)
	if supportsExtension == false {
		return "", fmt.Errorf("Extension not supported for orchestrator. Error: %s", err)
	}

	templateLinkBytes, err := getExtensionResource(rootURL, extensionName, version, "template-link.json")
	if err != nil {
		return "", err
	}

	return string(templateLinkBytes), nil
}

func orchestratorSupportsExtension(rootURL string, orchestrator string, extensionName string, version string) (bool, error) {
	orchestratorBytes, err := getExtensionResource(rootURL, extensionName, version, "supported-orchestrators.json")
	if err != nil {
		return false, err
	}

	var supportedOrchestrators []string
	err = json.Unmarshal(orchestratorBytes, &supportedOrchestrators)
	if err != nil {
		return false, fmt.Errorf("Unable to parse supported-orchestrators.json for Extension %s Version %s", extensionName, version)
	}

	if stringInSlice(orchestrator, supportedOrchestrators) != true {
		return false, fmt.Errorf("Orchestrator: %s not in list of supported orchestrators for Extension: %s Version %s", orchestrator, extensionName, version)
	}

	return true, nil
}

func getExtensionResource(rootURL string, extensionName string, version string, fileName string) ([]byte, error) {
	requestURL := getExtensionURL(rootURL, extensionName, version, fileName)

	res, err := http.Get(requestURL)
	if err != nil {
		return nil, fmt.Errorf("Unable to GET extension resource for extension: %s with version %s with filename %s at URL: %s Error: %s", extensionName, version, fileName, requestURL, err)
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("Unable to GET extension resource for extension: %s with version %s with filename %s at URL: %s StatusCode: %s: Status: %s", extensionName, version, fileName, requestURL, strconv.Itoa(res.StatusCode), res.Status)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("Unable to GET extension resource for extension: %s with version %s  with filename %s at URL: %s Error: %s", extensionName, version, fileName, requestURL, err)
	}

	return body, nil
}

func getExtensionURL(rootURL string, extensionName string, version string, fileName string) string {
	extensionsDir := "extensions"

	return rootURL + extensionsDir + "/" + extensionName + "/" + version + "/" + fileName
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func GetSwarmVersions(orchestratorVersion, dockerComposeVersion string) string {
	return fmt.Sprintf("\"orchestratorVersion\": \"%s\",\n\"dockerComposeVersion\": \"%s\",\n", orchestratorVersion, dockerComposeVersion)
}

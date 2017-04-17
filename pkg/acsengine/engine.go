package acsengine

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"hash/fnv"
	"math/rand"
	"regexp"
	"strings"
	"text/template"

	"github.com/Azure/acs-engine/pkg/api"
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
	dcosCustomData173 = "dcoscustomdata173.t"
	dcosCustomData184 = "dcoscustomdata184.t"
	dcosCustomData187 = "dcoscustomdata187.t"
	dcosCustomData188 = "dcoscustomdata188.t"
	dcosCustomData190 = "dcoscustomdata190.t"
	dcosProvision     = "dcosprovision.sh"
)

const (
	swarmProvision        = "configure-swarm-cluster.sh"
	swarmWindowsProvision = "Install-ContainerHost-And-Join-Swarm.ps1"

	swarmModeProvision        = "configure-swarmmode-cluster.sh"
	swarmModeWindowsProvision = "Join-SwarmMode-cluster.ps1"
)

const (
	agentOutputs                 = "agentoutputs.t"
	agentParams                  = "agentparams.t"
	classicParams                = "classicparams.t"
	dcosAgentResourcesVMAS       = "dcosagentresourcesvmas.t"
	dcosAgentResourcesVMSS       = "dcosagentresourcesvmss.t"
	dcosAgentVars                = "dcosagentvars.t"
	dcosBaseFile                 = "dcosbase.t"
	dcosParams                   = "dcosparams.t"
	dcosMasterResources          = "dcosmasterresources.t"
	dcosMasterVars               = "dcosmastervars.t"
	kubernetesBaseFile           = "kubernetesbase.t"
	kubernetesAgentResourcesVMAS = "kubernetesagentresourcesvmas.t"
	kubernetesAgentVars          = "kubernetesagentvars.t"
	kubernetesMasterResources    = "kubernetesmasterresources.t"
	kubernetesMasterVars         = "kubernetesmastervars.t"
	kubernetesParams             = "kubernetesparams.t"
	kubernetesWinAgentVars       = "kuberneteswinagentresourcesvmas.t"
	kubernetesKubeletService     = "kuberneteskubelet.service"
	masterOutputs                = "masteroutputs.t"
	masterParams                 = "masterparams.t"
	swarmBaseFile                = "swarmbase.t"
	swarmAgentResourcesVMAS      = "swarmagentresourcesvmas.t"
	swarmAgentResourcesVMSS      = "swarmagentresourcesvmss.t"
	swarmAgentResourcesClassic   = "swarmagentresourcesclassic.t"
	swarmAgentVars               = "swarmagentvars.t"
	swarmMasterResources         = "swarmmasterresources.t"
	swarmMasterVars              = "swarmmastervars.t"
	swarmWinAgentResourcesVMAS   = "swarmwinagentresourcesvmas.t"
	swarmWinAgentResourcesVMSS   = "swarmwinagentresourcesvmss.t"
	windowsParams                = "windowsparams.t"
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

var kubernetesAddonYamls = map[string]string{
	"MASTER_ADDON_HEAPSTER_DEPLOYMENT_B64_GZIP_STR":             "kubernetesmasteraddons-heapster-deployment.yaml",
	"MASTER_ADDON_HEAPSTER_SERVICE_B64_GZIP_STR":                "kubernetesmasteraddons-heapster-service.yaml",
	"MASTER_ADDON_KUBE_DNS_DEPLOYMENT_B64_GZIP_STR":             "kubernetesmasteraddons-kube-dns-deployment.yaml",
	"MASTER_ADDON_KUBE_DNS_SERVICE_B64_GZIP_STR":                "kubernetesmasteraddons-kube-dns-service.yaml",
	"MASTER_ADDON_KUBE_PROXY_DAEMONSET_B64_GZIP_STR":            "kubernetesmasteraddons-kube-proxy-daemonset.yaml",
	"MASTER_ADDON_KUBERNETES_DASHBOARD_DEPLOYMENT_B64_GZIP_STR": "kubernetesmasteraddons-kubernetes-dashboard-deployment.yaml",
	"MASTER_ADDON_KUBERNETES_DASHBOARD_SERVICE_B64_GZIP_STR":    "kubernetesmasteraddons-kubernetes-dashboard-service.yaml",
	"MASTER_ADDON_DEFAULT_STORAGE_CLASS_B64_GZIP_STR":           "kubernetesmasteraddons-default-storage-class.yaml",
}

var commonTemplateFiles = []string{agentOutputs, agentParams, classicParams, masterOutputs, masterParams, windowsParams}
var dcosTemplateFiles = []string{dcosAgentResourcesVMAS, dcosAgentResourcesVMSS, dcosAgentVars, dcosBaseFile, dcosMasterResources, dcosMasterVars, dcosParams}
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

type KeyVaultID struct {
	ID string `json:"id"`
}

type KeyVaultRef struct {
	KeyVault      KeyVaultID `json:"keyVault"`
	SecretName    string     `json:"secretName"`
	SecretVersion string     `json:"secretVersion,omitempty"`
}

var keyvaultSecretPath_re *regexp.Regexp

func init() {
	keyvaultSecretPath_re = regexp.MustCompile(`^(/subscriptions/\S+/resourceGroups/\S+/providers/Microsoft.KeyVault/vaults/\S+)/secrets/([^/\s]+)(/(\S+))?$`)
}

func (t *TemplateGenerator) verifyFiles() error {
	allFiles := append(commonTemplateFiles, dcosTemplateFiles...)
	allFiles = append(allFiles, kubernetesTemplateFiles...)
	allFiles = append(allFiles, swarmTemplateFiles...)
	for _, file := range allFiles {
		if _, err := Asset(file); err != nil {
			return fmt.Errorf("template file %s does not exist", file)
		}
	}
	return nil
}

// TemplateGenerator represents the object that performs the template generation.
type TemplateGenerator struct {
	ClassicMode bool
}

// InitializeTemplateGenerator creates a new template generator object
func InitializeTemplateGenerator(classicMode bool) (*TemplateGenerator, error) {
	t := &TemplateGenerator{
		ClassicMode: classicMode,
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
	certsGenerated = false
	err = nil

	var templ *template.Template

	properties := containerService.Properties

	if certsGenerated, err = SetPropertiesDefaults(containerService); err != nil {
		return templateRaw, parametersRaw, certsGenerated, err
	}

	templ = template.New("acs template").Funcs(t.getTemplateFuncMap(containerService))

	files, baseFile, e := prepareTemplateFiles(properties)
	if e != nil {
		return "", "", false, e
	}

	for _, file := range files {
		bytes, e := Asset(file)
		if e != nil {
			err = fmt.Errorf("Error reading file %s, Error: %s", file, e.Error())
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
			err = fmt.Errorf("%v", r)
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

	var parametersMap map[string]interface{}
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
	h.Write([]byte(properties.MasterProfile.DNSPrefix))
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
	kubeconfig = strings.Replace(kubeconfig, "<<<variables('caCertificate')>>>", base64.StdEncoding.EncodeToString([]byte(properties.CertificateProfile.CaCertificate)), -1)
	kubeconfig = strings.Replace(kubeconfig, "<<<reference(concat('Microsoft.Network/publicIPAddresses/', variables('masterPublicIPAddressName'))).dnsSettings.fqdn>>>", FormatAzureProdFQDN(properties.MasterProfile.DNSPrefix, location), -1)
	kubeconfig = strings.Replace(kubeconfig, "{{{resourceGroup}}}", properties.MasterProfile.DNSPrefix, -1)
	kubeconfig = strings.Replace(kubeconfig, "<<<variables('kubeConfigCertificate')>>>", base64.StdEncoding.EncodeToString([]byte(properties.CertificateProfile.KubeConfigCertificate)), -1)
	kubeconfig = strings.Replace(kubeconfig, "<<<variables('kubeConfigPrivateKey')>>>", base64.StdEncoding.EncodeToString([]byte(properties.CertificateProfile.KubeConfigPrivateKey)), -1)

	return kubeconfig, nil
}

func prepareTemplateFiles(properties *api.Properties) ([]string, string, error) {
	var files []string
	var baseFile string
	if properties.OrchestratorProfile.OrchestratorType == api.DCOS188 ||
		properties.OrchestratorProfile.OrchestratorType == api.DCOS190 ||
		properties.OrchestratorProfile.OrchestratorType == api.DCOS187 ||
		properties.OrchestratorProfile.OrchestratorType == api.DCOS184 ||
		properties.OrchestratorProfile.OrchestratorType == api.DCOS173 {
		files = append(commonTemplateFiles, dcosTemplateFiles...)
		baseFile = dcosBaseFile
	} else if properties.OrchestratorProfile.OrchestratorType == api.Swarm {
		files = append(commonTemplateFiles, swarmTemplateFiles...)
		baseFile = swarmBaseFile
	} else if properties.OrchestratorProfile.OrchestratorType == api.Kubernetes {
		files = append(commonTemplateFiles, kubernetesTemplateFiles...)
		baseFile = kubernetesBaseFile
	} else if properties.OrchestratorProfile.OrchestratorType == api.SwarmMode {
		files = append(commonTemplateFiles, swarmModeTemplateFiles...)
		baseFile = swarmBaseFile
	} else {
		return nil, "", fmt.Errorf("orchestrator '%s' is unsupported", properties.OrchestratorProfile.OrchestratorType)
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

func getParameters(cs *api.ContainerService, isClassicMode bool) (map[string]interface{}, error) {
	properties := cs.Properties
	location := cs.Location
	parametersMap := map[string]interface{}{}

	// Master Parameters
	addValue(parametersMap, "location", location)
	addValue(parametersMap, "targetEnvironment", GetCloudTargetEnv(location))
	addValue(parametersMap, "linuxAdminUsername", properties.LinuxProfile.AdminUsername)
	addValue(parametersMap, "masterEndpointDNSNamePrefix", properties.MasterProfile.DNSPrefix)
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
		addSecret(parametersMap, "apiServerCertificate", properties.CertificateProfile.APIServerCertificate, true)
		addSecret(parametersMap, "apiServerPrivateKey", properties.CertificateProfile.APIServerPrivateKey, true)
		addSecret(parametersMap, "caCertificate", properties.CertificateProfile.CaCertificate, true)
		addSecret(parametersMap, "clientCertificate", properties.CertificateProfile.ClientCertificate, true)
		addSecret(parametersMap, "clientPrivateKey", properties.CertificateProfile.ClientPrivateKey, true)
		addSecret(parametersMap, "kubeConfigCertificate", properties.CertificateProfile.KubeConfigCertificate, true)
		addSecret(parametersMap, "kubeConfigPrivateKey", properties.CertificateProfile.KubeConfigPrivateKey, true)
		addValue(parametersMap, "dockerEngineDownloadRepo", cloudSpecConfig.DockerSpecConfig.DockerEngineRepo)
		addValue(parametersMap, "kubeDnsServiceIP", properties.OrchestratorProfile.KubernetesConfig.DnsServiceIP)
		addValue(parametersMap, "kubeServiceCidr", properties.OrchestratorProfile.KubernetesConfig.ServiceCIDR)
		addValue(parametersMap, "kubeClusterCidr", properties.OrchestratorProfile.KubernetesConfig.ClusterCIDR)
		addValue(parametersMap, "kubernetesHyperkubeSpec", properties.OrchestratorProfile.KubernetesConfig.KubernetesImageBase+KubernetesHyperkubeImageName)
		addValue(parametersMap, "kubernetesAddonManagerSpec", cloudSpecConfig.KubernetesSpecConfig.KubernetesImageBase+KubernetesAddonManagerImageName)
		addValue(parametersMap, "kubernetesAddonResizerSpec", cloudSpecConfig.KubernetesSpecConfig.KubernetesImageBase+KubernetesAddonResizerImageName)
		addValue(parametersMap, "kubernetesDashboardSpec", cloudSpecConfig.KubernetesSpecConfig.KubernetesImageBase+KubernetesDashboardImageName)
		addValue(parametersMap, "kubernetesDNSMasqSpec", cloudSpecConfig.KubernetesSpecConfig.KubernetesImageBase+KubernetesDNSMasqImageName)
		addValue(parametersMap, "kubernetesExecHealthzSpec", cloudSpecConfig.KubernetesSpecConfig.KubernetesImageBase+KubernetesExechealthzImageName)
		addValue(parametersMap, "kubernetesHeapsterSpec", cloudSpecConfig.KubernetesSpecConfig.KubernetesImageBase+KubernetesHeapsterImageName)
		addValue(parametersMap, "kubernetesKubeDNSSpec", cloudSpecConfig.KubernetesSpecConfig.KubernetesImageBase+KubernetesDNSImageName)
		addValue(parametersMap, "kubernetesPodInfraContainerSpec", cloudSpecConfig.KubernetesSpecConfig.KubernetesImageBase+KubernetesPodInfraContainerImageName)
		addValue(parametersMap, "networkPolicy", properties.OrchestratorProfile.KubernetesConfig.NetworkPolicy)
		addValue(parametersMap, "servicePrincipalClientId", properties.ServicePrincipalProfile.ClientID)
		addSecret(parametersMap, "servicePrincipalClientSecret", properties.ServicePrincipalProfile.Secret, false)
	}

	if strings.HasPrefix(string(properties.OrchestratorProfile.OrchestratorType), string(api.DCOS)) {
		dcosBootstrapURL := cloudSpecConfig.DCOSSpecConfig.DCOS188_BootstrapDownloadURL
		switch properties.OrchestratorProfile.OrchestratorType {
		case api.DCOS:
			dcosBootstrapURL = cloudSpecConfig.DCOSSpecConfig.DCOS190_BootstrapDownloadURL
		case api.DCOS173:
			dcosBootstrapURL = cloudSpecConfig.DCOSSpecConfig.DCOS173_BootstrapDownloadURL
		case api.DCOS184:
			dcosBootstrapURL = cloudSpecConfig.DCOSSpecConfig.DCOS184_BootstrapDownloadURL
		case api.DCOS187:
			dcosBootstrapURL = cloudSpecConfig.DCOSSpecConfig.DCOS187_BootstrapDownloadURL
		case api.DCOS188:
			dcosBootstrapURL = cloudSpecConfig.DCOSSpecConfig.DCOS188_BootstrapDownloadURL
		case api.DCOS190:
			dcosBootstrapURL = cloudSpecConfig.DCOSSpecConfig.DCOS190_BootstrapDownloadURL
		}
		addValue(parametersMap, "dcosBootstrapURL", dcosBootstrapURL)
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
		addValue(parametersMap, "windowsAdminPassword", properties.WindowsProfile.AdminPassword)
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

func addValue(m map[string]interface{}, k string, v interface{}) {
	m[k] = map[string]interface{}{
		"value": v,
	}
}

func addSecret(m map[string]interface{}, k string, v interface{}, encode bool) {
	str, ok := v.(string)
	if !ok {
		addValue(m, k, v)
		return
	}
	parts := keyvaultSecretPath_re.FindStringSubmatch(str)
	if parts == nil || len(parts) != 5 {
		if encode {
			addValue(m, k, base64.StdEncoding.EncodeToString([]byte(str)))
		} else {
			addValue(m, k, str)
		}
		return
	}

	m[k] = map[string]interface{}{
		"reference": &KeyVaultRef{
			KeyVault: KeyVaultID{
				ID: parts[1],
			},
			SecretName:    parts[2],
			SecretVersion: parts[4],
		},
	}
}

// getTemplateFuncMap returns all functions used in template generation
func (t *TemplateGenerator) getTemplateFuncMap(cs *api.ContainerService) map[string]interface{} {
	return template.FuncMap{
		"IsDCOS173": func() bool {
			return cs.Properties.OrchestratorProfile.OrchestratorType == api.DCOS173
		},
		"IsDCOS184": func() bool {
			return cs.Properties.OrchestratorProfile.OrchestratorType == api.DCOS184
		},
		"IsDCOS187": func() bool {
			return cs.Properties.OrchestratorProfile.OrchestratorType == api.DCOS187
		},
		"IsDCOS188": func() bool {
			return cs.Properties.OrchestratorProfile.OrchestratorType == api.DCOS188
		},
		"IsDCOS190": func() bool {
			return cs.Properties.OrchestratorProfile.OrchestratorType == api.DCOS190
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
			str := getSingleLineDCOSCustomData(cs.Properties.OrchestratorProfile.OrchestratorType, cs.Properties.MasterProfile.Count, masterProvisionScript, masterAttributeContents)

			return fmt.Sprintf("\"customData\": \"[base64(concat('#cloud-config\\n\\n', '%s'))]\",", str)
		},
		"GetDCOSAgentCustomData": func(profile *api.AgentPoolProfile) string {
			agentProvisionScript := getDCOSAgentProvisionScript(profile)
			attributeContents := getDCOSAgentCustomNodeLabels(profile)
			str := getSingleLineDCOSCustomData(cs.Properties.OrchestratorProfile.OrchestratorType, cs.Properties.MasterProfile.Count, agentProvisionScript, attributeContents)

			return fmt.Sprintf("\"customData\": \"[base64(concat('#cloud-config\\n\\n', '%s'))]\",", str)
		},
		"GetMasterAllowedSizes": func() string {
			if t.ClassicMode {
				return GetClassicAllowedSizes()
			} else if cs.Properties.OrchestratorProfile.OrchestratorType == api.DCOS190 ||
				cs.Properties.OrchestratorProfile.OrchestratorType == api.DCOS188 ||
				cs.Properties.OrchestratorProfile.OrchestratorType == api.DCOS187 ||
				cs.Properties.OrchestratorProfile.OrchestratorType == api.DCOS184 ||
				cs.Properties.OrchestratorProfile.OrchestratorType == api.DCOS173 {
				return GetDCOSMasterAllowedSizes()
			}
			return GetMasterAgentAllowedSizes()
		},
		"GetAgentAllowedSizes": func() string {
			if t.ClassicMode {
				return GetClassicAllowedSizes()
			}
			return GetMasterAgentAllowedSizes()
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
		"GetKubernetesMasterCustomData": func() string {
			str, e := getSingleLineForTemplate(kubernetesMasterCustomDataYaml)
			if e != nil {
				return ""
			}

			for placeholder, filename := range kubernetesManifestYamls {
				manifestTextContents := getBase64CustomScript(filename)
				str = strings.Replace(str, placeholder, manifestTextContents, -1)
			}

			// add artifacts and addons
			for placeholder, filename := range kubernetesAritfacts {
				addonTextContents := getBase64CustomScript(filename)
				str = strings.Replace(str, placeholder, addonTextContents, -1)
			}

			for placeholder, filename := range kubernetesAddonYamls {
				addonTextContents := getBase64CustomScript(filename)
				str = strings.Replace(str, placeholder, addonTextContents, -1)
			}

			// return the custom data
			return fmt.Sprintf("\"customData\": \"[base64(concat('%s'))]\",", str)
		},
		"GetKubernetesAgentCustomData": func() string {
			str, e := getSingleLineForTemplate(kubernetesAgentCustomDataYaml)
			if e != nil {
				return ""
			}

			// add artifacts
			for placeholder, filename := range kubernetesAritfacts {
				addonTextContents := getBase64CustomScript(filename)
				str = strings.Replace(str, placeholder, addonTextContents, -1)
			}

			return fmt.Sprintf("\"customData\": \"[base64(concat('%s'))]\",", str)
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
		"GetKubernetesWindowsAgentCustomData": func() string {
			str, e := getSingleLineForTemplate(kubernetesWindowsAgentCustomDataPS1)
			if e != nil {
				return ""
			}
			return fmt.Sprintf("\"customData\": \"[base64(concat('%s'))]\",", str)
		},
		"GetKubernetesKubeConfig": func() string {
			str, e := getSingleLineForTemplate(kubeConfigJSON)
			if e != nil {
				return ""
			}
			return str
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
		"AnyAgentHasDisks": func() bool {
			for _, agentProfile := range cs.Properties.AgentPoolProfiles {
				if agentProfile.HasDisks() {
					return true
				}
			}
			return false
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

func getPackageGUID(orchestratorType api.OrchestratorType, masterCount int) string {
	if orchestratorType == api.DCOS190 {
		switch masterCount {
		case 1:
			return "bcc883b7a3191412cf41824bdee06c1142187a0b"
		case 3:
			return "dcff7e24c0c1827bebeb7f1a806f558054481b33"
		case 5:
			return "b41bfa84137a6374b2ff5eb1655364d7302bd257"
		}
	} else if orchestratorType == api.DCOS188 {
		switch masterCount {
		case 1:
			return "441385ce2f5942df7e29075c12fb38fa5e92cbba"
		case 3:
			return "b1cd359287504efb780257bd12cc3a63704e42d4"
		case 5:
			return "d9b61156dfcc9383e014851529738aa550ef57d9"
		}
	} else if orchestratorType == api.DCOS187 {
		switch masterCount {
		case 1:
			return "556978041b6ed059cc0f474501083e35ea5645b8"
		case 3:
			return "1eb387eda0403c7fd6f1dacf66e530be74c3c3de"
		case 5:
			return "2e38627207dc70f46296b9649f9ee2a43500ec15"
		}
	} else if orchestratorType == api.DCOS184 {
		switch masterCount {
		case 1:
			return "5ac6a7d060584c58c704e1f625627a591ecbde4e"
		case 3:
			return "42bd1d74e9a2b23836bd78919c716c20b98d5a0e"
		case 5:
			return "97947a91e2c024ed4f043bfcdad49da9418d3095"
		}
	} else if orchestratorType == api.DCOS173 {
		switch masterCount {
		case 1:
			return "6b604c1331c2b8b52bb23d1ea8a8d17e0f2b7428"
		case 3:
			return "6af5097e7956962a3d4318d28fbf280a47305485"
		case 5:
			return "376e07e0dbad2af3da2c03bc92bb07e84b3dafd5"
		}
	}
	return ""
}

func getDCOSCustomDataPublicIPStr(orchestratorType api.OrchestratorType, masterCount int) string {
	if orchestratorType == api.DCOS173 ||
		orchestratorType == api.DCOS184 ||
		orchestratorType == api.DCOS187 ||
		orchestratorType == api.DCOS188 ||
		orchestratorType == api.DCOS190 {
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
	for i := range properties.AgentPoolProfiles {
		profile := &properties.AgentPoolProfiles[i]
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
func getSingleLineForTemplate(textFilename string) (string, error) {
	b, err := Asset(textFilename)
	if err != nil {
		return "", fmt.Errorf("yaml file %s does not exist", textFilename)
	}

	textStr := escapeSingleLine(string(b))

	// variable replacement
	rVariable, e1 := regexp.Compile("{{{([^}]*)}}}")
	if e1 != nil {
		return "", e1
	}
	textStr = rVariable.ReplaceAllString(textStr, "',variables('$1'),'")
	// verbatim replacement
	rVerbatim, e2 := regexp.Compile("<<<([^>]*)>>>")
	if e2 != nil {
		return "", e2
	}
	textStr = rVerbatim.ReplaceAllString(textStr, "',$1,'")
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
	bp, err1 := Asset(dcosProvision)
	if err1 != nil {
		panic(fmt.Sprintf("BUG: %s", err1.Error()))
	}

	provisionScript := string(bp)
	if strings.Contains(provisionScript, "'") {
		panic(fmt.Sprintf("BUG: %s may not contain character '", dcosProvision))
	}

	// the embedded roleFileContents
	roleFileContents := ""
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
func getSingleLineDCOSCustomData(orchestratorType api.OrchestratorType, masterCount int, provisionContent string, attributeContents string) string {
	yamlFilename := ""
	switch orchestratorType {
	case api.DCOS190:
		yamlFilename = dcosCustomData190
	case api.DCOS188:
		yamlFilename = dcosCustomData188
	case api.DCOS187:
		yamlFilename = dcosCustomData187
	case api.DCOS184:
		yamlFilename = dcosCustomData184
	case api.DCOS173:
		yamlFilename = dcosCustomData173
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
	guid := getPackageGUID(orchestratorType, masterCount)
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

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
)

const (
	kubernetesMasterCustomDataYaml = "kubernetesmastercustomdata.yml"
	kubernetesMasterCustomScript   = "kubernetesmastercustomscript.sh"
	kubernetesAgentCustomDataYaml  = "kubernetesagentcustomdata.yml"
	kubernetesAgentCustomScript    = "kubernetesagentcustomscript.sh"
	kubeConfigJSON                 = "kubeconfig.json"
)

const (
	agentOutputs                 = "agentoutputs.t"
	agentParams                  = "agentparams.t"
	classicParams                = "classicparams.t"
	dcosAgentResourcesVMAS       = "dcosagentresourcesvmas.t"
	dcosAgentResourcesVMSS       = "dcosagentresourcesvmss.t"
	dcosAgentVars                = "dcosagentvars.t"
	dcosBaseFile                 = "dcosbase.t"
	dcosCustomData173            = "dcoscustomdata173.t"
	dcosCustomData184            = "dcoscustomdata184.t"
	dcosMasterResources          = "dcosmasterresources.t"
	dcosMasterVars               = "dcosmastervars.t"
	kubernetesBaseFile           = "kubernetesbase.t"
	kubernetesAgentResourcesVMAS = "kubernetesagentresourcesvmas.t"
	kubernetesAgentVars          = "kubernetesagentvars.t"
	kubernetesMasterResources    = "kubernetesmasterresources.t"
	kubernetesMasterVars         = "kubernetesmastervars.t"
	kubernetesParams             = "kubernetesparams.t"
	kubernetesWinAgentVars       = "kuberneteswinagentresourcesvmas.t"
	masterOutputs                = "masteroutputs.t"
	masterParams                 = "masterparams.t"
	swarmBaseFile                = "swarmbase.t"
	swarmAgentCustomData         = "swarmagentcustomdata.t"
	swarmAgentResourcesVMAS      = "swarmagentresourcesvmas.t"
	swarmAgentResourcesVMSS      = "swarmagentresourcesvmss.t"
	swarmAgentVars               = "swarmagentvars.t"
	swarmMasterCustomData        = "swarmmastercustomdata.t"
	swarmMasterResources         = "swarmmasterresources.t"
	swarmMasterVars              = "swarmmastervars.t"
	swarmWinAgentResourcesVMAS   = "swarmwinagentresourcesvmas.t"
	swarmWinAgentResourcesVMSS   = "swarmwinagentresourcesvmss.t"
	windowsParams                = "windowsparams.t"
)

var kubernetesAddonYamls = map[string]string{
	"MASTER_ADDON_HEAPSTER_DEPLOYMENT_B64_GZIP_STR":             "kubernetesmasteraddons-heapster-deployment.yaml",
	"MASTER_ADDON_HEAPSTER_SERVICE_B64_GZIP_STR":                "kubernetesmasteraddons-heapster-service.yaml",
	"MASTER_ADDON_KUBE_DNS_DEPLOYMENT_B64_GZIP_STR":             "kubernetesmasteraddons-kube-dns-deployment.yaml",
	"MASTER_ADDON_KUBE_DNS_SERVICE_B64_GZIP_STR":                "kubernetesmasteraddons-kube-dns-service.yaml",
	"MASTER_ADDON_KUBE_PROXY_DAEMONSET_B64_GZIP_STR":            "kubernetesmasteraddons-kube-proxy-daemonset.yaml",
	"MASTER_ADDON_KUBERNETES_DASHBOARD_DEPLOYMENT_B64_GZIP_STR": "kubernetesmasteraddons-kubernetes-dashboard-deployment.yaml",
	"MASTER_ADDON_KUBERNETES_DASHBOARD_SERVICE_B64_GZIP_STR":    "kubernetesmasteraddons-kubernetes-dashboard-service.yaml",
}

var commonTemplateFiles = []string{agentOutputs, agentParams, classicParams, masterOutputs, masterParams, windowsParams}
var dcosTemplateFiles = []string{dcosAgentResourcesVMAS, dcosAgentResourcesVMSS, dcosAgentVars, dcosBaseFile, dcosCustomData173, dcosCustomData184, dcosMasterResources, dcosMasterVars}
var kubernetesTemplateFiles = []string{kubernetesBaseFile, kubernetesAgentResourcesVMAS, kubernetesAgentVars, kubernetesMasterResources, kubernetesMasterVars, kubernetesParams, kubernetesWinAgentVars}
var swarmTemplateFiles = []string{swarmBaseFile, swarmAgentCustomData, swarmAgentResourcesVMAS, swarmAgentVars, swarmAgentResourcesVMSS, swarmBaseFile, swarmMasterCustomData, swarmMasterResources, swarmMasterVars, swarmWinAgentResourcesVMAS, swarmWinAgentResourcesVMSS}

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
func (t *TemplateGenerator) GenerateTemplate(containerService *api.ContainerService) (string, string, bool, error) {
	var err error
	var templ *template.Template
	certsGenerated := false

	properties := &containerService.Properties

	if certsGenerated, err = SetPropertiesDefaults(properties); err != nil {
		return "", "", certsGenerated, err
	}

	templ = template.New("acs template").Funcs(t.getTemplateFuncMap(properties))

	files, baseFile, e := prepareTemplateFiles(properties)
	if e != nil {
		return "", "", false, e
	}

	for _, file := range files {
		bytes, e := Asset(file)
		if e != nil {
			return "", "", certsGenerated, fmt.Errorf("Error reading file %s, Error: %s", file, e.Error())
		}
		if _, err = templ.New(file).Parse(string(bytes)); err != nil {
			return "", "", certsGenerated, err
		}
	}
	var b bytes.Buffer
	if err = templ.ExecuteTemplate(&b, baseFile, properties); err != nil {
		return "", "", certsGenerated, err
	}

	var parametersMap map[string]interface{}
	if parametersMap, err = getParameters(properties); err != nil {
		return "", "", certsGenerated, err
	}
	var parameterBytes []byte
	if parameterBytes, err = json.Marshal(parametersMap); err != nil {
		return "", "", certsGenerated, err
	}

	return b.String(), string(parameterBytes), certsGenerated, nil
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
	if properties.OrchestratorProfile.OrchestratorType == api.DCOS184 ||
		properties.OrchestratorProfile.OrchestratorType == api.DCOS ||
		properties.OrchestratorProfile.OrchestratorType == api.DCOS173 {
		files = append(commonTemplateFiles, dcosTemplateFiles...)
		baseFile = dcosBaseFile
	} else if properties.OrchestratorProfile.OrchestratorType == api.Swarm {
		files = append(commonTemplateFiles, swarmTemplateFiles...)
		baseFile = swarmBaseFile
	} else if properties.OrchestratorProfile.OrchestratorType == api.Kubernetes {
		files = append(commonTemplateFiles, kubernetesTemplateFiles...)
		baseFile = kubernetesBaseFile
	} else {
		return nil, "", fmt.Errorf("orchestrator '%s' is unsupported", properties.OrchestratorProfile.OrchestratorType)
	}

	return files, baseFile, nil
}

func getParameters(properties *api.Properties) (map[string]interface{}, error) {
	parametersMap := map[string]interface{}{}

	// Master Parameters
	addValue(parametersMap, "linuxAdminUsername", properties.LinuxProfile.AdminUsername)
	addValue(parametersMap, "masterEndpointDNSNamePrefix", properties.MasterProfile.DNSPrefix)
	if properties.MasterProfile.IsCustomVNET() {
		addValue(parametersMap, "masterVnetSubnetID", properties.MasterProfile.VnetSubnetID)
	} else {
		addValue(parametersMap, "masterSubnet", properties.MasterProfile.Subnet)
	}
	addValue(parametersMap, "firstConsecutiveStaticIP", properties.MasterProfile.FirstConsecutiveStaticIP)
	addValue(parametersMap, "masterVMSize", properties.MasterProfile.VMSize)
	addValue(parametersMap, "sshRSAPublicKey", properties.LinuxProfile.SSH.PublicKeys[0].KeyData)

	// Kubernetes Parameters
	if properties.OrchestratorProfile.OrchestratorType == api.Kubernetes {
		addValue(parametersMap, "apiServerCertificate", base64.StdEncoding.EncodeToString([]byte(properties.CertificateProfile.APIServerCertificate)))
		addValue(parametersMap, "apiServerPrivateKey", base64.StdEncoding.EncodeToString([]byte(properties.CertificateProfile.APIServerPrivateKey)))
		addValue(parametersMap, "caCertificate", base64.StdEncoding.EncodeToString([]byte(properties.CertificateProfile.CaCertificate)))
		addValue(parametersMap, "clientCertificate", base64.StdEncoding.EncodeToString([]byte(properties.CertificateProfile.ClientCertificate)))
		addValue(parametersMap, "clientPrivateKey", base64.StdEncoding.EncodeToString([]byte(properties.CertificateProfile.ClientPrivateKey)))
		addValue(parametersMap, "kubeConfigCertificate", base64.StdEncoding.EncodeToString([]byte(properties.CertificateProfile.KubeConfigCertificate)))
		addValue(parametersMap, "kubeConfigPrivateKey", base64.StdEncoding.EncodeToString([]byte(properties.CertificateProfile.KubeConfigPrivateKey)))
		addValue(parametersMap, "kubernetesHyperkubeSpec", properties.OrchestratorProfile.KubernetesConfig.KubernetesHyperkubeSpec)
		addValue(parametersMap, "kubectlVersion", properties.OrchestratorProfile.KubernetesConfig.KubectlVersion)
		addValue(parametersMap, "servicePrincipalClientId", properties.ServicePrincipalProfile.ClientID)
		addValue(parametersMap, "servicePrincipalClientSecret", properties.ServicePrincipalProfile.Secret)
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
	}

	return parametersMap, nil
}

func addValue(m map[string]interface{}, k string, v interface{}) {
	m[k] = map[string]interface{}{
		"value": v,
	}
}

// getTemplateFuncMap returns all functions used in template generation
func (t *TemplateGenerator) getTemplateFuncMap(properties *api.Properties) map[string]interface{} {
	return template.FuncMap{
		"IsDCOS173": func() bool {
			return properties.OrchestratorProfile.OrchestratorType == api.DCOS173
		},
		"IsDCOS184": func() bool {
			return properties.OrchestratorProfile.OrchestratorType == api.DCOS184 ||
				properties.OrchestratorProfile.OrchestratorType == api.DCOS
		},
		"RequiresFakeAgentOutput": func() bool {
			return properties.OrchestratorProfile.OrchestratorType == api.Kubernetes
		},
		"IsPublic": func(ports []int) bool {
			return len(ports) > 0
		},
		"GetVNETSubnetDependencies": func() string {
			return getVNETSubnetDependencies(properties)
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
		"GetMasterRolesFileContents": func() string {
			return getMasterRolesFileContents()
		},
		"GetAgentRolesFileContents": func(ports []int) string {
			return getAgentRolesFileContents(ports)
		},
		"GetDCOSCustomDataPublicIPStr": func() string {
			return getDCOSCustomDataPublicIPStr(properties.OrchestratorProfile.OrchestratorType, properties.MasterProfile.Count)
		},
		"GetDCOSGUID": func() string {
			return getPackageGUID(properties.OrchestratorProfile.OrchestratorType, properties.MasterProfile.Count)
		},
		"GetUniqueNameSuffix": func() string {
			return GenerateClusterID(properties)
		},
		"GetVNETAddressPrefixes": func() string {
			return getVNETAddressPrefixes(properties)
		},
		"GetVNETSubnets": func(addNSG bool) string {
			return getVNETSubnets(properties, addNSG)
		},
		"GetDataDisks": func(profile *api.AgentPoolProfile) string {
			return getDataDisks(profile)
		},
		"GetMasterAllowedSizes": func() string {
			if t.ClassicMode {
				return GetClassicAllowedSizes()
			}
			return GetMasterAllowedSizes()
		},
		"GetAgentAllowedSizes": func() string {
			if t.ClassicMode {
				return GetClassicAllowedSizes()
			}
			return GetAgentAllowedSizes()
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
		"GetKubernetesMasterCustomScript": func() string {
			return t.getBase64CustomScript(kubernetesMasterCustomScript)
		},
		"GetKubernetesMasterCustomData": func() string {
			str, e := t.getSingleLineForTemplate(kubernetesMasterCustomDataYaml)
			if e != nil {
				return ""
			}
			// add the master provisioning script
			masterProvisionB64GzipStr := t.getBase64CustomScript(kubernetesMasterCustomScript)
			str = strings.Replace(str, "MASTER_PROVISION_B64_GZIP_STR", masterProvisionB64GzipStr, -1)

			for placeholder, filename := range kubernetesAddonYamls {
				addonTextContents := t.getBase64CustomScript(filename)
				str = strings.Replace(str, placeholder, addonTextContents, -1)
			}

			// return the custom data
			return fmt.Sprintf("\"customData\": \"[base64(concat('%s'))]\",", str)
		},
		"GetKubernetesAgentCustomData": func(profile *api.AgentPoolProfile) string {
			str, e := t.getSingleLineForTemplate(kubernetesAgentCustomDataYaml)
			if e != nil {
				return ""
			}
			// add the agent provisioning script
			agentProvisionB64GzipStr := t.getBase64CustomScript(kubernetesAgentCustomScript)
			str = strings.Replace(str, "AGENT_PROVISION_B64_GZIP_STR", agentProvisionB64GzipStr, -1)

			return fmt.Sprintf("\"customData\": \"[base64(concat('%s'))]\",", str)
		},
		"GetKubernetesKubeConfig": func() string {
			str, e := t.getSingleLineForTemplate(kubeConfigJSON)
			if e != nil {
				return ""
			}
			return str
		},
		"AnyAgentHasDisks": func() bool {
			for _, agentProfile := range properties.AgentPoolProfiles {
				if agentProfile.HasDisks() {
					return true
				}
			}
			return false
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
	}
}

func getPackageGUID(orchestratorType api.OrchestratorType, masterCount int) string {
	if orchestratorType == api.DCOS || orchestratorType == api.DCOS184 {
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
	if orchestratorType == api.DCOS ||
		orchestratorType == api.DCOS173 ||
		orchestratorType == api.DCOS184 {
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

func getMasterRolesFileContents() string {
	return `{\"content\": \"\", \"path\": \"/etc/mesosphere/roles/master\"}, {\"content\": \"\", \"path\": \"/etc/mesosphere/roles/azure_master\"},`
}

func getAgentRolesFileContents(ports []int) string {
	if len(ports) > 0 {
		// public agents
		return `{\"content\": \"\", \"path\": \"/etc/mesosphere/roles/slave_public\"},`
	}
	// private agents
	return `{\"content\": \"\", \"path\": \"/etc/mesosphere/roles/slave\"},`
}

// getSingleLineForTemplate returns the file as a single line for embedding in an arm template
func (t *TemplateGenerator) getSingleLineForTemplate(yamlFilename string) (string, error) {
	b, err := Asset(yamlFilename)
	if err != nil {
		return "", fmt.Errorf("yaml file %s does not exist", yamlFilename)
	}
	// template.JSEscapeString leaves undesirable chars that don't work with pretty print
	yamlStr := string(b)
	yamlStr = strings.Replace(yamlStr, "\\", "\\\\", -1)
	yamlStr = strings.Replace(yamlStr, "\r\n", "\\n", -1)
	yamlStr = strings.Replace(yamlStr, "\n", "\\n", -1)
	yamlStr = strings.Replace(yamlStr, "\"", "\\\"", -1)

	// variable replacement
	rVariable, e1 := regexp.Compile("{{{([^}]*)}}}")
	if e1 != nil {
		return "", e1
	}
	yamlStr = rVariable.ReplaceAllString(yamlStr, "',variables('$1'),'")
	// verbatim replacement
	rVerbatim, e2 := regexp.Compile("<<<([^>]*)>>>")
	if e2 != nil {
		return "", e2
	}
	yamlStr = rVerbatim.ReplaceAllString(yamlStr, "',$1,'")
	return yamlStr, nil
}

// getBase64CustomScript will return a base64 of the CSE
func (t *TemplateGenerator) getBase64CustomScript(csFilename string) string {
	b, err := Asset(csFilename)
	if err != nil {
		// this should never happen and this is a bug
		panic(err.Error())
	}
	// translate the parameters
	csStr := string(b)
	csStr = strings.Replace(csStr, "\r\n", "\n", -1)

	var gzipB bytes.Buffer
	w := gzip.NewWriter(&gzipB)
	w.Write([]byte(csStr))
	w.Close()

	return base64.StdEncoding.EncodeToString(gzipB.Bytes())
}

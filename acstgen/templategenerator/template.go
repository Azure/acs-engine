package templategenerator

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"
	"text/template"

	"./../api/vlabs"
)

const (
	kubernetesMasterCustomDataYaml = "kubernetesmastercustomdata.yml"
	kubernetesAgentCustomDataYaml  = "kubernetesagentcustomdata.yml"
	kubeConfigJSON                 = "kubeconfig.json"
)

const (
	agentOutputs              = "agentoutputs.t"
	agentParams               = "agentparams.t"
	dcosAgentResources        = "dcosagentresources.t"
	dcosAgentResourcesDisks   = "dcosagentresourcesdisks.t"
	dcosAgentVars             = "dcosagentvars.t"
	dcosBaseFile              = "dcosbase.t"
	dcosCustomData173         = "dcoscustomdata173.t"
	dcosCustomData184         = "dcoscustomdata184.t"
	dcosMasterResources       = "dcosmasterresources.t"
	dcosMasterVars            = "dcosmastervars.t"
	kubernetesBaseFile        = "kubernetesbase.t"
	kubernetesAgentResources  = "kubernetesagentresources.t"
	kubernetesAgentVars       = "kubernetesagentvars.t"
	kubernetesMasterResources = "kubernetesmasterresources.t"
	kubernetesMasterVars      = "kubernetesmastervars.t"
	kubernetesParams          = "kubernetesparams.t"
	masterOutputs             = "masteroutputs.t"
	masterParams              = "masterparams.t"
	swarmBaseFile             = "swarmbase.t"
	swarmAgentCustomData      = "swarmagentcustomdata.t"
	swarmAgentResources       = "swarmagentresources.t"
	swarmAgentVars            = "swarmagentvars.t"
	swarmMasterCustomData     = "swarmmastercustomdata.t"
	swarmMasterResources      = "swarmmasterresources.t"
	swarmMasterVars           = "swarmmastervars.t"
)

var commonTemplateFiles = []string{agentOutputs, agentParams, masterOutputs, masterParams}
var dcosTemplateFiles = []string{dcosAgentResources, dcosAgentResourcesDisks, dcosAgentVars, dcosBaseFile, dcosCustomData173, dcosCustomData184, dcosMasterResources, dcosMasterVars}
var kubernetesTemplateFiles = []string{kubernetesBaseFile, kubernetesAgentResources, kubernetesAgentVars, kubernetesMasterResources, kubernetesMasterVars, kubernetesParams}
var swarmTemplateFiles = []string{swarmBaseFile, swarmAgentCustomData, swarmAgentResources, swarmAgentVars, swarmBaseFile, swarmMasterCustomData, swarmMasterResources, swarmMasterVars}

// VerifyFiles verifies that the required template files exist
func VerifyFiles(partsDirectory string) error {
	allFiles := append(commonTemplateFiles, dcosTemplateFiles...)
	allFiles = append(allFiles, kubernetesTemplateFiles...)
	allFiles = append(allFiles, swarmTemplateFiles...)
	for _, file := range allFiles {
		templateFile := path.Join(partsDirectory, file)
		if _, err := os.Stat(templateFile); os.IsNotExist(err) {
			return fmt.Errorf("template file %s does not exist, did you specify the correct template directory?", templateFile)
		}
	}
	return nil
}

// GenerateTemplate generates the template from the API Model
func GenerateTemplate(acsCluster *vlabs.AcsCluster, partsDirectory string) (string, error) {
	var err error
	var templ *template.Template

	templ = template.New("acs template").Funcs(getTemplateFuncMap(acsCluster, partsDirectory))

	var files []string
	var baseFile string
	if acsCluster.OrchestratorProfile.OrchestratorType == vlabs.DCOS184 ||
		acsCluster.OrchestratorProfile.OrchestratorType == vlabs.DCOS ||
		acsCluster.OrchestratorProfile.OrchestratorType == vlabs.DCOS173 {
		files = append(commonTemplateFiles, dcosTemplateFiles...)
		baseFile = dcosBaseFile
	} else if acsCluster.OrchestratorProfile.OrchestratorType == vlabs.Swarm {
		files = append(commonTemplateFiles, swarmTemplateFiles...)
		baseFile = swarmBaseFile
	} else if acsCluster.OrchestratorProfile.OrchestratorType == vlabs.Kubernetes {
		files = append(commonTemplateFiles, kubernetesTemplateFiles...)
		baseFile = kubernetesBaseFile
	} else {
		return "", fmt.Errorf("orchestrator '%s' is unsupported", acsCluster.OrchestratorProfile.OrchestratorType)
	}

	for _, file := range files {
		templateFile := path.Join(partsDirectory, file)
		bytes, e := ioutil.ReadFile(templateFile)
		if e != nil {
			return "", fmt.Errorf("Error reading file %s: %s", templateFile, e.Error())
		}
		if _, err = templ.New(file).Parse(string(bytes)); err != nil {
			return "", err
		}
	}
	var b bytes.Buffer
	if err = templ.ExecuteTemplate(&b, baseFile, acsCluster); err != nil {
		return "", err
	}

	return b.String(), nil
}

// getTemplateFuncMap returns all functions used in template generation
func getTemplateFuncMap(acsCluster *vlabs.AcsCluster, partsDirectory string) map[string]interface{} {
	return template.FuncMap{
		"IsDCOS173": func() bool {
			return acsCluster.OrchestratorProfile.OrchestratorType == vlabs.DCOS173
		},
		"IsDCOS184": func() bool {
			return acsCluster.OrchestratorProfile.OrchestratorType == vlabs.DCOS184 ||
				acsCluster.OrchestratorProfile.OrchestratorType == vlabs.DCOS
		},
		"IsPublic": func(ports []int) bool {
			return len(ports) > 0
		},
		"GetVNETSubnetDependencies": func() string {
			return getVNETSubnetDependencies(acsCluster)
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
			return getDCOSCustomDataPublicIPStr(acsCluster.OrchestratorProfile.OrchestratorType, acsCluster.MasterProfile.Count)
		},
		"GetDCOSGUID": func() string {
			return getPackageGUID(acsCluster.OrchestratorProfile.OrchestratorType, acsCluster.MasterProfile.Count)
		},
		"GetLinuxProfileFirstSSHPublicKey": func() string {
			return acsCluster.LinuxProfile.SSH.PublicKeys[0].KeyData
		},
		"GetUniqueNameSuffix": func() string {
			return acsCluster.OrchestratorProfile.ClusterID
		},
		"GetVNETAddressPrefixes": func() string {
			return getVNETAddressPrefixes(acsCluster)
		},
		"GetVNETSubnets": func(addNSG bool) string {
			return getVNETSubnets(acsCluster, addNSG)
		},
		"GetDataDisks": func(profile *vlabs.AgentPoolProfile) string {
			return getDataDisks(profile)
		},
		"GetMasterAllowedSizes": func() string {
			return GetMasterAllowedSizes()
		},
		"GetAgentAllowedSizes": func() string {
			return GetAgentAllowedSizes()
		},
		"GetSizeMap": func() string {
			return GetSizeMap()
		},
		"Base64": func(s string) string {
			return base64.URLEncoding.EncodeToString([]byte(s))
		},
		"GetKubernetesMasterCustomData": func() string {
			str, e := getSingleLineForTemplate(kubernetesMasterCustomDataYaml, partsDirectory)
			if e != nil {
				return ""
			}
			return fmt.Sprintf("\"customData\": \"[base64(concat('%s'))]\",", str)
		},
		"GetKubernetesAgentCustomData": func() string {
			str, e := getSingleLineForTemplate(kubernetesAgentCustomDataYaml, partsDirectory)
			if e != nil {
				return ""
			}
			return fmt.Sprintf("\"customData\": \"[base64(concat('%s'))]\",", str)
		},
		"GetKubernetesKubeConfig": func() string {
			str, e := getSingleLineForTemplate(kubeConfigJSON, partsDirectory)
			if e != nil {
				return ""
			}
			return str
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

func getPackageGUID(orchestratorType string, masterCount int) string {
	if orchestratorType == vlabs.DCOS || orchestratorType == vlabs.DCOS184 {
		switch masterCount {
		case 1:
			return "5ac6a7d060584c58c704e1f625627a591ecbde4e"
		case 3:
			return "42bd1d74e9a2b23836bd78919c716c20b98d5a0e"
		case 5:
			return "97947a91e2c024ed4f043bfcdad49da9418d3095"
		}
	} else if orchestratorType == vlabs.DCOS173 {
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

func getDCOSCustomDataPublicIPStr(orchestratorType string, masterCount int) string {
	if orchestratorType == vlabs.DCOS ||
		orchestratorType == vlabs.DCOS173 ||
		orchestratorType == vlabs.DCOS184 {
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

func getVNETAddressPrefixes(acsCluster *vlabs.AcsCluster) string {
	visitedSubnets := make(map[string]bool)
	var buf bytes.Buffer
	buf.WriteString(`"[variables('masterSubnet')]"`)
	visitedSubnets[acsCluster.MasterProfile.GetSubnet()] = true
	for i := range acsCluster.AgentPoolProfiles {
		profile := &acsCluster.AgentPoolProfiles[i]
		if _, ok := visitedSubnets[profile.GetSubnet()]; !ok {
			buf.WriteString(fmt.Sprintf(",\n            \"[variables('%sSubnet')]\"", profile.Name))
		}
	}
	return buf.String()
}

func getVNETSubnetDependencies(acsCluster *vlabs.AcsCluster) string {
	agentString := `        "[concat('Microsoft.Network/networkSecurityGroups/', variables('%sNSGName'))]"`
	var buf bytes.Buffer
	for index, agentProfile := range acsCluster.AgentPoolProfiles {
		if index > 0 {
			buf.WriteString(",\n")
		}
		buf.WriteString(fmt.Sprintf(agentString, agentProfile.Name))
	}
	return buf.String()
}

func getVNETSubnets(acsCluster *vlabs.AcsCluster, addNSG bool) string {
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
	for _, agentProfile := range acsCluster.AgentPoolProfiles {
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

func getDataDisks(a *vlabs.AgentPoolProfile) string {
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
	for i, diskSize := range a.DiskSizesGB {
		if i > 0 {
			buf.WriteString(",\n")
		}
		buf.WriteString(fmt.Sprintf(dataDisks, diskSize, i, a.Name, i, a.Name, a.Name, a.Name, a.Name, i))
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
func getSingleLineForTemplate(yamlFilename string, partsDirectory string) (string, error) {
	yamlFile := path.Join(partsDirectory, yamlFilename)
	if _, err := os.Stat(yamlFile); os.IsNotExist(err) {
		return "", fmt.Errorf("yaml file %s does not exist, did you specify the correct template directory?", yamlFile)
	}
	b, err := ioutil.ReadFile(yamlFile)
	if err != nil {
		return "", fmt.Errorf("error reading yaml file %s: %s", yamlFile, err.Error())
	}
	yamlStr := string(b)
	yamlStr = strings.Replace(yamlStr, "\\", "\\\\", -1)
	yamlStr = strings.Replace(yamlStr, "\r\n", "\\n", -1)
	yamlStr = strings.Replace(yamlStr, "\n", "\\n", -1)
	yamlStr = strings.Replace(yamlStr, "\"", "\\\"", -1)

	// variable replacement
	rVariable := regexp.MustCompile("{{{([^}]*)}}}")
	yamlStr = rVariable.ReplaceAllString(yamlStr, "',variables('$1'),'")
	// verbatim replacement
	rVerbatim := regexp.MustCompile("<<<([^>]*)>>>")
	yamlStr = rVerbatim.ReplaceAllString(yamlStr, "',$1,'")
	return yamlStr, nil
}

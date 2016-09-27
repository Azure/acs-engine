package clustertemplate

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"text/template"

	"./../api/vlabs"
)

const (
	agentOutputs         = "agentoutputs.t"
	agentParams          = "agentparams.t"
	dcosAgentResources   = "dcosagentresources.t"
	dcosAgentVars        = "dcosagentvars.t"
	dcosBaseFile         = "dcosbase.t"
	dcosCustomData173    = "dcoscustomdata173.t"
	dcosCustomData184    = "dcoscustomdata184.t"
	dcosMasterResources  = "dcosmasterresources.t"
	dcosMasterVars       = "dcosmastervars.t"
	masterOutputs        = "masteroutputs.t"
	masterParams         = "masterparams.t"
	swarmBaseFile        = "swarmbase.t"
	swarmCustomData      = "swarmcustomdata.t"
	swarmMasterResources = "swarmmasterresources.t"
	swarmMasterVars      = "swarmmastervars.t"
)

var dcosTemplateFiles = []string{agentOutputs, agentParams, dcosAgentResources, dcosAgentVars, dcosBaseFile, dcosCustomData173, dcosCustomData184, dcosMasterResources, dcosMasterVars, masterOutputs, masterParams}
var swarmTemplateFiles = []string{agentOutputs, agentParams, swarmBaseFile, masterOutputs, masterParams, swarmBaseFile, swarmCustomData, swarmMasterResources, swarmMasterVars}

// VerifyFiles verifies that the required template files exist
func VerifyFiles(partsDirectory string) error {
	for _, file := range append(dcosTemplateFiles, swarmTemplateFiles...) {
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

	templateMap := template.FuncMap{
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
			return generateUniqueNameSuffix(acsCluster)
		},
		"GetVNETAddressPrefixes": func() string {
			return getVNETAddressPrefixes(acsCluster)
		},
		"GetVNETSubnets": func() string {
			return getVNETSubnets(acsCluster)
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
	templ = template.New("acs template").Funcs(templateMap)

	var files []string
	var baseFile string
	if isDCOS(acsCluster) {
		files = dcosTemplateFiles
		baseFile = dcosBaseFile
	} else if isSwarm(acsCluster) {
		files = swarmTemplateFiles
		baseFile = swarmBaseFile
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

func generateUniqueNameSuffix(acsCluster *vlabs.AcsCluster) string {
	uniqueNameSuffixSize := 8
	var seed int64
	for _, c := range acsCluster.MasterProfile.DNSPrefix {
		seed += int64(c)
	}
	rand.Seed(seed)
	return fmt.Sprintf("%08d", rand.Uint32())[:uniqueNameSuffixSize]
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
	var buf bytes.Buffer
	buf.WriteString(`"[variables('masterSubnet')]"`)
	for _, agentProfile := range acsCluster.AgentPoolProfiles {
		buf.WriteString(fmt.Sprintf(",\n            \"[variables('%sAddressPrefix')]\"", agentProfile.Name))
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

func getVNETSubnets(acsCluster *vlabs.AcsCluster) string {
	masterString := `{
            "name": "[variables('masterSubnetName')]", 
            "properties": {
              "addressPrefix": "[variables('masterSubnet')]"
            }
          }`
	agentString := `          {
            "name": "[variables('%sSubnetName')]", 
            "properties": {
              "addressPrefix": "[variables('%sAddressPrefix')]", 
              "networkSecurityGroup": {
                "id": "[resourceId('Microsoft.Network/networkSecurityGroups', variables('%sNSGName'))]"
              }
            }
          }`
	var buf bytes.Buffer
	buf.WriteString(masterString)
	for _, agentProfile := range acsCluster.AgentPoolProfiles {
		buf.WriteString(",\n")
		buf.WriteString(fmt.Sprintf(agentString, agentProfile.Name, agentProfile.Name, agentProfile.Name))
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
          }`, port, port, port, vlabs.BaseLBPriority+portIndex)
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
	} else {
		// private agents
		return `{\"content\": \"\", \"path\": \"/etc/mesosphere/roles/slave\"},`
	}
}

func isDCOS(acsCluster *vlabs.AcsCluster) bool {
	return acsCluster.OrchestratorProfile.OrchestratorType == vlabs.DCOS184 ||
		acsCluster.OrchestratorProfile.OrchestratorType == vlabs.DCOS ||
		acsCluster.OrchestratorProfile.OrchestratorType == vlabs.DCOS173
}

func isSwarm(acsCluster *vlabs.AcsCluster) bool {
	return acsCluster.OrchestratorProfile.OrchestratorType == vlabs.SWARM
}

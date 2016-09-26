package clustertemplate

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"text/template"
	"time"

	"./../api/vlabs"
)

const (
	baseFile             = "base.t"
	masterParams         = "masterparams.t"
	masterOutputs        = "masteroutputs.t"
	dcosMasterVars       = "dcosmastervars.t"
	dcosMasterResources  = "dcosmasterresources.t"
	dcosCustomData173    = "dcoscustomdata173.t"
	dcosCustomData184    = "dcoscustomdata184.t"
	swarmMasterVars      = "swarmmastervars.t"
	swarmMasterResources = "swarmmasterresources.t"
	swarmCustomData      = "swarmcustomdata.t"
)

var templateFiles = []string{baseFile, masterParams, masterOutputs, dcosMasterVars, dcosMasterResources, dcosCustomData173, dcosCustomData184, swarmMasterVars, swarmMasterResources, swarmCustomData}

// ClusterContext contains the template context for binding during template generation
type ClusterContext struct {
	vlabs.AcsCluster
	LinuxProfileFirstSSHPublicKey string
	UniqueNameSuffix              string
	VNETAddressPrefixes           string
	VNETSubnets                   string
	DCOSGUID                      string
	DCOSCustomDataPublicIPStr     string
}

// VerifyFiles verifies that the required template files exist
func VerifyFiles(partsDirectory string) error {
	for _, file := range templateFiles {
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
	var clusterContext *ClusterContext

	// build the context for binding to the template
	clusterContext, err = generateContext(acsCluster)
	if err != nil {
		return "", err
	}
	templateMap := template.FuncMap{
		"IsDCOS173": func() bool {
			return clusterContext.OrchestratorProfile.OrchestratorType == vlabs.DCOS173
		},
		"IsDCOS184": func() bool {
			return clusterContext.OrchestratorProfile.OrchestratorType == vlabs.DCOS184 ||
				clusterContext.OrchestratorProfile.OrchestratorType == vlabs.DCOS
		},
		"IsDCOS": func() bool {
			return clusterContext.OrchestratorProfile.OrchestratorType == vlabs.DCOS184 ||
				clusterContext.OrchestratorProfile.OrchestratorType == vlabs.DCOS ||
				clusterContext.OrchestratorProfile.OrchestratorType == vlabs.DCOS173
		},
		"IsSwarm": func() bool {
			return clusterContext.OrchestratorProfile.OrchestratorType == vlabs.SWARM
		},
	}
	templ = template.New("acs template").Funcs(templateMap)

	for _, file := range templateFiles {
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
	if err = templ.ExecuteTemplate(&b, baseFile, clusterContext); err != nil {
		return "", err
	}

	return b.String(), nil
}

func generateContext(acsCluster *vlabs.AcsCluster) (*ClusterContext, error) {
	clusterContext := &ClusterContext{}
	clusterContext.AcsCluster = *acsCluster
	clusterContext.LinuxProfileFirstSSHPublicKey = acsCluster.LinuxProfile.SSH.PublicKeys[0].KeyData
	clusterContext.UniqueNameSuffix = generateUniqueNameSuffix()
	clusterContext.VNETAddressPrefixes = getVNETAddressPrefixes()
	clusterContext.VNETSubnets = getVNETSubnets()
	clusterContext.DCOSGUID = getPackageGUID(clusterContext.OrchestratorProfile.OrchestratorType, clusterContext.MasterProfile.Count)
	clusterContext.DCOSCustomDataPublicIPStr = getDCOSCustomDataPublicIPStr(clusterContext.OrchestratorProfile.OrchestratorType, clusterContext.MasterProfile.Count)
	return clusterContext, nil
}

func generateUniqueNameSuffix() string {
	uniqueNameSuffixSize := 8
	rand.Seed(time.Now().UTC().UnixNano())
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

func getVNETAddressPrefixes() string {
	return `"[variables('masterSubnet')]"`
}

func getVNETSubnets() string {
	return `{
            "name": "[variables('masterSubnetName')]", 
            "properties": {
              "addressPrefix": "[variables('masterSubnet')]"
            }
          }`
}

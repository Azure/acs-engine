package acsengine

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/Azure/acs-engine/pkg/api"
)

func kubernetesCustomFiles(profile *api.Properties) []api.CustomFile {
	if profile.OrchestratorProfile.KubernetesConfig.CustomFiles != nil {
		return *profile.OrchestratorProfile.KubernetesConfig.CustomFiles
	}
	return []api.CustomFile{}
}

func substituteConfigStringCustomFiles(input string, customFiles []api.CustomFile, placeholder string) string {

	var config string
	for _, customFile := range customFiles {
		config += buildConfigStringCustomFiles(
			customFile.Source,
			customFile.Dest)

	}
	return strings.Replace(input, placeholder, config, -1)
}

func buildConfigStringCustomFiles(sourceFile string, destinationFile string) string {
	contents := []string{
		fmt.Sprintf("- path: %s", destinationFile),
		"  permissions: \\\"0644\\\"",
		"  encoding: gzip",
		"  owner: \\\"root\\\"",
		"  content: !!binary |",
		fmt.Sprintf("    %s\\n\\n", getBase64CustomFile(sourceFile)),
	}

	return strings.Join(contents, "\\n")
}

func getBase64CustomFile(cfFilepath string) string {
	dat, err := ioutil.ReadFile(cfFilepath)
	if err != nil {
		panic(fmt.Sprintf("Could not read custom file: %s", err.Error()))
	}
	csStr := string(dat)
	csStr = strings.Replace(csStr, "\r\n", "\n", -1)
	return getBase64CustomScriptFromStr(csStr)
}

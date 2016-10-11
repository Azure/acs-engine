package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"./api/vlabs"
	"./tgen"
)

// loadAcsCluster loads an ACS Cluster API Model from a JSON file
func loadAcsCluster(jsonFile string) (*vlabs.AcsCluster, error) {
	contents, e := ioutil.ReadFile(jsonFile)
	if e != nil {
		return nil, fmt.Errorf("error reading file %s: %s", jsonFile, e.Error())
	}

	acsCluster := &vlabs.AcsCluster{}
	if e := json.Unmarshal(contents, &acsCluster); e != nil {
		return nil, fmt.Errorf("error unmarshalling file %s: %s", jsonFile, e.Error())
	}

	if e := acsCluster.Validate(); e != nil {
		return nil, fmt.Errorf("error validating acs cluster from file %s: %s", jsonFile, e.Error())
	}

	return acsCluster, nil
}

func translateJSON(content string, translateParams [][]string, reverseTranslate bool) string {
	for _, tuple := range translateParams {
		if len(tuple) != 2 {
			panic("string tuples must be of size 2")
		}
		a := tuple[0]
		b := tuple[1]
		if reverseTranslate {
			content = strings.Replace(content, b, a, -1)
		} else {
			content = strings.Replace(content, a, b, -1)
		}
	}
	return content
}

func prettyPrintJSON(content string) (string, error) {
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(content), &data); err != nil {
		return "", err
	}
	prettyprint, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", err
	}
	return string(prettyprint), nil
}

func prettyPrintArmTemplate(template string) (string, error) {
	translateParams := [][]string{
		{"\"parameters\"", "\"dparameters\""},
		{"\"variables\"", "\"evariables\""},
		{"\"resources\"", "\"fresources\""},
		{"\"outputs\"", "\"zoutputs\""},
		// there is a bug in ARM where it doesn't correctly translate back '\u003e' (>)
		{">", "GREATERTHAN"},
		{"<", "LESSTHAN"},
		{"&", "AMPERSAND"},
	}

	template = translateJSON(template, translateParams, false)
	var err error
	if template, err = prettyPrintJSON(template); err != nil {
		return "", err
	}
	template = translateJSON(template, translateParams, true)

	return template, nil
}

func writeArtifacts(acsCluster *vlabs.AcsCluster, template string, parameters, artifactsDir string, templateDirectory string, certsGenerated bool) error {
	if len(artifactsDir) == 0 {
		artifactsDir = fmt.Sprintf("%s-%s", acsCluster.OrchestratorProfile.OrchestratorType, tgen.GenerateClusterID(acsCluster))
		artifactsDir = path.Join("output", artifactsDir)
	}

	b, err := json.MarshalIndent(acsCluster, "", "  ")
	if err != nil {
		return err
	}

	if e := saveFile(artifactsDir, "apimodel.json", b); e != nil {
		return e
	}

	if e := saveFileString(artifactsDir, "azuredeploy.json", template); e != nil {
		return e
	}

	if e := saveFileString(artifactsDir, "azuredeploy.parameters.json", parameters); e != nil {
		return e
	}

	if certsGenerated {
		if acsCluster.OrchestratorProfile.OrchestratorType == vlabs.Kubernetes {
			directory := path.Join(artifactsDir, "kubeconfig")
			for _, location := range tgen.AzureLocations {
				b, gkcerr := tgen.GenerateKubeConfig(acsCluster, templateDirectory, location)
				if gkcerr != nil {
					return gkcerr
				}
				if e := saveFileString(directory, fmt.Sprintf("kubeconfig.%s.json", location), b); e != nil {
					return e
				}
			}
		}

		if e := saveFileString(artifactsDir, "ca.key", acsCluster.CertificateProfile.GetCAPrivateKey()); e != nil {
			return e
		}
		if e := saveFileString(artifactsDir, "ca.crt", acsCluster.CertificateProfile.CaCertificate); e != nil {
			return e
		}
		if e := saveFileString(artifactsDir, "apiserver.key", acsCluster.CertificateProfile.APIServerPrivateKey); e != nil {
			return e
		}
		if e := saveFileString(artifactsDir, "apiserver.crt", acsCluster.CertificateProfile.APIServerCertificate); e != nil {
			return e
		}
		if e := saveFileString(artifactsDir, "client.key", acsCluster.CertificateProfile.ClientPrivateKey); e != nil {
			return e
		}
		if e := saveFileString(artifactsDir, "client.crt", acsCluster.CertificateProfile.ClientCertificate); e != nil {
			return e
		}
		if e := saveFileString(artifactsDir, "kubectlClient.key", acsCluster.CertificateProfile.KubeConfigPrivateKey); e != nil {
			return e
		}
		if e := saveFileString(artifactsDir, "kubectlClient.crt", acsCluster.CertificateProfile.KubeConfigCertificate); e != nil {
			return e
		}
	}

	return nil
}

func saveFileString(dir string, file string, data string) error {
	return saveFile(dir, file, []byte(data))
}

func saveFile(dir string, file string, data []byte) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if e := os.MkdirAll(dir, 0700); e != nil {
			return fmt.Errorf("error creating directory '%s': %s", dir, e.Error())
		}
	}

	path := path.Join(dir, file)
	if err := ioutil.WriteFile(path, []byte(data), 0600); err != nil {
		return err
	}

	fmt.Fprintf(os.Stderr, "wrote %s\n", path)

	return nil
}

func usage(errs ...error) {
	for _, err := range errs {
		fmt.Fprintf(os.Stderr, "error: %s\n\n", err.Error())
	}
	fmt.Fprintf(os.Stderr, "usage: %s [OPTIONS] ClusterDefinitionFile\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "       read the ClusterDefinitionFile and output an arm template")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "options:\n")
	flag.PrintDefaults()
}

var templateDirectory = flag.String("templateDirectory", "./parts", "directory containing base template files")
var noPrettyPrint = flag.Bool("noPrettyPrint", false, "do not pretty print output")
var artifactsDir = flag.String("artifacts", "", "directory where artifacts will be written")
var classicMode = flag.Bool("classicMode", false, "enable classic parameters and outputs")

func main() {
	var acsCluster *vlabs.AcsCluster
	var template string
	var parameters string
	var err error

	flag.Parse()

	if argCount := len(flag.Args()); argCount == 0 {
		usage()
		os.Exit(1)
	}

	jsonFile := flag.Arg(0)
	if _, err = os.Stat(jsonFile); os.IsNotExist(err) {
		usage(fmt.Errorf("file %s does not exist", jsonFile))
		os.Exit(1)
	}

	if _, err = os.Stat(*templateDirectory); os.IsNotExist(err) {
		usage(fmt.Errorf("base templates directory %s does not exist", jsonFile))
		os.Exit(1)
	}

	if err = tgen.VerifyFiles(*templateDirectory); err != nil {
		fmt.Fprintf(os.Stderr, "verification failed: %s\n", err.Error())
		os.Exit(1)
	}

	if acsCluster, err = loadAcsCluster(jsonFile); err != nil {
		fmt.Fprintf(os.Stderr, "error while loading %s: %s", jsonFile, err.Error())
		os.Exit(1)
	}

	certsGenerated := false
	if certsGenerated, err = tgen.SetAcsClusterDefaults(acsCluster); err != nil {
		fmt.Fprintf(os.Stderr, "error while setting defaults %s: %s", jsonFile, err.Error())
		os.Exit(1)
	}

	if *classicMode {
		acsCluster.SetClassicMode(true)
	}

	if template, parameters, err = tgen.GenerateTemplate(acsCluster, *templateDirectory); err != nil {
		fmt.Fprintf(os.Stderr, "error generating template %s: %s", jsonFile, err.Error())
		os.Exit(1)
	}

	if !*noPrettyPrint {
		if template, err = prettyPrintArmTemplate(template); err != nil {
			fmt.Fprintf(os.Stderr, "error pretty printing template %s", err.Error())
			os.Exit(1)
		}
		if parameters, err = prettyPrintArmTemplate(parameters); err != nil {
			fmt.Fprintf(os.Stderr, "error pretty printing template %s", err.Error())
			os.Exit(1)
		}
	}

	if err = writeArtifacts(acsCluster, template, parameters, *artifactsDir, *templateDirectory, certsGenerated); err != nil {
		fmt.Fprintf(os.Stderr, "error writing artifacts %s", err.Error())
		os.Exit(1)
	}
}

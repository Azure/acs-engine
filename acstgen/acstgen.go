package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"time"

	"github.com/Azure/acs-labs/acstgen/pkg/api"
	"github.com/Azure/acs-labs/acstgen/pkg/api/v20160330"
	"github.com/Azure/acs-labs/acstgen/pkg/api/vlabs"
	"github.com/Azure/acs-labs/acstgen/pkg/tgen"
)

func writeArtifacts(containerService *api.ContainerService, template string, parameters, artifactsDir string, templateDirectory string, certsGenerated bool) error {
	if len(artifactsDir) == 0 {
		artifactsDir = fmt.Sprintf("%s-%s", containerService.Properties.OrchestratorProfile.OrchestratorType, tgen.GenerateClusterID(&containerService.Properties))
		artifactsDir = path.Join("_output", artifactsDir)
	}

	// convert back the API object, and write it
	var b []byte
	var err error
	switch containerService.APIVersion {
	case v20160330.APIVersion:
		v20160330ContainerService := &v20160330.ContainerService{}
		api.ConvertContainerServiceToV20160330(containerService, v20160330ContainerService)
		b, err = json.MarshalIndent(v20160330ContainerService, "", "  ")

	case vlabs.APIVersion:
		vlabsContainerService := &vlabs.ContainerService{}
		api.ConvertContainerServiceToVLabs(containerService, vlabsContainerService)
		b, err = json.MarshalIndent(vlabsContainerService, "", "  ")

	default:
		return fmt.Errorf("invalid version %s for conversion back from unversioned object", containerService.APIVersion)
	}

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
		properties := &containerService.Properties
		if properties.OrchestratorProfile.OrchestratorType == vlabs.Kubernetes {
			directory := path.Join(artifactsDir, "kubeconfig")
			for _, location := range tgen.AzureLocations {
				b, gkcerr := tgen.GenerateKubeConfig(properties, templateDirectory, location)
				if gkcerr != nil {
					return gkcerr
				}
				if e := saveFileString(directory, fmt.Sprintf("kubeconfig.%s.json", location), b); e != nil {
					return e
				}
			}
		}

		if e := saveFileString(artifactsDir, "ca.key", properties.CertificateProfile.GetCAPrivateKey()); e != nil {
			return e
		}
		if e := saveFileString(artifactsDir, "ca.crt", properties.CertificateProfile.CaCertificate); e != nil {
			return e
		}
		if e := saveFileString(artifactsDir, "apiserver.key", properties.CertificateProfile.APIServerPrivateKey); e != nil {
			return e
		}
		if e := saveFileString(artifactsDir, "apiserver.crt", properties.CertificateProfile.APIServerCertificate); e != nil {
			return e
		}
		if e := saveFileString(artifactsDir, "client.key", properties.CertificateProfile.ClientPrivateKey); e != nil {
			return e
		}
		if e := saveFileString(artifactsDir, "client.crt", properties.CertificateProfile.ClientCertificate); e != nil {
			return e
		}
		if e := saveFileString(artifactsDir, "kubectlClient.key", properties.CertificateProfile.KubeConfigPrivateKey); e != nil {
			return e
		}
		if e := saveFileString(artifactsDir, "kubectlClient.crt", properties.CertificateProfile.KubeConfigCertificate); e != nil {
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
	start := time.Now()
	defer func(s time.Time) {
		fmt.Fprintf(os.Stderr, "acstgen took %s\n", time.Since(s))
	}(start)
	var containerService *api.ContainerService
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

	if containerService, err = tgen.LoadContainerServiceFromFile(jsonFile); err != nil {
		fmt.Fprintf(os.Stderr, "error while loading %s: %s", jsonFile, err.Error())
		os.Exit(1)
	}

	if *classicMode {
		containerService.Properties.SetClassicMode(true)
	}

	certsGenerated := false
	if template, parameters, certsGenerated, err = tgen.GenerateTemplate(containerService, *templateDirectory); err != nil {
		fmt.Fprintf(os.Stderr, "error generating template %s: %s", jsonFile, err.Error())
		os.Exit(1)
	}

	if !*noPrettyPrint {
		if template, err = tgen.PrettyPrintArmTemplate(template); err != nil {
			fmt.Fprintf(os.Stderr, "error pretty printing template %s", err.Error())
			os.Exit(1)
		}
		if parameters, err = tgen.PrettyPrintJSON(parameters); err != nil {
			fmt.Fprintf(os.Stderr, "error pretty printing template %s", err.Error())
			os.Exit(1)
		}
	}

	if err = writeArtifacts(containerService, template, parameters, *artifactsDir, *templateDirectory, certsGenerated); err != nil {
		fmt.Fprintf(os.Stderr, "error writing artifacts %s", err.Error())
		os.Exit(1)
	}
}

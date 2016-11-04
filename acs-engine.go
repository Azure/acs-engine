package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"time"

	"github.com/Azure/acs-engine/pkg/acsengine"
	"github.com/Azure/acs-engine/pkg/api"
	"github.com/Azure/acs-engine/pkg/api/vlabs"
)

func writeArtifacts(containerService *api.ContainerService, template string, parameters, artifactsDir string, certsGenerated bool, parametersOnly bool) error {
	if len(artifactsDir) == 0 {
		artifactsDir = fmt.Sprintf("%s-%s", containerService.Properties.OrchestratorProfile.OrchestratorType, acsengine.GenerateClusterID(&containerService.Properties))
		artifactsDir = path.Join("_output", artifactsDir)
	}

	// convert back the API object, and write it
	var b []byte
	var err error
	if !parametersOnly {
		b, err = api.SerializeContainerService(containerService)

		if err != nil {
			return err
		}

		if e := saveFile(artifactsDir, "apimodel.json", b); e != nil {
			return e
		}

		if e := saveFileString(artifactsDir, "azuredeploy.json", template); e != nil {
			return e
		}
	}

	if e := saveFileString(artifactsDir, "azuredeploy.parameters.json", parameters); e != nil {
		return e
	}

	if certsGenerated {
		properties := &containerService.Properties
		if properties.OrchestratorProfile.OrchestratorType == vlabs.Kubernetes {
			directory := path.Join(artifactsDir, "kubeconfig")
			for _, location := range acsengine.AzureLocations {
				b, gkcerr := acsengine.GenerateKubeConfig(properties, location)
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

var noPrettyPrint = flag.Bool("noPrettyPrint", false, "do not pretty print output")
var artifactsDir = flag.String("artifacts", "", "directory where artifacts will be written")
var classicMode = flag.Bool("classicMode", false, "enable classic parameters and outputs")
var parametersOnly = flag.Bool("parametersOnly", false, "only output the parameters")

func main() {
	start := time.Now()
	defer func(s time.Time) {
		fmt.Fprintf(os.Stderr, "acsengine took %s\n", time.Since(s))
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

	templateGenerator, e := acsengine.InitializeTemplateGenerator(*classicMode)
	if e != nil {
		fmt.Fprintf(os.Stderr, "generator initialization failed: %s\n", e.Error())
		os.Exit(1)
	}

	if containerService, err = api.LoadContainerServiceFromFile(jsonFile); err != nil {
		fmt.Fprintf(os.Stderr, "error while loading %s: %s", jsonFile, err.Error())
		os.Exit(1)
	}

	certsGenerated := false
	if template, parameters, certsGenerated, err = templateGenerator.GenerateTemplate(containerService); err != nil {
		fmt.Fprintf(os.Stderr, "error generating template %s: %s", jsonFile, err.Error())
		os.Exit(1)
	}

	if !*noPrettyPrint {
		if template, err = acsengine.PrettyPrintArmTemplate(template); err != nil {
			fmt.Fprintf(os.Stderr, "error pretty printing template %s", err.Error())
			os.Exit(1)
		}
		if parameters, err = acsengine.PrettyPrintJSON(parameters); err != nil {
			fmt.Fprintf(os.Stderr, "error pretty printing template %s", err.Error())
			os.Exit(1)
		}
	}

	if err = writeArtifacts(containerService, template, parameters, *artifactsDir, certsGenerated, *parametersOnly); err != nil {
		fmt.Fprintf(os.Stderr, "error writing artifacts %s", err.Error())
		os.Exit(1)
	}
}

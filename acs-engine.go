package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"time"

	"github.com/Azure/acs-engine/pkg/acsengine"
	"github.com/Azure/acs-engine/pkg/api"
)

func writeArtifacts(containerService *api.ContainerService, apiVersion, template, parameters, artifactsDir string, certsGenerated bool, parametersOnly bool) error {
	if len(artifactsDir) == 0 {
		artifactsDir = fmt.Sprintf("%s-%s", containerService.Properties.OrchestratorProfile.OrchestratorType, acsengine.GenerateClusterID(containerService.Properties))
		artifactsDir = path.Join("_output", artifactsDir)
	}

	// convert back the API object, and write it
	var b []byte
	var err error
	if !parametersOnly {
		b, err = api.SerializeContainerService(containerService, apiVersion)

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
		properties := containerService.Properties
		if properties.OrchestratorProfile.OrchestratorType == api.Kubernetes {
			directory := path.Join(artifactsDir, "kubeconfig")
			var locations []string
			if containerService.Location != "" {
				locations = []string{containerService.Location}
			} else {
				locations = acsengine.AzureLocations
			}

			for _, location := range locations {
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
	fmt.Fprintf(os.Stderr, "\n%s --version or -v to get the build version \n", os.Args[0])
}

var noPrettyPrint = flag.Bool("noPrettyPrint", false, "do not pretty print output")
var artifactsDir = flag.String("artifacts", "", "directory where artifacts will be written")
var classicMode = flag.Bool("classicMode", false, "enable classic parameters and outputs")
var parametersOnly = flag.Bool("parametersOnly", false, "only output the parameters")

// AcsEngineBuildSHA is the Git SHA-1 of the last commit
var AcsEngineBuildSHA string

// AcsEngineBuildTime is the timestamp of when acs-engine was built
var AcsEngineBuildTime string

// acs-engine takes the caKey and caCert as args, since the caKey is stored separately
// from the api model since this cannot be easily revoked like the server and client key
var caCertificatePath = flag.String("caCertificatePath", "", "the path to the CA Certificate file")
var caKeyPath = flag.String("caKeyPath", "", "the path to the CA key file")

func main() {
	if (len(os.Args) == 2) &&
		((os.Args[1] == "--version") || (os.Args[1] == "-v")) {
		if len(AcsEngineBuildSHA) == 0 {
			fmt.Fprintf(os.Stderr, "No version set. Please run `make build`\n")
			os.Exit(1)
		}
		fmt.Fprintf(os.Stdout, "Git commit: %s\nBuild Timestamp: %s\n", AcsEngineBuildSHA, AcsEngineBuildTime)
		os.Exit(0)
	}

	start := time.Now()
	defer func(s time.Time) {
		fmt.Fprintf(os.Stderr, "acsengine took %s\n", time.Since(s))
	}(start)
	var containerService *api.ContainerService
	var caCertificateBytes []byte
	var caKeyBytes []byte
	var template string
	var parameters string
	var apiVersion string
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

	if (len(*caCertificatePath) > 0 && len(*caKeyPath) == 0) ||
		(len(*caCertificatePath) == 0 && len(*caKeyPath) > 0) {
		usage(errors.New("caKeyPath and caCertificatePath must be specified together"))
		os.Exit(1)
	}
	if len(*caCertificatePath) > 0 {
		if caCertificateBytes, err = ioutil.ReadFile(*caCertificatePath); err != nil {
			usage(err)
			os.Exit(1)
		}
		if caKeyBytes, err = ioutil.ReadFile(*caKeyPath); err != nil {
			usage(err)
			os.Exit(1)
		}
	}

	templateGenerator, e := acsengine.InitializeTemplateGenerator(*classicMode)
	if e != nil {
		fmt.Fprintf(os.Stderr, "generator initialization failed: %s\n", e.Error())
		os.Exit(1)
	}

	if containerService, apiVersion, err = api.LoadContainerServiceFromFile(jsonFile); err != nil {
		fmt.Fprintf(os.Stderr, "error while loading %s: %s", jsonFile, err.Error())
		os.Exit(1)
	}

	if len(caKeyBytes) != 0 {
		// the caKey is not in the api model, and should be stored separately from the model
		// we put these in the model after model is deserialized
		containerService.Properties.CertificateProfile.CaCertificate = string(caCertificateBytes)
		containerService.Properties.CertificateProfile.SetCAPrivateKey(string(caKeyBytes))
	}

	certsGenerated := false
	if template, parameters, certsGenerated, err = templateGenerator.GenerateTemplate(containerService); err != nil {
		fmt.Fprintf(os.Stderr, "error generating template %s: %s", jsonFile, err.Error())
		os.Exit(1)
	}

	if !*noPrettyPrint {
		if template, err = acsengine.PrettyPrintArmTemplate(template); err != nil {
			fmt.Fprintf(os.Stderr, "error pretty printing template: %s \n", err.Error())
			os.Exit(1)
		}
		if parameters, err = acsengine.PrettyPrintJSON(parameters); err != nil {
			fmt.Fprintf(os.Stderr, "error pretty printing template parameters: %s \n", err.Error())
			os.Exit(1)
		}
	}

	if err = writeArtifacts(containerService, apiVersion, template, parameters, *artifactsDir, certsGenerated, *parametersOnly); err != nil {
		fmt.Fprintf(os.Stderr, "error writing artifacts: %s \n", err.Error())
		os.Exit(1)
	}
}

package acsengine

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/Azure/acs-engine/pkg/api"
	log "github.com/Sirupsen/logrus"
)

func WriteArtifacts(containerService *api.ContainerService, apiVersion, template, parameters, artifactsDir string, certsGenerated bool, parametersOnly bool) error {
	if len(artifactsDir) == 0 {
		artifactsDir = fmt.Sprintf("%s-%s", containerService.Properties.OrchestratorProfile.OrchestratorType, GenerateClusterID(containerService.Properties))
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
				locations = AzureLocations
			}

			for _, location := range locations {
				b, gkcerr := GenerateKubeConfig(properties, location)
				if gkcerr != nil {
					return gkcerr
				}
				if e := saveFileString(directory, fmt.Sprintf("kubeconfig.%s.json", location), b); e != nil {
					return e
				}
			}

		}

		if e := saveFileString(artifactsDir, "ca.key", properties.CertificateProfile.CaPrivateKey); e != nil {
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

	log.Debugf("output: wrote %s", path)

	return nil
}

package acsengine

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/Azure/acs-engine/pkg/api"
	"github.com/Azure/acs-engine/pkg/i18n"
)

// ArtifactWriter represents the object that writes artifacts
type ArtifactWriter struct {
	Translator *i18n.Translator
}

func (w *ArtifactWriter) WriteArtifacts(containerService *api.ContainerService, apiVersion, template, parameters, artifactsDir string, certsGenerated bool, parametersOnly bool) error {
	if len(artifactsDir) == 0 {
		artifactsDir = fmt.Sprintf("%s-%s", containerService.Properties.OrchestratorProfile.OrchestratorType, GenerateClusterID(containerService.Properties))
		artifactsDir = path.Join("_output", artifactsDir)
	}

	// convert back the API object, and write it
	var b []byte
	var err error
	if !parametersOnly {
		apiloader := &api.Apiloader{
			Translator: w.Translator,
		}
		b, err = apiloader.SerializeContainerService(containerService, apiVersion)

		if err != nil {
			return err
		}

		if e := w.saveFile(artifactsDir, "apimodel.json", b); e != nil {
			return e
		}

		if e := w.saveFileString(artifactsDir, "azuredeploy.json", template); e != nil {
			return e
		}
	}

	if e := w.saveFileString(artifactsDir, "azuredeploy.parameters.json", parameters); e != nil {
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
				if e := w.saveFileString(directory, fmt.Sprintf("kubeconfig.%s.json", location), b); e != nil {
					return e
				}
			}

		}

		if e := w.saveFileString(artifactsDir, "ca.key", properties.CertificateProfile.CaPrivateKey); e != nil {
			return e
		}
		if e := w.saveFileString(artifactsDir, "ca.crt", properties.CertificateProfile.CaCertificate); e != nil {
			return e
		}
		if e := w.saveFileString(artifactsDir, "apiserver.key", properties.CertificateProfile.APIServerPrivateKey); e != nil {
			return e
		}
		if e := w.saveFileString(artifactsDir, "apiserver.crt", properties.CertificateProfile.APIServerCertificate); e != nil {
			return e
		}
		if e := w.saveFileString(artifactsDir, "client.key", properties.CertificateProfile.ClientPrivateKey); e != nil {
			return e
		}
		if e := w.saveFileString(artifactsDir, "client.crt", properties.CertificateProfile.ClientCertificate); e != nil {
			return e
		}
		if e := w.saveFileString(artifactsDir, "kubectlClient.key", properties.CertificateProfile.KubeConfigPrivateKey); e != nil {
			return e
		}
		if e := w.saveFileString(artifactsDir, "kubectlClient.crt", properties.CertificateProfile.KubeConfigCertificate); e != nil {
			return e
		}
	}

	return nil
}

func (w *ArtifactWriter) saveFileString(dir string, file string, data string) error {
	return w.saveFile(dir, file, []byte(data))
}

func (w *ArtifactWriter) saveFile(dir string, file string, data []byte) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if e := os.MkdirAll(dir, 0700); e != nil {
			return w.Translator.Errorf("error creating directory '%s': %s", dir, e.Error())
		}
	}

	path := path.Join(dir, file)
	if err := ioutil.WriteFile(path, []byte(data), 0600); err != nil {
		return err
	}

	fmt.Fprintf(os.Stderr, "wrote %s\n", path)

	return nil
}

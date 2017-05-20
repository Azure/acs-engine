package acsengine

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/Azure/acs-engine/pkg/api"
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

		if e := saveFileString(artifactsDir, "ca.key", properties.KubernetesCertificateProfile.GetKubernetesCAPrivateKey()); e != nil {
			return e
		}
		if e := saveFileString(artifactsDir, "ca.crt", properties.KubernetesCertificateProfile.CaCertificate); e != nil {
			return e
		}
		if e := saveFileString(artifactsDir, "apiserver.key", properties.KubernetesCertificateProfile.APIServerPrivateKey); e != nil {
			return e
		}
		if e := saveFileString(artifactsDir, "apiserver.crt", properties.KubernetesCertificateProfile.APIServerCertificate); e != nil {
			return e
		}
		if e := saveFileString(artifactsDir, "client.key", properties.KubernetesCertificateProfile.ClientPrivateKey); e != nil {
			return e
		}
		if e := saveFileString(artifactsDir, "client.crt", properties.KubernetesCertificateProfile.ClientCertificate); e != nil {
			return e
		}
		if e := saveFileString(artifactsDir, "kubectlClient.key", properties.KubernetesCertificateProfile.KubeConfigPrivateKey); e != nil {
			return e
		}
		if e := saveFileString(artifactsDir, "kubectlClient.crt", properties.KubernetesCertificateProfile.KubeConfigCertificate); e != nil {
			return e
		}
		// if e := saveFileString(artifactsDir, "ca.key", properties.SwarmModeCertificateProfile.GetSwarmModeCAPrivateKey()); e != nil {
		// 	return e
		// }
		// if e := saveFileString(artifactsDir, "ca.crt", properties.SwarmModeCertificateProfile.CaCertificate); e != nil {
		// 	return e
		// }
		// if e := saveFileString(artifactsDir, "server.key", properties.SwarmModeCertificateProfile.SwarmTLSServerPrivateKey); e != nil {
		// 	return e
		// }
		// if e := saveFileString(artifactsDir, "server.crt", properties.SwarmModeCertificateProfile.SwarmTLSServerCertificate); e != nil {
		// 	return e
		// }
		// if e := saveFileString(artifactsDir, "client.key", properties.SwarmModeCertificateProfile.SwarmTLSClientPrivateKey); e != nil {
		// 	return e
		// }
		// if e := saveFileString(artifactsDir, "client.crt", properties.SwarmModeCertificateProfile.SwarmTLSClientCertificate); e != nil {
		// 	return e
		// }
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

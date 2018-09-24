package acsengine

import (
	"fmt"
	"io/ioutil"
	"path"
	"path/filepath"
	"strconv"

	"github.com/Azure/acs-engine/pkg/api"
	"github.com/Azure/acs-engine/pkg/i18n"
	"github.com/pkg/errors"
)

// ArtifactWriter represents the object that writes artifacts
type ArtifactWriter struct {
	Translator *i18n.Translator
}

// WriteTLSArtifacts saves TLS certificates and keys to the server filesystem
func (w *ArtifactWriter) WriteTLSArtifacts(containerService *api.ContainerService, apiVersion, template, parameters, artifactsDir string, certsGenerated bool, parametersOnly bool) error {
	if len(artifactsDir) == 0 {
		artifactsDir = fmt.Sprintf("%s-%s", containerService.Properties.OrchestratorProfile.OrchestratorType, containerService.Properties.GetClusterID())
		artifactsDir = path.Join("_output", artifactsDir)
	}

	f := &FileSaver{
		Translator: w.Translator,
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

		if e := f.SaveFile(artifactsDir, "apimodel.json", b); e != nil {
			return e
		}

		if e := f.SaveFileString(artifactsDir, "azuredeploy.json", template); e != nil {
			return e
		}
	}

	if e := f.SaveFileString(artifactsDir, "azuredeploy.parameters.json", parameters); e != nil {
		return e
	}

	if !certsGenerated {
		return nil
	}

	properties := containerService.Properties
	if properties.OrchestratorProfile.IsKubernetes() {
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
			if e := f.SaveFileString(directory, fmt.Sprintf("kubeconfig.%s.json", location), b); e != nil {
				return e
			}
		}

		if e := f.SaveFileString(artifactsDir, "ca.key", properties.CertificateProfile.CaPrivateKey); e != nil {
			return e
		}
		if e := f.SaveFileString(artifactsDir, "ca.crt", properties.CertificateProfile.CaCertificate); e != nil {
			return e
		}
		if e := f.SaveFileString(artifactsDir, "apiserver.key", properties.CertificateProfile.APIServerPrivateKey); e != nil {
			return e
		}
		if e := f.SaveFileString(artifactsDir, "apiserver.crt", properties.CertificateProfile.APIServerCertificate); e != nil {
			return e
		}
		if e := f.SaveFileString(artifactsDir, "client.key", properties.CertificateProfile.ClientPrivateKey); e != nil {
			return e
		}
		if e := f.SaveFileString(artifactsDir, "client.crt", properties.CertificateProfile.ClientCertificate); e != nil {
			return e
		}
		if e := f.SaveFileString(artifactsDir, "kubectlClient.key", properties.CertificateProfile.KubeConfigPrivateKey); e != nil {
			return e
		}
		if e := f.SaveFileString(artifactsDir, "kubectlClient.crt", properties.CertificateProfile.KubeConfigCertificate); e != nil {
			return e
		}
		if e := f.SaveFileString(artifactsDir, "etcdserver.key", properties.CertificateProfile.EtcdServerPrivateKey); e != nil {
			return e
		}
		if e := f.SaveFileString(artifactsDir, "etcdserver.crt", properties.CertificateProfile.EtcdServerCertificate); e != nil {
			return e
		}
		if e := f.SaveFileString(artifactsDir, "etcdclient.key", properties.CertificateProfile.EtcdClientPrivateKey); e != nil {
			return e
		}
		if e := f.SaveFileString(artifactsDir, "etcdclient.crt", properties.CertificateProfile.EtcdClientCertificate); e != nil {
			return e
		}
		for i := 0; i < properties.MasterProfile.Count; i++ {
			if len(properties.CertificateProfile.EtcdPeerPrivateKeys) <= i || len(properties.CertificateProfile.EtcdPeerCertificates) <= i {
				return errors.New("missing etcd peer certificate/key pair")
			}
			k := "etcdpeer" + strconv.Itoa(i) + ".key"
			if e := f.SaveFileString(artifactsDir, k, properties.CertificateProfile.EtcdPeerPrivateKeys[i]); e != nil {
				return e
			}
			c := "etcdpeer" + strconv.Itoa(i) + ".crt"
			if e := f.SaveFileString(artifactsDir, c, properties.CertificateProfile.EtcdPeerCertificates[i]); e != nil {
				return e
			}
		}
	} else if properties.OrchestratorProfile.IsOpenShift() {
		masterTarballPath := filepath.Join(artifactsDir, "master.tar.gz")
		masterBundle := properties.OrchestratorProfile.OpenShiftConfig.ConfigBundles["master"]
		if err := ioutil.WriteFile(masterTarballPath, masterBundle, 0644); err != nil {
			return err
		}
		nodeTarballPath := filepath.Join(artifactsDir, "node.tar.gz")
		nodeBundle := properties.OrchestratorProfile.OpenShiftConfig.ConfigBundles["bootstrap"]
		return ioutil.WriteFile(nodeTarballPath, nodeBundle, 0644)
	}

	return nil
}

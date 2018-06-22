package acsengine

import (
	"os"
	"testing"

	"github.com/Azure/acs-engine/pkg/api"
	"github.com/Azure/acs-engine/pkg/helpers"
	"github.com/Azure/acs-engine/pkg/i18n"
	"github.com/satori/go.uuid"
)

func TestWriteTLSArtifacts(t *testing.T) {

	writer := &ArtifactWriter{
		Translator: &i18n.Translator{
			Locale: nil,
		},
	}
	dir := "_testoutputdir"
	defer os.Remove(dir)
	cs := CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2)
	err := writer.WriteTLSArtifacts(cs, "vlabs", "fake template", "fake parameters", dir, false, false)

	if err != nil {
		t.Fatalf("unexpected error trying to write TLS artifacts: %s", err.Error())
	}
}

// CreateMockContainerService returns a mock container service for testing purposes
func CreateMockContainerService(containerServiceName string, orchestratorVersion string, masterCount int, agentCount int) *api.ContainerService {
	cs := api.ContainerService{}
	cs.ID = uuid.NewV4().String()
	cs.Location = "eastus"
	cs.Name = containerServiceName

	cs.Properties = &api.Properties{}

	cs.Properties.MasterProfile = &api.MasterProfile{}
	cs.Properties.MasterProfile.Count = masterCount
	cs.Properties.MasterProfile.DNSPrefix = "testmaster"
	cs.Properties.MasterProfile.VMSize = "Standard_D2_v2"

	cs.Properties.AgentPoolProfiles = []*api.AgentPoolProfile{}
	agentPool := &api.AgentPoolProfile{}
	agentPool.Count = agentCount
	agentPool.Name = "agentpool1"
	agentPool.VMSize = "Standard_D2_v2"
	agentPool.OSType = "Linux"
	agentPool.AvailabilityProfile = "AvailabilitySet"
	agentPool.StorageProfile = "StorageAccount"

	cs.Properties.AgentPoolProfiles = append(cs.Properties.AgentPoolProfiles, agentPool)

	cs.Properties.LinuxProfile = &api.LinuxProfile{
		AdminUsername: "azureuser",
		SSH: struct {
			PublicKeys []api.PublicKey `json:"publicKeys"`
		}{},
	}

	cs.Properties.LinuxProfile.AdminUsername = "azureuser"
	cs.Properties.LinuxProfile.SSH.PublicKeys = append(
		cs.Properties.LinuxProfile.SSH.PublicKeys, api.PublicKey{KeyData: "test"})

	cs.Properties.ServicePrincipalProfile = &api.ServicePrincipalProfile{}
	cs.Properties.ServicePrincipalProfile.ClientID = "DEC923E3-1EF1-4745-9516-37906D56DEC4"
	cs.Properties.ServicePrincipalProfile.Secret = "DEC923E3-1EF1-4745-9516-37906D56DEC4"

	cs.Properties.OrchestratorProfile = &api.OrchestratorProfile{}
	cs.Properties.OrchestratorProfile.OrchestratorType = api.Kubernetes
	cs.Properties.OrchestratorProfile.OrchestratorVersion = orchestratorVersion
	cs.Properties.OrchestratorProfile.KubernetesConfig = &api.KubernetesConfig{
		EnableSecureKubelet: helpers.PointerToBool(api.DefaultSecureKubeletEnabled),
		EnableRbac:          helpers.PointerToBool(api.DefaultRBACEnabled),
		EtcdDiskSizeGB:      DefaultEtcdDiskSize,
		ServiceCIDR:         DefaultKubernetesServiceCIDR,
		DockerBridgeSubnet:  DefaultDockerBridgeSubnet,
		DNSServiceIP:        DefaultKubernetesDNSServiceIP,
		GCLowThreshold:      DefaultKubernetesGCLowThreshold,
		GCHighThreshold:     DefaultKubernetesGCHighThreshold,
		MaxPods:             DefaultKubernetesMaxPodsVNETIntegrated,
		ClusterSubnet:       DefaultKubernetesSubnet,
		ContainerRuntime:    DefaultContainerRuntime,
		NetworkPlugin:       DefaultNetworkPlugin,
		NetworkPolicy:       DefaultNetworkPolicy,
		EtcdVersion:         DefaultEtcdVersion,
		KubeletConfig:       make(map[string]string),
	}

	cs.Properties.CertificateProfile = &api.CertificateProfile{}
	cs.Properties.CertificateProfile.CaCertificate = "cacert"
	cs.Properties.CertificateProfile.KubeConfigCertificate = "kubeconfigcert"
	cs.Properties.CertificateProfile.KubeConfigPrivateKey = "kubeconfigkey"

	return &cs
}

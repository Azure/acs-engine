package acsengine

import (
	"github.com/Azure/acs-engine/pkg/api"
	"github.com/Azure/acs-engine/pkg/api/v20160330"
	"github.com/Azure/acs-engine/pkg/api/vlabs"
	"github.com/Azure/acs-engine/pkg/helpers"
	"github.com/Azure/acs-engine/pkg/i18n"
	uuid "github.com/satori/go.uuid"
)

// DCOSNodeType represents the type of DCOS Node
type DCOSNodeType string

// VlabsContainerService is the type we read and write from file
// needed because the json that is sent to ARM and acs-engine
// is different from the json that the ACS RP Api gets from ARM
type VlabsContainerService struct {
	api.TypeMeta
	*vlabs.ContainerService
}

// V20160330ContainerService is the type we read and write from file
// needed because the json that is sent to ARM and acs-engine
// is different from the json that the ACS RP Api gets from ARM
type V20160330ContainerService struct {
	api.TypeMeta
	*v20160330.ContainerService
}

//DockerSpecConfig is the configurations of docker
type DockerSpecConfig struct {
	DockerEngineRepo         string
	DockerComposeDownloadURL string
}

//DCOSSpecConfig is the configurations of DCOS
type DCOSSpecConfig struct {
	DCOS188BootstrapDownloadURL     string
	DCOS190BootstrapDownloadURL     string
	DCOS198BootstrapDownloadURL     string
	DCOS110BootstrapDownloadURL     string
	DCOS111BootstrapDownloadURL     string
	DCOSWindowsBootstrapDownloadURL string
	DcosRepositoryURL               string // For custom install, for example CI, need these three addributes
	DcosClusterPackageListID        string // the id of the package list file
	DcosProviderPackageID           string // the id of the dcos-provider-xxx package
}

//KubernetesSpecConfig is the kubernetes container images used.
type KubernetesSpecConfig struct {
	KubernetesImageBase              string
	TillerImageBase                  string
	ACIConnectorImageBase            string
	NVIDIAImageBase                  string
	AzureCNIImageBase                string
	EtcdDownloadURLBase              string
	KubeBinariesSASURLBase           string
	WindowsPackageSASURLBase         string
	WindowsTelemetryGUID             string
	CNIPluginsDownloadURL            string
	VnetCNILinuxPluginsDownloadURL   string
	VnetCNIWindowsPluginsDownloadURL string
}

//AzureEndpointConfig describes an Azure endpoint
type AzureEndpointConfig struct {
	ResourceManagerVMDNSSuffix string
}

//AzureOSImageConfig describes an Azure OS image
type AzureOSImageConfig struct {
	ImageOffer     string
	ImageSku       string
	ImagePublisher string
	ImageVersion   string
}

//AzureEnvironmentSpecConfig is the overall configuration differences in different cloud environments.
type AzureEnvironmentSpecConfig struct {
	DockerSpecConfig     DockerSpecConfig
	KubernetesSpecConfig KubernetesSpecConfig
	DCOSSpecConfig       DCOSSpecConfig
	EndpointConfig       AzureEndpointConfig
	OSImageConfig        map[api.Distro]AzureOSImageConfig
}

// Context represents the object that is passed to the package
type Context struct {
	Translator *i18n.Translator
}

// KeyVaultID represents a KeyVault instance on Azure
type KeyVaultID struct {
	ID string `json:"id"`
}

// KeyVaultRef represents a reference to KeyVault instance on Azure
type KeyVaultRef struct {
	KeyVault      KeyVaultID `json:"keyVault"`
	SecretName    string     `json:"secretName"`
	SecretVersion string     `json:"secretVersion,omitempty"`
}

type paramsMap map[string]interface{}

// CreateMockContainerService returns a mock container service for testing purposes
func CreateMockContainerService(containerServiceName, orchestratorVersion string, masterCount, agentCount int, certs bool) *api.ContainerService {
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
	if certs {
		cs.Properties.CertificateProfile.CaCertificate = "cacert"
		cs.Properties.CertificateProfile.CaPrivateKey = "cakey"
		cs.Properties.CertificateProfile.KubeConfigCertificate = "kubeconfigcert"
		cs.Properties.CertificateProfile.KubeConfigPrivateKey = "kubeconfigkey"
		cs.Properties.CertificateProfile.APIServerCertificate = "apiservercert"
		cs.Properties.CertificateProfile.APIServerPrivateKey = "apiserverkey"
		cs.Properties.CertificateProfile.ClientCertificate = "clientcert"
		cs.Properties.CertificateProfile.ClientPrivateKey = "clientkey"
		cs.Properties.CertificateProfile.EtcdServerCertificate = "etcdservercert"
		cs.Properties.CertificateProfile.EtcdServerPrivateKey = "etcdserverkey"
		cs.Properties.CertificateProfile.EtcdClientCertificate = "etcdclientcert"
		cs.Properties.CertificateProfile.EtcdClientPrivateKey = "etcdclientkey"
		cs.Properties.CertificateProfile.EtcdPeerCertificates = []string{"etcdpeercert1", "etcdpeercert2", "etcdpeercert3", "etcdpeercert4", "etcdpeercert5"}
		cs.Properties.CertificateProfile.EtcdPeerPrivateKeys = []string{"etcdpeerkey1", "etcdpeerkey2", "etcdpeerkey3", "etcdpeerkey4", "etcdpeerkey5"}

	}

	return &cs
}

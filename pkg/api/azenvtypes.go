package api

import "fmt"

//AzureEnvironmentSpecConfig is the overall configuration differences in different cloud environments.
type AzureEnvironmentSpecConfig struct {
	CloudName            string
	DockerSpecConfig     DockerSpecConfig
	KubernetesSpecConfig KubernetesSpecConfig
	DCOSSpecConfig       DCOSSpecConfig
	EndpointConfig       AzureEndpointConfig
	OSImageConfig        map[Distro]AzureOSImageConfig
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
	ContainerdDownloadURLBase        string
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

var (
	//DefaultKubernetesSpecConfig is the default Docker image source of Kubernetes
	DefaultKubernetesSpecConfig = KubernetesSpecConfig{
		KubernetesImageBase:              "k8s.gcr.io/",
		TillerImageBase:                  "gcr.io/kubernetes-helm/",
		ACIConnectorImageBase:            "microsoft/",
		NVIDIAImageBase:                  "nvidia/",
		AzureCNIImageBase:                "containernetworking/",
		EtcdDownloadURLBase:              "https://acs-mirror.azureedge.net/github-coreos",
		KubeBinariesSASURLBase:           "https://acs-mirror.azureedge.net/wink8s/",
		WindowsPackageSASURLBase:         "https://acs-mirror.azureedge.net/wink8s/",
		WindowsTelemetryGUID:             "fb801154-36b9-41bc-89c2-f4d4f05472b0",
		CNIPluginsDownloadURL:            "https://acs-mirror.azureedge.net/cni/cni-plugins-amd64-" + CNIPluginVer + ".tgz",
		VnetCNILinuxPluginsDownloadURL:   "https://acs-mirror.azureedge.net/cni/azure-vnet-cni-linux-amd64-" + AzureCniPluginVerLinux + ".tgz",
		VnetCNIWindowsPluginsDownloadURL: "https://acs-mirror.azureedge.net/cni/azure-vnet-cni-windows-amd64-" + AzureCniPluginVerWindows + ".zip",
		ContainerdDownloadURLBase:        "https://storage.googleapis.com/cri-containerd-release/",
	}

	//DefaultDCOSSpecConfig is the default DC/OS binary download URL.
	DefaultDCOSSpecConfig = DCOSSpecConfig{
		DCOS188BootstrapDownloadURL:     fmt.Sprintf(AzureEdgeDCOSBootstrapDownloadURL, "stable", "5df43052907c021eeb5de145419a3da1898c58a5"),
		DCOS190BootstrapDownloadURL:     fmt.Sprintf(AzureEdgeDCOSBootstrapDownloadURL, "stable", "58fd0833ce81b6244fc73bf65b5deb43217b0bd7"),
		DCOS198BootstrapDownloadURL:     fmt.Sprintf(AzureEdgeDCOSBootstrapDownloadURL, "stable/1.9.8", "f4ae0d20665fc68ee25282d6f78681b2773c6e10"),
		DCOS110BootstrapDownloadURL:     fmt.Sprintf(AzureEdgeDCOSBootstrapDownloadURL, "stable/1.10.0", "4d92536e7381176206e71ee15b5ffe454439920c"),
		DCOS111BootstrapDownloadURL:     fmt.Sprintf(AzureEdgeDCOSBootstrapDownloadURL, "stable/1.11.0", "a0654657903fb68dff60f6e522a7f241c1bfbf0f"),
		DCOSWindowsBootstrapDownloadURL: "http://dcos-win.westus.cloudapp.azure.com/dcos-windows/stable/",
		DcosRepositoryURL:               "https://dcosio.azureedge.net/dcos/stable/1.11.0",
		DcosClusterPackageListID:        "248a66388bba1adbcb14a52fd3b7b424ab06fa76",
	}

	//DefaultDockerSpecConfig is the default Docker engine repo.
	DefaultDockerSpecConfig = DockerSpecConfig{
		DockerEngineRepo:         "https://aptdocker.azureedge.net/repo",
		DockerComposeDownloadURL: "https://github.com/docker/compose/releases/download",
	}

	//DefaultUbuntuImageConfig is the default Linux distribution.
	DefaultUbuntuImageConfig = AzureOSImageConfig{
		ImageOffer:     "UbuntuServer",
		ImageSku:       "16.04-LTS",
		ImagePublisher: "Canonical",
		ImageVersion:   "latest",
	}

	//SovereignCloudsUbuntuImageConfig is the Linux distribution for Azure Sovereign Clouds.
	SovereignCloudsUbuntuImageConfig = AzureOSImageConfig{
		ImageOffer:     "UbuntuServer",
		ImageSku:       "16.04-LTS",
		ImagePublisher: "Canonical",
		ImageVersion:   "latest",
	}

	//GermanCloudUbuntuImageConfig is the Linux distribution for Azure Sovereign Clouds.
	GermanCloudUbuntuImageConfig = AzureOSImageConfig{
		ImageOffer:     "UbuntuServer",
		ImageSku:       "16.04-LTS",
		ImagePublisher: "Canonical",
		ImageVersion:   "16.04.201801050",
	}

	//DefaultRHELOSImageConfig is the RHEL Linux distribution.
	DefaultRHELOSImageConfig = AzureOSImageConfig{
		ImageOffer:     "RHEL",
		ImageSku:       "7.3",
		ImagePublisher: "RedHat",
		ImageVersion:   "latest",
	}

	//DefaultCoreOSImageConfig is the CoreOS Linux distribution.
	DefaultCoreOSImageConfig = AzureOSImageConfig{
		ImageOffer:     "CoreOS",
		ImageSku:       "Stable",
		ImagePublisher: "CoreOS",
		ImageVersion:   "latest",
	}

	// DefaultAKSOSImageConfig is the AKS image based on Ubuntu 16.04.
	DefaultAKSOSImageConfig = AzureOSImageConfig{
		ImageOffer:     "aks",
		ImageSku:       "aksbase",
		ImagePublisher: "microsoft-aks",
		ImageVersion:   "0.19.0",
	}

	//DefaultOpenShift39RHELImageConfig is the OpenShift on RHEL distribution.
	DefaultOpenShift39RHELImageConfig = AzureOSImageConfig{
		ImageOffer:     "acsengine-preview",
		ImageSku:       "rhel74",
		ImagePublisher: "redhat",
		ImageVersion:   "latest",
	}

	//DefaultOpenShift39CentOSImageConfig is the OpenShift on CentOS distribution.
	DefaultOpenShift39CentOSImageConfig = AzureOSImageConfig{
		ImageOffer:     "origin-acsengine-preview",
		ImageSku:       "centos7",
		ImagePublisher: "redhat",
		ImageVersion:   "latest",
	}

	//AzureCloudSpec is the default configurations for global azure.
	AzureCloudSpec = AzureEnvironmentSpecConfig{
		CloudName: azurePublicCloud,
		//DockerSpecConfig specify the docker engine download repo
		DockerSpecConfig: DefaultDockerSpecConfig,
		//KubernetesSpecConfig is the default kubernetes container image url.
		KubernetesSpecConfig: DefaultKubernetesSpecConfig,
		DCOSSpecConfig:       DefaultDCOSSpecConfig,

		EndpointConfig: AzureEndpointConfig{
			ResourceManagerVMDNSSuffix: "cloudapp.azure.com",
		},

		OSImageConfig: map[Distro]AzureOSImageConfig{
			Ubuntu: DefaultUbuntuImageConfig,
			RHEL:   DefaultRHELOSImageConfig,
			CoreOS: DefaultCoreOSImageConfig,
			AKS:    DefaultAKSOSImageConfig,
			// Image config supported for OpenShift
			OpenShift39RHEL: DefaultOpenShift39RHELImageConfig,
			OpenShiftCentOS: DefaultOpenShift39CentOSImageConfig,
		},
	}

	//AzureGermanCloudSpec is the German cloud config.
	AzureGermanCloudSpec = AzureEnvironmentSpecConfig{
		CloudName:            azureGermanCloud,
		DockerSpecConfig:     DefaultDockerSpecConfig,
		KubernetesSpecConfig: DefaultKubernetesSpecConfig,
		DCOSSpecConfig:       DefaultDCOSSpecConfig,
		EndpointConfig: AzureEndpointConfig{
			ResourceManagerVMDNSSuffix: "cloudapp.microsoftazure.de",
		},
		OSImageConfig: map[Distro]AzureOSImageConfig{
			Ubuntu: GermanCloudUbuntuImageConfig,
			RHEL:   DefaultRHELOSImageConfig,
			CoreOS: DefaultCoreOSImageConfig,
			AKS:    GermanCloudUbuntuImageConfig,
		},
	}

	//AzureUSGovernmentCloud is the US government config.
	AzureUSGovernmentCloud = AzureEnvironmentSpecConfig{
		CloudName:            azureUSGovernmentCloud,
		DockerSpecConfig:     DefaultDockerSpecConfig,
		KubernetesSpecConfig: DefaultKubernetesSpecConfig,
		DCOSSpecConfig:       DefaultDCOSSpecConfig,
		EndpointConfig: AzureEndpointConfig{
			ResourceManagerVMDNSSuffix: "cloudapp.usgovcloudapi.net",
		},
		OSImageConfig: map[Distro]AzureOSImageConfig{
			Ubuntu: SovereignCloudsUbuntuImageConfig,
			RHEL:   DefaultRHELOSImageConfig,
			CoreOS: DefaultCoreOSImageConfig,
			AKS:    SovereignCloudsUbuntuImageConfig,
		},
	}

	//AzureChinaCloudSpec is the configurations for Azure China (Mooncake)
	AzureChinaCloudSpec = AzureEnvironmentSpecConfig{
		CloudName: azureChinaCloud,
		//DockerSpecConfig specify the docker engine download repo
		DockerSpecConfig: DockerSpecConfig{
			DockerEngineRepo:         "https://mirror.azure.cn/docker-engine/apt/repo/",
			DockerComposeDownloadURL: "https://mirror.azure.cn/docker-toolbox/linux/compose",
		},
		//KubernetesSpecConfig - Due to Chinese firewall issue, the default containers from google is blocked, use the Chinese local mirror instead
		KubernetesSpecConfig: KubernetesSpecConfig{
			KubernetesImageBase:              "gcr.akscn.io/google_containers/",
			TillerImageBase:                  "gcr.akscn.io/kubernetes-helm/",
			ACIConnectorImageBase:            "dockerhub.akscn.io/microsoft/",
			NVIDIAImageBase:                  "dockerhub.akscn.io/nvidia/",
			AzureCNIImageBase:                "dockerhub.akscn.io/containernetworking/",
			EtcdDownloadURLBase:              DefaultKubernetesSpecConfig.EtcdDownloadURLBase,
			KubeBinariesSASURLBase:           DefaultKubernetesSpecConfig.KubeBinariesSASURLBase,
			WindowsPackageSASURLBase:         DefaultKubernetesSpecConfig.WindowsPackageSASURLBase,
			WindowsTelemetryGUID:             DefaultKubernetesSpecConfig.WindowsTelemetryGUID,
			CNIPluginsDownloadURL:            DefaultKubernetesSpecConfig.CNIPluginsDownloadURL,
			VnetCNILinuxPluginsDownloadURL:   DefaultKubernetesSpecConfig.VnetCNILinuxPluginsDownloadURL,
			VnetCNIWindowsPluginsDownloadURL: DefaultKubernetesSpecConfig.VnetCNIWindowsPluginsDownloadURL,
			ContainerdDownloadURLBase:        "https://mirror.azure.cn/kubernetes/containerd/",
		},
		DCOSSpecConfig: DCOSSpecConfig{
			DCOS188BootstrapDownloadURL:     fmt.Sprintf(AzureChinaCloudDCOSBootstrapDownloadURL, "5df43052907c021eeb5de145419a3da1898c58a5"),
			DCOSWindowsBootstrapDownloadURL: "https://dcosdevstorage.blob.core.windows.net/dcos-windows",
			DCOS190BootstrapDownloadURL:     fmt.Sprintf(AzureChinaCloudDCOSBootstrapDownloadURL, "58fd0833ce81b6244fc73bf65b5deb43217b0bd7"),
			DCOS198BootstrapDownloadURL:     fmt.Sprintf(AzureChinaCloudDCOSBootstrapDownloadURL, "f4ae0d20665fc68ee25282d6f78681b2773c6e10"),
		},

		EndpointConfig: AzureEndpointConfig{
			ResourceManagerVMDNSSuffix: "cloudapp.chinacloudapi.cn",
		},
		OSImageConfig: map[Distro]AzureOSImageConfig{
			Ubuntu: SovereignCloudsUbuntuImageConfig,
			RHEL:   DefaultRHELOSImageConfig,
			CoreOS: DefaultCoreOSImageConfig,
			AKS:    SovereignCloudsUbuntuImageConfig,
		},
	}

	// AzureCloudSpecEnvMap is the environment configuration map for all the Azure cloid environments.
	AzureCloudSpecEnvMap = map[string]AzureEnvironmentSpecConfig{
		azureChinaCloud:        AzureChinaCloudSpec,
		azureGermanCloud:       AzureGermanCloudSpec,
		azureUSGovernmentCloud: AzureUSGovernmentCloud,
		azurePublicCloud:       AzureCloudSpec,
	}
)

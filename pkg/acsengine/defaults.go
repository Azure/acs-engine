package acsengine

import (
	"fmt"
	"net"

	"github.com/Azure/acs-engine/pkg/api"
	"github.com/Azure/acs-engine/pkg/api/common"
	"github.com/Masterminds/semver"
)

var (
	//DefaultKubernetesSpecConfig is the default Docker image source of Kubernetes
	DefaultKubernetesSpecConfig = KubernetesSpecConfig{
		KubernetesImageBase:              "gcrio.azureedge.net/google_containers/",
		TillerImageBase:                  "gcrio.azureedge.net/kubernetes-helm/",
		KubeBinariesSASURLBase:           "https://acs-mirror.azureedge.net/wink8s/",
		WindowsTelemetryGUID:             "fb801154-36b9-41bc-89c2-f4d4f05472b0",
		CNIPluginsDownloadURL:            "https://acs-mirror.azureedge.net/cni/cni-plugins-amd64-latest.tgz",
		VnetCNILinuxPluginsDownloadURL:   "https://acs-mirror.azureedge.net/cni/azure-vnet-cni-linux-amd64-latest.tgz",
		VnetCNIWindowsPluginsDownloadURL: "https://acs-mirror.azureedge.net/cni/azure-vnet-cni-windows-amd64-latest.zip",
		CalicoConfigDownloadURL:          "https://raw.githubusercontent.com/projectcalico/calico/a4ebfbad55ab1b7f10fdf3b39585471f8012e898/v2.0/getting-started/kubernetes/installation/hosted/k8s-backend-addon-manager",
	}

	//DefaultDCOSSpecConfig is the default DC/OS binary download URL.
	DefaultDCOSSpecConfig = DCOSSpecConfig{
		DCOS188BootstrapDownloadURL:     fmt.Sprintf(AzureEdgeDCOSBootstrapDownloadURL, "stable", "5df43052907c021eeb5de145419a3da1898c58a5"),
		DCOS190BootstrapDownloadURL:     fmt.Sprintf(AzureEdgeDCOSBootstrapDownloadURL, "stable", "58fd0833ce81b6244fc73bf65b5deb43217b0bd7"),
		DCOS110BootstrapDownloadURL:     fmt.Sprintf(AzureEdgeDCOSBootstrapDownloadURL, "stable", "e38ab2aa282077c8eb7bf103c6fff7b0f08db1a4"),
		DCOSWindowsBootstrapDownloadURL: "http://dcos-win.westus.cloudapp.azure.com/dcos-windows/stable/",
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
		ImageVersion:   "16.04.201706191",
	}

	//DefaultRHELOSImageConfig is the RHEL Linux distribution.
	DefaultRHELOSImageConfig = AzureOSImageConfig{
		ImageOffer:     "RHEL",
		ImageSku:       "7.3",
		ImagePublisher: "RedHat",
		ImageVersion:   "latest",
	}

	//AzureCloudSpec is the default configurations for global azure.
	AzureCloudSpec = AzureEnvironmentSpecConfig{
		//DockerSpecConfig specify the docker engine download repo
		DockerSpecConfig: DefaultDockerSpecConfig,
		//KubernetesSpecConfig is the default kubernetes container image url.
		KubernetesSpecConfig: DefaultKubernetesSpecConfig,
		DCOSSpecConfig:       DefaultDCOSSpecConfig,

		EndpointConfig: AzureEndpointConfig{
			ResourceManagerVMDNSSuffix: "cloudapp.azure.com",
		},

		OSImageConfig: map[api.Distro]AzureOSImageConfig{
			api.Ubuntu: DefaultUbuntuImageConfig,
			api.RHEL:   DefaultRHELOSImageConfig,
		},
	}

	//AzureGermanCloudSpec is the German cloud config.
	AzureGermanCloudSpec = AzureEnvironmentSpecConfig{
		DockerSpecConfig:     DefaultDockerSpecConfig,
		KubernetesSpecConfig: DefaultKubernetesSpecConfig,
		DCOSSpecConfig:       DefaultDCOSSpecConfig,
		EndpointConfig: AzureEndpointConfig{
			ResourceManagerVMDNSSuffix: "cloudapp.microsoftazure.de",
		},
		OSImageConfig: map[api.Distro]AzureOSImageConfig{
			api.Ubuntu: {
				ImageOffer:     "UbuntuServer",
				ImageSku:       "16.04-LTS",
				ImagePublisher: "Canonical",
				ImageVersion:   "16.04.201701130",
			},
			api.RHEL: DefaultRHELOSImageConfig,
		},
	}

	//AzureUSGovernmentCloud is the US government config.
	AzureUSGovernmentCloud = AzureEnvironmentSpecConfig{
		DockerSpecConfig:     DefaultDockerSpecConfig,
		KubernetesSpecConfig: DefaultKubernetesSpecConfig,
		DCOSSpecConfig:       DefaultDCOSSpecConfig,
		EndpointConfig: AzureEndpointConfig{
			ResourceManagerVMDNSSuffix: "cloudapp.windowsazure.us",
		},
		OSImageConfig: map[api.Distro]AzureOSImageConfig{
			api.Ubuntu: DefaultUbuntuImageConfig,
			api.RHEL:   DefaultRHELOSImageConfig,
		},
	}

	//AzureChinaCloudSpec is the configurations for Azure China (Mooncake)
	AzureChinaCloudSpec = AzureEnvironmentSpecConfig{
		//DockerSpecConfig specify the docker engine download repo
		DockerSpecConfig: DockerSpecConfig{
			DockerEngineRepo:         "https://mirror.azure.cn/docker-engine/apt/repo/",
			DockerComposeDownloadURL: "https://mirror.azure.cn/docker-toolbox/linux/compose",
		},
		//KubernetesSpecConfig - Due to Chinese firewall issue, the default containers from google is blocked, use the Chinese local mirror instead
		KubernetesSpecConfig: KubernetesSpecConfig{
			KubernetesImageBase:              "crproxy.trafficmanager.net:6000/google_containers/",
			TillerImageBase:                  "mirror.azure.cn:5000/kubernetes-helm/",
			CNIPluginsDownloadURL:            "https://acsengine.blob.core.chinacloudapi.cn/cni/cni-plugins-amd64-latest.tgz",
			VnetCNILinuxPluginsDownloadURL:   "https://acsengine.blob.core.chinacloudapi.cn/cni/azure-vnet-cni-linux-amd64-latest.tgz",
			VnetCNIWindowsPluginsDownloadURL: "https://acsengine.blob.core.chinacloudapi.cn/cni/azure-vnet-cni-windows-amd64-latest.zip",
			CalicoConfigDownloadURL:          "https://acsengine.blob.core.chinacloudapi.cn/cni",
		},
		DCOSSpecConfig: DCOSSpecConfig{
			DCOS188BootstrapDownloadURL:     fmt.Sprintf(AzureChinaCloudDCOSBootstrapDownloadURL, "5df43052907c021eeb5de145419a3da1898c58a5"),
			DCOSWindowsBootstrapDownloadURL: "https://dcosdevstorage.blob.core.windows.net/dcos-windows",
			DCOS190BootstrapDownloadURL:     fmt.Sprintf(AzureChinaCloudDCOSBootstrapDownloadURL, "58fd0833ce81b6244fc73bf65b5deb43217b0bd7"),
		},

		EndpointConfig: AzureEndpointConfig{
			ResourceManagerVMDNSSuffix: "cloudapp.chinacloudapi.cn",
		},
		OSImageConfig: map[api.Distro]AzureOSImageConfig{
			api.Ubuntu: DefaultUbuntuImageConfig,
			api.RHEL:   DefaultRHELOSImageConfig,
		},
	}
)

// SetPropertiesDefaults for the container Properties, returns true if certs are generated
func SetPropertiesDefaults(cs *api.ContainerService) (bool, error) {
	properties := cs.Properties

	setOrchestratorDefaults(cs)

	setMasterNetworkDefaults(properties)

	setHostedMasterNetworkDefaults(properties)

	setAgentNetworkDefaults(properties)

	setStorageDefaults(properties)
	setExtensionDefaults(properties)

	certsGenerated, e := setDefaultCerts(properties)
	if e != nil {
		return false, e
	}
	return certsGenerated, nil
}

// setOrchestratorDefaults for orchestrators
func setOrchestratorDefaults(cs *api.ContainerService) {
	location := cs.Location
	a := cs.Properties

	cloudSpecConfig := GetCloudSpecConfig(location)
	if a.OrchestratorProfile == nil {
		return
	}
	if a.OrchestratorProfile.OrchestratorType == api.Kubernetes {
		k8sRelease := a.OrchestratorProfile.OrchestratorRelease
		if a.OrchestratorProfile.KubernetesConfig == nil {
			a.OrchestratorProfile.KubernetesConfig = &api.KubernetesConfig{}
		}
		if a.OrchestratorProfile.KubernetesConfig.KubernetesImageBase == "" {
			a.OrchestratorProfile.KubernetesConfig.KubernetesImageBase = cloudSpecConfig.KubernetesSpecConfig.KubernetesImageBase
		}
		if a.OrchestratorProfile.KubernetesConfig.NetworkPolicy == "" {
			a.OrchestratorProfile.KubernetesConfig.NetworkPolicy = DefaultNetworkPolicy
		}
		if a.OrchestratorProfile.KubernetesConfig.ClusterSubnet == "" {
			if a.OrchestratorProfile.IsVNETIntegrated() {
				// When VNET integration is enabled, all masters, agents and pods share the same large subnet.
				a.OrchestratorProfile.KubernetesConfig.ClusterSubnet = DefaultKubernetesSubnet
			} else {
				a.OrchestratorProfile.KubernetesConfig.ClusterSubnet = DefaultKubernetesClusterSubnet
			}
		}
		if a.OrchestratorProfile.KubernetesConfig.MaxPods == 0 {
			if a.OrchestratorProfile.IsVNETIntegrated() {
				a.OrchestratorProfile.KubernetesConfig.MaxPods = DefaultKubernetesMaxPodsVNETIntegrated
			} else {
				a.OrchestratorProfile.KubernetesConfig.MaxPods = DefaultKubernetesMaxPods
			}
		}
		if a.OrchestratorProfile.KubernetesConfig.GCHighThreshold == 0 {
			a.OrchestratorProfile.KubernetesConfig.GCHighThreshold = DefaultKubernetesGCHighThreshold
		}
		if a.OrchestratorProfile.KubernetesConfig.GCLowThreshold == 0 {
			a.OrchestratorProfile.KubernetesConfig.GCLowThreshold = DefaultKubernetesGCLowThreshold
		}
		if a.OrchestratorProfile.KubernetesConfig.DNSServiceIP == "" {
			a.OrchestratorProfile.KubernetesConfig.DNSServiceIP = DefaultKubernetesDNSServiceIP
		}
		if a.OrchestratorProfile.KubernetesConfig.DockerBridgeSubnet == "" {
			a.OrchestratorProfile.KubernetesConfig.DockerBridgeSubnet = DefaultDockerBridgeSubnet
		}
		if a.OrchestratorProfile.KubernetesConfig.ServiceCIDR == "" {
			a.OrchestratorProfile.KubernetesConfig.ServiceCIDR = DefaultKubernetesServiceCIDR
		}
		if a.OrchestratorProfile.KubernetesConfig.NodeStatusUpdateFrequency == "" {
			a.OrchestratorProfile.KubernetesConfig.NodeStatusUpdateFrequency = KubeConfigs[k8sRelease]["nodestatusfreq"]
		}
		if a.OrchestratorProfile.KubernetesConfig.CtrlMgrNodeMonitorGracePeriod == "" {
			a.OrchestratorProfile.KubernetesConfig.CtrlMgrNodeMonitorGracePeriod = KubeConfigs[k8sRelease]["nodegraceperiod"]
		}
		if a.OrchestratorProfile.KubernetesConfig.CtrlMgrPodEvictionTimeout == "" {
			a.OrchestratorProfile.KubernetesConfig.CtrlMgrPodEvictionTimeout = KubeConfigs[k8sRelease]["podeviction"]
		}
		if a.OrchestratorProfile.KubernetesConfig.CtrlMgrRouteReconciliationPeriod == "" {
			a.OrchestratorProfile.KubernetesConfig.CtrlMgrRouteReconciliationPeriod = KubeConfigs[k8sRelease]["routeperiod"]
		}
		// Enforce sane cloudprovider backoff defaults, if CloudProviderBackoff is true in KubernetesConfig
		if a.OrchestratorProfile.KubernetesConfig.CloudProviderBackoff == true {
			if a.OrchestratorProfile.KubernetesConfig.CloudProviderBackoffDuration == 0 {
				a.OrchestratorProfile.KubernetesConfig.CloudProviderBackoffDuration = DefaultKubernetesCloudProviderBackoffDuration
			}
			if a.OrchestratorProfile.KubernetesConfig.CloudProviderBackoffExponent == 0 {
				a.OrchestratorProfile.KubernetesConfig.CloudProviderBackoffExponent = DefaultKubernetesCloudProviderBackoffExponent
			}
			if a.OrchestratorProfile.KubernetesConfig.CloudProviderBackoffJitter == 0 {
				a.OrchestratorProfile.KubernetesConfig.CloudProviderBackoffJitter = DefaultKubernetesCloudProviderBackoffJitter
			}
			if a.OrchestratorProfile.KubernetesConfig.CloudProviderBackoffRetries == 0 {
				a.OrchestratorProfile.KubernetesConfig.CloudProviderBackoffRetries = DefaultKubernetesCloudProviderBackoffRetries
			}
		}
		k8sVersion, _ := semver.NewVersion(api.KubernetesReleaseToVersion[k8sRelease])
		constraint, _ := semver.NewConstraint(">= 1.6.6")
		// Enforce sane cloudprovider rate limit defaults, if CloudProviderRateLimit is true in KubernetesConfig
		// For k8s version greater or equal to 1.6.6, we will set the default CloudProviderRate* settings
		if a.OrchestratorProfile.KubernetesConfig.CloudProviderRateLimit == true && constraint.Check(k8sVersion) {
			if a.OrchestratorProfile.KubernetesConfig.CloudProviderRateLimitQPS == 0 {
				a.OrchestratorProfile.KubernetesConfig.CloudProviderRateLimitQPS = DefaultKubernetesCloudProviderRateLimitQPS
			}
			if a.OrchestratorProfile.KubernetesConfig.CloudProviderRateLimitBucket == 0 {
				a.OrchestratorProfile.KubernetesConfig.CloudProviderRateLimitBucket = DefaultKubernetesCloudProviderRateLimitBucket
			}
		}
	} else if a.OrchestratorProfile.OrchestratorType == api.DCOS {
		if a.OrchestratorProfile.DcosConfig == nil {
			a.OrchestratorProfile.DcosConfig = &api.DcosConfig{}
		}
		if a.OrchestratorProfile.DcosConfig.DcosWindowsBootstrapURL == "" {
			a.OrchestratorProfile.DcosConfig.DcosWindowsBootstrapURL = DefaultDCOSSpecConfig.DCOSWindowsBootstrapDownloadURL
		}
	}
}

func setExtensionDefaults(a *api.Properties) {
	if a.ExtensionProfiles == nil {
		return
	}
	for _, extension := range a.ExtensionProfiles {
		if extension.RootURL == "" {
			extension.RootURL = DefaultExtensionsRootURL
		}
	}
}

// SetHostedMasterNetworkDefaults for hosted masters
func setHostedMasterNetworkDefaults(a *api.Properties) {
	if a.HostedMasterProfile == nil {
		return
	}
	a.HostedMasterProfile.Subnet = DefaultKubernetesMasterSubnet
}

// SetMasterNetworkDefaults for masters
func setMasterNetworkDefaults(a *api.Properties) {
	if a.MasterProfile == nil {
		return
	}

	// Set default Distro to Ubuntu
	if a.MasterProfile.Distro == "" {
		a.MasterProfile.Distro = api.Ubuntu
	}

	if !a.MasterProfile.IsCustomVNET() {
		if a.OrchestratorProfile.OrchestratorType == api.Kubernetes {
			if a.OrchestratorProfile.IsVNETIntegrated() {
				// When VNET integration is enabled, all masters, agents and pods share the same large subnet.
				a.MasterProfile.Subnet = a.OrchestratorProfile.KubernetesConfig.ClusterSubnet
				a.MasterProfile.FirstConsecutiveStaticIP = getFirstConsecutiveStaticIPAddress(a.MasterProfile.Subnet)
			} else {
				a.MasterProfile.Subnet = DefaultKubernetesMasterSubnet
				a.MasterProfile.FirstConsecutiveStaticIP = DefaultFirstConsecutiveKubernetesStaticIP
			}
		} else if a.HasWindows() {
			a.MasterProfile.Subnet = DefaultSwarmWindowsMasterSubnet
			a.MasterProfile.FirstConsecutiveStaticIP = DefaultSwarmWindowsFirstConsecutiveStaticIP
		} else {
			a.MasterProfile.Subnet = DefaultMasterSubnet
			a.MasterProfile.FirstConsecutiveStaticIP = DefaultFirstConsecutiveStaticIP
		}
	}

	// Set the default number of IP addresses allocated for masters.
	if a.MasterProfile.IPAddressCount == 0 {
		// Allocate one IP address for the node.
		a.MasterProfile.IPAddressCount = 1

		// Allocate IP addresses for pods if VNET integration is enabled.
		if a.OrchestratorProfile.IsVNETIntegrated() {
			if a.OrchestratorProfile.OrchestratorType == api.Kubernetes {
				a.MasterProfile.IPAddressCount += a.OrchestratorProfile.KubernetesConfig.MaxPods
			}
		}
	}

	if a.MasterProfile.HTTPSourceAddressPrefix == "" {
		a.MasterProfile.HTTPSourceAddressPrefix = "*"
	}
}

// SetAgentNetworkDefaults for agents
func setAgentNetworkDefaults(a *api.Properties) {
	// configure the subnets if not in custom VNET
	if a.MasterProfile != nil && !a.MasterProfile.IsCustomVNET() {
		subnetCounter := 0
		for _, profile := range a.AgentPoolProfiles {
			if a.OrchestratorProfile.OrchestratorType == api.Kubernetes {
				profile.Subnet = a.MasterProfile.Subnet
			} else {
				profile.Subnet = fmt.Sprintf(DefaultAgentSubnetTemplate, subnetCounter)
			}

			subnetCounter++
		}
	}

	for _, profile := range a.AgentPoolProfiles {
		// set default OSType to Linux
		if profile.OSType == "" {
			profile.OSType = api.Linux
		}
		// set default Distro to Ubuntu
		if profile.Distro == "" {
			profile.Distro = api.Ubuntu
		}

		// Set the default number of IP addresses allocated for agents.
		if profile.IPAddressCount == 0 {
			// Allocate one IP address for the node.
			profile.IPAddressCount = 1

			// Allocate IP addresses for pods if VNET integration is enabled.
			if a.OrchestratorProfile.IsVNETIntegrated() {
				if a.OrchestratorProfile.OrchestratorType == api.Kubernetes {
					profile.IPAddressCount += a.OrchestratorProfile.KubernetesConfig.MaxPods
				}
			}
		}
	}
}

// setStorageDefaults for agents
func setStorageDefaults(a *api.Properties) {
	if a.MasterProfile != nil && len(a.MasterProfile.StorageProfile) == 0 {
		a.MasterProfile.StorageProfile = api.StorageAccount
	}
	for _, profile := range a.AgentPoolProfiles {
		if len(profile.StorageProfile) == 0 {
			profile.StorageProfile = api.StorageAccount
		}
		if len(profile.AvailabilityProfile) == 0 {
			profile.AvailabilityProfile = api.VirtualMachineScaleSets
		}
	}
}

func setDefaultCerts(a *api.Properties) (bool, error) {
	if !certGenerationRequired(a) {
		return false, nil
	}

	masterExtraFQDNs := FormatAzureProdFQDNs(a.MasterProfile.DNSPrefix)
	firstMasterIP := net.ParseIP(a.MasterProfile.FirstConsecutiveStaticIP).To4()

	if firstMasterIP == nil {
		return false, fmt.Errorf("MasterProfile.FirstConsecutiveStaticIP '%s' is an invalid IP address", a.MasterProfile.FirstConsecutiveStaticIP)
	}

	ips := []net.IP{firstMasterIP}

	// Add the Internal Loadbalancer IP which is always at at a known offset from the firstMasterIP
	ips = append(ips, net.IP{firstMasterIP[0], firstMasterIP[1], firstMasterIP[2], firstMasterIP[3] + byte(DefaultInternalLbStaticIPOffset)})

	// Include the Internal load balancer as well
	for i := 1; i < a.MasterProfile.Count; i++ {
		ip := net.IP{firstMasterIP[0], firstMasterIP[1], firstMasterIP[2], firstMasterIP[3] + byte(i)}
		ips = append(ips, ip)
	}

	if a.CertificateProfile == nil {
		a.CertificateProfile = &api.CertificateProfile{}
	}

	// use the specified Certificate Authority pair, or generate a new pair
	var caPair *PkiKeyCertPair
	if len(a.CertificateProfile.CaCertificate) != 0 && len(a.CertificateProfile.CaPrivateKey) != 0 {
		caPair = &PkiKeyCertPair{CertificatePem: a.CertificateProfile.CaCertificate, PrivateKeyPem: a.CertificateProfile.CaPrivateKey}
	} else {
		caCertificate, caPrivateKey, err := createCertificate("ca", nil, nil, false, nil, nil, nil)
		if err != nil {
			return false, err
		}
		caPair = &PkiKeyCertPair{CertificatePem: string(certificateToPem(caCertificate.Raw)), PrivateKeyPem: string(privateKeyToPem(caPrivateKey))}
		a.CertificateProfile.CaCertificate = caPair.CertificatePem
		a.CertificateProfile.CaPrivateKey = caPair.PrivateKeyPem
	}

	cidrFirstIP, err := common.CidrStringFirstIP(a.OrchestratorProfile.KubernetesConfig.ServiceCIDR)
	if err != nil {
		return false, err
	}
	ips = append(ips, cidrFirstIP)

	apiServerPair, clientPair, kubeConfigPair, err := CreatePki(masterExtraFQDNs, ips, DefaultKubernetesClusterDomain, caPair)
	if err != nil {
		return false, err
	}

	a.CertificateProfile.APIServerCertificate = apiServerPair.CertificatePem
	a.CertificateProfile.APIServerPrivateKey = apiServerPair.PrivateKeyPem
	a.CertificateProfile.ClientCertificate = clientPair.CertificatePem
	a.CertificateProfile.ClientPrivateKey = clientPair.PrivateKeyPem
	a.CertificateProfile.KubeConfigCertificate = kubeConfigPair.CertificatePem
	a.CertificateProfile.KubeConfigPrivateKey = kubeConfigPair.PrivateKeyPem

	return true, nil
}

func certGenerationRequired(a *api.Properties) bool {
	if certAlreadyPresent(a.CertificateProfile) {
		return false
	}
	if a.MasterProfile == nil {
		return false
	}

	switch a.OrchestratorProfile.OrchestratorType {
	case api.Kubernetes:
		return true
	default:
		return false
	}
}

// certAlreadyPresent determines if the passed in CertificateProfile includes certificate data
// TODO actually verify valid/useable certificate data
func certAlreadyPresent(c *api.CertificateProfile) bool {
	if c != nil {
		switch {
		case len(c.APIServerCertificate) > 0:
			return true
		case len(c.APIServerPrivateKey) > 0:
			return true
		case len(c.ClientCertificate) > 0:
			return true
		case len(c.ClientPrivateKey) > 0:
			return true
		default:
			return false
		}
	}
	return false
}

// getFirstConsecutiveStaticIPAddress returns the first static IP address of the given subnet.
func getFirstConsecutiveStaticIPAddress(subnetStr string) string {
	_, subnet, err := net.ParseCIDR(subnetStr)
	if err != nil {
		return DefaultFirstConsecutiveKubernetesStaticIP
	}

	// Find the first and last octet of the host bits.
	ones, bits := subnet.Mask.Size()
	firstOctet := ones / 8
	lastOctet := bits/8 - 1

	// Set the remaining host bits in the first octet.
	subnet.IP[firstOctet] |= (1 << byte((8 - (ones % 8)))) - 1

	// Fill the intermediate octets with 1s and last octet with offset. This is done so to match
	// the existing behavior of allocating static IP addresses from the last /24 of the subnet.
	for i := firstOctet + 1; i < lastOctet; i++ {
		subnet.IP[i] = 255
	}
	subnet.IP[lastOctet] = DefaultKubernetesFirstConsecutiveStaticIPOffset

	return subnet.IP.String()
}

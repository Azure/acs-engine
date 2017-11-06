package acsengine

import (
	"fmt"
	"net"

	"github.com/Azure/acs-engine/pkg/api"
	"github.com/Azure/acs-engine/pkg/api/common"
	"github.com/Masterminds/semver"
)

const (
	// CniPluginVer specifies version of CNI plugin, mirrored from
	// https://github.com/containernetworking/cni/releases/download/${CNI_PLUGIN_VER}/cni-amd64-${CNI_PLUGIN_VER}.tgz
	// to https://acs-mirror.azureedge.net/cni/
	CniPluginVer = "v0.6.0"

	// AzureCniPluginVer specifies version of Azure CNI plugin, which has been mirrored from
	// https://github.com/Azure/azure-container-networking/releases/download/${AZURE_PLUGIN_VER}/azure-vnet-cni-linux-amd64-${AZURE_PLUGIN_VER}.tgz
	// to https://acs-mirror.azureedge.net/cni/
	AzureCniPluginVer = "v0.91"
)

var (
	//DefaultKubernetesSpecConfig is the default Docker image source of Kubernetes
	DefaultKubernetesSpecConfig = KubernetesSpecConfig{
		KubernetesImageBase:              "gcrio.azureedge.net/google_containers/",
		TillerImageBase:                  "gcrio.azureedge.net/kubernetes-helm/",
		KubeBinariesSASURLBase:           "https://acs-mirror.azureedge.net/wink8s/",
		WindowsTelemetryGUID:             "fb801154-36b9-41bc-89c2-f4d4f05472b0",
		CNIPluginsDownloadURL:            "https://acs-mirror.azureedge.net/cni/cni-amd64-" + CniPluginVer + ".tgz",
		VnetCNILinuxPluginsDownloadURL:   "https://acs-mirror.azureedge.net/cni/azure-vnet-cni-linux-amd64-" + AzureCniPluginVer + ".tgz",
		VnetCNIWindowsPluginsDownloadURL: "https://acs-mirror.azureedge.net/cni/azure-vnet-cni-windows-amd64-" + AzureCniPluginVer + ".zip",
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
		ImageSku:       "16.04-DAILY-LTS",
		ImagePublisher: "Canonical",
		ImageVersion:   "16.04.201710110",
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
			ResourceManagerVMDNSSuffix: "cloudapp.usgovcloudapi.net",
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
			TillerImageBase:                  "crproxy.trafficmanager.net:6000/kubernetes-helm/",
			CNIPluginsDownloadURL:            "https://acsengine.blob.core.chinacloudapi.cn/cni/cni-plugins-amd64-latest.tgz",
			VnetCNILinuxPluginsDownloadURL:   "https://acsengine.blob.core.chinacloudapi.cn/cni/azure-vnet-cni-linux-amd64-latest.tgz",
			VnetCNIWindowsPluginsDownloadURL: "https://acsengine.blob.core.chinacloudapi.cn/cni/azure-vnet-cni-windows-amd64-latest.zip",
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
			api.Ubuntu: {
				ImageOffer:     "UbuntuServer",
				ImageSku:       "16.04-LTS",
				ImagePublisher: "Canonical",
				ImageVersion:   "latest",
			},
			api.RHEL: DefaultRHELOSImageConfig,
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
	o := a.OrchestratorProfile
	o.OrchestratorVersion = common.GetValidPatchVersion(
		o.OrchestratorType,
		o.OrchestratorVersion)
	if o.OrchestratorType == api.Kubernetes {
		k8sVersion := o.OrchestratorVersion

		if o.KubernetesConfig == nil {
			o.KubernetesConfig = &api.KubernetesConfig{}
		}
		if o.KubernetesConfig.KubernetesImageBase == "" {
			o.KubernetesConfig.KubernetesImageBase = cloudSpecConfig.KubernetesSpecConfig.KubernetesImageBase
		}
		if o.KubernetesConfig.NetworkPolicy == "" {
			o.KubernetesConfig.NetworkPolicy = DefaultNetworkPolicy
		}
		if o.KubernetesConfig.ClusterSubnet == "" {
			if o.IsVNETIntegrated() {
				// When VNET integration is enabled, all masters, agents and pods share the same large subnet.
				o.KubernetesConfig.ClusterSubnet = DefaultKubernetesSubnet
			} else {
				o.KubernetesConfig.ClusterSubnet = DefaultKubernetesClusterSubnet
			}
		}
		if o.KubernetesConfig.MaxPods == 0 {
			if o.IsVNETIntegrated() {
				o.KubernetesConfig.MaxPods = DefaultKubernetesMaxPodsVNETIntegrated
			} else {
				o.KubernetesConfig.MaxPods = DefaultKubernetesMaxPods
			}
		}
		if o.KubernetesConfig.GCHighThreshold == 0 {
			o.KubernetesConfig.GCHighThreshold = DefaultKubernetesGCHighThreshold
		}
		if o.KubernetesConfig.GCLowThreshold == 0 {
			o.KubernetesConfig.GCLowThreshold = DefaultKubernetesGCLowThreshold
		}
		if o.KubernetesConfig.DNSServiceIP == "" {
			o.KubernetesConfig.DNSServiceIP = DefaultKubernetesDNSServiceIP
		}
		if o.KubernetesConfig.DockerBridgeSubnet == "" {
			o.KubernetesConfig.DockerBridgeSubnet = DefaultDockerBridgeSubnet
		}
		if o.KubernetesConfig.ServiceCIDR == "" {
			o.KubernetesConfig.ServiceCIDR = DefaultKubernetesServiceCIDR
		}
		if o.KubernetesConfig.NonMasqueradeCidr == "" {
			o.KubernetesConfig.NonMasqueradeCidr = DefaultNonMasqueradeCidr
		}
		if o.KubernetesConfig.NodeStatusUpdateFrequency == "" {
			o.KubernetesConfig.NodeStatusUpdateFrequency = KubeConfigs[k8sVersion]["nodestatusfreq"]
		}
		if o.KubernetesConfig.CtrlMgrNodeMonitorGracePeriod == "" {
			o.KubernetesConfig.CtrlMgrNodeMonitorGracePeriod = KubeConfigs[k8sVersion]["nodegraceperiod"]
		}
		if o.KubernetesConfig.CtrlMgrPodEvictionTimeout == "" {
			o.KubernetesConfig.CtrlMgrPodEvictionTimeout = KubeConfigs[k8sVersion]["podeviction"]
		}
		if o.KubernetesConfig.CtrlMgrRouteReconciliationPeriod == "" {
			o.KubernetesConfig.CtrlMgrRouteReconciliationPeriod = KubeConfigs[k8sVersion]["routeperiod"]
		}
		// Enforce sane cloudprovider backoff defaults, if CloudProviderBackoff is true in KubernetesConfig
		if o.KubernetesConfig.CloudProviderBackoff == true {
			if o.KubernetesConfig.CloudProviderBackoffDuration == 0 {
				o.KubernetesConfig.CloudProviderBackoffDuration = DefaultKubernetesCloudProviderBackoffDuration
			}
			if o.KubernetesConfig.CloudProviderBackoffExponent == 0 {
				o.KubernetesConfig.CloudProviderBackoffExponent = DefaultKubernetesCloudProviderBackoffExponent
			}
			if o.KubernetesConfig.CloudProviderBackoffJitter == 0 {
				o.KubernetesConfig.CloudProviderBackoffJitter = DefaultKubernetesCloudProviderBackoffJitter
			}
			if o.KubernetesConfig.CloudProviderBackoffRetries == 0 {
				o.KubernetesConfig.CloudProviderBackoffRetries = DefaultKubernetesCloudProviderBackoffRetries
			}
		}
		k8sSemVer, _ := semver.NewVersion(k8sVersion)
		constraint, _ := semver.NewConstraint(">= 1.6.6")
		// Enforce sane cloudprovider rate limit defaults, if CloudProviderRateLimit is true in KubernetesConfig
		// For k8s version greater or equal to 1.6.6, we will set the default CloudProviderRate* settings
		if o.KubernetesConfig.CloudProviderRateLimit == true && constraint.Check(k8sSemVer) {
			if o.KubernetesConfig.CloudProviderRateLimitQPS == 0 {
				o.KubernetesConfig.CloudProviderRateLimitQPS = DefaultKubernetesCloudProviderRateLimitQPS
			}
			if o.KubernetesConfig.CloudProviderRateLimitBucket == 0 {
				o.KubernetesConfig.CloudProviderRateLimitBucket = DefaultKubernetesCloudProviderRateLimitBucket
			}
		}

		// default etcd version
		if "" == o.KubernetesConfig.EtcdVersion {
			o.KubernetesConfig.EtcdVersion = "2.5.2"
		}

		if "" == a.OrchestratorProfile.KubernetesConfig.TillerCPURequests {
			a.OrchestratorProfile.KubernetesConfig.TillerCPURequests = DefaultTillerCPURequests
		}

		if "" == a.OrchestratorProfile.KubernetesConfig.TillerCPULimit {
			a.OrchestratorProfile.KubernetesConfig.TillerCPULimit = DefaultTillerCPULimit
		}

		if "" == a.OrchestratorProfile.KubernetesConfig.TillerMemoryRequests {
			a.OrchestratorProfile.KubernetesConfig.TillerMemoryRequests = DefaultTillerMemoryRequests
		}

		if "" == a.OrchestratorProfile.KubernetesConfig.TillerMemoryLimit {
			a.OrchestratorProfile.KubernetesConfig.TillerMemoryLimit = DefaultTillerMemoryLimit
		}

	} else if o.OrchestratorType == api.DCOS {
		if o.DcosConfig == nil {
			o.DcosConfig = &api.DcosConfig{}
		}
		if o.DcosConfig.DcosWindowsBootstrapURL == "" {
			o.DcosConfig.DcosWindowsBootstrapURL = DefaultDCOSSpecConfig.DCOSWindowsBootstrapDownloadURL
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

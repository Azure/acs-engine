package acsengine

import (
	"fmt"
	"net"

	"github.com/Azure/acs-engine/pkg/api"
)

var (
	//AzureCloudSpec is the default configurations for global azure.
	AzureCloudSpec = AzureEnvironmentSpecConfig{
		//DockerConfigAzurePublicCloud specify the default script location of docker installer script
		DockerSpecConfig: DockerSpecConfig{
			DefaultDockerInstallScriptURL: "https://get.docker.com/",
		},
		//KubeConfigAzurePublicCloud is the default kubernetes container image url.
		KubernetesSpecConfig: KubernetesSpecConfig{
			DefaultKubernetesHyperkubeSpec:         "gcr.io/google_containers/hyperkube-amd64:v1.5.1",
			DefaultKubernetesDashboardSpec:         "gcr.io/google_containers/kubernetes-dashboard-amd64:v1.5.1",
			DefaultKubernetesExechealthzSpec:       "gcr.io/google_containers/exechealthz-amd64:1.2",
			DefaultKubernetesAddonResizerSpec:      "gcr.io/google_containers/addon-resizer:1.6",
			DefaultKubernetesHeapsterSpec:          "gcr.io/google_containers/heapster:v1.2.0",
			DefaultKubernetesDNSSpec:               "gcr.io/google_containers/kubedns-amd64:1.7",
			DefaultKubernetesAddonManagerSpec:      "gcr.io/google_containers/kube-addon-manager-amd64:v5.1",
			DefaultKubernetesDNSMasqSpec:           "gcr.io/google_containers/kube-dnsmasq-amd64:1.3",
			DefaultKubernetesPodInfraContainerSpec: "gcr.io/google_containers/pause-amd64:3.0",
			DefaultKubectlDownloadURL:              "https://storage.googleapis.com/kubernetes-release/release/" + DefaultKubectlVersion + "/bin/linux/amd64/kubectl",
		},

		DCOSSpecConfig: DCOSSpecConfig{
			DCOS173_BootstrapDownloadURL: "https://az837203.vo.msecnd.net/dcos/testing/bootstrap/${BOOTSTRAP_ID}.bootstrap.tar.xz",
			DCOS184_BootstrapDownloadURL: "https://dcosio.azureedge.net/dcos/testing/bootstrap/${BOOTSTRAP_ID}.bootstrap.tar.xz",
			DCOS187_BootstrapDownloadURL: "https://dcosio.azureedge.net/dcos/stable/bootstrap/e73ba2b1cd17795e4dcb3d6647d11a29b9c35084.bootstrap.tar.xz",
		},
	}

	//AzureChinaCloudSpec is the configurations for Azure China (Mooncake)
	AzureChinaCloudSpec = AzureEnvironmentSpecConfig{
		//DockerConfigAzureChinaCloud specify the docker install script download URL in China.
		DockerSpecConfig: DockerSpecConfig{
			DefaultDockerInstallScriptURL: "https://acsengine.blob.core.chinacloudapi.cn/docker/install-docker",
		},
		//KubeConfigAzureChinaCloud - Due to Chinese firewall issue, the default containers from google is blocked, use the Chinese local mirror instead
		KubernetesSpecConfig: KubernetesSpecConfig{
			DefaultKubernetesHyperkubeSpec:         "mirror.azure.cn:5000/google_containers/hyperkube-amd64:v1.5.1",
			DefaultKubernetesDashboardSpec:         "mirror.azure.cn:5000/google_containers/kubernetes-dashboard-amd64:v1.5.1",
			DefaultKubernetesExechealthzSpec:       "mirror.azure.cn:5000/google_containers/exechealthz-amd64:1.2",
			DefaultKubernetesAddonResizerSpec:      "mirror.azure.cn:5000/google_containers/addon-resizer:1.6",
			DefaultKubernetesHeapsterSpec:          "mirror.azure.cn:5000/google_containers/heapster:v1.2.0",
			DefaultKubernetesDNSSpec:               "mirror.azure.cn:5000/google_containers/kubedns-amd64:1.7",
			DefaultKubernetesAddonManagerSpec:      "mirror.azure.cn:5000/google_containers/kube-addon-manager-amd64:v5.1",
			DefaultKubernetesDNSMasqSpec:           "mirror.azure.cn:5000/google_containers/kube-dnsmasq-amd64:1.3",
			DefaultKubernetesPodInfraContainerSpec: "mirror.azure.cn:5000/google_containers/pause-amd64:3.0",
			DefaultKubectlDownloadURL:              "https://acsengine.blob.core.chinacloudapi.cn/kubernetes/kubectl/" + DefaultKubectlVersion + "/kubectl",
		},
		DCOSSpecConfig: DCOSSpecConfig{
			DCOS173_BootstrapDownloadURL: "https://acsengine.blob.core.chinacloudapi.cn/dcos/df308b6fc3bd91e1277baa5a3db928ae70964722.bootstrap.tar.xz",
			DCOS184_BootstrapDownloadURL: "https://acsengine.blob.core.chinacloudapi.cn/dcos/5b4aa43610c57ee1d60b4aa0751a1fb75824c083.bootstrap.tar.xz",
			DCOS187_BootstrapDownloadURL: "https://acsengine.blob.core.chinacloudapi.cn/dcos/e73ba2b1cd17795e4dcb3d6647d11a29b9c35084.bootstrap.tar.xz",
		},
	}

	//Set the AzureUSGovernment and AzureGermanCloud the same as the AzureCloud
	AzureUSGovernmentSpec = AzureCloud
	AzureGermanCloudSpec  = AzureCloud
)

// SetPropertiesDefaults for the container Properties, returns true if certs are generated
func SetPropertiesDefaults(properties *api.Properties, locationCode int) (bool, error) {

	setOrchestratorDefaults(properties, locationCode)

	setMasterNetworkDefaults(properties)

	setAgentNetworkDefaults(properties)

	setStorageDefaults(properties)

	certsGenerated, e := setDefaultCerts(properties)
	if e != nil {
		return false, e
	}
	return certsGenerated, nil
}

//GetKubeCloudSpecConfig returns the kubenernetes container images url configurations based on the deploy target environment
//for example: if the target is the public azure, then the default container image url should be gcr.io/google_container/...
//if the target is azure china, then the default container image should be mirror.azure.cn:5000/google_container/...
func GetKubeCloudSpecConfig(locationCode int) AzureEnvironmentSpecConfig {
	kubeSpecConfig := AzureCloudSpec
	switch locationCode {
	case AzureChinaCloud:
		kubeSpecConfig = AzureChinaCloudSpec
	default:
		kubeSpecConfig = AzureCloudSpec
	}

	return kubeSpecConfig
}

// setOrchestratorDefaults for orchestrators
func setOrchestratorDefaults(a *api.Properties, locationCode int) {
	cloudSpecConfig := GetKubeCloudSpecConfig(locationCode)
	if a.OrchestratorProfile.OrchestratorType == api.Kubernetes {
		a.OrchestratorProfile.KubernetesConfig.KubernetesHyperkubeSpec = cloudSpecConfig.KubernetesSpecConfig.DefaultKubernetesHyperkubeSpec
		a.OrchestratorProfile.KubernetesConfig.KubectlVersion = DefaultKubectlVersion
		a.OrchestratorProfile.KubernetesConfig.KubernetesAddonManagerSpec = cloudSpecConfig.KubernetesSpecConfig.DefaultKubernetesAddonManagerSpec
		a.OrchestratorProfile.KubernetesConfig.KubernetesAddonResizerSpec = cloudSpecConfig.KubernetesSpecConfig.DefaultKubernetesAddonResizerSpec
		a.OrchestratorProfile.KubernetesConfig.KubernetesDashboardSpec = cloudSpecConfig.KubernetesSpecConfig.DefaultKubernetesDashboardSpec
		a.OrchestratorProfile.KubernetesConfig.KubernetesExecHealthzSpec = cloudSpecConfig.KubernetesSpecConfig.DefaultKubernetesExechealthzSpec
		a.OrchestratorProfile.KubernetesConfig.KubernetesHeapsterSpec = cloudSpecConfig.KubernetesSpecConfig.DefaultKubernetesHeapsterSpec
		a.OrchestratorProfile.KubernetesConfig.KubernetesKubeDNSSpec = cloudSpecConfig.KubernetesSpecConfig.DefaultKubernetesDNSSpec
		a.OrchestratorProfile.KubernetesConfig.KubernetesDNSMasqSpec = cloudSpecConfig.KubernetesSpecConfig.DefaultKubernetesDNSMasqSpec
		a.OrchestratorProfile.KubernetesConfig.KubernetesPodInfraContainerSpec = cloudSpecConfig.KubernetesSpecConfig.DefaultKubernetesPodInfraContainerSpec
		a.OrchestratorProfile.KubernetesConfig.DockerInstallScriptURL = cloudSpecConfig.DockerSpecConfig.DefaultDockerInstallScriptURL
		a.OrchestratorProfile.KubernetesConfig.KubectlDownloadURL = cloudSpecConfig.KubernetesSpecConfig.DefaultKubectlDownloadURL
	}
	if a.OrchestratorProfile.OrchestratorType == api.DCOS {
		a.OrchestratorProfile.DCOSConfig.DCOS173_BootstrapDownloadURL = cloudSpecConfig.DCOSSpecConfig.DCOS173_BootstrapDownloadURL
		a.OrchestratorProfile.DCOSConfig.DCOS184_BootstrapDownloadURL = cloudSpecConfig.DCOSSpecConfig.DCOS184_BootstrapDownloadURL
		a.OrchestratorProfile.DCOSConfig.DCOS187_BootstrapDownloadURL = cloudSpecConfig.DCOSSpecConfig.DCOS187_BootstrapDownloadURL

		a.OrchestratorProfile.OrchestratorType = api.DCOS187
	}
}

// SetMasterNetworkDefaults for masters
func setMasterNetworkDefaults(a *api.Properties) {
	if !a.MasterProfile.IsCustomVNET() {
		if a.OrchestratorProfile.OrchestratorType == api.Kubernetes {
			a.MasterProfile.Subnet = DefaultKubernetesMasterSubnet
			a.MasterProfile.FirstConsecutiveStaticIP = DefaultFirstConsecutiveKubernetesStaticIP
		} else if a.HasWindows() {
			a.MasterProfile.Subnet = DefaultSwarmWindowsMasterSubnet
			a.MasterProfile.FirstConsecutiveStaticIP = DefaultSwarmWindowsFirstConsecutiveStaticIP
		} else {
			a.MasterProfile.Subnet = DefaultMasterSubnet
			a.MasterProfile.FirstConsecutiveStaticIP = DefaultFirstConsecutiveStaticIP
		}
	}
}

// SetAgentNetworkDefaults for agents
func setAgentNetworkDefaults(a *api.Properties) {
	// configure the subnets if not in custom VNET
	if !a.MasterProfile.IsCustomVNET() {
		subnetCounter := 0
		for i := range a.AgentPoolProfiles {
			profile := &a.AgentPoolProfiles[i]

			if a.OrchestratorProfile.OrchestratorType == api.Kubernetes {
				profile.Subnet = a.MasterProfile.Subnet
			} else {
				profile.Subnet = fmt.Sprintf(DefaultAgentSubnetTemplate, subnetCounter)
			}

			subnetCounter++
		}
	}
}

// setStorageDefaults for agents
func setStorageDefaults(a *api.Properties) {
	for i := range a.AgentPoolProfiles {
		profile := &a.AgentPoolProfiles[i]
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
	firstMasterIP := net.ParseIP(a.MasterProfile.FirstConsecutiveStaticIP)

	if firstMasterIP == nil {
		return false, fmt.Errorf("MasterProfile.FirstConsecutiveStaticIP '%s' is an invalid IP address", a.MasterProfile.FirstConsecutiveStaticIP)
	}

	ips := []net.IP{firstMasterIP}

	for i := 1; i < a.MasterProfile.Count; i++ {
		ips = append(ips, net.IP{firstMasterIP[12], firstMasterIP[13], firstMasterIP[14], firstMasterIP[15] + byte(i)})
	}

	// use the specified Certificate Authority pair, or generate a new pair
	var caPair *PkiKeyCertPair
	if len(a.CertificateProfile.CaCertificate) != 0 && len(a.CertificateProfile.GetCAPrivateKey()) != 0 {
		caPair = &PkiKeyCertPair{CertificatePem: a.CertificateProfile.CaCertificate, PrivateKeyPem: a.CertificateProfile.GetCAPrivateKey()}
	} else {
		caCertificate, caPrivateKey, err := createCertificate("ca", nil, nil, false, nil, nil)
		if err != nil {
			return false, err
		}
		caPair = &PkiKeyCertPair{CertificatePem: string(certificateToPem(caCertificate.Raw)), PrivateKeyPem: string(privateKeyToPem(caPrivateKey))}
		a.CertificateProfile.CaCertificate = caPair.CertificatePem
		a.CertificateProfile.SetCAPrivateKey(caPair.PrivateKeyPem)
	}

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
	if len(a.CertificateProfile.APIServerCertificate) > 0 || len(a.CertificateProfile.APIServerPrivateKey) > 0 ||
		len(a.CertificateProfile.ClientCertificate) > 0 || len(a.CertificateProfile.ClientPrivateKey) > 0 {
		return false
	}

	switch a.OrchestratorProfile.OrchestratorType {
	case api.DCOS:
		return false
	case api.DCOS184:
		return false
	case api.DCOS173:
		return false
	case api.Swarm:
		return false
	case api.Kubernetes:
		return true
	default:
		return false
	}
}

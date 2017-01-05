package acsengine

import (
	"fmt"
	"net"

	"github.com/Azure/acs-engine/pkg/api"
)

// SetPropertiesDefaults for the container Properties, returns true if certs are generated
func SetPropertiesDefaults(p *api.Properties) (bool, error) {

	setOrchestratorDefaults(p)

	setMasterNetworkDefaults(p)

	setAgentNetworkDefaults(p)

	setStorageDefaults(p)

	certsGenerated, e := setDefaultCerts(p)
	if e != nil {
		return false, e
	}
	return certsGenerated, nil
}

//DefaultCloudSpecConfigFromEnvironment returns the kubenernetes container images url configurations based on the deploy target environment
//for example: if the target is the public azure, then the default container image url should be gcr.io/google_container/...
//if the target is azure china, then the default container image should be mirror.azure.cn:5000/google_container/...
func DefaultCloudSpecConfigFromEnvironment(environment api.Environment) AzureEnvironmentSpecConfig {
	kubeSpecConfig := AzurePublicCloud
	if environment == "AzureChinaCloud" {
		kubeSpecConfig = AzureChinaCloud
	}

	return kubeSpecConfig
}

// setOrchestratorDefaults for orchestrators
func setOrchestratorDefaults(a *api.Properties) {
	if a.OrchestratorProfile.OrchestratorType == api.Kubernetes {
		cloudSpecConfig := DefaultCloudSpecConfigFromEnvironment(a.Environment)
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

	caPair, apiServerPair, clientPair, kubeConfigPair, err := CreatePki(masterExtraFQDNs, ips, DefaultKubernetesClusterDomain)
	if err != nil {
		return false, err
	}

	a.CertificateProfile.APIServerCertificate = apiServerPair.CertificatePem
	a.CertificateProfile.APIServerPrivateKey = apiServerPair.PrivateKeyPem
	a.CertificateProfile.CaCertificate = caPair.CertificatePem
	a.CertificateProfile.SetCAPrivateKey(caPair.PrivateKeyPem)
	a.CertificateProfile.ClientCertificate = clientPair.CertificatePem
	a.CertificateProfile.ClientPrivateKey = clientPair.PrivateKeyPem
	a.CertificateProfile.KubeConfigCertificate = kubeConfigPair.CertificatePem
	a.CertificateProfile.KubeConfigPrivateKey = kubeConfigPair.PrivateKeyPem

	return true, nil
}

func certGenerationRequired(a *api.Properties) bool {
	if len(a.CertificateProfile.APIServerCertificate) > 0 || len(a.CertificateProfile.APIServerPrivateKey) > 0 ||
		len(a.CertificateProfile.CaCertificate) > 0 ||
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

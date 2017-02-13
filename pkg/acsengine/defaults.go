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

// setOrchestratorDefaults for orchestrators
func setOrchestratorDefaults(a *api.Properties) {
	if a.OrchestratorProfile.OrchestratorType == api.Kubernetes {
		a.OrchestratorProfile.KubernetesConfig.KubernetesHyperkubeSpec = DefaultKubernetesHyperkubeSpec
		a.OrchestratorProfile.KubernetesConfig.KubectlVersion = DefaultKubectlVersion
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

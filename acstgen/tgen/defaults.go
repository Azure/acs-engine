package tgen

import (
	"fmt"
	"net"

	"./../api/vlabs"
	"./../util"
)

// SetAcsClusterDefaults for an AcsCluster, returns true if certs are generated
func SetAcsClusterDefaults(a *vlabs.AcsCluster) (bool, error) {

	setMasterNetworkDefaults(a)

	setAgentNetworkDefaults(a)

	setStorageDefaults(a)

	certsGenerated, e := setDefaultCerts(a)
	if e != nil {
		return false, e
	}
	return certsGenerated, nil
}

// SetMasterNetworkDefaults for masters
func setMasterNetworkDefaults(a *vlabs.AcsCluster) {
	if !a.MasterProfile.IsCustomVNET() {
		if a.OrchestratorProfile.OrchestratorType == vlabs.Kubernetes {
			a.MasterProfile.SetSubnet(DefaultKubernetesMasterSubnet)
			a.MasterProfile.FirstConsecutiveStaticIP = DefaultFirstConsecutiveKubernetesStaticIP
		} else if a.HasWindows() {
			a.MasterProfile.SetSubnet(DefaultSwarmWindowsMasterSubnet)
			a.MasterProfile.FirstConsecutiveStaticIP = DefaultSwarmWindowsFirstConsecutiveStaticIP
		} else {
			a.MasterProfile.SetSubnet(DefaultMasterSubnet)
			a.MasterProfile.FirstConsecutiveStaticIP = DefaultFirstConsecutiveStaticIP
		}
	}
}

// SetAgentNetworkDefaults for agents
func setAgentNetworkDefaults(a *vlabs.AcsCluster) {
	// configure the subnets if not in custom VNET
	if !a.MasterProfile.IsCustomVNET() {
		subnetCounter := 0
		for i := range a.AgentPoolProfiles {
			profile := &a.AgentPoolProfiles[i]

			if a.OrchestratorProfile.OrchestratorType == vlabs.Kubernetes {
				profile.SetSubnet(a.MasterProfile.GetSubnet())
			} else {
				profile.SetSubnet(fmt.Sprintf(DefaultAgentSubnetTemplate, subnetCounter))
			}

			subnetCounter++
		}
	}
}

// setStorageDefaults for agents
func setStorageDefaults(a *vlabs.AcsCluster) {
	for i := range a.AgentPoolProfiles {
		profile := &a.AgentPoolProfiles[i]
		if len(profile.StorageType) == 0 {
			profile.StorageType = vlabs.StorageExternal
		}
	}
}

func setDefaultCerts(a *vlabs.AcsCluster) (bool, error) {
	if !certGenerationRequired(a) {
		return false, nil
	}

	masterWildCardFQDN := FormatAzureProdFQDN(a.MasterProfile.DNSPrefix, "*")
	masterExtraFQDNs := FormatAzureProdFQDNs(a.MasterProfile.DNSPrefix)
	firstMasterIP := net.ParseIP(a.MasterProfile.FirstConsecutiveStaticIP)

	if firstMasterIP == nil {
		return false, fmt.Errorf("MasterProfile.FirstConsecutiveStaticIP '%s' is an invalid IP address", a.MasterProfile.FirstConsecutiveStaticIP)
	}

	ips := []net.IP{firstMasterIP}

	for i := 1; i < a.MasterProfile.Count; i++ {
		ips = append(ips, net.IP{firstMasterIP[12], firstMasterIP[13], firstMasterIP[14], firstMasterIP[15] + byte(i)})
	}

	caPair, apiServerPair, clientPair, kubeConfigPair, err := util.CreatePki(masterWildCardFQDN, masterExtraFQDNs, ips, DefaultKubernetesClusterDomain)
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

func certGenerationRequired(a *vlabs.AcsCluster) bool {
	if len(a.CertificateProfile.APIServerCertificate) > 0 || len(a.CertificateProfile.APIServerPrivateKey) > 0 ||
		len(a.CertificateProfile.CaCertificate) > 0 ||
		len(a.CertificateProfile.ClientCertificate) > 0 || len(a.CertificateProfile.ClientPrivateKey) > 0 {
		return false
	}

	switch a.OrchestratorProfile.OrchestratorType {
	case vlabs.DCOS:
		return false
	case vlabs.DCOS184:
		return false
	case vlabs.DCOS173:
		return false
	case vlabs.Swarm:
		return false
	case vlabs.Kubernetes:
		return true
	default:
		return false
	}
}

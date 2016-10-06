package templategenerator

import (
	"encoding/base64"
	"fmt"
	"net"

	"./../api/vlabs"
	"./../util"
)

// SetAcsClusterDefaults for an AcsCluster
func SetAcsClusterDefaults(a *vlabs.AcsCluster) error {
	setMasterNetworkDefaults(a)

	setAgentNetworkDefaults(a)

	if e := setDefaultCerts(&a.OrchestratorProfile, &a.MasterProfile); e != nil {
		return e
	}
	return nil
}

// SetMasterNetworkDefaults for masters
func setMasterNetworkDefaults(a *vlabs.AcsCluster) {
	if !a.MasterProfile.IsCustomVNET() {
		if a.OrchestratorProfile.OrchestratorType == vlabs.Kubernetes {
			a.MasterProfile.SetSubnet(DefaultKubernetesMasterSubnet)
			a.MasterProfile.FirstConsecutiveStaticIP = DefaultFirstConsecutiveKubernetesStaticIP
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

func setDefaultCerts(o *vlabs.OrchestratorProfile, m *vlabs.MasterProfile) error {
	// auto generate certs if none of them have been set by customer
	if len(o.ApiServerCertificate) == 0 && len(o.ApiServerPrivateKey) == 0 &&
		len(o.CaCertificate) == 0 &&
		len(o.ClientCertificate) == 0 && len(o.ClientPrivateKey) == 0 {
		masterWildCardFQDN := FormatAzureProdFQDN(m.DNSPrefix, "*")
		masterExtraFQDNs := FormatAzureProdFQDNs(m.DNSPrefix)

		firstMasterIP := net.ParseIP(m.FirstConsecutiveStaticIP)

		if firstMasterIP == nil {
			return fmt.Errorf("MasterProfile.FirstConsecutiveStaticIP '%s' is an invalid IP address", m.FirstConsecutiveStaticIP)
		}

		ips := []net.IP{firstMasterIP}

		for i := 1; i < m.Count; i++ {
			ips = append(ips, net.IP{firstMasterIP[12], firstMasterIP[13], firstMasterIP[14], firstMasterIP[15] + byte(i)})
		}

		caPair, apiServerPair, clientPair, err := util.CreatePki(masterWildCardFQDN, masterExtraFQDNs, ips, DefaultKubernetesClusterDomain)
		if err != nil {
			return err
		}

		o.ApiServerCertificate = base64.URLEncoding.EncodeToString([]byte(apiServerPair.CertificatePem))
		o.ApiServerPrivateKey = base64.URLEncoding.EncodeToString([]byte(apiServerPair.PrivateKeyPem))
		o.CaCertificate = base64.URLEncoding.EncodeToString([]byte(caPair.CertificatePem))
		o.SetCAPrivateKey(base64.URLEncoding.EncodeToString([]byte(caPair.PrivateKeyPem)))
		o.ClientCertificate = base64.URLEncoding.EncodeToString([]byte(clientPair.CertificatePem))
		o.ClientPrivateKey = base64.URLEncoding.EncodeToString([]byte(clientPair.PrivateKeyPem))
	}
	return nil
}

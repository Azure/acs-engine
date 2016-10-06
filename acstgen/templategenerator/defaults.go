package templategenerator

import (
	"fmt"
	"hash/fnv"
	"math/rand"
	"net"

	"./../api/vlabs"
	"./../util"
)

// SetAcsClusterDefaults for an AcsCluster, returns true if certs are generated
func SetAcsClusterDefaults(a *vlabs.AcsCluster) (bool, error) {

	if len(a.OrchestratorProfile.ClusterID) == 0 {
		a.OrchestratorProfile.ClusterID = generateClusterID(a)
	}

	setMasterNetworkDefaults(a)

	setAgentNetworkDefaults(a)

	certsGenerated, e := setDefaultCerts(&a.OrchestratorProfile, &a.MasterProfile)
	if e != nil {
		return false, e
	}
	return certsGenerated, nil
}

// generateClusterID creates a unique 8 string cluster ID
func generateClusterID(acsCluster *vlabs.AcsCluster) string {
	uniqueNameSuffixSize := 8
	// the name suffix uniquely identifies the cluster and is generated off a hash
	// from the master dns name
	h := fnv.New64a()
	h.Write([]byte(acsCluster.MasterProfile.DNSPrefix))
	rand.Seed(int64(h.Sum64()))
	return fmt.Sprintf("%08d", rand.Uint32())[:uniqueNameSuffixSize]
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

func setDefaultCerts(o *vlabs.OrchestratorProfile, m *vlabs.MasterProfile) (bool, error) {
	certsGenerated := false
	// auto generate certs if none of them have been set by customer
	if len(o.ApiServerCertificate) > 0 || len(o.ApiServerPrivateKey) > 0 ||
		len(o.CaCertificate) > 0 ||
		len(o.ClientCertificate) > 0 || len(o.ClientPrivateKey) > 0 {
		return certsGenerated, nil
	}
	masterWildCardFQDN := FormatAzureProdFQDN(m.DNSPrefix, "*")
	masterExtraFQDNs := FormatAzureProdFQDNs(m.DNSPrefix)
	firstMasterIP := net.ParseIP(m.FirstConsecutiveStaticIP)

	if firstMasterIP == nil {
		return false, fmt.Errorf("MasterProfile.FirstConsecutiveStaticIP '%s' is an invalid IP address", m.FirstConsecutiveStaticIP)
	}

	ips := []net.IP{firstMasterIP}

	for i := 1; i < m.Count; i++ {
		ips = append(ips, net.IP{firstMasterIP[12], firstMasterIP[13], firstMasterIP[14], firstMasterIP[15] + byte(i)})
	}

	caPair, apiServerPair, clientPair, err := util.CreatePki(masterWildCardFQDN, masterExtraFQDNs, ips, DefaultKubernetesClusterDomain)
	if err != nil {
		return false, err
	}

	certsGenerated = true

	o.ApiServerCertificate = apiServerPair.CertificatePem
	o.ApiServerPrivateKey = apiServerPair.PrivateKeyPem
	o.CaCertificate = caPair.CertificatePem
	o.SetCAPrivateKey(caPair.PrivateKeyPem)
	o.ClientCertificate = clientPair.CertificatePem
	o.ClientPrivateKey = clientPair.PrivateKeyPem

	return certsGenerated, nil
}

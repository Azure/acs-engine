package v20170727

import (
	"fmt"
	"strings"
)

// Validate will validate an agent pool
func (p *AgentPool) Validate() error {

	// -------------------------------
	// Kubernetes Version
	if p.Properties.KubernetesVersion == "" {
		return fmt.Errorf("Empty Kubernetes version")
	}
	if !strings.Contains(p.Properties.KubernetesVersion, ".") {
		return fmt.Errorf("Invalid Kubernetes version")
	}

	// -------------------------------
	// Kubernetes Endpoint
	if p.Properties.KubernetesEndpoint == "" {
		return fmt.Errorf("Empty Kubernetes endpoint")
	}
	if !strings.Contains(p.Properties.KubernetesEndpoint, ".") {
		return fmt.Errorf("Invalid Kubernetes endpoint")
	}
	// -------------------------------
	// Agent pools
	if len(p.Properties.AgentPoolProfiles) < 1 {
		return fmt.Errorf("Missing agent pool definition")
	}
	for _, agentPool := range p.Properties.AgentPoolProfiles {
		if agentPool.Count < 1 {
			return fmt.Errorf("Must using a value greater than 0 for agent pool count")
		}
		if agentPool.OSDiskSizeGb != 0 && agentPool.OSDiskSizeGb < 2 {
			return fmt.Errorf("Must use a disk size greater than or equal to 2gb for agent pool")
		}
		if agentPool.Name == "" {
			return fmt.Errorf("Empty name for agent pool")
		}
	}

	// -------------------------------
	// Network Profile
	if p.Properties.NetworkProfile.AgentCidr == "" {
		return fmt.Errorf("Empty service CIDR in network profile")
	}

	// -------------------------------
	// Certificate Profile
	if p.Properties.CertificateProfile.APIServerCertificate == "" {
		return fmt.Errorf("Empty API certificate")
	}
	if p.Properties.CertificateProfile.APIServerPrivateKey == "" {
		return fmt.Errorf("Empty API key")
	}
	if p.Properties.CertificateProfile.CaCertificate == "" {
		return fmt.Errorf("Empty CA certificate")
	}
	if p.Properties.CertificateProfile.CaPrivateKey == "" {
		return fmt.Errorf("Empty CA key")
	}
	if p.Properties.CertificateProfile.ClientCertificate == "" {
		return fmt.Errorf("Empty client certificate")
	}
	if p.Properties.CertificateProfile.ClientPrivateKey == "" {
		return fmt.Errorf("Empty client key")
	}
	if p.Properties.CertificateProfile.KubeConfigCertificate == "" {
		return fmt.Errorf("Empty kubeconfig certificate")
	}
	if p.Properties.CertificateProfile.KubeConfigPrivateKey == "" {
		return fmt.Errorf("Empty kubeconfig key")
	}

	// -------------------------------
	// Jump Box
	if p.Properties.JumpBoxProfile.VMSize == "" {
		return fmt.Errorf("Empty jumpbox vm size")
	}
	// TODO: validate and/or default p.Properties.JumpBoxProfile.InternalAdress

	// -------------------------------
	// SSH
	if len(p.Properties.LinuxProfile.SSH.PublicKeys) < 1 {
		return fmt.Errorf("Must specify SSH Public Key")
	}
	for _, key := range p.Properties.LinuxProfile.SSH.PublicKeys {
		if key.KeyData == "" {
			return fmt.Errorf("Empty key data for SSH key")
		}
	}

	return nil
}

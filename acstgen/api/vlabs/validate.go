package vlabs

import (
	"errors"
	"fmt"
	"net"
	"regexp"
)

// Validate implements APIObject
func (o *OrchestratorProfile) Validate() error {
	switch o.OrchestratorType {
	case DCOS:
	case DCOS184:
	case DCOS173:
	case Swarm:
	case Kubernetes:
	default:
		return fmt.Errorf("OrchestratorProfile has unknown orchestrator: %s", o.OrchestratorType)
	}

	if o.OrchestratorType == Kubernetes && len(o.ServicePrincipalClientID) == 0 {
		return fmt.Errorf("the service principal client ID must be specified with Orchestrator %s", o.OrchestratorType)
	}

	if o.OrchestratorType == Kubernetes && len(o.ServicePrincipalClientSecret) == 0 {
		return fmt.Errorf("the service principal client secrect must be specified with Orchestrator %s", o.OrchestratorType)
	}

	if o.OrchestratorType != Kubernetes && (len(o.ServicePrincipalClientID) > 0 || len(o.ServicePrincipalClientSecret) > 0) {
		return fmt.Errorf("Service principal and secret is not required for orchestrator %s", o.OrchestratorType)
	}

	if o.OrchestratorType == Kubernetes {
		if len(o.ApiServerCertificate) > 0 || len(o.ApiServerPrivateKey) > 0 || len(o.CaCertificate) > 0 || len(o.ClientCertificate) > 0 || len(o.ClientPrivateKey) > 0 {
			return fmt.Errorf("API, CA, and Client certs are required for orchestrator %s", o.OrchestratorType)
		}
	} else {
		if len(o.ApiServerCertificate) > 0 || len(o.ApiServerPrivateKey) > 0 || len(o.CaCertificate) > 0 || len(o.ClientCertificate) > 0 || len(o.ClientPrivateKey) > 0 {
			return fmt.Errorf("API, CA, and Client certs are not required for orchestrator %s", o.OrchestratorType)
		}
	}

	return nil
}

// Validate implements APIObject
func (m *MasterProfile) Validate() error {
	if m.Count != 1 && m.Count != 3 && m.Count != 5 {
		return fmt.Errorf("MasterProfile count needs to be 1, 3, or 5")
	}
	if e := validateName(m.DNSPrefix, "MasterProfile.DNSPrefix"); e != nil {
		return e
	}
	if e := validateDNSName(m.DNSPrefix); e != nil {
		return e
	}
	if e := validateName(m.VMSize, "MasterProfile.VMSize"); e != nil {
		return e
	}
	return nil
}

// Validate implements APIObject
func (a *AgentPoolProfile) Validate() error {
	if e := validateName(a.Name, "AgentPoolProfile.Name"); e != nil {
		return e
	}
	if e := validatePoolName(a.Name); e != nil {
		return e
	}
	if a.Count < MinAgentCount || a.Count > MaxAgentCount {
		return fmt.Errorf("AgentPoolProfile count needs to be in the range [%d,%d]", MinAgentCount, MaxAgentCount)
	}
	if e := validateName(a.VMSize, "AgentPoolProfile.VMSize"); e != nil {
		return e
	}
	if len(a.Ports) > 0 {
		if e := validateUniquePorts(a.Ports, a.Name); e != nil {
			return e
		}
		for _, port := range a.Ports {
			if port < MinPort || port > MaxPort {
				return fmt.Errorf("AgentPoolProfile Ports must be in the range[%d, %d]", MinPort, MaxPort)
			}
		}
		if e := validateName(a.DNSPrefix, "AgentPoolProfile.DNSPrefix when specifying AgentPoolProfile Ports"); e != nil {
			return e
		}
		if e := validateDNSName(a.DNSPrefix); e != nil {
			return e
		}
	}
	if len(a.DiskSizesGB) > 0 && !a.IsStateful {
		return fmt.Errorf("Disks were specified on a non stateful cluster named '%s'.  Ensure you add '\"isStateful\": true' to the model", a.Name)
	}
	if len(a.DiskSizesGB) > MaxDisks {
		return fmt.Errorf("A maximum of %d disks may be specified.  %d disks were specified for cluster named '%s'", MaxDisks, len(a.DiskSizesGB), a.Name)
	}
	if len(a.Ports) == 0 && len(a.DNSPrefix) > 0 {
		return fmt.Errorf("AgentPoolProfile.Ports must be non empty when AgentPoolProfile.DNSPrefix is specified")
	}
	return nil
}

// Validate implements APIObject
func (l *LinuxProfile) Validate() error {
	if e := validateName(l.AdminUsername, "LinuxProfile.AdminUsername"); e != nil {
		return e
	}
	if len(l.SSH.PublicKeys) != 1 {
		return errors.New("LinuxProfile.PublicKeys requires only 1 SSH Key")
	}
	if e := validateName(l.SSH.PublicKeys[0].KeyData, "LinuxProfile.PublicKeys.KeyData"); e != nil {
		return e
	}
	return nil
}

// Validate implements APIObject
func (a *AcsCluster) Validate() error {
	if e := a.OrchestratorProfile.Validate(); e != nil {
		return e
	}
	if e := a.MasterProfile.Validate(); e != nil {
		return e
	}
	if e := validateUniqueProfileNames(a.AgentPoolProfiles); e != nil {
		return e
	}
	for _, agentPoolProfile := range a.AgentPoolProfiles {
		if e := agentPoolProfile.Validate(); e != nil {
			return e
		}
		if a.OrchestratorProfile.OrchestratorType == Swarm && agentPoolProfile.IsStateful {
			return errors.New("stateful deployments are not supported with Swarm, please let us know if you want this feature")
		}
		if a.OrchestratorProfile.OrchestratorType == Kubernetes && !agentPoolProfile.IsStateful {
			return errors.New("stateless (VMSS) deployments are not supported with Kubernetes, Kubernetes requires the ability to attach/detach disks.  To fix specify 'isStateful=true'")
		}
		if a.OrchestratorProfile.OrchestratorType == Kubernetes && len(agentPoolProfile.DNSPrefix) > 0 {
			return errors.New("DNSPrefix not support for agent pools in Kubernetes - Kubernetes marks its own clusters public")
		}
	}
	if e := a.LinuxProfile.Validate(); e != nil {
		return e
	}
	if e := validateVNET(a); e != nil {
		return e
	}
	return nil
}

func validateName(name string, label string) error {
	if name == "" {
		return fmt.Errorf("%s must be a non-empty value", label)
	}
	return nil
}

func validatePoolName(poolName string) error {
	// we will cap at length of 12 and all lowercase letters since this makes up the VMName
	poolNameRegex := `^([a-z][a-z0-9]{0,11})$`
	re, err := regexp.Compile(poolNameRegex)
	if err != nil {
		return err
	}
	submatches := re.FindStringSubmatch(poolName)
	if len(submatches) != 2 {
		return fmt.Errorf("pool name '%s' is invalid. A pool name must start with a lowercase letter, have max length of 12, and only have characters a-z0-9", poolName)
	}
	return nil
}

func validateDNSName(dnsName string) error {
	dnsNameRegex := `^([a-z][a-z0-9-]{1,13}[a-z0-9])$`
	re, err := regexp.Compile(dnsNameRegex)
	if err != nil {
		return err
	}
	submatches := re.FindStringSubmatch(dnsName)
	if len(submatches) != 2 {
		return fmt.Errorf("DNS name '%s' is invalid. The DNS name must contain between 3 and 15 characters.  The name can contain only letters, numbers, and hyphens.  The name must start with a letter and must end with a letter or a number", dnsName)
	}
	return nil
}

func validateUniqueProfileNames(profiles []AgentPoolProfile) error {
	profileNames := make(map[string]bool)
	for _, profile := range profiles {
		if _, ok := profileNames[profile.Name]; ok {
			return fmt.Errorf("profile name '%s' already exists, profile names must be unique across pools", profile.Name)
		}
		profileNames[profile.Name] = true
	}
	return nil
}

func validateUniquePorts(ports []int, name string) error {
	portMap := make(map[int]bool)
	for _, port := range ports {
		if _, ok := portMap[port]; ok {
			return fmt.Errorf("agent profile '%s' has duplicate port '%d', ports must be unique", name, port)
		}
		portMap[port] = true
	}
	return nil
}

func validateVNET(a *AcsCluster) error {
	isCustomVNET := a.MasterProfile.IsCustomVNET()
	for _, agentPool := range a.AgentPoolProfiles {
		if agentPool.IsCustomVNET() != isCustomVNET {
			return fmt.Errorf("Multiple VNET Subnet configurations specified.  The master profile and each agent pool profile must all specify a custom VNET Subnet, or none at all.")
		}
	}
	if isCustomVNET {
		if a.OrchestratorProfile.OrchestratorType == Swarm {
			return errors.New("bring your own VNET is not supported with Swarm, please let us know if you want this feature")
		}
		subscription, resourcegroup, vnetname, _, e := GetVNETSubnetIDComponents(a.MasterProfile.VnetSubnetID)
		if e != nil {
			return e
		}

		for _, agentPool := range a.AgentPoolProfiles {
			agentSubID, agentRG, agentVNET, _, err := GetVNETSubnetIDComponents(agentPool.VnetSubnetID)
			if err != nil {
				return err
			}
			if agentSubID != subscription ||
				agentRG != resourcegroup ||
				agentVNET != vnetname {
				return errors.New("Multipe VNETS specified.  The master profile and each agent pool must reference the same VNET (but it is ok to reference different subnets on that VNET)")
			}
		}

		masterFirstIP := net.ParseIP(a.MasterProfile.FirstConsecutiveStaticIP)
		if masterFirstIP != nil {
			return fmt.Errorf("MasterProfile.FirstConsecutiveStaticIP (with VNET Subnet specification) '%s' is an invalid IP address", a.MasterProfile.FirstConsecutiveStaticIP)
		}
	}
	return nil
}

// GetVNETSubnetIDComponents extract subscription, resourcegroup, vnetname, subnetname from the vnetSubnetID
func GetVNETSubnetIDComponents(vnetSubnetID string) (string, string, string, string, error) {
	vnetSubnetIDRegex := `^\/subscriptions\/([^\/]*)\/resourceGroups\/([^\/]*)\/providers\/Microsoft.Network\/virtualNetworks\/([^\/]*)\/subnets\/([^\/]*)$`
	re, err := regexp.Compile(vnetSubnetIDRegex)
	if err != nil {
		return "", "", "", "", err
	}
	submatches := re.FindStringSubmatch(vnetSubnetID)
	if len(submatches) != 4 {
		return "", "", "", "", err
	}
	return submatches[1], submatches[2], submatches[3], submatches[4], nil
}

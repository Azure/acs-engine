package vlabs

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

// Validate implements APIObject
func (o *OrchestratorProfile) Validate() error {
	switch o.OrchestratorType {
	case DCOS:
	case DCOS184:
	case DCOS173:
	case SWARM:
	default:
		return fmt.Errorf("OrchestratorProfile has unknown orchestrator: %s", o.OrchestratorType)
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

func parseCIDR(cidr string) (octet1 int, octet2 int, octet3 int, octet4 int, subnet int, err error) {
	// verify cidr format and a /24 subnet
	// regular expression inspired by http://blog.markhatton.co.uk/2011/03/15/regular-expressions-for-ip-addresses-cidr-ranges-and-hostnames/
	cidrRegex := `^((?:[0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5]))\.((?:[0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5]))\.((?:[0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5]))\.([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\/((?:[0-9]|[1-2][0-9]|3[0-2]))$`
	var re *regexp.Regexp
	if re, err = regexp.Compile(cidrRegex); err != nil {
		return 0, 0, 0, 0, 0, err
	}
	submatches := re.FindStringSubmatch(cidr)
	if len(submatches) != 6 {
		return 0, 0, 0, 0, 0, fmt.Errorf("address %s is not specified as valid cidr", cidr)
	}
	if octet1, err = strconv.Atoi(submatches[1]); err != nil {
		return 0, 0, 0, 0, 0, err
	}
	if octet2, err = strconv.Atoi(submatches[2]); err != nil {
		return 0, 0, 0, 0, 0, err
	}
	if octet3, err = strconv.Atoi(submatches[3]); err != nil {
		return 0, 0, 0, 0, 0, err
	}
	if octet4, err = strconv.Atoi(submatches[4]); err != nil {
		return 0, 0, 0, 0, 0, err
	}
	if subnet, err = strconv.Atoi(submatches[5]); err != nil {
		return 0, 0, 0, 0, 0, err
	}
	return octet1, octet2, octet3, octet4, subnet, nil
}

func parseIP(ipaddress string) (octet1 int, octet2 int, octet3 int, octet4 int, err error) {
	// verify cidr format and a /24 subnet
	// regular expression inspired by http://blog.markhatton.co.uk/2011/03/15/regular-expressions-for-ip-addresses-cidr-ranges-and-hostnames/
	ipRegex := `^((?:[0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5]))\.((?:[0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5]))\.((?:[0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5]))\.([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`
	var re *regexp.Regexp
	if re, err = regexp.Compile(ipRegex); err != nil {
		return 0, 0, 0, 0, err
	}
	submatches := re.FindStringSubmatch(ipaddress)
	if len(submatches) != 5 {
		return 0, 0, 0, 0, fmt.Errorf("address %s is not specified as a valid ip address", ipaddress)
	}
	if octet1, err = strconv.Atoi(submatches[1]); err != nil {
		return 0, 0, 0, 0, err
	}
	if octet2, err = strconv.Atoi(submatches[2]); err != nil {
		return 0, 0, 0, 0, err
	}
	if octet3, err = strconv.Atoi(submatches[3]); err != nil {
		return 0, 0, 0, 0, err
	}
	if octet4, err = strconv.Atoi(submatches[4]); err != nil {
		return 0, 0, 0, 0, err
	}
	return octet1, octet2, octet3, octet4, nil
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
		if a.OrchestratorProfile.OrchestratorType == SWARM {
			return errors.New("bring your own VNET is not supported with SWARM")
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

		// validate that the first master IP address has been set
		if e = validateName(a.MasterProfile.FirstConsecutiveStaticIP, "MasterProfile.FirstConsecutiveStaticIP (with VNET Subnet specification)"); e != nil {
			return e
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

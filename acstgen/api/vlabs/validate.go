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
	_, _, _, _, subnet, err := parseCIDR(m.Subnet)
	if err != nil {
		return fmt.Errorf("MasterProfile.Subnet must be specified in cidr format, it returns error: %s", err.Error())
	}
	if subnet != 24 {
		return fmt.Errorf("MasterProfile.Subnet must have a /24 subnet")
	}

	return nil
}

// Validate implements APIObject
func (a *AgentPoolProfiles) Validate() error {
	if e := validateName(a.Name, "AgentPoolProfiles.Name"); e != nil {
		return e
	}
	if a.Count < MinAgentCount || a.Count > MaxAgentCount {
		return fmt.Errorf("AgentPoolProfiles count needs to be in the range [%d,%d]", MinAgentCount, MaxAgentCount)
	}
	if e := validateName(a.VMSize, "AgentPoolProfiles.VMSize"); e != nil {
		return e
	}
	if len(a.Ports) > 0 {
		for _, port := range a.Ports {
			if port < MinPort || port > MaxPort {
				return fmt.Errorf("AgentPoolProfiles Ports must be in the range[%d, %d]", MinPort, MaxPort)
			}
		}
		if e := validateName(a.Name, "AgentPoolProfiles.DNSPrefix when specifying AgentPoolProfiles Ports"); e != nil {
			return e
		}
	}
	return nil
}

// Validate implements APIObject
func (l *LinuxProfile) Validate() error {
	if e := validateName(l.AdminUsername, "LinuxProfile.AdminUsername"); e != nil {
		return e
	}
	if len(l.SSH.PublicKeys) != 1 {
		return errors.New("LinuxProfile.PublicKeys requires 1 SSH Key")
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
	for _, agentPoolProfile := range a.AgentPoolProfiles {
		if e := agentPoolProfile.Validate(); e != nil {
			return e
		}
	}
	if e := a.LinuxProfile.Validate(); e != nil {
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

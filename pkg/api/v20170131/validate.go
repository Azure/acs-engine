package v20170131

import (
	"errors"
	"fmt"
	"regexp"
)

// Validate implements APIObject
func (o *OrchestratorProfile) Validate() error {
	switch o.OrchestratorType {
	case DCOS:
	case Mesos:
	case Swarm:
	case Kubernetes:
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
	if e := validateDNSName(m.DNSPrefix); e != nil {
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
func (a *Properties) Validate() error {
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
		if agentPoolProfile.OSType == Windows {
			switch a.OrchestratorProfile.OrchestratorType {
			case Swarm:
			case Kubernetes:
			default:
				return fmt.Errorf("Orchestrator %s does not support Windows", a.OrchestratorProfile.OrchestratorType)
			}
			if len(a.WindowsProfile.AdminUsername) == 0 {
				return fmt.Errorf("WindowsProfile.AdminUsername must not be empty since agent pool '%s' specifies windows", agentPoolProfile.Name)
			}
			if len(a.WindowsProfile.AdminPassword) == 0 {
				return fmt.Errorf("WindowsProfile.AdminPassword must not be empty since  agent pool '%s' specifies windows", agentPoolProfile.Name)
			}
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
	dnsNameRegex := `^([A-Za-z][A-Za-z0-9-]{1,43}[A-Za-z0-9])$`
	re, err := regexp.Compile(dnsNameRegex)
	if err != nil {
		return err
	}
	if !re.MatchString(dnsName) {
		return fmt.Errorf("DNS name '%s' is invalid. The DNS name must contain between 3 and 45 characters.  The name can contain only letters, numbers, and hyphens.  The name must start with a letter and must end with a letter or a number. (length was %d)", dnsName, len(dnsName))
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

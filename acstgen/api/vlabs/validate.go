package vlabs

import (
	"errors"
	"fmt"
)

// Validate implements APIObject
func (o *OrchestratorProfile) Validate() error {
	switch o.OrchestratorType {
	case DCOS:
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

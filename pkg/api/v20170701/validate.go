package v20170701

import (
	"errors"
	"fmt"
	"net"
	"regexp"
)

// Validate implements APIObject
func (o *OrchestratorProfile) Validate() error {
	switch o.OrchestratorType {
	case Swarm:
	case DCOS:
		switch o.OrchestratorVersion {
		case DCOS187:
		case DCOS188:
		case DCOS190:
		case "":
		default:
			return fmt.Errorf("OrchestratorProfile has unknown orchestrator version: %s \n", o.OrchestratorVersion)
		}
	case DockerCE:
	case Kubernetes:
		switch o.OrchestratorVersion {
		case Kubernetes166:
		case Kubernetes157:
		case "":
		default:
			return fmt.Errorf("OrchestratorProfile has unknown orchestrator version: %s \n", o.OrchestratorVersion)
		}

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
	if e := validateName(m.VMSize, "MasterProfile.VMSize"); e != nil {
		return e
	}
	if m.OSDiskSizeGB != 0 && (m.OSDiskSizeGB < MinDiskSizeGB || m.OSDiskSizeGB > MaxDiskSizeGB) {
		return fmt.Errorf("Invalid master os disk size of %d specified.  The range of valid values are [%d, %d]", m.OSDiskSizeGB, MinDiskSizeGB, MaxDiskSizeGB)
	}
	if e := validateStorageProfile(m.StorageProfile); e != nil {
		return e
	}
	return nil
}

// Validate implements APIObject
func (a *AgentPoolProfile) Validate(orchestratorType OrchestratorType) error {
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
	// Kubernetes don't allow agent DNSPrefix
	if orchestratorType == Kubernetes {
		// The line below need to be removed after June 2017
		a.DNSPrefix = ""
		if e := validateNameEmpty(a.DNSPrefix, "AgentPoolProfile.DNSPrefix"); e != nil {
			return e
		}
	}
	if a.OSDiskSizeGB != 0 && (a.OSDiskSizeGB < MinDiskSizeGB || a.OSDiskSizeGB > MaxDiskSizeGB) {
		return fmt.Errorf("Invalid os disk size of %d specified.  The range of valid values are [%d, %d]", a.OSDiskSizeGB, MinDiskSizeGB, MaxDiskSizeGB)
	}
	if a.DNSPrefix != "" {
		if e := validateDNSName(a.DNSPrefix); e != nil {
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
		} else {
			a.Ports = []int{80, 443, 8080}
		}
	} else {
		if len(a.Ports) > 0 {
			return fmt.Errorf("AgentPoolProfile.Ports must be empty when AgentPoolProfile.DNSPrefix is empty")
		}
	}
	if e := validateStorageProfile(a.StorageProfile); e != nil {
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
	if a.OrchestratorProfile.OrchestratorType == Kubernetes && len(a.ServicePrincipalProfile.ClientID) == 0 {
		return fmt.Errorf("the service principal client ID must be specified with Orchestrator %s", a.OrchestratorProfile.OrchestratorType)
	}

	if a.OrchestratorProfile.OrchestratorType == Kubernetes && len(a.ServicePrincipalProfile.Secret) == 0 {
		return fmt.Errorf("the service principal client secrect must be specified with Orchestrator %s", a.OrchestratorProfile.OrchestratorType)
	}

	for _, agentPoolProfile := range a.AgentPoolProfiles {
		if e := agentPoolProfile.Validate(a.OrchestratorProfile.OrchestratorType); e != nil {
			return e
		}
		switch agentPoolProfile.StorageProfile {
		case StorageAccount:
		case ManagedDisks:
		case "":
		default:
			{
				return fmt.Errorf("unknown storage type '%s' for agent pool '%s'.  Specify either %s, or %s", agentPoolProfile.StorageProfile, agentPoolProfile.Name, StorageAccount, ManagedDisks)
			}
		}

		if agentPoolProfile.OSType == Windows {
			switch a.OrchestratorProfile.OrchestratorType {
			case Kubernetes:
			default:
				return fmt.Errorf("Orchestrator %s does not support Windows", a.OrchestratorProfile.OrchestratorType)
			}
			if a.WindowsProfile == nil {
				return fmt.Errorf("WindowsProfile must not be empty since agent pool '%s' specifies windows", agentPoolProfile.Name)
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
	if e := validateVNET(a); e != nil {
		return e
	}
	return nil
}

func validateNameEmpty(name string, label string) error {
	if name != "" {
		return fmt.Errorf("%s must be an empty value", label)
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

func validateUniqueProfileNames(profiles []*AgentPoolProfile) error {
	profileNames := make(map[string]bool)
	for _, profile := range profiles {
		if _, ok := profileNames[profile.Name]; ok {
			return fmt.Errorf("profile name '%s' already exists, profile names must be unique across pools", profile.Name)
		}
		profileNames[profile.Name] = true
	}
	return nil
}

func validateStorageProfile(storageProfile string) error {
	switch storageProfile {
	case StorageAccount:
	case ManagedDisks:
	case "":
	default:
		{
			return fmt.Errorf("Unknown storageProfile '%s'. Specify either %s or %s", storageProfile, StorageAccount, ManagedDisks)
		}
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

func validateVNET(a *Properties) error {
	isCustomVNET := a.MasterProfile.IsCustomVNET()
	for _, agentPool := range a.AgentPoolProfiles {
		if agentPool.IsCustomVNET() != isCustomVNET {
			return fmt.Errorf("Multiple VNET Subnet configurations specified.  The master profile and each agent pool profile must all specify a custom VNET Subnet, or none at all.")
		}
	}
	if isCustomVNET {
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
		if masterFirstIP == nil {
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

package v20170701

import (
	"errors"
	"fmt"
	"net"
	"regexp"
	"strings"

	validator "gopkg.in/go-playground/validator.v9"
)

var (
	validate                *validator.Validate
	keyvaultSecretPathRegex *regexp.Regexp
)

func init() {
	validate = validator.New()
	keyvaultSecretPathRegex = regexp.MustCompile(`^(/subscriptions/\S+/resourceGroups/\S+/providers/Microsoft.KeyVault/vaults/\S+)/secrets/([^/\s]+)(/(\S+))?$`)
}

// Validate implements APIObject
func (o *OrchestratorProfile) Validate() error {
	// Don't need to call validate.Struct(o)
	// It is handled by Properties.Validate()
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
	// Don't need to call validate.Struct(m)
	// It is handled by Properties.Validate()
	if e := validateDNSName(m.DNSPrefix); e != nil {
		return e
	}
	return nil
}

// Validate implements APIObject
func (a *AgentPoolProfile) Validate(orchestratorType OrchestratorType) error {
	// Don't need to call validate.Struct(a)
	// It is handled by Properties.Validate()
	if e := validatePoolName(a.Name); e != nil {
		return e
	}
	// Kubernetes don't allow agent DNSPrefix and ports
	if orchestratorType == Kubernetes {
		// The two lines below need to be removed after August 2017
		a.DNSPrefix = ""
		a.Ports = []int{}
		if e := validate.Var(a.DNSPrefix, "len=0"); e != nil {
			return fmt.Errorf("AgentPoolProfile.DNSPrefix must be empty for Kubernetes")
		}
		if e := validate.Var(a.Ports, "len=0"); e != nil {
			return fmt.Errorf("AgentPoolProfile.Ports must be empty for Kubernetes")
		}
	}
	if a.DNSPrefix != "" {
		if e := validateDNSName(a.DNSPrefix); e != nil {
			return e
		}
		if len(a.Ports) > 0 {
			if e := validateUniquePorts(a.Ports, a.Name); e != nil {
				return e
			}
		} else {
			a.Ports = []int{80, 443, 8080}
		}
	} else {
		if e := validate.Var(a.Ports, "len=0"); e != nil {
			return fmt.Errorf("AgentPoolProfile.Ports must be empty when AgentPoolProfile.DNSPrefix is empty for Orchestrator: %s", string(orchestratorType))
		}
	}
	return nil
}

// Validate implements APIObject
func (l *LinuxProfile) Validate() error {
	// Don't need to call validate.Struct(l)
	// It is handled by Properties.Validate()
	if e := validate.Var(l.SSH.PublicKeys[0].KeyData, "required"); e != nil {
		return fmt.Errorf("KeyData in LinuxProfile.SSH.PublicKeys cannot be empty string")
	}
	return nil
}

func handleValidationErrors(e validator.ValidationErrors) error {
	err := e[0]
	ns := err.Namespace()
	switch ns {
	case "Properties.OrchestratorProfile", "Properties.MasterProfile",
		"Properties.MasterProfile.DNSPrefix", "Properties.MasterProfile.VMSize",
		"Properties.LinuxProfile", "Properties.ServicePrincipalProfile.ClientID",
		"Properties.ServicePrincipalProfile.Secret", "Properties.WindowsProfile.AdminUsername",
		"Properties.WindowsProfile.AdminPassword":
		return fmt.Errorf("missing %s", ns)
	case "Properties.MasterProfile.Count":
		return fmt.Errorf("MasterProfile count needs to be 1, 3, or 5")
	case "Properties.MasterProfile.OSDiskSizeGB":
		return fmt.Errorf("Invalid os disk size of %d specified.  The range of valid values are [%d, %d]", err.Value().(int), MinDiskSizeGB, MaxDiskSizeGB)
	case "Properties.MasterProfile.StorageProfile":
		return fmt.Errorf("Unknown storageProfile '%s'. Specify either %s or %s", err.Value().(string), StorageAccount, ManagedDisks)
	default:
		if strings.HasPrefix(ns, "Properties.AgentPoolProfiles") {
			switch {
			case strings.HasSuffix(ns, ".Name") || strings.HasSuffix(ns, "VMSize"):
				return fmt.Errorf("missing %s", ns)
			case strings.HasSuffix(ns, ".Count"):
				return fmt.Errorf("AgentPoolProfile count needs to be in the range [%d,%d]", MinAgentCount, MaxAgentCount)
			case strings.HasSuffix(ns, ".OSDiskSizeGB"):
				return fmt.Errorf("Invalid os disk size of %d specified.  The range of valid values are [%d, %d]", err.Value().(int), MinDiskSizeGB, MaxDiskSizeGB)
			case strings.Contains(ns, ".Ports"):
				return fmt.Errorf("AgentPoolProfile Ports must be in the range[%d, %d]", MinPort, MaxPort)
			case strings.HasSuffix(ns, ".StorageProfile"):
				return fmt.Errorf("Unknown storageProfile '%s'. Specify either %s or %s", err.Value().(string), StorageAccount, ManagedDisks)
			default:
				break
			}
		}
	}
	return fmt.Errorf("Namespace %s is not caught, %+v", ns, e)
}

// Validate implements APIObject
func (a *Properties) Validate() error {
	if e := validate.Struct(a); e != nil {
		return handleValidationErrors(e.(validator.ValidationErrors))
	}
	if e := a.OrchestratorProfile.Validate(); e != nil {
		return e
	}
	if e := a.MasterProfile.Validate(); e != nil {
		return e
	}
	if e := validateUniqueProfileNames(a.AgentPoolProfiles); e != nil {
		return e
	}

	if a.OrchestratorProfile.OrchestratorType == Kubernetes {
		if (len(a.ServicePrincipalProfile.Secret) == 0 && len(a.ServicePrincipalProfile.KeyvaultSecretRef) == 0) ||
			(len(a.ServicePrincipalProfile.Secret) != 0 && len(a.ServicePrincipalProfile.KeyvaultSecretRef) != 0) {
			return fmt.Errorf("either the service principal client secret or keyvault secret reference must be specified with Orchestrator %s", a.OrchestratorProfile.OrchestratorType)
		}

		if len(a.ServicePrincipalProfile.KeyvaultSecretRef) != 0 {
			parts := keyvaultSecretPathRegex.FindStringSubmatch(a.ServicePrincipalProfile.KeyvaultSecretRef)
			if len(parts) != 5 {
				return fmt.Errorf("service principal client keyvault secret reference is of incorrect format")
			}
		}
	}

	for _, agentPoolProfile := range a.AgentPoolProfiles {
		if e := agentPoolProfile.Validate(a.OrchestratorProfile.OrchestratorType); e != nil {
			return e
		}

		if agentPoolProfile.OSType == Windows {
			if a.WindowsProfile == nil {
				return fmt.Errorf("missing WindowsProfile")
			}
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

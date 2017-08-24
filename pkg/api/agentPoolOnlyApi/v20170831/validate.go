package v20170831

import (
	"fmt"
	"regexp"
	"strings"

	validator "gopkg.in/go-playground/validator.v9"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// Validate implements APIObject
func (a *AgentPoolProfile) Validate() error {
	// Don't need to call validate.Struct(a)
	// It is handled by Properties.Validate()
	if e := validatePoolName(a.Name); e != nil {
		return e
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
	// TODO: Add more validation here
	case "Properties.LinuxProfile", "Properties.ServicePrincipalProfile.ClientID",
		"Properties.ServicePrincipalProfile.Secret", "Properties.WindowsProfile.AdminUsername",
		"Properties.WindowsProfile.AdminPassword":
		return fmt.Errorf("missing %s", ns)
	default:
		if strings.HasPrefix(ns, "Properties.AgentPoolProfiles") {
			switch {
			case strings.HasSuffix(ns, ".Name") || strings.HasSuffix(ns, "VMSize"):
				return fmt.Errorf("missing %s", ns)
			case strings.HasSuffix(ns, ".Count"):
				return fmt.Errorf("AgentPoolProfile count needs to be in the range [%d,%d]", MinAgentCount, MaxAgentCount)
			case strings.HasSuffix(ns, ".OSDiskSizeGB"):
				return fmt.Errorf("Invalid os disk size of %d specified.  The range of valid values are [%d, %d]", err.Value().(int), MinDiskSizeGB, MaxDiskSizeGB)
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

	// Don't need to call validate.Struct(m)
	// It is handled by Properties.Validate()
	if e := validateDNSName(a.DNSPrefix); e != nil {
		return e
	}

	if e := validateUniqueProfileNames(a.AgentPoolProfiles); e != nil {
		return e
	}

	if a.ServicePrincipalProfile == nil {
		return fmt.Errorf("missing ServicePrincipalProfile")
	}

	for _, agentPoolProfile := range a.AgentPoolProfiles {
		if agentPoolProfile.OSType == Windows {
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

func validateVNET(a *Properties) error {
	var customVNETCount int
	var isCustomVNET bool
	for _, agentPool := range a.AgentPoolProfiles {
		if agentPool.IsCustomVNET() {
			customVNETCount++
			isCustomVNET = agentPool.IsCustomVNET()
		}
	}

	if !(customVNETCount == 0 || customVNETCount == len(a.AgentPoolProfiles)) {
		return fmt.Errorf("Multiple VNET Subnet configurations specified.  Each agent pool profile must all specify a custom VNET Subnet, or none at all")
	}

	subIDMap := make(map[string]int)
	resourceGroupMap := make(map[string]int)
	agentVNETMap := make(map[string]int)
	if isCustomVNET {
		for _, agentPool := range a.AgentPoolProfiles {
			agentSubID, agentRG, agentVNET, _, err := GetVNETSubnetIDComponents(agentPool.VnetSubnetID)
			if err != nil {
				return err
			}

			subIDMap[agentSubID] = subIDMap[agentSubID] + 1
			resourceGroupMap[agentRG] = resourceGroupMap[agentRG] + 1
			agentVNETMap[agentVNET] = agentVNETMap[agentVNET] + 1
		}

		// TODO: Add more validation to ensure all agent pools belong to the same VNET, subscription, and resource group
		// 	if(len(subIDMap) != len(a.AgentPoolProfiles))

		// 	return errors.New("Multiple VNETS specified.  Each agent pool must reference the same VNET (but it is ok to reference different subnets on that VNET)")
		// }
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

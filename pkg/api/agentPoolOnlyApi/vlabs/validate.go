package vlabs

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/Azure/acs-engine/pkg/api/common"
	validator "gopkg.in/go-playground/validator.v9"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// Validate implements APIObject
func (a *AgentPoolProfile) Validate() error {
	if a.ImageRef != nil {
		if err := validateImageNameAndGroup(a.ImageRef.Name, a.ImageRef.ResourceGroup); err != nil {
			return err
		}
	}
	// Don't need to call validate.Struct(a)
	// It is handled by Properties.Validate()
	return validatePoolName(a.Name)
}

func validateImageNameAndGroup(name, resourceGroup string) error {
	if name == "" && resourceGroup != "" {
		return errors.New("imageName needs to be specified when imageResourceGroup is provided")
	}
	if name != "" && resourceGroup == "" {
		return errors.New("imageResourceGroup needs to be specified when imageName is provided")
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
	if e := common.ValidateDNSPrefix(a.DNSPrefix); e != nil {
		return e
	}

	if e := validateUniqueProfileNames(a.AgentPoolProfiles); e != nil {
		return e
	}

	if e := validateAgents(a.OrchestratorProfile, a.AgentPoolProfiles); e != nil {
		return e
	}

	if e := validateCertificateProfile(a.OrchestratorProfile, a.CertificateProfile); e != nil {
		return e
	}

	if e := a.LinuxProfile.Validate(); e != nil {
		return e
	}
	return validateVNET(a)
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

func validateAgents(orchestratorProfile *OrchestratorProfile, profiles []*AgentPoolProfile) error {
	orchestratorType := common.Kubernetes
	if orchestratorProfile != nil {
		orchestratorType = orchestratorProfile.OrchestratorType
	}

	profileNames := make(map[string]bool)
	for _, agentPoolProfile := range profiles {
		// validate that each AgentPoolProfile Name is unique
		if _, ok := profileNames[agentPoolProfile.Name]; ok {
			return fmt.Errorf("profile name '%s' already exists, profile names must be unique across pools", agentPoolProfile.Name)
		}
		profileNames[agentPoolProfile.Name] = true
		if err := agentPoolProfile.Validate(); err != nil {
			return err
		}
		if err := validateRoles(orchestratorType, agentPoolProfile.Role); err != nil {
			return err
		}
		if err := validateOpenShiftAgent(orchestratorType, agentPoolProfile); err != nil {
			return err
		}
	}
	if orchestratorType == common.OpenShift {
		if !reflect.DeepEqual(profileNames, map[string]bool{"compute": true, "infra": true}) {
			return fmt.Errorf("OpenShift requires exactly two agent pool profiles: compute and infra")
		}
	}
	return nil
}

func validateRoles(orchestratorType string, role AgentPoolProfileRole) error {
	validRoles := []AgentPoolProfileRole{AgentPoolProfileRoleEmpty}
	if orchestratorType == common.OpenShift {
		validRoles = append(validRoles, AgentPoolProfileRoleInfra)
	}
	var found bool
	for _, validRole := range validRoles {
		if role == validRole {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("role %q is not supported by orchestrator %q", role, orchestratorType)
	}
	return nil
}

func validateOpenShiftAgent(orchestratorType string, a *AgentPoolProfile) error {
	if orchestratorType != common.OpenShift {
		return nil
	}
	if a.AvailabilityProfile != common.AvailabilitySet {
		return fmt.Errorf("only AvailabilityProfile: AvailabilitySet is supported for Orchestrator 'OpenShift'")
	}
	if (a.Name == "infra") != (a.Role == "infra") {
		return errors.New("OpenShift requires that the 'infra' agent pool profile, and no other, should have role 'infra'")
	}
	return nil
}

func validateCertificateProfile(orchestratorProfile *OrchestratorProfile, certificateProfile *CertificateProfile) error {
	if certificateProfile == nil {
		return errors.New("certificateProfile is required")
	}
	if orchestratorProfile != nil && orchestratorProfile.OrchestratorType == common.OpenShift {
		// Invalidate missing master CA cert and key
		if certificateProfile.CaCertificate == "" {
			return errors.New("master CA certificate is required")
		}
		if certificateProfile.CaPrivateKey == "" {
			return errors.New("master CA private key is required")
		}
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
			agentSubID, agentRG, agentVNET, _, err := common.GetVNETSubnetIDComponents(agentPool.VnetSubnetID)
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

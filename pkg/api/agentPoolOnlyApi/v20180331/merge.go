package v20180331

import (
	"github.com/imdario/mergo"
	"github.com/pkg/errors"
)

// Merge existing ManagedCluster attribute into mc
func (mc *ManagedCluster) Merge(emc *ManagedCluster) error {

	// Merge properties.dnsPrefix
	if emc.Properties.DNSPrefix == "" {
		return errors.New("existing ManagedCluster expect properties.dnsPrefix not to be empty")
	}

	if mc.Properties.DNSPrefix == "" {
		mc.Properties.DNSPrefix = emc.Properties.DNSPrefix
	} else if mc.Properties.DNSPrefix != emc.Properties.DNSPrefix {
		return errors.Errorf("change dnsPrefix from %s to %s is not supported",
			emc.Properties.DNSPrefix,
			mc.Properties.DNSPrefix)
	}

	// Merge Properties.AgentPoolProfiles
	if mc.Properties.AgentPoolProfiles == nil {
		mc.Properties.AgentPoolProfiles = emc.Properties.AgentPoolProfiles
	}
	// Merge Properties.LinuxProfile
	if mc.Properties.LinuxProfile == nil {
		mc.Properties.LinuxProfile = emc.Properties.LinuxProfile
	}
	// Merge Properties.WindowsProfile
	if mc.Properties.WindowsProfile == nil {
		mc.Properties.WindowsProfile = emc.Properties.WindowsProfile
	}
	// Merge Properties.ServicePrincipalProfile
	if mc.Properties.ServicePrincipalProfile == nil {
		mc.Properties.ServicePrincipalProfile = emc.Properties.ServicePrincipalProfile
	}
	// Merge Properties.KubernetesVersion
	if mc.Properties.KubernetesVersion == "" {
		mc.Properties.KubernetesVersion = emc.Properties.KubernetesVersion
	}

	// Merge properties.enableRBAC
	if emc.Properties.EnableRBAC == nil {
		return errors.New("existing ManagedCluster expect properties.enableRBAC not to be nil")
	}

	if mc.Properties.EnableRBAC == nil {
		// For update scenario, the default behavior is to use existing behavior
		mc.Properties.EnableRBAC = emc.Properties.EnableRBAC
	} else if *mc.Properties.EnableRBAC != *emc.Properties.EnableRBAC {
		return errors.Errorf("existing ManagedCluster has properties.enableRBAC %v. update to %v is not supported",
			*emc.Properties.EnableRBAC,
			*mc.Properties.EnableRBAC)
	}
	if mc.Properties.NetworkProfile == nil {
		// For update scenario, the default behavior is to use existing behavior
		mc.Properties.NetworkProfile = emc.Properties.NetworkProfile
	}

	if emc.Properties.AADProfile != nil {
		if mc.Properties.AADProfile == nil {
			mc.Properties.AADProfile = &AADProfile{}
		}
		if err := mergo.Merge(mc.Properties.AADProfile,
			*emc.Properties.AADProfile); err != nil {
			return err
		}
	}
	return nil
}

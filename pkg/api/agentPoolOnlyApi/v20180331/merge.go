package v20180331

import (
	"fmt"
)

// Merge existing ManagedCluster attribute into mc
func (mc *ManagedCluster) Merge(emc *ManagedCluster) error {

	// Merge properties.dnsPrefix
	if emc.Properties.DNSPrefix == "" {
		return fmt.Errorf("existing ManagedCluster expect properties.dnsPrefix not to be empty")
	}

	if mc.Properties.DNSPrefix == "" {
		mc.Properties.DNSPrefix = emc.Properties.DNSPrefix
	} else if mc.Properties.DNSPrefix != emc.Properties.DNSPrefix {
		return fmt.Errorf("change dnsPrefix from %s to %s is not supported",
			emc.Properties.DNSPrefix,
			mc.Properties.DNSPrefix)
	}

	// Merge properties.enableRBAC
	if emc.Properties.EnableRBAC == nil {
		return fmt.Errorf("existing ManagedCluster expect properties.enableRBAC not to be nil")
	}

	if mc.Properties.EnableRBAC == nil {
		// For update scenario, the default behavior is to use existing behavior
		mc.Properties.EnableRBAC = emc.Properties.EnableRBAC
	} else if *mc.Properties.EnableRBAC != *emc.Properties.EnableRBAC {
		return fmt.Errorf("existing ManagedCluster has properties.enableRBAC %v. update to %v is not supported",
			*emc.Properties.EnableRBAC,
			*mc.Properties.EnableRBAC)
	}
	return nil
}

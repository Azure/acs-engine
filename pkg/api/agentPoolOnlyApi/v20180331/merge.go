package v20180331

import (
	"fmt"
)

// Merge existing ManagedCluster attribute into mc
func (mc *ManagedCluster) Merge(emc *ManagedCluster) error {
	if emc.Properties.EnableRBAC == nil {
		return fmt.Errorf("existing ManagedCluster expect properties.enableRBAC not to be nil")
	}

	if mc.Properties.EnableRBAC == nil {
		// For update scenario, the default behavior is to use existing behavior
		mc.Properties.EnableRBAC = emc.Properties.EnableRBAC
	}

	if *mc.Properties.EnableRBAC != *emc.Properties.EnableRBAC {
		return fmt.Errorf("existing ManagedCluster has properties.enableRBAC %v. update to %v is not supported",
			*emc.Properties.EnableRBAC,
			*mc.Properties.EnableRBAC)
	}
	return nil
}

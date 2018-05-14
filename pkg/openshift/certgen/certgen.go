package certgen

import (
	"github.com/Azure/acs-engine/pkg/api"
	"github.com/Azure/acs-engine/pkg/api/common"
	"github.com/Azure/acs-engine/pkg/openshift/certgen/release39"
	"github.com/Azure/acs-engine/pkg/openshift/certgen/unstable"
)

// OpenShiftSetDefaultCerts sets default certificate and configuration properties in the
// openshift orchestrator.
func OpenShiftSetDefaultCerts(a *api.Properties, orchestratorName, clusterID string) (bool, error) {
	switch a.OrchestratorProfile.OrchestratorVersion {
	case common.OpenShiftVersion3Dot9Dot0:
		return release39.OpenShiftSetDefaultCerts(a, orchestratorName, clusterID)
	default:
		return unstable.OpenShiftSetDefaultCerts(a, orchestratorName, clusterID)
	}
}

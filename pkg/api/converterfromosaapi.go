package api

import (
	"github.com/Azure/acs-engine/pkg/api/osa/vlabs"
)

// ConvertVLabsOpenShiftClusterToContainerService converts from a
// vlabs.OpenShiftCluster to a ContainerService.
func ConvertVLabsOpenShiftClusterToContainerService(oc *vlabs.OpenShiftCluster) *ContainerService {
	cs := &ContainerService{
		ID:       oc.ID,
		Location: oc.Location,
		Name:     oc.Name,
		Tags:     oc.Tags,
		Type:     oc.Type,
	}

	if oc.Plan != nil {
		cs.Plan = &ResourcePurchasePlan{
			Name:          oc.Plan.Name,
			Product:       oc.Plan.Product,
			PromotionCode: oc.Plan.PromotionCode,
			Publisher:     oc.Plan.Publisher,
		}
	}

	if oc.Properties != nil {
		cs.Properties = &Properties{
			ProvisioningState: ProvisioningState(oc.Properties.ProvisioningState),
			OrchestratorProfile: &OrchestratorProfile{
				OrchestratorVersion: oc.Properties.OpenShiftVersion,
				OpenShiftConfig: &OpenShiftConfig{
					PublicHostname:         oc.Properties.PublicHostname,
					RoutingConfigSubdomain: oc.Properties.RoutingConfigSubdomain,
					RoutingConfigFQDN:      oc.Properties.RoutingConfigFQDN,
				},
			},
			MasterProfile: &MasterProfile{
				FQDN: oc.Properties.FQDN,
			},
			ServicePrincipalProfile: &ServicePrincipalProfile{
				ClientID: oc.Properties.ServicePrincipalProfile.ClientID,
				Secret:   oc.Properties.ServicePrincipalProfile.Secret,
			},
		}

		cs.Properties.AgentPoolProfiles = make([]*AgentPoolProfile, 0, len(oc.Properties.AgentPoolProfiles))
		for _, app := range oc.Properties.AgentPoolProfiles {
			cs.Properties.AgentPoolProfiles = append(cs.Properties.AgentPoolProfiles,
				&AgentPoolProfile{
					Name:         app.Name,
					Count:        app.Count,
					VMSize:       app.VMSize,
					OSType:       OSType(app.OSType),
					VnetSubnetID: app.VnetSubnetID,
					Role:         AgentPoolProfileRole(app.Role),
				},
			)
		}
	}

	return cs
}

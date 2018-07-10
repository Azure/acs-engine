package vlabs

import (
	"github.com/Azure/acs-engine/pkg/api"
)

// AsContainerService returns an OpenShiftCluster converted to an
// api.ContainerService.
func (oc *OpenShiftCluster) AsContainerService() *api.ContainerService {
	cs := &api.ContainerService{
		ID:       oc.ID,
		Location: oc.Location,
		Name:     oc.Name,
		Tags:     oc.Tags,
		Type:     oc.Type,
		Properties: &api.Properties{
			ProvisioningState: api.ProvisioningState(oc.Properties.ProvisioningState),
			OrchestratorProfile: &api.OrchestratorProfile{
				OpenShiftConfig: &api.OpenShiftConfig{
					OpenShiftVersion:       oc.Properties.OpenShiftVersion,
					PublicHostname:         oc.Properties.PublicHostname,
					RoutingConfigSubdomain: oc.Properties.RoutingConfigSubdomain,
				},
			},
			ServicePrincipalProfile: &api.ServicePrincipalProfile{
				ClientID: oc.Properties.ServicePrincipalProfile.ClientID,
				Secret:   oc.Properties.ServicePrincipalProfile.Secret,
			},
		},
	}

	if oc.Plan != nil {
		cs.Plan = &api.ResourcePurchasePlan{
			Name:          oc.Plan.Name,
			Product:       oc.Plan.Product,
			PromotionCode: oc.Plan.PromotionCode,
			Publisher:     oc.Plan.Publisher,
		}
	}

	cs.Properties.AgentPoolProfiles = make([]*api.AgentPoolProfile, 0, len(oc.Properties.ComputePools)+1)
	for _, app := range oc.Properties.ComputePools {
		cs.Properties.AgentPoolProfiles = append(cs.Properties.AgentPoolProfiles,
			&api.AgentPoolProfile{
				Name:         app.Name,
				Role:         api.AgentPoolProfileRole(AgentPoolProfileRoleEmpty),
				Count:        app.Count,
				VMSize:       app.VMSize,
				VnetSubnetID: app.VnetSubnetID,
			},
		)
	}

	cs.Properties.AgentPoolProfiles = append(cs.Properties.AgentPoolProfiles,
		&api.AgentPoolProfile{
			Name:         oc.Properties.InfraPool.Name,
			Role:         api.AgentPoolProfileRole(AgentPoolProfileRoleInfra),
			Count:        oc.Properties.InfraPool.Count,
			VMSize:       oc.Properties.InfraPool.VMSize,
			VnetSubnetID: oc.Properties.InfraPool.VnetSubnetID,
		},
	)

	return cs
}

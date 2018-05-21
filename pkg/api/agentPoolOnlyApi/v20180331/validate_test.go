package v20180331

import (
	"testing"
)

func TestValidateVNET(t *testing.T) {

	serviceCidr := "10.0.0.0/16"
	serviceCidrBad := "10.0.0.0"
	dNSServiceIP := "10.0.0.10"
	dNSServiceIPBad := "10.0.0.257"
	dNSServiceIPOutOfRange := "10.1.0.1"
	dNSServiceIPFirstInServiceCidr := "10.0.0.1"
	dockerBridgeCidr := "127.17.0.1/16"
	dockerBridgeCidrBad := "127.17.0.1/50"

	vnetSubnetID1 := "/subscriptions/mySubscription/resourceGroups/myResourceGroup/providers/Microsoft.Network/virtualNetworks/myVnet/subnets/mySubnet1"
	vnetSubnetID2 := "/subscriptions/mySubscription/resourceGroups/myResourceGroup/providers/Microsoft.Network/virtualNetworks/myVnet/subnets/mySubnet1"
	vnetSubnetID1Bad := "/subscription/mySubscription/resourceGroups/myResourceGroup/providers/Microsoft.Network/virtualNetworks/myVnet/subnets/mySubnet1"
	vnetSubnetID1WrongSubscription := "/subscriptions/wrongSubscription/resourceGroups/myResourceGroup/providers/Microsoft.Network/virtualNetworks/myVnet/subnets/mySubnet1"
	vnetSubnetID1WrongResourceGroup := "/subscriptions/mySubscription/resourceGroups/wrongResourceGroup/providers/Microsoft.Network/virtualNetworks/myVnet/subnets/mySubnet1"
	vnetSubnetID1WrongVnet := "/subscriptions/mySubscription/resourceGroups/myResourceGroup/providers/Microsoft.Network/virtualNetworks/wrongVnet/subnets/mySubnet1"

	maxPods1 := 20
	maxPods2 := 50
	maxPodsTooSmall := 4

	// all network profile fields have values, should pass
	n := &NetworkProfile{
		NetworkPlugin:    NetworkPlugin("azure"),
		ServiceCidr:      serviceCidr,
		DNSServiceIP:     dNSServiceIP,
		DockerBridgeCidr: dockerBridgeCidr,
	}

	p := []*AgentPoolProfile{
		{
			VnetSubnetID: vnetSubnetID1,
			MaxPods:      &maxPods1,
		},
		{
			VnetSubnetID: vnetSubnetID2,
			MaxPods:      &maxPods2,
		},
	}

	a := &Properties{
		NetworkProfile:    n,
		AgentPoolProfiles: p,
	}

	if err := validateVNET(a); err != nil {
		t.Errorf("Failed to validate VNET: %s", err)
	}

	// no network profile, this is prior v20180331 case, should pass
	p = []*AgentPoolProfile{
		{
			VnetSubnetID: vnetSubnetID1,
			MaxPods:      &maxPods1,
		},
		{
			VnetSubnetID: vnetSubnetID2,
			MaxPods:      &maxPods2,
		},
	}

	a = &Properties{
		AgentPoolProfiles: p,
	}

	if err := validateVNET(a); err != nil {
		t.Errorf("Failed to validate VNET: %s", err)
	}

	// network profile has only NetworkPlugin field, should pass
	n = &NetworkProfile{
		NetworkPlugin: NetworkPlugin("azure"),
	}

	p = []*AgentPoolProfile{
		{
			VnetSubnetID: vnetSubnetID1,
			MaxPods:      &maxPods1,
		},
		{
			VnetSubnetID: vnetSubnetID2,
			MaxPods:      &maxPods2,
		},
	}

	a = &Properties{
		NetworkProfile:    n,
		AgentPoolProfiles: p,
	}

	if err := validateVNET(a); err != nil {
		t.Errorf("Failed to validate VNET: %s", err)
	}

	// network profile has NetworkPlugin and ServiceCidr field, should fail
	n = &NetworkProfile{
		NetworkPlugin: NetworkPlugin("azure"),
		ServiceCidr:   serviceCidr,
	}

	p = []*AgentPoolProfile{
		{
			VnetSubnetID: vnetSubnetID1,
			MaxPods:      &maxPods1,
		},
		{
			VnetSubnetID: vnetSubnetID2,
			MaxPods:      &maxPods2,
		},
	}

	a = &Properties{
		NetworkProfile:    n,
		AgentPoolProfiles: p,
	}

	if err := validateVNET(a); err != ErrorInvalidNetworkProfile {
		t.Errorf("Failed to validate VNET: %s", ErrorInvalidNetworkProfile)
	}

	// NetworkPlugin is not azure or kubenet
	n = &NetworkProfile{
		NetworkPlugin:    NetworkPlugin("none"),
		ServiceCidr:      serviceCidr,
		DNSServiceIP:     dNSServiceIP,
		DockerBridgeCidr: dockerBridgeCidr,
	}

	p = []*AgentPoolProfile{
		{
			VnetSubnetID: vnetSubnetID1,
			MaxPods:      &maxPods1,
		},
		{
			VnetSubnetID: vnetSubnetID2,
			MaxPods:      &maxPods2,
		},
	}

	a = &Properties{
		NetworkProfile:    n,
		AgentPoolProfiles: p,
	}

	if err := validateVNET(a); err != ErrorInvalidNetworkPlugin {
		t.Errorf("Failed to throw error, %s", ErrorInvalidNetworkPlugin)
	}

	// NetworkPlugin = Azure, bad serviceCidr
	n = &NetworkProfile{
		NetworkPlugin:    NetworkPlugin("azure"),
		ServiceCidr:      serviceCidrBad,
		DNSServiceIP:     dNSServiceIP,
		DockerBridgeCidr: dockerBridgeCidr,
	}

	p = []*AgentPoolProfile{
		{
			VnetSubnetID: vnetSubnetID1,
			MaxPods:      &maxPods1,
		},
		{
			VnetSubnetID: vnetSubnetID2,
			MaxPods:      &maxPods2,
		},
	}

	a = &Properties{
		NetworkProfile:    n,
		AgentPoolProfiles: p,
	}

	if err := validateVNET(a); err != ErrorInvalidServiceCidr {
		t.Errorf("Failed to throw error, %s", ErrorInvalidServiceCidr)
	}

	// NetworkPlugin = Azure, bad dNSServiceIP
	n = &NetworkProfile{
		NetworkPlugin:    NetworkPlugin("azure"),
		ServiceCidr:      serviceCidr,
		DNSServiceIP:     dNSServiceIPBad,
		DockerBridgeCidr: dockerBridgeCidr,
	}

	p = []*AgentPoolProfile{
		{
			VnetSubnetID: vnetSubnetID1,
			MaxPods:      &maxPods1,
		},
		{
			VnetSubnetID: vnetSubnetID2,
			MaxPods:      &maxPods2,
		},
	}

	a = &Properties{
		NetworkProfile:    n,
		AgentPoolProfiles: p,
	}

	if err := validateVNET(a); err != ErrorInvalidDNSServiceIP {
		t.Errorf("Failed to throw error, %s", ErrorInvalidDNSServiceIP)
	}

	// NetworkPlugin = Azure, bad dockerBridgeCidr
	n = &NetworkProfile{
		NetworkPlugin:    NetworkPlugin("azure"),
		ServiceCidr:      serviceCidr,
		DNSServiceIP:     dNSServiceIP,
		DockerBridgeCidr: dockerBridgeCidrBad,
	}

	p = []*AgentPoolProfile{
		{
			VnetSubnetID: vnetSubnetID1,
			MaxPods:      &maxPods1,
		},
		{
			VnetSubnetID: vnetSubnetID2,
			MaxPods:      &maxPods2,
		},
	}

	a = &Properties{
		NetworkProfile:    n,
		AgentPoolProfiles: p,
	}

	if err := validateVNET(a); err != ErrorInvalidDockerBridgeCidr {
		t.Errorf("Failed to throw error, %s", ErrorInvalidDockerBridgeCidr)
	}

	// NetworkPlugin = Azure, DNSServiceIP is not within ServiceCidr
	n = &NetworkProfile{
		NetworkPlugin:    NetworkPlugin("azure"),
		ServiceCidr:      serviceCidr,
		DNSServiceIP:     dNSServiceIPOutOfRange,
		DockerBridgeCidr: dockerBridgeCidr,
	}

	p = []*AgentPoolProfile{
		{
			VnetSubnetID: vnetSubnetID1,
			MaxPods:      &maxPods1,
		},
		{
			VnetSubnetID: vnetSubnetID2,
			MaxPods:      &maxPods2,
		},
	}

	a = &Properties{
		NetworkProfile:    n,
		AgentPoolProfiles: p,
	}

	if err := validateVNET(a); err != ErrorDNSServiceIPNotInServiceCidr {
		t.Errorf("Failed to throw error, %s", ErrorDNSServiceIPNotInServiceCidr)
	}

	// NetworkPlugin = Azure, DNSServiceIP is the first IP in ServiceCidr
	n = &NetworkProfile{
		NetworkPlugin:    NetworkPlugin("azure"),
		ServiceCidr:      serviceCidr,
		DNSServiceIP:     dNSServiceIPFirstInServiceCidr,
		DockerBridgeCidr: dockerBridgeCidr,
	}

	p = []*AgentPoolProfile{
		{
			VnetSubnetID: vnetSubnetID1,
			MaxPods:      &maxPods1,
		},
		{
			VnetSubnetID: vnetSubnetID2,
			MaxPods:      &maxPods2,
		},
	}

	a = &Properties{
		NetworkProfile:    n,
		AgentPoolProfiles: p,
	}

	if err := validateVNET(a); err != ErrorDNSServiceIPAlreadyUsed {
		t.Errorf("Failed to throw error, %s", ErrorDNSServiceIPAlreadyUsed)
	}

	// NetworkPlugin = Azure, at least one agent pool does not have subnet defined
	n = &NetworkProfile{
		NetworkPlugin:    NetworkPlugin("azure"),
		ServiceCidr:      serviceCidr,
		DNSServiceIP:     dNSServiceIP,
		DockerBridgeCidr: dockerBridgeCidr,
	}

	p = []*AgentPoolProfile{
		{
			VnetSubnetID: vnetSubnetID1,
			MaxPods:      &maxPods1,
		},
		{
			MaxPods: &maxPods2,
		},
	}

	a = &Properties{
		NetworkProfile:    n,
		AgentPoolProfiles: p,
	}

	if err := validateVNET(a); err != ErrorAtLeastAgentPoolNoSubnet {
		t.Errorf("Failed to throw error, %s", ErrorAtLeastAgentPoolNoSubnet)
	}

	// NetworkPlugin = Azure, max pods is less than 5
	n = &NetworkProfile{
		NetworkPlugin:    NetworkPlugin("azure"),
		ServiceCidr:      serviceCidr,
		DNSServiceIP:     dNSServiceIP,
		DockerBridgeCidr: dockerBridgeCidr,
	}

	p = []*AgentPoolProfile{
		{
			VnetSubnetID: vnetSubnetID1,
			MaxPods:      &maxPodsTooSmall,
		},
		{
			VnetSubnetID: vnetSubnetID2,
			MaxPods:      &maxPodsTooSmall,
		},
	}

	a = &Properties{
		NetworkProfile:    n,
		AgentPoolProfiles: p,
	}

	if err := validateVNET(a); err != ErrorInvalidMaxPods {
		t.Errorf("Failed to throw error, %s", ErrorInvalidMaxPods)
	}

	// NetworkPlugin = Azure, Failed to parse VnetSubnetID
	n = &NetworkProfile{
		NetworkPlugin:    NetworkPlugin("azure"),
		ServiceCidr:      serviceCidr,
		DNSServiceIP:     dNSServiceIP,
		DockerBridgeCidr: dockerBridgeCidr,
	}

	p = []*AgentPoolProfile{
		{
			VnetSubnetID: vnetSubnetID1Bad,
			MaxPods:      &maxPods1,
		},
		{
			VnetSubnetID: vnetSubnetID2,
			MaxPods:      &maxPods2,
		},
	}

	a = &Properties{
		NetworkProfile:    n,
		AgentPoolProfiles: p,
	}

	if err := validateVNET(a); err != ErrorParsingSubnetID {
		t.Errorf("Failed to validate VNET: %s", ErrorParsingSubnetID)
	}

	// NetworkPlugin = Azure, Subscription not match
	n = &NetworkProfile{
		NetworkPlugin:    NetworkPlugin("azure"),
		ServiceCidr:      serviceCidr,
		DNSServiceIP:     dNSServiceIP,
		DockerBridgeCidr: dockerBridgeCidr,
	}

	p = []*AgentPoolProfile{
		{
			VnetSubnetID: vnetSubnetID1WrongSubscription,
			MaxPods:      &maxPods1,
		},
		{
			VnetSubnetID: vnetSubnetID2,
			MaxPods:      &maxPods2,
		},
	}

	a = &Properties{
		NetworkProfile:    n,
		AgentPoolProfiles: p,
	}

	if err := validateVNET(a); err != ErrorSubscriptionNotMatch {
		t.Errorf("Failed to validate VNET: %s", ErrorSubscriptionNotMatch)
	}

	// NetworkPlugin = Azure, ResourceGroup not match
	n = &NetworkProfile{
		NetworkPlugin:    NetworkPlugin("azure"),
		ServiceCidr:      serviceCidr,
		DNSServiceIP:     dNSServiceIP,
		DockerBridgeCidr: dockerBridgeCidr,
	}

	p = []*AgentPoolProfile{
		{
			VnetSubnetID: vnetSubnetID1WrongResourceGroup,
			MaxPods:      &maxPods1,
		},
		{
			VnetSubnetID: vnetSubnetID2,
			MaxPods:      &maxPods2,
		},
	}

	a = &Properties{
		NetworkProfile:    n,
		AgentPoolProfiles: p,
	}

	if err := validateVNET(a); err != ErrorResourceGroupNotMatch {
		t.Errorf("Failed to validate VNET: %s", ErrorResourceGroupNotMatch)
	}

	// NetworkPlugin = Azure, Vnet not match
	n = &NetworkProfile{
		NetworkPlugin:    NetworkPlugin("azure"),
		ServiceCidr:      serviceCidr,
		DNSServiceIP:     dNSServiceIP,
		DockerBridgeCidr: dockerBridgeCidr,
	}

	p = []*AgentPoolProfile{
		{
			VnetSubnetID: vnetSubnetID1WrongVnet,
			MaxPods:      &maxPods1,
		},
		{
			VnetSubnetID: vnetSubnetID2,
			MaxPods:      &maxPods2,
		},
	}

	a = &Properties{
		NetworkProfile:    n,
		AgentPoolProfiles: p,
	}

	if err := validateVNET(a); err != ErrorVnetNotMatch {
		t.Errorf("Failed to validate VNET: %s", ErrorVnetNotMatch)
	}
}

func TestValidateAADProfile(t *testing.T) {
	mc := ManagedCluster{}
	mc.Properties = &Properties{}
	mc.Properties.EnableRBAC = nil
	mc.Properties.AADProfile = &AADProfile{
		ServerAppID: "ccbfaea3-7312-497e-81d9-9ad9b8a99853",
	}
	if err := mc.Properties.AADProfile.Validate(mc.Properties.EnableRBAC); err != ErrorRBACNotEnabledForAAD {
		t.Errorf("Expected to fail because RBAC is not enabled")
	}

	mc = ManagedCluster{}
	mc.Properties = &Properties{}
	enableRBAC := true
	mc.Properties.EnableRBAC = &enableRBAC
	mc.Properties.AADProfile = &AADProfile{
		ServerAppSecret: "ccbfaea3-7312-497e-81d9-9ad9b8a99853",
	}
	if err := mc.Properties.AADProfile.Validate(mc.Properties.EnableRBAC); err != ErrorAADServerAppIDNotSet {
		t.Errorf("Expected to fail because ServerAppID is not set")
	}

	mc = ManagedCluster{}
	mc.Properties = &Properties{}
	enableRBAC = true
	mc.Properties.EnableRBAC = &enableRBAC
	mc.Properties.AADProfile = &AADProfile{
		ServerAppID: "ccbfaea3-7312-497e-81d9-9ad9b8a99853",
	}
	if err := mc.Properties.AADProfile.Validate(mc.Properties.EnableRBAC); err != ErrorAADServerAppSecretNotSet {
		t.Errorf("Expected to fail because ServerAppSecret is not set")
	}

	mc = ManagedCluster{}
	mc.Properties = &Properties{}
	enableRBAC = true
	mc.Properties.EnableRBAC = &enableRBAC
	mc.Properties.AADProfile = &AADProfile{
		ServerAppID:     "ccbfaea3-7312-497e-81d9-9ad9b8a99853",
		ServerAppSecret: "bcbfaea3-7312-497e-81d9-9ad9b8a99853",
	}
	if err := mc.Properties.AADProfile.Validate(mc.Properties.EnableRBAC); err != ErrorAADClientAppIDNotSet {
		t.Errorf("Expected to fail because ClientAppID is not set")
	}

	mc = ManagedCluster{}
	mc.Properties = &Properties{}
	enableRBAC = true
	mc.Properties.EnableRBAC = &enableRBAC
	mc.Properties.AADProfile = &AADProfile{
		ServerAppID:     "ccbfaea3-7312-497e-81d9-9ad9b8a99853",
		ServerAppSecret: "bcbfaea3-7312-497e-81d9-9ad9b8a99853",
		ClientAppID:     "acbfaea3-7312-497e-81d9-9ad9b8a99853",
	}
	if err := mc.Properties.AADProfile.Validate(mc.Properties.EnableRBAC); err != ErrorAADTenantIDNotSet {
		t.Errorf("Expected to fail because TenantID is not set")
	}
}

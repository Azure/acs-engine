package v20180331

import (
	"testing"
)

func TestValidateVNET(t *testing.T) {

	serviceCidr := "10.0.0.0/16"
	serviceCidrBad := "10.0.0.0"
	serviceCidrTooLarge := "10.0.0.0/11"
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
		t.Errorf("Failed to validate VNET: %s", err.Error())
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
		t.Errorf("Failed to validate VNET: %s", err.Error())
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
		t.Errorf("Failed to validate VNET: %s", err.Error())
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
		if err == nil {
			t.Errorf("Failed to test validate VNET: expected %s but got no error", ErrorInvalidNetworkProfile)
		}
		t.Errorf("Failed to test validate VNET: expected %s but got %s", ErrorInvalidNetworkProfile, err.Error())
	}

	// network profile has NetworkPlugin set to azure and PodCidr set, should fail
	n = &NetworkProfile{
		NetworkPlugin: NetworkPlugin("azure"),
		PodCidr:       "a.b.c.d",
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

	if err := validateVNET(a); err != ErrorPodCidrNotSetableInAzureCNI {
		if err == nil {
			t.Errorf("Failed to test validate VNET: expected %s but got no error", ErrorPodCidrNotSetableInAzureCNI)
		}
		t.Errorf("Failed to test validate VNET: expected %s but got %s", ErrorPodCidrNotSetableInAzureCNI, err.Error())
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
		if err == nil {
			t.Errorf("Failed to test validate VNET: expected %s but got no error", ErrorInvalidNetworkPlugin)
		}
		t.Errorf("Failed to test validate VNET: expected %s but got %s", ErrorInvalidNetworkPlugin, err.Error())
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
		if err == nil {
			t.Errorf("Failed to test validate VNET: expected %s but got no error", ErrorInvalidServiceCidr)
		}
		t.Errorf("Failed to test validate VNET: expected %s but got %s", ErrorInvalidServiceCidr, err.Error())
	}

	// NetworkPlugin = Azure, serviceCidr too large
	n = &NetworkProfile{
		NetworkPlugin:    NetworkPlugin("azure"),
		ServiceCidr:      serviceCidrTooLarge,
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

	if err := validateVNET(a); err != ErrorServiceCidrTooLarge {
		if err == nil {
			t.Errorf("Failed to test validate VNET: expected %s but got no error", ErrorServiceCidrTooLarge)
		}
		t.Errorf("Failed to test validate VNET: expected %s but got %s", ErrorServiceCidrTooLarge, err.Error())
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
		if err == nil {
			t.Errorf("Failed to test validate VNET: expected %s but got no error", ErrorInvalidDNSServiceIP)
		}
		t.Errorf("Failed to test validate VNET: expected %s but got %s", ErrorInvalidDNSServiceIP, err.Error())
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
		if err == nil {
			t.Errorf("Failed to test validate VNET: expected %s but got no error", ErrorInvalidDockerBridgeCidr)
		}
		t.Errorf("Failed to test validate VNET: expected %s but got %s", ErrorInvalidDockerBridgeCidr, err.Error())
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
		if err == nil {
			t.Errorf("Failed to test validate VNET: expected %s but got no error", ErrorDNSServiceIPNotInServiceCidr)
		}
		t.Errorf("Failed to test validate VNET: expected %s but got %s", ErrorDNSServiceIPNotInServiceCidr, err.Error())
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
		if err == nil {
			t.Errorf("Failed to test validate VNET: expected %s but got no error", ErrorDNSServiceIPAlreadyUsed)
		}
		t.Errorf("Failed to test validate VNET: expected %s but got %s", ErrorDNSServiceIPAlreadyUsed, err.Error())
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
		if err == nil {
			t.Errorf("Failed to test validate VNET: expected %s but got no error", ErrorAtLeastAgentPoolNoSubnet)
		}
		t.Errorf("Failed to test validate VNET: expected %s but got %s", ErrorAtLeastAgentPoolNoSubnet, err.Error())
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
		if err == nil {
			t.Errorf("Failed to test validate VNET: expected %s but got no error", ErrorInvalidMaxPods)
		}
		t.Errorf("Failed to test validate VNET: expected %s but got %s", ErrorInvalidMaxPods, err.Error())
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
		if err == nil {
			t.Errorf("Failed to test validate VNET: expected %s but got no error", ErrorParsingSubnetID)
		}
		t.Errorf("Failed to test validate VNET: expected %s but got %s", ErrorParsingSubnetID, err.Error())

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
		if err == nil {
			t.Errorf("Failed to test validate VNET: expected %s but got no error", ErrorSubscriptionNotMatch)
		}
		t.Errorf("Failed to test validate VNET: expected %s but got %s", ErrorSubscriptionNotMatch, err.Error())
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
		if err == nil {
			t.Errorf("Failed to test validate VNET: expected %s but got no error", ErrorResourceGroupNotMatch)
		}
		t.Errorf("Failed to test validate VNET: expected %s but got %s", ErrorResourceGroupNotMatch, err.Error())
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
		if err == nil {
			t.Errorf("Failed to test validate VNET: expected %s but got no error", ErrorVnetNotMatch)
		}
		t.Errorf("Failed to test validate VNET: expected %s but got %s", ErrorVnetNotMatch, err.Error())
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

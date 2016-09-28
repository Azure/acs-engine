package vlabs

import "fmt"

// SetDefaults implements APIObject
func (o *OrchestratorProfile) SetDefaults() {
}

// SetDefaults implements APIObject
func (m *MasterProfile) SetDefaults() {
	if !m.IsCustomVNET() {
		m.subnet = DefaultMasterSubnet
		m.FirstConsecutiveStaticIP = DefaultFirstConsecutiveStaticIP
	}
}

// SetDefaults implements APIObject
func (a *AgentPoolProfile) SetDefaults() {
}

// SetDefaults implements APIObject
func (l *LinuxProfile) SetDefaults() {
}

// SetDefaults implements APIObject
func (a *AcsCluster) SetDefaults() {
	a.OrchestratorProfile.SetDefaults()
	a.MasterProfile.SetDefaults()

	// assign subnets if VNET not specified
	subnetCounter := 0
	for i := range a.AgentPoolProfiles {
		profile := &a.AgentPoolProfiles[i]
		profile.SetDefaults()
		if !profile.IsCustomVNET() {
			profile.subnet = fmt.Sprintf(DefaultAgentSubnetTemplate, subnetCounter)
			subnetCounter++
		}
	}
	a.LinuxProfile.SetDefaults()
}

package vlabs

// SetDefaults implements APIObject
func (o *OrchestratorProfile) SetDefaults() {
}

// SetDefaults implements APIObject
func (m *MasterProfile) SetDefaults() {
	if m.Subnet == "" {
		m.Subnet = DefaultMasterSubnet
	}
}

// SetDefaults implements APIObject
func (a *AgentPoolProfiles) SetDefaults() {
}

// SetDefaults implements APIObject
func (l *LinuxProfile) SetDefaults() {
}

// SetDefaults implements APIObject
func (a *AcsCluster) SetDefaults() {
	a.OrchestratorProfile.SetDefaults()
	a.MasterProfile.SetDefaults()
	for _, a := range a.AgentPoolProfiles {
		a.SetDefaults()
	}
	a.LinuxProfile.SetDefaults()
}

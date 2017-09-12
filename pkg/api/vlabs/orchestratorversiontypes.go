package vlabs

// OrchestratorVersionProfile contains orchestrator version info
type OrchestratorVersionProfile struct {
	OrchestratorProfile
	Default  bool                   `json:"default,omitempty"`
	Upgrades []*OrchestratorProfile `json:"upgrades,omitempty"`
}

// OrchestratorVersionProfileList contains list of version profiles for supported orchestrators
type OrchestratorVersionProfileList struct {
	Orchestrators []*OrchestratorVersionProfile `json:"orchestrators"`
}

package v20170930

// OrchestratorEdition contains version and release numbers
type OrchestratorEdition struct {
	OrchestratorRelease string `json:"orchestratorRelease,omitempty"`
	OrchestratorVersion string `json:"orchestratorVersion"`
}

// OrchestratorVersionProfile contains orchestrator version info
type OrchestratorVersionProfile struct {
	OrchestratorType string `json:"orchestratorType"`
	OrchestratorEdition
	Default     bool                   `json:"default,omitempty"`
	Upgradables []*OrchestratorEdition `json:"upgradables,omitempty"`
}

// OrchestratorVersionProfileList contains list of version profiles for supported orchestrators
type OrchestratorVersionProfileList struct {
	Orchestrators []*OrchestratorVersionProfile `json:"orchestrators"`
}

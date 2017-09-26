package api

// OrchestratorVersionProfile contains information of a supported orchestrator version:
//  - orchestrator type and version
//  - whether this orchestrator version is deployed by default if orchestrator release is not specified
//  - list of available upgrades for this orchestrator version
type OrchestratorVersionProfile struct {
	OrchestratorProfile
	Default  bool                   `json:"default,omitempty"`
	Upgrades []*OrchestratorProfile `json:"upgrades,omitempty"`
}

// OrchestratorVersionProfileList contains list of version profiles for supported orchestrators
type OrchestratorVersionProfileList struct {
	Orchestrators []*OrchestratorVersionProfile `json:"orchestrators"`
}

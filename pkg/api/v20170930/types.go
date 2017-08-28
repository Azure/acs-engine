package v20170930

// VersionInfo contains rersion and release numbers
type VersionInfo struct {
	Release string `json:"release,omitempty"`
	Version string `json:"version"`
}

// OrchestratorInfo contains orchestrator version info
type OrchestratorInfo struct {
	Orchestrator string `json:"orchestrator"`
	VersionInfo
	DockerComposeVersion string         `json:"docker-compose-version,omitempty"`
	Default              bool           `json:"default,omitempty"`
	Upgradable           []*VersionInfo `json:"upgradable,omitempty"`
}

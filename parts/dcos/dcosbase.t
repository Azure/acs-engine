{
  "$schema": "https://schema.management.azure.com/schemas/2015-01-01/deploymentTemplate.json#",
  "contentVersion": "1.0.0.0",
  "parameters": {
    {{range .AgentPoolProfiles}}{{template "agentparams.t" .}},{{end}}
    {{if .HasWindows}}
      "dcosBinariesURL": {
        {{PopulateClassicModeDefaultValue "dcosBinariesURL"}}
        "metadata": {
          "description": "The download url for dcos/mesos windows binaries."
        },
        "type": "string"
      },
      "dcosBinariesVersion": {
        {{PopulateClassicModeDefaultValue "dcosBinariesVersion"}}
        "metadata": {
          "description": "DCOS windows binaries version"
        },
        "type": "string"
      },
      {{template "windowsparams.t"}},
    {{end}}
    {{template "dcos/dcosparams.t" .}}
    {{template "masterparams.t" .}}
  },
  "variables": {
    {{range $index, $agent := .AgentPoolProfiles}}
        "{{.Name}}Index": {{$index}},
        {{template "dcos/dcosagentvars.t" .}}
        {{if .IsStorageAccount}}
          "{{.Name}}StorageAccountOffset": "[mul(variables('maxStorageAccountsPerAgent'),{{$index}})]",
          "{{.Name}}AccountName": "[concat(variables('storageAccountBaseName'), 'agnt{{$index}}')]",
          {{if .HasDisks}}
            "{{.Name}}DataAccountName": "[concat(variables('storageAccountBaseName'), 'data{{$index}}')]",
          {{end}}
        {{end}}
    {{end}}
    
    {{template "dcos/dcosmastervars.t" .}}
  },
  "resources": [
    {{range .AgentPoolProfiles}}
      {{if .IsWindows}}
        {{if .IsAvailabilitySets}}
          {{template "dcos/dcosWindowsAgentResourcesVmas.t" .}},
        {{else}}
          {{template "dcos/dcosWindowsAgentResourcesVmss.t" .}},
        {{end}}
      {{else}}
        {{if .IsAvailabilitySets}}
          {{template "dcos/dcosagentresourcesvmas.t" .}},
        {{else}}
          {{template "dcos/dcosagentresourcesvmss.t" .}},
        {{end}}
      {{end}}
    {{end}}
    {{template "dcos/dcosmasterresources.t" .}}
  ],
  "outputs": {
    {{range .AgentPoolProfiles}}{{template "agentoutputs.t" .}}
    {{end}}
    {{template "masteroutputs.t" .}}
  }
}

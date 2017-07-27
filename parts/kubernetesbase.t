{
  "$schema": "https://schema.management.azure.com/schemas/2015-01-01/deploymentTemplate.json#",
  "contentVersion": "1.0.0.0",
  "parameters": {
    {{range .AgentPoolProfiles}}{{template "agentparams.t" .}},{{end}}
    {{if .HasWindows}}
      "kubeBinariesSASURL": {
        {{PopulateClassicModeDefaultValue "kubeBinariesSASURL"}}
        "metadata": {
          "description": "The download url for kubernetes windows binaries."
        },
        "type": "string"
      },
      "kubeBinariesVersion": {
        {{PopulateClassicModeDefaultValue "kubeBinariesVersion"}}
        "metadata": {
          "description": "Kubernetes windows binaries version"
        },
        "type": "string"
      },
      {{template "windowsparams.t"}},
    {{end}}
    {{template "masterparams.t" .}},
    {{template "kubernetesparams.t" .}}
  },
  "variables": {
    {{range $index, $agent := .AgentPoolProfiles}}
        "{{.Name}}Index": {{$index}},
        {{template "kubernetesagentvars.t" .}}
        {{if .IsStorageAccount}}
          {{if .HasDisks}}
            "{{.Name}}DataAccountName": "[concat(variables('storageAccountBaseName'), 'data{{$index}}')]",
          {{end}}
          "{{.Name}}AccountName": "[concat(variables('storageAccountBaseName'), 'agnt{{$index}}')]", 
        {{end}}
    {{end}}
    {{template "kubernetesmastervars.t" .}}
  },
  "resources": [
    {{range .AgentPoolProfiles}}
      {{if .IsWindows}}
        {{template "kuberneteswinagentresourcesvmas.t" .}},
      {{else}}
        {{template "kubernetesagentresourcesvmas.t" .}},
      {{end}}
    {{end}}
    {{template "kubernetesmasterresources.t" .}}
  ],
  "outputs": {
    {{range .AgentPoolProfiles}}
      {{template "agentoutputs.t" .}}
    {{end}}
    {{template "masteroutputs.t" .}}
  }
}
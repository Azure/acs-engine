{
  "$schema": "https://schema.management.azure.com/schemas/2015-01-01/deploymentTemplate.json#",
  "contentVersion": "1.0.0.0",
  "parameters": {
    {{range .AgentPoolProfiles}}{{template "agentparams.t" .}},{{end}}
    {{if .HasWindows}}
      {{template "windowsparams.t"}},
    {{end}}
    {{template "masterparams.t" .}}
    {{template "swarmparams.t" .}}
  },
  "variables": {
    {{range $index, $agent := .AgentPoolProfiles}}
        "{{.Name}}Index": {{$index}},
        {{template "swarmagentvars.t" .}}
        {{if .IsStorageAccount}}
          "{{.Name}}StorageAccountOffset": "[mul(variables('maxStorageAccountsPerAgent'),{{$index}})]",
          "{{.Name}}AccountName": "[concat(variables('storageAccountBaseName'), 'agnt{{$index}}')]",
          {{if .HasDisks}}
            "{{.Name}}DataAccountName": "[concat(variables('storageAccountBaseName'), 'data{{$index}}')]",
          {{end}}
        {{end}}
    {{end}}

    {{template "swarmmastervars.t" .}}
  },
  "resources": [
    {{range .AgentPoolProfiles}}
      {{if .IsWindows}}
        {{if .IsAvailabilitySets}}
          {{template "swarmwinagentresourcesvmas.t" .}},
        {{else}}
          {{template "swarmwinagentresourcesvmss.t" .}},
        {{end}}
      {{else}}
        {{if .IsAvailabilitySets}}
          {{template "swarmagentresourcesvmas.t" .}},
        {{else}}
          {{template "swarmagentresourcesvmss.t" .}},
        {{end}}
      {{end}}      
    {{end}}
    {{template "swarmmasterresources.t" .}}
  ],
  "outputs": {
    {{range .AgentPoolProfiles}}{{template "agentoutputs.t" .}}
    {{end}}
    {{template "masteroutputs.t" .}}
  }
}

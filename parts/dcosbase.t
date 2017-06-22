{
  "$schema": "https://schema.management.azure.com/schemas/2015-01-01/deploymentTemplate.json#",
  "contentVersion": "1.0.0.0",
  "parameters": {
    {{range .AgentPoolProfiles}}{{template "agentparams.t" .}},{{end}}
    {{template "dcosparams.t" .}}
    {{template "masterparams.t" .}}
  },
  "variables": {
    {{range $index, $agent := .AgentPoolProfiles}}
        {{template "dcosagentvars.t" .}}
        {{if .IsStorageAccount}}
          "{{.Name}}StorageAccountOffset": "[mul(variables('maxStorageAccountsPerAgent'),{{$index}})]",
          "{{.Name}}AccountName": "[concat(variables('storageAccountBaseName'), 'agnt{{$index}}')]",
          {{if .HasDisks}}
            "{{.Name}}DataAccountName": "[concat(variables('storageAccountBaseName'), 'data{{$index}}')]",
          {{end}}
        {{end}}
    {{end}}
    
    {{template "dcosmastervars.t" .}}
  },
  "resources": [
    {{range .AgentPoolProfiles}}
      {{if .IsAvailabilitySets}}
        {{template "dcosagentresourcesvmas.t" .}},
      {{else}}
        {{template "dcosagentresourcesvmss.t" .}},
      {{end}}
    {{end}}
    {{template "dcosmasterresources.t" .}}
  ],
  "outputs": {
    {{range .AgentPoolProfiles}}{{template "agentoutputs.t" .}}
    {{end}}{{template "masteroutputs.t" .}}
  }
}

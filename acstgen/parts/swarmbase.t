{
  "$schema": "https://schema.management.azure.com/schemas/2015-01-01/deploymentTemplate.json#",
  "contentVersion": "1.0.0.0",
  "parameters": {
    {{range .AgentPoolProfiles}}{{template "agentparams.t" .}},{{end}}
    {{template "masterparams.t" .}}
  },
  "variables": {
    {{range $index, $agent := .AgentPoolProfiles}}
        {{template "swarmagentvars.t" .}}
        {{if .IsVolumeBasedStorage}}
          "{{.Name}}DataAccountName": "[concat(variables('storageAccountBaseName'), 'data{{$index}}')]",
        {{end}}
        "{{.Name}}Index": {{$index}},
        "{{.Name}}AccountName": "[concat(variables('storageAccountBaseName'), 'agnt{{$index}}')]",
    {{end}}

    {{template "swarmmastervars.t" .}},
    
    {{GetSizeMap}}
  },
  "resources": [
    {{range .AgentPoolProfiles}}
      {{if .IsVolumeBasedStorage}}
        {{template "swarmagentresourcesdisks.t" .}},
      {{else}}
        {{template "swarmagentresources.t" .}},
      {{end}}
    {{end}}
    {{template "swarmmasterresources.t" .}}
  ],
  "outputs": {
    {{range .AgentPoolProfiles}}{{template "agentoutputs.t" .}}
    {{end}}{{template "masteroutputs.t" .}}
  }
}

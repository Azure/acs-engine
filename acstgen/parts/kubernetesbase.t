{
  "$schema": "https://schema.management.azure.com/schemas/2015-01-01/deploymentTemplate.json#",
  "contentVersion": "1.0.0.0",
  "parameters": {
    {{range .AgentPoolProfiles}}{{template "agentparams.t" .}},{{end}}
    {{template "masterparams.t" .}},
    {{template "kubernetesparams.t" .}}
  },
  "variables": {
    {{range $index, $agent := .AgentPoolProfiles}}
        {{template "kubernetesagentvars.t" .}}
        {{if .HasDisks}}
        "{{.Name}}DataAccountName": "[concat(variables('storageAccountBaseName'), 'data{{$index}}')]",
        {{end}}
        "{{.Name}}Index": {{$index}},
        "{{.Name}}AccountName": "[concat(variables('storageAccountBaseName'), 'agnt{{$index}}')]", 
    {{end}}
    {{template "kubernetesmastervars.t" .}},
    
    {{GetSizeMap}}
  },
  "resources": [
    {{range .AgentPoolProfiles}}
      {{template "kubernetesagentresources.t" .}},
    {{end}}
    {{template "kubernetesmasterresources.t" .}}
  ],
  "outputs": {
    {{template "masteroutputs.t" .}},
    "kubeConfig": {
      "type": "string", 
      "value": "[concat('{{GetKubernetesKubeConfig}}')]"
    }
  }
}
{
  "$schema": "https://schema.management.azure.com/schemas/2015-01-01/deploymentTemplate.json#",
  "contentVersion": "1.0.0.0",
  "parameters": {
    {{template "masterparams.t" .}},
    {{GetSizeMap}}
  },
  "variables": {
    {{template "kubernetesmastervars.t" .}}
  },
  "resources": [
    {{template "kubernetesmasterresources.t" .}}
  ],
  "outputs": {
    {{range .AgentPoolProfiles}}{{template "agentoutputs.t" .}}
    {{end}}{{template "masteroutputs.t" .}}
  }
}
{
  "$schema": "https://schema.management.azure.com/schemas/2015-01-01/deploymentTemplate.json#",
  "contentVersion": "1.0.0.0",
  "parameters": {
    {{template "masterparams.t" .}}
  },
  "variables": {
{{if IsDCOS}}{{template "dcosmastervars.t" .}}{{else if IsSwarm}}{{template "swarmmastervars.t" .}}{{end}}
  },
  "resources": [
{{if IsDCOS}}{{template "dcosmasterresources.t" .}}{{else if IsSwarm}}{{template "swarmmasterresources.t" .}}{{end}}
  ],
  "outputs": {
    {{template "masteroutputs.t" .}}
  }
}

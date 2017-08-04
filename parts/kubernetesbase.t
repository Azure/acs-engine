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
        {{if $index}}, {{end}}
        "{{.Name}}Index": {{$index}},
        {{template "kubernetesagentvars.t" .}}
        {{if .IsStorageAccount}}
          {{if .HasDisks}}
            "{{.Name}}DataAccountName": "[concat(variables('storageAccountBaseName'), 'data{{$index}}')]"
          {{end}}
          "{{.Name}}AccountName": "[concat(variables('storageAccountBaseName'), 'agnt{{$index}}')]"
        {{end}}
    {{end}}
    , {{template "kubernetesmastervars.t" .}}
  },
  "resources": [
    {{ range $index, $element := .AgentPoolProfiles}}
      {{if $index}}, {{end}}
      {{if .IsWindows}}
        {{template "kuberneteswinagentresourcesvmas.t" .}}
      {{else}}
        {{template "kubernetesagentresourcesvmas.t" .}}
      {{end}}
    {{end}}
    {{if .HasMaster}}
      ,{{template "kubernetesmasterresources.t" .}}
    {{else}}
    ,{
      "apiVersion": "[variables('apiVersionDefault')]",
      "dependsOn": [
        "[concat('Microsoft.Network/networkSecurityGroups/', variables('nsgName'))]"
    {{if not IsVNETIntegrated}}
        ,
        "[concat('Microsoft.Network/routeTables/', variables('routeTableName'))]"
    {{end}}
      ],
      "location": "[variables('location')]",
      "name": "[variables('virtualNetworkName')]",
      "properties": {
        "addressSpace": {
          "addressPrefixes": [
            "[variables('vnetCidr')]"
          ]
        },
        "subnets": [
          {
            "name": "[variables('subnetName')]",
            "properties": {
              "addressPrefix": "[variables('subnet')]",
              "networkSecurityGroup": {
                "id": "[variables('nsgID')]"
              }
    {{if not IsVNETIntegrated}}
              ,
              "routeTable": {
                "id": "[variables('routeTableID')]"
              }
    {{end}}
            }
          }
        ]
      },
      "type": "Microsoft.Network/virtualNetworks"
    }
    {{end}}
  ],
  "outputs": {
    {{range .AgentPoolProfiles}}
      {{template "agentoutputs.t" .}}
    {{end}}
    {{if .HasMaster}}
      , {{template "masteroutputs.t" .}}
    {{end}}
  }
}
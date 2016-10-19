    "masterFQDN": {
      "type": "string", 
      "value": "[reference(concat('Microsoft.Network/publicIPAddresses/', variables('masterPublicIPAddressName'))).dnsSettings.fqdn]"
    }
{{if  GetClassicMode}}
    {{if RequiresFakeAgentOutput}}
    ,"agentFQDN": {
      "type": "string",
      "value": ""
    },
    {{end}}
    "diagnosticsStorageAccountUri": {
      "type": "string",
      "value": ""
    },
    "jumpboxFQDN": {
      "type": "string",
      "value": ""
    }
{{end}}
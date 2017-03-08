    "masterFQDN": {
      "type": "string", 
      "value": "[reference(concat('Microsoft.Network/publicIPAddresses/', variables('masterPublicIPAddressName'))).dnsSettings.fqdn]"
    }
{{if  GetClassicMode}}
    ,
    {{if RequiresFakeAgentOutput}}
    "agentFQDN": {
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
{{if AnyAgentUsesAvailablilitySets}}
    ,
    "agentStorageAccountSuffix": {
      "type": "string",
      "value": "[variables('storageAccountBaseName')]"
    },
    "agentStorageAccountPrefixes": {
      "type": "array",
      "value": "[variables('storageAccountPrefixes')]"
    }
{{end}}
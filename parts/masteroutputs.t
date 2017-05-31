    "masterFQDN": {
      "type": "string", 
      "value": "NEEDS UPDATING"
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
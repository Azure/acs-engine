{{if IsPublic .Ports}}
  {{ if and (not IsKubernetes) (not IsOpenShift)}}
    "{{.Name}}FQDN": {
        "type": "string",
        "value": "[reference(concat('Microsoft.Network/publicIPAddresses/', variables('{{.Name}}Config').IPAddressName)).dnsSettings.fqdn]"
    },
  {{end}}
{{end}}
{{if and .IsAvailabilitySets .IsStorageAccount}}
  "{{.Name}}StorageAccountOffset": {
      "type": "int",
      "value": "[variables('{{.Name}}Config').StorageAccountOffset]"
    },
    "{{.Name}}StorageAccountCount": {
      "type": "int",
      "value": "[variables('{{.Name}}Config').StorageAccountsCount]"
    },
    "{{.Name}}SubnetName": {
      "type": "string",
      "value": "[variables('{{.Name}}Config').SubnetName]"
    },
{{end}}
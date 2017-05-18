{{if IsPublic .Ports}}
  "{{.Name}}FQDN": {
      "type": "string", 
      "value": "[reference(concat('Microsoft.Network/publicIPAddresses/', variables('{{.Name}}IPAddressName'))).dnsSettings.fqdn]"
  },
{{end}}
{{if and .IsAvailabilitySets .IsStorageAccount}}
    "{{.Name}}StorageAccountOffset": {
      "type": "int",
      "value": "[variables('{{.Name}}StorageAccountOffset')]"
    },
    "{{.Name}}StorageAccountCount": {
      "type": "int",
      "value": "[variables('{{.Name}}StorageAccountsCount')]"
    },
{{end}}

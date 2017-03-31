{{if IsPublic .Ports}}
  "{{.Name}}FQDN": {
      "type": "string", 
      "value": "[reference(concat('Microsoft.Network/publicIPAddresses/', variables('{{.Name}}IPAddressName'))).dnsSettings.fqdn]"
  },
{{end}}
{{if .IsAvailabilitySets}}
    "{{.Name}}StorageAccountOffset": {
      "type": "int",
      "value": "[variables('{{.Name}}Offset')]"
    },
    "{{.Name}}StorageAccountCount": {
      "type": "int",
      "value": "[variables('{{.Name}}Count')]"
    },
{{end}}

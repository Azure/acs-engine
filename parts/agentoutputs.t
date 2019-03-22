{{if IsPublic .Ports}}
  {{ if and (not IsKubernetes) (not IsOpenShift)}}
    "{{.Name}}FQDN": {
        "type": "string",
        "value": "[reference(concat('Microsoft.Network/publicIPAddresses/', variables('{{.Name}}IPAddressName'))).dnsSettings.fqdn]"
    },
  {{end}}
{{end}}
{{if and .IsAvailabilitySets .IsStorageAccount}}
  "{{.Name}}StorageAccountOffset": {
      "type": "int",
      "value": "[variables('{{.Name}}Variables').StorageAccountOffset]"
    },
    "{{.Name}}StorageAccountCount": {
      "type": "int",
      "value": "[variables('{{.Name}}Variables').StorageAccountsCount]"
    },
    "{{.Name}}SubnetName": {
      "type": "string",
      "value": "[variables('{{.Name}}Variables').SubnetName]"
    },
{{end}}
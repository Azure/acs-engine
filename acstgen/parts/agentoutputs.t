{{if IsPublic .Ports}}
  "{{.Name}}FQDN": {
      "type": "string", 
      "value": "[reference(concat('Microsoft.Network/publicIPAddresses/', variables('{{.Name}}IPAddressName'))).dnsSettings.fqdn]"
  },
{{end}}
    "masterFQDN": {
      "type": "string", 
      "value": "[reference(concat('Microsoft.Network/publicIPAddresses/', variables('masterPublicIPAddressName'))).dnsSettings.fqdn]"
    }

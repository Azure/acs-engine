    "linuxAdminUsername": {
      "metadata": {
        "description": "User name for the Linux Virtual Machines (SSH or Password)."
      },
      "type": "string"
    },
    "bootstrapEndpointDNSNamePrefix": {
      "metadata": {
        "description": "Sets the Domain name label for the bootstrap IP Address.  The concatenation of the domain name label and the regional DNS zone make up the fully qualified domain name associated with the public IP address."
      },
      "type": "string"
    },
    {{range .ExtensionProfiles}}
      "{{.Name}}Parameters": {
        "metadata": {
        "description": "Parameters for the extension"
      },
      "type": "securestring"
      },
    {{end}}
{{if IsHostedBootstrap}}
    "bootstrapSubnet": {
      "defaultValue": "{{.HostedBootstrapProfile.Subnet}}",
      "metadata": {
        "description": "Sets the subnet for the VMs in the cluster."
      },
      "type": "string"
    },
    "bootstrapEndpoint": {
      "defaultValue": "{{.HostedBootstrapProfile.FQDN}}",
      "metadata": {
        "description": "Sets the static IP of the first bootstrap"
      },
      "type": "string"
    },
{{else}}
    "bootstrapFirstConsecutiveStaticIP": {
      "metadata": {
        "description": "Sets the static IP of the first bootstrap"
      },
      "type": "string"
    },
    "bootstrapVMSize": {
      {{GetMasterAllowedSizes}}
      "metadata": {
        "description": "The size of the Virtual Machine."
      },
      "type": "string"
    },
    "bootstrapCount": {
      "metadata": {
        "description": "The number of the Virtual Machines in the set."
      },
      "type": "int"
    },
{{end}}
    "sshRSAPublicKey": {
      "metadata": {
        "description": "SSH public key used for auth to all Linux machines.  Not Required.  If not set, you must provide a password key."
      },
      "type": "string"
    },
    "nameSuffix": {
      "defaultValue": "{{GetUniqueNameSuffix}}",
      "metadata": {
        "description": "A string hash of the bootstrap DNS name to uniquely identify the cluster."
      },
      "type": "string"
    },
    "osImageName": {
      "defaultValue": "",
      "metadata": {
        "description": "Name of a Linux OS image. Needs to be used in conjuction with osImageResourceGroup."
      },
      "type": "string"
    },
    "osImageResourceGroup": {
      "defaultValue": "",
      "metadata": {
        "description": "Resource group of a Linux OS image. Needs to be used in conjuction with osImageName."
      },
      "type": "string"
    },
    "osImageOffer": {
      "defaultValue": "UbuntuServer",
      "metadata": {
        "description": "Linux OS image type."
      },
      "type": "string"
    },
    "osImagePublisher": {
      "defaultValue": "Canonical",
      "metadata": {
        "description": "OS image publisher."
      },
      "type": "string"
    },
    "osImageSKU": {
      "defaultValue": "16.04-LTS",
      "metadata": {
        "description": "OS image SKU."
      },
      "type": "string"
    },
    "osImageVersion": {
      "defaultValue": "16.04.201804050",
      "metadata": {
        "description": "OS image version."
      },
      "type": "string"
    },
    "fqdnEndpointSuffix":{
      "defaultValue": "cloudapp.azure.com",
      "metadata": {
        "description": "Endpoint of FQDN."
      },
      "type": "string"
    },
    "targetEnvironment": {
      "defaultValue": "AzurePublicCloud",
      "metadata": {
        "description": "The azure deploy environment. Currently support: AzurePublicCloud, AzureChinaCloud"
      },
      "type": "string"
    },
    "location": {
      "defaultValue": "{{GetLocation}}",
      "metadata": {
        "description": "Sets the location for all resources in the cluster"
      },
      "type": "string"
    }
{{if  GetClassicMode}}
    ,{{template "classicparams.t" .}}
{{end}}
{{if .LinuxProfile.HasSecrets}}
  {{range  $vIndex, $vault := .LinuxProfile.Secrets}}
    ,
    "linuxKeyVaultID{{$vIndex}}": {
      "metadata": {
        "description": "KeyVaultId{{$vIndex}} to install certificates from on linux machines."
      },
      "type": "string"
    }
    {{range $cIndex, $cert := $vault.VaultCertificates}}
      ,
      "linuxKeyVaultID{{$vIndex}}CertificateURL{{$cIndex}}": {
        "metadata": {
          "description": "CertificateURL{{$cIndex}} to install from KeyVaultId{{$vIndex}} on linux machines."
        },
        "type": "string"
      }
    {{end}}
  {{end}}
{{end}}
{{if .HasWindows}}{{if .WindowsProfile.HasSecrets}}
  {{range  $vIndex, $vault := .WindowsProfile.Secrets}}
    ,
    "windowsKeyVaultID{{$vIndex}}": {
      "metadata": {
        "description": "KeyVaultId{{$vIndex}} to install certificates from on windows machines."
      },
      "type": "string"
    }
    {{range $cIndex, $cert := $vault.VaultCertificates}}
      ,
      "windowsKeyVaultID{{$vIndex}}CertificateURL{{$cIndex}}": {
        "metadata": {
          "description": "Url to retrieve Certificate{{$cIndex}} from KeyVaultId{{$vIndex}} to install on windows machines."
        },
        "type": "string"
      },
      "windowsKeyVaultID{{$vIndex}}CertificateStore{{$cIndex}}": {
        "metadata": {
          "description": "CertificateStore to install Certificate{{$cIndex}} from KeyVaultId{{$vIndex}} on windows machines."
        },
        "type": "string"
      }
    {{end}}
  {{end}}
{{end}} {{end}}

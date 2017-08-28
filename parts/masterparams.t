    "linuxAdminUsername": {
      "metadata": {
        "description": "User name for the Linux Virtual Machines (SSH or Password)."
      }, 
      "type": "string"
    },
    "masterEndpointDNSNamePrefix": {
      "metadata": {
        "description": "Sets the Domain name label for the master IP Address.  The concatenation of the domain name label and the regional DNS zone make up the fully qualified domain name associated with the public IP address."
      }, 
      "type": "string"
    },
{{if not IsHostedMaster }}
  {{if .MasterProfile.IsCustomVNET}}
    "masterVnetSubnetID": {
      "metadata": {
        "description": "Sets the vnet subnet of the master."
      }, 
      "type": "string"
    },
  {{else}}
    "masterSubnet": {
      "defaultValue": "{{.MasterProfile.Subnet}}",
      "metadata": {
        "description": "Sets the subnet of the master node(s)."
      }, 
      "type": "string"
    },
  {{end}}
{{end}}
{{if IsHostedMaster}}
    "kubernetesEndpoint": {
      "defaultValue": "{{.HostedMasterProfile.FQDN}}",
      "metadata": {
        "description": "Sets the static IP of the first master"
      },
      "type": "string"
    },
{{else}}
    "firstConsecutiveStaticIP": {
      "defaultValue": "{{.MasterProfile.FirstConsecutiveStaticIP}}",
      "metadata": {
        "description": "Sets the static IP of the first master"
      }, 
      "type": "string"
    },
    "masterVMSize": {
      {{GetMasterAllowedSizes}}
      "metadata": {
        "description": "The size of the Virtual Machine."
      }, 
      "type": "string"
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
        "description": "A string hash of the master DNS name to uniquely identify the cluster."
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

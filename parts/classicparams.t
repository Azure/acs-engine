   "agentEndpointDNSNamePrefix": {
      "defaultValue": "UNUSED",
      "metadata": {
        "description": "Sets the Domain name label for the agent pool IP Address.  The concatenation of the domain name label and the regional DNS zone make up the fully qualified domain name associated with the public IP address."
      },
      "type": "string"
    },
   "disablePasswordAuthentication": {
      "defaultValue": true,
      "metadata": {
        "description": "This setting controls whether password auth is disabled for Linux VMs provisioned by this template. Default is true which disables password and makes SSH key required."
      },
      "type": "bool"
    },
    "enableNewStorageAccountNaming": {
      "defaultValue": true,
      "metadata": {
        "description": "If true: uses DNS name prefix + Orchestrator name + Region to create storage account name to reduce name collision probability. If false: uses DNS name prefix + Orchestrator name to create storage account name to maintain template idempotency."
      },
      "type": "bool"
    },
    "enableVMDiagnostics": {
      "defaultValue": true,
      "metadata": {
        "description": "Allows user to enable/disable boot & vm diagnostics."
      },
      "type": "bool"
    },
    "isValidation": {
      "allowedValues": [
        0,
        1
      ],
      "defaultValue": 0,
      "metadata": {
        "description": "This is testing in the validation region"
      },
      "type": "int"
    },
    "jumpboxEndpointDNSNamePrefix": {
      "defaultValue": "UNUSED",
      "metadata": {
        "description": "Sets the Domain name label for the jumpbox.  The concatenation of the domain name label and the regionalized DNS zone make up the fully qualified domain name associated with the public IP address."
      },
      "type": "string"
    },
    "linuxAdminPassword": {
      "defaultValue": "UNUSED",
      "metadata": {
        "description": "Password for the Linux Virtual Machine.  Not Required.  If not set, you must provide a SSH key."
      },
      "type": "securestring"
    },
   "linuxOffer": {
      "defaultValue": "UNUSED",
      "metadata": {
        "description": "This is the offer of the image used by the linux cluster"
      },
      "type": "string"
    },
    "linuxPublisher": {
      "defaultValue": "UNUSED",
      "metadata": {
        "description": "This is the publisher of the image used by the linux cluster"
      },
      "type": "string"
    },
    "linuxSku": {
      "defaultValue": "UNUSED",
      "metadata": {
        "description": "This is the linux sku used by the linux cluster"
      },
      "type": "string"
    },
    "linuxVersion": {
      "defaultValue": "UNUSED",
      "metadata": {
        "description": "This is the linux version used by the linux cluster"
      },
      "type": "string"
    },
    "masterCount": {
      "allowedValues": [
        1,
        3,
        5
      ],
      "defaultValue": 1,
      "metadata": {
        "description": "The number of Mesos masters for the cluster."
      },
      "type": "int"
    },
   "oauthEnabled": {
      "allowedValues": [
        "true",
        "false"
      ],
      "defaultValue": "false",
      "metadata": {
        "description": "Enable OAuth authentication"
      },
      "type": "string"
    },
    "postInstallScriptURI": {
      "defaultValue": "disabled",
      "metadata": {
        "description": "After installation, this specifies a script to download and install.  To disabled, set value to 'disabled'."
      },
      "type": "string"
    },
   "setLinuxConfigurationForVMCreate": {
      "allowedValues": [
        0,
        1
      ],
      "defaultValue": 1,
      "metadata": {
        "description": "This setting controls whether Linux configuration with SSH Key is passed in VM PUT Payload.  Defaults to 1.  If SSH Key is blank, this must be set to 0."
      },
      "type": "int"
    },
   "vmsPerStorageAccount": {
      "defaultValue": 5,
      "metadata": {
        "description": "This specifies the number of VMs per storage accounts"
      },
      "type": "int"
    },
{{if not .HasWindows}}
    "windowsAdminPassword": {
      "defaultValue": "UNUSED",
      "metadata": {
        "description": "Password for the Windows Virtual Machine."
      },
      "type": "securestring"
    },
    "windowsAdminUsername": {
      "defaultValue": "UNUSED",
      "metadata": {
        "description": "User name for the Windows Virtual Machine (Password Only Supported)."
      },
      "type": "string"
    },
    "kubeBinariesSASURL": {
      "defaultValue": "UNUSED",
      "metadata": {
        "description": "The download url for kubernetes windows binaries."
      },
      "type": "string"
    },
{{end}}
    "windowsJumpboxOffer": {
      "defaultValue": "UNUSED",
      "metadata": {
        "description": "This is the windows offer used by the windows"
      },
      "type": "string"
    },
    "windowsJumpboxPublisher": {
      "defaultValue": "UNUSED",
      "metadata": {
        "description": "This is the windows publisher used by the windows"
      },
      "type": "string"
    },
    "windowsJumpboxSku": {
      "defaultValue": "UNUSED",
      "metadata": {
        "description": "This is the windows sku used by the windows"
      },
      "type": "string"
    }
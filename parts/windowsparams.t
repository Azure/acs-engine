    "windowsAdminUsername": {
      "type": "string",
      "metadata": {
        "description": "User name for the Windows Swarm Agent Virtual Machines (Password Only Supported)."
      }
    },
    "windowsAdminPassword": {
      "type": "securestring",
      "metadata": {
        "description": "Password for the Windows Swarm Agent Virtual Machines."
      }
    },
    "agentWindowsVersion": {
      "defaultValue": "latest",
      "metadata": {
        "description": "Version of the Windows Server 2016 OS image to use for the agent virtual machines."
      },
      "type": "string"
    },
    "agentWindowsSourceUrl": {
      "defaultValue": "",
      "metadata": {
        "description": "The source of the generalized blob which will be used to create a custom windows image for the agent virtual machines."
      },
      "type": "string"
    },
    "agentWindowsPublisher": {
      "defaultValue": "MicrosoftWindowsServer",
      "metadata": {
        "description": "The publisher of windows image for the agent virtual machines."
      },
      "type": "string"
    },
    "agentWindowsOffer": {
      "defaultValue": "WindowsServerSemiAnnual",
      "metadata": {
        "description": "The offer of windows image for the agent virtual machines."
      },
      "type": "string"
    },
    "agentWindowsSku": {
      "defaultValue": "Datacenter-Core-1803-with-Containers-smalldisk",
      "metadata": {
        "description": "The SKU of windows image for the agent virtual machines."
      },
      "type": "string"
    },
    "windowsDockerVersion": {
      "defaultValue": "17.06.2-ee-16",
      "metadata": {
        "description": "The version of Docker to be installed on Windows Nodes"
      },
      "type": "string"
    }

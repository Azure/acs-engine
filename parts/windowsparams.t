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
    }
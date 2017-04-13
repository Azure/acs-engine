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
    "kubeBinariesSASURL": {
      "defaultValue": "https://acsengine.blob.core.windows.net/wink8s/v1.6.0int.zip",
      "metadata": {
        "description": "The download url for kubernetes windows binaries."
      },
      "type": "string"
    }
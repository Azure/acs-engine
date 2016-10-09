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
{{if .MasterProfile.IsCustomVNET}}
    "masterVnetSubnetID": {
      "metadata": {
        "description": "Sets the vnet subnet of the master."
      }, 
      "type": "string"
    },
{{else}}
    "masterSubnet": {
      "metadata": {
        "description": "Sets the subnet of the master node(s)."
      }, 
      "type": "string"
    },
{{end}}
    "firstConsecutiveStaticIP": {
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
    "sshRSAPublicKey": {
      "metadata": {
        "description": "SSH public key used for auth to all Linux machines.  Not Required.  If not set, you must provide a password key."
      }, 
      "type": "string"
    }
    "linuxAdminUsername": {
      "defaultValue": "{{.LinuxProfile.AdminUsername}}", 
      "metadata": {
        "description": "User name for the Linux Virtual Machines (SSH or Password)."
      }, 
      "type": "string"
    },
    "masterEndpointDNSNamePrefix": {
      "defaultValue": "{{.MasterProfile.DNSPrefix}}",
      "metadata": {
        "description": "Sets the Domain name label for the master IP Address.  The concatenation of the domain name label and the regional DNS zone make up the fully qualified domain name associated with the public IP address."
      }, 
      "type": "string"
    },
{{if .MasterProfile.IsCustomVNET}}
    "masterVnetSubnetID": {
      "defaultValue": "{{.MasterProfile.VnetSubnetID}}",
      "metadata": {
        "description": "Sets the vnet subnet of the master."
      }, 
      "type": "string"
    },
{{else}}
    "masterSubnet": {
      "defaultValue": "{{.MasterProfile.GetSubnet}}",
      "metadata": {
        "description": "Sets the subnet of the master."
      }, 
      "type": "string"
    },
{{end}}
    "firstConsecutiveStaticIP": {
      "defaultValue": "{{.MasterProfile.FirstConsecutiveStaticIP}}",
      "metadata": {
        "description": "Sets the static IP of the first master"
      }, 
      "type": "string"
    },
    "masterVMSize": {
      {{GetMasterAllowedSizes}}
      "defaultValue": "{{.MasterProfile.VMSize}}", 
      "metadata": {
        "description": "The size of the Virtual Machine."
      }, 
      "type": "string"
    }, 
    "sshRSAPublicKey": {
      "defaultValue": "{{GetLinuxProfileFirstSSHPublicKey}}", 
      "metadata": {
        "description": "SSH public key used for auth to all Linux machines.  Not Required.  If not set, you must provide a password key."
      }, 
      "type": "string"
    }
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
    "masterSubnet": {
      "defaultValue": "{{.MasterProfile.Subnet}}",
      "metadata": {
        "description": "Sets the subnet of the master, must be specified in CIDR format with a /24 subnet."
      }, 
      "type": "string"
    },
    "masterVMSize": {
      "allowedValues": [
        "Standard_A0", 
        "Standard_A1", 
        "Standard_A2", 
        "Standard_A3", 
        "Standard_A4", 
        "Standard_A5", 
        "Standard_A6", 
        "Standard_A7", 
        "Standard_A8", 
        "Standard_A9", 
        "Standard_A10", 
        "Standard_A11", 
        "Standard_D1", 
        "Standard_D2", 
        "Standard_D3", 
        "Standard_D4", 
        "Standard_D11", 
        "Standard_D12", 
        "Standard_D13", 
        "Standard_D14", 
        "Standard_D1_v2", 
        "Standard_D2_v2", 
        "Standard_D3_v2", 
        "Standard_D4_v2", 
        "Standard_D5_v2", 
        "Standard_D11_v2", 
        "Standard_D12_v2", 
        "Standard_D13_v2", 
        "Standard_D14_v2", 
        "Standard_G1", 
        "Standard_G2", 
        "Standard_G3", 
        "Standard_G4", 
        "Standard_G5", 
        "Standard_DS1", 
        "Standard_DS2", 
        "Standard_DS3", 
        "Standard_DS4", 
        "Standard_DS11", 
        "Standard_DS12", 
        "Standard_DS13", 
        "Standard_DS14", 
        "Standard_GS1", 
        "Standard_GS2", 
        "Standard_GS3", 
        "Standard_GS4", 
        "Standard_GS5"
      ], 
      "defaultValue": "{{.MasterProfile.VMSize}}", 
      "metadata": {
        "description": "The size of the Virtual Machine."
      }, 
      "type": "string"
    }, 
    "sshRSAPublicKey": {
      "defaultValue": "{{.LinuxProfileFirstSSHPublicKey}}", 
      "metadata": {
        "description": "SSH public key used for auth to all Linux machines.  Not Required.  If not set, you must provide a password key."
      }, 
      "type": "string"
    }
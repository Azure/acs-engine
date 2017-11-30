    "{{.Name}}Count": {
      "allowedValues": [1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21,22,23,24,25,26,27,28,29,30,31,32,33,34,35,36,37,38,39,40,41,42,43,44,45,46,47,48,49,50,51,52,53,54,55,56,57,58,59,60,61,62,63,64,65,66,67,68,69,70,71,72,73,74,75,76,77,78,79,80,81,82,83,84,85,86,87,88,89,90,91,92,93,94,95,96,97,98,99,100], 
      "defaultValue": {{.Count}},
      "metadata": {
        "description": "The number of agents for the cluster.  This value can be from 1 to 100"
      }, 
      "type": "int"
    },
{{if .IsAvailabilitySets}}
    "{{.Name}}Offset": {
      "allowedValues": [0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21,22,23,24,25,26,27,28,29,30,31,32,33,34,35,36,37,38,39,40,41,42,43,44,45,46,47,48,49,50,51,52,53,54,55,56,57,58,59,60,61,62,63,64,65,66,67,68,69,70,71,72,73,74,75,76,77,78,79,80,81,82,83,84,85,86,87,88,89,90,91,92,93,94,95,96,97,98,99], 
      "defaultValue": 0,
      "metadata": {
        "description": "The offset into the agent pool where to start creating agents.  This value can be from 0 to 99, but must be less than agentCount"
      }, 
      "type": "int"
    },
{{end}}
    "{{.Name}}VMSize": {
      {{GetAgentAllowedSizes}}
      "defaultValue": "{{.VMSize}}",
      "metadata": {
        "description": "The size of the Virtual Machine."
      }, 
      "type": "string"
    },
    "{{.Name}}osImageOffer": {
      "defaultValue": "UbuntuServer",
      "metadata": {
        "description": "Linux OS image type."
      },
      "type": "string"
    }, 
    "{{.Name}}osImagePublisher": {
      "defaultValue": "Canonical",
      "metadata": {
        "description": "OS image publisher."
      },
      "type": "string"
    }, 
    "{{.Name}}osImageSKU": {
      "defaultValue": "16.04-LTS",
      "metadata": {
        "description": "OS image SKU."
      },
      "type": "string"
    }, 
    "{{.Name}}osImageVersion": {
      "defaultValue": "16.04.201711211",
      "metadata": {
        "description": "OS image version."
      },
      "type": "string"
    },
{{if .IsCustomVNET}}
    "{{.Name}}VnetSubnetID": {
      "metadata": {
        "description": "Sets the vnet subnet of agent pool '{{.Name}}'."
      }, 
      "type": "string"
    }
{{else}}
    "{{.Name}}Subnet": {
      "defaultValue": "{{.Subnet}}",
      "metadata": {
        "description": "Sets the subnet of agent pool '{{.Name}}'."
      }, 
      "type": "string"
    }
{{end}}
{{if IsPublic .Ports}}
  ,"{{.Name}}EndpointDNSNamePrefix": {
      "metadata": {
        "description": "Sets the Domain name label for the agent pool IP Address.  The concatenation of the domain name label and the regional DNS zone make up the fully qualified domain name associated with the public IP address."
      }, 
      "type": "string"
    }
{{end}} 

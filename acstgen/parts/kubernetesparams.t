    "apiServerCertificate": {
      "defaultValue": "{{Base64 .OrchestratorProfile.ApiServerCertificate}}", 
      "metadata": {
        "description": "The AD Tenant Id"
      }, 
      "type": "string"
    }, 
    "apiServerPrivateKey": {
      "defaultValue": "{{Base64 .OrchestratorProfile.ApiServerPrivateKey}}", 
      "metadata": {
        "description": "User name for the Linux Virtual Machines (SSH or Password)."
      }, 
      "type": "securestring"
    }, 
    "caCertificate": {
      "defaultValue": "{{Base64 .OrchestratorProfile.CaCertificate}}",  
      "metadata": {
        "description": "The certificate authority certificate"
      },  
      "type": "string"
    }, 
    "clientCertificate": {
      "defaultValue": "{{Base64 .OrchestratorProfile.ClientCertificate}}",  
      "metadata": {
        "description": "The client certificate used to communicate with the master"
      }, 
      "type": "string"
    }, 
    "clientPrivateKey": {
      "defaultValue": "{{Base64 .OrchestratorProfile.ClientPrivateKey}}", 
      "metadata": {
        "description": "The client private key used to communicate with the master"
      },  
      "type": "securestring"
    }, 
    "kubernetesHyperkubeSpec": {
      "defaultValue": "gcr.io/google_containers/hyperkube-amd64:v1.4.0-beta.10",
      "metadata": {
        "description": "The container spec for hyperkube."
      }, 
      "type": "string"
    }, 
    "servicePrincipalClientId": {
      "defaultValue": "{{.OrchestratorProfile.ServicePrincipalClientID}}",
      "metadata": {
        "description": "Client ID (used by cloudprovider)"
      }, 
      "type": "string"
    }, 
    "servicePrincipalClientSecret": {
      "defaultValue": "{{.OrchestratorProfile.ServicePrincipalClientSecret}}",
      "metadata": {
        "description": "The Service Principal Client Secret."
      },
      "type": "string"
    }
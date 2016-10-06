    "apiServerCertificate": {
      "defaultValue": "{{.OrchestratorProfile.ApiServerCertificate}}", 
      "metadata": {
        "description": "The AD Tenant Id"
      }, 
      "type": "string"
    }, 
    "apiServerPrivateKey": {
      "defaultValue": "{{.OrchestratorProfile.ApiServerPrivateKey}}", 
      "metadata": {
        "description": "User name for the Linux Virtual Machines (SSH or Password)."
      }, 
      "type": "securestring"
    }, 
    "caCertificate": {
      "defaultValue": "{{.OrchestratorProfile.CaCertificate}}",  
      "metadata": {
        "description": "The certificate authority certificate"
      },  
      "type": "string"
    }, 
    "caPrivateKey": {
      "defaultValue": "{{.OrchestratorProfile.GetCAPrivateKey}}",  
      "metadata": {
        "description": "The certificate authority private key, this is only generated once and not used in this file, save this file for future updates"
      },  
      "type": "string"
    }, 
    "clientCertificate": {
      "defaultValue": "{{.OrchestratorProfile.ClientCertificate}}",  
      "metadata": {
        "description": "The client certificate used to communicate with the master"
      }, 
      "type": "string"
    }, 
    "clientPrivateKey": {
      "defaultValue": "{{.OrchestratorProfile.ClientPrivateKey}}", 
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
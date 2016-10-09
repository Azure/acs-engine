    "apiServerCertificate": {
      "metadata": {
        "description": "The base 64 server certificate used on the master"
      }, 
      "type": "string"
    }, 
    "apiServerPrivateKey": {
      "metadata": {
        "description": "The base 64 server private key used on the master."
      }, 
      "type": "securestring"
    }, 
    "caCertificate": {
      "metadata": {
        "description": "The base 64 certificate authority certificate"
      },  
      "type": "string"
    }, 
    "clientCertificate": {
      "metadata": {
        "description": "The base 64 client certificate used to communicate with the master"
      }, 
      "type": "string"
    }, 
    "clientPrivateKey": {
      "metadata": {
        "description": "The base 64 client private key used to communicate with the master"
      },  
      "type": "securestring"
    },
    "kubeConfigCertificate": {
      "metadata": {
        "description": "The base 64 certificate used by cli to communicate with the master"
      }, 
      "type": "string"
    }, 
    "kubeConfigPrivateKey": {
      "metadata": {
        "description": "The base 64 private key used by cli to communicate with the master"
      },  
      "type": "securestring"
    }, 
    "kubernetesHyperkubeSpec": {
      "metadata": {
        "description": "The container spec for hyperkube."
      }, 
      "type": "string"
    }, 
    "servicePrincipalClientId": {
      "metadata": {
        "description": "Client ID (used by cloudprovider)"
      }, 
      "type": "securestring"
    }, 
    "servicePrincipalClientSecret": {
      "metadata": {
        "description": "The Service Principal Client Secret."
      },
      "type": "securestring"
    }
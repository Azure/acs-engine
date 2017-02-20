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
      "defaultValue": "{{.OrchestratorProfile.KubernetesConfig.KubernetesHyperkubeSpec}}",
      "metadata": {
        "description": "The container spec for hyperkube."
      },
      "type": "string"
    },
    "kubectlVersion": {
      "defaultValue": "{{.OrchestratorProfile.KubernetesConfig.KubectlVersion}}",
      "metadata": {
        "description": "The kubernetes version."
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
    },
    "kubeDnsServiceIp": {
        "metadata": {
            "description": "Kubernetes DNS IP"
        },
        "type": "string"
    },
    "kubeServiceCidr": {
        "metadata": {
            "description": "Kubernetes service address space"
        },
        "type": "string"
    },
    "kubeClusterCidr": {
        "metadata": {
            "description": "Kubernetes cluster address space"
        },
        "type": "string"
    }

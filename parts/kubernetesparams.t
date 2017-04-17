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
    "kubeDnsServiceIP": {
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
    },
    "kubernetesHyperkubeSpec": {
      "defaultValue": "",
      "metadata": {
        "description": "The container spec for hyperkube."
      },
      "type": "string"
    },
    "kubernetesAddonManagerSpec": {
      "defaultValue": "",
      "metadata": {
        "description": "The container spec for hyperkube."
      },
      "type": "string"
    },
    "kubernetesAddonResizerSpec": {
      "defaultValue": "",
      "metadata": {
        "description": "The container spec for addon-resizer."
      },
      "type": "string"
    },
    "kubernetesDashboardSpec": {
      "defaultValue": "",
      "metadata": {
        "description": "The container spec for kubernetes-dashboard-amd64."
      },
      "type": "string"
    },
    "kubernetesExecHealthzSpec": {
      "defaultValue": "",
      "metadata": {
        "description": "The container spec for exechealthz-amd64."
      },
      "type": "string"
    },
    "kubernetesHeapsterSpec": {
      "defaultValue": "",
      "metadata": {
        "description": "The container spec for heapster."
      },
      "type": "string"
    },
    "kubernetesPodInfraContainerSpec": {
      "defaultValue": "",
      "metadata": {
        "description": "The container spec for pod infra."
      },
      "type": "string"
    },
    "kubernetesKubeDNSSpec": {
      "defaultValue": "",
      "metadata": {
        "description": "The container spec for kubedns-amd64."
      },
      "type": "string"
    },
    "kubernetesDNSMasqSpec": {
      "defaultValue": "",
      "metadata": {
        "description": "The container spec for kube-dnsmasq-amd64."
      },
      "type": "string"
    },
    "dockerEngineDownloadRepo": {
      "defaultValue": "https://aptdocker.azureedge.net/repo",
      "metadata": {
        "description": "The docker engine download url for kubernetes."
      },
      "type": "string"
    },
    "networkPolicy": {
      "defaultValue": "{{.OrchestratorProfile.KubernetesConfig.NetworkPolicy}}",
      "metadata": {
        "description": "The network policy enforcement to use (none|azure|calico)"
      },
      "allowedValues": [
        "none",
        "azure",
        "calico"
      ],
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

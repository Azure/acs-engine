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
    "kubernetesAddonManagerSpec": {
      "defaultValue": "{{.OrchestratorProfile.KubernetesConfig.KubernetesAddonManagerSpec}}",
      "metadata": {
        "description": "The container spec for hyperkube."
      },
      "type": "string"
    },
    "kubernetesAddonResizerSpec": {
      "defaultValue": "{{.OrchestratorProfile.KubernetesConfig.KubernetesAddonResizerSpec}}",
      "metadata": {
        "description": "The container spec for addon-resizer."
      },
      "type": "string"
    },
    "kubernetesDashboardSpec": {
      "defaultValue": "{{.OrchestratorProfile.KubernetesConfig.KubernetesDashboardSpec}}",
      "metadata": {
        "description": "The container spec for kubernetes-dashboard-amd64."
      },
      "type": "string"
    },
    "kubernetesExecHealthzSpec": {
      "defaultValue": "{{.OrchestratorProfile.KubernetesConfig.KubernetesExecHealthzSpec}}",
      "metadata": {
        "description": "The container spec for exechealthz-amd64."
      },
      "type": "string"
    },
    "kubernetesHeapsterSpec": {
      "defaultValue": "{{.OrchestratorProfile.KubernetesConfig.KubernetesHeapsterSpec}}",
      "metadata": {
        "description": "The container spec for heapster."
      },
      "type": "string"
    },
    "kubernetesPodInfraContainerSpec": {
      "defaultValue": "{{.OrchestratorProfile.KubernetesConfig.KubernetesPodInfraContainerSpec}}",
      "metadata": {
        "description": "The container spec for pod infra."
      },
      "type": "string"
    },
    "kubernetesKubeDNSSpec": {
      "defaultValue": "{{.OrchestratorProfile.KubernetesConfig.KubernetesKubeDNSSpec}}",
      "metadata": {
        "description": "The container spec for kubedns-amd64."
      },
      "type": "string"
    },
    "kubernetesDNSMasqSpec": {
      "defaultValue": "{{.OrchestratorProfile.KubernetesConfig.KubernetesDNSMasqSpec}}",
      "metadata": {
        "description": "The container spec for kube-dnsmasq-amd64."
      },
      "type": "string"
    },
    "kubectlDownloadURL": {
      "defaultValue": "{{.OrchestratorProfile.KubernetesConfig.KubectlDownloadURL}}",
      "metadata": {
        "description": "This is the download URL of kubectl"
      },
      "type": "string"
    },
    "dockerInstallScriptURL": {
      "defaultValue": "{{.OrchestratorProfile.KubernetesConfig.DockerInstallScriptURL}}",
      "metadata": {
        "description": "This is the download URL of docker installer script"
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
    }
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
    "caPrivateKey": {
      {{PopulateClassicModeDefaultValue "caPrivateKey"}}
      "metadata": {
        "description": "The base 64 CA private key used on the master."
      },
      "type": "securestring"
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
    "dockerBridgeCidr": {
      {{PopulateClassicModeDefaultValue "dockerBridgeCidr"}}
      "metadata": {
        "description": "Docker bridge network IP address and subnet"
      },
      "type": "string"
    },
    "kubeClusterCidr": {
      {{PopulateClassicModeDefaultValue "kubeClusterCidr"}}
      "metadata": {
        "description": "Kubernetes cluster subnet"
      },
      "type": "string"
    },
    "kubernetesHyperkubeSpec": {
      {{PopulateClassicModeDefaultValue "kubernetesHyperkubeSpec"}}
      "metadata": {
        "description": "The container spec for hyperkube."
      },
      "type": "string"
    },
    "kubernetesAddonManagerSpec": {
      {{PopulateClassicModeDefaultValue "kubernetesAddonManagerSpec"}}
      "metadata": {
        "description": "The container spec for hyperkube."
      },
      "type": "string"
    },
    "kubernetesAddonResizerSpec": {
      {{PopulateClassicModeDefaultValue "kubernetesAddonResizerSpec"}}
      "metadata": {
        "description": "The container spec for addon-resizer."
      },
      "type": "string"
    },
    "kubernetesDashboardSpec": {
      {{PopulateClassicModeDefaultValue "kubernetesDashboardSpec"}}
      "metadata": {
        "description": "The container spec for kubernetes-dashboard-amd64."
      },
      "type": "string"
    },
    "kubernetesExecHealthzSpec": {
      {{PopulateClassicModeDefaultValue "kubernetesExecHealthzSpec"}}
      "metadata": {
        "description": "The container spec for exechealthz-amd64."
      },
      "type": "string"
    },
    "kubernetesHeapsterSpec": {
      {{PopulateClassicModeDefaultValue "kubernetesHeapsterSpec"}}
      "metadata": {
        "description": "The container spec for heapster."
      },
      "type": "string"
    },
    "kubernetesTillerSpec": {
      {{PopulateClassicModeDefaultValue "kubernetesTillerSpec"}}
      "metadata": {
        "description": "The container spec for Helm Tiller."
      },
      "type": "string"
    },
    "kubernetesPodInfraContainerSpec": {
      {{PopulateClassicModeDefaultValue "kubernetesPodInfraContainerSpec"}}
      "metadata": {
        "description": "The container spec for pod infra."
      },
      "type": "string"
    },
    "kubernetesNodeStatusUpdateFrequency": {
      {{PopulateClassicModeDefaultValue "kubernetesNodeStatusUpdateFrequency"}}
      "metadata": {
        "description": "Kubelet config for node status update frequency interval."
      },
      "type": "string"
    },
    "kubernetesCtrlMgrNodeMonitorGracePeriod": {
      {{PopulateClassicModeDefaultValue "kubernetesCtrlMgrNodeMonitorGracePeriod"}}
      "metadata": {
        "description": "Kubernetes controller manager grace period for node status updates."
      },
      "type": "string"
    },
    "kubernetesCtrlMgrPodEvictionTimeout": {
      {{PopulateClassicModeDefaultValue "kubernetesCtrlMgrPodEvictionTimeout"}}
      "metadata": {
        "description": "Kubernetes controller manager pod eviction timeout."
      },
      "type": "string"
    },
    "kubernetesCtrlMgrRouteReconciliationPeriod": {
      {{PopulateClassicModeDefaultValue "kubernetesCtrlMgrRouteReconciliationPeriod"}}
      "metadata": {
        "description": "Kubernetes controller manager route reconciliation period."
      },
      "type": "string"
    },
    "cloudProviderBackoff": {
      {{PopulateClassicModeDefaultValue "cloudProviderBackoff"}}
      "metadata": {
        "description": "Enable cloudprovider backoff?"
      },
      "type": "string"
    },
    "cloudProviderBackoffRetries": {
      {{PopulateClassicModeDefaultValue "cloudProviderBackoffRetries"}}
      "metadata": {
        "description": "If backoff enabled, how many times to retry"
      },
      "type": "string"
    },
    "cloudProviderBackoffExponent": {
      {{PopulateClassicModeDefaultValue "cloudProviderBackoffExponent"}}
      "metadata": {
        "description": "If backoff enabled, retry exponent"
      },
      "type": "string"
    },
    "cloudProviderBackoffDuration": {
      {{PopulateClassicModeDefaultValue "cloudProviderBackoffDuration"}}
      "metadata": {
        "description": "If backoff enabled, how long until timeout"
      },
      "type": "string"
    },
    "cloudProviderBackoffJitter": {
      {{PopulateClassicModeDefaultValue "cloudProviderBackoffJitter"}}
      "metadata": {
        "description": "If backoff enabled, jitter factor between retries"
      },
      "type": "string"
    },
    "cloudProviderRatelimit": {
      {{PopulateClassicModeDefaultValue "cloudProviderRatelimit"}}
      "metadata": {
        "description": "Enable cloudprovider rate limiting?"
      },
      "type": "string"
    },
    "cloudProviderRatelimitQPS": {
      {{PopulateClassicModeDefaultValue "cloudProviderRatelimitQPS"}}
      "metadata": {
        "description": "If rate limiting enabled, target maximum QPS"
      },
      "type": "string"
    },
    "cloudProviderRatelimitBucket": {
      {{PopulateClassicModeDefaultValue "cloudProviderRatelimitBucket"}}
      "metadata": {
        "description": "If rate limiting enabled, bucket size"
      },
      "type": "string"
    },
    "kubernetesKubeDNSSpec": {
      {{PopulateClassicModeDefaultValue "kubernetesKubeDNSSpec"}}
      "metadata": {
        "description": "The container spec for kubedns-amd64."
      },
      "type": "string"
    },
    "kubernetesDNSMasqSpec": {
      {{PopulateClassicModeDefaultValue "kubernetesDNSMasqSpec"}}
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
    },
    "masterOffset": {
      "defaultValue": 0,
      "allowedValues": [
        0,
        1,
        2,
        3,
        4
      ],
      "metadata": {
        "description": "The offset into the master pool where to start creating master VMs.  This value can be from 0 to 4, but must be less than masterCount."
      },
      "type": "int"
    }


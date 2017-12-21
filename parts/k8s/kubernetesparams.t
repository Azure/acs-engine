{{if .HasAadProfile}}
    "aadServerAppId": {
      "metadata": {
        "description": "The server AAD application ID"
      },
      "type": "string"
    },
    "aadTenantId": {
      "defaultValue": "",
      "metadata": {
        "description": "The AAD tenant ID to use for authentication. If not specified, will use the tenant of the deployment subscription."
      },
      "type": "string"
    },
{{end}}
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
    "etcdServerCertificate": {
      "metadata": {
        "description": "The base 64 server certificate used on the master"
      }, 
      "type": "string"
    }, 
    "etcdServerPrivateKey": {
      "metadata": {
        "description": "The base 64 server private key used on the master."
      }, 
      "type": "securestring"
    }, 
    "etcdClientCertificate": {
      "metadata": {
        "description": "The base 64 server certificate used on the master"
      }, 
      "type": "string"
    }, 
    "etcdClientPrivateKey": {
      "metadata": {
        "description": "The base 64 server private key used on the master."
      }, 
      "type": "securestring"
    },
    "etcdPeerCertificates": {
      "metadata": {
        "description": "The base 64 server certificates used on the master"
      }, 
      "type": "array"
    }, 
    "etcdPeerPrivateKeys": {
      "metadata": {
        "description": "The base 64 server private keys used on the master."
      }, 
      "type": "array"
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
{{if IsHostedMaster}}
    "kubernetesEndpoint": {
      "metadata": {
        "description": "The Kubernetes API endpoint https://<kubernetesEndpoint>:443"
      },
      "type": "string"
    },
{{end}}
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
    "generatorCode": {
      {{PopulateClassicModeDefaultValue "generatorCode"}}
      "metadata": {
        "description": "The generator code used to identify the generator"
      },
      "type": "string"
    },
    "orchestratorName": {
      {{PopulateClassicModeDefaultValue "orchestratorName"}}
      "metadata": {
        "description": "The orchestrator name used to identify the orchestrator.  This must be no more than 3 digits in length, otherwise it will exceed Windows Naming"
      },
      "minLength": 3,
      "maxLength": 3,
      "type": "string"
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
    "kubeDNSServiceIP": {
      {{PopulateClassicModeDefaultValue "kubeDNSServiceIP"}}
      "metadata": {
        "description": "Kubernetes DNS IP"
      },
      "type": "string"
    },
    "kubeServiceCidr": {
      {{PopulateClassicModeDefaultValue "kubeServiceCidr"}}
      "metadata": {
        "description": "Kubernetes service address space"
      },
      "type": "string"
    },
    "kubernetesNonMasqueradeCidr": {
      "metadata": {
        "description": "kubernetesNonMasqueradeCidr cluster subnet"
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
    "kubernetesCcmImageSpec": {
      "defaultValue": "",
      "metadata": {
        "description": "The container spec for cloud-controller-manager."
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
    "kubernetesDashboardCPURequests": {
      {{PopulateClassicModeDefaultValue "kubernetesDashboardCPURequests"}}
      "metadata": {
        "description": "Dashboard CPU Requests."
      },
      "type": "string"
    },
    "kubernetesDashboardMemoryRequests": {
      {{PopulateClassicModeDefaultValue "kubernetesDashboardMemoryRequests"}}
      "metadata": {
        "description": "Dashboard Memory Requests."
      },
      "type": "string"
    },
    "kubernetesDashboardCPULimit": {
      {{PopulateClassicModeDefaultValue "kubernetesDashboardCPULimit"}}
      "metadata": {
        "description": "Dashboard CPU Limit."
      },
      "type": "string"
    },
    "kubernetesDashboardMemoryLimit": {
      {{PopulateClassicModeDefaultValue "kubernetesDashboardMemoryLimit"}}
      "metadata": {
        "description": "Dashboard Memory Limit."
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
    "kubernetesTillerCPURequests": {
      {{PopulateClassicModeDefaultValue "kubernetesTillerCPURequests"}}
      "metadata": {
        "description": "Helm Tiller CPU Requests."
      },
      "type": "string"
    },
    "kubernetesTillerMemoryRequests": {
      {{PopulateClassicModeDefaultValue "kubernetesTillerMemoryRequests"}}
      "metadata": {
        "description": "Helm Tiller Memory Requests."
      },
      "type": "string"
    },
    "kubernetesTillerCPULimit": {
      {{PopulateClassicModeDefaultValue "kubernetesTillerCPULimit"}}
      "metadata": {
        "description": "Helm Tiller CPU Limit."
      },
      "type": "string"
    },
    "kubernetesTillerMemoryLimit": {
      {{PopulateClassicModeDefaultValue "kubernetesTillerMemoryLimit"}}
      "metadata": {
        "description": "Helm Tiller Memory Limit."
      },
      "type": "string"
    },
    "kubernetesACIConnectorSpec": {
      {{PopulateClassicModeDefaultValue "kubernetesACIConnectorSpec"}}
      "metadata": {
        "description": "The container spec for ACI Connector."
      },
      "type": "string"
    },
    "kubernetesACIConnectorClientId": {
      {{PopulateClassicModeDefaultValue "kubernetesACIConnectorClientId"}}
      "metadata": {
        "description": "Client id for ACI Connector."
      },
      "type": "string"
    },
    "kubernetesACIConnectorClientKey": {
      {{PopulateClassicModeDefaultValue "kubernetesACIConnectorClientKey"}}
      "metadata": {
        "description": "Client key for ACI Connector."
      },
      "type": "string"
    },
    "kubernetesACIConnectorTenantId": {
      {{PopulateClassicModeDefaultValue "kubernetesACIConnectorTenantId"}}
      "metadata": {
        "description": "Tenant id for ACI Connector."
      },
      "type": "string"
    },
    "kubernetesACIConnectorSubscriptionId": {
      {{PopulateClassicModeDefaultValue "kubernetesACIConnectorSubscriptionId"}}
      "metadata": {
        "description": "Subscription id for ACI Connector."
      },
      "type": "string"
    },
    "kubernetesACIConnectorResourceGroup": {
      {{PopulateClassicModeDefaultValue "kubernetesACIConnectorResourceGroup"}}
      "metadata": {
        "description": "Resource group for ACI Connector."
      },
      "type": "string"
    },
    "kubernetesACIConnectorNodeName": {
      {{PopulateClassicModeDefaultValue "kubernetesACIConnectorNodeName"}}
      "metadata": {
        "description": "Node name for ACI Connector."
      },
      "type": "string"
    },
    "kubernetesACIConnectorOS": {
      {{PopulateClassicModeDefaultValue "kubernetesACIConnectorOS"}}
      "metadata": {
        "description": "OS for ACI Connector."
      },
      "type": "string"
    },
    "kubernetesACIConnectorTaint": {
      {{PopulateClassicModeDefaultValue "kubernetesACIConnectorTaint"}}
      "metadata": {
        "description": "Taint for ACI Connector."
      },
      "type": "string"
    },
    "kubernetesACIConnectorRegion": {
      {{PopulateClassicModeDefaultValue "kubernetesACIConnectorRegion"}}
      "metadata": {
        "description": "Region for ACI Connector."
      },
      "type": "string"
    },
    "kubernetesACIConnectorCPURequests": {
      {{PopulateClassicModeDefaultValue "kubernetesACIConnectorCPURequests"}}
      "metadata": {
        "description": "ACI Connector CPU Requests"
      },
      "type": "string"
    },
    "kubernetesACIConnectorMemoryRequests": {
      {{PopulateClassicModeDefaultValue "kubernetesACIConnectorMemoryRequests"}}
      "metadata": {
        "description": "ACI Connector Memory Requests"
      },
      "type": "string"
    },
    "kubernetesACIConnectorCPULimit": {
      {{PopulateClassicModeDefaultValue "kubernetesACIConnectorCPULimit"}}
      "metadata": {
        "description": "ACI Connector CPU Limit"
      },
      "type": "string"
    },
    "kubernetesACIConnectorMemoryLimit": {
      {{PopulateClassicModeDefaultValue "kubernetesACIConnectorMemoryLimit"}}
      "metadata": {
        "description": "ACI Connector Memory Limit"
      },
      "type": "string"
    },
    "kubernetesReschedulerSpec": {
      {{PopulateClassicModeDefaultValue "kubernetesReschedulerSpec"}}
      "metadata": {
        "description": "The container spec for rescheduler."
      },
      "type": "string"
    },
    "kubernetesReschedulerCPURequests": {
      {{PopulateClassicModeDefaultValue "kubernetesReschedulerCPURequests"}}
      "metadata": {
        "description": "Rescheduler CPU Requests."
      },
      "type": "string"
    },
    "kubernetesReschedulerMemoryRequests": {
      {{PopulateClassicModeDefaultValue "kubernetesReschedulerMemoryRequests"}}
      "metadata": {
        "description": "Rescheduler Memory Requests."
      },
      "type": "string"
    },
    "kubernetesReschedulerCPULimit": {
      {{PopulateClassicModeDefaultValue "kubernetesReschedulerCPULimit"}}
      "metadata": {
        "description": "Rescheduler CPU Limit."
      },
      "type": "string"
    },
    "kubernetesReschedulerMemoryLimit": {
      {{PopulateClassicModeDefaultValue "kubernetesReschedulerMemoryLimit"}}
      "metadata": {
        "description": "Rescheduler Memory Limit."
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
    "dockerEngineVersion": {
      {{PopulateClassicModeDefaultValue "dockerEngineVersion"}}
      "metadata": {
        "description": "The docker engine version to install."
      },
      "allowedValues": [
         "17.05.*",
         "17.04.*",
         "17.03.*",
         "1.13.*",
         "1.12.*",
         "1.11.*"
       ],
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
    "cniPluginsURL": {
      "defaultValue": "https://acs-mirror.azureedge.net/cni/cni-plugins-amd64-latest.tgz",
      "type": "string"
    },
    "vnetCniLinuxPluginsURL": {
      "defaultValue": "https://acs-mirror.azureedge.net/cni/azure-vnet-cni-linux-amd64-latest.tgz",
      "type": "string"
    },
    "vnetCniWindowsPluginsURL": {
      "defaultValue": "https://acs-mirror.azureedge.net/cni/azure-vnet-cni-windows-amd64-latest.zip",
      "type": "string"
    },
    "maxPods": {
      "defaultValue": 110,
      "metadata": {
        "description": "The maximum number of pods per node."
      },
      "type": "int"
    },
    "vnetCidr": {
      "defaultValue": "10.0.0.0/8",
      "metadata": {
        "description": "Cluster vnet cidr"
      },
      "type": "string"
    },
    "gcHighThreshold": {
      "defaultValue": 85,
      "metadata": {
        "description": "High Threshold for Image Garbage collection on each node"
      },
      "type": "int"
    },
    "gcLowThreshold": {
      "defaultValue": 80,
      "metadata": {
        "description": "Low Threshold for Image Garbage collection on each node."
      },
      "type": "int"
    },
{{ if not UseManagedIdentity }}
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
{{ end }}
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
    },
    "etcdDiskSizeGB": {
      {{PopulateClassicModeDefaultValue "etcdDiskSizeGB"}}
      "metadata": {
        "description": "Size in GB to allocate for etcd volume"
      },
      "type": "string"
    },
    "etcdDownloadURLBase": {
      {{PopulateClassicModeDefaultValue "etcdDownloadURLBase"}}
      "metadata": {
        "description": "etcd image base URL"
      },
      "type": "string"
    },
    "etcdVersion": {
      {{PopulateClassicModeDefaultValue "etcdVersion"}}
      "metadata": {
        "description": "etcd version"
      },
      "type": "string"
    }

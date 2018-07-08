{{if .HasAadProfile}}
    "aadTenantId": {
      "defaultValue": "",
      "metadata": {
        "description": "The AAD tenant ID to use for authentication. If not specified, will use the tenant of the deployment subscription."
      },
      "type": "string"
    },
    "aadAdminGroupId": {
      "defaultValue": "",
      "metadata": {
        "description": "The AAD default Admin group Object ID used to create a cluster-admin RBAC role."
      },
      "type": "string"
    },
{{end}}
{{if IsHostedMaster}}
    "kubernetesEndpoint": {
      "metadata": {
        "description": "The Kubernetes API endpoint https://<kubernetesEndpoint>:443"
      },
      "type": "string"
    },
{{else}}
{{if not IsOpenShift}}
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
    "etcdPeerCertificate0": {
      "metadata": {
        "description": "The base 64 server certificates used on the master"
      },
      "type": "string"
    },
    "etcdPeerPrivateKey0": {
      "metadata": {
        "description": "The base 64 server private keys used on the master."
      },
      "type": "securestring"
    },
    {{if ge .MasterProfile.Count 3}}
      "etcdPeerCertificate1": {
        "metadata": {
          "description": "The base 64 server certificates used on the master"
        },
        "type": "string"
      },
      "etcdPeerCertificate2": {
        "metadata": {
          "description": "The base 64 server certificates used on the master"
        },
        "type": "string"
      },
      "etcdPeerPrivateKey1": {
        "metadata": {
          "description": "The base 64 server private keys used on the master."
        },
        "type": "securestring"
      },
      "etcdPeerPrivateKey2": {
        "metadata": {
          "description": "The base 64 server private keys used on the master."
        },
        "type": "securestring"
      },
      {{if ge .MasterProfile.Count 5}}
        "etcdPeerCertificate3": {
          "metadata": {
            "description": "The base 64 server certificates used on the master"
          },
          "type": "string"
        },
        "etcdPeerCertificate4": {
          "metadata": {
            "description": "The base 64 server certificates used on the master"
          },
          "type": "string"
        },
        "etcdPeerPrivateKey3": {
          "metadata": {
            "description": "The base 64 server private keys used on the master."
          },
          "type": "securestring"
        },
        "etcdPeerPrivateKey4": {
          "metadata": {
            "description": "The base 64 server private keys used on the master."
          },
          "type": "securestring"
        },
      {{end}}
    {{end}}
{{end}}
{{end}}
{{if not IsOpenShift}}
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
{{end}}
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
    "kubernetesKubeletClusterDomain": {
      "metadata": {
        "description": "--cluster-domain Kubelet config"
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
    "kubernetesMetricsServerSpec": {
      {{PopulateClassicModeDefaultValue "kubernetesMetricsServerSpec"}}
      "metadata": {
        "description": "The container spec for Metrics Server."
      },
      "type": "string"
    },
    "kubernetesNVIDIADevicePluginSpec": {
      {{PopulateClassicModeDefaultValue "kubernetesNVIDIADevicePluginSpec"}}
      "metadata": {
        "description": "The container spec for NVIDIA Device Plugin."
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
    "kubernetesTillerMaxHistory": {
      {{PopulateClassicModeDefaultValue "kubernetesTillerMaxHistory"}}
      "metadata": {
        "description": "Helm Tiller Max History to Store. '0' for no limit."
      },
      "type": "string"
    },
    "kubernetesACIConnectorEnabled": {
      "defaultValue": false,
      "metadata": {
        "description": "ACI Connector Status"
      },
      "type": "bool"
    },
    "kubernetesACIConnectorSpec": {
      {{PopulateClassicModeDefaultValue "kubernetesACIConnectorSpec"}}
      "metadata": {
        "description": "The container spec for ACI Connector."
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
    "kubernetesClusterAutoscalerSpec": {
      {{PopulateClassicModeDefaultValue "kubernetesClusterAutoscalerSpec"}}
      "metadata": {
        "description": "The container spec for the cluster autoscaler."
      },
      "type": "string"
    },
    "kubernetesClusterAutoscalerAzureCloud": {
      {{PopulateClassicModeDefaultValue "kubernetesClusterAutoscalerAzureCloud"}}
      "metadata": {
        "description": "Name of the Azure cloud for the cluster autoscaler."
      },
      "type": "string"
    },
    "kubernetesClusterAutoscalerCPULimit": {
      {{PopulateClassicModeDefaultValue "kubernetesClusterAutoscalerCPULimit"}}
      "metadata": {
        "description": "Cluster autoscaler cpu limit"
      },
      "type": "string"
    },
    "kubernetesClusterAutoscalerMemoryLimit": {
      {{PopulateClassicModeDefaultValue "kubernetesClusterAutoscalerMemoryLimit"}}
      "metadata": {
        "description": "Cluster autoscaler memory limit"
      },
      "type": "string"
    },
    "kubernetesClusterAutoscalerCPURequests": {
      {{PopulateClassicModeDefaultValue "kubernetesClusterAutoscalerCPURequests"}}
      "metadata": {
        "description": "Cluster autoscaler cpu requests"
      },
      "type": "string"
    },
    "kubernetesClusterAutoscalerMemoryRequests": {
      {{PopulateClassicModeDefaultValue "kubernetesClusterAutoscalerMemoryRequests"}}
      "metadata": {
        "description": "Cluster autoscaler memory requests"
      },
      "type": "string"
    },
    "kubernetesClusterAutoscalerMinNodes": {
      {{PopulateClassicModeDefaultValue "kubernetesClusterAutoscalerMinNodes"}}
      "metadata": {
        "description": "Cluster autoscaler min nodes"
      },
      "type": "string"
    },
    "kubernetesClusterAutoscalerMaxNodes": {
      {{PopulateClassicModeDefaultValue "kubernetesClusterAutoscalerMaxNodes"}}
      "metadata": {
        "description": "Cluster autoscaler max nodes"
      },
      "type": "string"
    },
    "kubernetesClusterAutoscalerEnabled": {
      "defaultValue": false,
      "metadata": {
        "description": "Cluster autoscaler status"
      },
      "type": "bool"
    },
    "kubernetesClusterAutoscalerUseManagedIdentity": {
      {{PopulateClassicModeDefaultValue "kubernetesClusterAutoscalerUseManagedIdentity"}}
      "metadata": {
        "description": "Managed identity for the cluster autoscaler addon"
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
    {{if not IsOpenShift}}
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
    {{end}}
    "networkPolicy": {
      "defaultValue": "{{.OrchestratorProfile.KubernetesConfig.NetworkPolicy}}",
      "metadata": {
        "description": "The network policy enforcement to use (calico|cilium); 'none' and 'azure' here for backwards compatibility"
      },
      "allowedValues": [
        "",
        "none",
        "azure",
        "calico",
        "cilium"
      ],
      "type": "string"
    },
    "networkPlugin": {
      "defaultValue": "{{.OrchestratorProfile.KubernetesConfig.NetworkPlugin}}",
      "metadata": {
        "description": "The network plugin to use for Kubernetes (kubenet|azure|flannel|cilium)"
      },
      "allowedValues": [
        "kubenet",
        "azure",
        "flannel",
        "cilium"
      ],
      "type": "string"
    },
    "containerRuntime": {
      "defaultValue": "{{.OrchestratorProfile.KubernetesConfig.ContainerRuntime}}",
      "metadata": {
        "description": "The container runtime to use (docker|clear-containers|containerd)"
      },
      "allowedValues": [
        "docker",
        "clear-containers",
        "containerd"
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
      "defaultValue": 30,
      "metadata": {
        "description": "This param has been deprecated."
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
    "omsAgentVersion": {
      "defaultValue": "",
      "metadata": {
        "description": "OMS agent version for Container Monitoring."
      },
      "type": "string"
    },
    "omsAgentDockerProviderVersion": {
      "defaultValue": "",
      "metadata": {
        "description": "Docker provider version for Container Monitoring."
      },
      "type": "string"
    },
    "omsAgentImage": {
      "defaultValue": "",
      "metadata": {
        "description": "OMS agent image for Container Monitoring."
      },
      "type": "string"
    },
    "omsAgentWorkspaceGuid": {
      "defaultValue": "",
      "metadata": {
        "description": "OMS workspace guid"
      },
      "type": "string"
    },
    "omsAgentWorkspaceKey": {
      "defaultValue": "",
      "metadata": {
        "description": "OMS workspace key"
      },
      "type": "string"
    },
    "kubernetesOMSAgentCPURequests": {
      "defaultValue": "",
      "metadata": {
        "description": "OMS Agent CPU requests resource limit"
      },
      "type": "string"
    },
    "kubernetesOMSAgentMemoryRequests": {
      "defaultValue": "",
      "metadata": {
        "description": "OMS Agent memory requests resource limit"
      },
      "type": "string"
    },
    "kubernetesOMSAgentCPULimit": {
      "defaultValue": "",
      "metadata": {
        "description": "OMS Agent CPU limit resource limit"
      },
      "type": "string"
    },
    "kubernetesOMSAgentMemoryLimit": {
      "defaultValue": "",
      "metadata": {
        "description": "OMS Agent memory limit resource limit"
      },
      "type": "string"
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
    },
    "etcdEncryptionKey": {
      "metadata": {
        "description": "Encryption at rest key for etcd"
      },
      "type": "string"
    }
{{if ProvisionJumpbox}}
    ,"jumpboxVMName": {
      "metadata": {
        "description": "jumpbox VM Name"
      },
      "type": "string"
    },
    "jumpboxVMSize": {
      {{GetMasterAllowedSizes}}
      "metadata": {
        "description": "The size of the Virtual Machine. Required"
      },
      "type": "string"
    },
    "jumpboxOSDiskSizeGB": {
      {{PopulateClassicModeDefaultValue "jumpboxOSDiskSizeGB"}}
      "metadata": {
        "description": "Size in GB to allocate to the private cluster jumpbox VM OS."
      },
      "type": "int"
    },
    "jumpboxPublicKey": {
      "metadata": {
        "description": "SSH public key used for auth to the private cluster jumpbox"
      },
      "type": "string"
    },
    "jumpboxUsername": {
      "metadata": {
        "description": "Username for the private cluster jumpbox"
      },
      "type": "string"
    },
    "jumpboxStorageProfile": {
      "metadata": {
        "description": "Storage Profile for the private cluster jumpbox"
      },
      "type": "string"
    }
{{end}}
{{if HasCustomSearchDomain}}
    ,"searchDomainName": {
      "defaultValue": "",
      "metadata": {
        "description": "Custom Search Domain name."
      },
      "type": "string"
    },
    "searchDomainRealmUser": {
      "defaultValue": "",
      "metadata": {
        "description": "Windows server AD user name to join the Linux Machines with active directory and be able to change dns registries."
      },
      "type": "string"
    },
    "searchDomainRealmPassword": {
      "defaultValue": "",
      "metadata": {
        "description": "Windows server AD user password to join the Linux Machines with active directory and be able to change dns registries."
      },
      "type": "securestring"
    }
{{end}}
{{if HasCustomNodesDNS}}
    ,"dnsServer": {
      "defaultValue": "",
      "metadata": {
        "description": "DNS Server IP"
      },
      "type": "string"
    }
{{end}}

{{if EnableEncryptionWithExternalKms}}
   ,
   {{if not UseManagedIdentity}}
   "servicePrincipalObjectId": {
      "metadata": {
        "description": "Object ID (used by cloudprovider)"
      },
      "type": "securestring"
    },
    {{end}}
    "clusterKeyVaultSku": {
       "type": "string",
       "defaultValue": "Standard",
       "allowedValues": [
         "Standard",
         "Premium"
       ],
       "metadata": {
         "description": "SKU for the key vault used by the cluster"
       }
     }
 {{end}}
 {{if IsAzureCNI}}
    ,"AzureCNINetworkMonitorImageURL": {
      "defaultValue": "",
      "metadata": {
        "description": "Azure CNI networkmonitor Image URL"
      },
      "type": "string"
    }
 {{end}}

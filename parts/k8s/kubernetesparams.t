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
      "metadata": {
        "description": "The generator code used to identify the generator"
      },
      "type": "string"
    },
    "orchestratorName": {
      "metadata": {
        "description": "The orchestrator name used to identify the orchestrator.  This must be no more than 3 digits in length, otherwise it will exceed Windows Naming"
      },
      "minLength": 3,
      "maxLength": 3,
      "type": "string"
    },
    "dockerBridgeCidr": {
      "metadata": {
        "description": "Docker bridge network IP address and subnet"
      },
      "type": "string"
    },
    "kubeClusterCidr": {
      "metadata": {
        "description": "Kubernetes cluster subnet"
      },
      "type": "string"
    },
    "kubeDNSServiceIP": {
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
{{if not IsHostedMaster}}
    "kubernetesNonMasqueradeCidr": {
      "metadata": {
        "description": "kubernetesNonMasqueradeCidr cluster subnet"
      },
      "defaultValue": "{{GetDefaultVNETCIDR}}",
      "type": "string"
    },
{{end}}
    "kubernetesKubeletClusterDomain": {
      "metadata": {
        "description": "--cluster-domain Kubelet config"
      },
      "type": "string"
    },
    "kubernetesHyperkubeSpec": {
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
      "metadata": {
        "description": "The container spec for hyperkube."
      },
      "type": "string"
    },
    "kubernetesAddonResizerSpec": {
      "metadata": {
        "description": "The container spec for addon-resizer."
      },
      "type": "string"
    },
{{if .OrchestratorProfile.KubernetesConfig.IsDashboardEnabled}}
    "kubernetesDashboardSpec": {
      "metadata": {
        "description": "The container spec for kubernetes-dashboard-amd64."
      },
      "type": "string"
    },
    "kubernetesDashboardCPURequests": {
      "metadata": {
        "description": "Dashboard CPU Requests."
      },
      "type": "string"
    },
    "kubernetesDashboardMemoryRequests": {
      "metadata": {
        "description": "Dashboard Memory Requests."
      },
      "type": "string"
    },
    "kubernetesDashboardCPULimit": {
      "metadata": {
        "description": "Dashboard CPU Limit."
      },
      "type": "string"
    },
    "kubernetesDashboardMemoryLimit": {
      "metadata": {
        "description": "Dashboard Memory Limit."
      },
      "type": "string"
    },
{{end}}
    "enableAggregatedAPIs": {
      "metadata": {
        "description": "Enable aggregated API on master nodes"
      },
      "defaultValue": false,
      "type": "bool"
    },
    "kubernetesExecHealthzSpec": {
      "metadata": {
        "description": "The container spec for exechealthz-amd64."
      },
      "type": "string"
    },
    "kubernetesDNSSidecarSpec": {
      "metadata": {
        "description": "The container spec for k8s-dns-sidecar-amd64."
      },
      "type": "string"
    },
    "kubernetesHeapsterSpec": {
      "metadata": {
        "description": "The container spec for heapster."
      },
      "type": "string"
    },
{{if .OrchestratorProfile.IsMetricsServerEnabled}}
    "kubernetesMetricsServerSpec": {
      "metadata": {
        "description": "The container spec for Metrics Server."
      },
      "type": "string"
    },
{{end}}
{{if .IsNVIDIADevicePluginEnabled}}
    "kubernetesNVIDIADevicePluginSpec": {
      "metadata": {
        "description": "The container spec for NVIDIA Device Plugin."
      },
      "type": "string"
    },
    "kubernetesNVIDIADevicePluginCPURequests": {
      "metadata": {
        "description": "NVIDIA Device Plugin CPU Requests"
      },
      "type": "string"
    },
    "kubernetesNVIDIADevicePluginMemoryRequests": {
      "metadata": {
        "description": "NVIDIA Device Plugin Memory Requests"
      },
      "type": "string"
    },
    "kubernetesNVIDIADevicePluginCPULimit": {
      "metadata": {
        "description": "NVIDIA Device Plugin CPU Limit"
      },
      "type": "string"
    },
    "kubernetesNVIDIADevicePluginMemoryLimit": {
      "metadata": {
        "description": "NVIDIA Device Plugin Memory Limit"
      },
      "type": "string"
    },
{{end}}
{{if .OrchestratorProfile.KubernetesConfig.IsTillerEnabled}}
    "kubernetesTillerSpec": {
      "metadata": {
        "description": "The container spec for Helm Tiller."
      },
      "type": "string"
    },
    "kubernetesTillerCPURequests": {
      "metadata": {
        "description": "Helm Tiller CPU Requests."
      },
      "type": "string"
    },
    "kubernetesTillerMemoryRequests": {
      "metadata": {
        "description": "Helm Tiller Memory Requests."
      },
      "type": "string"
    },
    "kubernetesTillerCPULimit": {
      "metadata": {
        "description": "Helm Tiller CPU Limit."
      },
      "type": "string"
    },
    "kubernetesTillerMemoryLimit": {
      "metadata": {
        "description": "Helm Tiller Memory Limit."
      },
      "type": "string"
    },
    "kubernetesTillerMaxHistory": {
      "metadata": {
        "description": "Helm Tiller Max History to Store. '0' for no limit."
      },
      "type": "string"
    },
{{end}}
{{if .OrchestratorProfile.KubernetesConfig.IsAADPodIdentityEnabled}}
    "kubernetesAADPodIdentityEnabled": {
      "defaultValue": false,
      "metadata": {
        "description": "AAD Pod Identity status"
      },
      "type": "bool"
    },
{{end}}
    "kubernetesACIConnectorEnabled": {
      "metadata": {
        "description": "ACI Connector Status"
      },
      "type": "bool"
    },
{{if .OrchestratorProfile.KubernetesConfig.IsACIConnectorEnabled}}
    "kubernetesACIConnectorSpec": {
      "metadata": {
        "description": "The container spec for ACI Connector."
      },
      "type": "string"
    },
    "kubernetesACIConnectorNodeName": {
      "metadata": {
        "description": "Node name for ACI Connector."
      },
      "type": "string"
    },
    "kubernetesACIConnectorOS": {
      "metadata": {
        "description": "OS for ACI Connector."
      },
      "type": "string"
    },
    "kubernetesACIConnectorTaint": {
      "metadata": {
        "description": "Taint for ACI Connector."
      },
      "type": "string"
    },
    "kubernetesACIConnectorRegion": {
      "metadata": {
        "description": "Region for ACI Connector."
      },
      "type": "string"
    },
    "kubernetesACIConnectorCPURequests": {
      "metadata": {
        "description": "ACI Connector CPU Requests"
      },
      "type": "string"
    },
    "kubernetesACIConnectorMemoryRequests": {
      "metadata": {
        "description": "ACI Connector Memory Requests"
      },
      "type": "string"
    },
    "kubernetesACIConnectorCPULimit": {
      "metadata": {
        "description": "ACI Connector CPU Limit"
      },
      "type": "string"
    },
    "kubernetesACIConnectorMemoryLimit": {
      "metadata": {
        "description": "ACI Connector Memory Limit"
      },
      "type": "string"
    },
{{end}}
    "kubernetesClusterAutoscalerEnabled": {
      "metadata": {
        "description": "Cluster autoscaler status"
      },
      "type": "bool"
    },
{{if .OrchestratorProfile.KubernetesConfig.IsClusterAutoscalerEnabled}}
    "kubernetesClusterAutoscalerSpec": {
      "metadata": {
        "description": "The container spec for the cluster autoscaler."
      },
      "type": "string"
    },
    "kubernetesClusterAutoscalerAzureCloud": {
      "metadata": {
        "description": "Name of the Azure cloud for the cluster autoscaler."
      },
      "type": "string"
    },
    "kubernetesClusterAutoscalerCPULimit": {
      "metadata": {
        "description": "Cluster autoscaler cpu limit"
      },
      "type": "string"
    },
    "kubernetesClusterAutoscalerMemoryLimit": {
      "metadata": {
        "description": "Cluster autoscaler memory limit"
      },
      "type": "string"
    },
    "kubernetesClusterAutoscalerCPURequests": {
      "metadata": {
        "description": "Cluster autoscaler cpu requests"
      },
      "type": "string"
    },
    "kubernetesClusterAutoscalerMemoryRequests": {
      "metadata": {
        "description": "Cluster autoscaler memory requests"
      },
      "type": "string"
    },
    "kubernetesClusterAutoscalerMinNodes": {
      "metadata": {
        "description": "Cluster autoscaler min nodes"
      },
      "type": "string"
    },
    "kubernetesClusterAutoscalerMaxNodes": {
      "metadata": {
        "description": "Cluster autoscaler max nodes"
      },
      "type": "string"
    },
    "kubernetesClusterAutoscalerUseManagedIdentity": {
      "metadata": {
        "description": "Managed identity for the cluster autoscaler addon"
      },
      "type": "string"
    },
{{end}}
     "flexVolumeDriverConfig": {
      "type": "object",
      "defaultValue": {
        "kubernetesBlobfuseFlexVolumeInstallerCPURequests": "50m",
        "kubernetesBlobfuseFlexVolumeInstallerMemoryRequests": "10Mi",
        "kubernetesBlobfuseFlexVolumeInstallerCPULimit": "50m",
        "kubernetesBlobfuseFlexVolumeInstallerMemoryLimit": "10Mi",
        "kubernetesSMBFlexVolumeInstallerCPURequests": "50m",
        "kubernetesSMBFlexVolumeInstallerMemoryRequests": "10Mi",
        "kubernetesSMBFlexVolumeInstallerCPULimit": "50m",
        "kubernetesSMBFlexVolumeInstallerMemoryLimit": "10Mi"
      }
    },
{{if .OrchestratorProfile.KubernetesConfig.IsKeyVaultFlexVolumeEnabled}}
    "kubernetesKeyVaultFlexVolumeInstallerCPURequests": {
      "metadata": {
        "description": "Key Vault FlexVolume Installer CPU Requests"
      },
      "type": "string"
    },
    "kubernetesKeyVaultFlexVolumeInstallerMemoryRequests": {
      "metadata": {
        "description": "Key Vault FlexVolume Installer Memory Requests"
      },
      "type": "string"
    },
    "kubernetesKeyVaultFlexVolumeInstallerCPULimit": {
      "metadata": {
        "description": "Key Vault FlexVolume Installer CPU Limit"
      },
      "type": "string"
    },
    "kubernetesKeyVaultFlexVolumeInstallerMemoryLimit": {
      "metadata": {
        "description": "Key Vault FlexVolume Installer Memory Limit"
      },
      "type": "string"
    },
{{end}}
{{if .OrchestratorProfile.KubernetesConfig.IsReschedulerEnabled}}
    "kubernetesReschedulerSpec": {
      "metadata": {
        "description": "The container spec for rescheduler."
      },
      "type": "string"
    },
    "kubernetesReschedulerCPURequests": {
      "metadata": {
        "description": "Rescheduler CPU Requests."
      },
      "type": "string"
    },
    "kubernetesReschedulerMemoryRequests": {
      "metadata": {
        "description": "Rescheduler Memory Requests."
      },
      "type": "string"
    },
    "kubernetesReschedulerCPULimit": {
      "metadata": {
        "description": "Rescheduler CPU Limit."
      },
      "type": "string"
    },
    "kubernetesReschedulerMemoryLimit": {
      "metadata": {
        "description": "Rescheduler Memory Limit."
      },
      "type": "string"
    },
{{end}}
{{if .OrchestratorProfile.KubernetesConfig.IsIPMasqAgentEnabled}}
    "kubernetesIPMasqAgentCPURequests": {
      "metadata": {
        "description": "IP Masq Agent CPU Requests"
      },
      "type": "string"
    },
    "kubernetesIPMasqAgentMemoryRequests": {
      "metadata": {
        "description": "IP Masq Agent Memory Requests"
      },
      "type": "string"
    },
    "kubernetesIPMasqAgentCPULimit": {
      "metadata": {
        "description": "IP Masq Agent CPU Limit"
      },
      "type": "string"
    },
    "kubernetesIPMasqAgentMemoryLimit": {
      "metadata": {
        "description": "IP Masq Agent Memory Limit"
      },
      "type": "string"
    },
{{end}}
    "kubernetesPodInfraContainerSpec": {
      "metadata": {
        "description": "The container spec for pod infra."
      },
      "type": "string"
    },
    "cloudproviderConfig": {
      "type": "object",
      "defaultValue": {
        "cloudProviderBackoff": true,
        "cloudProviderBackoffRetries": 10,
        "cloudProviderBackoffJitter": "0",
        "cloudProviderBackoffDuration": 0,
        "cloudProviderBackoffExponent": "0",
        "cloudProviderRateLimit": false,
        "cloudProviderRateLimitQPS": "0",
        "cloudProviderRateLimitBucket": 0
      }
    },
    "kubernetesKubeDNSSpec": {
      "metadata": {
        "description": "The container spec for kubedns-amd64."
      },
      "type": "string"
    },
    "kubernetesCoreDNSSpec": {
      "metadata": {
        "description": "The container spec for coredns"
      },
      "type": "string"
    },
    "kubernetesDNSMasqSpec": {
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
        "description": "The container runtime to use (docker|clear-containers|kata-containers|containerd)"
      },
      "allowedValues": [
        "docker",
        "clear-containers",
        "kata-containers",
        "containerd"
      ],
      "type": "string"
    },
    "containerdDownloadURLBase": {
      "defaultValue": "https://storage.googleapis.com/cri-containerd-release/",
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
      "defaultValue": "{{GetDefaultVNETCIDR}}",
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
    "kuberneteselbsvcname": {
      "defaultValue": "",
      "metadata": {
        "description": "elb service for standard lb"
      },
      "type": "string"
    },
{{if .OrchestratorProfile.KubernetesConfig.IsContainerMonitoringEnabled}}
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
{{end}}
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
      "metadata": {
        "description": "Size in GB to allocate for etcd volume"
      },
      "type": "string"
    },
    "etcdDownloadURLBase": {
      "metadata": {
        "description": "etcd image base URL"
      },
      "type": "string"
    },
    "etcdVersion": {
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

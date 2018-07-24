# Microsoft Azure Container Service Engine - Cluster Definition

## Cluster Defintions for apiVersion "vlabs"

Here are the cluster definitions for apiVersion "vlabs":

### apiVersion

| Name       | Required | Description                                                   |
| ---------- | -------- | ------------------------------------------------------------- |
| apiVersion | yes      | The version of the template. For "vlabs" the value is "vlabs" |

### orchestratorProfile

`orchestratorProfile` describes the orchestrator settings.

| Name             | Required | Description                                     |
| ---------------- | -------- | ----------------------------------------------- |
| orchestratorType | yes      | Specifies the orchestrator type for the cluster |

Here are the valid values for the orchestrator types:

1.  `DCOS` - this represents the [DC/OS orchestrator](dcos.md). [Older releases of DCOS 1.8 may be specified](../examples/dcos-releases).
2.  `Kubernetes` - this represents the [Kubernetes orchestrator](kubernetes.md).
3.  `Swarm` - this represents the [Swarm orchestrator](swarm.md).
4.  `Swarm Mode` - this represents the [Swarm Mode orchestrator](swarmmode.md).
5.  `OpenShift` - this represents the [OpenShift orchestrator](openshift.md).

### kubernetesConfig

`kubernetesConfig` describes Kubernetes specific configuration.

| Name                            | Required | Description                                                                                                                                                                                                                                                                                                                                                                                                   |
| ------------------------------- | -------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| addons                          | no       | Configure various Kubernetes addons configuration (currently supported: tiller, kubernetes-dashboard). See `addons` configuration below                                                                                                                                                                                                                                                                       |
| apiServerConfig                 | no       | Configure various runtime configuration for apiserver. See `apiServerConfig` [below](#feat-apiserver-config)                                                                                                                                                                                                                                                                                                  |
| cloudControllerManagerConfig    | no       | Configure various runtime configuration for cloud-controller-manager. See `cloudControllerManagerConfig` [below](#feat-cloud-controller-manager-config)                                                                                                                                                                                                                                                       |
| clusterSubnet                   | no       | The IP subnet used for allocating IP addresses for pod network interfaces. The subnet must be in the VNET address space. Default value is 10.244.0.0/16                                                                                                                                                                                                                                                       |
| containerRuntime                | no       | The container runtime to use as a backend. The default is `docker`. The other options are `clear-containers`, `kata-containers`, and `containerd`                                                                                                                                                                                                                                                             |
| controllerManagerConfig         | no       | Configure various runtime configuration for controller-manager. See `controllerManagerConfig` [below](#feat-controller-manager-config)                                                                                                                                                                                                                                                                        |
| customWindowsPackageURL         | no       | Configure custom windows Kubernetes release package URL for deployment on Windows                                                                                                                                                                                                                                                                                                                             |
| dnsServiceIP                    | no       | IP address for kube-dns to listen on. If specified must be in the range of `serviceCidr`                                                                                                                                                                                                                                                                                                                      |
| dockerBridgeSubnet              | no       | The specific IP and subnet used for allocating IP addresses for the docker bridge network created on the kubernetes master and agents. Default value is 172.17.0.1/16. This value is used to configure the docker daemon using the [--bip flag](https://docs.docker.com/engine/userguide/networking/default_network/custom-docker0)                                                                           |
| dockerEngineVersion             | no       | Which version of docker-engine to use in your cluster, e.g. "17.03.\*"                                                                                                                                                                                                                                                                                                                                        |
| enableAggregatedAPIs            | no       | Enable [Kubernetes Aggregated APIs](https://kubernetes.io/docs/concepts/api-extension/apiserver-aggregation/).This is required by [Service Catalog](https://github.com/kubernetes-incubator/service-catalog/blob/master/README.md). (boolean - default is true for k8s versions greater or equal to 1.9.0, false otherwise)                                                                                                                                              |
| enableDataEncryptionAtRest      | no       | Enable [kubernetes data encryption at rest](https://kubernetes.io/docs/tasks/administer-cluster/encrypt-data/).This is currently an alpha feature. (boolean - default == false)                                                                                                                                                                                                                               |
| enableEncryptionWithExternalKms | no       | Enable [kubernetes data encryption at rest with external KMS](https://kubernetes.io/docs/tasks/administer-cluster/encrypt-data/).This is currently an alpha feature. (boolean - default == false)                                                                                                                                                                                                             |
| enablePodSecurityPolicy         | no       | Enable [kubernetes pod security policy](https://kubernetes.io/docs/concepts/policy/pod-security-policy/).This is currently a beta feature. (boolean - default == false)                                                                                                                                                                                                                                       |
| enableRbac                      | no       | Enable [Kubernetes RBAC](https://kubernetes.io/docs/admin/authorization/rbac/) (boolean - default == true)                                                                                                                                                                                                                                                                                                    |
| etcdDiskSizeGB                  | no       | Size in GB to assign to etcd data volume. Defaults (if no user value provided) are: 256 GB for clusters up to 3 nodes; 512 GB for clusters with between 4 and 10 nodes; 1024 GB for clusters with between 11 and 20 nodes; and 2048 GB for clusters with more than 20 nodes                                                                                                                                   |
| etcdEncryptionKey               | no       | Enryption key to be used if enableDataEncryptionAtRest is enabled. Defaults to a random, generated, key                                                                                                                                                                                                                                                                                                       |
| gcHighThreshold                 | no       | Sets the --image-gc-high-threshold value on the kublet configuration. Default is 85. [See kubelet Garbage Collection](https://kubernetes.io/docs/concepts/cluster-administration/kubelet-garbage-collection/)                                                                                                                                                                                                 |
| gcLowThreshold                  | no       | Sets the --image-gc-low-threshold value on the kublet configuration. Default is 80. [See kubelet Garbage Collection](https://kubernetes.io/docs/concepts/cluster-administration/kubelet-garbage-collection/)                                                                                                                                                                                                  |
| kubeletConfig                   | no       | Configure various runtime configuration for kubelet. See `kubeletConfig` [below](#feat-kubelet-config)                                                                                                                                                                                                                                                                                                        |
| kubernetesImageBase             | no       | Specifies the base URL (everything preceding the actual image filename) of the kubernetes hyperkube image to use for cluster deployment, e.g., `k8s.gcr.io/`                                                                                                                                                                                                                                     |
| networkPlugin                   | no       | Specifies the network plugin implementation for the cluster. Valid values are:<br>`"azure"` (default), which provides an Azure native networking experience <br>`"kubenet"` for k8s software networking implementation. <br> `"flannel"` for using CoreOS Flannel <br> `"cilium"` for using the default Cilium CNI IPAM                                                                                       |
| networkPolicy                   | no       | Specifies the network policy enforcement tool for the cluster (currently Linux-only). Valid values are:<br>`"calico"` for Calico network policy.<br>`"cilium"` for cilium network policy (Lin), and `"azure"` (experimental) for Azure CNI-compliant network policy (note: Azure CNI-compliant network policy requires explicit `"networkPlugin": "azure"` configuration as well).<br>See [network policy examples](../examples/networkpolicy) for more information.                                                                                                                                  |
| privateCluster                  | no       | Build a cluster without public addresses assigned. See `privateClusters` [below](#feat-private-cluster).                                                                                                                                                                                                                                                                                                      |
| schedulerConfig                 | no       | Configure various runtime configuration for scheduler. See `schedulerConfig` [below](#feat-scheduler-config)                                                                                                                                                                                                                                                                                                  |
| serviceCidr                     | no       | IP range for Service IPs, Default is "10.0.0.0/16". This range is never routed outside of a node so does not need to lie within clusterSubnet or the VNET                                                                                                                                                                                                                                                     |
| useInstanceMetadata             | no       | Use the Azure cloudprovider instance metadata service for appropriate resource discovery operations. Default is `true`                                                                                                                                                                                                                                                                                        |
| useManagedIdentity              | no       | Includes and uses MSI identities for all interactions with the Azure Resource Manager (ARM) API. Instead of using a static service principal written to /etc/kubernetes/azure.json, Kubernetes will use a dynamic, time-limited token fetched from the MSI extension running on master and agent nodes. This support is currently alpha and requires Kubernetes v1.9.1 or newer. (boolean - default == false) |

#### addons

`addons` describes various addons configuration. It is a child property of `kubernetesConfig`. Below is a list of currently available addons:

| Name of addon                                                         | Enabled by default? | How many containers | Description                                                                                                                                                         |
| --------------------------------------------------------------------- | ------------------- | ------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| tiller                                                                | true                | 1                   | Delivers the Helm server-side component: tiller. See https://github.com/kubernetes/helm for more info                                                               |
| kubernetes-dashboard                                                  | true                | 1                   | Delivers the Kubernetes dashboard component. See https://github.com/kubernetes/dashboard for more info                                                              |
| rescheduler                                                           | false               | 1                   | Delivers the Kubernetes rescheduler component                                                                                                                       |
| [cluster-autoscaler](../examples/addons/cluster-autoscaler/README.md) | false               | 1                   | Delivers the Kubernetes cluster autoscaler component. See https://github.com/kubernetes/autoscaler/tree/master/cluster-autoscaler/cloudprovider/azure for more info |
| [nvidia-device-plugin](../examples/addons/nvidia-device-plugin/README.md) | true if using a Kubernetes cluster (v1.10+) with an N-series agent pool               | 1                   | Delivers the Kubernetes NVIDIA device plugin component. See https://github.com/NVIDIA/k8s-device-plugin for more info |
| container-monitoring                       | false               | 1                   | Delivers the Kubernetes container monitoring component |
| [keyvault-flexvolume](../examples/addons/keyvault-flexvolume/README.md)                        | false               | as many as linux agent nodes                   | Access secrets, keys, and certs in Azure Key Vault from pods |
| [aad-pod-identity](../examples/addons/aad-pod-identity/README.md)                        | false               | 1 + 1 on each linux agent nodes | Assign Azure Active Directory Identities to Kubernetes applications |

To give a bit more info on the `addons` property: We've tried to expose the basic bits of data that allow useful configuration of these cluster features. Here are some example usage patterns that will unpack what `addons` provide:

To enable an addon (using "tiller" as an example):

```
"kubernetesConfig": {
    "addons": [
        {
            "name": "tiller",
            "enabled" : true
        }
    ]
}
```

As you can see above, `addons` is an array child property of `kubernetesConfig`. Each addon that you want to add custom configuration to would be represented as an object item in the array. For example, to disable both tiller and dashboard:

```
"kubernetesConfig": {
    "addons": [
        {
            "name": "tiller",
            "enabled" : false
        },
        {
            "name": "kubernetes-dashboard",
            "enabled" : false
        }
    ]
}
```

More usefully, let's add some custom configuration to the above addons:

```
"kubernetesConfig": {
    "addons": [
        {
            "name": "tiller",
            "enabled": true,
            "containers": [
                {
                  "name": "tiller",
                  "image": "myDockerHubUser/tiller:v3.0.0-alpha",
                  "cpuRequests": "1",
                  "memoryRequests": "1024Mi",
                  "cpuLimits": "1",
                  "memoryLimits": "1024Mi"
                }
              ]
        },
        {
            "name": "kubernetes-dashboard",
            "enabled": true,
            "containers": [
                {
                  "name": "kubernetes-dashboard",
                  "cpuRequests": "50m",
                  "memoryRequests": "512Mi",
                  "cpuLimits": "50m",
                  "memoryLimits": "512Mi"
                }
              ]
        },
        {
            "name": "cluster-autoscaler",
            "enabled": true,
            "containers": [
              {
                "name": "cluster-autoscaler",
                "cpuRequests": "100m",
                "memoryRequests": "300Mi",
                "cpuLimits": "100m",
                "memoryLimits": "300Mi"
              }
            ],
            "config": {
              "maxNodes": "5",
              "minNodes": "1"
            }
        }
    ]
}
```

Above you see custom configuration for both tiller and kubernetes-dashboard. Both include specific resource limit values across the following dimensions:

- cpuRequests
- memoryRequests
- cpuLimits
- memoryLimits

See https://kubernetes.io/docs/concepts/configuration/manage-compute-resources-container/ for more on Kubernetes resource limits.

Additionally above, we specified a custom docker image for tiller, let's say we want to build a cluster and test an alpha version of tiller in it.

Finally, the `addons.enabled` boolean property was omitted above; that's by design. If you specify a `containers` configuration, acs-engine assumes you're enabling the addon. The very first example above demonstrates a simple "enable this addon with default configuration" declaration.

<a name="feat-kubelet-config"></a>

#### kubeletConfig

`kubeletConfig` declares runtime configuration for the kubelet running on all master and agent nodes. It is a generic key/value object, and a child property of `kubernetesConfig`. An example custom kubelet config:

```
"kubernetesConfig": {
    "kubeletConfig": {
        "--eviction-hard": "memory.available<250Mi,nodefs.available<20%,nodefs.inodesFree<10%"
    }
}
```

See [here](https://kubernetes.io/docs/reference/generated/kubelet/) for a reference of supported kubelet options.

Below is a list of kubelet options that acs-engine will configure by default:

| kubelet option                      | default value                                                                                                                                                 |
| ----------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| "--cloud-config"                    | "/etc/kubernetes/azure.json"                                                                                                                                  |
| "--cloud-provider"                  | "azure"                                                                                                                                                       |
| "--cluster-domain"                  | "cluster.local"                                                                                                                                               |
| "--pod-infra-container-image"       | "pause-amd64:_version_"                                                                                                                                       |
| "--max-pods"                        | "30", or "110" if using kubenet --network-plugin (i.e., `"networkPlugin": "kubenet"`)                                                                         |
| "--eviction-hard"                   | "memory.available<100Mi,nodefs.available<10%,nodefs.inodesFree<5%"                                                                                            |
| "--node-status-update-frequency"    | "10s"                                                                                                                                                         |
| "--image-gc-high-threshold"         | "85"                                                                                                                                                          |
| "--image-gc-low-threshold"          | "850"                                                                                                                                                         |
| "--non-masquerade-cidr"             | "10.0.0.0/8"                                                                                                                                                  |
| "--azure-container-registry-config" | "/etc/kubernetes/azure.json"                                                                                                                                  |
| "--pod-max-pids"                    | "100" (need to activate the feature in --feature-gates=SupportPodPidsLimit=true)                                                                              |
| "--image-pull-progress-deadline"    | "30m"                                                                                                                                                         |
| "--feature-gates"                   | No default (can be a comma-separated list). On agent nodes `Accelerators=true` will be applied in the `--feature-gates` option for k8s versions before 1.11.0 |

Below is a list of kubelet options that are _not_ currently user-configurable, either because a higher order configuration vector is available that enforces kubelet configuration, or because a static configuration is required to build a functional cluster:

| kubelet option                               | default value                                    |
| -------------------------------------------- | ------------------------------------------------ |
| "--address"                                  | "0.0.0.0"                                        |
| "--allow-privileged"                         | "true"                                           |
| "--pod-manifest-path"                        | "/etc/kubernetes/manifests"                      |
| "--network-plugin"                           | "cni"                                            |
| "--node-labels"                              | (based on Azure node metadata)                   |
| "--cgroups-per-qos"                          | "true"                                           |
| "--enforce-node-allocatable"                 | "pods"                                           |
| "--kubeconfig"                               | "/var/lib/kubelet/kubeconfig"                    |
| "--register-node" (master nodes only)        | "true"                                           |
| "--register-with-taints" (master nodes only) | "node-role.kubernetes.io/master=true:NoSchedule" |
| "--keep-terminated-pod-volumes"              | "false"                                          |

<a name="feat-controller-manager-config"></a>

#### controllerManagerConfig

`controllerManagerConfig` declares runtime configuration for the kube-controller-manager daemon running on all master nodes. Like `kubeletConfig` it is a generic key/value object, and a child property of `kubernetesConfig`. An example custom controller-manager config:

```
"kubernetesConfig": {
    "controllerManagerConfig": {
        "--node-monitor-grace-period": "40s",
        "--pod-eviction-timeout": "5m0s",
        "--route-reconciliation-period": "10s"
        "--terminated-pod-gc-threshold": "5000"
    }
}
```

See [here](https://kubernetes.io/docs/reference/generated/kube-controller-manager/) for a reference of supported controller-manager options.

Below is a list of controller-manager options that acs-engine will configure by default:

| controller-manager option       | default value                              |
| ------------------------------- | ------------------------------------------ |
| "--node-monitor-grace-period"   | "40s"                                      |
| "--pod-eviction-timeout"        | "5m0s"                                     |
| "--route-reconciliation-period" | "10s"                                      |
| "--terminated-pod-gc-threshold" | "5000"                                     |
| "--feature-gates"               | No default (can be a comma-separated list) |

Below is a list of controller-manager options that are _not_ currently user-configurable, either because a higher order configuration vector is available that enforces controller-manager configuration, or because a static configuration is required to build a functional cluster:

| controller-manager option            | default value                                           |
| ------------------------------------ | ------------------------------------------------------- |
| "--kubeconfig"                       | "/var/lib/kubelet/kubeconfig"                           |
| "--allocate-node-cidrs"              | "false"                                                 |
| "--cluster-cidr"                     | "10.240.0.0/12"                                         |
| "--cluster-name"                     | _auto-generated using api model properties_             |
| "--cloud-provider"                   | "azure"                                                 |
| "--cloud-config"                     | "/etc/kubernetes/azure.json"                            |
| "--root-ca-file"                     | "/etc/kubernetes/certs/ca.crt"                          |
| "--cluster-signing-cert-file"        | "/etc/kubernetes/certs/ca.crt"                          |
| "--cluster-signing-key-file"         | "/etc/kubernetes/certs/ca.key"                          |
| "--service-account-private-key-file" | "/etc/kubernetes/certs/apiserver.key"                   |
| "--leader-elect"                     | "true"                                                  |
| "--v"                                | "2"                                                     |
| "--profiling"                        | "false"                                                 |
| "--use-service-account-credentials"  | "false" ("true" if kubernetesConfig.enableRbac is true) |

<a name="feat-cloud-controller-manager-config"></a>

#### cloudControllerManagerConfig

`cloudControllerManagerConfig` declares runtime configuration for the cloud-controller-manager daemon running on all master nodes in a Cloud Controller Manager configuration. Like `kubeletConfig` it is a generic key/value object, and a child property of `kubernetesConfig`. An example custom cloud-controller-manager config:

```
"kubernetesConfig": {
    "cloudControllerManagerConfig": {
        "--route-reconciliation-period": "1m"
    }
}
```

See [here](https://kubernetes.io/docs/reference/generated/cloud-controller-manager/) for a reference of supported controller-manager options.

Below is a list of cloud-controller-manager options that acs-engine will configure by default:

| controller-manager option       | default value |
| ------------------------------- | ------------- |
| "--route-reconciliation-period" | "10s"         |

Below is a list of cloud-controller-manager options that are _not_ currently user-configurable, either because a higher order configuration vector is available that enforces controller-manager configuration, or because a static configuration is required to build a functional cluster:

| controller-manager option | default value                               |
| ------------------------- | ------------------------------------------- |
| "--kubeconfig"            | "/var/lib/kubelet/kubeconfig"               |
| "--allocate-node-cidrs"   | "false"                                     |
| "--cluster-cidr"          | "10.240.0.0/12"                             |
| "--cluster-name"          | _auto-generated using api model properties_ |
| "--cloud-provider"        | "azure"                                     |
| "--cloud-config"          | "/etc/kubernetes/azure.json"                |
| "--leader-elect"          | "true"                                      |
| "--v"                     | "2"                                         |

<a name="feat-apiserver-config"></a>

#### apiServerConfig

`apiServerConfig` declares runtime configuration for the kube-apiserver daemon running on all master nodes. Like `kubeletConfig` and `controllerManagerConfig` it is a generic key/value object, and a child property of `kubernetesConfig`. An example custom apiserver config:

```
"kubernetesConfig": {
    "apiServerConfig": {
        "--request-timeout": "30s"
    }
}
```

Or perhaps you want to customize/override the set of admission-control flags passed to the API Server by default, you can omit the options you don't want and specify only the ones you need as follows:

```
"orchestratorProfile": {
      "orchestratorType": "Kubernetes",
      "orchestratorRelease": "1.8",
      "kubernetesConfig": {
        "apiServerConfig": {
          "--admission-control":  "NamespaceLifecycle,LimitRanger,ServiceAccount,DefaultStorageClass,ResourceQuota,AlwaysPullImages"
        }
      }
    }
```

See [here](https://kubernetes.io/docs/reference/generated/kube-apiserver/) for a reference of supported apiserver options.

Below is a list of apiserver options that acs-engine will configure by default:

| apiserver option                | default value                                                                                                                                                                                                                           |
| ------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| "--admission-control"           | "NamespaceLifecycle,LimitRanger,ServiceAccount,DefaultStorageClass,ResourceQuota,AlwaysPullImages" (Kubernetes versions prior to 1.9.0                                                                               |
| "--enable-admission-plugins"`*` | "NamespaceLifecycle,LimitRanger,ServiceAccount,DefaultStorageClass,DefaultTolerationSeconds,MutatingAdmissionWebhook,ValidatingAdmissionWebhook,ResourceQuota,AlwaysPullImages" (Kubernetes versions 1.9.0 and later |
| "--authorization-mode"          | "Node", "RBAC" (_the latter if enabledRbac is true_)                                                                                                                                                                                    |
| "--audit-log-maxage"            | "30"                                                                                                                                                                                                                                    |
| "--audit-log-maxbackup"         | "10"                                                                                                                                                                                                                                    |
| "--audit-log-maxsize"           | "100"                                                                                                                                                                                                                                   |
| "--feature-gates"               | No default (can be a comma-separated list)                                                                                                                                                                                              |
| "--oidc-username-claim"         | "oid" (_if has AADProfile_)                                                                                                                                                                                                             |
| "--oidc-groups-claim"           | "groups" (_if has AADProfile_)                                                                                                                                                                                                          |
| "--oidc-client-id"              | _calculated value that represents OID client ID_ (_if has AADProfile_)                                                                                                                                                                  |
| "--oidc-issuer-url"             | _calculated value that represents OID issuer URL_ (_if has AADProfile_)                                                                                                                                                                 |

`*` In Kubernetes versions 1.10.0 and later the `--admission-control` flag is deprecated and `--enable-admission-plugins` is used in its stead.

Below is a list of apiserver options that are _not_ currently user-configurable, either because a higher order configuration vector is available that enforces apiserver configuration, or because a static configuration is required to build a functional cluster:

| apiserver option                            | default value                                                                           |
| ------------------------------------------- | --------------------------------------------------------------------------------------- |
| "--bind-address"                            | "0.0.0.0"                                                                               |
| "--advertise-address"                       | _calculated value that represents listening URI for API server_                         |
| "--allow-privileged"                        | "true"                                                                                  |
| "--anonymous-auth"                          | "false                                                                                  |
| "--audit-log-path"                          | "/var/log/apiserver/audit.log"                                                          |
| "--insecure-port"                           | "8080"                                                                                  |
| "--secure-port"                             | "443"                                                                                   |
| "--service-account-lookup"                  | "true"                                                                                  |
| "--etcd-cafile"                             | "/etc/kubernetes/certs/ca.crt"                                                          |
| "--etcd-certfile"                           | "/etc/kubernetes/certs/etcdclient.crt"                                                  |
| "--etcd-keyfile"                            | "/etc/kubernetes/certs/etcdclient.key"                                                  |
| "--etcd-servers"                            | _calculated value that represents etcd servers_                                         |
| "--profiling"                               | "false"                                                                                 |
| "--repair-malformed-updates"                | "false"                                                                                 |
| "--tls-cert-file"                           | "/etc/kubernetes/certs/apiserver.crt"                                                   |
| "--tls-private-key-file"                    | "/etc/kubernetes/certs/apiserver.key"                                                   |
| "--client-ca-file"                          | "/etc/kubernetes/certs/ca.crt"                                                          |
| "--service-account-key-file"                | "/etc/kubernetes/certs/apiserver.key"                                                   |
| "--kubelet-client-certificate"              | "/etc/kubernetes/certs/client.crt"                                                      |
| "--kubelet-client-key"                      | "/etc/kubernetes/certs/client.key"                                                      |
| "--service-cluster-ip-range"                | _see serviceCIDR_                                                                       |
| "--storage-backend"                         | _calculated value that represents etcd version_                                         |
| "--v"                                       | "4"                                                                                     |
| "--experimental-encryption-provider-config" | "/etc/kubernetes/encryption-config.yaml" (_if enableDataEncryptionAtRest is true_)      |
| "--experimental-encryption-provider-config" | "/etc/kubernetes/encryption-config.yaml" (_if enableEncryptionWithExternalKms is true_) |
| "--requestheader-client-ca-file"            | "/etc/kubernetes/certs/proxy-ca.crt" (_if enableAggregatedAPIs is true_)                |
| "--proxy-client-cert-file"                  | "/etc/kubernetes/certs/proxy.crt" (_if enableAggregatedAPIs is true_)                   |
| "--proxy-client-key-file"                   | "/etc/kubernetes/certs/proxy.key" (_if enableAggregatedAPIs is true_)                   |
| "--requestheader-allowed-names"             | "" (_if enableAggregatedAPIs is true_)                                                  |
| "--requestheader-extra-headers-prefix"      | "X-Remote-Extra-" (_if enableAggregatedAPIs is true_)                                   |
| "--requestheader-group-headers"             | "X-Remote-Group" (_if enableAggregatedAPIs is true_)                                    |
| "--requestheader-username-headers"          | "X-Remote-User" (_if enableAggregatedAPIs is true_)                                     |
| "--cloud-provider"                          | "azure" (_unless useCloudControllerManager is true_)                                    |
| "--cloud-config"                            | "/etc/kubernetes/azure.json" (_unless useCloudControllerManager is true_)               |

<a name="feat-scheduler-config"></a>

#### schedulerConfig

`schedulerConfig` declares runtime configuration for the kube-scheduler daemon running on all master nodes. Like `kubeletConfig`, `controllerManagerConfig`, and `apiServerConfig` it is a generic key/value object, and a child property of `kubernetesConfig`. An example custom apiserver config:

```
"kubernetesConfig": {
    "schedulerConfig": {
        "--v": "2"
    }
}
```

See [here](https://kubernetes.io/docs/reference/generated/kube-scheduler/) for a reference of supported kube-scheduler options.

Below is a list of scheduler options that acs-engine will configure by default:

| kube-scheduler option | default value                              |
| --------------------- | ------------------------------------------ |
| "--v"                 | "2"                                        |
| "--feature-gates"     | No default (can be a comma-separated list) |

Below is a list of kube-scheduler options that are _not_ currently user-configurable, either because a higher order configuration vector is available that enforces kube-scheduler configuration, or because a static configuration is required to build a functional cluster:

| kube-scheduler option | default value                 |
| --------------------- | ----------------------------- |
| "--kubeconfig"        | "/var/lib/kubelet/kubeconfig" |
| "--leader-elect"      | "true"                        |
| "--profiling"         | "false"                       |

We consider `kubeletConfig`, `controllerManagerConfig`, `apiServerConfig`, and `schedulerConfig` to be generic conveniences that add power/flexibility to cluster deployments. Their usage comes with no operational guarantees! They are manual tuning features that enable low-level configuration of a kubernetes cluster.

<a name="feat-private-cluster"></a>

#### privateCluster

`privateCluster` defines a cluster without public addresses assigned. It is a child property of `kubernetesConfig`.

| Name           | Required | Description                                                                                                                                          |
| -------------- | -------- | ---------------------------------------------------------------------------------------------------------------------------------------------------- |
| enabled        | no       | Enable [Private Cluster](./kubernetes/features.md/#feat-private-cluster) (boolean - default == false)                                                |
| jumpboxProfile | no       | Configure and auto-provision a jumpbox to access your private cluster. `jumpboxProfile` is ignored if enabled is `false`. See `jumpboxProfile` below |

#### jumpboxProfile

`jumpboxProfile` describes the settings for a jumpbox deployed via acs-engine to access a private cluster. It is a child property of `privateCluster`.

| Name           | Required | Description                                                                                                                                                                        |
| -------------- | -------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| name           | yes      | This is the unique name for the jumpbox VM. Some resources deployed with the jumpbox are derived from this name                                                                    |
| vmSize         | yes      | Describes a valid [Azure VM Sizes](https://azure.microsoft.com/en-us/documentation/articles/virtual-machines-windows-sizes/)                                                       |
| publicKey      | yes      | The public SSH key used for authenticating access to the jumpbox. Here are instructions for [generating a public/private key pair](ssh.md#ssh-key-generation)                      |
| osDiskSizeGB   | no       | Describes the OS Disk Size in GB. Defaults to `30`                                                                                                                                 |
| storageProfile | no       | Specifies the storage profile to use. Valid values are [ManagedDisks](../examples/disks-managed) or [StorageAccount](../examples/disks-storageaccount). Defaults to `ManagedDisks` |
| username       | no       | Describes the admin username to be used on the jumpbox. Defaults to `azureuser`                                                                                                    |

### masterProfile

`masterProfile` describes the settings for master configuration.

| Name                         | Required                                  | Description                                                                                                                                                                                                                                                                                                                                                                                                                |
| ---------------------------- | ----------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| count                        | yes                                       | Masters have count value of 1, 3, or 5 masters                                                                                                                                                                                                                                                                                                                                                                             |
| dnsPrefix                    | yes                                       | The dns prefix for the master FQDN. The master FQDN is used for SSH or commandline access. This must be a unique name. ([bring your own VNET examples](../examples/vnet))                                                                                                                                                                                                                                                  |
| subjectAltNames              | no                                        | An array of fully qualified domain names using which a user can reach API server. These domains are added as Subject Alternative Names to the generated API server certificate. **NOTE**: These domains **will not** be automatically provisioned.                                                                                                                                                                         |
| firstConsecutiveStaticIP     | only required when vnetSubnetId specified | The IP address of the first master. IP Addresses will be assigned consecutively to additional master nodes                                                                                                                                                                                                                                                                                                                 |
| vmsize                       | yes                                       | Describes a valid [Azure VM Sizes](https://azure.microsoft.com/en-us/documentation/articles/virtual-machines-windows-sizes/). These are restricted to machines with at least 2 cores and 100GB of ephemeral disk space                                                                                                                                                                                                     |
| osDiskSizeGB                 | no                                        | Describes the OS Disk Size in GB                                                                                                                                                                                                                                                                                                                                                                                           |
| vnetSubnetId                 | no                                        | Specifies the Id of an alternate VNET subnet. The subnet id must specify a valid VNET ID owned by the same subscription. ([bring your own VNET examples](../examples/vnet))                                                                                                                                                                                                                                                |
| extensions                   | no                                        | This is an array of extensions. This indicates that the extension be run on a single master. The name in the extensions array must exactly match the extension name in the extensionProfiles                                                                                                                                                                                                                               |
| vnetCidr                     | no                                        | Specifies the VNET cidr when using a custom VNET ([bring your own VNET examples](../examples/vnet))                                                                                                                                                                                                                                                                                                                        |
| imageReference.name          | no                                        | The name of the Linux OS image. Needs to be used in conjunction with resourceGroup, below                                                                                                                                                                                                                                                                                                                                  |
| imageReference.resourceGroup | no                                        | Resource group that contains the Linux OS image. Needs to be used in conjunction with name, above                                                                                                                                                                                                                                                                                                                          |
| distro                       | no                                        | Select Master(s) Operating System (Linux only). Currently supported values are: `ubuntu` and `coreos` (CoreOS support is currently experimental). Defaults to `ubuntu` if undefined. Currently supported OS and orchestrator configurations -- `ubuntu`: DCOS, Docker Swarm, Kubernetes; `RHEL`: OpenShift; `coreos`: Kubernetes. [Example of CoreOS Master with CoreOS Agents](../examples/coreos/kubernetes-coreos.json) |
| customFiles                  | no                                        | The custom files to be provisioned to the master nodes. Defined as an array of json objects with each defined as `"source":"absolute-local-path", "dest":"absolute-path-on-masternodes"`.[See examples](../examples/customfiles)                                                                                                                                                                                           |

### agentPoolProfiles

A cluster can have 0 to 12 agent pool profiles. Agent Pool Profiles are used for creating agents with different capabilities such as VMSizes, VMSS or Availability Set, Public/Private access, user-defined OS Images, [attached storage disks](../examples/disks-storageaccount), [attached managed disks](../examples/disks-managed), or [Windows](../examples/windows).

| Name                         | Required                                                             | Description                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                      |
| ---------------------------- | -------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| availabilityProfile          | no                                                                   | Supported values are `VirtualMachineScaleSets` (default, except for Kubernetes clusters before version 1.10) and `AvailabilitySet`.                                                                                                                                                                                                                                                                                                                                                                                              |
| count                        | yes                                                                  | Describes the node count                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                         |
| scaleSetPriority             | no                                                                   | Supported values are `Regular` (default) and `Low`. Only applies to clusters with availabilityProfile `VirtualMachineScaleSets`. Enables the usage of [Low-priority VMs on Scale Sets](https://docs.microsoft.com/en-us/azure/virtual-machine-scale-sets/virtual-machine-scale-sets-use-low-priority).                                                                                                                                                                                                                           |
| scaleSetEvictionPolicy       | no                                                                   | Supported values are `Delete` (default) and `Deallocate`. Only applies to clusters with availabilityProfile of `VirtualMachineScaleSets` and scaleSetPriority of `Low`.                                                                                                                                                                                                                                                                                                                                                          |
| diskSizesGB                  | no                                                                   | Describes an array of up to 4 attached disk sizes. Valid disk size values are between 1 and 1024                                                                                                                                                                                                                                                                                                                                                                                                                                 |
| dnsPrefix                    | Required if agents are to be exposed publically with a load balancer | The dns prefix that forms the FQDN to access the loadbalancer for this agent pool. This must be a unique name among all agent pools. Not supported for Kubernetes clusters                                                                                                                                                                                                                                                                                                                                                       |
| name                         | yes                                                                  | This is the unique name for the agent pool profile. The resources of the agent pool profile are derived from this name                                                                                                                                                                                                                                                                                                                                                                                                           |
| ports                        | only required if needed for exposing services publically             | Describes an array of ports need for exposing publically. A tcp probe is configured for each port and only opens to an agent node if the agent node is listening on that port. A maximum of 150 ports may be specified. Not supported for Kubernetes clusters                                                                                                                                                                                                                                                                    |
| storageProfile               | no                                                                   | Specifies the storage profile to use. Valid values are [ManagedDisks](../examples/disks-managed) or [StorageAccount](../examples/disks-storageaccount). Defaults to `ManagedDisks`                                                                                                                                                                                                                                                                                                                                               |
| vmsize                       | yes                                                                  | Describes a valid [Azure VM Sizes](https://azure.microsoft.com/en-us/documentation/articles/virtual-machines-windows-sizes/). These are restricted to machines with at least 2 cores                                                                                                                                                                                                                                                                                                                                             |
| osDiskSizeGB                 | no                                                                   | Describes the OS Disk Size in GB                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                 |
| vnetSubnetId                 | no                                                                   | Specifies the Id of an alternate VNET subnet. The subnet id must specify a valid VNET ID owned by the same subscription. ([bring your own VNET examples](../examples/vnet))                                                                                                                                                                                                                                                                                                                                                      |
| imageReference.name          | no                                                                   | The name of a a Linux OS image. Needs to be used in conjunction with resourceGroup, below                                                                                                                                                                                                                                                                                                                                                                                                                                        |
| imageReference.resourceGroup | no                                                                   | Resource group that contains the Linux OS image. Needs to be used in conjunction with name, above                                                                                                                                                                                                                                                                                                                                                                                                                                |
| osType                       | no                                                                   | Specifies the agent pool's Operating System. Supported values are `Windows` and `Linux`. Defaults to `Linux`                                                                                                                                                                                                                                                                                                                                                                                                                     |
| distro                       | no                                                                   | Specifies the agent pool's Linux distribution. Supported values are `ubuntu` and `coreos` (CoreOS support is currently experimental). Defaults to `ubuntu` if undefined, unless `osType` is defined as `Windows` (in which case `distro` is unused). Currently supported OS and orchestrator configurations -- `ubuntu`: DCOS, Docker Swarm, Kubernetes; `RHEL`: OpenShift; `coreos`: Kubernetes. [Example of CoreOS Master with Windows and Linux (CoreOS and Ubuntu) Agents](../examples/coreos/kubernetes-coreos-hybrid.json) |
| acceleratedNetworkingEnabled | no                                                                   | Use [Azure Accelerated Networking](https://azure.microsoft.com/en-us/blog/maximize-your-vm-s-performance-with-accelerated-networking-now-generally-available-for-both-windows-and-linux/) feature for agents (You must select a VM SKU that support Accelerated Networking)                                                                                                                                                                                                                                                      |

### linuxProfile

`linuxProfile` provides the linux configuration for each linux node in the cluster

| Name                             | Required | Description                                                                                                                                                                      |
| -------------------------------- | -------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| adminUsername                    | yes      | Describes the username to be used on all linux clusters                                                                                                                          |
| ssh.publicKeys.keyData           | yes      | The public SSH key used for authenticating access to all Linux nodes in the cluster. Here are instructions for [generating a public/private key pair](ssh.md#ssh-key-generation) |
| secrets                          | no       | Specifies an array of key vaults to pull secrets from and what secrets to pull from each                                                                                         |
| customSearchDomain.name          | no       | describes the search domain to be used on all linux clusters                                                                                                                     |
| customSearchDomain.realmUser     | no       | describes the realm user with permissions to update dns registries on Windows Server DNS                                                                                         |
| customSearchDomain.realmPassword | no       | describes the realm user password to update dns registries on Windows Server DNS                                                                                                 |
| customNodesDNS.dnsServer         | no       | describes the IP address of the DNS Server                                                                                                                                       |

#### secrets

`secrets` details which certificates to install on the masters and nodes in the cluster.

A cluster can have a list of key vaults to install certs from.

On linux boxes the certs are saved on under the directory "/var/lib/waagent/". 2 files are saved per certificate:

1.  `{thumbprint}.crt` : this is the full cert chain saved in PEM format
2.  `{thumbprint}.prv` : this is the private key saved in PEM format

| Name                             | Required | Description                                                         |
| -------------------------------- | -------- | ------------------------------------------------------------------- |
| sourceVault.id                   | yes      | The azure resource manager id of the key vault to pull secrets from |
| vaultCertificates.certificateUrl | yes      | Keyvault URL to this cert including the version                     |

format for `sourceVault.id`, can be obtained in cli, or found in the portal: /subscriptions/{subscription-id}/resourceGroups/{resource-group}/providers/Microsoft.KeyVault/vaults/{keyvaultname}

format for `vaultCertificates.certificateUrl`, can be obtained in cli, or found in the portal:
https://{keyvaultname}.vault.azure.net:443/secrets/{secretName}/{version}

### servicePrincipalProfile

`servicePrincipalProfile` describes an Azure Service credentials to be used by the cluster for self-configuration. See [service principal](serviceprincipal.md) for more details on creation.

| Name                         | Required                          | Description                                                                                                 |
| ---------------------------- | --------------------------------- | ----------------------------------------------------------------------------------------------------------- |
| clientId                     | yes, for Kubernetes clusters      | describes the Azure client id. It is recommended to use a separate client ID per cluster                    |
| secret                       | yes, for Kubernetes clusters      | describes the Azure client secret. It is recommended to use a separate client secret per client id          |
| objectId                     | optional, for Kubernetes clusters | describes the Azure service principal object id. It is required if enableEncryptionWithExternalKms is true  |
| keyvaultSecretRef.vaultId    | no, for Kubernetes clusters       | describes the vault id of the keyvault to retrieve the service principal secret from. See below for format. |
| keyvaultSecretRef.secretName | no, for Kubernetes clusters       | describes the name of the service principal secret in keyvault                                              |
| keyvaultSecretRef.version    | no, for Kubernetes clusters       | describes the version of the secret to use                                                                  |

format for `keyvaultSecretRef.vaultId`, can be obtained in cli, or found in the portal:
`/subscriptions/<SUB_ID>/resourceGroups/<RG_NAME>/providers/Microsoft.KeyVault/vaults/<KV_NAME>`. See [keyvault params](../examples/keyvault-params/README.md#service-principal-profile) for an example.

## Cluster Defintions for apiVersion "2016-03-30"

Here are the cluster definitions for apiVersion "2016-03-30". This matches the api version of the Azure Container Service Engine.

### apiVersion

| Name       | Required | Description                                                             |
| ---------- | -------- | ----------------------------------------------------------------------- |
| apiVersion | yes      | The version of the template. For "2016-03-30" the value is "2016-03-30" |

### orchestratorProfile

`orchestratorProfile` describes the orchestrator settings.

| Name             | Required | Description                                     |
| ---------------- | -------- | ----------------------------------------------- |
| orchestratorType | yes      | Specifies the orchestrator type for the cluster |

Here are the valid values for the orchestrator types:

1.  `DCOS` - this represents the [DC/OS orchestrator](dcos.md).
2.  `Swarm` - this represents the [Swarm orchestrator](swarm.md).
3.  `Kubernetes` - this represents the [Kubernetes orchestrator](kubernetes.md).
4.  `Swarm Mode` - this represents the [Swarm Mode orchestrator](swarmmode.md).
5.  `OpenShift` - this represents the [OpenShift orchestrator](openshift.md)

### masterProfile

`masterProfile` describes the settings for master configuration.

| Name      | Required | Description                                                                                                                                                                |
| --------- | -------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| count     | yes      | Masters have count value of 1, 3, or 5 masters                                                                                                                             |
| dnsPrefix | yes      | The dns prefix for the masters FQDN. The master FQDN is used for SSH or commandline access. This must be a unique name. ([bring your own VNET examples](../examples/vnet)) |

### agentPoolProfiles

For apiVersion "2016-03-30", a cluster may have only 1 agent pool profiles.

| Name      | Required                                                             | Description                                                                                                                                                                          |
| --------- | -------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| count     | yes                                                                  | Describes the node count                                                                                                                                                             |
| dnsPrefix | required if agents are to be exposed publically with a load balancer | The dns prefix that forms the FQDN to access the loadbalancer for this agent pool. This must be a unique name among all agent pools                                                  |
| name      | yes                                                                  | The unique name for the agent pool profile. The resources of the agent pool profile are derived from this name                                                                       |
| vmsize    | yes                                                                  | Describes a valid [Azure VM Sizes](https://azure.microsoft.com/en-us/documentation/articles/virtual-machines-windows-sizes/). These are restricted to machines with at least 2 cores |

### linuxProfile

`linuxProfile` provides the linux configuration for each linux node in the cluster

| Name                      | Required | Description                                                                                                                                                                      |
| ------------------------- | -------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| adminUsername             | yes      | Describes the username to be used on all linux clusters                                                                                                                          |
| ssh.publicKeys[0].keyData | yes      | The public SSH key used for authenticating access to all Linux nodes in the cluster. Here are instructions for [generating a public/private key pair](ssh.md#ssh-key-generation) |

### aadProfile

`aadProfile` provides [Azure Active Directory integration](kubernetes.aad.md) configuration for the cluster, currently only available for Kubernetes orchestrator.

| Name         | Required | Description                                                                                                                 |
| ------------ | -------- | --------------------------------------------------------------------------------------------------------------------------- |
| clientAppID  | yes      | Describes the client AAD application ID                                                                                     |
| serverAppID  | yes      | Describes the server AAD application ID                                                                                     |
| adminGroupID | no       | Describes the AAD Group Object ID that will be assigned the cluster-admin RBAC role                                         |
| tenantID     | no       | Describes the AAD tenant ID to use for authentication. If not specified, will use the tenant of the deployment subscription |

### extensionProfiles

A cluster can have 0 - N extensions in extension profiles. Extension profiles allow a user to easily add pre-packaged functionality into a cluster. An example would be configuring a monitoring solution on your cluster. You can think of extensions like a marketplace for acs clusters.

| Name                | Required | Description                                                                                                                                                                      |
| ------------------- | -------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| name                | yes      | The name of the extension. This has to exactly match the name of a folder under the extensions folder                                                                            |
| version             | yes      | The version of the extension. This has to exactly match the name of the folder under the extension name folder                                                                   |
| extensionParameters | optional | Extension parameters may be required by extensions. The format of the parameters is also extension dependant                                                                     |
| rootURL             | optional | URL to the root location of extensions. The rootURL must have an extensions child folder that follows the extensions convention. The rootURL is mainly used for testing purposes |

You can find more information, as well as a list of extensions on the [extensions documentation](extensions.md).

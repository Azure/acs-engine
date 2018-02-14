# Microsoft Azure Container Service Engine - Cluster Definition

## Cluster Defintions for apiVersion "vlabs"

Here are the cluster definitions for apiVersion "vlabs"

### apiVersion

|Name|Required|Description|
|---|---|---|
|apiVersion|yes|The version of the template.  For "vlabs" the value is "vlabs".|

### orchestratorProfile
`orchestratorProfile` describes the orchestrator settings.

|Name|Required|Description|
|---|---|---|
|orchestratorType|yes|This specifies the orchestrator type for the cluster.|

Here are the valid values for the orchestrator types:

1. `DCOS` - this represents the [DC/OS orchestrator](dcos.md).  [Older releases of DCOS 1.8 may be specified](../examples/dcos-releases).
2. `Kubernetes` - this represents the [Kubernetes orchestrator](kubernetes.md).
3. `Swarm` - this represents the [Swarm orchestrator](swarm.md).
4. `Swarm Mode` - this represents the [Swarm Mode orchestrator](swarmmode.md).

### kubernetesConfig

`kubernetesConfig` describes Kubernetes specific configuration.

|Name|Required|Description|
|---|---|---|
|kubernetesImageBase|no|This specifies the base URL (everything preceding the actual image filename) of the kubernetes hyperkube image to use for cluster deployment, e.g., `k8s-gcrio.azureedge.net/`.|
|dockerEngineVersion|no|Which version of docker-engine to use in your cluster, e.g.. "17.03.*"|
|networkPolicy|no|Specifies the network policy tool for the cluster. Valid values are:<br>`"azure"` (default), which provides an Azure native networking experience,<br>`none` for not enforcing any network policy,<br>`calico` for Calico network policy (clusters with Linux agents only).<br>See [network policy examples](../examples/networkpolicy) for more information.|
|containerRuntime|no|The container runtime to use as a backend. The default is `docker`. The only other option is `clear-containers`.|
|clusterSubnet|no|The IP subnet used for allocating IP addresses for pod network interfaces. The subnet must be in the VNET address space. Default value is 10.244.0.0/16.|
|dnsServiceIP|no|IP address for kube-dns to listen on. If specified must be in the range of `serviceCidr`.|
|dockerBridgeSubnet|no|The specific IP and subnet used for allocating IP addresses for the docker bridge network created on the kubernetes master and agents. Default value is 172.17.0.1/16. This value is used to configure the docker daemon using the [--bip flag](https://docs.docker.com/engine/userguide/networking/default_network/custom-docker0).|
|serviceCidr|no|IP range for Service IPs, Default is "10.0.0.0/16". This range is never routed outside of a node so does not need to lie within clusterSubnet or the VNet.|
|enableRbac|no|Enable [Kubernetes RBAC](https://kubernetes.io/docs/admin/authorization/rbac/) (boolean - default == true) |
|enableAggregatedAPIs|no|Enable [Kubernetes Aggregated APIs](https://kubernetes.io/docs/concepts/api-extension/apiserver-aggregation/).This is required by [Service Catalog](https://github.com/kubernetes-incubator/service-catalog/blob/master/README.md). (boolean - default == false) |
|enableDataEncryptionAtRest|no|Enable [kuberetes data encryption at rest](https://kubernetes.io/docs/tasks/administer-cluster/encrypt-data/).This is currently an alpha feature. (boolean - default == false) |
|maxPods|no|The maximum number of pods per node. The minimum valid value, necessary for running kube-system pods, is 5. Default value is 30 when networkPolicy equals azure, 110 otherwise.|
|gcHighThreshold|no|Sets the --image-gc-high-threshold value on the kublet configuration. Default is 85. [See kubelet Garbage Collection](https://kubernetes.io/docs/concepts/cluster-administration/kubelet-garbage-collection/) |
|gcLowThreshold|no|Sets the --image-gc-low-threshold value on the kublet configuration. Default is 80. [See kubelet Garbage Collection](https://kubernetes.io/docs/concepts/cluster-administration/kubelet-garbage-collection/) |
|useInstanceMetadata|no|Use the Azure cloudprovider instance metadata service for appropriate resource discovery operations. Default is `true`.|
|addons|no|Configure various Kubernetes addons configuration (currently supported: tiller, kubernetes-dashboard). See `addons` configuration below.|
|kubeletConfig|no|Configure various runtime configuration for kubelet. See `kubeletConfig` [below](#feat-kubelet-config).|
|controllerManagerConfig|no|Configure various runtime configuration for controller-manager. See `controllerManagerConfig` [below](#feat-controller-manager-config).|
|cloudControllerManagerConfig|no|Configure various runtime configuration for cloud-controller-manager. See `cloudControllerManagerConfig` [below](#feat-cloud-controller-manager-config).|
|apiServerConfig|no|Configure various runtime configuration for apiserver. See `apiServerConfig` [below](#feat-apiserver-config).|

#### addons

`addons` describes various addons configuration. It is a child property of `kubernetesConfig`. Below is a list of currently available addons:

|Name of addon|Enabled by default?|How many containers|Description|
|---|---|---|---|
|tiller|true|1|Delivers the Helm server-side component: tiller. See https://github.com/kubernetes/helm for more info.|
|kubernetes-dashboard|true|1|Delivers the kubernetes dashboard component. See https://github.com/kubernetes/dashboard for more info.|
|rescheduler|false|1|Delivers the kubernetes rescheduler component.|

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

More usefully, let's add some custom configuration to both of the above addons:

```
"kubernetesConfig": {
    "addons": [
        {
            "name": "tiller",
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
            "containers": [
                {
                  "name": "kubernetes-dashboard",
                  "cpuRequests": "50m",
                  "memoryRequests": "512Mi",
                  "cpuLimits": "50m",
                  "memoryLimits": "512Mi"
                }
              ]
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

|kubelet option|default value|
|---|---|
|"--cloud-config"|"/etc/kubernetes/azure.json"|
|"--cloud-provider"|"azure"|
|"--cluster-domain"|"cluster.local"|
|"--pod-infra-container-image"|"pause-amd64:<version>"|
|"--max-pods"|"110"|
|"--eviction-hard"|"memory.available<100Mi,nodefs.available<10%,nodefs.inodesFree<5%"|
|"--node-status-update-frequency"|"10s"|
|"--image-gc-high-threshold"|"85"|
|"--image-gc-low-threshold"|"850"|
|"--non-masquerade-cidr"|"10.0.0.0/8"|
|"--feature-gates"|No default (can be a comma-separated list). On agent nodes `Accelerators=true` will be applied in the `--feature-gates` option.|

Below is a list of kubelet options that are *not* currently user-configurable, either because a higher order configuration vector is available that enforces kubelet configuration, or because a static configuration is required to build a functional cluster:

|kubelet option|default value|
|---|---|
|"--address"|"0.0.0.0"|
|"--azure-container-registry-config"|"/etc/kubernetes/azure.json"|
|"--allow-privileged"|"true"|
|"--pod-manifest-path"|"/etc/kubernetes/manifests"|
|"--network-plugin"|"cni"|
|"--node-labels"|(based on Azure node metadata)|
|"--cgroups-per-qos"|"false"|
|"--enforce-node-allocatable"|""|
|"--kubeconfig"|"/var/lib/kubelet/kubeconfig"|
|"--register-node" (master nodes only)|"true"|
|"--register-with-taints" (master nodes only)|"node-role.kubernetes.io/master=true:NoSchedule"|
|"--read-only-port"|"0"|
|"--keep-terminated-pod-volumes"|"false"|

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

|controller-manager option|default value|
|---|---|
|"--node-monitor-grace-period"|"40s"|
|"--pod-eviction-timeout"|"5m0s"|
|"--route-reconciliation-period"|"10s"|
|"--terminated-pod-gc-threshold"|"5000"|
|"--feature-gates"|No default (can be a comma-separated list)|


Below is a list of controller-manager options that are *not* currently user-configurable, either because a higher order configuration vector is available that enforces controller-manager configuration, or because a static configuration is required to build a functional cluster:

|controller-manager option|default value|
|---|---|
|"--kubeconfig"|"/var/lib/kubelet/kubeconfig"|
|"--allocate-node-cidrs"|"false"|
|"--cluster-cidr"|"10.240.0.0/12"|
|"--cluster-name"|<auto-generated using api model properties>|
|"--cloud-provider"|"azure"|
|"--cloud-config"|"/etc/kubernetes/azure.json"|
|"--root-ca-file"|"/etc/kubernetes/certs/ca.crt"|
|"--cluster-signing-cert-file"|"/etc/kubernetes/certs/ca.crt"|
|"--cluster-signing-key-file"|"/etc/kubernetes/certs/ca.key"|
|"--service-account-private-key-file"|"/etc/kubernetes/certs/apiserver.key"|
|"--leader-elect"|"true"|
|"--v"|"2"|
|"--profiling"|"false"|
|"--use-service-account-credentials"|"false" ("true" if kubernetesConfig.enableRbac is true)|

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

|controller-manager option|default value|
|---|---|
|"--route-reconciliation-period"|"10s"|


Below is a list of cloud-controller-manager options that are *not* currently user-configurable, either because a higher order configuration vector is available that enforces controller-manager configuration, or because a static configuration is required to build a functional cluster:

|controller-manager option|default value|
|---|---|
|"--kubeconfig"|"/var/lib/kubelet/kubeconfig"|
|"--allocate-node-cidrs"|"false"|
|"--cluster-cidr"|"10.240.0.0/12"|
|"--cluster-name"|<auto-generated using api model properties>|
|"--cloud-provider"|"azure"|
|"--cloud-config"|"/etc/kubernetes/azure.json"|
|"--leader-elect"|"true"|
|"--v"|"2"|

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

|apiserver option|default value|
|---|---|
|"--admission-control"|"NamespaceLifecycle, LimitRanger, ServiceAccount, DefaultStorageClass, ResourceQuota, DenyEscalatingExec, AlwaysPullImages, SecurityContextDeny"|
|"--authorization-mode"|"Node", "RBAC" (*the latter if enabledRbac is true*)|
|"--audit-log-maxage"|"30"|
|"--audit-log-maxbackup"|"10"|
|"--audit-log-maxsize"|"100"|
|"--feature-gates"|No default (can be a comma-separated list)|


Below is a list of apiserver options that are *not* currently user-configurable, either because a higher order configuration vector is available that enforces kubelet configuration, or because a static configuration is required to build a functional cluster:

|apiserver option|default value|
|---|---|
|"--bind-address"|"0.0.0.0"|
|"--advertise-address"|*calculated value that represents listening URI for API server*|
|"--allow-privileged"|"true"|
|"--anonymous-auth"|"false|
|"--audit-log-path"|"/var/log/apiserver/audit.log"|
|"--insecure-port"|"8080"|
|"--secure-port"|"443"|
|"--service-account-lookup"|"true"|
|"--etcd-cafile"|"/etc/kubernetes/certs/ca.crt"|
|"--etcd-certfile"|"/etc/kubernetes/certs/etcdclient.crt"|
|"--etcd-keyfile"|"/etc/kubernetes/certs/etcdclient.key"|
|"--etcd-servers"|*calculated value that represents etcd servers*|
|"--etcd-quorum-read"|"true"|
|"--profiling"|"false"|
|"--repair-malformed-updates"|"false"|
|"--tls-cert-file"|"/etc/kubernetes/certs/apiserver.crt"|
|"--tls-private-key-file"|"/etc/kubernetes/certs/apiserver.key"|
|"--client-ca-file"|"/etc/kubernetes/certs/ca.crt"|
|"--service-account-key-file"|"/etc/kubernetes/certs/apiserver.key"|
|"--kubelet-client-certificate"|"/etc/kubernetes/certs/client.crt"|
|"--kubelet-client-key"|"/etc/kubernetes/certs/client.key"|
|"--service-cluster-ip-range"|*see serviceCIDR*|
|"--storage-backend"|*calculated value that represents etcd version*|
|"--v"|"4"|
|"--experimental-encryption-provider-config"|"/etc/kubernetes/encryption-config.yaml" (*if enableDataEncryptionAtRest is true*)|
|"--requestheader-client-ca-file"|"/etc/kubernetes/certs/proxy-ca.crt" (*if enableAggregatedAPIs is true*)|
|"--proxy-client-cert-file"|"/etc/kubernetes/certs/proxy.crt" (*if enableAggregatedAPIs is true*)|
|"--proxy-client-key-file"|"/etc/kubernetes/certs/proxy.key" (*if enableAggregatedAPIs is true*)|
|"--requestheader-allowed-names"|"" (*if enableAggregatedAPIs is true*)|
|"--requestheader-extra-headers-prefix"|"X-Remote-Extra-" (*if enableAggregatedAPIs is true*)|
|"--requestheader-group-headers"|"X-Remote-Group" (*if enableAggregatedAPIs is true*)|
|"--requestheader-username-headers"|"X-Remote-User" (*if enableAggregatedAPIs is true*)|
|"--cloud-provider"|"azure" (*unless useCloudControllerManager is true*)|
|"--cloud-config"|"/etc/kubernetes/azure.json" (*unless useCloudControllerManager is true*)|
|"--oidc-username-claim"|"oid" (*if has AADProfile*)|
|"--oidc-groups-claim"|"groups" (*if has AADProfile*)|
|"--oidc-client-id"|*calculated value that represents OID client ID* (*if has AADProfile*)|
|"--oidc-issuer-url"|*calculated value that represents OID issuer URL* (*if has AADProfile*)|

We consider `kubeletConfig`, `controllerManagerConfig`, and `apiServerConfig` to be generic conveniences that add power/flexibility to cluster deployments. Their usage comes with no operational guarantees! They are manual tuning features that enable low-level configuration of a kubernetes cluster.

### masterProfile
`masterProfile` describes the settings for master configuration.

|Name|Required|Description|
|---|---|---|
|count|yes|Masters have count value of 1, 3, or 5 masters|
|dnsPrefix|yes|this is the dns prefix for the masters FQDN.  The master FQDN is used for SSH or commandline access. This must be a unique name. ([bring your own VNET examples](../examples/vnet))|
|firstConsecutiveStaticIP|only required when vnetSubnetId specified|this is the IP address of the first master.  IP Addresses will be assigned consecutively to additional master nodes.|
|vmsize|yes|Describes a valid [Azure VM Sizes](https://azure.microsoft.com/en-us/documentation/articles/virtual-machines-windows-sizes/).  These are restricted machines with at least 2 cores and 100GB of ephemeral disk space.|
|osDiskSizeGB|no|Describes the OS Disk Size in GB|
|vnetSubnetId|no|specifies the Id of an alternate VNET subnet.  The subnet id must specify a valid VNET ID owned by the same subscription. ([bring your own VNET examples](../examples/vnet))|
|extensions|no|This is an array of extensions.  This indicates that the extension be run on a single master.  The name in the extensions array must exactly match the extension name in the extensionProfiles.|
|vnetCidr|no| specifies the vnet cidr when using custom Vnets ([bring your own VNET examples](../examples/vnet))|
|distro|no| Select Master(s) Operating System (Linux). Currently supported values are: `ubuntu` and `coreos` (CoreOS support is currently experimental). Defaults to `ubuntu` if undefined. Currently supported OS and orchestrator configurations -- `ubuntu`: DCOS, Docker Swarm, Kubernetes; `coreos`: Kubernetes. [Example of CoreOS Master with CoreOS Agents](../examples/coreos/kubernetes-coreos.json)|

### agentPoolProfiles
A cluster can have 0 to 12 agent pool profiles. Agent Pool Profiles are used for creating agents with different capabilities such as VMSizes, VMSS or Availability Set, Public/Private access, [attached storage disks](../examples/disks-storageaccount), [attached managed disks](../examples/disks-managed), or [Windows](../examples/windows).

|Name|Required|Description|
|---|---|---|
|availabilityProfile|no, defaults to `VirtualMachineScaleSets`| You can choose between `VirtualMachineScaleSets` and `AvailabilitySet`.  As a rule of thumb always choose `VirtualMachineScaleSets` unless you need features such as dynamic attached disks or require Kubernetes|
|count|yes|Describes the node count|
|diskSizesGB|no|describes an array of up to 4 attached disk sizes.  Valid disk size values are between 1 and 1024.|
|dnsPrefix|required if agents are to be exposed publically with a load balancer|this is the dns prefix that forms the FQDN to access the loadbalancer for this agent pool.  This must be a unique name among all agent pools.|
|name|yes|This is the unique name for the agent pool profile. The resources of the agent pool profile are derived from this name.|
|ports|only required if needed for exposing services publically|Describes an array of ports need for exposing publically.  A tcp probe is configured for each port and only opens to an agent node if the agent node is listening on that port.  A maximum of 150 ports may be specified.|
|storageProfile|no, defaults to `StorageAccount`|specifies the storage profile to use.  Valid values are [StorageAccount](../examples/disks-storageaccount) or [ManagedDisks](../examples/disks-managed)|
|vmsize|yes|Describes a valid [Azure VM Sizes](https://azure.microsoft.com/en-us/documentation/articles/virtual-machines-windows-sizes/).  These are restricted to machines with at least 2 cores|
|osDiskSizeGB|no|Describes the OS Disk Size in GB|
|vnetSubnetId|no|specifies the Id of an alternate VNET subnet.  The subnet id must specify a valid VNET ID owned by the same subscription. ([bring your own VNET examples](../examples/vnet))|
|distro|no| Select Agent Pool(s) Operating System (Linux). Currently supported values are: `ubuntu` and `coreos` (CoreOS support is currently experimental). Defaults to `ubuntu` if undefined, unless `osType` is defined as `Windows` (in which case `distro` is unused). Currently supported OS and orchestrator configurations -- `ubuntu`: DCOS, Docker Swarm, Kubernetes; `coreos`: Kubernetes.  [Example of CoreOS Master with Windows and Linux (CoreOS and Ubuntu) Agents](../examples/coreos/kubernetes-coreos-hybrid.json) |

### linuxProfile

`linuxProfile` provides the linux configuration for each linux node in the cluster

|Name|Required|Description|
|---|---|---|
|adminUsername|yes|describes the username to be used on all linux clusters|
|ssh.publicKeys.keyData|yes|The public SSH key used for authenticating access to all Linux nodes in the cluster.  Here are instructions for [generating a public/private key pair](ssh.md#ssh-key-generation).|
|secrets|no|specifies an array of key vaults to pull secrets from and what secrets to pull from each|

#### secrets
`secrets` details which certificates to install on the masters and nodes in the cluster.

A cluster can have a list of key vaults to install certs from.

On linux boxes the certs are saved on under the directory "/var/lib/waagent/". 2 files are saved per certificate:

1. `{thumbprint}.crt` : this is the full cert chain saved in PEM format
2. `{thumbprint}.prv` : this is the private key saved in PEM format

|Name|Required|Description|
|---|---|---|
|sourceVault.id|yes|the azure resource manager id of the key vault to pull secrets from|
|vaultCertificates.certificateUrl|yes|key vault url to this cert including the version|
format for `sourceVault.id`, can be obtained in cli, or found in the portal: /subscriptions/{subscription-id}/resourceGroups/{resource-group}/providers/Microsoft.KeyVault/vaults/{keyvaultname}

format for `vaultCertificates.certificateUrl`, can be obtained in cli, or found in the portal:
https://{keyvaultname}.vault.azure.net:443/secrets/{secretName}/{version}

### servicePrincipalProfile

`servicePrincipalProfile` describes an Azure Service credentials to be used by the cluster for self-configuration.  See [service principal](serviceprincipal.md) for more details on creation.

|Name|Required|Description|
|---|---|---|
|clientId|yes, for Kubernetes clusters|describes the Azure client id.  It is recommended to use a separate client ID per cluster|
|secret|yes, for Kubernetes clusters|describes the Azure client secret.  It is recommended to use a separate client secret per client id|

## Cluster Defintions for apiVersion "2016-03-30"

Here are the cluster definitions for apiVersion "2016-03-30".  This matches the api version of the Azure Container Service Engine.

### apiVersion

|Name|Required|Description|
|---|---|---|
|apiVersion|yes|The version of the template.  For "2016-03-30" the value is "2016-03-30".|

### orchestratorProfile
`orchestratorProfile` describes the orchestrator settings.

|Name|Required|Description|
|---|---|---|
|orchestratorType|yes|This specifies the orchestrator type for the cluster.|

Here are the valid values for the orchestrator types:

1. `DCOS` - this represents the [DC/OS orchestrator](dcos.md).
2. `Swarm` - this represents the [Swarm orchestrator](swarm.md).
3. `Kubernetes` - this represents the [Kubernetes orchestrator](kubernetes.md).
4. `Swarm Mode` - this represents the [Swarm Mode orchestrator](swarmmode.md).

### masterProfile
`masterProfile` describes the settings for master configuration.

|Name|Required|Description|
|---|---|---|
|count|yes|Masters have count value of 1, 3, or 5 masters|
|dnsPrefix|yes|this is the dns prefix for the masters FQDN.  The master FQDN is used for SSH or commandline access. This must be a unique name. ([bring your own VNET examples](../examples/vnet))|

### agentPoolProfiles
For apiVersion "2016-03-30", a cluster may have only 1 agent pool profiles.

|Name|Required|Description|
|---|---|---|
|count|yes|Describes the node count|
|dnsPrefix|required if agents are to be exposed publically with a load balancer|this is the dns prefix that forms the FQDN to access the loadbalancer for this agent pool.  This must be a unique name among all agent pools.|
|name|yes|This is the unique name for the agent pool profile. The resources of the agent pool profile are derived from this name.|
|vmsize|yes|Describes a valid [Azure VM Sizes](https://azure.microsoft.com/en-us/documentation/articles/virtual-machines-windows-sizes/).  These are restricted to machines with at least 2 cores|

### linuxProfile

`linuxProfile` provides the linux configuration for each linux node in the cluster

|Name|Required|Description|
|---|---|---|
|adminUsername|yes|describes the username to be used on all linux clusters|
|ssh.publicKeys[0].keyData|yes|The public SSH key used for authenticating access to all Linux nodes in the cluster.  Here are instructions for [generating a public/private key pair](ssh.md#ssh-key-generation).|
### aadProfile

`linuxProfile` provides [AAD integration](kubernetes.aad.md) configuration for the cluster, currently only available for Kubernetes orchestrator.

|Name|Required|Description|
|---|---|---|
|clientAppID|yes|describes the client AAD application ID|
|serverAppID|yes|describes the server AAD application ID|
|tenantID|no|describes the AAD tenant ID to use for authentication. If not specified, will use the tenant of the deployment subscription.|
### extensionProfiles
A cluster can have 0 - N extensions in extension profiles.  Extension profiles allow a user to easily add pre-packaged functionality into a cluster.  An example would be configuring a monitoring solution on your cluster.  You can think of extensions like a marketplace for acs clusters.

|Name|Required|Description|
|---|---|---|
|name|yes|the name of the extension.  This has to exactly match the name of a folder under the extensions folder|
|version|yes|the version of the extension.  This has to exactly match the name of the folder under the extension name folder|
|extensionParameters|optional|extension parameters may be required by extensions.  The format of the parameters is also extension dependant.|
|rootURL|optional|url to the root location of extensions.  The rootURL must have an extensions child folder that follows the extensions convention.  The rootURL is mainly used for testing purposes.|

You can find more information, as well as a list of extensions on the [extensions documentation](extensions.md).

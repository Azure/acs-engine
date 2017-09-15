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
|kubernetesImageBase|no|This specifies the base URL (everything preceding the actual image filename) of the kubernetes hyperkube image to use for cluster deploymenbt, e.g., `gcrio.azureedge.net/google_containers/`.|
|networkPolicy|no|Specifies the network policy tool for the cluster. Valid values are:<br>`none` (default), which won't enforce any network policy,<br>`azure` for applying Azure VNET network policy,<br>`calico` for Calico network policy for clusters with Linux agents only.<br>See [network policy examples](../examples/networkpolicy) for more information.|
|clusterSubnet|no|The IP subnet used for allocating IP addresses for pod network interfaces. The subnet must be in the VNET address space. Default value is 10.244.0.0/16.|
|dnsServiceIP|no|IP address for kube-dns to listen on. If specified must be in the range of `serviceCidr`.|
|dockerBridgeSubnet|no|The specific IP and subnet used for allocating IP addresses for the docker bridge network created on the kubernetes master and agents. Default value is 172.17.0.1/16. This value is used to configure the docker daemon using the [--bip flag](https://docs.docker.com/engine/userguide/networking/default_network/custom-docker0).|
|serviceCidr|no|IP range for Service IPs, Default is "10.0.0.0/16". This range is never routed outside of a node so does not need to lie within clusterSubnet or the VNet.|
|enableRbac|no|Enable [Kubernetes RBAC](https://kubernetes.io/docs/admin/authorization/rbac/) (boolean - default == false) |
|maxPods|no|The maximum number of pods per node. The minimum valid value, necessary for running kube-system pods, is 5. Default value is 30 when networkPolicy equals azure, 110 otherwise.|
|gcHighThreshold|no|Sets the --image-gc-high-threshold value on the kublet configuration. Default is 85. [See kubelet Garbage Collection](https://kubernetes.io/docs/concepts/cluster-administration/kubelet-garbage-collection/) |
|gcLowThreshold|no|Sets the --image-gc-low-threshold value on the kublet configuration. Default is 80. [See kubelet Garbage Collection](https://kubernetes.io/docs/concepts/cluster-administration/kubelet-garbage-collection/) |

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

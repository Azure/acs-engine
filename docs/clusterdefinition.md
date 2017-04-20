# Microsoft Azure Container Service Engine - Cluster Definition

##Cluster Defintions for apiVersion "vlabs"

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

1. `DCOS` - this represents the [DC/OS orchestrator](dcos.md).  [Older versions of DCOS173 and DCOS184 may be specified](../examples/dcos-versions).
2. `Kubernetes` - this represents the [Kubernetes orchestrator](kubernetes.md).
3. `Swarm` - this represents the [Swarm orchestrator](swarm.md).
4. `Swarm Mode` - this represents the [Swarm Mode orchestrator](swarmmode.md).

### kubernetesConfig

`kubernetesConfig` describes Kubernetes specific configuration.

|Name|Required|Description|
|---|---|---|
|kubernetesImageBase|no|This specifies the image of kubernetes to use for the cluster.|
|networkPolicy|no|Specifies the network policy tool for the cluster. Default is `none`, which won't enforce network policy. This can be set to `calico` for clusters with Linux agents only.|
|clusterCidr|no|Pod IP address range if you wish to change the default.|
|dnsServiceIP|no|IP address for kube-dns to listen on. If specified must be in the range of `serivceCidr`.|
|serviceCidr|no|IP range for Service IPs, Default is "10.0.0.0/16". This range is never routed outside of a node so does not need to lie within clusterCidr or the VNet.|

### masterProfile
`masterProfile` describes the settings for master configuration.

|Name|Required|Description|
|---|---|---|
|count|yes|Masters have count value of 1, 3, or 5 masters|
|dnsPrefix|yes|this is the dns prefix for the masters FQDN.  The master FQDN is used for SSH or commandline access. This must be a unique name. ([bring your own VNET examples](../examples/vnet))|
|firstConsecutiveStaticIP|only required when vnetSubnetId specified|this is the IP address of the first master.  IP Addresses will be assigned consecutively to additional master nodes.|
|vmsize|yes|Describes a valid [Azure VM Sizes](https://azure.microsoft.com/en-us/documentation/articles/virtual-machines-windows-sizes/).  These are restricted machines with at least 2 cores and 100GB of ephemeral disk space.|
|vnetSubnetId|no|specifies the Id of an alternate VNET subnet.  The subnet id must specify a valid VNET ID owned by the same subscription. ([bring your own VNET examples](../examples/vnet))|

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
|servicePrincipalClientID|yes, for Kubernetes clusters|describes the Azure client id.  It is recommended to use a separate client ID per cluster|
|servicePrincipalClientSecret|yes, for Kubernetes clusters|describes the Azure client secret.  It is recommended to use a separate client secret per client id|

##Cluster Defintions for apiVersion "2016-03-30"

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

# Microsoft Azure Container Service Engine - Network Policy

There are 3 different Network Policy options :

- Azure Container Networking (default)
- Calico
- Kubenet (none)

Please note that only the `calico` network policy supports the Kubernetes notion
of network policies.

## Azure Container Networking (default)

By default (currently Linux clusters only), the `azure` network policy is applied. It is an open source implementation of [the CNI Network Plugin interface](https://github.com/containernetworking/cni/blob/master/SPEC.md) and [the CNI Ipam plugin interface](https://github.com/containernetworking/cni/blob/master/SPEC.md#ip-address-management-ipam-interface)

CNI brings the containers to a single flat L3 Azure subnet. This enables full integration with other SDN features such as network security groups and VNET peering. The plugin creates a bridge for each underlying Azure VNET. The bridge functions in L2 mode and is connected to the host network interface.

If the container host VM has multiple network interfaces, the primary network interface is reserved for management traffic. A secondary interface is used for container traffic whenever possible.

More detailed documentation can be found in [the Azure Container Networking Repository](https://github.com/Azure/azure-container-networking/tree/master/docs)

Example of templates enabling CNI:

```json
  "properties": {
    "orchestratorProfile": {
      "orchestratorType": "Kubernetes",
      "kubernetesConfig": {
        "networkPolicy": "azure"
      }
    }
    ...
  }
```

Or by not specifying any network policy, leaving the default :

```json
    "properties": {
    "orchestratorProfile": {
      "orchestratorType": "Kubernetes"
    }
    ...
  }
```

## Calico

The kubernetes-calico deployment template enables Calico networking and policies for the ACS-engine cluster via `"networkPolicy": "calico"` being present inside the `kubernetesConfig`.

```json
  "properties": {
    "orchestratorProfile": {
      "orchestratorType": "Kubernetes",
      "kubernetesConfig": {
        "networkPolicy": "calico"
      }
```

If `"orchestratorRelease": "1.8",` is set a K8s 1.8.x cluster will be provisioned.  If `orchestratorRelease` is not specified a K8s 1.7.x cluster will be deployed.  In either of these cases, this template will deploy the [v2.6 release](https://docs.projectcalico.org/v2.6/releases/) of [Kubernetes Datastore Install](https://docs.projectcalico.org/v2.6/getting-started/kubernetes/installation/hosted/kubernetes-datastore/) version of calico with the "Calico policy-only with user-supplied networking" which supports kubernetes ingress policies and has some limitations as denoted on the referenced page.  

> Note: If deploying on a K8s 1.8 cluster, then egress policies are also supported!

If `orchestratorRelease` is set to 1.5 or 1.6, then this template will deploy the [v2.4.1 release](https://github.com/projectcalico/calico/releases/tag/v2.4.1) of [Kubernetes Datastore Install](https://docs.projectcalico.org/v2.4/getting-started/kubernetes/installation/hosted/kubernetes-datastore/) version of calico with the "Calico policy-only with user-supplied networking" which supports kubernetes ingress policies and has some limitations as denoted on the referenced page.

To understand how to deploy this template, please read the baseline  [Kubernetes](../../docs/kubernetes.md) document and simply make sure to use the **kubernetes-calico.json** file in this folder which has the above referenced line to enable.

### Post installation

Once the template has been successfully deployed, following the [simple policy tutorial](https://docs.projectcalico.org/v2.6/getting-started/kubernetes/tutorials/simple-policy) or the [advanced policy tutorial](https://docs.projectcalico.org/v2.6/getting-started/kubernetes/tutorials/advanced-policy) will help to understand calico networking.

> Note: `ping` (ICMP) traffic is blocked on the cluster by default.  Wherever `ping` is used in any tutorial substitute testing access with something like `wget -q --timeout=5 google.com -O -` instead.

## Kubenet (none)

Also available is the Kubernetes-native kubenet implementation, which is declared as configuration thusly:

```json
  "properties": {
    "orchestratorProfile": {
      "orchestratorType": "Kubernetes",
      "kubernetesConfig": {
        "networkPolicy": "none"
      }
    }
    ...
  }
```

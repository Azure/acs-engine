# Microsoft Azure Container Service Engine - Network Policy

There are 2 different Network Policy options :

- Calico
- Cilium (docs are //TODO)

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

This template will deploy the [v3.1 release](https://docs.projectcalico.org/v3.1/releases/) of [Kubernetes Datastore Install](https://docs.projectcalico.org/v3.1/getting-started/kubernetes/installation/other) version of calico with the "Calico for policy" with user-supplied networking which supports kubernetes ingress policies.

> Note: The Typha service and deployment is installed on the cluster, but effectively disabled using the default settings of deployment replicas set to 0 and Typha service name not configured.  Typha is recommended to be enabled when scaling to 50+ nodes on the cluster to reduce the load on the Kubernetes API server.  If this functionality is desired to be configurable via the API model, please file an issue on Github requesting this feature be added.  Otherwise, this can be manually changed via modifying and applying changes with the `/etc/kubernetes/addons/calico-daemonset.yaml` file on every master node in the cluster.

If deploying on a K8s 1.8 or later cluster, then egress policies are also supported!

To understand how to deploy this template, please read the baseline [Kubernetes](../../docs/kubernetes.md) document, and use the example **kubernetes-calico.json** file in this folder as an api model reference.

### Post installation

Once the template has been successfully deployed, following the [simple policy tutorial](https://docs.projectcalico.org/v3.1/getting-started/kubernetes/tutorials/simple-policy) or the [advanced policy tutorial](https://docs.projectcalico.org/v3.1/getting-started/kubernetes/tutorials/advanced-policy) will help to understand calico networking.

> Note: `ping` (ICMP) traffic is blocked on the cluster by default.  Wherever `ping` is used in any tutorial substitute testing access with something like `wget -q --timeout=5 google.com -O -` instead.

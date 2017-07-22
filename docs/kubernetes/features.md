## Features

|Feature|Status|API Version|Example|Description|
|---|---|---|---|---|
|Managed Identity|Alpha|`vlabs`|[kubernetes-msi.json](../examples/managed-identity/kubernetes-msi.json)|[Description](#feat-kubernetes-msi)|

<a name="feat-kubernetes-msi"></a>
### Managed Identity

Enabling Managed Identity configures acs-engine to include and use MSI identities for all interactions with the Azure Resource Manager (ARM) API.

Instead of using a static servic principal written to `/etc/kubernetes/azure.json`, Kubernetes will use a dynamic, time-limited token fetched from the MSI extension running on master and agent nodes. This support is currently alpha and requires Kubernetes v1.7.2 or newer.

Enable Managed Identity by adding `useManagedIdentity` in `kubernetesConfig`.

```json
"kubernetesConfig": {
  "useManagedIdentity": true,
  "customHyperkubeImage": "docker.io/colemickens/hyperkube-amd64:3b15e8a446fa09d68a2056e2a5e650c90ae849ed"
}
```

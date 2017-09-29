# Microsoft Azure Container Service Engine - Network Policy

## Overview

The kubernetes-calico deployment template enables Calico networking and policies for the ACS-engine cluster via `"networkPolicy": "calico"` being present inside the `kubernetesConfig`.

```
  "properties": {
    "orchestratorProfile": {
      "orchestratorType": "Kubernetes",
      "kubernetesConfig": {
        "networkPolicy": "calico"
      }
```

This deploys the [v2.6.0 release](https://docs.projectcalico.org/v2.6/releases/) of [Kubernetes Datastore Install](https://docs.projectcalico.org/v2.6/getting-started/kubernetes/installation/hosted/kubernetes-datastore/) version of calico with the "Calico policy-only with user-supplied networking" which supports kubernetes ingress policies and has some limitations as denoted on the referenced page.

To understand how to deploy this template, please read the baseline  [Kubernetes](../../docs/kubernetes.md) document and simply make sure to use the **kubernetes-calico.json** file in this folder which has the above referenced line to enable.

> Note: this version of Calico _must_ be deployed on a 1.7.x or later version of K8s in order to function

## Post installation

Once the template has been successfully deployed, following the [simple policy tutorial](https://docs.projectcalico.org/v2.6/getting-started/kubernetes/tutorials/simple-policy) will help to understand calico networking.

> Note: `ping` (ICMP) traffic is blocked on the cluster.  Wherever `ping` is used in any tutorial substitute testing access with something like `wget -q --timeout=5 google.com -O -` instead.
# Microsoft Azure Container Service Engine - Network Policy

## Overview

The kubernetes-calico deployment template enables Calico networking and policies for the ACS-engine cluster via `"networkPolicy": "calico"` being present inside the `kubernetesConfig`.

```
  "properties": {
    "orchestratorProfile": {
      "orchestratorType": "Kubernetes",
      "orchestratorRelease": "1.6",
      "orchestratorVersion": "1.6.9",
      "kubernetesConfig": {
        "networkPolicy": "calico"
      }
```

This deploys the [v2.5.1 release](https://docs.projectcalico.org/v2.5/releases/) of [Standard Hosted Install](https://docs.projectcalico.org/v2.5/getting-started/kubernetes/installation/hosted/hosted) version of calico which has its own calico policy controller connected to the etcd hosted on the master K8s nodes in its own `/calico` etcd namespace without any backend IP configuration. This installation option allows the potential for full calico functionality including egress policy control compared to the "etcdless" [Kubernetes Datastore hosted installation](https://docs.projectcalico.org/v2.5/getting-started/kubernetes/installation/hosted/kubernetes-datastore/).

To understand how to deploy this template, please read the baseline  [Kubernetes](../../docs/kubernetes.md) document and simply make sure to use the **kubernetes-calico.json** file which has the above referenced line to enable.

> Note: this version of Calico has been verified to _not_ work with the 1.7.5 version of K8s on ACS-engine, hence why the template is coded to use 1.6.9.  K8s 1.7.x will be used with the improvements in Calico 2.6.x.

## Post installation

Once the template has been successfully deployed, following the [advanced policy tutorial](https://docs.projectcalico.org/v2.5/getting-started/kubernetes/tutorials/advanced-policy) will help to understand calico networking.

> Note: `ping` (ICMP) traffic is blocked on the cluster.  Wherever `ping` is used in the tutorial substitute testing access with `wget -q --timeout=5 google.com -O -` instead.

The `calicoctl` binary is required for the tutorial and to manage the policies. Calicoctl requires a connection to the `/calico` etcd namespace on the K8s master node, so here are options in order to make `calicoctl` function:

### Install and run on a master node

1. ssh into one of the K8s master nodes and install the `calicoctl` binary

  ```bash
  wget https://github.com/projectcalico/calicoctl/releases/download/v1.5.0/calicoctl
  chmod +x calicoctl
  ```
2. Now simply use the binary via `.\calicoctl`
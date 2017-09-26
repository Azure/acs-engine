# Deploy a Kubernetes Cluster

## Install Pre-requisites

All the commands in this guide require both the Azure CLI and `acs-engine`. Follow the [installation instructions to download acs-engine before continuing](../acsengine.md#install-acs-engine) or [compile from source](../acsengine.md#build-from-source).

For installation instructions see [the Azure CLI GitHub repository](https://github.com/Azure/azure-cli#installation) for the latest release.

## Overview

`acs-engine` reads a cluster definition which describes the size, shape, and configuration of your cluster. This guide takes the default configuration of one master and two linux agents. If you would like to change the configuration, edit `examples/kubernetes.json` before continuing.

The `acs-engine deploy` command automates creation of a Service Principal, Resource Group and SSH key for your cluster. If operators need more control or are are intersted in the individual steps see the ["Long Way" section below](#the-long-way).

## Gather Information

* The subscription in which you would like to provision the cluster. This is a uuid which can be found with `az account list -o table`.
* A `dnsPrefix` which forms part of the the hostname for your cluster (e.g. staging, prodwest, blueberry). The DNS prefix must be unique so pick a random name.
* A location to provision the cluster e.g. `westus2`.

```
$ az account list -o table
Name                                             CloudName    SubscriptionId                        State    IsDefault
-----------------------------------------------  -----------  ------------------------------------  -------  -----------
Contoso Subscription                             AzureCloud   51ac25de-afdg-9201-d923-8d8e8e8e8e8e  Enabled  True
```

## Deploy

For this example, the subscription id is `51ac25de-afdg-9201-d923-8d8e8e8e8e8e`, the DNS prefix is `contoso-apple`, and location is `westus2`.

Run `acs-engine deploy` with the appropriate argumets:

```
$ acs-engine deploy --subscription-id 51ac25de-afdg-9201-d923-8d8e8e8e8e8e \
    --dns-prefix contoso-apple --location westus2 \
    --auto-suffix --api-model examples/kubernetes.json

WARN[0005] apimodel: missing masterProfile.dnsPrefix will use "contoso-apple-59769a59"
WARN[0005] --resource-group was not specified. Using the DNS prefix from the apimodel as the resource group name: contoso-apple-59769a59
WARN[0008] apimodel: ServicePrincipalProfile was empty, creating application...
WARN[0017] created application with applicationID (7e2d433f-d039-48b8-87dc-83fa4dfa38d4) and servicePrincipalObjectID (db6167e1-aeed-407a-b218-086589759442).
WARN[0017] apimodel: ServicePrincipalProfile was empty, assigning role to application...
INFO[0034] Starting ARM Deployment (contoso-apple-59769a59-1423145182). This will take some time...
INFO[0393] Finished ARM Deployment (contoso-apple-59769a59-1423145182).
```

`acs-engine` will output Azure Resource Manager (ARM) templates, SSH keys, and a kubeconfig file in `_output/contoso-apple-59769a59` directory:

   * `_output/contoso-apple-59769a59/azureuser_rsa`
   * `_output/contoso-apple-59769a59/kubeconfig/kubeconfig.uswest2.json`

Acs-engine generates kubeconfig files for each possible region. Access the new cluster by using the kubeconfig generated for the cluster's location. This example used `uswest2`, so the kubeconfig is `_output/<clustername>/kubeconfig/kubeconfig.uswest2.json`:

```
$ KUBECONFIG=_output/contoso-apple-59769a59/kubeconfig/kubeconfig.westus2.json kubectl cluster-info
Kubernetes master is running at https://contoso-apple-59769a59.westus2.cloudapp.azure.com
Heapster is running at https://contoso-apple-59769a59.westus2.cloudapp.azure.com/api/v1/proxy/namespaces/kube-system/services/heapster
KubeDNS is running at https://contoso-apple-59769a59.westus2.cloudapp.azure.com/api/v1/proxy/namespaces/kube-system/services/kube-dns
kubernetes-dashboard is running at https://contoso-apple-59769a59.westus2.cloudapp.azure.com/api/v1/proxy/namespaces/kube-system/services/kubernetes-dashboard

To further debug and diagnose cluster problems, use 'kubectl cluster-info dump'.
```

**Note**: If the cluster is using an existing VNET please see the [Custom VNET](features.md#feat-custom-vnet) feature documentation for additional steps that must be completed after cluster provisioning.

<a href="#the-long-way"></a>

## ACS Engine the Long Way

### Step 1: Generate an SSH Key

In addition to using Kubernetes APIs to interact with the clusters, cluster operators may access the master and agent machines using SSH.

If you don't have an SSH key [cluster operators may generate a new one](../ssh.md#ssh-key-generation).

### Step 2: Create a Service Principal

Kubernetes clusters have integrated support for various cloud providers as core functionality. On Azure, acs-engine uses a Service Principal to interact with Azure Resource Manager (ARM). Follow the instructions to [create a new service principal](../serviceprincipal.md)

### Step 3: Edit your Cluster Definition

ACS Engine consumes a cluster definition which outlines the desired shape, size, and configuration of Kubernetes. There are a number of features that can be enabled through the cluster definition, check the `examples` directory for a number of... examples.

Edit the [simple Kubernetes cluster definition](/examples/kubernetes.json) and fill out the required values:

* `dnsPrefix`: must be a region-unique name and will form part of the hostname (e.g. myprod1, staging, leapinglama), be unique!
* `keyData`: must contain the public portion of an SSH key, this will be associated with the `adminUsername` value found in the same section of the cluster definition (e.g. 'ssh-rsa AAAAB3NzaC1yc2EAAAADAQABA....')
* `clientId`: this is the appId uuid or name from step 2
* `secret`: this is the password or randomly-generated password from step 2

### Step 4: Generate the Templates

The generate command takes a cluster definition and outputs a number of templates which describe your Kubernetes cluster. By default, `generate` will create a new directory named after your cluster nested in the `_output` directory. If my dnsPrefix was `larry` my cluster templates would be found in `_output/larry-`.

Run `acs-engine generate examples/kubernetes.json`

### Step 5: Submit your Templates to Azure Resource Manager (ARM)

[deploy the output azuredeploy.json and azuredeploy.parameters.json](../acsengine.md#deployment-usage)
  * To enable the optional network policy enforcement using calico, you have to
    set the parameter during this step according to this [guide](../kubernetes.md#optional-enable-network-policy-enforcement-using-calico)


**Note**: If the cluster is using an existing VNET please see the [Custom VNET](features.md#feat-custom-vnet) feature documentation for additional steps that must be completed after cluster provisioning.

# Troubleshooting

## VMExtensionProvisioningError or VMExtensionProvisioningTimeout

The two above VMExtensionProvisioning— errors tell us that a vm in the cluster failed installing required application prerequisites after CRP provisioned the VM into the resource group. When acs-engine creates a new Kubernetes cluster, a series of shell scripts runs to install prereq's like docker, etcd, Kubernetes runtime, and various other host OS packages that support the Kubernetes application layer. *Usually* this indicates one of the following:

1. Something about the cluster configuration is pathological. For example, perhaps the cluster config includes a custom version of a particular software dependency that doesn't exist. Or, another example, for a cluster created inside a custom VNET (i.e., a user-provided, pre-existing VNET), perhaps that custom VNET does not have general outbound internet access, and so apt, docker pull, etc is not able to execute successfully.
2. A transient Azure environmental error caused the shell script operation to timeout, or exceed its retry count. For example, the shell script may attempt to download a required package (e.g., etcd), and if the Azure networking environment for the newly provisioned vm is flaky for a period of time, then the shell script may retry several times, but eventually timeout and fail.

For classification #1 above, the appropriate strategic response is to figure out what about the cluster configuration is incorrect, and to fix it. We expect such scenarios to always fail in the above way: cluster deployments will not be successful until the cluster configuration is made to be correct.

For classification #2 above, the appropriate strategic response is to retry a few times. If a 2nd or 3rd attempt succeeds, it is a hint that a transient environmental condition is the cause of the initial failure.

### What is CSE?

CSE stands for CustomScriptExtension, and is just a way of expressing: "a script that executes as part of the VM provisioning process, and that must exit 0 (i.e., successfully) in order for that VM provisioning process to succeed". Basically it's another way of expressing the VMExtensionProvisioning— concept above.

To summarize, the way that acs-engine implements Kubernetes on Azure is a collection of (1) Azure VM configuration + (2) shell script execution. Both are implemented as a single operational unit, and when #2 fails, we consider the entire VM provisioning operation to be a failure; more importantly, if only one VM in the cluster deployment fails, we consider the entire cluster operation to be a failure.

### How To Debug CSE errors

In order to troubleshoot a cluster that failed in the above way(s), we need to grab the CSE logs from the host VM itself.

From a vm node that did not provision successfully:

- grab the entire file at `/var/log/azure/cluster-provision.log`

- grab the entire file at `/var/log/cloud-init-output.log`

How to determine the above?

1. Look at the deployment error message. The error should include which VM extension failed the deployment. For example, `cse-master-0` means that the CSE extension of VM master 0 failed.

2. From a master node: `kubectl get nodes`

- Are there any missing master or agent nodes?
  - if so, that node vm probably failed CSE: grab the log files above from that vm
- Are there no working nodes?
  - if so, grab the log files above from the master vm you are on

#### CSE Exit Codes

```
"code": "VMExtensionProvisioningError"
"message": "VM has reported a failure when processing extension 'cse1'. Error message: "Enable failed: failed to
execute command: command terminated with exit status=20\n[stdout]\n\n[stderr]\n"."
```

Look for the exit code. In the above example, the exit code is `20`. The list of exit codes and their meaning can be found [here](../../parts/k8s/kubernetescustomscript.sh).

If after following the above you are still unable to troubleshoot your deployment error, please open a Github issue with title "CSE error: exit code <INSERT_YOUR_EXIT_CODE>" and include the following in the description:

1. The apimodel json used to deploy the cluster (aka your cluster config). **Please make sure you remove all secrets and keys before posting it on GitHub.**

2. The output of `kubectl get nodes`

3. The content of `/var/log/azure/cluster-provision.log` and `/var/log/cloud-init-output.log`


# Misconfigured Service Principal

If your Service Principal is misconfigured, none of the Kubernetes components will come up in a healthy manner.
You can check to see if this the problem:

```shell
ssh -i ~/.ssh/id_rsa USER@MASTERFQDN sudo journalctl -u kubelet | grep --text autorest
```

If you see output that looks like the following, then you have **not** configured the Service Principal correctly.
You may need to check to ensure the credentials were provided accurately, and that the configured Service Principal has
read and **write** permissions to the target Subscription.

`Nov 10 16:35:22 k8s-master-43D6F832-0 docker[3177]: E1110 16:35:22.840688    3201 kubelet_node_status.go:69] Unable to construct api.Node object for kubelet: failed to get external ID from cloud provider: autorest#WithErrorUnlessStatusCode: POST https://login.microsoftonline.com/72f988bf-86f1-41af-91ab-2d7cd011db47/oauth2/token?api-version=1.0 failed with 400 Bad Request: StatusCode=400`

[This documentation](../serviceprincipal.md) explains how to create/configure a service principal for an ACS-Engine Kubernetes cluster.

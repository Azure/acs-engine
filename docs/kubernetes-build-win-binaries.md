# Building Windows Kubernetes Binaries and deploy to an Azure storage account

The following instructions show how to deploy the Windows Kubernetes Binaries and deploy them to an Azure Storage Account.

### Prerequisites
* Azure Storage Account and Azure Storage Container to store Windows binaries
* Access to [winnat.sys](https://blogs.technet.microsoft.com/virtualization/2016/05/25/windows-nat-winnat-capabilities-and-limitations/) stored in a storage container. (WinNAT) is used to provide required NAT networking functionality for Windows containers that will be included in a future Windows image update.
* Docker installed and running. MacOS users using Docker for Mac must have at [least 3GB of memory allocated to Docker](https://github.com/kubernetes/kubernetes/tree/master/build/#requirements) or building will likely fail.

## Background
Microsoft maintains a fork of the Kubernetes project at https://github.com/Azure/kubernetes which includes patches not yet included in upstream Kubernetes; these are needed for Windows containers to function.

[build-windows-k8s.sh](../scripts/build-windows-k8s.sh) does the following:
- Checks out the fork of Azure/kubernetes (includes Windows fixes not yet in upstream Kubernetes, needed for Windows containers to function)
- Builds kubelet.exe and kube-proxy.exe from source in a Docker container
- Downloads kubectl.exe for desired release
- Downloads [NSSM](https://nssm.cc) which is used to start kubelet and kube-proxy on Windows
- Downloads [Windows NAT](https://blogs.technet.microsoft.com/virtualization/2016/05/25/windows-nat-winnat-capabilities-and-limitations/)
- Creates an .zip archive of these Windows components
- Uploads archive to Azure Blob Storage

More information about building Kubernetes binaries from source here: https://github.com/kubernetes/kubernetes/tree/master/build/

### Set Azure Storage credentials and Container name

A storage container is used to download winnat.sys during the build phase and to upload the resulting archive artifact.
```
$ export AZURE_STORAGE_CONNECTION_STRING="DefaultEndpointsProtocol=https;AccountName=MyStorageAccountName;AccountKey=..." \
AZURE_STORAGE_CONTAINER_NAME=MyStorageContainerName
```

### Run `make build-windows-k8s`

Usage: `make build-windows-k8s K8S_VERSION=1.7.4 PATCH_VERSION=1`

This will build Kubernetes binaries for Windows agents based on Kubernetes release version 1.7.4 and place the zip artifact in Azure Blob Storage.

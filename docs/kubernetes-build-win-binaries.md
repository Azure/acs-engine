# Building Windows Kubernetes Binaries and deploy to an Azure storage account

The following instructions show how to deploy the Windows Kubernetes Binaries and deploy them to an Azure Storage Account.

### Prerequisites
* Azure Storage Account and Azure Storage Container to store Windows binaries
* Access to [winnat.sys](https://blogs.technet.microsoft.com/virtualization/2016/05/25/windows-nat-winnat-capabilities-and-limitations/) stored in a storage container. (WinNAT) is used to provide required NAT networking functionality for Windows containers that will be included in a future Windows image update.

### Set Azure Storage credentials and Container name

A storage container is used to download winnat.sys during the build phase and to upload the resulting archive artifact.
```
$ export AZURE_STORAGE_CONNECTION_STRING="DefaultEndpointsProtocol=https;AccountName=MyStorageAccountName;AccountKey=..." \
AZURE_STORAGE_CONTAINER_NAME=MyStorageContainerName
```

### Run `make build-windows-k8s`

The `make build-windows-k8s` will do the following:
- Clone the fork of Azure/kubernetes which include Windows fixes not yet in upstream Kubernetes
- Build kubelet.exe, kube-proxy.exe
- Download kubectl.exe for desired release
- Download [NSSM](https://nssm.cc) which is used to start kubelet and kube-proxy on Windows
- Download [Windows NAT](https://blogs.technet.microsoft.com/virtualization/2016/05/25/windows-nat-winnat-capabilities-and-limitations/)
- Create an .zip archive of these Windows components
- Upload archive to Azure Blob Storage

Usage: `make build-windows-k8s K8S_VERSION=1.7.4 PATCH_VERSION=1`

This will build Kubernetes binaries for Windows agents based on Kubernetes release version 1.7.4.

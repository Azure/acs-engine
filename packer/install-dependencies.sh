#!/bin/bash

source /home/packer/provision_installs.sh
source /home/packer/provision_source.sh

# TODO: deal with etcd versions
ETCD_VERSION="3.2.23"
ETCD_DOWNLOAD_URL="https://acs-mirror.azureedge.net/github-coreos"
installEtcd

installDeps

DOCKER_REPO="https://apt.dockerproject.org/repo"
DOCKER_ENGINE_VERSION="1.13.*"
installDocker

installClearContainersRuntime

VNET_CNI_PLUGINS_URL="https://acs-mirror.azureedge.net/cni/azure-vnet-cni-linux-amd64-v1.0.10.tgz"
CNI_PLUGINS_URL="https://acs-mirror.azureedge.net/cni/cni-plugins-amd64-v0.7.1.tgz"

installAzureCNI

CONTAINERD_DOWNLOAD_URL_BASE="https://storage.googleapis.com/cri-containerd-release/"
installContainerd

for KUBERNETES_VERSION in 1.8.15 1.9.10 1.10.6 1.11.2; do
    HYPERKUBE_URL="k8s.gcr.io/hyperkube-amd64:v${KUBERNETES_VERSION}"
    pullHyperkube
done

echo "Install completed successfully on " `date` > /var/log/azure/golden-image-install.complete
#!/bin/bash

source /home/packer/provision_installs.sh
source /home/packer/provision_source.sh

ETCD_VERSION="3.2.23"
ETCD_DOWNLOAD_URL="https://acs-mirror.azureedge.net/github-coreos"
installEtcd

installDeps

DOCKER_REPO="https://apt.dockerproject.org/repo"
DOCKER_ENGINE_VERSION="1.13.*"
installDocker

installClearContainersRuntime

VNET_CNI_VERSIONS="1.0.10 1.0.11"
CNI_PLUGIN_VERSIONS="0.7.1"

for VNET_CNI_VERSION in $VNET_CNI_VERSIONS; do
    VNET_CNI_PLUGINS_URL="https://acs-mirror.azureedge.net/cni/azure-vnet-cni-linux-amd64-v${VNET_CNI_VERSION}.tgz"
    downloadAzureCNI
done

for CNI_PLUGIN_VERSION in $CNI_PLUGIN_VERSIONS; do
    CNI_PLUGINS_URL="https://acs-mirror.azureedge.net/cni/cni-plugins-amd64-v${CNI_PLUGIN_VERSION}.tgz"
    downloadCNI
done

CONTAINERD_DOWNLOAD_URL_BASE="https://storage.googleapis.com/cri-containerd-release/"
installContainerd

installImg

# TODO: fetch supported k8s versions from an acs-engine command instead of hardcoding them here
K8S_VERSIONS="1.7.7 1.7.9 1.7.12 1.7.15 1.7.16 1.8.0 1.8.1 1.8.2 1.8.4 1.8.5 1.8.6 1.8.7 1.8.8 1.8.9 1.8.10 1.8.11 1.8.12 1.8.13 1.8.14 1.8.15 1.9.0 1.9.1 1.9.2 1.9.3 1.9.4 1.9.5 1.9.6 1.9.7 1.9.8 1.9.9 1.9.10 1.10.0 1.10.1 1.10.2 1.10.3 1.10.4 1.10.5 1.10.6 1.10.7 1.11.0 1.11.1 1.11.2"

for KUBERNETES_VERSION in ${K8S_VERSIONS}; do
    HYPERKUBE_URL="k8s.gcr.io/hyperkube-amd64:v${KUBERNETES_VERSION}"
    pullHyperkube
done

df -h

echo "Install completed successfully on " `date` > /var/log/azure/golden-image-install.complete
#!/bin/bash
set -x
source provision-source.sh
source kubernetescustomscript.sh

OS=$(cat /etc/*-release | grep ^ID= | tr -d 'ID="' | awk '{print toupper($0)}')
UBUNTU_OS_NAME="UBUNTU"
KUBECTL=/usr/local/bin/kubectl
DOCKER=/usr/bin/docker
CNI_BIN_DIR=/opt/cni/bin

DOCKER_REPO="https://apt.dockerproject.org/repo"
DOCKER_ENGINE_VERSION="1.13.*"
ADMINUSER="azureuser"
VNET_CNI_PLUGINS_URL="https://acs-mirror.azureedge.net/cni/azure-vnet-cni-linux-amd64-latest.tgz"
HYPERKUBE_URL="k8s-gcrio.azureedge.net/hyperkube-amd64:v1.10.5"

installDeps
installDocker
configAzureCNI
installContainerd
extractHyperkube

echo "Install complete successfully" > /var/log/azure/golden-image-install.complete

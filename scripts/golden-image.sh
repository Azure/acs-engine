#!/bin/bash
set -x
source provision-source.sh

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

function installDeps() {
    apt_get_update
    apt_get_install 20 30 300 apt-transport-https ca-certificates iptables iproute2 ebtables socat util-linux mount ethtool init-system-helpers nfs-common ceph-common conntrack glusterfs-client ipset jq cgroup-lite git pigz xz-utils || exit $ERR_APT_INSTALL_TIMEOUT
}

function installDocker() {
    retrycmd_if_failure_no_stats 20 1 5 curl -fsSL https://aptdocker.azureedge.net/gpg > /tmp/aptdocker.gpg
    retrycmd_if_failure 10 5 10 apt-key add /tmp/aptdocker.gpg
    echo "deb $DOCKER_REPO ubuntu-xenial main" | sudo tee /etc/apt/sources.list.d/docker.list
    printf "Package: docker-engine\nPin: version $DOCKER_ENGINE_VERSION\nPin-Priority: 550\n" > /etc/apt/preferences.d/docker.pref
    apt_get_update
    apt_get_install 20 30 120 docker-engine
    echo "ExecStartPost=/sbin/iptables -P FORWARD ACCEPT" >> /etc/systemd/system/docker.service.d/exec_start.conf
    usermod -aG docker ${ADMINUSER}
}

function installCNI() {
    mkdir -p $CNI_BIN_DIR
    CONTAINERNETWORKING_CNI_TGZ_TMP=/tmp/containernetworking_cni.tgz
    retrycmd_get_tarball 60 5 $CONTAINERNETWORKING_CNI_TGZ_TMP $CNI_PLUGINS_URL
    tar -xzf $CONTAINERNETWORKING_CNI_TGZ_TMP -C $CNI_BIN_DIR
    chown -R root:root $CNI_BIN_DIR
    chmod -R 755 $CNI_BIN_DIR
}

function configAzureCNI() {
    CNI_CONFIG_DIR=/etc/cni/net.d
    mkdir -p $CNI_CONFIG_DIR
    chown -R root:root $CNI_CONFIG_DIR
    chmod 755 $CNI_CONFIG_DIR
    mkdir -p $CNI_BIN_DIR
    AZURE_CNI_TGZ_TMP=/tmp/azure_cni.tgz
    retrycmd_get_tarball 60 5 $AZURE_CNI_TGZ_TMP $VNET_CNI_PLUGINS_URL
    tar -xzf $AZURE_CNI_TGZ_TMP -C $CNI_BIN_DIR
    installCNI
    mv $CNI_BIN_DIR/10-azure.conflist $CNI_CONFIG_DIR/
    chmod 600 $CNI_CONFIG_DIR/10-azure.conflist
    /sbin/ebtables -t nat --list
}

function configNetworkPlugin() {
    configAzureCNI
	#installCNI
}

function installContainerd() {
	CRI_CONTAINERD_VERSION="1.1.0"
	CONTAINERD_DOWNLOAD_URL="https://storage.googleapis.com/cri-containerd-release/cri-containerd-$CRI_CONTAINERD_VERSION.linux-amd64.tar.gz"

    CONTAINERD_TGZ_TMP=/tmp/containerd.tar.gz
    retrycmd_get_tarball 60 5 "$CONTAINERD_TGZ_TMP" "$CONTAINERD_DOWNLOAD_URL"
	tar -xzf "$CONTAINERD_TGZ_TMP" -C /
	rm -f "$CONTAINERD_TGZ_TMP"
	sed -i '/\[Service\]/a ExecStartPost=\/sbin\/iptables -P FORWARD ACCEPT' /etc/systemd/system/containerd.service

	echo "Successfully installed cri-containerd..."
}

function extractHyperkube(){
    TMP_DIR=$(mktemp -d)
    retrycmd_if_failure 100 1 30 curl -sSL -o /usr/local/bin/img "https://acs-mirror.azureedge.net/img/img-linux-amd64-v0.4.6"
    chmod +x /usr/local/bin/img
    retrycmd_if_failure 75 1 60 img pull "k8s-gcrio.azureedge.net/hyperkube-amd64:v1.10.5"
    path=$(find /tmp/img -name "hyperkube")

   
    cp "$path" "/usr/local/bin/kubelet"
    cp "$path" "/usr/local/bin/kubectl"
    chmod a+x /usr/local/bin/kubelet /usr/local/bin/kubectl

    rm -rf /tmp/hyperkube.tar "/tmp/img"
}


installDeps
installDocker
configNetworkPlugin
installContainerd
extractHyperkube

echo "Install complete successfully"

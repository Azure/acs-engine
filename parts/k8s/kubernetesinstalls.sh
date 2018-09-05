#!/bin/bash

CC_SERVICE_IN_TMP=/opt/azure/containers/cc-proxy.service.in
CC_SOCKET_IN_TMP=/opt/azure/containers/cc-proxy.socket.in
CNI_CONFIG_DIR="/etc/cni/net.d"
CNI_BIN_DIR="/opt/cni/bin"
CNI_DOWNLOADS_DIR="/opt/cni/downloads"

function removeEtcd() {
    rm -rf /usr/bin/etcd
}

function installEtcd() {
    CURRENT_VERSION=$(etcd --version | grep "etcd Version" | cut -d ":" -f 2 | tr -d '[:space:]')
    if [[ "$CURRENT_VERSION" == "${ETCD_VERSION}" ]]; then
        echo "etcd version ${ETCD_VERSION} is already installed, skipping download"
    else
        retrycmd_get_tarball 60 10 /tmp/etcd-v${ETCD_VERSION}-linux-amd64.tar.gz ${ETCD_DOWNLOAD_URL}/etcd-v${ETCD_VERSION}-linux-amd64.tar.gz || exit $ERR_ETCD_DOWNLOAD_TIMEOUT
        removeEtcd
        tar -xzvf /tmp/etcd-v${ETCD_VERSION}-linux-amd64.tar.gz -C /usr/bin/ --strip-components=1 || exit $ERR_ETCD_DOWNLOAD_TIMEOUT
    fi
}

function installDeps() {
    retrycmd_if_failure_no_stats 20 1 5 curl -fsSL https://packages.microsoft.com/config/ubuntu/16.04/packages-microsoft-prod.deb > /tmp/packages-microsoft-prod.deb || exit $ERR_MS_PROD_DEB_DOWNLOAD_TIMEOUT
    retrycmd_if_failure 60 5 10 dpkg -i /tmp/packages-microsoft-prod.deb || exit $ERR_MS_PROD_DEB_PKG_ADD_FAIL
    apt_get_update || exit $ERR_APT_UPDATE_TIMEOUT
    # See https://github.com/kubernetes/kubernetes/blob/master/build/debian-hyperkube-base/Dockerfile#L25-L44
    apt_get_install 20 30 300 apt-transport-https ca-certificates iptables iproute2 ebtables socat util-linux mount ethtool init-system-helpers nfs-common ceph-common conntrack glusterfs-client ipset jq cgroup-lite git pigz xz-utils blobfuse fuse cifs-utils || exit $ERR_APT_INSTALL_TIMEOUT
}

function installGPUDrivers() {
    nvidia-modprobe  --version
    if [ $? -eq 0 ]; then
        echo "GPU driver is already installed, skipping download"
    else
        # First we remove the nouveau drivers, which are the open source drivers for NVIDIA cards. Nouveau is installed on NV Series VMs by default.
        # We also installed needed dependencies.
        rmmod nouveau
        echo blacklist nouveau >> /etc/modprobe.d/blacklist.conf
        update-initramfs -u
        mkdir -p $GPU_DEST
        cd $GPU_DEST
        # Installing nvidia-docker, setting nvidia runtime as default and restarting docker daemon
        retrycmd_if_failure_no_stats 180 1 5 curl -fsSL https://nvidia.github.io/nvidia-docker/gpgkey > /tmp/aptnvidia.gpg || exit $ERR_GPU_DRIVERS_INSTALL_TIMEOUT
        cat /tmp/aptnvidia.gpg | apt-key add -
        retrycmd_if_failure_no_stats 180 1 5 curl -fsSL https://nvidia.github.io/nvidia-docker/ubuntu16.04/amd64/nvidia-docker.list > /tmp/nvidia-docker.list || exit  $ERR_GPU_DRIVERS_INSTALL_TIMEOUT
        cat /tmp/nvidia-docker.list | sudo tee /etc/apt/sources.list.d/nvidia-docker.list
        apt_get_update
        retrycmd_if_failure 30 5 300 apt-get install -y linux-headers-$(uname -r) gcc make dkms || exit $ERR_GPU_DRIVERS_INSTALL_TIMEOUT
        retrycmd_if_failure 30 5 300 apt-get -o Dpkg::Options::="--force-confold" install -y nvidia-docker2=$NVIDIA_DOCKER_VERSION+docker$DOCKER_VERSION nvidia-container-runtime=$NVIDIA_CONTAINER_RUNTIME_VERSION+docker$DOCKER_VERSION || exit $ERR_GPU_DRIVERS_INSTALL_TIMEOUT
        pkill -SIGHUP dockerd
        mkdir -p $GPU_DEST
        cd $GPU_DEST
        # Download the .run file from NVIDIA.
        # Nvidia libraries are always install in /usr/lib/x86_64-linux-gnu, and there is no option in the run file to change this.
        # Instead we use Overlayfs to move the newly installed libraries under /usr/local/nvidia/lib64
        retrycmd_if_failure 5 10 60 curl -fLS https://us.download.nvidia.com/tesla/$GPU_DV/NVIDIA-Linux-x86_64-$GPU_DV.run -o nvidia-drivers-$GPU_DV || exit $ERR_GPU_DRIVERS_INSTALL_TIMEOUT
        mkdir -p lib64 overlay-workdir
        mount -t overlay -o lowerdir=/usr/lib/x86_64-linux-gnu,upperdir=lib64,workdir=overlay-workdir none /usr/lib/x86_64-linux-gnu
    fi
}

function installContainerRuntime() {
    if [[ "$CONTAINER_RUNTIME" == "docker" ]]; then
        installDocker
    elif [[ "$CONTAINER_RUNTIME" == "clear-containers" ]]; then
	    # Ensure we can nest virtualization
        if grep -q vmx /proc/cpuinfo; then
            installClearContainersRuntime
        fi
    fi
}

function installDocker() {
    CURRENT_VERSION=$(docker --version | cut -d " " -f 3 | cut -d "," -f 1)
    if [[ "$CURRENT_VERSION" = ${DOCKER_ENGINE_VERSION} ]]; then
        echo "docker version ${DOCKER_ENGINE_VERSION} is already installed, skipping download"
    else
        retrycmd_if_failure_no_stats 20 1 5 curl -fsSL https://aptdocker.azureedge.net/gpg > /tmp/aptdocker.gpg || exit $ERR_DOCKER_KEY_DOWNLOAD_TIMEOUT
        retrycmd_if_failure 10 5 10 apt-key add /tmp/aptdocker.gpg || exit $ERR_DOCKER_APT_KEY_TIMEOUT
        echo "deb ${DOCKER_REPO} ubuntu-xenial main" | sudo tee /etc/apt/sources.list.d/docker.list
        printf "Package: docker-engine\nPin: version ${DOCKER_ENGINE_VERSION}\nPin-Priority: 550\n" > /etc/apt/preferences.d/docker.pref
        apt_get_update || exit $ERR_APT_UPDATE_TIMEOUT
        apt_get_install 20 30 120 docker-engine || exit $ERR_DOCKER_INSTALL_TIMEOUT
    fi
}

function installKataContainersRuntime() {
    # TODO incorporate this into packer CI so that it is pre-baked into the VHD image
    # Add Kata Containers repository key
    echo "Adding Kata Containers repository key..."
    KATA_RELEASE_KEY_TMP=/tmp/kata-containers-release.key
    KATA_URL=http://download.opensuse.org/repositories/home:/katacontainers:/release/xUbuntu_16.04/Release.key
    retrycmd_if_failure_no_stats 20 1 5 curl -fsSL $KATA_URL > $KATA_RELEASE_KEY_TMP || exit $ERR_KATA_KEY_DOWNLOAD_TIMEOUT
    retrycmd_if_failure 10 5 10 apt-key add $KATA_RELEASE_KEY_TMP || exit $ERR_KATA_APT_KEY_TIMEOUT

    # Add Kata Container repository
    echo "Adding Kata Containers repository..."
    echo 'deb http://download.opensuse.org/repositories/home:/katacontainers:/release/xUbuntu_16.04/ /' > /etc/apt/sources.list.d/kata-containers.list

    # Install Kata Containers runtime
    echo "Installing Kata Containers runtime..."
    apt_get_update || exit $ERR_APT_UPDATE_TIMEOUT
    apt_get_install 20 30 120 kata-runtime || exit $ERR_KATA_INSTALL_TIMEOUT
}

function installClearContainersRuntime() {
    cc-runtime --version
    if [ $? -eq 0 ]; then
        echo "cc-runtime is already installed, skipping download"
    else
        # Add Clear Containers repository key
        echo "Adding Clear Containers repository key..."
        CC_RELEASE_KEY_TMP=/tmp/clear-containers-release.key
        CC_URL=https://download.opensuse.org/repositories/home:clearcontainers:clear-containers-3/xUbuntu_16.04/Release.key
        retrycmd_if_failure_no_stats 20 1 5 curl -fsSL $CC_URL > $CC_RELEASE_KEY_TMP || exit $ERR_APT_INSTALL_TIMEOUT
        retrycmd_if_failure 10 5 10 apt-key add $CC_RELEASE_KEY_TMP || exit $ERR_APT_INSTALL_TIMEOUT

        # Add Clear Container repository
        echo "Adding Clear Containers repository..."
        echo 'deb http://download.opensuse.org/repositories/home:/clearcontainers:/clear-containers-3/xUbuntu_16.04/ /' > /etc/apt/sources.list.d/cc-runtime.list

        # Install Clear Containers runtime
        echo "Installing Clear Containers runtime..."
        apt_get_update || exit $ERR_APT_UPDATE_TIMEOUT
        apt_get_install 20 30 120 cc-runtime

        # Install the systemd service and socket files.
        local repo_uri="https://raw.githubusercontent.com/clearcontainers/proxy/3.0.23"
        retrycmd_if_failure_no_stats 20 1 5 curl -fsSL "${repo_uri}/cc-proxy.service.in" > $CC_SERVICE_IN_TMP
        retrycmd_if_failure_no_stats 20 1 5 curl -fsSL "${repo_uri}/cc-proxy.socket.in" > $CC_SOCKET_IN_TMP
    fi
}

function installNetworkPlugin() {
    if [[ "${NETWORK_PLUGIN}" = "azure" ]]; then
        installAzureCNI
    fi
    installCNI
    rm -rf $CNI_DOWNLOADS_DIR
}

function downloadCNI() {
    mkdir -p $CNI_DOWNLOADS_DIR
    CNI_TGZ_TMP=$(echo ${CNI_PLUGINS_URL} | cut -d "/" -f 5)
    retrycmd_get_tarball 60 5 "$CNI_DOWNLOADS_DIR/${CNI_TGZ_TMP}" ${CNI_PLUGINS_URL} || exit $ERR_CNI_DOWNLOAD_TIMEOUT
}

function downloadAzureCNI() {
    mkdir -p $CNI_DOWNLOADS_DIR
    CNI_TGZ_TMP=$(echo ${VNET_CNI_PLUGINS_URL} | cut -d "/" -f 5)
    retrycmd_get_tarball 60 5 "$CNI_DOWNLOADS_DIR/${CNI_TGZ_TMP}" ${VNET_CNI_PLUGINS_URL} || exit $ERR_CNI_DOWNLOAD_TIMEOUT
}

function installCNI() {
    CNI_TGZ_TMP=$(echo ${CNI_PLUGINS_URL} | cut -d "/" -f 5)
    if [[ ! -f "$CNI_DOWNLOADS_DIR/${CNI_TGZ_TMP}" ]]; then
        downloadCNI
    fi
    mkdir -p $CNI_BIN_DIR
    tar -xzf "$CNI_DOWNLOADS_DIR/${CNI_TGZ_TMP}" -C $CNI_BIN_DIR
    chown -R root:root $CNI_BIN_DIR
    chmod -R 755 $CNI_BIN_DIR
}

function installAzureCNI() {
    CNI_TGZ_TMP=$(echo ${VNET_CNI_PLUGINS_URL} | cut -d "/" -f 5)
    if [[ ! -f "$CNI_DOWNLOADS_DIR/${CNI_TGZ_TMP}" ]]; then
        downloadAzureCNI
    fi
    mkdir -p $CNI_CONFIG_DIR
    chown -R root:root $CNI_CONFIG_DIR
    chmod 755 $CNI_CONFIG_DIR
    mkdir -p $CNI_BIN_DIR
    tar -xzf "$CNI_DOWNLOADS_DIR/${CNI_TGZ_TMP}" -C $CNI_BIN_DIR
}

function installContainerd() {
    containerd --version
    if [ $? -eq 0 ]; then
        echo "containerd is already installed, skipping download"
    else
        CRI_CONTAINERD_VERSION="1.1.0"
        CONTAINERD_DOWNLOAD_URL="${CONTAINERD_DOWNLOAD_URL_BASE}cri-containerd-${CRI_CONTAINERD_VERSION}.linux-amd64.tar.gz"
        CONTAINERD_TGZ_TMP=/tmp/containerd.tar.gz
        retrycmd_get_tarball 60 5 "$CONTAINERD_TGZ_TMP" "$CONTAINERD_DOWNLOAD_URL" || exit $ERR_CONTAINERD_DOWNLOAD_TIMEOUT
        tar -xzf "$CONTAINERD_TGZ_TMP" -C /
        rm -f "$CONTAINERD_TGZ_TMP"
        sed -i '/\[Service\]/a ExecStartPost=\/sbin\/iptables -P FORWARD ACCEPT' /etc/systemd/system/containerd.service
        echo "Successfully installed cri-containerd..."
    fi
}

function installImg() { 
    img_filepath=/usr/local/bin/img
    retrycmd_get_executable 20 5 $img_filepath "https://acs-mirror.azureedge.net/img/img-linux-amd64-v0.4.6" ls || exit $ERR_IMG_DOWNLOAD_TIMEOUT
}

function pullHyperkube() {
    retrycmd_if_failure 75 1 60 img pull $HYPERKUBE_URL || exit $ERR_K8S_DOWNLOAD_TIMEOUT
    img unpack -o "/home/rootfs-${KUBERNETES_VERSION}" $HYPERKUBE_URL
    path=$(find /home/rootfs-${KUBERNETES_VERSION} -name "hyperkube")

    if [[ $OS == $COREOS_OS_NAME ]]; then
        cp "$path" "/opt/kubelet"
        cp "$path" "/opt/kubectl"
        chmod a+x /opt/kubelet /opt/kubectl
    else
        cp "$path" "/usr/local/bin/kubelet-${KUBERNETES_VERSION}"
        cp "$path" "/usr/local/bin/kubectl-${KUBERNETES_VERSION}"
    fi
}

function extractHyperkube() {
    if [[ ! -f "/usr/local/bin/kubelet-${KUBERNETES_VERSION}" ]]; then
        installImg
        pullHyperkube
    fi
    mv "/usr/local/bin/kubelet-${KUBERNETES_VERSION}" "/usr/local/bin/kubelet"
    mv "/usr/local/bin/kubectl-${KUBERNETES_VERSION}" "/usr/local/bin/kubectl"
    chmod a+x /usr/local/bin/kubelet /usr/local/bin/kubectl
    rm -rf /usr/local/bin/kubelet-* /usr/local/bin/kubectl-*
}

function pullImg() {
    DOCKER_IMAGE_URL=$1
    retrycmd_if_failure 75 1 60 img pull $DOCKER_IMAGE_URL || exit $ERR_IMG_DOWNLOAD_TIMEOUT
}
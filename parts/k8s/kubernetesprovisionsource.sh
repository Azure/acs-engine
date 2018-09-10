#!/bin/sh

#ERR_SYSTEMCTL_ENABLE_FAIL=3 # Service could not be enabled by systemctl -- DEPRECATED
ERR_SYSTEMCTL_START_FAIL=4 # Service could not be started or enabled by systemctl
ERR_CLOUD_INIT_TIMEOUT=5 # Timeout waiting for cloud-init runcmd to complete
ERR_FILE_WATCH_TIMEOUT=6 # Timeout waiting for a file
ERR_HOLD_WALINUXAGENT=7 # Unable to place walinuxagent apt package on hold during install
ERR_RELEASE_HOLD_WALINUXAGENT=8 # Unable to release hold on walinuxagent apt package after install
ERR_APT_INSTALL_TIMEOUT=9 # Timeout installing required apt packages
ERR_ETCD_DATA_DIR_NOT_FOUND=10 # Etcd data dir not found
ERR_ETCD_RUNNING_TIMEOUT=11 # Timeout waiting for etcd to be accessible
ERR_ETCD_DOWNLOAD_TIMEOUT=12 # Timeout waiting for etcd to download
ERR_ETCD_VOL_MOUNT_FAIL=13 # Unable to mount etcd disk volume
ERR_ETCD_START_TIMEOUT=14 # Unable to start etcd runtime
ERR_ETCD_CONFIG_FAIL=15 # Unable to configure etcd cluster
ERR_DOCKER_INSTALL_TIMEOUT=20 # Timeout waiting for docker install
ERR_DOCKER_DOWNLOAD_TIMEOUT=21 # Timout waiting for docker download(s)
ERR_DOCKER_KEY_DOWNLOAD_TIMEOUT=22 # Timeout waiting to download docker repo key
ERR_DOCKER_APT_KEY_TIMEOUT=23 # Timeout waiting for docker apt-key
ERR_DOCKER_START_FAIL=24 # Docker could not be started by systemctl
ERR_K8S_RUNNING_TIMEOUT=30 # Timeout waiting for k8s cluster to be healthy
ERR_K8S_DOWNLOAD_TIMEOUT=31 # Timeout waiting for Kubernetes download(s)
ERR_KUBECTL_NOT_FOUND=32 # kubectl client binary not found on local disk
ERR_KUBELET_START_FAIL=34 # kubelet could not be started by systemctl
ERR_IMG_DOWNLOAD_TIMEOUT=33 # Timeout waiting for img download
ERR_CNI_DOWNLOAD_TIMEOUT=41 # Timeout waiting for CNI download(s)
ERR_MS_PROD_DEB_DOWNLOAD_TIMEOUT=42 # Timeout waiting for https://packages.microsoft.com/config/ubuntu/16.04/packages-microsoft-prod.deb
ERR_MS_PROD_DEB_PKG_ADD_FAIL=43 # Failed to add repo pkg file
#ERR_FLEXVOLUME_DOWNLOAD_TIMEOUT=44 # Failed to add repo pkg file -- DEPRECATED
ERR_MODPROBE_FAIL=49 # Unable to load a kernel module using modprobe
ERR_OUTBOUND_CONN_FAIL=50 # Unable to establish outbound connection
ERR_KATA_KEY_DOWNLOAD_TIMEOUT=60 # Timeout waiting to download kata repo key
ERR_KATA_APT_KEY_TIMEOUT=61 # Timeout waiting for kata apt-key
ERR_KATA_INSTALL_TIMEOUT=62 # Timeout waiting for kata install
ERR_CONTAINERD_DOWNLOAD_TIMEOUT=70 # Timeout waiting for containerd download(s)
ERR_CUSTOM_SEARCH_DOMAINS_FAIL=80 # Unable to configure custom search domains
ERR_GPU_DRIVERS_START_FAIL=84 # nvidia-modprobe could not be started by systemctl
ERR_GPU_DRIVERS_INSTALL_TIMEOUT=85 # Timeout waiting for GPU drivers install
ERR_APT_DAILY_TIMEOUT=98 # Timeout waiting for apt daily updates
ERR_APT_UPDATE_TIMEOUT=99 # Timeout waiting for apt-get update to complete
ERR_CSE_PROVISION_SCRIPT_NOT_READY_TIMEOUT=100 # Timeout waiting for cloud-init to place this (!) script on the vm

OS=$(cat /etc/*-release | grep ^ID= | tr -d 'ID="' | awk '{print toupper($0)}')
UBUNTU_OS_NAME="UBUNTU"
RHEL_OS_NAME="RHEL"
COREOS_OS_NAME="COREOS"
KUBECTL=/usr/local/bin/kubectl
DOCKER=/usr/bin/docker
# latest version of the GPU drivers. Later this parameter could be bubbled up so that users can choose specific driver versions.
GPU_DV=396.26
GPU_DEST=/usr/local/nvidia
NVIDIA_DOCKER_VERSION=2.0.3
DOCKER_VERSION=1.13.1-1
NVIDIA_CONTAINER_RUNTIME_VERSION=2.0.0

retrycmd_if_failure() {
    retries=$1; wait_sleep=$2; timeout=$3; shift && shift && shift
    for i in $(seq 1 $retries); do
        timeout $timeout ${@}
        [ $? -eq 0  ] && break || \
        if [ $i -eq $retries ]; then
            echo Executed \"$@\" $i times;
            return 1
        else
            sleep $wait_sleep
        fi
    done
    echo Executed \"$@\" $i times;
}
retrycmd_if_failure_no_stats() {
    retries=$1; wait_sleep=$2; timeout=$3; shift && shift && shift
    for i in $(seq 1 $retries); do
        timeout $timeout ${@}
        [ $? -eq 0  ] && break || \
        if [ $i -eq $retries ]; then
            return 1
        else
            sleep $wait_sleep
        fi
    done
}
retrycmd_get_tarball() {
    tar_retries=$1; wait_sleep=$2; tarball=$3; url=$4
    echo "${tar_retries} retries"
    for i in $(seq 1 $tar_retries); do
        tar -tzf $tarball
        [ $? -eq 0  ] && break || \
        if [ $i -eq $tar_retries ]; then
            return 1
        else
            timeout 60 curl -fsSL $url -o $tarball
            sleep $wait_sleep
        fi
    done
}
retrycmd_get_executable() {
    retries=$1; wait_sleep=$2; filepath=$3; url=$4; validation_args=$5
    echo "${retries} retries"
    for i in $(seq 1 $retries); do
        $filepath $validation_args
        [ $? -eq 0  ] && break || \
        if [ $i -eq $retries ]; then
            return 1
        else
            timeout 30 curl -fsSL $url -o $filepath
            chmod +x $filepath
            sleep $wait_sleep
        fi
    done
}
wait_for_file() {
    retries=$1; wait_sleep=$2; filepath=$3
    for i in $(seq 1 $retries); do
        if [ -f $filepath ]; then
            break
        fi
        if [ $i -eq $retries ]; then
            return 1
        else
            sleep $wait_sleep
        fi
    done
}
apt_get_update() {
    retries=10
    apt_update_output=/tmp/apt-get-update.out
    for i in $(seq 1 $retries); do
        timeout 30 dpkg --configure -a
        timeout 30 apt-get -f -y install
        timeout 120 apt-get update 2>&1 | tee $apt_update_output | grep -E "^([WE]:.*)|([eE]rr.*)$"
        [ $? -ne 0  ] && cat $apt_update_output && break || \
        cat $apt_update_output
        if [ $i -eq $retries ]; then
            return 1
        else sleep 30
        fi
    done
    echo Executed apt-get update $i times
}
apt_get_install() {
    retries=$1; wait_sleep=$2; timeout=$3; shift && shift && shift
    for i in $(seq 1 $retries); do
        timeout 30 dpkg --configure -a
        timeout $timeout apt-get install --no-install-recommends -y ${@}
        [ $? -eq 0  ] && break || \
        if [ $i -eq $retries ]; then
            return 1
        else
            sleep $wait_sleep
            apt_get_update
        fi
    done
    echo Executed apt-get install --no-install-recommends -y \"$@\" $i times;
}
systemctl_restart() {
    retries=$1; wait_sleep=$2; timeout=$3 svcname=$4
    for i in $(seq 1 $retries); do
        timeout $timeout systemctl daemon-reload
        timeout $timeout systemctl restart $svcname
        [ $? -eq 0  ] && break || \
        if [ $i -eq $retries ]; then
            return 1
        else
            sleep $wait_sleep
        fi
    done
}

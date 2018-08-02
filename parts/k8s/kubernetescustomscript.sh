#!/bin/bash
set -x
echo `date`,`hostname`, startscript>>/opt/m
source /opt/azure/containers/provision_source.sh
# TODO standardize/generalize CSE exit codes
ERR_SYSTEMCTL_ENABLE_FAIL=3 # Service could not be enabled by systemctl
ERR_SYSTEMCTL_START_FAIL=4 # Service could not be started by systemctl
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
ERR_K8S_RUNNING_TIMEOUT=30 # Timeout waiting for k8s cluster to be healthy
ERR_K8S_DOWNLOAD_TIMEOUT=31 # Timeout waiting for Kubernetes download(s)
ERR_KUBECTL_NOT_FOUND=32 # kubectl client binary not found on local disk
ERR_CNI_DOWNLOAD_TIMEOUT=41 # Timeout waiting for CNI download(s)
ERR_MS_PROD_DEB_DOWNLOAD_TIMEOUT=42 # Timeout waiting for https://packages.microsoft.com/config/ubuntu/16.04/packages-microsoft-prod.deb
ERR_MS_PROD_DEB_PKG_ADD_FAIL=43 # Failed to add repo pkg file
#ERR_FLEXVOLUME_DOWNLOAD_TIMEOUT=44 # Failed to add repo pkg file -- DEPRECATED
ERR_MODPROBE_FAIL=49 # Unable to load a kernel module using modprobe
ERR_OUTBOUND_CONN_FAIL=50 # Unable to establish outbound connection
ERR_KATA_KEY_DOWNLOAD_TIMEOUT=60 # Timeout waiting to download kata repo key
ERR_KATA_APT_KEY_TIMEOUT=61 # Timeout waiting for kata apt-key
ERR_KATA_INSTALL_TIMEOUT=62 # Timeout waiting for kata install
ERR_CUSTOM_SEARCH_DOMAINS_FAIL=80 # Unable to configure custom search domains
ERR_APT_DAILY_TIMEOUT=98 # Timeout waiting for apt daily updates
ERR_APT_UPDATE_TIMEOUT=99 # Timeout waiting for apt-get update to complete

OS=$(cat /etc/*-release | grep ^ID= | tr -d 'ID="' | awk '{print toupper($0)}')
UBUNTU_OS_NAME="UBUNTU"
RHEL_OS_NAME="RHEL"
COREOS_OS_NAME="COREOS"
KUBECTL=/usr/local/bin/kubectl
DOCKER=/usr/bin/docker
CNI_BIN_DIR=/opt/cni/bin
CUSTOM_SEARCH_DOMAIN_SCRIPT=/opt/azure/containers/setup-custom-search-domains.sh

set +x
ETCD_PEER_CERT=$(echo ${ETCD_PEER_CERTIFICATES} | cut -d'[' -f 2 | cut -d']' -f 1 | cut -d',' -f $((${MASTER_INDEX}+1)))
ETCD_PEER_KEY=$(echo ${ETCD_PEER_PRIVATE_KEYS} | cut -d'[' -f 2 | cut -d']' -f 1 | cut -d',' -f $((${MASTER_INDEX}+1)))
set -x

if [[ $OS == $COREOS_OS_NAME ]]; then
    echo "Changing default kubectl bin location"
    KUBECTL=/opt/kubectl
fi

if [ -f /var/run/reboot-required ]; then
    REBOOTREQUIRED=true
else
    REBOOTREQUIRED=false
fi

function testOutboundConnection() {
    retrycmd_if_failure 20 1 3 nc -v www.google.com 443 || retrycmd_if_failure 20 1 3 nc -v www.1688.com 443 || exit $ERR_OUTBOUND_CONN_FAIL
}

function waitForCloudInit() {
    wait_for_file 900 1 /var/log/azure/cloud-init.complete || exit $ERR_CLOUD_INIT_TIMEOUT
}

function systemctlEnableAndStart() {
    systemctl_restart 100 5 30 $1
    RESTART_STATUS=$?
    systemctl status $1 --no-pager -l > /var/log/azure/$1-status.log
    if [ $RESTART_STATUS -ne 0 ]; then
        echo "$1 could not be started"
        exit $ERR_SYSTEMCTL_START_FAIL
    fi
    retrycmd_if_failure 10 5 3 systemctl enable $1
    if [ $? -ne 0 ]; then
        echo "$1 could not be enabled by systemctl"
        exit $ERR_SYSTEMCTL_ENABLE_FAIL
    fi
}

function installEtcd() {
    useradd -U "etcd"
    usermod -p "$(head -c 32 /dev/urandom | base64)" "etcd"
    passwd -u "etcd"
    id "etcd"

    APISERVER_PRIVATE_KEY_PATH="/etc/kubernetes/certs/apiserver.key"
    touch "${APISERVER_PRIVATE_KEY_PATH}"
    chmod 0600 "${APISERVER_PRIVATE_KEY_PATH}"
    chown root:root "${APISERVER_PRIVATE_KEY_PATH}"

    CA_PRIVATE_KEY_PATH="/etc/kubernetes/certs/ca.key"
    touch "${CA_PRIVATE_KEY_PATH}"
    chmod 0600 "${CA_PRIVATE_KEY_PATH}"
    chown root:root "${CA_PRIVATE_KEY_PATH}"

    ETCD_SERVER_PRIVATE_KEY_PATH="/etc/kubernetes/certs/etcdserver.key"
    touch "${ETCD_SERVER_PRIVATE_KEY_PATH}"
    chmod 0600 "${ETCD_SERVER_PRIVATE_KEY_PATH}"
    chown etcd:etcd "${ETCD_SERVER_PRIVATE_KEY_PATH}"

    ETCD_CLIENT_PRIVATE_KEY_PATH="/etc/kubernetes/certs/etcdclient.key"
    touch "${ETCD_CLIENT_PRIVATE_KEY_PATH}"
    chmod 0600 "${ETCD_CLIENT_PRIVATE_KEY_PATH}"
    chown root:root "${ETCD_CLIENT_PRIVATE_KEY_PATH}"

    ETCD_PEER_PRIVATE_KEY_PATH="/etc/kubernetes/certs/etcdpeer${MASTER_INDEX}.key"
    touch "${ETCD_PEER_PRIVATE_KEY_PATH}"
    chmod 0600 "${ETCD_PEER_PRIVATE_KEY_PATH}"
    chown etcd:etcd "${ETCD_PEER_PRIVATE_KEY_PATH}"

    ETCD_SERVER_CERTIFICATE_PATH="/etc/kubernetes/certs/etcdserver.crt"
    touch "${ETCD_SERVER_CERTIFICATE_PATH}"
    chmod 0644 "${ETCD_SERVER_CERTIFICATE_PATH}"
    chown root:root "${ETCD_SERVER_CERTIFICATE_PATH}"

    ETCD_CLIENT_CERTIFICATE_PATH="/etc/kubernetes/certs/etcdclient.crt"
    touch "${ETCD_CLIENT_CERTIFICATE_PATH}"
    chmod 0644 "${ETCD_CLIENT_CERTIFICATE_PATH}"
    chown root:root "${ETCD_CLIENT_CERTIFICATE_PATH}"

    ETCD_PEER_CERTIFICATE_PATH="/etc/kubernetes/certs/etcdpeer${MASTER_INDEX}.crt"
    touch "${ETCD_PEER_CERTIFICATE_PATH}"
    chmod 0644 "${ETCD_PEER_CERTIFICATE_PATH}"
    chown root:root "${ETCD_PEER_CERTIFICATE_PATH}"

    set +x
    echo "${APISERVER_PRIVATE_KEY}" | base64 --decode > "${APISERVER_PRIVATE_KEY_PATH}"
    echo "${CA_PRIVATE_KEY}" | base64 --decode > "${CA_PRIVATE_KEY_PATH}"
    echo "${ETCD_SERVER_PRIVATE_KEY}" | base64 --decode > "${ETCD_SERVER_PRIVATE_KEY_PATH}"
    echo "${ETCD_CLIENT_PRIVATE_KEY}" | base64 --decode > "${ETCD_CLIENT_PRIVATE_KEY_PATH}"
    echo "${ETCD_PEER_KEY}" | base64 --decode > "${ETCD_PEER_PRIVATE_KEY_PATH}"
    echo "${ETCD_SERVER_CERTIFICATE}" | base64 --decode > "${ETCD_SERVER_CERTIFICATE_PATH}"
    echo "${ETCD_CLIENT_CERTIFICATE}" | base64 --decode > "${ETCD_CLIENT_CERTIFICATE_PATH}"
    echo "${ETCD_PEER_CERT}" | base64 --decode > "${ETCD_PEER_CERTIFICATE_PATH}"
    set -x

    /opt/azure/containers/setup-etcd.sh > /opt/azure/containers/setup-etcd.log 2>&1
    RET=$?
    if [ $RET -ne 0 ]; then
        exit $RET
    fi

    /opt/azure/containers/mountetcd.sh || exit $ERR_ETCD_VOL_MOUNT_FAIL
    systemctlEnableAndStart etcd
    for i in $(seq 1 600); do
        MEMBER="$(sudo etcdctl member list | grep -E ${MASTER_VM_NAME} | cut -d':' -f 1)"
        if [ "$MEMBER" != "" ]; then
            break
        else
            sleep 1
        fi
    done
    retrycmd_if_failure 10 1 5 sudo etcdctl member update $MEMBER ${ETCD_PEER_URL} || exit $ERR_ETCD_CONFIG_FAIL
}

function installDeps() {
    retrycmd_if_failure_no_stats 20 1 5 curl -fsSL https://packages.microsoft.com/config/ubuntu/16.04/packages-microsoft-prod.deb > /tmp/packages-microsoft-prod.deb || exit $ERR_MS_PROD_DEB_DOWNLOAD_TIMEOUT
    retrycmd_if_failure 60 5 10 dpkg -i /tmp/packages-microsoft-prod.deb || exit $ERR_MS_PROD_DEB_PKG_ADD_FAIL
    echo `date`,`hostname`, apt-get_update_begin>>/opt/m
    apt_get_update || exit $ERR_APT_UPDATE_TIMEOUT
    echo `date`,`hostname`, apt-get_update_end>>/opt/m
    # make sure walinuxagent doesn't get updated in the middle of running this script
    retrycmd_if_failure 20 5 30 apt-mark hold walinuxagent || exit $ERR_HOLD_WALINUXAGENT
    # See https://github.com/kubernetes/kubernetes/blob/master/build/debian-hyperkube-base/Dockerfile#L25-L44
    apt_get_install 20 30 300 apt-transport-https ca-certificates iptables iproute2 ebtables socat util-linux mount ethtool init-system-helpers nfs-common ceph-common conntrack glusterfs-client ipset jq cgroup-lite git pigz xz-utils blobfuse fuse cifs-utils || exit $ERR_APT_INSTALL_TIMEOUT
    systemctlEnableAndStart rpcbind
    systemctlEnableAndStart rpc-statd
}

function installDocker() {
    retrycmd_if_failure_no_stats 20 1 5 curl -fsSL https://aptdocker.azureedge.net/gpg > /tmp/aptdocker.gpg || exit $ERR_DOCKER_KEY_DOWNLOAD_TIMEOUT
    retrycmd_if_failure 10 5 10 apt-key add /tmp/aptdocker.gpg || exit $ERR_DOCKER_APT_KEY_TIMEOUT
    echo "deb ${DOCKER_REPO} ubuntu-xenial main" | sudo tee /etc/apt/sources.list.d/docker.list
    printf "Package: docker-engine\nPin: version ${DOCKER_ENGINE_VERSION}\nPin-Priority: 550\n" > /etc/apt/preferences.d/docker.pref
    apt_get_update || exit $ERR_APT_UPDATE_TIMEOUT
    apt_get_install 20 30 120 docker-engine || exit $ERR_DOCKER_INSTALL_TIMEOUT
    echo "ExecStartPost=/sbin/iptables -P FORWARD ACCEPT" >> /etc/systemd/system/docker.service.d/exec_start.conf
    usermod -aG docker ${ADMINUSER}
    touch /var/log/azure/docker-install.complete
}

function runAptDaily() {
    /usr/lib/apt/apt.systemd.daily
}

function generateAggregatedAPICerts() {
    wait_for_file 1 1 /etc/kubernetes/generate-proxy-certs.sh && /etc/kubernetes/generate-proxy-certs.sh
}

function configureK8s() {
    KUBELET_PRIVATE_KEY_PATH="/etc/kubernetes/certs/client.key"
    touch "${KUBELET_PRIVATE_KEY_PATH}"
    chmod 0600 "${KUBELET_PRIVATE_KEY_PATH}"
    chown root:root "${KUBELET_PRIVATE_KEY_PATH}"

    APISERVER_PUBLIC_KEY_PATH="/etc/kubernetes/certs/apiserver.crt"
    touch "${APISERVER_PUBLIC_KEY_PATH}"
    chmod 0644 "${APISERVER_PUBLIC_KEY_PATH}"
    chown root:root "${APISERVER_PUBLIC_KEY_PATH}"

    AZURE_JSON_PATH="/etc/kubernetes/azure.json"
    touch "${AZURE_JSON_PATH}"
    chmod 0600 "${AZURE_JSON_PATH}"
    chown root:root "${AZURE_JSON_PATH}"

    set +x
    echo "${KUBELET_PRIVATE_KEY}" | base64 --decode > "${KUBELET_PRIVATE_KEY_PATH}"
    echo "${APISERVER_PUBLIC_KEY}" | base64 --decode > "${APISERVER_PUBLIC_KEY_PATH}"
    cat << EOF > "${AZURE_JSON_PATH}"
{
    "cloud":"${TARGET_ENVIRONMENT}",
    "tenantId": "${TENANT_ID}",
    "subscriptionId": "${SUBSCRIPTION_ID}",
    "aadClientId": "${SERVICE_PRINCIPAL_CLIENT_ID}",
    "aadClientSecret": "${SERVICE_PRINCIPAL_CLIENT_SECRET}",
    "resourceGroup": "${RESOURCE_GROUP}",
    "location": "${LOCATION}",
    "vmType": "${VM_TYPE}",
    "subnetName": "${SUBNET}",
    "securityGroupName": "${NETWORK_SECURITY_GROUP}",
    "vnetName": "${VIRTUAL_NETWORK}",
    "vnetResourceGroup": "${VIRTUAL_NETWORK_RESOURCE_GROUP}",
    "routeTableName": "${ROUTE_TABLE}",
    "primaryAvailabilitySetName": "${PRIMARY_AVAILABILITY_SET}",
    "primaryScaleSetName": "${PRIMARY_SCALE_SET}",
    "cloudProviderBackoff": ${CLOUDPROVIDER_BACKOFF},
    "cloudProviderBackoffRetries": ${CLOUDPROVIDER_BACKOFF_RETRIES},
    "cloudProviderBackoffExponent": ${CLOUDPROVIDER_BACKOFF_EXPONENT},
    "cloudProviderBackoffDuration": ${CLOUDPROVIDER_BACKOFF_DURATION},
    "cloudProviderBackoffJitter": ${CLOUDPROVIDER_BACKOFF_JITTER},
    "cloudProviderRatelimit": ${CLOUDPROVIDER_RATELIMIT},
    "cloudProviderRateLimitQPS": ${CLOUDPROVIDER_RATELIMIT_QPS},
    "cloudProviderRateLimitBucket": ${CLOUDPROVIDER_RATELIMIT_BUCKET},
    "useManagedIdentityExtension": ${USE_MANAGED_IDENTITY_EXTENSION},
    "useInstanceMetadata": ${USE_INSTANCE_METADATA},
    "providerVaultName": "${KMS_PROVIDER_VAULT_NAME}",
    "providerKeyName": "k8s",
    "providerKeyVersion": ""
}
EOF
    set -x
    generateAggregatedAPICerts
}

function setKubeletOpts () {
	sed -i "s#^KUBELET_OPTS=.*#KUBELET_OPTS=${1}#" /etc/default/kubelet
}

function installCNI() {
    mkdir -p $CNI_BIN_DIR
    CONTAINERNETWORKING_CNI_TGZ_TMP=/tmp/containernetworking_cni.tgz
    retrycmd_get_tarball 60 5 $CONTAINERNETWORKING_CNI_TGZ_TMP ${CNI_PLUGINS_URL} || exit $ERR_CNI_DOWNLOAD_TIMEOUT
    tar -xzf $CONTAINERNETWORKING_CNI_TGZ_TMP -C $CNI_BIN_DIR
    chown -R root:root $CNI_BIN_DIR
    chmod -R 755 $CNI_BIN_DIR
    # Turn on br_netfilter (needed for the iptables rules to work on bridges)
    # and permanently enable it
    retrycmd_if_failure 30 6 10 modprobe br_netfilter || exit $ERR_MODPROBE_FAIL
    # /etc/modules-load.d is the location used by systemd to load modules
    echo -n "br_netfilter" > /etc/modules-load.d/br_netfilter.conf
}

function configAzureCNI() {
    CNI_CONFIG_DIR=/etc/cni/net.d
    mkdir -p $CNI_CONFIG_DIR
    chown -R root:root $CNI_CONFIG_DIR
    chmod 755 $CNI_CONFIG_DIR
    mkdir -p $CNI_BIN_DIR
    AZURE_CNI_TGZ_TMP=/tmp/azure_cni.tgz
    retrycmd_get_tarball 60 5 $AZURE_CNI_TGZ_TMP ${VNET_CNI_PLUGINS_URL} || exit $ERR_CNI_DOWNLOAD_TIMEOUT
    tar -xzf $AZURE_CNI_TGZ_TMP -C $CNI_BIN_DIR
    installCNI
    mv $CNI_BIN_DIR/10-azure.conflist $CNI_CONFIG_DIR/
    chmod 600 $CNI_CONFIG_DIR/10-azure.conflist
    /sbin/ebtables -t nat --list
}

function configNetworkPlugin() {
    if [[ "${NETWORK_PLUGIN}" = "azure" ]]; then
        configAzureCNI
    elif [[ "${NETWORK_PLUGIN}" = "kubenet" ]]; then
		installCNI
	elif [[ "${NETWORK_PLUGIN}" = "flannel" ]]; then
        installCNI
    fi
}

function installKataContainersRuntime() {
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
    apt_get_update
    apt_get_install 20 30 120 cc-runtime

	# Install the systemd service and socket files.
	local repo_uri="https://raw.githubusercontent.com/clearcontainers/proxy/3.0.23"
    CC_SERVICE_IN_TMP=/tmp/cc-proxy.service.in
    CC_SOCKET_IN_TMP=/tmp/cc-proxy.socket.in
    retrycmd_if_failure_no_stats 20 1 5 curl -fsSL "${repo_uri}/cc-proxy.service.in" > $CC_SERVICE_IN_TMP
    retrycmd_if_failure_no_stats 20 1 5 curl -fsSL "${repo_uri}/cc-proxy.socket.in" > $CC_SOCKET_IN_TMP
    cat $CC_SERVICE_IN_TMP | sed 's#@libexecdir@#/usr/libexec#' > /etc/systemd/system/cc-proxy.service
    cat $CC_SOCKET_IN_TMP sed 's#@localstatedir@#/var#' > /etc/systemd/system/cc-proxy.socket

	# Enable and start Clear Containers proxy service
	echo "Enabling and starting Clear Containers proxy service..."
	systemctlEnableAndStart cc-proxy
}

function setupContainerd() {
	echo "Configuring cri-containerd..."

	mkdir -p "/etc/containerd"
	CRI_CONTAINERD_CONFIG="/etc/containerd/config.toml"
	echo "subreaper = false" > "$CRI_CONTAINERD_CONFIG"
	echo "oom_score = 0" >> "$CRI_CONTAINERD_CONFIG"

    echo "[plugins.cri]" >> "$CRI_CONTAINERD_CONFIG"
    echo "sandbox_image = \"$POD_INFRA_CONTAINER_SPEC\"" >> "$CRI_CONTAINERD_CONFIG"

	echo "[plugins.cri.containerd.untrusted_workload_runtime]" >> "$CRI_CONTAINERD_CONFIG"
	echo "runtime_type = 'io.containerd.runtime.v1.linux'" >> "$CRI_CONTAINERD_CONFIG"
	if [[ "$CONTAINER_RUNTIME" == "clear-containers" ]]; then
		echo "runtime_engine = '/usr/bin/cc-runtime'" >> "$CRI_CONTAINERD_CONFIG"
	elif [[ "$CONTAINER_RUNTIME" == "kata-containers" ]]; then
		echo "runtime_engine = '/usr/bin/kata-runtime'" >> "$CRI_CONTAINERD_CONFIG"
	else
		echo "runtime_engine = '/usr/local/sbin/runc'" >> "$CRI_CONTAINERD_CONFIG"
	fi
	echo "[plugins.cri.containerd.default_runtime]" >> "$CRI_CONTAINERD_CONFIG"
	echo "runtime_type = 'io.containerd.runtime.v1.linux'" >> "$CRI_CONTAINERD_CONFIG"
	echo "runtime_engine = '/usr/local/sbin/runc'" >> "$CRI_CONTAINERD_CONFIG"

	setKubeletOpts " --container-runtime=remote --runtime-request-timeout=15m --container-runtime-endpoint=unix:///run/containerd/containerd.sock"
}

function installContainerd() {
	CRI_CONTAINERD_VERSION="1.1.0"
	CONTAINERD_DOWNLOAD_URL="${CONTAINERD_DOWNLOAD_URL_BASE}cri-containerd-${CRI_CONTAINERD_VERSION}.linux-amd64.tar.gz"

    CONTAINERD_TGZ_TMP=/tmp/containerd.tar.gz
    retrycmd_get_tarball 60 5 "$CONTAINERD_TGZ_TMP" "$CONTAINERD_DOWNLOAD_URL"
	tar -xzf "$CONTAINERD_TGZ_TMP" -C /
	rm -f "$CONTAINERD_TGZ_TMP"
	sed -i '/\[Service\]/a ExecStartPost=\/sbin\/iptables -P FORWARD ACCEPT' /etc/systemd/system/containerd.service

	echo "Successfully installed cri-containerd..."
	if [[ "$CONTAINER_RUNTIME" == "clear-containers" ]] || [[ "$CONTAINER_RUNTIME" == "kata-containers" ]] || [[ "$CONTAINER_RUNTIME" == "containerd" ]]; then
		setupContainerd
	fi
}

function ensureContainerd() {
	if [[ "$CONTAINER_RUNTIME" == "clear-containers" ]] || [[ "$CONTAINER_RUNTIME" == "kata-containers" ]] || [[ "$CONTAINER_RUNTIME" == "containerd" ]]; then
		# Enable and start cri-containerd service
		# Make sure this is done after networking plugins are installed
		echo "Enabling and starting cri-containerd service..."
		systemctlEnableAndStart containerd
	fi
}

function ensureDocker() {
    wait_for_file 600 1 $DOCKER || exit $ERR_FILE_WATCH_TIMEOUT
    systemctlEnableAndStart docker
    retrycmd_if_failure 6 1 10 docker pull busybox # pre-pull busybox, but don't exit if fail
    systemctlEnableAndStart docker-health-probe
}
function ensureKMS() {
    systemctlEnableAndStart kms
}

function ensureKubelet() {
    systemctlEnableAndStart kubelet
}

function extractHyperkube(){
    TMP_DIR=$(mktemp -d)
    retrycmd_if_failure 100 1 30 curl -sSL -o /usr/local/bin/img "https://acs-mirror.azureedge.net/img/img-linux-amd64-v0.4.6"
    chmod +x /usr/local/bin/img
    retrycmd_if_failure 75 1 60 img pull $HYPERKUBE_URL || exit $ERR_K8S_DOWNLOAD_TIMEOUT
    path=$(find /tmp/img -name "hyperkube")

    if [[ $OS == $COREOS_OS_NAME ]]; then
        cp "$path" "/opt/kubelet"
        cp "$path" "/opt/kubectl"
        chmod a+x /opt/kubelet /opt/kubectl
    else
        cp "$path" "/usr/local/bin/kubelet"
        cp "$path" "/usr/local/bin/kubectl"
        chmod a+x /usr/local/bin/kubelet /usr/local/bin/kubectl
    fi
    rm -rf /tmp/hyperkube.tar "/tmp/img"
}

function ensureJournal(){
    echo "Storage=persistent" >> /etc/systemd/journald.conf
    echo "SystemMaxUse=1G" >> /etc/systemd/journald.conf
    echo "RuntimeMaxUse=1G" >> /etc/systemd/journald.conf
    echo "ForwardToSyslog=no" >> /etc/systemd/journald.conf
    systemctlEnableAndStart systemd-journald
}

function ensurePodSecurityPolicy() {
    POD_SECURITY_POLICY_FILE="/etc/kubernetes/manifests/pod-security-policy.yaml"
    if [ -f $POD_SECURITY_POLICY_FILE ]; then
        $KUBECTL create -f $POD_SECURITY_POLICY_FILE
    fi
}

function ensureK8sControlPlane() {
    if $REBOOTREQUIRED; then
        return
    fi
    wait_for_file 600 1 $KUBECTL || exit $ERR_FILE_WATCH_TIMEOUT
    retrycmd_if_failure 600 1 20 $KUBECTL 2>/dev/null cluster-info || exit $ERR_K8S_RUNNING_TIMEOUT
    ensurePodSecurityPolicy
}

function ensureEtcd() {
    retrycmd_if_failure 120 5 10 curl --cacert /etc/kubernetes/certs/ca.crt --cert /etc/kubernetes/certs/etcdclient.crt --key /etc/kubernetes/certs/etcdclient.key ${ETCD_CLIENT_URL}/v2/machines || exit $ERR_ETCD_RUNNING_TIMEOUT
}

function writeKubeConfig() {
    KUBECONFIGDIR=/home/$ADMINUSER/.kube
    KUBECONFIGFILE=$KUBECONFIGDIR/config
    mkdir -p $KUBECONFIGDIR
    touch $KUBECONFIGFILE
    chown $ADMINUSER:$ADMINUSER $KUBECONFIGDIR
    chown $ADMINUSER:$ADMINUSER $KUBECONFIGFILE
    chmod 700 $KUBECONFIGDIR
    chmod 600 $KUBECONFIGFILE

    # disable logging after secret output
    set +x
    echo "
---
apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: \"$CA_CERTIFICATE\"
    server: $KUBECONFIG_SERVER
  name: \"$MASTER_FQDN\"
contexts:
- context:
    cluster: \"$MASTER_FQDN\"
    user: \"$MASTER_FQDN-admin\"
  name: \"$MASTER_FQDN\"
current-context: \"$MASTER_FQDN\"
kind: Config
users:
- name: \"$MASTER_FQDN-admin\"
  user:
    client-certificate-data: \"$KUBECONFIG_CERTIFICATE\"
    client-key-data: \"$KUBECONFIG_KEY\"
" > $KUBECONFIGFILE
    # renable logging after secrets
    set -x
}

function configClusterAutoscalerAddon() {
    if [[ "${USE_MANAGED_IDENTITY_EXTENSION}" == true ]]; then
        CLUSTER_AUTOSCALER_MSI_VOLUME_MOUNT="- mountPath: /var/lib/waagent/\n\          name: waagent\n\          readOnly: true"
        CLUSTER_AUTOSCALER_MSI_VOLUME="- hostPath:\n\          path: /var/lib/waagent/\n\        name: waagent"
        CLUSTER_AUTOSCALER_MSI_HOST_NETWORK="hostNetwork: true"

        sed -i "s|<kubernetesClusterAutoscalerVolumeMounts>|${CLUSTER_AUTOSCALER_MSI_VOLUME_MOUNT}|g" "/etc/kubernetes/addons/cluster-autoscaler-deployment.yaml"
        sed -i "s|<kubernetesClusterAutoscalerVolumes>|${CLUSTER_AUTOSCALER_MSI_VOLUME}|g" "/etc/kubernetes/addons/cluster-autoscaler-deployment.yaml"
        sed -i "s|<kubernetesClusterAutoscalerHostNetwork>|$(echo "${CLUSTER_AUTOSCALER_MSI_HOST_NETWORK}")|g" "/etc/kubernetes/addons/cluster-autoscaler-deployment.yaml"
    elif [[ "${USE_MANAGED_IDENTITY_EXTENSION}" == false ]]; then
        sed -i "s|<kubernetesClusterAutoscalerVolumeMounts>|""|g" "/etc/kubernetes/addons/cluster-autoscaler-deployment.yaml"
        sed -i "s|<kubernetesClusterAutoscalerVolumes>|""|g" "/etc/kubernetes/addons/cluster-autoscaler-deployment.yaml"
        sed -i "s|<kubernetesClusterAutoscalerHostNetwork>|""|g" "/etc/kubernetes/addons/cluster-autoscaler-deployment.yaml"
    fi

    sed -i "s|<kubernetesClusterAutoscalerClientId>|$(echo $SERVICE_PRINCIPAL_CLIENT_ID | base64)|g" "/etc/kubernetes/addons/cluster-autoscaler-deployment.yaml"
    sed -i "s|<kubernetesClusterAutoscalerClientSecret>|$(echo $SERVICE_PRINCIPAL_CLIENT_SECRET | base64)|g" "/etc/kubernetes/addons/cluster-autoscaler-deployment.yaml"
    sed -i "s|<kubernetesClusterAutoscalerSubscriptionId>|$(echo $SUBSCRIPTION_ID | base64)|g" "/etc/kubernetes/addons/cluster-autoscaler-deployment.yaml"
    sed -i "s|<kubernetesClusterAutoscalerTenantId>|$(echo $TENANT_ID | base64)|g" "/etc/kubernetes/addons/cluster-autoscaler-deployment.yaml"
    sed -i "s|<kubernetesClusterAutoscalerResourceGroup>|$(echo $RESOURCE_GROUP | base64)|g" "/etc/kubernetes/addons/cluster-autoscaler-deployment.yaml"
    sed -i "s|<kubernetesClusterAutoscalerVmType>|$(echo $VM_TYPE | base64)|g" "/etc/kubernetes/addons/cluster-autoscaler-deployment.yaml"
    sed -i "s|<kubernetesClusterAutoscalerVMSSName>|$(echo $PRIMARY_SCALE_SET)|g" "/etc/kubernetes/addons/cluster-autoscaler-deployment.yaml"
}

configACIConnectorAddon() {
    ACI_CONNECTOR_CREDENTIALS=$(printf "{\"clientId\": \"$(echo $SERVICE_PRINCIPAL_CLIENT_ID)\", \"clientSecret\": \"$(echo $SERVICE_PRINCIPAL_CLIENT_SECRET)\", \"tenantId\": \"$(echo $TENANT_ID)\", \"subscriptionId\": \"$(echo $SUBSCRIPTION_ID)\", \"activeDirectoryEndpointUrl\": \"https://login.microsoftonline.com\",\"resourceManagerEndpointUrl\": \"https://management.azure.com/\", \"activeDirectoryGraphResourceId\": \"https://graph.windows.net/\", \"sqlManagementEndpointUrl\": \"https://management.core.windows.net:8443/\", \"galleryEndpointUrl\": \"https://gallery.azure.com/\", \"managementEndpointUrl\": \"https://management.core.windows.net/\"}" | base64 -w 0)

    openssl req -newkey rsa:4096 -new -nodes -x509 -days 3650 -keyout /etc/kubernetes/certs/aci-connector-key.pem -out /etc/kubernetes/certs/aci-connector-cert.pem -subj "/C=US/ST=CA/L=virtualkubelet/O=virtualkubelet/OU=virtualkubelet/CN=virtualkubelet"
    ACI_CONNECTOR_KEY=$(base64 /etc/kubernetes/certs/aci-connector-key.pem -w0)
    ACI_CONNECTOR_CERT=$(base64 /etc/kubernetes/certs/aci-connector-cert.pem -w0)

    sed -i "s|<kubernetesACIConnectorCredentials>|$ACI_CONNECTOR_CREDENTIALS|g" "/etc/kubernetes/addons/aci-connector-deployment.yaml"
    sed -i "s|<kubernetesACIConnectorResourceGroup>|$(echo $RESOURCE_GROUP)|g" "/etc/kubernetes/addons/aci-connector-deployment.yaml"
    sed -i "s|<kubernetesACIConnectorCert>|$(echo $ACI_CONNECTOR_CERT)|g" "/etc/kubernetes/addons/aci-connector-deployment.yaml"
    sed -i "s|<kubernetesACIConnectorKey>|$(echo $ACI_CONNECTOR_KEY)|g" "/etc/kubernetes/addons/aci-connector-deployment.yaml"
}

configAddons() {
    if [[ "${CLUSTER_AUTOSCALER_ADDON}" = True ]]; then
        configClusterAutoscalerAddon
    fi

    if [[ "${ACI_CONNECTOR_ADDON}" = True ]]; then
        configACIConnectorAddon
    fi
}

testOutboundConnection
waitForCloudInit

if [[ ! -z "${MASTER_NODE}" ]]; then
    echo "executing master node provision operations"
    installEtcd
else
    echo "skipping master node provision operations, this is an agent node"
fi

if [ -f $CUSTOM_SEARCH_DOMAIN_SCRIPT ]; then
    $CUSTOM_SEARCH_DOMAIN_SCRIPT > /opt/azure/containers/setup-custom-search-domain.log 2>&1 || exit $ERR_CUSTOM_SEARCH_DOMAINS_FAIL
fi

installDeps

if [[ "$CONTAINER_RUNTIME" == "docker" ]]; then
    installDocker
    ensureDocker
fi

configureK8s

configNetworkPlugin

if [[ ! -z "${MASTER_NODE}" ]]; then
    echo `date`,`hostname`, configAddonsStart>>/opt/m
    configAddons
    echo `date`,`hostname`, configAddonsDone>>/opt/m
fi

# containerd needs to be installed before extractHyperkube
# so runc is present.
echo `date`,`hostname`, installContainerdStart>>/opt/m
installContainerd
echo `date`,`hostname`, extractHyperkubeStart>>/opt/m
extractHyperkube
echo `date`,`hostname`, extractHyperkubeDone>>/opt/m

if [[ "$CONTAINER_RUNTIME" == "clear-containers" ]]; then
	# Ensure we can nest virtualization
	if grep -q vmx /proc/cpuinfo; then
		installClearContainersRuntime
	fi
fi

if [[ "$CONTAINER_RUNTIME" == "kata-containers" ]]; then
	# Ensure we can nest virtualization
	if grep -q vmx /proc/cpuinfo; then
		installKataContainersRuntime
	fi
fi

echo `date`,`hostname`, ensureContainerdStart>>/opt/m
ensureContainerd

if [[ ! -z "${MASTER_NODE}" && ! -z "${EnableEncryptionWithExternalKms}" ]]; then
    ensureKMS
fi

ensureKubelet
ensureJournal

if [[ ! -z "${MASTER_NODE}" ]]; then
    writeKubeConfig
    ensureEtcd
    ensureK8sControlPlane
fi

if [[ $OS == $UBUNTU_OS_NAME ]]; then
    # mitigation for bug https://bugs.launchpad.net/ubuntu/+source/linux/+bug/1676635
    echo 2dd1ce17-079e-403c-b352-a1921ee207ee > /sys/bus/vmbus/drivers/hv_util/unbind
    sed -i "13i\echo 2dd1ce17-079e-403c-b352-a1921ee207ee > /sys/bus/vmbus/drivers/hv_util/unbind\n" /etc/rc.local

    retrycmd_if_failure 20 5 30 apt-mark unhold walinuxagent || exit $ERR_RELEASE_HOLD_WALINUXAGENTs
fi

echo "Install complete successfully"

mkdir -p /opt/azure/containers && touch /opt/azure/containers/provision.complete
ps auxfww > /opt/azure/provision-ps.log &

if $REBOOTREQUIRED; then
  # wait 1 minute to restart node, so that the custom script extension can complete
  echo 'reboot required, rebooting node in 1 minute'
  /bin/bash -c "shutdown -r 1 &"
else
  runAptDaily &
fi

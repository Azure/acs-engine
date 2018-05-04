#!/bin/bash

# This script runs on every Kubernetes VM
# Exit codes represent the following:
# | exit code number | meaning |
# | 20 | Timeout waiting for docker install to finish |
# | 3 | Service could not be enabled by systemctl |
# | 4 | Service could not be started by systemctl |
# | 5 | Timeout waiting for cloud-init runcmd to complete |
# | 6 | Timeout waiting for a file |
# | 10 | Etcd data dir not found |
# | 11 | Timeout waiting for etcd to be accessible |
# | 30 | Timeout waiting for k8s cluster to be healthy|

set -x
source /opt/azure/containers/provision_source.sh

OS=$(cat /etc/*-release | grep ^ID= | tr -d 'ID="' | awk '{print toupper($0)}')
UBUNTU_OS_NAME="UBUNTU"
RHEL_OS_NAME="RHEL"
COREOS_OS_NAME="COREOS"
KUBECTL=/usr/local/bin/kubectl
DOCKER=/usr/bin/docker

set +x
ETCD_PEER_CERT=$(echo ${ETCD_PEER_CERTIFICATES} | cut -d'[' -f 2 | cut -d']' -f 1 | cut -d',' -f $((${MASTER_INDEX}+1)))
ETCD_PEER_KEY=$(echo ${ETCD_PEER_PRIVATE_KEYS} | cut -d'[' -f 2 | cut -d']' -f 1 | cut -d',' -f $((${MASTER_INDEX}+1)))
set -x

if [[ $OS == $COREOS_OS_NAME ]]; then
    echo "Changing default kubectl bin location"
    KUBECTL=/opt/kubectl
fi

ensureRunCommandCompleted()
{
    echo "waiting for runcmd to finish"
    wait_for_file 900 1 /opt/azure/containers/runcmd.complete
    if [ ! -f /opt/azure/containers/runcmd.complete ]; then
        echo "Timeout waiting for cloud-init runcmd to complete"
        exit 5
    fi
}

ensureDockerInstallCompleted()
{
    echo "waiting for docker install to finish"
    wait_for_file 3600 1 /opt/azure/containers/dockerinstall.complete
    if [ ! -f /opt/azure/containers/dockerinstall.complete ]; then
        echo "Timeout waiting for docker install to finish"
        exit 20
    fi
}

echo `date`,`hostname`, startscript>>/opt/m

if [ -f /var/run/reboot-required ]; then
    REBOOTREQUIRED=true
else
    REBOOTREQUIRED=false
fi

if [[ ! -z "${MASTER_NODE}" ]]; then
    echo "executing master node provision operations"

    useradd -U "etcd"
    usermod -p "$(head -c 32 /dev/urandom | base64)" "etcd"
    passwd -u "etcd"
    id "etcd"

    echo `date`,`hostname`, beginGettingEtcdCerts>>/opt/m
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

    echo `date`,`hostname`, endGettingEtcdCerts>>/opt/m
    mkdir -p /opt/azure/containers && touch /opt/azure/containers/certs.ready
else
    echo "skipping master node provision operations, this is an agent node"
fi

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

function ensureFilepath() {
    if $REBOOTREQUIRED; then
        return
    fi
    wait_for_file 600 1 $1
    if [ ! -f $1 ]; then
        echo "Timeout waiting for $1"
        exit 6
    fi
    
}

function setKubeletOpts () {
	sed -i "s#^KUBELET_OPTS=.*#KUBELET_OPTS=${1}#" /etc/default/kubelet
}

function installCNI() {
    CNI_BIN_DIR=/opt/cni/bin
    mkdir -p $CNI_BIN_DIR
    CONTAINERNETWORKING_CNI_TGZ_TMP=/tmp/containernetworking_cni.tgz
    retrycmd_get_tarball 60 1 $CONTAINERNETWORKING_CNI_TGZ_TMP ${CNI_PLUGINS_URL}
    tar -xzf $CONTAINERNETWORKING_CNI_TGZ_TMP -C $CNI_BIN_DIR
    chown -R root:root $CNI_BIN_DIR
    chmod -R 755 $CNI_BIN_DIR
}

function configAzureCNI() {
    CNI_CONFIG_DIR=/etc/cni/net.d
    mkdir -p $CNI_CONFIG_DIR
    chown -R root:root $CNI_CONFIG_DIR
    chmod 755 $CNI_CONFIG_DIR
    CNI_BIN_DIR=/opt/cni/bin
    mkdir -p $CNI_BIN_DIR
    AZURE_CNI_TGZ_TMP=/tmp/azure_cni.tgz
    retrycmd_get_tarball 60 1 $AZURE_CNI_TGZ_TMP ${VNET_CNI_PLUGINS_URL}
    tar -xzf $AZURE_CNI_TGZ_TMP -C $CNI_BIN_DIR
    installCNI
    mv $CNI_BIN_DIR/10-azure.conflist $CNI_CONFIG_DIR/
    chmod 600 $CNI_CONFIG_DIR/10-azure.conflist
    /sbin/ebtables -t nat --list
}

function configKubenet() {
    installCNI
}

function configFlannel() {
    installCNI
    setDockerOpts " --volume=/etc/cni/:/etc/cni:ro --volume=/opt/cni/:/opt/cni:ro"
}

function configNetworkPlugin() {
    if [[ "${NETWORK_PLUGIN}" = "azure" ]]; then
        configAzureCNI
    elif [[ "${NETWORK_PLUGIN}" = "kubenet" ]] ; then
        installCNI
    elif [[ "${NETWORK_POLICY}" = "flannel" ]] ; then
        configCNINetworkPolicy
    fi
}

function installClearContainersRuntime() {
	# Add Clear Containers repository key
	echo "Adding Clear Containers repository key..."
	curl -sSL --retry 5 --retry-delay 10 --retry-max-time 30 "https://download.opensuse.org/repositories/home:clearcontainers:clear-containers-3/xUbuntu_16.04/Release.key" | apt-key add -

	# Add Clear Container repository
	echo "Adding Clear Containers repository..."
	echo 'deb http://download.opensuse.org/repositories/home:/clearcontainers:/clear-containers-3/xUbuntu_16.04/ /' > /etc/apt/sources.list.d/cc-runtime.list

	# Install Clear Containers runtime
	echo "Installing Clear Containers runtime..."
	apt-get update && apt-get install --no-install-recommends -y \
		cc-runtime

	# Install the systemd service and socket files.
	local repo_uri="https://raw.githubusercontent.com/clearcontainers/proxy/3.0.23"
	curl -sSL --retry 5 --retry-delay 10 --retry-max-time 30 "${repo_uri}/cc-proxy.service.in" | sed 's#@libexecdir@#/usr/libexec#' > /etc/systemd/system/cc-proxy.service
	curl -sSL --retry 5 --retry-delay 10 --retry-max-time 30 "${repo_uri}/cc-proxy.socket.in" | sed 's#@localstatedir@#/var#' > /etc/systemd/system/cc-proxy.socket

	# Enable and start Clear Containers proxy service
	echo "Enabling and starting Clear Containers proxy service..."
	systemctlEnableAndStart cc-proxy

	setKubeletOpts " --container-runtime=remote --runtime-request-timeout=15m --container-runtime-endpoint=unix:///run/containerd/containerd.sock"
}

function installContainerd() {
	CRI_CONTAINERD_VERSION="1.1.0"
	local CONTAINERD_DOWNLOAD_URL="https://storage.googleapis.com/cri-containerd-release/cri-containerd-${CRI_CONTAINERD_VERSION}.linux-amd64.tar.gz"

    CONTAINERD_TGZ_TMP=/tmp/containerd.tar.gz
    retrycmd_get_tarball 60 1 "$CONTAINERD_TGZ_TMP" "$CONTAINERD_DOWNLOAD_URL"
	tar -xzf "$CONTAINERD_TGZ_TMP" -C /
	rm -f "$CONTAINERD_TGZ_TMP"

	echo "Successfully installed cri-containerd..."
	setupContainerd;
}

function setupContainerd() {
	echo "Configuring cri-containerd..."

	mkdir -p "/etc/containerd"
	CRI_CONTAINERD_CONFIG="/etc/containerd/config.toml"
	echo "subreaper = false" > "$CRI_CONTAINERD_CONFIG"
	echo "oom_score = 0" >> "$CRI_CONTAINERD_CONFIG"
	echo "[plugins.cri.containerd.untrusted_workload_runtime]" >> "$CRI_CONTAINERD_CONFIG"
	echo "runtime_type = 'io.containerd.runtime.v1.linux'" >> "$CRI_CONTAINERD_CONFIG"
	echo "runtime_engine = '/usr/bin/cc-runtime'" >> "$CRI_CONTAINERD_CONFIG"
	echo "[plugins.cri.containerd.default_runtime]" >> "$CRI_CONTAINERD_CONFIG"
	echo "runtime_type = 'io.containerd.runtime.v1.linux'" >> "$CRI_CONTAINERD_CONFIG"
	echo "runtime_engine = '/usr/local/sbin/runc'" >> "$CRI_CONTAINERD_CONFIG"
}

function ensureContainerd() {
	if [[ "$CONTAINER_RUNTIME" == "clear-containers" ]]; then
		# Make sure we can nest virtualization
		if grep -q vmx /proc/cpuinfo; then
			# Enable and start cri-containerd service
			# Make sure this is done after networking plugins are installed
			echo "Enabling and starting cri-containerd service..."
			systemctlEnableAndStart containerd
		fi
	fi
}

function systemctlEnableAndStart() {
    retrycmd_if_failure 10 1 3 systemctl daemon-reload
    systemctl enable $1
    systemctl is-enabled $1
    enabled=$?
    for i in {1..900}; do
        if [ $enabled -ne 0 ]; then
            systemctl enable $1
            systemctl is-enabled $1
            enabled=$?
        else
            echo "$1 took $i seconds to be enabled by systemctl"
            break
        fi
        sleep 1
    done
    if [ $enabled -ne 0 ]
    then
        echo "$1 could not be enabled by systemctl"
        exit 3
    fi
    systemctl_restart 100 1 10 $1
    retrycmd_if_failure 10 1 3 systemctl status $1 --no-pager -l > /var/log/azure/$1-status.log
    systemctl is-failed $1
    if [ $? -eq 0 ]
    then
        echo "$1 could not be started"
        exit 4
    fi
}

function ensureDocker() {
    systemctlEnableAndStart docker
}
function ensureKMS() {
    systemctlEnableAndStart kms
}

function ensureKubelet() {
    systemctlEnableAndStart kubelet
}

function extractHyperkube(){
    retrycmd_if_failure 100 1 60 docker pull $HYPERKUBE_URL
    systemctlEnableAndStart hyperkube-extract
}

function ensureJournal(){
    systemctlEnableAndStart systemd-journald
    echo "Storage=persistent" >> /etc/systemd/journald.conf
    echo "SystemMaxUse=1G" >> /etc/systemd/journald.conf
    echo "RuntimeMaxUse=1G" >> /etc/systemd/journald.conf
    echo "ForwardToSyslog=no" >> /etc/systemd/journald.conf
}

function ensureK8s() {
    if $REBOOTREQUIRED; then
        return
    fi
    k8sHealthy=1
    nodesActive=1
    nodesReady=1
    wait_for_file 600 1 $KUBECTL
    for i in {1..600}; do
        $KUBECTL 2>/dev/null cluster-info
            if [ "$?" = "0" ]
            then
                echo "k8s cluster is healthy, took $i seconds"
                k8sHealthy=0
                break
            fi
        sleep 1
    done
    if [ $k8sHealthy -ne 0 ]
    then
        echo "k8s cluster is not healthy after $i seconds"
        exit 30
    fi
    ensurePodSecurityPolicy
}

function ensureEtcd() {
    etcdIsRunning=1
    for i in {1..600}; do
        curl --cacert /etc/kubernetes/certs/ca.crt --cert /etc/kubernetes/certs/etcdclient.crt --key /etc/kubernetes/certs/etcdclient.key --max-time 60 https://127.0.0.1:2379/v2/machines;
        if [ $? -eq 0 ]
        then
            etcdIsRunning=0
            echo "Etcd setup successfully, took $i seconds"
            break
        fi
        sleep 1
    done
    if [ $etcdIsRunning -ne 0 ]
    then
        echo "Etcd not accessible after $i seconds"
        exit 11
    fi
}

function ensureEtcdDataDir() {
    mount | grep /dev/sdc1 | grep /var/lib/etcddisk
    if [ "$?" = "0" ]
    then
        echo "Etcd is running with data dir at: /var/lib/etcddisk"
        return
    else
        echo "/var/lib/etcddisk was not found at /dev/sdc1. Trying to mount all devices."
        s = 5
        for i in {1..60}; do
            sudo mount -a && mount | grep /dev/sdc1 | grep /var/lib/etcddisk;
            if [ "$?" = "0" ]
            then
                (( t = ${i} * ${s} ))
                echo "/var/lib/etcddisk mounted at: /dev/sdc1, took $t seconds"
                return
            fi
            sleep $s
        done
    fi

   echo "Etcd data dir was not found at: /var/lib/etcddisk"
   exit 10
}

function ensurePodSecurityPolicy(){
    if $REBOOTREQUIRED; then
        return
    fi
    POD_SECURITY_POLICY_FILE="/etc/kubernetes/manifests/pod-security-policy.yaml"
    if [ -f $POD_SECURITY_POLICY_FILE ]; then
        $KUBECTL create -f $POD_SECURITY_POLICY_FILE
    fi
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

if [[ "$CONTAINER_RUNTIME" == "clear-containers" ]]; then
	# If the container runtime is "clear-containers" we need to ensure the
	# run command is completed _before_ we start installing all the dependencies
	# for clear-containers to make sure there is not a dpkg lock.
	ensureRunCommandCompleted
	echo `date`,`hostname`, RunCmdCompleted>>/opt/m
fi

if [[ $OS == $UBUNTU_OS_NAME ]]; then
	# make sure walinuxagent doesn't get updated in the middle of running this script
	retrycmd_if_failure 20 5 5 apt-mark hold walinuxagent
fi

echo `date`,`hostname`, EnsureDockerStart>>/opt/m
ensureDockerInstallCompleted
ensureDocker
echo `date`,`hostname`, configNetworkPluginStart>>/opt/m
configNetworkPlugin
if [[ "$CONTAINER_RUNTIME" == "clear-containers" ]]; then
	# Ensure we can nest virtualization
	if grep -q vmx /proc/cpuinfo; then
		echo `date`,`hostname`, installClearContainersRuntimeStart>>/opt/m
		installClearContainersRuntime
		echo `date`,`hostname`, installContainerdStart>>/opt/m
		installContainerd
	fi
fi
echo `date`,`hostname`, ensureContainerdStart>>/opt/m
ensureContainerd
echo `date`,`hostname`, extractHyperkubeStart>>/opt/m
extractHyperkube
if [[ ! -z "${MASTER_NODE}" && ! -z "${EnableEncryptionWithExternalKms}" ]]; then
    echo `date`,`hostname`, ensureKMSStart>>/opt/m
    ensureKMS
fi
echo `date`,`hostname`, ensureKubeletStart>>/opt/m
ensureKubelet
echo `date`,`hostname`, ensureJournalStart>>/opt/m
ensureJournal
echo `date`,`hostname`, ensureJournalDone>>/opt/m
ensureRunCommandCompleted
echo `date`,`hostname`, RunCmdCompleted>>/opt/m

if [[ ! -z "${MASTER_NODE}" ]]; then
    writeKubeConfig
    ensureFilepath $KUBECTL
    ensureFilepath $DOCKER
    ensureEtcdDataDir
    ensureEtcd
    ensureK8s
fi

if [[ $OS == $UBUNTU_OS_NAME ]]; then
    # mitigation for bug https://bugs.launchpad.net/ubuntu/+source/linux/+bug/1676635
    echo 2dd1ce17-079e-403c-b352-a1921ee207ee > /sys/bus/vmbus/drivers/hv_util/unbind
    sed -i "13i\echo 2dd1ce17-079e-403c-b352-a1921ee207ee > /sys/bus/vmbus/drivers/hv_util/unbind\n" /etc/rc.local

    retrycmd_if_failure 20 5 5 apt-mark unhold walinuxagent
fi

echo "Install complete successfully"

if $REBOOTREQUIRED; then
  # wait 1 minute to restart node, so that the custom script extension can complete
  echo 'reboot required, rebooting node in 1 minute'
  /bin/bash -c "shutdown -r 1 &"
fi

echo `date`,`hostname`, endscript>>/opt/m

mkdir -p /opt/azure/containers && touch /opt/azure/containers/provision.complete
ps auxfww > /opt/azure/provision-ps.log &

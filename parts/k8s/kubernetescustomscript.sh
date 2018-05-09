#!/bin/bash

# This script runs on every Kubernetes VM
# Exit codes represent the following:
# | exit code number | meaning |
# | 3 | Service could not be enabled by systemctl |
# | 4 | Service could not be started by systemctl |
# | 5 | Timeout waiting for cloud-init runcmd to complete |
# | 6 | Timeout waiting for a file |
# | 7 | Error placing apt-mark hold on walinuxagent
# | 8 | Error releasing apt-mark hold on walinuxagent
# | 9 | Timeout installing apt packages
# | 10 | Etcd data dir not found |
# | 11 | Timeout waiting for etcd to be accessible |
# | 12 | Timeout waiting for etcd download(s) |
# | 13 | Unable to mount etcd volume |
# | 14 | Unable to start etcd service |
# | 15 | Unable to configure etcd membership |
# | 20 | Timeout waiting for docker install to finish |
# | 21 | Timeout waiting for docker download(s) |
# | 30 | Timeout waiting for k8s cluster to be healthy|
# | 31 | Timeout waiting for k8s download(s)|
# | 32 | Timeout waiting for kubectl|
# | 41 | Timeout waiting for CNI download(s)|

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
    /opt/azure/containers/setup-etcd.sh > /opt/azure/containers/setup-etcd.log 2>&1
    RET=$?
    if [ $RET -ne 0 ]; then
        exit $RET
    fi
    apt_get_update || exit 9
    retrycmd_if_failure 20 5 5 apt-mark hold walinuxagent  || exit 7
    /opt/azure/containers/mountetcd.sh || exit 13
    systemctl_restart 10 1 5 etcd || exit 14
    MEMBER="$(sudo etcdctl member list | grep -E ${MASTER_VM_NAME} | cut -d':' -f 1)"
    retrycmd_if_failure 10 1 5 sudo etcdctl member update $MEMBER ${ETCD_PEER_URL} || exit 15
    retrycmd_if_failure 5 1 10 curl --cacert /etc/kubernetes/certs/ca.crt --cert /etc/kubernetes/certs/etcdclient.crt --key /etc/kubernetes/certs/etcdclient.key --retry 5 --retry-delay 10 --retry-max-time 10 --max-time 60 ${ETCD_CLIENT_URL}/v2/machines || exit 11
else
    echo "skipping master node provision operations, this is an agent node"
    retrycmd_if_failure 10 1 3 systemctl enable rpcbind rpc-statd || exit 3
    systemctl_restart 20 1 10 rpcbind || exit 4
    systemctl_restart 20 1 10 rpc-statd || exit 4
fi

retrycmd_if_failure 20 5 5 apt-mark hold walinuxagent  || exit 7
apt_get_update || exit 9
retrycmd_if_failure 5 1 120 apt-get install -y apt-transport-https ca-certificates iptables iproute2 socat util-linux mount ebtables ethtool init-system-helpers || exit 9
retrycmd_if_failure_no_stats 180 1 5 curl -fsSL https://aptdocker.azureedge.net/gpg > /tmp/aptdocker.gpg || exit 21
retrycmd_if_failure 10 1 5 apt-key add /tmp/aptdocker.gpg || exit 9
echo "deb ${DOCKER_REPO} ubuntu-xenial main" | sudo tee /etc/apt/sources.list.d/docker.list
printf "Package: docker-engine\nPin: version ${DOCKER_ENGINE_VERSION}\nPin-Priority: 550\n" > /etc/apt/preferences.d/docker.pref
apt_get_update || exit 9
retrycmd_if_failure 20 1 120 apt-get install -y ebtables docker-engine || exit 20
echo "ExecStartPost=/sbin/iptables -P FORWARD ACCEPT" >> /etc/systemd/system/docker.service.d/exec_start.conf
mkdir -p /etc/kubernetes/manifests
usermod -aG docker ${ADMINUSER}
retrycmd_if_failure 20 1 10 /usr/lib/apt/apt.systemd.daily || exit 9
# TODO {{if EnableAggregatedAPIs}}
bash /etc/kubernetes/generate-proxy-certs.sh
retrycmd_if_failure 20 1 5 apt-mark unhold walinuxagent || exit 8

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
    if [ $? -ne 0 ]; then
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
    if [ $? -ne 0 ]; then
        echo "could not download required CNI artifact at ${CNI_PLUGINS_URL}"
        exit 41
    fi
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
    if [ $? -ne 0 ]; then
        echo "could not download required CNI artifact at ${VNET_CNI_PLUGINS_URL}"
        exit 41
    fi
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
}

function installContainerd() {
	CRI_CONTAINERD_VERSION="1.1.0"
	local CONTAINERD_DOWNLOAD_URL="https://storage.googleapis.com/cri-containerd-release/cri-containerd-${CRI_CONTAINERD_VERSION}.linux-amd64.tar.gz"

    CONTAINERD_TGZ_TMP=/tmp/containerd.tar.gz
    retrycmd_get_tarball 60 1 "$CONTAINERD_TGZ_TMP" "$CONTAINERD_DOWNLOAD_URL"
    if [ $? -ne 0 ]; then
        echo "could not download required CNI artifact at $CONTAINERD_DOWNLOAD_URL"
        exit 41
    fi
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
	if [[ "$CONTAINER_RUNTIME" == "clear-containers" ]]; then
		echo "runtime_engine = '/usr/bin/cc-runtime'" >> "$CRI_CONTAINERD_CONFIG"
	else
		echo "runtime_engine = '/usr/local/sbin/runc'" >> "$CRI_CONTAINERD_CONFIG"
	fi
	echo "[plugins.cri.containerd.default_runtime]" >> "$CRI_CONTAINERD_CONFIG"
	echo "runtime_type = 'io.containerd.runtime.v1.linux'" >> "$CRI_CONTAINERD_CONFIG"
	echo "runtime_engine = '/usr/local/sbin/runc'" >> "$CRI_CONTAINERD_CONFIG"

	setKubeletOpts " --container-runtime=remote --runtime-request-timeout=15m --container-runtime-endpoint=unix:///run/containerd/containerd.sock"
}

function ensureContainerd() {
	if [[ "$CONTAINER_RUNTIME" == "clear-containers" ]] || [[ "$CONTAINER_RUNTIME" == "containerd" ]]; then
		# Enable and start cri-containerd service
		# Make sure this is done after networking plugins are installed
		echo "Enabling and starting cri-containerd service..."
		systemctlEnableAndStart containerd
	fi
}

function systemctlEnableAndStart() {
    systemctl_restart 20 1 10 $1
    RESTART_STATUS=$?
    systemctl status $1 --no-pager -l > /var/log/azure/$1-status.log
    if [ $RESTART_STATUS -ne 0 ]; then
        echo "$1 could not be started"
        exit 4
    fi
    retrycmd_if_failure 10 1 3 systemctl enable $1
    if [ $? -ne 0 ]; then
        echo "$1 could not be enabled by systemctl"
        exit 3
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
    if [ $? -ne 0 ]; then
        echo "required kubernetes docker image could not be downloaded at $HYPERKUBE_URL"
        exit 31
    fi
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
    if [ $? -ne 0 ]; then
        echo "could not find kubectl at $KUBECTL"
        exit 32
    fi
    for i in {1..600}; do
        $KUBECTL 2>/dev/null cluster-info
            if [ $? -eq 0 ]; then
                echo "k8s cluster is healthy, took $i seconds"
                k8sHealthy=0
                break
            fi
        sleep 1
    done
    if [ $k8sHealthy -ne 0 ]; then
        echo "k8s cluster is not healthy after $i seconds"
        exit 30
    fi
    ensurePodSecurityPolicy
}

function ensureEtcd() {
    etcdIsRunning=1
    for i in {1..600}; do
        curl --cacert /etc/kubernetes/certs/ca.crt --cert /etc/kubernetes/certs/etcdclient.crt --key /etc/kubernetes/certs/etcdclient.key --max-time 60 https://127.0.0.1:2379/v2/machines
        if [ $? -eq 0 ]; then
            etcdIsRunning=0
            echo "Etcd setup successfully, took $i seconds"
            break
        fi
        sleep 1
    done
    if [ $etcdIsRunning -ne 0 ]; then
        echo "Etcd not accessible after $i seconds"
        exit 11
    fi
}

function ensureEtcdDataDir() {
    mount | grep /dev/sdc1 | grep /var/lib/etcddisk
    if [ $? -eq 0 ]; then
        echo "Etcd is running with data dir at: /var/lib/etcddisk"
        return
    else
        echo "/var/lib/etcddisk was not found at /dev/sdc1. Trying to mount all devices."
        s = 5
        for i in {1..60}; do
            sudo mount -a && mount | grep /dev/sdc1 | grep /var/lib/etcddisk;
            if [ "$?" = "0" ]; then
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

function configAddons() {
    if [[ "${CLUSTER_AUTOSCALER_ADDON}" = True ]]; then
        configClusterAutoscalerAddon
    fi
    echo `date`,`hostname`, configAddonsDone>>/opt/m
}

function configClusterAutoscalerAddon() {
    echo `date`,`hostname`, configClusterAutoscalerAddonStart>>/opt/m

    if [[ "${USE_MANAGED_IDENTITY_EXTENSION}" == true ]]; then
        echo `date`,`hostname`, configClusterAutoscalerAddonManagedIdentityStart>>/opt/m
        CLUSTER_AUTOSCALER_MSI_VOLUME_MOUNT="- mountPath: /var/lib/waagent/\n\          name: waagent\n\          readOnly: true"
        CLUSTER_AUTOSCALER_MSI_VOLUME="- hostPath:\n\          path: /var/lib/waagent/\n\        name: waagent"
        CLUSTER_AUTOSCALER_MSI_HOST_NETWORK="hostNetwork: true"

        sed -i "s|<kubernetesClusterAutoscalerVolumeMounts>|${CLUSTER_AUTOSCALER_MSI_VOLUME_MOUNT}|g" "/etc/kubernetes/addons/cluster-autoscaler-deployment.yaml"
        sed -i "s|<kubernetesClusterAutoscalerVolumes>|${CLUSTER_AUTOSCALER_MSI_VOLUME}|g" "/etc/kubernetes/addons/cluster-autoscaler-deployment.yaml"
        sed -i "s|<kubernetesClusterAutoscalerHostNetwork>|$(echo "${CLUSTER_AUTOSCALER_MSI_HOST_NETWORK}")|g" "/etc/kubernetes/addons/cluster-autoscaler-deployment.yaml"
        echo `date`,`hostname`, configClusterAutoscalerAddonManagedIdentityDone>>/opt/m
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
    echo `date`,`hostname`, configClusterAutoscalerAddonDone>>/opt/m
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
    if [ $? -ne 0 ]; then
        echo "error placing apt-mark hold on walinuxagent"
        exit 7
    fi
    
fi

ensureDocker
echo `date`,`hostname`, configNetworkPluginStart>>/opt/m
configNetworkPlugin
echo `date`,`hostname`, configAddonsStart>>/opt/m
configAddons
if [[ "$CONTAINER_RUNTIME" == "clear-containers" ]]; then
	# Ensure we can nest virtualization
	if grep -q vmx /proc/cpuinfo; then
		echo `date`,`hostname`, installClearContainersRuntimeStart>>/opt/m
		installClearContainersRuntime
	fi
fi
if [[ "$CONTAINER_RUNTIME" == "clear-containers" ]] || [[ "$CONTAINER_RUNTIME" == "containerd" ]]; then
	echo `date`,`hostname`, installContainerdStart>>/opt/m
	installContainerd
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
    if [ $? -ne 0 ]; then
        echo "error releasing apt-mark hold on walinuxagent"
        exit 8
    fi
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

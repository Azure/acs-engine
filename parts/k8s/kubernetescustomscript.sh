#!/bin/bash

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
    for i in {1..900}; do
        if [ -e /opt/azure/containers/runcmd.complete ]; then
            echo "runcmd finished, took $i seconds"
            break
        fi
        sleep 1
    done
}

ensureDockerInstallCompleted()
{
    echo "waiting for docker install to finish"
    for i in {1..900}; do
        if [ -e /opt/azure/containers/dockerinstall.complete ]; then
            echo "docker install finished, took $i seconds"
            break
        fi
        sleep 1
    done
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
    "subnetName": "${SUBNET}",
    "securityGroupName": "${NETWORK_SECURITY_GROUP}",
    "vnetName": "${VIRTUAL_NETWORK}",
    "vnetResourceGroup": "${VIRTUAL_NETWORK_RESOURCE_GROUP}",
    "routeTableName": "${ROUTE_TABLE}",
    "primaryAvailabilitySetName": "${PRIMARY_AVAILABILITY_SET}",
    "cloudProviderBackoff": ${CLOUDPROVIDER_BACKOFF},
    "cloudProviderBackoffRetries": ${CLOUDPROVIDER_BACKOFF_RETRIES},
    "cloudProviderBackoffExponent": ${CLOUDPROVIDER_BACKOFF_EXPONENT},
    "cloudProviderBackoffDuration": ${CLOUDPROVIDER_BACKOFF_DURATION},
    "cloudProviderBackoffJitter": ${CLOUDPROVIDER_BACKOFF_JITTER},
    "cloudProviderRatelimit": ${CLOUDPROVIDER_RATELIMIT},
    "cloudProviderRateLimitQPS": ${CLOUDPROVIDER_RATELIMIT_QPS},
    "cloudProviderRateLimitBucket": ${CLOUDPROVIDER_RATELIMIT_BUCKET},
    "useManagedIdentityExtension": ${USE_MANAGED_IDENTITY_EXTENSION},
    "useInstanceMetadata": ${USE_INSTANCE_METADATA}
}
EOF

set -x

function ensureFilepath() {
    if $REBOOTREQUIRED; then
        return
    fi
    found=1
    for i in {1..600}; do
        if [ -e $1 ]
        then
            found=0
            echo "$1 is present, took $i seconds to verify"
            break
        fi
        sleep 1
    done
    if [ $found -ne 0 ]
    then
        echo "$1 is not present after $i seconds of trying to verify"
        exit 1
    fi
}

function setKubeletOpts () {
	sed -i "s#^KUBELET_OPTS=.*#KUBELET_OPTS=${1}#" /etc/default/kubelet
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
    CONTAINERNETWORKING_CNI_TGZ_TMP=/tmp/containernetworking_cni.tgz
    retrycmd_get_tarball 60 1 $CONTAINERNETWORKING_CNI_TGZ_TMP ${CNI_PLUGINS_URL}
    tar -xzf $CONTAINERNETWORKING_CNI_TGZ_TMP -C $CNI_BIN_DIR ./loopback ./portmap
    chown -R root:root $CNI_BIN_DIR
    chmod -R 755 $CNI_BIN_DIR
    mv $CNI_BIN_DIR/10-azure.conflist $CNI_CONFIG_DIR/
    chmod 600 $CNI_CONFIG_DIR/10-azure.conflist
    /sbin/ebtables -t nat --list
}

function configKubenet() {
    CNI_BIN_DIR=/opt/cni/bin
    mkdir -p $CNI_BIN_DIR
    CONTAINERNETWORKING_CNI_TGZ_TMP=/tmp/containernetworking_cni.tgz
    retrycmd_get_tarball 60 1 $CONTAINERNETWORKING_CNI_TGZ_TMP ${CNI_PLUGINS_URL}
    tar -xzf $CONTAINERNETWORKING_CNI_TGZ_TMP -C $CNI_BIN_DIR ./loopback ./bridge ./host-local
    chown -R root:root $CNI_BIN_DIR
    chmod -R 755 $CNI_BIN_DIR
}

function configNetworkPolicy() {
    if [[ "${NETWORK_POLICY}" = "azure" ]]; then
        configAzureCNI
    elif [[ "${NETWORK_POLICY}" = "none" ]] ; then
        configKubenet
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

	# Load systemd changes
	echo "Loading changes to systemd service files..."
	systemctl daemon-reload

	# Enable and start Clear Containers proxy service
	echo "Enabling and starting Clear Containers proxy service..."
	systemctl enable cc-proxy
	systemctl start cc-proxy

	setKubeletOpts " --container-runtime=remote --runtime-request-timeout=15m --container-runtime-endpoint=unix:///run/containerd/containerd.sock"
}

function installContainerd() {
	CRI_CONTAINERD_VERSION="1.1.0-rc.0"
	local download_uri="https://storage.googleapis.com/cri-containerd-release/cri-containerd-${CRI_CONTAINERD_VERSION}.linux-amd64.tar.gz"

	curl -sSL "$download_uri" | tar -xz -C /

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

	systemctl daemon-reload
}

function ensureContainerd() {
	if [[ "$CONTAINER_RUNTIME" == "clear-containers" ]]; then
		# Make sure we can nest virtualization
		if grep -q vmx /proc/cpuinfo; then
			# Enable and start cri-containerd service
			# Make sure this is done after networking plugins are installed
			echo "Enabling and starting cri-containerd service..."
			systemctl enable containerd
			systemctl start containerd
		fi
	fi
}

function systemctlEnableAndCheck() {
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
        exit 5
    fi
}

function ensureDocker() {
    systemctlEnableAndCheck docker
    # only start if a reboot is not required
    if ! $REBOOTREQUIRED; then
        dockerStarted=1
        for i in {1..900}; do
            if ! timeout 10s $DOCKER info; then
                echo "status $?"
                timeout 60s /bin/systemctl restart docker
            else
                echo "docker started, took $i seconds"
                dockerStarted=0
                break
            fi
            sleep 1
        done
        if [ $dockerStarted -ne 0 ]
        then
            echo "docker did not start"
            exit 2
        fi
    fi
}

function ensureKubelet() {
    systemctlEnableAndCheck kubelet
    # only start if a reboot is not required
    if ! $REBOOTREQUIRED; then
        systemctl restart kubelet
    fi
}

function extractHyperkube(){
    retrycmd_if_failure 100 1 60 docker pull $HYPERKUBE_URL
    systemctlEnableAndCheck hyperkube-extract
    # only start if a reboot is not required
    if ! $REBOOTREQUIRED; then
        systemctl restart hyperkube-extract
    fi
}

function ensureJournal(){
    systemctl daemon-reload
    systemctlEnableAndCheck systemd-journald.service
    echo "Storage=persistent" >> /etc/systemd/journald.conf
    echo "SystemMaxUse=1G" >> /etc/systemd/journald.conf
    echo "RuntimeMaxUse=1G" >> /etc/systemd/journald.conf
    echo "ForwardToSyslog=no" >> /etc/systemd/journald.conf
    # only start if a reboot is not required
    if ! $REBOOTREQUIRED; then
        systemctl restart systemd-journald.service
    fi
}

function ensureK8s() {
    if $REBOOTREQUIRED; then
        return
    fi
    k8sHealthy=1
    nodesActive=1
    nodesReady=1
    for i in {1..600}; do
        if [ -e $KUBECTL ]
        then
            break
        fi
        sleep 1
    done
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
        exit 3
    fi
    for i in {1..1800}; do
        nodes=$(${KUBECTL} get nodes 2>/dev/null | grep 'Ready' | wc -l)
            if [ $nodes -eq $TOTAL_NODES ]
            then
                echo "all nodes are participating, took $i seconds"
                nodesActive=0
                break
            fi
        sleep 1
    done
    if [ $nodesActive -ne 0 ]
    then
        echo "still waiting for active nodes after $i seconds"
        exit 3
    fi
    for i in {1..600}; do
        notReady=$(${KUBECTL} get nodes 2>/dev/null | grep 'NotReady' | wc -l)
            if [ $notReady -eq 0 ]
            then
                echo "all nodes are Ready, took $i seconds"
                nodesReady=0
                break
            fi
        sleep 1
    done
    if [ $nodesReady -ne 0 ]
    then
        echo "still waiting for Ready nodes after $i seconds"
        exit 3
    fi
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
        exit 3
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
   exit 4
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
echo `date`,`hostname`, configNetworkPolicyStart>>/opt/m
configNetworkPolicy
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
    ensurePodSecurityPolicy
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

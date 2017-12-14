#!/bin/bash

###########################################################
# START SECRET DATA - ECHO DISABLED
###########################################################

# Following parameters now read from environment variable
# Fields for `azure.json`
# TENANT_ID SUBSCRIPTION_ID RESOURCE_GROUP LOCATION SUBNET
# NETWORK_SECURITY_GROUP VIRTUAL_NETWORK VIRTUAL_NETWORK_RESOURCE_GROUP ROUTE_TABLE PRIMARY_AVAILABILITY_SET
# SERVICE_PRINCIPAL_CLIENT_ID SERVICE_PRINCIPAL_CLIENT_SECRET KUBELET_PRIVATE_KEY TARGET_ENVIRONMENT NETWORK_POLICY
# FQDNSuffix VNET_CNI_PLUGINS_URL CNI_PLUGINS_URL MAX_PODS

# Default values for backoff configuration
# CLOUDPROVIDER_BACKOFF CLOUDPROVIDER_BACKOFF_RETRIES CLOUDPROVIDER_BACKOFF_EXPONENT CLOUDPROVIDER_BACKOFF_DURATION CLOUDPROVIDER_BACKOFF_JITTER
# Default values for rate limit configuration
# CLOUDPROVIDER_RATELIMIT CLOUDPROVIDER_RATELIMIT_QPS CLOUDPROVIDER_RATELIMIT_BUCKET

# USE_MANAGED_IDENTITY_EXTENSION USE_INSTANCE_METADATA

# Master only secrets
# APISERVER_PRIVATE_KEY CA_CERTIFICATE CA_PRIVATE_KEY MASTER_FQDN KUBECONFIG_CERTIFICATE
# KUBECONFIG_KEY ETCD_SERVER_CERTIFICATE ETCD_SERVER_PRIVATE_KEY ETCD_CLIENT_CERTIFICATE ETCD_CLIENT_PRIVATE_KEY 
# ETCD_PEER_CERTIFICATES ETCD_PEER_PRIVATE_KEYS ADMINUSER MASTER_INDEX

# Find distro name via ID value in releases files and upcase
OS=$(cat /etc/*-release | grep ^ID= | tr -d 'ID="' | awk '{print toupper($0)}')
UBUNTU_OS_NAME="UBUNTU"
RHEL_OS_NAME="RHEL"
COREOS_OS_NAME="COREOS"

# Set default kubectl
KUBECTL=/usr/local/bin/kubectl

ETCD_PEER_CERT=$(echo ${ETCD_PEER_CERTIFICATES} | cut -d'[' -f 2 | cut -d']' -f 1 | cut -d',' -f $((${MASTER_INDEX}+1)))
ETCD_PEER_KEY=$(echo ${ETCD_PEER_PRIVATE_KEYS} | cut -d'[' -f 2 | cut -d']' -f 1 | cut -d',' -f $((${MASTER_INDEX}+1)))

# CoreOS: /usr is read-only; therefore kubectl is installed at /opt/kubectl
#   Details on install at kubernetetsmastercustomdataforcoreos.yml
if [[ $OS == $COREOS_OS_NAME ]]; then
    echo "Changing default kubectl bin location"
    KUBECTL=/opt/kubectl
fi

# cloudinit runcmd and the extension will run in parallel, this is to ensure
# runcmd finishes
ensureRunCommandCompleted()
{
    echo "waiting for runcmd to finish"
    for i in {1..900}; do
        if [ -e /opt/azure/containers/runcmd.complete ]; then
            echo "runcmd finished"
            break
        fi
        sleep 1
    done
}

echo `date`,`hostname`, startscript>>/opt/m

# A delay to start the kubernetes processes is necessary
# if a reboot is required.  Otherwise, the agents will encounter issue: 
# https://github.com/kubernetes/kubernetes/issues/41185
if [ -f /var/run/reboot-required ]; then
    REBOOTREQUIRED=true
else
    REBOOTREQUIRED=false
fi

# If APISERVER_PRIVATE_KEY is empty, then we are not on the master
if [[ ! -z "${APISERVER_PRIVATE_KEY}" ]]; then
    echo "APISERVER_PRIVATE_KEY is non-empty, assuming master node"

    APISERVER_PRIVATE_KEY_PATH="/etc/kubernetes/certs/apiserver.key"
    touch "${APISERVER_PRIVATE_KEY_PATH}"
    chmod 0600 "${APISERVER_PRIVATE_KEY_PATH}"
    chown root:root "${APISERVER_PRIVATE_KEY_PATH}"
    echo "${APISERVER_PRIVATE_KEY}" | base64 --decode > "${APISERVER_PRIVATE_KEY_PATH}"
else
    echo "APISERVER_PRIVATE_KEY is empty, assuming worker node"
fi

# If CA_PRIVATE_KEY is empty, then we are not on the master
if [[ ! -z "${CA_PRIVATE_KEY}" ]]; then
    echo "CA_KEY is non-empty, assuming master node"

    CA_PRIVATE_KEY_PATH="/etc/kubernetes/certs/ca.key"
    touch "${CA_PRIVATE_KEY_PATH}"
    chmod 0600 "${CA_PRIVATE_KEY_PATH}"
    chown root:root "${CA_PRIVATE_KEY_PATH}"
    echo "${CA_PRIVATE_KEY}" | base64 --decode > "${CA_PRIVATE_KEY_PATH}"
else
    echo "CA_PRIVATE_KEY is empty, assuming worker node"
fi

ETCD_SERVER_PRIVATE_KEY_PATH="/etc/kubernetes/certs/etcdserver.key"
touch "${ETCD_SERVER_PRIVATE_KEY_PATH}"
chmod 0600 "${ETCD_SERVER_PRIVATE_KEY_PATH}"
chown root:root "${ETCD_SERVER_PRIVATE_KEY_PATH}"
echo "${ETCD_SERVER_PRIVATE_KEY}" | base64 --decode > "${ETCD_SERVER_PRIVATE_KEY_PATH}"

ETCD_CLIENT_PRIVATE_KEY_PATH="/etc/kubernetes/certs/etcdclient.key"
touch "${ETCD_CLIENT_PRIVATE_KEY_PATH}"
chmod 0600 "${ETCD_CLIENT_PRIVATE_KEY_PATH}"
chown root:root "${ETCD_CLIENT_PRIVATE_KEY_PATH}"
echo "${ETCD_CLIENT_PRIVATE_KEY}" | base64 --decode > "${ETCD_CLIENT_PRIVATE_KEY_PATH}"

ETCD_PEER_PRIVATE_KEY_PATH="/etc/kubernetes/certs/etcdpeer${MASTER_INDEX}.key"
touch "${ETCD_PEER_PRIVATE_KEY_PATH}"
chmod 0600 "${ETCD_PEER_PRIVATE_KEY_PATH}"
chown root:root "${ETCD_PEER_PRIVATE_KEY_PATH}"
echo "${ETCD_PEER_KEY}" | base64 --decode > "${ETCD_PEER_PRIVATE_KEY_PATH}"

ETCD_SERVER_CERTIFICATE_PATH="/etc/kubernetes/certs/etcdserver.crt"
touch "${ETCD_SERVER_CERTIFICATE_PATH}"
chmod 0644 "${ETCD_SERVER_CERTIFICATE_PATH}"
chown root:root "${ETCD_SERVER_CERTIFICATE_PATH}"
echo "${ETCD_SERVER_CERTIFICATE}" | base64 --decode > "${ETCD_SERVER_CERTIFICATE_PATH}"

ETCD_CLIENT_CERTIFICATE_PATH="/etc/kubernetes/certs/etcdclient.crt"
touch "${ETCD_CLIENT_CERTIFICATE_PATH}"
chmod 0644 "${ETCD_CLIENT_CERTIFICATE_PATH}"
chown root:root "${ETCD_CLIENT_CERTIFICATE_PATH}"
echo "${ETCD_CLIENT_CERTIFICATE}" | base64 --decode > "${ETCD_CLIENT_CERTIFICATE_PATH}"

ETCD_PEER_CERTIFICATE_PATH="/etc/kubernetes/certs/etcdpeer${MASTER_INDEX}.crt"
touch "${ETCD_PEER_CERTIFICATE_PATH}"
chmod 0644 "${ETCD_PEER_CERTIFICATE_PATH}"
chown root:root "${ETCD_PEER_CERTIFICATE_PATH}"
echo "${ETCD_PEER_CERT}" | base64 --decode > "${ETCD_PEER_CERTIFICATE_PATH}"

echo `date`,`hostname`, finishedGettingEtcdCerts>>/opt/m
mkdir -p /opt/azure/containers && touch /opt/azure/containers/etcdcerts.complete

KUBELET_PRIVATE_KEY_PATH="/etc/kubernetes/certs/client.key"
touch "${KUBELET_PRIVATE_KEY_PATH}"
chmod 0600 "${KUBELET_PRIVATE_KEY_PATH}"
chown root:root "${KUBELET_PRIVATE_KEY_PATH}"
echo "${KUBELET_PRIVATE_KEY}" | base64 --decode > "${KUBELET_PRIVATE_KEY_PATH}"

APISERVER_PUBLIC_KEY_PATH="/etc/kubernetes/certs/apiserver.crt"
touch "${APISERVER_PUBLIC_KEY_PATH}"
chmod 0644 "${APISERVER_PUBLIC_KEY_PATH}"
chown root:root "${APISERVER_PUBLIC_KEY_PATH}"
echo "${APISERVER_PUBLIC_KEY}" | base64 --decode > "${APISERVER_PUBLIC_KEY_PATH}"

AZURE_JSON_PATH="/etc/kubernetes/azure.json"
touch "${AZURE_JSON_PATH}"
chmod 0600 "${AZURE_JSON_PATH}"
chown root:root "${AZURE_JSON_PATH}"
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

###########################################################
# END OF SECRET DATA
###########################################################

set -x

# wait for kubectl to report successful cluster health
function ensureKubectl() {
    if $REBOOTREQUIRED; then
        return
    fi
    kubectlfound=1
    for i in {1..600}; do
        if [ -e $KUBECTL ]
        then
            kubectlfound=0
            break
        fi
        sleep 1
    done
    if [ $kubectlfound -ne 0 ]
    then
        if [ ! -e /usr/bin/docker ]
        then
            echo "kubectl nor docker did not install successfully"
            exit 1
        fi
    fi
}

function downloadUrl () {
	# Wrapper around curl to download blobs more reliably.
	# Workaround the --retry issues with a for loop and set a max timeout.
	for i in 1 2 3 4 5; do curl --max-time 60 -fsSL ${1}; [ $? -eq 0 ] && break || sleep 10; done
}

function setMaxPods () {
    sed -i "s/^KUBELET_MAX_PODS=.*/KUBELET_MAX_PODS=${1}/" /etc/default/kubelet
}

function setNetworkPlugin () {
    sed -i "s/^KUBELET_NETWORK_PLUGIN=.*/KUBELET_NETWORK_PLUGIN=${1}/" /etc/default/kubelet
}

function setDockerOpts () {
    sed -i "s#^DOCKER_OPTS=.*#DOCKER_OPTS=${1}#" /etc/default/kubelet
}

function configAzureNetworkPolicy() {
    CNI_CONFIG_DIR=/etc/cni/net.d
    mkdir -p $CNI_CONFIG_DIR

    chown -R root:root $CNI_CONFIG_DIR
    chmod 755 $CNI_CONFIG_DIR

    # Download Azure VNET CNI plugins.
    CNI_BIN_DIR=/opt/cni/bin
    mkdir -p $CNI_BIN_DIR

    # Mirror from https://github.com/Azure/azure-container-networking/releases/tag/$AZURE_PLUGIN_VER/azure-vnet-cni-linux-amd64-$AZURE_PLUGIN_VER.tgz
    downloadUrl ${VNET_CNI_PLUGINS_URL} | tar -xz -C $CNI_BIN_DIR
    # Mirror from https://github.com/containernetworking/cni/releases/download/$CNI_RELEASE_VER/cni-amd64-$CNI_RELEASE_VERSION.tgz
    downloadUrl ${CNI_PLUGINS_URL} | tar -xz -C $CNI_BIN_DIR ./loopback
    chown -R root:root $CNI_BIN_DIR
    chmod -R 755 $CNI_BIN_DIR

    # Copy config file
    mv $CNI_BIN_DIR/10-azure.conf $CNI_CONFIG_DIR/
    chmod 600 $CNI_CONFIG_DIR/10-azure.conf

    # Dump ebtables rules.
    /sbin/ebtables -t nat --list

    # Enable CNI.
    setNetworkPlugin cni
    setDockerOpts " --volume=/etc/cni/:/etc/cni:ro --volume=/opt/cni/:/opt/cni:ro"
}

# Configures Kubelet to use CNI and mount the appropriate hostpaths
function configCalicoNetworkPolicy() {
    setNetworkPlugin cni
    setDockerOpts " --volume=/etc/cni/:/etc/cni:ro --volume=/opt/cni/:/opt/cni:ro"
}

function configNetworkPolicy() {
    if [[ "${NETWORK_POLICY}" = "azure" ]]; then
        configAzureNetworkPolicy
    elif [[ "${NETWORK_POLICY}" = "calico" ]]; then
        configCalicoNetworkPolicy
    else
        # No policy, defaults to kubenet.
        setNetworkPlugin kubenet
        setDockerOpts ""
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
            break
        fi
        sleep 1
    done
    if [ $enabled -ne 0 ]
    then
        echo "$1 could not be enabled by systemctl"
        exit 5
    fi
    systemctl enable $1
}

function ensureDocker() {
    systemctlEnableAndCheck docker
    # only start if a reboot is not required
    if ! $REBOOTREQUIRED; then
        systemctl restart docker
        dockerStarted=1
        for i in {1..900}; do
            if ! /usr/bin/docker info; then
                echo "status $?"
                /bin/systemctl restart docker
            else
                echo "docker started"
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

function extractKubectl(){
    systemctlEnableAndCheck kubectl-extract
    # only start if a reboot is not required
    if ! $REBOOTREQUIRED; then
        systemctl restart kubectl-extract
    fi
}

function ensureJournal(){
    systemctl daemon-reload
    systemctlEnableAndCheck systemd-journald.service
    # only start if a reboot is not required
    if ! $REBOOTREQUIRED; then
        systemctl restart systemd-journald.service
    fi
}

function ensureApiserver() {
    if $REBOOTREQUIRED; then
        return
    fi
    kubernetesStarted=1
    for i in {1..600}; do
        if [ -e $KUBECTL ]
        then
            $KUBECTL cluster-info
            if [ "$?" = "0" ]
            then
                echo "kubernetes started"
                kubernetesStarted=0
                break
            fi
        else
            /usr/bin/docker ps | grep apiserver
            if [ "$?" = "0" ]
            then
                echo "kubernetes started"
                kubernetesStarted=0
                break
            fi
        fi
        sleep 1
    done
    if [ $kubernetesStarted -ne 0 ]
    then
        echo "kubernetes did not start"
        exit 3
    fi
}

function ensureEtcd() {
    for i in {1..600}; do
        curl --cacert /etc/kubernetes/certs/ca.crt --cert /etc/kubernetes/certs/etcdclient.crt --key /etc/kubernetes/certs/etcdclient.key --max-time 60 https://127.0.0.1:2379/v2/machines;
        if [ $? -eq 0 ]
        then
            echo "Etcd setup successfully"
            break
        fi
        sleep 5
    done
}

function ensureEtcdDataDir() {
    mount | grep /dev/sdc1 | grep /var/lib/etcddisk
    if [ "$?" = "0" ]
    then
        echo "Etcd is running with data dir at: /var/lib/etcddisk"
        return
    else
        echo "/var/lib/etcddisk was not found at /dev/sdc1. Trying to mount all devices."
        for i in {1..60}; do
            sudo mount -a && mount | grep /dev/sdc1 | grep /var/lib/etcddisk;
            if [ "$?" = "0" ]
            then
                echo "/var/lib/etcddisk mounted at: /dev/sdc1"
                return
            fi
            sleep 5
        done
    fi

   echo "Etcd data dir was not found at: /var/lib/etcddisk"
   exit 4
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
    server: https://$MASTER_FQDN.$LOCATION.$FQDNSuffix
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

# master and node
echo `date`,`hostname`, EnsureDockerStart>>/opt/m 
ensureDocker
echo `date`,`hostname`, configNetworkPolicyStart>>/opt/m 
configNetworkPolicy
echo `date`,`hostname`, setMaxPodsStart>>/opt/m 
setMaxPods ${MAX_PODS}
echo `date`,`hostname`, ensureKubeletStart>>/opt/m
ensureKubelet
echo `date`,`hostname`, extractKubctlStart>>/opt/m 
extractKubectl
echo `date`,`hostname`, ensureJournalStart>>/opt/m 
ensureJournal
echo `date`,`hostname`, ensureJournalDone>>/opt/m 

ensureRunCommandCompleted
echo `date`,`hostname`, RunCmdCompleted>>/opt/m 

if [[ $OS == $UBUNTU_OS_NAME ]]; then
    # make sure walinuxagent doesn't get updated in the middle of running this script
    apt-mark hold walinuxagent
fi

# master only
if [[ ! -z "${APISERVER_PRIVATE_KEY}" ]]; then
    writeKubeConfig
    ensureKubectl
    ensureEtcdDataDir
    ensureEtcd
    ensureApiserver
fi

if [[ $OS == $UBUNTU_OS_NAME ]]; then
    # mitigation for bug https://bugs.launchpad.net/ubuntu/+source/linux/+bug/1676635
    echo 2dd1ce17-079e-403c-b352-a1921ee207ee > /sys/bus/vmbus/drivers/hv_util/unbind
    sed -i "13i\echo 2dd1ce17-079e-403c-b352-a1921ee207ee > /sys/bus/vmbus/drivers/hv_util/unbind\n" /etc/rc.local

    # If APISERVER_PRIVATE_KEY is empty, then we are not on the master
    apt-mark unhold walinuxagent
fi

echo "Install complete successfully"

if $REBOOTREQUIRED; then
  # wait 1 minute to restart node, so that the custom script extension can complete
  echo 'reboot required, rebooting node in 1 minute'
  /bin/bash -c "shutdown -r 1 &"
fi

echo `date`,`hostname`, endscript>>/opt/m

mkdir -p /opt/azure/containers && touch /opt/azure/containers/provision.complete
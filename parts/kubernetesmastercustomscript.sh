#!/bin/bash

###########################################################
# START SECRET DATA - ECHO DISABLED
###########################################################

# Fields for `azure.json`
TENANT_ID="${1}"
SUBNETSCRIPTION_ID="${2}"
RESOURCE_GROUP="${3}"
LOCATION="${4}"
SUBNET="${5}"
NETWORK_SECURITY_GROUP="${6}"
VIRTUAL_NETWORK="${7}"
ROUTE_TABLE="${8}"
PRIMARY_AVAILABILITY_SET="${9}"
SERVICE_PRINCIPAL_CLIENT_ID="${10}"
SERVICE_PRINCIPAL_CLIENT_SECRET="${11}"
KUBELET_PRIVATE_KEY="${12}"
TARGET_ENVIRONMENT="${13}"
NETWORK_POLICY="${14}"

# Default values for backoff configuration
CLOUDPROVIDER_BACKOFF="${15}"
CLOUDPROVIDER_BACKOFF_RETRIES="${16}"
CLOUDPROVIDER_BACKOFF_EXPONENT="${17}"
CLOUDPROVIDER_BACKOFF_DURATION="${18}"
CLOUDPROVIDER_BACKOFF_JITTER="${19}"
# Default values for rate limit configuration
CLOUDPROVIDER_RATELIMIT="${20}"
CLOUDPROVIDER_RATELIMIT_QPS="${21}"
CLOUDPROVIDER_RATELIMIT_BUCKET="${22}"

USE_MANAGED_IDENTITY_EXTENSION="${23}"
USE_INSTANCE_METADATA="${24}"

# Master only secrets
APISERVER_PRIVATE_KEY="${25}"
CA_CERTIFICATE="${26}"
CA_PRIVATE_KEY="${27}"
MASTER_FQDN="${28}"
KUBECONFIG_CERTIFICATE="${29}"
KUBECONFIG_KEY="${30}"
ADMINUSER="${31}"

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
ensureRunCommandCompleted

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

KUBELET_PRIVATE_KEY_PATH="/etc/kubernetes/certs/client.key"
touch "${KUBELET_PRIVATE_KEY_PATH}"
chmod 0600 "${KUBELET_PRIVATE_KEY_PATH}"
chown root:root "${KUBELET_PRIVATE_KEY_PATH}"
echo "${KUBELET_PRIVATE_KEY}" | base64 --decode > "${KUBELET_PRIVATE_KEY_PATH}"

AZURE_JSON_PATH="/etc/kubernetes/azure.json"
touch "${AZURE_JSON_PATH}"
chmod 0600 "${AZURE_JSON_PATH}"
chown root:root "${AZURE_JSON_PATH}"
cat << EOF > "${AZURE_JSON_PATH}"
{
    "cloud":"${TARGET_ENVIRONMENT}",
    "tenantId": "${TENANT_ID}",
    "subscriptionId": "${SUBNETSCRIPTION_ID}",
    "aadClientId": "${SERVICE_PRINCIPAL_CLIENT_ID}",
    "aadClientSecret": "${SERVICE_PRINCIPAL_CLIENT_SECRET}",
    "resourceGroup": "${RESOURCE_GROUP}",
    "location": "${LOCATION}",
    "subnetName": "${SUBNET}",
    "securityGroupName": "${NETWORK_SECURITY_GROUP}",
    "vnetName": "${VIRTUAL_NETWORK}",
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
        if [ -e /usr/local/bin/kubectl ]
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
    downloadUrl https://acs-mirror.azureedge.net/cni/azure-vnet-cni-linux-amd64-latest.tgz | tar -xz -C $CNI_BIN_DIR
    # Mirror from https://github.com/containernetworking/cni/releases/download/$CNI_RELEASE_VER/cni-amd64-$CNI_RELEASE_VERSION.tgz
    downloadUrl https://acs-mirror.azureedge.net/cni/cni-amd64-latest.tgz | tar -xz -C $CNI_BIN_DIR ./loopback
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
        if [ -e /usr/local/bin/kubectl ]
        then
            /usr/local/bin/kubectl cluster-info
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
        curl --max-time 60 http://127.0.0.1:2379/v2/machines;
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

    FQDNSuffix="cloudapp.azure.com"
    if [ "$TARGET_ENVIRONMENT" = "AzureChinaCloud" ]
    then
        FQDNSuffix="cloudapp.chinacloudapi.cn"
    fi
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
ensureDocker
configNetworkPolicy
ensureKubelet
extractKubectl
ensureJournal

# master only
if [[ ! -z "${APISERVER_PRIVATE_KEY}" ]]; then
    writeKubeConfig
    ensureKubectl
    ensureEtcdDataDir
    ensureEtcd
    ensureApiserver
fi

# mitigation for bug https://bugs.launchpad.net/ubuntu/+source/linux/+bug/1676635
echo 2dd1ce17-079e-403c-b352-a1921ee207ee > /sys/bus/vmbus/drivers/hv_util/unbind
sed -i "13i\echo 2dd1ce17-079e-403c-b352-a1921ee207ee > /sys/bus/vmbus/drivers/hv_util/unbind\n" /etc/rc.local

# If APISERVER_PRIVATE_KEY is empty, then we are not on the master
echo "Install complete successfully"

if $REBOOTREQUIRED; then
  if [[ ! -z "${APISERVER_PRIVATE_KEY}" ]]; then
    # wait 1 minute to restart master
    echo 'reboot required, rebooting master in 1 minute'
    /bin/bash -c "shutdown -r 1 &"
  else
    echo 'reboot required, rebooting agent in 1 minute'
    shutdown -r now
  fi
fi

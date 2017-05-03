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

# Master only secrets
APISERVER_PRIVATE_KEY="${15}"
CA_CERTIFICATE="${16}"
MASTER_FQDN="${17}"
KUBECONFIG_CERTIFICATE="${18}"
KUBECONFIG_KEY="${19}"
ADMINUSER="${20}"

# If APISERVER_PRIVATE_KEY is empty, then we are not on the master
if [[ ! -z "${APISERVER_PRIVATE_KEY}" ]]; then
    echo "APISERVER_PRIVATE_KEY is non-empty, assuming master node"

    APISERVER_PRIVATE_KEY_PATH="/etc/kubernetes/certs/apiserver.key"
    touch "${APISERVER_PRIVATE_KEY_PATH}"
    chmod 0644 "${APISERVER_PRIVATE_KEY_PATH}"
    chown root:root "${APISERVER_PRIVATE_KEY_PATH}"
    echo "${APISERVER_PRIVATE_KEY}" | base64 --decode > "${APISERVER_PRIVATE_KEY_PATH}"

    AZURE_JSON_PATH="/etc/kubernetes/azure.json"
    touch "${AZURE_JSON_PATH}"
    chmod 0644 "${AZURE_JSON_PATH}"
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
    "primaryAvailabilitySetName": "${PRIMARY_AVAILABILITY_SET}"
}
EOF

else
    echo "APISERVER_PRIVATE_KEY is empty, assuming worker node"
fi

KUBELET_PRIVATE_KEY_PATH="/etc/kubernetes/certs/client.key"
touch "${KUBELET_PRIVATE_KEY_PATH}"
chmod 0644 "${KUBELET_PRIVATE_KEY_PATH}"
chown root:root "${KUBELET_PRIVATE_KEY_PATH}"
echo "${KUBELET_PRIVATE_KEY}" | base64 --decode > "${KUBELET_PRIVATE_KEY_PATH}"

###########################################################
# END OF SECRET DATA
###########################################################

set -x

# wait for kubectl to report successful cluster health
function ensureKubectl() {
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
    chmod 644 $CNI_CONFIG_DIR/10-azure.conf

    # Dump ebtables rules.
    /sbin/ebtables -t nat --list

    # Enable CNI.
    setNetworkPlugin cni
    setDockerOpts " --volume=/etc/cni/:/etc/cni:ro --volume=/opt/cni/:/opt/cni:ro"
}

function configCalicoNetworkPolicy() {
    if [[ ! -z "${APISERVER_PRIVATE_KEY}" ]]; then
        # on masters
        ADDONS="calico-configmap.yaml calico-daemonset.yaml"
        ADDONS_PATH=/etc/kubernetes/addons
        CALICO_URL="https://raw.githubusercontent.com/projectcalico/calico/a4ebfbad55ab1b7f10fdf3b39585471f8012e898/v2.0/getting-started/kubernetes/installation/hosted/k8s-backend-addon-manager"

        # download calico yamls
        for addon in ${ADDONS}; do
            downloadUrl "${CALICO_URL}/${addon}" > "${ADDONS_PATH}/${addon}"
        done
    else
        # on agents
        setNetworkPlugin cni
        setDockerOpts " --volume=/etc/cni/:/etc/cni:ro --volume=/opt/cni/:/opt/cni:ro"
    fi
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

function ensureDocker() {
    systemctl enable docker
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
        exit 1
    fi
}

function ensureKubelet() {
    systemctl enable kubelet
    systemctl restart kubelet
}

function extractKubectl(){
    systemctl enable kubectl-extract
    systemctl restart kubectl-extract
}

function ensureJournal(){
    systemctl daemon-reload
    systemctl enable systemd-journald.service
    systemctl restart systemd-journald.service
}

function ensureApiserver() {
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
        exit 1
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
    ensureEtcd
    ensureApiserver
fi

# If APISERVER_PRIVATE_KEY is empty, then we are not on the master
echo "Install complete successfully"


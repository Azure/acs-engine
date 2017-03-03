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

# Master only secrets
APISERVER_PRIVATE_KEY="${14}"
CA_CERTIFICATE="${15}"
MASTER_FQDN="${16}"
KUBECONFIG_CERTIFICATE="${17}"
KUBECONFIG_KEY="${18}"
ADMINUSER="${19}"

# If APISERVER_PRIVATE_KEY is empty, then we are not on the master
if [[ ! -z "${APISERVER_PRIVATE_KEY}" ]]; then
    echo "APISERVER_PRIVATE_KEY is non-empty, assuming master node"

    APISERVER_PRIVATE_KEY_PATH="/etc/kubernetes/certs/apiserver.key"
    touch "${APISERVER_PRIVATE_KEY_PATH}"
    chmod 0644 "${APISERVER_PRIVATE_KEY_PATH}"
    chown root:root "${APISERVER_PRIVATE_KEY_PATH}"
    echo "${APISERVER_PRIVATE_KEY}" | base64 --decode > "${APISERVER_PRIVATE_KEY_PATH}"
else
    echo "APISERVER_PRIVATE_KEY is empty, assuming worker node"
fi

KUBELET_PRIVATE_KEY_PATH="/etc/kubernetes/certs/client.key"
touch "${KUBELET_PRIVATE_KEY_PATH}"
chmod 0644 "${KUBELET_PRIVATE_KEY_PATH}"
chown root:root "${KUBELET_PRIVATE_KEY_PATH}"
echo "${KUBELET_PRIVATE_KEY}" | base64 --decode > "${KUBELET_PRIVATE_KEY_PATH}"

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

function ensureEtcd() {
    systemctl stop etcd
    rm -rf /var/lib/etcd/default
    systemctl restart etcd
}

function ensureDocker() {
    systemctl enable docker
    systemctl restart docker
    dockerStarted=1
    for i in {1..600}; do
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
ensureKubelet
extractKubectl

# master only 
if [[ ! -z "${APISERVER_PRIVATE_KEY}" ]]; then
    writeKubeConfig
    ensureKubectl
    ensureEtcd
    ensureApiserver
fi

# If APISERVER_PRIVATE_KEY is empty, then we are not on the master
echo "Install complete successfully"


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
SERVICE_PRINCIPAL_CLIENT_ID="${9}"
SERVICE_PRINCIPAL_CLIENT_SECRET="${10}"

# Extra secrets for Kubernetes biring-up
KUBELET_PRIVATE_KEY="${11}"
APISERVER_PRIVATE_KEY="${12}"

APISERVER_PRIVATE_KEY_PATH="/etc/kubernetes/certs/apiserver.key"
touch "${APISERVER_PRIVATE_KEY_PATH}"
chmod 0644 "${APISERVER_PRIVATE_KEY_PATH}"
chown root:root "${APISERVER_PRIVATE_KEY_PATH}"
echo "${APISERVER_PRIVATE_KEY}" | base64 --decode > "${APISERVER_PRIVATE_KEY_PATH}"

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
    "tenantId": "${TENANT_ID}",
    "subscriptionId": "${SUBNETSCRIPTION_ID}",
    "aadClientId": "${SERVICE_PRINCIPAL_CLIENT_ID}",
    "aadClientSecret": "${SERVICE_PRINCIPAL_CLIENT_SECRET}",
    "resourceGroup": "${RESOURCE_GROUP}",
    "location": "${LOCATION}",
    "subnetName": "${SUBNET}",
    "securityGroupName": "${NETWORK_SECURITY_GROUP}",
    "vnetName": "${VIRTUAL_NETWORK}",
    "routeTableName": "${ROUTE_TABLE}"
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
    systemctl restart etcd
}

function ensureDocker() {
    systemctl restart docker
    dockerStarted=1
    for i in {1..600}; do
        /usr/bin/docker ps 2>&1 | grep "daemon running"
        if [ "$?" = "0" ]
        then
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

function ensureKubernetes() {
    systemctl restart kubelet
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

ensureKubectl
ensureEtcd
ensureDocker
ensureKubernetes

echo "Install complete successfully"

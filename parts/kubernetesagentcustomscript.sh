#!/bin/bash

###########################################################
# START SECRET DATA - ECHO DISABLED
###########################################################

TID=$1
SID=$2
RGP=$3
LOC=$4
SUB=$5
NSG=$6
VNT=$7
RTB=$8
SVCPrincipalClientId=$9
SVCPrincipalClientSecret=${10}
CLIENTPRIVATEKEY=${11}

CLIENTKEY=/etc/kubernetes/certs/client.key
touch $CLIENTKEY
chmod 0644 $CLIENTKEY
chown root:root $CLIENTKEY
echo $CLIENTPRIVATEKEY | /usr/bin/base64 --decode > $CLIENTKEY 

AZUREJSON=/etc/kubernetes/azure.json
touch $AZUREJSON
chmod 0644 $AZUREJSON
chown root:root $AZUREJSON
AZURECONTENT=$(cat <<EOF
{
    "tenantId": "$TID",
    "subscriptionId": "$SID",
    "aadClientId": "$SVCPrincipalClientId",
    "aadClientSecret": "$SVCPrincipalClientSecret",
    "resourceGroup": "$RGP",
    "location": "$LOC",
    "subnetName": "$SUB",
    "securityGroupName": "$NSG",
    "vnetName": "$VNT",
    "routeTableName": "$RTB"
}
EOF
)
echo "$AZURECONTENT" > $AZUREJSON

###########################################################
# END OF SECRET DATA
###########################################################

set -x

# wait for docker to be available
ensuredockerbinary()
{
    dockerfound=1
    for i in {1..600}; do
        if [ -e /usr/bin/docker ]
        then
            dockerfound=0
            break
        fi
        sleep 1
    done
    if [ $dockerfound -ne 0 ]
    then
        echo "kubectl nor docker did not install successfully"
        exit 1
    fi
}
ensuredockerbinary

# start all the services
/bin/systemctl restart docker
ensureDocker()
{
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
ensureDocker

/bin/systemctl restart kubelet

ensureKubernetes()
{
    kubernetesStarted=1
    for i in {1..600}; do
        /usr/bin/docker ps | grep kubelet
        if [ "$?" = "0" ]
        then
            echo "kubernetes started"
            kubernetesStarted=0
            break
        else
            echo "kubernetes status $?"
        fi
        sleep 1
    done
    if [ $kubernetesStarted -ne 0 ]
    then
        echo "kubernetes did not start"
        exit 1
    fi
}
ensureKubernetes

echo "Install complete successfully"

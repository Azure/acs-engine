#!/bin/bash

###########################################################
# START SECRET DATA - ECHO DISABLED
###########################################################
APISERVERKEY=/etc/kubernetes/certs/apiserver.key
touch $APISERVERKEY
chmod 0644 $APISERVERKEY
chown root:root $APISERVERKEY
echo {{{apiServerPrivateKey}}} | /usr/bin/base64 --decode > $APISERVERKEY

CLIENTKEY=/etc/kubernetes/certs/client.key
touch $CLIENTKEY
chmod 0644 $CLIENTKEY
chown root:root $CLIENTKEY
echo {{{clientPrivateKey}}} | /usr/bin/base64 --decode > $CLIENTKEY 

AZUREJSON=/etc/kubernetes/azure.json
touch $AZUREJSON
chmod 0644 $AZUREJSON
chown root:root $AZUREJSON
AZURECONTENT=$(cat <<EOF
{
    "tenantId": "$TID",
    "subscriptionId": "$SID",
    "aadClientId": "{{{servicePrincipalClientId}}}",
    "aadClientSecret": "{{{servicePrincipalClientSecret}}}",
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

# wait for kubectl to report successful cluster health
ensurekubectl()
{
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
ensurekubectl

/bin/systemctl restart etcd

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
#wait for kubernetes to start 
ensureKubernetes()
{
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
ensureKubernetes

echo "Install complete successfully"

# Script file to install the docker hello-world container

#!/bin/bash

set -e

echo $(date) " - Starting Script"

echo $(date) " - Waiting for API Server to start"
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

# Deploy container
echo $(date) " - Deploying hello-world container"

kubectl run hello-world --quiet --image=busybox --restart=OnFailure -- echo "Hello Kubernetes!"

echo $(date) " - run kubectl get pods --show-all to list the pods"
echo $(date) " - run kubectl logs (passing the pod name gathered from kubectl get pods)"
echo $(date) " - Script complete"

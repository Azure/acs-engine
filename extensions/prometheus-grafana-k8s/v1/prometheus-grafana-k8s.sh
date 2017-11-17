# Script file to install prometheus and grafana

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
echo $(date) " - Downloading helm"

curl https://storage.googleapis.com/kubernetes-helm/helm-v2.6.2-linux-amd64.tar.gz > helm-v2.6.2-linux-amd64.tar.gz
tar -zxvf helm-v2.6.2-linux-amd64.tar.gz
mv linux-amd64/helm /usr/local/bin/helm
echo $(date) " - Downloading prometheus values"
curl https://raw.githubusercontent.com/ritazh/acs-engine/feat-monitor/extensions/prometheus-grafana-k8s/v1/prometheus_values.yaml > prometheus_values.yaml 
pwd

sleep 120

echo $(date) " - helm version"
helm version
helm init

echo $(date) " - helm installed"

echo $(date) " - Deploying prometheus chart"
helm install --name monitoring -f prometheus_values.yaml stable/prometheus

echo $(date) " - Deploying grafana chart"

helm install --name dashboard stable/grafana

echo $(date) " - installed grafana chart"
echo $(date) " - Script complete"

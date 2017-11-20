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

# TODO: replace this
curl https://raw.githubusercontent.com/ritazh/acs-engine/feat-monitor/extensions/prometheus-grafana-k8s/v1/prometheus_values.yaml > prometheus_values.yaml 
pwd

sleep 60

echo $(date) " - helm version"
helm version
helm init

echo $(date) " - helm installed"

NAMESPACE=default
K8S_SECRET_NAME=dashboard-grafana
DS_TYPE=prometheus
DS_NAME=prometheus1

PROM_URL=http://monitoring-prometheus-server

echo $(date) " - Deploying prometheus chart"
helm install --name monitoring -f prometheus_values.yaml stable/prometheus

echo $(date) " - Deploying grafana chart"

helm install --name dashboard stable/grafana

echo $(date) " - Installed grafana chart"

sleep 5

echo $(date) " - Creating the Prometheus datasource in Grafana"
GF_USER_NAME=$(kubectl get secret $K8S_SECRET_NAME -o jsonpath="{.data.grafana-admin-user}" | base64 --decode)
echo $GF_USER_NAME
GF_PASSWORD=$(kubectl get secret $K8S_SECRET_NAME -o jsonpath="{.data.grafana-admin-password}" | base64 --decode)
echo $GF_PASSWORD
GF_URL=$(kubectl get svc -l "app=dashboard-grafana,component=grafana" -o jsonpath="{.items[0].spec.clusterIP}")
echo $GF_URL

echo retrieving current data sources...
CURRENT_DS_LIST=$(curl -s --user "$GF_USER_NAME:$GF_PASSWORD" "$GF_URL/api/datasources")
echo $CURRENT_DS_LIST | grep -q "\"name\":\"$DS_NAME\""
if [[ $? -eq 0 ]]; then
    echo data source $DS_NAME already exists
    echo $CURRENT_DS_LIST | python -m json.tool
    exit 0
fi

echo data source $DS_NAME does not exist, creating...
DS_RAW=$(cat << EOF
{
    "name": "$DS_NAME",
    "type": "$DS_TYPE",
    "url": "$PROM_URL",
    "access": "proxy"
}
EOF
)

curl \
    -X POST \
    --user "$GF_USER_NAME:$GF_PASSWORD" \
    -H "Content-Type: application/json" \
    -d "$DS_RAW" \
    "$GF_URL/api/datasources"

echo $(date) " - Creating the Kubernetes dashboard in Grafana"

cat << EOF > sanitize_dashboard.py
#!/usr/bin/python3

import fileinput
import json

dashboard = json.loads(''.join(fileinput.input()))
dashboard.pop('__inputs')
dashboard.pop('__requires')
print(json.dumps(dashboard).replace('${DS_PROMETHEUS}', 'prometheus1'))

EOF

chmod u+x sanitize_dashboard.py

DB_RAW=$(cat << EOF
{
    "dashboard": $(curl -sL "https://grafana.com/api/dashboards/315/revisions/3/download" | ./sanitize_dashboard.py),
    "overwrite": false
}
EOF
)

curl \
    -X POST \
    --user "$GF_USER_NAME:$GF_PASSWORD" \
    -H "Content-Type: application/json" \
    -d "$DB_RAW" \
    "$GF_URL/api/dashboards/db"

echo $(date) " - Script complete"

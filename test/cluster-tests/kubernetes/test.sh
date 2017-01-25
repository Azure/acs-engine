#!/bin/bash

# exit on errors
set -e
# verbose logging
set -x

EXPECTED_NODE_COUNT=4
EXPECTED_DNS=2

namespace="namespace-${RANDOM}"
echo "Running test in namespace: ${namespace}"
trap teardown EXIT

function teardown {
  kubectl delete namespaces ${namespace}
}


echo "Testing number of nodes is ${EXPECTED_NODE_COUNT}"
node_count=$(kubectl get nodes --no-headers | wc | awk '{print $1}')
if [[ ${node_count} != ${EXPECTED_NODE_COUNT} ]]; then
  echo "Unexpected nodes: ${node_count}"
  kubectl get nodes
  exit 1
fi

echo "Testing system tools are running"
running=$(kubectl get pods --namespace=kube-system | grep kube-dns | grep Running | wc | awk '{print $1}')
if [[ ${running} != ${EXPECTED_DNS} ]]; then
  echo "Unexpected number of DNS servers: ${running}"
  kubectl get pods --namespace=kube-system
  exit 1
fi

echo "Testing proxies are running"
running=$(kubectl get pods --namespace=kube-system | grep kube-proxy | grep Running | wc | awk '{print $1}')
if [[ ${running} != ${EXPECTED_NODE_COUNT} ]]; then
  echo "Unexpected number of proxies running: ${running}"
  kubectl get pods --namespace=kube-system
  exit 1
fi

echo "Testing deployments"
kubectl create namespace ${namespace}

kubectl run --image=nginx nginx --namespace=${namespace}
count=0
while [[ ${count} < 10 ]]; do
  echo "Waiting for Pod to run"
  running=$(kubectl get pods --namespace=${namespace} | grep nginx | grep Running | wc | awk '{print $1}')
  if [[ ${running} == 1 ]]; then
    break
  fi
  count=(count+1)
  sleep 5
done

if [[ ${count} == 10 ]]; then
  echo "Deployment failed."
  kubectl get all --namespace=${namespace}
  exit 1
fi

kubectl expose deployments/nginx --namespace=${namespace} --port=80

# TODO actually check status here.
sleep 10

kubectl run busybox --image=busybox --attach=true --restart=Never --namespace=${namespace} -- wget nginx


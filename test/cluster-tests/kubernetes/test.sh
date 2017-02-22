#!/bin/bash

# exit on errors
set -e
# exit on unbound variables
set -u
# verbose logging
set -x

ENV_FILE="${CLUSTER_DEFINITION}.env"
if [ -e "${ENV_FILE}" ]; then
  source "${ENV_FILE}"
fi

EXPECTED_NODE_COUNT="${EXPECTED_NODE_COUNT:-4}"
EXPECTED_DNS="${EXPECTED_DNS:-2}"
EXPECTED_DASHBOARD="${EXPECTED_DASHBOARD:-1}"

# set TEST_ACR to "y" for ACR testing
TEST_ACR="${TEST_ACR:-n}"

namespace="namespace-${RANDOM}"
echo "Running test in namespace: ${namespace}"
trap teardown EXIT

function teardown {
  kubectl get all --all-namespaces
  kubectl get nodes
  kubectl get namespaces
  kubectl delete namespaces ${namespace}
}

# TODO: cleanup the loops more
# TODO: the wc|awk business can just be kubectl with an output format and wc -l

###### Deploy ACR
if [[ "${TEST_ACR}" == "y" ]]; then
	ACR_NAME="${INSTANCE_NAME//[-._]/}1"
	ACR_REGISTRY="${ACR_NAME}-microsoft.azurecr.io" # fix this for non-ms tenant users
	if ! az acr show --resource-group "${RESOURCE_GROUP}" --name "${ACR_NAME}" ; then
		az acr create --location "${LOCATION}" --resource-group "${RESOURCE_GROUP}" --name "${ACR_NAME}" &
	fi
fi

###### Check node count
function check_node_count() {
  wait=5
  count=12
  while (( $count > 0 )); do
    node_count=$(kubectl get nodes --no-headers | grep -v NotReady | grep Ready | wc | awk '{print $1}')
    if (( ${node_count} == ${EXPECTED_NODE_COUNT} )); then break; fi
    sleep 5; count=$((count-1))
  done
  if (( $node_count != ${EXPECTED_NODE_COUNT} )); then
    echo "gave up waiting for apiserver / node counts"; exit -1
  fi
}

check_node_count

###### Wait for no more container creating
wait=5
count=12
while (( $count > 0 )); do
  creating_count=$(kubectl get nodes --no-headers | grep 'CreatingContainer' | wc | awk '{print $1}')
  if (( ${creating_count} == 0 )); then break; fi
  sleep 5; count=$((count-1))
done
if (( ${creating_count} != 0 )); then
  echo "gave up waiting for creation to finish"; exit -1
fi


###### Check for Kube-DNS
wait=5
count=12
while (( $count > 0 )); do
  running=$(kubectl get pods --namespace=kube-system | grep kube-dns | grep Running | wc | awk '{print $1}')
  if (( ${running} == ${EXPECTED_DNS} )); then break; fi
  sleep 5; count=$((count-1))
done
if (( ${running} != ${EXPECTED_DNS} )); then
  echo "gave up waiting for kube-dns"; exit -1
fi

###### Check for Kube-Dashboard
wait=5
count=12
while (( $count > 0 )); do
  running=$(kubectl get pods --namespace=kube-system | grep kubernetes-dashboard | grep Running | wc | awk '{print $1}')
  if (( ${running} == ${EXPECTED_DASHBOARD} )); then break; fi
  sleep 5; count=$((count-1))
done
if (( ${running} != ${EXPECTED_DASHBOARD} )); then
  echo "gave up waiting for kubernetes-dashboard"; exit -1
fi

###### Check for Kube-Proxys
wait=5
count=12
while (( $count > 0 )); do
  nonrunning=$(kubectl get pods --namespace=kube-system | grep kube-proxy | grep -v Running | wc | awk '{print $1}')
  if (( ${nonrunning} == 0 )); then break; fi
  sleep 5; count=$((count-1))
done

# get master public hostname
master=$(kubectl config view | grep server | cut -f 3- -d "/" | tr -d " ")
# get dashboard port
port=$(kubectl get svc --namespace=kube-system | grep dashboard | awk '{print $4}' | sed -n 's/^80:\(.*\)\/TCP$/\1/p')
# get internal IPs of the nodes
ips=$(kubectl get nodes --all-namespaces -o yaml | grep -B 1 InternalIP | grep address | awk '{print $3}')

while read -r ip; do
  ssh -i "${OUTPUT}/id_rsa" -o "ConnectTimeout 3" -o "StrictHostKeyChecking no" -o "UserKnownHostsFile /dev/null" "azureuser@${master}" "curl http://${ip}:${port}"
done <<< "$ips"

###### Testing an nginx deployment
echo "Testing deployments"
kubectl create namespace ${namespace}

NGINX="docker.io/library/nginx:latest"
IMAGE="${NGINX}" # default to the library image unless we're in TEST_ACR mode
if [[ "${TEST_ACR}" == "y" ]]; then
	# force it to pull from ACR
	IMAGE="${ACR_REGISTRY}/test/nginx:latest"
	# wait for acr
	wait
	# TODO: how to do this without polluting user home dir?
	docker login --username="${SERVICE_PRINCIPAL_CLIENT_ID}" --password="${SERVICE_PRINCIPAL_CLIENT_SECRET}" "${ACR_REGISTRY}"
	docker pull "${NGINX}"
	docker tag "${NGINX}" "${IMAGE}"
	docker push "${IMAGE}"
fi

kubectl run --image="${IMAGE}" nginx --namespace=${namespace} --overrides='{ "apiVersion": "extensions/v1beta1", "spec":{"template":{"spec": {"nodeSelector":{"beta.kubernetes.io/os":"linux"}}}}}'
wait=5
count=12
while (( $count > 0 )); do
  running=$(kubectl get pods --namespace=${namespace} | grep nginx | grep Running | wc | awk '{print $1}')
  if (( ${running} == 1 )); then break; fi
  sleep 5; count=$((count-1))
done
if (( ${running} != 1 )); then
  echo "gave up waiting for deployment"
  kubectl get all --namespace=${namespace}
  exit -1
fi

kubectl expose deployments/nginx --type=LoadBalancer --namespace=${namespace} --port=80

wait=5
count=60
external_ip=""
while true; do
	external_ip=$(kubectl get svc --namespace ${namespace} nginx --template="{{range .status.loadBalancer.ingress}}{{.ip}}{{end}}")
	[[ ! -z "${external_ip}" ]] && break
	sleep 10
done
if [[ -z "${external_ip}" ]]; then
  echo "gave up waiting for loadbalancer to get an ingress ip"
  exit -1
fi

count=5
success="n"
while (( $count > 0 )); do
	curl -f "http://${external_ip}" | grep 'Welcome to nginx!'
	if [[ $? == 0 ]]; then
		success="y"
		break;
	fi
done
if [[ "${success}" != "y" ]]; then
  echo "failed to get expected response from nginx through the loadbalancer"
  exit -1
fi

check_node_count


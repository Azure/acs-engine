#!/bin/bash

####################################################
SOURCE="${BASH_SOURCE[0]}"
while [ -h "$SOURCE" ]; do # resolve $SOURCE until the file is no longer a symlink
  DIR="$( cd -P "$( dirname "$SOURCE" )" && pwd )"
  SOURCE="$(readlink "$SOURCE")"
  [[ $SOURCE != /* ]] && SOURCE="$DIR/$SOURCE" # if $SOURCE was a relative symlink, we need to resolve it relative to the path where the symlink file was located
done
DIR="$( cd -P "$( dirname "$SOURCE" )" && pwd )"
####################################################

# exit on errors
set -e
# exit on unbound variables
set -u
# verbose logging
set -x

source "$DIR/../utils.sh"

ENV_FILE="${CLUSTER_DEFINITION}.env"
if [ -e "${ENV_FILE}" ]; then
  source "${ENV_FILE}"
fi

EXPECTED_NODE_COUNT="${EXPECTED_NODE_COUNT:-4}"
EXPECTED_DNS="${EXPECTED_DNS:-2}"
EXPECTED_DASHBOARD="${EXPECTED_DASHBOARD:-1}"
EXPECTED_ORCHESTRATOR_VERSION="${EXPECTED_ORCHESTRATOR_VERSION:-}"

# set TEST_ACR to "y" for ACR testing
TEST_ACR="${TEST_ACR:-n}"

namespace="namespace-${RANDOM}"
log "Running test in namespace: ${namespace}"
trap teardown EXIT

function teardown {
  kubectl get all --all-namespaces || echo "teardown error"
  kubectl get nodes || echo "teardown error"
  kubectl get namespaces || echo "teardown error"
  kubectl delete namespaces ${namespace} || echo "teardown error"
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
  log "Checking node count"
  count=20
  while (( $count > 0 )); do
    log "  ... counting down $count"
    node_count=$(kubectl get nodes --no-headers | grep -v NotReady | grep Ready | wc | awk '{print $1}')
    if (( ${node_count} == ${EXPECTED_NODE_COUNT} )); then break; fi
    sleep 15; count=$((count-1))
  done
  if (( $node_count != ${EXPECTED_NODE_COUNT} )); then
    log "gave up waiting for apiserver / node counts"; exit -1
  fi
}

check_node_count

###### Validate Kubernetes version
log "Checking Kubernetes version. Expected: ${EXPECTED_ORCHESTRATOR_VERSION}"
if [ ! -z "${EXPECTED_ORCHESTRATOR_VERSION}" ]; then
  kubernetes_version=$(kubectl version --short)
  if [[ ${kubernetes_version} != *"Server Version: v${EXPECTED_ORCHESTRATOR_VERSION}"* ]]; then
    log "unexpected kubernetes version:\n${kubernetes_version}"; exit -1
  fi
fi

###### Wait for no more container creating
log "Checking containers being created"
count=12
while (( $count > 0 )); do
  log "  ... counting down $count"
  creating_count=$(kubectl get nodes --no-headers | grep 'CreatingContainer' | wc | awk '{print $1}')
  if (( ${creating_count} == 0 )); then break; fi
  sleep 5; count=$((count-1))
done
if (( ${creating_count} != 0 )); then
  log "gave up waiting for creation to finish"; exit -1
fi

###### Check existence and status of essential pods

# we test other essential pods (kube-dns, kube-proxy, kubernetes-dashboard) separately
pods="heapster kube-addon-manager kube-apiserver kube-controller-manager kube-scheduler tiller"
log "Checking $pods"

count=12
while (( $count > 0 )); do
  for pod in $pods; do
    running=$(kubectl get pods --all-namespaces | grep $pod | grep Running | wc -l)
    if (( $running > 0 )); then
      log "... $pod is Running"
      pods=$(echo $pods | sed -e "s/ *$pod */ /")
    fi
  done
  if [ -z "$(echo $pods | tr -d '[:space:]')" ]; then
    break
  fi
  sleep 5; count=$((count-1))
done

if [ ! -z "$(echo $pods | tr -d '[:space:]')" ]; then
  log "gave up waiting for running pods [$pods]"; exit -1
fi

###### Check for Kube-DNS
log "Checking Kube-DNS"
count=12
while (( $count > 0 )); do
  log "  ... counting down $count"
  running=$(kubectl get pods --namespace=kube-system | grep kube-dns | grep Running | wc | awk '{print $1}')
  if (( ${running} == ${EXPECTED_DNS} )); then break; fi
  sleep 5; count=$((count-1))
done
if (( ${running} != ${EXPECTED_DNS} )); then
  log "gave up waiting for kube-dns"; exit -1
fi

###### Check for Kube-Dashboard
log "Checking Kube-Dashboard"
count=12
while (( $count > 0 )); do
  log "  ... counting down $count"
  running=$(kubectl get pods --namespace=kube-system | grep kubernetes-dashboard | grep Running | wc | awk '{print $1}')
  if (( ${running} == ${EXPECTED_DASHBOARD} )); then break; fi
  sleep 5; count=$((count-1))
done
if (( ${running} != ${EXPECTED_DASHBOARD} )); then
  log "gave up waiting for kubernetes-dashboard"; exit -1
fi

###### Check for Kube-Proxys
log "Checking Kube-Proxys"
count=12
while (( $count > 0 )); do
  log "  ... counting down $count"
  running=$(kubectl get pods --namespace=kube-system | grep kube-proxy | grep Running | wc | awk '{print $1}')
  if (( ${running} == ${EXPECTED_NODE_COUNT} )); then break; fi
  sleep 5; count=$((count-1))
done

# get master public hostname
master=$(kubectl config view | grep server | cut -f 3- -d "/" | tr -d " ")
# get dashboard port
port=$(kubectl get svc --namespace=kube-system | grep dashboard | awk '{print $4}' | sed -n 's/^80:\(.*\)\/TCP$/\1/p')
# get internal IPs of the nodes
ips=$(kubectl get nodes --all-namespaces -o yaml | grep -B 1 InternalIP | grep address | awk '{print $3}')

for ip in $ips; do
  log "Probing IP address ${ip}"
  count=12
  success="n"
  while (( $count > 0 )); do
    log "  ... counting down $count"
    ret=$(ssh -i "${OUTPUT}/id_rsa" -o ConnectTimeout=30 -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null "azureuser@${master}" "curl --max-time 60 http://${ip}:${port}" || echo "curl_error")
    if [[ ! $ret =~ .*curl_error.* ]]; then
      success="y"
      break
    fi
    sleep 5; count=$((count-1))
  done
  if [[ "${success}" == "n" ]]; then
    log $ret; exit -1
  fi
done

###### Testing an nginx deployment
log "Testing deployments"
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
count=12
while (( $count > 0 )); do
  log "  ... counting down $count"
  running=$(kubectl get pods --namespace=${namespace} | grep nginx | grep Running | wc | awk '{print $1}')
  if (( ${running} == 1 )); then break; fi
  sleep 5; count=$((count-1))
done
if (( ${running} != 1 )); then
  log "gave up waiting for deployment"
  kubectl get all --namespace=${namespace}
  exit -1
fi

kubectl expose deployments/nginx --type=LoadBalancer --namespace=${namespace} --port=80

log "Checking Service External IP"
count=60
external_ip=""
while (( $count > 0 )); do
  log "  ... counting down $count"
	external_ip=$(kubectl get svc --namespace ${namespace} nginx --template="{{range .status.loadBalancer.ingress}}{{.ip}}{{end}}" || echo "")
	[[ ! -z "${external_ip}" ]] && break
	sleep 10; count=$((count-1))
done
if [[ -z "${external_ip}" ]]; then
  log "gave up waiting for loadbalancer to get an ingress ip"
  exit -1
fi

log "Checking Service"
count=5
success="n"
while (( $count > 0 )); do
  log "  ... counting down $count"
  ret=$(curl -f --max-time 60 "http://${external_ip}" | grep 'Welcome to nginx!' || echo "curl_error")
  if [[ $ret =~ .*'Welcome to nginx!'.* ]]; then
    success="y"
    break
	fi
  sleep 5; count=$((count-1))
done
if [[ "${success}" != "y" ]]; then
  log "failed to get expected response from nginx through the loadbalancer"
  exit -1
fi

check_node_count

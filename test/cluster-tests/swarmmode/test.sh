#!/bin/bash

set -x
set -e
set -u

ssh_args="-i ${SSH_KEY} -o StrictHostKeyChecking=no -p2200 azureuser@${INSTANCE_NAME}.${LOCATION}.cloudapp.azure.com"

function log {
    local message="$1"
    local caller="$(caller 0)"
	  now=$(date +"%D %T %Z")

	if [[ ! -z "${LOGFILE:-}" ]]; then
		echo "[${now}] [${caller}] ${message}" | tee -a ${LOGFILE}
	else
		echo "[${now}] [${caller}] ${message}"
    fi
}

function teardown {
  ssh ${ssh_args} docker service rm nginx || true
  sleep 10
  ssh ${ssh_args} docker network rm network || true
}

trap teardown EXIT

log "Starting swarmmode deployment validation in ${LOCATION}"
sleep 30
log "Creating network..."
wait=10
count=12
args="${ssh_args} docker network create --driver overlay --subnet 10.0.9.0/24 --opt encrypted network"
while (( $count > 0 )); do
  log "  ... counting down $count"
  ret=$(ssh $args || echo "ssh_error")
  if [[ "$ret" != "ssh_error" ]]; then break; fi
  sleep $wait
  count=$((count-1))
done
if [[ "$ret" == "ssh_error" ]]; then
  log "gave up waiting for network to be created"
  exit -1
fi

log "Creating service"
wait=5
count=12
args="${ssh_args} docker service create --replicas 3 --name nginx --network network --publish 80:80 nginx"
while (( $count > 0 )); do
  log "  ... counting down $count"
  ret=$(ssh $args || echo "ssh_error")
  if [[ "$ret" != "ssh_error" ]]; then break; fi
  sleep $wait
  count=$((count-1))
done
if [[ "$ret" == "ssh_error" ]]; then
  log "gave up waiting for service to be created"
  exit -1
fi

sleep 10
log "Testing service"
wait=5
count=12
while (( $count > 0 )); do
  log "  ... counting down $count"
  ret=$(curl --fail "http://${INSTANCE_NAME}0.${LOCATION}.cloudapp.azure.com:80/" || echo "curl_error")
  if [[ "$ret" != "curl_error" ]]; then break; fi
  sleep $wait
  count=$((count-1))
done
if [[ "$ret" == "curl_error" ]]; then
  log "gave up waiting for service to be externally reachable"
  exit -1
fi

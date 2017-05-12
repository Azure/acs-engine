#!/bin/bash

set -x
set -e
set -u

ssh_args="-i ${SSH_KEY} -o StrictHostKeyChecking=no -p2200 azureuser@${INSTANCE_NAME}.${LOCATION}.cloudapp.azure.com"

function teardown {
  ssh ${ssh_args} docker service rm nginx || true
  sleep 10
  ssh ${ssh_args} docker network rm network || true
}

trap teardown EXIT

sleep 30

wait=10
count=12
args="${ssh_args} docker network create --driver overlay --subnet 10.0.9.0/24 --opt encrypted network"
while (( $count > 0 )); do
  ret=$(ssh $args || echo "ssh_error")
  if [[ "$ret" != "ssh_error" ]]; then break; fi
  sleep $wait
  count=$((count-1))
done
if [[ "$ret" == "ssh_error" ]]; then
  echo "gave up waiting for network to be created"
  exit -1
fi

wait=5
count=12
args="${ssh_args} docker service create --replicas 3 --name nginx --network network --publish 80:80 nginx"
while (( $count > 0 )); do
  ret=$(ssh $args || echo "ssh_error")
  if [[ "$ret" != "ssh_error" ]]; then break; fi
  sleep $wait
  count=$((count-1))
done
if [[ "$ret" == "ssh_error" ]]; then
  echo "gave up waiting for service to be created"
  exit -1
fi

sleep 10

wait=5
count=12
while (( $count > 0 )); do
  ret=$(curl --fail "http://${INSTANCE_NAME}0.${LOCATION}.cloudapp.azure.com:80/" || echo "curl_error")
  if [[ "$ret" != "curl_error" ]]; then break; fi
  sleep $wait
  count=$((count-1))
done
if [[ "$ret" == "curl_error" ]]; then
  echo "gave up waiting for service to be externally reachable"
  exit -1
fi

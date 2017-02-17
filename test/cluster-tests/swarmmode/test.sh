#!/bin/bash

set -x
set -e
set -u

remote_exec="ssh -i ${SSH_KEY} -o StrictHostKeyChecking=no -p2200 azureuser@${INSTANCE_NAME}.${LOCATION}.cloudapp.azure.com"

function teardown {
  ${remote_exec} docker service rm nginx || true
  sleep 10
  ${remote_exec} docker network rm network || true
}

trap teardown EXIT

# TODO: investigate if this is really needed - replace with deterministic wait or retry loop
sleep 60

${remote_exec} docker network create \
	--driver overlay \
	--subnet 10.0.9.0/24 \
	--opt encrypted \
	network

${remote_exec} docker service create \
	--replicas 3 \
	--name nginx \
	--network network \
	--publish 80:80 \
	nginx

sleep 10

wait=5
count=12
success="n"
while (( $count > 0 )); do
  curl --fail "http://${INSTANCE_NAME}0.${LOCATION}.cloudapp.azure.com:80/"
  if [[ $? == 0 ]]; then
    success="y"
    break
  fi
done
if [[ "${success}" != "y" ]]; then
  echo "gave up waiting for service to be externally reachable"; exit -1
fi

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

# do not use 'set -e'
set -x
set -u


source "$DIR/../utils.sh"

ssh_args="-i ${SSH_KEY} -o ConnectTimeout=30 -o StrictHostKeyChecking=no -p2200 azureuser@${INSTANCE_NAME}.${LOCATION}.cloudapp.azure.com"

function teardown {
  ssh ${ssh_args} docker service rm nginx
  sleep 10
  ssh ${ssh_args} docker network rm network
}

trap teardown EXIT

log "Starting swarmmode deployment validation in ${LOCATION}"
sleep 30
log "Creating network"
wait=10
count=12
args="${ssh_args} docker network create --driver overlay --subnet 10.0.9.0/24 --opt encrypted network"
while (( $count > 0 )); do
  log "  ... counting down $count"
  timeout -k 60s -s KILL 60s ssh $args
  retval=$?
  if [[ "$retval" == "0" ]]; then break; fi
  sleep $wait
  count=$((count-1))
done
if [[ "$retval" != "0" ]]; then
  log "DockerCE: gave up waiting for network to be created"
  exit 1
fi

log "Creating service"
wait=5
count=12
args="${ssh_args} docker service create --replicas 3 --name nginx --network network --publish 80:80 nginx"
while (( $count > 0 )); do
  log "  ... counting down $count"
  ssh $args
  retval=$?
  if [[ "$retval" == "0" ]]; then break; fi
  sleep $wait
  count=$((count-1))
done
if [[ "$retval" != "0" ]]; then
  log "DockerCE: gave up waiting for service to be created"
  exit 1
fi

sleep 10
log "Testing service"
wait=5
count=12
while (( $count > 0 )); do
  log "  ... counting down $count"
  curl --fail "http://${INSTANCE_NAME}0.${LOCATION}.cloudapp.azure.com:80/"
  retval=$?
  if [[ "$retval" == "0" ]]; then break; fi
  sleep $wait
  count=$((count-1))
done
if [[ "$retval" != "0" ]]; then
  log "DockerCE: gave up waiting for service to be externally reachable"
  exit 1
fi

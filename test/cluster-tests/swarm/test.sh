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

function assert {
# $1 is the string to describe the value
# $2 is the actual value
# $3 is the expected value

  if [ $# -ne 3 ]; then
    log "$@" >&2
    exit 1
  fi

  if [ "$2" != "$3" ]; then
    log "[Actual] $1 $2, [Expected] $3" >&2
    exit 1
  fi
}

SSH="ssh -i ${SSH_KEY} -o ConnectTimeout=30 -o StrictHostKeyChecking=no -p2200 azureuser@${INSTANCE_NAME}.${LOCATION}.cloudapp.azure.com"

deploy="docker -H :2375 run -d -p 80:80 yeasy/simple-web"
wait_duration=10
total_loops=30
while true; do
  # || true is used to suppress the failure like "Error response from daemon: No elected primary cluster manager"
  # it should be gone after a few retries
  containerId="$($SSH $deploy 2>/dev/null )" || true
  [[ ! -z $containerId ]] && [[ "$(echo $containerId | grep '[0-9a-z]\{64\}')" ]] && log "container deployed! containerId is $containerId" && break
  log "Validation: Expected to get containerId. $(($total_loops*$wait_duration)) seconds remain"
  sleep $wait_duration
  total_loops=$((total_loops-1))
  if [ $total_loops -eq 0 ]; then 
      log "swarm Validation Failed, timeout"; exit 1;
  fi
done

result=$($SSH curl localhost:2375/containers/json)
assert "deployed container should be found." \
  "$(echo $result | jq "[.[].Id==\"$containerId\"] | any")" \
  "true" 

log "Finished swarm_deep_test"

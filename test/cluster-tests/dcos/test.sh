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

set -e
# -o pipefail
set -x

source "$DIR/../utils.sh"

remote_exec="ssh -i "${SSH_KEY}" -o ConnectTimeout=30 -o StrictHostKeyChecking=no azureuser@${INSTANCE_NAME}.${LOCATION}.cloudapp.azure.com -p2200"
agentFQDN="dcos-agent-ip-${INSTANCE_NAME}.${LOCATION}.cloudapp.azure.com"
remote_cp="scp -i "${SSH_KEY}" -P 2200 -o StrictHostKeyChecking=no"

function teardown {
  ${remote_exec} ./dcos marathon app remove /web
}

${remote_exec} curl -O https://downloads.dcos.io/binaries/cli/linux/x86-64/dcos-1.8/dcos
${remote_exec} chmod a+x ./dcos
${remote_exec} ./dcos config set core.dcos_url http://localhost:80

${remote_cp} "${DIR}/marathon.json" azureuser@${INSTANCE_NAME}.${LOCATION}.cloudapp.azure.com:marathon.json

# feed agentFQDN to marathon.json
${remote_exec} sed -i "s/{agentFQDN}/${agentFQDN}/g" marathon.json
${remote_exec} ./dcos marathon app add marathon.json

# only need to teardown if app added successfully
trap teardown EXIT

log "Validating app"
count=0
while [[ ${count} -lt 25 ]]; do
  count=$((count+1))
  log "  ... cycle $count"
  running=$(${remote_exec} ./dcos marathon app show /web | jq .tasksRunning)
  if [[ "${running}" == "3" ]]; then
    log "Found 3 running tasks"
    break
  fi
  sleep ${count}
done

if [[ "${running}" != "3" ]]; then
  log "Validation failed"
  ${remote_exec} ./dcos marathon app show /web
  exit 1
fi

# install marathon-lb
${remote_exec} ./dcos package install marathon-lb --yes

# curl simpleweb through external haproxy
log "Checking Service"
count=10
while (( $count > 0 )); do
  log "  ... counting down $count"
  [[ $(curl -sI --max-time 60 "http://${agentFQDN}" |head -n1 |cut -d$' ' -f2) -eq "200" ]] && echo "Successfully hitting simpleweb through external haproxy http://${agentFQDN}" && break
  if [[ "${count}" -le 0 ]]; then
    log "failed to get expected response from nginx through the loadbalancer"
    exit 1
  fi
  sleep 5; count=$((count-1))
done

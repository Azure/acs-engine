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
# -o pipefail
set -x

source "$DIR/../utils.sh"

remote_exec="ssh -i "${SSH_KEY}" -o ConnectTimeout=30 -o StrictHostKeyChecking=no azureuser@${INSTANCE_NAME}.${LOCATION}.cloudapp.azure.com -p2200"
agentFQDN="${INSTANCE_NAME}0.${LOCATION}.cloudapp.azure.com"
remote_cp="scp -i "${SSH_KEY}" -P 2200 -o StrictHostKeyChecking=no"

function teardown {
  ${remote_exec} ./dcos marathon app remove /web
}

log "Downloading dcos"
${remote_exec} curl -O https://downloads.dcos.io/binaries/cli/linux/x86-64/dcos-1.8/dcos || (log "Failed to download dcos"; exit 1)
log "Setting dcos permissions"
${remote_exec} chmod a+x ./dcos || (log "Failed to chmod dcos"; exit 1)
log "Configuring dcos"
${remote_exec} ./dcos config set core.dcos_url http://localhost:80 || (log "Failed to configure dcos"; exit 1)

log "Copying marathon.json"
${remote_cp} "${DIR}/marathon.json" azureuser@${INSTANCE_NAME}.${LOCATION}.cloudapp.azure.com:marathon.json || (log "Failed to copy marathon.json"; exit 1)

# feed agentFQDN to marathon.json
log "Configuring marathon.json"
${remote_exec} sed -i "s/{agentFQDN}/${agentFQDN}/g" marathon.json || (log "Failed to configure marathon.json"; exit 1)


log "Adding marathon app"
count=20
while (( $count > 0 )); do
  log "  ... counting down $count"
  ${remote_exec} ./dcos marathon app list | grep /web
  retval=$?
  if [[ $retval -eq 0 ]]; then log "Marathon App successfully installed" && break; fi
  ${remote_exec} ./dcos marathon app add marathon.json
  retval=$?
  if [[ "$retval" == "0" ]]; then break; fi
  sleep 15; count=$((count-1))
done
if [[ $retval -ne 0 ]]; then
  log "gave up waiting for marathon to be added"
  exit -1
fi

# only need to teardown if app added successfully
#trap teardown EXIT

log "Validating marathon app"
count=0
while [[ ${count} -lt 25 ]]; do
  count=$((count+1))
  log "  ... cycle $count"
  running=$(${remote_exec} ./dcos marathon app show /web | jq .tasksRunning)
  if [[ "${running}" == "3" ]]; then
    log "Found 3 running tasks"
    break
  else
    ${remote_exec} ./dcos marathon app list
  fi
  sleep ${count}
done

if [[ "${running}" != "3" ]]; then
  log "marathon validation failed"
  ${remote_exec} ./dcos marathon app show /web
  exit 1
fi

# install marathon-lb
${remote_exec} ./dcos package install marathon-lb --yes || (log "Failed to install marathon-lb"; exit 1)

# curl simpleweb through external haproxy
log "Checking Service"
count=10
while true; do
  log "  ... counting down $count"
  rc=$(curl -sI --max-time 60 "http://${agentFQDN}" | head -n1 | cut -d$' ' -f2)
  [[ "$rc" -eq "200" ]] && log "Successfully hitting simpleweb through external haproxy http://${agentFQDN}" && break
  if [[ "${count}" -le 1 ]]; then
    log "failed to get expected response from nginx through the loadbalancer: Error $rc"
    exit 1
  fi
  sleep 5; count=$((count-1))
done

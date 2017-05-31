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

ENV_FILE="${CLUSTER_DEFINITION}.env"
if [ -e "${ENV_FILE}" ]; then
  source "${ENV_FILE}"
fi

MARATHON_JSON="${MARATHON_JSON:-marathon.json}"
FQDNSuffix="cloudapp.azure.com"
if [ "$TARGET_ENVIRONMENT" = "AzureChinaCloud" ]; then
    FQDNSuffix="cloudapp.chinacloudapi.cn"
fi
remote_exec="ssh -i "${SSH_KEY}" -o ConnectTimeout=30 -o StrictHostKeyChecking=no azureuser@${INSTANCE_NAME}.${LOCATION}.${FQDNSuffix} -p2200"
agentFQDN="${INSTANCE_NAME}0.${LOCATION}.${FQDNSuffix}"
remote_cp="scp -i "${SSH_KEY}" -P 2200 -o StrictHostKeyChecking=no"

function teardown {
  ${remote_exec} ./dcos marathon app remove /web
}

###### Check node count
function check_node_count() {
  log "Checking node count"
  count=20
  while (( $count > 0 )); do
    log "  ... counting down $count"
    node_count=$(${remote_exec} curl -s http://localhost:1050/system/health/v1/nodes | jq '.nodes | length')
    [ $? -eq 0 ] && [ ! -z "$node_count" ] && [ $node_count -eq ${EXPECTED_NODE_COUNT} ] && log "Successfully got $EXPECTED_NODE_COUNT nodes" && break
    sleep 30; count=$((count-1))
  done
  if (( $node_count != ${EXPECTED_NODE_COUNT} )); then
    log "gave up waiting for DCOS nodes: $node_count available, ${EXPECTED_NODE_COUNT} expected"
    exit 1
  fi
}

check_node_count

log "Downloading dcos"
${remote_exec} curl -O https://downloads.dcos.io/binaries/cli/linux/x86-64/dcos-1.8/dcos
if [[ "$?" != "0" ]]; then log "Failed to download dcos"; exit 1; fi
log "Setting dcos permissions"
${remote_exec} chmod a+x ./dcos
if [[ "$?" != "0" ]]; then log "Failed to chmod dcos"; exit 1; fi
log "Configuring dcos"
${remote_exec} ./dcos config set core.dcos_url http://localhost:80
if [[ "$?" != "0" ]]; then log "Failed to configure dcos"; exit 1; fi

log "Copying marathon.json"

${remote_cp} "${DIR}/${MARATHON_JSON}" azureuser@${INSTANCE_NAME}.${LOCATION}.${FQDNSuffix}:marathon.json
if [[ "$?" != "0" ]]; then log "Failed to copy marathon.json"; exit 1; fi

# feed agentFQDN to marathon.json
log "Configuring marathon.json"
${remote_exec} sed -i "s/{agentFQDN}/${agentFQDN}/g" marathon.json
if [[ "$?" != "0" ]]; then log "Failed to configure marathon.json"; exit 1; fi


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
if [[ $retval -ne 0 ]]; then log "gave up waiting for marathon to be added"; exit 1; fi

# only need to teardown if app added successfully
trap teardown EXIT

log "Validating marathon app"
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
  log "marathon validation failed"
  ${remote_exec} ./dcos marathon app show /web
  ${remote_exec} ./dcos marathon app list
  exit 1
fi

# install marathon-lb
${remote_exec} ./dcos package install marathon-lb --yes
if [[ "$?" != "0" ]]; then log "Failed to install marathon-lb"; exit 1; fi

# curl simpleweb through external haproxy
log "Checking Service"
count=20
while true; do
  log "  ... counting down $count"
  rc=$(curl -sI --max-time 60 "http://${agentFQDN}" | head -n1 | cut -d$' ' -f2)
  [[ "$rc" -eq "200" ]] && log "Successfully hitting simpleweb through external haproxy http://${agentFQDN}" && break
  if [[ "${count}" -le 1 ]]; then
    log "failed to get expected response from nginx through the loadbalancer: Error $rc"
    exit 1
  fi
  sleep 15; count=$((count-1))
done

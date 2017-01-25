#!/bin/bash

set -x
set -e

usage() { echo "Usage: $0 [-h <hostname>] [-u <username>]" 1>&2; exit 1; }

while getopts ":h:u:" o; do
    case "${o}" in
        h)
            host=${OPTARG}
            ;;
        u)
            user=${OPTARG}
            ;;
        *)
            usage
            ;;
    esac
done
shift $((OPTIND-1))

if [[ ! -z $1 ]]; then
  usage
fi

if [[ -z $host ]]; then
  host=$(az acs show --resource-group=acs-weekly-dcos --name=weekly-test --query=masterProfile.fqdn | sed -e 's/^"//' -e 's/"$//')
fi

if [[ -z $user ]]; then
  user=$(az acs show --resource-group=acs-weekly-dcos --name=weekly-test --query=linuxProfile.adminUsername | sed -e 's/^"//' -e 's/"$//')
fi

echo $host

remote_exec="ssh -i ~/.ssh/id_rsa ${user}@${host}"

function teardown {
  ${remote_exec} dcos marathon app remove /web
}

scp -i ~/.ssh/id_rsa ${HOME}/marathon.json ${user}@${host}:marathon.json

trap teardown EXIT

${remote_exec} dcos marathon app add marathon.json

count=0
while [[ ${count} < 10 ]]; do
  count=(count + 1)
  running=$(${remote_exec} dcos marathon app show /web | jq .tasksRunning)
  if [[ "${running}" == "3" ]]; then
    echo "Found 3 running tasks"
    break
  fi
  sleep ${count}
done

if [[ "${running}" == "3" ]]; then
  echo "Deployment succeeded"
else
  echo "Deployment failed"
  exit 1
fi

#!/usr/bin/env bash

SOURCE="${BASH_SOURCE[0]}"
while [ -h "$SOURCE" ]; do # resolve $SOURCE until the file is no longer a symlink
  DIR="$( cd -P "$( dirname "$SOURCE" )" && pwd )"
  SOURCE="$(readlink "$SOURCE")"
  [[ $SOURCE != /* ]] && SOURCE="$DIR/$SOURCE" # if $SOURCE was a relative symlink, we need to resolve it relative to the path where the symlink file was located
done
DIR="$( cd -P "$( dirname "$SOURCE" )" && pwd )"

###############################################################################

set -e
set -u
set -o pipefail

ROOT="${DIR}/.."

source "${ROOT}/test/common.sh"

case $1 in

set_azure_account)
  set_azure_account
;;

create_resource_group)
  create_resource_group
;;

predeploy)
  ACSE_PREDEPLOY=${ACSE_PREDEPLOY:-}
  if [ ! -z "${ACSE_PREDEPLOY}" ] && [ -x "${ACSE_PREDEPLOY}" ]; then
      "${ACSE_PREDEPLOY}"
  fi
;;

postdeploy)
  ACSE_POSTDEPLOY=${ACSE_POSTDEPLOY:-}
  if [ ! -z "${ACSE_POSTDEPLOY}" ] && [ -x "${ACSE_POSTDEPLOY}" ]; then
      "${ACSE_POSTDEPLOY}"
  fi
;;

generate_template)
  export OUTPUT="${ROOT}/_output/${INSTANCE_NAME}"
  generate_template
;;

deploy_template)
  export OUTPUT="${ROOT}/_output/${INSTANCE_NAME}"
  deploy_template
;;

get_node_count)
  export OUTPUT="${ROOT}/_output/${INSTANCE_NAME}"
  get_node_count
;;

get_orchestrator_type)
  get_orchestrator_type
;;

get_orchestrator_release)
  get_orchestrator_release
;;

validate)
  export OUTPUT="${ROOT}/_output/${INSTANCE_NAME}"
  export SSH_KEY="${OUTPUT}/id_rsa"
  if [ ${ORCHESTRATOR} = "kubernetes" ]; then
    export KUBECONFIG="${OUTPUT}/kubeconfig/kubeconfig.${LOCATION}.json"
  fi
  "${ROOT}/test/cluster-tests/${ORCHESTRATOR}/test.sh"
;;

cleanup)
  export CLEANUP="${CLEANUP:-y}"
  cleanup
;;
esac

#!/usr/bin/env bash

source "${HOME}/test/common.sh"

function shunittest_validate_deployment {
  set -eu -o pipefail

  export OUTPUT="${HOME}/_output/${INSTANCE_NAME}"
  export SSH_KEY="${OUTPUT}/id_rsa"
  if [ ${ORCHESTRATOR} = "kubernetes" ]; then
    export KUBECONFIG="${OUTPUT}/kubeconfig/kubeconfig.${LOCATION}.json"
  fi

  "${HOME}/test/cluster-tests/${ORCHESTRATOR}/test.sh"
}

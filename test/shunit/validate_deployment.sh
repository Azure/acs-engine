#!/usr/bin/env bash

function shunittest_validate_deployment {
  set -eux -o pipefail

  export OUTPUT="${HOME}/_output/${INSTANCE_NAME}"
  export SSH_KEY="${OUTPUT}/id_rsa"
  if [[ "${ORCHESTRATOR}" == "kubernetes" ]]; then
    export KUBECONFIG="${OUTPUT}/kubeconfig/kubeconfig.${LOCATION}.json"
    export EXPECTED_NODE_COUNT=$(${HOME}/test/step.sh get_node_count)
  fi

  script="${HOME}/test/cluster-tests/${ORCHESTRATOR}/test.sh"

  if [ -x "${script}" ]; then
    "${script}"
  else
    echo "${script}: not an executable or no such file"
    exit 1
  fi
}

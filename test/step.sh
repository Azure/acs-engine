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

# Load any user set environment
if [[ -f "${ROOT}/test/user.env" ]]; then
	source "${ROOT}/test/user.env"
fi

# Set output directory
export OUTPUT="${ROOT}/_output/${INSTANCE_NAME}"

source "${ROOT}/test/common.sh"

case $1 in
generate_template)
  generate_template
;;

deploy_template)
#  trap cleanup EXIT
  deploy_template
;;

verify)
  export SSH_KEY="${OUTPUT}/id_rsa"
  export KUBECONFIG="${OUTPUT}/kubeconfig/kubeconfig.${LOCATION}.json"
  "${ROOT}/test/cluster-tests/${ORCHESTRATOR}/test.sh"
;;

cleanup)
  cleanup
;;
esac

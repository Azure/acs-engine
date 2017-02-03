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

function cleanup {
	if [[ "${CLEANUP:-}" == "y" ]]; then
		az group delete --no-wait --name "${INSTANCE_NAME}" || true
	fi
}

# Usage:
#
# Manual user usage (Specific name):
#   export INSTANCE_NAME=official-jenkins-infra
#   ./scripts/deploy.sh ./examples/kubernetes.json
#
# Manual user usage (Lots of rapid fire):
# In this mode, the user can repeat the same deploy
# command blindly and get new clusters each time.
#   unset INSTANCE_NAME
#   vim ./test/user.env (add your stuff)
#   ./scripts/deploy.sh ./examples.kubernetes.json
#   sleep 1
#   ./scripts/deploy.sh ./examples.kubernetes.json
#
# Prow:
#   export PULL_NUMBER=...
#   export VALIDATE=true
#   export CLUSTER_DEFIITION=examples/kubernetes.json
#   export CLUSTER_TYPE=kubernetes
#   ./scripts/deploy.sh

# Load any user set environment
if [[ -f "${ROOT}/test/user.env" ]]; then
	source "${ROOT}/test/user.env"
fi


# Ensure Cluster Definition
if [[ -z "${CLUSTER_DEFINITION:-}" ]]; then
	if [[ -z "${1:-}" ]]; then echo "You must specify a parameterized apimodel.json clusterdefinition"; exit -1; fi
	CLUSTER_DEFINITION="${1}"
fi

# Set Instance Name for PR or random run
if [[ ! -z "${PULL_NUMBER:-}" ]]; then
	export INSTANCE_NAME="${JOB_NAME}-${PULL_NUMBER}-$(printf "%x" $(date '+%s'))"
	# if we're running a pull request, assume we want to cleanup unless the user specified otherwise
	if [[ -z "${CLEANUP:-}" ]]; then
		export CLEANUP="y"
	fi
else
	export INSTANCE_NAME_DEFAULT="${INSTANCE_NAME_PREFIX}-$(printf "%x" $(date '+%s'))"
	export INSTANCE_NAME="${INSTANCE_NAME:-${INSTANCE_NAME_DEFAULT}}"
fi

make -C "${ROOT}"
trap cleanup EXIT
deploy

if [[ -z "${VALIDATE:-}" ]]; then
	exit 0
fi

export SSH_KEY="${ROOT}/_output/${INSTANCE_NAME}/id_rsa"
export KUBECONFIG="${ROOT}/_output/${INSTANCE_NAME}/kubeconfig/kubeconfig.${LOCATION}.json"

"${ROOT}/${VALIDATE}"

echo "post-test..."

# TODO: this shouldn't be necessary but this trap doesn't seem to fire
# so... manually call it (of course, this only works in the happy path)
cleanup

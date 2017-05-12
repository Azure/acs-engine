#!/usr/bin/env bash

SOURCE="${BASH_SOURCE[0]}"
while [ -h "$SOURCE" ]; do # resolve $SOURCE until the file is no longer a symlink
  DIR="$( cd -P "$( dirname "$SOURCE" )" && pwd )"
  SOURCE="$(readlink "$SOURCE")"
  [[ $SOURCE != /* ]] && SOURCE="$DIR/$SOURCE" # if $SOURCE was a relative symlink, we need to resolve it relative to the path where the symlink file was located
done
DIR="$( cd -P "$( dirname "$SOURCE" )" && pwd )"

###############################################################################

set -x
set -e
set -u
set -o pipefail

ROOT="${DIR}/.."

# Load any user set environment
if [[ -f "${ROOT}/test/user.env" ]]; then
	source "${ROOT}/test/user.env"
fi

# ./test/deploy.sh  ./examples/kubernetesversions/kubernetes1.5.3.json
# ./test/upgrade.sh [deploy-dir] [rg] [upgrademodel-file]
# ./test/upgrade.sh ./_output/colemick-<TAB> colemick-<SAME VAL> ./examples/operations/upgrade/kubernetes-1.6.2.json

# device auth:
# ./test/upgrade.sh ./_output/colemick-<TAB> colemick-<SAME VAL> ./examples/operations/upgrade/kubernetes-1.6.2.json device

API_MODEL="${1}"
RESOURCE_GROUP="${2}"

AUTH_METHOD="${3:-device}"

# make -C "${ROOT}" ci # this ensure you didn't accidentally forget to regen
make -C "${ROOT}" build
"${ROOT}/acs-engine" generate \
  --subscription-id="${SUBSCRIPTION_ID}" \
  --auth-method="${AUTH_METHOD}" \
  --client-id="${SERVICE_PRINCIPAL_CLIENT_ID}" \
  --client-secret "${SERVICE_PRINCIPAL_CLIENT_SECRET}" \
  --resource-group="${RESOURCE_GROUP}" \
  --deploy \
  --debug \
  "${API_MODEL}"

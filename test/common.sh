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

set -eu -o pipefail
set -x

# see: https://github.com/stedolan/jq/issues/105
# and: https://github.com/stedolan/jq/wiki/FAQ#general-questions
function jqi() {
	filename="${1}"
	jqexpr="${2}"
	jq "${jqexpr}" "${filename}" > "${filename}.tmp" && mv "${filename}.tmp" "${filename}"
}

function deploy() {
	# Check pre-requisites
	[[ ! -z "${INSTANCE_NAME:-}" ]] || (echo "Must specify INSTANCE_NAME" && exit -1)
	[[ ! -z "${LOCATION:-}" ]] || (echo "Must specify LOCATION" && exit -1)
	[[ ! -z "${CLUSTER_DEFINITION:-}" ]] || (echo "Must specify CLUSTER_DEFINITION" && exit -1)
	[[ ! -z "${SERVICE_PRINCIPAL_CLIENT_ID:-}" ]] || (echo "Must specify SERVICE_PRINCIPAL_CLIENT_ID" && exit -1)
	[[ ! -z "${SERVICE_PRINCIPAL_CLIENT_SECRET:-}" ]] || (echo "Must specify SERVICE_PRINCIPAL_CLIENT_SECRET" && exit -1)
	which kubectl || (echo "kubectl must be on PATH" && exit -1)
	which az || (echo "az must be on PATH" && exit -1)
	
	# Set output directory
	export OUTPUT="$(pwd)/_output/${INSTANCE_NAME}"
	mkdir -p "${OUTPUT}"

	# Prep SSH Key
	# (can't use ssh-keygen, no user info inside Jenkins build container env)
	openssl genpkey -algorithm RSA -out "${OUTPUT}/id_rsa" -pkeyopt rsa_keygen_bits:2048
	echo -n "ssh-rsa " > "${OUTPUT}/id_rsa.pub"
	grep -v -- ----- "${OUTPUT}/id_rsa" | base64 -d | dd bs=1 skip=32 count=257 status=none | xxd -p -c257 | sed s/^/00000007\ 7373682d727361\ 00000003\ 010001\ 00000101\ / | xxd -p -r | base64 -w0 >> "${OUTPUT}/id_rsa.pub"
	echo >> "${OUTPUT}/id_rsa.pub"
	export SSH_KEY_DATA="$(cat "${OUTPUT}/id_rsa.pub")"

	# Form the final cluster_definition file
	export FINAL_CLUSTER_DEFINITION="${OUTPUT}/clusterdefinition.json"
	cp "${CLUSTER_DEFINITION}" "${FINAL_CLUSTER_DEFINITION}"
	jqi "${FINAL_CLUSTER_DEFINITION}" ".properties.masterProfile.dnsPrefix = \"${INSTANCE_NAME}\""
	jqi "${FINAL_CLUSTER_DEFINITION}" ".properties.linuxProfile.ssh.publicKeys[0].keyData = \"${SSH_KEY_DATA}\"" t "${FINAL_CLUSTER_DEFINITION}" 
	jqi "${FINAL_CLUSTER_DEFINITION}" ".properties.servicePrincipalProfile.servicePrincipalClientID = \"${SERVICE_PRINCIPAL_CLIENT_ID}\""
	jqi "${FINAL_CLUSTER_DEFINITION}" ".properties.servicePrincipalProfile.servicePrincipalClientSecret = \"${SERVICE_PRINCIPAL_CLIENT_SECRET}\""

	# Generate template
	"${DIR}/../acs-engine" -artifacts "${OUTPUT}" "${FINAL_CLUSTER_DEFINITION}"

	# Fill in custom hyperkube spec, if it was set
	if [[ ! -z "${CUSTOM_HYPERKUBE_SPEC:-}" ]]; then
		jqi "${OUTPUT}/azuredeploy.parameters.json" ".kubernetesHyperkubeSpec.value = \"${CUSTOM_HYPERKUBE_SPEC}\""
	fi

	# Login to Azure-Cli
	az login --service-principal \
		--username "${SERVICE_PRINCIPAL_CLIENT_ID}" \
		--password "${SERVICE_PRINCIPAL_CLIENT_SECRET}" \
		--tenant "${TENANT_ID}"

	az account set --name "${SUBSCRIPTION_ID}"

	# Deploy the template
	az resource group create --name="${INSTANCE_NAME}" --location="${LOCATION}"
	sleep 10 # TODO: investigate why this is needed (eventual consistency in ARM)
	az resource group deployment create \
		--verbose \
		--name "${INSTANCE_NAME}" \
		--resource-group "${INSTANCE_NAME}" \
		--template-file "${OUTPUT}/azuredeploy.json" \
		--parameters "@${OUTPUT}/azuredeploy.parameters.json" \
			2>&1 > "${OUTPUT}/deployment-debug.log"
}

function cleanup() {
	timeout 10s az resource group delete --name="${INSTANCE_NAME}" || true
}

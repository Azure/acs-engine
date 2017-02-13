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

ROOT="${DIR}/.."

# see: https://github.com/stedolan/jq/issues/105 & https://github.com/stedolan/jq/wiki/FAQ#general-questions
function jqi() { filename="${1}"; jqexpr="${2}"; jq "${jqexpr}" "${filename}" > "${filename}.tmp" && mv "${filename}.tmp" "${filename}"; }

function generate_template() {
	# Check pre-requisites
	[[ ! -z "${INSTANCE_NAME:-}" ]] || (echo "Must specify INSTANCE_NAME" && exit -1)
	[[ ! -z "${CLUSTER_DEFINITION:-}" ]] || (echo "Must specify CLUSTER_DEFINITION" && exit -1)
	[[ ! -z "${SERVICE_PRINCIPAL_CLIENT_ID:-}" ]] || [[ ! -z "${CLUSTER_SERVICE_PRINCIPAL_CLIENT_ID:-}" ]] || (echo "Must specify SERVICE_PRINCIPAL_CLIENT_ID" && exit -1)
	[[ ! -z "${SERVICE_PRINCIPAL_CLIENT_SECRET:-}" ]] || [[ ! -z "${CLUSTER_SERVICE_PRINCIPAL_CLIENT_SECRET:-}" ]] || (echo "Must specify SERVICE_PRINCIPAL_CLIENT_SECRET" && exit -1)
	[[ ! -z "${OUTPUT:-}" ]] || (echo "Must specify OUTPUT" && exit -1)
	
	# Set output directory
	mkdir -p "${OUTPUT}"

	# Set custom dir so we don't clobber global 'az' config
	AZURE_CONFIG_DIR="$(mktemp -d)"
	trap 'rm -rf ${AZURE_CONFIG_DIR}' EXIT

	# Prep SSH Key
	ssh-keygen -b 2048 -t rsa -f "${OUTPUT}/id_rsa" -q -N ""
	ssh-keygen -y -f "${OUTPUT}/id_rsa" > "${OUTPUT}/id_rsa.pub"
	export SSH_KEY_DATA="$(cat "${OUTPUT}/id_rsa.pub")"

	# Allow different credentials for cluster vs the deployment
	export CLUSTER_SERVICE_PRINCIPAL_CLIENT_ID="${CLUSTER_SERVICE_PRINCIPAL_CLIENT_ID:-${SERVICE_PRINCIPAL_CLIENT_ID}}"
	export CLUSTER_SERVICE_PRINCIPAL_CLIENT_SECRET="${CLUSTER_SERVICE_PRINCIPAL_CLIENT_SECRET:-${SERVICE_PRINCIPAL_CLIENT_SECRET}}"

	# Form the final cluster_definition file
	export FINAL_CLUSTER_DEFINITION="${OUTPUT}/clusterdefinition.json"
	cp "${CLUSTER_DEFINITION}" "${FINAL_CLUSTER_DEFINITION}"
	jqi "${FINAL_CLUSTER_DEFINITION}" ".properties.masterProfile.dnsPrefix = \"${INSTANCE_NAME}\""
	jqi "${FINAL_CLUSTER_DEFINITION}" ".properties.agentPoolProfiles |= map(if .name==\"agentpublic\" then .dnsPrefix = \"${INSTANCE_NAME}0\" else . end)"
	jqi "${FINAL_CLUSTER_DEFINITION}" ".properties.linuxProfile.ssh.publicKeys[0].keyData = \"${SSH_KEY_DATA}\""
	jqi "${FINAL_CLUSTER_DEFINITION}" ".properties.servicePrincipalProfile.servicePrincipalClientID = \"${CLUSTER_SERVICE_PRINCIPAL_CLIENT_ID}\""
	jqi "${FINAL_CLUSTER_DEFINITION}" ".properties.servicePrincipalProfile.servicePrincipalClientSecret = \"${CLUSTER_SERVICE_PRINCIPAL_CLIENT_SECRET}\""

	# Generate template
	"${DIR}/../acs-engine" -artifacts "${OUTPUT}" "${FINAL_CLUSTER_DEFINITION}"

	# Fill in custom hyperkube spec, if it was set
	if [[ ! -z "${CUSTOM_HYPERKUBE_SPEC:-}" ]]; then
		# TODO: plumb hyperkube into the apimodel
		jqi "${OUTPUT}/azuredeploy.parameters.json" ".kubernetesHyperkubeSpec.value = \"${CUSTOM_HYPERKUBE_SPEC}\""
	fi
}

function set_azure_account() {
	# Check pre-requisites
	[[ ! -z "${SUBSCRIPTION_ID:-}" ]] || (echo "Must specify SUBSCRIPTION_ID" && exit -1)
	[[ ! -z "${TENANT_ID:-}" ]] || (echo "Must specify TENANT_ID" && exit -1)
	[[ ! -z "${SERVICE_PRINCIPAL_CLIENT_ID:-}" ]] || (echo "Must specify SERVICE_PRINCIPAL_CLIENT_ID" && exit -1)
	[[ ! -z "${SERVICE_PRINCIPAL_CLIENT_SECRET:-}" ]] || (echo "Must specify SERVICE_PRINCIPAL_CLIENT_SECRET" && exit -1)
	which kubectl || (echo "kubectl must be on PATH" && exit -1)
	which az || (echo "az must be on PATH" && exit -1)

	# Login to Azure-Cli
	az login --service-principal \
		--username "${SERVICE_PRINCIPAL_CLIENT_ID}" \
		--password "${SERVICE_PRINCIPAL_CLIENT_SECRET}" \
		--tenant "${TENANT_ID}" &>/dev/null

	az account set --subscription "${SUBSCRIPTION_ID}"
}

function deploy_template() {
	# Check pre-requisites
	[[ ! -z "${DEPLOYMENT_NAME:-}" ]] || (echo "Must specify DEPLOYMENT_NAME" && exit -1)
	[[ ! -z "${LOCATION:-}" ]] || (echo "Must specify LOCATION" && exit -1)
	[[ ! -z "${RESOURCE_GROUP:-}" ]] || (echo "Must specify RESOURCE_GROUP" && exit -1)
	[[ ! -z "${OUTPUT:-}" ]] || (echo "Must specify OUTPUT" && exit -1)

	which kubectl || (echo "kubectl must be on PATH" && exit -1)
	which az || (echo "az must be on PATH" && exit -1)

	# Deploy the template
	az group create --name="${RESOURCE_GROUP}" --location="${LOCATION}"

	sleep 3 # TODO: investigate why this is needed (eventual consistency in ARM)
	az group deployment create \
		--name "${DEPLOYMENT_NAME}" \
		--resource-group "${RESOURCE_GROUP}" \
		--template-file "${OUTPUT}/azuredeploy.json" \
		--parameters "@${OUTPUT}/azuredeploy.parameters.json"
}

function cleanup() {
	if [[ "${CLEANUP:-}" == "y" ]]; then
		az group delete --no-wait --name="${RESOURCE_GROUP}" --force || true
	fi
}

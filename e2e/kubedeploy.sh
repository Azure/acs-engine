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

# TODO: move this to ./common.sh and source it
# see: https://github.com/stedolan/jq/issues/105
# and: https://github.com/stedolan/jq/wiki/FAQ#general-questions
function jqi() {
	filename="${1}"
	jqexpr="${2}"
	jq "${jqexpr}" "${filename}" > "${filename}.tmp" && mv "${filename}.tmp" "${filename}"
}

####################################################
# USAGE:
#   $ export RESOURCE_GROUP_PREFIX=mickens
#   $ export CLUSTER_DEFINITION=$(pwd)/../examples/kubernetes.json
#   $ export LOCATION=eastus
#   $ export CUSTOM_HYPERKUBE_SPEC=docker.io/colemickens/hyperkube-amd64:v1.4.1-colemickens-azure-specify-availabilityset
#   $ export SERVICE_PRINCIPAL_CLIENT_ID=http://colemick-acs-kube-sp
#   $ export SERVICE_PRINCIPAL_CLIENT_SECRET=
#   $ ./kubedeploy.sh
####################################################

function verify_prereqs() {
	# Check pre-requisites
	[[ ! -z "${RESOURCE_GROUP_PREFIX:-}" ]] || (echo "Must specify RESOURCE_GROUP_PREFIX" && exit -1)
	[[ ! -z "${LOCATION:-}" ]] || (echo "Must specify LOCATION" && exit -1)
	[[ ! -z "${CLUSTER_DEFINITION:-}" ]] || (echo "Must specify CLUSTER_DEFINITION" && exit -1)
	[[ ! -z "${SERVICE_PRINCIPAL_CLIENT_ID:-}" ]] || (echo "Must specify SERVICE_PRINCIPAL_CLIENT_ID" && exit -1)
	[[ ! -z "${SERVICE_PRINCIPAL_CLIENT_SECRET:-}" ]] || (echo "Must specify SERVICE_PRINCIPAL_CLIENT_SECRET" && exit -1)
	which kubectl || (echo "kubectl must be on PATH" && exit -1)
	which az || (echo "az must be on PATH" && exit -1)
	
	# If they didn't specify SSH data, just load it from the user
	# TODO: forcibly create new SSH keys and put them in the output directory
	[[ ! -z "${SSH_KEY_DATA:-}" ]] || SSH_KEY_DATA="$(cat ~/.ssh/id_rsa.pub)"
	
	# Calculate unique, date-based suffix
	VERSION_DATE="$(printf '%x' $(date '+%s'))"
	
	# Check RG Prefix. We use it for the DNS name too.
	# VERSION_DATE is always 8 characters, leaving 7 left over for the prefix
	[[ ${#RESOURCE_GROUP_PREFIX} -le 40 ]] || (echo "RESOURCE_GROUP_PREFIX must be no longer than 7 chars" && exit -1)
	RESOURCE_GROUP="${RESOURCE_GROUP_PREFIX}${VERSION_DATE}"
	DNS_PREFIX="${RESOURCE_GROUP}"
	
	# Set output directory
	OUTPUT="$(pwd)/_output/${RESOURCE_GROUP}"
	mkdir -p "${OUTPUT}"
}

function prepare_cluster_definition() {
	# Form the final cluster_definition file
	FINAL_CLUSTER_DEFINITION="${OUTPUT}/clusterdefinition.json"
	cp "${CLUSTER_DEFINITION}" "${FINAL_CLUSTER_DEFINITION}"

	jqi "${FINAL_CLUSTER_DEFINITION}" ".properties.masterProfile.dnsPrefix = \"${RESOURCE_GROUP}\""
	jqi "${FINAL_CLUSTER_DEFINITION}" ".properties.linuxProfile.ssh.publicKeys[0].keyData = \"${SSH_KEY_DATA}\"" t "${FINAL_CLUSTER_DEFINITION}" 
	jqi "${FINAL_CLUSTER_DEFINITION}" ".properties.servicePrincipalProfile.servicePrincipalClientID = \"${SERVICE_PRINCIPAL_CLIENT_ID}\""
	jqi "${FINAL_CLUSTER_DEFINITION}" ".properties.servicePrincipalProfile.servicePrincipalClientSecret = \"${SERVICE_PRINCIPAL_CLIENT_SECRET}\""
}

function generate_template() {
	# Generate template
	./acsengine -artifacts "${OUTPUT}" "${FINAL_CLUSTER_DEFINITION}"

	# Fill in custom hyperkube spec, if it was set
	if [[ ! -z "${CUSTOM_HYPERKUBE_SPEC:-}" ]]; then
		sed -i "s|gcr.io/google_containers/hyperkube-amd64:v1.4.0|${CUSTOM_HYPERKUBE_SPEC}|g" \
			"${OUTPUT}/azuredeploy.parameters.json"
	fi
}

function deploy_template() {
	# Deploy the template
	az resource group create --name="${RESOURCE_GROUP}" --location="${LOCATION}"
	az resource group deployment create \
		--debug \
		--name "${RESOURCE_GROUP}" \
		--resource-group "${RESOURCE_GROUP}" \
		--template-file-path "${OUTPUT}/azuredeploy.json" \
		--parameters-file-path "${OUTPUT}/azuredeploy.parameters.json" \
			2>&1 | tee "${OUTPUT}/deployment-debug.log"
}

function kubernetes_validate_cluster() {
	# Get the kubeconfig file
	export KUBECONFIG="${OUTPUT}/kubeconfig/kubeconfig.${LOCATION}.json"

	# Wait for at least some nodes
	# TODO: this only checks for count of first agent pool
	total_time=0
	wait_duration=10
	num_nodes="$(( $(jq -r '.properties.agentPoolProfiles[0].count' "${FINAL_CLUSTER_DEFINITION}") + 1 ))"
	while true; do
		total_time=$(( ${total_time} + ${wait_duration} ))
		hcount=$(kubectl get nodes 2>/dev/null | grep 'Ready' | grep -v 'NotReady' | wc -l) || true
		echo "Validation: Expected ${num_nodes} healthy nodes; found ${hcount}. (${total_time}s elapsed)"
		[[ "${hcount}" -ge "${num_nodes}" ]] && echo "Validation: Success!" && break
		sleep ${wait_duration}
	done

	# Deploy nginx and give it a public IP
	kubectl run nginx --image="nginx"
	kubectl expose deployment nginx --port="80" --type="LoadBalancer"

	# Wait for the external IP to be populated
	externalip=""
	while : ; do
		externalip=$(kubectl get svc nginx --template="{{range .status.loadBalancer.ingress}}{{.ip}}{{end}}")
		[ -z "${externalip}" ] || break; sleep 5
	done

	# TODO: curl nginx to make sure it's /really/ working
	# TODO: run k8s conformance test
}

cd "${DIR}/.." && go build .
verify_prereqs
prepare_cluster_definition
generate_template
deploy_template
kubernetes_validate_cluster

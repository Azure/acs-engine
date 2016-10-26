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

source "${DIR}/common.sh"

function validate_k8s() {
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


trap cleanup EXIT
deploy
#export -f validate_k8s && timeout 10m bash -xc validate_k8s

# TODO: it shouldn't take anywhere near 5-10 minutes
# TODO: Kube is assigning the external IP and then taking ages to actually show it

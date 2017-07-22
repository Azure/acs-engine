function test_linux_deployment() {
  ###### Testing an nginx deployment
  log "Testing deployments"
  kubectl create namespace ${namespace}

  NGINX="docker.io/library/nginx:latest"
  IMAGE="${NGINX}" # default to the library image unless we're in TEST_ACR mode
  if [[ "${TEST_ACR}" == "y" ]]; then
    # force it to pull from ACR
    IMAGE="${ACR_REGISTRY}/test/nginx:latest"
    # wait for acr
    wait
    # TODO: how to do this without polluting user home dir?
    docker login --username="${SERVICE_PRINCIPAL_CLIENT_ID}" --password="${SERVICE_PRINCIPAL_CLIENT_SECRET}" "${ACR_REGISTRY}"
    docker pull "${NGINX}"
    docker tag "${NGINX}" "${IMAGE}"
    docker push "${IMAGE}"
  fi

  kubectl run --image="${IMAGE}" nginx --namespace=${namespace} --overrides='{ "apiVersion": "extensions/v1beta1", "spec":{"template":{"spec": {"nodeSelector":{"beta.kubernetes.io/os":"linux"}}}}}'
  count=12
  while (( $count > 0 )); do
    log "  ... counting down $count"
    running=$(kubectl get pods --namespace=${namespace} | grep nginx | grep Running | wc | awk '{print $1}')
    if (( ${running} == 1 )); then break; fi
      sleep 5; count=$((count-1))
  done
  if (( ${running} != 1 )); then
    log "K8S: gave up waiting for deployment"
    kubectl get all --namespace=${namespace}
    exit 1
  fi

  kubectl expose deployments/nginx --type=LoadBalancer --namespace=${namespace} --port=80

  log "Checking Service External IP"
  count=60
  external_ip=""
  while (( $count > 0 )); do
    log "  ... counting down $count"
    external_ip=$(kubectl get svc --namespace ${namespace} nginx --template="{{range .status.loadBalancer.ingress}}{{.ip}}{{end}}" || echo "")
    [[ ! -z "${external_ip}" ]] && break
    sleep 10; count=$((count-1))
  done
  if [[ -z "${external_ip}" ]]; then
    log "K8S: gave up waiting for loadbalancer to get an ingress ip"
    exit 1
  fi

  log "Checking Service"
  count=5
  success="n"
  while (( $count > 0 )); do
    log "  ... counting down $count"
    ret=$(curl -f --max-time 60 "http://${external_ip}" | grep 'Welcome to nginx!' || echo "curl_error")
    if [[ $ret =~ .*'Welcome to nginx!'.* ]]; then
      success="y"
      break
    fi
    sleep 5; count=$((count-1))
  done
  if [[ "${success}" != "y" ]]; then
    log "K8S: failed to get expected response from nginx through the loadbalancer"
    exit 1
  fi
}

function test_windows_deployment() {
  echo "coming soon"
}

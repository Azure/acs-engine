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
    log "K8S-Linux: gave up waiting for deployment"
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
    log "K8S-Linux: gave up waiting for loadbalancer to get an ingress ip"
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
    log "K8S-Linux: failed to get expected response from nginx through the loadbalancer"
    exit 1
  fi
}

function test_windows_deployment() {
  ###### Testing a simpleweb windows deployment
  log "Testing Windows deployments"

  log "Creating simpleweb service"
  kubectl apply -f "$DIR/simpleweb-windows.yaml"
  count=90
  while (( $count > 0 )); do
    log "  ... counting down $count"
    running=$(kubectl get pods --namespace=default | grep win-webserver | grep Running | wc | awk '{print $1}')
    if (( ${running} == 1 )); then break; fi
    sleep 10; count=$((count-1))
  done
  if (( ${running} != 1 )); then
    log "K8S-Windows: gave up waiting for deployment"
    kubectl get all --namespace=default
    exit 1
  fi

  log "Checking Service External IP"
  count=60
  external_ip=""
  while (( $count > 0 )); do
    log "  ... counting down $count"
    external_ip=$(kubectl get svc --namespace default win-webserver --template="{{range .status.loadBalancer.ingress}}{{.ip}}{{end}}" || echo "")
    [[ ! -z "${external_ip}" ]] && break
    sleep 10; count=$((count-1))
  done
  if [[ -z "${external_ip}" ]]; then
    log "K8S-Windows: gave up waiting for loadbalancer to get an ingress ip"
    exit 1
  fi

  log "Checking Service"
  count=5
  success="n"
  while (( $count > 0 )); do
    log "  ... counting down $count"
    ret=$(curl -f --max-time 60 "http://${external_ip}" | grep 'Windows Container Web Server' || echo "curl_error")
    if [[ $ret =~ .*'Windows Container Web Server'.* ]]; then
      success="y"
      break
    fi
    sleep 10; count=$((count-1))
  done
  if [[ "${success}" != "y" ]]; then
    log "K8S-Windows: failed to get expected response from simpleweb through the loadbalancer"
    exit 1
  fi

  log "Checking outbound connection"
  count=10
  while (( $count > 0 )); do
    log "  ... counting down $count"
    winpodname=$(kubectl get pods --namespace=default | grep win-webserver | awk '{print $1}')
    [[ ! -z "${winpodname}" ]] && break
    sleep 10; count=$((count-1))
  done
  if [[ -z "${winpodname}" ]]; then
    log "K8S-Windows: failed to get expected pod name for simpleweb"
    exit 1
  fi

  log "query DNS"
  count=10
  success="n"
  while (( $count > 0 )); do
    log "  ... counting down $count"
    query=$(kubectl exec $winpodname -- powershell nslookup www.bing.com)
    if [[ $(echo ${query} | grep "DNS request timed out" | wc -l) == 0 ]] && [[ $(echo ${query} | grep "UnKnown" | wc -l) == 0 ]]; then
      success="y"
      break
    fi
    sleep 10; count=$((count-1))
  done

  # temporarily disable breaking on errors to allow the retry
  set +e
  log "curl external website"
  count=10
  success="n"
  while (( $count > 0 )); do
    log "  ... counting down $count"
    # curl without getting status first and see the response. getting status sometimes has the problem to hang
    # and it doesn't repro when running kubectl from the node
    kubectl exec $winpodname -- powershell iwr -UseBasicParsing -TimeoutSec 60 www.bing.com
    statuscode=$(kubectl exec $winpodname -- powershell iwr -UseBasicParsing -TimeoutSec 60 www.bing.com | grep StatusCode)
    if [[ ${statuscode} != "" ]] && [[ $(echo ${statuscode} | grep 200 | awk '{print $3}' | tr -d '\r') -eq "200" ]]; then
      log "got 200 status code"
      log "${statuscode}"
      success="y"
      break
    fi
    log "curl failed, retrying..."
    ipconfig=$(kubectl exec $winpodname -- powershell ipconfig /all)
    log "$ipconfig"
    sleep 10; count=$((count-1))
  done
  set -e
  if [[ "${success}" != "y" ]]; then
    nslookup=$(kubectl exec $winpodname -- powershell nslookup www.bing.com)
    log "$nslookup"
    log "getting the last 50 events to check timeout failure"
    hdr=$(kubectl get events | head -n 1)
    log "$hdr"
    evt=$(kubectl get events | tail -n 50)
    log "$evt"
    log "K8S-Windows: failed to get outbound internet connection inside simpleweb container"
    exit 1
  else
    log "outbound connection succeeded!"
  fi
}

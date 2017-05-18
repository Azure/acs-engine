#!/bin/bash

set -e
[ "$DEBUG" == 'true' ] && set -x

parameters=$(echo $1 | base64 -d -)

log() {
  echo "`date +'[%Y-%m-%d %H:%M:%S:%N %Z]'` $1"
}

get_param() {
  local param=$1
  echo $(echo "$parameters" | jq ".$param" -r)
}

install_script_dependencies() {
  log ''
  log 'Installing script dependencies'
  log ''

  # Install jq to obtain the input parameters
  log 'Installing jq'
  log ''
  sudo apt-get -y install jq
  log ''

  log 'done'
  log ''
}

cleanup_script_dependencies() {
  log ''
  log 'Removing script dependencies'
  log ''

  log 'Removing jq'
  log ''
  sudo apt-get -y remove jq
  log ''

  log 'done'
  log ''
}

create_deployment_yaml() {
  log ''
  log 'Creating helm-deployment.yaml file'
  log ''

  cat > helm-deployment.yaml <<EOFDEPLOYMENT
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  creationTimestamp: null
  labels:
    app: helm
    name: tiller
  name: tiller-deploy
  namespace: kube-system
spec:
  replicas: 1
  strategy: {}
  template:
    metadata:
      creationTimestamp: null
      labels:
        app: helm
        name: tiller
    spec:
      containers:
      - env:
        - name: TILLER_NAMESPACE
          value: kube-system
        image: gcr.io/kubernetes-helm/tiller:v2.2.2
        imagePullPolicy: IfNotPresent
        livenessProbe:
          httpGet:
            path: /liveness
            port: 44135
          initialDelaySeconds: 1
          timeoutSeconds: 1
        name: tiller
        ports:
        - containerPort: 44134
          name: tiller
        readinessProbe:
          httpGet:
            path: /readiness
            port: 44135
          initialDelaySeconds: 1
          timeoutSeconds: 1
        resources: {}
      serviceAccountName: ""
      volumes: null
status: {}
EOFDEPLOYMENT

  log 'done'
  log ''
}

create_service_yaml() {
  log ''
  log 'Creating helm-service.yaml file'
  log ''

  cat > helm-service.yaml <<EOFSERVICE
apiVersion: v1
kind: Service
metadata:
  creationTimestamp: null
  labels:
    app: helm
    name: tiller
  name: tiller-deploy
  namespace: kube-system
spec:
  ports:
  - name: tiller
    nodePort: 0
    port: 44134
    protocol: ""
    targetPort: tiller
  selector:
    app: helm
    name: tiller
  type: ClusterIP
status:
  loadBalancer: {}
EOFSERVICE

  log 'done'
  log ''
}

deploy_yaml() {
  log ''
  log 'Deploying yaml files'
  log ''

  log 'Deploying - helm-deployment.yaml'
  kubectl create -f 'helm-deployment.yaml'

  log 'Deploying - helm-service.yaml' 
  kubectl create -f 'helm-service.yaml'

  log 'done'
  log ''
}

log 'ACS-Engine - installing Helm Tiller (k8s)'
log '--------------------------------------------------'

install_script_dependencies
create_deployment_yaml
create_service_yaml
deploy_yaml
cleanup_script_dependencies

log ''
log 'done'
log ''

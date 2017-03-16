# Script file to configure Sysdig Agent on all Kubernetes nodes

#!/bin/bash

set -e

ACCESS_KEY="$1"

echo $(date) " - Starting Script"

# Create sysdig-daemonset.yaml file
echo $(date) " - Creating sysdig-daemonset.yaml file"

cat > sysdig-daemonset.yaml <<EOF
#Use this sysdig.yaml when Daemon Sets are enabled on Kubernetes (minimum version 1.1.1). Otherwise use the RC method.

apiVersion: extensions/v1beta1
kind: DaemonSet                     
metadata:
  name: sysdig-agent
  labels:
    app: sysdig-agent
spec:
  template:
    metadata:
      labels:
        name: sysdig-agent
    spec:
      volumes:
      - name: docker-sock
        hostPath:
         path: /var/run/docker.sock
      - name: dev-vol
        hostPath:
         path: /dev
      - name: proc-vol
        hostPath:
         path: /proc
      - name: boot-vol
        hostPath:
         path: /boot
      - name: modules-vol
        hostPath:
         path: /lib/modules
      - name: usr-vol
        hostPath:
          path: /usr
      hostNetwork: true
      hostPID: true
      containers:
      - name: sysdig-agent
        image: sysdig/agent
#        imagePullPolicy: Always                            #OPTIONAL - Always pull the latest container image tag 
        securityContext:
         privileged: true
        env:
        - name: ACCESS_KEY                                  #REQUIRED - replace with your Sysdig Cloud access key
          value: "$ACCESS_KEY"
#        - name: TAGS                                       #OPTIONAL
#          value: linux:ubuntu,dept:dev,local:nyc 
#        - name: ADDITIONAL_CONF                            #OPTIONAL pass additional parameters to the agent such as authentication example provided here
#          value: "k8s_uri: https://myacct:mypass@localhost:4430\nk8s_ca_certificate: k8s-ca.crt\nk8s_ssl_verify_certificate: true"
        volumeMounts:
        - mountPath: /host/var/run/docker.sock
          name: docker-sock
          readOnly: false
        - mountPath: /host/dev
          name: dev-vol
          readOnly: false
        - mountPath: /host/proc
          name: proc-vol
          readOnly: true
        - mountPath: /host/boot
          name: boot-vol
          readOnly: true
        - mountPath: /host/lib/modules
          name: modules-vol
          readOnly: true
        - mountPath: /host/usr
          name: usr-vol
          readOnly: true
EOF

# Deploy Sysdig DaemonSet
echo $(date) " - Deploying Sysdig DaemonSet"

kubectl create -f 'sysdig-daemonset.yaml'

echo $(date) " - Script complete"

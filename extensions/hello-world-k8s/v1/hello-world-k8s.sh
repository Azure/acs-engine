# Script file to install the docker hello-world container

#!/bin/bash

set -e

echo $(date) " - Starting Script"

# Deploy container
echo $(date) " - Deploying hello-world container"

kubectl run hello-world --image=hello-world

echo $(date) " - run kubectl get pods to list the pods"
echo $(date) " - run kubectl logs (passing the pod name gathered from kubectl get pods)"
echo $(date) " - Script complete"


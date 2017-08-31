# Script file to run hello-world in dcos

#!/bin/bash

set -e

echo $(date) " - Starting Script"

# Deploy container
echo $(date) " - Deploying hello-world"

curl -X post http://localhost:8080/v2/apps -d "{ \"id\": \"hello-marathon\", \"cmd\": \"while [ true ] ; do echo 'Hello World' ; sleep 5 ; done\", \"cpus\": 0.1, \"mem\": 10.0, \"instances\": 1 }" -H "Content-type:application/json"

echo $(date) " - view applications in marathon UI to validate"
echo $(date) " - Script complete"


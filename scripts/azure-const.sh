#!/bin/bash
if [ -z "$CLIENT_ID" ]; then
    echo "must provide a CLIENT_ID env var"
    exit 1;
fi

if [ -z "$CLIENT_SECRET" ]; then
    echo "must provide a CLIENT_SECRET env var"
    exit 1;
fi

if [ -z "$TENANT_ID" ]; then
    echo "must provide a TENANT_ID env var"
    exit 1;
fi

if [ -z "$SUBSCRIPTION_ID_TO_CLEANUP" ]; then
    echo "must provide a SUBSCRIPTION_ID_TO_CLEANUP env var"
    exit 1;
fi

set -eu -o pipefail

az login --service-principal \
		--username "${CLIENT_ID}" \
		--password "${CLIENT_SECRET}" \
		--tenant "${TENANT_ID}" &>/dev/null

# set to the sub id we want to cleanup
az account set -s $SUBSCRIPTION_ID

python pkg/acsengine/Get-AzureConstants.py
git status | grep pkg/acsengine/azureconst.go
exit_code=$?
if [ $exit_code -gt "0" ]; then
  exit 0
else
  echo "File was modified, failing test"
  exit 1
fi 

#!/bin/bash

####################################################

if [ -z "$SERVICE_PRINCIPAL_CLIENT_ID" ]; then
    exit 1;
fi

if [ -z "$SERVICE_PRINCIPAL_CLIENT_SECRET" ]; then
    exit 1;
fi

if [ -z "$TENANT_ID" ]; then
    exit 1;
fi

if [ -z "$SUBSCRIPTION_ID_TO_CLEANUP" ]; then
    exit 1;
fi

if [ -z "$EXPIRATION_IN_HOURS" ]; then
    EXPIRATION_IN_HOURS=2
fi

set -eu -o pipefail

az login --service-principal \
		--username "${SERVICE_PRINCIPAL_CLIENT_ID}" \
		--password "${SERVICE_PRINCIPAL_CLIENT_SECRET}" \
		--tenant "${TENANT_ID}" &>/dev/null

# set to the sub id we want to cleanup
az account set -s $SUBSCRIPTION_ID_TO_CLEANUP

# convert to seconds so we can compare it against the "tags.now" property in the resource group metadata
expirationInSecs=$(echo "${EXPIRATION_IN_HOURS} * 60 * 60" | bc )
# deadline = the "date +%s" representation of the oldest age we're willing to keep
(( deadline=$(date +%s)-${expirationInSecs%.*} ))
# find resource groups created before our deadline
for resourceGroup in `az group list | jq -r ".[] | select((.tags.now|tonumber < $deadline)).name"`; do
    echo "Will delete resource group ${resourceGroup}..."
    # delete old resource groups
    az group delete -n $resourceGroup -y --no-wait
done
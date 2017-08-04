#!/bin/bash

set -x
set -e

export ACSENGINE_EXPERIMENTAL_FEATURES=1

OUTPUT="_output/${INSTANCE_NAME}"
K8S_UPGRADE_CONF="$OUTPUT/k8sUpgrade.json"

cat > $K8S_UPGRADE_CONF <<END
{
  "apiVersion": "vlabs",
  "orchestratorProfile": {
      "orchestratorType": "Kubernetes",
      "orchestratorVersion": "1.6.2"
    }
}
END

./bin/acs-engine upgrade \
  --subscription-id ${SUBSCRIPTION_ID} \
  --deployment-dir ${OUTPUT} \
  --resource-group ${RESOURCE_GROUP} \
  --upgrademodel-file $K8S_UPGRADE_CONF \
  --auth-method client_secret \
  --client-id ${SERVICE_PRINCIPAL_CLIENT_ID} \
  --client-secret ${SERVICE_PRINCIPAL_CLIENT_SECRET}

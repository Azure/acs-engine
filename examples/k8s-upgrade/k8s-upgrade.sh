#!/bin/bash

set -e

# some tests set EXPECTED_ORCHESTRATOR_RELEASE in .env files
ENV_FILE="${CLUSTER_DEFINITION}.env"
if [ -e "${ENV_FILE}" ]; then
  source "${ENV_FILE}"
fi

[[ ! -z "${EXPECTED_ORCHESTRATOR_RELEASE:-}" ]] || (echo "Must specify EXPECTED_ORCHESTRATOR_RELEASE" && exit 1)

OUTPUT="_output/${INSTANCE_NAME}"
K8S_UPGRADE_CONF="$OUTPUT/k8sUpgrade.json"

cat > $K8S_UPGRADE_CONF <<END
{
  "orchestratorType": "Kubernetes",
  "orchestratorRelease": "${EXPECTED_ORCHESTRATOR_RELEASE}"
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

# (temp) allow 5 minutes for cluster to 'settle up'
# TODO: ensure cluster operability by the end of the upgrade
sleep 300

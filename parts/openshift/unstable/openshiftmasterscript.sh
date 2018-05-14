#!/bin/bash -ex

# TODO: /etc/dnsmasq.d/origin-upstream-dns.conf is currently hardcoded; it
# probably shouldn't be

if [ -f "/etc/sysconfig/atomic-openshift-node" ]; then
	SERVICE_TYPE=atomic-openshift
else
	SERVICE_TYPE=origin
fi
VERSION="$(rpm -q $SERVICE_TYPE --queryformat %{VERSION})"

if [ -f "/etc/sysconfig/atomic-openshift-node" ]; then
	ANSIBLE_DEPLOY_TYPE="openshift-enterprise"
	IMAGE_TYPE=ose
	IMAGE_PREFIX="registry.access.redhat.com/openshift3"
	ANSIBLE_CONTAINER_VERSION="v${VERSION}"
	PROMETHEUS_EXPORTER_VERSION="v${VERSION}"
	COCKPIT_PREFIX="${IMAGE_PREFIX}"
	COCKPIT_BASENAME="registry-console"
	COCKPIT_VERSION="v${VERSION}"
else
	ANSIBLE_DEPLOY_TYPE="origin"
	IMAGE_TYPE="${SERVICE_TYPE}"
	IMAGE_PREFIX="openshift"
	# FIXME: These versions are set to deal with differences in how Origin and OCP
	#        components are versioned
	ANSIBLE_CONTAINER_VERSION="v${VERSION%.*}"
	COCKPIT_PREFIX="cockpit"
	COCKPIT_BASENAME="kubernetes"
	COCKPIT_VERSION="latest"
fi
systemctl restart docker.service

echo "BOOTSTRAP_CONFIG_NAME=node-config-master" >>/etc/sysconfig/${SERVICE_TYPE}-node

for dst in tcp,2379 tcp,2380 tcp,8443 tcp,8444 tcp,8053 udp,8053 tcp,9090; do
	proto=${dst%%,*}
	port=${dst##*,}
	iptables -A OS_FIREWALL_ALLOW -p $proto -m state --state NEW -m $proto --dport $port -j ACCEPT
done

iptables-save >/etc/sysconfig/iptables

sed -i -e "s#--master=.*#--master=https://$(hostname --fqdn):8443#" /etc/sysconfig/${SERVICE_TYPE}-master-api

rm -rf /etc/etcd/* /etc/origin/master/* /etc/origin/node/*

oc adm create-bootstrap-policy-file --filename=/etc/origin/master/policy.json

( cd / && base64 -d <<< {{ .ConfigBundle }} | tar -xz)

cp /etc/origin/node/ca.crt /etc/pki/ca-trust/source/anchors/openshift-ca.crt
update-ca-trust

# FIXME: It is horrible that we're installing az.  Try to avoid adding
# additional functionality in this script that requires it.  One route to remove
# this code is to bake this script into the base image, then pass in parameters
# such as the registry storage account name and key direct from ARM.
rpm -i https://packages.microsoft.com/yumrepos/azure-cli/azure-cli-2.0.31-1.el7.x86_64.rpm

set +x
. <(sed -e 's/: */=/' /etc/azure/azure.conf)
az login --service-principal -u "$aadClientId" -p "$aadClientSecret" --tenant "$aadTenantId" &>/dev/null
REGISTRY_STORAGE_AZURE_ACCOUNTNAME=$(az storage account list -g "$resourceGroup" --query "[?ends_with(name, 'registry')].name" -o tsv)
REGISTRY_STORAGE_AZURE_ACCOUNTKEY=$(az storage account keys list -g "$resourceGroup" -n "$REGISTRY_STORAGE_AZURE_ACCOUNTNAME" --query "[?keyName == 'key1'].value" -o tsv)
az logout
set -x

###
# retrieve the public ip via dns for the router public ip and sub it in for the routingConfig.subdomain
###
routerLBHost="{{.RouterLBHostname}}"
routerLBIP=$(dig +short $routerLBHost)

# NOTE: The version of openshift-ansible for origin defaults the ansible var
#       openshift_prometheus_node_exporter_image_version correctly as needed by
#       origin, but for OCP it does not.
#
#       This is fixed in openshift/openshift-ansible@c27a0f4, which is in
#       openshift-ansible >= 3.9.15, so once we're shipping OCP >= v3.9.15 we
#       can remove this and the definition of the cooresonding variable in the
#       ansible inventory file.
if [[ "${ANSIBLE_DEPLOY_TYPE}" == "origin" ]]; then
    sed -i "/PROMETHEUS_EXPORTER_VERSION/d" /tmp/ansible/azure-local-master-inventory.yml
else
    sed -i "s|PROMETHEUS_EXPORTER_VERSION|${PROMETHEUS_EXPORTER_VERSION}|g;" /tmp/ansible/azure-local-master-inventory.yml
fi

for i in /etc/origin/master/master-config.yaml /tmp/bootstrapconfigs/* /tmp/ansible/azure-local-master-inventory.yml; do
    sed -i "s/TEMPROUTERIP/${routerLBIP}/; s|IMAGE_PREFIX|$IMAGE_PREFIX|g; s|ANSIBLE_DEPLOY_TYPE|$ANSIBLE_DEPLOY_TYPE|g" $i
    sed -i "s|REGISTRY_STORAGE_AZURE_ACCOUNTNAME|${REGISTRY_STORAGE_AZURE_ACCOUNTNAME}|g; s|REGISTRY_STORAGE_AZURE_ACCOUNTKEY|${REGISTRY_STORAGE_AZURE_ACCOUNTKEY}|g" $i
    sed -i "s|COCKPIT_VERSION|${COCKPIT_VERSION}|g; s|COCKPIT_BASENAME|${COCKPIT_BASENAME}|g; s|COCKPIT_PREFIX|${COCKPIT_PREFIX}|g;" $i
    sed -i "s|VERSION|${VERSION}|g; s|SHORT_VER|${VERSION%.*}|g; s|SERVICE_TYPE|${SERVICE_TYPE}|g; s|IMAGE_TYPE|${IMAGE_TYPE}|g" $i
    sed -i "s|HOSTNAME|${HOSTNAME}|g;" $i
done

# note: ${SERVICE_TYPE}-node crash loops until master is up
for unit in etcd.service ${SERVICE_TYPE}-master-api.service ${SERVICE_TYPE}-master-controllers.service; do
	systemctl enable $unit
	systemctl start $unit
done

mkdir -p /root/.kube
cp /etc/origin/master/admin.kubeconfig /root/.kube/config

export KUBECONFIG=/etc/origin/master/admin.kubeconfig

while ! curl -o /dev/null -m 2 -kfs https://localhost:8443/healthz; do
	sleep 1
done

while ! oc get svc kubernetes &>/dev/null; do
	sleep 1
done

oc create -f - <<'EOF'
kind: StorageClass
apiVersion: storage.k8s.io/v1beta1
metadata:
  name: azure
  annotations:
    storageclass.kubernetes.io/is-default-class: "true"
provisioner: kubernetes.io/azure-disk
parameters:
  skuName: Premium_LRS
  location: {{ .Location }}
  kind: managed
EOF

oc create configmap node-config-master --namespace openshift-node --from-file=node-config.yaml=/tmp/bootstrapconfigs/master-config.yaml
oc create configmap node-config-compute --namespace openshift-node --from-file=node-config.yaml=/tmp/bootstrapconfigs/compute-config.yaml
oc create configmap node-config-infra --namespace openshift-node --from-file=node-config.yaml=/tmp/bootstrapconfigs/infra-config.yaml

# must start ${SERVICE_TYPE}-node after master is fully up and running
# otherwise the implicit dns change may cause master startup to fail
systemctl enable ${SERVICE_TYPE}-node.service
systemctl start ${SERVICE_TYPE}-node.service &

chmod +x /tmp/ansible/ansible.sh

docker run \
	--rm \
	-u "$(id -u)" \
	-v /etc/origin:/etc/origin:z \
	-v /tmp/ansible:/opt/app-root/src:z \
	-v /root/.kube:/opt/app-root/src/.kube:z \
	-w /opt/app-root/src \
	-e IMAGE_BASE="${IMAGE_PREFIX}/${IMAGE_TYPE}" \
	-e VERSION="$VERSION" \
	-e HOSTNAME="$(hostname)" \
	--network="host" \
	"${IMAGE_PREFIX}/${IMAGE_TYPE}-ansible:${ANSIBLE_CONTAINER_VERSION}" \
	/opt/app-root/src/ansible.sh

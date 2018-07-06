#!/bin/bash -ex

# TODO: /etc/dnsmasq.d/origin-upstream-dns.conf is currently hardcoded; it
# probably shouldn't be

if [ -f "/etc/sysconfig/atomic-openshift-node" ]; then
	SERVICE_TYPE=atomic-openshift
else
	SERVICE_TYPE=origin
fi
VERSION="$(rpm -q $SERVICE_TYPE --queryformat %{VERSION})"
IP_ADDRESS="{{ .MasterIP }}"

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

if grep -q ^ResourceDisk.Filesystem=xfs /etc/waagent.conf; then
	# Bad image: docker and waagent are racing.  Try to fix up.  Leave this code
	# until the bad images have gone away.
	set +e

	# stop docker if it hasn't failed already
	systemctl stop docker.service

	# wait until waagent has run mkfs and mounted /var/lib/docker
	while ! mountpoint -q /var/lib/docker; do
		sleep 1
	done

	# now roll us back. /var/lib/docker/* may be mounted if docker lost the
	# race.
	umount /var/lib/docker
	umount /var/lib/docker/*

	# disable waagent from racing again if we reboot.
	sed -i -e '/^ResourceDisk.Format=/ s/=.*/=n/' /etc/waagent.conf
	set -e
fi

systemctl stop docker.service
# Also a bad image: the umount should also go away.
umount /var/lib/docker || true
mkfs.xfs -f /dev/sdb1
echo '/dev/sdb1  /var/lib/docker  xfs  grpquota  0 0' >>/etc/fstab
mount /var/lib/docker
restorecon -R /var/lib/docker
systemctl start docker.service

echo "BOOTSTRAP_CONFIG_NAME=node-config-master" >>/etc/sysconfig/${SERVICE_TYPE}-node

for dst in tcp,2379 tcp,2380 tcp,8443 tcp,8444 tcp,8053 udp,8053 tcp,9090; do
	proto=${dst%%,*}
	port=${dst##*,}
	iptables -A OS_FIREWALL_ALLOW -p $proto -m state --state NEW -m $proto --dport $port -j ACCEPT
done

iptables-save >/etc/sysconfig/iptables

rm -rf /etc/etcd/* /etc/origin/master/*

mkdir -p /etc/origin/master

oc adm create-bootstrap-policy-file --filename=/etc/origin/master/policy.json

( cd / && base64 -d <<< {{ .ConfigBundle | shellQuote }} | tar -xz)

cp /etc/origin/node/ca.crt /etc/pki/ca-trust/source/anchors/openshift-ca.crt
update-ca-trust

# FIXME: It is horrible that we're installing az.  Try to avoid adding
# additional functionality in this script that requires it.  One route to remove
# this code is to bake this script into the base image, then pass in parameters
# such as the registry storage account name and key direct from ARM.
rpm -i https://packages.microsoft.com/yumrepos/azure-cli/azure-cli-2.0.31-1.el7.x86_64.rpm

set +x
. <(sed -e 's/: */=/' /etc/origin/cloudprovider/azure.conf)
az login --service-principal -u "$aadClientId" -p "$aadClientSecret" --tenant "$aadTenantId" &>/dev/null
REGISTRY_STORAGE_AZURE_ACCOUNTNAME=$(az storage account list -g "$resourceGroup" --query "[?ends_with(name, 'registry')].name" -o tsv)
REGISTRY_STORAGE_AZURE_ACCOUNTKEY=$(az storage account keys list -g "$resourceGroup" -n "$REGISTRY_STORAGE_AZURE_ACCOUNTNAME" --query "[?keyName == 'key1'].value" -o tsv)
az logout
set -x

###
# retrieve the public ip via dns for the router public ip and sub it in for the routingConfig.subdomain
###
routerLBHost={{ .RouterLBHostname | shellQuote }}
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

MASTER_OREG_URL="$IMAGE_PREFIX/$IMAGE_TYPE"
if [[ -f /etc/origin/oreg_url ]]; then
	MASTER_OREG_URL=$(cat /etc/origin/oreg_url)
fi

for i in /etc/origin/master/master-config.yaml /tmp/bootstrapconfigs/* /tmp/ansible/azure-local-master-inventory.yml; do
    sed -i "s/TEMPROUTERIP/${routerLBIP}/; s|IMAGE_PREFIX|$IMAGE_PREFIX|g; s|ANSIBLE_DEPLOY_TYPE|$ANSIBLE_DEPLOY_TYPE|g" $i
    sed -i "s|REGISTRY_STORAGE_AZURE_ACCOUNTNAME|${REGISTRY_STORAGE_AZURE_ACCOUNTNAME}|g; s|REGISTRY_STORAGE_AZURE_ACCOUNTKEY|${REGISTRY_STORAGE_AZURE_ACCOUNTKEY}|g" $i
    sed -i "s|COCKPIT_VERSION|${COCKPIT_VERSION}|g; s|COCKPIT_BASENAME|${COCKPIT_BASENAME}|g; s|COCKPIT_PREFIX|${COCKPIT_PREFIX}|g;" $i
    sed -i "s|VERSION|${VERSION}|g; s|SHORT_VER|${VERSION%.*}|g; s|SERVICE_TYPE|${SERVICE_TYPE}|g; s|IMAGE_TYPE|${IMAGE_TYPE}|g" $i
    sed -i "s|HOSTNAME|${HOSTNAME}|g;" $i
    sed -i "s|MASTER_OREG_URL|${MASTER_OREG_URL}|g" $i
done

mkdir -p /root/.kube

for loc in /root/.kube/config /etc/origin/node/bootstrap.kubeconfig /etc/origin/node/node.kubeconfig; do
  cp /etc/origin/master/admin.kubeconfig "$loc"
done


# Patch the etcd_ip address placed inside of the static pod definition from the node install
sed -i "s/ETCD_IP_REPLACE/${IP_ADDRESS}/g" /etc/origin/node/disabled/etcd.yaml

export KUBECONFIG=/etc/origin/master/admin.kubeconfig

# Move each static pod into place so the kubelet will run it.
# Pods: [apiserver, controller, etcd]
oc set env --local -f /etc/origin/node/disabled/apiserver.yaml DEBUG_LOGLEVEL=4 -o yaml --dry-run > /etc/origin/node/pods/apiserver.yaml
oc set env --local -f /etc/origin/node/disabled/controller.yaml DEBUG_LOGLEVEL=4 -o yaml --dry-run > /etc/origin/node/pods/controller.yaml
mv /etc/origin/node/disabled/etcd.yaml /etc/origin/node/pods/etcd.yaml
rm -rf /etc/origin/node/disabled

systemctl start ${SERVICE_TYPE}-node

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
  location: {{ .Location | quote }}
  kind: managed
EOF

oc create configmap node-config-master --namespace openshift-node --from-file=node-config.yaml=/tmp/bootstrapconfigs/master-config.yaml
oc create configmap node-config-compute --namespace openshift-node --from-file=node-config.yaml=/tmp/bootstrapconfigs/compute-config.yaml
oc create configmap node-config-infra --namespace openshift-node --from-file=node-config.yaml=/tmp/bootstrapconfigs/infra-config.yaml

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

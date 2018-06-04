#!/bin/bash -ex


### Variables 

# TODO: /etc/dnsmasq.d/origin-upstream-dns.conf is currently hardcoded; it
# probably shouldn't be
SERVICE_TYPE=origin
if [ -f "/etc/sysconfig/atomic-openshift-node" ]; then
    SERVICE_TYPE=atomic-openshift
fi

VERSION="$(rpm -q $SERVICE_TYPE --queryformat %{VERSION})"
IP_ADDRESS="{{ .MasterIP }}"
NODE_CONFIG_NAMESPACE=openshift-node
NODE_CONFIGMAP_LIST="master infra compute"

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

### Functions 

# validates if ip address is valid format. 
# TODO: Remove this when we can abstract some of the validation with more error prune code
function valid_ip() {
    local  ip=$1
    local  stat=1

    if [[ $ip =~ ^[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}$ ]]; then
        OIFS=$IFS
        IFS='.'
        ip=($ip)
        IFS=$OIFS
        [[ ${ip[0]} -le 255 && ${ip[1]} -le 255 \
            && ${ip[2]} -le 255 && ${ip[3]} -le 255 ]]
        stat=$?
    fi
    return $stat
}

# retrieve the public ip via dns for the router public ip and sub it in for the routingConfig.subdomain
function set_router_lb() {

ROUTER_LB_HOST={{.RouterLBHostname}}
ROUTER_LB_IP=$(dig +short ${ROUTER_LB_HOST})
valid_ip ${ROUTER_LB_IP}

	local n=0
	local try=5
        until [[ $n -ge $try ]]
        do
            ROUTER_LB_IP=$(dig +short ${ROUTER_LB_HOST})
            valid_ip ${ROUTER_LB_IP} && break || {
                        echo "Command Fail.."
                        ((n++))
                        echo "retry $n ::"
                        sleep 5;
                        }

        done
}

# Create/update nodes config maps.
function create_nodes_configmap () {
    for name in ${NODE_CONFIGMAP_LIST}; do
       if ! oc get configmap node-config-${name} --namespace ${NODE_CONFIG_NAMESPACE} 1> /dev/null 2>&1; then
  			 oc create configmap node-config-${name} --namespace ${NODE_CONFIG_NAMESPACE} --from-file=node-config.yaml=/tmp/bootstrapconfigs/${name}-config.yaml
		else
  			oc create configmap node-config-${name} --namespace ${NODE_CONFIG_NAMESPACE} --from-file=node-config.yaml=/tmp/bootstrapconfigs/${name}-config.yaml --dry-run -o yaml | oc replace -f - --namespace $NODE_CONFIG_NAMESPACE
		fi
    done
}

# update iptables rules
function configure_iptables() {
for dst in tcp,2379 tcp,2380 tcp,8443 tcp,8444 tcp,8053 udp,8053 tcp,9090; do
	proto=${dst%%,*}
	port=${dst##*,}
	iptables -A OS_FIREWALL_ALLOW -p $proto -m state --state NEW -m $proto --dport $port -j ACCEPT
done
iptables-save >/etc/sysconfig/iptables
}

# inject generated certificates
function configure_certificates() {
rm -rf /etc/etcd/* /etc/origin/master/*

mkdir -p /etc/origin/master

oc adm create-bootstrap-policy-file --filename=/etc/origin/master/policy.json

( cd / && base64 -d <<< {{ .ConfigBundle }} | tar -xz)

cp /etc/origin/node/ca.crt /etc/pki/ca-trust/source/anchors/openshift-ca.crt
update-ca-trust
}

# configure azure cloud provider
function configure_cloud_provider() {
set +x
. <(sed -e 's/: */=/' /etc/origin/cloudprovider/azure.conf)
az login --service-principal -u "$aadClientId" -p "$aadClientSecret" --tenant "$aadTenantId" &>/dev/null
REGISTRY_STORAGE_AZURE_ACCOUNTNAME=$(az storage account list -g "$resourceGroup" --query "[?ends_with(name, 'registry')].name" -o tsv)
REGISTRY_STORAGE_AZURE_ACCOUNTKEY=$(az storage account keys list -g "$resourceGroup" -n "$REGISTRY_STORAGE_AZURE_ACCOUNTNAME" --query "[?keyName == 'key1'].value" -o tsv)
az logout
set -x

}

# configure ansible inventory
function configure_inventory() {
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
    sed -i "s/TEMPROUTERIP/${ROUTER_LB_IP}/; s|IMAGE_PREFIX|$IMAGE_PREFIX|g; s|ANSIBLE_DEPLOY_TYPE|$ANSIBLE_DEPLOY_TYPE|g" $i
    sed -i "s|REGISTRY_STORAGE_AZURE_ACCOUNTNAME|${REGISTRY_STORAGE_AZURE_ACCOUNTNAME}|g; s|REGISTRY_STORAGE_AZURE_ACCOUNTKEY|${REGISTRY_STORAGE_AZURE_ACCOUNTKEY}|g" $i
    sed -i "s|COCKPIT_VERSION|${COCKPIT_VERSION}|g; s|COCKPIT_BASENAME|${COCKPIT_BASENAME}|g; s|COCKPIT_PREFIX|${COCKPIT_PREFIX}|g;" $i
    sed -i "s|VERSION|${VERSION}|g; s|SHORT_VER|${VERSION%.*}|g; s|SERVICE_TYPE|${SERVICE_TYPE}|g; s|IMAGE_TYPE|${IMAGE_TYPE}|g" $i
    sed -i "s|HOSTNAME|${HOSTNAME}|g;" $i
done
}

# configure static pods and wait for ready state
function configure_static_pods() {
mkdir -p /root/.kube
for loc in /root/.kube/config /etc/origin/node/bootstrap.kubeconfig /etc/origin/node/node.kubeconfig; do
  cp /etc/origin/master/admin.kubeconfig "$loc"
done


# Patch the etcd_ip address placed inside of the static pod definition from the node install
sed -i "s/ETCD_IP_REPLACE/${IP_ADDRESS}/g" /etc/origin/node/disabled/etcd.yaml

# Move each static pod into place so the kubelet will run it.
# Pods: [apiserver, controller, etcd]
if ls /etc/origin/node/disabled/* 1> /dev/null 2>&1; then
    mv /etc/origin/node/disabled/* /etc/origin/node/pods
fi
systemctl start ${SERVICE_TYPE}-node

export KUBECONFIG=/etc/origin/master/admin.kubeconfig
while ! curl -o /dev/null -m 2 -kfs https://localhost:8443/healthz; do
	sleep 1
done

while ! oc get svc kubernetes &>/dev/null; do
	sleep 1
done

}

# create/update storage class
function configure_storage() {

oc apply -f - <<'EOF'
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

}

# call ansible playbooks and install hosted components.
# FIXME: ansible image is not part of golden image. This can create flakes in the bootstrap.
function install_hosted_components() {

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
}

### Actions

systemctl restart docker.service

echo "BOOTSTRAP_CONFIG_NAME=node-config-master" >>/etc/sysconfig/${SERVICE_TYPE}-node

# FIXME: It is horrible that we're installing az.  Try to avoid adding
# additional functionality in this script that requires it.  One route to remove
# this code is to bake this script into the base image, then pass in parameters
# such as the registry storage account name and key direct from ARM.
if ! rpm -qa | grep -qw azure-cli; then
  rpm -i https://packages.microsoft.com/yumrepos/azure-cli/azure-cli-2.0.31-1.el7.x86_64.rpm
fi

configure_iptables
configure_certificates
configure_cloud_provider
set_router_lb
configure_inventory
configure_static_pods
create_nodes_configmap
configure_storage
install_hosted_components
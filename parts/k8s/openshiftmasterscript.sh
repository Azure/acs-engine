#!/bin/bash -x

# TODO: /etc/dnsmasq.d/origin-upstream-dns.conf is currently hardcoded; it
# probably shouldn't be
SERVICE_TYPE=origin
IMAGE_BASE=openshift/origin
if [ -f "/etc/sysconfig/atomic-openshift-node" ]; then
	SERVICE_TYPE=atomic-openshift
	IMAGE_BASE=registry.reg-aws.openshift.com:443/openshift3/ose
fi

# TODO: remove this once we generate the registry certificate
cat >>/etc/sysconfig/docker <<'EOF'
INSECURE_REGISTRY='--insecure-registry 172.30.0.0/16'
EOF

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

chown -R etcd:etcd /etc/etcd
chmod 0600 /etc/origin/master/htpasswd

cp /etc/origin/node/ca.crt /etc/pki/ca-trust/source/anchors/openshift-ca.crt
update-ca-trust

###
# retrieve the public ip via dns for the router public ip and sub it in for the routingConfig.subdomain
###
routerLBHost="{{.RouterLBHostname}}"
routerLBIP=$(dig +short $routerLBHost)

for i in /etc/origin/master/master-config.yaml /tmp/bootstrapconfigs/* /tmp/ansible/azure-local-master-inventory.yml; do
	sed -i "s/TEMPROUTERIP/${routerLBIP}/; s|TEMPIMAGEBASE|$IMAGE_BASE|" $i
done

# TODO: when enabling secure registry, may need:
# ln -s /etc/origin/node/node-client-ca.crt /etc/docker/certs.d/docker-registry.default.svc:5000

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
EOF

oc create configmap node-config-master --namespace openshift-node --from-file=node-config.yaml=/tmp/bootstrapconfigs/master-config.yaml
oc create configmap node-config-compute --namespace openshift-node --from-file=node-config.yaml=/tmp/bootstrapconfigs/compute-config.yaml
oc create configmap node-config-infra --namespace openshift-node --from-file=node-config.yaml=/tmp/bootstrapconfigs/infra-config.yaml

# must start ${SERVICE_TYPE}-node after master is fully up and running
# otherwise the implicit dns change may cause master startup to fail
systemctl enable ${SERVICE_TYPE}-node.service
systemctl start ${SERVICE_TYPE}-node.service &

# TODO: run a CSR auto-approver
# https://github.com/kargakis/acs-engine/issues/46
csrs=($(oc get csr -o name))
while [[ ${#csrs[@]} != "3" ]]; do
	sleep 2
	csrs=($(oc get csr -o name))
	if [[ ${#csrs[@]} == "3" ]]; then
		break
	fi
done

for csr in ${csrs[@]}; do
	oc adm certificate approve $csr
done

csrs=($(oc get csr -o name))
while [[ ${#csrs[@]} != "6" ]]; do
	sleep 2
	csrs=($(oc get csr -o name))
	if [[ ${#csrs[@]} == "6" ]]; then
		break
	fi
done

for csr in ${csrs[@]}; do
	oc adm certificate approve $csr
done

chmod +x /tmp/ansible/ansible.sh
docker run \
	--rm \
	-u "$(id -u)" \
	-v /etc/origin:/etc/origin:z \
	-v /tmp/ansible:/opt/app-root/src:z \
	-v /root/.kube:/opt/app-root/src/.kube:z \
	-w /opt/app-root/src \
	-e IMAGE_BASE="$IMAGE_BASE" \
	-e HOSTNAME="$(hostname)" \
	"$IMAGE_BASE-ansible:v3.9.11" \
	/opt/app-root/src/ansible.sh

exit 0

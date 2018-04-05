#!/bin/bash -x

# TODO: /etc/dnsmasq.d/origin-upstream-dns.conf is currently hardcoded; it
# probably shouldn't be
SERVICE_TYPE=origin
if [ -f "/etc/sysconfig/atomic-openshift-node" ]; then
    SERVICE_TYPE=atomic-openshift
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
sed -i "s/TEMPROUTERIP/${routerLBIP}/" /etc/origin/master/master-config.yaml
sed -i "s/TEMPROUTERIP/${routerLBIP}/" /tmp/ansible-playbooks/azure-local-master-inventory.yml

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

# TODO: do this, and more (registry console, asb), the proper way

oc patch project default -p '{"metadata":{"annotations":{"openshift.io/node-selector": ""}}}'

oc adm registry --images='registry.reg-aws.openshift.com:443/openshift3/ose-${component}:${version}' --selector='region=infra'

# Deploy the router reusing relevant parts from openshift-ansible
ANSIBLE_ROLES_PATH=/usr/share/ansible/openshift-ansible/roles/ ansible-playbook -c local /tmp/ansible-playbooks/deploy-router.yml -i /tmp/ansible-playbooks/azure-local-master-inventory.yml

oc create -f - <<'EOF'
kind: Project
apiVersion: v1
metadata:
  name: openshift-web-console
  annotations:
    openshift.io/node-selector: ""
EOF

oc process -f /usr/share/ansible/openshift-ansible/roles/openshift_web_console/files/console-template.yaml \
	-p API_SERVER_CONFIG="$(sed -e s/127.0.0.1/{{ .ExternalMasterHostname }}/g </usr/share/ansible/openshift-ansible/roles/openshift_web_console/files/console-config.yaml)" \
	-p NODE_SELECTOR='{"node-role.kubernetes.io/master":"true"}' \
	-p IMAGE='registry.reg-aws.openshift.com:443/openshift3/ose-web-console:v3.9.11' \
	| oc create -f -

oc create -f - <<'EOF'
kind: Project
apiVersion: v1
metadata:
  name: kube-service-catalog
  annotations:
    openshift.io/node-selector: ""
EOF

oc create secret generic -n kube-service-catalog apiserver-ssl \
  --from-file=tls.crt=/etc/origin/service-catalog/apiserver.crt \
  --from-file=tls.key=/etc/origin/service-catalog/apiserver.key

oc create secret generic -n kube-service-catalog service-catalog-ssl \
	--from-file=tls.crt=/etc/origin/service-catalog/apiserver.crt

oc create -f - <<EOF
apiVersion: apiregistration.k8s.io/v1beta1
kind: APIService
metadata:
  name: v1beta1.servicecatalog.k8s.io
spec:
  caBundle: $(base64 -w0 </etc/origin/service-catalog/ca.crt)
  group: servicecatalog.k8s.io
  groupPriorityMinimum: 20
  service:
    name: apiserver
    namespace: kube-service-catalog
  version: v1beta1
  versionPriority: 10
EOF

oc project kube-service-catalog
oc process -f /usr/share/ansible/openshift-ansible/roles/openshift_service_catalog/files/kubeservicecatalog_roles_bindings.yml | oc create -f -
oc project default
oc process -f /usr/share/ansible/openshift-ansible/roles/openshift_service_catalog/files/kubesystem_roles_bindings.yml | oc create -f -
oc auth reconcile -f /usr/share/ansible/openshift-ansible/roles/openshift_service_catalog/files/openshift_catalog_clusterroles.yml
oc adm policy add-scc-to-user hostmount-anyuid system:serviceaccount:kube-service-catalog:service-catalog-apiserver
oc adm policy add-cluster-role-to-user admin system:serviceaccount:kube-service-catalog:default
oc process -f /tmp/service-catalog/objects.yaml \
  -p CA_HASH="$(base64 -w0 </etc/origin/service-catalog/ca.crt | sha1sum | cut -d' ' -f1)" \
  -p ETCD_SERVER="$(hostname)" \
  | oc create -f -
oc rollout status -n kube-service-catalog daemonset apiserver

oc create -f - <<'EOF'
kind: Project
apiVersion: v1
metadata:
  name: openshift-template-service-broker
  annotations:
    openshift.io/node-selector: ""
EOF

oc process -f /usr/share/ansible/openshift-ansible/roles/template_service_broker/files/apiserver-template.yaml \
	-p IMAGE='registry.reg-aws.openshift.com:443/openshift3/ose-template-service-broker:v3.9.11' \
	-p NODE_SELECTOR='{"region":"infra"}' \
	| oc create -f -
oc process -f /usr/share/ansible/openshift-ansible/roles/template_service_broker/files/rbac-template.yaml | oc auth reconcile -f -

while true; do
  oc process -f /usr/share/ansible/openshift-ansible/roles/template_service_broker/files/template-service-broker-registration.yaml \
	  -p CA_BUNDLE=$(base64 -w 0 </etc/origin/master/service-signer.crt) \
	  | oc create -f - && break
  sleep 10
done

for file in /usr/share/ansible/openshift-ansible/roles/openshift_examples/files/examples/v3.9/db-templates/*.json \
    /usr/share/ansible/openshift-ansible/roles/openshift_examples/files/examples/v3.9/image-streams/*-rhel7.json \
	  /usr/share/ansible/openshift-ansible/roles/openshift_examples/files/examples/v3.9/quickstart-templates/*.json \
	  /usr/share/ansible/openshift-ansible/roles/openshift_examples/files/examples/v3.9/xpaas-streams/*.json \
	  /usr/share/ansible/openshift-ansible/roles/openshift_examples/files/examples/v3.9/xpaas-templates/*.json; do
	oc create -n openshift -f $file
done

# TODO: possibly wait here for convergence?

exit 0

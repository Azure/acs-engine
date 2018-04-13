#!/bin/bash -x

# TODO: do this, and more (registry console, asb), the proper way

# we get "dial tcp: lookup foo.eastus.cloudapp.azure.com on 10.0.0.11:53: read
# udp 172.17.0.2:56662->10.0.0.11:53: read: no route to host errors" at
# start-up: wait until these subside.
while ! oc version &>/dev/null; do
  sleep 1
done

oc patch project default -p '{"metadata":{"annotations":{"openshift.io/node-selector": ""}}}'

# FIXME - This should be handled by the openshift-ansible playbooks to ensure
#         a directory it needs to write to exists before attempting to write
#         to it
mkdir -p /etc/origin/master/named_certificates

# Deploy the router reusing relevant parts from openshift-ansible
ANSIBLE_ROLES_PATH=/usr/share/ansible/openshift-ansible/roles/ ansible-playbook -c local azure-ocp-deploy.yml -i azure-local-master-inventory.yml

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
oc process -f service-catalog.yaml \
  -p CA_HASH="$(base64 -w0 </etc/origin/service-catalog/ca.crt | sha1sum | cut -d' ' -f1)" \
  -p ETCD_SERVER="$HOSTNAME" \
	-p IMAGE="$IMAGE_BASE-service-catalog:v$VERSION" \
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
	-p IMAGE="$IMAGE_BASE-template-service-broker:v$VERSION" \
	-p NODE_SELECTOR='{"region":"infra"}' \
	| oc create -f -
oc process -f /usr/share/ansible/openshift-ansible/roles/template_service_broker/files/rbac-template.yaml | oc auth reconcile -f -

while true; do
  oc process -f /usr/share/ansible/openshift-ansible/roles/template_service_broker/files/template-service-broker-registration.yaml \
	  -p CA_BUNDLE=$(base64 -w 0 </etc/origin/master/service-signer.crt) \
	  | oc create -f - && break
  sleep 10
done

MAJORVERSION="${VERSION%.*}"
for file in /usr/share/ansible/openshift-ansible/roles/openshift_examples/files/examples/v$MAJORVERSION/db-templates/*.json \
    /usr/share/ansible/openshift-ansible/roles/openshift_examples/files/examples/v$MAJORVERSION/image-streams/*-rhel7.json \
	  /usr/share/ansible/openshift-ansible/roles/openshift_examples/files/examples/v$MAJORVERSION/quickstart-templates/*.json \
	  /usr/share/ansible/openshift-ansible/roles/openshift_examples/files/examples/v$MAJORVERSION/xpaas-streams/*.json \
	  /usr/share/ansible/openshift-ansible/roles/openshift_examples/files/examples/v$MAJORVERSION/xpaas-templates/*.json; do
	oc create -n openshift -f $file
done

# TODO: possibly wait here for convergence?

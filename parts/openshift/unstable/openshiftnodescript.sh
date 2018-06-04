#!/bin/bash -ex

### Variables 

# TODO: /etc/dnsmasq.d/origin-upstream-dns.conf is currently hardcoded; it
# probably shouldn't be
SERVICE_TYPE=origin
if [ -f "/etc/sysconfig/atomic-openshift-node" ]; then
    SERVICE_TYPE=atomic-openshift
fi

### Functions

# configure node certificates
configure_certificates() {

{{if eq .Role "infra"}}
echo "BOOTSTRAP_CONFIG_NAME=node-config-infra" >>/etc/sysconfig/${SERVICE_TYPE}-node
{{else}}
echo "BOOTSTRAP_CONFIG_NAME=node-config-compute" >>/etc/sysconfig/${SERVICE_TYPE}-node
{{end}}

rm -rf /etc/etcd/* /etc/origin/master/*

( cd / && base64 -d <<< {{ .ConfigBundle }} | tar -xz)

cp /etc/origin/node/ca.crt /etc/pki/ca-trust/source/anchors/openshift-ca.crt
update-ca-trust

}

# configure and start openshift-node process
configure_openshift_node(){
# note: ${SERVICE_TYPE}-node crash loops until master is up
systemctl enable ${SERVICE_TYPE}-node.service
systemctl start ${SERVICE_TYPE}-node.service &

while [[ $(KUBECONFIG=/etc/origin/node/node.kubeconfig oc get node $(hostname) -o template \
    --template '{{`{{range .status.conditions}}{{if eq .type "Ready"}}{{.status}}{{end}}{{end}}`}}') != True ]]; do
    sleep 1
done

}

### Actions
configure_certificates
configure_openshift_node
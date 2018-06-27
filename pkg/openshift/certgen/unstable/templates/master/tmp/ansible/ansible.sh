#!/bin/bash -ex

# TODO: do this, and more (registry console, asb), the proper way

# we get "dial tcp: lookup foo.eastus.cloudapp.azure.com on 10.0.0.11:53: read
# udp 172.17.0.2:56662->10.0.0.11:53: read: no route to host errors" at
# start-up: wait until these subside.
while ! oc version &>/dev/null; do
  sleep 1
done

for project in default openshift-infra; do
  oc patch project $project -p '{"metadata":{"annotations":{"openshift.io/node-selector": ""}}}'
done

# FIXME - This should be handled by the openshift-ansible playbooks to ensure
#         a directory it needs to write to exists before attempting to write
#         to it
mkdir -p /etc/origin/master/named_certificates

# Deploy all infra components reusing relevant parts from openshift-ansible
ANSIBLE_ROLES_PATH=/usr/share/ansible/openshift-ansible/roles/ \
  ansible-playbook\
  -M /usr/share/ansible/openshift-ansible/roles/lib_utils/library \
  -c local azure-ocp-deploy.yml \
  -i azure-local-master-inventory.yml

# TODO: possibly wait here for convergence?

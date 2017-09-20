#!/bin/bash

PROXY_CA_KEY="${PROXY_CA_KEY:=/tmp/proxy-client-ca.key}"
PROXY_CRT="${PROXY_CRT:=/tmp/proxy-client-ca.crt}"
PROXY_CLIENT_KEY="${PROXY_CLIENT_KEY:=/tmp/proxy-client.key}"
PROXY_CLIENT_CSR="${PROXY_CLIENT_CSR:=/tmp/proxy-client.csr}"
PROXY_CLIENT_CRT="${PROXY_CLIENT_CRT:=/tmp/proxy-client.crt}"
ETCD_REQUESTHEADER_CLIENT_CA="${ETCD_REQUESTHEADER_CLIENT_CA:=/proxycerts/requestheader-client-ca-file}"
ETCD_PROXY_CERT="${ETCD_PROXY_CERT:=/proxycerts/proxy-client-cert-file}"
ETCD_PROXY_KEY="${ETCD_PROXY_KEY:=/proxycerts/proxy-client-key-file}"
K8S_PROXY_CA_CRT_FILEPATH="${K8S_PROXY_CA_CRT_FILEPATH:=/etc/kubernetes/certs/proxy-ca.crt}"
K8S_PROXY_KEY_FILEPATH="${K8S_PROXY_KEY_FILEPATH:=/etc/kubernetes/certs/proxy.key}"
K8S_PROXY_CRT_FILEPATH="${K8S_PROXY_CRT_FILEPATH:=/etc/kubernetes/certs/proxy.crt}"

# generate root CA
openssl genrsa -out $PROXY_CA_KEY 2048
openssl req -new -x509 -days 1826 -key $PROXY_CA_KEY -out $PROXY_CRT -subj '/CN=proxyClientCA'
# generate new cert
openssl genrsa -out $PROXY_CLIENT_KEY 2048
openssl req -new -key $PROXY_CLIENT_KEY -out $PROXY_CLIENT_CSR -subj '/CN=aggregator/O=system:masters'
openssl x509 -req -days 730 -in $PROXY_CLIENT_CSR -CA $PROXY_CRT -CAkey $PROXY_CA_KEY -set_serial 02 -out $PROXY_CLIENT_CRT

retrycmd_if_failure() { for i in 1 2 3 4 5 6 7 8 9 10; do $@; [ $? -eq 0  ] && break || sleep 30; done ; }

write_certs_to_disk() {
    etcdctl get $ETCD_REQUESTHEADER_CLIENT_CA > $K8S_PROXY_CA_CRT_FILEPATH
    etcdctl get $ETCD_PROXY_KEY > $K8S_PROXY_KEY_FILEPATH
    etcdctl get $ETCD_PROXY_CERT > $K8S_PROXY_CRT_FILEPATH
    # Remove whitespace padding at beginning of 1st line
    sed -i '1s/\s//' $K8S_PROXY_CA_CRT_FILEPATH $K8S_PROXY_CRT_FILEPATH $K8S_PROXY_KEY_FILEPATH
    chmod 600 $K8S_PROXY_KEY_FILEPATH
}

write_certs_to_disk_with_retry() {
    for i in 1 2 3 4 5 6 7 8 9 10 11 12; do
        write_certs_to_disk
        [ $? -eq 0  ] && break || sleep 5
    done
}

# block until all etcd is ready
retrycmd_if_failure etcdctl cluster-health
# Make etcd keys, adding a leading whitespace because etcd won't accept a val that begins with a '-' (hyphen)!
if etcdctl mk $ETCD_REQUESTHEADER_CLIENT_CA " $(cat ${PROXY_CRT})"; then
    etcdctl mk $ETCD_PROXY_KEY " $(cat ${PROXY_CLIENT_KEY})"
    etcdctl mk $ETCD_PROXY_CERT " $(cat ${PROXY_CLIENT_CRT})"
    sleep 5
    write_certs_to_disk_with_retry
# If the etcdtl mk command failed, that means the key already exists
else
    sleep 5
    write_certs_to_disk_with_retry
fi
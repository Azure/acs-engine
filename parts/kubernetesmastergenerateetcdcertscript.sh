#!/bin/bash

# ca cert
ETCD_CA_KEY="${ETCD_CA_KEY:=/tmp/etcd-ca.key}"
ETCD_CA_CRT="${ETCD_CA_CRT:=/tmp/etcd-ca.crt}"
K8S_ETCD_CA_CRT_FILEPATH="${K8S_ETCD_CA_CRT_FILEPATH:=/etc/kubernetes/certs/etcd-ca.crt}"

# etcd tls server certs
ETCD_SERVER_KEY="${ETCD_SERVER_KEY:=/tmp/etcd-server.key}"
ETCD_SERVER_CSR="${ETCD_SERVER_CSR:=/tmp/etcd-server.csr}"
ETCD_SERVER_CRT="${ETCD_SERVER_CRT:=/tmp/etcd-server.crt}"
K8S_ETCD_SERVER_CRT_FILEPATH="${K8S_ETCD_SERVER_CRT_FILEPATH:=/etc/kubernetes/certs/etcd-server.crt}"
K8S_ETCD_SERVER_KEY_FILEPATH="${K8S_ETCD_SERVER_KEY_FILEPATH:=/etc/kubernetes/certs/etcd-server.key}"

# etcd tls peer certs
ETCD_PEER_KEY="${ETCD_PEER_KEY:=/tmp/etcd-peer.key}"
ETCD_PEER_CSR="${ETCD_PEER_CSR:=/tmp/etcd-peer.csr}"
ETCD_PEER_CRT="${ETCD_PEER_CRT:=/tmp/etcd-peer.crt}"
K8S_ETCD_PEER_CRT_FILEPATH="${K8S_ETCD_PEER_CRT_FILEPATH:=/etc/kubernetes/certs/etcd-peer.crt}"
K8S_ETCD_PEER_KEY_FILEPATH="${K8S_ETCD_PEER_KEY_FILEPATH:=/etc/kubernetes/certs/etcd-peer.key}"

# etcd tls client certs
ETCD_CLIENT_KEY="${ETCD_CLIENT_KEY:=/tmp/etcd-client.key}"
ETCD_CLIENT_CSR="${ETCD_CLIENT_CSR:=/tmp/etcd-client.csr}"
ETCD_CLIENT_CRT="${ETCD_CLIENT_CRT:=/tmp/etcd-client.crt}"
K8S_ETCD_CLIENT_CRT_FILEPATH="${K8S_ETCD_CLIENT_CRT_FILEPATH:=/etc/kubernetes/certs/etcd-client.crt}"
K8S_ETCD_CLIENT_KEY_FILEPATH="${K8S_ETCD_CLIENT_KEY_FILEPATH:=/etc/kubernetes/certs/etcd-client.key}"

# etcd files for writing certs to disk
ETCD_REQUESTHEADER_CA="${ETCD_REQUESTHEADER_CA:=/etcdcerts/requestheader-etcd-ca-file}"
ETCD_SERVER_CERT_FILE="${ETCD_SERVER_CERT_FILE:=/etcdcerts/etcd-server-cert-file}"
ETCD_SERVER_KEY_FILE="${ETCD_SERVER_KEY_FILE:=/etcdcerts/etcd-server-key-file}"
ETCD_CLIENT_CERT_FILE="${ETCD_CLIENT_CERT_FILE:=/etcdcerts/etcd-client-cert-file}"
ETCD_CLIENT_KEY_FILE="${ETCD_CLIENT_KEY_FILE:=/etcdcerts/etcd-client-key-file}"
ETCD_PEER_CERT_FILE="${ETCD_PEER_CERT_FILE:=/etcdcerts/etcd-peer-cert-file}"
ETCD_PEER_KEY_FILE="${ETCD_PEER_KEY_FILE:=/etcdcerts/etcd-peer-key-file}"

echo subjectAltName = IP:127.0.0.1 > extfile.cnf

# generate root CA
openssl genrsa -out $ETCD_CA_KEY 2048
openssl req -new -x509 -days 1826 -key $ETCD_CA_KEY -out $ETCD_CA_CRT -subj '/CN=etcdCA'
# generate new cert
openssl genrsa -out $ETCD_SERVER_KEY 2048
openssl req -new -key $ETCD_SERVER_KEY -out $ETCD_SERVER_CSR -subj '/CN=127.0.0.1/O=system:masters'
openssl x509 -req -days 730 -in $ETCD_SERVER_CSR -CA $ETCD_CA_CRT -CAkey $ETCD_CA_KEY -set_serial 02 -out $ETCD_SERVER_CRT -extfile extfile.cnf
openssl genrsa -out $ETCD_CLIENT_KEY 2048
openssl req -new -key $ETCD_CLIENT_KEY -out $ETCD_CLIENT_CSR -subj '/CN=127.0.0.1/O=system:masters'
openssl x509 -req -days 730 -in $ETCD_CLIENT_CSR -CA $ETCD_CA_CRT -CAkey $ETCD_CA_KEY -set_serial 02 -out $ETCD_CLIENT_CRT -extfile extfile.cnf
openssl genrsa -out $ETCD_PEER_KEY 2048
openssl req -new -key $ETCD_PEER_KEY -out $ETCD_PEER_CSR -subj '/CN=127.0.0.1/O=system:masters'
openssl x509 -req -days 730 -in $ETCD_PEER_CSR -CA $ETCD_CA_CRT -CAkey $ETCD_CA_KEY -set_serial 02 -out $ETCD_PEER_CRT -extfile extfile.cnf

retrycmd_if_failure() { for i in 1 2 3 4 5 6 7 8 9 10; do $@; [ $? -eq 0  ] && break || sleep 30; done ; }

write_certs_to_disk() {
    etcdctl get $ETCD_REQUESTHEADER_CA > $K8S_ETCD_CA_CRT_FILEPATH
    etcdctl get $ETCD_SERVER_CERT_FILE > $K8S_ETCD_SERVER_CRT_FILEPATH
    etcdctl get $ETCD_SERVER_KEY_FILE > $K8S_ETCD_SERVER_KEY_FILEPATH
    etcdctl get $ETCD_CLIENT_CERT_FILE > $K8S_ETCD_CLIENT_CRT_FILEPATH
    etcdctl get $ETCD_CLIENT_KEY_FILE > $K8S_ETCD_CLIENT_KEY_FILEPATH
    etcdctl get $ETCD_PEER_CERT_FILE > $K8S_ETCD_PEER_CRT_FILEPATH
    etcdctl get $ETCD_PEER_KEY_FILE > $K8S_ETCD_PEER_KEY_FILEPATH
    # Remove whitespace padding at beginning of 1st line
    sed -i '1s/\s//' $K8S_ETCD_CA_CRT_FILEPATH $K8S_ETCD_SERVER_CRT_FILEPATH $K8S_ETCD_SERVER_KEY_FILEPATH $K8S_ETCD_CLIENT_CRT_FILEPATH $K8S_ETCD_CLIENT_KEY_FILEPATH $K8S_ETCD_PEER_CRT_FILEPATH $K8S_ETCD_PEER_KEY_FILEPATH
    chmod 600 $K8S_ETCD_SERVER_KEY_FILEPATH
    chmod 600 $K8S_ETCD_CLIENT_KEY_FILEPATH
    chmod 600 $K8S_ETCD_PEER_KEY_FILEPATH
    chown etcd:etcd $K8S_ETCD_SERVER_KEY_FILEPATH
    chown etcd:etcd $K8S_ETCD_CLIENT_KEY_FILEPATH
    chown etcd:etcd $K8S_ETCD_PEER_KEY_FILEPATH
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
if etcdctl mk $ETCD_REQUESTHEADER_CA " $(cat ${ETCD_CA_CRT})"; then
    etcdctl mk $ETCD_SERVER_KEY_FILE " $(cat ${ETCD_SERVER_KEY})"
    etcdctl mk $ETCD_SERVER_CERT_FILE " $(cat ${ETCD_SERVER_CRT})"
    etcdctl mk $ETCD_CLIENT_KEY_FILE " $(cat ${ETCD_CLIENT_KEY})"
    etcdctl mk $ETCD_CLIENT_CERT_FILE " $(cat ${ETCD_CLIENT_CRT})"
    etcdctl mk $ETCD_PEER_KEY_FILE " $(cat ${ETCD_PEER_KEY})"
    etcdctl mk $ETCD_PEER_CERT_FILE " $(cat ${ETCD_PEER_CRT})"
    sleep 5
    write_certs_to_disk_with_retry
# If the etcdtl mk command failed, that means the key already exists
else
    sleep 5
    write_certs_to_disk_with_retry
fi
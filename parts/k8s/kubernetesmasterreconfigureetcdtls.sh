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
ETCD_PEER_KEYS=("${ETCD_PEER_KEY_0:=/tmp/etcd-peer0.key}" "${ETCD_PEER_KEY_1:=/tmp/etcd-peer1.key}" "${ETCD_PEER_KEY_2:=/tmp/etcd-peer2.key}" "${ETCD_PEER_KEY_3:=/tmp/etcd-peer3.key}" "${ETCD_PEER_KEY_4:=/tmp/etcd-peer4.key}")
ETCD_PEER_CSRS=("${ETCD_PEER_CSR_0:=/tmp/etcd-peer0.csr}" "${ETCD_PEER_CSR_1:=/tmp/etcd-peer1.csr}" "${ETCD_PEER_CSR_2:=/tmp/etcd-peer2.csr}" "${ETCD_PEER_CSR_3:=/tmp/etcd-peer3.csr}" "${ETCD_PEER_CSR_4:=/tmp/etcd-peer4.csr}")
ETCD_PEER_CRTS=("${ETCD_PEER_CRT_0:=/tmp/etcd-peer0.crt}" "${ETCD_PEER_CRT_1:=/tmp/etcd-peer1.crt}" "${ETCD_PEER_CRT_2:=/tmp/etcd-peer2.crt}" "${ETCD_PEER_CRT_3:=/tmp/etcd-peer3.crt}" "${ETCD_PEER_CRT_4:=/tmp/etcd-peer4.crt}")
K8S_ETCD_PEER_CRT_FILEPATH="${K8S_ETCD_PEER_CRT_FILEPATH:=/etc/kubernetes/certs/etcd-peer${1}.crt}"
K8S_ETCD_PEER_KEY_FILEPATH="${K8S_ETCD_PEER_KEY_FILEPATH:=/etc/kubernetes/certs/etcd-peer${1}.key}"

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
ETCD_PEER_CERT_FILES=("${ETCD_PEER_CERT_FILE_0:=/etcdcerts/etcd-peer-cert-file-0}" "${ETCD_PEER_CERT_FILE_1:=/etcdcerts/etcd-peer-cert-file-1}" "${ETCD_PEER_CERT_FILE_2:=/etcdcerts/etcd-peer-cert-file-2}" "${ETCD_PEER_CERT_FILE_3:=/etcdcerts/etcd-peer-cert-file-3}" "${ETCD_PEER_CERT_FILE_4:=/etcdcerts/etcd-peer-cert-file-4}")
ETCD_PEER_KEY_FILES=("${ETCD_PEER_KEY_FILE_0:=/etcdcerts/etcd-peer-key-file-0}" "${ETCD_PEER_KEY_FILE_1:=/etcdcerts/etcd-peer-key-file-1}" "${ETCD_PEER_KEY_FILE_2:=/etcdcerts/etcd-peer-key-file-2}" "${ETCD_PEER_KEY_FILE_3:=/etcdcerts/etcd-peer-key-file-3}" "${ETCD_PEER_KEY_FILE_4:=/etcdcerts/etcd-peer-key-file-4}")

# generate root CA
openssl genrsa -out $ETCD_CA_KEY 2048
openssl req -new -x509 -days 1826 -key $ETCD_CA_KEY -out $ETCD_CA_CRT -subj '/CN=etcdCA'
# generate new cert
openssl genrsa -out $ETCD_SERVER_KEY 2048
openssl req -new -key $ETCD_SERVER_KEY -out $ETCD_SERVER_CSR -subj '/CN=127.0.0.1/O=system:masters'
openssl x509 -req -days 730 -in $ETCD_SERVER_CSR -CA $ETCD_CA_CRT -CAkey $ETCD_CA_KEY -set_serial 02 -out $ETCD_SERVER_CRT -extfile <(printf "subjectAltName=${2}")
openssl genrsa -out $ETCD_CLIENT_KEY 2048
openssl req -new -key $ETCD_CLIENT_KEY -out $ETCD_CLIENT_CSR -subj '/CN=127.0.0.1/O=system:masters'
openssl x509 -req -days 730 -in $ETCD_CLIENT_CSR -CA $ETCD_CA_CRT -CAkey $ETCD_CA_KEY -set_serial 02 -out $ETCD_CLIENT_CRT -extfile <(printf "subjectAltName=${2}")
for ((i = 0; i < ${#ETCD_PEER_KEYS[@]}; ++i)); do
    openssl genrsa -out ${ETCD_PEER_KEYS[$i]} 2048
    openssl req -new -key ${ETCD_PEER_KEYS[$i]} -out ${ETCD_PEER_CSRS[$i]} -subj '/CN=127.0.0.1/O=system:masters'
    openssl x509 -req -days 730 -in ${ETCD_PEER_CSRS[$i]} -CA $ETCD_CA_CRT -CAkey $ETCD_CA_KEY -set_serial 02 -out ${ETCD_PEER_CRTS[$i]} -extfile <(printf "subjectAltName=${2}")
done

retrycmd_if_failure() { for i in {1..10}; do $@; [ $? -eq 0  ] && break || sleep 30; done ; }

write_certs_to_disk() {
    etcdctl get $ETCD_REQUESTHEADER_CA > $K8S_ETCD_CA_CRT_FILEPATH
    etcdctl get $ETCD_SERVER_CERT_FILE > $K8S_ETCD_SERVER_CRT_FILEPATH
    etcdctl get $ETCD_SERVER_KEY_FILE > $K8S_ETCD_SERVER_KEY_FILEPATH
    etcdctl get $ETCD_CLIENT_CERT_FILE > $K8S_ETCD_CLIENT_CRT_FILEPATH
    etcdctl get $ETCD_CLIENT_KEY_FILE > $K8S_ETCD_CLIENT_KEY_FILEPATH
    etcdctl get ${ETCD_PEER_CERT_FILES[${1}]} > $K8S_ETCD_PEER_CRT_FILEPATH
    etcdctl get ${ETCD_PEER_KEY_FILES[${1}]} > $K8S_ETCD_PEER_KEY_FILEPATH
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
   for i in {1..12}; do
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
    for ((i = 0; i < ${#ETCD_PEER_KEY_FILES[@]}; ++i)); do
        etcdctl mk ${ETCD_PEER_KEY_FILES[$i]} " $(cat ${ETCD_PEER_KEYS[$i]})"
        etcdctl mk ${ETCD_PEER_CERT_FILES[$i]} " $(cat ${ETCD_PEER_CRTS[$i]})"
    done
    sleep 5
    write_certs_to_disk_with_retry
# If the etcdtl mk command failed, that means the key already exists
else
    sleep 5
    write_certs_to_disk_with_retry
fi

cat /tmp/etcdtls > /etc/default/etcd
MEMBER="$(etcdctl member list | grep -E ${4} | cut -d':' -f 1)"
echo ${MEMBER} ${3} >> /opt/etcdtls
sleep 60 #TODO: fix this 
etcdctl member update ${MEMBER} ${3}
sed -i "11iEnvironment=ETCD_CA_FILE=$K8S_ETCD_CA_CRT_FILEPATH" /etc/systemd/system/etcd.service
sed -i "11iEnvironment=ETCD_CERT_FILE=$K8S_ETCD_CLIENT_CRT_FILEPATH" /etc/systemd/system/etcd.service
sed -i "11iEnvironment=ETCD_KEY_FILE=$K8S_ETCD_CLIENT_KEY_FILEPATH" /etc/systemd/system/etcd.service
sed -i "11iEnvironment=ETCD_ENDPOINTS=https://127.0.0.1:2379" /etc/systemd/system/etcd.service

#ETCDCTL_ENDPOINTS=https://127.0.0.1:2379
# azureuser@k8s-master-19135580-0:~$ ETCDCTL_CA_FILE=/etc/kubernetes/certs/etcd-ca.crt
# azureuser@k8s-master-19135580-0:~$ ETCDCTL_KEY_FILE=/etc/kubernetes/certs/etcd-client.key
# azureuser@k8s-master-19135580-0:~$ ETCDCTL_CERT_FILE=/etc/kubernetes/certs/etcd-client.crt

systemctl daemon-reload
systemctl restart etcd
rm /tmp/etcd*
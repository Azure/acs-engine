#!/bin/bash

#TODO PR1406

# generate root CA
openssl genrsa -out proxy-client-ca.key 2048
openssl req -new -x509 -days 1826 -key proxy-client-ca.key -out proxy-client-ca.crt -subj '/CN=proxyClientCA'
# generate new cert
openssl genrsa -out proxy-client.key 2048
openssl req -new -key proxy-client.key -out proxy-client.csr -subj '/CN=aggregator/O=system:masters'
openssl x509 -req -days 730 -in proxy-client.csr -CA proxy-client-ca.crt -CAkey proxy-client-ca.key -set_serial 02 -out proxy-client.crt

retrycmd_if_failure() { for i in 1 2 3 4 5 6 7 8 9 10; do $@; [ $? -eq 0  ] && break || sleep 30; done ; }

write_certs_to_disk() {
    etcdctl get /proxycerts/requestheader-client-ca-file > /etc/kubernetes/certs/proxy-ca.crt
    etcdctl get /proxycerts/proxy-client-cert-file > /etc/kubernetes/certs/proxy.crt
    etcdctl get /proxycerts/proxy-client-key-file > /etc/kubernetes/certs/proxy.key
}

# block until all etcd is ready
retrycmd_if_failure etcdctl cluster-health
if etcdctl mk /proxycerts/requestheader-client-ca-file " $(echo $(cat proxy-client-ca.key))"; then
    etcdctl mk /proxycerts/proxy-client-key-file " $(echo $(cat proxy-client.key))"
    etcdctl mk /proxycerts/proxy-client-cert-file " $(echo $(cat proxy-client.crt))"
    write_certs_to_disk
else
    sleep 30
    write_certs_to_disk
fi
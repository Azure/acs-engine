#!/bin/bash

#TODO write script

# PSEUDO CODE
# block until all etcd is ready
# for each cert:
#  get cert from central datastore on cluster (etcd?)
#  if success
#    cert = get response
#    write cert payload to local disk at appropriate filepath
#  if fail
#    generate new cert
#    put new cert if not exist into central datastore (etcdctl mk /example/key data)
#    if success
#      cert = new cert
#      write cert payload to local disk at appropriate filepath
#    if fail
#    retry everything in this enumeration
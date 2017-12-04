#!/bin/bash

set -e

# complete setting up custom VNET
examples/vnet/k8s-vnet-predeploy.sh

# scale
examples/azure-cni/k8s-scale.sh

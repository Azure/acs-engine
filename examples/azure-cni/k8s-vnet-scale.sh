#!/bin/bash

set -e

# complete setting up custom VNET
examples/vnet/k8s-vnet-postdeploy.sh

# scale
examples/azure-cni/k8s-scale.sh

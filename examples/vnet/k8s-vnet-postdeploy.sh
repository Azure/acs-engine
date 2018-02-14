#!/bin/bash

rt=$(az network route-table list -g ${RESOURCE_GROUP} | jq -r '.[].id')

az network vnet subnet update -n KubernetesSubnet -g ${RESOURCE_GROUP} --vnet-name KubernetesCustomVNET --route-table $rt

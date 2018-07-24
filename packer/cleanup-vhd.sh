#!/bin/bash -eux

## Cleanup packer SSH key and machine ID generated for this boot
rm /root/.ssh/authorized_keys
rm /home/packer/.ssh/authorized_keys
rm /etc/machine-id
touch /etc/machine-id
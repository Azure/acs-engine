#!/bin/bash
set -x
source /opt/azure/containers/provision_source.sh

sdName="<searchDomainName>"
sdRealmUser=$"<searchDomainRealmUser>"
sdRealmPassword=$"<searchDomainRealmPassword>"
sdComputerOU=$"<searchDomainComputerOU>"
ucDomainName=$(echo "${sdName}" | tr /a-z/ /A-Z/)

computerOUSwitch=""
if [[ ! -z "${sdComputerOU}" ]]; then
  computerOUSwitch="--computer-ou='${sdComputerOU}'"
fi

echo "  dns-search ${sdName}" >> /etc/network/interfaces.d/50-cloud-init.cfg
systemctl_restart 20 5 10 restart networking

retrycmd_if_failure 10 5 120 apt-get update
retrycmd_if_failure 10 5 120 apt-get -y install \
  realmd \
  sssd \
  sssd-tools \
  samba-common \
  samba \
  python2.7 \
  samba-libs \
  packagekit

echo "${sdRealmPassword}" | \
  realm join -U "${sdRealmUser}@${ucDomainName}" ${ucDomainName} ${computerOUSwitch}

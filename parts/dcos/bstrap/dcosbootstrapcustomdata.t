bootcmd:
- bash -c "if [ ! -f /var/lib/sdb-gpt ];then echo DCOS-5890;parted -s /dev/sdb mklabel
  gpt;touch /var/lib/sdb-gpt;fi"
disk_setup:
  ephemeral0:
    layout:
    - 50
    - 50
    overwrite: true
    table_type: gpt
fs_setup:
- device: ephemeral0.1
  filesystem: ext4
  overwrite: true
- device: ephemeral0.2
  filesystem: ext4
  overwrite: true
mounts:
- - ephemeral0.1
  - /var/lib/mesos
- - ephemeral0.2
  - /var/lib/docker
runcmd: PREPROVISION_EXTENSION
    - [ ln, -s, /bin/rm, /usr/bin/rm ]
    - [ ln, -s, /bin/mkdir, /usr/bin/mkdir ]
    - [ ln, -s, /bin/tar, /usr/bin/tar ]
    - [ ln, -s, /bin/ln, /usr/bin/ln ]
    - [ ln, -s, /bin/cp, /usr/bin/cp ]
    - [ ln, -s, /bin/systemctl, /usr/bin/systemctl ]
    - [ ln, -s, /bin/mount, /usr/bin/mount ]
    - [ ln, -s, /bin/bash, /usr/bin/bash ]
    - [ ln, -s, /usr/sbin/useradd, /usr/bin/useradd ]
    - [ systemctl, disable, --now, resolvconf.service ]
    - [ systemctl, mask, --now, lxc-net.service ]
    - [ systemctl, disable, --now, unscd.service ]
    - [ systemctl, stop, --now, unscd.service ]
    - /opt/azure/containers/provision.sh
write_files:
- content: |
    bootstrap_url: {{{dcosBootstrapURL}}}
    cluster_name: azure-dcos
    exhibitor_storage_backend: static
    master_discovery: static
    ip_detect_public_filename: /opt/azure/genconf/ip-detect.sh
    master_list:
MASTER_IP_LIST
    resolvers:
    - 198.51.100.1
    - 198.51.100.2
    - 198.51.100.3
  owner: root
  path: /opt/azure/genconf/config.yaml
  permissions: '0644'
- content: |
    #!/bin/sh
    curl -H Metadata:true "http://169.254.169.254/metadata/instance/network/interface/0/ipv4/ipAddress/0/privateIpAddress?api-version=2017-08-01&format=text"
  owner: root
  path: /opt/azure/genconf/ip-detect.sh
  permissions: '0755'

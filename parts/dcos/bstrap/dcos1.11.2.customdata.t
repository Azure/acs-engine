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
- ln -s /bin/rm /usr/bin/rm
- ln -s /bin/mkdir /usr/bin/mkdir
- ln -s /bin/tar /usr/bin/tar
- ln -s /bin/ln /usr/bin/ln
- ln -s /bin/cp /usr/bin/cp
- ln -s /bin/systemctl /usr/bin/systemctl
- ln -s /bin/mount /usr/bin/mount
- ln -s /bin/bash /usr/bin/bash
- ln -s /usr/sbin/useradd /usr/bin/useradd
- systemctl disable --now resolvconf.service
- systemctl mask --now lxc-net.service
- systemctl disable --now unscd.service
- systemctl stop --now unscd.service
- /opt/azure/containers/provision.sh
- bash /tmp/dcos/dcos_install.sh ROLENAME
- /opt/azure/dcos/postinstall-cond.sh
write_files:
- content: |
    [Service]
    Restart=always
    StartLimitInterval=0
    RestartSec=15
    ExecStartPre=-/sbin/ip link del docker0
    ExecStart=
    ExecStart=/usr/bin/dockerd --storage-driver=overlay
  path: /etc/systemd/system/docker.service.d/execstart.conf
  permissions: '0644'
- content: |
    [Unit]
    PartOf=docker.service
    [Socket]
    ListenStream=/var/run/docker.sock
    SocketMode=0660
    SocketUser=root
    SocketGroup=docker
    ListenStream=2375
    BindIPv6Only=both
    [Install]
    WantedBy=sockets.target
  path: /etc/systemd/system/docker.socket
  permissions: '0644'
- content: |
    DCOS_ENVIRONMENT={{{targetEnvironment}}}
  owner: root
  path: /opt/azure/dcos/environment
  permissions: '0644'
- path: /var/lib/dcos/mesos-slave-common
  content: 'ATTRIBUTES_STR'
  permissions: "0644"
  owner: "root"
- content: 'PROVISION_SOURCE_STR'
  path: /opt/azure/containers/provision_source.sh
  permissions: "0744"
  owner: "root"
- content: 'PROVISION_STR'
  path: /opt/azure/containers/provision.sh
  permissions: "0744"
  owner: "root"
- content: |
    #!/bin/bash
    if [ -f /opt/azure/dcos/postinstall.sh ]; then /opt/azure/dcos/postinstall.sh; fi
  path: /opt/azure/dcos/postinstall-cond.sh
  permissions: "0744"
  owner: "root"

bootcmd:
- bash -c "if [ ! -f /var/lib/sdb-gpt ];then echo DCOS-5890;parted -s /dev/sdb mklabel gpt;touch /var/lib/sdb-gpt;fi"
- bash -c "if [ ! -f /var/lib/sdc-gpt ];then echo DCOS-5890;parted -s /dev/sdc mklabel gpt&&touch /var/lib/sdc-gpt;fi"
- bash -c "if [ ! -f /var/lib/sdd-gpt ];then echo DCOS-5890;parted -s /dev/sdd mklabel gpt&&touch /var/lib/sdd-gpt;fi"
- bash -c "if [ ! -f /var/lib/sde-gpt ];then echo DCOS-5890;parted -s /dev/sde mklabel gpt&&touch /var/lib/sde-gpt;fi"
- bash -c "if [ ! -f /var/lib/sdf-gpt ];then echo DCOS-5890;parted -s /dev/sdf mklabel gpt&&touch /var/lib/sdf-gpt;fi"
- bash -c "mkdir -p /dcos/volume{0,1,2,3}"
disk_setup:
  ephemeral0:
    layout:
    - 45
    - 45
    - 10
    overwrite: true
    table_type: gpt
  /dev/sdc:
    layout: true
    overwrite: true
    table_type: gpt
  /dev/sdd:
    layout: true
    overwrite: true
    table_type: gpt
  /dev/sde:
    layout: true
    overwrite: true
    table_type: gpt
  /dev/sdf:
    layout: true
    overwrite: true
    table_type: gpt
fs_setup:
- device: ephemeral0.1
  filesystem: ext4
  overwrite: true
- device: ephemeral0.2
  filesystem: ext4
  overwrite: true
- device: ephemeral0.3
  filesystem: ext4
  overwrite: true
- device: /dev/sdc1
  filesystem: ext4
  overwrite: true
- device: /dev/sdd1
  filesystem: ext4
  overwrite: true
- device: /dev/sde1
  filesystem: ext4
  overwrite: true
- device: /dev/sdf1
  filesystem: ext4
  overwrite: true
mounts:
- - ephemeral0.1
  - /var/lib/mesos
- - ephemeral0.2
  - /var/lib/docker
- - ephemeral0.3
  - /var/tmp
- - /dev/sdc1
  - /dcos/volume0
- - /dev/sdd1
  - /dcos/volume1
- - /dev/sde1
  - /dcos/volume2
- - /dev/sdf1
  - /dcos/volume3
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
    - service ssh restart
    - [ cp, -p, /etc/resolv.conf, /tmp/resolv.conf ]
    - [ rm, -f, /etc/resolv.conf ]
    - [ cp, -p, /tmp/resolv.conf, /etc/resolv.conf ]
    - [ systemctl, start, dcos-docker-install.service ]
    - [ systemctl, start, dcos-config-writer.service ]
    - [ systemctl restart systemd-journald.service ]
    - [ systemctl restart docker.service ]
    - [ systemctl start dcos-link-env.service ]
    - [ systemctl enable dcos-setup.service ]
    - [ systemctl --no-block start dcos-setup.service ]

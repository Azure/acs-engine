bootcmd:
- bash -c "if [ ! -f /var/lib/sdb-gpt ];then echo DCOS-5890;parted -s /dev/sdb mklabel
  gpt;touch /var/lib/sdb-gpt;fi"
- bash -c "if [ ! -f /var/lib/sdc-gpt ];then echo DCOS-5890;parted -s /dev/sdc mklabel
  gpt&&touch /var/lib/sdc-gpt;fi"
- bash -c "if [ ! -f /var/lib/sdd-gpt ];then echo DCOS-5890;parted -s /dev/sdd mklabel
  gpt&&touch /var/lib/sdd-gpt;fi"
- bash -c "if [ ! -f /var/lib/sde-gpt ];then echo DCOS-5890;parted -s /dev/sde mklabel
  gpt&&touch /var/lib/sde-gpt;fi"
- bash -c "if [ ! -f /var/lib/sdf-gpt ];then echo DCOS-5890;parted -s /dev/sdf mklabel
  gpt&&touch /var/lib/sdf-gpt;fi"
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
runcmd:
- - ln
  - -s
  - /bin/rm
  - /usr/bin/rm
- - ln
  - -s
  - /bin/mkdir
  - /usr/bin/mkdir
- - ln
  - -s
  - /bin/tar
  - /usr/bin/tar
- - ln
  - -s
  - /bin/ln
  - /usr/bin/ln
- - ln
  - -s
  - /bin/cp
  - /usr/bin/cp
- - ln
  - -s
  - /bin/systemctl
  - /usr/bin/systemctl
- - ln
  - -s
  - /bin/mount
  - /usr/bin/mount
- - ln
  - -s
  - /bin/bash
  - /usr/bin/bash
- - ln
  - -s
  - /usr/sbin/useradd
  - /usr/bin/useradd
- - systemctl
  - disable
  - --now
  - resolvconf.service
- - systemctl
  - mask
  - --now
  - lxc-net.service
- - systemctl
  - disable
  - --now
  - unscd.service
- - systemctl
  - stop
  - --now
  - unscd.service
- /opt/azure/containers/provision.sh
- - cp
  - -p
  - /etc/resolv.conf
  - /tmp/resolv.conf
- - rm
  - -f
  - /etc/resolv.conf
- - cp
  - -p
  - /tmp/resolv.conf
  - /etc/resolv.conf
- - systemctl
  - start
  - dcos-docker-install.service
- - systemctl
  - start
  - dcos-config-writer.service
- - systemctl
  - restart
  - systemd-journald.service
- - systemctl
  - restart
  - docker.service
- - systemctl
  - start
  - dcos-link-env.service
- - systemctl
  - enable
  - dcos-setup.service
- - systemctl
  - --no-block
  - start
  - dcos-setup.service
write_files:
- content: 'https://dcosio.azureedge.net/dcos/stable/1.10.0

'
  owner: root
  path: /etc/mesosphere/setup-flags/repository-url
  permissions: '0644'
- content: '["adminrouter--1166a3736442e7963a68d1d644bf5f54ca3cb01d", "avro-cpp--9cb0ee14e3cd5bbdb171efcc72a84d16862ea02d",
    "boost-libs--8d515c2f703c666ae1b6c5ccc35cc0f8fa36677f", "bootstrap--c1bc86593e212cf9fe83db2246bacd129a6b3adc",
    "boto--3890cb2817c00b874ba033abe784b5b343caa3c7", "check-time--79e3f6ab99125471e1d94d5f6bc0fea88446831c",
    "cni--7a8572e385c3f5262945c52c8003d1bbb22cf7aa", "cosmos--e84c5bf3259405df90d682536ba445cc4839a324",
    "curl--17866a8ae9305826aa5f357a09db2c1f2b2c2ad0", "dcos-checks--8fd33919e6f163dba1bd13e4c7e4e0523919a719",
    "dcos-cni--12a77c1e9bebd4cbd600524a864c2bd8483330d3", "dcos-config--setup_DCOSGUID",
    "dcos-diagnostics--e3b557b0ec8e98617d0cd0fdf136ef9dded96316", "dcos-history--23de88ddc1a5f9018dd11b279c5be6a768a18de4",
    "dcos-image--df630d8e930d6650ce3d0ade519660142233d862", "dcos-image-deps--81d23d00b1acddb316c9b15fd8499c2b10f6b697",
    "dcos-integration-test--9ec173650d4e73ba494603324e7583d23970e4b8", "dcos-log--d2af4b1a47d3755a51823e95fbc6c366cf0f9269",
    "dcos-metadata--setup_DCOSGUID", "dcos-metrics--2a26c0b50b0b6564f86c48d50aa86f681c9af93c",
    "dcos-oauth--445bb1388670981c6acc667b2529fc32d4c1fbd4", "dcos-signal--4366023212ea49a64c5c9aef1965e5a3133c4b61",
    "dcos-test-utils--1066d896d25f4c1e3f6d9a5e7f9c1c6e8c675bb7", "dcos-ui--cc2e3d26537ea190efacd6f899dd4cc2210d45b7",
    "dnspython--0be432372a3820eafcfa66975943c9536dbe1164", "docker-gc--89f5535aea154dca504f84cd60eac6f61836aef9",
    "dvdcli--ee85411e3cb9f0988ed54b5cc0789172b887f12f", "erlang--d693172f6f033707c7f07ff78fc18ac543d66b41",
    "exhibitor--c3e48bbae19c0ed9c30d7f9396305d1e77130658", "flask--6d0f985ad677e8422c7190cbe207424acd813c3b",
    "java--ce5ff19502fca31eaf4a9af86d50a10a8c212a5b", "libevent--05dc18bc0ab7434b2738318c5ebaa2e61a311f50",
    "libffi--0e5b99b94f296b2a9a1b75e9fe5f74f5446f5e9b", "libsodium--e7056355f1fe160ade83aac0d11352a2bf3844e6",
    "logrotate--877aece1fd506af3b9167b6938c316adfa79d4f5", "marathon--accdc43bafeca02da1be340baba4b55011eadf63",
    "mesos--0677ce2b7d2e8c45091f6481884542f1f765c3d5", "mesos-dns--600da87080b7634f2380594499004a7ff0b34662",
    "mesos-modules--1f5c4860450949db92ed27326c3146526041e681", "metronome--2ec6f56be44ed822e7228cb66c4dae6a78345789",
    "navstar--c66f92f01d837433de3e2b19d221c64d26cc54b1", "ncurses--030fd6b08ed46a7ecce001c36901f5b4ad5d2af5",
    "octarine--4e37c062d2f145f9c2ce01d30dadf72c2aac5c4a", "openssl--44777d19d54a3c33cc19543f2201cb20bf085d98",
    "pkgpanda-api--30cb1e68f92ed5d4b89d57ca526f8a69b44132c8", "pkgpanda-role--612a6734567cc0c7c2ae1d508f03172f4bc7beed",
    "pytest--5e26c8ed9fd2c325672d56fe558299bfbd0f7018", "python--5a4285ff7296548732203950bf73d360ea67f6ab",
    "python-azure-mgmt-resource--26cbe8349f3fe139f7dc8bff7f0cb735382314fc", "python-cryptography--0d83d8afef4a8faddf0d8b713619d9d76e510a9e",
    "python-dateutil--519201adebeba186049ecd79a9f358f614173b10", "python-docopt--0af809c220a922f7f6c58f15beafebaa043477c7",
    "python-gunicorn--2ceb53716237da0736f67f4004682083f6ac68e1", "python-isodate--c9efb5859a0cfb06d82f25220cc5b387914af85d",
    "python-jinja2--601a1443aa4c649ab1da10c2a6d7a4477a263fb3", "python-kazoo--0ff8e6ef528f58c6f36f0a9df6dc27d3871e5c27",
    "python-markupsafe--1388c95920b4eb920c7a753d620a1ad07fc8b64d", "python-passlib--4691268be760073188b555dc436f836c6706b37a",
    "python-pyyaml--d8a775d6e43da5eb239af5cccdf1d3fceeb0335f", "python-requests--db0474fab16019ba29a609a354285f221c1a2859",
    "python-retrying--37dd25bf69bcbefe0c50139085d6bb2e22ccf439", "python-tox--322c468e2a75c5b143cb06af460b5e801ee34342",
    "rexray--da7f17f8a4b772c0bac3f8d289a08abd4ff272b4", "six--93734bac9907087744815f9cb5b6152e9a198fae",
    "spartan--c3d8005b1340bcbc3a00496861745b2d0bb2d697", "strace--9be573456909e3931a890785eb6474af7e0dcce4",
    "teamcity-messages--073793b16cf369e58ebdb6348b93ed14b0e5e59a", "toybox--0c49f879bfe2f99e6f99b397136894fa5096fa0c"]

'
  owner: root
  path: /etc/mesosphere/setup-flags/cluster-packages.json
  permissions: '0644'
- content: |
    [Journal]
    MaxLevelConsole=warning
    RateLimitInterval=1s
    RateLimitBurst=20000
  owner: root
  path: /etc/systemd/journald.conf.d/dcos.conf
  permissions: '0644'
- content: |
    rexray:
      loglevel: info
      modules:
        default-docker:
          disabled: true
      service: vfs
  path: /etc/rexray/config.yml
  permissions: '0644'
- content: |
    [Unit]
    After=network-online.target
    Wants=network-online.target
    [Service]
    Type=oneshot
    Environment=DEBIAN_FRONTEND=noninteractive
    StandardOutput=journal+console
    StandardError=journal+console
    ExecStartPre=/usr/bin/curl -fLsSv --retry 20 -Y 100000 -y 60 -o /var/tmp/d.deb https://mesosphere.blob.core.windows.net/dcos-deps/docker-engine_1.13.1-0-ubuntu-xenial_amd64.deb
    ExecStart=/usr/bin/bash -c "try=1;until dpkg -D3 -i /var/tmp/d.deb || ((try>9));do echo retry $((try++));sleep $((try*try));done;systemctl --now start docker;systemctl restart docker.socket"
  path: /etc/systemd/system/dcos-docker-install.service
  permissions: '0644'
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
      [Unit]
      Requires=dcos-setup.service
      After=dcos-setup.service
      [Service]
      Type=oneshot
      EnvironmentFile=/etc/environment
      EnvironmentFile=/opt/mesosphere/environment
      ExecStart=/usr/bin/bash -c "echo $(detect_ip) $(hostname) > /etc/hosts"
  path: /etc/systemd/system/dcos-config-writer.service
  permissions: '0644'
- content: |
    "bound_values":
      "adminrouter_auth_enabled": |-
        {{{oauthEnabled}}}
      "cluster_name": |-
        {{{masterPublicIPAddressName}}}
      "exhibitor_azure_account_key": |-
        ', listKeys(resourceId('Microsoft.Storage/storageAccounts', variables('masterStorageAccountExhibitorName')), '2015-06-15').key1, '
      "exhibitor_azure_account_name": |-
        {{{masterStorageAccountExhibitorName}}}
      "exhibitor_azure_prefix": |-
        {{{masterPublicIPAddressName}}}
      "master_list": |-
        ["', DCOSCUSTOMDATAPUBLICIPSTR'"]
      "oauth_enabled": |-
        {{{oauthEnabled}}}
    "late_bound_package_id": |-
      dcos-provider-DCOSGUID-azure--setup
  owner: root
  path: /etc/mesosphere/setup-flags/late-config.yaml
  permissions: '0644'
- content: |
    [Unit]
    Before=dcos.target
    [Service]
    Type=oneshot
    StandardOutput=journal+console
    StandardError=journal+console
    ExecStartPre=/usr/bin/mkdir -p /etc/profile.d
    ExecStart=/usr/bin/ln -sf /opt/mesosphere/bin/add_dcos_path.sh /etc/profile.d/dcos.sh
  path: /etc/systemd/system/dcos-link-env.service
  permissions: '0644'
- content: |
    [Unit]
    Description=Pkgpanda: Download DC/OS to this host.
    After=network-online.target
    Wants=network-online.target
    ConditionPathExists=!/opt/mesosphere/
    [Service]
    Type=oneshot
    StandardOutput=journal+console
    StandardError=journal+console
    ExecStartPre=/usr/bin/curl --keepalive-time 2 -fLsSv --retry 20 -Y 100000 -y 60 -o //var/tmp/bootstrap.tar.xz https://dcosio.azureedge.net/dcos/stable/1.10.0/bootstrap/4d92536e7381176206e71ee15b5ffe454439920c.bootstrap.tar.xz
    ExecStartPre=/usr/bin/mkdir -p /opt/mesosphere
    ExecStart=/usr/bin/tar -axf //var/tmp/bootstrap.tar.xz -C /opt/mesosphere
    ExecStartPost=-/usr/bin/rm -f //var/tmp/bootstrap.tar.xz
  path: /etc/systemd/system/dcos-download.service
  permissions: '0644'
- content: |
    [Unit]
    Description=Pkgpanda: Specialize DC/OS for this host.
    Requires=dcos-download.service
    After=dcos-download.service
    [Service]
    Type=oneshot
    StandardOutput=journal+console
    StandardError=journal+console
    EnvironmentFile=/opt/mesosphere/environment
    ExecStart=/opt/mesosphere/bin/pkgpanda setup --no-block-systemd
    [Install]
    WantedBy=multi-user.target
  path: /etc/systemd/system/dcos-setup.service
  permissions: '0644'
- content: ''
  path: /etc/mesosphere/roles/azure
- path: /var/lib/dcos/mesos-slave-common
  content: 'ATTRIBUTES_STR'
  permissions: "0644"
  owner: "root"
- content: 'PROVISION_STR'
  path: /opt/azure/containers/provision.sh
  permissions: "0744"
  owner: "root"

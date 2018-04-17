agent_list:
{{getDCOSAgentList() .}}
# DC/OS Enterprise Only
auth_cookie_secure_flag: `<true|false>`
bootstrap_url: {{getDCOSBootstrapURL .}}
# DC/OS Enterprise Only
bouncer_expiration_auth_token_days: `<time>`
cluster_docker_credentials:
  auths:
    '<path-to-credentials>':
      auth: <username>
      email: <email>
  cluster_docker_credentials_dcos_owned: <true|false>
    cluster_docker_credentials_write_to_etc: <true|false>
cluster_docker_credentials_enabled: <true|false>
cluster_docker_registry_url: <url>
cluster_name: {{getDCOSClusterName .}}
cosmos_config:
staged_package_storage_uri: <temp-path-to-files>
package_storage_uri: <permanent-path-to-files>
# DC/OS Enterprise Only
ca_certificate: <path-to-certificate>
ca_certificate_key: <path-to-private-key>
ca_certificate_chain: <path-to-certificate-chain>
customer_key: <customer-key>
custom_checks:
  cluster_checks:
    custom-check-1:
      description: Foobar cluster service is healthy
      cmd:
        - echo
        - hello
      timeout: 1s
  node_checks:
    checks:
      custom-check-2:
        description: Foobar node service is healthy
        cmd:
          - echo
          - hello
        timeout: 1s
        roles:
          - agent
    poststart:
      - custom-check-2
dcos_overlay_enable: `<true|false>`
dcos_overlay_config_attempts: <num-failed-attempts>
dcos_overlay_mtu: <mtu>
dcos_overlay_network:
  vtep_subnet: <address>
  vtep_mac_oui: <mac-address>
  overlays:
    - name: <name>
      subnet: <address>
      prefix: <size>
dns_search: <domain1 domain2 domain3>  
docker_remove_delay: <num>hrs
enable_docker_gc: `<true|false>`
exhibitor_storage_backend: static
exhibitor_storage_backend: zookeeper
exhibitor_zk_hosts: `<list-of-ip-port>`
exhibitor_zk_path: <filepath-to-data>
exhibitor_storage_backend: aws_s3
aws_access_key_id: <key-id>
aws_region: <bucket-region>
aws_secret_access_key: <secret-access-key>
exhibitor_explicit_keys: <true|false>
s3_bucket: <s3-bucket>
s3_prefix: <s3-prefix>
exhibitor_storage_backend: azure
exhibitor_azure_account_name: <storage-account-name>
exhibitor_azure_account_key: <storage-account-key>
exhibitor_azure_prefix: <blob-prefix>
gc_delay: <num>days
log_directory: `<path-to-install-logs>`
master_discovery: static
master_list:
{{getDcosMasterIPList .}}
master_discovery: master_http_loadbalancer
exhibitor_address: <loadbalancer-ip>
master_dns_bindall: `<true|false>`
num_masters: <num-of-masters>
# DC/OS only
oauth_enabled: `<true|false>`  
public_agent_list:
{{getDcosPublicAgentList . }}
platform: <platform>
process_timeout: <num-seconds>
rexray_config:
    rexray:
      loglevel:
      service:
    libstorage:
      integration:
        volume:
          operations:
            unmount:
              ignoreusedcount:
      server:
        tasks:
          logTimeout: 5m
# DC/OS Enterprise Only
security: <security-mode>
# DC/OS Enterprise Only
superuser_username: <username>
ssh_key_path: <path-to-ssh-key>
ssh_port: '<port-number>'
ssh_user: <username>
# DC/OS Enterprise Only
superuser_password_hash: <hashed-password>
# DC/OS Enterprise Only
superuser_username: <username>
telemetry_enabled: `<true|false>`
use_proxy: `<true|false>`
http_proxy: http://<proxy_host>:<http_proxy_port>
https_proxy: https://<proxy_host>:<https_proxy_port>
no_proxy:
- '<blocked.address1.com>'
- '<blocked.address2.com>'
# DC/OS Enterprise Only
zk_super_credentials: 'super:<long, random string>'
zk_master_credentials: 'dcos-master:<long, random string>'
zk_agent_credentials: 'dcos-agent:<long, random string>'

package dcosupgrade

import (
	"fmt"
	"strings"

	"github.com/Azure/acs-engine/pkg/operations"
)

var bootstrapNodeConfigScript = `#!/bin/bash

echo "Starting upgrade configuration"
if [ ! -e /opt/azure/dcos/upgrade/DCOS_VERSION/upgrade_url ]; then
  echo "Setting up bootstrap node"
  rm -rf /opt/azure/dcos/upgrade/DCOS_VERSION
  mkdir -p /opt/azure/dcos/upgrade/DCOS_VERSION/genconf
  cp /opt/azure/dcos/genconf/config.yaml /opt/azure/dcos/genconf/ip-detect /opt/azure/dcos/upgrade/DCOS_VERSION/genconf/
  cd /opt/azure/dcos/upgrade/DCOS_VERSION/
  curl -s -O https://dcos-mirror.azureedge.net/dcos-DCOS_DASHED_VERSION/dcos_generate_config.sh
  bash dcos_generate_config.sh --generate-node-upgrade-script DCOS_VERSION | tee /opt/azure/dcos/upgrade/DCOS_VERSION/log
  process=\$(docker ps -f ancestor=nginx -q)
  if [ ! -z "\$process" ]; then
    echo "Stopping nginx service \$process"
    docker kill \$process
  fi
  echo "Starting nginx service \$process"
  docker run -d -p 8086:80 -v \$PWD/genconf/serve:/usr/share/nginx/html:ro nginx
  docker ps
  grep 'Node upgrade script URL' /opt/azure/dcos/upgrade/DCOS_VERSION/log | awk -F ': ' '{print \$2}' | cat > /opt/azure/dcos/upgrade/DCOS_VERSION/upgrade_url
fi
upgrade_url=\$(cat /opt/azure/dcos/upgrade/DCOS_VERSION/upgrade_url)
if [ -z \${upgrade_url} ]; then
  rm -f /opt/azure/dcos/upgrade/DCOS_VERSION/upgrade_url
  echo "Failed to set up bootstrap node. Please try again"
  exit 1
else
  echo "Setting up bootstrap node completed. Node upgrade script URL \${upgrade_url}"
fi
`

var clusterNodeUpgradeScript = `#!/bin/bash

echo "Starting node upgrade"
mkdir -p /opt/azure/dcos/upgrade/DCOS_VERSION
cd /opt/azure/dcos/upgrade/DCOS_VERSION
curl -O \$(cat /opt/azure/dcos/upgrade/DCOS_VERSION/upgrade_url)
bash ./dcos_node_upgrade.sh

`

func (uc *UpgradeCluster) setBootstrapNode() error {
	uc.Logger.Info("Configuring bootstrap node for upgrade")
	ip := "23.101.124.31"

	orchestratorVersion := uc.ClusterTopology.DataModel.Properties.OrchestratorProfile.OrchestratorVersion
	dashedVersion := strings.Replace(orchestratorVersion, ".", "-", -1)
	configScript := strings.Replace(bootstrapNodeConfigScript, "DCOS_VERSION", orchestratorVersion, -1)
	configScript = strings.Replace(configScript, "DCOS_DASHED_VERSION", dashedVersion, -1)

	// copy the script over to the bootstrap node
	out, err := operations.RemoteRun("azureuser", ip, "private_key", fmt.Sprintf("cat << END > run_upgrade.sh\n%s\nEND\n", configScript))
	if err != nil {
		uc.Logger.Errorf(out)
		return err
	}
	// set script permissions
	out, err = operations.RemoteRun("azureuser", ip, "private_key", "chmod 755 ./run_upgrade.sh")
	if err != nil {
		uc.Logger.Errorf(out)
		return err
	}
	// run the script
	out, err = operations.RemoteRun("azureuser", ip, "private_key", "sudo ./run_upgrade.sh")
	if err != nil {
		uc.Logger.Errorf(out)
		return err
	}
	uc.Logger.Info(out)

	return nil
}

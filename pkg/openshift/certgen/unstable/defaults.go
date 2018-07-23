package unstable

import (
	"bytes"
	"fmt"
	"net"

	"github.com/Azure/acs-engine/pkg/api"
	"github.com/Azure/acs-engine/pkg/openshift/filesystem"
)

// OpenShiftSetDefaultCerts sets default certificate and configuration properties in the
// openshift orchestrator.
func OpenShiftSetDefaultCerts(a *api.Properties, orchestratorName, clusterID string) (bool, error) {
	if len(a.OrchestratorProfile.OpenShiftConfig.ConfigBundles["master"]) > 0 &&
		len(a.OrchestratorProfile.OpenShiftConfig.ConfigBundles["bootstrap"]) > 0 {
		return true, nil
	}

	c := Config{
		Master: &Master{
			Hostname: fmt.Sprintf("%s-master-%s-0", orchestratorName, clusterID),
			IPs: []net.IP{
				net.ParseIP(a.MasterProfile.FirstConsecutiveStaticIP),
			},
			Port: 8443,
		},
		ExternalMasterHostname:  fmt.Sprintf("%s.%s.cloudapp.azure.com", a.MasterProfile.DNSPrefix, a.AzProfile.Location),
		ClusterUsername:         a.OrchestratorProfile.OpenShiftConfig.ClusterUsername,
		ClusterPassword:         a.OrchestratorProfile.OpenShiftConfig.ClusterPassword,
		EnableAADAuthentication: a.OrchestratorProfile.OpenShiftConfig.EnableAADAuthentication,
		AzureConfig: AzureConfig{
			TenantID:                   a.AzProfile.TenantID,
			SubscriptionID:             a.AzProfile.SubscriptionID,
			AADClientID:                a.ServicePrincipalProfile.ClientID,
			AADClientSecret:            a.ServicePrincipalProfile.Secret,
			ResourceGroup:              a.AzProfile.ResourceGroup,
			Location:                   a.AzProfile.Location,
			SecurityGroupName:          fmt.Sprintf("%s-master-%s-nsg", orchestratorName, clusterID),
			PrimaryAvailabilitySetName: fmt.Sprintf("compute-availabilityset-%s", clusterID),
		},
	}

	err := c.PrepareMasterCerts()
	if err != nil {
		return false, err
	}
	err = c.PrepareMasterKubeConfigs()
	if err != nil {
		return false, err
	}
	err = c.PrepareMasterFiles()
	if err != nil {
		return false, err
	}

	err = c.PrepareBootstrapKubeConfig()
	if err != nil {
		return false, err
	}

	if a.OrchestratorProfile.OpenShiftConfig.ConfigBundles == nil {
		a.OrchestratorProfile.OpenShiftConfig.ConfigBundles = make(map[string][]byte)
	}

	masterBundle, err := getConfigBundle(c.WriteMaster)
	if err != nil {
		return false, err
	}
	a.OrchestratorProfile.OpenShiftConfig.ConfigBundles["master"] = masterBundle

	nodeBundle, err := getConfigBundle(c.WriteNode)
	if err != nil {
		return false, err
	}
	a.OrchestratorProfile.OpenShiftConfig.ConfigBundles["bootstrap"] = nodeBundle

	return true, nil
}

type writeFn func(filesystem.Writer) error

func getConfigBundle(write writeFn) ([]byte, error) {
	b := &bytes.Buffer{}

	fs, err := filesystem.NewTGZWriter(b)
	if err != nil {
		return nil, err
	}

	err = write(fs)
	if err != nil {
		return nil, err
	}

	err = fs.Close()
	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

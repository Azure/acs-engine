package runner

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/Azure/acs-engine/test/e2e/azure"
	"github.com/Azure/acs-engine/test/e2e/config"
	"github.com/Azure/acs-engine/test/e2e/dcos"
	"github.com/Azure/acs-engine/test/e2e/engine"
	"github.com/Azure/acs-engine/test/e2e/kubernetes/node"
	"github.com/Azure/acs-engine/test/e2e/metrics"
	"github.com/Azure/acs-engine/test/e2e/remote"
	"github.com/kelseyhightower/envconfig"
)

// CLIProvisioner holds the configuration needed to provision a clusters
type CLIProvisioner struct {
	ClusterDefinition string `envconfig:"CLUSTER_DEFINITION" required:"true" default:"examples/kubernetes.json"` // ClusterDefinition is the path on disk to the json template these are normally located in examples/
	ProvisionRetries  int    `envcofnig:"PROVISION_RETRIES" default:"3"`
	CreateVNET        bool   `envconfig:"CREATE_VNET" default:"false"`
	Config            *config.Config
	Account           *azure.Account
	Point             *metrics.Point
	ResourceGroups    []string
	Engine            *engine.Engine
}

// BuildCLIProvisioner will return a ProvisionerConfig object which is used to run a provision
func BuildCLIProvisioner(cfg *config.Config, acct *azure.Account, pt *metrics.Point) (*CLIProvisioner, error) {
	p := new(CLIProvisioner)
	if err := envconfig.Process("provisioner", p); err != nil {
		return nil, err
	}
	p.Config = cfg
	p.Account = acct
	p.Point = pt
	return p, nil
}

// Run will provision a cluster using the azure cli
func (cli *CLIProvisioner) Run() error {
	rgs := make([]string, 0)
	for i := 1; i <= cli.ProvisionRetries; i++ {
		cli.Point.SetProvisionStart()
		err := cli.provision()
		rgs = append(rgs, cli.Config.Name)
		cli.ResourceGroups = rgs
		if err != nil {
			if i < cli.ProvisionRetries {
				cli.Point.RecordProvisionError()
			} else if i == cli.ProvisionRetries {
				cli.Point.RecordProvisionError()
				return fmt.Errorf("Exceeded provision retry count")
			}
		} else {
			cli.Point.RecordProvisionSuccess()
			cli.Point.SetNodeWaitStart()
			err := cli.waitForNodes()
			cli.Point.RecordNodeWait(err)
			if err != nil {

				return err
			}
			return nil
		}
	}
	return fmt.Errorf("Unable to run provisioner")
}

func (cli *CLIProvisioner) provision() error {
	cli.Config.Name = cli.generateName()
	os.Setenv("NAME", cli.Config.Name)
	log.Printf("Cluster name:%s\n", cli.Config.Name)

	outputPath := filepath.Join(cli.Config.CurrentWorkingDir, "_output")
	os.Mkdir(outputPath, 0755)

	out, err := exec.Command("ssh-keygen", "-f", cli.Config.GetSSHKeyPath(), "-q", "-N", "", "-b", "2048", "-t", "rsa").CombinedOutput()
	if err != nil {
		return fmt.Errorf("Error while trying to generate ssh key:%s\nOutput:%s", err, out)
	}
	exec.Command("chmod", "0600", cli.Config.GetSSHKeyPath()+"*")

	publicSSHKey, err := cli.Config.ReadPublicSSHKey()
	if err != nil {
		return fmt.Errorf("Error while trying to read public ssh key: %s", err)
	}
	os.Setenv("PUBLIC_SSH_KEY", publicSSHKey)
	os.Setenv("DNS_PREFIX", cli.Config.Name)

	err = cli.Account.CreateGroup(cli.Config.Name, cli.Config.Location)
	if err != nil {
		return fmt.Errorf("Error while trying to create resource group: %s", err)
	}

	subnetID := ""
	vnetName := fmt.Sprintf("%sCustomVnet", cli.Config.Name)
	subnetName := fmt.Sprintf("%sCustomSubnet", cli.Config.Name)
	if cli.CreateVNET {
		err = cli.Account.CreateVnet(vnetName, "10.239.0.0/16", subnetName, "10.239.0.0/16")
		if err != nil {
			return fmt.Errorf("Error trying to create vnet:%s", err)
		}
		subnetID = fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/virtualNetworks/%s/subnets/%s", cli.Account.SubscriptionID, cli.Account.ResourceGroup.Name, vnetName, subnetName)
	}

	// Lets modify our template and call acs-engine generate on it
	eng, err := engine.Build(cli.Config, subnetID)
	if err != nil {
		return fmt.Errorf("Error while trying to build cluster definition:%s", err)
	}
	cli.Engine = eng

	err = cli.Engine.Write()
	if err != nil {
		return fmt.Errorf("Error while trying to write Engine Template to disk:%s", err)
	}

	err = cli.Engine.Generate()
	if err != nil {
		return fmt.Errorf("Error while trying to generate acs-engine template:%s", err)
	}

	// Lets start by just using the normal az group deployment cli for creating a cluster
	err = cli.Account.CreateDeployment(cli.Config.Name, eng)
	if err != nil {
		return fmt.Errorf("Error while trying to create deployment:%s", err)
	}

	return nil
}

// GenerateName will generate a new name if one has not been set
func (cli *CLIProvisioner) generateName() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	suffix := r.Intn(99999)
	prefix := fmt.Sprintf("%s-%s", cli.Config.Orchestrator, cli.Config.Location)
	return fmt.Sprintf("%s-%v", prefix, suffix)
}

func (cli *CLIProvisioner) waitForNodes() error {
	if cli.Config.IsKubernetes() {
		cli.Config.SetKubeConfig()
		log.Println("Waiting on nodes to go into ready state...")
		ready := node.WaitOnReady(cli.Engine.NodeCount(), 10*time.Second, cli.Config.Timeout)
		if ready == false {
			return errors.New("Error: Not all nodes in a healthy state")
		}
	}

	if cli.Config.IsDCOS() {
		host := fmt.Sprintf("%s.%s.cloudapp.azure.com", cli.Config.Name, cli.Config.Location)
		user := cli.Engine.ClusterDefinition.Properties.LinuxProfile.AdminUsername
		log.Printf("SSH Key: %s\n", cli.Config.GetSSHKeyPath())
		log.Printf("Master Node: %s@%s\n", user, host)
		log.Printf("SSH Command: ssh -i %s -p 2200 %s@%s", cli.Config.GetSSHKeyPath(), user, host)
		cluster, err := dcos.NewCluster(cli.Config, cli.Engine)
		if err != nil {
			return err
		}
		err = cluster.InstallDCOSClient()
		if err != nil {
			return fmt.Errorf("Error trying to install dcos client:%s", err)
		}
		ready := cluster.WaitForNodes(cli.Engine.NodeCount(), 10*time.Second, cli.Config.Timeout)
		if ready == false {
			return errors.New("Error: Not all nodes in a healthy state")
		}
	}
	return nil
}

func (cli *CLIProvisioner) fetchProvisioningMetrics(path string) error {
	conn, err := remote.NewConnection(fmt.Sprintf("%s.%s.cloudapp.azure.com", cli.Config.Name, cli.Config.Location), "2200", cli.Engine.ClusterDefinition.Properties.LinuxProfile.AdminUsername, cli.Config.GetSSHKeyPath())
	if err != nil {
		return err
	}
	data, err := conn.Read(path)
	if err != nil {
		return fmt.Errorf("Error reading file from path (%s):%s", path, err)
	}
	log.Printf("Data:%s\n", data)
	return nil
}

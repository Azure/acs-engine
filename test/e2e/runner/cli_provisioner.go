package runner

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/Azure/acs-engine/test/e2e/azure"
	"github.com/Azure/acs-engine/test/e2e/config"
	"github.com/Azure/acs-engine/test/e2e/dcos"
	"github.com/Azure/acs-engine/test/e2e/engine"
	"github.com/Azure/acs-engine/test/e2e/kubernetes/node"
	"github.com/Azure/acs-engine/test/e2e/kubernetes/util"
	"github.com/Azure/acs-engine/test/e2e/metrics"
	"github.com/Azure/acs-engine/test/e2e/remote"
	"github.com/kelseyhightower/envconfig"
)

// CLIProvisioner holds the configuration needed to provision a clusters
type CLIProvisioner struct {
	ClusterDefinition string `envconfig:"CLUSTER_DEFINITION" required:"true" default:"examples/kubernetes.json"` // ClusterDefinition is the path on disk to the json template these are normally located in examples/
	ProvisionRetries  int    `envconfig:"PROVISION_RETRIES" default:"0"`
	CreateVNET        bool   `envconfig:"CREATE_VNET" default:"false"`
	Config            *config.Config
	Account           *azure.Account
	Point             *metrics.Point
	ResourceGroups    []string
	Engine            *engine.Engine
	Masters           []azure.VM
	Agents            []azure.VM
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
	for i := 0; i <= cli.ProvisionRetries; i++ {
		cli.Point.SetProvisionStart()
		err := cli.provision()
		rgs = append(rgs, cli.Config.Name)
		cli.ResourceGroups = rgs
		if err != nil {
			if i < cli.ProvisionRetries {
				cli.Point.RecordProvisionError()
			} else if i == cli.ProvisionRetries {
				cli.Point.RecordProvisionError()
				return fmt.Errorf("Exceeded provision retry count: %s", err)
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
	if cli.Config.SoakClusterName != "" {
		cli.Config.Name = cli.Config.SoakClusterName
	}
	os.Setenv("NAME", cli.Config.Name)

	outputPath := filepath.Join(cli.Config.CurrentWorkingDir, "_output")
	os.Mkdir(outputPath, 0755)

	if cli.Config.SoakClusterName == "" {
		cmd := exec.Command("ssh-keygen", "-f", cli.Config.GetSSHKeyPath(), "-q", "-N", "", "-b", "2048", "-t", "rsa")
		util.PrintCommand(cmd)
		out, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("Error while trying to generate ssh key:%s\nOutput:%s", err, out)
		}
	}

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

	c, err := config.ParseConfig()
	if err != nil {
		return fmt.Errorf("unable to parse base config")
	}
	engCfg, err := engine.ParseConfig(cli.Config.CurrentWorkingDir, c.ClusterDefinition, c.Name)
	if err != nil {
		return fmt.Errorf("unable to parse config")
	}
	csGenerated, err := engine.ParseOutput(engCfg.GeneratedDefinitionPath + "/apimodel.json")
	if err != nil {
		return fmt.Errorf("unable to parse output")
	}
	cli.Engine.ExpandedDefinition = csGenerated

	// Lets start by just using the normal az group deployment cli for creating a cluster
	err = cli.Account.CreateDeployment(cli.Config.Name, eng)
	if err != nil {
		return fmt.Errorf("Error while trying to create deployment:%s", err)
	}

	// Store the hosts for future introspection
	hosts, err := cli.Account.GetHosts(cli.Config.Name)
	if err != nil {
		return err
	}
	var masters, agents []azure.VM
	for _, host := range hosts {
		if strings.Contains(host.Name, "master") {
			masters = append(masters, host)
		} else if strings.Contains(host.Name, "agent") {
			agents = append(agents, host)
		}
	}
	cli.Masters = masters
	cli.Agents = agents

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
		version, err := node.Version()
		if err != nil {
			log.Printf("Ready nodes did not return a version: %s", err)
		}
		log.Printf("Testing a Kubernetes %s cluster...\n", version)
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

// FetchProvisioningMetrics gets provisioning files from all hosts in a cluster
func (cli *CLIProvisioner) FetchProvisioningMetrics(path string, cfg *config.Config, acct *azure.Account) error {
	agentFiles := []string{"/var/log/azure/cluster-provision.log", "/var/log/cloud-init.log",
		"/var/log/cloud-init-output.log", "/var/log/syslog", "/var/log/azure/custom-script/handler.log",
		"/opt/m", "/opt/azure/containers/kubelet.sh", "/opt/azure/containers/provision.sh",
		"/opt/azure/provision-ps.log", "/var/log/azure/dnsdump.pcap"}
	masterFiles := agentFiles
	masterFiles = append(masterFiles, "/opt/azure/containers/mountetcd.sh", "/opt/azure/containers/setup-etcd.sh", "/opt/azure/containers/setup-etcd.log")
	hostname := fmt.Sprintf("%s.%s.cloudapp.azure.com", cli.Config.Name, cli.Config.Location)
	conn, err := remote.NewConnection(hostname, "22", cli.Engine.ClusterDefinition.Properties.LinuxProfile.AdminUsername, cli.Config.GetSSHKeyPath())
	if err != nil {
		return err
	}
	for _, master := range cli.Masters {
		for _, fp := range masterFiles {
			err := conn.CopyRemote(master.Name, fp)
			if err != nil {
				log.Printf("Error reading file from path (%s):%s", path, err)
			}
		}
	}

	for _, agent := range cli.Agents {
		for _, fp := range agentFiles {
			err := conn.CopyRemote(agent.Name, fp)
			if err != nil {
				log.Printf("Error reading file from path (%s):%s", path, err)
			}
		}
	}
	connectString := fmt.Sprintf("%s@%s:/tmp/k8s-*", conn.User, hostname)
	logsPath := filepath.Join(cfg.CurrentWorkingDir, "_logs", hostname)
	cmd := exec.Command("scp", "-i", conn.PrivateKeyPath, "-o", "ConnectTimeout=30", "-o", "StrictHostKeyChecking=no", connectString, logsPath)
	util.PrintCommand(cmd)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error output:%s\n", out)
		return err
	}

	return nil
}

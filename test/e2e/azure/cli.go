package azure

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"time"

	"github.com/Azure/acs-engine/test/e2e/engine"
	"github.com/Azure/acs-engine/test/e2e/kubernetes/util"

	"github.com/kelseyhightower/envconfig"
)

// Account holds the values needed to talk to the Azure API
type Account struct {
	User           *User  `json:"user"`
	TenantID       string `json:"tenantId" envconfig:"TENANT_ID" required:"true"`
	SubscriptionID string `json:"id" envconfig:"SUBSCRIPTION_ID" required:"true"`
	ResourceGroup  ResourceGroup
	Deployment     Deployment
}

// ResourceGroup represents a collection of azure resources
type ResourceGroup struct {
	Name     string
	Location string
}

// VM represents an azure vm
type VM struct {
	Name string `json:"name"`
}

// Deployment represents a deployment of an acs cluster
type Deployment struct {
	Name              string // Name of the deployment
	TemplateDirectory string // engine.GeneratedDefinitionPath
}

// User represents the user currently logged into an Account
type User struct {
	ID     string `json:"name" envconfig:"CLIENT_ID" required:"true"`
	Secret string `envconfig:"CLIENT_SECRET" required:"true"`
	Type   string `json:"type"`
}

// NewAccount will parse env vars and return a new struct
func NewAccount() (*Account, error) {
	a := new(Account)
	if err := envconfig.Process("account", a); err != nil {
		return nil, err
	}
	u := new(User)
	if err := envconfig.Process("user", u); err != nil {
		return nil, err
	}
	a.User = u
	return a, nil
}

// Login will login to a given subscription
func (a *Account) Login() error {
	output, err := exec.Command("az", "login",
		"--service-principal",
		"--username", a.User.ID,
		"--password", a.User.Secret,
		"--tenant", a.TenantID).CombinedOutput()
	if err != nil {
		log.Printf("output:%s\n", output)
		return err
	}
	return nil
}

// SetSubscription will call az account set --subscription for the given Account
func (a *Account) SetSubscription() error {
	cmd := exec.Command("az", "account", "set", "--subscription", a.SubscriptionID)
	util.PrintCommand(cmd)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error while trying to set subscription (%s):%s", a.SubscriptionID, err)
		log.Printf("Output:%s\n", out)
		return err
	}
	return nil
}

// CreateGroup will create a resource group in a given location
//--tags "type=${RESOURCE_GROUP_TAG_TYPE:-}" "now=$(date +%s)" "job=${JOB_BASE_NAME:-}" "buildno=${BUILD_NUM:-}"
func (a *Account) CreateGroup(name, location string) error {
	now := fmt.Sprintf("now=%v", time.Now().Unix())
	cmd := exec.Command("az", "group", "create", "--name", name, "--location", location, "--tags", now)
	util.PrintCommand(cmd)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error while trying to create resource group (%s) in %s:%s", name, location, err)
		log.Printf("Output:%s\n", out)
		return err
	}
	r := ResourceGroup{
		Name:     name,
		Location: location,
	}
	a.ResourceGroup = r
	return nil
}

// DeleteGroup deletes a given resource group by name
func (a *Account) DeleteGroup(name string, wait bool) error {
	var cmd *exec.Cmd
	if !wait {
		cmd = exec.Command("az", "group", "delete", "--name", name, "--no-wait", "--yes")
	} else {
		cmd = exec.Command("az", "group", "delete", "--name", name, "--yes")
	}
	util.PrintCommand(cmd)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error while trying to delete resource group (%s):%s", name, out)
		return err
	}
	return nil
}

// CreateDeployment will deploy a cluster to a given resource group using the template and parameters on disk
func (a *Account) CreateDeployment(name string, e *engine.Engine) error {
	d := Deployment{
		Name:              name,
		TemplateDirectory: e.Config.GeneratedDefinitionPath,
	}

	ticker := time.NewTicker(1 * time.Minute)
	quit := make(chan bool)
	go func() {
		for {
			select {
			case <-ticker.C:
				fmt.Print(".")
			case <-quit:
				fmt.Print("\n")
				ticker.Stop()
				return
			}
		}
	}()

	cmd := exec.Command("az", "group", "deployment", "create",
		"--name", d.Name,
		"--resource-group", a.ResourceGroup.Name,
		"--template-file", e.Config.GeneratedTemplatePath,
		"--parameters", e.Config.GeneratedParametersPath)
	util.PrintCommand(cmd)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("\nError from deployment for %s in resource group %s:%s\n", d.Name, a.ResourceGroup.Name, err)
		log.Printf("Command Output: %s\n", out)
		return err
	}
	quit <- true
	a.Deployment = d
	return nil
}

// GetCurrentAccount will run an az account show and parse that into an account strcut
func GetCurrentAccount() (*Account, error) {
	cmd := exec.Command("az", "account", "show")
	util.PrintCommand(cmd)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error trying to run 'account show':%s\n", err)
		return nil, err
	}
	a := Account{}
	err = json.Unmarshal(out, &a)
	if err != nil {
		log.Printf("Error unmarshalling account json:%s\n", err)
		log.Printf("JSON:%s\n", out)
	}
	return &a, nil
}

// CreateVnet will create a vnet in a resource group
func (a *Account) CreateVnet(vnet, addressPrefixes, subnetName, subnetPrefix string) error {
	cmd := exec.Command("az", "network", "vnet", "create", "-g", a.ResourceGroup.Name, "-n", vnet, "--address-prefixes", addressPrefixes, "--subnet-name", subnetName, "--subnet-prefix", subnetPrefix)
	util.PrintCommand(cmd)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error while trying to create vnet with the following command:\n az network vnet create -g %s -n %s --address-prefixes %s --subnet-name %s --subnet-prefix %s \n Output:%s\n", a.ResourceGroup.Name, vnet, addressPrefixes, subnetName, subnetPrefix, out)
		return err
	}
	return nil
}

// RouteTable holds information from running az network route-table list
type RouteTable struct {
	ID                string `json:"id"`
	Location          string `json:"location"`
	Name              string `json:"name"`
	ProvisioningState string `json:"provisioningState"`
	ResourceGroup     string `json:"resourceGroup"`
}

// UpdateRouteTables is used to updated a vnet with the appropriate route tables
func (a *Account) UpdateRouteTables(subnet, vnet string) error {
	cmd := exec.Command("az", "network", "route-table", "list", "-g", a.ResourceGroup.Name)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error while trying to get route table list!\n Output:%s\n", out)
		return err
	}
	rts := []RouteTable{}
	json.Unmarshal(out, &rts)

	cmd = exec.Command("az", "network", "vnet", "subnet", "update", "-n", subnet, "-g", a.ResourceGroup.Name, "--vnet-name", vnet, "--route-table", rts[0].Name)
	util.PrintCommand(cmd)
	out, err = cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error while trying to update vnet route tables:%s\n", out)
		return err
	}
	return nil
}

// GetHosts will get a list of vms in the resource group
func (a *Account) GetHosts(name string) ([]VM, error) {
	var resourceGroup string
	if name != "" {
		resourceGroup = name
	} else {
		resourceGroup = a.ResourceGroup.Name
	}
	cmd := exec.Command("az", "vm", "list", "-g", resourceGroup)
	util.PrintCommand(cmd)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error while trying to get vm list:%s\n", out)
		return nil, err
	}
	v := []VM{{}}
	err = json.Unmarshal(out, &v)
	if err != nil {
		log.Printf("Error unmarshalling account json:%s\n", err)
		log.Printf("JSON:%s\n", out)
	}
	return v, nil
}

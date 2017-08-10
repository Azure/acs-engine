package azure

import (
	"encoding/json"
	"log"
	"os/exec"

	"github.com/Azure/acs-engine/test/e2e/engine"

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
	cmd := exec.Command("az", "login",
		"--service-principal",
		"--username", a.User.ID,
		"--password", a.User.Secret,
		"--tenant", a.TenantID)
	err := cmd.Start()
	if err != nil {
		log.Printf("Error while trying to start login:%s\n", err)
		return err
	}
	err = cmd.Wait()
	if err != nil {
		log.Printf("Error occurred while waiting for login to complete:%s\n", err)
		return err
	}
	return nil
}

// SetSubscription will call az account set --subscription for the given Account
func (a *Account) SetSubscription() error {
	cmd := exec.Command("az", "account", "set", "--subscription", a.SubscriptionID)
	err := cmd.Start()
	if err != nil {
		log.Printf("Error while trying to start account set for subscription %s:%s\n", a.SubscriptionID, err)
		return err
	}
	err = cmd.Wait()
	if err != nil {
		log.Printf("Error occurred while waiting for account set for subscription %s to complete:%s\n", a.SubscriptionID, err)
		return err
	}
	return nil
}

// CreateGroup will create a resource group in a given location
func (a *Account) CreateGroup(name, location string) error {
	cmd := exec.Command("az", "group", "create", "--name", name, "--location", location)
	err := cmd.Start()
	if err != nil {
		log.Printf("Error while trying to start command to create resource group (%s) in %s:%s", name, location, err)
		return err
	}
	err = cmd.Wait()
	if err != nil {
		log.Printf("Error occurred while waiting for resource group (%s) in %s:%s", name, location, err)
		return err
	}
	r := ResourceGroup{
		Name:     name,
		Location: location,
	}
	a.ResourceGroup = r
	return nil
}

// DeleteGroup delets a given resource group by name
func (a *Account) DeleteGroup() error {
	cmd := exec.Command("az", "group", "delete", "--name", a.Deployment.Name, "--no-wait", "--yes")
	err := cmd.Start()

	if err != nil {
		log.Printf("Error while trying to start command to delete resource group (%s):%s", a.Deployment.Name, err)
		return err
	}
	err = cmd.Wait()
	if err != nil {
		log.Printf("Error occurred while waiting for resource group (%s):%s", a.Deployment.Name, err)
		return err
	}
	return nil
}

// CreateDeployment will deploy a cluster to a given resource group using the template and parameters on disk
func (a *Account) CreateDeployment(name string, e *engine.Engine) error {
	d := Deployment{
		Name:              name,
		TemplateDirectory: e.GeneratedDefinitionPath,
	}
	cmd := exec.Command("az", "group", "deployment", "create",
		"--name", d.Name,
		"--resource-group", a.ResourceGroup.Name,
		"--template-file", e.GeneratedTemplatePath,
		"--parameters", e.GeneratedParametersPath)

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error while trying to start deployment for %s in resource group %s:%s", d.Name, a.ResourceGroup.Name, err)
		log.Printf("Command Output: %s\n", output)
		return err
	}

	a.Deployment = d
	return nil
}

// GetCurrentAccount will run an az account show and parse that into an account strcut
func GetCurrentAccount() (*Account, error) {
	out, err := exec.Command("az", "account", "show").CombinedOutput()
	if err != nil {
		log.Printf("Error trying to run 'account show':%s\n", err)
		return nil, err
	}
	a := Account{}
	err = json.Unmarshal(out, &a)
	if err != nil {
		log.Printf("Error unmarshalling account json:%s\n", err)
	}
	return &a, nil
}

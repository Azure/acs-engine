package e2e

import (
	"encoding/json"
	"log"
	"os/exec"

	"github.com/kelseyhightower/envconfig"
)

// Account holds the values needed to talk to the Azure API
type Account struct {
	User           *User  `json:"user"`
	TenantID       string `json:"tenantId" envconfig:"TENANT_ID" required:"true"`
	SubscriptionID string `json:"id" envconfig:"SUBSCRIPTION_ID" required:"true"`
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
	return nil
}

// CreateDeployment will deploy a cluster to a given resource group using the template and parameters on disk
func (a *Account) CreateDeployment(name, resourceGroup, templateFilePath, parametersFilePath string) error {
	cmd := exec.Command("az", "group", "deployment", "create",
		"--name", name,
		"--resource_group", resourceGroup,
		"--tempalte-file", templateFilePath,
		"--parameters", parametersFilePath)
	err := cmd.Start()
	if err != nil {
		log.Printf("Error while trying to start deployment for %s in resource group %s:%s", name, resourceGroup, err)
		return err
	}
	err = cmd.Wait()
	if err != nil {
		log.Printf("Error occured while waiting for deployment %s in resource group %s:%s", name, resourceGroup, err)
		return err
	}
	return nil
}

// GetCurrentAccount will run an az account show and parse that into an account strcut
func GetCurrentAccount() (*Account, error) {
	out, err := exec.Command("az", "account", "show").Output()
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

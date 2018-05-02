package cmd

import (
	"fmt"
	"strconv"
	"testing"

	"os"

	"github.com/Azure/acs-engine/pkg/api"
	"github.com/Azure/acs-engine/pkg/armhelpers"
	"github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

const ExampleAPIModel = `{
  "apiVersion": "vlabs",
  "properties": {
		"orchestratorProfile": { "orchestratorType": "Kubernetes", "kubernetesConfig": { "useManagedIdentity": %s, "etcdVersion" : "2.3.8" } },
    "masterProfile": { "count": 1, "dnsPrefix": "", "vmSize": "Standard_D2_v2" },
    "agentPoolProfiles": [ { "name": "linuxpool1", "count": 2, "vmSize": "Standard_D2_v2", "availabilityProfile": "AvailabilitySet" } ],
    "windowsProfile": { "adminUsername": "azureuser", "adminPassword": "replacepassword1234$" },
    "linuxProfile": { "adminUsername": "azureuser", "ssh": { "publicKeys": [ { "keyData": "" } ] }
    },
    "servicePrincipalProfile": { "clientId": "%s", "secret": "%s" }
  }
}
`

const ExampleAPIModelWithDNSPrefix = `{
	"apiVersion": "vlabs",
	"properties": {
		  "orchestratorProfile": { "orchestratorType": "Kubernetes", "kubernetesConfig": { "useManagedIdentity": %s, "etcdVersion" : "2.3.8" } },
	  "masterProfile": { "count": 1, "dnsPrefix": "mytestcluster", "vmSize": "Standard_D2_v2" },
	  "agentPoolProfiles": [ { "name": "linuxpool1", "count": 2, "vmSize": "Standard_D2_v2", "availabilityProfile": "AvailabilitySet" } ],
	  "windowsProfile": { "adminUsername": "azureuser", "adminPassword": "replacepassword1234$" },
	  "linuxProfile": { "adminUsername": "azureuser", "ssh": { "publicKeys": [ { "keyData": "" } ] }
	  },
	  "servicePrincipalProfile": { "clientId": "%s", "secret": "%s" }
	}
  }
  `

const ExampleAPIModelWithoutServicePrincipalProfile = `{
	"apiVersion": "vlabs",
	"properties": {
		  "orchestratorProfile": { "orchestratorType": "Kubernetes", "kubernetesConfig": { "useManagedIdentity": %s, "etcdVersion" : "2.3.8" } },
	  "masterProfile": { "count": 1, "dnsPrefix": "mytestcluster", "vmSize": "Standard_D2_v2" },
	  "agentPoolProfiles": [ { "name": "linuxpool1", "count": 2, "vmSize": "Standard_D2_v2", "availabilityProfile": "AvailabilitySet" } ],
	  "windowsProfile": { "adminUsername": "azureuser", "adminPassword": "replacepassword1234$" },
	  "linuxProfile": { "adminUsername": "azureuser", "ssh": { "publicKeys": [ { "keyData": "" } ] }
	  }
	}
  }
  `

func getExampleAPIModel(useManagedIdentity bool, clientID, clientSecret string) string {
	return getAPIModel(ExampleAPIModel, useManagedIdentity, clientID, clientSecret)
}

func getAPIModel(baseAPIModel string, useManagedIdentity bool, clientID, clientSecret string) string {
	return fmt.Sprintf(
		baseAPIModel,
		strconv.FormatBool(useManagedIdentity),
		clientID,
		clientSecret)
}

func getAPIModelWithoutServicePrincipalProfile(baseAPIModel string, useManagedIdentity bool) string {
	return fmt.Sprintf(
		baseAPIModel,
		strconv.FormatBool(useManagedIdentity))
}

func TestAutofillApimodelWithoutManagedIdentityCreatesCreds(t *testing.T) {
	testAutodeployCredentialHandling(t, false, "", "")
}

func TestAutofillApimodelWithManagedIdentitySkipsCreds(t *testing.T) {
	testAutodeployCredentialHandling(t, true, "", "")
}

func TestAutofillApimodelAllowsPrespecifiedCreds(t *testing.T) {
	testAutodeployCredentialHandling(t, false, "clientID", "clientSecret")
}

func TestAutoSufixWithDnsPrefixInApiModel(t *testing.T) {
	apiloader := &api.Apiloader{
		Translator: nil,
	}

	apimodel := getAPIModel(ExampleAPIModelWithDNSPrefix, false, "clientID", "clientSecret")
	cs, ver, err := apiloader.DeserializeContainerService([]byte(apimodel), false, false, nil)
	if err != nil {
		t.Fatalf("unexpected error deserializing the example apimodel: %s", err)
	}
	deployCmd := &deployCmd{
		apimodelPath:     "./this/is/unused.json",
		outputDirectory:  "_test_output",
		forceOverwrite:   true,
		location:         "westus",
		autoSuffix:       true,
		containerService: cs,
		apiVersion:       ver,

		client: &armhelpers.MockACSEngineClient{},
	}
	autofillApimodel(deployCmd)

	defer os.RemoveAll(deployCmd.outputDirectory)

	if deployCmd.containerService.Properties.MasterProfile.DNSPrefix == "mytestcluster" {
		t.Fatalf("expected %s-{timestampsuffix} but got %s", "mytestcluster", deployCmd.containerService.Properties.MasterProfile.DNSPrefix)
	}

}

func TestAPIModelWithoutServicePrincipalProfileAndClientIdAndSecretInCmd(t *testing.T) {
	apiloader := &api.Apiloader{
		Translator: nil,
	}

	apimodel := getAPIModelWithoutServicePrincipalProfile(ExampleAPIModelWithoutServicePrincipalProfile, false)
	TestClientIDInCmd, err := uuid.FromString("DEC923E3-1EF1-4745-9516-37906D56DEC4")
	if err != nil {
		t.Fatalf("Invalid ClientID in Test: %s", err)
	}

	TestClientSecretInCmd := "DEC923E3-1EF1-4745-9516-37906D56DEC4"

	cs, ver, err := apiloader.DeserializeContainerService([]byte(apimodel), false, false, nil)
	if err != nil {
		t.Fatalf("unexpected error deserializing the example apimodel: %s", err)
	}
	deployCmd := &deployCmd{
		apimodelPath:     "./this/is/unused.json",
		outputDirectory:  "_test_output",
		forceOverwrite:   true,
		location:         "westus",
		containerService: cs,
		apiVersion:       ver,

		client: &armhelpers.MockACSEngineClient{},
	}
	deployCmd.ClientID = TestClientIDInCmd
	deployCmd.ClientSecret = TestClientSecretInCmd
	autofillApimodel(deployCmd)

	defer os.RemoveAll(deployCmd.outputDirectory)

	if deployCmd.containerService.Properties.ServicePrincipalProfile == nil || deployCmd.containerService.Properties.ServicePrincipalProfile.ClientID == "" || deployCmd.containerService.Properties.ServicePrincipalProfile.Secret == "" {
		t.Fatalf("expected service principal profile to be populated from deployment command arguments")
	}

	if deployCmd.containerService.Properties.ServicePrincipalProfile.ClientID != TestClientIDInCmd.String() {
		t.Fatalf("expected service principal profile client id to be %s but got %s", TestClientIDInCmd.String(), deployCmd.containerService.Properties.ServicePrincipalProfile.ClientID)
	}

	if deployCmd.containerService.Properties.ServicePrincipalProfile.Secret != TestClientSecretInCmd {
		t.Fatalf("expected service principal profile client secret to be %s but got %s", TestClientSecretInCmd, deployCmd.containerService.Properties.ServicePrincipalProfile.Secret)
	}
}

func TestAPIModelWithEmptyServicePrincipalProfileAndClientIdAndSecretInCmd(t *testing.T) {
	apiloader := &api.Apiloader{
		Translator: nil,
	}

	apimodel := getAPIModel(ExampleAPIModelWithDNSPrefix, false, "", "")
	TestClientIDInCmd, err := uuid.FromString("DEC923E3-1EF1-4745-9516-37906D56DEC4")
	if err != nil {
		t.Fatalf("Invalid ClientID in Test: %s", err)
	}

	TestClientSecretInCmd := "DEC923E3-1EF1-4745-9516-37906D56DEC4"

	cs, ver, err := apiloader.DeserializeContainerService([]byte(apimodel), false, false, nil)
	if err != nil {
		t.Fatalf("unexpected error deserializing the example apimodel: %s", err)
	}
	deployCmd := &deployCmd{
		apimodelPath:     "./this/is/unused.json",
		outputDirectory:  "_test_output",
		forceOverwrite:   true,
		location:         "westus",
		containerService: cs,
		apiVersion:       ver,

		client: &armhelpers.MockACSEngineClient{},
	}
	deployCmd.ClientID = TestClientIDInCmd
	deployCmd.ClientSecret = TestClientSecretInCmd
	autofillApimodel(deployCmd)

	defer os.RemoveAll(deployCmd.outputDirectory)

	if deployCmd.containerService.Properties.ServicePrincipalProfile == nil || deployCmd.containerService.Properties.ServicePrincipalProfile.ClientID == "" || deployCmd.containerService.Properties.ServicePrincipalProfile.Secret == "" {
		t.Fatalf("expected service principal profile to be populated from deployment command arguments")
	}

	if deployCmd.containerService.Properties.ServicePrincipalProfile.ClientID != TestClientIDInCmd.String() {
		t.Fatalf("expected service principal profile client id to be %s but got %s", TestClientIDInCmd.String(), deployCmd.containerService.Properties.ServicePrincipalProfile.ClientID)
	}

	if deployCmd.containerService.Properties.ServicePrincipalProfile.Secret != TestClientSecretInCmd {
		t.Fatalf("expected service principal profile client secret to be %s but got %s", TestClientSecretInCmd, deployCmd.containerService.Properties.ServicePrincipalProfile.Secret)
	}
}

func TestAPIModelWithoutServicePrincipalProfileAndWithoutClientIdAndSecretInCmd(t *testing.T) {
	apiloader := &api.Apiloader{
		Translator: nil,
	}

	apimodel := getAPIModelWithoutServicePrincipalProfile(ExampleAPIModelWithoutServicePrincipalProfile, false)

	cs, ver, err := apiloader.DeserializeContainerService([]byte(apimodel), false, false, nil)
	if err != nil {
		t.Fatalf("unexpected error deserializing the example apimodel: %s", err)
	}
	deployCmd := &deployCmd{
		apimodelPath:     "./this/is/unused.json",
		outputDirectory:  "_test_output",
		forceOverwrite:   true,
		location:         "westus",
		containerService: cs,
		apiVersion:       ver,

		client: &armhelpers.MockACSEngineClient{},
	}
	autofillApimodel(deployCmd)

	defer os.RemoveAll(deployCmd.outputDirectory)

	if deployCmd.containerService.Properties.ServicePrincipalProfile != nil {
		t.Fatalf("expected service principal profile to be nil for unmanaged identity, where client id and secret are not supplied in api model and deployment command")
	}

}

func TestAPIModelWithEmptyServicePrincipalProfileAndWithoutClientIdAndSecretInCmd(t *testing.T) {
	apiloader := &api.Apiloader{
		Translator: nil,
	}

	apimodel := getAPIModel(ExampleAPIModelWithDNSPrefix, false, "", "")

	cs, ver, err := apiloader.DeserializeContainerService([]byte(apimodel), false, false, nil)
	if err != nil {
		t.Fatalf("unexpected error deserializing the example apimodel: %s", err)
	}
	deployCmd := &deployCmd{
		apimodelPath:     "./this/is/unused.json",
		outputDirectory:  "_test_output",
		forceOverwrite:   true,
		location:         "westus",
		containerService: cs,
		apiVersion:       ver,

		client: &armhelpers.MockACSEngineClient{},
	}
	autofillApimodel(deployCmd)

	defer os.RemoveAll(deployCmd.outputDirectory)

	if deployCmd.containerService.Properties.ServicePrincipalProfile == nil {
		t.Fatalf("expected service principal profile to be Empty and not nil for unmanaged identity, where client id and secret are not supplied in api model and deployment command")
	}

	// mockclient returns "app-id" for ClientID when empty
	if deployCmd.containerService.Properties.ServicePrincipalProfile.ClientID != "app-id" {
		t.Fatalf("expected service principal profile client id to be empty but got %s", deployCmd.containerService.Properties.ServicePrincipalProfile.ClientID)
	}

	// mockcliet returns "client-secret" when empty
	if deployCmd.containerService.Properties.ServicePrincipalProfile.Secret != "client-secret" {
		t.Fatalf("expected service principal profile client secret to be empty but got %s", deployCmd.containerService.Properties.ServicePrincipalProfile.Secret)
	}

}

func testAutodeployCredentialHandling(t *testing.T, useManagedIdentity bool, clientID, clientSecret string) {
	apiloader := &api.Apiloader{
		Translator: nil,
	}

	apimodel := getExampleAPIModel(useManagedIdentity, clientID, clientSecret)
	cs, ver, err := apiloader.DeserializeContainerService([]byte(apimodel), false, false, nil)
	if err != nil {
		t.Fatalf("unexpected error deserializing the example apimodel: %s", err)
	}

	// deserialization happens in validate(), but we are testing just the default
	// setting that occurs in autofillApimodel (which is called from validate)
	// Thus, it assumes that containerService/apiVersion are already populated
	deployCmd := &deployCmd{
		apimodelPath:    "./this/is/unused.json",
		dnsPrefix:       "dnsPrefix1",
		outputDirectory: "_test_output",
		forceOverwrite:  true,
		location:        "westus",

		containerService: cs,
		apiVersion:       ver,

		client: &armhelpers.MockACSEngineClient{},
	}

	autofillApimodel(deployCmd)

	// cleanup, since auto-populations creates dirs and saves the SSH private key that it might create
	defer os.RemoveAll(deployCmd.outputDirectory)

	cs, _, err = revalidateApimodel(apiloader, cs, ver)
	if err != nil {
		log.Fatalf("unexpected error validating apimodel after populating defaults: %s", err)
	}

	if useManagedIdentity {
		if cs.Properties.ServicePrincipalProfile != nil &&
			(cs.Properties.ServicePrincipalProfile.ClientID != "" || cs.Properties.ServicePrincipalProfile.Secret != "") {
			log.Fatalf("Unexpected credentials were populated even though MSI was active.")
		}
	} else {
		if cs.Properties.ServicePrincipalProfile == nil ||
			cs.Properties.ServicePrincipalProfile.ClientID == "" || cs.Properties.ServicePrincipalProfile.Secret == "" {
			log.Fatalf("Credentials were missing even though MSI was not active.")
		}
	}
}

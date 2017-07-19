package cmd

import (
	"fmt"
	"strconv"
	"testing"

	"os"

	"github.com/Azure/acs-engine/pkg/api"
	"github.com/Azure/acs-engine/pkg/armhelpers"
	log "github.com/Sirupsen/logrus"
)

const ExampleAPIModel = `{
  "apiVersion": "vlabs",
  "properties": {
    "orchestratorProfile": { "orchestratorType": "Kubernetes", "kubernetesConfig": { "useManagedIdentity": %s } },
    "masterProfile": { "count": 1, "dnsPrefix": "", "vmSize": "Standard_D2_v2" },
    "agentPoolProfiles": [ { "name": "linuxpool1", "count": 2, "vmSize": "Standard_D2_v2", "availabilityProfile": "AvailabilitySet" } ],
    "windowsProfile": { "adminUsername": "azureuser", "adminPassword": "replacepassword1234$" },
    "linuxProfile": { "adminUsername": "azureuser", "ssh": { "publicKeys": [ { "keyData": "" } ] }
    },
    "servicePrincipalProfile": { "servicePrincipalClientID": "", "servicePrincipalClientSecret": "" }
  }
}
`

func getExampleAPIModel(useManagedIdentity bool) string {
	return fmt.Sprintf(ExampleAPIModel, strconv.FormatBool(useManagedIdentity))
}

func TestAutofillApimodelWithoutManagedIdentityCreatesCreds(t *testing.T) {
	testMSIPopulatedInternal(t, false)
}

func TestAutofillApimodelWithManagedIdentitySkipsCreds(t *testing.T) {
	testMSIPopulatedInternal(t, true)
}

func testMSIPopulatedInternal(t *testing.T, useManagedIdentity bool) {
	apimodel := getExampleAPIModel(useManagedIdentity)
	cs, ver, err := api.DeserializeContainerService([]byte(apimodel), false)
	if err != nil {
		t.Fatalf("unexpected error deserializing the example apimodel: %s", err)
	}

	// deserialization happens in validate(), but we are testing just the default
	// setting that occurs in autofillApimodel (which is called from validate)
	// Thus, it assumes that containerService/apiVersion are already populated
	deployCmd := &deployCmd{
		apimodelPath:    "./this/is/unused.json",
		dnsPrefix:       "dnsPrefix1",
		outputDirectory: "dummy/path/",
		location:        "westus",

		containerService: cs,
		apiVersion:       ver,

		client: &armhelpers.MockACSEngineClient{},
	}

	autofillApimodel(deployCmd)

	cs, ver, err = revalidateApimodel(cs, ver)
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

	// cleanup, since auto-populations creates dirs and saves the SSH private key that it might create
	os.RemoveAll(deployCmd.outputDirectory)
}

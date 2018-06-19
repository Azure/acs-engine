package openshift

import (
	"encoding/json"
	"io/ioutil"
	"path"
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"

	"github.com/Azure/acs-engine/pkg/api"
	"github.com/Azure/acs-engine/pkg/api/agentPoolOnlyApi/vlabs"
	"github.com/Azure/acs-engine/pkg/i18n"
)

func TestExamplesInSync(t *testing.T) {
	baseExampleDir := "../../examples"
	baseExample := "openshift.json"

	tests := []string{}

	locale, err := i18n.LoadTranslations()
	if err != nil {
		t.Fatalf("error loading translations %v", err)
	}

	apiloader := &api.Apiloader{
		Translator: &i18n.Translator{
			Locale: locale,
		},
	}

	baseExampleCS, _, err := apiloader.LoadContainerServiceFromFile(
		path.Join(baseExampleDir, baseExample),
		false, //don't validate
		false, // not update
		nil,
	)
	if err != nil {
		t.Fatalf("error parsing the api model: %s", err.Error())
	}

	for _, test := range tests {
		testCS, _, err := apiloader.LoadContainerServiceFromFile(
			path.Join(baseExampleDir, test),
			false,
			false,
			nil,
		)
		if err != nil {
			t.Errorf("failed parsing %s: %#v", test, err)
			continue
		}

		// todo normalize where necessary (seems easier than reflect right now)

		if !reflect.DeepEqual(baseExampleCS.Properties, testCS.Properties) {
			t.Errorf(spew.Sprintf("Testing %s\nExpected:\n%+v\nGot:\n%+v", test, baseExampleCS.Properties, testCS.Properties))
		}
	}
}

func TestAgentsOnlyExample(t *testing.T) {
	example := "../../examples/openshift-agents-only.json"

	contents, err := ioutil.ReadFile(example)
	if err != nil {
		t.Fatalf("cannot read file %s: %v", example, err)
	}

	m := &vlabs.ManagedCluster{}
	if err := json.Unmarshal(contents, &m); err != nil {
		t.Fatalf("cannot unmarshal file %s: %v", example, err)
	}

	m.Properties.DNSPrefix = "mycluster"
	m.Properties.ServicePrincipalProfile = &vlabs.ServicePrincipalProfile{
		ClientID: "CLIENT_ID",
		Secret:   "TOP_SECRET",
	}
	m.Properties.LinuxProfile.SSH = struct {
		PublicKeys []vlabs.PublicKey `json:"publicKeys" validate:"required,len=1"`
	}{
		PublicKeys: []vlabs.PublicKey{
			{KeyData: "KEY_DATA"},
		},
	}
	m.Properties.CertificateProfile = &vlabs.CertificateProfile{
		CaCertificate: "CA_CERT",
		CaPrivateKey:  "CA_PRIV_KEY",
	}

	if err := m.Properties.Validate(); err != nil {
		t.Fatalf("cannot validate file %s: %v", example, err)
	}
}

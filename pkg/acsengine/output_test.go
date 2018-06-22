package acsengine

import (
	"os"
	"testing"

	"github.com/Azure/acs-engine/pkg/i18n"
)

func TestWriteTLSArtifacts(t *testing.T) {

	writer := &ArtifactWriter{
		Translator: &i18n.Translator{
			Locale: nil,
		},
	}
	dir := "_testoutputdir"
	defer os.Remove(dir + "/apimodel.json")
	defer os.Remove(dir + "/azuredeploy.json")
	defer os.Remove(dir + "/azuredeploy.parameters.json")
	cs := CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2)
	err := writer.WriteTLSArtifacts(cs, "vlabs", "fake template", "fake parameters", dir, false, false)

	if err != nil {
		t.Fatalf("unexpected error trying to write TLS artifacts: %s", err.Error())
	}
}

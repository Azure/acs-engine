package acsengine

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/Azure/acs-engine/pkg/i18n"
)

func TestWriteTLSArtifacts(t *testing.T) {

	cs := CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2)
	writer := &ArtifactWriter{
		Translator: &i18n.Translator{
			Locale: nil,
		},
	}
	dir := "_testoutputdir"
	defaultDir := fmt.Sprintf("%s-%s", cs.Properties.OrchestratorProfile.OrchestratorType, GenerateClusterID(cs.Properties))
	defaultDir = path.Join("_output", defaultDir)
	defer os.RemoveAll(dir)
	defer os.RemoveAll(defaultDir)

	err := writer.WriteTLSArtifacts(cs, "vlabs", "fake template", "fake parameters", dir, false, false)

	if err != nil {
		t.Fatalf("unexpected error trying to write TLS artifacts: %s", err.Error())
	}

	if _, err := os.Stat(dir + "/apimodel.json"); os.IsNotExist(err) {
		t.Fatalf("expected file %s/apimodel.json to be generated by WriteTLSArtifacts", dir)
	}

	if _, err := os.Stat(dir + "/azuredeploy.json"); os.IsNotExist(err) {
		t.Fatalf("expected file %s/azuredeploy.json to be generated by WriteTLSArtifacts", dir)
	}

	if _, err := os.Stat(dir + "/azuredeploy.parameters.json"); os.IsNotExist(err) {
		t.Fatalf("expected file %s/azuredeploy.parameters.json to be generated by WriteTLSArtifacts", dir)
	}

	os.RemoveAll(dir)

	err = writer.WriteTLSArtifacts(cs, "vlabs", "fake template", "fake parameters", "", true, true)

	if err != nil {
		t.Fatalf("unexpected error trying to write TLS artifacts: %s", err.Error())
	}

	if _, err := os.Stat(defaultDir + "/apimodel.json"); !os.IsNotExist(err) {
		t.Fatalf("expected file %s/apimodel.json not to be generated by WriteTLSArtifacts with parametersOnly set to true", defaultDir)
	}

	if _, err := os.Stat(defaultDir + "/azuredeploy.json"); !os.IsNotExist(err) {
		t.Fatalf("expected file %s/azuredeploy.json not to be generated by WriteTLSArtifacts with parametersOnly set to true", defaultDir)
	}

	if _, err := os.Stat(defaultDir + "/azuredeploy.parameters.json"); os.IsNotExist(err) {
		t.Fatalf("expected file %s/azuredeploy.parameters.json to be generated by WriteTLSArtifacts with parametersOnly set to true", defaultDir)
	}

}

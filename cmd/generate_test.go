package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestGenerateCmdValidate(t *testing.T) {
	g := &generateCmd{}
	r := &cobra.Command{}

	// validate cmd with 1 arg
	err := g.validate(r, []string{"../pkg/acsengine/testdata/simple/kubernetes.json"})
	if err != nil {
		t.Fatalf("unexpected error validating 1 arg: %s", err.Error())
	}

	g = &generateCmd{}

	// validate cmd with 0 args
	err = g.validate(r, []string{})
	t.Logf(err.Error())
	if err == nil {
		t.Fatalf("expected error validating 0 args")
	}

	g = &generateCmd{}

	// validate cmd with more than 1 arg
	err = g.validate(r, []string{"../pkg/acsengine/testdata/simple/kubernetes.json", "arg1"})
	t.Logf(err.Error())
	if err == nil {
		t.Fatalf("expected error validating multiple args")
	}

}

func TestGenerateCmdMergeAPIModel(t *testing.T) {
	g := &generateCmd{}
	g.apimodelPath = "../pkg/acsengine/testdata/simple/kubernetes.json"
	err := g.mergeAPIModel()
	if err != nil {
		t.Fatalf("unexpected error calling mergeAPIModel with no --set flag defined: %s", err.Error())
	}

	g = &generateCmd{}
	g.apimodelPath = "../pkg/acsengine/testdata/simple/kubernetes.json"
	g.set = []string{"masterProfile.count=3,linuxProfile.adminUsername=testuser"}
	err = g.mergeAPIModel()
	if err != nil {
		t.Fatalf("unexpected error calling mergeAPIModel with one --set flag: %s", err.Error())
	}

	g = &generateCmd{}
	g.apimodelPath = "../pkg/acsengine/testdata/simple/kubernetes.json"
	g.set = []string{"masterProfile.count=3", "linuxProfile.adminUsername=testuser"}
	err = g.mergeAPIModel()
	if err != nil {
		t.Fatalf("unexpected error calling mergeAPIModel with multiple --set flags: %s", err.Error())
	}

	g = &generateCmd{}
	g.apimodelPath = "../pkg/acsengine/testdata/simple/kubernetes.json"
	g.set = []string{"agentPoolProfiles[0].count=1"}
	err = g.mergeAPIModel()
	if err != nil {
		t.Fatalf("unexpected error calling mergeAPIModel with one --set flag to override an array property: %s", err.Error())
	}
}

func TestGenerateCmdMLoadAPIModel(t *testing.T) {
	g := &generateCmd{}
	r := &cobra.Command{}

	g.apimodelPath = "../pkg/acsengine/testdata/simple/kubernetes.json"
	g.set = []string{"agentPoolProfiles[0].count=1"}

	g.validate(r, []string{"../pkg/acsengine/testdata/simple/kubernetes.json"})
	g.mergeAPIModel()
	err := g.loadAPIModel(r, []string{"../pkg/acsengine/testdata/simple/kubernetes.json"})
	if err != nil {
		t.Fatalf("unexpected error loading api model: %s", err.Error())
	}
}

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

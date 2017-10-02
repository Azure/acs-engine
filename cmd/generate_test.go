package cmd

import (
	"testing"
)

func testGenerateCmdValidate(t *testing.T) {
	g := generateCmd{}
	r := NewRootCmd()

	// 1 arg
	err := g.validate(r, []string{"arg0"})
	if err != nil {
		t.Fatalf("unexpected error validating 1 arg")
	}
	// 0 args
	err = g.validate(r, []string{""})
	t.Logf(err.Error())
	if err == nil {
		t.Fatalf("expected error validating 0 args")
	}

	// more than 1 arg
	err = g.validate(r, []string{"arg0", "arg1"})
	t.Logf(err.Error())
	if err == nil {
		t.Fatalf("expected error validating multiple args")
	}
}

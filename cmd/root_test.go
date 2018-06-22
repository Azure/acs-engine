package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestNewRootCmd(t *testing.T) {
	output := NewRootCmd()
	if output.Use != rootName || output.Short != rootShortDescription || output.Long != rootLongDescription {
		t.Fatalf("root command should have use %s equal %s, short %s equal %s and long %s equal to %s", output.Use, rootName, output.Short, rootShortDescription, output.Long, rootLongDescription)
	}
	expectedCommands := []*cobra.Command{newDcosUpgradeCmd(), newDeployCmd(), newGenerateCmd(), newOrchestratorsCmd(), newScaleCmd(), newUpgradeCmd(), newVersionCmd()}
	rc := output.Commands()
	for i, c := range expectedCommands {
		if rc[i].Short != c.Short {
			t.Fatalf("root command should have command %s", c.Use)
		}
	}
}

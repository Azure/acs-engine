package cmd_test

import (
	"testing"

	. "github.com/Azure/acs-engine/pkg/test"
)

func TestCmd(t *testing.T) {

	RunSpecsWithReporters(t, "cmd", "Cmd Suite")

}

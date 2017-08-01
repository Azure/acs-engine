package cmd_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/reporters"
	. "github.com/onsi/gomega"
)

func TestCmd(t *testing.T) {
	RegisterFailHandler(Fail)
	//TODO: get the absolute path instead of ../test/junit
	junitReporter := reporters.NewJUnitReporter("../test/junit/cmd-junit.xml")
	RunSpecsWithDefaultAndCustomReporters(t, "Cmd Suite", []Reporter{junitReporter})
}

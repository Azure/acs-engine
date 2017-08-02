package cmd_test

import (
	"fmt"
	"testing"

	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/reporters"
	. "github.com/onsi/gomega"
)

var (
	JUnitOutDir = ""
)

func TestCmd(t *testing.T) {

	RegisterFailHandler(Fail)
	dir, _ := os.Getwd()
	localdir := filepath.Base(dir)
	junitReporter := reporters.NewJUnitReporter(filepath.Join(JUnitOutDir, fmt.Sprintf("%s-junit.xml", localdir)))
	RunSpecsWithDefaultAndCustomReporters(t, "Cmd Suite", []Reporter{junitReporter})
}

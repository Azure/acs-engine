package cmd

import (
	"fmt"

	logtest "github.com/Sirupsen/logrus/hooks/test"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("the version command", func() {
	It("should print the version of ACS-Engine", func() {
		command := newVersionCmd()
		hook := logtest.NewGlobal()
		command.Run(command, nil)
		Expect(hook.LastEntry().Message).To(Equal(fmt.Sprintf("ACS-Engine Version: %s (%s)", BuildTag, BuildSHA)))
	})
})

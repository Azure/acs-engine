package cmd

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/pflag"
)

var _ = Describe("the upgrade command", func() {

	It("should create an upgrade command", func() {

		flags := &pflag.FlagSet{
			pflag.formal: {
				"location": &pflag.Flag{
					Name:      "location",
					Shorthand: "l",
					Usage:     "location the cluster is deployed in",
					Value:     "",
					DefValue:  "",
				},
				"resource-group": &pflag.Flag{
					Name:      "resource-group",
					Shorthand: "g",
					Usage:     "the resource group where the cluster is deployed",
					Value:     "",
					DefValue:  "",
				},
				"deployment-dir": &pflag.Flag{
					Name:      "deployment-dir",
					Shorthand: "",
					Usage:     "the location of the output from `generate`",
					Value:     "",
					DefValue:  "",
				},
				"upgrade-version": &pflag.Flag{
					Name:      "upgrade-version",
					Shorthand: "",
					Usage:     "desired kubernetes version",
					Value:     "",
					DefValue:  "",
				},
			},
		}

		output := newUpgradeCmd()

		Expect(output.Use).Should(Equal(upgradeName))
		Expect(output.Short).Should(Equal(upgradeShortDescription))
		Expect(output.Long).Should(Equal(upgradeLongDescription))
		Expect(output.Flag("location")).Should(Equal(flags.Lookup("location")))
		Expect(output.Flag("location")).Should(Equal(flags.Lookup("location")))
		Expect(output.Flag("resource-group")).Should(Equal(flags.Lookup("resource-group")))
	})

})

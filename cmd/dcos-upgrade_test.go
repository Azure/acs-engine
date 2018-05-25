package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"
)

var _ = Describe("the upgrade command", func() {

	It("should create a DCOS upgrade command", func() {
		output := newDcosUpgradeCmd()

		Expect(output.Use).Should(Equal(dcosUpgradeName))
		Expect(output.Short).Should(Equal(dcosUpgradeShortDescription))
		Expect(output.Long).Should(Equal(dcosUpgradeLongDescription))
		Expect(output.Flags().Lookup("location")).NotTo(BeNil())
		Expect(output.Flags().Lookup("resource-group")).NotTo(BeNil())
		Expect(output.Flags().Lookup("deployment-dir")).NotTo(BeNil())
		Expect(output.Flags().Lookup("ssh-private-key-path")).NotTo(BeNil())
		Expect(output.Flags().Lookup("upgrade-version")).NotTo(BeNil())
	})

	It("should validate DCOS upgrade command", func() {
		r := &cobra.Command{}
		privKey, err := ioutil.TempFile("", "id_rsa")
		Expect(err).To(BeNil())
		defer os.Remove(privKey.Name())
		cases := []struct {
			uc          *dcosUpgradeCmd
			expectedErr error
		}{
			{
				uc: &dcosUpgradeCmd{
					resourceGroupName:   "",
					deploymentDirectory: "_output/test",
					upgradeVersion:      "1.8.9",
					location:            "centralus",
					sshPrivateKeyPath:   privKey.Name(),
					authArgs: authArgs{
						rawSubscriptionID: "99999999-0000-0000-0000-000000000000",
					},
				},
				expectedErr: fmt.Errorf("--resource-group must be specified"),
			},
			{
				uc: &dcosUpgradeCmd{
					resourceGroupName:   "test",
					deploymentDirectory: "_output/test",
					upgradeVersion:      "1.8.9",
					location:            "",
					sshPrivateKeyPath:   privKey.Name(),
					authArgs: authArgs{
						rawSubscriptionID: "99999999-0000-0000-0000-000000000000",
					},
				},
				expectedErr: fmt.Errorf("--location must be specified"),
			},
			{
				uc: &dcosUpgradeCmd{
					resourceGroupName:   "test",
					deploymentDirectory: "_output/test",
					upgradeVersion:      "",
					location:            "southcentralus",
					sshPrivateKeyPath:   privKey.Name(),
					authArgs: authArgs{
						rawSubscriptionID: "99999999-0000-0000-0000-000000000000",
					},
				},
				expectedErr: fmt.Errorf("--upgrade-version must be specified"),
			},
			{
				uc: &dcosUpgradeCmd{
					resourceGroupName:   "test",
					deploymentDirectory: "",
					upgradeVersion:      "1.9.0",
					location:            "southcentralus",
					sshPrivateKeyPath:   privKey.Name(),
					authArgs: authArgs{
						rawSubscriptionID: "99999999-0000-0000-0000-000000000000",
					},
				},
				expectedErr: fmt.Errorf("--deployment-dir must be specified"),
			},
			{
				uc: &dcosUpgradeCmd{
					resourceGroupName:   "test",
					deploymentDirectory: "",
					upgradeVersion:      "1.9.0",
					location:            "southcentralus",
					sshPrivateKeyPath:   privKey.Name(),
					authArgs: authArgs{
						rawSubscriptionID: "99999999-0000-0000-0000-000000000000",
					},
				},
				expectedErr: fmt.Errorf("--deployment-dir must be specified"),
			},
			{
				uc: &dcosUpgradeCmd{
					resourceGroupName:   "test",
					deploymentDirectory: "_output/mydir",
					upgradeVersion:      "1.9.0",
					location:            "southcentralus",
					sshPrivateKeyPath:   privKey.Name(),
					authArgs:            authArgs{},
				},
				expectedErr: fmt.Errorf("--subscription-id is required (and must be a valid UUID)"),
			},
			{
				uc: &dcosUpgradeCmd{
					resourceGroupName:   "test",
					deploymentDirectory: "_output/mydir",
					upgradeVersion:      "1.9.0",
					location:            "southcentralus",
					authArgs:            authArgs{},
				},
				expectedErr: fmt.Errorf("ssh-private-key-path must be specified: open _output/mydir/id_rsa: no such file or directory"),
			},
			{
				uc: &dcosUpgradeCmd{
					resourceGroupName:   "test",
					deploymentDirectory: "_output/mydir",
					upgradeVersion:      "1.9.0",
					location:            "southcentralus",
					sshPrivateKeyPath:   privKey.Name(),
					authArgs: authArgs{
						rawSubscriptionID:   "99999999-0000-0000-0000-000000000000",
						RawAzureEnvironment: "AzurePublicCloud",
						AuthMethod:          "device",
					},
				},
				expectedErr: nil,
			},
		}

		for _, c := range cases {
			err := c.uc.validate(r)

			if c.expectedErr != nil && err != nil {
				Expect(err.Error()).To(Equal(c.expectedErr.Error()))
			} else {
				Expect(err).To(BeNil())
				Expect(c.expectedErr).To(BeNil())
			}
		}

	})
})

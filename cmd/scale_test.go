package cmd

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func TestNewScaleCmd(t *testing.T) {
	output := newScaleCmd()
	if output.Use != scaleName || output.Short != scaleShortDescription || output.Long != scaleLongDescription {
		t.Fatalf("scale command should have use %s equal %s, short %s equal %s and long %s equal to %s", output.Use, scaleName, output.Short, scaleShortDescription, output.Long, scaleLongDescription)
	}

	expectedFlags := []string{"location", "resource-group", "deployment-dir", "new-node-count", "node-pool", "master-FQDN"}
	for _, f := range expectedFlags {
		if output.Flags().Lookup(f) == nil {
			t.Fatalf("scale command should have flag %s", f)
		}
	}
}

func TestScaleCmdValidate(t *testing.T) {
	r := &cobra.Command{}

	cases := []struct {
		sc          *scaleCmd
		expectedErr error
	}{
		{
			sc: &scaleCmd{
				location:             "centralus",
				resourceGroupName:    "",
				deploymentDirectory:  "_output/test",
				agentPoolToScale:     "agentpool1",
				newDesiredAgentCount: 5,
				masterFQDN:           "test",
			},
			expectedErr: errors.New("--resource-group must be specified"),
		},
		{
			sc: &scaleCmd{
				location:             "",
				resourceGroupName:    "testRG",
				deploymentDirectory:  "_output/test",
				agentPoolToScale:     "agentpool1",
				newDesiredAgentCount: 5,
				masterFQDN:           "test",
			},
			expectedErr: errors.New("--location must be specified"),
		},
		{
			sc: &scaleCmd{
				location:            "centralus",
				resourceGroupName:   "testRG",
				deploymentDirectory: "_output/test",
				agentPoolToScale:    "agentpool1",
				masterFQDN:          "test",
			},
			expectedErr: errors.New("--new-node-count must be specified"),
		},
		{
			sc: &scaleCmd{
				location:             "centralus",
				resourceGroupName:    "testRG",
				deploymentDirectory:  "",
				agentPoolToScale:     "agentpool1",
				newDesiredAgentCount: 5,
				masterFQDN:           "test",
			},
			expectedErr: errors.New("--deployment-dir must be specified"),
		},
		{
			sc: &scaleCmd{
				location:             "centralus",
				resourceGroupName:    "testRG",
				deploymentDirectory:  "_output/test",
				agentPoolToScale:     "agentpool1",
				newDesiredAgentCount: 5,
				masterFQDN:           "test",
			},
			expectedErr: nil,
		},
	}

	for _, c := range cases {
		err := c.sc.validate(r)
		if err != nil && c.expectedErr != nil {
			if err.Error() != c.expectedErr.Error() {
				t.Fatalf("expected validate scale command to return error %s, but instead got %s", c.expectedErr.Error(), err.Error())
			}
		} else {
			if c.expectedErr != nil {
				t.Fatalf("expected validate scale command to return error %s, but instead got no error", c.expectedErr.Error())
			} else if err != nil {
				t.Fatalf("expected validate scale command to return no error, but instead got %s", err.Error())
			}
		}
	}
}

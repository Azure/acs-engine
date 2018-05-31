package vlabs

import (
	"errors"
	"reflect"
	"testing"

	"github.com/Azure/acs-engine/pkg/api/common"
)

func TestValidateAgentPoolProfile(t *testing.T) {
	tests := []struct {
		name string

		agent *AgentPoolProfile

		expectedErr error
	}{
		{
			name: "valid agentpoolprofile",

			agent: &AgentPoolProfile{
				Name: "validname",
			},

			expectedErr: nil,
		},
		{
			name: "valid agentpoolprofile with imageref",

			agent: &AgentPoolProfile{
				Name: "validname",
				ImageRef: &ImageReference{
					Name:          "rhel7",
					ResourceGroup: "images",
				},
			},

			expectedErr: nil,
		},
		{
			name: "invalid name",

			agent: &AgentPoolProfile{
				Name: "invalid-name",
			},

			expectedErr: errors.New("pool name 'invalid-name' is invalid. A pool name must start with a lowercase letter, have max length of 12, and only have characters a-z0-9"),
		},
		{
			name: "invalid imageref - missing name",

			agent: &AgentPoolProfile{
				Name: "validname",
				ImageRef: &ImageReference{
					ResourceGroup: "images",
				},
			},

			expectedErr: errors.New("imageName needs to be specified when imageResourceGroup is provided"),
		},
		{
			name: "invalid imageref - missing resource group",

			agent: &AgentPoolProfile{
				Name: "validname",
				ImageRef: &ImageReference{
					Name: "rhel7",
				},
			},

			expectedErr: errors.New("imageResourceGroup needs to be specified when imageName is provided"),
		},
	}

	for _, test := range tests {
		err := test.agent.Validate()
		if !reflect.DeepEqual(err, test.expectedErr) {
			t.Logf("scenario %q:", test.name)
			t.Errorf("unexpected error: %v\nexpected error: %v", err, test.expectedErr)
		}
	}
}

func TestValidateAgents(t *testing.T) {
	tests := []struct {
		name string

		orchestratorProfile *OrchestratorProfile
		profiles            []*AgentPoolProfile

		expectedErr error
	}{
		{
			name: "valid agents",

			profiles: []*AgentPoolProfile{
				{
					Name: "foo",
					ImageRef: &ImageReference{
						Name:          "rhel7",
						ResourceGroup: "images",
					},
					Role: AgentPoolProfileRoleEmpty,
				},
				{
					Name: "bar",
					ImageRef: &ImageReference{
						Name:          "ubuntu",
						ResourceGroup: "images",
					},
					Role: AgentPoolProfileRoleEmpty,
				},
			},

			expectedErr: nil,
		},
		{
			name: "valid openshift agents",

			orchestratorProfile: &OrchestratorProfile{
				OrchestratorType: common.OpenShift,
			},
			profiles: []*AgentPoolProfile{
				{
					Name: "foo",
					ImageRef: &ImageReference{
						Name:          "rhel7",
						ResourceGroup: "images",
					},
					Role:                AgentPoolProfileRoleInfra,
					AvailabilityProfile: common.AvailabilitySet,
				},
				{
					Name: "bar",
					ImageRef: &ImageReference{
						Name:          "rhel7",
						ResourceGroup: "images",
					},
					Role:                AgentPoolProfileRoleEmpty,
					AvailabilityProfile: common.AvailabilitySet,
				},
			},

			expectedErr: nil,
		},
		{
			name: "invalid role",

			profiles: []*AgentPoolProfile{
				{
					Name: "foo",
					ImageRef: &ImageReference{
						Name:          "rhel7",
						ResourceGroup: "images",
					},
					Role: AgentPoolProfileRoleInfra,
				},
				{
					Name: "bar",
					ImageRef: &ImageReference{
						Name:          "ubuntu",
						ResourceGroup: "images",
					},
					Role: AgentPoolProfileRoleEmpty,
				},
			},

			expectedErr: errors.New(`role "infra" is not supported by orchestrator "Kubernetes"`),
		},
		{
			name: "invalid openshift availability profile",

			orchestratorProfile: &OrchestratorProfile{
				OrchestratorType: common.OpenShift,
			},
			profiles: []*AgentPoolProfile{
				{
					Name: "foo",
					ImageRef: &ImageReference{
						Name:          "rhel7",
						ResourceGroup: "images",
					},
					Role: AgentPoolProfileRoleInfra,
				},
				{
					Name: "bar",
					ImageRef: &ImageReference{
						Name:          "ubuntu",
						ResourceGroup: "images",
					},
					Role: AgentPoolProfileRoleEmpty,
				},
			},

			expectedErr: errors.New("only AvailabilityProfile: AvailabilitySet is supported for Orchestrator 'OpenShift'"),
		},
	}

	for _, test := range tests {
		t.Logf("scenario %q", test.name)

		err := validateAgents(test.orchestratorProfile, test.profiles)
		if !reflect.DeepEqual(err, test.expectedErr) {
			t.Errorf("unexpected error: %v\nexpected error: %v", err, test.expectedErr)
		}
	}
}

func TestValidateCertificateProfile(t *testing.T) {
	tests := []struct {
		name string

		orchestratorProfile *OrchestratorProfile
		certificateProfile  *CertificateProfile

		expectedErr error
	}{
		{
			name: "openshift orchestrator does not use certificate profile",

			orchestratorProfile: &OrchestratorProfile{
				OrchestratorType: common.OpenShift,
			},
			certificateProfile: nil,

			expectedErr: nil,
		},
		{
			name: "kubernetes orchestrator requires certificate profile",

			orchestratorProfile: &OrchestratorProfile{
				OrchestratorType: common.Kubernetes,
			},
			certificateProfile: &CertificateProfile{
				APIServerCertificate: "CERT",
			},

			expectedErr: nil,
		},
		{
			name: "invalid kubernetes orchestrator",

			orchestratorProfile: &OrchestratorProfile{
				OrchestratorType: common.Kubernetes,
			},
			certificateProfile: nil,

			expectedErr: errors.New("certificateProfile is required"),
		},
	}

	for _, test := range tests {
		t.Logf("scenario %q", test.name)

		err := validateCertificateProfile(test.orchestratorProfile, test.certificateProfile)
		if !reflect.DeepEqual(err, test.expectedErr) {
			t.Errorf("unexpected error: %v\nexpected error: %v", err, test.expectedErr)
		}
	}
}

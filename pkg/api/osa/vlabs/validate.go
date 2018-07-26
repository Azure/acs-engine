package vlabs

import (
	"fmt"
	"regexp"

	"k8s.io/apimachinery/pkg/util/errors"
)

var regexRfc1123 = regexp.MustCompile(`(?i)` +
	`^([a-z0-9]|[a-z0-9][-a-z0-9]{0,61}[a-z0-9])` +
	`(\.([a-z0-9]|[a-z0-9][-a-z0-9]{0,61}[a-z0-9]))*$`)

func validateHostname(h string, f string) error {
	if h == "" {
		return fmt.Errorf("%s must not be empty", f)
	}

	if !(len(h) <= 255 && regexRfc1123.MatchString(h)) {
		return fmt.Errorf("invalid %s %q", f, h)
	}

	return nil
}

var regexAgentPoolName = regexp.MustCompile(`^[a-z][a-z0-9]{0,11}$`)

// Validate validates an OpenShiftCluster struct
func (oc *OpenShiftCluster) Validate() error {
	if oc.Location == "" {
		return fmt.Errorf("location must not be empty")
	}
	if oc.Name == "" {
		return fmt.Errorf("name must not be empty")
	}

	return oc.Properties.Validate()
}

// Validate validates a Properties struct
func (p *Properties) Validate() error {
	switch p.ProvisioningState {
	case "", Creating, Updating, Failed, Succeeded, Deleting, Migrating, Upgrading:
	default:
		return fmt.Errorf("invalid provisioningState %q", p.ProvisioningState)
	}

	if p.OpenShiftVersion == "" {
		return fmt.Errorf("openShiftVersion must not be empty")
	}

	if err := validateHostname(p.PublicHostname, "publicHostname"); err != nil {
		return err
	}

	if err := validateHostname(p.RoutingConfigSubdomain, "routingConfigSubdomain"); err != nil {
		return err
	}

	if err := ValidatePools(p.ComputePools, p.InfraPool); err != nil {
		return err
	}

	return p.ServicePrincipalProfile.Validate()
}

// ValidatePools validate pools validates infra and compute pools.
func ValidatePools(compute AgentPoolProfiles, infra *InfraPoolProfile) error {
	if infra == nil || len(compute) == 0 {
		return fmt.Errorf("both infra and compute pools are required")
	}

	var errs []error

	if len(compute) > 1 {
		errs = append(errs, fmt.Errorf("only one compute pool is currently supported"))
	}

	names := map[string]struct{}{}
	names[infra.Name] = struct{}{}

	for _, app := range compute {
		if _, found := names[app.Name]; found {
			errs = append(errs, fmt.Errorf("duplicate name %q", app.Name))
		}
		names[app.Name] = struct{}{}

		if err := app.Validate(); err != nil {
			errs = append(errs, err)
		}
	}

	if err := infra.Validate(); err != nil {
		errs = append(errs, err)
	}

	return errors.NewAggregate(errs)
}

// Validate validates an InfraPoolProfile.
func (infra *InfraPoolProfile) Validate() error {
	var errs []error

	if !regexAgentPoolName.MatchString(infra.Name) {
		errs = append(errs, fmt.Errorf("invalid name %q", infra.Name))
	}

	if infra.Count < 3 {
		errs = append(errs, fmt.Errorf("must have at least 3 infra nodes"))
	}

	if infra.VMSize == "" {
		errs = append(errs, fmt.Errorf("vmSize must not be empty"))
	}

	return errors.NewAggregate(errs)
}

// Validate validates an AgentPoolProfile struct
func (app *AgentPoolProfile) Validate() error {
	if !regexAgentPoolName.MatchString(app.Name) {
		return fmt.Errorf("invalid name %q", app.Name)
	}

	if app.Count < 1 {
		return fmt.Errorf("invalid count %q", app.Count)
	}

	if app.VMSize == "" {
		return fmt.Errorf("vmSize must not be empty")
	}

	return nil
}

// Validate validates a ServicePrincipalProfile struct
func (spp *ServicePrincipalProfile) Validate() error {
	if spp.ClientID == "" {
		return fmt.Errorf("clientId must not be empty")
	}

	if spp.Secret == "" {
		return fmt.Errorf("secret must not be empty")
	}

	return nil
}

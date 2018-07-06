package vlabs

import (
	"fmt"
	"regexp"
)

var regexRfc1123 = regexp.MustCompile(`(?i)` +
	`^([a-z0-9]|[a-z0-9][-a-z0-9]{0,61}[a-z0-9])` +
	`(\.([a-z0-9]|[a-z0-9][-a-z0-9]{0,61}[a-z0-9]))*$`)

func isValidHostname(h string) bool {
	return len(h) <= 255 && regexRfc1123.MatchString(h)
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

	if oc.Properties == nil {
		return nil
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

	if p.PublicHostname != "" && !isValidHostname(p.PublicHostname) {
		return fmt.Errorf("invalid publicHostname %q", p.PublicHostname)
	}

	if p.FQDN != "" && !isValidHostname(p.FQDN) {
		return fmt.Errorf("invalid fqdn %q", p.FQDN)
	}

	if p.RoutingConfigSubdomain == "" && !isValidHostname(p.RoutingConfigSubdomain) {
		return fmt.Errorf("invalid routingConfigSubdomain %q", p.RoutingConfigSubdomain)
	}

	if p.RoutingConfigFQDN == "" && !isValidHostname(p.RoutingConfigFQDN) {
		return fmt.Errorf("invalid routingConfigFqdn %q", p.RoutingConfigFQDN)
	}

	if err := p.AgentPoolProfiles.Validate(); err != nil {
		return err
	}

	return p.ServicePrincipalProfile.Validate()
}

// Validate validates an AgentPoolProfiles slice
func (apps AgentPoolProfiles) Validate() error {
	names := map[string]struct{}{}

	for _, app := range apps {
		if _, found := names[app.Name]; found {
			return fmt.Errorf("duplicate name %q", app.Name)
		}
		names[app.Name] = struct{}{}

		if err := app.Validate(); err != nil {
			return err
		}
	}

	return nil
}

// Validate validates an AgentPoolProfile struct
func (app *AgentPoolProfile) Validate() error {
	if !regexAgentPoolName.MatchString(app.Name) {
		return fmt.Errorf("invalid name %q", app.Name)
	}

	switch app.Role {
	case AgentPoolProfileRoleCompute, AgentPoolProfileRoleInfra, AgentPoolProfileRoleMaster:
	default:
		return fmt.Errorf("invalid role %q", app.Role)
	}

	if app.Count < 1 {
		return fmt.Errorf("invalid count %q", app.Count)
	}

	if app.VMSize == "" {
		return fmt.Errorf("vmSize must not be empty")
	}

	switch app.OSType {
	case OSTypeLinux, OSTypeWindows:
	default:
		return fmt.Errorf("invalid osType %q", app.OSType)
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

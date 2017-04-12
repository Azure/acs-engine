package vlabs

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"regexp"
)

// Validate implements APIObject
func (o *OrchestratorProfile) Validate() error {
	switch o.OrchestratorType {
	case DCOS:
	case DCOS190:
	case DCOS188:
	case DCOS187:
	case DCOS184:
	case DCOS173:
	case Swarm:
	case Kubernetes:
		if o.KubernetesConfig != nil {
			err := o.KubernetesConfig.Validate()
			if err != nil {
				return err
			}
		}
	case SwarmMode:
	default:
		return fmt.Errorf("OrchestratorProfile has unknown orchestrator: %q", o.OrchestratorType)
	}

	if o.OrchestratorType != Kubernetes && o.KubernetesConfig != nil &&
		(*o.KubernetesConfig != KubernetesConfig{}) {
		return fmt.Errorf("KubernetesConfig can be specified only when OrchestratorType is Kubernetes")
	}

	return nil
}

// Validate implements APIObject
func (m *MasterProfile) Validate() error {
	if m.Count != 1 && m.Count != 3 && m.Count != 5 {
		return fmt.Errorf("MasterProfile count needs to be 1, 3, or 5")
	}
	if e := validateName(m.DNSPrefix, "MasterProfile.DNSPrefix"); e != nil {
		return e
	}
	if e := validateDNSName(m.DNSPrefix); e != nil {
		return e
	}
	if e := validateName(m.VMSize, "MasterProfile.VMSize"); e != nil {
		return e
	}
	if m.IPAddressCount != 0 && (m.IPAddressCount < MinIPAddressCount || m.IPAddressCount > MaxIPAddressCount) {
		return fmt.Errorf("MasterProfile.IPAddressCount needs to be in the range [%d,%d]", MinIPAddressCount, MaxIPAddressCount)
	}
	return nil
}

// Validate implements APIObject
func (a *AgentPoolProfile) Validate() error {
	if e := validateName(a.Name, "AgentPoolProfile.Name"); e != nil {
		return e
	}
	if e := validatePoolName(a.Name); e != nil {
		return e
	}
	if a.Count < MinAgentCount || a.Count > MaxAgentCount {
		return fmt.Errorf("AgentPoolProfile count needs to be in the range [%d,%d]", MinAgentCount, MaxAgentCount)
	}
	if e := validateName(a.VMSize, "AgentPoolProfile.VMSize"); e != nil {
		return e
	}
	if len(a.Ports) > 0 {
		if e := validateUniquePorts(a.Ports, a.Name); e != nil {
			return e
		}
		for _, port := range a.Ports {
			if port < MinPort || port > MaxPort {
				return fmt.Errorf("AgentPoolProfile Ports must be in the range[%d, %d]", MinPort, MaxPort)
			}
		}
		if e := validateName(a.DNSPrefix, "AgentPoolProfile.DNSPrefix when specifying AgentPoolProfile Ports"); e != nil {
			return e
		}
		if e := validateDNSName(a.DNSPrefix); e != nil {
			return e
		}
	}
	if len(a.DiskSizesGB) > 0 {
		if len(a.StorageProfile) == 0 {
			return fmt.Errorf("property 'StorageProfile' must be set to either '%s' or '%s' when attaching disks", StorageAccount, ManagedDisks)
		}
		if len(a.AvailabilityProfile) == 0 {
			return fmt.Errorf("property 'AvailabilityProfile' must be set to either '%s' or '%s' when attaching disks", VirtualMachineScaleSets, AvailabilitySet)
		}
		if a.StorageProfile == StorageAccount && (a.AvailabilityProfile == VirtualMachineScaleSets) {
			return fmt.Errorf("VirtualMachineScaleSets does not support storage account attached disks.  Instead specify 'StorageAccount': '%s' or specify AvailabilityProfile '%s'", ManagedDisks, AvailabilitySet)
		}
	}
	for _, s := range a.DiskSizesGB {
		if s < MinDiskSizeGB || s > MaxDiskSizeGB {
			return fmt.Errorf("Invalid disk size %d specified for cluster cluster named '%s'.  The range of valid values are [%d, %d]", s, a.Name, MinDiskSizeGB, MaxDiskSizeGB)
		}
	}
	if len(a.DiskSizesGB) > MaxDisks {
		return fmt.Errorf("A maximum of %d disks may be specified.  %d disks were specified for cluster named '%s'", MaxDisks, len(a.DiskSizesGB), a.Name)
	}
	if len(a.Ports) == 0 && len(a.DNSPrefix) > 0 {
		return fmt.Errorf("AgentPoolProfile.Ports must be non empty when AgentPoolProfile.DNSPrefix is specified")
	}
	if a.IPAddressCount != 0 && (a.IPAddressCount < MinIPAddressCount || a.IPAddressCount > MaxIPAddressCount) {
		return fmt.Errorf("AgentPoolProfile.IPAddressCount needs to be in the range [%d,%d]", MinIPAddressCount, MaxIPAddressCount)
	}
	if e := validateStorageProfile(a.StorageProfile); e != nil {
		return e
	}
	return nil
}

func validateKeyVaultSecrets(secrets []KeyVaultSecrets, requireCertificateStore bool) error {
	for _, s := range secrets {
		if len(s.VaultCertificates) == 0 {
			return fmt.Errorf("Invalid KeyVaultSecrets must have no empty VaultCertificates")
		}
		if s.SourceVault.ID == "" {
			return fmt.Errorf("KeyVaultSecrets must have a SourceVault.ID")
		}
		for _, c := range s.VaultCertificates {
			if _, e := url.Parse(c.CertificateURL); e != nil {
				return fmt.Errorf("Certificate url was invalid. recieved error %s", e)
			}
			if e := validateName(c.CertificateStore, "KeyVaultCertificate.CertificateStore"); requireCertificateStore && e != nil {
				return fmt.Errorf("%s for certificates in a WindowsProfile", e)
			}
		}
	}
	return nil
}

// Validate implements APIObject
func (l *LinuxProfile) Validate() error {
	if e := validateName(l.AdminUsername, "LinuxProfile.AdminUsername"); e != nil {
		return e
	}
	if len(l.SSH.PublicKeys) != 1 {
		return errors.New("LinuxProfile.PublicKeys requires only 1 SSH Key")
	}
	if e := validateName(l.SSH.PublicKeys[0].KeyData, "LinuxProfile.PublicKeys.KeyData"); e != nil {
		return e
	}
	if e := validateKeyVaultSecrets(l.Secrets, false); e != nil {
		return e
	}
	return nil
}

// Validate implements APIObject
func (a *Properties) Validate() error {
	if e := a.OrchestratorProfile.Validate(); e != nil {
		return e
	}
	if e := a.validateNetworkPolicy(); e != nil {
		return e
	}
	if e := a.MasterProfile.Validate(); e != nil {
		return e
	}
	if e := validateUniqueProfileNames(a.AgentPoolProfiles); e != nil {
		return e
	}
	if a.OrchestratorProfile.OrchestratorType == Kubernetes && len(a.ServicePrincipalProfile.ClientID) == 0 {
		return fmt.Errorf("the service principal client ID must be specified with Orchestrator %s", a.OrchestratorProfile.OrchestratorType)
	}

	if a.OrchestratorProfile.OrchestratorType == Kubernetes && len(a.ServicePrincipalProfile.Secret) == 0 {
		return fmt.Errorf("the service principal client secrect must be specified with Orchestrator %s", a.OrchestratorProfile.OrchestratorType)
	}

	for _, agentPoolProfile := range a.AgentPoolProfiles {
		if e := agentPoolProfile.Validate(); e != nil {
			return e
		}
		switch agentPoolProfile.AvailabilityProfile {
		case AvailabilitySet:
		case VirtualMachineScaleSets:
		case "":
		default:
			{
				return fmt.Errorf("unknown availability profile type '%s' for agent pool '%s'.  Specify either %s, or %s", agentPoolProfile.AvailabilityProfile, agentPoolProfile.Name, AvailabilitySet, VirtualMachineScaleSets)
			}
		}
		switch agentPoolProfile.StorageProfile {
		case StorageAccount:
		case ManagedDisks:
		case "":
		default:
			{
				return fmt.Errorf("unknown storage type '%s' for agent pool '%s'.  Specify either %s, or %s", agentPoolProfile.StorageProfile, agentPoolProfile.Name, StorageAccount, ManagedDisks)
			}
		}
		/* this switch statement is left to protect newly added orchestrators until they support Managed Disks*/

		if agentPoolProfile.StorageProfile == ManagedDisks {
			switch a.OrchestratorProfile.OrchestratorType {
			case DCOS:
			case DCOS173:
			case DCOS184:
			case DCOS187:
			case DCOS188:
			case DCOS190:
			case Swarm:
			case Kubernetes:
			case SwarmMode:
			default:
				return fmt.Errorf("HA volumes are currently unsupported for Orchestrator %s", a.OrchestratorProfile.OrchestratorType)
			}
		}

		if len(agentPoolProfile.CustomNodeLabels) > 0 {
			switch a.OrchestratorProfile.OrchestratorType {
			case DCOS:
			case DCOS173:
			case DCOS184:
			case DCOS187:
			case DCOS188:
			case DCOS190:
			default:
				return fmt.Errorf("Agent Type attributes are only supported for DCOS.")
			}
		}
		if a.OrchestratorProfile.OrchestratorType == Kubernetes && (agentPoolProfile.AvailabilityProfile == VirtualMachineScaleSets || len(agentPoolProfile.AvailabilityProfile) == 0) {
			return fmt.Errorf("VirtualMachineScaleSets are not supported with Kubernetes since Kubernetes requires the ability to attach/detach disks.  To fix specify \"AvailabilityProfile\":\"%s\"", AvailabilitySet)
		}
		if a.OrchestratorProfile.OrchestratorType == Kubernetes && len(agentPoolProfile.DNSPrefix) > 0 {
			return errors.New("DNSPrefix not support for agent pools in Kubernetes - Kubernetes marks its own clusters public")
		}
		if agentPoolProfile.OSType == Windows {
			switch a.OrchestratorProfile.OrchestratorType {
			case Swarm:
			case SwarmMode:
			case Kubernetes:
			default:
				return fmt.Errorf("Orchestrator %s does not support Windows", a.OrchestratorProfile.OrchestratorType)
			}
			if len(a.WindowsProfile.AdminUsername) == 0 {
				return fmt.Errorf("WindowsProfile.AdminUsername must not be empty since agent pool '%s' specifies windows", agentPoolProfile.Name)
			}
			if len(a.WindowsProfile.AdminPassword) == 0 {
				return fmt.Errorf("WindowsProfile.AdminPassword must not be empty since  agent pool '%s' specifies windows", agentPoolProfile.Name)
			}
			if e := validateKeyVaultSecrets(a.WindowsProfile.Secrets, true); e != nil {
				return e
			}
		}
	}
	if e := a.LinuxProfile.Validate(); e != nil {
		return e
	}
	if e := validateVNET(a); e != nil {
		return e
	}
	return nil
}

func (a *Properties) validateNetworkPolicy() error {
	var networkPolicy string

	switch a.OrchestratorProfile.OrchestratorType {
	case Kubernetes:
		if a.OrchestratorProfile.KubernetesConfig != nil {
			networkPolicy = a.OrchestratorProfile.KubernetesConfig.NetworkPolicy
		}
	default:
		return nil
	}

	// Check NetworkPolicy has a valid value.
	valid := false
	for _, policy := range NetworkPolicyValues {
		if networkPolicy == policy {
			valid = true
			break
		}
	}
	if !valid {
		return fmt.Errorf("unknown networkPolicy '%s' specified", networkPolicy)
	}

	// Temporary safety check, to be removed when Windows support is added.
	if (networkPolicy == "calico" || networkPolicy == "azure") && a.HasWindows() {
		return fmt.Errorf("networkPolicy '%s' is not supporting windows agents", networkPolicy)
	}

	return nil
}

func validateName(name string, label string) error {
	if name == "" {
		return fmt.Errorf("%s must be a non-empty value", label)
	}
	return nil
}

func validatePoolName(poolName string) error {
	// we will cap at length of 12 and all lowercase letters since this makes up the VMName
	poolNameRegex := `^([a-z][a-z0-9]{0,11})$`
	re, err := regexp.Compile(poolNameRegex)
	if err != nil {
		return err
	}
	submatches := re.FindStringSubmatch(poolName)
	if len(submatches) != 2 {
		return fmt.Errorf("pool name '%s' is invalid. A pool name must start with a lowercase letter, have max length of 12, and only have characters a-z0-9", poolName)
	}
	return nil
}

func validateDNSName(dnsName string) error {
	dnsNameRegex := `^([A-Za-z][A-Za-z0-9-]{1,43}[A-Za-z0-9])$`
	re, err := regexp.Compile(dnsNameRegex)
	if err != nil {
		return err
	}
	if !re.MatchString(dnsName) {
		return fmt.Errorf("DNS name '%s' is invalid. The DNS name must contain between 3 and 45 characters.  The name can contain only letters, numbers, and hyphens.  The name must start with a letter and must end with a letter or a number (length was %d)", dnsName, len(dnsName))
	}
	return nil
}

func validateUniqueProfileNames(profiles []AgentPoolProfile) error {
	profileNames := make(map[string]bool)
	for _, profile := range profiles {
		if _, ok := profileNames[profile.Name]; ok {
			return fmt.Errorf("profile name '%s' already exists, profile names must be unique across pools", profile.Name)
		}
		profileNames[profile.Name] = true
	}
	return nil
}

func validateStorageProfile(storageProfile string) error {
	switch storageProfile {
	case StorageAccount:
	case ManagedDisks:
	case "":
	default:
		{
			return fmt.Errorf("Unknown storage type '%s' for agent pool. Specify either %s or %s", storageProfile, StorageAccount, ManagedDisks)
		}
	}
	return nil
}

func validateUniquePorts(ports []int, name string) error {
	portMap := make(map[int]bool)
	for _, port := range ports {
		if _, ok := portMap[port]; ok {
			return fmt.Errorf("agent profile '%s' has duplicate port '%d', ports must be unique", name, port)
		}
		portMap[port] = true
	}
	return nil
}

func validateVNET(a *Properties) error {
	isCustomVNET := a.MasterProfile.IsCustomVNET()
	for _, agentPool := range a.AgentPoolProfiles {
		if agentPool.IsCustomVNET() != isCustomVNET {
			return fmt.Errorf("Multiple VNET Subnet configurations specified.  The master profile and each agent pool profile must all specify a custom VNET Subnet, or none at all.")
		}
	}
	if isCustomVNET {
		subscription, resourcegroup, vnetname, _, e := GetVNETSubnetIDComponents(a.MasterProfile.VnetSubnetID)
		if e != nil {
			return e
		}

		for _, agentPool := range a.AgentPoolProfiles {
			agentSubID, agentRG, agentVNET, _, err := GetVNETSubnetIDComponents(agentPool.VnetSubnetID)
			if err != nil {
				return err
			}
			if agentSubID != subscription ||
				agentRG != resourcegroup ||
				agentVNET != vnetname {
				return errors.New("Multipe VNETS specified.  The master profile and each agent pool must reference the same VNET (but it is ok to reference different subnets on that VNET)")
			}
		}

		masterFirstIP := net.ParseIP(a.MasterProfile.FirstConsecutiveStaticIP)
		if masterFirstIP == nil {
			return fmt.Errorf("MasterProfile.FirstConsecutiveStaticIP (with VNET Subnet specification) '%s' is an invalid IP address", a.MasterProfile.FirstConsecutiveStaticIP)
		}
	}
	return nil
}

// GetVNETSubnetIDComponents extract subscription, resourcegroup, vnetname, subnetname from the vnetSubnetID
func GetVNETSubnetIDComponents(vnetSubnetID string) (string, string, string, string, error) {
	vnetSubnetIDRegex := `^\/subscriptions\/([^\/]*)\/resourceGroups\/([^\/]*)\/providers\/Microsoft.Network\/virtualNetworks\/([^\/]*)\/subnets\/([^\/]*)$`
	re, err := regexp.Compile(vnetSubnetIDRegex)
	if err != nil {
		return "", "", "", "", err
	}
	submatches := re.FindStringSubmatch(vnetSubnetID)
	if len(submatches) != 4 {
		return "", "", "", "", err
	}
	return submatches[1], submatches[2], submatches[3], submatches[4], nil
}

func (a *KubernetesConfig) Validate() error {
	if a.ClusterCidr != "" {
		_, _, err := net.ParseCIDR(a.ClusterCidr)
		if err != nil {
			return fmt.Errorf("OrchestratorProfile.KubernetesConfig.ClusterCidr '%s' is an invalid CIDR subnet", a.ClusterCidr)
		}
	}

	if a.DnsServiceIP != "" || a.ServiceCidr != "" {
		if a.DnsServiceIP == "" {
			return errors.New("OrchestratorProfile.KubernetesConfig.ServiceCidr must be specified when DnsServiceIP is")
		}
		if a.ServiceCidr == "" {
			return errors.New("OrchestratorProfile.KubernetesConfig.DnsServiceIP must be specified when ServiceCidr is")
		}

		dnsIp := net.ParseIP(a.DnsServiceIP)
		if dnsIp == nil {
			return fmt.Errorf("OrchestratorProfile.KubernetesConfig.DnsServiceIP '%s' is an invalid IP address", a.DnsServiceIP)
		}

		_, serviceCidr, err := net.ParseCIDR(a.ServiceCidr)
		if err != nil {
			return fmt.Errorf("OrchestratorProfile.KubernetesConfig.ServiceCidr '%s' is an invalid CIDR subnet", a.ServiceCidr)
		}

		// Finally validate that the DNS ip is within the subnet, and _not_ that subnet broadcast address, otherwise it won't work
		if !serviceCidr.Contains(dnsIp) {
			return fmt.Errorf("OrchestratorProfile.KubernetesConfig.DnsServiceIP '%s' is not within the ServiceCidr '%s'", a.DnsServiceIP, a.ServiceCidr)
		}

		broadcast := ip4BroadcastAddress(serviceCidr)
		if dnsIp.Equal(broadcast) {
			return fmt.Errorf("OrchestratorProfile.KubernetesConfig.DnsServiceIP '%s' cannot be the broadcast address of ServiceCidr '%s'", a.DnsServiceIP, a.ServiceCidr)
		}
	}

	return nil
}

// ip4BroadcastAddress returns the broadcase address for the given IP subnet.
func ip4BroadcastAddress(n *net.IPNet) net.IP {
	ip4 := n.IP.To4()
	if ip4 == nil {
		return nil
	}
	last := make(net.IP, len(ip4))
	copy(last, ip4)
	for i := range ip4 {
		last[i] |= ^n.Mask[i]
	}
	return last
}

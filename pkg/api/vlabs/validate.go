package vlabs

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"regexp"
	"time"
)

var keyvaultSecretPathRegex *regexp.Regexp

func init() {
	keyvaultSecretPathRegex = regexp.MustCompile(`^(/subscriptions/\S+/resourceGroups/\S+/providers/Microsoft.KeyVault/vaults/\S+)/secrets/([^/\s]+)(/(\S+))?$`)
}

// Validate implements APIObject
func (o *OrchestratorProfile) Validate() error {
	switch o.OrchestratorType {
	case DCOS:
		switch o.OrchestratorVersion {
		case DCOS173:
		case DCOS184:
		case DCOS187:
		case DCOS188:
		case DCOS190:
		case "":
		default:
			return fmt.Errorf("OrchestratorProfile has unknown orchestrator version: %s", o.OrchestratorVersion)
		}

	case Swarm:
	case SwarmMode:

	case Kubernetes:
		switch o.OrchestratorVersion {
		case Kubernetes170:
		case Kubernetes166:
		case Kubernetes162:
		case Kubernetes160:
		case Kubernetes157:
		case Kubernetes153:
		case "":
		default:
			return fmt.Errorf("OrchestratorProfile has unknown orchestrator version: %s", o.OrchestratorVersion)
		}

		if o.KubernetesConfig != nil {
			err := o.KubernetesConfig.Validate(o.OrchestratorVersion)
			if err != nil {
				return err
			}
		}

	default:
		return fmt.Errorf("OrchestratorProfile has unknown orchestrator: %s", o.OrchestratorType)
	}

	if o.OrchestratorType != Kubernetes && o.KubernetesConfig != nil && (*o.KubernetesConfig != KubernetesConfig{}) {
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
	if m.OSDiskSizeGB != 0 && (m.OSDiskSizeGB < MinDiskSizeGB || m.OSDiskSizeGB > MaxDiskSizeGB) {
		return fmt.Errorf("Invalid master os disk size of %d specified.  The range of valid values are [%d, %d]", m.OSDiskSizeGB, MinDiskSizeGB, MaxDiskSizeGB)
	}
	if m.IPAddressCount != 0 && (m.IPAddressCount < MinIPAddressCount || m.IPAddressCount > MaxIPAddressCount) {
		return fmt.Errorf("MasterProfile.IPAddressCount needs to be in the range [%d,%d]", MinIPAddressCount, MaxIPAddressCount)
	}
	if e := validateStorageProfile(m.StorageProfile); e != nil {
		return e
	}
	return nil
}

// Validate implements APIObject
func (a *AgentPoolProfile) Validate(orchestratorType OrchestratorType) error {
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
	if a.OSDiskSizeGB != 0 && (a.OSDiskSizeGB < MinDiskSizeGB || a.OSDiskSizeGB > MaxDiskSizeGB) {
		return fmt.Errorf("Invalid os disk size of %d specified.  The range of valid values are [%d, %d]", a.OSDiskSizeGB, MinDiskSizeGB, MaxDiskSizeGB)
	}
	if a.DNSPrefix != "" {
		if e := validateDNSName(a.DNSPrefix); e != nil {
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
		} else {
			a.Ports = []int{80, 443, 8080}
		}
	} else {
		if len(a.Ports) > 0 {
			return fmt.Errorf("AgentPoolProfile.Ports must be empty when AgentPoolProfile.DNSPrefix is empty")
		}
	}

	// for Kubernetes, we don't support AgentPoolProfile.DNSPrefix
	if orchestratorType == Kubernetes {
		if e := validateNameEmpty(a.DNSPrefix, "AgentPoolProfile.DNSPrefix"); e != nil {
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
		if s.SourceVault == nil {
			return fmt.Errorf("missing SourceVault in KeyVaultSecrets")
		}
		if s.SourceVault.ID == "" {
			return fmt.Errorf("KeyVaultSecrets must have a SourceVault.ID")
		}
		for _, c := range s.VaultCertificates {
			if _, e := url.Parse(c.CertificateURL); e != nil {
				return fmt.Errorf("Certificate url was invalid. received error %s", e)
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
	if a.OrchestratorProfile == nil {
		return fmt.Errorf("missing OrchestratorProfile")
	}
	if a.MasterProfile == nil {
		return fmt.Errorf("missing MasterProfile")
	}
	if a.LinuxProfile == nil {
		return fmt.Errorf("missing LinuxProfile")
	}
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

	if a.OrchestratorProfile.OrchestratorType == Kubernetes {
		useManagedIdentity := (a.OrchestratorProfile.KubernetesConfig != nil &&
			a.OrchestratorProfile.KubernetesConfig.UseManagedIdentity)

		if !useManagedIdentity {
			if len(a.ServicePrincipalProfile.ClientID) == 0 {
				return fmt.Errorf("the service principal client ID must be specified with Orchestrator %s", a.OrchestratorProfile.OrchestratorType)
			}

			if (len(a.ServicePrincipalProfile.Secret) == 0 && len(a.ServicePrincipalProfile.KeyvaultSecretRef) == 0) ||
				(len(a.ServicePrincipalProfile.Secret) != 0 && len(a.ServicePrincipalProfile.KeyvaultSecretRef) != 0) {
				return fmt.Errorf("either the service principal client secret or keyvault secret reference must be specified with Orchestrator %s", a.OrchestratorProfile.OrchestratorType)
			}

			if len(a.ServicePrincipalProfile.KeyvaultSecretRef) != 0 {
				parts := keyvaultSecretPathRegex.FindStringSubmatch(a.ServicePrincipalProfile.KeyvaultSecretRef)
				if len(parts) != 5 {
					return fmt.Errorf("service principal client keyvault secret reference is of incorrect format")
				}
			}
		}
	}

	for _, agentPoolProfile := range a.AgentPoolProfiles {
		if e := agentPoolProfile.Validate(a.OrchestratorProfile.OrchestratorType); e != nil {
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
			case Kubernetes:
			default:
				return fmt.Errorf("Agent Type attributes are only supported for DCOS and Kubernetes")
			}
		}
		if a.OrchestratorProfile.OrchestratorType == Kubernetes && (agentPoolProfile.AvailabilityProfile == VirtualMachineScaleSets || len(agentPoolProfile.AvailabilityProfile) == 0) {
			return fmt.Errorf("VirtualMachineScaleSets are not supported with Kubernetes since Kubernetes requires the ability to attach/detach disks.  To fix specify \"AvailabilityProfile\":\"%s\"", AvailabilitySet)
		}
		if a.OrchestratorProfile.OrchestratorType == Kubernetes && len(agentPoolProfile.DNSPrefix) > 0 {
			return errors.New("DNSPrefix not support for agent pools in Kubernetes - Kubernetes marks its own clusters public")
		}
		if agentPoolProfile.OSType == Windows {
			if a.WindowsProfile == nil {
				return fmt.Errorf("missing WindowsProfile")
			}
			switch a.OrchestratorProfile.OrchestratorType {
			case Swarm:
			case SwarmMode:
			case Kubernetes:
			default:
				return fmt.Errorf("Orchestrator %s does not support Windows", a.OrchestratorProfile.OrchestratorType)
			}

			if a.WindowsProfile == nil {
				return fmt.Errorf("WindowsProfile must not be empty since agent pool '%s' specifies windows", agentPoolProfile.Name)
			}

			if len(a.WindowsProfile.AdminUsername) == 0 {
				return fmt.Errorf("WindowsProfile.AdminUsername must not be empty since agent pool '%s' specifies windows", agentPoolProfile.Name)
			}
			if len(a.WindowsProfile.AdminPassword) == 0 {
				return fmt.Errorf("WindowsProfile.AdminPassword must not be empty since agent pool '%s' specifies windows", agentPoolProfile.Name)
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

// Validate validates the KubernetesConfig.
func (a *KubernetesConfig) Validate(k8sVersion OrchestratorVersion) error {
	// number of minimum retries allowed for kubelet to post node status
	const minKubeletRetries = 4
	// k8s versions that have cloudprovider backoff enabled
	var backoffEnabledVersions = map[OrchestratorVersion]bool{
		Kubernetes166: true,
		Kubernetes170: true,
	}
	// k8s versions that have cloudprovider rate limiting enabled (currently identical with backoff enabled versions)
	ratelimitEnabledVersions := backoffEnabledVersions

	if a.ClusterSubnet != "" {
		_, _, err := net.ParseCIDR(a.ClusterSubnet)
		if err != nil {
			return fmt.Errorf("OrchestratorProfile.KubernetesConfig.ClusterSubnet '%s' is an invalid subnet", a.ClusterSubnet)
		}
	}

	if a.DockerBridgeSubnet != "" {
		_, _, err := net.ParseCIDR(a.DockerBridgeSubnet)
		if err != nil {
			return fmt.Errorf("OrchestratorProfile.KubernetesConfig.DockerBridgeSubnet '%s' is an invalid subnet", a.DockerBridgeSubnet)
		}
	}

	if a.NodeStatusUpdateFrequency != "" {
		_, err := time.ParseDuration(a.NodeStatusUpdateFrequency)
		if err != nil {
			return fmt.Errorf("OrchestratorProfile.KubernetesConfig.NodeStatusUpdateFrequency '%s' is not a valid duration", a.NodeStatusUpdateFrequency)
		}
		if a.CtrlMgrNodeMonitorGracePeriod == "" {
			return fmt.Errorf("OrchestratorProfile.KubernetesConfig.NodeStatusUpdateFrequency was set to '%s' but OrchestratorProfile.KubernetesConfig.CtrlMgrNodeMonitorGracePeriod was not set", a.NodeStatusUpdateFrequency)
		}
	}

	if a.CtrlMgrNodeMonitorGracePeriod != "" {
		_, err := time.ParseDuration(a.CtrlMgrNodeMonitorGracePeriod)
		if err != nil {
			return fmt.Errorf("OrchestratorProfile.KubernetesConfig.CtrlMgrNodeMonitorGracePeriod '%s' is not a valid duration", a.CtrlMgrNodeMonitorGracePeriod)
		}
		if a.NodeStatusUpdateFrequency == "" {
			return fmt.Errorf("OrchestratorProfile.KubernetesConfig.CtrlMgrNodeMonitorGracePeriod was set to '%s' but OrchestratorProfile.KubernetesConfig.NodeStatusUpdateFrequency was not set", a.NodeStatusUpdateFrequency)
		}
	}

	if a.NodeStatusUpdateFrequency != "" && a.CtrlMgrNodeMonitorGracePeriod != "" {
		nodeStatusUpdateFrequency, _ := time.ParseDuration(a.NodeStatusUpdateFrequency)
		ctrlMgrNodeMonitorGracePeriod, _ := time.ParseDuration(a.CtrlMgrNodeMonitorGracePeriod)
		kubeletRetries := ctrlMgrNodeMonitorGracePeriod.Seconds() / nodeStatusUpdateFrequency.Seconds()
		if kubeletRetries < minKubeletRetries {
			return fmt.Errorf("acs-engine requires that ctrlMgrNodeMonitorGracePeriod(%f)s be larger than nodeStatusUpdateFrequency(%f)s by at least a factor of %d; ", ctrlMgrNodeMonitorGracePeriod.Seconds(), nodeStatusUpdateFrequency.Seconds(), minKubeletRetries)
		}
	}

	if a.CtrlMgrPodEvictionTimeout != "" {
		_, err := time.ParseDuration(a.CtrlMgrPodEvictionTimeout)
		if err != nil {
			return fmt.Errorf("OrchestratorProfile.KubernetesConfig.CtrlMgrPodEvictionTimeout '%s' is not a valid duration", a.CtrlMgrPodEvictionTimeout)
		}
	}

	if a.CtrlMgrRouteReconciliationPeriod != "" {
		_, err := time.ParseDuration(a.CtrlMgrRouteReconciliationPeriod)
		if err != nil {
			return fmt.Errorf("OrchestratorProfile.KubernetesConfig.CtrlMgrRouteReconciliationPeriod '%s' is not a valid duration", a.CtrlMgrRouteReconciliationPeriod)
		}
	}

	if a.CloudProviderBackoff {
		if !backoffEnabledVersions[k8sVersion] {
			return fmt.Errorf("cloudprovider backoff functionality not available in kubernetes version %s", k8sVersion)
		}
	}

	if a.CloudProviderRateLimit {
		if !ratelimitEnabledVersions[k8sVersion] {
			return fmt.Errorf("cloudprovider rate limiting functionality not available in kubernetes version %s", k8sVersion)
		}
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

func validateNameEmpty(name string, label string) error {
	if name != "" {
		return fmt.Errorf("%s must be an empty value", label)
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

func validateUniqueProfileNames(profiles []*AgentPoolProfile) error {
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
			return fmt.Errorf("Unknown storageProfile '%s'. Specify either %s or %s", storageProfile, StorageAccount, ManagedDisks)
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
			return fmt.Errorf("Multiple VNET Subnet configurations specified.  The master profile and each agent pool profile must all specify a custom VNET Subnet, or none at all")
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
				return errors.New("Multiple VNETS specified.  The master profile and each agent pool must reference the same VNET (but it is ok to reference different subnets on that VNET)")
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

package vlabs

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/Azure/acs-engine/pkg/api/common"
	"github.com/satori/uuid"
	validator "gopkg.in/go-playground/validator.v9"
)

var (
	validate        *validator.Validate
	keyvaultIDRegex *regexp.Regexp
	labelValueRegex *regexp.Regexp
	labelKeyRegex   *regexp.Regexp
	// Any version has to be mirrored in https://acs-mirror.azureedge.net/github-coreos/etcd-v[Version]-linux-amd64.tar.gz
	etcdValidVersions = [...]string{"2.2.5", "2.3.0", "2.3.1", "2.3.2", "2.3.3", "2.3.4", "2.3.5", "2.3.6", "2.3.7", "2.3.8",
		"3.0.0", "3.0.1", "3.0.2", "3.0.3", "3.0.4", "3.0.5", "3.0.6", "3.0.7", "3.0.8", "3.0.9", "3.0.10", "3.0.11", "3.0.12", "3.0.13", "3.0.14", "3.0.15", "3.0.16", "3.0.17",
		"3.1.0", "3.1.1", "3.1.2", "3.1.2", "3.1.3", "3.1.4", "3.1.5", "3.1.6", "3.1.7", "3.1.8", "3.1.9", "3.1.10",
		"3.2.0", "3.2.1", "3.2.2", "3.2.3", "3.2.4", "3.2.5", "3.2.6", "3.2.7", "3.2.8", "3.2.9", "3.2.11"}
)

const (
	labelKeyPrefixMaxLength = 253
	labelValueFormat        = "^([A-Za-z0-9][-A-Za-z0-9_.]{0,61})?[A-Za-z0-9]$"
	labelKeyFormat          = "^(([a-zA-Z0-9-]+[.])*[a-zA-Z0-9-]+[/])?([A-Za-z0-9][-A-Za-z0-9_.]{0,61})?[A-Za-z0-9]$"
)

func init() {
	validate = validator.New()
	keyvaultIDRegex = regexp.MustCompile(`^/subscriptions/\S+/resourceGroups/\S+/providers/Microsoft.KeyVault/vaults/[^/\s]+$`)
	labelValueRegex = regexp.MustCompile(labelValueFormat)
	labelKeyRegex = regexp.MustCompile(labelKeyFormat)
}

func isValidEtcdVersion(etcdVersion string) error {
	// "" is a valid etcdVersion that maps to DefaultEtcdVersion
	if etcdVersion == "" {
		return nil
	}
	for _, ver := range etcdValidVersions {
		if ver == etcdVersion {
			return nil
		}
	}
	return fmt.Errorf("Invalid etcd version(%s), valid versions are%s", etcdVersion, etcdValidVersions)
}

// Validate implements APIObject
func (o *OrchestratorProfile) Validate(isUpdate bool) error {
	// Don't need to call validate.Struct(o)
	// It is handled by Properties.Validate()
	// On updates we only need to make sure there is a supported patch version for the minor version
	if !isUpdate {
		switch o.OrchestratorType {
		case DCOS:
			version := common.RationalizeReleaseAndVersion(
				o.OrchestratorType,
				o.OrchestratorRelease,
				o.OrchestratorVersion)
			if version == "" {
				return fmt.Errorf("OrchestratorProfile is not able to be rationalized, check supported Release or Version")
			}
		case Swarm:
		case SwarmMode:
		case Kubernetes:
			version := common.RationalizeReleaseAndVersion(
				o.OrchestratorType,
				o.OrchestratorRelease,
				o.OrchestratorVersion)
			if version == "" {
				return fmt.Errorf("OrchestratorProfile is not able to be rationalized, check supported Release or Version")
			}

			if o.KubernetesConfig != nil {
				err := o.KubernetesConfig.Validate(version)
				if err != nil {
					return err
				}
				if o.KubernetesConfig.EnableAggregatedAPIs {
					if o.OrchestratorVersion == common.KubernetesVersion1Dot5Dot7 ||
						o.OrchestratorVersion == common.KubernetesVersion1Dot5Dot8 ||
						o.OrchestratorVersion == common.KubernetesVersion1Dot6Dot6 ||
						o.OrchestratorVersion == common.KubernetesVersion1Dot6Dot9 ||
						o.OrchestratorVersion == common.KubernetesVersion1Dot6Dot11 {
						return fmt.Errorf("enableAggregatedAPIs is only available in Kubernetes version %s or greater; unable to validate for Kubernetes version %s",
							"1.7.0", o.OrchestratorVersion)
					}

					if o.KubernetesConfig.EnableRbac != nil {
						if !*o.KubernetesConfig.EnableRbac {
							return fmt.Errorf("enableAggregatedAPIs requires the enableRbac feature as a prerequisite")
						}
					}
				}
			}

		default:
			return fmt.Errorf("OrchestratorProfile has unknown orchestrator: %s", o.OrchestratorType)
		}
	} else {
		switch o.OrchestratorType {
		case DCOS, Kubernetes:

			version := common.RationalizeReleaseAndVersion(
				o.OrchestratorType,
				o.OrchestratorRelease,
				o.OrchestratorVersion)
			if version == "" {
				patchVersion := common.GetValidPatchVersion(o.OrchestratorType, o.OrchestratorVersion)
				// if there isn't a supported patch version for this version fail
				if patchVersion == "" {
					return fmt.Errorf("OrchestratorProfile is not able to be rationalized, check supported Release or Version")
				}
			}

		}
	}

	if o.OrchestratorType != Kubernetes && o.KubernetesConfig != nil {
		return fmt.Errorf("KubernetesConfig can be specified only when OrchestratorType is Kubernetes")
	}

	if o.OrchestratorType != DCOS && o.DcosConfig != nil && (*o.DcosConfig != DcosConfig{}) {
		return fmt.Errorf("DcosConfig can be specified only when OrchestratorType is DCOS")
	}
	return nil
}

// Validate implements APIObject
func (m *MasterProfile) Validate() error {
	if e := validateDNSName(m.DNSPrefix); e != nil {
		return e
	}
	return nil
}

// Validate implements APIObject
func (a *AgentPoolProfile) Validate(orchestratorType string) error {
	// Don't need to call validate.Struct(a)
	// It is handled by Properties.Validate()
	if e := validatePoolName(a.Name); e != nil {
		return e
	}

	// for Kubernetes, we don't support AgentPoolProfile.DNSPrefix
	if orchestratorType == Kubernetes {
		if e := validate.Var(a.DNSPrefix, "len=0"); e != nil {
			return fmt.Errorf("AgentPoolProfile.DNSPrefix must be empty for Kubernetes")
		}
		if e := validate.Var(a.Ports, "len=0"); e != nil {
			return fmt.Errorf("AgentPoolProfile.Ports must be empty for Kubernetes")
		}
	}

	if a.DNSPrefix != "" {
		if e := validateDNSName(a.DNSPrefix); e != nil {
			return e
		}
		if len(a.Ports) > 0 {
			if e := validateUniquePorts(a.Ports, a.Name); e != nil {
				return e
			}
		} else {
			a.Ports = []int{80, 443, 8080}
		}
	} else {
		if e := validate.Var(a.Ports, "len=0"); e != nil {
			return fmt.Errorf("AgentPoolProfile.Ports must be empty when AgentPoolProfile.DNSPrefix is empty for Orchestrator: %s", string(orchestratorType))
		}
	}

	if len(a.DiskSizesGB) > 0 {
		if e := validate.Var(a.StorageProfile, "eq=StorageAccount|eq=ManagedDisks"); e != nil {
			return fmt.Errorf("property 'StorageProfile' must be set to either '%s' or '%s' when attaching disks", StorageAccount, ManagedDisks)
		}
		if e := validate.Var(a.AvailabilityProfile, "eq=VirtualMachineScaleSets|eq=AvailabilitySet"); e != nil {
			return fmt.Errorf("property 'AvailabilityProfile' must be set to either '%s' or '%s' when attaching disks", VirtualMachineScaleSets, AvailabilitySet)
		}
		if a.StorageProfile == StorageAccount && (a.AvailabilityProfile == VirtualMachineScaleSets) {
			return fmt.Errorf("VirtualMachineScaleSets does not support storage account attached disks.  Instead specify 'StorageAccount': '%s' or specify AvailabilityProfile '%s'", ManagedDisks, AvailabilitySet)
		}
	}
	if len(a.Ports) == 0 && len(a.DNSPrefix) > 0 {
		return fmt.Errorf("AgentPoolProfile.Ports must be non empty when AgentPoolProfile.DNSPrefix is specified")
	}
	return nil
}

// Validate implements APIObject
func (o *OrchestratorVersionProfile) Validate() error {
	// The only difference compared with OrchestratorProfile.Validate is
	// Here we use strings.EqualFold, the other just string comparison.
	// Rationalize orchestrator type should be done from versioned to unversioned
	// I will go ahead to simplify this
	return o.OrchestratorProfile.Validate(false)
}

// ValidateForUpgrade validates upgrade input data
func (o *OrchestratorProfile) ValidateForUpgrade() error {
	switch o.OrchestratorType {
	case DCOS, SwarmMode, Swarm:
		return fmt.Errorf("Upgrade is not supported for orchestrator %s", o.OrchestratorType)
	case Kubernetes:
		switch o.OrchestratorVersion {
		case common.KubernetesVersion1Dot6Dot13:
		case common.KubernetesVersion1Dot7Dot10:
		default:
			return fmt.Errorf("Upgrade to Kubernetes version %s is not supported", o.OrchestratorVersion)
		}
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
	// Don't need to call validate.Struct(l)
	// It is handled by Properties.Validate()
	if e := validate.Var(l.SSH.PublicKeys[0].KeyData, "required"); e != nil {
		return fmt.Errorf("KeyData in LinuxProfile.SSH.PublicKeys cannot be empty string")
	}
	if e := validateKeyVaultSecrets(l.Secrets, false); e != nil {
		return e
	}
	return nil
}

func handleValidationErrors(e validator.ValidationErrors) error {
	// Override any version specific validation error message

	// common.HandleValidationErrors if the validation error message is general
	return common.HandleValidationErrors(e)
}

// Validate implements APIObject
func (w *WindowsProfile) Validate() error {
	if e := validate.Var(w.AdminUsername, "required"); e != nil {
		return fmt.Errorf("WindowsProfile.AdminUsername is required, when agent pool specifies windows")
	}
	if e := validate.Var(w.AdminPassword, "required"); e != nil {
		return fmt.Errorf("WindowsProfile.AdminPassword is required, when agent pool specifies windows")
	}
	if e := validateKeyVaultSecrets(w.Secrets, true); e != nil {
		return e
	}
	return nil
}

// Validate implements APIObject
func (profile *AADProfile) Validate() error {
	if _, err := uuid.FromString(profile.ClientAppID); err != nil {
		return fmt.Errorf("clientAppID '%v' is invalid", profile.ClientAppID)
	}
	if _, err := uuid.FromString(profile.ServerAppID); err != nil {
		return fmt.Errorf("serverAppID '%v' is invalid", profile.ServerAppID)
	}
	if len(profile.TenantID) > 0 {
		if _, err := uuid.FromString(profile.TenantID); err != nil {
			return fmt.Errorf("tenantID '%v' is invalid", profile.TenantID)
		}
	}
	return nil
}

// Validate implements APIObject
func (a *Properties) Validate(isUpdate bool) error {
	if e := validate.Struct(a); e != nil {
		return handleValidationErrors(e.(validator.ValidationErrors))
	}
	if e := a.OrchestratorProfile.Validate(isUpdate); e != nil {
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
			if a.ServicePrincipalProfile == nil {
				return fmt.Errorf("ServicePrincipalProfile must be specified with Orchestrator %s", a.OrchestratorProfile.OrchestratorType)
			}
			if e := validate.Var(a.ServicePrincipalProfile.ClientID, "required"); e != nil {
				return fmt.Errorf("the service principal client ID must be specified with Orchestrator %s", a.OrchestratorProfile.OrchestratorType)
			}
			if (len(a.ServicePrincipalProfile.Secret) == 0 && a.ServicePrincipalProfile.KeyvaultSecretRef == nil) ||
				(len(a.ServicePrincipalProfile.Secret) != 0 && a.ServicePrincipalProfile.KeyvaultSecretRef != nil) {
				return fmt.Errorf("either the service principal client secret or keyvault secret reference must be specified with Orchestrator %s", a.OrchestratorProfile.OrchestratorType)
			}

			if a.ServicePrincipalProfile.KeyvaultSecretRef != nil {
				if e := validate.Var(a.ServicePrincipalProfile.KeyvaultSecretRef.VaultID, "required"); e != nil {
					return fmt.Errorf("the Keyvault ID must be specified for the Service Principle with Orchestrator %s", a.OrchestratorProfile.OrchestratorType)
				}
				if e := validate.Var(a.ServicePrincipalProfile.KeyvaultSecretRef.SecretName, "required"); e != nil {
					return fmt.Errorf("the Keyvault Secret must be specified for the Service Principle with Orchestrator %s", a.OrchestratorProfile.OrchestratorType)
				}
				if !keyvaultIDRegex.MatchString(a.ServicePrincipalProfile.KeyvaultSecretRef.VaultID) {
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
				for k, v := range agentPoolProfile.CustomNodeLabels {
					if e := validateKubernetesLabelKey(k); e != nil {
						return e
					}
					if e := validateKubernetesLabelValue(v); e != nil {
						return e
					}
				}
			default:
				return fmt.Errorf("Agent Type attributes are only supported for DCOS and Kubernetes")
			}
		}
		if a.OrchestratorProfile.OrchestratorType == Kubernetes && (agentPoolProfile.AvailabilityProfile == VirtualMachineScaleSets || len(agentPoolProfile.AvailabilityProfile) == 0) {
			return fmt.Errorf("VirtualMachineScaleSets are not supported with Kubernetes since Kubernetes requires the ability to attach/detach disks.  To fix specify \"AvailabilityProfile\":\"%s\"", AvailabilitySet)
		}
		if agentPoolProfile.OSType == Windows {
			if e := validate.Var(a.WindowsProfile, "required"); e != nil {
				return fmt.Errorf("WindowsProfile must not be empty since agent pool '%s' specifies windows", agentPoolProfile.Name)
			}
			switch a.OrchestratorProfile.OrchestratorType {
			case DCOS:
			case Swarm:
			case SwarmMode:
			case Kubernetes:
				version := common.RationalizeReleaseAndVersion(
					a.OrchestratorProfile.OrchestratorType,
					a.OrchestratorProfile.OrchestratorRelease,
					a.OrchestratorProfile.OrchestratorVersion)
				if version == "" {
					return fmt.Errorf("OrchestratorProfile is not able to be rationalized, check supported Release or Version")
				}
				if _, ok := common.AllKubernetesWindowsSupportedVersions[version]; !ok {
					return fmt.Errorf("Orchestrator %s version %s does not support Windows", a.OrchestratorProfile.OrchestratorType, version)
				}
			default:
				return fmt.Errorf("Orchestrator %s does not support Windows", a.OrchestratorProfile.OrchestratorType)
			}
			if e := a.WindowsProfile.Validate(); e != nil {
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

	if a.AADProfile != nil {
		if a.OrchestratorProfile.OrchestratorType != Kubernetes {
			return fmt.Errorf("'aadProfile' is only supported by orchestrator '%v'", Kubernetes)
		}
		if e := a.AADProfile.Validate(); e != nil {
			return e
		}
	}

	for _, extension := range a.ExtensionProfiles {
		if extension.ExtensionParametersKeyVaultRef != nil {
			if e := validate.Var(extension.ExtensionParametersKeyVaultRef.VaultID, "required"); e != nil {
				return fmt.Errorf("the Keyvault ID must be specified for Extension %s", extension.Name)
			}
			if e := validate.Var(extension.ExtensionParametersKeyVaultRef.SecretName, "required"); e != nil {
				return fmt.Errorf("the Keyvault Secret must be specified for Extension %s", extension.Name)
			}
			if !keyvaultIDRegex.MatchString(extension.ExtensionParametersKeyVaultRef.VaultID) {
				return fmt.Errorf("Extension %s's keyvault secret reference is of incorrect format", extension.Name)
			}
		}
	}

	return nil
}

// Validate validates the KubernetesConfig.
func (a *KubernetesConfig) Validate(k8sVersion string) error {
	// number of minimum retries allowed for kubelet to post node status
	const minKubeletRetries = 4
	// k8s versions that have cloudprovider backoff enabled
	var backoffEnabledVersions = map[string]bool{
		common.KubernetesVersion1Dot8Dot0:  true,
		common.KubernetesVersion1Dot8Dot1:  true,
		common.KubernetesVersion1Dot8Dot2:  true,
		common.KubernetesVersion1Dot8Dot4:  true,
		common.KubernetesVersion1Dot7Dot0:  true,
		common.KubernetesVersion1Dot7Dot1:  true,
		common.KubernetesVersion1Dot7Dot2:  true,
		common.KubernetesVersion1Dot7Dot4:  true,
		common.KubernetesVersion1Dot7Dot5:  true,
		common.KubernetesVersion1Dot7Dot7:  true,
		common.KubernetesVersion1Dot7Dot9:  true,
		common.KubernetesVersion1Dot7Dot10: true,
		common.KubernetesVersion1Dot6Dot6:  true,
		common.KubernetesVersion1Dot6Dot9:  true,
		common.KubernetesVersion1Dot6Dot11: true,
		common.KubernetesVersion1Dot6Dot12: true,
		common.KubernetesVersion1Dot6Dot13: true,
	}
	// k8s versions that have cloudprovider rate limiting enabled (currently identical with backoff enabled versions)
	ratelimitEnabledVersions := backoffEnabledVersions

	if a.ClusterSubnet != "" {
		_, subnet, err := net.ParseCIDR(a.ClusterSubnet)
		if err != nil {
			return fmt.Errorf("OrchestratorProfile.KubernetesConfig.ClusterSubnet '%s' is an invalid subnet", a.ClusterSubnet)
		}

		if a.NetworkPolicy == "azure" {
			ones, bits := subnet.Mask.Size()
			if bits-ones <= 8 {
				return fmt.Errorf("OrchestratorProfile.KubernetesConfig.ClusterSubnet '%s' must reserve at least 9 bits for nodes", a.ClusterSubnet)
			}
		}
	}

	if a.DockerBridgeSubnet != "" {
		_, _, err := net.ParseCIDR(a.DockerBridgeSubnet)
		if err != nil {
			return fmt.Errorf("OrchestratorProfile.KubernetesConfig.DockerBridgeSubnet '%s' is an invalid subnet", a.DockerBridgeSubnet)
		}
	}

	if a.MaxPods != 0 {
		if a.MaxPods < KubernetesMinMaxPods {
			return fmt.Errorf("OrchestratorProfile.KubernetesConfig.MaxPods '%v' must be at least %v", a.MaxPods, KubernetesMinMaxPods)
		}
	}

	if a.KubeletConfig != nil {
		if _, ok := a.KubeletConfig["--node-status-update-frequency"]; ok {
			val := a.KubeletConfig["--node-status-update-frequency"]
			_, err := time.ParseDuration(val)
			if err != nil {
				return fmt.Errorf("--node-status-update-frequency '%s' is not a valid duration", val)
			}
			if a.CtrlMgrNodeMonitorGracePeriod == "" {
				return fmt.Errorf("--node-status-update-frequency was set to '%s' but OrchestratorProfile.KubernetesConfig.CtrlMgrNodeMonitorGracePeriod was not set", val)
			}
		}
	}

	if a.CtrlMgrNodeMonitorGracePeriod != "" {
		_, err := time.ParseDuration(a.CtrlMgrNodeMonitorGracePeriod)
		if err != nil {
			return fmt.Errorf("OrchestratorProfile.KubernetesConfig.CtrlMgrNodeMonitorGracePeriod '%s' is not a valid duration", a.CtrlMgrNodeMonitorGracePeriod)
		}
		if a.KubeletConfig != nil {
			if _, ok := a.KubeletConfig["--node-status-update-frequency"]; !ok {
				return fmt.Errorf("OrchestratorProfile.KubernetesConfig.CtrlMgrNodeMonitorGracePeriod was set to '%s' but kubelet config --node-status-update-frequency was not set", a.CtrlMgrNodeMonitorGracePeriod)
			}
		}
	}

	if a.KubeletConfig != nil {
		if _, ok := a.KubeletConfig["--node-status-update-frequency"]; ok {
			if a.CtrlMgrNodeMonitorGracePeriod != "" {
				nodeStatusUpdateFrequency, _ := time.ParseDuration(a.KubeletConfig["--node-status-update-frequency"])
				ctrlMgrNodeMonitorGracePeriod, _ := time.ParseDuration(a.CtrlMgrNodeMonitorGracePeriod)
				kubeletRetries := ctrlMgrNodeMonitorGracePeriod.Seconds() / nodeStatusUpdateFrequency.Seconds()
				if kubeletRetries < minKubeletRetries {
					return fmt.Errorf("acs-engine requires that ctrlMgrNodeMonitorGracePeriod(%f)s be larger than nodeStatusUpdateFrequency(%f)s by at least a factor of %d; ", ctrlMgrNodeMonitorGracePeriod.Seconds(), nodeStatusUpdateFrequency.Seconds(), minKubeletRetries)
				}
			}
		}
		if _, ok := a.KubeletConfig["--non-masquerade-cidr"]; ok {
			if _, _, err := net.ParseCIDR(a.KubeletConfig["--non-masquerade-cidr"]); err != nil {
				return fmt.Errorf("--non-masquerade-cidr kubelet config '%s' is an invalid CIDR string", a.KubeletConfig["--non-masquerade-cidr"])
			}
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

	if a.DNSServiceIP != "" || a.ServiceCidr != "" {
		if a.DNSServiceIP == "" {
			return errors.New("OrchestratorProfile.KubernetesConfig.ServiceCidr must be specified when DNSServiceIP is")
		}
		if a.ServiceCidr == "" {
			return errors.New("OrchestratorProfile.KubernetesConfig.DNSServiceIP must be specified when ServiceCidr is")
		}

		dnsIP := net.ParseIP(a.DNSServiceIP)
		if dnsIP == nil {
			return fmt.Errorf("OrchestratorProfile.KubernetesConfig.DNSServiceIP '%s' is an invalid IP address", a.DNSServiceIP)
		}

		_, serviceCidr, err := net.ParseCIDR(a.ServiceCidr)
		if err != nil {
			return fmt.Errorf("OrchestratorProfile.KubernetesConfig.ServiceCidr '%s' is an invalid CIDR subnet", a.ServiceCidr)
		}

		// Finally validate that the DNS ip is within the subnet
		if !serviceCidr.Contains(dnsIP) {
			return fmt.Errorf("OrchestratorProfile.KubernetesConfig.DNSServiceIP '%s' is not within the ServiceCidr '%s'", a.DNSServiceIP, a.ServiceCidr)
		}

		// and that the DNS IP is _not_ the subnet broadcast address
		broadcast := common.IP4BroadcastAddress(serviceCidr)
		if dnsIP.Equal(broadcast) {
			return fmt.Errorf("OrchestratorProfile.KubernetesConfig.DNSServiceIP '%s' cannot be the broadcast address of ServiceCidr '%s'", a.DNSServiceIP, a.ServiceCidr)
		}

		// and that the DNS IP is _not_ the first IP in the service subnet
		firstServiceIP := common.CidrFirstIP(serviceCidr.IP)
		if firstServiceIP.Equal(dnsIP) {
			return fmt.Errorf("OrchestratorProfile.KubernetesConfig.DNSServiceIP '%s' cannot be the first IP of ServiceCidr '%s'", a.DNSServiceIP, a.ServiceCidr)
		}
	}

	// Validate that we have a valid etcd version
	if e := isValidEtcdVersion(a.EtcdVersion); e != nil {
		return e
	}

	var ccmEnabledVersions = map[string]bool{
		common.KubernetesVersion1Dot8Dot0: true,
		common.KubernetesVersion1Dot8Dot1: true,
		common.KubernetesVersion1Dot8Dot2: true,
		common.KubernetesVersion1Dot8Dot4: true,
	}

	if a.UseCloudControllerManager != nil && *a.UseCloudControllerManager || a.CustomCcmImage != "" {
		if !ccmEnabledVersions[k8sVersion] {
			return fmt.Errorf("OrchestratorProfile.KubernetesConfig.UseCloudControllerManager and OrchestratorProfile.KubernetesConfig.CustomCcmImage not available in kubernetes version %s", k8sVersion)
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

func validateKubernetesLabelValue(v string) error {
	if !(len(v) == 0) && !labelValueRegex.MatchString(v) {
		return fmt.Errorf("Label value '%s' is invalid. Valid label values must be 63 characters or less and must be empty or begin and end with an alphanumeric character ([a-z0-9A-Z]) with dashes (-), underscores (_), dots (.), and alphanumerics between", v)
	}
	return nil
}

func validateKubernetesLabelKey(k string) error {
	if !labelKeyRegex.MatchString(k) {
		return fmt.Errorf("Label key '%s' is invalid. Valid label keys have two segments: an optional prefix and name, separated by a slash (/). The name segment is required and must be 63 characters or less, beginning and ending with an alphanumeric character ([a-z0-9A-Z]) with dashes (-), underscores (_), dots (.), and alphanumerics between. The prefix is optional. If specified, the prefix must be a DNS subdomain: a series of DNS labels separated by dots (.), not longer than 253 characters in total, followed by a slash (/)", k)
	}
	prefix := strings.Split(k, "/")
	if len(prefix) != 1 && len(prefix[0]) > labelKeyPrefixMaxLength {
		return fmt.Errorf("Label key prefix '%s' is invalid. If specified, the prefix must be no longer than 253 characters in total", k)
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

		if a.MasterProfile.VnetCidr != "" {
			_, _, err := net.ParseCIDR(a.MasterProfile.VnetCidr)
			if err != nil {
				return fmt.Errorf("MasterProfile.VnetCidr '%s' contains invalid cidr notation", a.MasterProfile.VnetCidr)
			}
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

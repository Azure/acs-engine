package v20180331

import "fmt"

// ErrorInvalidNetworkProfile error
var ErrorInvalidNetworkProfile = fmt.Errorf("ServiceCidr, DNSServiceIP, DockerBridgeCidr should all be empty or neither should be empty")

// ErrorInvalidNetworkPlugin error
var ErrorInvalidNetworkPlugin = fmt.Errorf("Network plugin should be either Azure or Kubenet")

// ErrorInvalidServiceCidr error
var ErrorInvalidServiceCidr = fmt.Errorf("ServiceCidr is not a valid CIDR")

// ErrorInvalidDNSServiceIP error
var ErrorInvalidDNSServiceIP = fmt.Errorf("DNSServiceIP is not a valid IP address")

// ErrorInvalidDockerBridgeCidr error
var ErrorInvalidDockerBridgeCidr = fmt.Errorf("DockerBridgeCidr is not a valid IP address")

// ErrorDNSServiceIPNotInServiceCidr error
var ErrorDNSServiceIPNotInServiceCidr = fmt.Errorf("DNSServiceIP is not within ServiceCidr")

// ErrorDNSServiceIPAlreadyUsed error
var ErrorDNSServiceIPAlreadyUsed = fmt.Errorf("DNSServiceIP can not be the first IP address in ServiceCidr")

// ErrorAtLeastAgentPoolNoSubnet error
var ErrorAtLeastAgentPoolNoSubnet = fmt.Errorf("At least one agent pool does not have subnet defined")

// ErrorInvalidMaxPods error
var ErrorInvalidMaxPods = fmt.Errorf("Max pods per node needs to be at least 5")

// ErrorParsingSubnetID error
var ErrorParsingSubnetID = fmt.Errorf("Failed to parse VnetSubnetID")

// ErrorSubscriptionNotMatch error
var ErrorSubscriptionNotMatch = fmt.Errorf("Subscription for subnet does not match with other subnet")

// ErrorResourceGroupNotMatch error
var ErrorResourceGroupNotMatch = fmt.Errorf("ResourceGroup for subnet does not match with other subnet")

// ErrorVnetNotMatch error
var ErrorVnetNotMatch = fmt.Errorf("Vnet for subnet does not match with other subnet")

// ErrorRBACNotEnabledForAAD error
var ErrorRBACNotEnabledForAAD = fmt.Errorf("RBAC must be enabled for AAD to be enabled")

// ErrorAADServerAppIDNotSet error
var ErrorAADServerAppIDNotSet = fmt.Errorf("ServerAppID in AADProfile cannot be empty string")

// ErrorAADServerAppSecretNotSet error
var ErrorAADServerAppSecretNotSet = fmt.Errorf("ServerAppSecret in AADProfile cannot be empty string")

// ErrorAADClientAppIDNotSet error
var ErrorAADClientAppIDNotSet = fmt.Errorf("ClientAppID in AADProfile cannot be empty string")

// ErrorAADTenantIDNotSet error
var ErrorAADTenantIDNotSet = fmt.Errorf("TenantID in AADProfile cannot be empty string")

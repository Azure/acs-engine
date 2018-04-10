package v20180331

import "fmt"

// ErrorNilNetworkProfile error
var ErrorNilNetworkProfile = fmt.Errorf("Network profile can not be nil")

// ErrorNilAgentPoolProfile error
var ErrorNilAgentPoolProfile = fmt.Errorf("Agent pool profile can not be nil")

// ErrorInvalidNetworkPlugin error
var ErrorInvalidNetworkPlugin = fmt.Errorf("Network plugin should be either \"azure\" or \"kubenet\"")

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

// ErrorAgentPoolNoSubnet error
var ErrorAgentPoolNoSubnet = fmt.Errorf("Agent pool does not have subnet defined")

// ErrorKubenetNoCustomization error
var ErrorKubenetNoCustomization = fmt.Errorf("Kubenet does not support ServiceCidr or DNSServiceIP or DockerBridgeCidr customization")

// ErrorParsingSubnetID error
var ErrorParsingSubnetID = fmt.Errorf("Failed to parse VnetSubnetID")

// ErrorSubscriptionNotMatch error
var ErrorSubscriptionNotMatch = fmt.Errorf("Subscription for subnet does not match with other subnet")

// ErrorResourceGroupNotMatch error
var ErrorResourceGroupNotMatch = fmt.Errorf("ResourceGroup for subnet does not match with other subnet")

// ErrorVnetNotMatch error
var ErrorVnetNotMatch = fmt.Errorf("Vnet for subnet does not match with other subnet")

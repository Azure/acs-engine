package util

import "fmt"

const (
	// AzureProdFQDNFormat specifies the format for a prod dns name
	AzureProdFQDNFormat = "%s.%s.cloudapp.azure.com"
)

// AzureLocations provides all azure regions in prod.
// Related powershell to refresh this list:
//   Get-AzureRmLocation | Select-Object -Property Location
var AzureLocations = []string{
	"eastasia",
	"southeastasia",
	"centralus",
	"eastus",
	"eastus2",
	"westus",
	"northcentralus",
	"southcentralus",
	"northeurope",
	"westeurope",
	"japanwest",
	"japaneast",
	"brazilsouth",
	"australiaeast",
	"australiasoutheast",
	"southindia",
	"centralindia",
	"westindia",
	"canadacentral",
	"canadaeast",
	"uksouth",
	"ukwest",
	"westcentralus",
	"westus2",
}

// FormatAzureProdFQDNs constructs all possible Azure prod fqdn
func FormatAzureProdFQDNs(fqdnPrefix string) []string {
	var fqdns []string
	for _, location := range AzureLocations {
		fqdns = append(fqdns, FormatAzureProdFQDN(fqdnPrefix, location))
	}
	return fqdns
}

// FormatAzureProdFQDN constructs an Azure prod fqdn
func FormatAzureProdFQDN(fqdnPrefix string, location string) string {
	return fmt.Sprintf(AzureProdFQDNFormat, fqdnPrefix, location)
}

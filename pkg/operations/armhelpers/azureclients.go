package armhelpers

import (
	"github.com/Azure/azure-sdk-for-go/arm/compute"
	"github.com/Azure/azure-sdk-for-go/arm/resources/resources"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/Azure/go-autorest/autorest/azure"
)

// AzureClients contains Azure clients
type AzureClients struct {
	SubscriptionID   string
	AzureEnvironment azure.Environment
	TenantID         string

	GroupsClient *resources.GroupsClient
	VMClient     *compute.VirtualMachinesClient
}

// NewAzureClients creates various Azure clients
func NewAzureClients(token *adal.ServicePrincipalToken, subscriptionID string) AzureClients {
	azureClients := AzureClients{}
	azureClients.SubscriptionID = subscriptionID

	gc := resources.NewGroupsClient(azureClients.SubscriptionID)
	gc.Authorizer = autorest.NewBearerAuthorizer(token)
	azureClients.GroupsClient = &gc

	vmc := compute.NewVirtualMachinesClient(azureClients.SubscriptionID)
	vmc.Authorizer = autorest.NewBearerAuthorizer(token)
	azureClients.VMClient = &vmc

	return azureClients
}

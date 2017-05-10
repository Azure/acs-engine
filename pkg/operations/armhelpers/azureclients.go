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

// Create method creates various Azure clients
func (ac *AzureClients) Create(token *adal.ServicePrincipalToken) (*AzureClients, error) {
	gc := resources.NewGroupsClient(ac.SubscriptionID)
	gc.Authorizer = autorest.NewBearerAuthorizer(token)
	ac.GroupsClient = &gc

	vmc := compute.NewVirtualMachinesClient(ac.SubscriptionID)
	vmc.Authorizer = autorest.NewBearerAuthorizer(token)
	ac.VMClient = &vmc

	return ac, nil
}

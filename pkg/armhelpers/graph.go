package armhelpers

import (
	"fmt"
	"regexp"
	"time"

	"github.com/Azure/azure-sdk-for-go/arm/authorization"
	"github.com/Azure/azure-sdk-for-go/arm/graphrbac"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

const (
	// AADContributorRoleID is the role id that exists in every subscription for 'Contributor'
	AADContributorRoleID = "b24988ac-6180-42a0-ab88-20f7382dd24c"
	// AADRoleReferenceTemplate is a template for a roleDefinitionId
	AADRoleReferenceTemplate = "/subscriptions/%s/providers/Microsoft.Authorization/roleDefinitions/%s"
	// AADRoleResourceGroupScopeTemplate is a template for a roleDefinition scope
	AADRoleResourceGroupScopeTemplate = "/subscriptions/%s/resourceGroups/%s"
)

// CreateGraphApplication creates an application via the graphrbac client
func (az *AzureClient) CreateGraphApplication(applicationCreateParameters graphrbac.ApplicationCreateParameters) (graphrbac.Application, error) {
	return az.applicationsClient.Create(applicationCreateParameters)
}

// CreateGraphPrincipal creates a service principal via the graphrbac client
func (az *AzureClient) CreateGraphPrincipal(servicePrincipalCreateParameters graphrbac.ServicePrincipalCreateParameters) (graphrbac.ServicePrincipal, error) {
	return az.servicePrincipalsClient.Create(servicePrincipalCreateParameters)
}

// CreateRoleAssignment creates a role assignment via the authorization client
func (az *AzureClient) CreateRoleAssignment(scope string, roleAssignmentName string, parameters authorization.RoleAssignmentCreateParameters) (authorization.RoleAssignment, error) {
	return az.authorizationClient.Create(scope, roleAssignmentName, parameters)
}

// CreateApp is a simpler method for creating an application
func (az *AzureClient) CreateApp(appName, appURL string) (applicationID, servicePrincipalObjectID, servicePrincipalClientSecret string, err error) {
	notBefore := time.Now()
	notAfter := time.Now().Add(10000 * 24 * time.Hour)

	startDate := date.Time{Time: notBefore}
	endDate := date.Time{Time: notAfter}

	servicePrincipalClientSecret = uuid.NewV4().String()

	log.Debugf("ad: creating application with name=%q identifierURL=%q", appName, appURL)
	applicationReq := graphrbac.ApplicationCreateParameters{
		AvailableToOtherTenants: to.BoolPtr(false),
		DisplayName:             to.StringPtr(appName),
		Homepage:                to.StringPtr(appURL),
		IdentifierUris:          to.StringSlicePtr([]string{appURL}),
		PasswordCredentials: &[]graphrbac.PasswordCredential{
			{
				KeyID:     to.StringPtr(uuid.NewV4().String()),
				StartDate: &startDate,
				EndDate:   &endDate,
				Value:     to.StringPtr(servicePrincipalClientSecret),
			},
		},
	}
	applicationResp, err := az.CreateGraphApplication(applicationReq)
	if err != nil {
		return "", "", "", err
	}
	applicationID = to.String(applicationResp.AppID)

	log.Debugf("ad: creating servicePrincipal for applicationID: %q", applicationID)

	servicePrincipalReq := graphrbac.ServicePrincipalCreateParameters{
		AppID:          applicationResp.AppID,
		AccountEnabled: to.BoolPtr(true),
	}
	servicePrincipalResp, err := az.servicePrincipalsClient.Create(servicePrincipalReq)
	if err != nil {
		return "", "", "", err
	}

	servicePrincipalObjectID = to.String(servicePrincipalResp.ObjectID)

	return applicationID, servicePrincipalObjectID, servicePrincipalClientSecret, nil
}

// CreateRoleAssignmentSimple is a wrapper around RoleAssignmentsClient.Create
func (az *AzureClient) CreateRoleAssignmentSimple(resourceGroup, servicePrincipalObjectID string) error {
	roleAssignmentName := uuid.NewV4().String()

	roleDefinitionID := fmt.Sprintf(AADRoleReferenceTemplate, az.subscriptionID, AADContributorRoleID)
	scope := fmt.Sprintf(AADRoleResourceGroupScopeTemplate, az.subscriptionID, resourceGroup)

	roleAssignmentParameters := authorization.RoleAssignmentCreateParameters{
		Properties: &authorization.RoleAssignmentProperties{
			RoleDefinitionID: to.StringPtr(roleDefinitionID),
			PrincipalID:      to.StringPtr(servicePrincipalObjectID),
		},
	}

	re := regexp.MustCompile("(?i)status=(\\d+)")
	for {
		_, err := az.CreateRoleAssignment(
			scope,
			roleAssignmentName,
			roleAssignmentParameters,
		)
		if err != nil {
			match := re.FindStringSubmatch(err.Error())
			if match != nil && (match[1] == "403") {
				//insufficient permissions. stop now
				log.Debugf("Failed to create role assignment (will abort now): %q", err)
				return err
			}
			log.Debugf("Failed to create role assignment (will retry): %q", err)
			time.Sleep(3 * time.Second)
			continue
		}
		break
	}

	return nil
}

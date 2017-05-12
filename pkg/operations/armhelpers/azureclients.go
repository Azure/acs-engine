package util

// TODO: refactor a bunch of this out of dockermachine and this into a better azure package
// TODO: See if SDK folks want to own this... it's super generic
// TODO: fix the token cache - cache by authority + client_id

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/Azure/acs-engine/pkg/acsengine"

	"github.com/Azure/azure-sdk-for-go/arm/authorization"
	"github.com/Azure/azure-sdk-for-go/arm/compute"
	"github.com/Azure/azure-sdk-for-go/arm/resources/resources"
	"github.com/Azure/azure-sdk-for-go/arm/resources/subscriptions"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/to"
	log "github.com/Sirupsen/logrus"
	"github.com/mitchellh/go-homedir"
)

const (
	// AcsEngineClientID is the AAD ClientID for the CLI native application
	AcsEngineClientID = "76e0feec-6b7f-41f0-81a7-b1b944520261"

	// ApplicationDir is the name of the dir where the token is cached
	ApplicationDir = ".acsengine"
)

var (
	// RequiredResourceProviders is the list of Azure Resource Providers needed for ACS-Engine to function
	RequiredResourceProviders = []string{"Microsoft.Compute", "Microsoft.Storage", "Microsoft.Network"}
)

// AzureClient is the uber client
// If done right, we really shouldn't need SubscriptionID or TenantID in the client anywhere else
// they're only used for setting up the clients, we can add them back later if needed
type AzureClient struct {
	DeploymentsClient     resources.DeploymentsClient
	GroupsClient          resources.GroupsClient
	RoleAssignmentsClient authorization.RoleAssignmentsClient
	ResourcesClient       resources.GroupClient
	ProvidersClient       resources.ProvidersClient
	SubscriptionsClient   subscriptions.GroupClient
	VirtualMachinesClient compute.VirtualMachinesClient
}

// NewAzureClientWithDeviceAuth returns an AzureClient by having a user complete a device authentication flow
func NewAzureClientWithDeviceAuth(env azure.Environment, subscriptionID string) (*AzureClient, error) {
	oauthConfig, tenantID, err := getOAuthConfig(env, subscriptionID)
	if err != nil {
		return nil, err
	}

	home, err := homedir.Dir()
	if err != nil {
		return nil, fmt.Errorf("Failed to get user home directory to look for cached token: %q", err)
	}
	cachePath := filepath.Join(home, ApplicationDir, "cache", fmt.Sprintf("%s_%s.token.json", tenantID, AcsEngineClientID))

	rawToken, err := tryLoadCachedToken(cachePath)
	if err != nil {
		return nil, err
	}

	var armSpt *adal.ServicePrincipalToken
	if rawToken != nil {
		armSpt, err = adal.NewServicePrincipalTokenFromManualToken(*oauthConfig, AcsEngineClientID, env.ServiceManagementEndpoint, *rawToken, tokenCallback(cachePath))
		if err != nil {
			return nil, err
		}
		err = armSpt.Refresh()
		if err != nil {
			log.Warnf("Refresh token failed. Will fallback to device auth. %q", err)
		} else {
			adSpt, err := adal.NewServicePrincipalTokenFromManualToken(*oauthConfig, AcsEngineClientID, env.GraphEndpoint, armSpt.Token)
			if err != nil {
				return nil, err
			}
			return getClient(env, subscriptionID, armSpt, adSpt)
		}
	}

	client := &autorest.Client{}

	deviceCode, err := adal.InitiateDeviceAuth(client, *oauthConfig, AcsEngineClientID, env.ServiceManagementEndpoint)
	if err != nil {
		return nil, err
	}
	log.Warnln(*deviceCode.Message)
	deviceToken, err := adal.WaitForUserCompletion(client, deviceCode)
	if err != nil {
		return nil, err
	}

	armSpt, err = adal.NewServicePrincipalTokenFromManualToken(*oauthConfig, AcsEngineClientID, env.ServiceManagementEndpoint, *deviceToken, tokenCallback(cachePath))
	if err != nil {
		return nil, err
	}
	armSpt.Refresh()

	adRawToken := armSpt.Token
	adRawToken.Resource = env.GraphEndpoint
	adSpt, err := adal.NewServicePrincipalTokenFromManualToken(*oauthConfig, AcsEngineClientID, env.GraphEndpoint, adRawToken)
	if err != nil {
		return nil, err
	}

	return getClient(env, subscriptionID, armSpt, adSpt)
}

// NewAzureClientWithClientSecret returns an AzureClient via client_id and client_secret
func NewAzureClientWithClientSecret(env azure.Environment, subscriptionID, clientID, clientSecret string) (*AzureClient, error) {
	oauthConfig, _, err := getOAuthConfig(env, subscriptionID)
	if err != nil {
		return nil, err
	}

	armSpt, err := adal.NewServicePrincipalToken(*oauthConfig, clientID, clientSecret, env.ServiceManagementEndpoint)
	if err != nil {
		return nil, err
	}
	adSpt, err := adal.NewServicePrincipalToken(*oauthConfig, clientID, clientSecret, env.GraphEndpoint)
	if err != nil {
		return nil, err
	}

	return getClient(env, subscriptionID, armSpt, adSpt)
}

// NewAzureClientWithClientCertificate returns an AzureClient via client_id and jwt certificate assertion
func NewAzureClientWithClientCertificate(env azure.Environment, subscriptionID, clientID, certificatePath, privateKeyPath string) (*AzureClient, error) {
	oauthConfig, _, err := getOAuthConfig(env, subscriptionID)
	if err != nil {
		return nil, err
	}

	certificateData, err := ioutil.ReadFile(certificatePath)
	if err != nil {
		return nil, fmt.Errorf("Failed to read certificate: %q", err)
	}

	block, _ := pem.Decode(certificateData)
	if block == nil {
		return nil, fmt.Errorf("Failed to decode pem block from certificate")
	}

	certificate, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse certificate: %q", err)
	}

	privateKey, err := parseRsaPrivateKey(privateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse rsa private key: %q", err)
	}

	armSpt, err := adal.NewServicePrincipalTokenFromCertificate(*oauthConfig, clientID, certificate, privateKey, env.ServiceManagementEndpoint)
	if err != nil {
		return nil, err
	}
	adSpt, err := adal.NewServicePrincipalTokenFromCertificate(*oauthConfig, clientID, certificate, privateKey, env.GraphEndpoint)
	if err != nil {
		return nil, err
	}

	return getClient(env, subscriptionID, armSpt, adSpt)
}

func tokenCallback(path string) func(t adal.Token) error {
	return func(token adal.Token) error {
		err := adal.SaveToken(path, 0600, token)
		if err != nil {
			return err
		}
		log.Debugf("Saved token to cache. path=%q", path)
		return nil
	}
}

func tryLoadCachedToken(cachePath string) (*adal.Token, error) {
	log.Debugf("Attempting to load token from cache. path=%q", cachePath)

	// Check for file not found so we can suppress the file not found error
	// LoadToken doesn't discern and returns error either way
	if _, err := os.Stat(cachePath); err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	token, err := adal.LoadToken(cachePath)
	if err != nil {
		return nil, fmt.Errorf("Failed to load token from file: %v", err)
	}

	return token, nil
}

func getOAuthConfig(env azure.Environment, subscriptionID string) (*adal.OAuthConfig, string, error) {
	tenantID, err := acsengine.GetTenantID(env, subscriptionID)
	if err != nil {
		return nil, "", err
	}

	oauthConfig, err := adal.NewOAuthConfig(env.ActiveDirectoryEndpoint, tenantID)
	if err != nil {
		return nil, "", err
	}

	return oauthConfig, tenantID, nil
}

func getClient(env azure.Environment, subscriptionID string, armSpt *adal.ServicePrincipalToken, adSpt *adal.ServicePrincipalToken) (*AzureClient, error) {
	c := &AzureClient{
		DeploymentsClient:     resources.NewDeploymentsClientWithBaseURI(env.ResourceManagerEndpoint, subscriptionID),
		GroupsClient:          resources.NewGroupsClientWithBaseURI(env.ResourceManagerEndpoint, subscriptionID),
		RoleAssignmentsClient: authorization.NewRoleAssignmentsClientWithBaseURI(env.ResourceManagerEndpoint, subscriptionID),
		ResourcesClient:       resources.NewGroupClientWithBaseURI(env.ResourceManagerEndpoint, subscriptionID),
		ProvidersClient:       resources.NewProvidersClientWithBaseURI(env.ResourceManagerEndpoint, subscriptionID),
		VirtualMachinesClient: compute.NewVirtualMachinesClientWithBaseURI(env.ResourceManagerEndpoint, subscriptionID),
	}

	authorizer := autorest.NewBearerAuthorizer(armSpt)
	c.DeploymentsClient.Authorizer = authorizer
	c.GroupsClient.Authorizer = authorizer
	c.RoleAssignmentsClient.Authorizer = authorizer
	c.ResourcesClient.Authorizer = authorizer
	c.ProvidersClient.Authorizer = authorizer
	c.VirtualMachinesClient.Authorizer = authorizer

	err := c.ensureProvidersRegistered(subscriptionID)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (azureClient *AzureClient) ensureProvidersRegistered(subscriptionID string) error {
	registeredProviders, err := azureClient.ProvidersClient.List(to.Int32Ptr(100), "")
	if err != nil {
		return err
	}
	if registeredProviders.Value == nil {
		return fmt.Errorf("Providers list was nil. subscription=%q", subscriptionID)
	}

	m := make(map[string]bool)
	for _, provider := range *registeredProviders.Value {
		m[strings.ToLower(to.String(provider.Namespace))] = to.String(provider.RegistrationState) == "Registered"
	}

	for _, provider := range RequiredResourceProviders {
		registered, ok := m[strings.ToLower(provider)]
		if !ok {
			return fmt.Errorf("Unknown resource provider %q", provider)
		}
		if registered {
			log.Debugf("Already registered for %q", provider)
		} else {
			log.Info("Registering subscription to resource provider. provider=%q subscription=%q", provider, subscriptionID)
			if _, err := azureClient.ProvidersClient.Register(provider); err != nil {
				return err
			}
		}
	}
	return nil
}

func parseRsaPrivateKey(path string) (*rsa.PrivateKey, error) {
	privateKeyData, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(privateKeyData)
	if block == nil {
		return nil, fmt.Errorf("Failed to decode a pem block from private key")
	}

	privatePkcs1Key, errPkcs1 := x509.ParsePKCS1PrivateKey(block.Bytes)
	if errPkcs1 == nil {
		return privatePkcs1Key, nil
	}

	privatePkcs8Key, errPkcs8 := x509.ParsePKCS8PrivateKey(block.Bytes)
	if errPkcs8 == nil {
		privatePkcs8RsaKey, ok := privatePkcs8Key.(*rsa.PrivateKey)
		if !ok {
			return nil, fmt.Errorf("pkcs8 contained non-RSA key. Expected RSA key")
		}
		return privatePkcs8RsaKey, nil
	}

	return nil, fmt.Errorf("failed to parse private key as Pkcs#1 or Pkcs#8. (%s). (%s)", errPkcs1, errPkcs8)
}

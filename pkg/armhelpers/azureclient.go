package armhelpers

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/arm/compute"
	"github.com/Azure/azure-sdk-for-go/arm/network"
	"github.com/Azure/azure-sdk-for-go/arm/resources/resources"
	"github.com/Azure/azure-sdk-for-go/arm/resources/subscriptions"
	"github.com/Azure/azure-sdk-for-go/arm/storage"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/to"
	log "github.com/Sirupsen/logrus"
	"github.com/mitchellh/go-homedir"

	"github.com/Azure/acs-engine/pkg/acsengine"
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

// AzureClient implements the `ACSEngineClient` interface.
// This client is backed by real Azure clients talking to an ARM endpoint.
type AzureClient struct {
	environment azure.Environment

	deploymentsClient             resources.DeploymentsClient
	resourcesClient               resources.GroupClient
	storageAccountsClient         storage.AccountsClient
	interfacesClient              network.InterfacesClient
	groupsClient                  resources.GroupsClient
	providersClient               resources.ProvidersClient
	subscriptionsClient           subscriptions.GroupClient
	virtualMachinesClient         compute.VirtualMachinesClient
	virtualMachineScaleSetsClient compute.VirtualMachineScaleSetsClient
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
		environment:                   env,
		deploymentsClient:             resources.NewDeploymentsClientWithBaseURI(env.ResourceManagerEndpoint, subscriptionID),
		resourcesClient:               resources.NewGroupClientWithBaseURI(env.ResourceManagerEndpoint, subscriptionID),
		storageAccountsClient:         storage.NewAccountsClientWithBaseURI(env.ResourceManagerEndpoint, subscriptionID),
		interfacesClient:              network.NewInterfacesClientWithBaseURI(env.ResourceManagerEndpoint, subscriptionID),
		groupsClient:                  resources.NewGroupsClientWithBaseURI(env.ResourceManagerEndpoint, subscriptionID),
		providersClient:               resources.NewProvidersClientWithBaseURI(env.ResourceManagerEndpoint, subscriptionID),
		virtualMachinesClient:         compute.NewVirtualMachinesClientWithBaseURI(env.ResourceManagerEndpoint, subscriptionID),
		virtualMachineScaleSetsClient: compute.NewVirtualMachineScaleSetsClientWithBaseURI(env.ResourceManagerEndpoint, subscriptionID),
	}

	authorizer := autorest.NewBearerAuthorizer(armSpt)
	c.deploymentsClient.Authorizer = authorizer
	c.resourcesClient.Authorizer = authorizer
	c.storageAccountsClient.Authorizer = authorizer
	c.interfacesClient.Authorizer = authorizer
	c.groupsClient.Authorizer = authorizer
	c.providersClient.Authorizer = authorizer
	c.virtualMachinesClient.Authorizer = authorizer
	c.virtualMachineScaleSetsClient.Authorizer = authorizer

	c.deploymentsClient.PollingDelay = time.Second * 5

	err := c.ensureProvidersRegistered(subscriptionID)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (az *AzureClient) ensureProvidersRegistered(subscriptionID string) error {
	registeredProviders, err := az.providersClient.List(to.Int32Ptr(100), "")
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
			if _, err := az.providersClient.Register(provider); err != nil {
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

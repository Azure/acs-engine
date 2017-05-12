package util

// TODO: refactor a bunch of this out of dockermachine and this into a better azure package

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/Azure/azure-sdk-for-go/arm/authorization"
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
	AzkubeClientID = "a87032a7-203c-4bf7-913c-44c50d23409a"
)

var (
	RequiredResourceProviders = []string{"Microsoft.Compute", "Microsoft.Storage", "Microsoft.Network"}
)

type AzureClient struct {
	Environment    azure.Environment
	OAuthConfig    adal.OAuthConfig
	SubscriptionID string
	TenantID       string
	ClientID       string

	DeploymentsClient     resources.DeploymentsClient
	GroupsClient          resources.GroupsClient
	RoleAssignmentsClient authorization.RoleAssignmentsClient
	ResourcesClient       resources.GroupClient
	ProvidersClient       resources.ProvidersClient
	SubscriptionsClient   subscriptions.GroupClient
}

func NewClientWithDeviceAuth(azureEnvironment azure.Environment, subscriptionID, tenantID string) (*AzureClient, error) {
	oauthConfig, err := adal.OAuthConfigForTenant(azureEnvironment.ActiveDirectoryEndpoint, tenantID)
	if err != nil {
		return nil, err
	}

	azureClient := AzureClient{
		Environment:    azureEnvironment,
		OAuthConfig:    *oauthConfig,
		TenantID:       tenantID,
		SubscriptionID: subscriptionID,
		ClientID:       AzkubeClientID,
	}

	home, err := homedir.Dir()
	if err != nil {
		return nil, fmt.Errorf("Failed to get user home directory to look for cached token: %q", err)
	}
	cachePath := filepath.Join(home, ".azkube", fmt.Sprintf("token-cache-%s.json", tenantID))

	armSpt, err := azureClient.tryLoadToken(cachePath)
	if err != nil {
		return nil, err
	}
	if armSpt != nil {
		err = armSpt.Refresh()
		if err != nil {
			log.Warnf("Refresh token failed. Will fallback to device auth. %q", err)
		} else {
			adSpt, err := adal.NewServicePrincipalTokenFromManualToken(azureClient.OAuthConfig, azureClient.ClientID, azureClient.Environment.GraphEndpoint, armSpt.Token)
			if err != nil {
				return nil, err
			}
			return azureClient.build(armSpt, adSpt)
		}
	}

	client := &autorest.Client{}

	deviceCode, err := azure.InitiateDeviceAuth(client, *oauthConfig, AzkubeClientID, azureEnvironment.ServiceManagementEndpoint)
	if err != nil {
		return nil, err
	}
	fmt.Println(*deviceCode.Message)
	deviceToken, err := azure.WaitForUserCompletion(client, deviceCode)
	if err != nil {
		return nil, err
	}

	armSpt, err = adal.NewServicePrincipalTokenFromManualToken(*oauthConfig, AzkubeClientID, azureEnvironment.ServiceManagementEndpoint, *deviceToken, tokenCallback(cachePath))
	if err != nil {
		return nil, err
	}
	armSpt.Refresh()

	rawToken := armSpt.Token
	rawToken.Resource = azureClient.Environment.GraphEndpoint
	adSpt, err := adal.NewServicePrincipalTokenFromManualToken(azureClient.OAuthConfig, azureClient.ClientID, azureClient.Environment.GraphEndpoint, rawToken)
	if err != nil {
		return nil, err
	}

	return azureClient.build(armSpt, adSpt)
}

func NewClientWithClientSecret(azureEnvironment azure.Environment, subscriptionID, tenantID, clientID, clientSecret string) (*AzureClient, error) {
	oauthConfig, err := azureEnvironment.OAuthConfigForTenant(tenantID)
	if err != nil {
		return nil, err
	}

	azureClient := AzureClient{
		Environment:    azureEnvironment,
		OAuthConfig:    *oauthConfig,
		TenantID:       tenantID,
		SubscriptionID: subscriptionID,
		ClientID:       clientID,
	}

	armSpt, err := adal.NewServicePrincipalToken(*oauthConfig, clientID, clientSecret, azureEnvironment.ServiceManagementEndpoint)
	if err != nil {
		return nil, err
	}
	adSpt, err := adal.NewServicePrincipalToken(*oauthConfig, clientID, clientSecret, azureEnvironment.GraphEndpoint)
	if err != nil {
		return nil, err
	}

	return azureClient.build(armSpt, adSpt)
}

func NewClientWithClientCertificate(azureEnvironment azure.Environment, subscriptionID, tenantID, clientID, certificatePath, privateKeyPath string) (*AzureClient, error) {
	oauthConfig, err := azureEnvironment.OAuthConfigForTenant(tenantID)
	if err != nil {
		return nil, err
	}

	azureClient := AzureClient{
		Environment:    azureEnvironment,
		OAuthConfig:    *oauthConfig,
		TenantID:       tenantID,
		SubscriptionID: subscriptionID,
		ClientID:       clientID,
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

	armSpt, err := adal.NewServicePrincipalTokenFromCertificate(*oauthConfig, clientID, certificate, privateKey, azureEnvironment.ServiceManagementEndpoint)
	if err != nil {
		return nil, err
	}
	adSpt, err := adal.NewServicePrincipalTokenFromCertificate(*oauthConfig, clientID, certificate, privateKey, azureEnvironment.GraphEndpoint)
	if err != nil {
		return nil, err
	}

	return azureClient.build(armSpt, adSpt)
}

func tokenCallback(path string) func(t azure.Token) error {
	return func(token azure.Token) error {
		err := azure.SaveToken(path, 0600, token)
		if err != nil {
			return err
		}
		log.Debugf("Saved token to cache. path=%q", path)
		return nil
	}
}

func (azureClient *AzureClient) tryLoadToken(cachePath string) (*azure.ServicePrincipalToken, error) {
	log.Debugf("Attempting to load token from cache. path=%q", cachePath)

	if _, err := os.Stat(cachePath); err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	token, err := azure.LoadToken(cachePath)
	if err != nil {
		return nil, fmt.Errorf("Failed to load token from file: %v", err)
	}

	armSpt, err := adal.NewServicePrincipalTokenFromManualToken(azureClient.OAuthConfig, azureClient.ClientID, azureClient.Environment.ServiceManagementEndpoint, *token, tokenCallback(cachePath))
	if err != nil {
		return nil, fmt.Errorf("Error constructing service principal token: %v", err)
	}
	return armSpt, nil
}

func buildSpt(subscriptionID string, armSpt *azure.ServicePrincipalToken) (*AzureClient, error) {
	azureClient.DeploymentsClient = resources.NewDeploymentsClientWithBaseURI(env.ResourceManagerEndpoint, subscriptionID)
	azureClient.GroupsClient = resources.NewGroupsClient(env.ResourceManagerEndpoint, subscriptionID)
	azureClient.RoleAssignmentsClient = authorization.NewRoleAssignmentsClient(env.ResourceManagerEndpoint, subscriptionID)
	azureClient.ResourcesClient = resources.NewClient(env.ResourceManagerEndpoint, subscriptionID)
	azureClient.ProvidersClient = resources.NewProvidersClient(env.ResourceManagerEndpoint, subscriptionID)

	azureClient.DeploymentsClient.Authorizer = armSpt
	azureClient.GroupsClient.Authorizer = armSpt
	azureClient.RoleAssignmentsClient.Authorizer = armSpt
	azureClient.ResourcesClient.Authorizer = armSpt
	azureClient.ProvidersClient.Authorizer = armSpt

	err := azureClient.ensureProvidersRegistered(subscriptionID)
	if err != nil {
		return nil, err
	}

	return azureClient, nil
}

func (azureClient *AzureClient) ensureProvidersRegistered(subscriptionID string) error {
	registeredProviders, err := azureClient.ProvidersClient.List(nil)
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
			return nil, fmt.Errorf("Pkcs8 contained non-RSA key. Expected RSA key.")
		}
		return privatePkcs8RsaKey, nil
	}

	return nil, fmt.Errorf("Failed to parse private key as Pkcs#1 or Pkcs#8. (%s). (%s).", errPkcs1, errPkcs8)
}

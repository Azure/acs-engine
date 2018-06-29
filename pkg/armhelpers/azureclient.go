package armhelpers

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/arm/authorization"
	"github.com/Azure/azure-sdk-for-go/arm/compute"
	"github.com/Azure/azure-sdk-for-go/arm/graphrbac"
	"github.com/Azure/azure-sdk-for-go/arm/network"
	"github.com/Azure/azure-sdk-for-go/arm/resources/resources"
	"github.com/Azure/azure-sdk-for-go/arm/storage"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/Azure/acs-engine/pkg/acsengine"
	"github.com/Azure/azure-sdk-for-go/arm/disk"
)

const (
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
	acceptLanguages []string
	environment     azure.Environment
	subscriptionID  string

	authorizationClient             authorization.RoleAssignmentsClient
	deploymentsClient               resources.DeploymentsClient
	deploymentOperationsClient      resources.DeploymentOperationsClient
	resourcesClient                 resources.GroupClient
	storageAccountsClient           storage.AccountsClient
	interfacesClient                network.InterfacesClient
	groupsClient                    resources.GroupsClient
	providersClient                 resources.ProvidersClient
	virtualMachinesClient           compute.VirtualMachinesClient
	virtualMachineScaleSetsClient   compute.VirtualMachineScaleSetsClient
	virtualMachineScaleSetVMsClient compute.VirtualMachineScaleSetVMsClient
	disksClient                     disk.DisksClient

	applicationsClient      graphrbac.ApplicationsClient
	servicePrincipalsClient graphrbac.ServicePrincipalsClient
}

// NewAzureClientWithDeviceAuth returns an AzureClient by having a user complete a device authentication flow
func NewAzureClientWithDeviceAuth(env azure.Environment, subscriptionID string) (*AzureClient, error) {
	oauthConfig, tenantID, err := getOAuthConfig(env, subscriptionID)
	if err != nil {
		return nil, err
	}

	// AcsEngineClientID is the AAD ClientID for the CLI native application
	acsEngineClientID := getAcsEngineClientID(env.Name)

	home, err := homedir.Dir()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get user home directory to look for cached token")
	}
	cachePath := filepath.Join(home, ApplicationDir, "cache", fmt.Sprintf("%s_%s.token.json", tenantID, acsEngineClientID))

	rawToken, err := tryLoadCachedToken(cachePath)
	if err != nil {
		return nil, err
	}

	var armSpt *adal.ServicePrincipalToken
	if rawToken != nil {
		armSpt, err = adal.NewServicePrincipalTokenFromManualToken(*oauthConfig, acsEngineClientID, env.ServiceManagementEndpoint, *rawToken, tokenCallback(cachePath))
		if err != nil {
			return nil, err
		}
		err = armSpt.Refresh()
		if err != nil {
			log.Warnf("Refresh token failed. Will fallback to device auth. %q", err)
		} else {
			graphSpt, err := adal.NewServicePrincipalTokenFromManualToken(*oauthConfig, acsEngineClientID, env.GraphEndpoint, armSpt.Token)
			if err != nil {
				return nil, err
			}
			graphSpt.Refresh()

			return getClient(env, subscriptionID, tenantID, armSpt, graphSpt), nil
		}
	}

	client := &autorest.Client{}

	deviceCode, err := adal.InitiateDeviceAuth(client, *oauthConfig, acsEngineClientID, env.ServiceManagementEndpoint)
	if err != nil {
		return nil, err
	}
	log.Warnln(*deviceCode.Message)
	deviceToken, err := adal.WaitForUserCompletion(client, deviceCode)
	if err != nil {
		return nil, err
	}

	armSpt, err = adal.NewServicePrincipalTokenFromManualToken(*oauthConfig, acsEngineClientID, env.ServiceManagementEndpoint, *deviceToken, tokenCallback(cachePath))
	if err != nil {
		return nil, err
	}
	armSpt.Refresh()

	adRawToken := armSpt.Token
	adRawToken.Resource = env.GraphEndpoint
	graphSpt, err := adal.NewServicePrincipalTokenFromManualToken(*oauthConfig, acsEngineClientID, env.GraphEndpoint, adRawToken)
	if err != nil {
		return nil, err
	}
	graphSpt.Refresh()

	return getClient(env, subscriptionID, tenantID, armSpt, graphSpt), nil
}

// NewAzureClientWithClientSecret returns an AzureClient via client_id and client_secret
func NewAzureClientWithClientSecret(env azure.Environment, subscriptionID, clientID, clientSecret string) (*AzureClient, error) {
	oauthConfig, tenantID, err := getOAuthConfig(env, subscriptionID)
	if err != nil {
		return nil, err
	}

	armSpt, err := adal.NewServicePrincipalToken(*oauthConfig, clientID, clientSecret, env.ServiceManagementEndpoint)
	if err != nil {
		return nil, err
	}
	graphSpt, err := adal.NewServicePrincipalToken(*oauthConfig, clientID, clientSecret, env.GraphEndpoint)
	if err != nil {
		return nil, err
	}
	graphSpt.Refresh()

	return getClient(env, subscriptionID, tenantID, armSpt, graphSpt), nil
}

// NewAzureClientWithClientCertificateFile returns an AzureClient via client_id and jwt certificate assertion
func NewAzureClientWithClientCertificateFile(env azure.Environment, subscriptionID, clientID, certificatePath, privateKeyPath string) (*AzureClient, error) {
	certificateData, err := ioutil.ReadFile(certificatePath)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to read certificate")
	}

	block, _ := pem.Decode(certificateData)
	if block == nil {
		return nil, errors.New("Failed to decode pem block from certificate")
	}

	certificate, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to parse certificate")
	}

	privateKey, err := parseRsaPrivateKey(privateKeyPath)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to parse rsa private key")
	}

	return NewAzureClientWithClientCertificate(env, subscriptionID, clientID, certificate, privateKey)
}

// NewAzureClientWithClientCertificate returns an AzureClient via client_id and jwt certificate assertion
func NewAzureClientWithClientCertificate(env azure.Environment, subscriptionID, clientID string, certificate *x509.Certificate, privateKey *rsa.PrivateKey) (*AzureClient, error) {
	oauthConfig, tenantID, err := getOAuthConfig(env, subscriptionID)
	if err != nil {
		return nil, err
	}

	if certificate == nil {
		return nil, errors.New("certificate should not be nil")
	}

	if privateKey == nil {
		return nil, errors.New("privateKey should not be nil")
	}

	armSpt, err := adal.NewServicePrincipalTokenFromCertificate(*oauthConfig, clientID, certificate, privateKey, env.ServiceManagementEndpoint)
	if err != nil {
		return nil, err
	}
	graphSpt, err := adal.NewServicePrincipalTokenFromCertificate(*oauthConfig, clientID, certificate, privateKey, env.GraphEndpoint)
	if err != nil {
		return nil, err
	}
	graphSpt.Refresh()

	return getClient(env, subscriptionID, tenantID, armSpt, graphSpt), nil
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
		return nil, errors.Wrap(err, "Failed to load token from file")
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

func getAcsEngineClientID(envName string) string {
	switch envName {
	case "AzureUSGovernmentCloud":
		return "e8b7f94b-85c9-47f4-964a-98dafd7fc2d8"
	default:
		return "76e0feec-6b7f-41f0-81a7-b1b944520261"
	}
}

func getClient(env azure.Environment, subscriptionID, tenantID string, armSpt *adal.ServicePrincipalToken, graphSpt *adal.ServicePrincipalToken) *AzureClient {
	c := &AzureClient{
		environment:    env,
		subscriptionID: subscriptionID,

		authorizationClient:             authorization.NewRoleAssignmentsClientWithBaseURI(env.ResourceManagerEndpoint, subscriptionID),
		deploymentsClient:               resources.NewDeploymentsClientWithBaseURI(env.ResourceManagerEndpoint, subscriptionID),
		deploymentOperationsClient:      resources.NewDeploymentOperationsClientWithBaseURI(env.ResourceManagerEndpoint, subscriptionID),
		resourcesClient:                 resources.NewGroupClientWithBaseURI(env.ResourceManagerEndpoint, subscriptionID),
		storageAccountsClient:           storage.NewAccountsClientWithBaseURI(env.ResourceManagerEndpoint, subscriptionID),
		interfacesClient:                network.NewInterfacesClientWithBaseURI(env.ResourceManagerEndpoint, subscriptionID),
		groupsClient:                    resources.NewGroupsClientWithBaseURI(env.ResourceManagerEndpoint, subscriptionID),
		providersClient:                 resources.NewProvidersClientWithBaseURI(env.ResourceManagerEndpoint, subscriptionID),
		virtualMachinesClient:           compute.NewVirtualMachinesClientWithBaseURI(env.ResourceManagerEndpoint, subscriptionID),
		virtualMachineScaleSetsClient:   compute.NewVirtualMachineScaleSetsClientWithBaseURI(env.ResourceManagerEndpoint, subscriptionID),
		virtualMachineScaleSetVMsClient: compute.NewVirtualMachineScaleSetVMsClientWithBaseURI(env.ResourceManagerEndpoint, subscriptionID),
		disksClient:                     disk.NewDisksClientWithBaseURI(env.ResourceManagerEndpoint, subscriptionID),

		applicationsClient:      graphrbac.NewApplicationsClientWithBaseURI(env.GraphEndpoint, tenantID),
		servicePrincipalsClient: graphrbac.NewServicePrincipalsClientWithBaseURI(env.GraphEndpoint, tenantID),
	}

	authorizer := autorest.NewBearerAuthorizer(armSpt)
	c.authorizationClient.Authorizer = authorizer
	c.deploymentsClient.Authorizer = authorizer
	c.deploymentOperationsClient.Authorizer = authorizer
	c.resourcesClient.Authorizer = authorizer
	c.storageAccountsClient.Authorizer = authorizer
	c.interfacesClient.Authorizer = authorizer
	c.groupsClient.Authorizer = authorizer
	c.providersClient.Authorizer = authorizer
	c.virtualMachinesClient.Authorizer = authorizer
	c.virtualMachineScaleSetsClient.Authorizer = authorizer
	c.virtualMachineScaleSetVMsClient.Authorizer = authorizer
	c.disksClient.Authorizer = authorizer

	c.deploymentsClient.PollingDelay = time.Second * 5
	c.resourcesClient.PollingDelay = time.Second * 5

	graphAuthorizer := autorest.NewBearerAuthorizer(graphSpt)
	c.applicationsClient.Authorizer = graphAuthorizer
	c.servicePrincipalsClient.Authorizer = graphAuthorizer

	return c
}

// EnsureProvidersRegistered checks if the AzureClient is registered to required resource providers and, if not, register subscription to providers
func (az *AzureClient) EnsureProvidersRegistered(subscriptionID string) error {
	registeredProviders, err := az.providersClient.List(to.Int32Ptr(100), "")
	if err != nil {
		return err
	}
	if registeredProviders.Value == nil {
		return errors.Errorf("Providers list was nil. subscription=%q", subscriptionID)
	}

	m := make(map[string]bool)
	for _, provider := range *registeredProviders.Value {
		m[strings.ToLower(to.String(provider.Namespace))] = to.String(provider.RegistrationState) == "Registered"
	}

	for _, provider := range RequiredResourceProviders {
		registered, ok := m[strings.ToLower(provider)]
		if !ok {
			return errors.Errorf("Unknown resource provider %q", provider)
		}
		if registered {
			log.Debugf("Already registered for %q", provider)
		} else {
			log.Infof("Registering subscription to resource provider. provider=%q subscription=%q", provider, subscriptionID)
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
		return nil, errors.New("Failed to decode a pem block from private key")
	}

	privatePkcs1Key, errPkcs1 := x509.ParsePKCS1PrivateKey(block.Bytes)
	if errPkcs1 == nil {
		return privatePkcs1Key, nil
	}

	privatePkcs8Key, errPkcs8 := x509.ParsePKCS8PrivateKey(block.Bytes)
	if errPkcs8 == nil {
		privatePkcs8RsaKey, ok := privatePkcs8Key.(*rsa.PrivateKey)
		if !ok {
			return nil, errors.New("pkcs8 contained non-RSA key. Expected RSA key")
		}
		return privatePkcs8RsaKey, nil
	}

	return nil, errors.Errorf("failed to parse private key as Pkcs#1 or Pkcs#8. (%s). (%s)", errPkcs1, errPkcs8)
}

//AddAcceptLanguages sets the list of languages to accept on this request
func (az *AzureClient) AddAcceptLanguages(languages []string) {
	az.acceptLanguages = languages
	az.authorizationClient.ManagementClient.Client.RequestInspector = az.addAcceptLanguages()
	az.deploymentOperationsClient.ManagementClient.Client.RequestInspector = az.addAcceptLanguages()
	az.deploymentsClient.ManagementClient.Client.RequestInspector = az.addAcceptLanguages()
	az.deploymentsClient.ManagementClient.Client.RequestInspector = az.addAcceptLanguages()
	az.deploymentOperationsClient.ManagementClient.Client.RequestInspector = az.addAcceptLanguages()
	az.resourcesClient.ManagementClient.Client.RequestInspector = az.addAcceptLanguages()
	az.storageAccountsClient.ManagementClient.Client.RequestInspector = az.addAcceptLanguages()
	az.interfacesClient.ManagementClient.Client.RequestInspector = az.addAcceptLanguages()
	az.groupsClient.ManagementClient.Client.RequestInspector = az.addAcceptLanguages()
	az.providersClient.ManagementClient.Client.RequestInspector = az.addAcceptLanguages()
	az.virtualMachinesClient.ManagementClient.Client.RequestInspector = az.addAcceptLanguages()
	az.virtualMachineScaleSetsClient.ManagementClient.Client.RequestInspector = az.addAcceptLanguages()
	az.disksClient.ManagementClient.Client.RequestInspector = az.addAcceptLanguages()

	az.applicationsClient.ManagementClient.Client.RequestInspector = az.addAcceptLanguages()
	az.servicePrincipalsClient.ManagementClient.Client.RequestInspector = az.addAcceptLanguages()
}

func (az *AzureClient) addAcceptLanguages() autorest.PrepareDecorator {
	return func(p autorest.Preparer) autorest.Preparer {
		return autorest.PreparerFunc(func(r *http.Request) (*http.Request, error) {
			r, err := p.Prepare(r)
			if err != nil {
				return r, err
			}
			if az.acceptLanguages != nil {
				for _, language := range az.acceptLanguages {
					r.Header.Add("Accept-Language", language)
				}
			}
			return r, nil
		})
	}
}

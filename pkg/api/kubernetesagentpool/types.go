package kubernetesagentpool

const (
	APIVersion = "kubernetesagentpool"
)

// ContainerService complies with the ARM model of
// resource definition in a JSON template.
type ContainerService struct {
	ID       string                `json:"id,omitempty"`
	Location string                `json:"location,omitempty"`
	Name     string                `json:"name,omitempty"`
	Plan     *ResourcePurchasePlan `json:"plan,omitempty"`
	Tags     map[string]string     `json:"tags,omitempty"`
	Type     string                `json:"type,omitempty"`

	Properties *Properties `json:"properties"`
}

type Properties struct {
	DnsPrefix               string                   `json:"dnsPrefix,omitempty"`
	Version                 string                   `json:"version,omitempty"`
	AgentPoolProfiles       []*AgentPoolProfile      `json:"agentPoolProfiles,omitempty"`
	WindowsProfile          *WindowsProfile          `json:"windowsProfile,omitempty"`
	LinuxProfile            *LinuxProfile            `json:"linuxProfile,omitempty"`
	NetworkProfile          *NetworkProfile          `json:"networkProfile,omitempty"`
	JumpBoxProfile          *JumpBoxProfile          `json:"jumpboxProfile,omitempty"`
	ServicePrincipalProfile *ServicePrincipalProfile `json:"servicePrincipalProfile,omitempty"`
}

// ServicePrincipalProfile contains the client and secret used by the cluster for Azure Resource CRUD
// The 'Secret' parameter could be either a plain text, or referenced to a secret in a keyvault.
// In the latter case, the format of the parameter's value should be
// "/subscriptions/<SUB_ID>/resourceGroups/<RG_NAME>/providers/Microsoft.KeyVault/vaults/<KV_NAME>/secrets/<NAME>[/<VERSION>]"
// where:
//    <SUB_ID> is the subscription ID of the keyvault
//    <RG_NAME> is the resource group of the keyvault
//    <KV_NAME> is the name of the keyvault
//    <NAME> is the name of the secret.
//    <VERSION> (optional) is the version of the secret (default: the latest version)
type ServicePrincipalProfile struct {
	ClientID string `json:"servicePrincipalClientID,omitempty"`
	Secret   string `json:"servicePrincipalClientSecret,omitempty"`
}

type JumpBoxProfile struct {
	PublicIpAddressId string `json:"publicIpAddressId,omitempty"`
}

type NetworkProfile struct {
	PodCIDR          string `json:"podCidr,omitempty"`
	ServiceCIDR      string `json:"serviceCidr,omitempty"`
	VnetSubnetId     string `json:"vnetSubnetID,omitempty"`
	KubeDnsServiceIp string `json:"kubeDnsServiceIP,omitempty"`
}

type AgentPoolProfile struct {
	Name         string `json:"name,omitempty"`
	Count        int    `json:"count,omitempty"`
	VMSize       string `json:"vmSize,omitempty"`
	OSType       string `json:"osType,omitempty"`
	OSDiskSizeGb string `json:"osDiskSizeGb,omitempty"`
}

// WindowsProfile represents the windows parameters passed to the cluster
type WindowsProfile struct {
	AdminUsername string            `json:"adminUsername,omitempty"`
	AdminPassword string            `json:"adminPassword,omitempty"`
	Secrets       []KeyVaultSecrets `json:"secrets,omitempty"`
}

// LinuxProfile represents the linux parameters passed to the cluster
type LinuxProfile struct {
	AdminUsername string `json:"adminUsername"`
	SSH           struct {
		PublicKeys []struct {
			KeyData string `json:"keyData"`
		} `json:"publicKeys"`
	} `json:"ssh"`
	Secrets []KeyVaultSecrets `json:"secrets,omitempty"`
}

// KeyVaultSecrets specifies certificates to install on the pool
// of machines from a given key vault
// the key vault specified must have been granted read permissions to CRP
type KeyVaultSecrets struct {
	SourceVault       *KeyVaultID           `json:"sourceVault,omitempty"`
	VaultCertificates []KeyVaultCertificate `json:"vaultCertificates,omitempty"`
}

// KeyVaultID specifies a key vault
type KeyVaultID struct {
	ID string `json:"id,omitempty"`
}

// KeyVaultCertificate specifies a certificate to install
// On Linux, the certificate file is placed under the /var/lib/waagent directory
// with the file name <UppercaseThumbprint>.crt for the X509 certificate file
// and <UppercaseThumbprint>.prv for the private key. Both of these files are .pem formatted.
// On windows the certificate will be saved in the specified store.
type KeyVaultCertificate struct {
	CertificateURL   string `json:"certificateUrl,omitempty"`
	CertificateStore string `json:"certificateStore,omitempty"`
}

// ResourcePurchasePlan defines resource plan as required by ARM
// for billing purposes.
type ResourcePurchasePlan struct {
	Name          string `json:"name,omitempty"`
	Product       string `json:"product,omitempty"`
	PromotionCode string `json:"promotionCode,omitempty"`
	Publisher     string `json:"publisher,omitempty"`
}

//{
//  "location": "westus",
//  "tags": {
//    "key": "value"
//  },
//  "properties": {
//    "dnsPrefix": "masterdns1",
//    "version": "1.6.6",
//    "agentPoolProfiles": [
//      {
//        "name": "agentpool1",
//        "count": 3,
//        "vmSize": "Standard_D2_v2",
//        "osType": "Linux",
//        "osDiskSizeGB": 128
//      },
//      {
//        "name": "agentpool2",
//        "count": 3,
//        "vmSize": "Standard_D2_v2",
//        "osType": "Windows"
//      }
//    ],
//    "windowsProfile": {
//      "adminUsername": "azureuser",
//      "adminPassword": "replacepassword1234$"
//    },
//    "linuxProfile": {
//      "adminUsername": "azureuser",
//      "ssh": {
//        "publicKeys": [
//          {
//            "keyData": "ssh-rsa PUBLICKEY azureuser@linuxvm"
//          }
//        ]
//      }
//    },
//    "networkProfile" : {
//      "podCidr" : "10.0.0.0/24",
//      "serviceCidr" : "10.0.0.0/16",
//      "vnetSubnetID" : "/subscriptions/c1239427-83d3-12346-9f35-5af546a6eb67/resourceGroups/testresourcegroup/providers/Microsoft.Network/virtualNetworks/vnet1/subnets/subnet1",
//      "kubeDnsServiceIP" : "10.240.0.4"
//    },
//    "jumpboxProfile": {
//      "publicIpAddressId" : "/subscriptions/c1239427-83d3-12346-9f35-5af546a6eb67/resourceGroups/testresourcegroup/providers/Microsoft.Network/publicIPAddresses/ubuntu-ip"
//    },
//    "servicePrincipalProfile": {
//      "clientId": "ServicePrincipalClientID",
//      "secret": "ServicePrincipalClientSecret"
//    }
//  }
//}

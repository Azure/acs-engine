package api

import (
	"log"
	"regexp"
	"testing"

	"github.com/Azure/acs-engine/pkg/helpers"
)

const exampleCustomHyperkubeImage = `example.azurecr.io/example/hyperkube-amd64:custom`

const exampleAPIModel = `{
		"apiVersion": "vlabs",
	"properties": {
		"orchestratorProfile": {
			"orchestratorType": "Kubernetes",
			"kubernetesConfig": {
				"customHyperkubeImage": "` + exampleCustomHyperkubeImage + `"
			}
		},
		"masterProfile": { "count": 1, "dnsPrefix": "", "vmSize": "Standard_D2_v2" },
		"agentPoolProfiles": [ { "name": "linuxpool1", "count": 2, "vmSize": "Standard_D2_v2", "availabilityProfile": "AvailabilitySet" } ],
		"windowsProfile": { "adminUsername": "azureuser", "adminPassword": "replacepassword1234$" },
		"linuxProfile": { "adminUsername": "azureuser", "ssh": { "publicKeys": [ { "keyData": "" } ] }
		},
		"servicePrincipalProfile": { "clientId": "", "secret": "" }
	}
}
`

const exampleSystemMSIModel = `{
	"apiVersion": "vlabs",
"properties": {
	"orchestratorProfile": {
		"orchestratorType": "Kubernetes",
		"kubernetesConfig": {
			"useManagedIdentity": true
		}
	},
	"masterProfile": { "count": 1, "dnsPrefix": "", "vmSize": "Standard_D2_v2" },
	"agentPoolProfiles": [ { "name": "linuxpool1", "count": 2, "vmSize": "Standard_D2_v2", "availabilityProfile": "AvailabilitySet" } ],
	"windowsProfile": { "adminUsername": "azureuser", "adminPassword": "replacepassword1234$" },
	"linuxProfile": { "adminUsername": "azureuser", "ssh": { "publicKeys": [ { "keyData": "" } ] }
	},
	"servicePrincipalProfile": { "clientId": "", "secret": "" }
}
}
`

const exampleUserMSI = "/subscriptions/<subscription>/resourcegroups/<rg_name>/providers/Microsoft.ManagedIdentity/userAssignedIdentities/<identityName>"

const exampleUserMSIModel = `{
	"apiVersion": "vlabs",
"properties": {
	"orchestratorProfile": {
		"orchestratorType": "Kubernetes",
		"kubernetesConfig": {
			"useManagedIdentity": true,
			"userAssignedID": "` + exampleUserMSI + `"
		}
	},
	"masterProfile": { "count": 1, "dnsPrefix": "", "vmSize": "Standard_D2_v2" },
	"agentPoolProfiles": [ { "name": "linuxpool1", "count": 2, "vmSize": "Standard_D2_v2", "availabilityProfile": "AvailabilitySet" } ],
	"windowsProfile": { "adminUsername": "azureuser", "adminPassword": "replacepassword1234$" },
	"linuxProfile": { "adminUsername": "azureuser", "ssh": { "publicKeys": [ { "keyData": "" } ] }
	},
	"servicePrincipalProfile": { "clientId": "", "secret": "" }
}
}
`

func TestOSType(t *testing.T) {
	p := Properties{
		MasterProfile: &MasterProfile{
			Distro: RHEL,
		},
		AgentPoolProfiles: []*AgentPoolProfile{
			{
				OSType: Linux,
			},
			{
				OSType: Linux,
				Distro: RHEL,
			},
		},
	}

	if p.HasWindows() {
		t.Fatalf("expected HasWindows() to return false but instead returned true")
	}
	if p.AgentPoolProfiles[0].IsWindows() {
		t.Fatalf("expected IsWindows() to return false but instead returned true")
	}

	if !p.AgentPoolProfiles[0].IsLinux() {
		t.Fatalf("expected IsLinux() to return true but instead returned false")
	}

	if p.AgentPoolProfiles[0].IsRHEL() {
		t.Fatalf("expected IsRHEL() to return false but instead returned true")
	}

	if p.AgentPoolProfiles[0].IsCoreOS() {
		t.Fatalf("expected IsCoreOS() to return false but instead returned true")
	}

	if !p.AgentPoolProfiles[1].IsRHEL() {
		t.Fatalf("expected IsRHEL() to return true but instead returned false")
	}

	if p.AgentPoolProfiles[1].IsCoreOS() {
		t.Fatalf("expected IsCoreOS() to return false but instead returned true")
	}

	if !p.MasterProfile.IsRHEL() {
		t.Fatalf("expected IsRHEL() to return true but instead returned false")
	}

	if p.MasterProfile.IsCoreOS() {
		t.Fatalf("expected IsCoreOS() to return false but instead returned true")
	}

	p.MasterProfile.Distro = CoreOS
	p.AgentPoolProfiles[0].OSType = Windows
	p.AgentPoolProfiles[1].Distro = CoreOS

	if !p.HasWindows() {
		t.Fatalf("expected HasWindows() to return true but instead returned false")
	}

	if !p.AgentPoolProfiles[0].IsWindows() {
		t.Fatalf("expected IsWindows() to return true but instead returned false")
	}

	if p.AgentPoolProfiles[0].IsLinux() {
		t.Fatalf("expected IsLinux() to return false but instead returned true")
	}

	if p.AgentPoolProfiles[0].IsRHEL() {
		t.Fatalf("expected IsRHEL() to return false but instead returned true")
	}

	if p.AgentPoolProfiles[0].IsCoreOS() {
		t.Fatalf("expected IsCoreOS() to return false but instead returned true")
	}

	if p.AgentPoolProfiles[1].IsRHEL() {
		t.Fatalf("expected IsRHEL() to return false but instead returned true")
	}

	if !p.AgentPoolProfiles[1].IsCoreOS() {
		t.Fatalf("expected IsCoreOS() to return true but instead returned false")
	}

	if p.MasterProfile.IsRHEL() {
		t.Fatalf("expected IsRHEL() to return false but instead returned true")
	}

	if !p.MasterProfile.IsCoreOS() {
		t.Fatalf("expected IsCoreOS() to return true but instead returned false")
	}
}

func TestHasStorageProfile(t *testing.T) {
	cases := []struct {
		p                 Properties
		expectedHasMD     bool
		expectedHasSA     bool
		expectedMasterMD  bool
		expectedAgent0MD  bool
		expectedPrivateJB bool
		expectedHasDisks  bool
	}{
		{
			p: Properties{
				MasterProfile: &MasterProfile{
					StorageProfile: StorageAccount,
				},
				AgentPoolProfiles: []*AgentPoolProfile{
					{
						StorageProfile: StorageAccount,
						DiskSizesGB:    []int{5},
					},
					{
						StorageProfile: StorageAccount,
					},
				},
			},
			expectedHasMD:    false,
			expectedHasSA:    true,
			expectedMasterMD: false,
			expectedAgent0MD: false,
			expectedHasDisks: true,
		},
		{
			p: Properties{
				MasterProfile: &MasterProfile{
					StorageProfile: ManagedDisks,
				},
				AgentPoolProfiles: []*AgentPoolProfile{
					{
						StorageProfile: StorageAccount,
					},
					{
						StorageProfile: StorageAccount,
					},
				},
			},
			expectedHasMD:    true,
			expectedHasSA:    true,
			expectedMasterMD: true,
			expectedAgent0MD: false,
		},
		{
			p: Properties{
				MasterProfile: &MasterProfile{
					StorageProfile: StorageAccount,
				},
				AgentPoolProfiles: []*AgentPoolProfile{
					{
						StorageProfile: ManagedDisks,
					},
					{
						StorageProfile: StorageAccount,
					},
				},
			},
			expectedHasMD:    true,
			expectedHasSA:    true,
			expectedMasterMD: false,
			expectedAgent0MD: true,
		},
		{
			p: Properties{
				OrchestratorProfile: &OrchestratorProfile{
					OrchestratorType: Kubernetes,
				},
				MasterProfile: &MasterProfile{
					StorageProfile: ManagedDisks,
				},
				AgentPoolProfiles: []*AgentPoolProfile{
					{
						StorageProfile: ManagedDisks,
					},
					{
						StorageProfile: ManagedDisks,
					},
				},
			},
			expectedHasMD:     true,
			expectedHasSA:     false,
			expectedMasterMD:  true,
			expectedAgent0MD:  true,
			expectedPrivateJB: false,
		},
		{
			p: Properties{
				OrchestratorProfile: &OrchestratorProfile{
					OrchestratorType: Kubernetes,
					KubernetesConfig: &KubernetesConfig{
						PrivateCluster: &PrivateCluster{
							Enabled: helpers.PointerToBool(true),
							JumpboxProfile: &PrivateJumpboxProfile{
								StorageProfile: ManagedDisks,
							},
						},
					},
				},
				MasterProfile: &MasterProfile{
					StorageProfile: StorageAccount,
				},
				AgentPoolProfiles: []*AgentPoolProfile{
					{
						StorageProfile: StorageAccount,
					},
				},
			},
			expectedHasMD:     true,
			expectedHasSA:     true,
			expectedMasterMD:  false,
			expectedAgent0MD:  false,
			expectedPrivateJB: true,
		},

		{
			p: Properties{
				OrchestratorProfile: &OrchestratorProfile{
					OrchestratorType: Kubernetes,
					KubernetesConfig: &KubernetesConfig{
						PrivateCluster: &PrivateCluster{
							Enabled: helpers.PointerToBool(true),
							JumpboxProfile: &PrivateJumpboxProfile{
								StorageProfile: StorageAccount,
							},
						},
					},
				},
				MasterProfile: &MasterProfile{
					StorageProfile: ManagedDisks,
				},
				AgentPoolProfiles: []*AgentPoolProfile{
					{
						StorageProfile: ManagedDisks,
					},
				},
			},
			expectedHasMD:     true,
			expectedHasSA:     true,
			expectedMasterMD:  true,
			expectedAgent0MD:  true,
			expectedPrivateJB: true,
		},
	}

	for _, c := range cases {
		if c.p.HasManagedDisks() != c.expectedHasMD {
			t.Fatalf("expected HasManagedDisks() to return %t but instead returned %t", c.expectedHasMD, c.p.HasManagedDisks())
		}
		if c.p.HasStorageAccountDisks() != c.expectedHasSA {
			t.Fatalf("expected HasStorageAccountDisks() to return %t but instead returned %t", c.expectedHasSA, c.p.HasStorageAccountDisks())
		}
		if c.p.MasterProfile.IsManagedDisks() != c.expectedMasterMD {
			t.Fatalf("expected IsManagedDisks() to return %t but instead returned %t", c.expectedMasterMD, c.p.MasterProfile.IsManagedDisks())
		}
		if c.p.MasterProfile.IsStorageAccount() == c.expectedMasterMD {
			t.Fatalf("expected IsStorageAccount() to return %t but instead returned %t", !c.expectedMasterMD, c.p.MasterProfile.IsStorageAccount())
		}
		if c.p.AgentPoolProfiles[0].IsManagedDisks() != c.expectedAgent0MD {
			t.Fatalf("expected IsManagedDisks() to return %t but instead returned %t", c.expectedAgent0MD, c.p.AgentPoolProfiles[0].IsManagedDisks())
		}
		if c.p.AgentPoolProfiles[0].IsStorageAccount() == c.expectedAgent0MD {
			t.Fatalf("expected IsStorageAccount() to return %t but instead returned %t", !c.expectedAgent0MD, c.p.AgentPoolProfiles[0].IsStorageAccount())
		}
		if c.p.OrchestratorProfile != nil && c.p.OrchestratorProfile.KubernetesConfig.PrivateJumpboxProvision() != c.expectedPrivateJB {
			t.Fatalf("expected PrivateJumpboxProvision() to return %t but instead returned %t", c.expectedPrivateJB, c.p.OrchestratorProfile.KubernetesConfig.PrivateJumpboxProvision())
		}
		if c.p.AgentPoolProfiles[0].HasDisks() != c.expectedHasDisks {
			t.Fatalf("expected HasDisks() to return %t but instead returned %t", c.expectedHasDisks, c.p.AgentPoolProfiles[0].HasDisks())
		}
	}
}

func TestTotalNodes(t *testing.T) {
	cases := []struct {
		p        Properties
		expected int
	}{
		{
			p: Properties{
				MasterProfile: &MasterProfile{
					Count: 1,
				},
				AgentPoolProfiles: []*AgentPoolProfile{
					{
						Count: 1,
					},
				},
			},
			expected: 2,
		},
		{
			p: Properties{
				AgentPoolProfiles: []*AgentPoolProfile{
					{
						Count: 3,
					},
					{
						Count: 4,
					},
				},
			},
			expected: 7,
		},
		{
			p: Properties{
				MasterProfile: &MasterProfile{
					Count: 5,
				},
				AgentPoolProfiles: []*AgentPoolProfile{
					{
						Count: 6,
					},
				},
			},
			expected: 11,
		},
	}

	for _, c := range cases {
		if c.p.TotalNodes() != c.expected {
			t.Fatalf("expected TotalNodes() to return %d but instead returned %d", c.expected, c.p.TotalNodes())
		}
	}
}
func TestMasterAvailabilityProfile(t *testing.T) {
	cases := []struct {
		p              Properties
		expectedISVMSS bool
	}{
		{
			p: Properties{
				MasterProfile: &MasterProfile{},
			},
			expectedISVMSS: false,
		},
		{
			p: Properties{
				MasterProfile: &MasterProfile{
					AvailabilityProfile: AvailabilitySet,
				},
			},
			expectedISVMSS: false,
		},
		{
			p: Properties{
				MasterProfile: &MasterProfile{
					AvailabilityProfile: VirtualMachineScaleSets,
				},
			},
			expectedISVMSS: true,
		},
	}

	for _, c := range cases {
		if c.p.MasterProfile.IsVirtualMachineScaleSets() != c.expectedISVMSS {
			t.Fatalf("expected MasterProfile.IsVirtualMachineScaleSets() to return %t but instead returned %t", c.expectedISVMSS, c.p.MasterProfile.IsVirtualMachineScaleSets())
		}
	}
}
func TestAvailabilityProfile(t *testing.T) {
	cases := []struct {
		p               Properties
		expectedHasVMSS bool
		expectedISVMSS  bool
		expectedIsAS    bool
		expectedLowPri  bool
	}{
		{
			p: Properties{
				AgentPoolProfiles: []*AgentPoolProfile{
					{
						AvailabilityProfile: VirtualMachineScaleSets,
						ScaleSetPriority:    ScaleSetPriorityLow,
					},
				},
			},
			expectedHasVMSS: true,
			expectedISVMSS:  true,
			expectedIsAS:    false,
			expectedLowPri:  true,
		},
		{
			p: Properties{
				AgentPoolProfiles: []*AgentPoolProfile{
					{
						AvailabilityProfile: VirtualMachineScaleSets,
						ScaleSetPriority:    ScaleSetPriorityRegular,
					},
					{
						AvailabilityProfile: AvailabilitySet,
					},
				},
			},
			expectedHasVMSS: true,
			expectedISVMSS:  true,
			expectedIsAS:    false,
			expectedLowPri:  false,
		},
		{
			p: Properties{
				AgentPoolProfiles: []*AgentPoolProfile{
					{
						AvailabilityProfile: AvailabilitySet,
					},
				},
			},
			expectedHasVMSS: false,
			expectedISVMSS:  false,
			expectedIsAS:    true,
			expectedLowPri:  false,
		},
	}

	for _, c := range cases {
		if c.p.HasVMSSAgentPool() != c.expectedHasVMSS {
			t.Fatalf("expected HasVMSSAgentPool() to return %t but instead returned %t", c.expectedHasVMSS, c.p.HasVMSSAgentPool())
		}
		if c.p.AgentPoolProfiles[0].IsVirtualMachineScaleSets() != c.expectedISVMSS {
			t.Fatalf("expected IsVirtualMachineScaleSets() to return %t but instead returned %t", c.expectedISVMSS, c.p.AgentPoolProfiles[0].IsVirtualMachineScaleSets())
		}
		if c.p.AgentPoolProfiles[0].IsAvailabilitySets() != c.expectedIsAS {
			t.Fatalf("expected IsAvailabilitySets() to return %t but instead returned %t", c.expectedIsAS, c.p.AgentPoolProfiles[0].IsAvailabilitySets())
		}
		if c.p.AgentPoolProfiles[0].IsLowPriorityScaleSet() != c.expectedLowPri {
			t.Fatalf("expected IsLowPriorityScaleSet() to return %t but instead returned %t", c.expectedLowPri, c.p.AgentPoolProfiles[0].IsLowPriorityScaleSet())
		}
	}
}

func TestIsCustomVNET(t *testing.T) {
	cases := []struct {
		p              Properties
		expectedMaster bool
		expectedAgent  bool
	}{
		{
			p: Properties{
				MasterProfile: &MasterProfile{
					VnetSubnetID: "testSubnet",
				},
				AgentPoolProfiles: []*AgentPoolProfile{
					{
						VnetSubnetID: "testSubnet",
					},
				},
			},
			expectedMaster: true,
			expectedAgent:  true,
		},
		{
			p: Properties{
				MasterProfile: &MasterProfile{
					Count: 1,
				},
				AgentPoolProfiles: []*AgentPoolProfile{
					{
						Count: 1,
					},
					{
						Count: 1,
					},
				},
			},
			expectedMaster: false,
			expectedAgent:  false,
		},
	}

	for _, c := range cases {
		if c.p.MasterProfile.IsCustomVNET() != c.expectedMaster {
			t.Fatalf("expected IsCustomVnet() to return %t but instead returned %t", c.expectedMaster, c.p.MasterProfile.IsCustomVNET())
		}
		if c.p.AgentPoolProfiles[0].IsCustomVNET() != c.expectedAgent {
			t.Fatalf("expected IsCustomVnet() to return %t but instead returned %t", c.expectedAgent, c.p.AgentPoolProfiles[0].IsCustomVNET())
		}
	}

}

func TestHasAvailabilityZones(t *testing.T) {
	cases := []struct {
		p                Properties
		expectedMaster   bool
		expectedAgent    bool
		expectedAllZones bool
	}{
		{
			p: Properties{
				MasterProfile: &MasterProfile{
					Count:             1,
					AvailabilityZones: []string{"1", "2"},
				},
				AgentPoolProfiles: []*AgentPoolProfile{
					{
						Count:             1,
						AvailabilityZones: []string{"1", "2"},
					},
					{
						Count:             1,
						AvailabilityZones: []string{"1", "2"},
					},
				},
			},
			expectedMaster:   true,
			expectedAgent:    true,
			expectedAllZones: true,
		},
		{
			p: Properties{
				MasterProfile: &MasterProfile{
					Count: 1,
				},
				AgentPoolProfiles: []*AgentPoolProfile{
					{
						Count: 1,
					},
					{
						Count:             1,
						AvailabilityZones: []string{"1", "2"},
					},
				},
			},
			expectedMaster:   false,
			expectedAgent:    false,
			expectedAllZones: false,
		},
		{
			p: Properties{
				MasterProfile: &MasterProfile{
					Count: 1,
				},
				AgentPoolProfiles: []*AgentPoolProfile{
					{
						Count:             1,
						AvailabilityZones: []string{},
					},
					{
						Count:             1,
						AvailabilityZones: []string{"1", "2"},
					},
				},
			},
			expectedMaster:   false,
			expectedAgent:    false,
			expectedAllZones: false,
		},
	}

	for _, c := range cases {
		if c.p.MasterProfile.HasAvailabilityZones() != c.expectedMaster {
			t.Fatalf("expected HasAvailabilityZones() to return %t but instead returned %t", c.expectedMaster, c.p.MasterProfile.HasAvailabilityZones())
		}
		if c.p.AgentPoolProfiles[0].HasAvailabilityZones() != c.expectedAgent {
			t.Fatalf("expected HasAvailabilityZones() to return %t but instead returned %t", c.expectedAgent, c.p.AgentPoolProfiles[0].HasAvailabilityZones())
		}
		if c.p.HasZonesForAllAgentPools() != c.expectedAllZones {
			t.Fatalf("expected HasZonesForAllAgentPools() to return %t but instead returned %t", c.expectedAllZones, c.p.HasZonesForAllAgentPools())
		}
	}
}

func TestRequireRouteTable(t *testing.T) {
	cases := []struct {
		p        Properties
		expected bool
	}{
		{
			p: Properties{
				OrchestratorProfile: &OrchestratorProfile{
					OrchestratorType: DCOS,
				},
			},
			expected: false,
		},
		{
			p: Properties{
				OrchestratorProfile: &OrchestratorProfile{
					OrchestratorType: Kubernetes,
					KubernetesConfig: &KubernetesConfig{
						NetworkPolicy: "",
					},
				},
			},
			expected: true,
		},
		{
			p: Properties{
				OrchestratorProfile: &OrchestratorProfile{
					OrchestratorType: Kubernetes,
					KubernetesConfig: &KubernetesConfig{
						NetworkPlugin: "azure",
					},
				},
			},
			expected: false,
		},
		{
			p: Properties{
				OrchestratorProfile: &OrchestratorProfile{
					OrchestratorType: Kubernetes,
					KubernetesConfig: &KubernetesConfig{
						NetworkPolicy: "cilium",
					},
				},
			},
			expected: false,
		},
	}

	for _, c := range cases {
		if c.p.OrchestratorProfile.RequireRouteTable() != c.expected {
			t.Fatalf("expected RequireRouteTable() to return %t but instead got %t", c.expected, c.p.OrchestratorProfile.RequireRouteTable())
		}
	}
}

func TestIsAzureCNI(t *testing.T) {
	k := &KubernetesConfig{
		NetworkPlugin: "azure",
	}

	o := &OrchestratorProfile{
		KubernetesConfig: k,
	}
	if !o.IsAzureCNI() {
		t.Fatalf("unable to detect orchestrator profile is using Azure CNI from NetworkPlugin=%s", o.KubernetesConfig.NetworkPlugin)
	}

	k = &KubernetesConfig{
		NetworkPlugin: "none",
	}

	o = &OrchestratorProfile{
		KubernetesConfig: k,
	}
	if o.IsAzureCNI() {
		t.Fatalf("unable to detect orchestrator profile is not using Azure CNI from NetworkPlugin=%s", o.KubernetesConfig.NetworkPlugin)
	}

	o = &OrchestratorProfile{}
	if o.IsAzureCNI() {
		t.Fatalf("unable to detect orchestrator profile is not using Azure CNI from nil KubernetesConfig")
	}
}

func TestOrchestrator(t *testing.T) {
	cases := []struct {
		p                    Properties
		expectedIsDCOS       bool
		expectedIsKubernetes bool
		expectedIsOpenShift  bool
		expectedIsSwarmMode  bool
	}{
		{
			p: Properties{
				OrchestratorProfile: &OrchestratorProfile{
					OrchestratorType: DCOS,
				},
			},
			expectedIsDCOS:       true,
			expectedIsKubernetes: false,
			expectedIsOpenShift:  false,
			expectedIsSwarmMode:  false,
		},
		{
			p: Properties{
				OrchestratorProfile: &OrchestratorProfile{
					OrchestratorType: Kubernetes,
				},
			},
			expectedIsDCOS:       false,
			expectedIsKubernetes: true,
			expectedIsOpenShift:  false,
			expectedIsSwarmMode:  false,
		},
		{
			p: Properties{
				OrchestratorProfile: &OrchestratorProfile{
					OrchestratorType: OpenShift,
				},
			},
			expectedIsDCOS:       false,
			expectedIsKubernetes: false,
			expectedIsOpenShift:  true,
			expectedIsSwarmMode:  false,
		},
		{
			p: Properties{
				OrchestratorProfile: &OrchestratorProfile{
					OrchestratorType: SwarmMode,
				},
			},
			expectedIsDCOS:       false,
			expectedIsKubernetes: false,
			expectedIsOpenShift:  false,
			expectedIsSwarmMode:  true,
		},
	}

	for _, c := range cases {
		if c.expectedIsDCOS != c.p.OrchestratorProfile.IsDCOS() {
			t.Fatalf("Expected IsDCOS() to be %t with OrchestratorType=%s", c.expectedIsDCOS, c.p.OrchestratorProfile.OrchestratorType)
		}
		if c.expectedIsKubernetes != c.p.OrchestratorProfile.IsKubernetes() {
			t.Fatalf("Expected IsKubernetes() to be %t with OrchestratorType=%s", c.expectedIsKubernetes, c.p.OrchestratorProfile.OrchestratorType)
		}
		if c.expectedIsOpenShift != c.p.OrchestratorProfile.IsOpenShift() {
			t.Fatalf("Expected IsOpenShift() to be %t with OrchestratorType=%s", c.expectedIsOpenShift, c.p.OrchestratorProfile.OrchestratorType)
		}
		if c.expectedIsSwarmMode != c.p.OrchestratorProfile.IsSwarmMode() {
			t.Fatalf("Expected IsSwarmMode() to be %t with OrchestratorType=%s", c.expectedIsSwarmMode, c.p.OrchestratorProfile.OrchestratorType)
		}
		if c.expectedIsOpenShift && !c.p.HasStorageAccountDisks() {
			t.Fatalf("Expected HasStorageAccountDisks() to return true when OrchestratorType is OpenShift")
		}
	}
}

func TestWindowsProfile(t *testing.T) {
	w := WindowsProfile{}

	if w.HasSecrets() || w.HasCustomImage() {
		t.Fatalf("Expected HasSecrets() and HasCustomImage() to return false when WindowsProfile is empty")
	}

	w = WindowsProfile{
		Secrets: []KeyVaultSecrets{
			{
				SourceVault: &KeyVaultID{"testVault"},
				VaultCertificates: []KeyVaultCertificate{
					{
						CertificateURL:   "testURL",
						CertificateStore: "testStore",
					},
				},
			},
		},
		WindowsImageSourceURL: "testCustomImage",
	}

	if !(w.HasSecrets() && w.HasCustomImage()) {
		t.Fatalf("Expected HasSecrets() and HasCustomImage() to return true")
	}
}

func TestLinuxProfile(t *testing.T) {
	l := LinuxProfile{}

	if l.HasSecrets() || l.HasSearchDomain() || l.HasCustomNodesDNS() {
		t.Fatalf("Expected HasSecrets(), HasSearchDomain() and HasCustomNodesDNS() to return false when LinuxProfile is empty")
	}

	l = LinuxProfile{
		Secrets: []KeyVaultSecrets{
			{
				SourceVault: &KeyVaultID{"testVault"},
				VaultCertificates: []KeyVaultCertificate{
					{
						CertificateURL:   "testURL",
						CertificateStore: "testStore",
					},
				},
			},
		},
		CustomNodesDNS: &CustomNodesDNS{
			DNSServer: "testDNSServer",
		},
		CustomSearchDomain: &CustomSearchDomain{
			Name:          "testName",
			RealmPassword: "testRealmPassword",
			RealmUser:     "testRealmUser",
		},
	}

	if !(l.HasSecrets() && l.HasSearchDomain() && l.HasCustomNodesDNS()) {
		t.Fatalf("Expected HasSecrets(), HasSearchDomain() and HasCustomNodesDNS() to return true")
	}
}

func TestGetAPIServerEtcdAPIVersion(t *testing.T) {
	o := OrchestratorProfile{}

	if o.GetAPIServerEtcdAPIVersion() != "" {
		t.Fatalf("Expected GetAPIServerEtcdAPIVersion() to return \"\" but instead got %s", o.GetAPIServerEtcdAPIVersion())
	}

	o.KubernetesConfig = &KubernetesConfig{
		EtcdVersion: "3.2.1",
	}

	if o.GetAPIServerEtcdAPIVersion() != "etcd3" {
		t.Fatalf("Expected GetAPIServerEtcdAPIVersion() to return \"etcd3\" but instead got %s", o.GetAPIServerEtcdAPIVersion())
	}

	// invalid version string
	o.KubernetesConfig.EtcdVersion = "2.3.8"
	if o.GetAPIServerEtcdAPIVersion() != "etcd2" {
		t.Fatalf("Expected GetAPIServerEtcdAPIVersion() to return \"etcd2\" but instead got %s", o.GetAPIServerEtcdAPIVersion())
	}
}

func TestHasAadProfile(t *testing.T) {
	p := Properties{}

	if p.HasAadProfile() {
		t.Fatalf("Expected HasAadProfile() to return false")
	}

	p.AADProfile = &AADProfile{
		ClientAppID: "test",
		ServerAppID: "test",
	}

	if !p.HasAadProfile() {
		t.Fatalf("Expected HasAadProfile() to return true")
	}

}

func TestCustomHyperkubeImageField(t *testing.T) {
	log.Println(exampleAPIModel)
	apiloader := &Apiloader{
		Translator: nil,
	}
	apimodel, _, err := apiloader.DeserializeContainerService([]byte(exampleAPIModel), false, false, nil)
	if err != nil {
		t.Fatalf("unexpectedly error deserializing the example apimodel: %s", err)
	}

	actualCustomHyperkubeImage := apimodel.Properties.OrchestratorProfile.KubernetesConfig.CustomHyperkubeImage
	if actualCustomHyperkubeImage != exampleCustomHyperkubeImage {
		t.Fatalf("kubernetesConfig->customHyperkubeImage field value was unexpected: got(%s), expected(%s)", actualCustomHyperkubeImage, exampleCustomHyperkubeImage)
	}
}

func TestUserAssignedMSI(t *testing.T) {
	// Test1: With just System MSI
	log.Println(exampleSystemMSIModel)
	apiloader := &Apiloader{
		Translator: nil,
	}
	apiModel, _, err := apiloader.DeserializeContainerService([]byte(exampleSystemMSIModel), false, false, nil)
	if err != nil {
		t.Fatalf("unexpected error deserailizing the example user msi api model: %s", err)
	}
	systemMSI := apiModel.Properties.OrchestratorProfile.KubernetesConfig.UseManagedIdentity
	actualUserMSI := apiModel.Properties.OrchestratorProfile.KubernetesConfig.UserAssignedID
	if !systemMSI || actualUserMSI != "" {
		t.Fatalf("found user msi: %t and usermsi: %s", systemMSI, actualUserMSI)
	}

	// Test2: With user assigned MSI
	log.Println(exampleUserMSIModel)
	apiloader = &Apiloader{
		Translator: nil,
	}
	apiModel, _, err = apiloader.DeserializeContainerService([]byte(exampleUserMSIModel), false, false, nil)
	if err != nil {
		t.Fatalf("unexpected error deserailizing the example user msi api model: %s", err)
	}
	systemMSI = apiModel.Properties.OrchestratorProfile.KubernetesConfig.UseManagedIdentity
	actualUserMSI = apiModel.Properties.OrchestratorProfile.KubernetesConfig.UserAssignedID
	if !systemMSI && actualUserMSI != exampleUserMSI {
		t.Fatalf("found user msi: %t and usermsi: %s", systemMSI, actualUserMSI)
	}
}

func TestKubernetesAddon(t *testing.T) {
	addon := getMockAddon("addon")
	if !addon.IsEnabled(true) {
		t.Fatalf("KubernetesAddon.IsEnabled(true) should always return true when Enabled property is not specified")
	}

	if addon.IsEnabled(false) {
		t.Fatalf("KubernetesAddon.IsEnabled(false) should always return false when Enabled property is not specified")
	}
	e := true
	addon.Enabled = &e
	if !addon.IsEnabled(false) {
		t.Fatalf("KubernetesAddon.IsEnabled(false) should always return true when Enabled property is set to true")
	}
	if !addon.IsEnabled(true) {
		t.Fatalf("KubernetesAddon.IsEnabled(true) should always return true when Enabled property is set to true")
	}
	e = false
	addon.Enabled = &e
	if addon.IsEnabled(false) {
		t.Fatalf("KubernetesAddon.IsEnabled(false) should always return false when Enabled property is set to false")
	}
	if addon.IsEnabled(true) {
		t.Fatalf("KubernetesAddon.IsEnabled(true) should always return false when Enabled property is set to false")
	}
}

func TestIsTillerEnabled(t *testing.T) {
	c := KubernetesConfig{
		Addons: []KubernetesAddon{
			getMockAddon("addon"),
		},
	}
	enabled := c.IsTillerEnabled()
	enabledDefault := DefaultTillerAddonEnabled
	if enabled != enabledDefault {
		t.Fatalf("KubernetesConfig.IsTillerEnabled() should return %t when no tiller addon has been specified, instead returned %t", enabledDefault, enabled)
	}
	c.Addons = append(c.Addons, getMockAddon(DefaultTillerAddonName))
	enabled = c.IsTillerEnabled()
	if !enabled {
		t.Fatalf("KubernetesConfig.IsTillerEnabled() should return true when a custom tiller addon has been specified, instead returned %t", enabled)
	}
	b := false
	c = KubernetesConfig{
		Addons: []KubernetesAddon{
			{
				Name:    DefaultTillerAddonName,
				Enabled: &b,
			},
		},
	}
	enabled = c.IsTillerEnabled()
	if enabled {
		t.Fatalf("KubernetesConfig.IsTillerEnabled() should return false when a custom tiller addon has been specified as disabled, instead returned %t", enabled)
	}
}

func TestIsAADPodIdentityEnabled(t *testing.T) {
	c := KubernetesConfig{
		Addons: []KubernetesAddon{
			getMockAddon("addon"),
		},
	}
	enabled := c.IsAADPodIdentityEnabled()
	enabledDefault := DefaultAADPodIdentityAddonEnabled
	if enabled != enabledDefault {
		t.Fatalf("KubernetesConfig.IsAADPodIdentityEnabled() should return %t when no aad pod identity addon has been specified, instead returned %t", enabledDefault, enabled)
	}
	c.Addons = append(c.Addons, getMockAddon(DefaultAADPodIdentityAddonName))
	enabled = c.IsAADPodIdentityEnabled()
	if enabled {
		t.Fatalf("KubernetesConfig.IsAADPodIdentityEnabled() should return true when aad pod identity addon has been specified, instead returned %t", enabled)
	}
	b := true
	c = KubernetesConfig{
		Addons: []KubernetesAddon{
			{
				Name:    DefaultAADPodIdentityAddonName,
				Enabled: &b,
			},
		},
	}
	enabled = c.IsAADPodIdentityEnabled()
	if !enabled {
		t.Fatalf("KubernetesConfig.IsAADPodIdentityEnabled() should return false when aad pod identity addon has been specified as disabled, instead returned %t", enabled)
	}
}

func TestIsACIConnectorEnabled(t *testing.T) {
	c := KubernetesConfig{
		Addons: []KubernetesAddon{
			getMockAddon("addon"),
		},
	}
	enabled := c.IsACIConnectorEnabled()
	enabledDefault := DefaultACIConnectorAddonEnabled
	if enabled != enabledDefault {
		t.Fatalf("KubernetesConfig.IsACIConnectorEnabled() should return %t when no ACI connector addon has been specified, instead returned %t", enabledDefault, enabled)
	}
	c.Addons = append(c.Addons, getMockAddon(DefaultACIConnectorAddonName))
	enabled = c.IsACIConnectorEnabled()
	if enabled {
		t.Fatalf("KubernetesConfig.IsACIConnectorEnabled() should return true when ACI connector has been specified, instead returned %t", enabled)
	}
	b := true
	c = KubernetesConfig{
		Addons: []KubernetesAddon{
			{
				Name:    DefaultACIConnectorAddonName,
				Enabled: &b,
			},
		},
	}
	enabled = c.IsACIConnectorEnabled()
	if !enabled {
		t.Fatalf("KubernetesConfig.IsACIConnectorEnabled() should return false when ACI connector addon has been specified as disabled, instead returned %t", enabled)
	}
}

func TestIsClusterAutoscalerEnabled(t *testing.T) {
	c := KubernetesConfig{
		Addons: []KubernetesAddon{
			getMockAddon("addon"),
		},
	}
	enabled := c.IsClusterAutoscalerEnabled()
	enabledDefault := DefaultClusterAutoscalerAddonEnabled
	if enabled != enabledDefault {
		t.Fatalf("KubernetesConfig.IsClusterAutoscalerEnabled() should return %t when no cluster autoscaler addon has been specified, instead returned %t", enabledDefault, enabled)
	}
	c.Addons = append(c.Addons, getMockAddon(DefaultClusterAutoscalerAddonName))
	enabled = c.IsClusterAutoscalerEnabled()
	if enabled {
		t.Fatalf("KubernetesConfig.IsClusterAutoscalerEnabled() should return true when cluster autoscaler has been specified, instead returned %t", enabled)
	}
	b := true
	c = KubernetesConfig{
		Addons: []KubernetesAddon{
			{
				Name:    DefaultClusterAutoscalerAddonName,
				Enabled: &b,
			},
		},
	}
	enabled = c.IsClusterAutoscalerEnabled()
	if !enabled {
		t.Fatalf("KubernetesConfig.IsClusterAutoscalerEnabled() should return false when cluster autoscaler addon has been specified as disabled, instead returned %t", enabled)
	}
}

func TestIsBlobfuseFlexVolumeEnabled(t *testing.T) {
	c := KubernetesConfig{
		Addons: []KubernetesAddon{
			getMockAddon("addon"),
		},
	}
	enabled := c.IsBlobfuseFlexVolumeEnabled()
	enabledDefault := DefaultBlobfuseFlexVolumeAddonEnabled
	if enabled != enabledDefault {
		t.Fatalf("KubernetesConfig.IsBlobfuseFlexVolumeEnabled() should return %t when no blob fuse flex volume addon has been specified, instead returned %t", enabledDefault, enabled)
	}
	c.Addons = append(c.Addons, getMockAddon(DefaultBlobfuseFlexVolumeAddonName))
	enabled = c.IsBlobfuseFlexVolumeEnabled()
	if !enabled {
		t.Fatalf("KubernetesConfig.IsBlobfuseFlexVolumeEnabled() should return true when blob fuse flex volume has been specified, instead returned %t", enabled)
	}
	b := true
	c = KubernetesConfig{
		Addons: []KubernetesAddon{
			{
				Name:    DefaultBlobfuseFlexVolumeAddonName,
				Enabled: &b,
			},
		},
	}
	enabled = c.IsBlobfuseFlexVolumeEnabled()
	if !enabled {
		t.Fatalf("KubernetesConfig.IsBlobfuseFlexVolumeEnabled() should return false when blob fuse flex volume addon has been specified as disabled, instead returned %t", enabled)
	}
}

func TestIsSMBFlexVolumeEnabled(t *testing.T) {
	c := KubernetesConfig{
		Addons: []KubernetesAddon{
			getMockAddon("addon"),
		},
	}
	enabled := c.IsSMBFlexVolumeEnabled()
	enabledDefault := DefaultSMBFlexVolumeAddonEnabled
	if enabled != enabledDefault {
		t.Fatalf("KubernetesConfig.IsSMBFlexVolumeEnabled() should return %t when no SMB flex volume addon has been specified, instead returned %t", enabledDefault, enabled)
	}
	c.Addons = append(c.Addons, getMockAddon(DefaultSMBFlexVolumeAddonName))
	enabled = c.IsSMBFlexVolumeEnabled()
	if !enabled {
		t.Fatalf("KubernetesConfig.IsSMBFlexVolumeEnabled() should return true when SMB flex volume has been specified, instead returned %t", enabled)
	}
	b := true
	c = KubernetesConfig{
		Addons: []KubernetesAddon{
			{
				Name:    DefaultSMBFlexVolumeAddonName,
				Enabled: &b,
			},
		},
	}
	enabled = c.IsSMBFlexVolumeEnabled()
	if !enabled {
		t.Fatalf("KubernetesConfig.IsSMBFlexVolumeEnabled() should return false when SMB flex volume addon has been specified as disabled, instead returned %t", enabled)
	}
}

func TestIsKeyVaultFlexVolumeEnabled(t *testing.T) {
	c := KubernetesConfig{
		Addons: []KubernetesAddon{
			getMockAddon("addon"),
		},
	}
	enabled := c.IsKeyVaultFlexVolumeEnabled()
	enabledDefault := DefaultKeyVaultFlexVolumeAddonEnabled
	if enabled != enabledDefault {
		t.Fatalf("KubernetesConfig.IsKeyVaultFlexVolumeEnabled() should return %t when no key vault flex volume addon has been specified, instead returned %t", enabledDefault, enabled)
	}
	c.Addons = append(c.Addons, getMockAddon(DefaultKeyVaultFlexVolumeAddonName))
	enabled = c.IsKeyVaultFlexVolumeEnabled()
	if !enabled {
		t.Fatalf("KubernetesConfig.IsKeyVaultFlexVolumeEnabled() should return true when no keyvault flex volume has been specified, instead returned %t", enabled)
	}
	b := false
	c = KubernetesConfig{
		Addons: []KubernetesAddon{
			{
				Name:    DefaultKeyVaultFlexVolumeAddonName,
				Enabled: &b,
			},
		},
	}
	enabled = c.IsKeyVaultFlexVolumeEnabled()
	if enabled {
		t.Fatalf("KubernetesConfig.IsKeyVaultFlexVolumeEnabled() should return false when keyvault flex volume addon has been specified as disabled, instead returned %t", enabled)
	}
}

func TestIsNVIDIADevicePluginEnabled(t *testing.T) {
	p := Properties{
		AgentPoolProfiles: []*AgentPoolProfile{
			{
				Name:   "agentpool",
				VMSize: "Standard_N",
				Count:  1,
			},
		},
		OrchestratorProfile: &OrchestratorProfile{
			OrchestratorType:    Kubernetes,
			OrchestratorVersion: "1.9.0",
			KubernetesConfig: &KubernetesConfig{
				Addons: []KubernetesAddon{
					getMockAddon("addon"),
				},
			},
		},
	}

	if !IsNSeriesSKU(&p) {
		t.Fatalf("IsNSeriesSKU should return true when explicitly using VM Size %s", p.AgentPoolProfiles[0].VMSize)
	}
	if p.IsNVIDIADevicePluginEnabled() {
		t.Fatalf("KubernetesConfig.IsNVIDIADevicePluginEnabled() should return false with N-series VMs with < k8s 1.10, instead returned %t", p.IsNVIDIADevicePluginEnabled())
	}

	p.OrchestratorProfile.OrchestratorVersion = "1.10.0"
	if !p.IsNVIDIADevicePluginEnabled() {
		t.Fatalf("KubernetesConfig.IsNVIDIADevicePluginEnabled() should return true with N-series VMs with k8s >= 1.10, instead returned %t", p.IsNVIDIADevicePluginEnabled())
	}

	p.AgentPoolProfiles[0].VMSize = "Standard_D2_v2"
	p.OrchestratorProfile.KubernetesConfig.Addons = []KubernetesAddon{
		{
			Name:    NVIDIADevicePluginAddonName,
			Enabled: helpers.PointerToBool(false),
		},
	}

	if IsNSeriesSKU(&p) {
		t.Fatalf("IsNSeriesSKU should return false when explicitly using VM Size %s", p.AgentPoolProfiles[0].VMSize)
	}
	if p.IsNVIDIADevicePluginEnabled() {
		t.Fatalf("KubernetesConfig.IsNVIDIADevicePluginEnabled() should return false when explicitly disabled")
	}
}

func TestIsContainerMonitoringEnabled(t *testing.T) {
	v := "1.9.0"
	o := OrchestratorProfile{
		OrchestratorType:    "Kubernetes",
		OrchestratorVersion: v,
		KubernetesConfig: &KubernetesConfig{Addons: []KubernetesAddon{
			getMockAddon("addon"),
		},
		},
	}
	enabled := o.KubernetesConfig.IsContainerMonitoringEnabled()
	enabledDefault := DefaultContainerMonitoringAddonEnabled
	if enabled != enabledDefault {
		t.Fatalf("KubernetesConfig.IsContainerMonitoringEnabled() should return %t for kubernetes version %s when no container-monitoring addon has been specified, instead returned %t", enabledDefault, v, enabled)
	}

	b := true
	cm := getMockAddon(ContainerMonitoringAddonName)
	cm.Enabled = &b
	o.KubernetesConfig.Addons = append(o.KubernetesConfig.Addons, cm)
	enabled = o.KubernetesConfig.IsContainerMonitoringEnabled()
	if !enabled {
		t.Fatalf("KubernetesConfig.IsContainerMonitoringEnabled() should return %t for kubernetes version %s when the container-monitoring addon has been specified, instead returned %t", true, v, enabled)
	}

	b = false
	o = OrchestratorProfile{
		OrchestratorType:    "Kubernetes",
		OrchestratorVersion: v,
		KubernetesConfig: &KubernetesConfig{Addons: []KubernetesAddon{
			{
				Name:    ContainerMonitoringAddonName,
				Enabled: &b,
			},
		},
		},
	}
	enabled = o.KubernetesConfig.IsContainerMonitoringEnabled()
	if enabled {
		t.Fatalf("KubernetesConfig.IsContainerMonitoringEnabled() should return false when a custom container monitoring addon has been specified as disabled, instead returned %t", enabled)
	}
}

func TestIsDashboardEnabled(t *testing.T) {
	c := KubernetesConfig{
		Addons: []KubernetesAddon{
			getMockAddon("addon"),
		},
	}
	enabled := c.IsDashboardEnabled()
	enabledDefault := DefaultDashboardAddonEnabled
	if enabled != enabledDefault {
		t.Fatalf("KubernetesConfig.IsDashboardEnabled() should return %t when no kubernetes-dashboard addon has been specified, instead returned %t", enabledDefault, enabled)
	}
	c.Addons = append(c.Addons, getMockAddon(DefaultDashboardAddonName))
	enabled = c.IsDashboardEnabled()
	if !enabled {
		t.Fatalf("KubernetesConfig.IsDashboardEnabled() should return true when a custom kubernetes-dashboard addon has been specified, instead returned %t", enabled)
	}
	b := false
	c = KubernetesConfig{
		Addons: []KubernetesAddon{
			{
				Name:    DefaultDashboardAddonName,
				Enabled: &b,
			},
		},
	}
	enabled = c.IsDashboardEnabled()
	if enabled {
		t.Fatalf("KubernetesConfig.IsDashboardEnabled() should return false when a custom kubernetes-dashboard addon has been specified as disabled, instead returned %t", enabled)
	}
}

func TestIsReschedulerEnabled(t *testing.T) {
	c := KubernetesConfig{
		Addons: []KubernetesAddon{
			getMockAddon("addon"),
		},
	}
	enabled := c.IsReschedulerEnabled()
	enabledDefault := DefaultReschedulerAddonEnabled
	if enabled != enabledDefault {
		t.Fatalf("KubernetesConfig.IsReschedulerEnabled() should return %t when no rescheduler addon has been specified, instead returned %t", enabledDefault, enabled)
	}
	c.Addons = append(c.Addons, getMockAddon(DefaultReschedulerAddonName))
	enabled = c.IsReschedulerEnabled()
	if enabled {
		t.Fatalf("KubernetesConfig.IsReschedulerEnabled() should return true when a custom rescheduler addon has been specified, instead returned %t", enabled)
	}
	b := true
	c = KubernetesConfig{
		Addons: []KubernetesAddon{
			{
				Name:    DefaultReschedulerAddonName,
				Enabled: &b,
			},
		},
	}
	enabled = c.IsReschedulerEnabled()
	if !enabled {
		t.Fatalf("KubernetesConfig.IsReschedulerEnabled() should return false when a custom rescheduler addon has been specified as enabled, instead returned %t", enabled)
	}
}

func TestIsMetricsServerEnabled(t *testing.T) {
	v := "1.8.0"
	o := OrchestratorProfile{
		OrchestratorType:    "Kubernetes",
		OrchestratorVersion: v,
		KubernetesConfig: &KubernetesConfig{Addons: []KubernetesAddon{
			getMockAddon("addon"),
		},
		},
	}
	enabled := o.IsMetricsServerEnabled()
	enabledDefault := DefaultMetricsServerAddonEnabled
	if enabled != enabledDefault {
		t.Fatalf("KubernetesConfig.IsMetricsServerEnabled() should return %t for kubernetes version %s when no metrics-server addon has been specified, instead returned %t", enabledDefault, v, enabled)
	}

	o.KubernetesConfig.Addons = append(o.KubernetesConfig.Addons, getMockAddon(DefaultMetricsServerAddonName))
	enabled = o.IsMetricsServerEnabled()
	enabledDefault = DefaultMetricsServerAddonEnabled
	if enabled != enabledDefault {
		t.Fatalf("KubernetesConfig.IsMetricsServerEnabled() should return %t for kubernetes version %s when the metrics-server addon has been specified, instead returned %t", enabledDefault, v, enabled)
	}

	b := true
	o = OrchestratorProfile{
		OrchestratorType:    "Kubernetes",
		OrchestratorVersion: v,
		KubernetesConfig: &KubernetesConfig{Addons: []KubernetesAddon{
			{
				Name:    DefaultMetricsServerAddonName,
				Enabled: &b,
			},
		},
		},
	}
	enabled = o.IsMetricsServerEnabled()
	if !enabled {
		t.Fatalf("KubernetesConfig.IsMetricsServerEnabled() should return true for kubernetes version %s when the metrics-server addon has been specified as enabled, instead returned %t", v, enabled)
	}
}

func TestCloudProviderDefaults(t *testing.T) {
	// Test cloudprovider defaults when no user-provided values
	v := "1.8.0"
	o := OrchestratorProfile{
		OrchestratorType:    "Kubernetes",
		OrchestratorVersion: v,
		KubernetesConfig:    &KubernetesConfig{},
	}
	o.KubernetesConfig.SetCloudProviderBackoffDefaults()
	o.KubernetesConfig.SetCloudProviderRateLimitDefaults()

	intCases := []struct {
		defaultVal  int
		computedVal int
	}{
		{
			defaultVal:  DefaultKubernetesCloudProviderBackoffRetries,
			computedVal: o.KubernetesConfig.CloudProviderBackoffRetries,
		},
		{
			defaultVal:  DefaultKubernetesCloudProviderBackoffDuration,
			computedVal: o.KubernetesConfig.CloudProviderBackoffDuration,
		},
		{
			defaultVal:  DefaultKubernetesCloudProviderRateLimitBucket,
			computedVal: o.KubernetesConfig.CloudProviderRateLimitBucket,
		},
	}

	for _, c := range intCases {
		if c.computedVal != c.defaultVal {
			t.Fatalf("KubernetesConfig empty cloudprovider configs should reflect default values after SetCloudProviderBackoffDefaults(), expected %d, got %d", c.defaultVal, c.computedVal)
		}
	}

	floatCases := []struct {
		defaultVal  float64
		computedVal float64
	}{
		{
			defaultVal:  DefaultKubernetesCloudProviderBackoffJitter,
			computedVal: o.KubernetesConfig.CloudProviderBackoffJitter,
		},
		{
			defaultVal:  DefaultKubernetesCloudProviderBackoffExponent,
			computedVal: o.KubernetesConfig.CloudProviderBackoffExponent,
		},
		{
			defaultVal:  DefaultKubernetesCloudProviderRateLimitQPS,
			computedVal: o.KubernetesConfig.CloudProviderRateLimitQPS,
		},
	}

	for _, c := range floatCases {
		if c.computedVal != c.defaultVal {
			t.Fatalf("KubernetesConfig empty cloudprovider configs should reflect default values after SetCloudProviderBackoffDefaults(), expected %f, got %f", c.defaultVal, c.computedVal)
		}
	}

	customCloudProviderBackoffDuration := 99
	customCloudProviderBackoffExponent := 10.0
	customCloudProviderBackoffJitter := 11.9
	customCloudProviderBackoffRetries := 9
	customCloudProviderRateLimitBucket := 37
	customCloudProviderRateLimitQPS := 9.9

	// Test cloudprovider defaults when user provides configuration
	v = "1.8.0"
	o = OrchestratorProfile{
		OrchestratorType:    "Kubernetes",
		OrchestratorVersion: v,
		KubernetesConfig: &KubernetesConfig{
			CloudProviderBackoffDuration: customCloudProviderBackoffDuration,
			CloudProviderBackoffExponent: customCloudProviderBackoffExponent,
			CloudProviderBackoffJitter:   customCloudProviderBackoffJitter,
			CloudProviderBackoffRetries:  customCloudProviderBackoffRetries,
			CloudProviderRateLimitBucket: customCloudProviderRateLimitBucket,
			CloudProviderRateLimitQPS:    customCloudProviderRateLimitQPS,
		},
	}
	o.KubernetesConfig.SetCloudProviderBackoffDefaults()
	o.KubernetesConfig.SetCloudProviderRateLimitDefaults()

	intCasesCustom := []struct {
		customVal   int
		computedVal int
	}{
		{
			customVal:   customCloudProviderBackoffRetries,
			computedVal: o.KubernetesConfig.CloudProviderBackoffRetries,
		},
		{
			customVal:   customCloudProviderBackoffDuration,
			computedVal: o.KubernetesConfig.CloudProviderBackoffDuration,
		},
		{
			customVal:   customCloudProviderRateLimitBucket,
			computedVal: o.KubernetesConfig.CloudProviderRateLimitBucket,
		},
	}

	for _, c := range intCasesCustom {
		if c.computedVal != c.customVal {
			t.Fatalf("KubernetesConfig empty cloudprovider configs should reflect default values after SetCloudProviderBackoffDefaults(), expected %d, got %d", c.customVal, c.computedVal)
		}
	}

	floatCasesCustom := []struct {
		customVal   float64
		computedVal float64
	}{
		{
			customVal:   customCloudProviderBackoffJitter,
			computedVal: o.KubernetesConfig.CloudProviderBackoffJitter,
		},
		{
			customVal:   customCloudProviderBackoffExponent,
			computedVal: o.KubernetesConfig.CloudProviderBackoffExponent,
		},
		{
			customVal:   customCloudProviderRateLimitQPS,
			computedVal: o.KubernetesConfig.CloudProviderRateLimitQPS,
		},
	}

	for _, c := range floatCasesCustom {
		if c.computedVal != c.customVal {
			t.Fatalf("KubernetesConfig empty cloudprovider configs should reflect default values after SetCloudProviderBackoffDefaults(), expected %f, got %f", c.customVal, c.computedVal)
		}
	}

	// Test cloudprovider defaults when user provides *some* config values
	v = "1.8.0"
	o = OrchestratorProfile{
		OrchestratorType:    "Kubernetes",
		OrchestratorVersion: v,
		KubernetesConfig: &KubernetesConfig{
			CloudProviderBackoffDuration: customCloudProviderBackoffDuration,
			CloudProviderRateLimitBucket: customCloudProviderRateLimitBucket,
			CloudProviderRateLimitQPS:    customCloudProviderRateLimitQPS,
		},
	}
	o.KubernetesConfig.SetCloudProviderBackoffDefaults()
	o.KubernetesConfig.SetCloudProviderRateLimitDefaults()

	intCasesMixed := []struct {
		expectedVal int
		computedVal int
	}{
		{
			expectedVal: DefaultKubernetesCloudProviderBackoffRetries,
			computedVal: o.KubernetesConfig.CloudProviderBackoffRetries,
		},
		{
			expectedVal: customCloudProviderBackoffDuration,
			computedVal: o.KubernetesConfig.CloudProviderBackoffDuration,
		},
		{
			expectedVal: customCloudProviderRateLimitBucket,
			computedVal: o.KubernetesConfig.CloudProviderRateLimitBucket,
		},
	}

	for _, c := range intCasesMixed {
		if c.computedVal != c.expectedVal {
			t.Fatalf("KubernetesConfig empty cloudprovider configs should reflect default values after SetCloudProviderBackoffDefaults(), expected %d, got %d", c.expectedVal, c.computedVal)
		}
	}

	floatCasesMixed := []struct {
		expectedVal float64
		computedVal float64
	}{
		{
			expectedVal: DefaultKubernetesCloudProviderBackoffJitter,
			computedVal: o.KubernetesConfig.CloudProviderBackoffJitter,
		},
		{
			expectedVal: DefaultKubernetesCloudProviderBackoffExponent,
			computedVal: o.KubernetesConfig.CloudProviderBackoffExponent,
		},
		{
			expectedVal: customCloudProviderRateLimitQPS,
			computedVal: o.KubernetesConfig.CloudProviderRateLimitQPS,
		},
	}

	for _, c := range floatCasesMixed {
		if c.computedVal != c.expectedVal {
			t.Fatalf("KubernetesConfig empty cloudprovider configs should reflect default values after SetCloudProviderBackoffDefaults(), expected %f, got %f", c.expectedVal, c.computedVal)
		}
	}
}

func getMockAddon(name string) KubernetesAddon {
	return KubernetesAddon{
		Name: name,
		Containers: []KubernetesContainerSpec{
			{
				Name:           name,
				CPURequests:    "50m",
				MemoryRequests: "150Mi",
				CPULimits:      "50m",
				MemoryLimits:   "150Mi",
			},
		},
	}
}

func TestAreAgentProfilesCustomVNET(t *testing.T) {
	p := Properties{}
	p.AgentPoolProfiles = []*AgentPoolProfile{
		{
			VnetSubnetID: "subnetlink1",
		},
		{
			VnetSubnetID: "subnetlink2",
		},
	}

	if !p.AreAgentProfilesCustomVNET() {
		t.Fatalf("Expected isCustomVNET to be true when subnet exists for all agent pool profile")
	}

	p.AgentPoolProfiles = []*AgentPoolProfile{
		{
			VnetSubnetID: "subnetlink1",
		},
		{
			VnetSubnetID: "",
		},
	}

	if p.AreAgentProfilesCustomVNET() {
		t.Fatalf("Expected isCustomVNET to be false when subnet exists for some agent pool profile")
	}

	p.AgentPoolProfiles = nil

	if p.AreAgentProfilesCustomVNET() {
		t.Fatalf("Expected isCustomVNET to be false when agent pool profiles is nil")
	}
}

func TestGenerateClusterID(t *testing.T) {
	p := &Properties{
		AgentPoolProfiles: []*AgentPoolProfile{
			{
				Name: "foo_agent",
			},
		},
	}

	clusterID := p.GetClusterID()

	r, err := regexp.Compile("[0-9]{8}")

	if err != nil {
		t.Errorf("unexpected error while parsing regex : %s", err.Error())
	}

	if !r.MatchString(clusterID) {
		t.Fatal("ClusterID should be an 8 digit integer string")
	}

	p = &Properties{
		HostedMasterProfile: &HostedMasterProfile{
			DNSPrefix: "foodnsprefx",
		},
	}

	clusterID = p.GetClusterID()

	r, err = regexp.Compile("[0-9]{8}")

	if err != nil {
		t.Errorf("unexpected error while parsing regex : %s", err.Error())
	}

	if !r.MatchString(clusterID) {
		t.Fatal("ClusterID should be an 8 digit integer string")
	}
}

func TestGetPrimaryAvailabilitySetName(t *testing.T) {
	p := &Properties{
		OrchestratorProfile: &OrchestratorProfile{
			OrchestratorType: Kubernetes,
		},
		MasterProfile: &MasterProfile{
			Count:     1,
			DNSPrefix: "foo",
			VMSize:    "Standard_DS2_v2",
		},
		AgentPoolProfiles: []*AgentPoolProfile{
			{
				Name:                "agentpool",
				VMSize:              "Standard_D2_v2",
				Count:               1,
				AvailabilityProfile: AvailabilitySet,
			},
		},
	}

	expected := "agentpool-availabilitySet-28513887"
	got := p.GetPrimaryAvailabilitySetName()
	if got != expected {
		t.Errorf("expected primary availability set name %s, but got %s", expected, got)
	}
}

func TestGetPrimaryScaleSetName(t *testing.T) {
	p := &Properties{
		OrchestratorProfile: &OrchestratorProfile{
			OrchestratorType: Kubernetes,
		},
		MasterProfile: &MasterProfile{
			Count:     1,
			DNSPrefix: "foo",
			VMSize:    "Standard_DS2_v2",
		},
		AgentPoolProfiles: []*AgentPoolProfile{
			{
				Name:                "agentpool",
				VMSize:              "Standard_D2_v2",
				Count:               1,
				AvailabilityProfile: VirtualMachineScaleSets,
			},
		},
	}

	expected := "k8s-agentpool-28513887-vmss"
	got := p.GetPrimaryScaleSetName()
	if got != expected {
		t.Errorf("expected primary availability set name %s, but got %s", expected, got)
	}
}

func TestGetRouteTableName(t *testing.T) {
	p := &Properties{
		OrchestratorProfile: &OrchestratorProfile{
			OrchestratorType: Kubernetes,
		},
		HostedMasterProfile: &HostedMasterProfile{
			FQDN:      "fqdn",
			DNSPrefix: "foo",
			Subnet:    "mastersubnet",
		},
		AgentPoolProfiles: []*AgentPoolProfile{
			{
				Name:                "agentpool",
				VMSize:              "Standard_D2_v2",
				Count:               1,
				AvailabilityProfile: VirtualMachineScaleSets,
			},
		},
	}

	actualRTName := p.GetRouteTableName()
	expectedRTName := "aks-agentpool-28513887-routetable"

	actualNSGName := p.GetNSGName()
	expectedNSGName := "aks-agentpool-28513887-nsg"

	if actualRTName != expectedRTName {
		t.Errorf("expected route table name %s, but got %s", expectedRTName, actualRTName)
	}

	if actualNSGName != expectedNSGName {
		t.Errorf("expected route table name %s, but got %s", expectedNSGName, actualNSGName)
	}

	p = &Properties{
		OrchestratorProfile: &OrchestratorProfile{
			OrchestratorType: Kubernetes,
		},
		MasterProfile: &MasterProfile{
			Count:     1,
			DNSPrefix: "foo",
			VMSize:    "Standard_DS2_v2",
		},
		AgentPoolProfiles: []*AgentPoolProfile{
			{
				Name:                "agentpool",
				VMSize:              "Standard_D2_v2",
				Count:               1,
				AvailabilityProfile: VirtualMachineScaleSets,
			},
		},
	}

	actualRTName = p.GetRouteTableName()
	expectedRTName = "k8s-master-28513887-routetable"

	actualNSGName = p.GetNSGName()
	expectedNSGName = "k8s-master-28513887-nsg"

	if actualRTName != expectedRTName {
		t.Errorf("expected route table name %s, but got %s", actualRTName, expectedRTName)
	}

	if actualNSGName != expectedNSGName {
		t.Errorf("expected route table name %s, but got %s", actualNSGName, expectedNSGName)
	}
}

func TestGetSubnetName(t *testing.T) {
	tests := []struct {
		name               string
		properties         *Properties
		expectedSubnetName string
	}{
		{
			name: "Cluster with HosterMasterProfile",
			properties: &Properties{
				OrchestratorProfile: &OrchestratorProfile{
					OrchestratorType: Kubernetes,
				},
				HostedMasterProfile: &HostedMasterProfile{
					FQDN:      "fqdn",
					DNSPrefix: "foo",
					Subnet:    "mastersubnet",
				},
				AgentPoolProfiles: []*AgentPoolProfile{
					{
						Name:                "agentpool",
						VMSize:              "Standard_D2_v2",
						Count:               1,
						AvailabilityProfile: VirtualMachineScaleSets,
					},
				},
			},
			expectedSubnetName: "aks-subnet",
		},
		{
			name: "Cluster with HosterMasterProfile and custom VNET",
			properties: &Properties{
				OrchestratorProfile: &OrchestratorProfile{
					OrchestratorType: Kubernetes,
				},
				HostedMasterProfile: &HostedMasterProfile{
					FQDN:      "fqdn",
					DNSPrefix: "foo",
					Subnet:    "mastersubnet",
				},
				AgentPoolProfiles: []*AgentPoolProfile{
					{
						Name:                "agentpool",
						VMSize:              "Standard_D2_v2",
						Count:               1,
						AvailabilityProfile: VirtualMachineScaleSets,
						VnetSubnetID:        "/subscriptions/SUBSCRIPTION_ID/resourceGroups/RESOURCE_GROUP_NAME/providers/Microsoft.Network/virtualNetworks/ExampleCustomVNET/subnets/BazAgentSubnet",
					},
				},
			},
			expectedSubnetName: "BazAgentSubnet",
		},
		{
			name: "Cluster with MasterProfile",
			properties: &Properties{
				OrchestratorProfile: &OrchestratorProfile{
					OrchestratorType: Kubernetes,
				},
				MasterProfile: &MasterProfile{
					Count:     1,
					DNSPrefix: "foo",
					VMSize:    "Standard_DS2_v2",
				},
				AgentPoolProfiles: []*AgentPoolProfile{
					{
						Name:                "agentpool",
						VMSize:              "Standard_D2_v2",
						Count:               1,
						AvailabilityProfile: VirtualMachineScaleSets,
					},
				},
			},
			expectedSubnetName: "k8s-subnet",
		},
		{
			name: "Cluster with MasterProfile and custom VNET",
			properties: &Properties{
				OrchestratorProfile: &OrchestratorProfile{
					OrchestratorType: Kubernetes,
				},
				MasterProfile: &MasterProfile{
					Count:        1,
					DNSPrefix:    "foo",
					VMSize:       "Standard_DS2_v2",
					VnetSubnetID: "/subscriptions/SUBSCRIPTION_ID/resourceGroups/RESOURCE_GROUP_NAME/providers/Microsoft.Network/virtualNetworks/ExampleCustomVNET/subnets/BazAgentSubnet",
				},
				AgentPoolProfiles: []*AgentPoolProfile{
					{
						Name:                "agentpool",
						VMSize:              "Standard_D2_v2",
						Count:               1,
						AvailabilityProfile: VirtualMachineScaleSets,
					},
				},
			},
			expectedSubnetName: "BazAgentSubnet",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			actual := test.properties.GetSubnetName()

			if actual != test.expectedSubnetName {
				t.Errorf("expected subnet name %s, but got %s", test.expectedSubnetName, actual)
			}
		})
	}
}

func TestProperties_GetVirtualNetworkName(t *testing.T) {
	tests := []struct {
		name                       string
		properties                 *Properties
		expectedVirtualNetworkName string
	}{
		{
			name: "Cluster with HostedMasterProfile and Custom VNET AgentProfiles",
			properties: &Properties{
				HostedMasterProfile: &HostedMasterProfile{
					FQDN:      "fqdn",
					DNSPrefix: "foo",
					Subnet:    "mastersubnet",
				},
				AgentPoolProfiles: []*AgentPoolProfile{
					{
						Name:                "agentpool",
						VMSize:              "Standard_D2_v2",
						Count:               1,
						AvailabilityProfile: VirtualMachineScaleSets,
						VnetSubnetID:        "/subscriptions/SUBSCRIPTION_ID/resourceGroups/RESOURCE_GROUP_NAME/providers/Microsoft.Network/virtualNetworks/ExampleCustomVNET/subnets/BazAgentSubnet",
					},
				},
			},
			expectedVirtualNetworkName: "ExampleCustomVNET",
		},
		{
			name: "Cluster with HostedMasterProfile and AgentProfiles",
			properties: &Properties{
				OrchestratorProfile: &OrchestratorProfile{
					OrchestratorType: Kubernetes,
				},
				HostedMasterProfile: &HostedMasterProfile{
					FQDN:      "fqdn",
					DNSPrefix: "foo",
					Subnet:    "mastersubnet",
				},
				AgentPoolProfiles: []*AgentPoolProfile{
					{
						Name:                "agentpool",
						VMSize:              "Standard_D2_v2",
						Count:               1,
						AvailabilityProfile: VirtualMachineScaleSets,
					},
				},
			},
			expectedVirtualNetworkName: "aks-vnet-28513887",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			actual := test.properties.GetVirtualNetworkName()

			if actual != test.expectedVirtualNetworkName {
				t.Errorf("expected virtual network name %s, but got %s", test.expectedVirtualNetworkName, actual)
			}
		})
	}
}

func TestProperties_GetVNetResourceGroupName(t *testing.T) {
	p := &Properties{
		HostedMasterProfile: &HostedMasterProfile{
			FQDN:      "fqdn",
			DNSPrefix: "foo",
			Subnet:    "mastersubnet",
		},
		AgentPoolProfiles: []*AgentPoolProfile{
			{
				Name:                "agentpool",
				VMSize:              "Standard_D2_v2",
				Count:               1,
				AvailabilityProfile: VirtualMachineScaleSets,
				VnetSubnetID:        "/subscriptions/SUBSCRIPTION_ID/resourceGroups/RESOURCE_GROUP_NAME/providers/Microsoft.Network/virtualNetworks/ExampleCustomVNET/subnets/BazAgentSubnet",
			},
		},
	}
	expectedVNETResourceGroupName := "RESOURCE_GROUP_NAME"

	actual := p.GetVNetResourceGroupName()

	if expectedVNETResourceGroupName != actual {
		t.Errorf("expected vnet resource group name name %s, but got %s", expectedVNETResourceGroupName, actual)
	}
}

func TestProperties_GetClusterMetadata(t *testing.T) {
	p := &Properties{
		OrchestratorProfile: &OrchestratorProfile{
			OrchestratorType: Kubernetes,
		},
		MasterProfile: &MasterProfile{
			Count:        1,
			DNSPrefix:    "foo",
			VMSize:       "Standard_DS2_v2",
			VnetSubnetID: "/subscriptions/SUBSCRIPTION_ID/resourceGroups/SAMPLE_RESOURCE_GROUP_NAME/providers/Microsoft.Network/virtualNetworks/ExampleCustomVNET/subnets/BazAgentSubnet",
		},
		AgentPoolProfiles: []*AgentPoolProfile{
			{
				Name:                "agentpool",
				VMSize:              "Standard_D2_v2",
				Count:               1,
				AvailabilityProfile: AvailabilitySet,
			},
		},
	}

	metadata := p.GetClusterMetadata()

	if metadata == nil {
		t.Error("did not expect cluster metadata to be nil")
	}

	expectedSubnetName := "BazAgentSubnet"
	if metadata.SubnetName != expectedSubnetName {
		t.Errorf("expected subnet name %s, but got %s", expectedSubnetName, metadata.SubnetName)
	}

	expectedVNetResourceGroupName := "SAMPLE_RESOURCE_GROUP_NAME"
	if metadata.VNetResourceGroupName != expectedVNetResourceGroupName {
		t.Errorf("expected vNetResourceGroupName name %s, but got %s", expectedVNetResourceGroupName, metadata.VNetResourceGroupName)
	}

	expectedVirtualNetworkName := "ExampleCustomVNET"
	if metadata.VirtualNetworkName != expectedVirtualNetworkName {
		t.Errorf("expected VirtualNetworkName name %s, but got %s", expectedVirtualNetworkName, metadata.VirtualNetworkName)
	}

	expectedRouteTableName := "k8s-master-28513887-routetable"
	if metadata.RouteTableName != expectedRouteTableName {
		t.Errorf("expected RouteTableName name %s, but got %s", expectedVirtualNetworkName, metadata.RouteTableName)
	}

	expectedSecurityGroupName := "k8s-master-28513887-nsg"
	if metadata.SecurityGroupName != expectedSecurityGroupName {
		t.Errorf("expected SecurityGroupName name %s, but got %s", expectedSecurityGroupName, metadata.SecurityGroupName)
	}

	expectedPrimaryAvailabilitySetName := "agentpool-availabilitySet-28513887"
	if metadata.PrimaryAvailabilitySetName != expectedPrimaryAvailabilitySetName {
		t.Errorf("expected PrimaryAvailabilitySetName name %s, but got %s", expectedPrimaryAvailabilitySetName, metadata.PrimaryAvailabilitySetName)
	}

	expectedPrimaryScaleSetName := "k8s-agentpool-28513887-vmss"
	if metadata.PrimaryScaleSetName != expectedPrimaryScaleSetName {
		t.Errorf("expected PrimaryScaleSetName name %s, but got %s", expectedPrimaryScaleSetName, metadata.PrimaryScaleSetName)
	}

}

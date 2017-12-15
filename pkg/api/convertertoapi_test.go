package api

import (
	"testing"

	"github.com/Azure/acs-engine/pkg/api/common"
	"github.com/Azure/acs-engine/pkg/api/v20170701"
	"github.com/Azure/acs-engine/pkg/api/vlabs"
)

func TestAddDCOSPublicAgentPool(t *testing.T) {
	expectedNumPools := 2
	for _, masterCount := range [2]int{1, 3} {
		profiles := []*AgentPoolProfile{}
		profile := makeAgentPoolProfile(1, "agentprivate", "test-dcos-pool", "Standard_D2_v2", "Linux")
		profiles = append(profiles, profile)
		master := makeMasterProfile(masterCount, "test-dcos", "Standard_D2_v2")
		props := getProperties(profiles, master)
		expectedPublicPoolName := props.AgentPoolProfiles[0].Name + publicAgentPoolSuffix
		expectedPublicDNSPrefix := props.AgentPoolProfiles[0].DNSPrefix
		expectedPrivateDNSPrefix := ""
		expectedPublicOSType := props.AgentPoolProfiles[0].OSType
		expectedPublicVMSize := props.AgentPoolProfiles[0].VMSize
		addDCOSPublicAgentPool(props)
		if len(props.AgentPoolProfiles) != expectedNumPools {
			t.Fatalf("incorrect agent pools count. expected=%d actual=%d", expectedNumPools, len(props.AgentPoolProfiles))
		}
		if props.AgentPoolProfiles[1].Name != expectedPublicPoolName {
			t.Fatalf("incorrect public pool name. expected=%s actual=%s", expectedPublicPoolName, props.AgentPoolProfiles[1].Name)
		}
		if props.AgentPoolProfiles[1].DNSPrefix != expectedPublicDNSPrefix {
			t.Fatalf("incorrect public pool DNS prefix. expected=%s actual=%s", expectedPublicDNSPrefix, props.AgentPoolProfiles[1].DNSPrefix)
		}
		if props.AgentPoolProfiles[0].DNSPrefix != expectedPrivateDNSPrefix {
			t.Fatalf("incorrect private pool DNS prefix. expected=%s actual=%s", expectedPrivateDNSPrefix, props.AgentPoolProfiles[0].DNSPrefix)
		}
		if props.AgentPoolProfiles[1].OSType != expectedPublicOSType {
			t.Fatalf("incorrect public pool OS type. expected=%s actual=%s", expectedPublicOSType, props.AgentPoolProfiles[1].OSType)
		}
		if props.AgentPoolProfiles[1].VMSize != expectedPublicVMSize {
			t.Fatalf("incorrect public pool VM size. expected=%s actual=%s", expectedPublicVMSize, props.AgentPoolProfiles[1].VMSize)
		}
		for i, port := range [3]int{80, 443, 8080} {
			if props.AgentPoolProfiles[1].Ports[i] != port {
				t.Fatalf("incorrect public pool port assignment. expected=%d actual=%d", port, props.AgentPoolProfiles[1].Ports[i])
			}
		}
		if props.AgentPoolProfiles[1].Count != masterCount {
			t.Fatalf("incorrect public pool VM size. expected=%d actual=%d", masterCount, props.AgentPoolProfiles[1].Count)
		}
	}
}

func makeAgentPoolProfile(count int, name, dNSPrefix, vMSize string, oSType OSType) *AgentPoolProfile {
	return &AgentPoolProfile{
		Name:      name,
		Count:     count,
		DNSPrefix: dNSPrefix,
		OSType:    oSType,
		VMSize:    vMSize,
	}
}

func makeMasterProfile(count int, dNSPrefix, vMSize string) *MasterProfile {
	return &MasterProfile{
		Count:     count,
		DNSPrefix: "test-dcos",
		VMSize:    "Standard_D2_v2",
	}
}

func getProperties(profiles []*AgentPoolProfile, master *MasterProfile) *Properties {
	return &Properties{
		AgentPoolProfiles: profiles,
		MasterProfile:     master,
	}
}

func TestOrchestratorVersion(t *testing.T) {
	// test v20170701
	v20170701cs := &v20170701.ContainerService{
		Properties: &v20170701.Properties{
			OrchestratorProfile: &v20170701.OrchestratorProfile{
				OrchestratorType: v20170701.Kubernetes,
			},
		},
	}
	cs := ConvertV20170701ContainerService(v20170701cs)
	if cs.Properties.OrchestratorProfile.OrchestratorVersion != common.KubernetesDefaultVersion {
		t.Fatalf("incorrect OrchestratorVersion '%s'", cs.Properties.OrchestratorProfile.OrchestratorVersion)
	}

	v20170701cs = &v20170701.ContainerService{
		Properties: &v20170701.Properties{
			OrchestratorProfile: &v20170701.OrchestratorProfile{
				OrchestratorType:    v20170701.Kubernetes,
				OrchestratorVersion: common.KubernetesVersion1Dot6Dot11,
			},
		},
	}
	cs = ConvertV20170701ContainerService(v20170701cs)
	if cs.Properties.OrchestratorProfile.OrchestratorVersion != common.KubernetesVersion1Dot6Dot11 {
		t.Fatalf("incorrect OrchestratorVersion '%s'", cs.Properties.OrchestratorProfile.OrchestratorVersion)
	}
	// test vlabs
	vlabscs := &vlabs.ContainerService{
		Properties: &vlabs.Properties{
			OrchestratorProfile: &vlabs.OrchestratorProfile{
				OrchestratorType: vlabs.Kubernetes,
			},
		},
	}
	cs = ConvertVLabsContainerService(vlabscs)
	if cs.Properties.OrchestratorProfile.OrchestratorVersion != common.KubernetesDefaultVersion {
		t.Fatalf("incorrect OrchestratorVersion '%s'", cs.Properties.OrchestratorProfile.OrchestratorVersion)
	}

	vlabscs = &vlabs.ContainerService{
		Properties: &vlabs.Properties{
			OrchestratorProfile: &vlabs.OrchestratorProfile{
				OrchestratorType:    vlabs.Kubernetes,
				OrchestratorVersion: common.KubernetesVersion1Dot6Dot11,
			},
		},
	}
	cs = ConvertVLabsContainerService(vlabscs)
	if cs.Properties.OrchestratorProfile.OrchestratorVersion != common.KubernetesVersion1Dot6Dot11 {
		t.Fatalf("incorrect OrchestratorVersion '%s'", cs.Properties.OrchestratorProfile.OrchestratorVersion)
	}
}

func TestKubernetesVlabsDefaults(t *testing.T) {
	vp := makeKubernetesPropertiesVlabs()
	ap := makeKubernetesProperties()
	setVlabsKubernetesDefaults(vp, ap.OrchestratorProfile)
	if ap.OrchestratorProfile.KubernetesConfig == nil {
		t.Fatalf("KubernetesConfig cannot be nil after vlabs default conversion")
	}
	if ap.OrchestratorProfile.KubernetesConfig.NetworkPolicy != vlabs.DefaultNetworkPolicy {
		t.Fatalf("vlabs defaults not applied, expected NetworkPolicy: %s, instead got: %s", vlabs.DefaultNetworkPolicy, ap.OrchestratorProfile.KubernetesConfig.NetworkPolicy)
	}

	vp = makeKubernetesPropertiesVlabs()
	vp.WindowsProfile = &vlabs.WindowsProfile{}
	vp.AgentPoolProfiles = append(vp.AgentPoolProfiles, &vlabs.AgentPoolProfile{OSType: "Windows"})
	ap = makeKubernetesProperties()
	setVlabsKubernetesDefaults(vp, ap.OrchestratorProfile)
	if ap.OrchestratorProfile.KubernetesConfig == nil {
		t.Fatalf("KubernetesConfig cannot be nil after vlabs default conversion")
	}
	if ap.OrchestratorProfile.KubernetesConfig.NetworkPolicy != vlabs.DefaultNetworkPolicyWindows {
		t.Fatalf("vlabs defaults not applied, expected NetworkPolicy: %s, instead got: %s", vlabs.DefaultNetworkPolicyWindows, ap.OrchestratorProfile.KubernetesConfig.NetworkPolicy)
	}
}

func makeKubernetesProperties() *Properties {
	ap := &Properties{}
	ap.OrchestratorProfile = &OrchestratorProfile{}
	ap.OrchestratorProfile.OrchestratorType = "Kubernetes"
	return ap
}

func makeKubernetesPropertiesVlabs() *vlabs.Properties {
	vp := &vlabs.Properties{}
	vp.OrchestratorProfile = &vlabs.OrchestratorProfile{}
	vp.OrchestratorProfile.OrchestratorType = "Kubernetes"
	return vp
}

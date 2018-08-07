package vlabs

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"
)

var testOpenShiftCluster = &OpenShiftCluster{
	ID:       "id",
	Location: "location",
	Name:     "name",
	Plan: &ResourcePurchasePlan{
		Name:          "plan.name",
		Product:       "plan.product",
		PromotionCode: "plan.promotionCode",
		Publisher:     "plan.publisher",
	},
	Tags: map[string]string{
		"tags.k1": "v1",
		"tags.k2": "v2",
	},
	Type: "type",
	Properties: &Properties{
		ProvisioningState: "properties.provisioningState",
		OpenShiftVersion:  "properties.openShiftVersion",
		PublicHostname:    "properties.publicHostname",
		FQDN:              "properties.fqdn",
		RouterProfiles: []RouterProfile{
			{
				Name:            "properties.routerProfiles.0.name",
				PublicSubdomain: "properties.routerProfiles.0.publicSubdomain",
				FQDN:            "properties.routerProfiles.0.fqdn",
			},
			{
				Name:            "properties.routerProfiles.1.name",
				PublicSubdomain: "properties.routerProfiles.1.publicSubdomain",
				FQDN:            "properties.routerProfiles.1.fqdn",
			},
		},
		AgentPoolProfiles: []AgentPoolProfile{
			{
				Name:         "properties.agentPoolProfiles.0.name",
				Role:         "properties.agentPoolProfiles.0.role",
				Count:        1,
				VMSize:       "properties.agentPoolProfiles.0.vmSize",
				VnetSubnetID: "properties.agentPoolProfiles.0.vnetSubnetID",
				OSType:       "properties.agentPoolProfiles.0.osType",
			},
			{
				Name:         "properties.agentPoolProfiles.0.name",
				Role:         "properties.agentPoolProfiles.0.role",
				Count:        2,
				VMSize:       "properties.agentPoolProfiles.0.vmSize",
				VnetSubnetID: "properties.agentPoolProfiles.0.vnetSubnetID",
				OSType:       "properties.agentPoolProfiles.0.osType",
			},
		},
		ServicePrincipalProfile: ServicePrincipalProfile{
			ClientID: "properties.servicePrincipalProfile.clientID",
			Secret:   "properties.servicePrincipalProfile.secret",
		},
	},
}

var testOpenShiftClusterJSON = []byte(`{
	"id": "id",
	"location": "location",
	"name": "name",
	"plan": {
		"name": "plan.name",
		"product": "plan.product",
		"promotionCode": "plan.promotionCode",
		"publisher": "plan.publisher"
	},
	"tags": {
		"tags.k1": "v1",
		"tags.k2": "v2"
	},
	"type": "type",
	"properties": {
		"provisioningState": "properties.provisioningState",
		"openShiftVersion": "properties.openShiftVersion",
		"publicHostname": "properties.publicHostname",
		"fqdn": "properties.fqdn",
		"routerProfiles": [
			{
				"name": "properties.routerProfiles.0.name",
				"publicSubdomain": "properties.routerProfiles.0.publicSubdomain",
				"fqdn": "properties.routerProfiles.0.fqdn"
			},
			{
				"name": "properties.routerProfiles.1.name",
				"publicSubdomain": "properties.routerProfiles.1.publicSubdomain",
				"fqdn": "properties.routerProfiles.1.fqdn"
			}
		],
		"agentPoolProfiles": [
			{
				"name": "properties.agentPoolProfiles.0.name",
				"role": "properties.agentPoolProfiles.0.role",
				"count": 1,
				"vmSize": "properties.agentPoolProfiles.0.vmSize",
				"vnetSubnetID": "properties.agentPoolProfiles.0.vnetSubnetID",
				"osType": "properties.agentPoolProfiles.0.osType"
			},
			{
				"name": "properties.agentPoolProfiles.0.name",
				"role": "properties.agentPoolProfiles.0.role",
				"count": 2,
				"vmSize": "properties.agentPoolProfiles.0.vmSize",
				"vnetSubnetID": "properties.agentPoolProfiles.0.vnetSubnetID",
				"osType": "properties.agentPoolProfiles.0.osType"
			}
		],
		"servicePrincipalProfile": {
			"clientId": "properties.servicePrincipalProfile.clientID",
			"secret": "properties.servicePrincipalProfile.secret"
		}
	}
}`)

func TestMarshal(t *testing.T) {
	b, err := json.MarshalIndent(testOpenShiftCluster, "", "\t")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(b, testOpenShiftClusterJSON) {
		t.Errorf("json.MarshalIndent returned unexpected result\n%s\n", string(b))
	}
}

func TestUnmarshal(t *testing.T) {
	var oc *OpenShiftCluster
	err := json.Unmarshal(testOpenShiftClusterJSON, &oc)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(oc, testOpenShiftCluster) {
		t.Errorf("json.Unmarshal returned unexpected result\n%#v\n", oc)
	}
}

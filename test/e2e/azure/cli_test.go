package azure

import (
	"fmt"
	"testing"
	"time"
)

func TestHasClusterExpired(t *testing.T) {
	rg := ResourceGroup{
		Name:     "testRG",
		Location: "test",
		Tags: map[string]string{
			"now": "799786800",
		},
	}
	a := Account{
		User:           new(User),
		TenantID:       "1234",
		SubscriptionID: "1234",
		ResourceGroup:  rg,
		Deployment:     Deployment{},
	}

	d, err := time.ParseDuration("1h")
	if err != nil {
		t.Fatalf("unexpected error parsing duration: %s", err)
	}
	expected := true
	result := a.HasClusterExpired(d)
	if expected != result {
		t.Fatalf("Resource group should be older than 1h: expected %t but got %t", expected, result)
	}

	a.ResourceGroup.Tags["now"] = fmt.Sprintf("%v", time.Now().Unix())

	d, err = time.ParseDuration("300h")
	if err != nil {
		t.Fatalf("unexpected error parsing duration: %s", err)
	}
	expected = false
	result = a.HasClusterExpired(d)
	if expected != result {
		t.Fatalf("Resource group should not be older than 300h: expected %t but got %t", expected, result)
	}

	a.ResourceGroup.Name = "thisrgdoesntexist"
	d, err = time.ParseDuration("1s")
	if err != nil {
		t.Fatalf("unexpected error parsing duration: %s", err)
	}
	expected = true
	result = a.HasClusterExpired(d)
	if expected != result {
		t.Fatalf("Resource group does not exist: expected %t but got %t", expected, result)
	}
}

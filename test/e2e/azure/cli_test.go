package azure

import (
	"fmt"
	"testing"
	"time"
)

func TestIsClusterExpired(t *testing.T) {
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
	result := a.IsClusterExpired(d)
	if expected != result {
		t.Fatalf("Resource group should be older than 1h: expected %t but got %t", expected, result)
	}

	a.ResourceGroup.Tags["now"] = fmt.Sprintf("%v", time.Now().Unix())

	d, err = time.ParseDuration("300h")
	if err != nil {
		t.Fatalf("unexpected error parsing duration: %s", err)
	}
	expected = false
	result = a.IsClusterExpired(d)
	if expected != result {
		t.Fatalf("Resource group should not be older than 300h: expected %t but got %t", expected, result)
	}

	a.ResourceGroup.Name = "thisrgdoesntexist"
	a.ResourceGroup.Tags = map[string]string{}
	d, err = time.ParseDuration("1s")
	if err != nil {
		t.Fatalf("unexpected error parsing duration: %s", err)
	}
	expected = true
	result = a.IsClusterExpired(d)
	if expected != result {
		t.Fatalf("Resource group does not exist: expected %t but got %t", expected, result)
	}
}

// TODO
// func TestStorageAccount(t *testing.T) {}

// func TestUploadFiles(t *testing.T) {}

// func TestDownloadFiles(t *testing.T) {}

// func TestDeleteFiles(t *testing.T) {}

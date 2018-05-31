package main

import "testing"

func TestGetTrackingID(t *testing.T) {
	output := "ERROR: The template deployment 'openshift-eastus-47708' is not valid according to the validation procedure. The tracking id is 'a111f13b-86bb-4838-a557-fdce07111971'. See inner errors for details. Please see https://aka.ms/arm-deploy for usage details."
	id := getTrackingID(output)
	if id != "a111f13b-86bb-4838-a557-fdce07111971" {
		t.Fatalf("unexpected id: %q\nexpected id: %q", id, "a111f13b-86bb-4838-a557-fdce07111971")
	}
}

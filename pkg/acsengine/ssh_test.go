package acsengine

import (
	"os"
	"testing"

	"github.com/Azure/acs-engine/pkg/i18n"
)

func TestCreateSaveSSH(t *testing.T) {
	translator := &i18n.Translator{
		Locale: nil,
	}
	username := "test_user"
	outputDirectory := "unit_tests"
	expectedFile := outputDirectory + "/" + username + "_rsa"

	defer os.Remove(expectedFile)

	_, _, err := CreateSaveSSH(username, outputDirectory, translator)

	if err != nil {
		t.Fatalf("Unexpected error creating and saving ssh key: %s", err)
	}

	if _, err := os.Stat(expectedFile); os.IsNotExist(err) {
		t.Fatalf("ssh file was not created")
	}
}

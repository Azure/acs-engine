package acsengine

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"

	"github.com/Azure/acs-engine/pkg/helpers"
	"github.com/Azure/acs-engine/pkg/i18n"
)

// CreateSaveSSH generates and stashes an SSH key pair.
func CreateSaveSSH(username, outputDirectory string, s *i18n.Translator) (privateKey *rsa.PrivateKey, publicKeyString string, err error) {

	privateKey, publicKeyString, err = helpers.CreateSSH(rand.Reader, s)
	if err != nil {
		return nil, "", err
	}

	privateKeyPem := privateKeyToPem(privateKey)

	f := &FileSaver{
		Translator: s,
	}

	err = f.SaveFile(outputDirectory, fmt.Sprintf("%s_rsa", username), privateKeyPem)
	if err != nil {
		return nil, "", err
	}

	return privateKey, publicKeyString, nil
}

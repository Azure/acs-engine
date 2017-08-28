package acsengine

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"io"

	"github.com/Azure/acs-engine/pkg/i18n"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
)

// SSHCreator represents the object that creates SSH key pair
type SSHCreator struct {
	Translator *i18n.Translator
}

const (
	// SSHKeySize is the size (in bytes) of SSH key to create
	SSHKeySize = 4096
)

// CreateSaveSSH generates and stashes an SSH key pair.
func (s *SSHCreator) CreateSaveSSH(username, outputDirectory string) (privateKey *rsa.PrivateKey, publicKeyString string, err error) {
	privateKey, publicKeyString, err = s.CreateSSH(rand.Reader)
	if err != nil {
		return nil, "", err
	}

	privateKeyPem := privateKeyToPem(privateKey)

	f := &FileSaver{
		Translator: s.Translator,
	}

	err = f.SaveFile(outputDirectory, fmt.Sprintf("%s_rsa", username), privateKeyPem)
	if err != nil {
		return nil, "", err
	}

	return privateKey, publicKeyString, nil
}

// CreateSSH creates an SSH key pair.
func (s *SSHCreator) CreateSSH(rg io.Reader) (privateKey *rsa.PrivateKey, publicKeyString string, err error) {
	log.Debugf("ssh: generating %dbit rsa key", SSHKeySize)
	privateKey, err = rsa.GenerateKey(rg, SSHKeySize)
	if err != nil {
		return nil, "", s.Translator.Errorf("failed to generate private key for ssh: %q", err)
	}

	publicKey := privateKey.PublicKey
	sshPublicKey, err := ssh.NewPublicKey(&publicKey)
	if err != nil {
		return nil, "", s.Translator.Errorf("failed to create openssh public key string: %q", err)
	}
	authorizedKeyBytes := ssh.MarshalAuthorizedKey(sshPublicKey)
	authorizedKey := string(authorizedKeyBytes)

	return privateKey, authorizedKey, nil
}

package acsengine

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"io"

	log "github.com/Sirupsen/logrus"
	"golang.org/x/crypto/ssh"
)

const (
	// SshKeySize is the size of SSH key to create
	SSHKeySize = 4096
)

// CreateSaveSSH generates and stashes an SSH key pair.
func CreateSaveSSH(username, outputDirectory string) (privateKey *rsa.PrivateKey, publicKeyString string, err error) {
	privateKey, publicKeyString, err = CreateSSH(rand.Reader)
	if err != nil {
		return nil, "", err
	}

	privateKeyPem := privateKeyToPem(privateKey)

	err = saveFile(outputDirectory, fmt.Sprintf("%s_rsa", username), privateKeyPem)
	if err != nil {
		return nil, "", err
	}

	return privateKey, publicKeyString, nil
}

// CreateSSH creates an SSH key pair.
func CreateSSH(rg io.Reader) (privateKey *rsa.PrivateKey, publicKeyString string, err error) {
	log.Debugf("ssh: generating %dbit rsa key", SSHKeySize)
	privateKey, err = rsa.GenerateKey(rg, SSHKeySize)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate private key for ssh: %q", err)
	}

	publicKey := privateKey.PublicKey
	sshPublicKey, err := ssh.NewPublicKey(&publicKey)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create openssh public key string: %q", err)
	}
	authorizedKeyBytes := ssh.MarshalAuthorizedKey(sshPublicKey)
	authorizedKey := string(authorizedKeyBytes)

	return privateKey, authorizedKey, nil
}

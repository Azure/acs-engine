package helpers

import (
	// "fmt"
	"bytes"
	"crypto/rsa"
	"encoding/json"
	"io"
	"strings"

	"github.com/Azure/acs-engine/pkg/i18n"
	"golang.org/x/crypto/ssh"
)

const (
	// SSHKeySize is the size (in bytes) of SSH key to create
	SSHKeySize = 4096
)

// NormalizeAzureRegion returns a normalized Azure region with white spaces removed and converted to lower case
func NormalizeAzureRegion(name string) string {
	return strings.ToLower(strings.Replace(name, " ", "", -1))
}

// JSONMarshalIndent marshals formatted JSON w/ optional SetEscapeHTML
func JSONMarshalIndent(content interface{}, prefix, indent string, escape bool) ([]byte, error) {
	b, err := JSONMarshal(content, escape)
	if err != nil {
		return nil, err
	}

	var bufIndent bytes.Buffer
	if err := json.Indent(&bufIndent, b, prefix, indent); err != nil {
		return nil, err
	}

	return bufIndent.Bytes(), nil
}

// JSONMarshal marshals JSON w/ optional SetEscapeHTML
func JSONMarshal(content interface{}, escape bool) ([]byte, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(escape)
	if err := enc.Encode(content); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// IsTrueBoolPointer is a simple boolean helper function for boolean pointers
func IsTrueBoolPointer(b *bool) bool {
	if b != nil && *b {
		return true
	}
	return false
}

// PointerToBool returns a pointer to a bool
func PointerToBool(b bool) *bool {
	p := b
	return &p
}

// CreateSSH creates an SSH key pair.
func CreateSSH(rg io.Reader, s *i18n.Translator) (privateKey *rsa.PrivateKey, publicKeyString string, err error) {
	privateKey, err = rsa.GenerateKey(rg, SSHKeySize)
	if err != nil {
		return nil, "", s.Errorf("failed to generate private key for ssh: %q", err)
	}

	publicKey := privateKey.PublicKey
	sshPublicKey, err := ssh.NewPublicKey(&publicKey)
	if err != nil {
		return nil, "", s.Errorf("failed to create openssh public key string: %q", err)
	}
	authorizedKeyBytes := ssh.MarshalAuthorizedKey(sshPublicKey)
	authorizedKey := string(authorizedKeyBytes)

	return privateKey, authorizedKey, nil
}

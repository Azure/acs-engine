package helpers

import (
	// "fmt"
	"bytes"
	"crypto/rsa"
	"encoding/json"
	"io"
	"os"
	"runtime"
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

// PointerToInt returns a pointer to a int
func PointerToInt(i int) *int {
	p := i
	return &p
}

// EqualError is a ni;-safe method which reports whether errors a and b are considered equal.
// They're equal if both are nil, or both are not nil and a.Error() == b.Error().
func EqualError(a, b error) bool {
	return a == nil && b == nil || a != nil && b != nil && a.Error() == b.Error()
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

// AcceleratedNetworkingSupported check if the VmSKU support the Accelerated Networking
func AcceleratedNetworkingSupported(sku string) bool {
	if strings.Contains(sku, "Standard_D2s_v3") {
		return false
	}
	if strings.Contains(sku, "Standard_DS3") {
		return false
	}
	if strings.Contains(sku, "Standard_D2_v3") {
		return false
	}
	if strings.Contains(sku, "Standard_A") {
		return false
	}
	if strings.Contains(sku, "Standard_B") {
		return false
	}
	if strings.Contains(sku, "Standard_G") {
		return false
	}
	if strings.Contains(sku, "Standard_H") {
		return false
	}
	if strings.Contains(sku, "Standard_L") {
		return false
	}
	if strings.Contains(sku, "Standard_N") {
		return false
	}
	if strings.EqualFold(sku, "Standard_D1") || strings.Contains(sku, "Standard_D1_") {
		return false
	}
	if strings.EqualFold(sku, "Standard_DS1") || strings.Contains(sku, "Standard_DS1_") {
		return false
	}
	if strings.EqualFold(sku, "Standard_F1") || strings.EqualFold(sku, "Standard_F1s") {
		return false
	}
	return true
}

// GetHomeDir attempts to get the home dir from env
func GetHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}

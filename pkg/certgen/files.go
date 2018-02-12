package certgen

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"regexp"
	"strings"
	"text/template"

	"github.com/Azure/acs-engine/pkg/certgen/templates"
	"github.com/Azure/acs-engine/pkg/filesystem"
)

// PrepareMasterFiles creates the shared authentication and encryption secrets
func (c *Config) PrepareMasterFiles() error {
	b := make([]byte, 24)
	_, err := rand.Read(b)
	if err != nil {
		return err
	}
	c.AuthSecret = base64.StdEncoding.EncodeToString(b)

	_, err = rand.Read(b)
	if err != nil {
		return err
	}
	c.EncSecret = base64.StdEncoding.EncodeToString(b)

	return nil
}

// WriteMasterFiles writes the templated master config
func (c *Config) WriteMasterFiles(fs filesystem.Filesystem) error {
	for _, name := range templates.AssetNames() {
		if !strings.HasPrefix(name, "master/") {
			continue
		}
		tb := templates.MustAsset(name)

		t, err := template.New("template").Funcs(template.FuncMap{
			"QuoteMeta": regexp.QuoteMeta,
		}).Parse(string(tb))
		if err != nil {
			return err
		}

		b := &bytes.Buffer{}
		err = t.Execute(b, c)
		if err != nil {
			return err
		}

		err = fs.WriteFile(strings.TrimPrefix(name, "master/"), b.Bytes(), 0666)
		if err != nil {
			return err
		}
	}

	return nil
}

// WriteNodeFiles writes the templated node config
func (c *Config) WriteNodeFiles(fs filesystem.Filesystem) error {
	for _, name := range templates.AssetNames() {
		if !strings.HasPrefix(name, "node/") {
			continue
		}

		tb := templates.MustAsset(name)

		t, err := template.New("template").Funcs(template.FuncMap{
			"QuoteMeta": regexp.QuoteMeta,
		}).Parse(string(tb))
		if err != nil {
			return err
		}

		b := &bytes.Buffer{}
		err = t.Execute(b, struct{}{})
		if err != nil {
			return err
		}

		err = fs.WriteFile(strings.TrimPrefix(name, "node/"), b.Bytes(), 0666)
		if err != nil {
			return err
		}
	}

	return nil
}

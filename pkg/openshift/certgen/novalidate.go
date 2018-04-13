package certgen

// novalidate.go is split out of files.go to avoid static validation.
// `make test-style` is failing non-deterministically (flakying) with
// the following message:
//
// pkg/certgen/files.go:36:23:warning: AssetNames not declared by package templates (unused)

import (
	"github.com/Azure/acs-engine/pkg/openshift/certgen/templates"
)

func getAssets() []string {
	return templates.AssetNames()
}

func assetMustExist(name string) []byte {
	return templates.MustAsset(name)
}

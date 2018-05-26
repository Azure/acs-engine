//go:generate go-bindata -nometadata -pkg $GOPACKAGE master/... node/...
//go:generate gofmt -s -l -w bindata.go

package templates

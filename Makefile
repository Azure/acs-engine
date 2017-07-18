.NOTPARALLEL:

.PHONY: prereqs build test test_fmt validate-generated fmt lint ci devenv

VERSION=`git describe --always --long --dirty`
BUILD=`date +%FT%T%z`

# this isn't particularly pleasant, but it works with the least amount
# of requirements around $GOPATH. The extra sed is needed because `gofmt`
# operates on paths, go list returns package names, and `go fmt` always rewrites
# which is not what we need to do in the `test_fmt` target.
GOFILES=`go list ./... | grep -v "github.com/Azure/acs-engine/vendor" | sed 's|github.com/Azure/acs-engine|.|g' | grep -v -w '^.$$'`

all: build

prereqs:
	go get github.com/Masterminds/glide
	go get github.com/jteeuwen/go-bindata/...
	glide install

build:
	go generate -v $(GOFILES)
	go build -v -ldflags="-X github.com/Azure/acs-engine/cmd.BuildSHA=${VERSION} -X github.com/Azure/acs-engine/cmd.BuildTime=${BUILD}"
	cd test/acs-engine-test; go build -v

test: test_fmt
	go test -v $(GOFILES)

.PHONY: test-style
test-style:
	@scripts/validate-go.sh

ci: prereqs build test lint

devenv:
	./scripts/devenv.sh

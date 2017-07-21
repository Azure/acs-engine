.NOTPARALLEL:

.PHONY: bootstrap build test test_fmt validate-generated fmt lint ci devenv

VERSION=`git describe --always --long --dirty`
BUILD=`date +%FT%T%z`

# this isn't particularly pleasant, but it works with the least amount
# of requirements around $GOPATH. The extra sed is needed because `gofmt`
# operates on paths, go list returns package names, and `go fmt` always rewrites
# which is not what we need to do in the `test_fmt` target.
GOFILES=`go list ./... | grep -v "github.com/Azure/acs-engine/vendor" | sed 's|github.com/Azure/acs-engine|.|g' | grep -v -w '^.$$'`

all: build

HAS_GLIDE := $(shell command -v glide;)
HAS_GOX := $(shell command -v gox;)
HAS_GIT := $(shell command -v git;)
HAS_GOBINDATA := $(shell command -v go-bindata;)

.PHONY: bootstrap
bootstrap:
ifndef HAS_GLIDE
	go get -u github.com/Masterminds/glide
endif
ifndef HAS_GOX
	go get -u github.com/mitchellh/gox
endif
ifndef HAS_GOBINDATA
	go get github.com/jteeuwen/go-bindata/...
endif
ifndef HAS_GIT
	$(error You must install Git)
endif
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

ci: bootstrap build test lint

devenv:
	./scripts/devenv.sh

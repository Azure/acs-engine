.NOTPARALLEL:

.PHONY: prereqs build test test_fmt validate-generated fmt lint ci devenv

VERSION=`git describe --always --long --dirty`
BUILD=`date +%FT%T%z`

all: build

prereqs:
	go get github.com/jteeuwen/go-bindata
	go get github.com/Masterminds/glide
	glide install

_build:
	go generate -v ./pkg/...
	go build -v -ldflags="-X github.com/Azure/acs-engine/cmd.BuildSHA=${VERSION} -X github.com/Azure/acs-engine/cmd.BuildTime=${BUILD}"
	cd test/acs-engine-test; go build -v

build: prereqs _build

test: prereqs test_fmt
	go test -v ./...

test_fmt: prereqs
	test -z "$$(gofmt -s -l . | tee /dev/stderr)"

validate-generated: prereqs
	./scripts/validate-generated.sh

fmt:
	gofmt -s -l -w .

lint: prereqs
	go get -u github.com/golang/lint/golint
	# TODO: fix lint errors, enable linting
	# golint -set_exit_status

ci: validate-generated build test lint

devenv:
	./scripts/devenv.sh
.NOTPARALLEL:

.PHONY: prereqs build test test_fmt validate-generated fmt lint ci devenv

VERSION=`git describe --always --long --dirty`
BUILD=`date +%FT%T%z`

prereqs:
	go get github.com/jteeuwen/go-bindata/...
	go get github.com/satori/uuid/...

build: prereqs
	go generate -v ./...
	go get .
	go build -v -ldflags="-X main.AcsEngineBuildSHA=${VERSION} -X main.AcsEngineBuildTime=${BUILD}"
	cd test/acs-engine-test; go build -v

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

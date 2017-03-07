.NOTPARALLEL:

.PHONY: build test validate-generated lint ci devenv

prereqs:
	go get github.com/jteeuwen/go-bindata/...

build: prereqs
	go generate -v ./...
	go get .
	go build -v
	cd test/acs-engine-test; go build -v

test: prereqs
	go test -v ./...

validate-generated: prereqs
	./scripts/validate-generated.sh

lint: prereqs
	go get -u github.com/golang/lint/golint
	# TODO: fix lint errors, enable linting
	# golint -set_exit_status

ci: validate-generated build test lint

devenv:
	./scripts/devenv.sh

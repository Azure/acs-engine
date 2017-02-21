.NOTPARALLEL:

.PHONY: build test validate-generated lint ci devenv

build:
	go get github.com/jteeuwen/go-bindata/...
	go generate -v ./...
	go get .
	go build -v

test:
	go test -v ./...

validate-generated:
	./scripts/validate-generated.sh

lint:
	go get -u github.com/golang/lint/golint
	# TODO: fix lint errors, enable linting
	# golint -set_exit_status

ci: validate-generated build test lint

devenv:
	./scripts/devenv.sh

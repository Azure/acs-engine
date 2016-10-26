.NOTPARALLEL:

build:
	go get .
	go generate -v ./...
	go build -v

test:
	go test -v ./...

validate-generated:
	./scripts/validate-generated.sh

lint:
	# TODO: fix lint errors, enable linting
	# golint -set_exit_status

ci: validate-generated build test lint

devenv:
	./scripts/devenv.sh


.NOTPARALLEL:

build:
	go generate -v
	go build -v

test:
	go test -v ./...

validate-generated:
	./scripts/validate-generated.sh

lint:
	# TODO: fix lint errors, enable linting
	# golint -set_exit_status

ci: validate-generated build test lint

dev:
	docker build -t acs-engine .
	docker run -it -v `pwd`:/acs-engine -w /acs-engine acs-engine /bin/bash


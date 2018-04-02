TARGETS           = darwin/amd64 linux/amd64 windows/amd64
DIST_DIRS         = find * -type d -exec

.NOTPARALLEL:

.PHONY: bootstrap build test test_fmt validate-generated fmt lint ci devenv

ifdef DEBUG
GOFLAGS   := -gcflags="-N -l"
else
GOFLAGS   :=
endif

# go option
GO              ?= go
TAGS            :=
LDFLAGS         :=
BINDIR          := $(CURDIR)/bin
BINARIES        := acs-engine
VERSION         ?= $(shell git rev-parse HEAD)
VERSION_SHORT   ?= $(shell git rev-parse --short HEAD)
GITTAG          := $(shell git describe --exact-match --tags $(shell git log -n1 --pretty='%h') 2> /dev/null)
ifeq ($(GITTAG),)
GITTAG := $(VERSION_SHORT)
endif

REPO_PATH := github.com/Azure/acs-engine
DEV_ENV_IMAGE := quay.io/deis/go-dev:v1.9.1
DEV_ENV_WORK_DIR := /go/src/${REPO_PATH}
DEV_ENV_OPTS := --rm -v ${CURDIR}:${DEV_ENV_WORK_DIR} -w ${DEV_ENV_WORK_DIR} ${DEV_ENV_VARS}
DEV_ENV_CMD := docker run ${DEV_ENV_OPTS} ${DEV_ENV_IMAGE}
DEV_ENV_CMD_IT := docker run -it ${DEV_ENV_OPTS} ${DEV_ENV_IMAGE}
DEV_CMD_RUN := docker run ${DEV_ENV_OPTS}
ifdef DEBUG
LDFLAGS := -X main.version=${VERSION}
else
LDFLAGS := -s -X main.version=${VERSION}
endif
BINARY_DEST_DIR ?= bin

all: build

.PHONY: dev
dev:
	$(DEV_ENV_CMD_IT) bash

.PHONY: generate
generate: bootstrap
	go generate $(GOFLAGS) -v `glide novendor | xargs go list`

.PHONY: generate-azure-constants
generate-azure-constants:
	python pkg/acsengine/Get-AzureConstants.py

.PHONY: build
build: generate
	GOBIN=$(BINDIR) $(GO) install $(GOFLAGS) -ldflags '$(LDFLAGS)'
	cd test/acs-engine-test; go build $(GOFLAGS)

build-binary: generate
	go build $(GOFLAGS) -v -ldflags "${LDFLAGS}" -o ${BINARY_DEST_DIR}/acs-engine .

# usage: make clean build-cross dist VERSION=v0.4.0
.PHONY: build-cross
build-cross: build
build-cross: LDFLAGS += -extldflags "-static"
build-cross:
	CGO_ENABLED=0 gox -output="_dist/acs-engine-${GITTAG}-{{.OS}}-{{.Arch}}/{{.Dir}}" -osarch='$(TARGETS)' $(GOFLAGS) -tags '$(TAGS)' -ldflags '$(LDFLAGS)'

.PHONY: build-windows-k8s
build-windows-k8s:
	./scripts/build-windows-k8s.sh -v ${K8S_VERSION} -p ${PATCH_VERSION}

.PHONY: dist
dist: build-cross
	( \
		cd _dist && \
		$(DIST_DIRS) cp ../LICENSE {} \; && \
		$(DIST_DIRS) cp ../README.md {} \; && \
		$(DIST_DIRS) tar -zcf {}.tar.gz {} \; && \
		$(DIST_DIRS) zip -r {}.zip {} \; \
	)

.PHONY: checksum
checksum:
	for f in _dist/*.{gz,zip} ; do \
		shasum -a 256 "$${f}"  | awk '{print $$1}' > "$${f}.sha256" ; \
	done

.PHONY: clean
clean:
	@rm -rf $(BINDIR) ./_dist

GIT_BASEDIR    = $(shell git rev-parse --show-toplevel 2>/dev/null)
ifneq ($(GIT_BASEDIR),)
	LDFLAGS += -X github.com/Azure/acs-engine/pkg/test.JUnitOutDir=${GIT_BASEDIR}/test/junit
endif

test: generate
	ginkgo -skipPackage test/e2e -r .

.PHONY: test-style
test-style:
	@scripts/validate-go.sh

.PHONY: test-e2e
test-e2e:
	@test/e2e.sh

HAS_GLIDE := $(shell command -v glide;)
HAS_GOX := $(shell command -v gox;)
HAS_GIT := $(shell command -v git;)
HAS_GOBINDATA := $(shell command -v go-bindata;)
HAS_GOMETALINTER := $(shell command -v gometalinter;)
HAS_GINKGO := $(shell command -v ginkgo;)

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
ifndef HAS_GOMETALINTER
	go get -u github.com/alecthomas/gometalinter
	gometalinter --install
endif
ifndef HAS_GINKGO
	go get -u github.com/onsi/ginkgo/ginkgo
endif

build-vendor:
	${DEV_ENV_CMD} rm -f glide.lock && rm -Rf vendor/ && glide --debug install --force --strip-vendor
	rm -rf vendor/github.com/docker/distribution/contrib/docker-integration/generated_certs.d

ci: bootstrap test-style build test lint
	./scripts/coverage.sh --coveralls

.PHONY: coverage
coverage:
	@scripts/ginkgo.coverage.sh

devenv:
	./scripts/devenv.sh

include versioning.mk
include test.mk

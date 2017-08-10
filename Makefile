TARGETS           = darwin/amd64 linux/amd64 windows/amd64
DIST_DIRS         = find * -type d -exec

.NOTPARALLEL:

.PHONY: bootstrap build test test_fmt validate-generated fmt lint ci devenv

# go option
GO        ?= go
PKG       := $(shell glide novendor)
TAGS      :=
LDFLAGS   :=
GOFLAGS   :=
BINDIR    := $(CURDIR)/bin
BINARIES  := acs-engine
VERSION   := $(shell git rev-parse HEAD)

# this isn't particularly pleasant, but it works with the least amount
# of requirements around $GOPATH. The extra sed is needed because `gofmt`
# operates on paths, go list returns package names, and `go fmt` always rewrites
# which is not what we need to do in the `test_fmt` target.
GOFILES=`go list ./... | grep -v "github.com/Azure/acs-engine/vendor" | sed 's|github.com/Azure/acs-engine|.|g' | grep -v -w '^.$$'`

all: build

.PHONY: generate
generate:
	go generate -v $(GOFILES)

.PHONY: build
build: generate
	GOBIN=$(BINDIR) $(GO) install $(GOFLAGS) -ldflags '$(LDFLAGS)'
	cd test/acs-engine-test; go build

# usage: make clean build-cross dist VERSION=v0.4.0
.PHONY: build-cross
build-cross: LDFLAGS += -extldflags "-static"
build-cross:
	CGO_ENABLED=0 gox -output="_dist/acs-engine-${VERSION}-{{.OS}}-{{.Arch}}/{{.Dir}}" -osarch='$(TARGETS)' $(GOFLAGS) -tags '$(TAGS)' -ldflags '$(LDFLAGS)'

.PHONY: build-windows-k8s
build-windows-k8s:
	./scripts/build-windows-k8s.sh -v ${K8S_VERSION}

.PHONY: dist
dist:
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

test:
	ginkgo -r -ldflags='$(LDFLAGS)' .

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
	glide install
ifndef HAS_GINKGO
	go get -u github.com/onsi/ginkgo/ginkgo
endif


ci: bootstrap test-style build test lint
	./scripts/coverage.sh --coveralls

.PHONY: coverage
coverage:
	@scripts/coverage.sh

devenv:
	./scripts/devenv.sh

include versioning.mk

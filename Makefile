CLYMENE_IMPORT_PATH=github.com/Clymene-project/Clymene

include docker/Makefile

# all .go files that are not auto-generated and should be auto-formatted and linted.
ALL_SRC := $(shell find . -name '*.go' \
				   -not -name 'doc.go' \
				   -not -name '_*' \
				   -not -name '.*' \
				   -not -name 'model.pb.go' \
				   -not -path './examples/*' \
				   -not -path './vendor/*' \
				   -not -path '*/mocks/*' \
				   -type f | \
				sort)

# ALL_PKGS is used with 'golint'
ALL_PKGS := $(shell echo $(dir $(ALL_SRC)) | tr ' ' '\n' | sort -u)

UNAME := $(shell uname -m)
#Race flag is not supported on s390x architecture
ifeq ($(UNAME), s390x)
	RACE=
else
	RACE=-race
endif
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
GOBUILD=CGO_ENABLED=0 installsuffix=cgo go build -trimpath
GOTEST=go test -v $(RACE)
GOLINT=golint
GOVET=go vet
GOFMT=gofmt
FMT_LOG=.fmt.log
LINT_LOG=.lint.log
IMPORT_LOG=.import.log

GIT_SHA=$(shell git rev-parse HEAD)
GIT_BRANCH=$(shell git branch)

DATE=$(shell date -u +'%Y-%m-%dT%H:%M:%SZ')

BUILD_INFO=-ldflags "-X 'main.Version=$(GIT_BRANCH)($(GIT_SHA))' -X 'main.BuildTime=$(DATE)'"

.PHONY: build-agent
build-agent :
	$(GOBUILD) -o ./out/clymene-agent-$(GOOS)-$(GOARCH) $(BUILD_INFO) ./cmd/agent/main.go

.PHONY: build-ingester
build-ingester :
	$(GOBUILD) -o ./out/clymene-ingester-$(GOOS)-$(GOARCH) $(BUILD_INFO) ./cmd/ingester/main.go

.PHONY: build-gateway
build-gateway :
	$(GOBUILD) -o ./out/clymene-gateway-$(GOOS)-$(GOARCH) $(BUILD_INFO) ./cmd/gateway/main.go

.PHONY: build-promtail
build-promtail :
	$(GOBUILD) -o ./out/clymene-promtail-$(GOOS)-$(GOARCH) $(BUILD_INFO) ./cmd/promtail/main.go

.PHONY: docker
docker: build-binaries-linux docker-images-only

.PHONY: build-binaries-linux
build-binaries-linux:
	GOOS=linux GOARCH=amd64 $(MAKE) build-platform-binaries

.PHONY: build-binaries-windows
build-binaries-windows:
	GOOS=windows GOARCH=amd64 $(MAKE) build-platform-binaries

.PHONY: build-binaries-darwin
build-binaries-darwin:
	GOOS=darwin GOARCH=amd64 $(MAKE) build-platform-binaries

.PHONY: build-binaries-s390x
build-binaries-s390x:
	GOOS=linux GOARCH=s390x $(MAKE) build-platform-binaries

.PHONY: build-binaries-arm64
build-binaries-arm64:
	GOOS=linux GOARCH=arm64 $(MAKE) build-platform-binaries

.PHONY: build-binaries-ppc64le
build-binaries-ppc64le:
	GOOS=linux GOARCH=ppc64le $(MAKE) build-platform-binaries

.PHONY: build-platform-binaries
build-platform-binaries: build-agent \
	build-ingester \
	build-gateway \
	build-promtail \

.PHONY: build-all-platforms
build-all-platforms: build-binaries-linux build-binaries-windows build-binaries-darwin build-binaries-s390x build-binaries-arm64 build-binaries-ppc64le
GOCMD=$(shell command -v go)
GOBUILD=$(GOCMD) build
GOINSTALL=$(GOCMD) install
GORUN=$(GOCMD) run
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOTOOL=$(GOCMD) tool
GOMOD=$(GOCMD) mod
GOPATH=$(shell $(GOCMD) env GOPATH)
GOBIN=$(shell $(GOCMD) env GOBIN)

ifeq ($(GOBIN),)
	GOBIN := $(GOPATH)/bin
endif

BINARY=filescript
BINARY_DARWIN=$(BINARY)_darwin
BINARY_LINUX=$(BINARY)_linux
BINARY_PI=$(BINARY)_pi
BINARY_WINDOWS=$(BINARY).exe

BUILD_DATE=`date`
BUILD_COMMIT=`git rev-parse HEAD`

CURRENT_DIR = $(shell pwd)


.PHONY: prepare
prepare: download
	@echo "[*] $@"
	GOPATH="$(GOPATH)" \
 	$(GOCMD) generate ./...

.PHONY: download
download:
	@echo "[*] $@"
	@$(GOMOD) download

.PHONY: build
build: prepare
	@echo "[*] $@"
	$(GOBUILD) -o $(BINARY) -v -ldflags "-X 'main.commit=${BUILD_COMMIT}' -X 'main.date=${BUILD_DATE}' -X 'main.builtBy=make'"

.PHONY: build-mac
build-mac: prepare
	@echo "[*] $@"
	GOOS="darwin" GOARCH="amd64" $(GOBUILD) -o $(BINARY_DARWIN) -v -ldflags "-X 'main.commit=${BUILD_COMMIT}' -X 'main.date=${BUILD_DATE}' -X 'main.builtBy=make'"

.PHONY: build-linux
build-linux: prepare
	@echo "[*] $@"
	GOOS="linux" GOARCH="amd64" $(GOBUILD) -o $(BINARY_LINUX) -v -ldflags "-X 'main.commit=${BUILD_COMMIT}' -X 'main.date=${BUILD_DATE}' -X 'main.builtBy=make'"

.PHONY: build-pi
build-pi: prepare
	@echo "[*] $@"
	GOOS="linux" GOARCH="arm" GOARM="6" $(GOBUILD) -o $(BINARY_PI) -v -ldflags "-X 'main.commit=${BUILD_COMMIT}' -X 'main.date=${BUILD_DATE}' -X 'main.builtBy=make'"

.PHONY: build-windows
build-windows: prepare
	@echo "[*] $@"
	GOOS="windows" GOARCH="amd64" $(GOBUILD) -o $(BINARY_WINDOWS) -v -ldflags "-X 'main.commit=${BUILD_COMMIT}' -X 'main.date=${BUILD_DATE}' -X 'main.builtBy=make'"

.PHONY: build-all
build-all: build-mac build-linux build-pi build-windows

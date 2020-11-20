EXECUTABLE=check_https_go
WINDOWS=$(EXECUTABLE)_windows_amd64.exe
LINUX=$(EXECUTABLE)_linux_amd64
MACOS=$(EXECUTABLE)_macos_amd64
VERSION=$(shell git describe --tags --always --long --dirty)
BINDIR=bin/

.PHONY: all test clean

all: clean build ## Clean and build

build: windows linux macos ## Build binaries
	@echo version: $(VERSION)
	@echo built to: $(shell pwd)/$(BINDIR)

windows: $(WINDOWS) ## Build for Windows

linux: $(LINUX) ## Build for Linux

macos: $(MACOS) ## Build for macOS (Darwin)

$(WINDOWS):
	env GOOS=windows GOARCH=amd64 go build -v -o $(BINDIR)$(WINDOWS) -ldflags="-s -w -X main.version=$(VERSION)" main.go

$(LINUX):
	env GOOS=linux GOARCH=amd64 go build -v -o $(BINDIR)$(LINUX) -ldflags="-s -w -X main.version=$(VERSION)" main.go

$(MACOS):
	env GOOS=darwin GOARCH=amd64 go build -v -o $(BINDIR)$(MACOS) -ldflags="-s -w -X main.version=$(VERSION)" main.go

clean: ## Remove previous build
	rm -f bin/$(WINDOWS) bin/$(LINUX) bin/$(MACOS)

help: ## Display available commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
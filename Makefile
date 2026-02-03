APP_NAME := promdigger
CMD := main.go
DIST := dist
LDFLAGS := -s -w

# Detect host OS/ARCH by default
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
CGO_ENABLED ?= 0

.PHONY: all build run test clean help build-all

all: build

help: ## Show this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | \
	awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

build: ## Build for current OS/ARCH
	@echo "Building for $(GOOS)/$(GOARCH)"
	CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) \
	go build -ldflags "$(LDFLAGS)" -o $(APP_NAME) $(CMD)

run: build ## Build and run locally
	./$(APP_NAME)

test: ## Run unit tests
	go test -v ./...

clean: ## Remove generated files
	rm -rf $(APP_NAME) $(DIST)
	rm -rf data
	go clean


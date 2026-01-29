# Define the application name (optional, but useful for consistency)
APP_NAME := promdigger

.PHONY: build run test clean help

# Default target
all: build

help: ## Show this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

build: ## Compile the Go program
	go build -o $(APP_NAME) main.go

run: build ## Build and run the Go program
	./$(APP_NAME)

test: ## Run unit tests
	go test -v ./...

clean: ## Remove generated files
	rm -f $(APP_NAME)
	go clean

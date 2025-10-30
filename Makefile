.PHONY: help build test clean run docker-build docker-run lint fmt vet deps coverage integration-test

# Variables
BINARY_NAME=dbus-controller
DOCKER_IMAGE=mesbrj/dbus-controller
VERSION=$(shell git describe --tags --always --dirty)
BUILD_DIR=build
COVERAGE_DIR=coverage

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

help: ## Display this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

deps: ## Download dependencies
	$(GOMOD) download
	$(GOMOD) verify

fmt: ## Format Go code
	$(GOCMD) fmt ./...

vet: ## Run go vet
	$(GOCMD) vet ./...

lint: ## Run linter
	golangci-lint run

test: ## Run unit tests
	mkdir -p $(COVERAGE_DIR)
	$(GOTEST) -v -race -coverprofile=$(COVERAGE_DIR)/coverage.out ./...
	$(GOCMD) tool cover -html=$(COVERAGE_DIR)/coverage.out -o $(COVERAGE_DIR)/coverage.html

test-short: ## Run unit tests (short mode)
	$(GOTEST) -short -v ./...

integration-test: ## Run integration tests (requires D-Bus)
	$(GOTEST) -v -tags=integration ./...

coverage: test ## Generate test coverage report
	$(GOCMD) tool cover -func=$(COVERAGE_DIR)/coverage.out

build: deps fmt vet ## Build the application
	mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=linux $(GOBUILD) -a -installsuffix cgo \
		-ldflags "-X main.version=$(VERSION)" \
		-o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/server

build-all: ## Build for multiple platforms
	mkdir -p $(BUILD_DIR)
	# Linux
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) \
		-ldflags "-X main.version=$(VERSION)" \
		-o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 ./cmd/server
	# macOS
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) \
		-ldflags "-X main.version=$(VERSION)" \
		-o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 ./cmd/server
	# Windows
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) \
		-ldflags "-X main.version=$(VERSION)" \
		-o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe ./cmd/server

run: build ## Build and run the application
	./$(BUILD_DIR)/$(BINARY_NAME)

run-dev: ## Run in development mode with live reload (requires air)
	air

docker-build: ## Build Docker image
	docker build -t $(DOCKER_IMAGE):$(VERSION) .
	docker tag $(DOCKER_IMAGE):$(VERSION) $(DOCKER_IMAGE):latest

docker-run: docker-build ## Build and run Docker container
	docker run --rm -p 8080:8080 \
		--privileged \
		-v /var/run/dbus:/var/run/dbus:ro \
		$(DOCKER_IMAGE):latest

docker-compose-up: ## Start with docker-compose
	docker-compose up -d

docker-compose-down: ## Stop docker-compose services
	docker-compose down

clean: ## Clean build artifacts
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)
	rm -rf $(COVERAGE_DIR)
	docker rmi $(DOCKER_IMAGE):$(VERSION) $(DOCKER_IMAGE):latest || true

install: ## Install the binary to $GOPATH/bin
	$(GOCMD) install ./cmd/server

benchmark: ## Run benchmarks
	$(GOTEST) -bench=. -benchmem ./...

security: ## Run security checks
	gosec ./...

check: fmt vet lint test ## Run all checks (format, vet, lint, test)

ci: deps check build ## Run CI pipeline locally

# Development helpers
dev-setup: ## Setup development environment
	$(GOGET) github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	$(GOGET) github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
	$(GOGET) github.com/cosmtrek/air@latest

generate: ## Generate code (if using go:generate)
	$(GOCMD) generate ./...

mod-tidy: ## Tidy go modules
	$(GOMOD) tidy

mod-upgrade: ## Upgrade dependencies
	$(GOGET) -u all
	$(GOMOD) tidy

# API testing
test-api: ## Test API endpoints (requires running server)
	@echo "Testing API endpoints..."
	curl -f http://localhost:8080/buses || echo "Server not running on localhost:8080"

# Database/D-Bus testing
test-dbus: ## Test D-Bus connectivity
	@echo "Testing D-Bus connectivity..."
	dbus-send --system --print-reply --dest=org.freedesktop.DBus /org/freedesktop/DBus org.freedesktop.DBus.ListNames || echo "D-Bus system bus not available"
	dbus-send --session --print-reply --dest=org.freedesktop.DBus /org/freedesktop/DBus org.freedesktop.DBus.ListNames || echo "D-Bus session bus not available"

# Release helpers
tag: ## Create a new git tag (usage: make tag VERSION=v1.0.0)
	git tag -a $(VERSION) -m "Release $(VERSION)"
	git push origin $(VERSION)

release: build-all ## Prepare release artifacts
	@echo "Building release for version $(VERSION)"
	mkdir -p release
	cp $(BUILD_DIR)/* release/
	cd release && sha256sum * > checksums.txt
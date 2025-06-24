# Makefile for Go project
help:
    @echo ""
    @echo "📦 Available commands:"
    @grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'
    @echo ""

run: ## Run the application
    @env $(shell cat local.env | xargs) go run app/main.go

debug: ## Run the application in debug mode
    @env $(shell cat local.env | xargs) dlv debug app/main.go

build: ## Build the application
    go build -o bin/app app/main.go

clean: ## Clean up build artifacts
    rm -rf bin/*

vendor: ## Get vendor
    go mod vendor

tidy: ## Clean up go.mod and go.sum
    go mod tidy

generate: ## Run go generate
    go generate ./...

unit: ## Run unit tests
    go test ./... -cover -short

integration: ## Run integration tests
    go test ./test/integration/... -v

cover: ## Run tests with coverage
    go test ./... -coverprofile=reports/coverage.out
    go tool cover -func=reports/coverage.out

mock: ## Generate mocks
    mockgen -source=repositories/entity/repository.go \
            -destination=test/mocks/entity/repository.go \
            -package=mocks

lint: ## Run linter to check code quality
    golangci-lint run

test-all: ## Run all tests (unit + integration)
    go test ./... -cover -v

docker-build: ## Build Docker image
    docker build -t my-app .

docker-run: ## Run Docker container
    docker run --rm -p 8080:8080 my-app

install-deps: ## Install required dependencies
	go get -u github.com/swaggo/swag/cmd/swag
	go get -u github.com/securego/gosec/v2/cmd/gosec@latest

swagger-gen: ## Generate Swagger documentation using swaggo
    @if ! command -v swag &> /dev/null; then \
        echo "Swaggo is not installed. Run 'make install-deps' first."; \
        exit 1; \
    fi
	swag init -o docs

gosec: ## Run security analysis using gosec
	@if ! command -v gosec &> /dev/null; then \
        echo "Gosec is not installed. Run 'make install-deps' first."; \
        exit 1; \
    fi
	gosec ./...


# Makefile for Go project
help:
	@echo ""
	@echo "📦 Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'
	@echo ""

run: ## Run the application
	@env $(shell cat local.env | xargs) go run main.go

debug: ## Run the application in debug mode
	@env $(shell cat local.env | xargs) dlv debug main.go

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

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
GOTEST=$(GOCMD) test
GOCLEAN=$(GOCMD) clean
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=api

# Docker parameters
DOCKER_COMPOSE=docker-compose

.PHONY: all build clean test coverage run deps docker-up docker-down docker-build lint swagger help

all: clean build

build: ## Build the application
	$(GOBUILD) -o $(BINARY_NAME) ./cmd/api

clean: ## Clean build files
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f coverage.txt

test: ## Run tests
	$(GOTEST) -v ./...

coverage: ## Run tests with coverage
	$(GOTEST) -v -race -coverprofile=coverage.txt -covermode=atomic ./...

run: ## Run the application
	$(GORUN) cmd/api/main.go

deps: ## Download dependencies
	$(GOMOD) download
	$(GOMOD) tidy

docker-up: ## Start docker containers
	$(DOCKER_COMPOSE) up -d

docker-down: ## Stop docker containers
	$(DOCKER_COMPOSE) down

docker-build: ## Build docker images
	$(DOCKER_COMPOSE) build

lint: ## Run linter
	golangci-lint run

swagger: ## Generate swagger documentation
	swag init -g cmd/api/main.go -o docs

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
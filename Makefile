# Makefile for Go Hexagonal Architecture Template

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=go-hexagonal-template
MAIN_PATH=cmd/app/main.go

# Docker parameters
DOCKER_COMPOSE=docker-compose

.PHONY: all build clean test coverage deps run docker-up docker-down docker-restart

all: test build

build:
	$(GOBUILD) -o $(BINARY_NAME) $(MAIN_PATH)

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f coverage.out

test:
	$(GOTEST) ./...

test-verbose:
	$(GOTEST) -v ./...

coverage:
	$(GOTEST) -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out

coverage-func:
	$(GOTEST) -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -func=coverage.out

deps:
	$(GOMOD) download

tidy:
	$(GOMOD) tidy

run:
	$(GOCMD) run $(MAIN_PATH)

docker-up:
	$(DOCKER_COMPOSE) up -d

docker-down:
	$(DOCKER_COMPOSE) down

docker-restart:
	$(DOCKER_COMPOSE) restart

# Create .env file from example if it doesn't exist
.env:
	cp .env.example .env

setup: deps .env
	@echo "Setup complete. Edit .env file with your configuration."
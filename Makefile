# Makefile for Go Hexagonal Architecture Template

GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOMOD=$(GOCMD) mod
BINARY_NAME=go-hexagonal-template
MAIN_PATH=cmd/app/main.go
MIGRATIONS_DIR=internal/adapters/database/migrations

DOCKER_COMPOSE=docker-compose

.PHONY: all build clean test test-verbose coverage coverage-func deps tidy run \
        docker-up docker-down docker-restart \
        migrate-create migrate-up migrate-down migrate-version migrate-force \
        setup

all: test build

# ── Build ────────────────────────────────────────────────────────────────────

build:
	$(GOBUILD) -o $(BINARY_NAME) $(MAIN_PATH)

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f coverage.out

# ── Test ─────────────────────────────────────────────────────────────────────

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

# ── Dependencies ──────────────────────────────────────────────────────────────

deps:
	$(GOMOD) download

tidy:
	$(GOMOD) tidy

# ── Run ───────────────────────────────────────────────────────────────────────

run:
	$(GOCMD) run $(MAIN_PATH)

# ── Docker ────────────────────────────────────────────────────────────────────

docker-up:
	$(DOCKER_COMPOSE) up -d

docker-down:
	$(DOCKER_COMPOSE) down

docker-restart:
	$(DOCKER_COMPOSE) restart

# ── Migrations ────────────────────────────────────────────────────────────────
# Requires the migrate CLI: https://github.com/golang-migrate/migrate/tree/master/cmd/migrate
# Install: go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

DB_URL ?= postgres://$(DB_USER):$(DB_PASS)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSL_MODE)

migrate-create:
	@test -n "$(NAME)" || (echo "Usage: make migrate-create NAME=<migration_name>" && exit 1)
	migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $(NAME)

migrate-up:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" up

migrate-down:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" down 1

migrate-version:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" version

migrate-force:
	@test -n "$(VERSION)" || (echo "Usage: make migrate-force VERSION=<version>" && exit 1)
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" force $(VERSION)

# ── Setup ─────────────────────────────────────────────────────────────────────

.env:
	cp .env.example .env

setup: deps .env
	@echo "Setup complete. Edit .env with your configuration."
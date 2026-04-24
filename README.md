# Go Hexagonal Architecture Template

A production-ready template for building Go applications using Hexagonal Architecture (Ports and Adapters).

[![Go version](https://img.shields.io/badge/Go-≥1.24-blue)](https://go.dev/)
[![codecov](https://codecov.io/gh/DeSouzaRafael/go-hexagonal-template/graph/badge.svg?token=Z1GX03OUB2)](https://codecov.io/gh/DeSouzaRafael/go-hexagonal-template)
[![Go Report Card](https://goreportcard.com/badge/github.com/DeSouzaRafael/go-hexagonal-template)](https://goreportcard.com/report/github.com/DeSouzaRafael/go-hexagonal-template)
[![License](https://img.shields.io/github/license/DeSouzaRafael/go-hexagonal-template)](https://github.com/DeSouzaRafael/go-hexagonal-template/blob/main/LICENSE)

## Overview

This template demonstrates Hexagonal Architecture in Go. The domain has zero framework dependencies — all external concerns (HTTP, database, JWT) are adapters that plug into port interfaces. Swap any adapter without touching business logic.

## Stack

| Concern | Library |
|---------|---------|
| HTTP | [Echo v4](https://echo.labstack.com/) |
| Database | [GORM](https://gorm.io/) + PostgreSQL |
| Auth | [golang-jwt/jwt v5](https://github.com/golang-jwt/jwt) |
| Validation | [go-playground/validator v10](https://github.com/go-playground/validator) |
| API Docs | [Swaggo](https://github.com/swaggo/swag) |

## Architecture

```
cmd/app/
└── main.go                  # Entry point — wires config, DB, container, server

internal/
├── container.go             # Manual DI — wires repos → services → handlers
├── config/                  # Env-based config loading
├── core/
│   ├── domain/              # Entities and sentinel errors (no external deps)
│   ├── port/                # Interfaces (contracts between layers)
│   └── service/             # Business logic, depends only on ports
└── adapters/
    ├── database/
    │   ├── database.go      # GORM + Postgres connection
    │   └── repositories/    # Repository implementations
    └── web/
        ├── handler/         # HTTP handlers (bind → validate → service → JSON)
        ├── middleware/       # JWT validation middleware
        ├── router/          # Route registration
        ├── token/           # JWT token generation (implements TokenGenerator port)
        └── validator/       # Echo validator adapter (go-playground/validator)

pkg/util/                    # Generic utilities (bcrypt wrapper, env helpers)
docs/                        # Auto-generated Swagger docs (swag init)
```

**Dependency rule:** arrows point inward. `domain` imports nothing. `service` imports only `domain` and `port`. Adapters import `port`, never each other.

## Getting Started

### Prerequisites

- Go 1.24+
- PostgreSQL
- Docker (optional)

### Setup

```bash
git clone https://github.com/DeSouzaRafael/go-hexagonal-template.git
cd go-hexagonal-template
make setup
```

Or manually:

```bash
go mod download
cp .env.example .env
```

### Configuration

Edit `.env`:

```env
APP_NAME=go-hexagonal-template
APP_ENV=development

WEB_PORT=8086
WEB_DOMAIN=localhost

JWT_SECRET=your-secret-key
JWT_EXPIRATION=3600

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASS=postgres
DB_NAME=hexagonal
DB_SSL_MODE=disable
DB_LOG_LEVEL=4
```

> `AutoMigrate` only runs when `APP_ENV` is not `production`. Use proper migrations (e.g. [golang-migrate](https://github.com/golang-migrate/migrate)) for production deployments.

### Run

```bash
# With real database
make run

# Without database (in-memory mock — useful for quick local testing)
go run cmd/app/main.go --mock-db
```

### Docker

```bash
make docker-up      # Start app + Postgres
make docker-down    # Stop containers
make docker-restart # Restart containers
```

## API

Base path: `http://localhost:8086/api`

Interactive docs: `http://localhost:8086/swagger/index.html`

### Auth

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| POST | `/v1/auth/register` | — | Create account |
| POST | `/v1/auth/login` | — | Get JWT token |
| GET | `/v1/auth/profile` | Bearer | Current user ID |

### Users

All endpoints require `Authorization: Bearer <token>`.

| Method | Path | Description |
|--------|------|-------------|
| GET | `/v1/users` | List all users |
| GET | `/v1/users?name=<name>` | Find user by name |
| GET | `/v1/users/:id` | Get user by ID |
| PUT | `/v1/users/:id` | Update user |
| DELETE | `/v1/users/:id` | Delete user |

## Testing

```bash
make test           # Run all tests
make test-verbose   # Verbose output
make coverage       # HTML coverage report
make coverage-func  # Function-level coverage
```

## Regenerating Swagger Docs

After modifying handler annotations:

```bash
swag init -g cmd/app/main.go -o docs --parseDependency --parseInternal
```

## License

MIT — see [LICENSE](LICENSE).

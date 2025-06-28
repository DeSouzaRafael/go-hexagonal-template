
# Go Hexagonal Architecture Template

A robust and scalable template for building Go applications using Hexagonal Architecture (also known as Ports and Adapters).

[![Go version](https://img.shields.io/badge/Go-≥1.24-blue)](https://go.dev/)
[![codecov](https://codecov.io/gh/DeSouzaRafael/go-hexagonal-template/graph/badge.svg?token=Z1GX03OUB2)](https://codecov.io/gh/DeSouzaRafael/go-hexagonal-template)
[![Go Report Card](https://goreportcard.com/badge/github.com/DeSouzaRafael/go-hexagonal-template)](https://goreportcard.com/report/github.com/DeSouzaRafael/go-hexagonal-template)
[![License](https://img.shields.io/github/license/evrone/go-clean-template.svg)](https://github.com/DeSouzaRafael/go-hexagonal-template/blob/main/LICENSE)
[![GitHub Release](https://img.shields.io/github/v/release/DeSouzaRafael/go-hexagonal-template)](https://github.com/DeSouzaRafael/go-hexagonal-template/releases/)

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Architecture](#architecture)
- [Project Structure](#project-structure)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
  - [Configuration](#configuration)
  - [Running the Application](#running-the-application)
  - [Docker](#docker)
- [API Documentation](#api-documentation)
- [Testing](#testing)
- [Contributing](#contributing)
- [License](#license)

## Overview

This template provides a foundation for building Go applications using Hexagonal Architecture, which promotes separation of concerns and testability. It includes a complete implementation of user authentication, JWT token generation, and CRUD operations for a user entity.

## Features

- **Hexagonal Architecture**: Clear separation between domain logic and external dependencies
- **JWT Authentication**: Secure API endpoints with JWT tokens
- **PostgreSQL Integration**: Database persistence using GORM
- **RESTful API**: Clean API design with Echo framework
- **Comprehensive Testing**: High test coverage with mocking
- **Docker Support**: Easy deployment with Docker and Docker Compose
- **Configuration Management**: Environment-based configuration
- **Error Handling**: Consistent error handling throughout the application

## Architecture

The application follows the Hexagonal Architecture pattern, which consists of three main layers:

1. **Domain Layer** (`internal/core/domain`): Contains the business entities and logic
2. **Application Layer** (`internal/core/service`): Implements use cases and orchestrates domain objects
3. **Infrastructure Layer** (`internal/adapters`): Provides implementations for external dependencies

The communication between these layers is facilitated through ports (interfaces) defined in the `internal/core/port` package.

### Key Architectural Components:

- **Ports**: Interfaces that define how the application interacts with external systems
- **Adapters**: Implementations of ports that connect to external systems
- **Domain Models**: Business entities that encapsulate business rules
- **Services**: Application logic that orchestrates domain objects
- **Dependency Injection**: Container that wires up all components

## Project Structure

```
.
├── cmd/                  # Application entry points
│   └── app/              # Main application
├── internal/             # Private application code
│   ├── adapters/         # Implementations of ports (infrastructure layer)
│   │   ├── database/     # Database adapters
│   │   │   └── repositories/ # Repository implementations
│   │   └── web/          # Web adapters (HTTP handlers, routers)
│   │       ├── handler/  # HTTP handlers
│   │       ├── middleware/ # HTTP middleware
│   │       ├── router/   # HTTP routers
│   │       └── token/    # Token generation
│   ├── core/             # Core business logic
│   │   ├── domain/       # Domain models and errors
│   │   ├── port/         # Interfaces (ports)
│   │   └── service/      # Application services
│   ├── config/           # Configuration
│   └── container.go      # Dependency injection container
├── pkg/                  # Public libraries
│   └── util/             # Utility functions
├── .env.example          # Example environment variables
├── docker-compose.yml    # Docker Compose configuration
├── Dockerfile            # Docker configuration
└── README.md             # Project documentation
```

## Getting Started

### Prerequisites

- Go 1.24 or higher
- PostgreSQL
- Docker and Docker Compose (optional)

### Installation

1. Clone the repository:

```bash
git clone https://github.com/DeSouzaRafael/go-hexagonal-template.git
cd go-hexagonal-template
```

2. Set up the project using the Makefile:

```bash
make setup
```

This will download dependencies and create a `.env` file from the example.

Alternatively, you can do these steps manually:

```bash
go mod download
cp .env.example .env
```

### Configuration

1. Edit the `.env` file to match your environment:

```
# Application Settings
APP_NAME=go-hexagonal-template
APP_ENV=development

# Web Server Settings
WEB_PORT=8080
WEB_DOMAIN=localhost

# JWT Settings
JWT_SECRET=your-secret-key
JWT_EXPIRATION=3600 # 1 hour

# Database Settings 
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASS=postgres
DB_NAME=hexagonal
DB_SSL_MODE=disable
DB_LOG_LEVEL=4 # 1 = Silent, 2 = Error, 3 = Warn, 4 = Info
```

### Running the Application

To run the application locally:

```bash
make run
```

Or manually:

```bash
go run cmd/app/main.go
```

The API will be available at `http://localhost:8080/api`.

### Docker

To run the application using Docker:

```bash
make docker-up
```

Or manually:

```bash
docker-compose up -d
```

This will start the application and a PostgreSQL database.

Other useful Docker commands:

```bash
make docker-down    # Stop the containers
make docker-restart # Restart the containers
```

## API Documentation

### Authentication Endpoints

- **Register User**
  - `POST /api/v1/auth/register`
  - Request Body: `{ "name": "username", "password": "password" }`

- **Login**
  - `POST /api/v1/auth/login`
  - Request Body: `{ "name": "username", "password": "password" }`
  - Response: `{ "token": "jwt-token" }`

- **Get User Profile**
  - `GET /api/v1/auth/profile`
  - Headers: `Authorization: Bearer jwt-token`

### User Endpoints (Protected)

- **Get All Users**
  - `GET /api/v1/users`
  - Headers: `Authorization: Bearer jwt-token`

- **Get User by ID**
  - `GET /api/v1/users/:id`
  - Headers: `Authorization: Bearer jwt-token`

- **Update User**
  - `PUT /api/v1/users/:id`
  - Headers: `Authorization: Bearer jwt-token`
  - Request Body: `{ "name": "new-username", "password": "new-password" }`

- **Delete User**
  - `DELETE /api/v1/users/:id`
  - Headers: `Authorization: Bearer jwt-token`

## Testing

The project includes a Makefile with several useful commands for testing and development.

To run all tests:

```bash
make test
```

To run tests with verbose output:

```bash
make test-verbose
```

To generate a coverage report and view it in a browser:

```bash
make coverage
```

To see function-level coverage statistics:

```bash
make coverage-func
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

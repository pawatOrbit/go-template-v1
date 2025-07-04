# go-template-v1

A modern, scalable Go project template designed for rapid backend service development. This template provides a robust foundation with best practices for configuration, logging, error handling, HTTP server/client, database integration, and code generation.

## Features

- **Modular Project Structure**: Clean separation of concerns for configuration, core logic, transport, and utilities.
- **Configurable Environments**: Supports multiple profiles (local, dev, sit, stg, prd) with YAML-based configuration.
- **HTTP Server & Middleware**: Built-in HTTP server with middleware stack, routing, and OpenTelemetry tracing.
- **Database Integration**: PostgreSQL support with read/write pools, migrations, and schema management.
- **Code Generation**: Uses `sqlc` for type-safe database access and schema management.
- **Structured Logging**: Canonical, context-aware logging with zap and slog, supporting OpenTelemetry.
- **Error Handling**: Centralized, structured error handling and response formatting.
- **Validation Utilities**: Simple struct validation for request/response models.
- **HTTP Client Utilities**: Typed, traceable HTTP client helpers for external service calls.

## Directory Structure

```
.
├── cmd/            # CLI entrypoints and command setup (Cobra)
├── config/         # Configuration files and examples
├── core/           # Core libraries: config, logger, db, http, validation, etc.
│   ├── config/     # Configuration structs and profile logic
│   ├── exception/  # Error and exception handling
│   ├── httpclient/ # HTTP client utilities and service clients
│   ├── logger/     # Logging setup and canonical logger
│   ├── pgdb/       # PostgreSQL integration, migrations, SQL helpers
│   ├── transport/  # HTTP server, middleware, and routing
│   └── validation/ # Struct validation utilities
├── internal/       # Application-specific business logic, models, services
│   ├── db/         # (empty or for internal DB helpers)
│   ├── model/      # Data models
│   ├── repository/ # Data access layer
│   ├── server/     # HTTP server setup
│   ├── service/    # Business logic/services
│   └── sqlc/       # sqlc-generated code and SQL schemas
├── utils/          # Utility packages (checksum, string conversion, etc.)
├── main.go         # Main entrypoint
├── go.mod          # Go module definition
├── sqlc.yaml       # sqlc configuration for code generation
└── README.md       # Project documentation
```

## Getting Started

### Prerequisites

- Go 1.23+
- PostgreSQL (for database features)
- [sqlc](https://docs.sqlc.dev/en/stable/) (for code generation)

### Configuration

Copy the example config and adjust as needed:

```sh
cp config/example.config.yaml config/config.local.yaml
```

Example config (`config/example.config.yaml`):

```yaml
env: local

restServer:
  port: 8080

cors:
  allowOrigins: ["*"]
  AllowedMethods: ["HEAD", "GET", "POST", "PUT", "PATCH", "DELETE"]
  allowHeaders: ["*"]
  exposeHeaders: ["Accept", "Content-Encoding", ...]
  maxAge: 7200

postgres:
  read:
    host: ""
    port: 5432
    username: ""
    password: ""
    database: ""
    schema: "public"
    maxConnections: 20
  write:
    host: ""
    port: 5432
    username: ""
    password: ""
    database: ""
    schema: "public"
    maxConnections: 20
```

### Running the Server

```sh
go run main.go serve:all-api
```

### Code Generation

To generate type-safe database code from SQL:

```sh
sqlc generate
```

### Project Structure Overview

- **cmd/**: CLI entrypoints using Cobra.
- **core/**: Reusable libraries for config, logging, DB, HTTP, etc.
- **internal/**: Your business logic, models, and services.
- **utils/**: Helper utilities.

### Key Technologies

- [Cobra](https://github.com/spf13/cobra) for CLI
- [Viper](https://github.com/spf13/viper) for configuration
- [Zap](https://github.com/uber-go/zap) and [slog](https://pkg.go.dev/log/slog) for logging
- [OpenTelemetry](https://opentelemetry.io/) for tracing
- [sqlc](https://docs.sqlc.dev/) for database code generation
- [pgx](https://github.com/jackc/pgx) for PostgreSQL driver

## Customization

- Add your own models, repositories, and services under `internal/`.
- Define new HTTP endpoints and middleware in `core/transport/httpserver/`.
- Extend configuration as needed in `config/` and `core/config/`.

## License

MIT
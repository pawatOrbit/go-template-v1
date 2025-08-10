# Go API Template

A modern, production-ready Go project template designed for rapid backend API development. This template provides a robust foundation with enterprise-grade features including authentication, observability, testing, and containerization.

## 🚀 Features

### 🏗️ **Architecture & Structure**
- **Clean Architecture**: Separation of concerns with layered structure (handlers, services, repositories)
- **Modular Design**: Reusable core components and easy-to-extend business logic
- **Type Safety**: Full Go type safety with structured request/response models

### 🔐 **Security & Authentication**
- **JWT Authentication**: Complete JWT-based auth with refresh tokens
- **Role-Based Access Control**: Flexible RBAC system with middleware
- **Request ID Tracking**: Full request tracing with correlation IDs
- **Security Headers**: CORS, rate limiting, and security middleware

### 📊 **Observability & Monitoring**
- **Structured Logging**: Canonical logging with zap/slog and OpenTelemetry integration
- **Health Checks**: Comprehensive health endpoints (liveness, readiness, detailed)
- **Database Monitoring**: Connection pool monitoring and health checks
- **Request Tracing**: OpenTelemetry integration for distributed tracing

### 🗄️ **Database & Persistence**
- **PostgreSQL Integration**: Read/write connection pools with pgx driver
- **Type-Safe Queries**: SQL code generation with sqlc
- **Migration Support**: Database initialization scripts
- **Connection Pooling**: Optimized connection management

### 🧪 **Testing & Quality**
- **Comprehensive Testing**: Unit and integration test suites with testify
- **Test Coverage**: Coverage reporting and quality gates
- **CI/CD Ready**: GitHub Actions workflows (coming soon)

### 🐳 **Containerization & Deployment**
- **Docker Support**: Multi-stage Dockerfile with security best practices
- **Docker Compose**: Development and production-ready compose files
- **Health Checks**: Container health monitoring
- **Development Tools**: Live reload and debugging support

### ⚙️ **Configuration & Environment**
- **YAML Configuration**: Environment-based configuration management
- **Multiple Profiles**: Support for local, dev, staging, production environments
- **Configuration Validation**: Type-safe configuration with validation

## 📁 Project Structure

```
.
├── cmd/                    # CLI entrypoints (Cobra)
├── config/                 # Configuration files
│   ├── example.config.yaml # Example configuration
│   └── config.docker.yaml  # Docker environment config
├── core/                   # Core libraries and shared components
│   ├── auth/               # Authentication service
│   ├── config/             # Configuration management
│   ├── exception/          # Error handling
│   ├── health/             # Health check system
│   ├── httpclient/         # HTTP client utilities
│   ├── logger/             # Structured logging
│   ├── pgdb/               # PostgreSQL integration
│   ├── transport/          # HTTP server and middleware
│   └── validation/         # Request validation
├── internal/               # Application-specific code
│   ├── model/              # Data models and DTOs
│   ├── repository/         # Data access layer
│   ├── server/             # HTTP server setup
│   └── service/            # Business logic
├── tests/                  # Test suites
│   ├── unit/               # Unit tests
│   ├── integration/        # Integration tests
│   └── README.md           # Testing guide
├── migrations/             # Database migrations
│   ├── 001_create_users_table.up.sql
│   ├── 001_create_users_table.down.sql
│   └── README.md           # Migration guide
├── utils/                  # Utility packages
├── scripts/                # Database and deployment scripts
├── Dockerfile              # Container definition
├── docker-compose.yml      # Container orchestration
├── Makefile                # Development commands
└── main.go                 # Application entrypoint
```

## 🚀 Quick Start

### Prerequisites

- **Go 1.23+**
- **Docker & Docker Compose** (recommended)
- **PostgreSQL** (if running locally)
- **Make** (optional, for convenient commands)

### 1. Clone and Setup

```bash
git clone <your-repo-url>
cd go-api-template
```

### 2. Configuration

Copy and customize the configuration:

```bash
cp config/example.config.yaml config/config.local.yaml
# Edit config.local.yaml with your settings
```

### 3. Run with Docker (Recommended)

```bash
# Start all services (API, PostgreSQL, Redis)
make docker-run

# Or for development with tools
make docker-dev
```

### 4. Run Locally

```bash
# Install dependencies
go mod download

# Start database
make db-up

# Run the application
make run
# or
go run main.go serve:all-api
```

## 🔧 Development

### Available Commands

```bash
# Development
make run              # Run application locally
make run-dev          # Run with live reload (requires air)
make build            # Build application
make test             # Run all tests
make test-unit        # Run unit tests only
make coverage         # Generate test coverage report

# Code Quality
make lint             # Run linter
make fmt              # Format code
make security         # Run security scan

# Database
make db-up            # Start database
make db-down          # Stop database
make db-reset         # Reset database
make sqlc-generate    # Generate code from SQL

# Migrations
make migrate          # Run all pending migrations
make migrate-status   # Check migration status
make migrate-create name=migration_name # Create new migration
make migrate-down     # Rollback last migration

# Docker
make docker-build     # Build Docker image
make docker-run       # Run with Docker Compose
make docker-dev       # Run in development mode
make docker-logs      # View logs
```

### Configuration Examples

**Local Development** (`config/config.local.yaml`):
```yaml
env: local
restServer:
  port: 8080
postgres:
  read:
    host: "localhost"
    port: 5432
    username: "postgres"
    password: "your-password"
    database: "your-database"
auth:
  jwtSecretKey: "your-secret-key"
  skipAuthPaths:
    - "/health"
    - "/api/v1/auth/login"
```

**Docker Environment** (`config/config.docker.yaml`):
```yaml
env: docker
restServer:
  port: 8080
postgres:
  read:
    host: "postgres"
    port: 5432
    username: "postgres"
    password: "postgres"
    database: "go_api_template"
```

## 📝 API Endpoints

### Health Checks
- `GET /health` - Comprehensive health check
- `GET /health/liveness` - Liveness probe
- `GET /health/readiness` - Readiness probe

### Example Endpoints (Replace with your APIs)
- `GET /api/v1/examples/{id}` - Get example by ID
- `POST /api/v1/examples` - Create new example

### Authentication
- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/refresh` - Refresh token
- Protected endpoints require `Authorization: Bearer <token>` header

## 🧪 Testing

```bash
# Run all tests
make test

# Run with coverage
make coverage

# Run specific test types
make test-unit
make test-integration

# Run specific test
go test -run TestAuthService ./tests/unit/...
```

## 🗄️ Database Migrations

The template includes a comprehensive migration system using [golang-migrate](https://github.com/golang-migrate/migrate).

### Quick Start with Migrations

```bash
# Start database
make db-up

# Run all pending migrations
make migrate

# Check migration status
make migrate-status
```

### Migration Commands

```bash
# Create a new migration
make migrate-create name=create_posts_table

# Run migrations
make migrate                  # Run all pending
make migrate-up              # Same as above

# Rollback migrations
make migrate-down            # Rollback last migration
make migrate-down-all        # Rollback all (DANGEROUS!)

# Migration status
make migrate-status          # Show current status
make migrate-version         # Show current version

# Advanced (use with caution)
make migrate-force version=2 # Force to specific version
```

### Example Migration Workflow

1. **Create a new migration**:
   ```bash
   make migrate-create name=add_user_preferences
   ```

2. **Edit the generated files**:
   - `migrations/YYYYMMDDHHMMSS_add_user_preferences.up.sql`
   - `migrations/YYYYMMDDHHMMSS_add_user_preferences.down.sql`

3. **Run the migration**:
   ```bash
   make migrate
   ```

The template comes with example migrations for `users`, `api_keys`, and `products` tables. See the [Migration Guide](./migrations/README.md) for detailed documentation and best practices.

## 🐳 Docker Deployment

### Production Deployment
```bash
# Build and run production containers
docker-compose up -d

# Scale API instances
docker-compose up -d --scale api=3

# View logs
docker-compose logs -f api
```

### Development with Docker
```bash
# Run with development tools (pgAdmin, Redis Commander)
docker-compose -f docker-compose.yml -f docker-compose.dev.yml up -d
```

## 🔧 Customization Guide

### Adding New Services

1. **Create Model** in `internal/model/`:
```go
type YourRequest struct {
    Name string `json:"name" validate:"required"`
}
```

2. **Create Service** in `internal/service/`:
```go
type YourService interface {
    DoSomething(ctx context.Context, req *model.YourRequest) (*model.YourResponse, error)
}
```

3. **Register Service** in `internal/service/service.go`:
```go
YourService: NewYourService(repo, errors),
```

4. **Add Routes** in `internal/server/route.go`:
```go
r.Post("/api/v1/your-endpoint", httpserver.NewTransport(
    &model.YourRequest{},
    httpserver.NewEndpoint(service.YourService.DoSomething),
))
```

### Adding Authentication to Endpoints

```go
// In route.go, wrap endpoints with auth middleware
authMiddleware := middleware.AuthMiddleware(authConfig)
protectedHandler := authMiddleware(yourHandler)

// For role-based access
adminOnly := middleware.RequireRoles("admin")(yourHandler)
```

## 🛠️ Technology Stack

- **Framework**: Go 1.23+ with standard library
- **Router**: Native Go HTTP mux with custom routing layer
- **Database**: PostgreSQL with pgx driver
- **Authentication**: JWT with golang-jwt/jwt
- **Validation**: Custom validation with struct tags
- **Logging**: Structured logging with slog and zap
- **Testing**: testify framework
- **Containerization**: Docker with multi-stage builds
- **Configuration**: YAML with Viper
- **CLI**: Cobra for command-line interface
- **Observability**: OpenTelemetry integration

## 📚 Additional Resources

- [Testing Guide](./tests/README.md)
- [Architecture Documentation](./docs/architecture.md) (coming soon)
- [Deployment Guide](./docs/deployment.md) (coming soon)
- [API Documentation](./docs/api.md) (coming soon)

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

**Ready to build something amazing?** 🚀

This template provides everything you need to get started with a production-ready Go API. Simply replace the example services with your business logic and you're ready to deploy!
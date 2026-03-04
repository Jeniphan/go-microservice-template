# Go Microservice Template

A production-ready Go microservice template built with **Echo v5**, **GORM**, and **PostgreSQL**. Designed as a clean, opinionated starting point for building REST API microservices with a layered architecture, structured logging, token-based authentication, and a powerful generic filter system.

## Tech Stack

| Layer | Technology |
|---|---|
| Language | Go 1.25 |
| HTTP Framework | [Echo v5](https://github.com/labstack/echo) |
| ORM | [GORM](https://gorm.io) |
| Database | PostgreSQL (pgx driver) |
| Validation | [go-playground/validator](https://github.com/go-playground/validator) |
| Hot Reload | [Air](https://github.com/air-verse/air) |
| Containerization | Docker (multi-stage, Alpine) |
| CI/CD | GitLab CI + Dokploy |

## Project Structure

```
.
├── cmd/
│   └── main.go                         # Application entrypoint
├── configs/
│   ├── database.config.go              # PostgreSQL connection (GORM)
│   ├── env.config.go                   # .env loader (godotenv)
│   └── validation.config.go            # Request validator setup
├── internal/
│   ├── bootstrap/
│   │   └── app.go                      # Dependency injection & app wiring
│   ├── common/
│   │   ├── error_handlers/
│   │   │   └── http_error_handler.go   # Global HTTP error handler
│   │   ├── filter/                     # Advanced generic GORM filter system
│   │   │   ├── builders/               # Filter builder implementations
│   │   │   │   ├── basic.go            # Basic field filtering
│   │   │   │   ├── nested.go           # Nested relation filtering
│   │   │   │   ├── parent.go           # Parent relation filtering
│   │   │   │   ├── m2m.go              # Many-to-many filtering
│   │   │   │   ├── search.go           # Full-text search
│   │   │   │   ├── range.go            # Date/numeric range filtering
│   │   │   │   ├── sort.go             # Sorting
│   │   │   │   ├── group.go            # Grouping
│   │   │   │   └── pagination.go       # Pagination
│   │   │   ├── utils/                  # Validators, parsers, JSON helpers
│   │   │   ├── filter.go               # AdvanceFilter[T] core
│   │   │   ├── query.go                # Query parameter structs
│   │   │   ├── result.go               # Paginated result struct
│   │   │   ├── interfaces.go           # Filterable interface
│   │   │   └── options.go              # Filter options
│   │   └── loggers/
│   │       └── logger.go               # Structured JSON logger
│   ├── controllers/                    # HTTP handlers (thin layer)
│   ├── middlewares/
│   │   ├── application.middleware.go   # Bearer token auth (calls external auth MS)
│   │   └── request_logging.middleware.go # Request/Response logging middleware
│   ├── models/
│   │   ├── entities/
│   │   │   └── base.entity.go          # Base GORM entity (UUID, timestamps, soft delete)
│   │   └── log.model.go                # Log models
│   ├── routers/
│   │   └── router.go                   # Route registration
│   └── services/                       # Business logic layer
├── .air.toml                           # Air hot-reload configuration
├── .env.example                        # Environment variable template
├── .gitlab-ci.yml                      # GitLab CI/CD pipeline
├── docker-compose.yml                  # Docker Compose for local/prod
├── Dockerfile                          # Multi-stage production build
└── go.mod
```

## Architecture

This template follows a **layered architecture** (Controller → Service → Repository) with manual dependency injection wired in `bootstrap/app.go`.

```
Request
  │
  ▼
Middleware (Auth, Request Logging)
  │
  ▼
Router (Echo)
  │
  ▼
Controller  ──→  Service  ──→  Repository (GORM)
                   │
                   ▼
              AdvanceFilter[T]
```

### Key Design Decisions

- **Dependency Injection** — All dependencies (DB, services, controllers, middlewares) are constructed once in `Bootstrap()` and injected via structs. No global singletons.
- **Interface-based layers** — Each service and controller is defined as an interface, making it easy to swap implementations or write unit tests.
- **Generic filter system** — `AdvanceFilter[T any]` provides a reusable, composable query builder for any GORM entity without duplicating filter logic.
- **Base entity** — All GORM entities embed `entities.Base` to get UUID primary key, `created_at`, `updated_at`, and soft delete (`deleted_at`) automatically.

## Getting Started

### Prerequisites

- Go 1.25+
- PostgreSQL
- [Air](https://github.com/air-verse/air) (for hot reload in development)
- Docker (optional)

### 1. Clone and configure

```bash
git clone <your-repo-url>
cd <project-name>

cp .env.example .env
# Edit .env with your configuration
```

### 2. Environment Variables

```env
# Application
APP_ENV=development
APP_PORT=8080

# PostgreSQL
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_user
DB_PASSWORD=your_password
DB_NAME=your_db
DB_SCHEMA=public

# Authentication Middleware (external auth service)
METHOD_APPLICATION_VERIFY=POST
ENDPOINT_APPLICATION_MS=http://your-auth-service
URL_APPLICATION_VERIFY=/verify
```

### 3. Run with hot reload (development)

```bash
air
```

### 4. Run without hot reload

```bash
go run ./cmd/main.go
```

### 5. Run with Docker

```bash
docker compose up --build
```

## API Endpoints

| Method | Path | Auth | Description |
|---|---|---|---|
| `GET` | `/order/healthcheck` | No | Health check — returns DB connection status |
| `*` | `/order/*` | Bearer Token | All other routes require authentication |

## Middleware

### `RequestLogging`
Applied globally. Captures and logs every request and response as structured JSON including: timestamp, method, URI, headers, request body, response body, status code, and latency.

```json
{
  "timestamp": "2025-01-01T00:00:00Z",
  "service": "RequestLoggingMiddleware",
  "request_id": "abc-123",
  "logging_type": "RequestLogs",
  "method": "GET",
  "uri": "/order/healthcheck",
  "status": 200,
  "latency": "1.2ms"
}
```

### `AuthenticatedToken`
Applied to the `/order` route group. Validates `Authorization: Bearer <token>` by forwarding the token to an external application microservice. On success, the resolved `AppID` is set in the request header for downstream handlers.

## Advanced Filter System

`AdvanceFilter[T]` is a reusable generic query builder that composes multiple filter strategies. Use it in any service to build dynamic, paginated queries.

```go
filter := filter.NewAdvanceFilter[YourEntity](db, &builders.FilterOptions{
    SoftDelete: true,
    AppID:      &appID,
})

result, err := filter.Apply(ctx, &filter.AdvanceFilterQuery{
    FilterBy:    []string{"status"},
    Filter:      []string{"active"},
    SearchBy:    []string{"name"},
    Search:      "john",
    SortBy:      "created_at",
    Sort:        "desc",
    Page:        1,
    PerPage:     20,
})
```

### Supported Filter Types

| Builder | Description |
|---|---|
| `BasicFilter` | Exact/partial match on direct columns |
| `NestedFilter` | Filter by columns on joined relations |
| `ParentFilter` | Filter by parent relation columns |
| `M2MFilter` | Filter through many-to-many join tables |
| `SearchFilter` | ILIKE search across multiple columns |
| `RangeFilter` | Date or numeric range (`start`/`end`) |
| `SortFilter` | Multi-column ordering |
| `GroupFilter` | GROUP BY with group-level sorting |
| `PaginationFilter` | Offset-based pagination with total pages |

### Filter Result

```go
type FilterResult[T any] struct {
    Data      []T   `json:"data"`
    Total     int64 `json:"total"`
    TotalPage int64 `json:"total_page"`
    Page      int   `json:"page"`
    PerPage   int   `json:"per_page"`
}
```

## Base Entity

Embed `entities.Base` in any GORM model to get UUID PK, timestamps, and soft delete for free:

```go
import "order-v2-microservice/internal/models/entities"

type Order struct {
    entities.Base
    AppID  uuid.UUID `gorm:"type:uuid;not null" json:"app_id"`
    Status string    `gorm:"not null"           json:"status"`
}
```

## Adding a New Feature

Follow this pattern when adding a new domain (e.g. `orders`):

1. **Model** — Create `internal/models/entities/order.entity.go` embedding `entities.Base`
2. **Service** — Create `internal/services/order.service.go` with an interface and implementation
3. **Controller** — Create `internal/controllers/order.controller.go` using the service interface
4. **Wire** — Register the service and controller in `internal/bootstrap/app.go`
5. **Route** — Add routes in `internal/routers/router.go`

## Docker

The Dockerfile uses a **multi-stage build**:

- **Stage 1 (builder)** — `golang:1.25-alpine` — compiles a static binary
- **Stage 2 (production)** — `alpine:latest` — minimal image, non-root user, ~15MB final image

```bash
# Build image
docker build -t my-service .

# Run container
docker run -p 8080:8080 --env-file .env my-service
```

## CI/CD (GitLab)

The `.gitlab-ci.yml` defines a 3-stage pipeline:

| Stage | Job | Trigger |
|---|---|---|
| `build` | Compile Go binary | Push to `main`/`dev`, Merge Requests |
| `publish` | Build & push Docker image to GitLab Registry | Push to `main`/`dev` only |
| `deploy` | Trigger [Dokploy](https://dokploy.com) webhook | Push to `main`/`dev` only |

### Required GitLab CI/CD Variables

| Variable | Description |
|---|---|
| `DOKPLOY_WEBHOOK_URL` | Dokploy production webhook URL |
| `DOKPLOY_TOKEN` | Dokploy production token (masked, protected) |
| `DOKPLOY_WEBHOOK_URL_DEV` | Dokploy dev webhook URL |
| `DOKPLOY_TOKEN_DEV` | Dokploy dev token (masked) |
| `DOKPLOY_AUTH_TYPE` | Auth type: `x-gitlab-token`, `basic`, or empty |

GitLab automatically provides `CI_REGISTRY`, `CI_REGISTRY_USER`, and `CI_REGISTRY_PASSWORD`.

## Structured Logging

All application logs are emitted as **JSON** for easy ingestion by log aggregators (Loki, Datadog, etc.):

```json
{
  "time_stamp": "2025-01-01T00:00:00Z",
  "service": "HealthCheckController",
  "request_id": "abc-123",
  "logging_type": "ApiLogs",
  "log_level": "Info",
  "message": "HealthCheck Ctl"
}
```

Use the logger in any service or controller:

```go
var log = loggers.NewLogger("MyService")

log.Info(c, "Processing request", "id", id)
log.Error(c, "Something failed", "error", err)
```

## Customizing This Template

When using this as a base, replace the following:

- [ ] Module name in `go.mod` (`order-v2-microservice` → your service name)
- [ ] Route prefix (`/order` → your service prefix)
- [ ] Docker healthcheck URL in `Dockerfile`
- [ ] Container name in `docker-compose.yml`
- [ ] `.env.example` with your actual environment variables

## License

MIT

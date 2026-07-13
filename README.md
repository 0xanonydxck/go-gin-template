# Golang Gin Template

A modular, testable, and containerized backend for a simple bookstore, built with Go, Gin, GORM, and Docker. The project follows best practices for configuration, authentication, and integration testing.

## Table of Contents

- [Simple Bookstore Backend](#simple-bookstore-backend)
  - [Table of Contents](#table-of-contents)
  - [Features](#features)
  - [Project Structure](#project-structure)
  - [Setup \& Installation](#setup--installation)
  - [Configuration](#configuration)
  - [Running the Application](#running-the-application)
  - [Environment Configuration](#environment-configuration)
  - [Docker Compose Services](#docker-compose-services)
  - [Observability](#observability)
  - [Testing](#testing)
  - [API Documentation](#api-documentation)
  - [Key Components](#key-components)
  - [Contributing](#contributing)
  - [License](#license)

---

## Features

- RESTful API for book and user management
- JWT-based authentication and authorization
- Integration tests with isolated Dockerized PostgreSQL
- Configurable via environment variables
- Modular package structure
- Structured JSON logging with request and trace correlation
- OpenTelemetry tracing with Jaeger support

---

## Project Structure

```
backend/
├── api/               # API route definitions and Swagger docs
├── benchmark/         # (Reserved for benchmarks)
├── cmd/               # Application entrypoints (main.go, etc.)
├── config/            # Configuration loading and management
├── db/                # Database migrations
├── docs/              # Generated API documentation (Swagger)
├── infrastructure/    # Auth, DB, Docker, Logger, etc.
├── internal/          # Business logic (books(with unittest), users(with unittest), middleware, utils)
├── pkg/               # Shared packages (crypto, migration, validator)
├── test/              # Integration tests
├── .env               # Environment variables (local)
├── .example.env       # Example environment variables
├── go.mod, go.sum     # Go dependencies
└── README.md          # Project documentation
```

---

## Setup & Installation

1. **Clone the repository:**
   ```sh
   git clone https://github.com/chai-rs/go-gin-template.git
   cd go-gin-template
   ```

2. **Copy and configure environment variables:**
   ```sh
   cp .example.env .env
   # Edit .env as needed
   ```

3. **Install Go dependencies:**
   ```sh
   go mod download
   ```

4. **(Optional) Install Docker**  
   Required for running integration tests.

---

## Configuration

- All configuration is managed via the `config` package and `.env` file.
- Secrets (e.g., `ACCESS_SECRET`, `REFRESH_SECRET`) are set in the environment and can be overridden in test setup.
- Observability is configured with `LOG_LEVEL`, `OTEL_ENABLED`, `OTEL_SERVICE_NAME`, `OTEL_ENVIRONMENT`, `OTEL_EXPORTER_OTLP_ENDPOINT`, and `OTEL_SAMPLE_RATIO`.

---

## Running the Application

```sh
go run ./cmd/server
```

Or use Docker Compose:

```sh
docker-compose --env-file .docker.env up --build
```

- The application will be available at http://localhost:8000
- Jaeger UI will be available at http://localhost:16686 when tracing is enabled.
- Only the server port is exposed; Redis and Postgres are accessible only within the Docker network.
- Environment variables for Compose are managed in `.docker.env` (see below).

---

## Environment Configuration

- `.env` is for local development (e.g., running with `go run` on your host).
- `.docker.env` is for Docker Compose. Service hostnames must match the Compose service names (`backend-postgres`, `backend-redis`).
- When running the server inside Docker Compose, use `OTEL_EXPORTER_OTLP_ENDPOINT=jaeger:4317`.
- When running the server on your host with Jaeger exposed from Docker, use `OTEL_EXPORTER_OTLP_ENDPOINT=localhost:4317`.
- See `.example.env` for all available variables and their expected values.

---

## Docker Compose Services

- **backend-server**: Your Go app, built from `/cmd/server`.
- **backend-postgres**: PostgreSQL 17, data stored in `pgdata` volume.
- **backend-redis**: Redis 7, data stored in `redisdata` volume.
- **backend-jaeger**: Jaeger all-in-one collector and UI for OpenTelemetry traces.

All services are connected via the `bookstore-network` for secure internal communication.

---

## Observability

The application uses `zerolog` for structured JSON logs. Request logs include:

- `request_id` from `X-Request-ID`, or a generated UUID when the header is missing
- `trace_id` and `span_id` from OpenTelemetry
- `method`, `path`, `route`, `status`, `latency`, `client_ip`, and `user_agent`
- `user_id` for authenticated routes
- `service` and `environment` from the observability configuration

Incoming W3C trace headers are propagated automatically. When `OTEL_ENABLED=true`, Gin requests, GORM queries, and Redis commands are traced and exported to the OTLP endpoint.

To view traces locally with Docker Compose:

```sh
cp .example.env .docker.env
docker-compose --env-file .docker.env up --build
```

Then call an API endpoint and open http://localhost:16686. Select the `simple-bookstore` service to inspect traces.

To add custom spans in application code, use the request context that already flows through handlers, services, and repositories:

```go
ctx, span := otel.Tracer("simple-bookstore").Start(ctx, "operation-name")
defer span.End()
```

---

## Testing

- Integration tests use `testify/suite` and `ory/dockertest` to spin up a PostgreSQL container, run migrations, and clean up automatically.
- Example test suite: `test/base_test.go`, `test/user_test.go`

Run all tests:

```sh
go test ./...
```

---

## API Documentation

- Swagger/OpenAPI docs are available in `docs/` and generated from code comments.
- To view locally, run the server and visit `/swagger/index.html` (if route is enabled).

---

## Key Components

- **Authentication:** JWT-based, with token generation and verification in `infrastructure/auth`.
- **Authorization:** Enforced via Casbin with Gorm adapter, configured in `auth_model.conf`.
- **Database:** GORM ORM with migrations in `db/migrations`.
- **Observability:** JSON logging in `infrastructure/logger`, tracing setup in `infrastructure/tracing`, and request correlation middleware in `internal/middleware`.
- **Testing:** Uses Dockerized PostgreSQL for isolation, see `BaseSuite` in `test/base_test.go`.
- **Handlers & Services:** Business logic in `internal/`, separated by domain (books, users).

---

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/fooBar`)
3. Commit your changes
4. Push to the branch (`git push origin feature/fooBar`)
5. Open a pull request

---

## License

MIT License

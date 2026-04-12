# school-api

A RESTful API for managing school students, built with Go. Follows clean layered architecture using Gin, PostgreSQL, GORM, and zerolog.

## Stack

- **Language:** Go 1.23
- **Web Framework:** [Gin](https://github.com/gin-gonic/gin)
- **Database:** PostgreSQL via [GORM](https://gorm.io) + [pgx](https://github.com/jackc/pgx)
- **Logging:** [zerolog](https://github.com/rs/zerolog)
- **Containerization:** Docker + Docker Compose

## Project Structure

```
school-api/
├── main.go                     # Entry point, graceful shutdown
├── internal/
│   ├── config/config.go        # Environment-based configuration
│   ├── db/postgres.go          # GORM connection, AutoMigrate
│   ├── handlers/students.go    # HTTP request handlers
│   ├── models/student.go       # Data models with GORM tags
│   └── server/server.go        # Router, middleware, server lifecycle
└── pkg/
    └── logger/logger.go        # Structured logger setup
```

## Endpoints

| Method   | Path                    | Description         |
|----------|-------------------------|---------------------|
| `POST`   | `/api/v1/students`      | Add a new student   |
| `GET`    | `/api/v1/students`      | List all students   |
| `GET`    | `/api/v1/students/:id`  | Get student by ID   |
| `DELETE` | `/api/v1/students/:id`  | Delete a student    |
| `GET`    | `/healthz`              | Health check        |

## Getting Started

### Prerequisites

- [Docker](https://www.docker.com/) and Docker Compose
- Go 1.23+ (for local development without Docker)

### Run with Docker Compose

```bash
docker compose up --build
```

This starts both the API (port `8080`) and PostgreSQL (port `5432`).

### Run Locally

1. Copy the example env file and edit as needed:
   ```bash
   cp .env.example .env
   ```

2. Start PostgreSQL (requires Docker):
   ```bash
   docker compose up postgres
   ```

3. Run the API:
   ```bash
   make run
   ```

The database schema is created automatically on startup via GORM `AutoMigrate` — no manual migrations needed.

## Configuration

All configuration is via environment variables.

| Variable      | Default     | Description                |
|---------------|-------------|----------------------------|
| `ENVIRONMENT` | `local`     | Deployment environment     |
| `PORT`        | `8080`      | HTTP server port           |
| `LOG_LEVEL`   | `info`      | Logging level              |
| `DB_HOST`     | `localhost` | PostgreSQL host            |
| `DB_PORT`     | `5432`      | PostgreSQL port            |
| `DB_NAME`     | `school`    | PostgreSQL database name   |
| `DB_USER`     | *(empty)*   | PostgreSQL username        |
| `DB_PASSWORD` | *(empty)*   | PostgreSQL password        |

## API Usage

### Add a Student

```bash
curl -X POST http://localhost:8080/api/v1/students \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Alice Johnson",
    "email": "alice@school.com",
    "age": 15,
    "grade": "10th"
  }'
```

### Get All Students

```bash
curl http://localhost:8080/api/v1/students
```

### Get Student by ID

```bash
curl http://localhost:8080/api/v1/students/<id>
```

### Delete a Student

```bash
curl -X DELETE http://localhost:8080/api/v1/students/<id>
```

## Development

```bash
make build        # Compile binary to bin/
make test         # Run tests with coverage
make lint         # Run golangci-lint
make tidy         # Run go mod tidy
make docker-up    # Start all services
make docker-down  # Stop and remove all services
```

## License

MIT

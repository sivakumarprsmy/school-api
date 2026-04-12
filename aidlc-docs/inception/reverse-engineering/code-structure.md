# Code Structure

## Build System
- **Type**: Go Modules
- **Configuration**: `go.mod` (module: `github.com/periasamy/school-api`, Go 1.23)
- **Key Build Command**: `go build ./...`
- **Run Command**: `go run main.go`
- **Test Command**: `go test ./...`

## Project Layout

```
school-api/
+-- main.go                            # Entry point — wires and starts the app
+-- go.mod                             # Go module definition and dependencies
+-- go.sum                             # Dependency checksums
+-- internal/
|   +-- config/
|   |   +-- config.go                  # Env-var based config loader
|   |   +-- config_test.go             # Unit tests: defaults, env override, DSN
|   +-- db/
|   |   +-- postgres.go                # PostgreSQL connection + GORM auto-migrate
|   +-- handlers/
|   |   +-- students.go                # HTTP handlers for student CRUD
|   |   +-- students_test.go           # Unit tests: all handlers via sqlmock + httptest
|   +-- models/
|   |   +-- student.go                 # Student struct + CreateStudentRequest
|   +-- server/
|       +-- server.go                  # Gin router setup, middleware, HTTP lifecycle
+-- pkg/
    +-- logger/
        +-- logger.go                  # zerolog logger factory
```

## Existing Files Inventory

- [main.go](../../../../../main.go) — Application entry point; wires config, logger, DB, and HTTP server; handles OS signals for graceful shutdown
- [internal/config/config.go](../../../../../internal/config/config.go) — Loads runtime configuration from environment variables; builds PostgreSQL DSN
- [internal/config/config_test.go](../../../../../internal/config/config_test.go) — 4 unit tests for config loading (defaults, env vars, invalid port, DSN format)
- [internal/db/postgres.go](../../../../../internal/db/postgres.go) — Opens GORM/PostgreSQL connection; runs AutoMigrate for Student model
- [internal/server/server.go](../../../../../internal/server/server.go) — Configures Gin router, registers middleware (recovery, gzip, request logger), defines routes, manages HTTP server start/shutdown
- [internal/handlers/students.go](../../../../../internal/handlers/students.go) — Implements CreateStudent, GetStudentByID, GetAllStudents, DeleteStudent HTTP handlers
- [internal/handlers/students_test.go](../../../../../internal/handlers/students_test.go) — 14 unit tests covering all handler paths; uses sqlmock + httptest + testify
- [internal/models/student.go](../../../../../internal/models/student.go) — Student GORM model and CreateStudentRequest with binding validation
- [pkg/logger/logger.go](../../../../../pkg/logger/logger.go) — Creates a zerolog.Logger with console writer, level, timestamp, and caller

## Key Classes/Modules

```
Config
+-- Environment string
+-- Port int
+-- LogLevel string
+-- DBHost, DBPort, DBName, DBUser, DBPassword string
+-- Load() (*Config, error)
+-- PostgresDSN() string

Student (GORM Model)
+-- ID uint (PK, auto-increment)
+-- Name string (not null)
+-- Email string (unique index, not null)
+-- Age int (not null)
+-- Grade string (not null)
+-- CreatedAt time.Time
+-- UpdatedAt time.Time

CreateStudentRequest
+-- Name string (required)
+-- Email string (required, email format)
+-- Age int (required, 1-120)
+-- Grade string (required)

StudentHandler
+-- db *gorm.DB
+-- log zerolog.Logger
+-- CreateStudent(c *gin.Context)
+-- GetStudentByID(c *gin.Context)
+-- GetAllStudents(c *gin.Context)
+-- DeleteStudent(c *gin.Context)
```

## Design Patterns

### Handler Pattern
- **Location**: `internal/handlers/students.go`
- **Purpose**: Separates HTTP concerns (request parsing, response writing) from business logic
- **Implementation**: `StudentHandler` struct holds dependencies (db, log); each method is a `gin.HandlerFunc`

### Repository-less Direct DB Access
- **Location**: `internal/handlers/students.go`
- **Purpose**: Simple direct GORM usage without a separate repository layer (appropriate for this scale)
- **Implementation**: Handlers call `h.db.WithContext(...).Create/Find/First/Delete`

### Graceful Shutdown
- **Location**: `main.go`, `internal/server/server.go`
- **Purpose**: Clean server shutdown on SIGTERM/SIGINT
- **Implementation**: `signal.NotifyContext` + `http.Server.Shutdown` with 5-second timeout

### Environment-based Configuration
- **Location**: `internal/config/config.go`
- **Purpose**: Twelve-factor app style configuration
- **Implementation**: `os.Getenv` with defaults; DSN built from discrete fields

### Test Helper Functions
- **Location**: `internal/handlers/students_test.go`
- **Purpose**: Reduce test boilerplate and ensure consistent test setup
- **Implementation**: `newTestDB(t)` creates a sqlmock-backed GORM DB; `newTestRouter(db)` wires a Gin test router; both use `t.Helper()` and `t.Cleanup()`

## Critical Dependencies

### gin-gonic/gin v1.9.1
- **Version**: v1.9.1
- **Usage**: HTTP router and middleware framework
- **Purpose**: Request routing, parameter binding, JSON responses, recovery middleware

### gorm.io/gorm v1.25.9 + gorm.io/driver/postgres v1.5.7
- **Version**: v1.25.9 / v1.5.7
- **Usage**: ORM for PostgreSQL; used in db package and handlers
- **Purpose**: Database connectivity, auto-migration, CRUD operations

### rs/zerolog v1.33.0
- **Version**: v1.33.0
- **Usage**: Structured logging in logger package and throughout handlers/server
- **Purpose**: Fast, structured, levelled JSON/console logging

### gin-contrib/gzip v1.0.1
- **Version**: v1.0.1
- **Usage**: Gin middleware in server.go
- **Purpose**: Automatic gzip compression of HTTP responses

### stretchr/testify v1.9.0
- **Version**: v1.9.0
- **Usage**: Test assertions in all `_test.go` files
- **Purpose**: Fluent assert/require helpers for test readability

### DATA-DOG/go-sqlmock v1.5.2
- **Version**: v1.5.2
- **Usage**: Mock SQL driver in handler tests
- **Purpose**: Simulate PostgreSQL responses without a real database

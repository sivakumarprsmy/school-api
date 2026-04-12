# System Architecture

## System Overview

School API is a single-binary Go HTTP service that exposes a RESTful API for student record management. It uses the Gin web framework for routing, GORM as an ORM layer over PostgreSQL, and zerolog for structured logging. The service is configured entirely through environment variables and supports graceful shutdown via OS signals.

## Architecture Diagram

```
+------------------+
|   API Client     |
| (HTTP requests)  |
+--------+---------+
         |
         | HTTP (port 8080)
         v
+--------+------------------------------------------+
|                   school-api                      |
|                                                   |
|  +--------------------------------------------+  |
|  |          Gin HTTP Router                   |  |
|  |  GET  /healthz                             |  |
|  |  POST /api/v1/students                     |  |
|  |  GET  /api/v1/students                     |  |
|  |  GET  /api/v1/students/:id                 |  |
|  |  DELETE /api/v1/students/:id               |  |
|  +--------------------+-----------------------+  |
|                       |                          |
|  +--------------------v-----------------------+  |
|  |         StudentHandler                     |  |
|  |  CreateStudent / GetAllStudents            |  |
|  |  GetStudentByID / DeleteStudent            |  |
|  +--------------------+-----------------------+  |
|                       |                          |
|  +--------------------v-----------------------+  |
|  |     GORM ORM (gorm.io/gorm)                |  |
|  |  AutoMigrate, Create, First, Find, Delete  |  |
|  +--------------------+-----------------------+  |
|                                                   |
|  +------+  +----------+  +---------+             |
|  |config|  |  logger  |  | server  |             |
|  +------+  +----------+  +---------+             |
+---------------------------+-----------------------+
                            |
                            | PostgreSQL wire protocol
                            v
               +------------+------------+
               |        PostgreSQL        |
               |    (students table)      |
               +-------------------------+
```

## Component Descriptions

### main.go
- **Purpose**: Application entry point
- **Responsibilities**: Signal handling, wires config + logger + DB + server together, graceful shutdown
- **Dependencies**: config, db, server, logger packages
- **Type**: Application

### internal/server
- **Purpose**: HTTP server setup and routing
- **Responsibilities**: Gin router configuration, middleware registration (recovery, gzip, request logging), route definitions, HTTP server lifecycle
- **Dependencies**: config, handlers, gin framework
- **Type**: Application

### internal/handlers
- **Purpose**: HTTP request handlers
- **Responsibilities**: Request binding and validation, business logic orchestration, response serialization, error mapping to HTTP status codes
- **Dependencies**: models, gorm.DB, zerolog
- **Type**: Application

### internal/models
- **Purpose**: Data model definitions
- **Responsibilities**: Student struct (GORM tags, JSON tags), CreateStudentRequest validation rules
- **Dependencies**: time (stdlib)
- **Type**: Model

### internal/db
- **Purpose**: Database connectivity and migration
- **Responsibilities**: PostgreSQL connection via GORM, auto-migration of Student schema
- **Dependencies**: gorm, gorm/driver/postgres, models
- **Type**: Application (infrastructure concern)

### internal/config
- **Purpose**: Application configuration
- **Responsibilities**: Env var loading with defaults, PostgreSQL DSN construction, validation
- **Dependencies**: os, strconv (stdlib)
- **Type**: Application

### pkg/logger
- **Purpose**: Structured logger factory
- **Responsibilities**: Create zerolog.Logger with console output, configurable level, timestamp, caller info
- **Dependencies**: zerolog, rs/zerolog
- **Type**: Application (shared utility)

## Data Flow

```
Sequence: Create Student

Client          Gin Router      StudentHandler      GORM         PostgreSQL
  |                |                 |                |               |
  |--POST /api/v1/students---------->|                |               |
  |                |---ShouldBindJSON|                |               |
  |                |                 |--db.Create()--->|               |
  |                |                 |                |--INSERT------->|
  |                |                 |                |<--OK-----------|
  |                |                 |<--student obj--|               |
  |<--201 Created (student JSON)-----|                |               |
```

## Integration Points

- **External APIs**: None
- **Databases**: PostgreSQL (via GORM) — stores student records in the `students` table
- **Third-party Services**: None

## Infrastructure Components

- **CDK Stacks**: None
- **Deployment Model**: Single binary Go service, runs on any host/container with PostgreSQL accessible
- **Networking**: Listens on configurable port (default 8080), connects to PostgreSQL on configurable host/port

# Component Inventory

## Application Packages

| Package | Path | Purpose |
|---------|------|---------|
| main | `main.go` | Entry point; wires all components and handles OS signals |
| server | `internal/server` | Gin HTTP router setup, middleware, route registration, server lifecycle |
| handlers | `internal/handlers` | HTTP handlers for student CRUD operations |
| models | `internal/models` | Student data model and request DTOs |
| db | `internal/db` | PostgreSQL connection via GORM, auto-migration |
| config | `internal/config` | Environment-variable based configuration loading |
| logger | `pkg/logger` | zerolog structured logger factory |

## Infrastructure Packages
- None — no CDK, Terraform, or CloudFormation found

## Shared Packages

| Package | Path | Type | Purpose |
|---------|------|------|---------|
| logger | `pkg/logger` | Utility | Reusable structured logger factory |

## Test Packages

| Package | Path | Type | Purpose |
|---------|------|------|---------|
| handlers_test | `internal/handlers/students_test.go` | Unit | Tests all 4 student HTTP handlers using sqlmock + httptest; 14 test cases |
| config_test | `internal/config/config_test.go` | Unit | Tests config loading with defaults, env vars, invalid inputs, DSN; 4 test cases |

## Total Count
- **Total Packages**: 9
- **Application**: 6 (main, server, handlers, models, db, config)
- **Infrastructure**: 0
- **Shared/Utility**: 1 (logger)
- **Test**: 2 (handlers_test, config_test)

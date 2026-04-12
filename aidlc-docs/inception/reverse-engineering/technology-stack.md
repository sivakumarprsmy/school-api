# Technology Stack

## Programming Languages
| Language | Version | Usage |
|----------|---------|-------|
| Go | 1.23 (go.mod) | Entire application and tests |

## Frameworks
| Framework | Version | Purpose |
|-----------|---------|---------|
| gin-gonic/gin v1.9.1 | v1.9.1 | HTTP web framework — routing, middleware, request binding |
| gorm.io/gorm | v1.25.9 | ORM — database abstraction, auto-migration, CRUD |
| gorm.io/driver/postgres | v1.5.7 | PostgreSQL driver for GORM |
| gin-contrib/gzip | v1.0.1 | Gin middleware for gzip response compression |
| rs/zerolog | v1.33.0 | Structured, levelled logging |

## Infrastructure
| Service | Purpose |
|---------|---------|
| PostgreSQL | Primary data store for student records |

## Build Tools
| Tool | Version | Purpose |
|------|---------|---------|
| Go toolchain | 1.23 | Build, test, module management |
| go modules | N/A | Dependency management (go.mod / go.sum) |

## Testing Tools
| Tool | Version | Purpose |
|------|---------|---------|
| github.com/stretchr/testify | v1.9.0 | Assertions (`assert`, `require`) for unit tests |
| github.com/DATA-DOG/go-sqlmock | v1.5.2 | SQL mock driver — simulates DB in handler unit tests |
| net/http/httptest | stdlib | HTTP recorder for testing Gin handlers without a real server |
| testing | stdlib | Go standard test runner |

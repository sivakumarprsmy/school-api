# Dependencies

## Internal Dependencies

```
main.go
  +-- internal/config    (loads Config)
  +-- internal/db        (creates *gorm.DB)
  +-- internal/server    (starts HTTP server)
  +-- pkg/logger         (creates zerolog.Logger)

internal/server
  +-- internal/config    (reads Port, Environment)
  +-- internal/handlers  (creates StudentHandler)

internal/handlers
  +-- internal/models    (Student, CreateStudentRequest)

internal/db
  +-- internal/models    (Student — for AutoMigrate)

internal/models
  (no internal deps)

internal/config
  (no internal deps)

pkg/logger
  (no internal deps)

--- Test Files ---

internal/handlers/students_test.go
  +-- internal/handlers  (StudentHandler under test)
  +-- internal/models    (Student, CreateStudentRequest)

internal/config/config_test.go
  +-- internal/config    (Config under test)
```

### internal/server depends on internal/handlers
- **Type**: Compile
- **Reason**: Server registers StudentHandler routes

### internal/handlers depends on internal/models
- **Type**: Compile
- **Reason**: Handlers use Student and CreateStudentRequest types

### internal/db depends on internal/models
- **Type**: Compile
- **Reason**: AutoMigrate requires the Student model struct

### main depends on internal/{config,db,server} and pkg/logger
- **Type**: Compile
- **Reason**: Wiring all application components together

## External Dependencies

### Production
| Dependency | Version | Purpose | License |
|------------|---------|---------|---------|
| github.com/gin-gonic/gin | v1.9.1 | HTTP framework | MIT |
| gorm.io/gorm | v1.25.9 | ORM | MIT |
| gorm.io/driver/postgres | v1.5.7 | PostgreSQL GORM driver | MIT |
| github.com/gin-contrib/gzip | v1.0.1 | Gzip middleware for Gin | MIT |
| github.com/rs/zerolog | v1.33.0 | Structured logger | MIT |
| github.com/jackc/pgx/v5 | v5.5.5 | PostgreSQL wire driver (indirect) | MIT |

### Test-only
| Dependency | Version | Purpose | License |
|------------|---------|---------|---------|
| github.com/stretchr/testify | v1.9.0 | Test assertions (assert + require) | MIT |
| github.com/DATA-DOG/go-sqlmock | v1.5.2 | SQL mock driver for unit tests | BSD-style |

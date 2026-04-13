# school-api

A RESTful API for managing school students, built with Go. Follows clean layered architecture using Gin, PostgreSQL, GORM, and zerolog. Includes an MCP server so AI assistants can interact with the API directly.

## Stack

- **Language:** Go 1.23
- **Web Framework:** [Gin](https://github.com/gin-gonic/gin)
- **Database:** PostgreSQL via [GORM](https://gorm.io) + [pgx](https://github.com/jackc/pgx)
- **Logging:** [zerolog](https://github.com/rs/zerolog)
- **MCP:** [mcp-go](https://github.com/mark3labs/mcp-go)
- **AI Agent:** [Anthropic Go SDK](https://github.com/anthropics/anthropic-sdk-go)
- **Containerization:** Docker + Docker Compose

## Project Structure

```
school-api/
├── main.go                          # Entry point, graceful shutdown
├── cmd/
│   ├── mcp/main.go                  # MCP server entry point
│   └── agent/main.go                # Enrollment agent (Anthropic SDK)
├── internal/
│   ├── config/config.go             # Environment-based configuration
│   ├── db/postgres.go               # GORM connection, AutoMigrate
│   ├── handlers/students.go         # HTTP request handlers
│   ├── mcptools/client.go           # HTTP client (shared by MCP + agent)
│   ├── mcptools/validator.go        # Email domain validator
│   ├── models/student.go            # Data models with GORM tags
│   └── server/server.go             # Router, middleware, server lifecycle
└── pkg/
    └── logger/logger.go             # Structured logger setup
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

## MCP Server

The MCP server exposes the school-api as tools for AI assistants (Claude Code, Claude Desktop, GitHub Copilot, Cursor, etc.).

### Available Tools

| Tool               | Description                    |
|--------------------|--------------------------------|
| `list_students`    | List all students              |
| `get_student`      | Get a student by ID            |
| `add_student`      | Add a new student              |
| `delete_student`   | Delete a student by ID         |

### Build the MCP binary

```bash
make build-mcp
# or: go build -o bin/school-mcp ./cmd/mcp
```

### Claude Code (VS Code extension)

Register the server via the Claude Code CLI:

```bash
claude mcp add school-api \
  --transport stdio \
  /path/to/school-api/bin/school-api-mcp \
  --env SCHOOL_API_URL=http://localhost:8080
```

Then reload the VS Code window (`Cmd+Shift+P` → `Developer: Reload Window`) and type `/mcp` in the chat to confirm it's connected.

### Claude Desktop

Add to `~/Library/Application Support/Claude/claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "school-api": {
      "type": "stdio",
      "command": "/path/to/school-api/bin/school-api-mcp",
      "env": {
        "SCHOOL_API_URL": "http://localhost:8080"
      }
    }
  }
}
```

Restart Claude Desktop.

### GitHub Copilot (VS Code)

Add to `.vscode/mcp.json` in your project:

```json
{
  "servers": {
    "school-api": {
      "type": "stdio",
      "command": "/path/to/school-api/bin/school-api-mcp",
      "env": {
        "SCHOOL_API_URL": "http://localhost:8080"
      }
    }
  }
}
```

Switch to **Agent mode** in Copilot Chat to use the tools.

### Usage examples

Once connected, just talk naturally to your AI assistant:

```
Add a student named Alice, email alice@school.com, age 15, grade 10th
Show me all students
Get student with ID 1
Delete student with ID 2
```

The `SCHOOL_API_URL` environment variable controls which instance of the API the MCP server talks to. Defaults to `http://localhost:8080`.

## Enrollment Agent

The enrollment agent is a standalone CLI program powered by the Anthropic API. It accepts a natural-language enrollment request via stdin, reasons about it, and enrolls the student — handling edge cases automatically.

### What it does

| Situation                     | Agent behaviour                                          |
|-------------------------------|----------------------------------------------------------|
| All fields present, valid email, no duplicate | Enrolls and confirms                   |
| Missing name / age / grade    | Tells you what's missing                                 |
| Invalid email format          | Reports the error, asks for a corrected email            |
| Email domain not in whitelist | Names the accepted domains, asks for a corrected email   |
| Student already enrolled      | Informs you and does **not** create a duplicate          |

### Build

```bash
make build-agent
# or: go build -o bin/enrollment-agent ./cmd/agent
```

### Run

```bash
# Piped input (typical usage)
echo "Enroll Alice Johnson, alice@school.com, age 15, grade 10th" | \
  ANTHROPIC_API_KEY=<key> ALLOWED_EMAIL_DOMAINS=school.com ./bin/enrollment-agent

# Interactive (no pipe — agent prompts you)
ANTHROPIC_API_KEY=<key> ./bin/enrollment-agent
```

### Environment variables

| Variable                | Required | Description                                         |
|-------------------------|----------|-----------------------------------------------------|
| `ANTHROPIC_API_KEY`     | Yes      | Your Anthropic API key                              |
| `SCHOOL_API_URL`        | No       | School API base URL (default: `http://localhost:8080`) |
| `ALLOWED_EMAIL_DOMAINS` | No       | Comma-separated allowed domains (e.g. `school.com,students.school.com`). Leave empty to allow any domain. |

### Example output

```
✅ Alice Johnson has been successfully enrolled!

  ID:    42
  Name:  Alice Johnson
  Email: alice@school.com
  Age:   15
  Grade: 10th
```

## Development

```bash
make build        # Compile REST API binary to bin/
make build-mcp    # Compile MCP server binary to bin/
make build-agent  # Compile enrollment agent binary to bin/
make test         # Run tests with coverage
make lint         # Run golangci-lint
make tidy         # Run go mod tidy
make docker-up    # Start all services
make docker-down  # Stop and remove all services
```

## License

MIT

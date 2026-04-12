# Requirements — Update Student Endpoint

## Intent Analysis

| Attribute | Value |
|-----------|-------|
| **User Request** | Implement the HLD: add a PUT endpoint to update student information |
| **HLD Source** | https://purpleblockinc.atlassian.net/wiki/spaces/~712020ee70bdfcd33a4731ad1f741f1fbb6c3f/pages/393217/HLD+School+API |
| **Request Type** | New Feature |
| **Scope** | Single Component — `internal/handlers`, `internal/models`, `internal/server` |
| **Complexity** | Simple — one new endpoint following existing patterns |

---

## Functional Requirements

### FR-01: New HTTP Endpoint
The system MUST expose a new HTTP endpoint:
- **Method**: `PUT`
- **Path**: `/api/v1/students/:id`
- **Purpose**: Update an existing student's information by ID

### FR-02: Partial Update Semantics
The endpoint MUST use PATCH-like update semantics:
- Only fields present in the request body are updated
- Fields omitted from the request body remain unchanged
- An empty request body is valid (no-op update)

### FR-03: Updatable Fields
The following student fields MUST be updatable:
- `name` (string)
- `email` (string)
- `age` (int)
- `grade` (string)

The fields `id`, `created_at`, and `updated_at` MUST NOT be updatable via this endpoint.

### FR-04: Input Validation
The request body MUST be validated:
- `name`: if present, must be non-empty string
- `email`: if present, must be a valid email format
- `age`: if present, must be an integer between 1 and 120
- `grade`: if present, must be non-empty string
- Invalid field values MUST return `400 Bad Request`

### FR-05: Student Existence Check
The endpoint MUST verify the student exists before updating:
- If the student ID is not a valid positive integer → `400 Bad Request`
- If no student with the given ID exists → `404 Not Found` with `{"error": "student not found"}`

### FR-06: Email Uniqueness
If `email` is provided in the request body:
- The system MUST verify the new email is not already used by a **different** student
- If the email is already taken → `409 Conflict` with `{"error": "student with this email already exists"}`
- Updating a student with their own existing email MUST NOT trigger a conflict

### FR-07: Success Response
On successful update:
- HTTP status: `200 OK`
- Body: full updated `Student` object in JSON format (same shape as `GET /api/v1/students/:id`)

### FR-08: Route Registration
The new route MUST be registered in `internal/server/server.go` within the existing `/api/v1/students` route group alongside other student routes.

### FR-09: UpdateStudentRequest DTO
A new `UpdateStudentRequest` struct MUST be defined in `internal/models/student.go`:
- All fields optional (pointer types or omitempty)
- Validation tags applied only to fields when present

---

## Non-Functional Requirements

### NFR-01: Consistency with Existing Patterns
The implementation MUST follow the same patterns as existing handlers:
- Use `h.db.WithContext(c.Request.Context())` for all database operations
- Use `h.log.Error()` for server-side error logging
- Use `c.JSON(...)` for all responses
- Use `strconv.ParseUint` for ID parsing (consistent with GetStudentByID / DeleteStudent)

### NFR-02: Test Coverage
Unit tests MUST be added to `internal/handlers/students_test.go` covering:
- `TestUpdateStudent_Success` — valid partial update, 200 OK with updated student
- `TestUpdateStudent_NotFound` — non-existent ID, 404 response
- `TestUpdateStudent_InvalidID` — non-numeric ID param, 400 response
- `TestUpdateStudent_ValidationError` — invalid email or out-of-range age, 400 response
- `TestUpdateStudent_DuplicateEmail` — conflicting email, 409 response
- `TestUpdateStudent_DBError` — database failure on update, 500 response

Tests MUST use the existing `newTestDB` / `newTestRouter` helpers. `newTestRouter` MUST be extended to register the new PUT route.

### NFR-03: Logging
The handler MUST log database errors using the structured logger (zerolog), consistent with existing handlers:
```go
h.log.Error().Err(err).Uint64("id", id).Msg("failed to update student")
```

### NFR-04: No Breaking Changes
The new endpoint MUST NOT modify any existing handler, model, or route. All changes are purely additive.

---

## Security Requirements

### SEC-01: Input Validation (SECURITY-05)
- All request body fields MUST be validated before processing (type, format, range)
- GORM parameterized queries MUST be used — no raw SQL string concatenation
- String inputs MUST be validated via binding tags

### SEC-02: Error Message Safety (SECURITY-09, SECURITY-15)
- Error responses MUST NOT expose database internals, stack traces, or framework details
- All error responses MUST use generic user-facing messages consistent with existing handlers

### SEC-03: Structured Logging — No PII in Logs (SECURITY-03)
- The handler MUST NOT log student PII (name, email) in error or info log entries
- Only non-sensitive metadata (student ID, HTTP status) may appear in logs

---

## Property-Based Testing Requirements

### PBT-01: Testable Properties to Identify at Functional Design
The following properties MUST be analysed and documented during Functional Design:
- **Idempotence (PBT-04)**: `PUT /api/v1/students/:id` with the same payload twice should produce the same final state
- **Invariant (PBT-03)**: Fields not included in the update request MUST remain unchanged
- **Round-trip (PBT-02)**: A student written via PUT and immediately read via GET MUST return the same field values

### PBT-09: Framework Selection
- PBT framework for Go: **`pgregory.net/rapid`** — lightweight, idiomatic, integrates with `testing.T`
- Must be added to `go.mod` as a test dependency

---

## Acceptance Criteria

| ID | Criterion |
|----|-----------|
| AC-01 | `PUT /api/v1/students/1` with `{"name": "Bob"}` updates only the name; email, age, grade are unchanged |
| AC-02 | `PUT /api/v1/students/999` returns `404 Not Found` |
| AC-03 | `PUT /api/v1/students/abc` returns `400 Bad Request` |
| AC-04 | `PUT /api/v1/students/1` with `{"email": "invalid"}` returns `400 Bad Request` |
| AC-05 | `PUT /api/v1/students/1` with `{"email": "other@student.com"}` where that email belongs to student 2 returns `409 Conflict` |
| AC-06 | `PUT /api/v1/students/1` with `{"email": "same@email.com"}` where student 1 already has that email succeeds with `200 OK` |
| AC-07 | Successful update returns full student JSON with `200 OK` |
| AC-08 | All 6 unit test cases pass |
| AC-09 | No existing tests broken |

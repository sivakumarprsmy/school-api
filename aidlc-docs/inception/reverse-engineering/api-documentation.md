# API Documentation

## REST APIs

### Health Check
- **Method**: GET
- **Path**: `/healthz`
- **Purpose**: Liveness probe for the service
- **Request**: No body required
- **Response**:
  ```json
  { "status": "ok" }
  ```
  HTTP 200

---

### Create Student
- **Method**: POST
- **Path**: `/api/v1/students`
- **Purpose**: Enrol a new student in the system
- **Request Body**:
  ```json
  {
    "name":  "Alice Smith",
    "email": "alice@school.edu",
    "age":   15,
    "grade": "10A"
  }
  ```
  | Field | Type   | Required | Validation         |
  |-------|--------|----------|--------------------|
  | name  | string | Yes      | Non-empty          |
  | email | string | Yes      | Valid email format |
  | age   | int    | Yes      | 1–120              |
  | grade | string | Yes      | Non-empty          |

- **Response**:
  - **201 Created** — Student object with assigned ID
    ```json
    {
      "id": 1,
      "name": "Alice Smith",
      "email": "alice@school.edu",
      "age": 15,
      "grade": "10A",
      "created_at": "2026-04-12T00:00:00Z",
      "updated_at": "2026-04-12T00:00:00Z"
    }
    ```
  - **400 Bad Request** — Invalid or missing fields
  - **409 Conflict** — Email already exists
  - **500 Internal Server Error** — Database failure

---

### Get All Students
- **Method**: GET
- **Path**: `/api/v1/students`
- **Purpose**: Retrieve the complete roster of enrolled students
- **Request**: No body required
- **Response**:
  - **200 OK** — Array of student objects (empty array if no students)
    ```json
    [
      { "id": 1, "name": "Alice Smith", "email": "alice@school.edu", "age": 15, "grade": "10A", "created_at": "...", "updated_at": "..." }
    ]
    ```
  - **500 Internal Server Error** — Database failure

---

### Get Student by ID
- **Method**: GET
- **Path**: `/api/v1/students/:id`
- **Purpose**: Retrieve a single student's details by system ID
- **Path Parameter**: `id` — Positive integer (uint64)
- **Response**:
  - **200 OK** — Student object
  - **400 Bad Request** — Non-numeric or invalid ID
  - **404 Not Found** — Student does not exist
  - **500 Internal Server Error** — Database failure

---

### Delete Student
- **Method**: DELETE
- **Path**: `/api/v1/students/:id`
- **Purpose**: Remove a student record from the system
- **Path Parameter**: `id` — Positive integer (uint64)
- **Response**:
  - **204 No Content** — Student successfully deleted
  - **400 Bad Request** — Non-numeric or invalid ID
  - **404 Not Found** — Student does not exist
  - **500 Internal Server Error** — Database failure

---

### Update Student *(MISSING — HLD Requirement)*
- **Method**: PUT
- **Path**: `/api/v1/students/:id`
- **Purpose**: Update an existing student's information
- **Status**: **NOT YET IMPLEMENTED** — this is the feature described in the HLD
- **Expected Design**:
  - Path param: `id` (student ID)
  - Body: Fields to be updated (name, email, age, grade)

---

## Data Models

### Student
- **Fields**:
  | Field      | Type      | DB Constraint         | JSON Key     |
  |------------|-----------|-----------------------|--------------|
  | ID         | uint      | PK, auto-increment    | `id`         |
  | Name       | string    | NOT NULL              | `name`       |
  | Email      | string    | UNIQUE INDEX, NOT NULL| `email`      |
  | Age        | int       | NOT NULL              | `age`        |
  | Grade      | string    | NOT NULL              | `grade`      |
  | CreatedAt  | time.Time | managed by GORM       | `created_at` |
  | UpdatedAt  | time.Time | managed by GORM       | `updated_at` |

- **Relationships**: None (standalone entity)
- **Validation**: Enforced at HTTP layer via `CreateStudentRequest` binding tags

### CreateStudentRequest
- **Purpose**: Input validation DTO for student creation
- **Fields**: name (required), email (required, email), age (required, 1-120), grade (required)
- **Note**: An equivalent `UpdateStudentRequest` DTO does not yet exist — needed for the HLD feature

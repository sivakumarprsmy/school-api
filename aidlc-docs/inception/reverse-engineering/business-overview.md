# Business Overview

## Business Context Diagram

```
+--------------------------------------------------+
|                  School API                      |
|                                                  |
|  Manages student records for a school system.   |
|  Provides CRUD operations over HTTP REST API.   |
+--------------------------------------------------+
         |
         | REST HTTP
         v
+------------------+       +-------------------+
|   API Clients    |       |   PostgreSQL DB   |
| (web apps, CLI,  | ----> |  (student data)   |
|   mobile, etc.)  |       +-------------------+
+------------------+
```

## Business Description

- **Business Description**: The School API is a RESTful backend service that manages student records for a school system. It enables creation, retrieval, and deletion of student profiles, each identified uniquely by their email address and a system-assigned ID. The system stores core academic and personal details per student.

- **Business Transactions**:
  1. **Enrol Student** - Register a new student in the system with name, email, age, and grade information
  2. **View Student** - Retrieve a single student's details by their unique ID
  3. **List All Students** - Retrieve a complete roster of all enrolled students
  4. **Remove Student** - Remove a student record from the system
  5. **Update Student** *(missing — the HLD requirement)* - Modify an existing student's information

- **Business Dictionary**:
  | Term | Meaning |
  |------|---------|
  | Student | A person enrolled in the school, identified by a unique ID and email |
  | Grade | The academic year/class the student belongs to (e.g., "10A", "Grade 5") |
  | Enrol | The act of registering a new student in the system |
  | Roster | The complete list of all enrolled students |

## Component Level Business Descriptions

### school-api (Main Application)
- **Purpose**: Serves as the HTTP API gateway for all student management operations
- **Responsibilities**: Accept REST requests, validate inputs, persist data to PostgreSQL, return structured JSON responses

### Student Handler
- **Purpose**: Implements the business logic for each student operation
- **Responsibilities**: Request parsing, validation, database interaction, response formatting, error handling

### Student Model
- **Purpose**: Defines the canonical data structure for a student record
- **Responsibilities**: Database schema definition, JSON serialization, input validation rules

### Configuration
- **Purpose**: Manages runtime configuration from environment variables
- **Responsibilities**: Load and validate all settings (database connection, port, environment, log level)

### Database Layer
- **Purpose**: Manages PostgreSQL connectivity and schema migrations
- **Responsibilities**: Connection establishment, auto-migration of schema, query execution via GORM

### Logger
- **Purpose**: Provides structured, levelled logging across the application
- **Responsibilities**: Console output with timestamps, configurable log levels, request logging middleware

# User Stories — Update Student Endpoint

**Feature**: Update Student Information
**Persona**: Alex — School Administrator
**Breakdown Approach**: Feature-based — two stories per feature (happy path + error/edge cases)

---

## US-01: Update Student Information

**As a** School Administrator,
**I want to** update an existing student's information by their ID,
**So that** student records remain accurate when details change.

### INVEST Check
- **Independent**: Does not depend on any other story
- **Negotiable**: Which fields are included can be adjusted
- **Valuable**: Directly addresses the gap identified in the HLD — no way to update students
- **Estimable**: Clear scope — one PUT endpoint with partial update logic
- **Small**: Single endpoint addition; fits in one sprint
- **Testable**: Concrete acceptance criteria below

### Acceptance Criteria

**AC-01**: I can update a student's name, email, age, or grade — or any combination of these fields — by providing the student's ID and only the fields I wish to change.

**AC-02**: Fields I do not include in my update request are left unchanged on the student record.

**AC-03**: After a successful update, I receive the complete, updated student record in the response so I can confirm my changes were applied.

**AC-04**: I can update a student's email to the same email address they already have, and the update succeeds without error.

### Persona Mapping
- **Alex (School Administrator)**: Primary user — corrects student data errors through an admin portal

---

## US-02: Receive Clear Feedback When a Student Update Cannot Be Completed

**As a** School Administrator,
**I want to** receive a clear and specific error message when my update request cannot be completed,
**So that** I understand what went wrong and can take the right corrective action.

### INVEST Check
- **Independent**: Can be implemented alongside US-01 (same handler, different code paths)
- **Negotiable**: Specific error message wording can be refined
- **Valuable**: Prevents confusion and support requests caused by silent failures or cryptic errors
- **Estimable**: Error paths are well-defined; bounded scope
- **Small**: No new infrastructure — error handling within the same endpoint
- **Testable**: Each error scenario has a distinct, verifiable outcome

### Acceptance Criteria

**AC-05**: If I provide a student ID that does not exist in the system, I am notified that the student was not found.

**AC-06**: If I provide a student ID that is not a valid number, I am notified that my request is invalid.

**AC-07**: If I attempt to update a student's email to one that is already used by a different student, I am notified that the email is already taken.

**AC-08**: If I provide a field value that does not meet the system's data rules (for example, a malformed email address or an age outside the accepted range), I am told my request contains invalid data.

**AC-09**: If a technical error occurs on the server while processing my update, I receive a generic error message that does not expose any internal system details.

### Persona Mapping
- **Alex (School Administrator)**: Primary user — needs actionable feedback when a record cannot be updated

---

## Story Summary

| ID | Title | Priority | Acceptance Criteria |
|----|-------|----------|---------------------|
| US-01 | Update Student Information | High | AC-01, AC-02, AC-03, AC-04 |
| US-02 | Receive Clear Feedback When Update Fails | High | AC-05, AC-06, AC-07, AC-08, AC-09 |

**Total Stories**: 2
**Total Acceptance Criteria**: 9
**Personas Covered**: School Administrator (Alex)

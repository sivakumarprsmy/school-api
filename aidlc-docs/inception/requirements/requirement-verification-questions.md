# Requirements Clarification Questions

**Feature**: Update Student endpoint — `PUT /api/v1/students/:id`

Please answer each question by filling in the letter choice after the `[Answer]:` tag.
If none of the options match, choose the last option (Other/X) and describe your preference.
Let me know when you're done.

---

## Question 1
Which student fields should be updatable via the PUT endpoint?

A) All editable fields: name, email, age, grade (full update — all fields required in request)
B) All editable fields: name, email, age, grade (partial update — only provided fields updated)
C) Only non-sensitive fields: name, age, grade (email excluded from updates)
X) Other (please describe after [Answer]: tag below)

[Answer]: B

---

## Question 2
What HTTP semantics should the update follow?

A) Strict PUT — all fields required in request body; missing fields cause a 400 error
B) PATCH-like PUT — only fields present in request body are updated; omitted fields unchanged
X) Other (please describe after [Answer]: tag below)

[Answer]: B

---

## Question 3
What should the endpoint return on a successful update?

A) 200 OK with the full updated student object (JSON)
B) 204 No Content (empty body)
X) Other (please describe after [Answer]: tag below)

[Answer]: A

---

## Question 4
What should happen if the student ID does not exist?

A) 404 Not Found with an error message (consistent with existing GET/DELETE behaviour)
B) 200 OK with empty body (upsert — create if not found)
X) Other (please describe after [Answer]: tag below)

[Answer]: A

---

## Question 5
If email is updatable (Question 1), should email uniqueness be re-validated?

A) Yes — return 409 Conflict if the new email is already used by another student
B) No — email updates are not supported regardless of Question 1 answer
X) Other (please describe after [Answer]: tag below)

[Answer]: A

---

## Question 6
Should a unit test be added for the new UpdateStudent handler (consistent with the existing test suite in `internal/handlers/students_test.go`)?

A) Yes — add tests covering success, not-found, invalid ID, validation error, and DB error paths
B) Yes — add a basic success test only
C) No — skip tests for this endpoint
X) Other (please describe after [Answer]: tag below)

[Answer]: A

---

## Question 7 — Security Extension
Should security extension rules be enforced for this project?

A) Yes — enforce all SECURITY rules as blocking constraints (recommended for production-grade applications)
B) No — skip all SECURITY rules (suitable for PoCs, prototypes, and experimental projects)
X) Other (please describe after [Answer]: tag below)

[Answer]: A

---

## Question 8 — Property-Based Testing Extension
Should property-based testing (PBT) rules be enforced for this project?

A) Yes — enforce all PBT rules as blocking constraints (recommended for projects with business logic, data transformations, serialization, or stateful components)
B) Partial — enforce PBT rules only for pure functions and serialization round-trips
C) No — skip all PBT rules (suitable for simple CRUD applications or thin integration layers)
X) Other (please describe after [Answer]: tag below)

[Answer]: A

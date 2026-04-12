# User Personas — School API

## Persona 1: School Administrator

| Attribute | Detail |
|-----------|--------|
| **Name** | Alex (School Administrator) |
| **Role** | Administrative staff responsible for maintaining student records |
| **Environment** | Works via an admin portal or internal management tool that calls the School API |
| **Technical Literacy** | Non-technical end user; interacts through a UI backed by the API |
| **Goals** | Keep student records accurate and up-to-date; correct data entry mistakes quickly |
| **Pain Points** | No way to fix incorrect student information once a record is created; must delete and re-create to correct an error |
| **Motivation** | Efficiency and accuracy — a simple update saves time and prevents data quality issues downstream |

### Typical Scenarios
- A student changes their email address and Alex needs to update the record
- A student was assigned the wrong grade at enrolment and needs correction
- A student's legal name changes and the record must be updated

### What Alex Needs
- Ability to update any student field (name, email, age, grade) individually or together
- Confidence that partial updates won't accidentally wipe out other fields
- Clear, actionable feedback when something goes wrong (duplicate email, student not found)
- The updated record returned immediately to confirm the change was applied

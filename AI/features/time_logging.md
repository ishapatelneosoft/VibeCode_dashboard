You are extending an existing backend system.

---

## FEATURE

Time Logging

---

## OBJECTIVE

Allow users to log time spent on tasks with immutable worklog entries.

---

## REQUIREMENTS

1. Log Work

* Users can log time in decimal hours (e.g., 2.5)
* Description is optional

2. Ownership

* Worklog must be associated with logged-in user

3. Immutability

* Worklogs cannot be edited or deleted

4. Multiple Entries

* Multiple worklogs allowed per task

---

## DATA MODEL

Worklog:

* id
* task_id
* user_id
* time_spent (decimal)
* description (optional)
* created_at

---

## API CONTRACT

POST /tasks/{id}/worklogs

GET /tasks/{id}/worklogs

---

## IMPLEMENTATION RULES

* Validate time_spent > 0
* Allow decimal values
* Do not allow update/delete APIs
* Do not overwrite existing records
* Use authenticated user as owner

---

## QUERY REQUIREMENTS

* Fetch worklogs ordered by created_at DESC

---

## PERFORMANCE

* Indexed queries
* Avoid unnecessary joins
* Keep insert lightweight

---

## TESTING (MANDATORY)

* Log time → success
* Decimal input → accepted
* Multiple logs → stored correctly
* Correct user association
* No edit/delete allowed
* Fetch logs → correct order

---

## OUTPUT

1. Change Summary
2. Model
3. DTO
4. Service Logic
5. Repository Logic
6. Controller Changes
7. Test Cases

---

## GOAL

Provide reliable, immutable time tracking per task.

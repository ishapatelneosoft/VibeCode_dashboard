You are extending an existing backend system.

---

## FEATURE

Shared Board Visibility

---

## OBJECTIVE

Ensure all authenticated users access the same shared board and tasks with consistent structure and visibility.

---

## REQUIREMENTS

1. Shared Board Access

* All authenticated users must access the same board
* Tasks created by one user must be visible to all users

2. Board Structure

* Board must contain exactly 8 columns
* Columns must be returned in a fixed, predefined order

3. Task Representation
   Each task must include:

* Title
* Assignee (nullable)
* Due Date (nullable)
* Created By (user reference)

---

## IMPLEMENTATION RULES

* Do NOT introduce per-user scoping
* Board is global (single source of truth)
* Use indexed queries for fetching tasks
* Ensure ordering at DB level (not in-memory)

---

## DATA MODEL EXTENSION

Task:

* id
* title
* description (optional)
* column (enum / reference)
* assignee_id (nullable)
* due_date (nullable)
* created_by (required)
* position (for ordering within column)
* created_at

Column:

* id
* name
* order (integer, fixed)

---

## API REQUIREMENTS

GET /board

* Returns all columns in order
* Each column contains ordered tasks

Response:
{
"success": true,
"data": {
"columns": [
{
"name": "",
"order": 0,
"tasks": []
}
]
}
}

---

## PERFORMANCE

* Avoid N+1 queries (use joins or batching)
* Index: column, created_at, position
* Fetch board in single optimized query

---

## TESTING (MANDATORY)

* User A creates task → User B can view it
* Board returns exactly 8 columns
* Columns are in correct order
* Task contains all required fields
* Null handling for optional fields

---

## OUTPUT

1. Change Summary
2. Updated Models
3. Repository Queries
4. Service Logic
5. Controller अपडेट
6. Test Cases (unit + integration)

---

## GOAL

Provide a consistent shared board view across all users with efficient data retrieval.

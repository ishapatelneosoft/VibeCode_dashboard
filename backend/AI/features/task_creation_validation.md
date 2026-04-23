You are extending an existing backend system.

---

## FEATURE

Task Creation & Validation

---

## OBJECTIVE

Implement task creation with strict validation, default behaviors, and correct ordering.

---

## REQUIREMENTS

1. Minimal Input

* Only Title is required

2. Title Validation

* Max length: 255 characters
* Reject if exceeded

3. Default Values

* Column defaults to "Backlog"
* Assignee defaults to null
* created_by must be current authenticated user

4. Ordering

* New tasks must appear at the bottom of the column
* Maintain position index per column

---

## IMPLEMENTATION RULES

* Validate input at DTO level
* Do not rely on frontend validation
* Position calculation must be atomic (avoid race conditions)
* Use DB transaction for task creation

---

## DATA LOGIC

On task creation:

1. Validate title length
2. Fetch max(position) for Backlog column
3. New position = max + 1
4. Insert task with:

   * column = Backlog
   * assignee = null
   * created_by = current user

---

## API REQUIREMENTS

POST /tasks

Request:
{
"title": "string"
}

Response:
{
"success": true,
"data": {
"task_id": "",
"column": "Backlog",
"position": 0
},
"message": "Task created"
}

---

## EDGE CASES

* Empty title → reject
* Title > 255 → reject
* Concurrent inserts → no duplicate position
* Missing auth → reject

---

## PERFORMANCE

* Index on (column, position)
* Avoid full table scan for position
* Use transaction for consistency

---

## TESTING (MANDATORY)

* Create task with only title → success
* Title > 255 → error
* Default column = Backlog
* Assignee = null
* created_by = logged-in user
* Tasks ordered correctly (new at bottom)
* Concurrent task creation maintains order

---

## OUTPUT

1. Change Summary
2. DTO Validation
3. Service Logic
4. Repository Logic
5. Controller Changes
6. Test Cases (unit + integration)

---

## GOAL

Ensure reliable, validated, and consistently ordered task creation.

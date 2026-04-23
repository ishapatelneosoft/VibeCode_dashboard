
## FEATURE

Task Movement & Reordering

---

## OBJECTIVE

Enable tasks to:

* Move across columns
* Reorder within the same column
* Persist position and status reliably

---

## REQUIREMENTS

1. Cross Column Movement

* Task can move from one column to another
* Column change must persist after refresh

2. Reordering Within Column

* Tasks can be reordered vertically
* Order must persist after refresh

3. Persistence

* All changes must be stored in DB
* No in-memory ordering

---

## POSITIONING STRATEGY (MANDATORY)

Use fractional indexing for ordering.

Rules:

* position is a float
* Insert between:
  new_position = (before.position + after.position) / 2
* Insert at top:
  new_position = after.position / 2
* Insert at bottom:
  new_position = before.position + 1

DO NOT:

* Shift entire column
* Recalculate all positions

---

## API CONTRACT

PATCH /tasks/{id}/move

Request:
{
"destination_column_id": "string",
"before_task_id": "string|null",
"after_task_id": "string|null"
}

Response:
{
"success": true,
"data": {
"task_id": "",
"column_id": "",
"position": 0
},
"message": "Task moved successfully"
}

---

## SERVICE LOGIC

1. Validate task exists
2. Validate destination column
3. Fetch adjacent tasks (before/after)
4. Compute new position using fractional indexing
5. Update:

   * column_id
   * position
6. Use DB transaction

---

## REPOSITORY REQUIREMENTS

* Fetch task by ID
* Fetch adjacent tasks
* Update task position and column
* Use indexed queries

Indexes:
(column_id, position)

---

## EDGE CASES

* Move to empty column
* Move to top
* Move to bottom
* Move within same column
* Invalid task IDs
* Missing auth

---

## PERFORMANCE

* O(1) reordering
* No bulk updates
* Avoid N+1 queries

---

## TESTING (MANDATORY)

* Move task across columns → persists after refresh
* Reorder within column → order correct after refresh
* Move to top/bottom
* Move to empty column
* Concurrent moves maintain consistency

---

## OUTPUT

1. Change Summary
2. DTO
3. Service Logic
4. Repository Logic
5. Controller Changes
6. Test Cases

---

## GOAL

Efficient, scalable task movement with consistent ordering.

You are extending an existing backend system.

---

## FEATURE

Assignment Management & History

---

## OBJECTIVE

Enable assigning tasks to users and maintain a complete history of assignment changes.

---

## REQUIREMENTS

1. Assign User

* Task can be assigned to any registered user
* Assignee can be changed

2. Unassign

* Assignee can be set to null

3. History Tracking

* Every change must create a history record
* Store old and new assignee
* Store who made the change
* Store timestamp

4. History Retrieval

* Return assignment history for a task
* Most recent first

---

## DATA MODEL

AssignmentHistory:

* id
* task_id
* old_assignee_id (nullable)
* new_assignee_id (nullable)
* changed_by
* created_at

---

## API CONTRACT

PATCH /tasks/{id}/assignee

GET /tasks/{id}/history

GET /users

---

## IMPLEMENTATION RULES

* Use DB transaction for update + history insert
* Do not overwrite history
* Avoid duplicate history if no change
* Validate assignee exists
* Maintain null support

---

## PERFORMANCE

* Indexed queries for history
* Avoid N+1 user fetch
* Batch user resolution if needed

---

## TESTING (MANDATORY)

* Assign user → task updated
* Change assignee → history created
* Unassign → works correctly
* Multiple changes → correct order
* Same assignee → no new history
* Fetch users → list available

---

## OUTPUT

1. Change Summary
2. Model Changes
3. Service Logic
4. Repository Logic
5. Controller Changes
6. Test Cases

---

## GOAL

Provide reliable assignment tracking with complete audit history.

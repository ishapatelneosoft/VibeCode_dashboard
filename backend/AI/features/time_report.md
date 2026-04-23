You are extending an existing backend system.

---

## FEATURE

Time Report View

---

## OBJECTIVE

Provide aggregated time tracking across tasks and the entire project.

---

## REQUIREMENTS

1. Task-Level Report

* Each task must include:

  * Title
  * Status (column)
  * Assignee
  * Total hours (sum of worklogs)

2. Project-Level Total

* Return grand total across all tasks

3. Accessibility

* Available to any authenticated user

---

## API CONTRACT

GET /reports/time

---

## DATA AGGREGATION

* Use SUM(time_spent)
* Group by task_id
* Include tasks with zero worklogs

---

## IMPLEMENTATION RULES

* Use optimized SQL query (JOIN + GROUP BY)
* Avoid N+1 queries
* Do not compute totals in loops
* Do not store totals in task table

---

## PERFORMANCE

* Use indexed queries
* Avoid full scans where possible
* Keep query minimal

---

## EDGE CASES

* Tasks with no worklogs → total = 0
* No tasks → empty response
* No worklogs → grand total = 0

---

## TESTING (MANDATORY)

* Task total matches sum of worklogs
* Grand total matches sum of all worklogs
* Tasks without logs return 0
* Multiple tasks aggregated correctly
* Unauthorized access rejected

---

## OUTPUT

1. Change Summary
2. Repository Query
3. Service Logic
4. Controller
5. Response DTO
6. Test Cases

---

## GOAL

Efficient and accurate time aggregation for reporting.

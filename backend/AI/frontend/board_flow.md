You are extending an existing Next.js frontend application.

---

## OBJECTIVE

Implement full board functionality using existing backend APIs.

Features included:

1. Task Modal (view only)
2. Task Creation
3. Drag & Drop (movement)
4. Assignment Management
5. Worklogs (time logging)
6. Reports Page

---

## PRIMARY RULE

Follow existing frontend architecture strictly:
API → Service → Hook → Component

Do NOT bypass layers.

---

## MANDATORY ANALYSIS

Before implementation:

* Reuse existing API client
* Reuse existing hooks pattern (login/board)
* Reuse response structure
* Do NOT create new patterns

---

# 1. TASK MODAL (VIEW ONLY)

---

## OBJECTIVE

Display task details in a modal.

---

## DATA SOURCES

GET /tasks/{id}/history
GET /tasks/{id}/worklogs
GET /users

---

## UI STRUCTURE

Modal Sections:

1. Task Info

   * Title
   * Assignee
   * Created by
   * Due date

2. Assignment History

   * List (latest first)
   * old → new
   * changed_by
   * timestamp

3. Worklogs

   * time_spent
   * description
   * user
   * timestamp

---

## RULES

* Read-only (no mutation)
* Fetch data using React Query
* Lazy load modal data (on open)
* Show loading state per section

---

# 2. TASK CREATION

---

## OBJECTIVE

Create new tasks from UI.

---

## API

POST /tasks

---

## RULES

* Only Title required
* On success:

  * Close modal/form
  * Refetch board
* Default column handled by backend

---

# 3. DRAG & DROP

---

## OBJECTIVE

Enable moving tasks across columns and reordering.

---

## LIBRARY

Use dnd-kit

---

## FLOW

1. Drag start
2. Optimistic UI update
3. Call PATCH /tasks/{id}/move
4. On error → rollback

---

## RULES

* Do NOT calculate position in frontend
* Use backend contract (before_task_id / after_task_id)

---

# 4. ASSIGNMENT MANAGEMENT

---

## OBJECTIVE

Update task assignee.

---

## API

PATCH /tasks/{id}/assignee

---

## UI

* Dropdown with users
* Include "Unassigned"

---

## RULES

* Refetch:

  * task data
  * assignment history

---

# 5. WORKLOGS

---

## OBJECTIVE

Log time for tasks.

---

## API

POST /tasks/{id}/worklogs
GET /tasks/{id}/worklogs

---

## UI

* Input: time_spent (decimal)
* Input: description
* List existing logs

---

## RULES

* Immutable logs
* Refetch after insert

---

# 6. REPORTS PAGE

---

## OBJECTIVE

Display aggregated time report.

---

## API

GET /reports/time

---

## UI

* Table view:

  * Task
  * Status
  * Assignee
  * Total hours
* Footer:

  * Grand total

---

## RULES

* No client-side aggregation
* Render backend data as-is

---

# SHARED COMPONENTS

---

## CREATE/EXTEND

/components
/modal
TaskModal.tsx
/task
TaskCard.tsx
TaskForm.tsx
/board
Board.tsx
Column.tsx
/report
TimeReport.tsx

---

# STATE MANAGEMENT

---

## RULES

* Server state → React Query
* UI state (modal, drag state) → local/Zustand
* Do NOT duplicate server data

---

# PERFORMANCE

---

## RULES

* Use optimistic updates for drag-drop
* Avoid full board re-renders
* Memoize TaskCard

---

# ERROR HANDLING

---

## RULES

* Show API errors
* Rollback optimistic updates on failure

---

# OUTPUT

---

Return:

1. File structure (new + updated)
2. Services (if extended)
3. Hooks (per feature)
4. Components
5. Integration flow

---

# STRICT RULES

---

* Do NOT fetch inside components directly
* Do NOT hardcode API URLs
* Do NOT compute backend logic in frontend
* Do NOT mix concerns (UI vs data)
* Do NOT break existing patterns

---

# GOAL

---

Build a scalable, maintainable, production-grade frontend aligned with backend architecture.

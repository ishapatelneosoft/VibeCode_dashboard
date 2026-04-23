You are extending an existing Next.js frontend application.

---

## FEATURE

Board View (Read-Only)

---

## OBJECTIVE

Display the shared board with columns and tasks using existing backend API.

---

## API CONTRACT (STRICT)

GET /board

Response:
{
"success": true,
"data": {
"columns": [
{
"name": "string",
"order": number,
"tasks": [
{
"id": "string",
"title": "string",
"assignee": "string|null",
"due_date": "string|null",
"created_by": "string"
}
]
}
]
}
}

---

## IMPLEMENTATION RULES

* Use existing API client
* Use React Query for fetching
* Do NOT transform backend structure unnecessarily
* Maintain backend ordering (do NOT sort in frontend)

---

## FOLDER STRUCTURE

/app/board/page.tsx
/services/board.service.ts
/hooks/useBoard.ts
/types/board.ts
/components/board/
Board.tsx
Column.tsx
TaskCard.tsx

---

## SERVICE LAYER

board.service.ts:

* getBoard()

---

## HOOK

useBoard:

* React Query (useQuery)
* Handles:

  * loading
  * error
  * data

---

## UI COMPONENTS

Board:

* Renders all columns

Column:

* Displays column name
* Renders task list

TaskCard:

* Displays:

  * Title
  * Assignee (if present)
  * Due date (if present)
  * Created by

---

## UI REQUIREMENTS

* Columns displayed in order (from API)
* Tasks displayed as provided (no reordering)
* Handle empty columns
* Show loading state
* Show error state

---

## DATA HANDLING

* Do NOT flatten structure
* Use columns → tasks as-is
* Null-safe rendering for optional fields

---

## PERFORMANCE

* Use key properly for lists
* Avoid re-rendering entire board
* Memoize TaskCard if needed

---

## ERROR HANDLING

* Display API error message
* Fallback UI for failure

---

## STRICT RULES

* Do NOT use local state for board data
* Do NOT fetch inside components
* Do NOT sort/filter client-side
* Do NOT assume data shape beyond API contract

---

## OUTPUT

1. Page (board/page.tsx)
2. Service
3. Hook
4. Components (Board, Column, TaskCard)
5. Types

---

## GOAL

Render a consistent, backend-driven board view with clean data flow.

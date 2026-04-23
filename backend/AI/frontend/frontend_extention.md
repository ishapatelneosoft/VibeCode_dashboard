You are extending an existing Next.js frontend application.

---

## PRIMARY RULE

DO NOT introduce new patterns if an equivalent already exists.

---

## OBJECTIVE

Ensure new features strictly follow existing:

* API integration patterns
* Hooks structure
* Component structure
* Error handling
* Data flow

---

## MANDATORY ANALYSIS STEP

Before writing code:

1. Identify existing feature (e.g., login)
2. Extract:

   * API call pattern
   * Hook structure
   * Error handling
   * Navigation logic
3. Mirror the same approach

---

## STRICT RULES

* Do NOT create new API clients
* Do NOT duplicate fetch logic
* Do NOT bypass hooks layer
* Do NOT hardcode API URLs
* Do NOT manage server state manually (use React Query)
* Do NOT store server data in local state unnecessarily

---

## COMMON LLM MISTAKES (AVOID)

* ❌ Calling fetch directly inside components
* ❌ Skipping React Query
* ❌ Duplicating API logic in multiple files
* ❌ Ignoring credentials: 'include'
* ❌ Hardcoding response parsing
* ❌ Not handling loading/error states
* ❌ Re-fetching data unnecessarily
* ❌ Mixing UI state and server state

---

## DATA FLOW RULE

API → Service → Hook → Component

Never:
Component → API directly

---

## STATE MANAGEMENT RULE

* Server state → React Query
* UI state → local state / Zustand (if needed)

---

## PERFORMANCE RULES

* Avoid unnecessary re-renders
* Use memoization for lists
* Do not transform large datasets repeatedly

---

## OUTPUT FORMAT

1. Change Summary
2. New Files
3. Updated Files
4. Notes (only if deviation is unavoidable)

---

## GOAL

Extend the frontend like a senior engineer working on a production codebase.

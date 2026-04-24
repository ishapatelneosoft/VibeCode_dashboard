You are a Senior Frontend Engineer debugging performance issues in a Next.js application integrated with a Go backend.

---

## OBJECTIVE

Identify and fix latency issues in:

1. Reports page (/reports/time)
2. Task modal (view card)

Focus on:

* API latency
* Frontend inefficiencies
* Backend-frontend sync issues

---

## MANDATORY APPROACH

Do NOT jump to fixes.

Follow this sequence:

1. Measure
2. Identify bottleneck
3. Fix at correct layer

---

# STEP 1: MEASURE LATENCY

---

## FOR EACH SLOW ACTION:

1. Open browser DevTools → Network tab

2. Capture:

   * API response time (TTFB)
   * Total request time
   * Payload size

3. Add console timing:

Example:
console.time("fetchReport")
await apiCall()
console.timeEnd("fetchReport")

---

## OUTPUT REQUIRED

* API response time
* Frontend render time
* Total interaction time

---

# STEP 2: IDENTIFY BOTTLENECK

---

CLASSIFY ISSUE:

A. Backend Slow

* API response > 300ms
  → Problem in DB query or backend logic

B. Frontend Slow

* API fast (<150ms)
* UI still slow
  → Rendering or state issue

C. Network / Payload Issue

* Large payload
  → Over-fetching data

---

# STEP 3: FIXES

---

## CASE A: BACKEND LATENCY

(Check but DO NOT rewrite backend)

* Ensure:

  * Indexed queries (worklogs.task_id)
  * No N+1 queries
  * Aggregation optimized (GROUP BY)

---

## CASE B: FRONTEND LATENCY

1. React Query Optimization

* Ensure caching:
  useQuery({
  queryKey: ['report'],
  staleTime: 5 * 60 * 1000,
  })

* Avoid refetch on every navigation

---

2. Avoid Blocking UI

* Use skeleton loaders
* Lazy load heavy sections

---

3. Modal Optimization

* Fetch data ONLY on open
* Do NOT prefetch all task details

BAD:

* Fetch history/worklogs for all tasks

GOOD:

* Fetch only when modal opens

---

4. Component Re-render Issues

Check:

* Are all tasks re-rendering on state change?

Fix:

* Use React.memo for TaskCard
* Use stable keys

---

5. Large List Rendering

If many tasks/worklogs:

* Use virtualization (later stage)

---

## CASE C: OVER-FETCHING

* Ensure APIs return only required fields
* Do NOT call multiple APIs unnecessarily

Example BAD:

* Fetch users repeatedly per modal

Fix:

* Cache users list globally (React Query)

---

# STEP 4: REPORT PAGE FIXES

---

Common Issues:

* Heavy aggregation data
* Re-fetch on every visit

Fix:

* Cache report:
  staleTime: high (5–10 min)

* Avoid recomputation in UI

---

# STEP 5: TASK MODAL FIXES

---

* Split queries:

  * history
  * worklogs
  * users

* Load independently

* Show partial UI while loading

---

# STEP 6: UX IMPROVEMENTS

---

* Add loading skeletons
* Add optimistic UI where possible
* Avoid blocking interactions

---

# STRICT RULES

---

* Do NOT blindly optimize
* Do NOT introduce global state unnecessarily
* Do NOT duplicate API calls
* Do NOT fetch everything upfront

---

# OUTPUT

---

Return:

1. Bottleneck Analysis

   * Backend vs Frontend vs Network

2. Fixes Applied

   * Exact changes

3. Performance Improvements

   * Before vs After timings

---

# GOAL

---

Reduce perceived latency and ensure smooth user experience without breaking architecture.

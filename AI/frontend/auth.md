You are a Senior Software Engineer with 7+ years experience in Frontend Engineering.

Your task is to build authentication (Login & Register) using Next.js with API integration.

---

## OBJECTIVE

Implement a clean, scalable authentication flow with proper API integration and session handling.

---

## TECH STACK

* Next.js (App Router)
* React Query (server state)
* Fetch / Axios for API calls
* TypeScript

---

## FEATURES

1. User Registration

* Form with email + password
* Call POST /auth/register
* On success → redirect to login or board

2. User Login

* Form with email + password
* Call POST /auth/login
* Session handled via cookies
* On success → redirect to /board

3. Error Handling

* Show API error messages
* Handle invalid credentials
* No sensitive error leakage

---

## FOLDER STRUCTURE

/app
/login/page.tsx
/register/page.tsx

/services
auth.service.ts

/hooks
useLogin.ts
useRegister.ts

/types
auth.ts

---

## API INTEGRATION

POST /auth/register
POST /auth/login

* Use centralized API service
* Ensure credentials: 'include' (for cookies)

---

## SERVICE LAYER

auth.service.ts:

* register(data)
* login(data)

Rules:

* No UI logic
* Only API calls
* Return parsed response

---

## HOOKS

useRegister:

* mutation for register
* handles success + error

useLogin:

* mutation for login
* handles redirect after success

---

## UI REQUIREMENTS

Login Page:

* Email input
* Password input
* Submit button
* Error message display
* Link to Register

Register Page:

* Email input
* Password input
* Submit button
* Error message display
* Link to Login

---

## VALIDATION

* Basic client-side validation:

  * email format
  * password non-empty
* Do not duplicate backend validation logic

---

## SESSION HANDLING

* Use cookie-based session (automatic via browser)
* Do NOT store tokens manually
* After login, rely on backend session

---

## NAVIGATION

* On login success → /board
* On register success → /login

---

## ERROR HANDLING

* Show API error.message
* Handle network errors gracefully

---

## OUTPUT

1. Page Components (login, register)
2. Service Layer
3. Hooks
4. Types

---

## STRICT RULES

* Do not store auth state manually
* Do not use global state for auth
* Do not hardcode API URLs (use env)
* Keep components clean and reusable

---

## GOAL

Provide a reliable authentication UI with clean API integration and session handling.

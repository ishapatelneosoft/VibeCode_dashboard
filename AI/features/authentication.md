# Authentication Module

## Overview

Secure user authentication using session-based architecture with PostgreSQL storage. Implements registration, login, session management, and logout following clean architecture patterns (Controller → Service → Repository). Built with Go, Gin, and PostgreSQL.

---

## Core Components

### 1. Authentication Flow

HTTP Request → AuthController → AuthService → Repository → PostgreSQL  
Response → Standardized JSON (success/error format)

### 2. Session Management

* Session-based authentication with secure HttpOnly cookies
* Sessions stored in PostgreSQL `sessions` table
* Session TTL configurable (default 7 days)
* Automatic session validation via middleware

### 3. Password Security

* Passwords hashed using bcrypt with cost 12
* Email validation via regex pattern
* Password strength enforcement (min 8 characters)

---

## Feature Behavior

### User Registration (`POST /api/auth/register`)

* Accepts `{email, password}` JSON
* Validates email format and password strength
* Checks for duplicate email via unique constraint
* Hashes password with bcrypt before storage
* Creates user record in `users` table
* Returns `201 Created` with user ID

### User Login (`POST /api/auth/login`)

* Validates `{email, password}` credentials
* Fetches user by email from repository
* Compares hashed password via bcrypt.CompareHashAndPassword
* Creates session record with UUID session ID
* Sets HttpOnly cookie `session_id`
* Returns `200 OK` with user data

### Invalid Login Handling

* Returns `401 Unauthorized` with generic error message
* Prevents user enumeration via consistent timing
* Uses `ErrInvalidCredentials` service error

### Session Persistence

* Session stored in PostgreSQL with `user_id`, `session_id`, `expires_at`
* Middleware extracts `session_id` from cookie, validates against store
* Session survives browser restart (persistent cookie)

### Logout (`POST /api/auth/logout`)

* Invalidates session in PostgreSQL store
* Clears client cookie by setting empty value with past expiration
* Returns `200 OK` success response

---

## Middleware

### Authentication Middleware (`auth_middleware.go`)

* Extracts `session_id` from cookie
* Validates session in repository
* Attaches `user_id` to request context via `SetUserID`
* Rejects unauthorized access with `401 Unauthorized`
* Used via `router.Use(middleware.RequireAuth)`

---

## Data Handling

### User Entity (`domain/user.go`)

```go
type User struct {
    ID           string    `db:"id"`
    Email        string    `db:"email"`
    PasswordHash string    `db:"password_hash"`
    CreatedAt    time.Time `db:"created_at"`
    UpdatedAt    time.Time `db:"updated_at"`
}
```

### Session Entity (`domain/session.go`)

```go
type Session struct {
    SessionID string    `db:"session_id"`
    UserID    string    `db:"user_id"`
    ExpiresAt time.Time `db:"expires_at"`
    CreatedAt time.Time `db:"created_at"`
}
```

### Database Schema (migrations/001_create_auth_tables.sql)

* `users` table with `email` unique index
* `sessions` table with `session_id` primary key and `user_id` foreign key
* Automatic timestamp columns

---

## Constraints

* No business logic in controllers (delegated to services)
* Services contain core validation and business rules
* Repositories handle raw SQL/database interaction only
* Consistent response format via `response.go` helpers
* Horizontal scalability supported (stateless services + shared PostgreSQL)

---

## Security Considerations

* HttpOnly, Secure cookies (configurable for production)
* BCrypt hashing with appropriate cost
* Generic error messages to prevent enumeration
* Input validation at service layer
* Session expiration enforced
* SQL injection prevention via parameterized queries

## Testing

* Unit tests for `AuthService` methods
* Integration tests for HTTP endpoints
* Repository tests with test database
* Middleware behavior verification

## Dependencies

* Go 1.21+
* Gin web framework
* PostgreSQL driver (`lib/pq` or `pgx`)
* BCrypt (`golang.org/x/crypto/bcrypt`)
* UUID generation (`github.com/google/uuid`)

## Configuration

* Session TTL via environment variable `SESSION_TTL`
* Cookie settings (Secure, HttpOnly, Path, Domain)
* Database connection pooling parameters

---

## Implementation Notes

Follows the existing project pattern:
1. Controller receives HTTP request, binds DTO
2. Service validates business rules, orchestrates domain logic
3. Repository performs database operations
4. Domain entities represent business objects
5. Middleware handles cross-cutting concerns

All error handling uses centralized response helpers (`BadRequestResponse`, `InternalServerErrorResponse`, etc.) defined in `internal/controller/response.go`.

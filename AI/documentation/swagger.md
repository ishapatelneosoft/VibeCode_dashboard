# Swagger (OpenAPI 3.0) Documentation

Generated OpenAPI specification for backend APIs.

## OpenAPI YAML

```yaml
openapi: 3.0.0
info:
  title: Backend API
  version: 1.0.0
  description: API for authentication, board, and task management

servers:
  - url: /api/v1
    description: API v1

paths:
  /health:
    get:
      summary: Health check
      tags:
        - Health
      responses:
        '200':
          description: OK
          content:
            application/json:
              schema:
                type: object
                properties:
                  status:
                    type: string
                    example: "ok"

  /auth/register:
    post:
      summary: Register a new user
      tags:
        - Auth
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/RegisterRequest'
      responses:
        '201':
          description: User created successfully
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/StandardResponse'
        '400':
          description: Invalid input or user already exists
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'
        '500':
          description: Internal server error
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'

  /auth/login:
    post:
      summary: Authenticate user and create session
      tags:
        - Auth
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/LoginRequest'
      responses:
        '200':
          description: Login successful, session cookie set
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/StandardResponse'
          headers:
            Set-Cookie:
              schema:
                type: string
                example: "session_id=abc123; HttpOnly; Path=/; Max-Age=604800"
        '401':
          description: Invalid credentials
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'
        '500':
          description: Internal server error
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'

  /auth/logout:
    post:
      summary: Logout user (invalidate session)
      tags:
        - Auth
      responses:
        '200':
          description: Logout successful, session cookie cleared
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/StandardResponse'
          headers:
            Set-Cookie:
              schema:
                type: string
                example: "session_id=; HttpOnly; Path=/; Max-Age=0"
        '500':
          description: Internal server error
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'

  /auth/me:
    get:
      summary: Get current authenticated user info
      tags:
        - Auth
      security:
        - sessionAuth: []
      responses:
        '200':
          description: User info retrieved
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/StandardResponse'
        '401':
          description: Not authenticated
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'
        '500':
          description: Internal server error
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'

  /auth/validate:
    get:
      summary: Validate session (no authentication required)
      tags:
        - Auth
      responses:
        '200':
          description: Session is valid
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/StandardResponse'
        '401':
          description: Session invalid or expired
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'

  /board:
    get:
      summary: Get shared board with columns and tasks
      tags:
        - Board
      security:
        - sessionAuth: []
      responses:
        '200':
          description: Board retrieved successfully
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/StandardResponse'
        '401':
          description: Not authenticated
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'
        '500':
          description: Internal server error
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'

  /tasks:
    post:
      summary: Create a new task
      tags:
        - Tasks
      security:
        - sessionAuth: []
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/CreateTaskRequest'
      responses:
        '201':
          description: Task created successfully
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/StandardResponse'
        '400':
          description: Invalid input (e.g., title too long)
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'
        '401':
          description: Not authenticated
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'
        '500':
          description: Internal server error
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'

  /protected/test:
    get:
      summary: Example protected route
      tags:
        - Protected
      security:
        - sessionAuth: []
      responses:
        '200':
          description: Access granted
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/StandardResponse'
        '401':
          description: Not authenticated
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'

components:
  securitySchemes:
    sessionAuth:
      type: apiKey
      in: cookie
      name: session_id
      description: Session ID stored in HttpOnly cookie

  schemas:
    StandardResponse:
      type: object
      properties:
        success:
          type: boolean
          example: true
        data:
          type: object
          nullable: true
        message:
          type: string
          nullable: true
      required:
        - success

    ErrorResponse:
      type: object
      properties:
        success:
          type: boolean
          example: false
        error:
          type: object
          properties:
            code:
              type: string
              nullable: true
            message:
              type: string
          required:
            - message
      required:
        - success
        - error

    RegisterRequest:
      type: object
      properties:
        email:
          type: string
          format: email
          example: "user@example.com"
        password:
          type: string
          format: password
          minLength: 8
          example: "securepassword123"
      required:
        - email
        - password

    LoginRequest:
      type: object
      properties:
        email:
          type: string
          format: email
          example: "user@example.com"
        password:
          type: string
          format: password
          example: "securepassword123"
      required:
        - email
        - password

    CreateTaskRequest:
      type: object
      properties:
        title:
          type: string
          maxLength: 255
          example: "Implement feature X"
        description:
          type: string
          nullable: true
          example: "Detailed description of the task"
        assignee_id:
          type: string
          format: uuid
          nullable: true
          example: "123e4567-e89b-12d3-a456-426614174000"
        due_date:
          type: string
          format: date-time
          nullable: true
          example: "2024-12-31T23:59:59Z"
      required:
        - title

    BoardColumn:
      type: object
      properties:
        id:
          type: string
          format: uuid
        name:
          type: string
        order:
          type: integer
        tasks:
          type: array
          items:
            $ref: '#/components/schemas/Task'

    Task:
      type: object
      properties:
        id:
          type: string
          format: uuid
        title:
          type: string
        description:
          type: string
          nullable: true
        column_id:
          type: string
          format: uuid
        assignee_id:
          type: string
          format: uuid
          nullable: true
        due_date:
          type: string
          format: date-time
          nullable: true
        created_by:
          type: string
          format: uuid
        position:
          type: integer
        created_at:
          type: string
          format: date-time
        updated_at:
          type: string
          format: date-time

    BoardResponse:
      type: object
      properties:
        columns:
          type: array
          items:
            $ref: '#/components/schemas/BoardColumn'

    CreateTaskResponse:
      type: object
      properties:
        task_id:
          type: string
          format: uuid
        column_id:
          type: string
          format: uuid
        column:
          type: string
        position:
          type: integer
        title:
          type: string

security:
  - sessionAuth: []
```

## Notes

- Authentication is cookie‑based (`session_id` cookie).
- All protected routes require the `sessionAuth` security scheme.
- Response format follows the standard `{success, data, message}` pattern.
- Error responses include an `error` object with `code` and `message`.
- The YAML is valid OpenAPI 3.0 and can be used with Swagger UI or other tools.

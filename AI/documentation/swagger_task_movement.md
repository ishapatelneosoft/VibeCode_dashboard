# Task Movement API (OpenAPI 3.0)

## Endpoint

`PATCH /api/v1/tasks/{id}/move`

Moves a task to a different column and/or reorders it within a column using fractional indexing.

## Request

```yaml
openapi: 3.0.0
info:
  title: Task Movement API
  version: 1.0.0
servers:
  - url: /api/v1
paths:
  /tasks/{id}/move:
    patch:
      summary: Move or reorder a task
      description: |
        Moves a task to a destination column and positions it relative to other tasks.
        Uses fractional indexing for ordering – provide either `before_task_id`, `after_task_id`, both, or none.
        
        **Rules:**
        - If both `before_task_id` and `after_task_id` are provided, the task is placed between them.
        - If only `before_task_id` is provided, the task is placed after that task (bottom of column).
        - If only `after_task_id` is provided, the task is placed before that task (top of column).
        - If neither is provided, the task is moved to an empty column (position 0).
      tags:
        - Tasks
      security:
        - BearerAuth: []
      parameters:
        - in: path
          name: id
          required: true
          schema:
            type: string
            format: uuid
          description: Task ID
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/MoveTaskRequest'
      responses:
        '200':
          description: Task moved successfully
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Response'
        '400':
          description: Invalid request (missing fields, invalid IDs)
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'
        '401':
          description: Unauthorized (missing or invalid session)
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'
        '404':
          description: Task, column, or referenced task not found
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

components:
  securitySchemes:
    BearerAuth:
      type: http
      scheme: bearer
      bearerFormat: JWT
  schemas:
    MoveTaskRequest:
      type: object
      required:
        - destination_column_id
      properties:
        destination_column_id:
          type: string
          format: uuid
          description: ID of the column where the task should be moved
        before_task_id:
          type: string
          format: uuid
          nullable: true
          description: ID of the task that should be immediately before the moved task (optional)
        after_task_id:
          type: string
          format: uuid
          nullable: true
          description: ID of the task that should be immediately after the moved task (optional)
      example:
        destination_column_id: "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
        before_task_id: "b2c3d4e5-f6g7-8901-bcde-f12345678901"
        after_task_id: null

    MoveTaskResponseData:
      type: object
      properties:
        task_id:
          type: string
          format: uuid
        column_id:
          type: string
          format: uuid
        position:
          type: number
          format: float
          description: New fractional position of the task

    Response:
      type: object
      properties:
        success:
          type: boolean
          example: true
        data:
          $ref: '#/components/schemas/MoveTaskResponseData'
        message:
          type: string
          example: "Task moved successfully"

    ErrorResponse:
      type: object
      properties:
        success:
          type: boolean
          example: false
        error:
          type: string
          example: "Task not found"
```

## Example Request

```json
{
  "destination_column_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
  "before_task_id": "b2c3d4e5-f6g7-8901-bcde-f12345678901",
  "after_task_id": null
}
```

## Example Response

```json
{
  "success": true,
  "data": {
    "task_id": "c3d4e5f6-g7h8-9012-cdef-123456789012",
    "column_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
    "position": 1.75
  },
  "message": "Task moved successfully"
}
```

## Notes

- The `position` field is a floating‑point number enabling fractional indexing.
- The endpoint is idempotent: moving a task to the same location has no effect.
- All referenced tasks must belong to the destination column; otherwise, a validation error is returned.
- The operation is atomic within a database transaction.

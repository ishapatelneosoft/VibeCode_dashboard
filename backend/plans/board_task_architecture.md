# Board & Task Features Architecture

## Overview
Extension of existing authentication system to add shared board visibility and task creation with validation.

## Existing Patterns Analysis

### Current Structure
1. **Domain Layer**: Entities + Repository interfaces
2. **Repository Layer**: PostgreSQL implementations
3. **Service Layer**: Business logic with DTOs
4. **Controller Layer**: HTTP handlers with standardized responses
5. **Middleware**: Authentication middleware

### Key Patterns to Follow
- UUID primary keys
- `CreatedAt`/`UpdatedAt` timestamps
- Repository interfaces in domain layer
- Service methods with request/response DTOs
- Error constants in service layer
- Standardized response format from `internal/controller/response.go`

## Database Schema Design

### Columns Table
```sql
CREATE TABLE columns (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    "order" INTEGER NOT NULL UNIQUE CHECK ("order" >= 0 AND "order" <= 7),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Predefined 8 columns (to be seeded)
-- 0: Backlog, 1: To Do, 2: In Progress, 3: Review, 4: Testing, 5: Done, 6: Blocked, 7: Archived
```

### Tasks Table
```sql
CREATE TABLE tasks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    column_id UUID NOT NULL REFERENCES columns(id) ON DELETE CASCADE,
    assignee_id UUID REFERENCES users(id) ON DELETE SET NULL,
    due_date TIMESTAMP WITH TIME ZONE,
    created_by UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    position INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for performance
CREATE INDEX idx_tasks_column_id ON tasks(column_id);
CREATE INDEX idx_tasks_assignee_id ON tasks(assignee_id);
CREATE INDEX idx_tasks_created_by ON tasks(created_by);
CREATE INDEX idx_tasks_column_position ON tasks(column_id, position);
CREATE INDEX idx_tasks_due_date ON tasks(due_date);
```

## Domain Models

### Column Entity
- `ID`, `Name`, `Order`, `CreatedAt`
- Fixed order (0-7) for 8 columns

### Task Entity
- `ID`, `Title`, `Description`, `ColumnID`, `AssigneeID`, `DueDate`, `CreatedBy`, `Position`, `CreatedAt`, `UpdatedAt`
- Validation: Title max 255 chars

## Repository Interfaces

### ColumnRepository
- `FindAll()` - Get all 8 columns in order
- `FindByID(id)` - Get column by ID
- `FindByName(name)` - Get column by name (e.g., "Backlog")

### TaskRepository
- `Create(task)` - Create new task with position calculation
- `FindByID(id)` - Get task by ID
- `FindByColumn(columnID)` - Get tasks for column ordered by position
- `FindAllWithColumns()` - Get all tasks with column info (for board view)
- `Update(task)` - Update task
- `Delete(id)` - Delete task
- `GetMaxPosition(columnID)` - Get max position in column (for new task placement)

## Service Layer Design

### BoardService
- `GetBoard()` - Returns all columns with tasks
- Follows existing AuthService pattern with DTOs

### TaskService
- `CreateTask(req)` - Validates title, sets defaults, calculates position
- `ValidateTitle(title)` - Title validation logic
- Uses transactions for position calculation to avoid race conditions

## API Endpoints

### GET /api/v1/board
- Returns board with 8 columns in order
- Each column contains ordered tasks
- Authentication required

### POST /api/v1/tasks
- Creates new task with only title required
- Defaults: column="Backlog", assignee=null
- Sets created_by to authenticated user
- Returns task details with position

## Integration Points

### Authentication
- Use existing `AuthMiddleware.RequireAuth()`
- Extract user ID from context for `created_by`

### Response Format
- Use existing `SuccessResponse`, `ErrorResponse` from controller package
- Consistent error handling

## Testing Strategy

### Unit Tests
- Service layer validation logic
- Position calculation
- Title length validation

### Integration Tests
- API endpoints with authentication
- Board returns 8 columns
- Task creation with defaults
- Concurrent task creation handling

## Migration Strategy
1. Create new migration file
2. Seed 8 columns with fixed order
3. Add foreign key to users table

## Performance Considerations
- Indexed queries for column/position
- Single query for board with joins (avoid N+1)
- Transaction for position calculation
- Pagination ready for future scaling
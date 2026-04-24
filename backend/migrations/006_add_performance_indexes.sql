-- Migration: 006_add_performance_indexes.sql
-- Adds performance indexes for better query execution and concurrent handling

-- Index for worklogs.user_id (used in joins and filtering)
CREATE INDEX IF NOT EXISTS idx_worklogs_user_id ON worklogs(user_id);

-- Composite index for worklogs (task_id, created_at) for time report queries
CREATE INDEX IF NOT EXISTS idx_worklogs_task_created ON worklogs(task_id, created_at);

-- Index for tasks.position queries (used in FindTaskBefore/After)
CREATE INDEX IF NOT EXISTS idx_tasks_column_position_filtered ON tasks(column_id, position) WHERE position IS NOT NULL;

-- Index for assignment_history.task_id (used in GetAssignmentHistory)
CREATE INDEX IF NOT EXISTS idx_assignment_history_task_id ON assignment_history(task_id);

-- Index for assignment_history.created_at (for ordering)
CREATE INDEX IF NOT EXISTS idx_assignment_history_created_at ON assignment_history(created_at DESC);

-- Index for sessions.user_id (used in session validation)
CREATE INDEX IF NOT EXISTS idx_sessions_user_id ON sessions(user_id);

-- Index for sessions.expires_at (for cleanup queries)
CREATE INDEX IF NOT EXISTS idx_sessions_expires_at ON sessions(expires_at);

-- Add index for users.email (used in authentication)
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

-- Add index for tasks.created_at (for reporting and filtering)
CREATE INDEX IF NOT EXISTS idx_tasks_created_at ON tasks(created_at);

-- Add index for tasks.updated_at (for activity tracking)
CREATE INDEX IF NOT EXISTS idx_tasks_updated_at ON tasks(updated_at);

-- Add composite index for time report query optimization
CREATE INDEX IF NOT EXISTS idx_tasks_column_assignee ON tasks(column_id, assignee_id) WHERE assignee_id IS NOT NULL;

-- Analyze tables to update statistics
ANALYZE tasks;
ANALYZE worklogs;
ANALYZE assignment_history;
ANALYZE sessions;
ANALYZE users;
-- Migration: 005_create_worklogs_table.sql
-- Creates worklogs table for immutable time logging per task

CREATE TABLE IF NOT EXISTS worklogs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    task_id UUID NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    time_spent DECIMAL NOT NULL CHECK (time_spent > 0),
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_worklogs_task_id ON worklogs(task_id);
CREATE INDEX IF NOT EXISTS idx_worklogs_created_at ON worklogs(created_at);
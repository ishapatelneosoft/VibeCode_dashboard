-- Migration: 004_create_assignment_history_table.sql
-- Creates assignment_history table for tracking task assignee changes

CREATE TABLE IF NOT EXISTS assignment_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    task_id UUID NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    old_assignee_id UUID REFERENCES users(id) ON DELETE SET NULL,
    new_assignee_id UUID REFERENCES users(id) ON DELETE SET NULL,
    changed_by UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_assignment_history_task_id ON assignment_history(task_id);
CREATE INDEX IF NOT EXISTS idx_assignment_history_created_at ON assignment_history(created_at);